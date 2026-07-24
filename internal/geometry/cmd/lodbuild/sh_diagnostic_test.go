package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

type shComponentEvidence struct {
	Index         int
	Digest        string
	Rings, Points int
	Area          float64
	Bounds        [4]float64
}

type shOverlapEvidence struct {
	Candidate int
	IoU       float64
	Recall    float64
}

type shVisibilityEvidence struct {
	Source, Contribution, CandidateContribution, IdentityRank int
	Dominant, Visible, Subscale, Retained                     bool
	IdentityReason, OmissionReason                            string
}

type shResolutionEvidence struct {
	Resolution                         int
	RawType                            string
	Source, Candidate                  []shComponentEvidence
	Overlap                            [][]shOverlapEvidence
	Nearest                            []shOverlapEvidence
	Visibility                         []shVisibilityEvidence
	OracleError                        string
	OutputPolygons, OutputParserParts  int
	OutputParserRings, OutputParserPts int
}

func TestSHCandidateLineageDiagnostic(t *testing.T) {
	const geometryID = "geo-0a336e9d4ea30c295a4611bf7f13a0c6bbd8cad41c0440ad503c1c68c63e4ec7"
	resolutions := []int{512, 320, 256, 224, 192, 160, 128, 127, 126, 125, 124, 123, 122, 121, 120, 116, 112, 108, 104, 100, 96, 80, 64, 48, 32, 24, 16, 12, 8}
	c := embeddedCorpus(t)
	in, err := geometry.InputFromCatalog(c, "SH", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	if in.Geometry.ID != geometryID || len(in.Entity.Protected) != 0 {
		t.Fatalf("unexpected SH source: geometry=%s protected=%d", in.Geometry.ID, len(in.Entity.Protected))
	}
	in, err = geometry.ApplyPreset(in)
	if err != nil {
		t.Fatal(err)
	}
	in.Preset, in.MaxPathBytes = "", 0
	oracle, err := geometry.EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	band := oracle.Bands[0]

	baseRaw, err := os.ReadFile(filepath.Join(sourceRoot(), "internal/geometry/lod/v1.recipe.json"))
	if err != nil {
		t.Fatal(err)
	}
	var base recipe
	if err := json.Unmarshal(baseRaw, &base); err != nil {
		t.Fatal(err)
	}
	v2Raw, err := os.ReadFile(filepath.Join(sourceRoot(), "internal/geometry/lod/v2.recipe.json"))
	if err != nil {
		t.Fatal(err)
	}
	v2Sum := sha256.Sum256(v2Raw)
	v2SHA := hex.EncodeToString(v2Sum[:])
	sourceRecord, err := geometry.ProjectLODGeometry(in.Geometry, base.Projection.Flatness, base.Projection.CoordinatePrecision)
	if err != nil {
		t.Fatal(err)
	}
	rawSourceEvidence := shSummarize(in.Geometry.Coordinates)
	if got := shPointCounts(rawSourceEvidence); !reflect.DeepEqual(got, []int{16, 18, 28, 21}) {
		t.Fatalf("P1 source point lineage changed: %v", got)
	}
	source := sourceRecord.Geometry.Coordinates
	sourceEvidence := shSummarize(source)
	if len(sourceEvidence) != 4 {
		t.Fatalf("source components=%d want=4", len(sourceEvidence))
	}

	var rows []shResolutionEvidence
	for _, resolution := range resolutions {
		dir := filepath.Join(t.TempDir(), strconv.Itoa(resolution))
		record, buildErr := buildDiagnosticCandidate(c, dir, geometryID, resolution, .7, base)
		if buildErr != nil {
			t.Fatal(buildErr)
		}
		record.RecipeSHA256 = v2SHA
		rawType, rawParts, rawRings, rawPoints := shRawOutputShape(t, filepath.Join(dir, "output.geojson"))
		candidate := record.Geometry.Coordinates
		candidateEvidence := shSummarize(candidate)
		if rawParts != len(candidate) {
			t.Fatalf("resolution=%d parser loss raw=%d parsed=%d", resolution, rawParts, len(candidate))
		}
		overlap, nearest := shOverlapMatrix(source, candidate, band.GridLongSide)
		visibility := shVisibility(source, candidate, sourceEvidence, nearest, band)
		evaluation, evaluateErr := geometry.EvaluateSilhouetteCandidate(in, record, v2SHA, band)
		oracleError := ""
		if evaluateErr != nil {
			oracleError = evaluateErr.Error()
		} else {
			visibility = make([]shVisibilityEvidence, len(evaluation.Visibility))
			for i, component := range evaluation.Visibility {
				visibility[i] = shVisibilityEvidence{
					Source: component.SourceOrder, Contribution: component.Contribution,
					CandidateContribution: component.CandidateContribution, IdentityRank: component.IdentityRank,
					Dominant: component.Dominant, Visible: component.Visible, Subscale: component.Subscale,
					Retained: component.Retained, IdentityReason: component.IdentityReason,
					OmissionReason: component.OmissionReason,
				}
			}
		}
		rows = append(rows, shResolutionEvidence{
			Resolution: resolution, RawType: rawType, Source: sourceEvidence, Candidate: candidateEvidence,
			Overlap: overlap, Nearest: nearest, Visibility: visibility, OracleError: oracleError,
			OutputPolygons: rawParts, OutputParserParts: len(candidate),
			OutputParserRings: rawRings, OutputParserPts: rawPoints,
		})
	}

	// Mutation teeth: output order cannot define lineage; parser, projection,
	// absence, and raster failures must remain independently distinguishable.
	probeResolution := resolutions[0]
	probeDir := filepath.Join(t.TempDir(), "mutation")
	probe, err := buildDiagnosticCandidate(c, probeDir, geometryID, probeResolution, .7, base)
	if err != nil {
		t.Fatal(err)
	}
	probe.RecipeSHA256 = v2SHA
	_, originalNearest := shOverlapMatrix(source, probe.Geometry.Coordinates, band.GridLongSide)
	reordered := probe
	reordered.Geometry.Coordinates = append(catalog.MultiPolygon(nil), probe.Geometry.Coordinates...)
	for left, right := 0, len(reordered.Geometry.Coordinates)-1; left < right; left, right = left+1, right-1 {
		reordered.Geometry.Coordinates[left], reordered.Geometry.Coordinates[right] = reordered.Geometry.Coordinates[right], reordered.Geometry.Coordinates[left]
	}
	_, reorderedNearest := shOverlapMatrix(source, reordered.Geometry.Coordinates, band.GridLongSide)
	if !shSameNearestRecall(originalNearest, reorderedNearest) {
		t.Fatal("component reorder changed overlap lineage")
	}
	originalEval, originalErr := geometry.EvaluateSilhouetteCandidate(in, probe, v2SHA, band)
	reorderedEval, reorderedErr := geometry.EvaluateSilhouetteCandidate(in, reordered, v2SHA, band)
	if (originalErr == nil) != (reorderedErr == nil) ||
		(originalErr != nil && originalErr.Error() != reorderedErr.Error()) ||
		(originalErr == nil && !reflect.DeepEqual(originalEval.Metrics, reorderedEval.Metrics)) {
		t.Fatalf("reorder changed oracle result: original=%v reordered=%v", originalErr, reorderedErr)
	}
	absent := probe
	absent.Geometry.Coordinates = append(catalog.MultiPolygon(nil), probe.Geometry.Coordinates...)
	absent.Geometry.Coordinates = append(absent.Geometry.Coordinates[:originalNearest[0].Candidate], absent.Geometry.Coordinates[originalNearest[0].Candidate+1:]...)
	_, absentNearest := shOverlapMatrix(source, absent.Geometry.Coordinates, band.GridLongSide)
	if absentNearest[0].Recall >= originalNearest[0].Recall {
		t.Fatalf("absent-component mutation was not distinguished: before=%g after=%g", originalNearest[0].Recall, absentNearest[0].Recall)
	}
	centerMismatch := probe
	centerMismatch.Projection.CenterLongitude++
	if _, centerErr := geometry.EvaluateSilhouetteCandidate(in, centerMismatch, v2SHA, band); centerErr == nil || !strings.Contains(centerErr.Error(), "lod_center") {
		t.Fatalf("projection-center mutation err=%v", centerErr)
	}
	sourceGeometry := shGeometry(source)
	self, err := geometry.CompareSilhouettes(sourceGeometry[0:1], sourceGeometry[0:1], band.GridLongSide)
	if err != nil || self.IoU != 1 || self.Recall != 1 {
		t.Fatalf("raster self-comparison failed: metrics=%+v err=%v", self, err)
	}

	raw, err := json.Marshal(struct {
		GeometryID                          string
		RawSourceType, ProjectedSourceType string
		Band                                geometry.SilhouetteBand
		RawSource                           []shComponentEvidence
		Rows                                []shResolutionEvidence
	}{geometryID, "MultiPolygon", "MultiPolygon", band, rawSourceEvidence, rows})
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	t.Logf("semantic_sha256=%s semantic=%s", hex.EncodeToString(sum[:]), raw)
}

func shRawOutputShape(t *testing.T, path string) (string, int, int, int) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fc struct {
		Features []struct {
			Geometry struct {
				Type        string
				Coordinates json.RawMessage
			}
		}
	}
	if err := json.Unmarshal(raw, &fc); err != nil || len(fc.Features) != 1 {
		t.Fatalf("raw output parse: features=%d err=%v", len(fc.Features), err)
	}
	var coordinates catalog.MultiPolygon
	switch fc.Features[0].Geometry.Type {
	case "MultiPolygon":
		if err := json.Unmarshal(fc.Features[0].Geometry.Coordinates, &coordinates); err != nil {
			t.Fatal(err)
		}
	case "Polygon":
		var polygon catalog.Polygon
		if err := json.Unmarshal(fc.Features[0].Geometry.Coordinates, &polygon); err != nil {
			t.Fatal(err)
		}
		coordinates = catalog.MultiPolygon{polygon}
	default:
		t.Fatalf("raw geometry type=%q", fc.Features[0].Geometry.Type)
	}
	rings, points := 0, 0
	for _, polygon := range coordinates {
		rings += len(polygon)
		for _, ring := range polygon {
			points += len(ring)
		}
	}
	return fc.Features[0].Geometry.Type, len(coordinates), rings, points
}

