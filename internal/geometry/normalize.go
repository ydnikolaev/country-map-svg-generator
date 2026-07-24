package geometry

import (
	"math"
	"sort"
)

func normalize(src MultiPolygon) (MultiPolygon, error) {
	out := make(MultiPolygon, 0, len(src))
	for pi, poly := range src {
		if len(poly) == 0 {
			return nil, fail(ErrInvalidSource, "", "geometry", "polygon %d has no rings", pi)
		}
		np := make(Polygon, 0, len(poly))
		for ri, ring := range poly {
			if len(ring) < 4 {
				return nil, fail(ErrInvalidSource, "", "geometry", "ring %d/%d has fewer than four points", pi, ri)
			}
			clean := make(Ring, 0, len(ring))
			for _, p := range ring {
				if !finite(p.X) || !finite(p.Y) {
					return nil, fail(ErrInvalidSource, "", "geometry", "non-finite coordinate")
				}
				if len(clean) == 0 || p != clean[len(clean)-1] {
					clean = append(clean, p)
				}
			}
			if clean[0] != clean[len(clean)-1] {
				clean = append(clean, clean[0])
			}
			if len(clean) < 4 || math.Abs(signedArea(clean)) < 1e-12 {
				return nil, fail(ErrInvalidSource, "", "geometry", "collapsed ring %d/%d", pi, ri)
			}
			wantPositive := ri == 0
			if (signedArea(clean) > 0) != wantPositive {
				reverseClosed(clean)
			}
			np = append(np, clean)
		}
		out = append(out, np)
	}
	sort.SliceStable(out, func(i, j int) bool { return math.Abs(signedArea(out[i][0])) > math.Abs(signedArea(out[j][0])) })
	return out, nil
}

func signedArea(r Ring) float64 {
	var a float64
	for i := 0; i < len(r)-1; i++ {
		a += r[i].X*r[i+1].Y - r[i+1].X*r[i].Y
	}
	return a / 2
}
func reverseClosed(r Ring) {
	for i, j := 0, len(r)-2; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	r[len(r)-1] = r[0]
}
func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
