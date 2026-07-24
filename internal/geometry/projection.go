package geometry

import (
	"math"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

type projector struct{ lon0, lat0, rotation float64 }

func LODProjectionContractV1() LODProjection {
	return LODProjection{
		Version: "v1", CoordinateSpace: "centered_laea",
		BuildRotation: 0, YAxis: "down", Flatness: .10, CoordinatePrecision: .000001,
	}
}

func centeredProjector(g MultiPolygon, rotation float64) projector {
	var x, y, z float64
	for _, p := range g {
		for _, r := range p {
			for i := 0; i < len(r)-1; i++ {
				lon, lat := rad(r[i].X), rad(r[i].Y)
				c := math.Cos(lat)
				x += c * math.Cos(lon)
				y += c * math.Sin(lon)
				z += math.Sin(lat)
			}
		}
	}
	return projector{lon0: math.Atan2(y, x), lat0: math.Atan2(z, math.Hypot(x, y)), rotation: rad(rotation)}
}

func (p projector) project(lonDeg, latDeg float64) (Point, error) {
	lon, lat := rad(lonDeg), rad(latDeg)
	dlon := wrap(lon - p.lon0)
	s0, c0 := math.Sin(p.lat0), math.Cos(p.lat0)
	s, c := math.Sin(lat), math.Cos(lat)
	den := 1 + s0*s + c0*c*math.Cos(dlon)
	if den <= 1e-14 {
		return Point{}, fail(ErrProjection, "", "coordinate", "antipodal projection singularity")
	}
	k := math.Sqrt(2 / den)
	x := k * c * math.Sin(dlon)
	y := k * (c0*s - s0*c*math.Cos(dlon))
	if p.rotation != 0 {
		cr, sr := math.Cos(p.rotation), math.Sin(p.rotation)
		x, y = x*cr-y*sr, x*sr+y*cr
	}
	return Point{X: x, Y: -y}, nil
}

func projectGeometry(g MultiPolygon, p projector, flatness float64) (MultiPolygon, int, error) {
	out := make(MultiPolygon, len(g))
	count := 0
	angular := math.Max(rad(.05), math.Min(rad(8), flatness/8))
	for pi, poly := range g {
		out[pi] = make(Polygon, len(poly))
		for ri, r := range poly {
			dst := make(Ring, 0, len(r)*2)
			for i := 0; i < len(r)-1; i++ {
				seg, err := subdivideGeo(r[i], r[i+1], angular, 0)
				if err != nil {
					return nil, 0, err
				}
				for j, q := range seg {
					if i > 0 && j == 0 {
						continue
					}
					v, err := p.project(q.X, q.Y)
					if err != nil {
						return nil, 0, err
					}
					dst = append(dst, v)
					count++
					if count > MaxPoints {
						return nil, 0, fail(ErrPointLimit, "", "geometry", "more than %d projected points", MaxPoints)
					}
				}
			}
			dst = append(dst, dst[0])
			count++
			out[pi][ri] = dst
		}
	}
	return out, count, nil
}

func subdivideGeo(a, b Point, maxAngle float64, depth int) ([]Point, error) {
	if depth > 16 {
		return []Point{a, b}, nil
	}
	va, vb := sphere(a), sphere(b)
	dot := clamp(va[0]*vb[0]+va[1]*vb[1]+va[2]*vb[2], -1, 1)
	angle := math.Acos(dot)
	if angle <= maxAngle {
		return []Point{a, b}, nil
	}
	m := Point{X: deg(math.Atan2(va[1]+vb[1], va[0]+vb[0])), Y: deg(math.Atan2(va[2]+vb[2], math.Hypot(va[0]+vb[0], va[1]+vb[1])))}
	l, _ := subdivideGeo(a, m, maxAngle, depth+1)
	r, _ := subdivideGeo(m, b, maxAngle, depth+1)
	return append(l[:len(l)-1], r...), nil
}
func sphere(p Point) [3]float64 {
	lon, lat := rad(p.X), rad(p.Y)
	c := math.Cos(lat)
	return [3]float64{c * math.Cos(lon), c * math.Sin(lon), math.Sin(lat)}
}

func ProjectLODGeometry(source catalog.Geometry, flatness, precision float64) (ProjectedLODGeometry, error) {
	contract := LODProjectionContractV1()
	if flatness != contract.Flatness || precision != contract.CoordinatePrecision {
		return ProjectedLODGeometry{}, fail(ErrProjection, "", "lod_projection_contract", "flatness/precision %g/%g do not match %s %g/%g", flatness, precision, contract.Version, contract.Flatness, contract.CoordinatePrecision)
	}
	normalized, err := normalize(fromCatalog(source.Coordinates))
	if err != nil {
		return ProjectedLODGeometry{}, err
	}
	prj := centeredProjector(normalized, 0)
	projected, _, err := projectGeometry(normalized, prj, flatness)
	if err != nil {
		return ProjectedLODGeometry{}, err
	}
	return ProjectedLODGeometry{
		Geometry: catalog.Geometry{ID: source.ID, Coordinates: toCatalogGeometry(projected)},
		Projection: LODProjection{
			Version: contract.Version, CoordinateSpace: contract.CoordinateSpace,
			CenterLongitude: deg(prj.lon0), CenterLatitude: deg(prj.lat0),
			BuildRotation: contract.BuildRotation, YAxis: contract.YAxis,
			Flatness: contract.Flatness, CoordinatePrecision: contract.CoordinatePrecision,
		},
	}, nil
}

func toCatalogGeometry(g MultiPolygon) catalog.MultiPolygon {
	out := make(catalog.MultiPolygon, len(g))
	for pi, polygon := range g {
		out[pi] = make(catalog.Polygon, len(polygon))
		for ri, ring := range polygon {
			out[pi][ri] = make(catalog.Ring, len(ring))
			for vi, point := range ring {
				out[pi][ri][vi] = catalog.Point{point.X, point.Y}
			}
		}
	}
	return out
}

func rotateProjectedLOD(g MultiPolygon, rotationDegrees float64) MultiPolygon {
	if rotationDegrees == 0 {
		return g
	}
	angle := rad(rotationDegrees)
	cosine, sine := math.Cos(angle), math.Sin(angle)
	out := make(MultiPolygon, len(g))
	for pi, polygon := range g {
		out[pi] = make(Polygon, len(polygon))
		for ri, ring := range polygon {
			out[pi][ri] = make(Ring, len(ring))
			for vi, point := range ring {
				out[pi][ri][vi] = Point{
					X: point.X*cosine + point.Y*sine,
					Y: -point.X*sine + point.Y*cosine,
				}
			}
		}
	}
	return out
}
func rad(v float64) float64 { return v * math.Pi / 180 }
func deg(v float64) float64 { return v * 180 / math.Pi }
func wrap(v float64) float64 {
	for v > math.Pi {
		v -= 2 * math.Pi
	}
	for v < -math.Pi {
		v += 2 * math.Pi
	}
	return v
}
