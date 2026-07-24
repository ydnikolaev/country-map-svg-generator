// Ladder build mode implements AM-005/DEC-009: per geometry x band candidate
// selection over a resolution ladder judged by the frozen silhouette oracle,
// with an identity fallback tried only after every declared/intermediate
// resolution rung fails, and a typed no-artifact outcome when nothing on the
// ladder — including identity — can satisfy topology at the actual fitted
// scale. It generalizes the T0 feasibility sweep (t0sweep/t0identity) into a
// committed-artifact build.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

// ladderRecipe mirrors internal/geometry/lod/ladder.recipe.json. It binds the
// frozen v2 recipe, oracle and corpus by digest and declares the fine-to-coarse
// resolution list — the frozen 29 plus the DEC-009 intermediate rung(s) — plus
// the identity-fallback ordering. It changes no threshold, tolerance, weighting,
// byte cap or q=0.01 quantization; those remain owned by the oracle and v2.
type ladderRecipe struct {
	SchemaVersion      int       `json:"schema_version"`
	Version            string    `json:"version"`
	SourceCorpus       string    `json:"source_corpus"`
	ProjectionVersion  string    `json:"projection_version"`
	OracleVersion      string    `json:"oracle_version"`
	OracleSHA256       string    `json:"oracle_sha256"`
	BaseRecipeVersion  string    `json:"base_recipe_version"`
	BaseRecipeSHA256   string    `json:"base_recipe_sha256"`
	MapshaperVersion   string    `json:"mapshaper_version"`
	MapshaperIntegrity string    `json:"mapshaper_integrity"`
	MapshaperShasum    string    `json:"mapshaper_shasum"`
	Algorithm          string    `json:"algorithm"`
	Weighting          float64   `json:"weighting"`
	CandidateOrder     string    `json:"candidate_order"`
	Resolutions        []float64 `json:"resolutions"`
	IdentityFallback   string    `json:"identity_fallback"`
	VisibilityPolicy   string    `json:"visibility_policy"`
	Bands              []struct {
		ID                    string  `json:"id"`
		MaximumEffectiveScale float64 `json:"maximum_effective_scale"`
		PathCap               int     `json:"path_cap"`
		CompleteFileMaximum   int     `json:"complete_file_maximum"`
	} `json:"bands"`
	Notes string `json:"notes"`
}

// ladderAttempt records one rejected candidate for full inspectability
// (DEC-006 item 7). Selection is either a formatted resolution or "identity".
type ladderAttempt struct {
	Selection string `json:"selection"`
	Reason    string `json:"reason"`
}

// ladderRow is one geometry x profile x band outcome.
type ladderRow struct {
	Alpha2, Profile, Band string
	GeometryID            string
	Status                string // "pass" | "no_artifact"
	Selection             string `json:",omitempty"` // formatted resolution, or "identity"
	IoU, Recall           float64
	PathBytes             int                                     `json:",omitempty"`
	CompleteEstimate      int                                     `json:",omitempty"`
	Omissions             int                                     `json:",omitempty"`
	Points, Parts         int                                     `json:",omitempty"`
	Visibility            []geometry.ProtectedVisibilityComponent `json:",omitempty"`
	GroupAnchors          []geometry.GroupAnchorResolution        `json:",omitempty"`
	NoArtifactReason      string                                  `json:",omitempty"`
	Attempts              []ladderAttempt                         `json:",omitempty"`
}

