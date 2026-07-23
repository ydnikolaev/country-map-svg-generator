package catalog

import (
	"strings"
	"testing"
)

func TestValidationAcceptsCompiledCorpus(t *testing.T) {
	if err := Validate(compiledCorpus(t)); err != nil {
		t.Fatal(err)
	}
}

func TestValidationRejectsDanglingGeometry(t *testing.T) {
	corpus := compiledCorpus(t)
	corpus.Manifest.Entities[0].Profiles.UN.GeometryID = "geo-missing"
	if err := Validate(corpus); err == nil || !strings.Contains(err.Error(), "geometry-reference") {
		t.Fatalf("expected dangling geometry failure, got %v", err)
	}
}

func TestValidationRejectsInvalidCapitalPrimaryAndCoordinate(t *testing.T) {
	for name, mutate := range map[string]func(*Corpus){
		"primary": func(c *Corpus) {
			for i := range c.Manifest.Entities {
				if len(c.Manifest.Entities[i].Capitals) > 0 {
					c.Manifest.Entities[i].Capitals[0].Primary = false
					return
				}
			}
		},
		"coordinate": func(c *Corpus) {
			for i := range c.Manifest.Entities {
				if len(c.Manifest.Entities[i].Capitals) > 0 {
					c.Manifest.Entities[i].Capitals[0].Point[0] = 181
					return
				}
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := Validate(func() *Corpus { c := compiledCorpus(t); mutate(c); return c }())
			if err == nil || !strings.Contains(err.Error(), "capital-") {
				t.Fatalf("expected capital invariant failure, got %v", err)
			}
		})
	}
}
