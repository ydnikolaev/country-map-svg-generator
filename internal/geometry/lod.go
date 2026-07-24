package geometry

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

type LODTable struct {
	Version                                   string
	RecipeSHA256                              string
	Compact, Standard                         map[string]ProjectedLODGeometry
	CompactMaximumScale, StandardMaximumScale float64
	// Ladder is set on a published table and nil on a hand-built one. When it
	// is present, selection is a lookup against the committed DEC-006/DEC-009
	// ladder rather than a fine-to-coarse tier walk, and the typed no-artifact
	// outcome becomes reachable.
	Ladder *LadderTable
	// Band path caps, carried from the oracle the ladder was judged against.
	// The ladder guarantees the linear path fits these; the runtime has to hold
	// the guarantee for whatever representation it finally emits.
	CompactPathCap, StandardPathCap int
	mu                              sync.Mutex
	prepared                        map[string]preparedLOD
}

// bandPathCap is the frozen path budget for a band, or 0 when the table carries
// no ladder and therefore makes no such promise.
func (t *LODTable) bandPathCap(band string) int {
	switch band {
	case "compact":
		return t.CompactPathCap
	case "standard":
		return t.StandardPathCap
	}
	return 0
}

type preparedLOD struct {
	full, projected MultiPolygon
	prj             projector
	points          int
}
type LODProvenance struct {
	RequestedTier, SelectedTier string
	CoordinateSpace             string
	Finalization                string
	Fallbacks                   []string
	Restored                    []Removal
	RawDeviation                float64
	FinalDeviation              float64
	MaximumDeviation            float64
	ResolvedTolerance           float64
	Quantization                float64
	SourceAttemptedTolerance    float64
	SourceSelectedTolerance     float64
	AttemptedPhase              GridPhase
	SelectedPhase               GridPhase
	PhaseAttempts               int
	OutputPoints                int
	ParserRoundTrip             bool
	Projection                  LODProjection
}

// GenerateWithLOD runs the production full-source-first pipeline against an
// injected table. Published builds inject the generated table; an absent table
// uses the exact source tier without any external process.
func GenerateWithLOD(raw Input, table *LODTable) (Result, error) {
	return generateWithLOD(raw, table, false)
}

// SelectLODProvenance runs projection, fitting, retention, fixed-grid
// canonicalization, and candidate validation — including the frozen byte-cap
// check via exact serialization of every non-final candidate — then returns
// without commands or a rendered path for the selected candidate.
func SelectLODProvenance(raw Input, table *LODTable) (LODProvenance, error) {
	result, err := generateWithLOD(raw, table, true)
	return result.LOD, err
}

