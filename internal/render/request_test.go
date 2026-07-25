package render

import (
	"strings"
	"testing"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
	"github.com/ydnikolaev/country-map-svg-generator/internal/geometry"
)

func corpusOrSkip(t *testing.T) *catalog.Corpus {
	t.Helper()
	corpus, err := catalog.Embedded()
	if err != nil {
		t.Fatalf("embedded corpus: %v", err)
	}
	return corpus
}

func strptr(v string) *string { return &v }
func f64(v float64) *float64  { return &v }

// TestTheTwoAxesAreNotTransposed is the tooth this whole file exists for.
//
// The configuration's `profile` is geometry's Preset and the configuration's
// `boundary` is geometry's Profile — the words cross. Swapped, the request would
// ask for the wrong boundary posture, the output would satisfy every byte budget
// and every structural check, and nothing but a human looking at a disputed
// border would ever notice.
func TestTheTwoAxesAreNotTransposed(t *testing.T) {
	corpus := corpusOrSkip(t)
	iso := corpus.Manifest.Entities[0].Alpha2

	input, err := GeometryRequest(corpus, iso, config.Settings{
		Profile: strptr("hero"), Boundary: strptr("de_facto"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if input.Preset != "hero" {
		t.Errorf("geometry Preset = %q, want the configuration's profile (hero)", input.Preset)
	}
	if input.Profile != "de_facto" {
		t.Errorf("geometry Profile = %q, want the configuration's boundary (de_facto)", input.Profile)
	}
	// Stated as the failure rather than as two assertions, because this is the
	// exact swap that would otherwise pass everything downstream.
	if input.Preset == "de_facto" || input.Profile == "hero" {
		t.Fatal("the two axes are transposed: profile and boundary were swapped on the way into geometry")
	}
}

// TestTheHyphenatedSpellingIsRefused covers the trap the corpus manifest sets.
// It declares its profile list as {"un", "de-facto"} while geometry accepts the
// underscore, so an author reading the manifest writes the spelling that does
// not work.
func TestTheHyphenatedSpellingIsRefused(t *testing.T) {
	corpus := corpusOrSkip(t)
	iso := corpus.Manifest.Entities[0].Alpha2

	_, err := GeometryRequest(corpus, iso, config.Settings{
		Profile: strptr("card"), Boundary: strptr("de-facto"),
	})
	if err == nil {
		t.Fatal("the hyphenated boundary spelling reached geometry")
	}
	if !strings.Contains(err.Error(), "unknown profile") {
		t.Errorf("refused for the wrong reason: %v", err)
	}
}

// TestVocabularyComesFromGeometryNotTheManifest guards the other half of the
// same trap: building the boundary enum from the corpus manifest would offer
// authors "de-facto", which the pipeline refuses.
func TestVocabularyComesFromGeometryNotTheManifest(t *testing.T) {
	corpus := corpusOrSkip(t)
	vocab, err := Vocabulary(corpus)
	if err != nil {
		t.Fatal(err)
	}
	for _, boundary := range vocab.Boundaries {
		if strings.Contains(boundary, "-") {
			t.Errorf("boundary vocabulary contains the manifest's hyphenated spelling %q", boundary)
		}
	}
	for _, offered := range vocab.Boundaries {
		if _, err := geometry.InputFromCatalog(corpus, corpus.Manifest.Entities[0].Alpha2, offered, "card"); err != nil {
			t.Errorf("the vocabulary offers %q, which geometry refuses: %v", offered, err)
		}
	}
	// Profiles are enumerated from the geometry presets, so a new preset appears
	// in the CLI with no edit to the config layer.
	presets, err := geometry.Presets()
	if err != nil {
		t.Fatal(err)
	}
	if len(vocab.Profiles) != len(presets) {
		t.Errorf("profiles = %v, want one per geometry preset (%d)", vocab.Profiles, len(presets))
	}
	if len(vocab.ISOCodes) != corpus.Manifest.EntityCount {
		t.Errorf("ISO codes = %d, want the corpus entity count %d", len(vocab.ISOCodes), corpus.Manifest.EntityCount)
	}
}

// TestAModeWithoutSizingLeavesTheLayoutToThePreset is the subtle rule.
// geometry.ApplyPreset fills a layout only when its mode is empty, so returning
// a mode with no size would suppress the preset's long side and produce a card
// sized by nothing at all.
func TestAModeWithoutSizingLeavesTheLayoutToThePreset(t *testing.T) {
	corpus := corpusOrSkip(t)
	iso := corpus.Manifest.Entities[0].Alpha2

	input, err := GeometryRequest(corpus, iso, config.Settings{
		Profile: strptr("card"), Boundary: strptr("un"),
		Layout: &config.Layout{Mode: strptr("tight")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if input.Layout.Mode != "" {
		t.Fatalf("layout mode = %q; a mode with no sizing must stay empty so ApplyPreset can fill it", input.Layout.Mode)
	}

	// And the preset really does fill it, which is the property being protected.
	applied, err := geometry.ApplyPreset(input)
	if err != nil {
		t.Fatal(err)
	}
	if applied.Layout.LongSide <= 0 {
		t.Errorf("the preset did not supply a long side: %+v", applied.Layout)
	}
	if applied.MaxPathBytes <= 0 {
		t.Errorf("the preset did not supply a byte budget: %d", applied.MaxPathBytes)
	}
}

// TestExplicitSizingSurvivesIntoTheRequest is the other side: when the author
// does say a size, it must reach geometry unchanged.
func TestExplicitSizingSurvivesIntoTheRequest(t *testing.T) {
	corpus := corpusOrSkip(t)
	iso := corpus.Manifest.Entities[0].Alpha2

	tight, err := GeometryRequest(corpus, iso, config.Settings{
		Profile: strptr("card"), Boundary: strptr("un"),
		Layout: &config.Layout{Mode: strptr("tight"), LongSide: f64(160), Padding: &config.Padding{Uniform: f64(8)}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if tight.Layout.Mode != geometry.LayoutTight || tight.Layout.LongSide != 160 {
		t.Errorf("tight layout = %+v", tight.Layout)
	}
	if tight.Layout.Padding != (geometry.Insets{Top: 8, Right: 8, Bottom: 8, Left: 8}) {
		t.Errorf("uniform padding did not expand to four sides: %+v", tight.Layout.Padding)
	}

	contain, err := GeometryRequest(corpus, iso, config.Settings{
		Profile: strptr("hero"), Boundary: strptr("un"),
		Layout: &config.Layout{Mode: strptr("contain"), Width: f64(720), Height: f64(420)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if contain.Layout.Mode != geometry.LayoutContain || contain.Layout.Width != 720 || contain.Layout.Height != 420 {
		t.Errorf("contain layout = %+v", contain.Layout)
	}
}

// TestTheBudgetIsNeverRestated keeps the frozen byte limits in one place. A
// request that carried its own MaxPathBytes would fork them.
func TestTheBudgetIsNeverRestated(t *testing.T) {
	corpus := corpusOrSkip(t)
	input, err := GeometryRequest(corpus, corpus.Manifest.Entities[0].Alpha2, config.Settings{
		Profile: strptr("card"), Boundary: strptr("un"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if input.MaxPathBytes != 0 {
		t.Errorf("MaxPathBytes = %d; the request must leave the budget to the geometry preset", input.MaxPathBytes)
	}
}

// TestAnUnresolvedRequestFailsRatherThanGuessing keeps a half-resolved
// configuration from silently rendering something plausible.
func TestAnUnresolvedRequestFailsRatherThanGuessing(t *testing.T) {
	corpus := corpusOrSkip(t)
	iso := corpus.Manifest.Entities[0].Alpha2

	if _, err := GeometryRequest(corpus, iso, config.Settings{Boundary: strptr("un")}); err == nil {
		t.Error("a request with no profile was accepted")
	}
	if _, err := GeometryRequest(corpus, iso, config.Settings{Profile: strptr("card")}); err == nil {
		t.Error("a request with no boundary was accepted")
	}
}
