package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBothFormatsAgreeOnEveryDocument is the whole reason both formats funnel
// through one typed decode. If YAML and JSON disagreed about which keys exist or
// how a value is spelled, a config would validate in one form and not the other,
// and the author would have no way to tell which behaviour was intended.
func TestBothFormatsAgreeOnEveryDocument(t *testing.T) {
	yamlDoc := `
schema: country-map/v1
extends: site-default
profile: card
boundary: de_facto
delivery: themed-inline
style: bold-soft
layout:
  mode: tight
  longSide: 160
  padding: 8
tokens:
  fill: var(--map-fill, currentColor)
  fillOpacity: 0.14
marker:
  mode: custom
  custom: [us-washington]
countries:
  US:
    profile: hero
    layout: { mode: contain, width: 720, height: 420 }
`
	jsonDoc := `{
  "schema": "country-map/v1",
  "extends": "site-default",
  "profile": "card",
  "boundary": "de_facto",
  "delivery": "themed-inline",
  "style": "bold-soft",
  "layout": {"mode": "tight", "longSide": 160, "padding": 8},
  "tokens": {"fill": "var(--map-fill, currentColor)", "fillOpacity": 0.14},
  "marker": {"mode": "custom", "custom": ["us-washington"]},
  "countries": {"US": {"profile": "hero", "layout": {"mode": "contain", "width": 720, "height": 420}}}
}`

	fromYAML, err := Decode([]byte(yamlDoc), FormatYAML)
	if err != nil {
		t.Fatalf("YAML: %v", err)
	}
	fromJSON, err := Decode([]byte(jsonDoc), FormatJSON)
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}

	if *fromYAML.Profile != *fromJSON.Profile || *fromYAML.Boundary != *fromJSON.Boundary {
		t.Error("the two formats produced different scalars")
	}
	if *fromYAML.Layout.LongSide != *fromJSON.Layout.LongSide {
		t.Error("the two formats produced different numbers")
	}
	if *fromYAML.Layout.Padding.Uniform != *fromJSON.Layout.Padding.Uniform {
		t.Error("the two formats disagreed about the uniform padding spelling")
	}
	if len(fromYAML.Countries) != len(fromJSON.Countries) {
		t.Fatal("the two formats produced different country sets")
	}
	if *fromYAML.Countries["US"].Layout.Width != *fromJSON.Countries["US"].Layout.Width {
		t.Error("the two formats disagreed inside a country overlay")
	}
	if fromYAML.Marker.Custom[0] != fromJSON.Marker.Custom[0] {
		t.Error("the two formats produced different lists")
	}
}

// TestUnknownKeysAreRejectedInYAMLToo is the specific asymmetry this design
// exists to prevent. A YAML decoder with its own unknown-field option would be a
// second implementation of the schema, free to differ.
func TestUnknownKeysAreRejectedInYAMLToo(t *testing.T) {
	for name, doc := range map[string]string{
		"unknown top level": "schema: country-map/v1\nprofil: card\n",
		"wrong case":        "schema: country-map/v1\nlayout:\n  longside: 160\n",
		"unknown nested":    "schema: country-map/v1\ntokens:\n  fillOpacty: 0.1\n",
		"reserved key":      "schema: country-map/v1\nselection:\n  components: [mainland]\n",
	} {
		if _, err := Decode([]byte(doc), FormatYAML); err == nil {
			t.Errorf("%s: accepted in YAML", name)
		}
	}
}

// TestNonFiniteYAMLNumbersAreRefusedAtTheBoundary covers values YAML can spell
// and JSON cannot. Unchecked, an infinity reaches the layout as a dimension and
// a viewBox full of NaN is the first anyone hears about it.
func TestNonFiniteYAMLNumbersAreRefusedAtTheBoundary(t *testing.T) {
	for name, doc := range map[string]string{
		"infinity":     "schema: country-map/v1\nlayout:\n  mode: tight\n  longSide: .inf\n",
		"negative inf": "schema: country-map/v1\nlayout:\n  mode: tight\n  longSide: -.inf\n",
		"nan":          "schema: country-map/v1\nlayout:\n  mode: tight\n  longSide: .nan\n",
	} {
		_, err := Decode([]byte(doc), FormatYAML)
		if err == nil {
			t.Errorf("%s: accepted", name)
			continue
		}
		if !strings.Contains(err.Error(), "longSide") {
			t.Errorf("%s: diagnostic %q does not name the offending key", name, err)
		}
	}
}

// TestMalformedYAMLNamesItsLocation is why goccy is the parser: an author fixing
// a config needs the line, not just "invalid YAML".
func TestMalformedYAMLNamesItsLocation(t *testing.T) {
	_, err := Decode([]byte("schema: country-map/v1\nlayout:\n  mode: tight\n   longSide: 160\n"), FormatYAML)
	if err == nil {
		t.Fatal("malformed YAML was accepted")
	}
	if !strings.Contains(err.Error(), "invalid YAML") {
		t.Fatalf("error does not identify the failure: %v", err)
	}
	// The formatted error carries a line marker; without it the diagnostic is
	// no better than a parse failure with no position.
	if !strings.Contains(err.Error(), "4 |") && !strings.Contains(err.Error(), "[4:") {
		t.Errorf("error does not carry a source location: %v", err)
	}
}

