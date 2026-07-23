package catalog

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGeoJSONPolygonNormalizesToMultiPolygon(t *testing.T) {
	f := feature{Geometry: rawGeometry{Type: "Polygon", Coordinates: json.RawMessage(`[[[1.1234567,2],[3,2],[3,4],[1.1234567,2]]]`)}}
	geometry, err := geometryFromFeature("fixture.geojson", "AA", f)
	if err != nil {
		t.Fatal(err)
	}
	if len(geometry) != 1 || geometry[0][0][0][0] != 1.123457 {
		t.Fatalf("unexpected normalized geometry: %#v", geometry)
	}
}

func TestGeoJSONRejectsEmptyAndInvalidGeometry(t *testing.T) {
	for name, f := range map[string]feature{
		"type":  {Geometry: rawGeometry{Type: "Point", Coordinates: json.RawMessage(`[0,0]`)}},
		"range": {Geometry: rawGeometry{Type: "Polygon", Coordinates: json.RawMessage(`[[[181,0],[0,0],[0,1],[181,0]]]`)}},
		"ring":  {Geometry: rawGeometry{Type: "Polygon", Coordinates: json.RawMessage(`[[[0,0],[1,0],[0,0]]]`)}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := geometryFromFeature("fixture.geojson", "AA", f)
			if err == nil || !strings.Contains(err.Error(), "invariant=") {
				t.Fatalf("expected typed geometry diagnostic, got %v", err)
			}
		})
	}
}