func shSummarize(g catalog.MultiPolygon) []shComponentEvidence {
	out := make([]shComponentEvidence, len(g))
	for i, polygon := range g {
		raw, _ := json.Marshal(polygon)
		sum := sha256.Sum256(raw)
		points := 0
		bounds := [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
		area := 0.0
		for ringIndex, ring := range polygon {
			points += len(ring)
			signed := 0.0
			for pointIndex, point := range ring {
				if point[0] < bounds[0] {
					bounds[0] = point[0]
				}
				if point[1] < bounds[1] {
					bounds[1] = point[1]
				}
				if point[0] > bounds[2] {
					bounds[2] = point[0]
				}
				if point[1] > bounds[3] {
					bounds[3] = point[1]
				}
				next := ring[(pointIndex+1)%len(ring)]
				signed += point[0]*next[1] - next[0]*point[1]
			}
			if ringIndex == 0 {
				area += math.Abs(signed) / 2
			} else {
				area -= math.Abs(signed) / 2
			}
		}
		out[i] = shComponentEvidence{
			Index: i, Digest: hex.EncodeToString(sum[:]), Rings: len(polygon),
			Points: points, Area: area, Bounds: bounds,
		}
	}
	return out
}

func shPointCounts(items []shComponentEvidence) []int {
	out := make([]int, len(items))
	for i, item := range items {
		out[i] = item.Points
	}
	return out
}

func shGeometry(g catalog.MultiPolygon) geometry.MultiPolygon {
	out := make(geometry.MultiPolygon, len(g))
	for pi, polygon := range g {
		out[pi] = make(geometry.Polygon, len(polygon))
		for ri, ring := range polygon {
			out[pi][ri] = make(geometry.Ring, len(ring))
			for vi, point := range ring {
				out[pi][ri][vi] = geometry.Point{X: point[0], Y: point[1]}
			}
		}
	}
	return out
}

func shOverlapMatrix(source, candidate catalog.MultiPolygon, grid int) ([][]shOverlapEvidence, []shOverlapEvidence) {
	sourceGeometry, candidateGeometry := shGeometry(source), shGeometry(candidate)
	matrix := make([][]shOverlapEvidence, len(sourceGeometry))
	nearest := make([]shOverlapEvidence, len(sourceGeometry))
	for sourceIndex := range sourceGeometry {
		matrix[sourceIndex] = make([]shOverlapEvidence, len(candidateGeometry))
		nearest[sourceIndex] = shOverlapEvidence{Candidate: -1}
		for candidateIndex := range candidateGeometry {
			metrics, err := geometry.CompareSilhouettes(sourceGeometry[sourceIndex:sourceIndex+1], candidateGeometry[candidateIndex:candidateIndex+1], grid)
			if err != nil {
				panic(err)
			}
			item := shOverlapEvidence{Candidate: candidateIndex, IoU: metrics.IoU, Recall: metrics.Recall}
			matrix[sourceIndex][candidateIndex] = item
			if item.Recall > nearest[sourceIndex].Recall ||
				(item.Recall == nearest[sourceIndex].Recall && (nearest[sourceIndex].Candidate < 0 || candidateIndex < nearest[sourceIndex].Candidate)) {
				nearest[sourceIndex] = item
			}
		}
	}
	return matrix, nearest
}

func shVisibility(source, candidate catalog.MultiPolygon, evidence []shComponentEvidence, nearest []shOverlapEvidence, band geometry.SilhouetteBand) []shVisibilityEvidence {
	sourceGeometry, candidateGeometry := shGeometry(source), shGeometry(candidate)
	bounds := shBounds(sourceGeometry)
	contributions := make([]int, len(evidence))
	candidateContributions := make([]int, len(evidence))
	retained := make([]bool, len(evidence))
	ranked := make([]int, len(evidence))
	for i := range evidence {
		ranked[i] = i
		raster, err := geometry.RasterizeSilhouette(sourceGeometry[i:i+1], bounds, band.GridLongSide)
		if err != nil {
			panic(err)
		}
		for _, filled := range raster.Pixels {
			if filled {
				contributions[i]++
			}
		}
		metrics, err := geometry.CompareSilhouettes(sourceGeometry[i:i+1], candidateGeometry, band.GridLongSide)
		if err != nil {
			panic(err)
		}
		candidateContributions[i], retained[i] = metrics.Intersection, metrics.Recall >= .5
	}
	dominant := 0
	for i := range evidence {
		if evidence[i].Area > evidence[dominant].Area {
			dominant = i
		}
	}
	sort.Slice(ranked, func(i, j int) bool {
		if contributions[ranked[i]] != contributions[ranked[j]] {
			return contributions[ranked[i]] > contributions[ranked[j]]
		}
		if ranked[i] != ranked[j] {
			return ranked[i] < ranked[j]
		}
		return evidence[ranked[i]].Digest < evidence[ranked[j]].Digest
	})
	out := make([]shVisibilityEvidence, len(evidence))
	for i := range evidence {
		visible := contributions[i] >= band.ContributionThreshold || i == dominant
		out[i] = shVisibilityEvidence{
			Source: i, Contribution: contributions[i], CandidateContribution: candidateContributions[i],
			Dominant: i == dominant, Visible: visible, Subscale: !visible, Retained: retained[i],
		}
		if i == ranked[0] {
			out[i].IdentityRank, out[i].IdentityReason = 1, "contribution_rank"
		}
		if !retained[i] && !out[i].Visible {
			if out[i].IdentityRank > 0 {
				out[i].OmissionReason = "protected_subscale"
			} else {
				out[i].OmissionReason = "subscale"
			}
		}
		_ = nearest
	}
	return out
}

func shBounds(g geometry.MultiPolygon) geometry.Bounds {
	bounds := geometry.Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)}
	for _, polygon := range g {
		for _, ring := range polygon {
			for _, point := range ring {
				if point.X < bounds.MinX {
					bounds.MinX = point.X
				}
				if point.Y < bounds.MinY {
					bounds.MinY = point.Y
				}
				if point.X > bounds.MaxX {
					bounds.MaxX = point.X
				}
				if point.Y > bounds.MaxY {
					bounds.MaxY = point.Y
				}
			}
		}
	}
	return bounds
}

func shSameNearestRecall(a, b []shOverlapEvidence) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].IoU != b[i].IoU || a[i].Recall != b[i].Recall {
			return false
		}
	}
	return true
}
