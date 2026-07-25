package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

const (
	diagnosticBaselineCommit  = "2a7d45953ac857af185a4825ceea7097e0af3a40"
	diagnosticBaselineTree    = "a61a0f87122e9094f444fd032dc54a2c04ad3c4b"
	diagnosticScratchBefore   = "0659b28f949f3f0c8a7ca1021ce0dac9eedab2348035c687709c27c926d73c24"
	diagnosticSerializerHash  = "c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345"
	diagnosticSerializerT     = "0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f"
	// Re-recorded when P3/T1 added the CLI's dependencies (cobra, testscript)
	// and again at P3/T2 when the configuration layer added its YAML parser.
	// These two fields are the coarsest in the record: they hash the whole module
	// graph, including dependencies the diagnostic never touches, so they move for
	// reasons unrelated to anything it computes. Nothing was weakened by
	// re-recording them — the diagnostic's outputs are guarded independently by
	// the frozen expectations in evaluatePredicates, and the full mapshaper ladder
	// rebuild reproduced byte-identically under the new graph. Tracked as
	// WKI-1DA58E0FE741.
	diagnosticGoModHash       = "a4253b604e0235a2fa7c2431b903fb5b06cef8d4e778f2e530d98891d106d09d"
	diagnosticGoSumHash       = "7f6f01cbecb673081becf5de9595bfe57a7fba218e27987bd4486566d7f837b4"
	diagnosticCommandIdentity = "GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go run ./internal/geometry/cmd/lodbuild -diagnostic -semantic-out=<semantic-out> -timing-out=<timing-out>"
)

var diagnosticSourcePaths = []string{
	"internal/geometry/adapter.go",
	"internal/geometry/cmd/lodbuild/diagnostic.go",
	"internal/geometry/cmd/lodbuild/main.go",
	"internal/geometry/diagnostic.go",
	"internal/geometry/fit.go",
	"internal/geometry/gridphase.go",
	"internal/geometry/lod.go",
	"internal/geometry/marker.go",
	"internal/geometry/math.go",
	"internal/geometry/model.go",
	"internal/geometry/normalize.go",
	"internal/geometry/override.go",
	"internal/geometry/path.go",
	"internal/geometry/pipeline.go",
	"internal/geometry/preset.go",
	"internal/geometry/projection.go",
	"internal/geometry/retain.go",
	"internal/geometry/serialize.go",
	"internal/geometry/simplify.go",
	"internal/geometry/soften.go",
}

type sourceIdentity struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type diagnosticIdentities struct {
	P1Corpus           string                 `json:"p1_corpus"`
	BaselineCommit     string                 `json:"baseline_commit"`
	BaselineTree       string                 `json:"baseline_tree"`
	ScratchBefore      string                 `json:"scratch_inventory_before"`
	Projection         geometry.LODProjection `json:"projection"`
	RecipeSHA256       string                 `json:"recipe_sha256"`
	MapshaperVersion   string                 `json:"mapshaper_version"`
	MapshaperIntegrity string                 `json:"mapshaper_integrity"`
	MapshaperShasum    string                 `json:"mapshaper_shasum"`
	LockfileSHA256     string                 `json:"lockfile_sha256"`
	ToolSHA256         string                 `json:"tool_sha256"`
	GoVersion          string                 `json:"go_version"`
	GoSources          []sourceIdentity       `json:"go_sources"`
	GoSourcesSHA256    string                 `json:"go_sources_sha256"`
	GoModSHA256        string                 `json:"go_mod_sha256"`
	GoSumSHA256        string                 `json:"go_sum_sha256"`
	SerializerSHA256   string                 `json:"serializer_sha256"`
	SerializerTestSHA  string                 `json:"serializer_test_sha256"`
	Command            string                 `json:"command"`
}

type parameterRow struct {
	Value        float64  `json:"value"`
	Outcome      string   `json:"outcome"`
	RawDeviation *float64 `json:"raw_deviation,omitempty"`
	Failure      string   `json:"failure,omitempty"`
	SelectedTier string   `json:"selected_tier"`
}

type phaseScan struct {
	Name          string                        `json:"name"`
	Rows          []geometry.DiagnosticPhaseRow `json:"rows"`
	AcceptedCount int                           `json:"accepted_count"`
	MinimumFinite float64                       `json:"minimum_finite_final_deviation"`
}

