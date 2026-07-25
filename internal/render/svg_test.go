package render

import (
	"strings"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/config"
	"github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

// sampleResult generates once and is reused across every style and delivery
// combination below. The styles are a serialization concern — nothing about
// them needs a second pipeline run, and generating per combination would
// multiply P3's share of `make check` by the size of the matrix.
func sampleResult(t *testing.T) (geometry.Result, string) {
	t.Helper()
	corpus := corpusOrSkip(t)
	entity := corpus.Manifest.Entities[0]
	result, err := generate(corpus, entity.Alpha2, config.Settings{
		Profile: strptr("card"), Boundary: strptr("un"),
	})
	if err != nil {
		t.Fatalf("%s: %v", entity.Alpha2, err)
	}
	return result, entity.Alpha2
}

// settingsFor builds a fully resolved configuration the way the CLI does, so a
// test can never assert against a token set the resolution chain would not
// actually produce.
func settingsFor(t *testing.T, style, delivery string, extra config.Settings) config.Settings {
	t.Helper()
	flags := extra
	flags.Style = &style
	flags.Delivery = &delivery
	resolved, err := config.Resolve(nil, "", flags)
	if err != nil {
		t.Fatal(err)
	}
	return resolved.Settings
}

// TestEveryStyleSerializesThroughTokens is VAL-3's core claim. If a style ever
// became a branch in the serializer, this would still pass — so it also asserts
// the negative: the emitted markup never mentions a style name.
func TestEveryStyleSerializesThroughTokens(t *testing.T) {
	result, iso := sampleResult(t)

	for _, style := range config.StyleNames {
		for _, delivery := range config.Deliveries {
			document, err := SVG(result, iso, "Test Entity", settingsFor(t, style, delivery, config.Settings{}))
			if err != nil {
				t.Fatalf("%s/%s: %v", style, delivery, err)
			}
			if err := Structure(document); err != nil {
				t.Errorf("%s/%s: %v", style, delivery, err)
			}
			if strings.Contains(string(document), style) {
				t.Errorf("%s/%s: the style name appears in the output; a style must reach the serializer as tokens, never as a name",
					style, delivery)
			}
			if !strings.Contains(string(document), `d="`) {
				t.Errorf("%s/%s: no path data", style, delivery)
			}
		}
	}
}

// TestStylesProduceDistinctOutput keeps the five styles from collapsing into
// each other. A style that renders identically to another is a name with no
// meaning, and nothing else in the suite would notice.
func TestStylesProduceDistinctOutput(t *testing.T) {
	result, iso := sampleResult(t)

	seen := map[string]string{}
	for _, style := range config.StyleNames {
		document, err := SVG(result, iso, "Test Entity", settingsFor(t, style, "standalone", config.Settings{}))
		if err != nil {
			t.Fatal(err)
		}
		// The path is identical across styles by construction; only presentation
		// may differ, so the comparison is over everything but the path.
		presentation := withoutPathData(string(document))
		if other, clash := seen[presentation]; clash {
			t.Errorf("styles %q and %q serialize identically; one of them means nothing", other, style)
		}
		seen[presentation] = style
	}
}

// TestStandaloneNeedsNoHostCSS is ARCH-INV-6. A standalone asset that carried a
// var() reference would render as nothing in an img tag, which looks like a
// broken file rather than a configuration mistake.
func TestStandaloneNeedsNoHostCSS(t *testing.T) {
	result, iso := sampleResult(t)

	document, err := SVG(result, iso, "Test Entity", settingsFor(t, "bold-soft", "standalone", config.Settings{}))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(document), "var(") {
		t.Errorf("a standalone asset depends on host CSS:\n%s", document)
	}
	if strings.Contains(string(document), RootClass) {
		t.Error("a standalone asset carries theming hooks nothing can target")
	}

	// An author-supplied var() collapses to its fallback rather than being
	// emitted as something that renders as nothing.
	withFallback := settingsFor(t, "bold-soft", "standalone", config.Settings{
		Tokens: &config.Tokens{Fill: strptr("var(--brand, rebeccapurple)")},
	})
	document, err = SVG(result, iso, "Test Entity", withFallback)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(document), `fill="rebeccapurple"`) {
		t.Errorf("the variable did not collapse to its fallback:\n%s", document)
	}

	// One with no fallback is refused, and the diagnostic says what to do.
	noFallback := settingsFor(t, "bold-soft", "standalone", config.Settings{
		Tokens: &config.Tokens{Fill: strptr("var(--brand)")},
	})
	_, err = SVG(result, iso, "Test Entity", noFallback)
	if err == nil {
		t.Fatal("a standalone asset accepted a variable with no fallback")
	}
	if !strings.Contains(err.Error(), "themed-inline") {
		t.Errorf("diagnostic %q does not name the alternative", err)
	}
}

