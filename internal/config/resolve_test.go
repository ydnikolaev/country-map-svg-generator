package config

import (
	"strings"
	"testing"
)

func decodeOrFail(t *testing.T, doc string) *Document {
	t.Helper()
	parsed, err := Decode([]byte(doc), FormatYAML)
	if err != nil {
		t.Fatalf("fixture did not decode: %v", err)
	}
	return parsed
}

// TestPrecedenceRunsInTheDocumentedOrder walks the whole chain in one fixture,
// with a different layer winning each key. Any reordering breaks at least one
// assertion.
func TestPrecedenceRunsInTheDocumentedOrder(t *testing.T) {
	doc := decodeOrFail(t, `
schema: country-map/v1
extends: site-default
style: ghost
profiles:
  hero:
    tokens: { strokeWidth: 2 }
countries:
  US:
    profile: hero
    tokens: { fill: "var(--us)" }
`)
	resolved, err := Resolve(doc, "US", Settings{Delivery: strptr("standalone")})
	if err != nil {
		t.Fatal(err)
	}
	origins := resolved.Provenance

	// Nobody but the defaults mentions the boundary.
	if *resolved.Settings.Boundary != "un" || origins["boundary"] != LayerNames[0] {
		t.Errorf("boundary = %v from %q", *resolved.Settings.Boundary, origins["boundary"])
	}
	// The preset beats the defaults.
	// The style the document selects supplies the tokens nobody overrode.
	if !strings.HasPrefix(origins["tokens.fillOpacity"], LayerNames[1]) {
		t.Errorf("fillOpacity origin = %q, want the style defaults", origins["tokens.fillOpacity"])
	}
	// The document beats the preset.
	if *resolved.Settings.Style != "ghost" || origins["style"] != LayerNames[3] {
		t.Errorf("style = %v from %q", *resolved.Settings.Style, origins["style"])
	}
	// The profile block applies, and it is the block for the profile the country
	// override selected — not the one the document's globals implied.
	if resolved.Settings.Tokens.StrokeWidth == nil || *resolved.Settings.Tokens.StrokeWidth != 2 {
		t.Fatalf("the hero profile block did not apply: %+v", resolved.Settings.Tokens)
	}
	if !strings.HasPrefix(origins["tokens.strokeWidth"], LayerNames[4]) {
		t.Errorf("strokeWidth origin = %q, want the profile block", origins["tokens.strokeWidth"])
	}
	// The country override beats the profile block.
	if *resolved.Settings.Tokens.Fill != "var(--us)" || !strings.HasPrefix(origins["tokens.fill"], LayerNames[5]) {
		t.Errorf("fill = %v from %q", *resolved.Settings.Tokens.Fill, origins["tokens.fill"])
	}
	// Flags beat everything.
	if *resolved.Settings.Delivery != "standalone" || origins["delivery"] != LayerNames[6] {
		t.Errorf("delivery = %v from %q", *resolved.Settings.Delivery, origins["delivery"])
	}
}

