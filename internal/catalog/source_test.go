package catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testDataRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", "data"))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestSourceISORealSnapshot(t *testing.T) {
	entities, err := loadISO(filepath.Join(testDataRoot(t), filepath.FromSlash(isoPath)))
	if err != nil {
		t.Fatal(err)
	}
	if len(entities) != 249 || entities[0].Alpha2 != "AD" || entities[len(entities)-1].Alpha2 != "ZW" {
		t.Fatalf("unexpected ISO snapshot bounds/count: %s..%s (%d)", entities[0].Alpha2, entities[len(entities)-1].Alpha2, len(entities))
	}
}

func TestSourceISOMissingAndDuplicateFail(t *testing.T) {
	path := filepath.Join(testDataRoot(t), filepath.FromSlash(isoPath))
	var snapshot isoSnapshot
	if err := decodeStrict(path, &snapshot); err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*isoSnapshot){
		"missing":   func(s *isoSnapshot) { s.Entities = s.Entities[:248] },
		"duplicate": func(s *isoSnapshot) { s.Entities[248] = s.Entities[247] },
	} {
		t.Run(name, func(t *testing.T) {
			clone := snapshot
			clone.Entities = append([]isoEntity(nil), snapshot.Entities...)
			mutate(&clone)
			target := filepath.Join(t.TempDir(), "iso.json")
			if err := writeJSON(target, clone); err != nil {
				t.Fatal(err)
			}
			_, err := loadISO(target)
			if err == nil || (!strings.Contains(err.Error(), "iso-count-249") && !strings.Contains(err.Error(), "iso-unique")) {
				t.Fatalf("expected ISO invariant diagnostic, got %v", err)
			}
		})
	}
}

func TestSourceReceiptDigestAndCompleteness(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "sources"), 0o755); err != nil {
		t.Fatal(err)
	}
	input := []byte("pinned source\n")
	if err := os.WriteFile(filepath.Join(root, "sources", "input"), input, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(input)
	file := receiptFile{Receipts: []Receipt{{ID: "one", Path: "sources/input", Origin: "https://example.test/input", UpstreamVersion: "1", License: "public", SHA256: hex.EncodeToString(sum[:]), Transformation: "none"}}}
	if err := writeJSON(filepath.Join(root, "sources", "receipts.json"), file); err != nil {
		t.Fatal(err)
	}
	if _, err := loadAndVerifyReceipts(root); err != nil {
		t.Fatal(err)
	}
	file.Receipts[0].SHA256 = strings.Repeat("0", 64)
	if err := writeJSON(filepath.Join(root, "sources", "receipts.json"), file); err != nil {
		t.Fatal(err)
	}
	if _, err := loadAndVerifyReceipts(root); err == nil || !strings.Contains(err.Error(), "source-digest") {
		t.Fatalf("expected source digest failure, got %v", err)
	}
}

func TestSourceUndocumentedInputFails(t *testing.T) {
	if err := requireReceipt([]Receipt{{Path: "sources/a"}}, "sources/b"); err == nil || !strings.Contains(err.Error(), "source-documented") {
		t.Fatalf("expected missing receipt diagnostic, got %v", err)
	}
}

func TestSourceStrictJSONRejectsTrailingValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, []byte(`{"observed_at":"x","authority":"x","entities":[]} {}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var snapshot isoSnapshot
	if err := decodeStrict(path, &snapshot); err == nil || !strings.Contains(err.Error(), "strict-json") {
		t.Fatalf("expected trailing JSON rejection, got %v", err)
	}
}
