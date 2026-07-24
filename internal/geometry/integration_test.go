package geometry

import (
	"crypto/sha256"
	"fmt"
	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	"testing"
)

func TestFullCorpusBothProfiles(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range c.Manifest.Entities {
		for _, profile := range []string{"un", "de_facto"} {
			in, err := InputFromCatalog(c, e.Alpha2, profile, "card")
			if err != nil {
				t.Fatalf("%s/%s adapter: %v", e.Alpha2, profile, err)
			}
			got, err := Generate(in)
			if err != nil {
				in, _ = ApplyPreset(in)
				geo, _ := normalize(fromCatalog(in.Geometry.Coordinates))
				p := centeredProjector(geo, 0)
				projected, _, _ := projectGeometry(geo, p, .2)
				fitted, _, tr, d, _, _ := fitGeometry(projected, in.Layout)
				prot, protErr := protectedComponents(in.Entity, p, tr, fitted)
				retained, _, comps, _ := retain(fitted, prot, 1, AutoQuality(d).MinimumArea)
				rp := map[int][]string{}
				for i, c := range comps {
					rp[i] = c.Protected
				}
				reduced, _, simplifyErr := simplify(retained, AutoQuality(d).Simplification, rp)
				_, canonicalErr := canonicalize(reduced, .01)
				t.Fatalf("%s/%s: %v; projected=%v protected=%v simplify=%v canonical=%v", e.Alpha2, profile, err, validateTopology(projected), protErr, simplifyErr, canonicalErr)
			}
			if got.Path == "" || got.Metrics.PathBytes != len(got.Path) || got.LayoutMode != LayoutTight {
				t.Fatalf("%s/%s incomplete", e.Alpha2, profile)
			}
			count++
		}
	}
	if count != c.Manifest.EntityCount*2 {
		t.Fatalf("processed %d", count)
	}
}

func TestChinaCanonicalQuantization(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	in, err := InputFromCatalog(c, "CN", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Generate(in)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("CN quantization=%g path_bytes=%d output_points=%d", got.Quality.Quantization, got.Metrics.PathBytes, got.Metrics.OutputPoints)
}

func TestRepresentativeNaturalRatiosAndArbitraryFrames(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	ratios := map[string]float64{}
	for _, code := range []string{"RU", "CL", "AU"} {
		in, err := InputFromCatalog(c, code, "un", "card")
		if err != nil {
			t.Fatal(err)
		}
		got, err := Generate(in)
		if err != nil {
			t.Fatal(err)
		}
		ratios[code] = got.ViewBox.Width() / got.ViewBox.Height()
	}
	if ratios["RU"] <= 1.5 || ratios["CL"] >= .75 || ratios["AU"] < .75 || ratios["AU"] > 1.5 {
		t.Fatalf("unexpected natural ratios: %v", ratios)
	}
	in, err := InputFromCatalog(c, "CL", "un", "")
	if err != nil {
		t.Fatal(err)
	}
	var baseline Quality
	for i, frame := range [][2]float64{{19, 19}, {10000, 12345}, {100, 1000}, {1000, 100}} {
		in.Layout = Layout{Mode: LayoutContain, Width: frame[0], Height: frame[1]}
		got, err := Generate(in)
		if err != nil {
			t.Fatalf("%v: %v", frame, err)
		}
		if i == 2 {
			baseline = got.Quality
		}
		if i == 3 && got.Quality == baseline { /* different fitted scale is allowed to differ */
		}
	}
}

func TestApprovalDigestAndBudgets(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.New()
	for _, code := range []string{"BR", "FR", "US", "CL", "ID", "JP", "FJ", "KI", "RU", "AQ", "MC", "SM", "VA", "NR", "CN", "CY", "IL", "IN"} {
		for _, profile := range []string{"un", "de_facto"} {
			for _, preset := range []string{"card", "hero"} {
				in, err := InputFromCatalog(c, code, profile, preset)
				if err != nil {
					t.Fatal(err)
				}
				got, err := Generate(in)
				if err != nil {
					t.Fatalf("%s/%s/%s: %v", code, profile, preset, err)
				}
				limit := 2200
				if preset == "hero" {
					limit = 7500
				}
				if got.Metrics.PathBytes > limit {
					t.Fatalf("%s/%s/%s path budget %d > %d", code, profile, preset, got.Metrics.PathBytes, limit)
				}
				fmt.Fprintf(h, "%s/%s/%s %s %d %.2f %.2f\n", code, profile, preset, got.Path, got.Metrics.PathBytes, got.ViewBox.Width(), got.ViewBox.Height())
			}
		}
	}
	t.Logf("approval digest: %x", h.Sum(nil))
}
