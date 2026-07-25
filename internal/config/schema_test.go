package config

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// decodeStrict is the one decode path both formats funnel through, so
// unknown-field behaviour cannot differ between YAML and JSON. It is defined in
// the test for now and moves into the package when the YAML front end lands.
//
// The exact-case key check runs first and is not redundant with
// DisallowUnknownFields: that option matches field names case-insensitively, so
// it accepts `longside` for `longSide`. CheckKeys is what makes the schema's
// spelling the only accepted spelling.
func decodeStrict(raw []byte, target any) error {
	var generic any
	if err := json.Unmarshal(raw, &generic); err != nil {
		return err
	}
	if err := CheckKeys(generic); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

// TestEmbeddedSettingsFlattenIntoDocumentKeys pins the assumption the whole
// layering design rests on: Settings is embedded anonymously and untagged, so
// its fields are document-level keys rather than nested under a `settings`
// object. If encoding/json ever stopped flattening, every config in the world
// would break, and it would break as a confusing unknown-field error rather than
// as anything that names the cause.
func TestEmbeddedSettingsFlattenIntoDocumentKeys(t *testing.T) {
	var doc Document
	if err := decodeStrict([]byte(`{"schema":"country-map/v1","profile":"card","boundary":"de_facto"}`), &doc); err != nil {
		t.Fatalf("flattened keys did not decode: %v", err)
	}
	if doc.Profile == nil || *doc.Profile != "card" {
		t.Fatalf("profile = %v, want card", doc.Profile)
	}
	if doc.Boundary == nil || *doc.Boundary != "de_facto" {
		t.Fatalf("boundary = %v, want de_facto", doc.Boundary)
	}

	// The other direction: a nested `settings` object must NOT be accepted, or
	// two spellings of the same document would exist and only one would be
	// documented.
	if err := decodeStrict([]byte(`{"schema":"country-map/v1","settings":{"profile":"card"}}`), &Document{}); err == nil {
		t.Error("a nested settings object was accepted; the flattening is ambiguous")
	}
}

// TestUnsetIsDistinguishableFromZero is why every field is a pointer. A layer
// that explicitly sets fillOpacity to 0 (a legitimate value — an invisible fill
// under a visible stroke) must override an earlier non-zero layer, and a layer
// that says nothing must not.
func TestUnsetIsDistinguishableFromZero(t *testing.T) {
	var explicit Document
	if err := decodeStrict([]byte(`{"schema":"country-map/v1","tokens":{"fillOpacity":0}}`), &explicit); err != nil {
		t.Fatal(err)
	}
	if explicit.Tokens == nil || explicit.Tokens.FillOpacity == nil {
		t.Fatal("an explicit zero decoded as unset; layering would silently drop it")
	}
	if *explicit.Tokens.FillOpacity != 0 {
		t.Fatalf("fillOpacity = %v, want 0", *explicit.Tokens.FillOpacity)
	}

	var silent Document
	if err := decodeStrict([]byte(`{"schema":"country-map/v1","tokens":{}}`), &silent); err != nil {
		t.Fatal(err)
	}
	if silent.Tokens == nil {
		t.Fatal("an empty tokens object decoded as absent")
	}
	if silent.Tokens.FillOpacity != nil {
		t.Fatal("an unmentioned field decoded as set")
	}
}

// TestUnknownFieldsAreRejectedAtEveryDepth is REQ-3's core promise. A typo in a
// nested block is the likeliest real mistake and the one most likely to be
// silently ignored by a permissive decoder.
func TestUnknownFieldsAreRejectedAtEveryDepth(t *testing.T) {
	cases := map[string]string{
		"top level":       `{"schema":"country-map/v1","profil":"card"}`,
		"layout block":    `{"schema":"country-map/v1","layout":{"longside":160}}`,
		"tokens block":    `{"schema":"country-map/v1","tokens":{"fillOpacty":0.1}}`,
		"advanced block":  `{"schema":"country-map/v1","tokens":{"advanced":{"mask":"url(#m)"}}}`,
		"country overlay": `{"schema":"country-map/v1","countries":{"US":{"profil":"hero"}}}`,
		"profile overlay": `{"schema":"country-map/v1","profiles":{"card":{"styl":"ghost"}}}`,
	}
	for name, raw := range cases {
		if err := decodeStrict([]byte(raw), &Document{}); err == nil {
			t.Errorf("%s: unknown field was accepted", name)
		}
	}
}

// TestSelectionIsUnclaimedSoDEC015CanReserveIt checks the reservation two ways:
// the name is not spelled by any field today, and a document using it is
// refused. Only the second would be caught by the decoder alone; the first is
// what stops a future v1 field from quietly taking the name.
func TestSelectionIsUnclaimedSoDEC015CanReserveIt(t *testing.T) {
	if err := decodeStrict([]byte(`{"schema":"country-map/v1","selection":{"components":["mainland"]}}`), &Document{}); err == nil {
		t.Error("a reserved-key document was accepted; DEC-015's namespace is not reserved")
	}

	// The reservation is only meaningful while no field claims the name. Marshal
	// an all-fields-set document and confirm the key is absent from the wire
	// form, at any nesting depth.
	value := 1.0
	text := "x"
	flag := true
	full := Document{
		Schema: SchemaVersion, Extends: &text,
		Settings: Settings{
			Profile: &text, Boundary: &text, Delivery: &text, Style: &text,
			Layout: &Layout{
				Mode: &text, LongSide: &value, MaxWidth: &value, MaxHeight: &value,
				Width: &value, Height: &value,
				Padding: &Padding{Top: &value, Right: &value, Bottom: &value, Left: &value},
			},
			Tokens: &Tokens{
				Fill: &text, FillOpacity: &value, Stroke: &text, StrokeOpacity: &value,
				StrokeWidth: &value, LineCap: &text, LineJoin: &text,
				MarkerFill: &text, MarkerStroke: &text, MarkerRadius: &value,
				Advanced: &AdvancedTokens{Gradient: &text, Pattern: &text, Filter: &text},
			},
			Marker:    &Marker{Mode: &text, Custom: []string{"x"}},
			Animation: &Animation{Enabled: &flag, Hook: &text},
			Output:    &Output{Dir: &text, Filename: &text},
		},
		Profiles:  map[string]Settings{"card": {Style: &text}},
		Countries: map[string]Settings{"US": {Style: &text}},
	}
	encoded, err := json.Marshal(full)
	if err != nil {
		t.Fatal(err)
	}
	for _, reserved := range ReservedKeys {
		if strings.Contains(string(encoded), `"`+reserved+`"`) {
			t.Errorf("a v1 field claims the reserved key %q", reserved)
		}
	}
}