// TestThemedInlineEmitsTheContractedHooks is CTR-005. These names are a public
// compatibility surface: a host stylesheet written against them keeps working
// only if they do not move.
func TestThemedInlineEmitsTheContractedHooks(t *testing.T) {
	result, iso := sampleResult(t)

	document, err := SVG(result, iso, "Test Entity", settingsFor(t, "bold-soft", "themed-inline", config.Settings{
		Marker: &config.Marker{Mode: strptr("all-capitals")},
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := string(document)

	for _, hook := range []string{
		`class="country-map"`, `data-country="` + iso + `"`,
		`class="country-map__shape"`,
		"var(--country-map-fill,", "var(--country-map-stroke,",
		"var(--country-map-stroke-width,",
	} {
		if !strings.Contains(text, hook) {
			t.Errorf("CTR-005 hook %q is missing:\n%s", hook, text)
		}
	}
	// Every custom property carries a fallback, or a page that sets none renders
	// an invisible map.
	if strings.Contains(text, "var(--country-map-fill)") {
		t.Error("a custom property was emitted with no fallback")
	}
	if !strings.Contains(text, "currentColor") {
		t.Error("no fallback resolves to currentColor; REQ-6 names it specifically")
	}
}

// TestAccessibilityModesSerializePredictably is AC-5's second half. The default
// is decorative because the accepted composition is cards beside text that
// already names the country, and announcing each one makes the page worse.
func TestAccessibilityModesSerializePredictably(t *testing.T) {
	result, iso := sampleResult(t)

	decorative, err := SVG(result, iso, "France", settingsFor(t, "outline", "standalone", config.Settings{}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(decorative), `aria-hidden="true"`) {
		t.Errorf("the default is not decorative:\n%s", decorative)
	}
	if strings.Contains(string(decorative), "aria-label") {
		t.Error("a decorative asset is also labelled; assistive technology gets two contradictory signals")
	}
	if !strings.Contains(string(decorative), `focusable="false"`) {
		t.Error("a decorative inline asset can still take keyboard focus, which is a trap on a decoration")
	}

	labelled, err := SVG(result, iso, "France", settingsFor(t, "outline", "standalone", config.Settings{
		Accessibility: &config.Accessibility{Mode: strptr("labelled")},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(labelled), `role="img"`) || !strings.Contains(string(labelled), `aria-label="France"`) {
		t.Errorf("labelled mode did not announce the entity:\n%s", labelled)
	}
	if strings.Contains(string(labelled), "aria-hidden") {
		t.Error("a labelled asset is also hidden")
	}

	custom, err := SVG(result, iso, "France", settingsFor(t, "outline", "standalone", config.Settings{
		Accessibility: &config.Accessibility{Mode: strptr("labelled"), Label: strptr("Map of France")},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(custom), `aria-label="Map of France"`) {
		t.Error("an explicit label did not override the entity name")
	}
}

// TestMarkerModesSelectWhatTheyName is REQ-7. `none` is the default because a
// decorative card should not sprout a dot because the corpus happens to know a
// capital.
func TestMarkerModesSelectWhatTheyName(t *testing.T) {
	corpus := corpusOrSkip(t)
	// An entity with more than one capital is the only one that distinguishes
	// `capital` from `all-capitals`, so it is found rather than named.
	var iso string
	var result geometry.Result
	for _, entity := range corpus.Manifest.Entities {
		if len(entity.Capitals) < 2 {
			continue
		}
		candidate, err := generate(corpus, entity.Alpha2, config.Settings{
			Profile: strptr("card"), Boundary: strptr("un"),
			Marker: &config.Marker{Mode: strptr("all-capitals")},
		})
		if err != nil || len(candidate.Markers) < 2 {
			continue
		}
		iso, result = entity.Alpha2, candidate
		break
	}
	if iso == "" {
		t.Skip("no sampled entity carries two rendered capitals")
	}

	count := func(mode string, custom []string) int {
		settings := settingsFor(t, "filled", "themed-inline", config.Settings{
			Marker: &config.Marker{Mode: &mode, Custom: custom},
		})
		document, err := SVG(result, iso, "Test Entity", settings)
		if err != nil {
			t.Fatal(err)
		}
		return strings.Count(string(document), "<circle")
	}

	if got := count("none", nil); got != 0 {
		t.Errorf("mode none emitted %d markers", got)
	}
	if got := count("capital", nil); got != 1 {
		t.Errorf("mode capital emitted %d markers, want 1", got)
	}
	if got := count("all-capitals", nil); got != len(result.Markers) {
		t.Errorf("mode all-capitals emitted %d markers, want %d", got, len(result.Markers))
	}
	if got := count("custom", []string{result.Markers[1].ID}); got != 1 {
		t.Errorf("mode custom emitted %d markers, want the 1 it named", got)
	}
	if got := count("custom", []string{"not-a-capital-id"}); got != 0 {
		t.Errorf("mode custom emitted %d markers for an id the corpus does not carry", got)
	}
}

// TestAnimationIsOptInAndHookBased is REQ-8. Motion lives in the host
// stylesheet, so the asset carries a class and nothing else — no keyframes, no
// script, nothing that runs.
func TestAnimationIsOptInAndHookBased(t *testing.T) {
	result, iso := sampleResult(t)

	off, err := SVG(result, iso, "Test Entity", settingsFor(t, "outline", "themed-inline", config.Settings{}))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(off), AnimatedClass) {
		t.Error("animation is on by default")
	}

	enabled := true
	on, err := SVG(result, iso, "Test Entity", settingsFor(t, "outline", "themed-inline", config.Settings{
		Animation: &config.Animation{Enabled: &enabled, Hook: strptr("drift")},
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := string(on)
	if !strings.Contains(text, AnimatedClass) || !strings.Contains(text, "drift") {
		t.Errorf("the animation hooks are missing:\n%s", text)
	}
	// Nothing that animates by itself: the host stylesheet owns the motion, which
	// is what makes prefers-reduced-motion the host's to honour and P4's to prove.
	for _, forbidden := range []string{"<animate", "<animateTransform", "@keyframes", "<style"} {
		if strings.Contains(text, forbidden) {
			t.Errorf("the asset carries %q; motion belongs in the host stylesheet", forbidden)
		}
	}
	if err := Structure(on); err != nil {
		t.Errorf("an animated asset is structurally invalid: %v", err)
	}
}

// TestOutputIsByteIdentical is ARCH-INV-4 at the serializer. Attribute order
// comes from insertion, never from a map, and this is what proves it.
func TestOutputIsByteIdentical(t *testing.T) {
	result, iso := sampleResult(t)
	settings := settingsFor(t, "bold-soft", "themed-inline", config.Settings{
		Marker: &config.Marker{Mode: strptr("all-capitals")},
	})

	first, err := SVG(result, iso, "Test Entity", settings)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		again, err := SVG(result, iso, "Test Entity", settings)
		if err != nil {
			t.Fatal(err)
		}
		if string(again) != string(first) {
			t.Fatalf("serialization is not deterministic; run %d differs:\n%s\n%s", i, first, again)
		}
	}
}

// TestNoRedundantMarkup is the REQ-9 clause nobody ever fails, made checkable.
func TestNoRedundantMarkup(t *testing.T) {
	result, iso := sampleResult(t)

	// `silhouette` has no stroke, so every stroke attribute is meaningless on it.
	document, err := SVG(result, iso, "Test Entity", settingsFor(t, "silhouette", "standalone", config.Settings{}))
	if err != nil {
		t.Fatal(err)
	}
	text := string(document)
	for _, absent := range []string{"stroke-width", "stroke-linecap", "stroke-linejoin", "stroke-opacity"} {
		if strings.Contains(text, absent) {
			t.Errorf("a style with no stroke still emits %q:\n%s", absent, text)
		}
	}
	// A default-valued attribute says nothing a renderer does not already assume.
	if strings.Contains(text, `fill-opacity="1"`) {
		t.Errorf("a default-valued attribute is emitted:\n%s", text)
	}
	for _, absent := range []string{"<g", "<metadata", "<!--", "<?xml", "<desc"} {
		if strings.Contains(text, absent) {
			t.Errorf("output carries %q:\n%s", absent, text)
		}
	}
}

func withoutPathData(document string) string {
	start := strings.Index(document, ` d="`)
	if start < 0 {
		return document
	}
	end := strings.Index(document[start+4:], `"`)
	if end < 0 {
		return document
	}
	return document[:start] + document[start+4+end:]
}
