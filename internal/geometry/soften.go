package geometry

import "math"

// softRing stores a bounded quadratic representation while retaining the
// validated polygon as the topology authority.
func softRing(r Ring, tolerance float64) []Command {
	if len(r) < 4 || tolerance <= 0 {
		return linearCommands(r)
	}
	cmd := []Command{{Op: "M", Values: []float64{r[0].X, r[0].Y}}}
	for i := 1; i < len(r)-1; i++ {
		prev, cur, next := r[i-1], r[i], r[i+1]
		a := distance(prev, cur)
		b := distance(cur, next)
		cut := clamp(tolerance, 0, mathMin(a, b)/4)
		if cut < 1e-6 || mathAbs(orient(prev, cur, next)) < 1e-8 {
			cmd = append(cmd, Command{Op: "L", Values: []float64{cur.X, cur.Y}})
			continue
		}
		entry := toward(cur, prev, cut/a)
		exit := toward(cur, next, cut/b)
		cmd = append(cmd, Command{Op: "L", Values: []float64{entry.X, entry.Y}}, Command{Op: "Q", Values: []float64{cur.X, cur.Y, exit.X, exit.Y}})
	}
	cmd = append(cmd, Command{Op: "Z"})
	return cmd
}
func linearCommands(r Ring) []Command {
	if len(r) == 0 {
		return nil
	}
	c := []Command{{Op: "M", Values: []float64{r[0].X, r[0].Y}}}
	for _, p := range r[1 : len(r)-1] {
		c = append(c, Command{Op: "L", Values: []float64{p.X, p.Y}})
	}
	return append(c, Command{Op: "Z"})
}
func distance(a, b Point) float64        { x, y := a.X-b.X, a.Y-b.Y; return sqrt(x*x + y*y) }
func toward(a, b Point, t float64) Point { return Point{a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t} }
func mathMin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func validateSoftenedGeometry(canonical MultiPolygon, tolerance float64) (float64, error) {
	sampled := make(MultiPolygon, len(canonical))
	for pi, polygon := range canonical {
		sampled[pi] = make(Polygon, len(polygon))
		for ri, ring := range polygon {
			sampled[pi][ri] = sampleSoftRing(ring, tolerance, .05)
		}
	}
	if err := validateTopology(sampled); err != nil {
		return 0, err
	}
	return matchedBoundaryDeviation(canonical, sampled, tolerance), nil
}

func sampleSoftRing(r Ring, tolerance, maximumStep float64) Ring {
	commands := softRing(r, tolerance)
	var sampled Ring
	var current, start Point
	for _, command := range commands {
		switch command.Op {
		case "M":
			current = Point{command.Values[0], command.Values[1]}
			start = current
			sampled = append(sampled, current)
		case "L":
			current = Point{command.Values[0], command.Values[1]}
			sampled = append(sampled, current)
		case "Q":
			control := Point{command.Values[0], command.Values[1]}
			end := Point{command.Values[2], command.Values[3]}
			steps := int(math.Ceil(2 * math.Max(distance(current, control), distance(control, end)) / maximumStep))
			if steps < 1 {
				steps = 1
			}
			origin := current
			for step := 1; step <= steps; step++ {
				t := float64(step) / float64(steps)
				u := 1 - t
				sampled = append(sampled, Point{
					X: u*u*origin.X + 2*u*t*control.X + t*t*end.X,
					Y: u*u*origin.Y + 2*u*t*control.Y + t*t*end.Y,
				})
			}
			current = end
		case "Z":
			if len(sampled) > 0 && sampled[len(sampled)-1] != start {
				sampled = append(sampled, start)
			}
		}
	}
	return sampled
}
