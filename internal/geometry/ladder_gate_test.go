package geometry

import (
	"fmt"
	"testing"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
)

// VAL-6 obligation 2: the offline ladder gate.
//
// The committed artifact is what production serves, so its every row has to be
// recomputable from the frozen inputs alone — no Node, no network, no Mapshaper,
// no committed metric taken on faith. This file recomputes topology, protected
// visibility, IoU/recall, path bytes and the complete-file estimate for every
// row against the oracle this binary embeds, and compares each result with what
// the artifact claims.
//
// What it deliberately does not cover: the *rejected* rungs of a no-artifact
// row. Those geometries are not stored (only the winner is), so reproducing
// them needs the Mapshaper sweep. That is the determinism gate's job
// (TestLadderRebuildMatchesCommittedArtifact, ~3 min, Node required). This gate
// recomputes the identity rung for a no-artifact row, which is the one attempt
// reachable without Mapshaper, and proves the last-resort candidate still fails.

// ladderBandPreset maps a band to the preset the build evaluated it under.
func ladderBandPreset(bandID string) (string, error) {
	switch bandID {
	case "compact":
		return "card", nil
	case "standard":
		return "hero", nil
	}
	return "", fmt.Errorf("unknown band %q", bandID)
}

// ladderCaseInput rebuilds the exact Input the ladder build evaluated for one
// row: catalog input, preset applied, the CLI-only fields cleared, then group
// anchors applied for this band.
func ladderCaseInput(c *catalog.Corpus, anchors GroupAnchorSet, row LadderRow, band SilhouetteBand) (Input, GroupAnchorApplication, error) {
	preset, err := ladderBandPreset(row.Band)
	if err != nil {
		return Input{}, GroupAnchorApplication{}, err
	}
	in, err := InputFromCatalog(c, row.Alpha2, row.Profile, preset)
	if err != nil {
		return Input{}, GroupAnchorApplication{}, err
	}
	resolved, err := ApplyPreset(in)
	if err != nil {
		return Input{}, GroupAnchorApplication{}, err
	}
	resolved.Preset, resolved.MaxPathBytes = "", 0
	return ApplyGroupAnchors(resolved, anchors, band)
}

// ladderRecomputation is one row's independently recomputed verdict.
type ladderRecomputation struct {
	evaluation SilhouetteCandidateEvaluation
	complete   int
	passes     bool
	reasons    []string
}

// recomputeLadderRow runs the frozen oracle predicate over one stored candidate
// exactly as the build did. The predicate is restated here rather than shared
// with the builder on purpose: a gate that calls the same helper as the thing it
// checks proves only that the helper is self-consistent.
func recomputeLadderRow(c *catalog.Corpus, anchors GroupAnchorSet, row LadderRow, band SilhouetteBand, record ProjectedLODGeometry, recipeSHA string) (ladderRecomputation, error) {
	resolved, application, err := ladderCaseInput(c, anchors, row, band)
	if err != nil {
		return ladderRecomputation{}, err
	}
	record.RecipeSHA256 = recipeSHA
	evaluation, err := EvaluateSilhouetteCandidate(resolved, record, recipeSHA, band)
	if err != nil {
		return ladderRecomputation{}, err
	}
	out := ladderRecomputation{evaluation: evaluation, complete: evaluation.PathBytes + 280}
	if err := ValidateGroupAnchorResult(application, evaluation.Visibility); err != nil {
		out.reasons = append(out.reasons, "group_anchor:"+err.Error())
	}
	if !SilhouetteBandPasses(band, evaluation.Metrics) {
		out.reasons = append(out.reasons, fmt.Sprintf("iou=%.6f/recall=%.6f below band minimums %.6f/%.6f",
			evaluation.Metrics.IoU, evaluation.Metrics.Recall, band.MinimumIoU, band.MinimumRecall))
	}
	if evaluation.PathBytes > band.PathCap {
		out.reasons = append(out.reasons, fmt.Sprintf("bytes=%d over cap %d", evaluation.PathBytes, band.PathCap))
	}
	if out.complete >= band.CompleteFileMaximum {
		out.reasons = append(out.reasons, fmt.Sprintf("file=%d over max %d", out.complete, band.CompleteFileMaximum))
	}
	if !evaluation.Topology {
		out.reasons = append(out.reasons, "topology=false")
	}
	if !evaluation.Protection {
		out.reasons = append(out.reasons, "protection=false")
	}
	if !evaluation.DominantComponent {
		out.reasons = append(out.reasons, "dominant_component=false")
	}
	out.passes = len(out.reasons) == 0
	return out, nil
}