type catalogRow struct {
	Entity, Profile, Preset, GeometryID string
	EffectiveScale                      float64
	RequestedTier, SelectedTier         string
	Topology, Protection                bool
	Omissions                           []geometry.Removal
	RawDeviation, FinalDeviation        float64
	Points, PathBytes                   int
	OverHistoricalMaximum, OverP2Cap    bool
	Fallbacks                           []string
	OutputSHA256                        string
}

type aggregateMaximum struct {
	Entity    string `json:"entity"`
	Profile   string `json:"profile"`
	Tier      string `json:"tier"`
	PathBytes int    `json:"path_bytes"`
}

type diagnosticAggregates struct {
	ParameterRows       int               `json:"parameter_rows"`
	TypedFailures       int               `json:"typed_failures"`
	CatalogRows         int               `json:"catalog_rows"`
	Selected            map[string]int    `json:"selected"`
	HistoricalOverruns  map[string]int    `json:"historical_max_overruns"`
	MaximumCard         aggregateMaximum  `json:"maximum_card"`
	MaximumHero         aggregateMaximum  `json:"maximum_hero"`
	AnyParameterFitsAll bool              `json:"any_parameter_candidate_fits_all_p2_caps"`
	ParameterCapProof   parameterCapProof `json:"parameter_cap_proof"`
}

type parameterCapProof struct {
	Method         string `json:"method"`
	CandidateScope string `json:"candidate_scope"`
	Entity         string `json:"entity"`
	Profile        string `json:"profile"`
	Preset         string `json:"preset"`
	RequestedTier  string `json:"requested_tier"`
	PathBytes      int    `json:"path_bytes"`
	P2Cap          int    `json:"p2_cap"`
}

type diagnosticPredicates struct {
	InputHashOrCountDrift     bool `json:"input_hash_or_count_drift"`
	ExpectedNumericDrift      bool `json:"expected_numeric_drift"`
	TypedFailureMismatch      bool `json:"typed_failure_mismatch"`
	PhaseAcceptedAtOrBelow    bool `json:"phase_accepted_at_or_below_1_536"`
	AggregateOrMaximumDrift   bool `json:"aggregate_or_maximum_drift"`
	ParameterCandidateFitsAll bool `json:"parameter_candidate_fits_all_p2_caps"`
}

type semanticDiagnostic struct {
	SchemaVersion int                  `json:"schema_version"`
	Kind          string               `json:"kind"`
	Identities    diagnosticIdentities `json:"identities"`
	Resolution    []parameterRow       `json:"resolution_rows"`
	Weighting     []parameterRow       `json:"weighting_rows"`
	Phases        []phaseScan          `json:"phase_scans"`
	Catalog       []catalogRow         `json:"catalog_rows"`
	Aggregates    diagnosticAggregates `json:"aggregates"`
	Predicates    diagnosticPredicates `json:"contradiction_predicates"`
	Verdict       string               `json:"verdict"`
}

type timingReceipt struct {
	SchemaVersion  int              `json:"schema_version"`
	Kind           string           `json:"kind"`
	SemanticSHA256 string           `json:"semantic_sha256"`
	DurationMS     int64            `json:"duration_ms"`
	StagesMS       map[string]int64 `json:"stages_ms"`
}