func generateWithLOD(raw Input, table *LODTable, selectionOnly bool) (Result, error) {
	in, err := ApplyPreset(raw)
	if err != nil {
		return Result{}, err
	}
	if !selectionOnly && in.MaxPathBytes < 0 {
		return Result{}, fail(ErrInvalidLayout, in.Entity.Alpha2, "max_path_bytes", "must be nonnegative")
	}
	rotation := 0.0
	minimumParts := 1
	offsets := map[string]Point{}
	if in.Override != nil {
		if err := ValidateOverride(*in.Override, in.CorpusID, in.Profile); err != nil {
			return Result{}, err
		}
		if in.Override.Rotation != nil {
			rotation = *in.Override.Rotation
		}
		if in.Override.Padding != nil {
			in.Layout.Padding = *in.Override.Padding
		}
		if in.Override.Quality != nil {
			in.Quality = *in.Override.Quality
		}
		if in.Override.MinimumParts > 0 {
			minimumParts = in.Override.MinimumParts
		}
		offsets = in.Override.MarkerOffsets
	}
	for _, feature := range in.Entity.Protected {
		if feature.MinimumParts > minimumParts {
			minimumParts = feature.MinimumParts
		}
	}
	full, err := normalize(fromCatalog(in.Geometry.Coordinates))
	if err != nil {
		return Result{}, err
	}
	projectionContract := LODProjectionContractV1()
	flat := projectionContract.Flatness
	projectionPolicy := projectionContract
	artifactProjectionCompatible := true
	if !in.Quality.Auto && in.Quality.Quantization > 0 && in.Quality.Flatness > 0 {
		flat = in.Quality.Flatness
		if flat != projectionContract.Flatness {
			projectionPolicy.Version = "explicit"
			projectionPolicy.Flatness = flat
			artifactProjectionCompatible = false
		}
	}
	prj := centeredProjector(full, rotation)
	var projected MultiPolygon
	var projectedPoints int
	if table != nil {
		key, hashErr := lodPreparationKey(in, rotation, flat)
		if hashErr != nil {
			return Result{}, hashErr
		}
		table.mu.Lock()
		if table.prepared == nil {
			table.prepared = map[string]preparedLOD{}
		}
		cached, ok := table.prepared[key]
		table.mu.Unlock()
		if ok {
			full, projected, prj, projectedPoints = cached.full, cached.projected, cached.prj, cached.points
		} else {
			projected, projectedPoints, err = projectGeometry(full, prj, flat)
			if err != nil {
				return Result{}, err
			}
			table.mu.Lock()
			table.prepared[key] = preparedLOD{full, projected, prj, projectedPoints}
			table.mu.Unlock()
		}
	} else {
		projected, projectedPoints, err = projectGeometry(full, prj, flat)
		if err != nil {
			return Result{}, err
		}
	}
	fitted, vb, tr, effective, diags, err := fitGeometry(projected, in.Layout)
	if err != nil {
		return Result{}, err
	}
	tr.CenterLon, tr.CenterLat = deg(prj.lon0), deg(prj.lat0)
	tr.Rotation = rotation
	quality := in.Quality
	if quality.Auto || quality.Quantization == 0 {
		quality = AutoQuality(effective)
	}
	if err := validateQuality(quality); err != nil {
		return Result{}, err
	}
	protected, err := protectedComponents(in.Entity, prj, tr, fitted)
	if err != nil {
		return Result{}, err
	}
	fullRequired, fullRemovals, fullComponents, err := retain(fitted, protected, minimumParts, quality.MinimumArea)
	if err != nil {
		return Result{}, err
	}
	requested := "source"
	if table != nil && artifactProjectionCompatible && effective <= table.StandardMaximumScale {
		requested = "standard"
		if effective <= table.CompactMaximumScale {
			requested = "compact"
		}
	}
	type candidate struct {
		name   string
		record *ProjectedLODGeometry
		ladder bool
	}
	var candidates []candidate
	switch {
	case table == nil:
		// No table at all: the exact source tier, unchanged.
	case table.Ladder != nil:
		// The committed ladder already picked the finest oracle-passing
		// candidate for this exact band, so exactly one derived candidate is
		// legal. Serving the standard rung inside a card viewport — which the
		// pre-DEC-006 tier walk did — is precisely the substitution DEC-006
		// removed.
		if requested != "source" && ladderVerdictApplies(in, rotation, quality, artifactProjectionCompatible) {
			selection, outcome := table.Ladder.Lookup(in.Geometry.ID, requested)
			switch outcome {
			case LadderNoArtifact:
				// DEC-009/DEC-010: a decision the build recorded, not a
				// failure to route around. Falling back to source here would
				// emit the oversized silhouette the oracle refused.
				return Result{}, fail(ErrNoArtifact, in.Entity.Alpha2, requested,
					"no ladder rung including identity satisfies the %s band: %s", requested, selection.Row.NoArtifactReason)
			case LadderPass:
				record := selection.Record
				candidates = append(candidates, candidate{requested, &record, true})
			}
		}
	default:
		if requested == "compact" {
			if g, ok := table.Compact[in.Geometry.ID]; ok {
				candidates = append(candidates, candidate{"compact", &g, false})
			}
			if g, ok := table.Standard[in.Geometry.ID]; ok {
				candidates = append(candidates, candidate{"standard", &g, false})
			}
		}
		if requested == "standard" {
			if g, ok := table.Standard[in.Geometry.ID]; ok {
				candidates = append(candidates, candidate{"standard", &g, false})
			}
		}
	}
	candidates = append(candidates, candidate{name: "source"})
	prov := LODProvenance{
		RequestedTier: requested, ResolvedTolerance: quality.Simplification,
		Quantization: quality.Quantization, Projection: projectionPolicy,
	}
	var selected, canonical MultiPolygon
	selectedCeiling := in.MaxPathBytes
	selectedTransform, selectedViewBox := tr, vb
	retainedProtected := map[int][]string{}
	for i, component := range fullComponents {
		if len(component.Protected) > 0 {
			retainedProtected[i] = component.Protected
		}
	}
	for candidateIndex, c := range candidates {
		var candidateGeometry MultiPolygon
		var candidateCanonical MultiPolygon
		var candidateTransform Transform
		var candidateViewBox Bounds
		var restorations []Removal
		rawDeviation, finalDeviation := 0.0, 0.0
		coordinateSpace, finalization := "centered_laea", "precomputed"
		sourceAttempted, sourceSelected := 0.0, 0.0
		var attemptedPhase, selectedPhase GridPhase
		phaseAttempts := 0
		if c.name == "source" {
			coordinateSpace, finalization = "centered_laea", "source_reduced"
			finalized, finalizationErr := finalizeSourceLOD(fullRequired, retainedProtected, quality, tr, vb, in, prj, minimumParts)
			prov.Fallbacks = append(prov.Fallbacks, finalized.Fallbacks...)
			if finalizationErr != nil {
				prov.Fallbacks = append(prov.Fallbacks, c.name+":finalization:"+errorIdentity(finalizationErr))
				continue
			}
			candidateGeometry, candidateCanonical = finalized.Geometry, finalized.Canonical
			candidateTransform, candidateViewBox = finalized.Transform, finalized.ViewBox
			rawDeviation, finalDeviation = finalized.RawDeviation, finalized.FinalDeviation
			sourceAttempted, sourceSelected = finalized.AttemptedTolerance, finalized.SelectedTolerance
			attemptedPhase, selectedPhase = finalized.AttemptedPhase, finalized.SelectedPhase
			phaseAttempts = finalized.PhaseAttempts
		} else {
			projectedLOD, bindingErr := bindProjectedLOD(*c.record, in, table, prj, rotation)
			if bindingErr != nil {
				prov.Fallbacks = append(prov.Fallbacks, c.name+":binding:"+errorIdentity(bindingErr))
				continue
			}
			candidateGeometry = transformGeometry(projectedLOD, tr.Scale, tr.TranslateX, tr.TranslateY)
			if !c.ladder {
				// The pre-DEC-006 derived path: inject full-detail source
				// components the candidate dropped, then judge the result by
				// raw boundary deviation. Both are exactly what DEC-006
				// replaced with the silhouette oracle, so a ladder candidate
				// gets neither — its omissions were judged deliberately at
				// build time and restoring them would serve geometry no oracle
				// ever saw.
				candidateGeometry, restorations = restoreRequiredComponents(fullRequired, fullComponents, candidateGeometry)
			}
			shared := sharedSourceVertices(fullRequired, 1e-12)
			// Same boundary as the oracle's: topology is judged on the canonical
			// geometry, which canonicalizeShared validates on the way out, not
			// on the pre-canonical fitted coordinates. Simplification can leave a
			// self-crossing four orders of magnitude below the q=0.01 output grid
			// that canonicalization then erases; rejecting here would discard a
			// candidate the build already accepted and fall back to source —
			// which for Russia meant 101605 bytes against a 2500 ceiling.
			//
			// The build and the runtime have to agree on this or the ladder is
			// selected against one predicate and served against another, which is
			// the whole class of failure this spec exists to end.
			if !c.ladder {
				if topologyErr := validateTopologyShared(candidateGeometry, shared); topologyErr != nil {
					prov.Fallbacks = append(prov.Fallbacks, c.name+":topology:"+errorIdentity(topologyErr))
					continue
				}
			}
			if !c.ladder {
				rawDeviation = matchedBoundaryDeviation(fullRequired, candidateGeometry, quality.Simplification)
				if !finite(rawDeviation) || rawDeviation > quality.Simplification {
					prov.Fallbacks = append(prov.Fallbacks, fmt.Sprintf("%s:raw_deviation:%g", c.name, rawDeviation))
					continue
				}
			}
			phaseResult, phaseErr := finalizeGridPhases(
				unfitGeometry(candidateGeometry, tr),
				unfitGeometry(fullRequired, tr),
				tr, vb, in, prj, quality, minimumParts, c.ladder,
			)
			if phaseErr != nil {
				prov.Fallbacks = append(prov.Fallbacks, c.name+":phase:"+errorIdentity(phaseErr))
				continue
			}
			candidateGeometry, candidateCanonical = phaseResult.Geometry, phaseResult.Canonical
			candidateTransform, candidateViewBox = phaseResult.Transform, phaseResult.ViewBox
			finalDeviation = phaseResult.FinalDeviation
			attemptedPhase, selectedPhase = phaseResult.AttemptedPhase, phaseResult.SelectedPhase
			phaseAttempts = phaseResult.Attempts
		}
		// A ladder candidate carries the oracle's promise that its linear path
		// fits the band's frozen path cap. Softening is a runtime-only
		// representation choice the oracle never judged, and on small compact
		// geometries it inflates the path several-fold — enough to clear the
		// band cap while still sitting under the preset's larger complete-file
		// maximum. Holding the softened variant to the band cap keeps the
		// guarantee the ladder was built to make; the DEC-005 source path,
		// which promises nothing of the kind, is unaffected.
		renderCeiling := in.MaxPathBytes
		if c.ladder {
			if cap := table.bandPathCap(c.name); cap > 0 && (renderCeiling == 0 || cap < renderCeiling) {
				renderCeiling = cap
			}
		}
		if candidateIndex != len(candidates)-1 {
			if _, _, _, budgetErr := renderLODRepresentation(candidateCanonical, quality, renderCeiling); budgetErr != nil {
				if pipelineErr, ok := budgetErr.(*PipelineError); ok && pipelineErr.Code == ErrBudget {
					prov.Fallbacks = append(prov.Fallbacks, c.name+":hard_budget:"+pipelineErr.Message)
				} else {
					prov.Fallbacks = append(prov.Fallbacks, c.name+":render:"+errorIdentity(budgetErr))
				}
				continue
			}
		}
		selected, canonical = candidateGeometry, candidateCanonical
		selectedCeiling = renderCeiling
		selectedTransform, selectedViewBox = candidateTransform, candidateViewBox
		prov.SelectedTier, prov.Restored = c.name, restorations
		prov.RawDeviation, prov.FinalDeviation, prov.MaximumDeviation = rawDeviation, finalDeviation, finalDeviation
		prov.CoordinateSpace, prov.Finalization = coordinateSpace, finalization
		prov.SourceAttemptedTolerance, prov.SourceSelectedTolerance = sourceAttempted, sourceSelected
		prov.AttemptedPhase, prov.SelectedPhase = attemptedPhase, selectedPhase
		prov.PhaseAttempts = phaseAttempts
		prov.OutputPoints = countPoints(candidateCanonical)
		if selectionOnly {
			break
		}
		break
	}
	if len(selected) == 0 {
		return Result{}, fail(ErrTopology, in.Entity.Alpha2, "lod", "no valid tier: %v", prov.Fallbacks)
	}
	if selectionOnly {
		return Result{Entity: in.Entity.Alpha2, Profile: in.Profile, Preset: in.Preset, EffectiveScale: effective, LOD: prov}, nil
	}
	cmds, path, renderDiagnostics, err := renderLODRepresentation(canonical, quality, selectedCeiling)
	if err != nil {
		return Result{}, err
	}
	prov.ParserRoundTrip = true
	diags = append(diags, renderDiagnostics...)
	markers, err := projectMarkers(in.Markers, prj, selectedTransform, selectedViewBox, offsets)
	if err != nil {
		return Result{}, err
	}
	for _, marker := range markers {
		if marker.Anomaly != "" {
			diags = append(diags, Diagnostic{"marker_" + marker.Anomaly, "warning", marker.ID + " is outside fitted result"})
		}
	}
	sort.Slice(diags, func(i, j int) bool {
		if diags[i].Code == diags[j].Code {
			return diags[i].Message < diags[j].Message
		}
		return diags[i].Code < diags[j].Code
	})
	sort.Slice(fullRemovals, func(i, j int) bool {
		if fullRemovals[i].SourcePolygon == fullRemovals[j].SourcePolygon {
			return fullRemovals[i].SourceRing < fullRemovals[j].SourceRing
		}
		return fullRemovals[i].SourcePolygon < fullRemovals[j].SourcePolygon
	})
	return Result{SchemaVersion: SchemaVersion, AlgorithmVersion: AlgorithmVersion, CorpusID: in.CorpusID, Entity: in.Entity.Alpha2, Profile: in.Profile, Preset: in.Preset, LayoutMode: in.Layout.Mode, ViewBox: selectedViewBox, NaturalAspect: geometryBounds(projected).Width() / geometryBounds(projected).Height(), EffectiveScale: effective, Transform: selectedTransform, Quality: quality, Commands: cmds, Path: path, Components: fullComponents, Removals: fullRemovals, Markers: markers, Diagnostics: diags, Metrics: Metrics{InputPoints: countPoints(full), ProjectedPoints: projectedPoints, OutputPoints: countPoints(canonical), PathBytes: len(path)}, LOD: prov}, nil
}