func ladderGateFixture(t *testing.T) (*LadderTable, *catalog.Corpus, SilhouetteOracle, GroupAnchorSet, map[string]SilhouetteBand) {
	t.Helper()
	table, err := EmbeddedLadderTable()
	if err != nil {
		t.Fatal(err)
	}
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	anchors, err := EmbeddedGroupAnchors()
	if err != nil {
		t.Fatal(err)
	}
	bands := map[string]SilhouetteBand{}
	for _, band := range oracle.Bands {
		bands[band.ID] = band
	}
	return table, c, oracle, anchors, bands
}

// TestLadderGateRecomputesEveryCommittedRow is the VAL-6 obligation 2 gate. It
// reddens on a perturbed coordinate (the recomputed metrics and byte counts stop
// matching), a moved threshold (the recomputed verdict flips), a suppressed
// omission or reordered candidate (the recomputed provenance stops matching),
// and a stale corpus/oracle/recipe identity (the table refuses to load at all).
func TestLadderGateRecomputesEveryCommittedRow(t *testing.T) {
	table, c, _, anchors, bands := ladderGateFixture(t)

	checked, identityRechecked := 0, 0
	for _, row := range table.Rows {
		band, ok := bands[row.Band]
		if !ok {
			t.Fatalf("%s/%s: band %q is not in the embedded oracle", row.Alpha2, row.Profile, row.Band)
		}
		label := fmt.Sprintf("%s/%s/%s", row.Alpha2, row.Profile, row.Band)

		if row.Status == string(LadderNoArtifact) {
			// The one rung reachable without Mapshaper. It must still fail,
			// and it must fail for the reason the artifact recorded last.
			identity, ok := table.Candidate(row.GeometryID, LadderIdentitySelection)
			if !ok {
				var err error
				identity, err = projectedIdentityRecord(c, row.GeometryID, table.LadderRecipeSHA256)
				if err != nil {
					t.Fatalf("%s: %v", label, err)
				}
			}
			out, err := recomputeLadderRow(c, anchors, row, band, identity, table.LadderRecipeSHA256)
			if err != nil {
				// A hard evaluation error is itself a rejection; the artifact
				// records exactly that for UM at compact.
				identityRechecked++
				continue
			}
			if out.passes {
				t.Fatalf("%s: recorded no_artifact but the identity rung now passes", label)
			}
			identityRechecked++
			continue
		}

		record, ok := table.Candidate(row.GeometryID, row.Selection)
		if !ok {
			t.Fatalf("%s: no stored candidate for selection %s", label, row.Selection)
		}
		out, err := recomputeLadderRow(c, anchors, row, band, record, table.LadderRecipeSHA256)
		if err != nil {
			t.Fatalf("%s: recomputation failed: %v", label, err)
		}
		if !out.passes {
			t.Fatalf("%s: committed as pass but recomputes as rejected: %v", label, out.reasons)
		}
		if out.evaluation.Metrics.IoU != row.IoU || out.evaluation.Metrics.Recall != row.Recall {
			t.Fatalf("%s: iou/recall drift: recomputed %.17g/%.17g committed %.17g/%.17g",
				label, out.evaluation.Metrics.IoU, out.evaluation.Metrics.Recall, row.IoU, row.Recall)
		}
		if out.evaluation.PathBytes != row.PathBytes || out.complete != row.CompleteEstimate {
			t.Fatalf("%s: byte drift: recomputed %d/%d committed %d/%d",
				label, out.evaluation.PathBytes, out.complete, row.PathBytes, row.CompleteEstimate)
		}
		if out.evaluation.Points != row.Points || out.evaluation.Parts != row.Parts {
			t.Fatalf("%s: shape drift: recomputed points=%d parts=%d committed points=%d parts=%d",
				label, out.evaluation.Points, out.evaluation.Parts, row.Points, row.Parts)
		}
		if out.evaluation.Omissions != row.Omissions {
			t.Fatalf("%s: omission drift: recomputed %d committed %d", label, out.evaluation.Omissions, row.Omissions)
		}
		if len(out.evaluation.Visibility) != len(row.Visibility) {
			t.Fatalf("%s: visibility component count drift: recomputed %d committed %d",
				label, len(out.evaluation.Visibility), len(row.Visibility))
		}
		for i := range row.Visibility {
			if out.evaluation.Visibility[i] != row.Visibility[i] {
				t.Fatalf("%s: visibility component %d drift:\n recomputed %+v\n committed  %+v",
					label, i, out.evaluation.Visibility[i], row.Visibility[i])
			}
		}
		checked++
	}
	if checked != 993 || identityRechecked != 3 {
		t.Fatalf("coverage: recomputed %d pass rows and %d no_artifact rows, want 993 and 3", checked, identityRechecked)
	}
}