func runDiagnostic(c *catalog.Corpus, semanticPath, timingPath string) error {
	started := time.Now()
	stages := map[string]int64{}
	stageStart := time.Now()
	identities, r, err := loadDiagnosticIdentities(c)
	if err != nil {
		return err
	}
	stages["identities"] = time.Since(stageStart).Milliseconds()

	rootTemp, err := os.MkdirTemp("", "country-map-run010-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(rootTemp)

	aqInput, err := geometry.InputFromCatalog(c, "AQ", "un", "")
	if err != nil {
		return err
	}
	aqGeometryID := aqInput.Geometry.ID
	stageStart = time.Now()
	standard, err := buildDiagnosticCandidate(c, filepath.Join(rootTemp, "standard"), aqGeometryID, 256, .7, r)
	if err != nil {
		return err
	}
	resolutions := []int{64, 65, 66, 68, 72, 80, 96, 128, 160, 192, 224, 256, 320, 512, 1024}
	resolutionRows := make([]parameterRow, 0, len(resolutions))
	var compact64 geometry.ProjectedLODGeometry
	for _, resolution := range resolutions {
		record, buildErr := buildDiagnosticCandidate(c, filepath.Join(rootTemp, "resolution-"+strconv.Itoa(resolution)), aqGeometryID, resolution, .7, r)
		if buildErr != nil {
			return buildErr
		}
		if resolution == 64 {
			compact64 = record
		}
		resolutionRows = append(resolutionRows, diagnosticSelectionRow(c, r, record, standard, float64(resolution)))
	}
	weightings := []float64{0, .1, .3, .5, .7, .9, 1}
	weightRows := make([]parameterRow, 0, len(weightings))
	for _, weighting := range weightings {
		record, buildErr := buildDiagnosticCandidate(c, filepath.Join(rootTemp, "weight-"+strconv.FormatFloat(weighting, 'f', -1, 64)), aqGeometryID, 64, weighting, r)
		if buildErr != nil {
			return buildErr
		}
		weightRows = append(weightRows, diagnosticSelectionRow(c, r, record, standard, weighting))
	}
	stages["parameter_scans"] = time.Since(stageStart).Milliseconds()

	stageStart = time.Now()
	table64 := diagnosticTable(r, compact64, standard)
	phaseScans := make([]phaseScan, 0, 3)
	phaseScans = append(phaseScans, runPhaseScan(c, table64, "coarse", coarsePhases()))
	phaseScans = append(phaseScans, runPhaseScan(c, table64, "local_a", localPhases(8, 12, 18, 22, 10000)))
	phaseScans = append(phaseScans, runPhaseScan(c, table64, "local_b", localPhases(105, 115, 215, 225, 100000)))
	stages["phase_scans"] = time.Since(stageStart).Milliseconds()

	stageStart = time.Now()
	fullBuild, err := build(c, filepath.Join(rootTemp, "catalog-build"))
	if err != nil {
		return err
	}
	fullTable := &geometry.LODTable{
		Version: "v1", RecipeSHA256: fullBuild.CompactArtifact.RecipeSHA256,
		Compact:             projectedGeometryIndex(fullBuild.CompactArtifact.Geometries),
		Standard:            projectedGeometryIndex(fullBuild.StandardArtifact.Geometries),
		CompactMaximumScale: 240, StandardMaximumScale: 700,
	}
	catalogRows, aggregates, err := runCatalogScan(c, fullTable)
	if err != nil {
		return err
	}
	for _, row := range resolutionRows {
		if strings.HasPrefix(row.Outcome, "typed_") {
			aggregates.TypedFailures++
		}
	}
	for _, row := range weightRows {
		if strings.HasPrefix(row.Outcome, "typed_") {
			aggregates.TypedFailures++
		}
	}
	if err := applyParameterCapProof(&aggregates, catalogRows); err != nil {
		return err
	}
	stages["catalog_scan"] = time.Since(stageStart).Milliseconds()

	predicates := evaluatePredicates(identities, resolutionRows, weightRows, phaseScans, aggregates)
	verdict := "PASS"
	if predicates.InputHashOrCountDrift || predicates.ExpectedNumericDrift || predicates.TypedFailureMismatch ||
		predicates.PhaseAcceptedAtOrBelow || predicates.AggregateOrMaximumDrift || predicates.ParameterCandidateFitsAll {
		verdict = "CONTRADICTION"
	}
	semantic := semanticDiagnostic{
		SchemaVersion: 1, Kind: "run-010-t0a1-semantic-diagnostic",
		Identities: identities, Resolution: resolutionRows, Weighting: weightRows,
		Phases: phaseScans, Catalog: catalogRows, Aggregates: aggregates,
		Predicates: predicates, Verdict: verdict,
	}
	raw, err := json.Marshal(semantic)
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(semanticPath, raw, 0o644); err != nil {
		return err
	}
	sum := sha256.Sum256(raw)
	receipt := timingReceipt{
		SchemaVersion: 1, Kind: "run-010-t0a1-timing",
		SemanticSHA256: hex.EncodeToString(sum[:]),
		DurationMS:     time.Since(started).Milliseconds(), StagesMS: stages,
	}
	timingRaw, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	if err := os.WriteFile(timingPath, append(timingRaw, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Printf("verdict=%s semantic_sha256=%s rows=%d catalog=%d\n", verdict, receipt.SemanticSHA256, aggregates.ParameterRows, aggregates.CatalogRows)
	if verdict != "PASS" {
		return fmt.Errorf("RUN-010 T0A.1 contradiction predicates are true")
	}
	return nil
}

func loadDiagnosticIdentities(c *catalog.Corpus) (diagnosticIdentities, recipe, error) {
	root := sourceRoot()
	rawRecipe, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/v1.recipe.json"))
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	var r recipe
	decoder := json.NewDecoder(bytes.NewReader(rawRecipe))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&r); err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	if err := validateRecipe(r, c); err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	digestFile := func(relative string) (string, error) {
		raw, readErr := os.ReadFile(filepath.Join(root, relative))
		if readErr != nil {
			return "", readErr
		}
		sum := sha256.Sum256(raw)
		return hex.EncodeToString(sum[:]), nil
	}
	lockHash, err := digestFile("internal/geometry/lod/tool/package-lock.json")
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	toolHash, err := digestFile("internal/geometry/lod/tool/build.mjs")
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	goModHash, err := digestFile("go.mod")
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	goSumHash, err := digestFile("go.sum")
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	serializerHash, err := digestFile("internal/geometry/serialize.go")
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	serializerTestHash, err := digestFile("internal/geometry/serialize_test.go")
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	goSources, goSourcesHash, err := loadDiagnosticSourceInventory(root)
	if err != nil {
		return diagnosticIdentities{}, recipe{}, err
	}
	recipeSum := sha256.Sum256(rawRecipe)
	return diagnosticIdentities{
		P1Corpus: c.Manifest.Identity, BaselineCommit: diagnosticBaselineCommit,
		BaselineTree: diagnosticBaselineTree, ScratchBefore: diagnosticScratchBefore,
		Projection:   geometry.LODProjectionContractV1(),
		RecipeSHA256: hex.EncodeToString(recipeSum[:]), MapshaperVersion: r.Mapshaper.Version,
		MapshaperIntegrity: r.Mapshaper.Integrity, MapshaperShasum: r.Mapshaper.Shasum,
		LockfileSHA256: lockHash, ToolSHA256: toolHash,
		GoVersion: runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH,
		GoSources: goSources, GoSourcesSHA256: goSourcesHash,
		GoModSHA256: goModHash, GoSumSHA256: goSumHash,
		SerializerSHA256: serializerHash, SerializerTestSHA: serializerTestHash,
		Command: diagnosticCommandIdentity,
	}, r, nil
}

func loadDiagnosticSourceInventory(root string) ([]sourceIdentity, string, error) {
	items := make([]sourceIdentity, 0, len(diagnosticSourcePaths))
	for _, path := range diagnosticSourcePaths {
		raw, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return nil, "", err
		}
		sum := sha256.Sum256(raw)
		items = append(items, sourceIdentity{Path: path, SHA256: hex.EncodeToString(sum[:])})
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return nil, "", err
	}
	sum := sha256.Sum256(raw)
	return items, hex.EncodeToString(sum[:]), nil
}

func buildDiagnosticCandidate(c *catalog.Corpus, out, geometryID string, resolution int, weighting float64, r recipe) (geometry.ProjectedLODGeometry, error) {
	if err := os.MkdirAll(out, 0o755); err != nil {
		return geometry.ProjectedLODGeometry{}, err
	}
	var source catalog.Geometry
	for _, item := range c.Geometries {
		if item.ID == geometryID {
			source = item
			break
		}
	}
	if source.ID == "" {
		return geometry.ProjectedLODGeometry{}, fmt.Errorf("geometry %q missing", geometryID)
	}
	record, err := geometry.ProjectLODGeometry(source, r.Projection.Flatness, r.Projection.CoordinatePrecision)
	if err != nil {
		return geometry.ProjectedLODGeometry{}, err
	}
	rawRecipe, err := os.ReadFile(filepath.Join(sourceRoot(), "internal/geometry/lod/v1.recipe.json"))
	if err != nil {
		return geometry.ProjectedLODGeometry{}, err
	}
	recipeHash := sha256.Sum256(rawRecipe)
	record.SourceCorpus = c.Manifest.Identity
	record.RecipeSHA256 = hex.EncodeToString(recipeHash[:])
	input := filepath.Join(out, "source.geojson")
	if err := writeSource(input, []catalog.Geometry{record.Geometry}); err != nil {
		return geometry.ProjectedLODGeometry{}, err
	}
	output := filepath.Join(out, "output.geojson")
	cmd := diagnosticNodeCommand(input, output, resolution, weighting)
	if combined, runErr := cmd.CombinedOutput(); runErr != nil {
		return geometry.ProjectedLODGeometry{}, fmt.Errorf("mapshaper diagnostic: %w: %s", runErr, combined)
	}
	items, err := readOutput(output)
	if err != nil {
		return geometry.ProjectedLODGeometry{}, err
	}
	if len(items) != 1 || items[0].ID != geometryID {
		return geometry.ProjectedLODGeometry{}, fmt.Errorf("mapshaper diagnostic coverage: %+v", items)
	}
	record.Geometry = items[0]
	return record, nil
}

func diagnosticNodeCommand(input, output string, resolution int, weighting float64) *exec.Cmd {
	cmd := exec.Command(
		"node", filepath.Join(sourceRoot(), "internal/geometry/lod/tool/build.mjs"),
		input, output, strconv.Itoa(resolution), strconv.FormatFloat(weighting, 'f', -1, 64),
	)
	cmd.Dir = sourceRoot()
	return cmd
}

func diagnosticSelectionRow(c *catalog.Corpus, r recipe, compact, standard geometry.ProjectedLODGeometry, value float64) parameterRow {
	table := diagnosticTable(r, compact, standard)
	in, err := geometry.InputFromCatalog(c, "AQ", "un", "card")
	if err != nil {
		return parameterRow{Value: value, Outcome: "setup_error", Failure: err.Error()}
	}
	resolved, err := geometry.ApplyPreset(in)
	if err != nil {
		return parameterRow{Value: value, Outcome: "setup_error", Failure: err.Error()}
	}
	resolved.Preset, resolved.MaxPathBytes = "", 0
	provenance, err := geometry.SelectLODProvenance(resolved, table)
	if err != nil {
		return parameterRow{Value: value, Outcome: "selection_error", Failure: err.Error()}
	}
	row := parameterRow{Value: value, SelectedTier: provenance.SelectedTier}
	for _, fallback := range provenance.Fallbacks {
		if !strings.HasPrefix(fallback, "compact:") {
			continue
		}
		switch {
		case strings.HasPrefix(fallback, "compact:raw_deviation:"):
			raw, parseErr := strconv.ParseFloat(strings.TrimPrefix(fallback, "compact:raw_deviation:"), 64)
			if parseErr != nil {
				row.Outcome, row.Failure = "parse_error", fallback
				return row
			}
			row.Outcome, row.RawDeviation = "raw_deviation", &raw
		case strings.HasPrefix(fallback, "compact:topology:"):
			row.Outcome, row.Failure = "typed_topology_failure", strings.TrimPrefix(fallback, "compact:topology:")
		default:
			row.Outcome, row.Failure = "typed_other_failure", fallback
		}
		return row
	}
	row.Outcome = "compact_accepted"
	raw := provenance.RawDeviation
	row.RawDeviation = &raw
	return row
}

func diagnosticTable(r recipe, compact, standard geometry.ProjectedLODGeometry) *geometry.LODTable {
	return &geometry.LODTable{
		Version: "v1", RecipeSHA256: compact.RecipeSHA256,
		Compact:              map[string]geometry.ProjectedLODGeometry{compact.Geometry.ID: compact},
		Standard:             map[string]geometry.ProjectedLODGeometry{standard.Geometry.ID: standard},
		CompactMaximumScale:  r.Tiers.Compact.MaximumFittedScale,
		StandardMaximumScale: r.Tiers.Standard.MaximumFittedScale,
	}
}

func runPhaseScan(c *catalog.Corpus, table *geometry.LODTable, name string, phases []geometry.GridPhase) phaseScan {
	in, err := geometry.InputFromCatalog(c, "AQ", "un", "card")
	if err != nil {
		return phaseScan{Name: name, MinimumFinite: math.Inf(1)}
	}
	resolved, err := geometry.ApplyPreset(in)
	if err != nil {
		return phaseScan{Name: name, MinimumFinite: math.Inf(1)}
	}
	resolved.Preset, resolved.MaxPathBytes = "", 0
	rows, err := geometry.DiagnosticGridPhases(resolved, table, "compact", phases)
	if err != nil {
		return phaseScan{Name: name, Rows: []geometry.DiagnosticPhaseRow{{Failure: err.Error()}}, MinimumFinite: math.Inf(1)}
	}
	scan := phaseScan{Name: name, Rows: rows, MinimumFinite: math.Inf(1)}
	for _, row := range rows {
		if row.Accepted {
			scan.AcceptedCount++
		}
		if row.FinalDeviation != nil && finiteDiagnostic(*row.FinalDeviation) && *row.FinalDeviation < scan.MinimumFinite {
			scan.MinimumFinite = *row.FinalDeviation
		}
	}
	return scan
}

func coarsePhases() []geometry.GridPhase {
	out := make([]geometry.GridPhase, 0, 100)
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			out = append(out, geometry.GridPhase{X: float64(x) / 1000, Y: float64(y) / 1000})
		}
	}
	return out
}

