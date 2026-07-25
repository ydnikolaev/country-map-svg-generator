package config

import (
	"encoding/json"
	"strings"
	"testing"
)

func f64(v float64) *float64 { return &v }
func boolptr(v bool) *bool   { return &v }

// TestLaterLayersWinAndProvenanceNamesThem is the precedence contract itself.
func TestLaterLayersWinAndProvenanceNamesThem(t *testing.T) {
	resolved, origins := Merge([]Layer{
		{Name: "embedded defaults", Settings: Settings{
			Profile: strptr("card"), Style: strptr("outline"),
			Layout: &Layout{Mode: strptr("tight"), LongSide: f64(128)},
		}},
		{Name: "preset ancestry", Settings: Settings{Style: strptr("bold-soft")}},
		{Name: "document globals", Settings: Settings{Layout: &Layout{LongSide: f64(160)}}},
		{Name: "country override", Settings: Settings{Profile: strptr("hero")}},
	})

	if *resolved.Profile != "hero" || origins["profile"] != "country override" {
		t.Errorf("profile = %v from %q", *resolved.Profile, origins["profile"])
	}
	if *resolved.Style != "bold-soft" || origins["style"] != "preset ancestry" {
		t.Errorf("style = %v from %q", *resolved.Style, origins["style"])
	}
	if *resolved.Layout.LongSide != 160 || origins["layout.longSide"] != "document globals" {
		t.Errorf("longSide = %v from %q", *resolved.Layout.LongSide, origins["layout.longSide"])
	}
	// The value nobody overrode keeps both its value and its original origin.
	if *resolved.Layout.Mode != "tight" || origins["layout.mode"] != "embedded defaults" {
		t.Errorf("mode = %v from %q", *resolved.Layout.Mode, origins["layout.mode"])
	}
}

// TestANestedBlockDoesNotEraseItsSiblings is the bug a naive pointer-swap merge
// would have: a country override that sets one token would drop every inherited
// token, and the result would still look like a valid config.
func TestANestedBlockDoesNotEraseItsSiblings(t *testing.T) {
	resolved, origins := Merge([]Layer{
		{Name: "preset ancestry", Settings: Settings{Tokens: &Tokens{
			Fill: strptr("currentColor"), FillOpacity: f64(0.14), Stroke: strptr("currentColor"),
		}}},
		{Name: "country override", Settings: Settings{Tokens: &Tokens{Fill: strptr("var(--brand)")}}},
	})

	if *resolved.Tokens.Fill != "var(--brand)" || origins["tokens.fill"] != "country override" {
		t.Errorf("fill = %v from %q", *resolved.Tokens.Fill, origins["tokens.fill"])
	}
	if resolved.Tokens.FillOpacity == nil || *resolved.Tokens.FillOpacity != 0.14 {
		t.Fatal("an inherited sibling token was erased by an override that never mentioned it")
	}
	if origins["tokens.fillOpacity"] != "preset ancestry" {
		t.Errorf("fillOpacity origin = %q, want the layer that actually set it", origins["tokens.fillOpacity"])
	}
	if resolved.Tokens.Stroke == nil {
		t.Error("stroke was erased")
	}
}

// TestAnExplicitZeroOverridesANonZero is the reason for pointers everywhere. An
// invisible fill under a visible stroke is a real style, and a merge that
// treated zero as "unset" would silently ignore it.
func TestAnExplicitZeroOverridesANonZero(t *testing.T) {
	resolved, origins := Merge([]Layer{
		{Name: "preset ancestry", Settings: Settings{Tokens: &Tokens{FillOpacity: f64(0.14)}}},
		{Name: "country override", Settings: Settings{Tokens: &Tokens{FillOpacity: f64(0)}}},
	})
	if resolved.Tokens.FillOpacity == nil || *resolved.Tokens.FillOpacity != 0 {
		t.Fatalf("explicit zero did not win: %v", resolved.Tokens.FillOpacity)
	}
	if origins["tokens.fillOpacity"] != "country override" {
		t.Errorf("origin = %q, want country override", origins["tokens.fillOpacity"])
	}
}

// TestSilenceIsNotAnOpinion is the same property from the other side.
func TestSilenceIsNotAnOpinion(t *testing.T) {
	resolved, _ := Merge([]Layer{
		{Name: "preset ancestry", Settings: Settings{Style: strptr("ghost")}},
		{Name: "country override", Settings: Settings{}},
	})
	if resolved.Style == nil || *resolved.Style != "ghost" {
		t.Fatalf("an empty layer erased an inherited value: %v", resolved.Style)
	}
}

// TestListsReplaceRatherThanAppend pins the choice. Appending would make a
// country override unable to shorten an inherited marker list, and would make
// the result depend on how many ancestors happened to mention it.
func TestListsReplaceRatherThanAppend(t *testing.T) {
	resolved, origins := Merge([]Layer{
		{Name: "preset ancestry", Settings: Settings{Marker: &Marker{Custom: []string{"a", "b", "c"}}}},
		{Name: "country override", Settings: Settings{Marker: &Marker{Custom: []string{"a"}}}},
	})
	if len(resolved.Marker.Custom) != 1 || resolved.Marker.Custom[0] != "a" {
		t.Fatalf("custom = %v, want the override's list verbatim", resolved.Marker.Custom)
	}
	if origins["marker.custom"] != "country override" {
		t.Errorf("origin = %q", origins["marker.custom"])
	}
}

