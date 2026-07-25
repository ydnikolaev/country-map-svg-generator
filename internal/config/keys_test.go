package config

import (
	"encoding/json"
	"strings"
	"testing"
)

func check(t *testing.T, raw string) error {
	t.Helper()
	var generic any
	if err := json.Unmarshal([]byte(raw), &generic); err != nil {
		t.Fatalf("fixture is not valid JSON: %v", err)
	}
	return CheckKeys(generic)
}

// TestWrongCaseKeysAreRejected is the tooth for the defect that motivated this
// file. encoding/json binds `longside` to LongSide case-insensitively, so
// DisallowUnknownFields alone accepts all of these. Each one would silently
// become part of the v1 contract.
func TestWrongCaseKeysAreRejected(t *testing.T) {
	cases := []string{
		`{"schema":"country-map/v1","layout":{"longside":160}}`,
		`{"schema":"country-map/v1","layout":{"LongSide":160}}`,
		`{"schema":"country-map/v1","layout":{"maxwidth":160}}`,
		`{"schema":"country-map/v1","tokens":{"fillopacity":0.1}}`,
		`{"schema":"country-map/v1","tokens":{"strokeWIDTH":2}}`,
		`{"schema":"country-map/v1","Schema":"country-map/v1"}`,
		`{"schema":"country-map/v1","countries":{"US":{"markerradius":3}}}`,
	}
	for _, raw := range cases {
		if err := check(t, raw); err == nil {
			t.Errorf("accepted a wrong-case key: %s", raw)
		}
	}
	// The documented spelling still works, or the check has simply broken
	// everything rather than fixed anything.
	if err := check(t, `{"schema":"country-map/v1","layout":{"longSide":160,"maxWidth":200},"tokens":{"fillOpacity":0.1,"strokeWidth":2}}`); err != nil {
		t.Errorf("rejected the documented spelling: %v", err)
	}
}

// TestAuthorChosenMapKeysAreNotSchemaKeys guards the other direction. ISO codes
// and profile names are the author's to choose; treating them as schema keys
// would reject every real config.
func TestAuthorChosenMapKeysAreNotSchemaKeys(t *testing.T) {
	raw := `{"schema":"country-map/v1","countries":{"US":{"style":"ghost"},"zz-not-an-iso":{"style":"ghost"}},"profiles":{"card":{"style":"ghost"},"anything":{"style":"ghost"}}}`
	if err := check(t, raw); err != nil {
		t.Errorf("author-chosen map keys were treated as schema keys: %v", err)
	}
	// Whether `zz-not-an-iso` names a real entity is a separate question, decided
	// against the corpus rather than against the schema shape. Key checking must
	// not pre-empt that, or the diagnostic would blame the wrong thing.
}

// TestEveryUnknownKeyIsReportedAtOnce matters for an agent: fixing one key per
// run turns a five-typo config into five round trips.
func TestEveryUnknownKeyIsReportedAtOnce(t *testing.T) {
	err := check(t, `{"schema":"country-map/v1","alpha":1,"beta":2,"layout":{"gamma":3},"countries":{"US":{"delta":4}}}`)
	if err == nil {
		t.Fatal("unknown keys were accepted")
	}
	unknown, ok := err.(*UnknownKeyError)
	if !ok {
		t.Fatalf("error type = %T, want *UnknownKeyError", err)
	}
	want := []string{"alpha", "beta", "countries.US.delta", "layout.gamma"}
	if len(unknown.Keys) != len(want) {
		t.Fatalf("keys = %v, want %v", unknown.Keys, want)
	}
	for i := range want {
		if unknown.Keys[i] != want[i] {
			t.Fatalf("keys = %v, want %v", unknown.Keys, want)
		}
	}
}

// TestSchemaKeysEnumeratesTheContract is the discovery surface, and it is also
// how the reservation in DEC-015 stays checkable: the reserved name must not
// appear anywhere in the enumeration, at any depth.
func TestSchemaKeysEnumeratesTheContract(t *testing.T) {
	keys := SchemaKeys()
	present := map[string]bool{}
	for _, key := range keys {
		present[key] = true
	}
	for _, want := range []string{
		"schema", "extends", "profile", "boundary", "delivery", "style",
		"layout", "layout.mode", "layout.longSide", "layout.width", "layout.padding.top",
		"tokens.fillOpacity", "tokens.advanced.gradient",
		"marker.mode", "animation.enabled", "output.filename",
		"countries.<key>.style", "profiles.<key>.layout.mode",
	} {
		if !present[want] {
			t.Errorf("schema does not enumerate %q", want)
		}
	}
	for _, key := range keys {
		for _, reserved := range ReservedKeys {
			if key == reserved || strings.HasSuffix(key, "."+reserved) {
				t.Errorf("schema key %q claims the reserved name %q", key, reserved)
			}
		}
	}
	// A field tagged "-" is not on the wire and must not be advertised as a key.
	if present["layout.padding.Uniform"] || present["layout.padding.uniform"] {
		t.Error("an off-wire field is advertised as a schema key")
	}
}
