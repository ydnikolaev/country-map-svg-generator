package geometry

import "testing"

func TestMarkerUsesExactGeometryTransform(t *testing.T) {
	p := projector{}
	got, err := projectMarkers([]MarkerInput{{ID: "x", Lon: 0, Lat: 0}}, p, Transform{Scale: 10, TranslateX: 5, TranslateY: 7}, Bounds{MaxX: 20, MaxY: 20}, nil)
	if err != nil {
		t.Fatal(err)
	}
	q, _ := p.project(0, 0)
	if got[0].X != q.X*10+5 || got[0].Y != q.Y*10+7 {
		t.Fatal("divergent marker transform")
	}
}
