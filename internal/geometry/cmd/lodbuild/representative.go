package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/ydnikolaev/country-map-svg-generator/internal/geometry"
)

type representativeRecipe struct {
	SchemaVersion     int      `json:"schema_version"`
	Version           string   `json:"version"`
	SourceCorpus      string   `json:"source_corpus"`
	ProjectionVersion string   `json:"projection_version"`
	OracleVersion     string   `json:"oracle_version"`
	MapshaperVersion  string   `json:"mapshaper_version"`
	Algorithm         string   `json:"algorithm"`
	Weighting         float64  `json:"weighting"`
	CandidateOrder    string   `json:"candidate_order"`
	Resolutions       []int    `json:"resolutions"`
	VisibilityPolicy  string   `json:"visibility_policy"`
	Representatives   []string `json:"representatives"`
	Bands             []struct {
		ID                    string  `json:"id"`
		MaximumEffectiveScale float64 `json:"maximum_effective_scale"`
		PathCap               int     `json:"path_cap"`
		CompleteFileMaximum   int     `json:"complete_file_maximum"`
	} `json:"bands"`
}

type representativeRow struct {
	Entity, Profile, Band, CandidateID, OutputSHA256 string
	Resolution                                       int
	EffectiveScale, IoU, Recall                      float64
	Points, PathBytes, CompleteEstimate, Omissions   int
	Phase                                            geometry.GridPhase
	ContributionThreshold                            int
	RetainedIdentity, OmittedIdentity                int
	Visibility                                       []geometry.ProtectedVisibilityComponent
	GroupAnchors                                     []geometry.GroupAnchorResolution `json:",omitempty"`
	ViewBox                                          geometry.Bounds                  `json:"-"`
	Path                                             string                           `json:"-"`
}

type representativeManifest struct {
	SchemaVersion    int                 `json:"schema_version"`
	Kind             string              `json:"kind"`
	Corpus           string              `json:"corpus"`
	OracleSHA256     string              `json:"oracle_sha256"`
	RecipeSHA256     string              `json:"recipe_sha256"`
	CandidateOrder   string              `json:"candidate_order"`
	VisibilityPolicy string              `json:"visibility_policy"`
	Rows             []representativeRow `json:"rows"`
}

type representativeReceipt struct {
	SchemaVersion, Cases                                    int
	OracleSHA256, RecipeSHA256, ManifestSHA256, SheetSHA256 string
	Verdict                                                 string
	Thresholds                                              []geometry.SilhouetteBand
}

func buildIndonesiaBoundary(c *catalog.Corpus, output string) error {
	root := sourceRoot()
	recipeRaw, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/v2.recipe.json"))
	if err != nil {
		return err
	}
	recipeSum := sha256.Sum256(recipeRaw)
	baseRaw, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/v1.recipe.json"))
	if err != nil {
		return err
	}
	var base recipe
	if err := json.Unmarshal(baseRaw, &base); err != nil {
		return err
	}
	oracle, err := geometry.EmbeddedSilhouetteOracle()
	if err != nil {
		return err
	}
	in, err := geometry.InputFromCatalog(c, "ID", "un", "card")
	if err != nil {
		return err
	}
	in, err = geometry.ApplyPreset(in)
	if err != nil {
		return err
	}
	in.Preset, in.MaxPathBytes = "", 0
	type row struct {
		Resolution, Parts, PathBytes, CompleteEstimate int
	}
	rows := make([]row, 0, 2)
	temp, err := os.MkdirTemp("", "run012-id-boundary-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	for _, resolution := range []int{121, 120} {
		record, err := buildDiagnosticCandidate(c, filepath.Join(temp, strconv.Itoa(resolution)), in.Geometry.ID, resolution, .7, base)
		if err != nil {
			return err
		}
		record.RecipeSHA256 = hex.EncodeToString(recipeSum[:])
		evaluation, err := geometry.EvaluateSilhouetteCandidate(in, record, record.RecipeSHA256, oracle.Bands[0])
		if err != nil {
			return err
		}
		rows = append(rows, row{resolution, evaluation.Parts, evaluation.PathBytes, evaluation.PathBytes + 280})
	}
	if rows[0] != (row{121, 20, 2713, 2993}) || rows[1].Parts != 19 {
		return fmt.Errorf("Indonesia boundary drift: %+v", rows)
	}
	raw, err := json.Marshal(struct {
		SchemaVersion int
		Rows          []row
	}{1, rows})
	if err != nil {
		return err
	}
	return os.WriteFile(output, append(raw, '\n'), 0o644)
}

