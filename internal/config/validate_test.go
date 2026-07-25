package config

import (
	"encoding/json"
	"strings"
	"testing"
)

// testVocab is deliberately not the real one. Injecting values here proves the
// package holds no literals of its own: if `card` or `de_facto` were hard-coded
// anywhere in validation, these fixtures would disagree with it.
func testVocab() Vocabulary {
	return Vocabulary{
		Profiles:   []string{"card", "hero"},
		Boundaries: []string{"un", "de_facto"},
		ISOCodes:   []string{"US", "FR", "CL"},
		CapitalIDs: map[string][]string{"US": {"us-washington"}, "FR": {"fr-paris"}},
	}
}

func validateRaw(t *testing.T, raw string) error {
	t.Helper()
	var generic any
	if err := json.Unmarshal([]byte(raw), &generic); err != nil {
		t.Fatalf("fixture is not valid JSON: %v", err)
	}
	if err := CheckKeys(generic); err != nil {
		return err
	}
	var doc Document
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return err
	}
	return ValidateDocument(&doc, testVocab())
}

func mustReject(t *testing.T, name, raw, wantSubstring string) {
	t.Helper()
	err := validateRaw(t, raw)
	if err == nil {
		t.Errorf("%s: accepted", name)
		return
	}
	if wantSubstring != "" && !strings.Contains(err.Error(), wantSubstring) {
		t.Errorf("%s: diagnostic %q does not mention %q", name, err.Error(), wantSubstring)
	}
}

func mustAccept(t *testing.T, name, raw string) {
	t.Helper()
	if err := validateRaw(t, raw); err != nil {
		t.Errorf("%s: rejected: %v", name, err)
	}
}

// TestLayoutModesTakeDisjointFields is REQ-13's core. Each rejection below is a
// value that would otherwise be silently ignored, producing a plausible card at
// the wrong size — the failure shape that is hardest to notice.
func TestLayoutModesTakeDisjointFields(t *testing.T) {
	mustReject(t, "tight with width",
		`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":160,"width":300}}`,
		"not accepted in tight mode")
	mustReject(t, "tight with height",
		`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":160,"height":300}}`,
		"not accepted in tight mode")
	mustReject(t, "contain with longSide",
		`{"schema":"country-map/v1","layout":{"mode":"contain","width":720,"height":420,"longSide":160}}`,
		"not accepted in contain mode")
	mustReject(t, "contain with maxWidth",
		`{"schema":"country-map/v1","layout":{"mode":"contain","width":720,"height":420,"maxWidth":160}}`,
		"not accepted in contain mode")
	mustReject(t, "two ways to size a tight layout",
		`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":160,"maxWidth":200}}`,
		"two ways to size")

	mustAccept(t, "tight by long side", `{"schema":"country-map/v1","layout":{"mode":"tight","longSide":160}}`)
	mustAccept(t, "tight by box", `{"schema":"country-map/v1","layout":{"mode":"tight","maxWidth":200,"maxHeight":140}}`)
	mustAccept(t, "contain", `{"schema":"country-map/v1","layout":{"mode":"contain","width":720,"height":420}}`)
}

// TestDimensionsMustBeSafeFiniteAndPositive covers "any safe finite positive
// frame". Each of these reaches the viewBox if unchecked.
func TestDimensionsMustBeSafeFiniteAndPositive(t *testing.T) {
	for _, value := range []string{"0", "-1", "1e300"} {
		mustReject(t, "long side "+value,
			`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":`+value+`}}`, "")
	}
	// JSON has no literal for NaN or Infinity, so a non-finite value can only
	// arrive through YAML (.inf, .nan) or a flag. The check exists for those
	// routes; here the boundary is exercised through the ceiling instead.
	mustAccept(t, "at the ceiling",
		`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":1000000}}`)
	mustReject(t, "just over the ceiling",
		`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":1000001}}`, "must not exceed")
}

// TestPaddingAcceptsZeroButNotNegative separates padding from the other
// dimensions: zero padding is ordinary, zero width is not.
func TestPaddingAcceptsZeroButNotNegative(t *testing.T) {
	mustAccept(t, "zero padding",
		`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":160,"padding":{"top":0,"right":0,"bottom":0,"left":0}}}`)
	mustReject(t, "negative padding",
		`{"schema":"country-map/v1","layout":{"mode":"tight","longSide":160,"padding":{"top":-1}}}`,
		"must not be negative")
}

