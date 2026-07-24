package geometry

import (
	"sort"
	"sync"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

// The published table is the committed DEC-006/DEC-009 ladder, resolved on
// first use. Tests inject their own tables through GenerateWithLOD and never
// mutate this seam.
//
// Resolution is lazy for the same reason the ladder parse is: it decodes ~10 MB.
// A caller that never generates never pays for it, and a failure to resolve is
// returned rather than panicking at init — a build shipping a ladder that does
// not bind its own corpus and oracle must fail loudly at the call, not silently
// fall back to serving source geometry.
var (
	publishedOnce  sync.Once
	publishedTable *LODTable
	publishedErr   error
)

func publishedLODTable() (*LODTable, error) {
	publishedOnce.Do(func() {
		ladder, err := EmbeddedLadderTable()
		if err != nil {
			publishedErr = err
			return
		}
		oracle, err := EmbeddedSilhouetteOracle()
		if err != nil {
			publishedErr = err
			return
		}
		publishedTable, publishedErr = ladder.RuntimeTable(oracle)
	})
	return publishedTable, publishedErr
}

func Generate(raw Input) (Result, error) {
	in, err := validatePublicInput(raw)
	if err != nil {
		return Result{}, err
	}
	table, err := publishedLODTable()
	if err != nil {
		return Result{}, err
	}
	return GenerateWithLOD(in, table)
}

func validatePublicInput(raw Input) (Input, error) {
	in, err := ApplyPreset(raw)
	if err != nil {
		return Input{}, err
	}
	if in.Entity.Alpha2 == "" || len(in.Geometry.Coordinates) == 0 {
		return Input{}, fail(ErrInvalidSource, in.Entity.Alpha2, "input", "entity and geometry are required")
	}
	if in.Geometry.ID != "" {
		expected := in.Entity.Profiles.UN.GeometryID
		if in.Profile == "de_facto" {
			expected = in.Entity.Profiles.DeFacto.GeometryID
			if expected == "" && in.Entity.Profiles.DeFacto.IdenticalTo == "un" {
				expected = in.Entity.Profiles.UN.GeometryID
			}
		}
		if expected != "" && in.Geometry.ID != expected {
			return Input{}, fail(ErrInvalidSource, in.Entity.Alpha2, "geometry", "profile references %q, got %q", expected, in.Geometry.ID)
		}
	}
	if in.Override != nil && in.Override.ISO != in.Entity.Alpha2 {
		return Input{}, fail(ErrOverride, in.Entity.Alpha2, "iso", "override targets %q", in.Override.ISO)
	}
	return in, nil
}

func fromCatalog(g catalog.MultiPolygon) MultiPolygon {
	out := make(MultiPolygon, len(g))
	for i, p := range g {
		out[i] = make(Polygon, len(p))
		for j, r := range p {
			out[i][j] = make(Ring, len(r))
			for k, q := range r {
				out[i][j][k] = Point{q[0], q[1]}
			}
		}
	}
	return out
}

func countPoints(g MultiPolygon) int {
	n := 0
	for _, p := range g {
		for _, r := range p {
			n += len(r)
		}
	}
	return n
}

func validateQuality(q Quality) error {
	vals := []float64{q.Flatness, q.Simplification, q.Softening, q.MinimumArea, q.Quantization}
	for _, v := range vals {
		if !finite(v) || v < 0 {
			return fail(ErrInvalidLayout, "", "quality", "quality values must be finite and nonnegative")
		}
	}
	if q.Quantization <= 0 || q.Softening > .35 || (!q.Auto && q.Simplification > 1.25) || q.Flatness > .30 || q.MinimumArea > 1.50 {
		return fail(ErrInvalidLayout, "", "quality", "quality exceeds bounded policy")
	}
	return nil
}

func protectedComponents(e catalog.Entity, p projector, t Transform, g MultiPolygon) (map[int][]string, error) {
	out := map[int][]string{}
	for _, f := range e.Protected {
		q, err := p.project(f.Anchor[0], f.Anchor[1])
		if err != nil {
			return nil, fail(ErrProtected, e.Alpha2, "protected", "%s", f.Name)
		}
		q = Point{q.X*t.Scale + t.TranslateX, q.Y*t.Scale + t.TranslateY}
		found := -1
		for i, poly := range g {
			if pointInRing(q, poly[0]) {
				found = i
				break
			}
		}
		if found < 0 {
			return nil, fail(ErrProtected, e.Alpha2, "protected", "%s anchor is outside geometry", f.Name)
		}
		out[found] = append(out[found], f.Name)
	}
	for i := range out {
		sort.Strings(out[i])
	}
	return out, nil
}