func renderLODRepresentation(canonical MultiPolygon, quality Quality, maximumBytes int) ([]Command, string, []Diagnostic, error) {
	var diagnostics []Diagnostic
	if quality.Softening > 0 {
		deviation, topologyErr := validateSoftenedGeometry(canonical, quality.Softening)
		switch {
		case topologyErr != nil:
			diagnostics = append(diagnostics, Diagnostic{"softening_topology_fallback", "warning", topologyErr.Error()})
		case !finite(deviation) || deviation > quality.Softening:
			diagnostics = append(diagnostics, Diagnostic{"softening_tolerance_fallback", "warning", fmt.Sprintf("softened displacement %g exceeds tolerance %g", deviation, quality.Softening)})
		default:
			commands, path, serializationErr := commandsAndPath(canonical, quality.Softening)
			if serializationErr != nil {
				diagnostics = append(diagnostics, Diagnostic{"softening_serialization_fallback", "warning", serializationErr.Error()})
			} else if maximumBytes == 0 || len(path) <= maximumBytes {
				return commands, path, diagnostics, nil
			} else {
				diagnostics = append(diagnostics, Diagnostic{"softening_budget_fallback", "warning", fmt.Sprintf("softened path bytes %d exceed hard maximum %d", len(path), maximumBytes)})
			}
		}
	}
	commands, path, err := commandsAndPath(canonical, 0)
	if err != nil {
		return nil, "", diagnostics, err
	}
	if maximumBytes > 0 && len(path) > maximumBytes {
		return nil, "", diagnostics, fail(ErrBudget, "", "max_path_bytes", "linear path bytes %d exceed hard maximum %d", len(path), maximumBytes)
	}
	return commands, path, diagnostics, nil
}

