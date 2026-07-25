package geometry

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
)

func TestLODSelectionIsScaleOnly(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	in, err := InputFromCatalog(c, "FR", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	in.MaxPathBytes = 1 << 20
	record := testProjectedLOD(t, in.Geometry, in.CorpusID, "test-recipe")
	table := &LODTable{Version: "test", RecipeSHA256: "test-recipe", Compact: map[string]ProjectedLODGeometry{in.Geometry.ID: record}, Standard: map[string]ProjectedLODGeometry{in.Geometry.ID: record}, CompactMaximumScale: 240, StandardMaximumScale: 700}
	got, err := SelectLODProvenance(in, table)
	if err != nil {
		t.Fatal(err)
	}
	if got.RequestedTier != "compact" {
		t.Fatalf("%+v", got)
	}
	in.Preset = "hero"
	in.Layout = Layout{}
	got, err = SelectLODProvenance(in, table)
	if err != nil {
		t.Fatal(err)
	}
	if got.RequestedTier != "standard" {
		t.Fatalf("%+v", got)
	}
}

func TestLODPreparationCacheParityAndSourceMutation(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	in, err := InputFromCatalog(c, "FR", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	in.MaxPathBytes = 1 << 20
	table := &LODTable{Version: "test"}
	first, err := SelectLODProvenance(in, table)
	if err != nil {
		t.Fatal(err)
	}
	cached, err := SelectLODProvenance(in, table)
	if err != nil {
		t.Fatal(err)
	}
	uncached, err := SelectLODProvenance(in, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, cached) || !reflect.DeepEqual(first, uncached) {
		t.Fatalf("cached=%+v uncached=%+v first=%+v", cached, uncached, first)
	}
	mutated := in
	mutated.Geometry.Coordinates = append(catalog.MultiPolygon(nil), in.Geometry.Coordinates...)
	mutated.Geometry.Coordinates[0] = append(catalog.Polygon(nil), in.Geometry.Coordinates[0]...)
	mutated.Geometry.Coordinates[0][0] = append(catalog.Ring(nil), in.Geometry.Coordinates[0][0]...)
	mutated.Geometry.Coordinates[0][0][0][0] += .000001
	changed, err := SelectLODProvenance(mutated, table)
	if err != nil {
		t.Fatal(err)
	}
	_ = changed
	if len(table.prepared) != 2 {
		t.Fatal("source mutation aliased cached preparation")
	}
}

func TestLODHardBudgetAdvancesToFittingTier(t *testing.T) {
	dense := testRectangleGeometry("test", 16)
	standard := testRectangleGeometry("test", 1)
	in := Input{
		Entity:       catalog.Entity{Alpha2: "ZZ", Profiles: catalog.BoundaryProfiles{UN: catalog.ProfileResolution{GeometryID: "test"}}},
		Geometry:     dense,
		CorpusID:     "test-corpus",
		Profile:      "un",
		Layout:       Layout{Mode: LayoutTight, LongSide: 100},
		Quality:      Quality{Flatness: .1, Simplification: 1.25, Quantization: .01},
		MaxPathBytes: 100,
	}
	table := &LODTable{
		Version:              "test",
		RecipeSHA256:         "test-recipe",
		Compact:              map[string]ProjectedLODGeometry{"test": testProjectedLOD(t, dense, "test-corpus", "test-recipe")},
		Standard:             map[string]ProjectedLODGeometry{"test": testProjectedLODUsingSourceCenter(t, standard, dense, "test-corpus", "test-recipe")},
		CompactMaximumScale:  1e6,
		StandardMaximumScale: 1e6,
	}
	// The compact candidate is provably over the frozen cap: production code
	// used to commit it as SelectedTier and only discover the overrun later,
	// terminally, inside renderLODRepresentation. That committed-then-hard-fail
	// contract is exactly the defect PLAN-013 replaces. A coarser tier (standard)
	// fits this frozen cap and preserves mandatory fidelity, so selection must
	// advance to it instead of stalling on compact.
	selected, err := SelectLODProvenance(in, table)
	if err != nil {
		t.Fatal(err)
	}
	if selected.SelectedTier != "standard" {
		t.Fatalf("selected=%s fallbacks=%v", selected.SelectedTier, selected.Fallbacks)
	}
	foundTruthfulReason := false
	for _, fallback := range selected.Fallbacks {
		if strings.HasPrefix(fallback, "compact:hard_budget:") {
			foundTruthfulReason = true
		}
	}
	if !foundTruthfulReason {
		t.Fatalf("expected a truthful compact:hard_budget fallback reason, got %v", selected.Fallbacks)
	}
	got, err := GenerateWithLOD(in, table)
	if err != nil {
		t.Fatalf("fallback tier should satisfy the frozen budget, got err=%v", err)
	}
	if got.LOD.SelectedTier != "standard" || got.Metrics.PathBytes > in.MaxPathBytes {
		t.Fatalf("bytes=%d budget=%d tier=%s", got.Metrics.PathBytes, in.MaxPathBytes, got.LOD.SelectedTier)
	}
	// Selection stays a pure, byte-identical function of geometry and the
	// frozen ladder: rendering the winning tier must not re-enter or mutate
	// a later selection.
	again, err := SelectLODProvenance(in, table)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(selected, again) {
		t.Fatalf("budget rendering changed later selection: before=%+v after=%+v", selected, again)
	}
}

func TestProjectedLODBindingRejectsGeographicIdentity(t *testing.T) {
	source := testRectangleGeometry("test", 4)
	record := testProjectedLOD(t, source, "test-corpus", "test-recipe")
	record.Projection.CoordinateSpace = "geographic"
	normalized, err := normalize(fromCatalog(source.Coordinates))
	if err != nil {
		t.Fatal(err)
	}
	in := Input{
		Entity:   catalog.Entity{Alpha2: "ZZ"},
		Geometry: source,
		CorpusID: "test-corpus",
	}
	_, err = bindProjectedLOD(record, in, &LODTable{RecipeSHA256: "test-recipe"}, centeredProjector(normalized, 0), 0)
	if err == nil || !strings.Contains(errorIdentity(err), "invalid_source:lod_coordinate_space") {
		t.Fatalf("err=%v, want typed coordinate-space rejection", err)
	}
}

func TestSourceFinalizationReservesQuantizationAndGuardsFinalDeviation(t *testing.T) {
	source := testRectangleGeometry("test", 32)
	normalized, err := normalize(fromCatalog(source.Coordinates))
	if err != nil {
		t.Fatal(err)
	}
	quality := Quality{Simplification: 1, Quantization: .01}
	got, err := finalizeSourceLOD(
		normalized, nil, quality,
		Transform{Scale: 1}, geometryBounds(normalized),
		Input{Entity: catalog.Entity{Alpha2: "ZZ"}}, projector{}, 1,
	)
	if err != nil {
		t.Fatal(err)
	}
	wantAttempt := quality.Simplification - quality.Quantization/math.Sqrt2
	if math.Abs(got.AttemptedTolerance-wantAttempt) > 1e-12 {
		t.Fatalf("attempted=%g want=%g", got.AttemptedTolerance, wantAttempt)
	}
	if got.SelectedTolerance > got.AttemptedTolerance || got.FinalDeviation > quality.Simplification || !finite(got.FinalDeviation) {
		t.Fatalf("selected=%g attempted=%g final=%g tolerance=%g", got.SelectedTolerance, got.AttemptedTolerance, got.FinalDeviation, quality.Simplification)
	}
	for _, polygon := range got.Canonical {
		for _, ring := range polygon {
			for _, point := range ring {
				if math.Abs(point.X/quality.Quantization-math.Round(point.X/quality.Quantization)) > 1e-9 ||
					math.Abs(point.Y/quality.Quantization-math.Round(point.Y/quality.Quantization)) > 1e-9 {
					t.Fatalf("point=%+v is not on q=%g grid", point, quality.Quantization)
				}
			}
		}
	}
}

func TestProjectionContractAlignsAutomaticRuntimeAndArtifacts(t *testing.T) {
	source := testRectangleGeometry("test", 8)
	in := Input{
		Entity:   catalog.Entity{Alpha2: "ZZ", Profiles: catalog.BoundaryProfiles{UN: catalog.ProfileResolution{GeometryID: "test"}}},
		Geometry: source, CorpusID: "test-corpus", Profile: "un",
		Layout: Layout{Mode: LayoutTight, LongSide: 128}, Quality: Quality{Auto: true},
	}
	got, err := GenerateWithLOD(in, nil)
	if err != nil {
		t.Fatal(err)
	}
	contract := LODProjectionContractV1()
	if got.LOD.Projection != contract || got.LOD.Projection.Flatness != .10 {
		t.Fatalf("runtime projection=%+v want=%+v", got.LOD.Projection, contract)
	}
	if _, err := ProjectLODGeometry(source, .20, contract.CoordinatePrecision); err == nil {
		t.Fatal("maintainer projection accepted a non-contract flatness")
	}
}

func TestProjectionContractPreservesExplicitSourceFlatness(t *testing.T) {
	source := testRectangleGeometry("test", 8)
	in := Input{
		Entity:   catalog.Entity{Alpha2: "ZZ", Profiles: catalog.BoundaryProfiles{UN: catalog.ProfileResolution{GeometryID: "test"}}},
		Geometry: source, CorpusID: "test-corpus", Profile: "un",
		Layout:  Layout{Mode: LayoutTight, LongSide: 128},
		Quality: Quality{Flatness: .20, Simplification: 1.25, Quantization: .01},
	}
	record := testProjectedLOD(t, source, "test-corpus", "test-recipe")
	table := &LODTable{
		Version: "test", RecipeSHA256: "test-recipe",
		Compact:             map[string]ProjectedLODGeometry{"test": record},
		Standard:            map[string]ProjectedLODGeometry{"test": record},
		CompactMaximumScale: 1e6, StandardMaximumScale: 1e6,
	}
	got, err := SelectLODProvenance(in, table)
	if err != nil {
		t.Fatal(err)
	}
	if got.RequestedTier != "source" || got.SelectedTier != "source" ||
		got.Projection.Version != "explicit" || got.Projection.Flatness != .20 {
		t.Fatalf("explicit projection did not deterministically bypass incompatible table: %+v", got)
	}
}

func TestSelectionRenderSeparationAndTypedSourceBudget(t *testing.T) {
	source := testRectangleGeometry("test", 16)
	in := Input{
		Entity:   catalog.Entity{Alpha2: "ZZ", Profiles: catalog.BoundaryProfiles{UN: catalog.ProfileResolution{GeometryID: "test"}}},
		Geometry: source, CorpusID: "test-corpus", Profile: "un",
		Layout:       Layout{Mode: LayoutTight, LongSide: 128},
		Quality:      Quality{Flatness: .10, Simplification: 1.25, Quantization: .01},
		MaxPathBytes: 1,
	}
	selected, err := SelectLODProvenance(in, nil)
	if err != nil {
		t.Fatal(err)
	}
	if selected.ParserRoundTrip {
		t.Fatal("geometry-only selection unexpectedly rendered")
	}
	in.MaxPathBytes = -1
	withoutBudget, err := SelectLODProvenance(in, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(selected, withoutBudget) {
		t.Fatalf("selection depends on budget: one=%+v negative=%+v", selected, withoutBudget)
	}
	in.MaxPathBytes = 1
	_, err = GenerateWithLOD(in, nil)
	var pipelineErr *PipelineError
	if !errors.As(err, &pipelineErr) || pipelineErr.Code != ErrBudget || pipelineErr.Field != "max_path_bytes" {
		t.Fatalf("err=%v, want typed hard budget failure", err)
	}
	after, err := SelectLODProvenance(in, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(selected, after) {
		t.Fatalf("budget failure re-entered or mutated selection: before=%+v after=%+v", selected, after)
	}
}

func testProjectedLOD(t *testing.T, source catalog.Geometry, corpus, recipe string) ProjectedLODGeometry {
	t.Helper()
	record, err := ProjectLODGeometry(source, .10, .000001)
	if err != nil {
		t.Fatal(err)
	}
	record.SourceCorpus = corpus
	record.RecipeSHA256 = recipe
	return record
}

func testProjectedLODUsingSourceCenter(t *testing.T, geometry, centerSource catalog.Geometry, corpus, recipe string) ProjectedLODGeometry {
	t.Helper()
	centerGeometry, err := normalize(fromCatalog(centerSource.Coordinates))
	if err != nil {
		t.Fatal(err)
	}
	projector := centeredProjector(centerGeometry, 0)
	candidate, err := normalize(fromCatalog(geometry.Coordinates))
	if err != nil {
		t.Fatal(err)
	}
	projected, _, err := projectGeometry(candidate, projector, .10)
	if err != nil {
		t.Fatal(err)
	}
	return ProjectedLODGeometry{
		Geometry: catalog.Geometry{ID: geometry.ID, Coordinates: toCatalogGeometry(projected)},
		Projection: LODProjection{
			Version: "v1", CoordinateSpace: "centered_laea", CenterLongitude: deg(projector.lon0), CenterLatitude: deg(projector.lat0),
			BuildRotation: 0, YAxis: "down", Flatness: .10, CoordinatePrecision: .000001,
		},
		SourceCorpus: corpus, RecipeSHA256: recipe,
	}
}

func testRectangleGeometry(id string, steps int) catalog.Geometry {
	ring := make(catalog.Ring, 0, steps*4+1)
	for i := 0; i <= steps; i++ {
		ring = append(ring, catalog.Point{float64(i) / float64(steps), 0})
	}
	for i := 1; i <= steps; i++ {
		ring = append(ring, catalog.Point{1, float64(i) / float64(steps)})
	}
	for i := 1; i <= steps; i++ {
		ring = append(ring, catalog.Point{1 - float64(i)/float64(steps), 1})
	}
	for i := 1; i <= steps; i++ {
		ring = append(ring, catalog.Point{0, 1 - float64(i)/float64(steps)})
	}
	return catalog.Geometry{ID: id, Coordinates: catalog.MultiPolygon{{ring}}}
}
