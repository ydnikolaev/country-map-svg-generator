package catalog

import (
	"bytes"
	"strings"
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

func TestDecodeCorpusBindsEveryAuthoritativePublishedComponent(t *testing.T) {
	base, err := MarshalCorpus(compiledCorpus(t))
	if err != nil {
		t.Fatal(err)
	}
	tests := map[string]func(map[string][]byte){
		"manifest byte": func(files map[string][]byte) { files["manifest.json"] = append([]byte(" "), files["manifest.json"]...) },
		"geometry byte": func(files map[string][]byte) { files["geometries.json"] = append(files["geometries.json"], '\n') },
		"receipt content": func(files map[string][]byte) {
			files["receipts.json"] = bytes.Replace(files["receipts.json"], []byte(`"license":"`), []byte(`"license":"changed `), 1)
		},
		"coverage content": func(files map[string][]byte) {
			files["coverage.json"] = bytes.Replace(files["coverage.json"], []byte(`"capital_records":212`), []byte(`"capital_records":211`), 1)
		},
		"extra file": func(files map[string][]byte) { files["unexpected.json"] = []byte("{}\n") },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			files := cloneFiles(base)
			mutate(files)
			if _, err := DecodeCorpus(files); err == nil {
				t.Fatal("authoritative byte mutation was accepted")
			}
		})
	}
}

func TestDecodeCorpusRejectsReceiptAndCoverageRowMutations(t *testing.T) {
	tests := map[string]func(*Corpus){
		"deleted receipt":      func(c *Corpus) { c.Receipts = c.Receipts[1:] },
		"added receipt":        func(c *Corpus) { c.Receipts = append(c.Receipts, c.Receipts[len(c.Receipts)-1]) },
		"mis-digested receipt": func(c *Corpus) { c.Receipts[0].SHA256 = strings.Repeat("0", 64) },
		"deleted coverage row": func(c *Corpus) { c.Coverage.Protected = c.Coverage.Protected[1:] },
		"added coverage row":   func(c *Corpus) { c.Coverage.Protected = append(c.Coverage.Protected, "ZZ") },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			corpus := compiledCorpus(t)
			mutate(corpus)
			files, err := MarshalCorpus(corpus)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeCorpus(files); err == nil {
				t.Fatal("receipt/coverage mutation was accepted")
			}
		})
	}
}

func cloneFiles(files map[string][]byte) map[string][]byte {
	out := make(map[string][]byte, len(files))
	for name, data := range files {
		out[name] = append([]byte(nil), data...)
	}
	return out
}

func TestGeometryIdentityChangesWithCoordinates(t *testing.T) {
	a, _ := geometryIdentity(MultiPolygon{{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}})
	b, _ := geometryIdentity(MultiPolygon{{{{0, 0}, {2, 0}, {1, 1}, {0, 0}}}})
	if a == b {
		t.Fatal("different geometry has identical content identity")
	}
}
