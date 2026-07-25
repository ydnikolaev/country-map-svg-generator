package geometry

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

// The shipped-path gate. T2's ladder gate proves the committed artifact is
// internally sound; this proves the artifact is what Generate() actually
// serves. They are different claims, and the twenty runs that preceded this
// spec failed precisely in the gap between them: a pipeline can hold every
// internal invariant and still emit source-only geometry.
//
// Five things are guarded, over every committed row:
//   - a pass row renders, and renders from its own band rather than falling
//     back to the DEC-005 source path;
//   - the emitted path fits the band's frozen cap — including the softened
//     representation, which the oracle never judged and which inflates small
//     compact geometries several-fold;
//   - a no-artifact row surfaces as the typed DEC-009 outcome, not an error and
//     not a silent substitution;
//   - the silhouette fills a sane share of its frame, so a card cannot go back
//     to being framed for geometry it does not draw;
//   - nothing errors.
//
// internal/geometry/cmd/svgproof runs the same sweep and writes real SVG files
// and a contact sheet. Use it when a human needs to look; this is the gate.
func TestShippedCatalogServesTheLadder(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	ladder, err := EmbeddedLadderTable()
	if err != nil {
		t.Fatal(err)
	}
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	caps := map[string]int{}
	for _, band := range oracle.Bands {
		caps[band.ID] = band.PathCap
	}
	presets := map[string]string{"compact": "card", "standard": "hero"}

	rendered, typedAbsent := 0, 0
	for _, row := range ladder.Rows {
		outcome, err := checkShippedRow(c, row, presets, caps, Generate)
		if err != nil {
			t.Fatal(err)
		}
		switch outcome {
		case shippedRendered:
			rendered++
		case shippedTypedAbsent:
			typedAbsent++
		}
	}
	if rendered != 993 || typedAbsent != 3 {
		t.Fatalf("shipped catalog: rendered=%d typed_absent=%d want 993 and 3", rendered, typedAbsent)
	}
}

type shippedOutcome int

const (
	shippedRendered shippedOutcome = iota
	shippedTypedAbsent
)

// checkShippedRow holds the assertions for one committed row and returns them
// as an error rather than failing a test directly. That is what lets the
// mutation tooth below drive the same assertions against a deliberately
// ladder-less pipeline and require that they reject it — an assertion that only
// ever runs against the healthy path cannot be shown to bite.
func checkShippedRow(
	c *catalog.Corpus,
	row LadderRow,
	presets map[string]string,
	caps map[string]int,
	generate func(Input) (Result, error),
) (shippedOutcome, error) {
	preset, ok := presets[row.Band]
	if !ok {
		return 0, fmt.Errorf("%s/%s: band %q has no preset", row.Alpha2, row.Profile, row.Band)
	}
	label := fmt.Sprintf("%s/%s/%s", row.Alpha2, row.Profile, row.Band)
	in, err := InputFromCatalog(c, row.Alpha2, row.Profile, preset)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", label, err)
	}
	got, genErr := generate(in)

	if row.Status == string(LadderNoArtifact) {
		if !IsNoArtifact(genErr) {
			return 0, fmt.Errorf("%s: committed no_artifact but Generate returned err=%v", label, genErr)
		}
		return shippedTypedAbsent, nil
	}
	if genErr != nil {
		return 0, fmt.Errorf("%s: committed pass but Generate failed: %w", label, genErr)
	}
	if got.LOD.SelectedTier != row.Band {
		return 0, fmt.Errorf("%s: served tier %q, not the committed band — the ladder is not reaching the shipped path (fallbacks=%v)",
			label, got.LOD.SelectedTier, got.LOD.Fallbacks)
	}
	// Serving the right band is not enough. The pre-DEC-006 tier walk reads the
	// same stored candidates, so it can report the same band while having
	// restored full-detail source components into the geometry and re-judged it
	// by raw boundary deviation — both superseded, and both producing something
	// no oracle ever saw. Provenance distinguishes the two exactly: the ladder
	// path computes no raw deviation and restores nothing, so a non-zero value
	// in either field means the superseded path ran.
	if got.LOD.RawDeviation != 0 || len(got.LOD.Restored) != 0 {
		return 0, fmt.Errorf("%s: provenance shows the superseded derived path (raw_deviation=%g restored=%d); a ladder candidate is judged by the oracle and must be neither re-judged by deviation nor have source components injected",
			label, got.LOD.RawDeviation, len(got.LOD.Restored))
	}
	if got.Path == "" || got.Metrics.PathBytes != len(got.Path) {
		return 0, fmt.Errorf("%s: path bytes=%d metric=%d", label, len(got.Path), got.Metrics.PathBytes)
	}
	if len(got.Path) > caps[row.Band] {
		return 0, fmt.Errorf("%s: emitted %d path bytes over the frozen %s cap of %d (committed %d)",
			label, len(got.Path), row.Band, caps[row.Band], row.PathBytes)
	}
	if !strings.HasPrefix(got.Path, "m") && !strings.HasPrefix(got.Path, "M") {
		return 0, fmt.Errorf("%s: path does not start with a move command: %.20q", label, got.Path)
	}
	if got.ViewBox.Width() <= 0 || got.ViewBox.Height() <= 0 {
		return 0, fmt.Errorf("%s: degenerate viewBox %v", label, got.ViewBox)
	}
	// The card must be mostly silhouette, not mostly empty. This guards DEC-013:
	// while the layout was fitted to the whole territorial claim, a card framed
	// for components it does not draw — France filled 11.9% of its frame, and
	// twenty cards sat under a fifth. Every other assertion here passed
	// throughout, which is exactly why this one exists.
	//
	// The floor is deliberately far below the current catalog. After DEC-013 the
	// thinnest compact card is Marshall Islands at 41%; 20% cannot fire on
	// today's data and would have caught France by a wide margin.
	if coverage := shippedFrameCoverage(got); coverage < 0.20 {
		return 0, fmt.Errorf("%s: the silhouette fills %.1f%% of its frame — the card is framed for geometry it does not draw",
			label, 100*coverage)
	}
	return shippedRendered, nil
}