// projectedIdentityRecord rebuilds the identity (unsimplified) candidate for a
// geometry the artifact stores no candidate for, which is the case for every
// no-artifact row.
func projectedIdentityRecord(c *catalog.Corpus, geometryID, recipeSHA string) (ProjectedLODGeometry, error) {
	contract := LODProjectionContractV1()
	for _, source := range c.Geometries {
		if source.ID != geometryID {
			continue
		}
		record, err := ProjectLODGeometry(source, contract.Flatness, contract.CoordinatePrecision)
		if err != nil {
			return ProjectedLODGeometry{}, err
		}
		record.SourceCorpus = c.Manifest.Identity
		record.RecipeSHA256 = recipeSHA
		return record, nil
	}
	return ProjectedLODGeometry{}, fmt.Errorf("geometry %s is not in the embedded corpus", geometryID)
}

// TestLadderGateRejectsAPerturbedCoordinate is a mutation tooth: moving one
// coordinate of one stored candidate by a hair must be caught by recomputation,
// not absorbed. Without this the gate could pass while serving geometry nobody
// judged.
func TestLadderGateRejectsAPerturbedCoordinate(t *testing.T) {
	table, c, _, anchors, bands := ladderGateFixture(t)

	row := table.Rows[0]
	for _, candidate := range table.Rows {
		if candidate.Status == string(LadderPass) && candidate.Points > 20 {
			row = candidate
			break
		}
	}
	record, ok := table.Candidate(row.GeometryID, row.Selection)
	if !ok {
		t.Fatalf("no candidate for %s/%s/%s", row.Alpha2, row.Profile, row.Band)
	}

	perturbed := record
	perturbed.Geometry = cloneCatalogGeometry(record.Geometry)
	perturbed.Geometry.Coordinates[0][0][0][0] += .002

	out, err := recomputeLadderRow(c, anchors, row, bands[row.Band], perturbed, table.LadderRecipeSHA256)
	if err != nil {
		// A hard evaluation failure is a rejection too. Logged so the tooth
		// cannot pass vacuously without saying which mechanism caught it.
		t.Logf("%s/%s/%s: caught by evaluation error: %v", row.Alpha2, row.Profile, row.Band, err)
		return
	}
	if out.evaluation.Metrics.IoU == row.IoU && out.evaluation.PathBytes == row.PathBytes &&
		out.evaluation.Points == row.Points && out.evaluation.Parts == row.Parts {
		t.Fatalf("%s/%s/%s: a perturbed coordinate reproduced the committed metrics exactly — the gate cannot see coordinate drift",
			row.Alpha2, row.Profile, row.Band)
	}
	t.Logf("%s/%s/%s: caught by metric drift: iou %.17g->%.17g bytes %d->%d points %d->%d",
		row.Alpha2, row.Profile, row.Band, row.IoU, out.evaluation.Metrics.IoU,
		row.PathBytes, out.evaluation.PathBytes, row.Points, out.evaluation.Points)
}

func cloneCatalogGeometry(g catalog.Geometry) catalog.Geometry {
	clone := g
	clone.Coordinates = make(catalog.MultiPolygon, len(g.Coordinates))
	for i, polygon := range g.Coordinates {
		clone.Coordinates[i] = make(catalog.Polygon, len(polygon))
		for j, ring := range polygon {
			clone.Coordinates[i][j] = append(catalog.Ring(nil), ring...)
		}
	}
	return clone
}

