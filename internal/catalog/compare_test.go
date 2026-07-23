package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompareAddedRemovedChanged(t *testing.T) {
	current := compiledCorpus(t)
	previous := CloneCorpus(current)
	previous.Manifest.Identity = "sha256:previous"
	previous.Manifest.Entities = previous.Manifest.Entities[1:]
	previous.Manifest.Entities[0].Name = "changed"
	previous.Manifest.Entities = append(previous.Manifest.Entities, Entity{Alpha2: "ZZ"})
	result := Compare(previous, current)
	if len(result.Added) != 1 || result.Added[0] != "AD" || len(result.Removed) != 1 || result.Removed[0] != "ZZ" || len(result.Changed) == 0 {
		t.Fatalf("unexpected comparison: %+v", result)
	}
}

func TestPublishCreateOnceAndRejectChangedBytes(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "v1")
	corpus := compiledCorpus(t)
	if err := Publish(directory, corpus); err != nil {
		t.Fatal(err)
	}
	if err := Publish(directory, corpus); err != nil {
		t.Fatalf("same bytes should be idempotent: %v", err)
	}
	before, err := os.ReadFile(filepath.Join(directory, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	changed := CloneCorpus(corpus)
	changed.Manifest.Entities[0].Name += " changed"
	changed.Manifest.Identity, _ = CorpusIdentity(changed.Manifest, changed.Geometries)
	if err := Publish(directory, changed); err == nil || !strings.Contains(err.Error(), "immutable corpus") {
		t.Fatalf("expected immutable publication failure, got %v", err)
	}
	after, err := os.ReadFile(filepath.Join(directory, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("failed publication altered accepted corpus")
	}
}
