package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

func TestLODSpike(t *testing.T) {
	c := embeddedCorpus(t)
	result, err := build(c, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range []tierArtifact{result.CompactArtifact, result.StandardArtifact} {
		if err := validateArtifact(a, c); err != nil {
			t.Fatalf("%s: %v", a.Tier, err)
		}
	}
	if total := len(result.Compact) + len(result.Standard); total > 16<<20 {
		t.Fatalf("artifact brake: %d bytes", total)
	}
	if result.Duration > 30*time.Second {
		t.Fatalf("runtime brake: %s", result.Duration)
	}
	verifyPackageLock(t)
	t.Logf("compact_bytes=%d standard_bytes=%d compact_sha256=%s standard_sha256=%s duration=%s compact_omissions=%d standard_omissions=%d",
		len(result.Compact), len(result.Standard), digest(result.Compact), digest(result.Standard), result.Duration, len(result.CompactArtifact.Omissions), len(result.StandardArtifact.Omissions))
}

// TestLODSpikeFullCorpus characterizes the **superseded** pre-DEC-006 path: the
// v1 tier artifact selected through raw boundary deviation. It is deliberately
// not a coverage gate any more, because the claim it used to assert — that this
// path covers the catalog within budget — is exactly the claim DEC-006 was
// accepted for refuting.
//
// What it still guards is worth keeping. Determinism through the v1 table is
// asserted strictly. And the over-budget set is asserted to be non-empty: that
// set is DEC-006's motivating evidence, so if the v1 path ever started covering
// the catalog, DEC-006's premise would need revisiting rather than being
// silently outlived by a green test.
//
// Live production coverage is TestShippedCatalogServesTheLadder in
// internal/geometry, which sweeps the same catalog through the shipped path.
func TestLODSpikeFullCorpus(t *testing.T) {
	c := embeddedCorpus(t)
	var hardOverruns []string
	result, err := build(c, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	compact := geometryIndex(result.CompactArtifact.Geometries)
	standard := geometryIndex(result.StandardArtifact.Geometries)
	table := &geometry.LODTable{Version: "v1", RecipeSHA256: result.CompactArtifact.RecipeSHA256, Compact: compact, Standard: standard, CompactMaximumScale: 240, StandardMaximumScale: 700}
	type maximum struct {
		Key, Tier string
		Bytes     int
	}
	maxima := map[string]maximum{"card": {}, "hero": {}}
	advisoryOverruns := 0
	cases := 0
	start := time.Now()
	// Deterministic sample, not the whole corpus, and said out loud rather than
	// hidden: every 8th entity plus BD, which is the entity whose card the v1
	// path cannot fit and therefore the one this characterization is about.
	//
	// The full sweep used to be affordable only because it aborted on the first
	// over-budget entity. Now that the over-budget outcome is characterized
	// instead of fatal, the run completes — and completing it means falling back
	// to the full-detail source path for most of the catalog, which pushed this
	// single package past the 10-minute go test timeout. A superseded path does
	// not deserve that share of every `make check`.
	const stride = 8
	sampled := []string{}
	for entityIndex, entity := range c.Manifest.Entities {
		if entityIndex%stride != 0 && entity.Alpha2 != "BD" {
			continue
		}
		sampled = append(sampled, entity.Alpha2)
		for _, profile := range []string{"un", "de_facto"} {
			for _, preset := range []string{"card", "hero"} {
				in, err := geometry.InputFromCatalog(c, entity.Alpha2, profile, preset)
				if err != nil {
					t.Fatal(err)
				}
				hard, advisory := 2500, 2200
				if preset == "hero" {
					hard, advisory = 8000, 7500
				}
				got, err := geometry.GenerateWithLOD(in, table)
				if err != nil {
					// The v1 path cannot fit this entity at all: its tier
					// candidate loses to raw deviation and the source fallback
					// blows the budget. That is the characterized outcome, not a
					// test failure — see this test's doc comment.
					var pipelineErr *geometry.PipelineError
					if errors.As(err, &pipelineErr) && pipelineErr.Code == geometry.ErrBudget {
						hardOverruns = append(hardOverruns, fmt.Sprintf("%s/%s/%s %v", entity.Alpha2, profile, preset, err))
						cases++
						continue
					}
					t.Fatalf("%s/%s/%s: %v", entity.Alpha2, profile, preset, err)
				}
				again, err := geometry.GenerateWithLOD(in, table)
				if err != nil {
					t.Fatal(err)
				}
				if got.Path != again.Path {
					t.Fatalf("%s/%s/%s: nondeterministic path", entity.Alpha2, profile, preset)
				}
				if got.Metrics.PathBytes > hard {
					hardOverruns = append(hardOverruns, fmt.Sprintf("%s/%s/%s tier=%s bytes=%d hard=%d",
						entity.Alpha2, profile, preset, got.LOD.SelectedTier, got.Metrics.PathBytes, hard))
				}
				if got.Metrics.PathBytes > advisory {
					advisoryOverruns++
				}
				if got.Metrics.PathBytes > maxima[preset].Bytes {
					maxima[preset] = maximum{entity.Alpha2 + "/" + profile + "/" + preset, got.LOD.SelectedTier, got.Metrics.PathBytes}
				}
				cases++
			}
		}
	}
	if cases != len(sampled)*4 {
		t.Fatalf("cases=%d want=%d over %d sampled entities", cases, len(sampled)*4, len(sampled))
	}
	if len(hardOverruns) == 0 {
		t.Fatal("no sampled entity exceeded its budget through the v1 tier path. BD is force-included in the sample precisely because it does not fit, so what this guard actually verifies is that BD still does not fit — and if that changed, DEC-006's motivating evidence must be re-examined rather than left silently outlived by a green characterization")
	}
	t.Logf("sampled %d of %d entities (stride %d plus BD), cases=%d runtime=%s card_max=%+v hero_max=%+v advisory_overruns=%d hard_overruns=%d (superseded v1 path; first: %s)",
		len(sampled), c.Manifest.EntityCount, stride, cases, time.Since(start),
		maxima["card"], maxima["hero"], advisoryOverruns, len(hardOverruns), hardOverruns[0])
}

func TestLODSpikeRebuild(t *testing.T) {
	c := embeddedCorpus(t)
	first, err := build(c, filepath.Join(t.TempDir(), "first"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := build(c, filepath.Join(t.TempDir(), "second"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first.Compact, second.Compact) || !bytes.Equal(first.Standard, second.Standard) {
		t.Fatalf("rebuild drift compact=%s/%s standard=%s/%s", digest(first.Compact), digest(second.Compact), digest(first.Standard), digest(second.Standard))
	}
	t.Logf("compact_sha256=%s standard_sha256=%s first=%s second=%s", digest(first.Compact), digest(first.Standard), first.Duration, second.Duration)
}

func TestLODAQSelectionOnly(t *testing.T) {
	c := embeddedCorpus(t)
	buildStart := time.Now()
	result, err := build(c, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	table := &geometry.LODTable{Version: "v1", RecipeSHA256: result.CompactArtifact.RecipeSHA256, Compact: geometryIndex(result.CompactArtifact.Geometries), Standard: geometryIndex(result.StandardArtifact.Geometries), CompactMaximumScale: 240, StandardMaximumScale: 700}
	t.Logf("stage=build duration=%s compact_resolution=%d compact_sha256=%s standard_resolution=%d standard_sha256=%s", time.Since(buildStart), result.CompactArtifact.Resolution, digest(result.Compact), result.StandardArtifact.Resolution, digest(result.Standard))
	for _, preset := range []string{"card", "hero"} {
		t.Run(preset, func(t *testing.T) {
			in, err := geometry.InputFromCatalog(c, "AQ", "un", preset)
			if err != nil {
				t.Fatal(err)
			}
			start := time.Now()
			provenance, err := geometry.SelectLODProvenance(in, table)
			elapsed := time.Since(start)
			if err != nil {
				t.Fatalf("AQ/un/%s duration=%s: %v", preset, elapsed, err)
			}
			t.Logf("case=AQ/un/%s duration=%s requested=%s selected=%s deviation=%g tolerance=%g fallbacks=%v restorations=%d", preset, elapsed, provenance.RequestedTier, provenance.SelectedTier, provenance.MaximumDeviation, provenance.ResolvedTolerance, provenance.Fallbacks, len(provenance.Restored))
		})
	}
}

func TestLODAQRUProjectionAlignedCheckpoint(t *testing.T) {
	c := embeddedCorpus(t)
	buildStart := time.Now()
	result, err := build(c, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	second, err := build(c, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result.Compact, second.Compact) || !bytes.Equal(result.Standard, second.Standard) {
		t.Fatalf("temporary rebuild drift compact=%s/%s standard=%s/%s", digest(result.Compact), digest(second.Compact), digest(result.Standard), digest(second.Standard))
	}
	table := &geometry.LODTable{Version: "v1", RecipeSHA256: result.CompactArtifact.RecipeSHA256, Compact: geometryIndex(result.CompactArtifact.Geometries), Standard: geometryIndex(result.StandardArtifact.Geometries), CompactMaximumScale: 240, StandardMaximumScale: 700}
	buildDuration := time.Since(buildStart)
	t.Logf("stage=build_twice duration=%s first=%s/%s second=%s/%s recipe=%s lock=%s tool=%s", buildDuration,
		digest(result.Compact), digest(result.Standard), digest(second.Compact), digest(second.Standard),
		fileDigest(t, filepath.Join(sourceRoot(), "internal/geometry/lod/v1.recipe.json")),
		fileDigest(t, filepath.Join(sourceRoot(), "internal/geometry/lod/tool/package-lock.json")),
		fileDigest(t, filepath.Join(sourceRoot(), "internal/geometry/lod/tool/build.mjs")))
	type brakeCase struct {
		entity, preset string
	}
	forcedCases := []brakeCase{
		{entity: "RU", preset: "card"},
		{entity: "AQ", preset: "card"},
		{entity: "AQ", preset: "hero"},
	}
	presets, err := geometry.Presets()
	if err != nil {
		t.Fatal(err)
	}
	microStart := time.Now()
	unlimited := make(map[brakeCase]geometry.Result, len(forcedCases))
	var representationFailures []string
	for _, tc := range forcedCases {
		in, inputErr := geometry.InputFromCatalog(c, tc.entity, "un", "")
		if inputErr != nil {
			t.Fatal(inputErr)
		}
		in.Layout = presets[tc.preset].Layout
		start := time.Now()
		got, generationErr := geometry.GenerateWithLOD(in, nil)
		elapsed := time.Since(start)
		if elapsed > 10*time.Second {
			representationFailures = append(representationFailures, fmt.Sprintf("%s/un/%s runtime=%s exceeds 10s", tc.entity, tc.preset, elapsed))
			continue
		}
		if generationErr != nil {
			representationFailures = append(representationFailures, fmt.Sprintf("%s/un/%s runtime=%s error=%v", tc.entity, tc.preset, elapsed, generationErr))
			continue
		}
		repeat, repeatErr := geometry.GenerateWithLOD(in, nil)
		if repeatErr != nil || !reflect.DeepEqual(got, repeat) {
			representationFailures = append(representationFailures, fmt.Sprintf("%s/un/%s nondeterministic repeat err=%v", tc.entity, tc.preset, repeatErr))
			continue
		}
		if guardErr := validateGridPhaseResult(got, presets[tc.preset], 0); guardErr != nil {
			representationFailures = append(representationFailures, fmt.Sprintf("%s/un/%s guards: %v", tc.entity, tc.preset, guardErr))
			continue
		}
		unlimited[tc] = got
		t.Logf("stage=unlimited case=%s/un/%s runtime=%s per_phase=%s requested=%s selected=%s coordinate_space=%s finalization=%s raw_deviation=%g final_deviation=%g tolerance=%g q=%g source_attempted=%g source_selected=%g attempted_phase=(%.3f,%.3f) selected_phase=(%.3f,%.3f) phase_count=%d viewbox=%+v transform=%+v markers=%v components=%v restorations=%d fallbacks=%v points=%d path_bytes=%d parser_roundtrip=%t repeat=true",
			tc.entity, tc.preset, elapsed, elapsed/time.Duration(got.LOD.PhaseAttempts), got.LOD.RequestedTier, got.LOD.SelectedTier, got.LOD.CoordinateSpace, got.LOD.Finalization,
			got.LOD.RawDeviation, got.LOD.FinalDeviation, got.LOD.ResolvedTolerance, got.LOD.Quantization,
			got.LOD.SourceAttemptedTolerance, got.LOD.SourceSelectedTolerance,
			got.LOD.AttemptedPhase.X, got.LOD.AttemptedPhase.Y, got.LOD.SelectedPhase.X, got.LOD.SelectedPhase.Y, got.LOD.PhaseAttempts,
			got.ViewBox, got.Transform, got.Markers, got.Components, len(got.LOD.Restored), got.LOD.Fallbacks,
			got.Metrics.OutputPoints, got.Metrics.PathBytes, got.LOD.ParserRoundTrip)
	}
	if len(representationFailures) > 0 {
		t.Fatalf("unlimited representation brake failed; budget stage skipped: %v", representationFailures)
	}
	type productionCase struct {
		entity, preset, tier string
	}
	for _, tc := range []productionCase{
		{entity: "AQ", preset: "card", tier: "source"},
		{entity: "AQ", preset: "hero", tier: "standard"},
		{entity: "RU", preset: "card", tier: "compact"},
		{entity: "RU", preset: "hero", tier: "standard"},
	} {
		in, inputErr := geometry.InputFromCatalog(c, tc.entity, "un", tc.preset)
		if inputErr != nil {
			t.Fatal(inputErr)
		}
		start := time.Now()
		got, generationErr := geometry.GenerateWithLOD(in, table)
		elapsed := time.Since(start)
		if elapsed > 10*time.Second {
			t.Fatalf("%s production-table %s runtime brake: %s", tc.entity, tc.preset, elapsed)
		}
		if generationErr != nil {
			// Same characterization as TestLODSpikeFullCorpus: through the v1
			// table Russia's card has no representation inside the frozen
			// budget. The bound and the typing are still asserted; the coverage
			// claim moved to the ladder gate.
			var pipelineErr *geometry.PipelineError
			if errors.As(generationErr, &pipelineErr) && pipelineErr.Code == geometry.ErrBudget {
				t.Logf("stage=production_table case=%s/un/%s runtime=%s superseded_v1_path_cannot_fit: %v",
					tc.entity, tc.preset, elapsed, generationErr)
				continue
			}
			t.Fatalf("%s production-table %s runtime=%s: %v", tc.entity, tc.preset, elapsed, generationErr)
		}
		// The tier expectation is deliberately gone. It asserted that the v1
		// artifact wins selection under raw boundary deviation, which DEC-006
		// superseded; production now selects from the committed ladder and is
		// covered by TestShippedCatalogServesTheLadder. What still matters here
		// — and is still asserted — is that the v1 records bind cleanly, that
		// the coordinate space is the contract one, and that the run is bounded
		// and deterministic.
		if strings.Contains(strings.Join(got.LOD.Fallbacks, " "), ":binding:") ||
			got.LOD.CoordinateSpace != "centered_laea" {
			t.Fatalf("%s/%s production binding: %+v", tc.entity, tc.preset, got.LOD)
		}
		if guardErr := validateGridPhaseResult(got, presets[tc.preset], presets[tc.preset].MaxPathBytes); guardErr != nil {
			t.Fatalf("%s/%s production guards: %v", tc.entity, tc.preset, guardErr)
		}
		repeat, repeatErr := geometry.GenerateWithLOD(in, table)
		if repeatErr != nil || !reflect.DeepEqual(got, repeat) {
			t.Fatalf("%s/%s production nondeterminism: %v", tc.entity, tc.preset, repeatErr)
		}
		t.Logf("stage=production_table case=%s/un/%s runtime=%s requested=%s selected=%s projection=%+v phase=(%.3f,%.3f) phase_count=%d raw=%g final=%g tolerance=%g points=%d bytes=%d budget=%d roundtrip=%t repeat=true fallbacks=%v",
			tc.entity, tc.preset, elapsed, got.LOD.RequestedTier, got.LOD.SelectedTier, got.LOD.Projection,
			got.LOD.SelectedPhase.X, got.LOD.SelectedPhase.Y, got.LOD.PhaseAttempts,
			got.LOD.RawDeviation, got.LOD.FinalDeviation, got.LOD.ResolvedTolerance,
			got.Metrics.OutputPoints, got.Metrics.PathBytes, presets[tc.preset].MaxPathBytes, got.LOD.ParserRoundTrip, got.LOD.Fallbacks)
	}
	ruBudget, err := geometry.InputFromCatalog(c, "RU", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	before := unlimited[brakeCase{entity: "RU", preset: "card"}].LOD
	start := time.Now()
	_, budgetErr := geometry.GenerateWithLOD(ruBudget, nil)
	budgetElapsed := time.Since(start)
	var pipelineErr *geometry.PipelineError
	if !errors.As(budgetErr, &pipelineErr) || pipelineErr.Code != geometry.ErrBudget || pipelineErr.Field != "max_path_bytes" {
		t.Fatalf("forced-source budget err=%v, want typed hard budget", budgetErr)
	}
	if budgetElapsed > 10*time.Second {
		t.Fatalf("forced-source budget runtime=%s exceeds 10s", budgetElapsed)
	}
	after, err := geometry.SelectLODProvenance(ruBudget, nil)
	if err != nil {
		t.Fatal(err)
	}
	// ParserRoundTrip is a render-stage flag, not a selection fact:
	// SelectLODProvenance returns before any path is serialized, so it is always
	// false there and always true after a full generate. Comparing the two on
	// that field is a category error — the claim under test is that the budget
	// failure did not re-enter the search, which is about selection only. This
	// assertion was previously unreachable because the tier expectation above
	// failed first.
	before.ParserRoundTrip = false
	after.ParserRoundTrip = false
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("forced-source budget re-entered search: before=%+v after=%+v", before, after)
	}
	t.Logf("stage=source_budget runtime=%s error_identity=%s:%s attempts_before=%d attempts_after=%d selected_tolerance=%g phase=(%.3f,%.3f)",
		budgetElapsed, pipelineErr.Code, pipelineErr.Field, before.PhaseAttempts, after.PhaseAttempts,
		after.SourceSelectedTolerance, after.SelectedPhase.X, after.SelectedPhase.Y)
	if elapsed := time.Since(microStart); elapsed > 30*time.Second {
		t.Fatalf("post-build RU/AQ micro-brake runtime=%s exceeds 30s", elapsed)
	} else {
		t.Logf("stage=micro_brake runtime=%s limit=30s", elapsed)
	}
}

func validateGridPhaseResult(got geometry.Result, preset geometry.Preset, maximumBytes int) error {
	if got.LOD.PhaseAttempts < 1 || got.LOD.PhaseAttempts > 1400 || !got.LOD.ParserRoundTrip {
		return fmt.Errorf("phase/provenance invalid: %+v", got.LOD)
	}
	if got.LOD.FinalDeviation > got.LOD.ResolvedTolerance || got.LOD.Quantization != .01 {
		return fmt.Errorf("deviation/q invalid: %+v", got.LOD)
	}
	if got.ViewBox.MinX != got.LOD.SelectedPhase.X || got.ViewBox.MinY != got.LOD.SelectedPhase.Y {
		return fmt.Errorf("viewBox phase mismatch: viewBox=%+v phase=%+v", got.ViewBox, got.LOD.SelectedPhase)
	}
	if maximumBytes > 0 && got.Metrics.PathBytes > maximumBytes {
		return fmt.Errorf("bytes=%d maximum=%d", got.Metrics.PathBytes, maximumBytes)
	}
	parsed, err := geometry.ParsePath(got.Path)
	if err != nil {
		return err
	}
	rings := 0
	for _, command := range parsed {
		for _, value := range command.Values {
			if math.Abs(value/.01-math.Round(value/.01)) > 1e-9 {
				return fmt.Errorf("off-grid path value %g", value)
			}
		}
		if command.Op == "Z" {
			rings++
		}
	}
	expectedRings := 0
	for _, component := range got.Components {
		expectedRings += component.Rings
	}
	if rings != expectedRings {
		return fmt.Errorf("ring identity changed: path=%d components=%d", rings, expectedRings)
	}
	if got.ViewBox.Width() <= 0 || got.ViewBox.Height() <= 0 || got.Transform.Scale <= 0 || preset.Layout.LongSide <= 0 {
		return fmt.Errorf("invalid viewBox/transform: %+v %+v", got.ViewBox, got.Transform)
	}
	return nil
}

func embeddedCorpus(t *testing.T) *catalog.Corpus {
	t.Helper()
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func digest(raw []byte) string { sum := sha256.Sum256(raw); return hex.EncodeToString(sum[:]) }
func fileDigest(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return digest(raw)
}
func geometryIndex(items []geometry.ProjectedLODGeometry) map[string]geometry.ProjectedLODGeometry {
	out := make(map[string]geometry.ProjectedLODGeometry, len(items))
	for _, g := range items {
		out[g.Geometry.ID] = g
	}
	return out
}

func sourceGeometryIndex(items []catalog.Geometry) map[string]catalog.Geometry {
	out := make(map[string]catalog.Geometry, len(items))
	for _, g := range items {
		out[g.ID] = g
	}
	return out
}

func verifyPackageLock(t *testing.T) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(sourceRoot(), "internal/geometry/lod/tool/package-lock.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lock struct {
		Packages map[string]struct{ Version, Integrity string } `json:"packages"`
	}
	if err := json.Unmarshal(raw, &lock); err != nil {
		t.Fatal(err)
	}
	m := lock.Packages["node_modules/mapshaper"]
	if m.Version != "0.7.44" || m.Integrity != "sha512-3Cx+IABMXt1G28Y8J7oalW5P5VYyt1vHz5FO+KkV51EHnJioo9h9maO+u+4IyCPZ9Mh1hqchgomFmc68GtFwQQ==" {
		t.Fatalf("mapshaper lock mismatch: %+v", m)
	}
}

func validateArtifact(a tierArtifact, c *catalog.Corpus) error {
	if a.SchemaVersion != 1 || a.LODVersion != "v1" || a.SourceCorpus != c.Manifest.Identity || a.ToolVersion != "0.7.44" {
		return fmt.Errorf("identity mismatch")
	}
	if len(a.Geometries) != expectedGeometryCount {
		return fmt.Errorf("coverage=%d", len(a.Geometries))
	}
	source := sourceGeometryIndex(c.Geometries)
	projectionContract := geometry.LODProjectionContractV1()
	seen := map[string]bool{}
	for _, record := range a.Geometries {
		g := record.Geometry
		if seen[g.ID] {
			return fmt.Errorf("duplicate %s", g.ID)
		}
		seen[g.ID] = true
		if _, ok := source[g.ID]; !ok {
			return fmt.Errorf("unknown %s", g.ID)
		}
		if record.SourceCorpus != c.Manifest.Identity || record.RecipeSHA256 != a.RecipeSHA256 ||
			record.Projection.Version != projectionContract.Version ||
			record.Projection.CoordinateSpace != projectionContract.CoordinateSpace ||
			record.Projection.BuildRotation != projectionContract.BuildRotation ||
			record.Projection.YAxis != projectionContract.YAxis ||
			record.Projection.Flatness != projectionContract.Flatness ||
			record.Projection.CoordinatePrecision != projectionContract.CoordinatePrecision {
			return fmt.Errorf("%s projection identity mismatch: %+v", g.ID, record)
		}
		for pi, p := range g.Coordinates {
			if len(p) == 0 {
				return fmt.Errorf("%s polygon %d empty", g.ID, pi)
			}
			for ri, r := range p {
				if len(r) < 4 || r[0] != r[len(r)-1] {
					return fmt.Errorf("%s ring %d/%d open", g.ID, pi, ri)
				}
				for _, q := range r {
					if math.IsNaN(q[0]) || math.IsNaN(q[1]) || math.IsInf(q[0], 0) || math.IsInf(q[1], 0) {
						return fmt.Errorf("%s nonfinite", g.ID)
					}
				}
			}
		}
	}
	expected := make([]string, 0, len(source))
	actual := make([]string, 0, len(seen))
	for id := range source {
		expected = append(expected, id)
	}
	for id := range seen {
		actual = append(actual, id)
	}
	sort.Strings(expected)
	sort.Strings(actual)
	if fmt.Sprint(expected) != fmt.Sprint(actual) {
		return fmt.Errorf("coverage ids differ")
	}
	for _, o := range a.Omissions {
		if o.GeometryID == "" || o.SourcePolygon < 0 || o.Reason == "" {
			return fmt.Errorf("invalid omission %+v", o)
		}
	}
	return nil
}
