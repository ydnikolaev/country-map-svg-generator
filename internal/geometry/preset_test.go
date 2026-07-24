package geometry

import (
	"math"
	"strings"
	"testing"
)

func TestAutoQualityVersionedRelativePolicy(t *testing.T) {
	for _, test := range []struct {
		scale, simplification float64
	}{
		{90, 1.08},
		{128, 1.536},
		{160, 1.92},
		{240, 2.88},
		{640, 7.68},
		{1024, 12.288},
		{2048, 24.576},
	} {
		got := AutoQuality(test.scale)
		if math.Abs(got.Simplification-test.simplification) > 1e-12 || !got.Auto {
			t.Fatalf("scale=%g quality=%+v want simplification=%g auto", test.scale, got, test.simplification)
		}
	}
}

func TestPresetPolicyRejectsMissingAndOutOfRangeCoefficients(t *testing.T) {
	for _, raw := range []string{
		`{"version":"v1","quality_policy":{"version":"v1"}}`,
		strings.Replace(string(presetJSON), `"floor_px": 0.80`, `"floor_px": 1.01`, 1),
		strings.Replace(string(presetJSON), `"relative_ratio": 0.012`, `"relative_ratio": 0.009`, 1),
		strings.Replace(string(presetJSON), `"relative_ratio": 0.012`, `"relative_ratio": 0.013`, 1),
		strings.Replace(string(presetJSON), `"quantization": 0.01`, `"quantization": 0.02`, 1),
	} {
		if _, err := parsePresetFile([]byte(raw)); err == nil {
			t.Fatalf("accepted invalid policy: %s", raw)
		}
	}
}

func TestExplicitQualityRetainsAbsoluteBound(t *testing.T) {
	explicit := Quality{Flatness: .2, Simplification: 1.25, Softening: .2, MinimumArea: 1, Quantization: .01}
	if err := validateQuality(explicit); err != nil {
		t.Fatal(err)
	}
	explicit.Simplification = 1.26
	if err := validateQuality(explicit); err == nil {
		t.Fatal("explicit simplification exceeded 1.25 without failure")
	}
	automatic := AutoQuality(640)
	if math.Abs(automatic.Simplification-7.68) > 1e-12 {
		t.Fatalf("automatic simplification=%g", automatic.Simplification)
	}
	if err := validateQuality(automatic); err != nil {
		t.Fatalf("automatic scale-relative quality rejected: %v", err)
	}
}
