package catalog

import (
	"strings"
	"testing"
)

func TestModelDiagnosticIncludesContext(t *testing.T) {
	err := diagnostic("fixture.json", "US", "alpha2", "iso-code", "bad value")
	for _, want := range []string{"source=fixture.json", "entity=US", "field=alpha2", "invariant=iso-code", "bad value"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("diagnostic %q does not contain %q", err, want)
		}
	}
}

func TestModelCoordinateRange(t *testing.T) {
	if !validCoordinate(Point{-180, 90}) || validCoordinate(Point{181, 0}) || validCoordinate(Point{0, -91}) {
		t.Fatal("coordinate range contract is not fail-closed")
	}
}