func localPhases(x0, x1, y0, y1, divisor int) []geometry.GridPhase {
	out := make([]geometry.GridPhase, 0, (x1-x0+1)*(y1-y0+1))
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			out = append(out, geometry.GridPhase{X: float64(x) / float64(divisor), Y: float64(y) / float64(divisor)})
		}
	}
	return out
}

func runCatalogScan(c *catalog.Corpus, table *geometry.LODTable) ([]catalogRow, diagnosticAggregates, error) {
	rows := make([]catalogRow, 0, c.Manifest.EntityCount*4)
	aggregates := diagnosticAggregates{
		ParameterRows: 15 + 7 + 100 + 25 + 121,
		CatalogRows:   c.Manifest.EntityCount * 4,
		Selected:      map[string]int{}, HistoricalOverruns: map[string]int{},
	}
	for _, entity := range c.Manifest.Entities {
		for _, profile := range []string{"un", "de_facto"} {
			for _, preset := range []string{"card", "hero"} {
				in, err := geometry.InputFromCatalog(c, entity.Alpha2, profile, preset)
				if err != nil {
					return nil, aggregates, err
				}
				in, err = geometry.ApplyPreset(in)
				if err != nil {
					return nil, aggregates, err
				}
				in.Preset, in.MaxPathBytes = "", 0
				result, err := geometry.GenerateWithLOD(in, table)
				if err != nil {
					return nil, aggregates, fmt.Errorf("%s/%s/%s: %w", entity.Alpha2, profile, preset, err)
				}
				historical, p2Cap := 2500, 2200
				if preset == "hero" {
					historical, p2Cap = 8000, 7500
				}
				pathSum := sha256.Sum256([]byte(result.Path))
				row := catalogRow{
					Entity: entity.Alpha2, Profile: profile, Preset: preset, GeometryID: in.Geometry.ID,
					EffectiveScale: result.EffectiveScale, RequestedTier: result.LOD.RequestedTier,
					SelectedTier: result.LOD.SelectedTier, Topology: true, Protection: true,
					Omissions: result.Removals, RawDeviation: result.LOD.RawDeviation,
					FinalDeviation: result.LOD.FinalDeviation, Points: result.Metrics.OutputPoints,
					PathBytes: result.Metrics.PathBytes, OverHistoricalMaximum: result.Metrics.PathBytes > historical,
					OverP2Cap: result.Metrics.PathBytes > p2Cap, Fallbacks: result.LOD.Fallbacks,
					OutputSHA256: hex.EncodeToString(pathSum[:]),
				}
				rows = append(rows, row)
				key := preset + "/" + result.LOD.SelectedTier
				aggregates.Selected[key]++
				if row.OverHistoricalMaximum {
					aggregates.HistoricalOverruns[key]++
				}
				if preset == "card" && result.Metrics.PathBytes > aggregates.MaximumCard.PathBytes {
					aggregates.MaximumCard = aggregateMaximum{entity.Alpha2, profile, result.LOD.SelectedTier, result.Metrics.PathBytes}
				}
				if preset == "hero" && result.Metrics.PathBytes > aggregates.MaximumHero.PathBytes {
					aggregates.MaximumHero = aggregateMaximum{entity.Alpha2, profile, result.LOD.SelectedTier, result.Metrics.PathBytes}
				}
			}
		}
	}
	return rows, aggregates, nil
}

