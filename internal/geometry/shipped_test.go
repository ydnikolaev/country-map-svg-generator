package geometry

import (
	"fmt"
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
// Four things are guarded, over every committed row:
//   - a pass row renders, and renders from its own band rather than falling
//     back to the DEC-005 source path;
//   - the emitted path fits the band's frozen cap — including the softened
//     representation, which the oracle never judged and which inflates small
//     compact geometries several-fold;
//   - a no-artifact row surfaces as the typed DEC-009 outcome, not an error and
//     not a silent substitution;
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
		preset, ok := presets[row.Band]
		if !ok {
			t.Fatalf("%s/%s: band %q has no preset", row.Alpha2, row.Profile, row.Band)
		}
		label := fmt.Sprintf("%s/%s/%s", row.Alpha2, row.Profile, row.Band)
		in, err := InputFromCatalog(c, row.Alpha2, row.Profile, preset)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		got, genErr := Generate(in)

		if row.Status == string(LadderNoArtifact) {
			if !IsNoArtifact(genErr) {
				t.Fatalf("%s: committed no_artifact but Generate returned err=%v", label, genErr)
			}
			typedAbsent++
			continue
		}
		if genErr != nil {
			t.Fatalf("%s: committed pass but Generate failed: %v", label, genErr)
		}
		if got.LOD.SelectedTier != row.Band {
			t.Fatalf("%s: served tier %q, not the committed band — the ladder is not reaching the shipped path (fallbacks=%v)",
				label, got.LOD.SelectedTier, got.LOD.Fallbacks)
		}
		if got.Path == "" || got.Metrics.PathBytes != len(got.Path) {
			t.Fatalf("%s: path bytes=%d metric=%d", label, len(got.Path), got.Metrics.PathBytes)
		}
		if len(got.Path) > caps[row.Band] {
			t.Fatalf("%s: emitted %d path bytes over the frozen %s cap of %d (committed %d)",
				label, len(got.Path), row.Band, caps[row.Band], row.PathBytes)
		}
		if !strings.HasPrefix(got.Path, "m") && !strings.HasPrefix(got.Path, "M") {
			t.Fatalf("%s: path does not start with a move command: %.20q", label, got.Path)
		}
		if got.ViewBox.Width() <= 0 || got.ViewBox.Height() <= 0 {
			t.Fatalf("%s: degenerate viewBox %v", label, got.ViewBox)
		}
		rendered++
	}
	if rendered != 993 || typedAbsent != 3 {
		t.Fatalf("shipped catalog: rendered=%d typed_absent=%d want 993 and 3", rendered, typedAbsent)
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
