package catalog

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestFullCatalogReproducibility(t *testing.T) {
	first, err := Compile(testDataRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	second, err := Compile(testDataRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	a, err := MarshalCorpus(first)
	if err != nil {
		t.Fatal(err)
	}
	b, err := MarshalCorpus(second)
	if err != nil {
		t.Fatal(err)
	}
	if !equalFiles(a, b) || first.Manifest.Identity != second.Manifest.Identity {
		t.Fatal("two clean compilations differ")
	}
	sourceA, err := EmbeddedSource(first)
	if err != nil {
		t.Fatal(err)
	}
	sourceB, err := EmbeddedSource(second)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(sourceA, sourceB) {
		t.Fatal("embedded bundle generation differs")
	}
}

func TestFullCatalogCommittedAndEmbeddedIdentity(t *testing.T) {
	committed, err := DecodeDirectory(filepath.Join(testDataRoot(t), "corpus", "v1"))
	if err != nil {
		t.Fatal(err)
	}
	embedded, err := Embedded()
	if err != nil {
		t.Fatal(err)
	}
	if committed.Manifest.Identity != embedded.Manifest.Identity {
		t.Fatalf("committed=%s embedded=%s", committed.Manifest.Identity, embedded.Manifest.Identity)
	}
	if committed.Manifest.Identity != compiledCorpus(t).Manifest.Identity {
		t.Fatal("committed corpus differs from current pinned compilation")
	}
}

func TestFullCatalogResolutionReturnsCopy(t *testing.T) {
	corpus := compiledCorpus(t)
	geometry, err := ResolveGeometry(corpus, "US", "un")
	if err != nil {
		t.Fatal(err)
	}
	original := geometry.Coordinates[0][0][0]
	geometry.Coordinates[0][0][0] = Point{}
	again, err := ResolveGeometry(corpus, "US", "un")
	if err != nil {
		t.Fatal(err)
	}
	if again.Coordinates[0][0][0] != original {
		t.Fatal("runtime decoder exposed mutable corpus geometry")
	}
}
