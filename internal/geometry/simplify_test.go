package geometry

import (
	"testing"
	"time"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

func TestTopologyAQSourceCanonicalProbe(t *testing.T) {
	r := aqSourceCanonicalRing(t)
	for _, segments := range []int{500, 1000} {
		if len(r)-1 < segments {
			t.Fatalf("AQ ring has %d segments, need %d", len(r)-1, segments)
		}
		probe := append(Ring(nil), r[:segments]...)
		probe = append(probe, probe[0])
		start := time.Now()
		candidates, crossings := bruteRingCrossings(probe, nil)
		t.Logf("segments=%d brute_candidates=%d crossings=%d duration=%s", segments, candidates, crossings, time.Since(start))
		want := segments * (segments - 3) / 2
		if candidates != want {
			t.Fatalf("segments=%d candidates=%d want=%d", segments, candidates, want)
		}
		sweepStart := time.Now()
		_, _, sweepCandidates, _ := ringIntersectionSweep(probe, nil)
		t.Logf("segments=%d sweep_candidates=%d duration=%s", segments, sweepCandidates, time.Since(sweepStart))
		if sweepCandidates >= candidates {
			t.Fatalf("segments=%d sweep candidates=%d did not reduce brute candidates=%d", segments, sweepCandidates, candidates)
		}
	}
}

func TestTopologySweepMatchesBrute(t *testing.T) {
	rings := []Ring{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
		{{0, 0}, {10, 10}, {0, 10}, {10, 0}, {0, 0}},
		{{0, 0}, {5, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	}
	aq := aqSourceCanonicalRing(t)
	probe := append(Ring(nil), aq[:1000]...)
	probe = append(probe, probe[0])
	rings = append(rings, probe)
	for index, r := range rings {
		wantI, wantJ, wantFound := bruteFirstIntersection(r, nil)
		gotI, gotJ, _, gotFound := ringIntersectionSweep(r, nil)
		if gotFound != wantFound || gotI != wantI || gotJ != wantJ {
			t.Fatalf("ring %d sweep=(%d,%d,%v) brute=(%d,%d,%v)", index, gotI, gotJ, gotFound, wantI, wantJ, wantFound)
		}
	}
}

func TestTopologyAQSourceValidationTiming(t *testing.T) {
	g := aqSourceFittedGeometry(t)
	start := time.Now()
	err := validateTopologyShared(g, sharedSourceVertices(g, 1e-12))
	t.Logf("points=%d duration=%s err=%v", countPoints(g), time.Since(start), err)
	if err != nil {
		t.Fatal(err)
	}
}

func bruteRingCrossings(r Ring, shared map[Point]bool) (candidates, crossings int) {
	for i := 0; i < len(r)-1; i++ {
		for j := i + 1; j < len(r)-1; j++ {
			if j == i+1 || (i == 0 && j == len(r)-2) {
				continue
			}
			candidates++
			if segmentsCrossShared(r[i], r[i+1], r[j], r[j+1], shared) {
				crossings++
			}
		}
	}
	return candidates, crossings
}

func bruteFirstIntersection(r Ring, shared map[Point]bool) (firstI, firstJ int, found bool) {
	for i := 0; i < len(r)-1; i++ {
		for j := i + 1; j < len(r)-1; j++ {
			if j == i+1 || (i == 0 && j == len(r)-2) {
				continue
			}
			if segmentsCrossShared(r[i], r[i+1], r[j], r[j+1], shared) {
				return i, j, true
			}
		}
	}
	return len(r) - 1, len(r) - 1, false
}

func aqSourceCanonicalRing(t *testing.T) Ring {
	t.Helper()
	fitted := aqSourceFittedGeometry(t)
	var largest Ring
	for _, polygon := range fitted {
		for _, ring := range polygon {
			if len(ring) > len(largest) {
				largest = ring
			}
		}
	}
	canonical, err := quantizePolygon(Polygon{largest}, .01)
	if err != nil {
		t.Fatal(err)
	}
	return canonical[0]
}

func aqSourceFittedGeometry(t *testing.T) MultiPolygon {
	t.Helper()
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	in, err := InputFromCatalog(c, "AQ", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	in, err = ApplyPreset(in)
	if err != nil {
		t.Fatal(err)
	}
	source, err := normalize(fromCatalog(in.Geometry.Coordinates))
	if err != nil {
		t.Fatal(err)
	}
	projected, _, err := projectGeometry(source, centeredProjector(source, 0), .2)
	if err != nil {
		t.Fatal(err)
	}
	fitted, _, _, _, _, err := fitGeometry(projected, in.Layout)
	if err != nil {
		t.Fatal(err)
	}
	return fitted
}