// ladderVerdictApplies reports whether the committed ladder's verdict was
// earned under conditions this request actually reproduces. The build judged
// every row through EvaluateSilhouetteCandidate at rotation 0, under the v1
// projection contract, at the automatic quality for the fitted scale, and with
// the caller's byte ceiling cleared. A request that changes any of those is a
// different question than the one the oracle answered, so it takes the
// unchanged DEC-005 source path instead of inheriting an answer that was never
// about it.
//
// This is deliberately one named predicate rather than conditions scattered
// through the candidate loop: it is the seam that decides whether a committed
// oracle verdict may be trusted, and it has to be reviewable in one place.
func ladderVerdictApplies(in Input, rotation float64, quality Quality, artifactProjectionCompatible bool) bool {
	if !artifactProjectionCompatible {
		return false // an explicit flatness rebuilds the projection the ladder was built in
	}
	if rotation != 0 {
		return false // the ladder was built and judged at rotation 0
	}
	if in.Override != nil && in.Override.Quality != nil {
		return false // an explicit quality replaces the AutoQuality the oracle used
	}
	if !in.Quality.Auto && in.Quality.Quantization > 0 && in.Quality.Flatness > 0 {
		return false // likewise for an explicit request quality
	}
	return quality.Auto || quality.Quantization > 0
}

