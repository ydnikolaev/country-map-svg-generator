package catalog

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type featureCollection struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}

type feature struct {
	Type       string                     `json:"type"`
	Properties map[string]json.RawMessage `json:"properties"`
	Geometry   rawGeometry                `json:"geometry"`
}

type rawGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

func loadFeatures(path string) ([]feature, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, diagnostic(path, "-", "-", "source-readable", "%v", err)
	}
	var collection featureCollection
	if err := json.Unmarshal(b, &collection); err != nil {
		return nil, diagnostic(path, "-", "-", "geojson", "%v", err)
	}
	if collection.Type != "FeatureCollection" || len(collection.Features) == 0 {
		return nil, diagnostic(path, "-", "features", "geojson-feature-collection", "empty or invalid feature collection")
	}
	return collection.Features, nil
}

func propertyString(f feature, key string) string {
	var s string
	_ = json.Unmarshal(f.Properties[key], &s)
	return s
}

func propertyInt64(f feature, key string) int64 {
	var n int64
	_ = json.Unmarshal(f.Properties[key], &n)
	return n
}

func propertyFloat(f feature, key string) float64 {
	var n float64
	_ = json.Unmarshal(f.Properties[key], &n)
	return n
}

func geometryFromFeature(path, entity string, f feature) (MultiPolygon, error) {
	var out MultiPolygon
	switch f.Geometry.Type {
	case "Polygon":
		var polygon Polygon
		if err := json.Unmarshal(f.Geometry.Coordinates, &polygon); err != nil {
			return nil, diagnostic(path, entity, "geometry.coordinates", "geometry-decode", "%v", err)
		}
		out = MultiPolygon{polygon}
	case "MultiPolygon":
		if err := json.Unmarshal(f.Geometry.Coordinates, &out); err != nil {
			return nil, diagnostic(path, entity, "geometry.coordinates", "geometry-decode", "%v", err)
		}
	default:
		return nil, diagnostic(path, entity, "geometry.type", "geometry-polygonal", "got %q", f.Geometry.Type)
	}
	for pi := range out {
		if len(out[pi]) == 0 {
			return nil, diagnostic(path, entity, "geometry", "geometry-nonempty", "polygon %d has no rings", pi)
		}
		for ri := range out[pi] {
			if len(out[pi][ri]) < 4 {
				return nil, diagnostic(path, entity, "geometry", "geometry-ring", "ring has fewer than four positions")
			}
			for i, point := range out[pi][ri] {
				if !validCoordinate(point) {
					return nil, diagnostic(path, entity, fmt.Sprintf("geometry[%d][%d][%d]", pi, ri, i), "coordinate-range", "invalid longitude/latitude")
				}
				out[pi][ri][i][0] = round6(point[0])
				out[pi][ri][i][1] = round6(point[1])
			}
		}
	}
	return out, nil
}

func validCoordinate(p Point) bool {
	return !math.IsNaN(p[0]) && !math.IsNaN(p[1]) && !math.IsInf(p[0], 0) && !math.IsInf(p[1], 0) && p[0] >= -180 && p[0] <= 180 && p[1] >= -90 && p[1] <= 90
}

func round6(v float64) float64 { return math.Round(v*1e6) / 1e6 }
