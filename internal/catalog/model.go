package catalog

import "fmt"

const (
	SchemaVersion = 1
	CorpusVersion = 1
)

type Point [2]float64
type Ring []Point
type Polygon []Ring
type MultiPolygon []Polygon

type Geometry struct {
	ID          string       `json:"id"`
	Coordinates MultiPolygon `json:"coordinates"`
}

type ProfileResolution struct {
	GeometryID  string `json:"geometry_id,omitempty"`
	IdenticalTo string `json:"identical_to,omitempty"`
}

type BoundaryProfiles struct {
	UN      ProfileResolution `json:"un"`
	DeFacto ProfileResolution `json:"de_facto"`
}

type Capital struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Point   Point    `json:"point"`
	Roles   []string `json:"roles"`
	Primary bool     `json:"primary,omitempty"`
}

type ProtectedFeature struct {
	Name         string `json:"name"`
	Anchor       Point  `json:"anchor"`
	MinimumParts int    `json:"minimum_parts"`
}

type Entity struct {
	Alpha2    string             `json:"alpha2"`
	Alpha3    string             `json:"alpha3"`
	Name      string             `json:"name"`
	Profiles  BoundaryProfiles   `json:"profiles"`
	Capitals  []Capital          `json:"capitals"`
	Protected []ProtectedFeature `json:"protected_features"`
}

type Receipt struct {
	ID              string `json:"id"`
	Path            string `json:"path"`
	Origin          string `json:"origin"`
	UpstreamVersion string `json:"upstream_version"`
	License         string `json:"license"`
	SHA256          string `json:"sha256"`
	Transformation  string `json:"transformation"`
}

type Manifest struct {
	SchemaVersion int      `json:"schema_version"`
	CorpusVersion int      `json:"corpus_version"`
	Identity      string   `json:"identity"`
	EntityCount   int      `json:"entity_count"`
	Profiles      []string `json:"profiles"`
	Entities      []Entity `json:"entities"`
}

type Coverage struct {
	EntityCount      int      `json:"entity_count"`
	UNExplicit       int      `json:"un_explicit"`
	DeFactoExplicit  int      `json:"de_facto_explicit"`
	DeFactoIdentical int      `json:"de_facto_identical"`
	CapitalEntities  int      `json:"capital_entities"`
	CapitalRecords   int      `json:"capital_records"`
	Protected        []string `json:"protected_entities"`
}

type Corpus struct {
	Manifest   Manifest   `json:"manifest"`
	Geometries []Geometry `json:"geometries"`
	Receipts   []Receipt  `json:"receipts"`
	Coverage   Coverage   `json:"coverage"`
}

type DiagnosticError struct {
	Source    string
	Entity    string
	Field     string
	Invariant string
	Message   string
}

func (e *DiagnosticError) Error() string {
	return fmt.Sprintf("source=%s entity=%s field=%s invariant=%s: %s", e.Source, e.Entity, e.Field, e.Invariant, e.Message)
}

func diagnostic(source, entity, field, invariant, format string, args ...any) error {
	return &DiagnosticError{Source: source, Entity: entity, Field: field, Invariant: invariant, Message: fmt.Sprintf(format, args...)}
}