// ladderArtifact is the committed T1 deliverable: per (alpha2, profile, band)
// selection provenance plus the deduplicated selected candidate geometries
// (keyed by geometry id, then by selection), so the same simplified or
// identity geometry is stored once even when both un/de_facto or both bands
// select it.
type ladderArtifact struct {
	SchemaVersion      int                                                 `json:"schema_version"`
	Kind               string                                              `json:"kind"`
	SourceCorpus       string                                              `json:"source_corpus"`
	OracleSHA256       string                                              `json:"oracle_sha256"`
	LadderRecipeSHA256 string                                              `json:"ladder_recipe_sha256"`
	BaseRecipeSHA256   string                                              `json:"base_recipe_sha256"`
	CandidateOrder     string                                              `json:"candidate_order"`
	VisibilityPolicy   string                                              `json:"visibility_policy"`
	Candidates         map[string]map[string]geometry.ProjectedLODGeometry `json:"candidates"`
	Rows               []ladderRow                                         `json:"rows"`
}

func loadLadderRecipe(root string, c *catalog.Corpus, oracle geometry.SilhouetteOracle) (ladderRecipe, string, error) {
	raw, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/ladder.recipe.json"))
	if err != nil {
		return ladderRecipe{}, "", err
	}
	var r ladderRecipe
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&r); err != nil {
		return ladderRecipe{}, "", fmt.Errorf("ladder recipe: %w", err)
	}
	oracleRaw, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/silhouette-oracle.v1.json"))
	if err != nil {
		return ladderRecipe{}, "", err
	}
	oracleSum := sha256.Sum256(oracleRaw)
	base, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/v2.recipe.json"))
	if err != nil {
		return ladderRecipe{}, "", err
	}
	baseSum := sha256.Sum256(base)
	if r.SchemaVersion != 1 || r.Version != "ladder-v1" || r.SourceCorpus != c.Manifest.Identity ||
		r.ProjectionVersion != "v1" || r.OracleVersion != "v1" || r.OracleSHA256 != hex.EncodeToString(oracleSum[:]) ||
		r.BaseRecipeVersion != "v2" || r.BaseRecipeSHA256 != hex.EncodeToString(baseSum[:]) ||
		r.MapshaperVersion != "0.7.44" || r.Algorithm != "weighted_visvalingam" || r.Weighting != .7 ||
		r.CandidateOrder != "fine_to_coarse" || r.IdentityFallback != "tried_last" ||
		r.VisibilityPolicy != oracle.VisibilityPolicy || len(r.Bands) != 2 || len(r.Resolutions) == 0 {
		return ladderRecipe{}, "", fmt.Errorf("ladder recipe does not bind the embedded corpus/oracle/v2-recipe contract")
	}
	for i := 1; i < len(r.Resolutions); i++ {
		if r.Resolutions[i] >= r.Resolutions[i-1] {
			return ladderRecipe{}, "", fmt.Errorf("ladder resolutions are not strictly fine-to-coarse at index %d", i)
		}
	}
	recipeSum := sha256.Sum256(raw)
	return r, hex.EncodeToString(recipeSum[:]), nil
}

func formatResolution(r float64) string { return strconv.FormatFloat(r, 'f', -1, 64) }

const ladderIdentitySelection = "identity"