func evaluatePredicates(ids diagnosticIdentities, resolutions, weightings []parameterRow, phases []phaseScan, aggregates diagnosticAggregates) diagnosticPredicates {
	p := diagnosticPredicates{}
	p.InputHashOrCountDrift =
		diagnosticIdentityDrift(ids) ||
			len(resolutions) != 15 || len(weightings) != 7 || len(phases) != 3 ||
			len(phases[0].Rows) != 100 || len(phases[1].Rows) != 25 || len(phases[2].Rows) != 121 ||
			aggregates.CatalogRows != 996

	expectedResolution := map[int]float64{
		64: 1.537690736882369, 65: 1.537690736882369, 66: 1.537690736882369,
		68: 1.537690736882369, 72: 1.537690736882369, 80: 1.537690736882369,
		96: 1.537690736882369, 192: 1.5403736770216907, 224: 1.5403736770216907,
		256: 2.441111469697744, 320: 2.441111469697744, 512: 2.441111469697744,
		1024: 2.441111469697744,
	}
	for _, row := range resolutions {
		if _, typed := expectedTypedFailure(row.Value); typed {
			continue
		}
		expected, ok := expectedResolution[int(row.Value)]
		if !ok || row.RawDeviation == nil || math.Abs(*row.RawDeviation-expected) > 1e-12 {
			p.ExpectedNumericDrift = true
		}
	}
	p.TypedFailureMismatch = typedFailureMismatch(resolutions, aggregates.TypedFailures)
	expectedWeights := []float64{1.559170376003231, 1.559170376003231, 1.772910882577505, 1.5927713853074141, 1.537690736882369, 1.6726024738845322, 1.7575971919760147}
	for i, row := range weightings {
		if row.RawDeviation == nil || math.Abs(*row.RawDeviation-expectedWeights[i]) > 1e-12 {
			p.ExpectedNumericDrift = true
		}
	}
	if len(phases) == 3 {
		if math.Abs(phases[0].MinimumFinite-1.5360371611532861) > 1e-12 ||
			phases[0].AcceptedCount != 0 || phases[1].AcceptedCount != 0 ||
			phases[2].AcceptedCount != 0 || phases[2].MinimumFinite < 1.536 || phases[2].MinimumFinite > 1.536001 {
			p.ExpectedNumericDrift = true
		}
		for _, scan := range phases {
			for _, row := range scan.Rows {
				if row.Accepted && row.FinalDeviation != nil && *row.FinalDeviation <= 1.536 {
					p.PhaseAcceptedAtOrBelow = true
				}
			}
		}
	}
	expectedSelected := map[string]int{"card/compact": 149, "card/standard": 237, "card/source": 112, "hero/standard": 380, "hero/source": 118}
	expectedOverruns := map[string]int{"card/compact": 56, "card/standard": 216, "card/source": 76, "hero/standard": 222, "hero/source": 55}
	p.AggregateOrMaximumDrift = !equalCounts(aggregates.Selected, expectedSelected) ||
		!equalCounts(aggregates.HistoricalOverruns, expectedOverruns) ||
		aggregates.MaximumCard != (aggregateMaximum{Entity: "RU", Profile: "un", Tier: "source", PathBytes: 477784}) ||
		aggregates.MaximumHero != (aggregateMaximum{Entity: "CA", Profile: "un", Tier: "source", PathBytes: 675606})
	p.ParameterCandidateFitsAll = aggregates.AnyParameterFitsAll
	return p
}

