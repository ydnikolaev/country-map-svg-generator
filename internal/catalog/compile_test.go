package catalog

import (
	"sync"
	"testing"
)

var (
	fullOnce   sync.Once
	fullCorpus *Corpus
	fullErr    error
)

func compiledCorpus(t *testing.T) *Corpus {
	t.Helper()
	fullOnce.Do(func() { fullCorpus, fullErr = Compile(testDataRoot(t)) })
	if fullErr != nil {
		t.Fatal(fullErr)
	}
	return CloneCorpus(fullCorpus)
}

func TestCompileRealSourceSet(t *testing.T) {
	corpus := compiledCorpus(t)
	if corpus.Manifest.EntityCount != 249 || len(corpus.Geometries) < 249 {
		t.Fatalf("unexpected corpus size: entities=%d geometries=%d", corpus.Manifest.EntityCount, len(corpus.Geometries))
	}
	if corpus.Coverage.UNExplicit != 249 || corpus.Coverage.DeFactoExplicit+corpus.Coverage.DeFactoIdentical != 249 {
		t.Fatalf("profile coverage is incomplete: %+v", corpus.Coverage)
	}
}

func TestCompileCapitalCasesAndOverrides(t *testing.T) {
	corpus := compiledCorpus(t)
	byCode := map[string]Entity{}
	for _, entity := range corpus.Manifest.Entities {
		byCode[entity.Alpha2] = entity
	}
	if len(byCode["AQ"].Capitals) != 0 {
		t.Fatal("Antarctica must exercise zero-capital support")
	}
	if len(byCode["US"].Capitals) != 1 || !byCode["US"].Capitals[0].Primary {
		t.Fatal("United States must exercise single-capital support")
	}
	if len(byCode["ZA"].Capitals) < 3 {
		t.Fatal("South Africa must exercise multi-capital support")
	}
	var primary string
	for _, capital := range byCode["ZA"].Capitals {
		if capital.Primary {
			primary = capital.Name
		}
	}
	if primary != "Pretoria" {
		t.Fatalf("South Africa reviewed primary override did not apply: %q", primary)
	}
}

func TestCompileProtectedCoverage(t *testing.T) {
	corpus := compiledCorpus(t)
	if len(corpus.Coverage.Protected) != len(requiredProtected) {
		t.Fatalf("protected coverage mismatch: %v", corpus.Coverage.Protected)
	}
}
