package geometry

import "math"

const gridPhaseSteps = 10

type GridPhase struct {
	X, Y float64
}

func gridPhaseSchedule() []GridPhase {
	phases := make([]GridPhase, 0, gridPhaseSteps*gridPhaseSteps)
	for x := 0; x < gridPhaseSteps; x++ {
		for y := 0; y < gridPhaseSteps; y++ {
			phases = append(phases, GridPhase{X: float64(x) / 1000, Y: float64(y) / 1000})
		}
	}
	return phases
}

func composeGridPhase(base Transform, phase GridPhase) Transform {
	out := base
	out.TranslateX += phase.X
	out.TranslateY += phase.Y
	return out
}

func shiftViewBox(base Bounds, phase GridPhase) Bounds {
	return Bounds{
		MinX: base.MinX + phase.X,
		MinY: base.MinY + phase.Y,
		MaxX: base.MaxX + phase.X,
		MaxY: base.MaxY + phase.Y,
	}
}

func unfitGeometry(fitted MultiPolygon, base Transform) MultiPolygon {
	out := make(MultiPolygon, len(fitted))
	for pi, polygon := range fitted {
		out[pi] = make(Polygon, len(polygon))
		for ri, ring := range polygon {
			out[pi][ri] = make(Ring, len(ring))
			for vi, point := range ring {
				out[pi][ri][vi] = Point{
					X: (point.X - base.TranslateX) / base.Scale,
					Y: (point.Y - base.TranslateY) / base.Scale,
				}
			}
		}
	}
	return out
}

func materializeGridPhase(unfitted MultiPolygon, transform Transform) MultiPolygon {
	return transformGeometry(unfitted, transform.Scale, transform.TranslateX, transform.TranslateY)
}

func geometryOnGrid(g MultiPolygon, q float64) bool {
	for _, polygon := range g {
		for _, ring := range polygon {
			for _, point := range ring {
				if math.Abs(point.X/q-math.Round(point.X/q)) > 1e-9 ||
					math.Abs(point.Y/q-math.Round(point.Y/q)) > 1e-9 {
					return false
				}
			}
		}
	}
	return true
}

func geometryContainedBy(g MultiPolygon, bounds Bounds) bool {
	const epsilon = 1e-9
	for _, polygon := range g {
		for _, ring := range polygon {
			for _, point := range ring {
				if point.X < bounds.MinX-epsilon || point.X > bounds.MaxX+epsilon ||
					point.Y < bounds.MinY-epsilon || point.Y > bounds.MaxY+epsilon {
					return false
				}
			}
		}
	}
	return true
}

func validatePhaseProtection(in Input, prj projector, transform Transform, canonical MultiPolygon, minimumParts int) error {
	if len(canonical) < minimumParts {
		return fail(ErrProtected, in.Entity.Alpha2, "minimum_parts", "got %d want at least %d", len(canonical), minimumParts)
	}
	for _, feature := range in.Entity.Protected {
		point, err := prj.project(feature.Anchor[0], feature.Anchor[1])
		if err != nil {
			return fail(ErrProtected, in.Entity.Alpha2, "protected", "%s: %v", feature.Name, err)
		}
		point = Point{
			X: point.X*transform.Scale + transform.TranslateX,
			Y: point.Y*transform.Scale + transform.TranslateY,
		}
		covered := false
		for _, polygon := range canonical {
			if pointInRing(point, polygon[0]) {
				covered = true
				break
			}
		}
		if !covered {
			return fail(ErrProtected, in.Entity.Alpha2, "protected", "%s anchor is outside final geometry", feature.Name)
		}
	}
	return nil
}