func buildLadder(c *catalog.Corpus, scratch string) (ladderArtifact, error) {
	root := sourceRoot()
	oracle, err := geometry.EmbeddedSilhouetteOracle()
	if err != nil {
		return ladderArtifact{}, err
	}
	recipe, recipeSHA, err := loadLadderRecipe(root, c, oracle)
	if err != nil {
		return ladderArtifact{}, err
	}
	groupAnchors, err := geometry.EmbeddedGroupAnchors()
	if err != nil {
		return ladderArtifact{}, err
	}
	if err := os.MkdirAll(scratch, 0o755); err != nil {
		return ladderArtifact{}, err
	}

	contract := geometry.LODProjectionContractV1()
	projectedByID := make(map[string]geometry.ProjectedLODGeometry, len(c.Geometries))
	projectedCatalog := make([]catalog.Geometry, len(c.Geometries))
	for i, source := range c.Geometries {
		record, projectErr := geometry.ProjectLODGeometry(source, contract.Flatness, contract.CoordinatePrecision)
		if projectErr != nil {
			return ladderArtifact{}, fmt.Errorf("project %s: %w", source.ID, projectErr)
		}
		record.SourceCorpus = c.Manifest.Identity
		record.RecipeSHA256 = recipeSHA
		projectedByID[source.ID] = record
		projectedCatalog[i] = record.Geometry
	}
	sourcePath := filepath.Join(scratch, "source.geojson")
	if err := writeSource(sourcePath, projectedCatalog); err != nil {
		return ladderArtifact{}, fmt.Errorf("write projected source: %w", err)
	}

	// Stage 1: run Mapshaper once per declared/intermediate resolution over
	// the whole projected catalog, fine to coarse.
	ladder := make(map[string]map[string]catalog.Geometry, len(recipe.Resolutions))
	stageStart := time.Now()
	for i, resolution := range recipe.Resolutions {
		key := formatResolution(resolution)
		geoPath := filepath.Join(scratch, fmt.Sprintf("res-%s.geojson", key))
		script := filepath.Join(root, "internal/geometry/lod/tool/build.mjs")
		cmd := exec.Command("node", script, sourcePath, geoPath, key, formatResolution(recipe.Weighting))
		cmd.Dir = root
		resStart := time.Now()
		combined, runErr := cmd.CombinedOutput()
		if runErr != nil {
			return ladderArtifact{}, fmt.Errorf("mapshaper resolution=%s: %w: %s", key, runErr, combined)
		}
		outputGeometries, readErr := readOutput(geoPath)
		if readErr != nil {
			return ladderArtifact{}, fmt.Errorf("read mapshaper output resolution=%s: %w", key, readErr)
		}
		byID := make(map[string]catalog.Geometry, len(outputGeometries))
		for _, g := range outputGeometries {
			byID[g.ID] = g
		}
		ladder[key] = byID
		fmt.Printf("[ladder] mapshaper %2d/%d resolution=%-8s geometries=%d duration=%s\n",
			i+1, len(recipe.Resolutions), key, len(byID), time.Since(resStart).Round(time.Millisecond))
	}
	fmt.Printf("[ladder] mapshaper stage complete duration=%s\n", time.Since(stageStart).Round(time.Second))

	// Stage 2: pure-Go per-entity, per-profile, per-band fine-to-coarse
	// search, identity fallback tried last, over the frozen oracle primitive
	// exactly as written.
	entities := append([]catalog.Entity(nil), c.Manifest.Entities...)
	sort.Slice(entities, func(i, j int) bool { return entities[i].Alpha2 < entities[j].Alpha2 })

	bandPresets := []struct {
		Band   geometry.SilhouetteBand
		Preset string
	}{
		{oracle.Bands[0], "card"},
		{oracle.Bands[1], "hero"},
	}

	candidates := map[string]map[string]geometry.ProjectedLODGeometry{}
	rows := make([]ladderRow, 0, len(entities)*4)
	evalStart := time.Now()
	for ei, entity := range entities {
		for _, profile := range []string{"un", "de_facto"} {
			for _, bp := range bandPresets {
				row, candidate, evalErr := evaluateLadderCase(c, entity.Alpha2, profile, bp.Preset, bp.Band, recipe.Resolutions, ladder, projectedByID, groupAnchors, recipeSHA)
				if evalErr != nil {
					return ladderArtifact{}, fmt.Errorf("%s/%s/%s: %w", entity.Alpha2, profile, bp.Band.ID, evalErr)
				}
				if candidate != nil {
					byResolution, ok := candidates[row.GeometryID]
					if !ok {
						byResolution = map[string]geometry.ProjectedLODGeometry{}
						candidates[row.GeometryID] = byResolution
					}
					byResolution[row.Selection] = *candidate
				}
				rows = append(rows, row)
			}
		}
		if (ei+1)%25 == 0 || ei == len(entities)-1 {
			fmt.Printf("[ladder] evaluated %d/%d entities duration=%s\n", ei+1, len(entities), time.Since(evalStart).Round(time.Second))
		}
	}

	return ladderArtifact{
		SchemaVersion: 1, Kind: "silhouette-ladder/v1",
		SourceCorpus: c.Manifest.Identity, OracleSHA256: recipe.OracleSHA256,
		LadderRecipeSHA256: recipeSHA, BaseRecipeSHA256: recipe.BaseRecipeSHA256,
		CandidateOrder: recipe.CandidateOrder, VisibilityPolicy: recipe.VisibilityPolicy,
		Candidates: candidates, Rows: rows,
	}, nil
}