// shippedFrameCoverage is the share of the viewBox area covered by the drawn
// silhouette's own bounding box, walked from the emitted commands so it measures
// what the file draws rather than what the pipeline intended.
func shippedFrameCoverage(r Result) float64 {
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, command := range r.Commands {
		for i := 0; i+1 < len(command.Values); i += 2 {
			x, y := command.Values[i], command.Values[i+1]
			minX, maxX = math.Min(minX, x), math.Max(maxX, x)
			minY, maxY = math.Min(minY, y), math.Max(maxY, y)
		}
	}
	area := r.ViewBox.Width() * r.ViewBox.Height()
	if area <= 0 || minX > maxX || minY > maxY {
		return 0
	}
	return ((maxX - minX) * (maxY - minY)) / area
}

// TestShippedGateRejectsALadderlessPipeline is the mutation tooth for the gate
// above, and it is the one that matters most: the twenty runs that preceded this
// spec all reported success while emitting source-only geometry, so a gate that
// cannot distinguish "ladder served" from "source served" would have passed
// through the entire failure.
//
// Two regression shapes are driven, both of which were reproduced by hand before
// being automated here:
//
//   - no published table at all, the shape of publishedLODTable regressing to
//     nil;
//   - a published table with the ladder stripped, the shape of the ladder
//     loading but selection no longer consulting it — which falls back to the
//     superseded pre-DEC-006 tier walk.
//
// A deterministic 60-row slice is used rather than the full 996: the assertion
// under test fires on the first row, so the remaining coverage would only cost
// time. The full sweep is the gate's own job.
func TestShippedGateRejectsALadderlessPipeline(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	ladder, err := EmbeddedLadderTable()
	if err != nil {
		t.Fatal(err)
	}
	published, err := publishedLODTable()
	if err != nil {
		t.Fatal(err)
	}
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	caps := map[string]int{}
	for _, band := range oracle.Bands {
		caps[band.ID] = band.PathCap
	}
	presets := map[string]string{"compact": "card", "standard": "hero"}

	ladderless := &LODTable{
		Version: published.Version, RecipeSHA256: published.RecipeSHA256,
		Compact: published.Compact, Standard: published.Standard,
		CompactMaximumScale:  published.CompactMaximumScale,
		StandardMaximumScale: published.StandardMaximumScale,
		CompactPathCap:       published.CompactPathCap,
		StandardPathCap:      published.StandardPathCap,
		// Ladder deliberately absent.
	}

	const sample = 60
	for _, mutation := range []struct {
		name     string
		generate func(Input) (Result, error)
	}{
		{"no published table", func(in Input) (Result, error) { return GenerateWithLOD(in, nil) }},
		{"published table with the ladder stripped", func(in Input) (Result, error) {
			return GenerateWithLOD(in, ladderless)
		}},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			checked, rejected := 0, 0
			var first error
			for _, row := range ladder.Rows {
				if checked >= sample {
					break
				}
				if row.Status != string(LadderPass) {
					continue
				}
				checked++
				if _, err := checkShippedRow(c, row, presets, caps, mutation.generate); err != nil {
					rejected++
					if first == nil {
						first = err
					}
				}
			}
			if checked == 0 {
				t.Fatal("no rows checked")
			}
			if rejected != checked {
				t.Fatalf("%d of %d rows still satisfied the shipped assertions with %s — the gate does not see this regression",
					checked-rejected, checked, mutation.name)
			}
			t.Logf("all %d sampled rows rejected; first: %v", checked, first)
		})
	}
}

// TestShippedPathHoldsTheBandCapAgainstSoftening is the mutation tooth for the
// cap fix above. Grenada's compact card is the smallest committed margin: its
// linear path is 575 bytes but softening inflates it to 2230, which clears the
// 2200 band cap while still fitting the preset's 2500 complete-file maximum.
// Serving that would silently void the guarantee the ladder was built to make.
func TestShippedPathHoldsTheBandCapAgainstSoftening(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	table, err := publishedLODTable()
	if err != nil {
		t.Fatal(err)
	}
	in, err := InputFromCatalog(c, "GD", "un", "card")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Generate(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Path) > table.CompactPathCap {
		t.Fatalf("GD/un/card emitted %d bytes over the %d compact cap", len(got.Path), table.CompactPathCap)
	}

	// Prove the tooth bites: with only the preset's larger ceiling in play, the
	// softened representation is chosen and does clear the band cap.
	resolved, err := ApplyPreset(in)
	if err != nil {
		t.Fatal(err)
	}
	selection, outcome := table.Ladder.Lookup(resolved.Geometry.ID, "compact")
	if outcome != LadderPass {
		t.Fatalf("GD/un/compact outcome=%s", outcome)
	}
	if selection.Row.PathBytes >= table.CompactPathCap {
		t.Fatalf("fixture no longer has margin: committed %d vs cap %d", selection.Row.PathBytes, table.CompactPathCap)
	}
	if len(got.Path) == selection.Row.PathBytes {
		return // the linear path was emitted, which is already within the cap
	}
	if len(got.Path) <= selection.Row.PathBytes {
		t.Fatalf("GD/un/card emitted %d bytes, below the committed linear %d — unexpected representation",
			len(got.Path), selection.Row.PathBytes)
	}
}
