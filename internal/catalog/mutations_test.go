package catalog

import (
	"strings"
	"testing"
)

func TestMutationsISOAndProfilesFailClosed(t *testing.T) {
	tests := map[string]struct {
		mutate    func(*Corpus)
		invariant string
	}{
		"missing ISO": {
			mutate:    func(c *Corpus) { c.Manifest.Entities = c.Manifest.Entities[:248] },
			invariant: "iso-count-249",
		},
		"duplicate ISO": {
			mutate:    func(c *Corpus) { c.Manifest.Entities[248] = c.Manifest.Entities[247] },
			invariant: "iso-unique-sorted",
		},
		"removed profile": {
			mutate:    func(c *Corpus) { c.Manifest.Entities[0].Profiles.DeFacto = ProfileResolution{} },
			invariant: "profile-resolution",
		},
		"selected dispute collapsed": {
			mutate: func(c *Corpus) {
				for i := range c.Manifest.Entities {
					if c.Manifest.Entities[i].Alpha2 == "CN" {
						c.Manifest.Entities[i].Profiles.DeFacto = ProfileResolution{IdenticalTo: "un"}
					}
				}
			},
			invariant: "dispute-profile-expectation",
		},
		"empty geometry": {
			mutate:    func(c *Corpus) { c.Geometries[0].Coordinates = nil },
			invariant: "geometry-unique-nonempty",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			corpus := compiledCorpus(t)
			test.mutate(corpus)
			err := Validate(corpus)
			if err == nil || !strings.Contains(err.Error(), "invariant="+test.invariant) {
				t.Fatalf("expected %s failure, got %v", test.invariant, err)
			}
		})
	}
}

func TestMutationProtectedFeatureRemovalNamesEntityAndFeature(t *testing.T) {
	corpus := compiledCorpus(t)
	for i := range corpus.Manifest.Entities {
		if corpus.Manifest.Entities[i].Alpha2 == "MC" {
			corpus.Manifest.Entities[i].Protected = nil
		}
	}
	err := Validate(corpus)
	if err == nil || !strings.Contains(err.Error(), "entity=MC") || !strings.Contains(err.Error(), "Monaco microstate") {
		t.Fatalf("expected named protected coverage failure, got %v", err)
	}
}