// evaluateLadderCase runs the fine-to-coarse resolution search for one
// (alpha2, profile, band) case, then the identity fallback last, exactly as
// DEC-009 orders it. It returns the winning row and its candidate geometry, or
// a no_artifact row carrying the typed reason and every rejected attempt.
func evaluateLadderCase(c *catalog.Corpus, alpha2, profile, preset string, band geometry.SilhouetteBand, resolutions []float64, ladder map[string]map[string]catalog.Geometry, projectedByID map[string]geometry.ProjectedLODGeometry, groupAnchors geometry.GroupAnchorSet, recipeSHA string) (ladderRow, *geometry.ProjectedLODGeometry, error) {
	in, err := geometry.InputFromCatalog(c, alpha2, profile, preset)
	if err != nil {
		return ladderRow{}, nil, err
	}
	resolved, err := geometry.ApplyPreset(in)
	if err != nil {
		return ladderRow{}, nil, err
	}
	resolved.Preset, resolved.MaxPathBytes = "", 0
	resolved, groupApplication, err := geometry.ApplyGroupAnchors(resolved, groupAnchors, band)
	if err != nil {
		return ladderRow{}, nil, err
	}

	var attempts []ladderAttempt
	tryCandidate := func(selection string, record geometry.ProjectedLODGeometry) (ladderRow, geometry.ProjectedLODGeometry, bool) {
		record.RecipeSHA256 = recipeSHA
		evaluation, evalErr := geometry.EvaluateSilhouetteCandidate(resolved, record, recipeSHA, band)
		if evalErr != nil {
			attempts = append(attempts, ladderAttempt{Selection: selection, Reason: evalErr.Error()})
			return ladderRow{}, geometry.ProjectedLODGeometry{}, false
		}
		if groupErr := geometry.ValidateGroupAnchorResult(groupApplication, evaluation.Visibility); groupErr != nil {
			attempts = append(attempts, ladderAttempt{Selection: selection, Reason: groupErr.Error()})
			return ladderRow{}, geometry.ProjectedLODGeometry{}, false
		}
		complete := evaluation.PathBytes + 280
		var reasons []string
		if !geometry.SilhouetteBandPasses(band, evaluation.Metrics) {
			reasons = append(reasons, fmt.Sprintf("iou=%.6f(min %.6f)", evaluation.Metrics.IoU, band.MinimumIoU))
			if evaluation.Metrics.Recall < band.MinimumRecall {
				reasons = append(reasons, fmt.Sprintf("recall=%.6f(min %.6f)", evaluation.Metrics.Recall, band.MinimumRecall))
			}
		}
		if evaluation.PathBytes > band.PathCap {
			reasons = append(reasons, fmt.Sprintf("bytes=%d(cap %d)", evaluation.PathBytes, band.PathCap))
		}
		if complete >= band.CompleteFileMaximum {
			reasons = append(reasons, fmt.Sprintf("file=%d(max %d)", complete, band.CompleteFileMaximum))
		}
		if !evaluation.Topology {
			reasons = append(reasons, "topology=false")
		}
		if !evaluation.Protection {
			reasons = append(reasons, "protection=false")
		}
		if !evaluation.DominantComponent {
			reasons = append(reasons, "dominant_component=false")
		}
		if len(reasons) > 0 {
			attempts = append(attempts, ladderAttempt{Selection: selection, Reason: strings.Join(reasons, ",")})
			return ladderRow{}, geometry.ProjectedLODGeometry{}, false
		}
		row := ladderRow{
			Alpha2: alpha2, Profile: profile, Band: band.ID, Status: "pass",
			GeometryID: resolved.Geometry.ID, Selection: selection,
			IoU: evaluation.Metrics.IoU, Recall: evaluation.Metrics.Recall,
			PathBytes: evaluation.PathBytes, CompleteEstimate: complete, Omissions: evaluation.Omissions,
			Points: evaluation.Points, Parts: evaluation.Parts,
			Visibility: evaluation.Visibility, GroupAnchors: groupApplication.Groups,
		}
		return row, record, true
	}

	for _, resolution := range resolutions {
		key := formatResolution(resolution)
		byID, ok := ladder[key]
		if !ok {
			attempts = append(attempts, ladderAttempt{Selection: key, Reason: "ladder_stage_missing"})
			continue
		}
		candidateGeometry, ok := byID[resolved.Geometry.ID]
		if !ok {
			attempts = append(attempts, ladderAttempt{Selection: key, Reason: "mapshaper_output_missing"})
			continue
		}
		base := projectedByID[resolved.Geometry.ID]
		base.Geometry = candidateGeometry
		if row, record, ok := tryCandidate(key, base); ok {
			return row, &record, nil
		}
	}
	identity := projectedByID[resolved.Geometry.ID]
	if row, record, ok := tryCandidate(ladderIdentitySelection, identity); ok {
		return row, &record, nil
	}

	last := attempts[len(attempts)-1]
	row := ladderRow{
		Alpha2: alpha2, Profile: profile, Band: band.ID, Status: "no_artifact",
		GeometryID: resolved.Geometry.ID, NoArtifactReason: last.Reason, Attempts: attempts,
	}
	return row, nil, nil
}

