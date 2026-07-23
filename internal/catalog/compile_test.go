package catalog

import (
	"encoding/json"
	"strings"
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

func TestCapitalOverrideAddReplaceRolesAndPrimary(t *testing.T) {
	base := []feature{capitalFeature("US", 1, "Washington, D.C.", Point{-77.0369, 38.9072}, true, false)}
	tests := map[string]struct {
		features []feature
		policy   capitalOverridePolicy
		check    func(t *testing.T, capitals map[string][]Capital)
	}{
		"zero to one addition": {
			policy: capitalOverridePolicy{Overrides: []capitalOverride{{
				Alpha2: "AQ", Capitals: []Capital{{ID: "AQ-reviewed", Name: "Research Seat", Point: Point{0, -80}, Roles: []string{"administrative"}}}, PrimaryID: "AQ-reviewed",
			}}},
			check: func(t *testing.T, capitals map[string][]Capital) {
				if len(capitals["AQ"]) != 1 || !capitals["AQ"][0].Primary {
					t.Fatalf("addition was not primary: %+v", capitals["AQ"])
				}
			},
		},
		"replacement assigns roles": {
			features: base,
			policy: capitalOverridePolicy{Overrides: []capitalOverride{{
				Alpha2: "US", Replace: true, Capitals: []Capital{{ID: "US-reviewed", Name: "Washington", Point: Point{-77.04, 38.9}, Roles: []string{"constitutional", "administrative"}}}, PrimaryID: "US-reviewed",
			}}},
			check: func(t *testing.T, capitals map[string][]Capital) {
				got := capitals["US"]
				if len(got) != 1 || got[0].ID != "US-reviewed" || len(got[0].Roles) != 2 || !got[0].Primary {
					t.Fatalf("replacement mismatch: %+v", got)
				}
			},
		},
		"additive primary selection": {
			features: base,
			policy: capitalOverridePolicy{Overrides: []capitalOverride{{
				Alpha2: "US", Capitals: []Capital{{ID: "US-secondary", Name: "Secondary", Point: Point{-76, 39}, Roles: []string{"ceremonial"}}}, PrimaryID: "US-secondary",
			}}},
			check: func(t *testing.T, capitals map[string][]Capital) {
				got := capitals["US"]
				primary := ""
				for _, capital := range got {
					if capital.Primary {
						primary = capital.ID
					}
				}
				if len(got) != 2 || primary != "US-secondary" {
					t.Fatalf("additive primary mismatch: %+v", capitals["US"])
				}
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := compileCapitals(test.features, test.policy)
			if err != nil {
				t.Fatal(err)
			}
			test.check(t, got)
		})
	}
}

func TestCapitalOverrideRejectsInvalidRecords(t *testing.T) {
	base := []feature{capitalFeature("US", 1, "Washington, D.C.", Point{-77.0369, 38.9072}, true, false)}
	tests := map[string]capitalOverride{
		"invalid coordinates": {Alpha2: "US", Capitals: []Capital{{ID: "bad", Name: "Bad", Point: Point{181, 0}, Roles: []string{"national"}}}, PrimaryName: "Washington, D.C."},
		"duplicate IDs":       {Alpha2: "US", Capitals: []Capital{{ID: "US-ne-1", Name: "Duplicate", Point: Point{0, 0}, Roles: []string{"national"}}}, PrimaryName: "Washington, D.C."},
		"duplicate roles":     {Alpha2: "US", Replace: true, Capitals: []Capital{{ID: "bad", Name: "Bad", Point: Point{0, 0}, Roles: []string{"national", "national"}}}, PrimaryID: "bad"},
		"dangling primary":    {Alpha2: "US", PrimaryID: "missing"},
	}
	for name, override := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := compileCapitals(base, capitalOverridePolicy{Overrides: []capitalOverride{override}})
			if err == nil || !strings.Contains(err.Error(), "capital-") {
				t.Fatalf("expected capital override rejection, got %v", err)
			}
		})
	}
}

func capitalFeature(code string, id int64, name string, point Point, adm0, alternate bool) feature {
	raw := func(value any) json.RawMessage {
		data, _ := json.Marshal(value)
		return data
	}
	return feature{Properties: map[string]json.RawMessage{
		"iso_a2": raw(code), "ne_id": raw(id), "name": raw(name),
		"longitude": raw(point[0]), "latitude": raw(point[1]),
		"adm0cap": raw(map[bool]int64{true: 1}[adm0]), "capalt": raw(map[bool]int64{true: 1}[alternate]),
	}}
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
