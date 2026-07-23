package catalog

import (
	"bytes"
	"testing"
)

func TestCanonicalJSONStable(t *testing.T) {
	value := Manifest{SchemaVersion: 1, CorpusVersion: 1, Profiles: []string{"un", "de-facto"}, Entities: []Entity{}}
	a, err := CanonicalJSON(value)
	if err != nil {
		t.Fatal(err)
	}
	b, err := CanonicalJSON(value)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) || len(a) == 0 || a[len(a)-1] != '\n' {
		t.Fatal("canonical JSON is not byte-stable newline-terminated JSON")
	}
}

func TestGeometryIdentityChangesWithCoordinates(t *testing.T) {
	a, _ := geometryIdentity(MultiPolygon{{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}})
	b, _ := geometryIdentity(MultiPolygon{{{{0, 0}, {2, 0}, {1, 1}, {0, 0}}}})
	if a == b {
		t.Fatal("different geometry has identical content identity")
	}
}