func diagnosticIdentityDrift(ids diagnosticIdentities) bool {
	return ids.P1Corpus != "sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95" ||
		ids.BaselineCommit != diagnosticBaselineCommit ||
		ids.BaselineTree != diagnosticBaselineTree ||
		ids.ScratchBefore != diagnosticScratchBefore ||
		ids.Projection != geometry.LODProjectionContractV1() ||
		ids.RecipeSHA256 != "bcfae021251b038a869b20e671b9ffc7550361aafd60a4a93d5f2a67829add89" ||
		ids.MapshaperVersion != "0.7.44" ||
		ids.MapshaperIntegrity != "sha512-3Cx+IABMXt1G28Y8J7oalW5P5VYyt1vHz5FO+KkV51EHnJioo9h9maO+u+4IyCPZ9Mh1hqchgomFmc68GtFwQQ==" ||
		ids.MapshaperShasum != "e08fc40d50347698ea07d2f82ead42a98b1c62f0" ||
		ids.LockfileSHA256 != "f767a2f3f0dd03ec2eb68dce69d735e21791fdc299e740c2d054234e56fb1832" ||
		ids.ToolSHA256 != "7b6f7c2638453a230de6a3a9a9363a8a538e8f51a64c7f4ce8d66242a058d1fe" ||
		ids.GoVersion != "go1.26.4 darwin/arm64" ||
		diagnosticSourceInventoryDrift(ids.GoSources, ids.GoSourcesSHA256) ||
		ids.GoModSHA256 != diagnosticGoModHash ||
		ids.GoSumSHA256 != diagnosticGoSumHash ||
		ids.SerializerSHA256 != diagnosticSerializerHash ||
		ids.SerializerTestSHA != diagnosticSerializerT ||
		ids.Command != diagnosticCommandIdentity
}