func buildRepresentative(c *catalog.Corpus, out string) error {
	root := sourceRoot()
	oracleRaw, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/silhouette-oracle.v1.json"))
	if err != nil {
		return err
	}
	oracle, err := geometry.ParseSilhouetteOracle(oracleRaw)
	if err != nil {
		return err
	}
	recipeRaw, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/v2.recipe.json"))
	if err != nil {
		return err
	}
	var v2 representativeRecipe
	if err := json.Unmarshal(recipeRaw, &v2); err != nil {
		return err
	}
	if v2.SchemaVersion != 1 || v2.Version != "v2" || v2.SourceCorpus != c.Manifest.Identity ||
		v2.ProjectionVersion != "v1" || v2.OracleVersion != "v1" ||
		v2.MapshaperVersion != "0.7.44" || v2.Algorithm != "weighted_visvalingam" ||
		v2.Weighting != .7 || v2.CandidateOrder != "fine_to_coarse" ||
		v2.VisibilityPolicy != oracle.VisibilityPolicy || len(v2.Bands) != 2 {
		return fmt.Errorf("invalid representative v2 recipe")
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	oracleSum, recipeSum := sha256.Sum256(oracleRaw), sha256.Sum256(recipeRaw)
	baseRecipeRaw, err := os.ReadFile(filepath.Join(root, "internal/geometry/lod/v1.recipe.json"))
	if err != nil {
		return err
	}
	var base recipe
	if err := json.Unmarshal(baseRecipeRaw, &base); err != nil {
		return err
	}
	groupAnchors, err := geometry.EmbeddedGroupAnchors()
	if err != nil {
		return err
	}
	var rows []representativeRow
	for _, entity := range v2.Representatives {
		for bandIndex, preset := range []string{"card", "hero"} {
			band := oracle.Bands[bandIndex]
			in, inputErr := geometry.InputFromCatalog(c, entity, "un", preset)
			if inputErr != nil {
				return inputErr
			}
			resolved, presetErr := geometry.ApplyPreset(in)
			if presetErr != nil {
				return presetErr
			}
			resolved.Preset, resolved.MaxPathBytes = "", 0
			resolved, groupApplication, groupErr := geometry.ApplyGroupAnchors(resolved, groupAnchors, band)
			if groupErr != nil {
				return groupErr
			}
			var accepted *representativeRow
			var rejected []string
			for _, resolution := range v2.Resolutions {
				dir := filepath.Join(out, "candidates", entity, band.ID, strconv.Itoa(resolution))
				record, buildErr := buildDiagnosticCandidate(c, dir, resolved.Geometry.ID, resolution, v2.Weighting, base)
				if buildErr != nil {
					return buildErr
				}
				record.RecipeSHA256 = hex.EncodeToString(recipeSum[:])
				evaluation, evaluateErr := geometry.EvaluateSilhouetteCandidate(resolved, record, record.RecipeSHA256, band)
				if evaluateErr != nil {
					rejected = append(rejected, fmt.Sprintf("%d:error:%v", resolution, evaluateErr))
					continue
				}
				if groupErr := geometry.ValidateGroupAnchorResult(groupApplication, evaluation.Visibility); groupErr != nil {
					rejected = append(rejected, fmt.Sprintf("%d:error:%v", resolution, groupErr))
					continue
				}
				complete := evaluation.PathBytes + 280
				if !geometry.SilhouetteBandPasses(band, evaluation.Metrics) ||
					evaluation.PathBytes > band.PathCap || complete >= band.CompleteFileMaximum ||
					!evaluation.Topology || !evaluation.Protection || !evaluation.DominantComponent {
					rejected = append(rejected, fmt.Sprintf("%d:iou=%.6f:recall=%.6f:bytes=%d:file=%d", resolution, evaluation.Metrics.IoU, evaluation.Metrics.Recall, evaluation.PathBytes, complete))
					continue
				}
				geometryRaw, _ := json.Marshal(record.Geometry)
				candidateSum := sha256.Sum256(geometryRaw)
				outputSum := sha256.Sum256([]byte(evaluation.Path))
				row := representativeRow{
					Entity: entity, Profile: "un", Band: band.ID,
					CandidateID: hex.EncodeToString(candidateSum[:]), Resolution: resolution,
					EffectiveScale: band.MaximumEffectiveScale, IoU: evaluation.Metrics.IoU,
					Recall: evaluation.Metrics.Recall, Points: evaluation.Points,
					PathBytes: evaluation.PathBytes, CompleteEstimate: complete,
					Omissions: evaluation.Omissions, Phase: evaluation.Phase,
					ContributionThreshold: band.ContributionThreshold,
					Visibility:            evaluation.Visibility,
					GroupAnchors:          groupApplication.Groups,
					ViewBox:               evaluation.ViewBox,
					OutputSHA256:          hex.EncodeToString(outputSum[:]), Path: evaluation.Path,
				}
				for _, component := range evaluation.Visibility {
					if component.IdentityRank > 0 && component.Retained {
						row.RetainedIdentity++
					}
					if component.IdentityRank > 0 && !component.Retained {
						row.OmittedIdentity++
					}
				}
				accepted = &row
				break
			}
			if accepted == nil {
				return fmt.Errorf("%s/%s: no global-threshold candidate fits oracle and budgets: %s", entity, band.ID, strings.Join(rejected, "; "))
			}
			if bandIndex == 1 {
				if err := validateMonotonicIdentity(rows[len(rows)-1], *accepted); err != nil {
					return fmt.Errorf("%s: %w", entity, err)
				}
			}
			rows = append(rows, *accepted)
		}
	}
	manifest := representativeManifest{
		SchemaVersion: 1, Kind: "silhouette-oracle-v1-representative",
		Corpus: c.Manifest.Identity, OracleSHA256: hex.EncodeToString(oracleSum[:]),
		RecipeSHA256: hex.EncodeToString(recipeSum[:]), CandidateOrder: v2.CandidateOrder,
		VisibilityPolicy: v2.VisibilityPolicy, Rows: rows,
	}
	manifestRaw, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	manifestRaw = append(manifestRaw, '\n')
	sheetRaw := renderRepresentativeSheet(rows)
	manifestSum, sheetSum := sha256.Sum256(manifestRaw), sha256.Sum256(sheetRaw)
	receipt := representativeReceipt{
		SchemaVersion: 1, Cases: len(rows), OracleSHA256: hex.EncodeToString(oracleSum[:]),
		RecipeSHA256: hex.EncodeToString(recipeSum[:]), ManifestSHA256: hex.EncodeToString(manifestSum[:]),
		SheetSHA256: hex.EncodeToString(sheetSum[:]), Verdict: "MACHINE_PASS", Thresholds: oracle.Bands,
	}
	receiptRaw, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	receiptRaw = append(receiptRaw, '\n')
	for name, raw := range map[string][]byte{
		"silhouette-oracle.v1.json":           oracleRaw,
		"v2.recipe.json":                      recipeRaw,
		"representative.manifest.json":        manifestRaw,
		"representative.contact-sheet.svg":    sheetRaw,
		"representative.machine-receipt.json": receiptRaw,
	} {
		if err := os.WriteFile(filepath.Join(out, name), raw, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func validateMonotonicIdentity(compact, standard representativeRow) error {
	standardByOrder := map[int]geometry.ProtectedVisibilityComponent{}
	for _, component := range standard.Visibility {
		standardByOrder[component.SourceOrder] = component
	}
	for _, component := range compact.Visibility {
		if component.IdentityRank == 0 || !component.Retained {
			continue
		}
		larger, ok := standardByOrder[component.SourceOrder]
		if !ok || !larger.Retained {
			return fmt.Errorf("non-monotonic identity source_order=%d", component.SourceOrder)
		}
	}
	return nil
}

func renderRepresentativeSheet(rows []representativeRow) []byte {
	const columns, cellW, cellH = 4, 300, 230
	height := ((len(rows) + columns - 1) / columns) * cellH
	var out strings.Builder
	fmt.Fprintf(&out, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d"><rect width="100%%" height="100%%" fill="#f4f0e8"/>`, columns*cellW, height, columns*cellW, height)
	for i, row := range rows {
		x, y := (i%columns)*cellW, (i/columns)*cellH
		const visualX, visualY, visualW, visualH = 18.0, 42.0, 264.0, 120.0
		scale := math.Min(visualW/row.ViewBox.Width(), visualH/row.ViewBox.Height())
		tx := visualX + (visualW-row.ViewBox.Width()*scale)/2 - row.ViewBox.MinX*scale
		ty := visualY + (visualH-row.ViewBox.Height()*scale)/2 - row.ViewBox.MinY*scale
		fmt.Fprintf(&out, `<g transform="translate(%d %d)"><rect x="8" y="8" width="284" height="214" rx="10" fill="white" stroke="#c8c1b4"/><text x="18" y="30" font-family="system-ui" font-size="14" font-weight="700">%s/un · %s</text><g transform="translate(%g %g) scale(%g)"><path d="%s" fill="#19324d"/></g><text x="18" y="178" font-family="ui-monospace,monospace" font-size="10">res=%d bytes=%d file=%d</text><text x="18" y="194" font-family="ui-monospace,monospace" font-size="10">IoU=%.4f recall=%.4f threshold=%d</text><text x="18" y="210" font-family="ui-monospace,monospace" font-size="10">identity kept=%d subscale=%d</text></g>`,
			x, y, html.EscapeString(row.Entity), html.EscapeString(row.Band), tx, ty, scale, html.EscapeString(row.Path),
			row.Resolution, row.PathBytes, row.CompleteEstimate, row.IoU, row.Recall, row.ContributionThreshold, row.RetainedIdentity, row.OmittedIdentity)
	}
	out.WriteString("</svg>\n")
	return []byte(out.String())
}