func errorIdentity(err error) string {
	if pipelineErr, ok := err.(*PipelineError); ok {
		return string(pipelineErr.Code) + ":" + pipelineErr.Field + ":" + pipelineErr.Message
	}
	return err.Error()
}

func bindProjectedLOD(record ProjectedLODGeometry, in Input, table *LODTable, prj projector, rotation float64) (MultiPolygon, error) {
	const metadataTolerance = 1e-12
	if table == nil || table.RecipeSHA256 == "" {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_recipe", "missing runtime recipe identity")
	}
	if record.Geometry.ID != in.Geometry.ID {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_geometry_id", "got %q want %q", record.Geometry.ID, in.Geometry.ID)
	}
	if record.SourceCorpus == "" || record.SourceCorpus != in.CorpusID {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_source_corpus", "got %q want %q", record.SourceCorpus, in.CorpusID)
	}
	if record.RecipeSHA256 == "" || record.RecipeSHA256 != table.RecipeSHA256 {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_recipe", "candidate recipe %q does not match table recipe %q", record.RecipeSHA256, table.RecipeSHA256)
	}
	metadata := record.Projection
	contract := LODProjectionContractV1()
	if metadata.Version != contract.Version {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_projection_version", "got %q want %q", metadata.Version, contract.Version)
	}
	if metadata.CoordinateSpace != contract.CoordinateSpace {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_coordinate_space", "got %q want %q", metadata.CoordinateSpace, contract.CoordinateSpace)
	}
	if metadata.BuildRotation != contract.BuildRotation {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_build_rotation", "got %g want %g", metadata.BuildRotation, contract.BuildRotation)
	}
	if metadata.YAxis != contract.YAxis {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_y_axis", "got %q want %q", metadata.YAxis, contract.YAxis)
	}
	if math.Abs(metadata.Flatness-contract.Flatness) > metadataTolerance {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_flatness", "got %g want %g", metadata.Flatness, contract.Flatness)
	}
	if math.Abs(metadata.CoordinatePrecision-contract.CoordinatePrecision) > metadataTolerance {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_coordinate_precision", "got %g want %g", metadata.CoordinatePrecision, contract.CoordinatePrecision)
	}
	if math.Abs(metadata.CenterLongitude-deg(prj.lon0)) > metadataTolerance || math.Abs(metadata.CenterLatitude-deg(prj.lat0)) > metadataTolerance {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_center", "got (%g,%g) want (%g,%g)", metadata.CenterLongitude, metadata.CenterLatitude, deg(prj.lon0), deg(prj.lat0))
	}
	projected, err := normalize(fromCatalog(record.Geometry.Coordinates))
	if err != nil {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "lod_geometry", "%v", err)
	}
	return rotateProjectedLOD(projected, rotation), nil
}

