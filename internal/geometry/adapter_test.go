package geometry

import (
	"strings"
	"testing"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
)

// TestAcceptedBoundaryProfilesMatchesWhatTheAdapterTakes keeps the exported
// vocabulary and the switch that implements it from drifting apart. A consumer
// enumerating the list would otherwise offer a value the adapter refuses, or
// hide one it accepts — and the CLI builds its config enum from exactly this
// list.
func TestAcceptedBoundaryProfilesMatchesWhatTheAdapterTakes(t *testing.T) {
	corpus, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	alpha2 := corpus.Manifest.Entities[0].Alpha2

	for _, profile := range AcceptedBoundaryProfiles {
		if _, err := InputFromCatalog(corpus, alpha2, profile, "card"); err != nil {
			t.Errorf("advertised profile %q was refused by the adapter: %v", profile, err)
		}
	}

	// The rejections that matter are the near misses, not nonsense: the
	// hyphenated spelling is the one an author reading the corpus manifest would
	// write, and casing is what a hand-typed flag gets wrong.
	for _, profile := range []string{"de-facto", "De_Facto", "UN", "defacto", ""} {
		if _, err := InputFromCatalog(corpus, alpha2, profile, "card"); err == nil {
			t.Errorf("the adapter accepted %q, which is not in AcceptedBoundaryProfiles", profile)
		} else if !strings.Contains(err.Error(), "unknown profile") {
			t.Errorf("%q was refused for the wrong reason: %v", profile, err)
		}
	}
}

// TestTheCorpusManifestSpellsDeFactoDifferently is not a complaint about the
// corpus — it pins the mismatch so that if either side is ever normalized, the
// consumer guidance that exists because of it gets revisited rather than
// quietly becoming wrong.
func TestTheCorpusManifestSpellsDeFactoDifferently(t *testing.T) {
	corpus, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	manifest := corpus.Manifest.Profiles
	if len(manifest) != 2 {
		t.Fatalf("manifest profiles = %v", manifest)
	}
	same := len(manifest) == len(AcceptedBoundaryProfiles)
	for i := range manifest {
		if same && manifest[i] != AcceptedBoundaryProfiles[i] {
			same = false
		}
	}
	if same {
		t.Fatalf("the manifest and the adapter now agree on %v; the spelling mismatch this guards is gone, so revisit the consumer guidance that works around it", manifest)
	}
}