// TestExternalReferencesAreRefused is REQ-9 reaching the config boundary. A
// paint token is author text that lands in an SVG attribute, so it is checked
// where the author can still be told which key was wrong.
func TestExternalReferencesAreRefused(t *testing.T) {
	for name, value := range map[string]string{
		"http":       "url(https://cdn.example/x.svg#g)",
		"protocol":   "url(//cdn.example/x.svg#g)",
		"data":       "data:image/png;base64,AAAA",
		"javascript": "javascript:alert(1)",
		"markup":     "<script>",
	} {
		mustReject(t, "fill "+name,
			`{"schema":"country-map/v1","tokens":{"fill":`+quote(value)+`}}`, "")
	}
	mustAccept(t, "currentColor", `{"schema":"country-map/v1","tokens":{"fill":"currentColor"}}`)
	mustAccept(t, "css variable with fallback",
		`{"schema":"country-map/v1","tokens":{"fill":"var(--map-fill, currentColor)"}}`)
	mustAccept(t, "local fragment",
		`{"schema":"country-map/v1","tokens":{"fill":"url(#brand-gradient)"}}`)
}

// TestAdvancedReferencesMustBeLocalFragments bounds REQ-5's escape hatch. The
// fragment id is constrained too: it is emitted into markup, so it may not carry
// anything that could close an attribute.
func TestAdvancedReferencesMustBeLocalFragments(t *testing.T) {
	mustReject(t, "bare id",
		`{"schema":"country-map/v1","tokens":{"advanced":{"gradient":"brand"}}}`, "url(#id)")
	mustReject(t, "external",
		`{"schema":"country-map/v1","tokens":{"advanced":{"filter":"url(https://x/y#f)"}}}`, "url(#id)")
	mustReject(t, "unsafe id",
		`{"schema":"country-map/v1","tokens":{"advanced":{"pattern":"url(#a\"onload=x)"}}}`, "may contain only")
	mustAccept(t, "safe local",
		`{"schema":"country-map/v1","tokens":{"advanced":{"gradient":"url(#brand_gradient-1)"}}}`)
}

// TestOpacityIsBounded keeps a token that is a ratio from being written as a
// percentage, which would render as fully opaque and look intentional.
func TestOpacityIsBounded(t *testing.T) {
	mustReject(t, "percentage", `{"schema":"country-map/v1","tokens":{"fillOpacity":14}}`, "between 0 and 1")
	mustReject(t, "negative", `{"schema":"country-map/v1","tokens":{"strokeOpacity":-0.1}}`, "between 0 and 1")
	mustAccept(t, "ratio", `{"schema":"country-map/v1","tokens":{"fillOpacity":0.14,"strokeOpacity":1}}`)
}

// TestUnknownEnumValuesNameTheAcceptedSet is REQ-12's "name remediation
// context" for the case an agent hits most: a plausible value that is not the
// accepted spelling.
func TestUnknownEnumValuesNameTheAcceptedSet(t *testing.T) {
	// The hyphen is the specific trap: the corpus manifest declares its profile
	// list as "de-facto" while geometry accepts "de_facto", so an author reading
	// the manifest writes the spelling that does not work.
	err := validateRaw(t, `{"schema":"country-map/v1","boundary":"de-facto"}`)
	if err == nil {
		t.Fatal("the hyphenated boundary spelling was accepted")
	}
	if !strings.Contains(err.Error(), "de_facto") {
		t.Errorf("diagnostic %q does not name the accepted spelling", err.Error())
	}

	mustReject(t, "unknown style", `{"schema":"country-map/v1","style":"neon"}`, "outline")
	mustReject(t, "unknown delivery", `{"schema":"country-map/v1","delivery":"inline"}`, "themed-inline")
	mustReject(t, "unknown profile", `{"schema":"country-map/v1","profile":"poster"}`, "card")
	mustReject(t, "unknown marker mode", `{"schema":"country-map/v1","marker":{"mode":"capitals"}}`, "all-capitals")
}

// TestSchemaVersionIsExact keeps a v2 document from being reinterpreted by a v1
// binary rather than refused.
func TestSchemaVersionIsExact(t *testing.T) {
	mustReject(t, "future major", `{"schema":"country-map/v2"}`, "country-map/v1")
	mustReject(t, "missing", `{}`, "country-map/v1")
	mustReject(t, "prefix match", `{"schema":"country-map/v1beta"}`, "country-map/v1")
}

