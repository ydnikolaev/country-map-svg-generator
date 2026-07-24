package geometry

import (
	"sort"

	sf "github.com/peterstace/simplefeatures/geom"
)

func simplify(g MultiPolygon, tolerance float64, protected map[int][]string) (MultiPolygon, []Removal, error) {
	return simplifyShared(g, tolerance, protected, nil)
}
func simplifyShared(g MultiPolygon, tolerance float64, protected map[int][]string, shared map[Point]bool) (MultiPolygon, []Removal, error) {
	out, removals, _, err := simplifySharedResolved(g, tolerance, protected, shared)
	return out, removals, err
}
func simplifySharedResolved(g MultiPolygon, tolerance float64, protected map[int][]string, shared map[Point]bool) (MultiPolygon, []Removal, float64, error) {
	_ = protected // component collapse is rejected below, including protected components.
	if tolerance <= 0 {
		return g, nil, 0, validateTopologyShared(g, shared)
	}
	for t := tolerance; t > 0.0001; t /= 2 {
		out := make(MultiPolygon, len(g))
		removed := []Removal{}
		ok := true
		for i, p := range g {
			rings := make([]sf.LineString, len(p))
			for j, r := range p {
				f := make([]float64, 0, len(r)*2)
				for _, q := range r {
					f = append(f, q.X, q.Y)
				}
				rings[j] = sf.NewLineString(sf.NewSequence(f, sf.DimXY))
			}
			poly := sf.NewPolygon(rings)
			var simplified sf.Polygon
			var err error
			for local := t; local > 0.0001; local /= 2 {
				simplified, err = poly.Simplify(local)
				if err == nil && !simplified.IsEmpty() {
					break
				}
			}
			if err != nil || simplified.IsEmpty() {
				candidate, candidateOK := simplifyRings(p, t)
				if candidateOK {
					out[i] = candidate
				} else {
					out[i] = p
				}
				continue
			}
			out[i] = make(Polygon, simplified.NumRings())
			simpleRings := simplified.DumpRings()
			for j, sr := range simpleRings {
				seq := sr.Coordinates()
				nr := make(Ring, seq.Length())
				for k := range nr {
					x := seq.GetXY(k)
					nr[k] = Point{x.X, x.Y}
				}
				out[i][j] = nr
			}
			for j := simplified.NumRings(); j < len(p); j++ {
				removed = append(removed, Removal{SourcePolygon: i, SourceRing: j, Reason: "collapsed_during_simplification", Area: mathAbs(signedArea(p[j]))})
			}
		}
		if ok && validateTopologyShared(out, shared) == nil {
			return out, removed, t, nil
		}
	}
	if validateTopologyShared(g, shared) == nil {
		return g, nil, 0, nil
	}
	return nil, nil, 0, fail(ErrTopology, "", "simplification", "no validated tolerance remained")
}

func simplifyRings(p Polygon, tolerance float64) (Polygon, bool) {
	for local := tolerance; local > 0.0001; local /= 2 {
		out := make(Polygon, len(p))
		ok := true
		for i, r := range p {
			f := make([]float64, 0, len(r)*2)
			for _, q := range r {
				f = append(f, q.X, q.Y)
			}
			seq := sf.NewLineString(sf.NewSequence(f, sf.DimXY)).Simplify(local).Coordinates()
			if seq.Length() < 4 {
				ok = false
				break
			}
			out[i] = make(Ring, seq.Length())
			for j := range out[i] {
				q := seq.GetXY(j)
				out[i][j] = Point{q.X, q.Y}
			}
		}
		if ok && validateTopology(MultiPolygon{out}) == nil {
			return out, true
		}
	}
	return nil, false
}

func validateTopology(g MultiPolygon) error {
	return validateTopologyShared(g, nil)
}
func validateTopologyShared(g MultiPolygon, shared map[Point]bool) error {
	for pi, p := range g {
		if len(p) == 0 {
			return fail(ErrTopology, "", "polygon", "polygon %d vanished", pi)
		}
		for ri, r := range p {
			if len(r) < 4 || r[0] != r[len(r)-1] || mathAbs(signedArea(r)) < 1e-8 {
				return fail(ErrTopology, "", "ring", "ring %d/%d collapsed", pi, ri)
			}
			if i, j, _, found := ringIntersectionSweep(r, shared); found {
				return fail(ErrTopology, "", "ring", "ring %d/%d self-intersects segments %d/%d (%g,%g)-(%g,%g) and (%g,%g)-(%g,%g)", pi, ri, i, j, r[i].X, r[i].Y, r[i+1].X, r[i+1].Y, r[j].X, r[j].Y, r[j+1].X, r[j+1].Y)
			}
			if ri > 0 && !pointInRing(r[0], p[0]) {
				return fail(ErrTopology, "", "hole", "hole %d/%d escaped exterior", pi, ri)
			}
		}
	}
	return nil
}

type topologySegment struct {
	index                  int
	a, b                   Point
	minX, minY, maxX, maxY float64
}

func ringIntersectionSweep(r Ring, shared map[Point]bool) (firstI, firstJ, candidates int, found bool) {
	segments := make([]topologySegment, len(r)-1)
	for i := range segments {
		a, b := r[i], r[i+1]
		segments[i] = topologySegment{index: i, a: a, b: b, minX: min(a.X, b.X), minY: min(a.Y, b.Y), maxX: max(a.X, b.X), maxY: max(a.Y, b.Y)}
	}
	sort.SliceStable(segments, func(i, j int) bool {
		if segments[i].minX == segments[j].minX {
			return segments[i].index < segments[j].index
		}
		return segments[i].minX < segments[j].minX
	})
	active := make([]topologySegment, 0, len(segments))
	firstI, firstJ = len(segments), len(segments)
	for _, current := range segments {
		next := active[:0]
		for _, other := range active {
			if other.maxX < current.minX {
				continue
			}
			next = append(next, other)
			if other.maxY < current.minY || current.maxY < other.minY {
				continue
			}
			i, j := other.index, current.index
			if i > j {
				i, j = j, i
			}
			if j == i+1 || (i == 0 && j == len(r)-2) {
				continue
			}
			candidates++
			if segmentsCrossShared(r[i], r[i+1], r[j], r[j+1], shared) && (i < firstI || i == firstI && j < firstJ) {
				firstI, firstJ, found = i, j, true
			}
		}
		active = append(next, current)
	}
	return firstI, firstJ, candidates, found
}

func segmentsCross(a, b, c, d Point) bool {
	return segmentsCrossShared(a, b, c, d, nil)
}
func segmentsCrossShared(a, b, c, d Point, shared map[Point]bool) bool {
	o1 := orient(a, b, c)
	o2 := orient(a, b, d)
	o3 := orient(c, d, a)
	o4 := orient(c, d, b)
	if !(o1*o2 < 0 && o3*o4 < 0) {
		return false
	}
	if shared != nil {
		for _, pair := range [][2]Point{{a, c}, {a, d}, {b, c}, {b, d}} {
			if (pair[0] == pair[1] && shared[pair[0]]) || (distance(pair[0], pair[1]) <= 1e-12 && shared[Point{quant(pair[0].X, 1e-12), quant(pair[0].Y, 1e-12)}]) {
				return false
			}
		}
	}
	return true
}
func orient(a, b, c Point) float64 { return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X) }
func mathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
