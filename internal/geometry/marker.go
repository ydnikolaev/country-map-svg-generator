package geometry

import (
	"math"
	"sort"
)

func projectMarkers(inputs []MarkerInput, p projector, t Transform, vb Bounds, offsets map[string]Point) ([]Marker, error) {
	out := make([]Marker, 0, len(inputs))
	for _, m := range inputs {
		if m.ID == "" || !finite(m.Lon) || !finite(m.Lat) || m.Lat < -90 || m.Lat > 90 {
			return nil, fail(ErrMarker, "", "marker", "invalid marker %q", m.ID)
		}
		q, err := p.project(m.Lon, m.Lat)
		if err != nil {
			return nil, fail(ErrMarker, "", "marker", "%s: %v", m.ID, err)
		}
		x, y := q.X*t.Scale+t.TranslateX, q.Y*t.Scale+t.TranslateY
		if o, ok := offsets[m.ID]; ok {
			x += o.X
			y += o.Y
		}
		if !finite(x) || !finite(y) {
			return nil, fail(ErrMarker, "", "marker", "non-finite result for %q", m.ID)
		}
		anomaly := ""
		margin := math.Max(vb.Width(), vb.Height()) * .05
		if x < vb.MinX-margin || x > vb.MaxX+margin || y < vb.MinY-margin || y > vb.MaxY+margin {
			anomaly = "outside_fitted_result"
		}
		out = append(out, Marker{ID: m.ID, X: x, Y: y, Anomaly: anomaly})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