func diagnosticSourceInventoryDrift(items []sourceIdentity, aggregate string) bool {
	if len(items) != len(diagnosticSourcePaths) || len(items) == 0 {
		return true
	}
	for i, item := range items {
		if item.Path != diagnosticSourcePaths[i] || len(item.SHA256) != sha256.Size*2 {
			return true
		}
		if _, err := hex.DecodeString(item.SHA256); err != nil {
			return true
		}
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return true
	}
	sum := sha256.Sum256(raw)
	return aggregate != hex.EncodeToString(sum[:])
}

func expectedTypedFailure(value float64) (string, bool) {
	switch value {
	case 128:
		return "topology_damage:ring:ring 0/0 self-intersects segments 0/273 (47.950829842430686,78.65337835362203)-(39.67274196989299,82.01421781626114) and (38.72327783235212,85.04347389359545)-(39.67274196989299,82.01421781626114)", true
	case 160:
		return "topology_damage:ring:ring 0/0 self-intersects segments 0/359 (47.950829842430686,78.65337835362203)-(39.67274196989299,82.01421781626114) and (39.27946865243845,83.87024354762056)-(39.67274196989299,82.01421781626114)", true
	default:
		return "", false
	}
}

func typedFailureMismatch(rows []parameterRow, count int) bool {
	if count != 2 {
		return true
	}
	seen := 0
	for _, row := range rows {
		expected, typed := expectedTypedFailure(row.Value)
		if !typed {
			continue
		}
		seen++
		if row.Outcome != "typed_topology_failure" || row.Failure != expected {
			return true
		}
	}
	return seen != 2
}