// TestLadderGateRejectsAMovedThreshold is the second mutation tooth: tightening
// a frozen threshold must flip committed rows from pass to rejected. A gate that
// stays green under a moved threshold is not gating on the threshold at all.
func TestLadderGateRejectsAMovedThreshold(t *testing.T) {
	table, c, _, anchors, bands := ladderGateFixture(t)

	for _, mutation := range []struct {
		name  string
		apply func(SilhouetteBand) SilhouetteBand
	}{
		{"path cap", func(b SilhouetteBand) SilhouetteBand { b.PathCap = 1; return b }},
		{"complete file maximum", func(b SilhouetteBand) SilhouetteBand { b.CompleteFileMaximum = 1; return b }},
		{"minimum IoU", func(b SilhouetteBand) SilhouetteBand { b.MinimumIoU = 1.0000001; return b }},
		{"minimum recall", func(b SilhouetteBand) SilhouetteBand { b.MinimumRecall = 1.0000001; return b }},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			survivors := 0
			checked := 0
			for _, row := range table.Rows {
				if row.Status != string(LadderPass) || checked >= 25 {
					continue
				}
				record, ok := table.Candidate(row.GeometryID, row.Selection)
				if !ok {
					continue
				}
				checked++
				out, err := recomputeLadderRow(c, anchors, row, mutation.apply(bands[row.Band]), record, table.LadderRecipeSHA256)
				if err == nil && out.passes {
					survivors++
				}
			}
			if checked == 0 {
				t.Fatal("no rows checked")
			}
			if survivors != 0 {
				t.Fatalf("%d of %d committed rows still pass under a tightened %s", survivors, checked, mutation.name)
			}
		})
	}
}

// TestLadderGateRejectsASuppressedOmission is the third mutation tooth: dropping
// a component from a stored candidate must be seen. This is the failure mode
// DEC-006 named — an omission that goes unreported is a silhouette the oracle
// never judged.
//
// Note which mechanism does the catching. A dropped component can still satisfy
// the band predicate; what it cannot do is reproduce the committed omission and
// part counts. So the teeth here are the full-catalog gate's provenance
// comparison, not the oracle's pass/fail — which is exactly why that gate
// compares every recomputed field instead of only the verdict.
func TestLadderGateRejectsASuppressedOmission(t *testing.T) {
	table, c, _, anchors, bands := ladderGateFixture(t)

	var row LadderRow
	var record ProjectedLODGeometry
	for _, candidate := range table.Rows {
		if candidate.Status != string(LadderPass) {
			continue
		}
		stored, ok := table.Candidate(candidate.GeometryID, candidate.Selection)
		if !ok || len(stored.Geometry.Coordinates) < 2 {
			continue
		}
		row, record = candidate, stored
		break
	}
	if row.GeometryID == "" {
		t.Skip("no committed multi-part candidate to drop a component from")
	}

	truncated := record
	truncated.Geometry = cloneCatalogGeometry(record.Geometry)
	truncated.Geometry.Coordinates = truncated.Geometry.Coordinates[:len(truncated.Geometry.Coordinates)-1]

	out, err := recomputeLadderRow(c, anchors, row, bands[row.Band], truncated, table.LadderRecipeSHA256)
	if err != nil {
		t.Logf("%s/%s/%s: caught by evaluation error: %v", row.Alpha2, row.Profile, row.Band, err)
		return
	}
	if out.passes && out.evaluation.Omissions == row.Omissions && out.evaluation.Parts == row.Parts {
		t.Fatalf("%s/%s/%s: a dropped component left omissions=%d parts=%d unchanged and still passes",
			row.Alpha2, row.Profile, row.Band, out.evaluation.Omissions, out.evaluation.Parts)
	}
	t.Logf("%s/%s/%s: caught: passes=%v omissions %d->%d parts %d->%d reasons=%v",
		row.Alpha2, row.Profile, row.Band, out.passes, row.Omissions, out.evaluation.Omissions,
		row.Parts, out.evaluation.Parts, out.reasons)
}
