package geometry

import (
	"errors"
	"math"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

func TestGridPhaseScheduleIsExactLexicographicHundred(t *testing.T) {
	phases := gridPhaseSchedule()
	if len(phases) != 100 || phases[0] != (GridPhase{}) || phases[99] != (GridPhase{X: .009, Y: .009}) {
		t.Fatalf("schedule endpoints/count: len=%d first=%+v last=%+v", len(phases), phases[0], phases[len(phases)-1])
	}
	for i, phase := range phases {
		want := GridPhase{X: float64(i/10) / 1000, Y: float64(i%10) / 1000}
		if phase != want {
			t.Fatalf("phase[%d]=%+v want=%+v", i, phase, want)
		}
	}
}

func TestGridPhaseComposesTransformViewBoxAndMarkerExactlyOnce(t *testing.T) {
	base := Transform{Scale: 2, TranslateX: 10, TranslateY: 20}
	phase := GridPhase{X: .003, Y: .007}
	unfitted := MultiPolygon{{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}}}
	fitted := materializeGridPhase(unfitted, base)
	composed := composeGridPhase(base, phase)
	shifted := materializeGridPhase(unfitGeometry(fitted, base), composed)
	if shifted[0][0][0] != (Point{X: 10.003, Y: 20.007}) {
		t.Fatalf("phase was not composed exactly once: %+v", shifted[0][0][0])
	}
	viewBox := Bounds{MaxX: 100, MaxY: 50}
	shiftedViewBox := shiftViewBox(viewBox, phase)
	if shiftedViewBox.Width() != viewBox.Width() || shiftedViewBox.Height() != viewBox.Height() {
		t.Fatalf("viewBox dimensions changed: before=%+v after=%+v", viewBox, shiftedViewBox)
	}
	markers, err := projectMarkers([]MarkerInput{{ID: "origin", Lon: 0, Lat: 0}}, projector{}, composed, shiftedViewBox, nil)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(markers[0].X-shifted[0][0][0].X) > 1e-12 || math.Abs(markers[0].Y-shifted[0][0][0].Y) > 1e-12 {
		t.Fatalf("marker/geometry relative placement changed: marker=%+v point=%+v", markers[0], shifted[0][0][0])
	}
}

func TestGridPhaseExhaustionIsTyped(t *testing.T) {
	collapsed := MultiPolygon{{{{0, 0}, {0, 0}, {0, 0}, {0, 0}}}}
	_, err := finalizeGridPhases(
		collapsed, collapsed,
		Transform{Scale: 1}, Bounds{MaxX: 1, MaxY: 1},
		Input{Entity: catalog.Entity{Alpha2: "ZZ"}}, projector{},
		Quality{Simplification: 1, Quantization: .01}, 1, false,
	)
	var pipelineErr *PipelineError
	if !errors.As(err, &pipelineErr) || pipelineErr.Code != ErrTopology || pipelineErr.Field != "grid_phase" {
		t.Fatalf("err=%v, want typed grid-phase exhaustion", err)
	}
}

func TestDiagnosticGridPhasesCannotChangeProductionSchedule(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	in, err := InputFromCatalog(c, "FR", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	record := testProjectedLOD(t, in.Geometry, in.CorpusID, "diagnostic-recipe")
	table := &LODTable{
		Version: "v1", RecipeSHA256: "diagnostic-recipe",
		Compact:             map[string]ProjectedLODGeometry{in.Geometry.ID: record},
		Standard:            map[string]ProjectedLODGeometry{in.Geometry.ID: record},
		CompactMaximumScale: 240, StandardMaximumScale: 700,
	}
	explicit := []GridPhase{{X: .001, Y: .002}, {X: .0011, Y: .0022}}
	rows, err := DiagnosticGridPhases(in, table, "compact", explicit)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != len(explicit) || rows[0].Phase != explicit[0] || rows[1].Phase != explicit[1] {
		t.Fatalf("diagnostic rows=%+v want phases=%+v", rows, explicit)
	}
	production := gridPhaseSchedule()
	if len(production) != 100 || production[0] != (GridPhase{}) || production[99] != (GridPhase{X: .009, Y: .009}) {
		t.Fatalf("diagnostic changed production schedule: %+v", production)
	}
}
