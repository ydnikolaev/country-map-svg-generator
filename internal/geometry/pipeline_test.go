package geometry

import (
	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	"reflect"
	"strings"
	"testing"
)

func squareInput() Input {
	return Input{
		Entity:   catalog.Entity{Alpha2: "ZZ", Profiles: catalog.BoundaryProfiles{UN: catalog.ProfileResolution{GeometryID: "zz"}}},
		Geometry: catalog.Geometry{ID: "zz", Coordinates: catalog.MultiPolygon{{{{-10, -5}, {10, -5}, {10, 5}, {-10, 5}, {-10, -5}}}}},
		Profile:  "un", CorpusID: "test", Layout: Layout{Mode: LayoutContain, Width: 200, Height: 100},
		Quality: Quality{Auto: true},
		Markers: []MarkerInput{{ID: "center", Lon: 0, Lat: 0}},
	}
}

func TestPipelineDeterministicPresentationFree(t *testing.T) {
	in := squareInput()
	a, err := Generate(in)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Generate(in)
	if err != nil {
		t.Fatal(err)
	}
	if a.Path != b.Path || a.EffectiveScale != b.EffectiveScale {
		t.Fatal("equal input was not deterministic")
	}
	if a.LayoutMode != LayoutContain || a.Transform.Scale <= 0 || a.NaturalAspect <= 1 {
		t.Fatalf("missing CTR-002 fields: %+v", a)
	}
	if len(a.Markers) != 1 || a.Markers[0].Anomaly != "" {
		t.Fatalf("bad marker: %+v", a.Markers)
	}
	for _, bad := range []string{"fill", "stroke", "class", "<svg"} {
		if strings.Contains(a.Path, bad) {
			t.Fatalf("presentation leaked: %s", bad)
		}
	}
}

func TestProtectedFeatureCannotDisappear(t *testing.T) {
	in := squareInput()
	in.Entity.Protected = []catalog.ProtectedFeature{{Name: "island", Anchor: catalog.Point{80, 80}, MinimumParts: 1}}
	_, err := Generate(in)
	if err == nil || !strings.Contains(err.Error(), "protected_feature_loss") {
		t.Fatalf("expected protected loss, got %v", err)
	}
}

func TestMarkerAnomalyIsVisible(t *testing.T) {
	in := squareInput()
	in.Markers = []MarkerInput{{ID: "far", Lon: 170, Lat: 0}}
	got, err := Generate(in)
	if err != nil {
		t.Fatal(err)
	}
	if got.Markers[0].Anomaly == "" || len(got.Diagnostics) == 0 {
		t.Fatalf("suppressed anomaly: %+v", got.Markers)
	}
}

func TestPublicGenerateNilTableMatchesIsolatedSourceInjection(t *testing.T) {
	in := squareInput()
	got, err := Generate(in)
	if err != nil {
		t.Fatal(err)
	}
	injected, err := GenerateWithLOD(in, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, injected) {
		t.Fatalf("public nil seam differs from isolated source injection:\npublic=%+v\ninjected=%+v", got, injected)
	}
	if got.LOD.RequestedTier != "source" || got.LOD.SelectedTier != "source" || !got.LOD.ParserRoundTrip {
		t.Fatalf("public seam did not finalize source: %+v", got.LOD)
	}
}

func TestPublicGeneratePreservesBoundaryValidations(t *testing.T) {
	missing := squareInput()
	missing.Entity.Alpha2 = ""
	if _, err := Generate(missing); err == nil || !strings.Contains(err.Error(), "entity and geometry are required") {
		t.Fatalf("missing input err=%v", err)
	}
	mismatched := squareInput()
	mismatched.Geometry.ID = "other"
	if _, err := Generate(mismatched); err == nil || !strings.Contains(err.Error(), "profile references") {
		t.Fatalf("profile binding err=%v", err)
	}
	override := squareInput()
	override.Override = &Override{ISO: "YY"}
	if _, err := Generate(override); err == nil || !strings.Contains(err.Error(), "override targets") {
		t.Fatalf("override target err=%v", err)
	}
}