// TestCountryAndProfileOverlaysAreChecked proves overlays are not a validation
// blind spot. A per-country block is where the unusual values live, so it is
// exactly where a permissive check would hurt.
func TestCountryAndProfileOverlaysAreChecked(t *testing.T) {
	mustReject(t, "unknown iso",
		`{"schema":"country-map/v1","countries":{"XX":{"style":"ghost"}}}`, "unknown ISO")
	mustReject(t, "bad style in overlay",
		`{"schema":"country-map/v1","countries":{"US":{"style":"neon"}}}`, "outline")
	mustReject(t, "bad layout in overlay",
		`{"schema":"country-map/v1","countries":{"US":{"layout":{"mode":"contain","width":720,"height":420,"longSide":160}}}}`,
		"not accepted in contain mode")
	mustReject(t, "profile block re-selecting a profile",
		`{"schema":"country-map/v1","profiles":{"card":{"profile":"hero"}}}`,
		"cannot select a different profile")
	mustReject(t, "unknown profile block",
		`{"schema":"country-map/v1","profiles":{"poster":{"style":"ghost"}}}`, "unknown profile")

	mustAccept(t, "the spec's example shape",
		`{"schema":"country-map/v1","profile":"card","layout":{"mode":"tight","longSide":160},"delivery":"themed-inline","style":"bold-soft","tokens":{"fill":"var(--map-fill, currentColor)","fillOpacity":0.14,"stroke":"var(--map-stroke, currentColor)"},"countries":{"US":{"profile":"hero","layout":{"mode":"contain","width":720,"height":420},"marker":{"mode":"capital"}}}}`)
}

// TestCustomMarkersResolveAgainstTheRegistry keeps a marker from disagreeing
// with the corpus about where a place is: the author selects an id, never a
// coordinate.
func TestCustomMarkersResolveAgainstTheRegistry(t *testing.T) {
	mustReject(t, "unknown capital id",
		`{"schema":"country-map/v1","countries":{"US":{"marker":{"mode":"custom","custom":["us-newyork"]}}}}`,
		"unknown capital id")
	mustReject(t, "custom list without custom mode",
		`{"schema":"country-map/v1","countries":{"US":{"marker":{"mode":"capital","custom":["us-washington"]}}}}`,
		"only meaningful when marker.mode is custom")
	mustAccept(t, "known capital id",
		`{"schema":"country-map/v1","countries":{"US":{"marker":{"mode":"custom","custom":["us-washington"]}}}}`)
}

// TestFilenameTemplateIsClosed keeps an unknown placeholder from appearing
// literally in a file name, and keeps a batch from collapsing onto one path.
func TestFilenameTemplateIsClosed(t *testing.T) {
	mustReject(t, "unknown placeholder",
		`{"schema":"country-map/v1","output":{"filename":"{iso}-{country}.svg"}}`, "unknown placeholder")
	mustReject(t, "unclosed",
		`{"schema":"country-map/v1","output":{"filename":"{iso.svg"}}`, "unclosed")
	mustReject(t, "path not a name",
		`{"schema":"country-map/v1","output":{"filename":"maps/{iso}.svg"}}`, "not a path")
	// A template without {iso} is accepted here on purpose: it is a legitimate
	// single-entity configuration, and for a batch the collision check refuses it
	// with a diagnostic that names which entities actually collided.
	mustAccept(t, "single-entity template",
		`{"schema":"country-map/v1","output":{"filename":"map-{profile}.svg"}}`)
	mustAccept(t, "closed set",
		`{"schema":"country-map/v1","output":{"filename":"{iso}-{profile}-{boundary}.svg"}}`)
}

// TestNoVocabularyMeansNoEntityChecks keeps `validate` usable without a corpus:
// an empty vocabulary skips the entity-level questions rather than failing them
// all closed, which would make the command useless exactly when an author is
// drafting.
func TestNoVocabularyMeansNoEntityChecks(t *testing.T) {
	var doc Document
	raw := `{"schema":"country-map/v1","profile":"anything","boundary":"anything","countries":{"XX":{"style":"ghost"}}}`
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDocument(&doc, Vocabulary{}); err != nil {
		t.Errorf("an empty vocabulary failed entity checks instead of skipping them: %v", err)
	}
	// Shape checks still bite without a corpus.
	if err := ValidateDocument(&Document{Schema: SchemaVersion, Settings: Settings{Style: strptr("neon")}}, Vocabulary{}); err == nil {
		t.Error("an empty vocabulary also disabled the package's own enums")
	}
}

// TestAllProblemsAreReportedInStableOrder matters for an agent and for the
// golden fixtures a later task will pin.
func TestAllProblemsAreReportedInStableOrder(t *testing.T) {
	err := validateRaw(t, `{"schema":"country-map/v9","style":"neon","delivery":"inline","tokens":{"fillOpacity":5}}`)
	if err == nil {
		t.Fatal("accepted")
	}
	problems, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("error type = %T, want ValidationErrors", err)
	}
	if len(problems) != 4 {
		t.Fatalf("reported %d problems, want 4: %v", len(problems), err)
	}
	for i := 1; i < len(problems); i++ {
		if problems[i-1].Path > problems[i].Path {
			t.Fatalf("problems are not in path order: %v", err)
		}
	}
}

func quote(value string) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func strptr(value string) *string { return &value }
