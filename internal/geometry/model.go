package geometry

import "github.com/yuranikolaev/country-map-svg-generator/internal/catalog"

const (
	SchemaVersion    = 2
	AlgorithmVersion = "laea-organic-v1"
	MaxPoints        = 2_000_000
)

type Point struct{ X, Y float64 }
type Ring []Point
type Polygon []Ring
type MultiPolygon []Polygon

type LODProjection struct {
	Version             string  `json:"version"`
	CoordinateSpace     string  `json:"coordinate_space"`
	CenterLongitude     float64 `json:"center_longitude"`
	CenterLatitude      float64 `json:"center_latitude"`
	BuildRotation       float64 `json:"build_rotation"`
	YAxis               string  `json:"y_axis"`
	Flatness            float64 `json:"flatness"`
	CoordinatePrecision float64 `json:"coordinate_precision"`
}

type ProjectedLODGeometry struct {
	Geometry     catalog.Geometry `json:"geometry"`
	Projection   LODProjection    `json:"projection"`
	SourceCorpus string           `json:"source_corpus"`
	RecipeSHA256 string           `json:"recipe_sha256"`
}

type Insets struct{ Top, Right, Bottom, Left float64 }

type LayoutMode string

const (
	LayoutTight   LayoutMode = "tight"
	LayoutContain LayoutMode = "contain"
)

type Layout struct {
	Mode                LayoutMode
	Width, Height       float64
	LongSide            float64
	MaxWidth, MaxHeight float64
	Padding             Insets
}

type Quality struct {
	Flatness, Simplification, Softening, MinimumArea, Quantization float64
	Auto                                                           bool
}

type MarkerInput struct {
	ID       string
	Lon, Lat float64
}

type Input struct {
	Entity       catalog.Entity
	Geometry     catalog.Geometry
	Profile      string
	CorpusID     string
	Preset       string
	Layout       Layout
	Quality      Quality
	Override     *Override
	Markers      []MarkerInput
	MaxPathBytes int
}

type Bounds struct{ MinX, MinY, MaxX, MaxY float64 }

func (b Bounds) Width() float64  { return b.MaxX - b.MinX }
func (b Bounds) Height() float64 { return b.MaxY - b.MinY }

type Transform struct {
	Scale, TranslateX, TranslateY float64
	CenterLon, CenterLat          float64
	Rotation                      float64
}

type Command struct {
	Op     string
	Values []float64
}

type Component struct {
	SourcePolygon int
	Rings         int
	Area          float64
	Protected     []string
}

type Removal struct {
	SourcePolygon, SourceRing int
	Reason                    string
	Area                      float64
}

type Marker struct {
	ID      string
	X, Y    float64
	Anomaly string
}

type Metrics struct {
	InputPoints, ProjectedPoints, OutputPoints, PathBytes int
}

type Result struct {
	SchemaVersion                                       int
	AlgorithmVersion, CorpusID, Entity, Profile, Preset string
	LayoutMode                                          LayoutMode
	ViewBox                                             Bounds
	NaturalAspect                                       float64
	EffectiveScale                                      float64
	Transform                                           Transform
	Quality                                             Quality
	Commands                                            []Command
	Path                                                string
	Components                                          []Component
	Removals                                            []Removal
	Markers                                             []Marker
	Diagnostics                                         []Diagnostic
	Metrics                                             Metrics
	LOD                                                 LODProvenance
}
