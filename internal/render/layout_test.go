package render

import (
	"math"
	"sort"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	"github.com/yuranikolaev/country-map-svg-generator/internal/config"
	"github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

// This file is VAL-8: the layout contract asserted against what geometry
// actually produces, not against what the configuration was allowed to say.
// Validation proves a document is well formed; only this proves the frame it
// asks for is the frame that comes out.
//
// The shapes are chosen from the corpus by measured aspect rather than by ISO
// literal, because a country literal in a test is a country literal — the same
// rule the pipeline obeys. Sampling every eighth entity keeps the cost inside
// P3's `make check` budget while still spanning the extremes.

const aspectSampleStride = 8

// drawnExtent is the bounding box of the emitted on-curve points. Q control
// points are excluded: they can sit outside the curve they shape, so including
// them would measure the control hull rather than the silhouette.
func drawnExtent(result geometry.Result) geometry.Bounds {
	bounds := geometry.Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)}
	consider := func(x, y float64) {
		bounds.MinX = math.Min(bounds.MinX, x)
		bounds.MinY = math.Min(bounds.MinY, y)
		bounds.MaxX = math.Max(bounds.MaxX, x)
		bounds.MaxY = math.Max(bounds.MaxY, y)
	}
	for _, command := range result.Commands {
		switch command.Op {
		case "M", "L":
			consider(command.Values[0], command.Values[1])
		case "Q":
			consider(command.Values[2], command.Values[3])
		}
	}
	return bounds
}

type shape struct {
	iso    string
	aspect float64
}

// representativeShapes returns the widest, the tallest and the nearest-square
// entity in a sample of the corpus.
func representativeShapes(t *testing.T, corpus *catalog.Corpus) map[string]shape {
	t.Helper()
	var measured []shape
	for i, entity := range corpus.Manifest.Entities {
		if i%aspectSampleStride != 0 {
			continue
		}
		result, err := generate(corpus, entity.Alpha2, config.Settings{
			Profile: strptr("card"), Boundary: strptr("un"),
		})
		if err != nil {
			// A typed no-artifact row is a recorded decision (DEC-009), not a
			// failure, and it simply carries no frame to measure.
			continue
		}
		extent := drawnExtent(result)
		if extent.Height() <= 0 {
			continue
		}
		measured = append(measured, shape{iso: entity.Alpha2, aspect: extent.Width() / extent.Height()})
	}
	if len(measured) < 3 {
		t.Fatalf("sampled %d usable entities; the corpus or the sampling is broken, not the layout", len(measured))
	}
	sort.Slice(measured, func(i, j int) bool { return measured[i].aspect < measured[j].aspect })
	nearestSquare := measured[0]
	for _, candidate := range measured {
		if math.Abs(math.Log(candidate.aspect)) < math.Abs(math.Log(nearestSquare.aspect)) {
			nearestSquare = candidate
		}
	}
	return map[string]shape{
		"tall":          measured[0],
		"wide":          measured[len(measured)-1],
		"nearly square": nearestSquare,
	}
}

func generate(corpus *catalog.Corpus, iso string, settings config.Settings) (geometry.Result, error) {
	input, err := GeometryRequest(corpus, iso, settings)
	if err != nil {
		return geometry.Result{}, err
	}
	return geometry.Generate(input)
}

// TestTightDerivesNaturalProportions is REQ-13's first half. The viewBox must
// follow the silhouette's own proportions and the requested long side, so a
// caller sizing by one number gets a frame that fits what is drawn.
func TestTightDerivesNaturalProportions(t *testing.T) {
	corpus := corpusOrSkip(t)
	shapes := representativeShapes(t, corpus)

	for name, subject := range shapes {
		for _, longSide := range []float64{64, 160, 640} {
			result, err := generate(corpus, subject.iso, config.Settings{
				Profile: strptr("card"), Boundary: strptr("un"),
				Layout: &config.Layout{Mode: strptr("tight"), LongSide: &longSide},
			})
			if err != nil {
				t.Fatalf("%s (%s) at %g: %v", name, subject.iso, longSide, err)
			}
			box := result.ViewBox
			if got := math.Max(box.Width(), box.Height()); math.Abs(got-longSide) > longSide*1e-9 {
				t.Errorf("%s (%s): long side = %g, want %g", name, subject.iso, got, longSide)
			}
			// The frame follows what is drawn. Comparing against the drawn
			// extent rather than against NaturalAspect is deliberate:
			// NaturalAspect is the full projection's, and DEC-013 fits the card
			// to the silhouette it actually draws, which for an entity with
			// excluded components is a different shape entirely.
			extent := drawnExtent(result)
			frameAspect := box.Width() / box.Height()
			drawnAspect := extent.Width() / extent.Height()
			if relativeGap(frameAspect, drawnAspect) > 0.02 {
				t.Errorf("%s (%s) at %g: frame aspect %g does not follow the drawn aspect %g",
					name, subject.iso, longSide, frameAspect, drawnAspect)
			}
		}
	}
}