type sourceLODResult struct {
	Geometry, Canonical                   MultiPolygon
	Transform                             Transform
	ViewBox                               Bounds
	RawDeviation, FinalDeviation          float64
	AttemptedTolerance, SelectedTolerance float64
	AttemptedPhase, SelectedPhase         GridPhase
	PhaseAttempts                         int
	Fallbacks                             []string
}

func finalizeSourceLOD(full MultiPolygon, protected map[int][]string, quality Quality, base Transform, viewBox Bounds, in Input, prj projector, minimumParts int) (sourceLODResult, error) {
	var result sourceLODResult
	rawShared := sharedSourceVertices(full, 1e-12)
	if err := validateTopologyShared(full, rawShared); err != nil {
		return result, err
	}
	remaining := quality.Simplification - quality.Quantization/math.Sqrt2
	result.AttemptedTolerance = remaining
	if !finite(remaining) || remaining <= 0 {
		return result, fail(ErrTopology, "", "source_finalization", "quantization reserve leaves no simplification tolerance")
	}
	unfittedFull := unfitGeometry(full, base)
	for attempt := remaining; attempt > .0001; {
		reduced, _, selectedTolerance, err := simplifySharedResolved(full, attempt, protected, rawShared)
		if err != nil {
			result.Fallbacks = append(result.Fallbacks, fmt.Sprintf("source:simplification:%g:%s", attempt, errorIdentity(err)))
			attempt /= 2
			continue
		}
		if !sameLODStructure(full, reduced) {
			result.Fallbacks = append(result.Fallbacks, fmt.Sprintf("source:retention:%g:component_or_ring_loss", selectedTolerance))
			attempt = nextSourceTolerance(attempt, selectedTolerance)
			continue
		}
		rawDeviation := matchedBoundaryDeviation(full, reduced, quality.Simplification)
		if !finite(rawDeviation) || rawDeviation > quality.Simplification {
			result.Fallbacks = append(result.Fallbacks, fmt.Sprintf("source:raw_deviation:%g:%g", selectedTolerance, rawDeviation))
			attempt = nextSourceTolerance(attempt, selectedTolerance)
			continue
		}
		phaseResult, err := finalizeGridPhases(
			unfitGeometry(reduced, base), unfittedFull,
			base, viewBox, in, prj, quality, minimumParts, false,
		)
		result.PhaseAttempts += phaseResult.Attempts
		result.AttemptedPhase = phaseResult.AttemptedPhase
		if err != nil {
			result.Fallbacks = append(result.Fallbacks, fmt.Sprintf("source:phase:%g:%s", selectedTolerance, errorIdentity(err)))
			attempt = nextSourceTolerance(attempt, selectedTolerance)
			continue
		}
		result.Geometry, result.Canonical = phaseResult.Geometry, phaseResult.Canonical
		result.Transform, result.ViewBox = phaseResult.Transform, phaseResult.ViewBox
		result.RawDeviation, result.FinalDeviation = rawDeviation, phaseResult.FinalDeviation
		result.SelectedTolerance = selectedTolerance
		result.SelectedPhase = phaseResult.SelectedPhase
		return result, nil
	}
	return result, fail(ErrTopology, "", "source_finalization", "no q-aware source tolerance remained after %d rejected candidates", len(result.Fallbacks))
}

type phaseLODResult struct {
	Geometry, Canonical           MultiPolygon
	Transform                     Transform
	ViewBox                       Bounds
	FinalDeviation                float64
	AttemptedPhase, SelectedPhase GridPhase
	Attempts                      int
}

func finalizeGridPhases(candidateUnfitted, referenceUnfitted MultiPolygon, base Transform, viewBox Bounds, in Input, prj projector, quality Quality, minimumParts int, oracleJudged bool) (phaseLODResult, error) {
	var result phaseLODResult
	var failures []string
	for _, phase := range gridPhaseSchedule() {
		result.Attempts++
		result.AttemptedPhase = phase
		attempt, err := finalizeGridPhase(candidateUnfitted, referenceUnfitted, base, viewBox, in, prj, quality, minimumParts, phase, oracleJudged)
		if err != nil {
			failures = append(failures, fmt.Sprintf("(%0.3f,%0.3f):%s", phase.X, phase.Y, gridPhaseFailure(err)))
			continue
		}
		attempt.Attempts = result.Attempts
		attempt.AttemptedPhase = phase
		result = attempt
		return result, nil
	}
	return result, fail(ErrTopology, in.Entity.Alpha2, "grid_phase", "exhausted %d phases: %v", result.Attempts, failures)
}