// marshalLadderArtifact is the single canonical encoding used both to write
// the committed artifact and, in tests, to prove a clean rebuild reproduces
// it byte-for-byte.
func marshalLadderArtifact(a ladderArtifact) ([]byte, error) {
	sort.Slice(a.Rows, func(i, j int) bool {
		if a.Rows[i].Alpha2 != a.Rows[j].Alpha2 {
			return a.Rows[i].Alpha2 < a.Rows[j].Alpha2
		}
		if a.Rows[i].Profile != a.Rows[j].Profile {
			return a.Rows[i].Profile < a.Rows[j].Profile
		}
		return a.Rows[i].Band < a.Rows[j].Band
	})
	// Compact, not indented: the committed artifact carries ~475 selected
	// candidate geometries across 283 source geometries; indentation roughly
	// doubles size for a machine-consumed file with no readability benefit.
	raw, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return append(raw, '\n'), nil
}

func writeLadderArtifact(path string, a ladderArtifact) error {
	raw, err := marshalLadderArtifact(a)
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

func runLadderBuild(scratch, artifactPath string) error {
	start := time.Now()
	c, err := catalog.Embedded()
	if err != nil {
		return err
	}
	artifact, err := buildLadder(c, scratch)
	if err != nil {
		return err
	}
	if err := writeLadderArtifact(artifactPath, artifact); err != nil {
		return err
	}
	pass, noArtifact, identity := 0, 0, 0
	for _, row := range artifact.Rows {
		switch row.Status {
		case "pass":
			pass++
			if row.Selection == ladderIdentitySelection {
				identity++
			}
		case "no_artifact":
			noArtifact++
		}
	}
	fmt.Printf("[ladder] rows=%d pass=%d identity=%d no_artifact=%d geometries=%d duration=%s\n",
		len(artifact.Rows), pass, identity, noArtifact, len(artifact.Candidates), time.Since(start).Round(time.Second))
	return nil
}
