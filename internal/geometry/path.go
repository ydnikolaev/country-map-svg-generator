package geometry

import (
	"math"
	"sort"
	"strconv"
)

func canonicalize(g MultiPolygon, q float64) (MultiPolygon, error) {
	return canonicalizeShared(g, q, nil)
}
func canonicalizeShared(g MultiPolygon, q float64, shared map[Point]bool) (MultiPolygon, error) {
	if q <= 0 || !finite(q) {
		return nil, fail(ErrInvalidLayout, "", "quantization", "must be finite and positive")
	}
	out := make(MultiPolygon, len(g))
	for pi, p := range g {
		np, err := quantizePolygon(p, q)
		if err != nil {
			return nil, err
		}
		out[pi] = np
		sort.SliceStable(out[pi][1:], func(i, j int) bool { return ringLess(out[pi][i+1], out[pi][j+1]) })
	}
	sort.SliceStable(out, func(i, j int) bool { return ringLess(out[i][0], out[j][0]) })
	if err := validateTopologyShared(out, shared); err != nil {
		return nil, err
	}
	return out, nil
}
func sharedSourceVertices(g MultiPolygon, q float64) map[Point]bool {
	counts := map[Point]int{}
	for _, p := range g {
		for _, r := range p {
			for _, v := range r[:len(r)-1] {
				key := Point{quant(v.X, q), quant(v.Y, q)}
				counts[key]++
			}
		}
	}
	out := map[Point]bool{}
	for key, n := range counts {
		if n > 1 {
			out[key] = true
		}
	}
	return out
}
func quantizePolygon(p Polygon, q float64) (Polygon, error) {
	out := make(Polygon, len(p))
	for ri, r := range p {
		n := make(Ring, len(r))
		for i, v := range r {
			n[i] = Point{quant(v.X, q), quant(v.Y, q)}
		}
		n = dedupe(n)
		if len(n) < 4 {
			return nil, fail(ErrTopology, "", "quantization", "ring collapsed")
		}
		if n[0] != n[len(n)-1] {
			n = append(n, n[0])
		}
		if (signedArea(n) > 0) != (ri == 0) {
			reverseClosed(n)
		}
		n = rotateMin(n)
		out[ri] = n
	}
	return out, nil
}
func quant(v, q float64) float64 {
	x := math.Round(v/q) * q
	if math.Abs(x) < q/2 {
		return 0
	}
	return x
}
func dedupe(r Ring) Ring {
	n := make(Ring, 0, len(r))
	for _, p := range r {
		if len(n) == 0 || p != n[len(n)-1] {
			n = append(n, p)
		}
	}
	return n
}
func rotateMin(r Ring) Ring {
	n := len(r) - 1
	at := 0
	for i := 1; i < n; i++ {
		if r[i].X < r[at].X || (r[i].X == r[at].X && r[i].Y < r[at].Y) {
			at = i
		}
	}
	out := make(Ring, 0, len(r))
	out = append(out, r[at:n]...)
	out = append(out, r[:at]...)
	out = append(out, out[0])
	return out
}
func ringLess(a, b Ring) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i].X != b[i].X {
			return a[i].X < b[i].X
		}
		if a[i].Y != b[i].Y {
			return a[i].Y < b[i].Y
		}
	}
	return len(a) < len(b)
}

func commandsAndPath(g MultiPolygon, softening float64) ([]Command, string, error) {
	cmds := []Command{}
	for _, p := range g {
		for _, r := range p {
			for _, c := range softRing(r, softening) {
				cmds = append(cmds, c)
			}
		}
	}
	path, err := SerializeCommands(cmds)
	if err != nil {
		return cmds, "", fail(ErrSerialization, "", "path", "%v", err)
	}
	if path == "" {
		return cmds, "", fail(ErrSerialization, "", "path", "serializer returned an empty path")
	}
	return cmds, path, nil
}
func number(v float64) string {
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}
