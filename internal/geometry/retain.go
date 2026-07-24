package geometry

import (
	"math"
	"sort"
)

func retain(g MultiPolygon, protected map[int][]string, minimumParts int, threshold float64) (MultiPolygon, []Removal, []Component, error) {
	if len(g) == 0 {
		return nil, nil, nil, fail(ErrEmpty, "", "geometry", "no components")
	}
	type row struct {
		i    int
		area float64
	}
	rows := make([]row, len(g))
	for i, p := range g {
		rows[i] = row{i, polygonArea(p)}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].area == rows[j].area {
			return rows[i].i < rows[j].i
		}
		return rows[i].area > rows[j].area
	})
	keep := map[int]bool{rows[0].i: true}
	for i, n := range protected {
		if len(n) > 0 {
			keep[i] = true
		}
	}
	for i := 0; i < minimumParts && i < len(rows); i++ {
		keep[rows[i].i] = true
	}
	out := MultiPolygon{}
	removals := []Removal{}
	components := []Component{}
	for i, p := range g {
		a := polygonArea(p)
		if !keep[i] && a < threshold {
			removals = append(removals, Removal{SourcePolygon: i, SourceRing: -1, Reason: "below_visible_area", Area: a})
			continue
		}
		out = append(out, p)
		components = append(components, Component{SourcePolygon: i, Rings: len(p), Area: a, Protected: append([]string(nil), protected[i]...)})
	}
	if len(out) == 0 {
		return nil, nil, nil, fail(ErrEmpty, "", "geometry", "retention removed all components")
	}
	return out, removals, components, nil
}
func polygonArea(p Polygon) float64 {
	if len(p) == 0 {
		return 0
	}
	a := math.Abs(signedArea(p[0]))
	for _, h := range p[1:] {
		a -= math.Abs(signedArea(h))
	}
	return math.Max(0, a)
}
func pointInRing(q Point, r Ring) bool {
	inside := false
	for i, j := 0, len(r)-1; i < len(r); j, i = i, i+1 {
		a, b := r[i], r[j]
		if ((a.Y > q.Y) != (b.Y > q.Y)) && q.X < (b.X-a.X)*(q.Y-a.Y)/(b.Y-a.Y)+a.X {
			inside = !inside
		}
	}
	return inside
}