// oracleJudged marks a candidate whose fidelity was already decided at build
// time by the DEC-006 silhouette oracle. Such a candidate still has to pass
// every structural check here — canonicalization, structure, on-grid,
// containment and protection — but not the raw boundary-deviation gate, which
// DEC-006 replaced as the fidelity boundary. The DEC-005 source path keeps that
// gate: it is the path's own acceptance criterion, and nothing else judges it.
func finalizeGridPhase(candidateUnfitted, referenceUnfitted MultiPolygon, base Transform, viewBox Bounds, in Input, prj projector, quality Quality, minimumParts int, phase GridPhase, oracleJudged bool) (phaseLODResult, error) {
	result := phaseLODResult{Attempts: 1, AttemptedPhase: phase, SelectedPhase: phase}
	transform := composeGridPhase(base, phase)
	candidate := materializeGridPhase(candidateUnfitted, transform)
	reference := materializeGridPhase(referenceUnfitted, transform)
	canonical, err := canonicalizeShared(candidate, quality.Quantization, sharedSourceVertices(reference, quality.Quantization))
	if err != nil {
		return result, fail(ErrTopology, in.Entity.Alpha2, "canonicalization", "%s", errorIdentity(err))
	}
	if !sameLODStructure(candidate, canonical) {
		return result, fail(ErrTopology, in.Entity.Alpha2, "structure", "component or ring identity changed")
	}
	if !geometryOnGrid(canonical, quality.Quantization) {
		return result, fail(ErrTopology, in.Entity.Alpha2, "off_grid", "canonical geometry is not on the fixed grid")
	}
	shiftedViewBox := shiftViewBox(viewBox, phase)
	if !geometryContainedBy(canonical, shiftedViewBox) {
		return result, fail(ErrTopology, in.Entity.Alpha2, "containment", "canonical geometry exceeds shifted viewBox")
	}
	if err := validatePhaseProtection(in, prj, transform, canonical, minimumParts, oracleJudged); err != nil {
		return result, fail(ErrProtected, in.Entity.Alpha2, "protection", "%s", errorIdentity(err))
	}
	finalDeviation := matchedBoundaryDeviation(reference, canonical, quality.Simplification)
	result.FinalDeviation = finalDeviation
	if !oracleJudged && (!finite(finalDeviation) || finalDeviation > quality.Simplification) {
		return result, fail(ErrTopology, in.Entity.Alpha2, "final_deviation", "%g", finalDeviation)
	}
	result.Geometry, result.Canonical = candidate, canonical
	result.Transform, result.ViewBox = transform, shiftedViewBox
	return result, nil
}

func gridPhaseFailure(err error) string {
	pipelineErr, ok := err.(*PipelineError)
	if !ok {
		return errorIdentity(err)
	}
	switch pipelineErr.Field {
	case "structure", "off_grid", "containment":
		return pipelineErr.Field
	case "canonicalization", "protection", "final_deviation":
		return pipelineErr.Field + ":" + pipelineErr.Message
	default:
		return errorIdentity(err)
	}
}

func nextSourceTolerance(attempt, selected float64) float64 {
	if selected > 0 && selected < attempt {
		return selected / 2
	}
	return attempt / 2
}

func sameLODStructure(a, b MultiPolygon) bool {
	if len(a) != len(b) {
		return false
	}
	aRings := make([]int, len(a))
	bRings := make([]int, len(b))
	for i := range a {
		aRings[i] = len(a[i])
		bRings[i] = len(b[i])
	}
	sort.Ints(aRings)
	sort.Ints(bRings)
	for i := range aRings {
		if aRings[i] != bRings[i] {
			return false
		}
	}
	return true
}

