package catalog

import (
	"bytes"
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
	changed.Manifest.Identity, _ = CorpusIdentity(changed.Manifest, changed.Geometries, changed.Receipts, changed.Coverage)
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

func TestDecodeDirectoryRejectsUnexpectedPublishedFile(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "v1")
	if err := Publish(directory, compiledCorpus(t)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "extra.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeDirectory(directory); err == nil || !strings.Contains(err.Error(), "exactly") {
		t.Fatalf("expected exact directory coverage failure, got %v", err)
	}
}

func TestPublishGeneratedSetStagesBeforeReplacingAnyOutput(t *testing.T) {
	root := t.TempDir()
	corpusDir := filepath.Join(root, "corpus")
	report := filepath.Join(root, "report.json")
	embedded := filepath.Join(root, "embedded_gen.go")
	if err := os.MkdirAll(corpusDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(corpusDir, "old"), []byte("old corpus"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(report, []byte("old report"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(embedded, []byte("old embedded"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := publishGeneratedSet(corpusDir, report, embedded, compiledCorpus(t), func(stages generatedStages) error {
		return os.WriteFile(stages.Report, []byte("corrupt staged report"), 0o644)
	})
	if err == nil {
		t.Fatal("expected staged validation failure")
	}
	assertFileContents(t, filepath.Join(corpusDir, "old"), "old corpus")
	assertFileContents(t, report, "old report")
	assertFileContents(t, embedded, "old embedded")
}

func TestPublishGeneratedSetReplacesCompleteOutputSet(t *testing.T) {
	root := t.TempDir()
	corpusDir := filepath.Join(root, "corpus")
	report := filepath.Join(root, "report.json")
	embedded := filepath.Join(root, "embedded_gen.go")
	for _, path := range []string{report, embedded} {
		if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(corpusDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(corpusDir, "old"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := PublishGeneratedSet(corpusDir, report, embedded, compiledCorpus(t)); err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeDirectory(corpusDir); err != nil {
		t.Fatalf("published corpus is invalid: %v", err)
	}
	wantReport, _ := CanonicalJSON(Compare(nil, compiledCorpus(t)))
	gotReport, _ := os.ReadFile(report)
	if !bytes.Equal(wantReport, gotReport) {
		t.Fatal("report was not replaced")
	}
	wantEmbedded, _ := EmbeddedSource(compiledCorpus(t))
	gotEmbedded, _ := os.ReadFile(embedded)
	if !bytes.Equal(wantEmbedded, gotEmbedded) {
		t.Fatal("embedded output was not replaced")
	}
}

func assertFileContents(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("%s changed: %q", path, got)
	}
}