// TestOriginsAreStablySorted matters because `explain` output is compared in
// golden fixtures and read by agents.
func TestOriginsAreStablySorted(t *testing.T) {
	_, provenance := Merge([]Layer{{Name: "document globals", Settings: Settings{
		Style: strptr("ghost"), Profile: strptr("card"),
		Animation: &Animation{Enabled: boolptr(true)},
		Layout:    &Layout{Mode: strptr("tight")},
	}}})
	origins := provenance.Origins()
	if len(origins) != 4 {
		t.Fatalf("origins = %v", origins)
	}
	for i := 1; i < len(origins); i++ {
		if origins[i-1].Path > origins[i].Path {
			t.Fatalf("origins are not sorted: %v", origins)
		}
	}
}

// TestPaddingAcceptsBothSpellings covers the union type, including the
// round-trip: a resolved config must re-serialize to something the author
// recognises as what they wrote.
func TestPaddingAcceptsBothSpellings(t *testing.T) {
	var uniform Layout
	if err := json.Unmarshal([]byte(`{"padding":8}`), &uniform); err != nil {
		t.Fatalf("a uniform padding was rejected: %v", err)
	}
	if uniform.Padding.Uniform == nil || *uniform.Padding.Uniform != 8 {
		t.Fatalf("uniform padding = %+v", uniform.Padding)
	}
	encoded, err := json.Marshal(uniform.Padding)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != "8" {
		t.Errorf("uniform padding re-serialized as %s, want 8", encoded)
	}

	var sides Layout
	if err := json.Unmarshal([]byte(`{"padding":{"top":1,"left":2}}`), &sides); err != nil {
		t.Fatalf("per-side padding was rejected: %v", err)
	}
	if sides.Padding.Top == nil || *sides.Padding.Top != 1 || sides.Padding.Left == nil || *sides.Padding.Left != 2 {
		t.Fatalf("per-side padding = %+v", sides.Padding)
	}
	if sides.Padding.Uniform != nil {
		t.Error("per-side padding also set the uniform carrier; the two spellings are not exclusive")
	}

	if err := json.Unmarshal([]byte(`{"padding":"wide"}`), &Layout{}); err == nil {
		t.Error("a string padding was accepted")
	}
}

// TestUniformPaddingMergesEvenThoughItHasNoKey guards the one off-wire field.
// It carries a value the author wrote, so it must merge like any other; it just
// has no path of its own to report.
func TestUniformPaddingMergesEvenThoughItHasNoKey(t *testing.T) {
	resolved, _ := Merge([]Layer{
		{Name: "preset ancestry", Settings: Settings{Layout: &Layout{Padding: &Padding{Uniform: f64(8)}}}},
		{Name: "country override", Settings: Settings{Layout: &Layout{Padding: &Padding{Uniform: f64(32)}}}},
	})
	if resolved.Layout.Padding.Uniform == nil || *resolved.Layout.Padding.Uniform != 32 {
		t.Fatalf("uniform padding did not merge: %+v", resolved.Layout.Padding)
	}
}

// TestMergeCoversEveryField is the tooth. Every test above checks a handful of
// fields by hand, and all of them would still pass if the merge silently skipped
// a field nobody thought to name. This walks the whole schema instead: it sets
// every leaf in one layer and requires the merge to carry all of them.
func TestMergeCoversEveryField(t *testing.T) {
	text, value, flag := "x", 1.0, true
	everything := Settings{
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
		Marker:    &Marker{Mode: &text, Custom: []string{"a"}},
		Animation: &Animation{Enabled: &flag, Hook: &text},
		Output:    &Output{Dir: &text, Filename: &text},
	}

	_, provenance := Merge([]Layer{{Name: "document globals", Settings: everything}})

	// Every wire key the schema enumerates for Settings must have an origin.
	// Container paths (layout, tokens) are not leaves and carry no value of
	// their own, so they are expected to be absent.
	containers := map[string]bool{
		"layout": true, "tokens": true, "tokens.advanced": true,
		"marker": true, "animation": true, "output": true, "layout.padding": true,
	}
	for _, path := range SchemaKeys() {
		switch {
		case containers[path]:
			continue
		// Document-level keys, not part of Settings and so not merged here.
		case path == "schema" || path == "extends":
			continue
		case path == "profiles" || strings.HasPrefix(path, "profiles."):
			continue
		case path == "countries" || strings.HasPrefix(path, "countries."):
			continue
		}
		if _, ok := provenance[path]; !ok {
			t.Errorf("merge did not carry %q; it is in the schema but not in the merge", path)
		}
	}
}