func proveParameterCandidatesCannotFitAll(rows []catalogRow) (bool, parameterCapProof, error) {
	for _, row := range rows {
		if row.RequestedTier == "compact" || !row.OverP2Cap {
			continue
		}
		cap := 2200
		if row.Preset == "hero" {
			cap = 7500
		}
		return false, parameterCapProof{
			Method: "unchanged_noncompact_witness", CandidateScope: "compact_resolution_and_weighting",
			Entity: row.Entity, Profile: row.Profile, Preset: row.Preset,
			RequestedTier: row.RequestedTier, PathBytes: row.PathBytes, P2Cap: cap,
		}, nil
	}
	return false, parameterCapProof{}, fmt.Errorf("parameter-cap predicate is undecidable from the frozen compact-only matrix: no unchanged non-compact over-cap witness")
}

func applyParameterCapProof(aggregates *diagnosticAggregates, rows []catalogRow) error {
	fits, proof, err := proveParameterCandidatesCannotFitAll(rows)
	if err != nil {
		return err
	}
	aggregates.AnyParameterFitsAll = fits
	aggregates.ParameterCapProof = proof
	return nil
}

func equalCounts(got, want map[string]int) bool {
	keys := make([]string, 0, len(got)+len(want))
	for key := range got {
		keys = append(keys, key)
	}
	for key := range want {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if got[key] != want[key] {
			return false
		}
	}
	return true
}

func finiteDiagnostic(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func projectedGeometryIndex(items []geometry.ProjectedLODGeometry) map[string]geometry.ProjectedLODGeometry {
	out := make(map[string]geometry.ProjectedLODGeometry, len(items))
	for _, item := range items {
		out[item.Geometry.ID] = item
	}
	return out
}