func lodPreparationKey(in Input, rotation, flatness float64) (string, error) {
	raw, err := json.Marshal(struct {
		Corpus, GeometryID string
		Geometry           catalog.MultiPolygon
		Rotation, Flatness float64
	}{in.CorpusID, in.Geometry.ID, in.Geometry.Coordinates, rotation, flatness})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func restoreRequiredComponents(full MultiPolygon, components []Component, candidate MultiPolygon) (MultiPolygon, []Removal) {
	out := make(MultiPolygon, 0, len(full))
	used := map[int]bool{}
	var restored []Removal
	for fi, fp := range full {
		center := ringCentroid(fp[0])
		matched := -1
		for ci, cp := range candidate {
			if used[ci] || len(cp) == 0 {
				continue
			}
			if pointInRing(center, cp[0]) || pointInRing(ringCentroid(cp[0]), fp[0]) {
				matched = ci
				break
			}
		}
		if matched < 0 || len(candidate[matched]) < len(fp) {
			out = append(out, fp)
			restored = append(restored, Removal{SourcePolygon: components[fi].SourcePolygon, SourceRing: -1, Reason: "lod_required_restoration", Area: polygonArea(fp)})
			continue
		}
		used[matched] = true
		out = append(out, candidate[matched])
	}
	return out, restored
}
func ringCentroid(r Ring) Point {
	var x, y float64
	n := len(r) - 1
	if n <= 0 {
		return Point{}
	}
	for _, p := range r[:n] {
		x += p.X
		y += p.Y
	}
	return Point{x / float64(n), y / float64(n)}
}
func matchedBoundaryDeviation(full, candidate MultiPolygon, limit float64) float64 {
	if len(full) != len(candidate) {
		return math.Inf(1)
	}
	used := make([]bool, len(candidate))
	maximum := 0.0
	for _, sourcePolygon := range full {
		match := -1
		best := math.Inf(1)
		sourceCenter := ringCentroid(sourcePolygon[0])
		for candidateIndex, candidatePolygon := range candidate {
			if used[candidateIndex] || len(sourcePolygon) != len(candidatePolygon) {
				continue
			}
			d := distance(sourceCenter, ringCentroid(candidatePolygon[0]))
			if d < best {
				best, match = d, candidateIndex
			}
		}
		if match < 0 {
			return math.Inf(1)
		}
		used[match] = true
		candidatePolygon := candidate[match]
		d := symmetricRingDeviation(sourcePolygon[0], candidatePolygon[0], limit)
		if d > maximum {
			maximum = d
		}
		if maximum > limit {
			return maximum
		}
		holeUsed := make([]bool, len(candidatePolygon))
		holeUsed[0] = true
		for sourceRingIndex := 1; sourceRingIndex < len(sourcePolygon); sourceRingIndex++ {
			holeMatch := -1
			holeBest := math.Inf(1)
			sourceRingCenter := ringCentroid(sourcePolygon[sourceRingIndex])
			for candidateRingIndex := 1; candidateRingIndex < len(candidatePolygon); candidateRingIndex++ {
				if holeUsed[candidateRingIndex] {
					continue
				}
				d := distance(sourceRingCenter, ringCentroid(candidatePolygon[candidateRingIndex]))
				if d < holeBest {
					holeBest, holeMatch = d, candidateRingIndex
				}
			}
			if holeMatch < 0 {
				return math.Inf(1)
			}
			holeUsed[holeMatch] = true
			d := symmetricRingDeviation(sourcePolygon[sourceRingIndex], candidatePolygon[holeMatch], limit)
			if d > maximum {
				maximum = d
			}
			if maximum > limit {
				return maximum
			}
		}
	}
	return maximum
}

type lodSegment struct{ a, b Point }
type lodSegmentIndex struct {
	cell  float64
	cells map[[2]int][]lodSegment
}

func newLODIndex(r Ring, limit float64) lodSegmentIndex {
	cell := math.Max(limit, .05)
	idx := lodSegmentIndex{cell: cell, cells: map[[2]int][]lodSegment{}}
	for i := 0; i < len(r)-1; i++ {
		s := lodSegment{r[i], r[i+1]}
		minX, maxX := math.Min(s.a.X, s.b.X)-limit, math.Max(s.a.X, s.b.X)+limit
		minY, maxY := math.Min(s.a.Y, s.b.Y)-limit, math.Max(s.a.Y, s.b.Y)+limit
		for x := int(math.Floor(minX / cell)); x <= int(math.Floor(maxX/cell)); x++ {
			for y := int(math.Floor(minY / cell)); y <= int(math.Floor(maxY/cell)); y++ {
				idx.cells[[2]int{x, y}] = append(idx.cells[[2]int{x, y}], s)
			}
		}
	}
	return idx
}
func symmetricRingDeviation(a, b Ring, limit float64) float64 {
	ai, bi := newLODIndex(a, limit), newLODIndex(b, limit)
	max := 0.0
	for _, p := range a[:len(a)-1] {
		d := bi.distance(p)
		if d > max {
			max = d
		}
		if max > limit {
			return max
		}
	}
	for _, p := range b[:len(b)-1] {
		d := ai.distance(p)
		if d > max {
			max = d
		}
		if max > limit {
			return max
		}
	}
	return max
}
func (i lodSegmentIndex) distance(p Point) float64 {
	best := math.Inf(1)
	key := [2]int{int(math.Floor(p.X / i.cell)), int(math.Floor(p.Y / i.cell))}
	for _, s := range i.cells[key] {
		d := pointSegmentDistanceLOD(p, s.a, s.b)
		if d < best {
			best = d
		}
	}
	return best
}
func pointSegmentDistanceLOD(p, a, b Point) float64 {
	dx, dy := b.X-a.X, b.Y-a.Y
	if dx == 0 && dy == 0 {
		return distance(p, a)
	}
	t := clamp(((p.X-a.X)*dx+(p.Y-a.Y)*dy)/(dx*dx+dy*dy), 0, 1)
	return distance(p, Point{a.X + t*dx, a.Y + t*dy})
}