// TestOnlyKnownExtensionsAreAccepted keeps format detection from guessing. A
// sniffed format turns a genuinely malformed file into a confusing report about
// the wrong syntax.
func TestOnlyKnownExtensionsAreAccepted(t *testing.T) {
	dir := t.TempDir()
	body := "schema: country-map/v1\nprofile: card\n"

	for _, name := range []string{"config.yaml", "config.yml"} {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeFile(path); err != nil {
			t.Errorf("%s was refused: %v", name, err)
		}
	}
	jsonPath := filepath.Join(dir, "config.json")
	if err := os.WriteFile(jsonPath, []byte(`{"schema":"country-map/v1"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeFile(jsonPath); err != nil {
		t.Errorf("config.json was refused: %v", err)
	}

	badPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(badPath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	err := DecodeFileErr(t, badPath)
	if err == nil {
		t.Fatal("an unsupported extension was accepted")
	}
	if !strings.Contains(err.Error(), ".yaml") || !strings.Contains(err.Error(), ".json") {
		t.Errorf("diagnostic %q does not name the accepted extensions", err)
	}
}

func DecodeFileErr(t *testing.T, path string) error {
	t.Helper()
	_, err := DecodeFile(path)
	return err
}

// TestEmbeddedPresetsParseAndChain covers the shipped presets and the
// inheritance walk. The chain order is the part worth pinning: walked
// child-first and reversed, so the nearest definition wins.
func TestEmbeddedPresetsParseAndChain(t *testing.T) {
	names, err := PresetNames()
	if err != nil {
		t.Fatalf("the embedded presets do not parse: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("no embedded presets")
	}
	// The spec's own example names this one; removing it would break every
	// config written against the documentation.
	if !contains(names, "site-default") {
		t.Errorf("preset names = %v, want site-default among them", names)
	}

	chain, err := PresetChain("standalone-default")
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 2 || chain[0].Name != "site-default" || chain[1].Name != "standalone-default" {
		t.Fatalf("chain = %v, want the ancestor first", chainNames(chain))
	}

	// Applied in that order, the child's opinion wins and the ancestor's
	// untouched values survive.
	layers := make([]Layer, 0, len(chain))
	for _, preset := range chain {
		layers = append(layers, Layer{Name: "preset " + preset.Name, Settings: preset.Settings})
	}
	resolved, origins := Merge(layers)
	if *resolved.Delivery != "standalone" || origins["delivery"] != "preset standalone-default" {
		t.Errorf("delivery = %v from %q", *resolved.Delivery, origins["delivery"])
	}
	// A value the ancestor sets and the child never mentions survives, and is
	// still attributed to the ancestor.
	if resolved.Tokens == nil || resolved.Tokens.LineCap == nil || *resolved.Tokens.LineCap != "round" {
		t.Fatalf("the ancestor's lineCap did not survive: %+v", resolved.Tokens)
	}
	if origins["tokens.lineCap"] != "preset site-default" {
		t.Errorf("lineCap origin = %q, want the ancestor that set it", origins["tokens.lineCap"])
	}
}

// TestPresetsCarryOnlyWhatTheyChange keeps provenance honest. A preset that
// restates an embedded default makes `explain` attribute the value to the
// preset, sending an author to edit the wrong place.
func TestPresetsCarryOnlyWhatTheyChange(t *testing.T) {
	defaults, err := Defaults()
	if err != nil {
		t.Fatal(err)
	}
	all, err := Presets()
	if err != nil {
		t.Fatal(err)
	}
	for name, preset := range all {
		// Only root presets are compared against the defaults; a child
		// legitimately restates an ancestor's value in order to override it.
		if preset.Extends != nil {
			continue
		}
		baseline := flatten(t, defaults)
		claimed := flatten(t, preset.Settings)
		var redundant []string
		for path, value := range claimed {
			if base, inDefaults := baseline[path]; inDefaults && base == value {
				redundant = append(redundant, path)
			}
		}
		if len(redundant) != 0 {
			t.Errorf("preset %q restates defaults at %v", name, sorted(redundant))
		}
	}
}

// flatten renders settings as dotted path to serialized value, so two layers
// can be compared key by key without hand-listing the schema.
func flatten(t *testing.T, settings Settings) map[string]string {
	t.Helper()
	encoded, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	var generic map[string]any
	if err := json.Unmarshal(encoded, &generic); err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	var walk func(any, string)
	walk = func(value any, prefix string) {
		if object, ok := value.(map[string]any); ok {
			for key, child := range object {
				walk(child, join(prefix, key))
			}
			return
		}
		rendered, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		out[prefix] = string(rendered)
	}
	walk(generic, "")
	return out
}

// TestUnknownPresetNamesTheKnownOnes is REQ-12 for the likeliest preset mistake.
func TestUnknownPresetNamesTheKnownOnes(t *testing.T) {
	_, err := PresetChain("site-defualt")
	if err == nil {
		t.Fatal("an unknown preset was accepted")
	}
	if !strings.Contains(err.Error(), "site-default") {
		t.Errorf("diagnostic %q does not list the known presets", err)
	}
}

// TestEmbeddedPresetsSurviveTheirOwnValidation keeps a shipped preset from being
// the one document nobody checks. A preset that fails the rules it inherits into
// would be a defect every config using it would carry.
func TestEmbeddedPresetsSurviveTheirOwnValidation(t *testing.T) {
	all, err := Presets()
	if err != nil {
		t.Fatal(err)
	}
	vocab := testVocab()
	for name, preset := range all {
		doc := Document{Schema: SchemaVersion, Settings: preset.Settings}
		if err := ValidateDocument(&doc, vocab); err != nil {
			t.Errorf("embedded preset %q does not satisfy the schema it ships with: %v", name, err)
		}
	}
}

func chainNames(chain []Preset) []string {
	out := make([]string, len(chain))
	for i, preset := range chain {
		out[i] = preset.Name
	}
	return out
}