// TestTheProfileBlockFollowsTheCountryOverride is the two-pass resolution. In a
// single pass the profile block would be chosen from the document's globals,
// and `countries: {US: {profile: hero}}` — the spec's own example — would apply
// the card block to a hero card.
func TestTheProfileBlockFollowsTheCountryOverride(t *testing.T) {
	doc := decodeOrFail(t, `
schema: country-map/v1
profile: card
profiles:
  card:
    tokens: { strokeWidth: 1 }
  hero:
    tokens: { strokeWidth: 4 }
countries:
  US:
    profile: hero
`)
	us, err := Resolve(doc, "US", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Settings.Tokens.StrokeWidth != 4 {
		t.Errorf("US strokeWidth = %v, want the hero block's 4", *us.Settings.Tokens.StrokeWidth)
	}

	fr, err := Resolve(doc, "FR", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	if *fr.Settings.Tokens.StrokeWidth != 1 {
		t.Errorf("FR strokeWidth = %v, want the card block's 1", *fr.Settings.Tokens.StrokeWidth)
	}
}

// TestAFlagCanRedirectTheProfileBlockToo is the same property from the strongest
// layer.
func TestAFlagCanRedirectTheProfileBlockToo(t *testing.T) {
	doc := decodeOrFail(t, `
schema: country-map/v1
profile: card
profiles:
  card: { tokens: { strokeWidth: 1 } }
  hero: { tokens: { strokeWidth: 4 } }
`)
	resolved, err := Resolve(doc, "FR", Settings{Profile: strptr("hero")})
	if err != nil {
		t.Fatal(err)
	}
	if *resolved.Settings.Tokens.StrokeWidth != 4 {
		t.Errorf("strokeWidth = %v, want the hero block's 4", *resolved.Settings.Tokens.StrokeWidth)
	}
}

// TestEachPresetInAChainIsNamedSeparately matters for `explain`: with a chain of
// three, "which preset set this" is the actual question.
func TestEachPresetInAChainIsNamedSeparately(t *testing.T) {
	doc := decodeOrFail(t, "schema: country-map/v1\nextends: standalone-default\n")
	resolved, err := Resolve(doc, "", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	if got := resolved.Provenance["delivery"]; got != LayerNames[2]+" standalone-default" {
		t.Errorf("delivery origin = %q, want the exact preset that set it", got)
	}
	// A token the preset never mentions is attributed to the style it selected,
	// not swept into the preset.
	if got := resolved.Provenance["tokens.lineCap"]; got != LayerNames[1]+" filled" {
		t.Errorf("lineCap origin = %q, want the style the preset selected", got)
	}
	// A value neither preset touches is still attributed to the defaults, not
	// swept into the nearest preset.
	if got := resolved.Provenance["profile"]; got != LayerNames[0] {
		t.Errorf("profile origin = %q, want the embedded defaults", got)
	}
}

// TestSwitchingModeDiscardsTheOtherModesSizing is the specification's own
// example: globals size a tight layout, and a country override reframes that one
// entity with contain. Without pruning, the merged layout would carry a contain
// mode and an inherited long side, and the documented shape would be refused.
func TestSwitchingModeDiscardsTheOtherModesSizing(t *testing.T) {
	doc := decodeOrFail(t, `
schema: country-map/v1
layout: { mode: tight, longSide: 160 }
countries:
  US:
    layout: { mode: contain, width: 720, height: 420 }
`)
	if err := ValidateDocument(doc, testVocab()); err != nil {
		t.Fatalf("the documented shape failed the per-layer check: %v", err)
	}

	us, err := Resolve(doc, "US", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateResolved("US", us, testVocab()); err != nil {
		t.Fatalf("the documented shape failed after merging: %v", err)
	}
	if us.Settings.Layout.LongSide != nil {
		t.Error("the superseded long side survived into a contain layout")
	}
	// It disappears from the report too: a value that no longer applies must not
	// be shown as if it did.
	if _, reported := us.Provenance["layout.longSide"]; reported {
		t.Error("explain would still report a long side that no longer applies")
	}

	// Every other entity keeps the tight layout the globals set.
	fr, err := Resolve(doc, "FR", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	if fr.Settings.Layout.LongSide == nil || *fr.Settings.Layout.LongSide != 160 {
		t.Errorf("pruning leaked into an entity that never overrode the mode: %+v", fr.Settings.Layout)
	}
}

// TestValidateResolvedCatchesCrossLayerContradictions is why resolution has its
// own validation pass. Pruning only drops sizing set *earlier* than the mode, so
// a foreign field arriving later is still a genuine contradiction — the author
// sized a frame for a mode that is not in effect.
func TestValidateResolvedCatchesCrossLayerContradictions(t *testing.T) {
	doc := decodeOrFail(t, `
schema: country-map/v1
layout: { mode: tight, longSide: 160 }
countries:
  US:
    layout: { width: 300 }
`)
	// Each layer is legal alone: the override names no mode, so nothing in it
	// contradicts anything in it.
	if err := ValidateDocument(doc, testVocab()); err != nil {
		t.Fatalf("the per-layer check should pass: %v", err)
	}
	resolved, err := Resolve(doc, "US", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	err = ValidateResolved("US", resolved, testVocab())
	if err == nil {
		t.Fatal("a width sized for a mode that is not in effect was accepted")
	}
	if !strings.Contains(err.Error(), "layout.width") {
		t.Errorf("diagnostic %q does not name the offending field", err)
	}

	// And the mirror: contain without a frame.
	other := decodeOrFail(t, "schema: country-map/v1\nlayout: { mode: contain }\n")
	resolvedOther, err := Resolve(other, "", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	err = ValidateResolved("", resolvedOther, testVocab())
	if err == nil || !strings.Contains(err.Error(), "width and height is missing") {
		t.Errorf("contain without a frame was accepted or misreported: %v", err)
	}
}

// TestAnEmptyDocumentResolvesToTheDefaults keeps the zero-config path working:
// running the tool with nothing configured must produce a complete, valid
// configuration rather than a pile of "not set" diagnostics.
func TestAnEmptyDocumentResolvesToTheDefaults(t *testing.T) {
	resolved, err := Resolve(nil, "FR", Settings{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateResolved("FR", resolved, testVocab()); err != nil {
		t.Fatalf("the zero-config path does not resolve to a complete configuration: %v", err)
	}
	if *resolved.Settings.Profile != "card" || *resolved.Settings.Boundary != "un" {
		t.Errorf("defaults = profile %v boundary %v", *resolved.Settings.Profile, *resolved.Settings.Boundary)
	}
	// The long side is deliberately absent: the geometry preset owns it, and
	// restating it here would fork the byte budget.
	if resolved.Settings.Layout.LongSide != nil {
		t.Error("the config defaults restate a long side the geometry preset owns")
	}
}

// TestOutputCollisionsAreRefusedBeforeAnythingIsWritten covers the failure whose
// symptom is a catalog quietly short by however many entities collided.
func TestOutputCollisionsAreRefusedBeforeAnythingIsWritten(t *testing.T) {
	err := CheckOutputCollisions(map[string]string{
		"US": "out/map.svg", "FR": "out/map.svg", "CL": "out/cl.svg",
	})
	if err == nil {
		t.Fatal("a collision was accepted")
	}
	if !strings.Contains(err.Error(), "FR, US") {
		t.Errorf("diagnostic %q does not name both claimants in stable order", err)
	}
	if !strings.Contains(err.Error(), "{iso}") {
		t.Errorf("diagnostic %q does not say how to fix it", err)
	}
	if err := CheckOutputCollisions(map[string]string{"US": "a.svg", "FR": "b.svg"}); err != nil {
		t.Errorf("distinct paths were reported as a collision: %v", err)
	}
}

// TestOutputPathExpandsExactlyTheValidatedPlaceholders keeps the substitution
// and the closed-set check from drifting: every placeholder validation accepts
// must expand, and nothing else may.
func TestOutputPathExpandsExactlyTheValidatedPlaceholders(t *testing.T) {
	settings := Settings{
		Profile: strptr("hero"), Boundary: strptr("de_facto"),
		Style: strptr("ghost"), Delivery: strptr("standalone"),
		Output: &Output{Dir: strptr("out"), Filename: strptr("{iso}-{profile}-{boundary}-{style}-{delivery}.svg")},
	}
	path, err := OutputPath("US", settings)
	if err != nil {
		t.Fatal(err)
	}
	want := "out/US-hero-de_facto-ghost-standalone.svg"
	if path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
	if strings.Contains(path, "{") {
		t.Error("a validated placeholder was left unexpanded")
	}

	// Every placeholder the validator accepts must be one this expands, or a
	// config could pass validation and then emit a literal brace in a file name.
	for _, placeholder := range FilenamePlaceholders {
		settings.Output.Filename = strptr("{" + placeholder + "}-{iso}.svg")
		expanded, err := OutputPath("US", settings)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(expanded, "{") {
			t.Errorf("placeholder {%s} is accepted by validation but not expanded", placeholder)
		}
	}
}
