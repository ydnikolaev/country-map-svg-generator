package geometry

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

//go:embed lod/silhouette-oracle.v1.json
var silhouetteOracleJSON []byte

type SilhouetteBand struct {
	ID                    string  `json:"id"`
	MaximumEffectiveScale float64 `json:"maximum_effective_scale"`
	GridLongSide          int     `json:"grid_long_side"`
	MinimumIoU            float64 `json:"minimum_iou"`
	MinimumRecall         float64 `json:"minimum_recall"`
	PathCap               int     `json:"path_cap"`
	CompleteFileMaximum   int     `json:"complete_file_maximum"`
	ContributionThreshold int     `json:"contribution_threshold"`
}

type SilhouetteOracle struct {
	SchemaVersion     int              `json:"schema_version"`
	Version           string           `json:"version"`
	P1Corpus          string           `json:"p1_corpus"`
	ProjectionVersion string           `json:"projection_version"`
	RecipeVersion     string           `json:"recipe_version"`
	RasterizerPath    string           `json:"rasterizer_path"`
	RasterizerSHA256  string           `json:"rasterizer_sha256"`
	Sampling          string           `json:"sampling"`
	FillRule          string           `json:"fill_rule"`
	Antialias         bool             `json:"antialias"`
	Stroke            bool             `json:"stroke"`
	Color             bool             `json:"color"`
	Softening         bool             `json:"softening"`
	Fit               string           `json:"fit"`
	GridDerivation    string           `json:"grid_derivation"`
	CandidateOrder    string           `json:"candidate_order"`
	TieBreak          string           `json:"tie_break"`
	VisibilityPolicy  string           `json:"visibility_policy"`
	Bands             []SilhouetteBand `json:"bands"`
}

type SilhouetteMetrics struct {
	Intersection, Union, ReferencePixels, CandidatePixels int
	IoU, Recall                                           float64
}

type SilhouetteRaster struct {
	Width, Height int
	Pixels        []bool
}

type SilhouetteCandidateEvaluation struct {
	Canonical                           MultiPolygon
	Path                                string
	Metrics                             SilhouetteMetrics
	Phase                               GridPhase
	ViewBox                             Bounds
	Points, Parts, PathBytes, Omissions int
	Topology, Protection                bool
	DominantComponent                   bool
	Visibility                          []ProtectedVisibilityComponent
}

type ProtectedVisibilityComponent struct {
	SourceOrder, Contribution, IdentityRank int
	SourceDigest, IdentityReason            string
	Dominant, Anchor, Visible, Subscale     bool
	Retained                                bool
	OmissionReason, TieBreak                string
	CandidateContribution                   int
}

func ParseSilhouetteOracle(raw []byte) (SilhouetteOracle, error) {
	var oracle SilhouetteOracle
	if err := json.Unmarshal(raw, &oracle); err != nil {
		return oracle, err
	}
	if oracle.SchemaVersion != 1 || oracle.Version != "v1" ||
		oracle.P1Corpus == "" || oracle.ProjectionVersion != "v1" || oracle.RecipeVersion != "v2" ||
		oracle.RasterizerPath != "internal/geometry/silhouette.go" || len(oracle.RasterizerSHA256) != 64 ||
		oracle.Sampling != "pixel_center" || oracle.FillRule != "nonzero" ||
		oracle.Antialias || oracle.Stroke || oracle.Color || oracle.Softening ||
		oracle.Fit != "natural_uniform" || oracle.GridDerivation != "band_long_side_from_reference_aspect" ||
		oracle.CandidateOrder != "fine_to_coarse" || oracle.TieBreak != "candidate_digest_ascending" ||
		oracle.VisibilityPolicy != "protected-visibility/v1" || len(oracle.Bands) != 2 {
		return oracle, fmt.Errorf("invalid silhouette-oracle/v1 contract")
	}
	want := []SilhouetteBand{
		{ID: "compact", MaximumEffectiveScale: 240, GridLongSide: 240, PathCap: 2200, CompleteFileMaximum: 2500},
		{ID: "standard", MaximumEffectiveScale: 700, GridLongSide: 700, PathCap: 7500, CompleteFileMaximum: 8000},
	}
	for i, band := range oracle.Bands {
		if band.ID != want[i].ID || band.MaximumEffectiveScale != want[i].MaximumEffectiveScale ||
			band.GridLongSide != want[i].GridLongSide || band.PathCap != want[i].PathCap ||
			band.CompleteFileMaximum != want[i].CompleteFileMaximum ||
			band.ContributionThreshold < 1 ||
			band.MinimumIoU <= 0 || band.MinimumIoU > 1 || band.MinimumRecall <= 0 || band.MinimumRecall > 1 {
			return oracle, fmt.Errorf("invalid %s silhouette band", band.ID)
		}
	}
	return oracle, nil
}

