package geometry

import (
	"fmt"
	"math"
)

type ErrorCode string

const (
	ErrInvalidSource ErrorCode = "invalid_source"
	ErrInvalidLayout ErrorCode = "invalid_layout"
	ErrProjection    ErrorCode = "projection_failure"
	ErrTopology      ErrorCode = "topology_damage"
	ErrEmpty         ErrorCode = "empty_result"
	ErrProtected     ErrorCode = "protected_feature_loss"
	ErrMarker        ErrorCode = "marker_anomaly"
	ErrOverride      ErrorCode = "invalid_override"
	ErrPointLimit    ErrorCode = "point_limit"
	ErrSerialization ErrorCode = "serialization_failure"
	ErrBudget        ErrorCode = "hard_budget_exceeded"
	// ErrNoArtifact is DEC-009's typed outcome for a band that cannot be drawn:
	// no ladder rung, identity included, satisfies topology, protected
	// visibility and the frozen byte caps at the fitted scale. It is a decision
	// the build recorded, not a pipeline failure, and callers are expected to
	// test for it and render absence deliberately (DEC-010).
	ErrNoArtifact ErrorCode = "no_artifact"
)

// IsNoArtifact reports whether err is the DEC-009 typed no-artifact outcome.
// Consumers use this to tell "this band is deliberately undrawable" apart from
// "generation failed", which every other error code means.
func IsNoArtifact(err error) bool {
	pipelineErr, ok := err.(*PipelineError)
	return ok && pipelineErr.Code == ErrNoArtifact
}

type PipelineError struct {
	Code                   ErrorCode
	Entity, Field, Message string
}

func (e *PipelineError) Error() string {
	return fmt.Sprintf("geometry %s entity=%s field=%s: %s", e.Code, e.Entity, e.Field, e.Message)
}

type Diagnostic struct {
	Code, Severity, Message string
}

// DiagnosticPhaseRow is maintainer-only evidence for one explicitly requested
// fixed-grid phase. Production selection never consumes this type or its helper.
type DiagnosticPhaseRow struct {
	Phase          GridPhase `json:"phase"`
	Accepted       bool      `json:"accepted"`
	Failure        string    `json:"failure,omitempty"`
	FinalDeviation *float64  `json:"final_deviation,omitempty"`
}

// DiagnosticGridPhases evaluates an explicit diagnostic phase list through the
// canonicalization, topology, protection, containment and final-deviation
// primitive over the shared 100-phase schedule.
//
// It reproduces the *superseded* pre-DEC-006 derived path — component
// restoration plus a deviation verdict — because that is the path the T0A.1
// spike measured and the evidence it produced is only comparable against it.
// Production no longer selects that way: a ladder candidate is judged by the
// silhouette oracle at build time and takes neither the restoration nor the
// deviation gate. Do not read this helper as a mirror of current selection.
func DiagnosticGridPhases(raw Input, table *LODTable, tier string, phases []GridPhase) ([]DiagnosticPhaseRow, error) {
	in, err := ApplyPreset(raw)
	if err != nil {
		return nil, err
	}
	full, err := normalize(fromCatalog(in.Geometry.Coordinates))
	if err != nil {
		return nil, err
	}
	contract := LODProjectionContractV1()
	prj := centeredProjector(full, 0)
	projected, _, err := projectGeometry(full, prj, contract.Flatness)
	if err != nil {
		return nil, err
	}
	fitted, viewBox, transform, effective, _, err := fitGeometry(projected, in.Layout)
	if err != nil {
		return nil, err
	}
	transform.CenterLon, transform.CenterLat = deg(prj.lon0), deg(prj.lat0)
	quality := in.Quality
	if quality.Auto || quality.Quantization == 0 {
		quality = AutoQuality(effective)
	}
	if err := validateQuality(quality); err != nil {
		return nil, err
	}
	minimumParts := 1
	for _, feature := range in.Entity.Protected {
		if feature.MinimumParts > minimumParts {
			minimumParts = feature.MinimumParts
		}
	}
	protected, err := protectedComponents(in.Entity, prj, transform, fitted)
	if err != nil {
		return nil, err
	}
	required, _, components, err := retain(fitted, protected, minimumParts, quality.MinimumArea)
	if err != nil {
		return nil, err
	}
	if table == nil {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "diagnostic_table", "table is required")
	}
	var record ProjectedLODGeometry
	switch tier {
	case "compact":
		record = table.Compact[in.Geometry.ID]
	case "standard":
		record = table.Standard[in.Geometry.ID]
	default:
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "diagnostic_tier", "unsupported tier %q", tier)
	}
	if record.Geometry.ID == "" {
		return nil, fail(ErrInvalidSource, in.Entity.Alpha2, "diagnostic_tier", "missing %s candidate", tier)
	}
	projectedLOD, err := bindProjectedLOD(record, in, table, prj, 0)
	if err != nil {
		return nil, err
	}
	candidate := transformGeometry(projectedLOD, transform.Scale, transform.TranslateX, transform.TranslateY)
	candidate, _ = restoreRequiredComponents(required, components, candidate)
	shared := sharedSourceVertices(required, 1e-12)
	if err := validateTopologyShared(candidate, shared); err != nil {
		return nil, err
	}
	candidateUnfitted := unfitGeometry(candidate, transform)
	referenceUnfitted := unfitGeometry(required, transform)
	rows := make([]DiagnosticPhaseRow, 0, len(phases))
	for _, phase := range phases {
		result, err := finalizeGridPhase(candidateUnfitted, referenceUnfitted, transform, viewBox, in, prj, quality, minimumParts, phase, false)
		row := DiagnosticPhaseRow{Phase: phase}
		if err != nil {
			row.Failure = errorIdentity(err)
			if result.FinalDeviation != 0 && finite(result.FinalDeviation) {
				value := result.FinalDeviation
				row.FinalDeviation = &value
			}
		} else {
			row.Accepted = true
			if finite(result.FinalDeviation) && !math.IsNaN(result.FinalDeviation) {
				value := result.FinalDeviation
				row.FinalDeviation = &value
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func fail(code ErrorCode, entity, field, format string, args ...any) error {
	return &PipelineError{Code: code, Entity: entity, Field: field, Message: fmt.Sprintf(format, args...)}
}