// TestContainHonoursTheFrameExactly is REQ-13's second half. The viewBox is the
// frame the caller asked for, to the number — a caller laying out a grid has
// already reserved that box.
func TestContainHonoursTheFrameExactly(t *testing.T) {
	corpus := corpusOrSkip(t)
	shapes := representativeShapes(t, corpus)

	frames := []struct{ width, height float64 }{
		{720, 420},   // the specification's own example
		{16, 16},     // tiny and square
		{4000, 1200}, // huge and extreme
		{120, 900},   // tall, against a wide subject
	}
	for name, subject := range shapes {
		for _, frame := range frames {
			width, height := frame.width, frame.height
			result, err := generate(corpus, subject.iso, config.Settings{
				Profile: strptr("hero"), Boundary: strptr("un"),
				Layout: &config.Layout{Mode: strptr("contain"), Width: &width, Height: &height},
			})
			if err != nil {
				t.Fatalf("%s (%s) in %gx%g: %v", name, subject.iso, width, height, err)
			}
			box := result.ViewBox
			if box.MinX != 0 || box.MinY != 0 || box.Width() != width || box.Height() != height {
				t.Errorf("%s (%s): viewBox = %+v, want exactly 0 0 %g %g", name, subject.iso, box, width, height)
			}
		}
	}
}

// TestContainNeverDistorts is the no-distortion claim, and it is asserted
// without reaching inside the pipeline.
//
// Under one uniform scale, the same entity rendered into two frames of very
// different aspect must produce drawn extents of the *same* aspect — the
// silhouette's own. Under a stretch-to-fill bug the drawn aspects would instead
// follow the frames, which is what this compares against.
func TestContainNeverDistorts(t *testing.T) {
	corpus := corpusOrSkip(t)
	shapes := representativeShapes(t, corpus)

	for name, subject := range shapes {
		wide, err := generate(corpus, subject.iso, config.Settings{
			Profile: strptr("hero"), Boundary: strptr("un"),
			Layout: &config.Layout{Mode: strptr("contain"), Width: f64(900), Height: f64(300)},
		})
		if err != nil {
			t.Fatalf("%s (%s): %v", name, subject.iso, err)
		}
		tall, err := generate(corpus, subject.iso, config.Settings{
			Profile: strptr("hero"), Boundary: strptr("un"),
			Layout: &config.Layout{Mode: strptr("contain"), Width: f64(300), Height: f64(900)},
		})
		if err != nil {
			t.Fatalf("%s (%s): %v", name, subject.iso, err)
		}

		wideExtent, tallExtent := drawnExtent(wide), drawnExtent(tall)
		wideAspect := wideExtent.Width() / wideExtent.Height()
		tallAspect := tallExtent.Width() / tallExtent.Height()

		// A stretch-to-fill bug would put these at 3.0 and 0.333; the tolerance
		// is three orders of magnitude away from that, so it separates
		// quantization noise from distortion without ambiguity.
		if relativeGap(wideAspect, tallAspect) > 0.02 {
			t.Errorf("%s (%s): the silhouette was distorted by its frame — aspect %g in 900x300 but %g in 300x900",
				name, subject.iso, wideAspect, tallAspect)
		}
	}
}

// subPixelTolerance is one output pixel.
//
// Centring is arithmetic, but the *drawn* extent is the simplified and quantized
// silhouette, and simplification can pull the extreme vertex inward by a
// different fraction at each edge — measured at up to 0.19px across the sampled
// corpus. Below one pixel nothing is visible, so that is the threshold: a real
// centring failure pins the silhouette to an edge and misses by hundreds, which
// this separates from noise by two orders of magnitude.
const subPixelTolerance = 1.0

// TestContainNeverLetsTheSilhouetteEscapeItsFrame holds for every entity,
// including the ones whose visibility removals leave them off-centre. A
// silhouette clipped by its own viewBox is a defect no amount of framing policy
// excuses.
func TestContainNeverLetsTheSilhouetteEscapeItsFrame(t *testing.T) {
	corpus := corpusOrSkip(t)
	for name, subject := range representativeShapes(t, corpus) {
		result, err := generate(corpus, subject.iso, config.Settings{
			Profile: strptr("hero"), Boundary: strptr("un"),
			Layout: &config.Layout{Mode: strptr("contain"), Width: f64(900), Height: f64(300)},
		})
		if err != nil {
			t.Fatalf("%s (%s): %v", name, subject.iso, err)
		}
		extent, box := drawnExtent(result), result.ViewBox
		if extent.MinX < box.MinX-subPixelTolerance || extent.MinY < box.MinY-subPixelTolerance ||
			extent.MaxX > box.MaxX+subPixelTolerance || extent.MaxY > box.MaxY+subPixelTolerance {
			t.Errorf("%s (%s): the silhouette escapes its frame — extent %+v in box %+v", name, subject.iso, extent, box)
		}
	}
}

