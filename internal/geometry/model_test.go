package geometry

import (
	"errors"
	"math"
	"testing"
)

func TestAutoQualityUsesFittedGeometryScale(t *testing.T) {
	a, b := AutoQuality(128), AutoQuality(128)
	if a != b {
		t.Fatal("equal fitted scale must produce equal quality")
	}
	if a == AutoQuality(64) {
		t.Fatal("changed fitted scale must change quality")
	}
	if a.Flatness != .225 || a.Simplification != 1.536 || a.Softening != .3075 || a.MinimumArea != 1.375 || a.Quantization != .01 {
		t.Fatalf("unexpected formulas: %+v", a)
	}
}

func TestLayoutValidationArbitraryFrames(t *testing.T) {
	g := MultiPolygon{{{{0, 0}, {2, 0}, {2, 1}, {0, 1}, {0, 0}}}}
	for _, l := range []Layout{
		{Mode: LayoutContain, Width: 19, Height: 19},
		{Mode: LayoutContain, Width: 10_000, Height: 12_345},
		{Mode: LayoutContain, Width: 9, Height: 100},
		{Mode: LayoutTight, LongSide: 128, Padding: Insets{8, 8, 8, 8}},
	} {
		_, _, tr, _, _, err := fitGeometry(g, l)
		if err != nil {
			t.Fatalf("%+v: %v", l, err)
		}
		if tr.Scale <= 0 || !finite(tr.Scale) {
			t.Fatalf("invalid uniform scale: %+v", tr)
		}
	}
	_, _, _, _, _, err := fitGeometry(g, Layout{Mode: LayoutContain, Width: 10, Height: 10, Padding: Insets{6, 6, 6, 6}})
	var pe *PipelineError
	if !errors.As(err, &pe) || pe.Code != ErrInvalidLayout {
		t.Fatalf("expected typed layout error, got %v", err)
	}
}

func TestTightNaturalAspectAndContainNeverStretch(t *testing.T) {
	g := MultiPolygon{{{{-4, -1}, {4, -1}, {4, 1}, {-4, 1}, {-4, -1}}}}
	_, tight, tr, d, _, err := fitGeometry(g, Layout{Mode: LayoutTight, LongSide: 128})
	if err != nil {
		t.Fatal(err)
	}
	if got := tight.Width() / tight.Height(); math.Abs(got-4) > 1e-12 {
		t.Fatalf("natural ratio=%g", got)
	}
	if d != 128 || tr.Scale != 16 {
		t.Fatalf("scale=%g d=%g", tr.Scale, d)
	}
	_, frame, tr2, d2, _, err := fitGeometry(g, Layout{Mode: LayoutContain, Width: 200, Height: 200})
	if err != nil {
		t.Fatal(err)
	}
	if frame.Width() != 200 || frame.Height() != 200 || d2 != 200 || tr2.Scale != 25 {
		t.Fatalf("bad contain: %+v %+v %g", frame, tr2, d2)
	}
}
