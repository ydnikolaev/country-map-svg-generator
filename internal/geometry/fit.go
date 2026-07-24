package geometry

import (
	"fmt"
	"math"
)

func geometryBounds(g MultiPolygon) Bounds {
	b := Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)}
	for _, p := range g {
		for _, r := range p {
			for _, q := range r {
				if q.X < b.MinX {
					b.MinX = q.X
				}
				if q.X > b.MaxX {
					b.MaxX = q.X
				}
				if q.Y < b.MinY {
					b.MinY = q.Y
				}
				if q.Y > b.MaxY {
					b.MaxY = q.Y
				}
			}
		}
	}
	return b
}

func validateLayout(l Layout) error {
	vals := []float64{l.Width, l.Height, l.LongSide, l.MaxWidth, l.MaxHeight, l.Padding.Top, l.Padding.Right, l.Padding.Bottom, l.Padding.Left}
	for _, v := range vals {
		if !finite(v) || v < 0 {
			return fail(ErrInvalidLayout, "", "layout", "dimensions and padding must be finite and nonnegative")
		}
	}
	switch l.Mode {
	case LayoutTight:
		if l.LongSide <= 0 && (l.MaxWidth <= 0 || l.MaxHeight <= 0) {
			return fail(ErrInvalidLayout, "", "layout", "tight requires long_side or positive maximum box")
		}
	case LayoutContain:
		if l.Width <= 0 || l.Height <= 0 {
			return fail(ErrInvalidLayout, "", "layout", "contain requires positive width and height")
		}
		if l.Width-l.Padding.Left-l.Padding.Right <= 0 || l.Height-l.Padding.Top-l.Padding.Bottom <= 0 {
			return fail(ErrInvalidLayout, "", "padding", "padding leaves no drawable area")
		}
	default:
		return fail(ErrInvalidLayout, "", "mode", "unknown layout mode %q", l.Mode)
	}
	return nil
}

func fitGeometry(g MultiPolygon, l Layout) (MultiPolygon, Bounds, Transform, float64, []Diagnostic, error) {
	if err := validateLayout(l); err != nil {
		return nil, Bounds{}, Transform{}, 0, nil, err
	}
	raw := geometryBounds(g)
	if raw.Width() <= 0 || raw.Height() <= 0 {
		return nil, Bounds{}, Transform{}, 0, nil, fail(ErrEmpty, "", "geometry", "projected bounds collapsed")
	}
	var scale, tx, ty float64
	var vb Bounds
	if l.Mode == LayoutTight {
		availW, availH := l.MaxWidth-l.Padding.Left-l.Padding.Right, l.MaxHeight-l.Padding.Top-l.Padding.Bottom
		if l.LongSide > 0 {
			scale = l.LongSide / math.Max(raw.Width(), raw.Height())
		} else {
			scale = math.Min(availW/raw.Width(), availH/raw.Height())
		}
		w, h := raw.Width()*scale+l.Padding.Left+l.Padding.Right, raw.Height()*scale+l.Padding.Top+l.Padding.Bottom
		vb = Bounds{MaxX: w, MaxY: h}
		tx = l.Padding.Left - raw.MinX*scale
		ty = l.Padding.Top - raw.MinY*scale
	} else {
		availW, availH := l.Width-l.Padding.Left-l.Padding.Right, l.Height-l.Padding.Top-l.Padding.Bottom
		scale = math.Min(availW/raw.Width(), availH/raw.Height())
		tx = l.Padding.Left + (availW-raw.Width()*scale)/2 - raw.MinX*scale
		ty = l.Padding.Top + (availH-raw.Height()*scale)/2 - raw.MinY*scale
		vb = Bounds{MaxX: l.Width, MaxY: l.Height}
	}
	if !finite(scale) || scale <= 0 {
		return nil, Bounds{}, Transform{}, 0, nil, fail(ErrInvalidLayout, "", "layout", "unsafe fit arithmetic")
	}
	out := transformGeometry(g, scale, tx, ty)
	d := math.Max(raw.Width()*scale, raw.Height()*scale)
	diags := []Diagnostic{}
	if d < 32 {
		diags = append(diags, Diagnostic{"low_resolution", "warning", fmt.Sprintf("fitted geometry long side %.2fpx is below 32px", d)})
	}
	if d > 2048 {
		diags = append(diags, Diagnostic{"generation_cost", "warning", fmt.Sprintf("fitted geometry long side %.2fpx exceeds 2048px", d)})
	}
	return out, vb, Transform{Scale: scale, TranslateX: tx, TranslateY: ty}, d, diags, nil
}

func transformGeometry(g MultiPolygon, s, tx, ty float64) MultiPolygon {
	out := make(MultiPolygon, len(g))
	for i, p := range g {
		out[i] = make(Polygon, len(p))
		for j, r := range p {
			out[i][j] = make(Ring, len(r))
			for k, q := range r {
				out[i][j][k] = Point{q.X*s + tx, q.Y*s + ty}
			}
		}
	}
	return out
}
