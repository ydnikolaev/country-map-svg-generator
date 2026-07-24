package geometry

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

func TestSilhouetteOracleContractAndRasterizerIdentity(t *testing.T) {
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile("silhouette.go")
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	if oracle.RasterizerSHA256 != hex.EncodeToString(sum[:]) {
		t.Fatalf("rasterizer identity=%s want=%x", oracle.RasterizerSHA256, sum)
	}
}

func TestSilhouetteOraclePixelCenterNonzeroMetrics(t *testing.T) {
	reference := MultiPolygon{{{
		{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
	}}}
	identical, err := CompareSilhouettes(reference, reference, 10)
	if err != nil {
		t.Fatal(err)
	}
	if identical.IoU != 1 || identical.Recall != 1 || identical.ReferencePixels != 100 {
		t.Fatalf("identical metrics=%+v", identical)
	}
	half := MultiPolygon{{{
		{0, 0}, {5, 0}, {5, 10}, {0, 10}, {0, 0},
	}}}
	metrics, err := CompareSilhouettes(reference, half, 10)
	if err != nil {
		t.Fatal(err)
	}
	if metrics.IoU != .5 || metrics.Recall != .5 {
		t.Fatalf("half metrics=%+v", metrics)
	}
}

func TestSilhouetteOracleMutationTeeth(t *testing.T) {
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	mutations := map[string]func(*SilhouetteOracle){
		"identity":        func(v *SilhouetteOracle) { v.RasterizerSHA256 = "" },
		"fill":            func(v *SilhouetteOracle) { v.FillRule = "evenodd" },
		"sample":          func(v *SilhouetteOracle) { v.Sampling = "pixel_corner" },
		"grid":            func(v *SilhouetteOracle) { v.Bands[0].GridLongSide++ },
		"threshold":       func(v *SilhouetteOracle) { v.Bands[0].MinimumIoU = 0 },
		"contribution":    func(v *SilhouetteOracle) { v.Bands[0].ContributionThreshold = 0 },
		"candidate_order": func(v *SilhouetteOracle) { v.CandidateOrder = "coarse_to_fine" },
		"tie_break":       func(v *SilhouetteOracle) { v.TieBreak = "input_order" },
		"antialias":       func(v *SilhouetteOracle) { v.Antialias = true },
		"softening":       func(v *SilhouetteOracle) { v.Softening = true },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			got := oracle
			got.Bands = append([]SilhouetteBand(nil), oracle.Bands...)
			mutate(&got)
			raw, marshalErr := json.Marshal(got)
			if marshalErr != nil {
				t.Fatal(marshalErr)
			}
			if _, parseErr := ParseSilhouetteOracle(raw); parseErr == nil {
				t.Fatalf("%s mutation passed", name)
			}
		})
	}
}

func TestProtectedVisibilityRejectsVisibleLossAndAllowsProvenSubscale(t *testing.T) {
	reference := MultiPolygon{
		{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
		{{{20, 0}, {22, 0}, {22, 2}, {20, 2}, {20, 0}}},
	}
	candidate := MultiPolygon{reference[0]}
	in := Input{Entity: catalog.Entity{Alpha2: "ZZ"}}
	visibleBand := SilhouetteBand{GridLongSide: 100, ContributionThreshold: 1}
	if _, err := buildProtectedVisibility(in, reference, candidate, projector{}, visibleBand); err == nil {
		t.Fatal("visible component loss passed")
	}
	subscaleBand := visibleBand
	subscaleBand.ContributionThreshold = 10_000
	got, err := buildProtectedVisibility(in, reference, candidate, projector{}, subscaleBand)
	if err != nil {
		t.Fatal(err)
	}
	if got[1].Retained || !got[1].Subscale || got[1].OmissionReason != "subscale" {
		t.Fatalf("subscale provenance=%+v", got[1])
	}
}

func TestSilhouetteOracleRejectsDominantAndProtectedLoss(t *testing.T) {
	reference := MultiPolygon{
		{{{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5}}},
		{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}},
	}
	band := SilhouetteBand{GridLongSide: 140, ContributionThreshold: 10_000}
	if _, err := buildProtectedVisibility(
		Input{Entity: catalog.Entity{Alpha2: "ZZ"}},
		reference, MultiPolygon{reference[1]}, projector{}, band,
	); err == nil {
		t.Fatal("dominant component loss passed")
	}
	protected := Input{Entity: catalog.Entity{
		Alpha2: "ZZ",
		Protected: []catalog.ProtectedFeature{{
			Name: "protected-island", Anchor: catalog.Point{0, 0}, MinimumParts: 2,
		}},
	}}
	if _, err := buildProtectedVisibility(
		protected, reference, MultiPolygon{reference[0]}, projector{}, band,
	); err == nil {
		t.Fatal("protected anchor loss passed")
	}
}