func EmbeddedSilhouetteOracle() (SilhouetteOracle, error) {
	return ParseSilhouetteOracle(silhouetteOracleJSON)
}

func RasterizeSilhouette(g MultiPolygon, referenceBounds Bounds, longSide int) (SilhouetteRaster, error) {
	if longSide <= 0 || referenceBounds.Width() <= 0 || referenceBounds.Height() <= 0 {
		return SilhouetteRaster{}, fmt.Errorf("invalid silhouette grid")
	}
	width, height := longSide, int(math.Round(float64(longSide)*referenceBounds.Height()/referenceBounds.Width()))
	if referenceBounds.Height() > referenceBounds.Width() {
		height, width = longSide, int(math.Round(float64(longSide)*referenceBounds.Width()/referenceBounds.Height()))
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	type event struct {
		x     float64
		delta int
	}
	pixels := make([]bool, width*height)
	for py := 0; py < height; py++ {
		y := referenceBounds.MinY + (float64(py)+.5)*referenceBounds.Height()/float64(height)
		var events []event
		for _, polygon := range g {
			for _, ring := range polygon {
				for i := 0; i+1 < len(ring); i++ {
					a, b := ring[i], ring[i+1]
					if a.Y == b.Y || !((a.Y <= y && y < b.Y) || (b.Y <= y && y < a.Y)) {
						continue
					}
					x := a.X + (y-a.Y)*(b.X-a.X)/(b.Y-a.Y)
					delta := -1
					if b.Y > a.Y {
						delta = 1
					}
					events = append(events, event{x, delta})
				}
			}
		}
		sort.Slice(events, func(i, j int) bool {
			if events[i].x == events[j].x {
				return events[i].delta < events[j].delta
			}
			return events[i].x < events[j].x
		})
		winding := 0
		e := 0
		for px := 0; px < width; px++ {
			x := referenceBounds.MinX + (float64(px)+.5)*referenceBounds.Width()/float64(width)
			for e < len(events) && events[e].x <= x {
				winding += events[e].delta
				e++
			}
			pixels[py*width+px] = winding != 0
		}
	}
	return SilhouetteRaster{Width: width, Height: height, Pixels: pixels}, nil
}

func CompareSilhouettes(reference, candidate MultiPolygon, longSide int) (SilhouetteMetrics, error) {
	bounds := geometryBounds(reference)
	ref, err := RasterizeSilhouette(reference, bounds, longSide)
	if err != nil {
		return SilhouetteMetrics{}, err
	}
	got, err := RasterizeSilhouette(candidate, bounds, longSide)
	if err != nil {
		return SilhouetteMetrics{}, err
	}
	if ref.Width != got.Width || ref.Height != got.Height {
		return SilhouetteMetrics{}, fmt.Errorf("silhouette grid mismatch")
	}
	var metrics SilhouetteMetrics
	for i := range ref.Pixels {
		if ref.Pixels[i] {
			metrics.ReferencePixels++
		}
		if got.Pixels[i] {
			metrics.CandidatePixels++
		}
		if ref.Pixels[i] && got.Pixels[i] {
			metrics.Intersection++
		}
		if ref.Pixels[i] || got.Pixels[i] {
			metrics.Union++
		}
	}
	if metrics.Union > 0 {
		metrics.IoU = float64(metrics.Intersection) / float64(metrics.Union)
	}
	if metrics.ReferencePixels > 0 {
		metrics.Recall = float64(metrics.Intersection) / float64(metrics.ReferencePixels)
	}
	return metrics, nil
}

func SilhouetteBandPasses(band SilhouetteBand, metrics SilhouetteMetrics) bool {
	return metrics.IoU >= band.MinimumIoU && metrics.Recall >= band.MinimumRecall
}

// EvaluateSilhouetteCandidate is a maintainer-only T0A.2 primitive. It shares
// projection, fixed-grid canonicalization, topology and protected-anchor
// validation with production, but automatic acceptance is exclusively oracle
// and budget driven; DEC-005 raw deviation is intentionally not consulted.
func EvaluateSilhouetteCandidate(raw Input, record ProjectedLODGeometry, recipeSHA string, band SilhouetteBand) (SilhouetteCandidateEvaluation, error) {
	in, err := ApplyPreset(raw)
	if err != nil {
		return SilhouetteCandidateEvaluation{}, err
	}
	full, err := normalize(fromCatalog(in.Geometry.Coordinates))
	if err != nil {
		return SilhouetteCandidateEvaluation{}, err
	}
	contract := LODProjectionContractV1()
	prj := centeredProjector(full, 0)
	projected, _, err := projectGeometry(full, prj, contract.Flatness)
	if err != nil {
		return SilhouetteCandidateEvaluation{}, err
	}
	minimumParts := 1
	for _, feature := range in.Entity.Protected {
		if feature.MinimumParts > minimumParts {
			minimumParts = feature.MinimumParts
		}
	}
	table := &LODTable{RecipeSHA256: recipeSHA}
	projectedCandidate, err := bindProjectedLOD(record, in, table, prj, 0)
	if err != nil {
		return SilhouetteCandidateEvaluation{}, err
	}
	// The card is fitted to the candidate — the geometry that will be drawn —
	// not to the whole territorial claim. Fidelity is already judged against the
	// visible reference, so fitting the layout to the full projection framed the
	// card for components the oracle had already decided are not drawn. France
	// is the extreme: 11 components spanning 119 degrees of longitude, of which
	// the card draws metropolitan France and Corsica, leaving the silhouette at
	// 11.9% of the frame width.
	//
	// Fitting to the *visible* reference instead is wrong, and the failure is
	// worth recording: a component can be present in the candidate while sitting
	// below the band's visibility threshold, so a frame sized to the visible set
	// leaves it outside the viewBox, containment fails, and the search walks to a
	// coarser rung until the candidate loses it. Measured on France, that traded
	// Corsica and 936 bytes of detail for a single 35-point blob.
	_, viewBox, transform, effective, _, err := fitGeometry(projectedCandidate, in.Layout)
	if err != nil {
		return SilhouetteCandidateEvaluation{}, err
	}
	transform.CenterLon, transform.CenterLat = deg(prj.lon0), deg(prj.lat0)
	quality := AutoQuality(effective)
	// The phase reference stays the full projection at the new transform: it
	// supplies shared vertices and the omission count, both of which are facts
	// about the source rather than about the frame.
	fitted := transformGeometry(projected, transform.Scale, transform.TranslateX, transform.TranslateY)
	visibility, err := buildProtectedVisibility(in, projected, projectedCandidate, prj, band)
	if err != nil {
		return SilhouetteCandidateEvaluation{}, err
	}
	visibleReference := make(MultiPolygon, 0, len(projected))
	for i, component := range visibility {
		if component.Visible {
			visibleReference = append(visibleReference, projected[i])
		}
	}
	metrics, err := CompareSilhouettes(visibleReference, projectedCandidate, band.GridLongSide)
	if err != nil {
		return SilhouetteCandidateEvaluation{}, err
	}
	dominant := true
	fittedCandidate := transformGeometry(projectedCandidate, transform.Scale, transform.TranslateX, transform.TranslateY)
	// Topology is judged on the canonical geometry inside the phase loop below,
	// not here. A pre-canonical check runs at a precision the output never has:
	// simplification can leave a self-crossing of ~1e-4 in fitted units, four
	// orders of magnitude below the q=0.01 grid every emitted path is snapped
	// to, and canonicalization removes it. Rejecting the candidate at this point
	// discards a representation that would have been valid.
	//
	// Measured on RU/un, which is what surfaced this: rungs 512 down to 32 were
	// all rejected here by the same self-crossing, leaving Russia on rung 24 as a
	// single 45-point part using 508 of its 7500 hero bytes. Replayed through
	// canonicalization those same rungs are clean, on-grid and contained, at 39,
	// 9 and 3 parts. The source geometry is valid at 209 and 214 parts, so the
	// defect is introduced by simplification and erased by canonicalization —
	// this check only ever saw it in between.
	omissions := len(fitted) - len(fittedCandidate)
	if omissions < 0 {
		omissions = 0
	}
	var lastErr error
	for _, phase := range gridPhaseSchedule() {
		phaseTransform := composeGridPhase(transform, phase)
		candidate := materializeGridPhase(unfitGeometry(fittedCandidate, transform), phaseTransform)
		reference := materializeGridPhase(unfitGeometry(fitted, transform), phaseTransform)
		canonical, canonicalErr := canonicalizeShared(candidate, quality.Quantization, sharedSourceVertices(reference, quality.Quantization))
		if canonicalErr != nil {
			lastErr = canonicalErr
			continue
		}
		if !geometryOnGrid(canonical, quality.Quantization) || !geometryContainedBy(canonical, shiftViewBox(viewBox, phase)) {
			lastErr = fail(ErrTopology, in.Entity.Alpha2, "silhouette_grid", "candidate is off-grid or outside viewBox")
			continue
		}
		if topologyErr := validateTopology(canonical); topologyErr != nil {
			lastErr = topologyErr
			continue
		}
		_ = minimumParts
		_, path, serializationErr := commandsAndPath(canonical, 0)
		if serializationErr != nil {
			lastErr = serializationErr
			continue
		}
		return SilhouetteCandidateEvaluation{
			Canonical: canonical, Path: path, Metrics: metrics, Phase: phase,
			ViewBox: shiftViewBox(viewBox, phase),
			Points:  countPoints(canonical), Parts: len(canonical), PathBytes: len(path), Omissions: omissions,
			Topology: true, Protection: true, DominantComponent: dominant, Visibility: visibility,
		}, nil
	}
	if lastErr == nil {
		lastErr = fail(ErrTopology, in.Entity.Alpha2, "silhouette_phase", "no phase accepted")
	}
	return SilhouetteCandidateEvaluation{}, lastErr
}

func buildProtectedVisibility(in Input, reference, candidate MultiPolygon, prj projector, band SilhouetteBand) ([]ProtectedVisibilityComponent, error) {
	if len(reference) == 0 {
		return nil, fail(ErrProtected, in.Entity.Alpha2, "visibility", "missing reference components")
	}
	bounds := geometryBounds(reference)
	items := make([]ProtectedVisibilityComponent, len(reference))
	largest := 0
	for i, polygon := range reference {
		raw, _ := json.Marshal(polygon)
		sum := sha256.Sum256(raw)
		raster, err := RasterizeSilhouette(MultiPolygon{polygon}, bounds, band.GridLongSide)
		if err != nil {
			return nil, err
		}
		contribution := 0
		for _, filled := range raster.Pixels {
			if filled {
				contribution++
			}
		}
		items[i] = ProtectedVisibilityComponent{
			SourceOrder: i, SourceDigest: hex.EncodeToString(sum[:]), Contribution: contribution,
			TieBreak: fmt.Sprintf("%012d:%s", i, hex.EncodeToString(sum[:])),
		}
		if polygonArea(polygon) > polygonArea(reference[largest]) {
			largest = i
		}
	}
	items[largest].Dominant = true
	anchorPoints := map[int][]Point{}
	for _, feature := range in.Entity.Protected {
		anchor, err := prj.project(feature.Anchor[0], feature.Anchor[1])
		if err != nil {
			return nil, err
		}
		found := -1
		for i, polygon := range reference {
			if len(polygon) > 0 && pointInRing(anchor, polygon[0]) {
				found = i
				break
			}
		}
		if found < 0 {
			return nil, fail(ErrProtected, in.Entity.Alpha2, "anchor_lineage", "anchor %q has no source component", feature.Name)
		}
		items[found].Anchor = true
		anchorPoints[found] = append(anchorPoints[found], anchor)
	}
	ranked := make([]int, len(items))
	for i := range ranked {
		ranked[i] = i
	}
	sort.Slice(ranked, func(i, j int) bool {
		a, b := items[ranked[i]], items[ranked[j]]
		if a.Contribution != b.Contribution {
			return a.Contribution > b.Contribution
		}
		if a.SourceOrder != b.SourceOrder {
			return a.SourceOrder < b.SourceOrder
		}
		return a.SourceDigest < b.SourceDigest
	})
	minimum := 1
	for _, feature := range in.Entity.Protected {
		if feature.MinimumParts > minimum {
			minimum = feature.MinimumParts
		}
	}
	if minimum > len(items) {
		minimum = len(items)
	}
	base := map[int]bool{}
	for rank, index := range ranked[:minimum] {
		base[index] = true
		items[index].IdentityRank, items[index].IdentityReason = rank+1, "contribution_rank"
	}
	for index := range items {
		if !items[index].Anchor || base[index] {
			continue
		}
		evict := -1
		for rank := minimum - 1; rank >= 0; rank-- {
			candidateIndex := ranked[rank]
			if base[candidateIndex] && !items[candidateIndex].Anchor && !items[candidateIndex].Dominant {
				evict = candidateIndex
				break
			}
		}
		if evict < 0 {
			return nil, fail(ErrProtected, in.Entity.Alpha2, "anchor_rank", "no replaceable base member")
		}
		delete(base, evict)
		items[evict].IdentityRank, items[evict].IdentityReason = 0, ""
		base[index] = true
		items[index].IdentityRank, items[index].IdentityReason = minimum, "anchor_replacement"
	}
	for i := range items {
		items[i].Visible = items[i].Contribution >= band.ContributionThreshold || items[i].Dominant || items[i].Anchor
		items[i].Subscale = !items[i].Visible
		componentMetrics, err := CompareSilhouettes(MultiPolygon{reference[i]}, candidate, band.GridLongSide)
		if err != nil {
			return nil, err
		}
		items[i].CandidateContribution = componentMetrics.Intersection
		items[i].Retained = componentMetrics.Recall >= .5
		if items[i].Anchor {
			items[i].Retained = true
			for _, anchor := range anchorPoints[i] {
				covered := false
				for _, polygon := range candidate {
					if len(polygon) > 0 && pointInRing(anchor, polygon[0]) {
						covered = true
						break
					}
				}
				if !covered {
					items[i].Retained = false
					break
				}
			}
		}
		if !items[i].Retained {
			if items[i].Visible {
				return nil, fail(ErrProtected, in.Entity.Alpha2, "visible_component", "source component %d contribution=%d threshold=%d was lost", i, items[i].Contribution, band.ContributionThreshold)
			}
			if base[i] {
				items[i].OmissionReason = "protected_subscale"
			} else {
				items[i].OmissionReason = "subscale"
			}
		}
	}
	return items, nil
}

func dominantComponentRetained(reference, candidate MultiPolygon, longSide int, minimumRecall float64) (bool, error) {
	if len(reference) == 0 || len(candidate) == 0 {
		return false, nil
	}
	largest := 0
	for i := 1; i < len(reference); i++ {
		if polygonArea(reference[i]) > polygonArea(reference[largest]) {
			largest = i
		}
	}
	metrics, err := CompareSilhouettes(MultiPolygon{reference[largest]}, candidate, longSide)
	if err != nil {
		return false, err
	}
	return metrics.Recall >= minimumRecall, nil
}