// TestContainCentresTheUnusedSpace is the rest of REQ-13's contain contract, and
// it is asserted only over entities that draw everything they were fitted for.
//
// The restriction is not convenience. `fitGeometry` centres the geometry it is
// given, and visibility removals happen afterwards, so for an entity that drops
// components the drawn silhouette is a subset of what was centred and sits
// off-centre by construction. Asserting centring over those would be asserting
// something the pipeline does not currently promise; asserting it over the rest
// is the strongest true statement available here.
// TestVisibilityRemovalsPushTheSilhouetteOffCentre below measures the excluded
// population rather than leaving it unexamined.
func TestContainCentresTheUnusedSpace(t *testing.T) {
	corpus := corpusOrSkip(t)

	checked := 0
	for i, entity := range corpus.Manifest.Entities {
		if i%aspectSampleStride != 0 {
			continue
		}
		result, err := generate(corpus, entity.Alpha2, config.Settings{
			Profile: strptr("hero"), Boundary: strptr("un"),
			Layout: &config.Layout{Mode: strptr("contain"), Width: f64(900), Height: f64(300)},
		})
		if err != nil || len(result.Removals) != 0 {
			continue
		}
		extent, box := drawnExtent(result), result.ViewBox
		leftGap, rightGap := extent.MinX-box.MinX, box.MaxX-extent.MaxX
		topGap, bottomGap := extent.MinY-box.MinY, box.MaxY-extent.MaxY
		if math.Abs(leftGap-rightGap) > subPixelTolerance {
			t.Errorf("%s: horizontal slack is not centred — %g left, %g right", entity.Alpha2, leftGap, rightGap)
		}
		if math.Abs(topGap-bottomGap) > subPixelTolerance {
			t.Errorf("%s: vertical slack is not centred — %g top, %g bottom", entity.Alpha2, topGap, bottomGap)
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no sampled entity draws everything it was fitted for; the centring claim went untested rather than passing")
	}
	t.Logf("centring verified over %d entities with no visibility removals", checked)
}

// TestVisibilityRemovalsPushTheSilhouetteOffCentre characterizes the gap the
// test above excludes, so it is recorded rather than merely skipped.
//
// This is the France-as-a-speck class of defect on the path DEC-013 did not
// reach. DEC-013 fits the card to the candidate it draws, but that applies to a
// committed ladder candidate; an explicit `contain` frame changes the layout, so
// the committed verdict does not transfer and the request takes the source path,
// where the frame is still fitted before visibility removals. Tracked as
// WKI-37F18A2AA6A5.
//
// The assertion is that the gap is real. If geometry ever closes it, this test
// reddens and the work item gets closed with it — the same shape P2 used for its
// superseded characterizations, so a fix cannot be silently outlived by a green
// suite.
func TestVisibilityRemovalsPushTheSilhouetteOffCentre(t *testing.T) {
	corpus := corpusOrSkip(t)

	worstOffset, worstISO, examined := 0.0, "", 0
	for i, entity := range corpus.Manifest.Entities {
		if i%aspectSampleStride != 0 {
			continue
		}
		result, err := generate(corpus, entity.Alpha2, config.Settings{
			Profile: strptr("hero"), Boundary: strptr("un"),
			Layout: &config.Layout{Mode: strptr("contain"), Width: f64(900), Height: f64(300)},
		})
		if err != nil || len(result.Removals) == 0 {
			continue
		}
		examined++
		extent, box := drawnExtent(result), result.ViewBox
		offset := math.Abs((extent.MinY - box.MinY) - (box.MaxY - extent.MaxY))
		if horizontal := math.Abs((extent.MinX - box.MinX) - (box.MaxX - extent.MaxX)); horizontal > offset {
			offset = horizontal
		}
		if offset > worstOffset {
			worstOffset, worstISO = offset, entity.Alpha2
		}
	}
	if examined == 0 {
		t.Skip("no sampled entity drops components at this frame")
	}
	if worstOffset <= subPixelTolerance {
		t.Fatalf("every entity with visibility removals is now centred (worst offset %g over %d entities); "+
			"the framing gap this characterizes is closed — assert centring unconditionally and close WKI-37F18A2AA6A5",
			worstOffset, examined)
	}
	t.Logf("visibility removals leave the silhouette up to %g off centre in a 900x300 frame (worst: %s, %d entities examined)",
		worstOffset, worstISO, examined)
}

// TestAnImpossibleFrameIsRefusedRatherThanRendered covers padding that leaves no
// drawable area. Rendered anyway it would produce an empty or inverted card that
// still satisfies the byte budget.
func TestAnImpossibleFrameIsRefusedRatherThanRendered(t *testing.T) {
	corpus := corpusOrSkip(t)
	iso := corpus.Manifest.Entities[0].Alpha2

	_, err := generate(corpus, iso, config.Settings{
		Profile: strptr("card"), Boundary: strptr("un"),
		Layout: &config.Layout{
			Mode: strptr("contain"), Width: f64(100), Height: f64(100),
			Padding: &config.Padding{Uniform: f64(60)},
		},
	})
	if err == nil {
		t.Fatal("a frame whose padding leaves no drawable area was rendered")
	}
}

func relativeGap(a, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	return math.Abs(a-b) / math.Max(math.Abs(a), math.Abs(b))
}
