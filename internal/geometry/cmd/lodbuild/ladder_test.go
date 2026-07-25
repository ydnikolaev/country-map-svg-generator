package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	geometry "github.com/ydnikolaev/country-map-svg-generator/internal/geometry"
)

func committedLadderArtifactPath() string {
	return filepath.Join(sourceRoot(), "internal/geometry/lod/ladder.artifact.json")
}

// TestLadderRecipeValidates is a fast check that the committed ladder recipe
// still binds the currently embedded oracle and the currently committed v2
// recipe by digest, without running Mapshaper.
func TestLadderRecipeValidates(t *testing.T) {
	c := embeddedCorpus(t)
	oracle, err := geometry.EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	r, recipeSHA, err := loadLadderRecipe(sourceRoot(), c, oracle)
	if err != nil {
		t.Fatal(err)
	}
	if recipeSHA == "" || r.IdentityFallback != "tried_last" || r.CandidateOrder != "fine_to_coarse" {
		t.Fatalf("recipe: %+v sha=%s", r, recipeSHA)
	}
	if len(r.Resolutions) != 30 {
		t.Fatalf("resolutions=%d want 30 (frozen 29 + one DEC-009 intermediate rung)", len(r.Resolutions))
	}
	foundIntermediate := false
	for _, res := range r.Resolutions {
		if res < 64 || res > 80 {
			continue
		}
		if res != 64 && res != 80 {
			foundIntermediate = true
		}
	}
	if !foundIntermediate {
		t.Fatal("no intermediate rung found strictly between the declared 64 and 80")
	}
}

// TestLadderArtifactSchemaAndCoverage is a fast structural check of the
// committed artifact: exact 996-row coverage, digest binding to the currently
// embedded oracle/recipe, and the DEC-009 resolution of SH, UM and HR. It does
// not rebuild anything.
func TestLadderArtifactSchemaAndCoverage(t *testing.T) {
	c := embeddedCorpus(t)
	oracle, err := geometry.EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	_, ladderRecipeSHA, err := loadLadderRecipe(sourceRoot(), c, oracle)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(committedLadderArtifactPath())
	if err != nil {
		t.Fatal(err)
	}
	var artifact ladderArtifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		t.Fatal(err)
	}
	if artifact.SchemaVersion != 1 || artifact.Kind != "silhouette-ladder/v1" ||
		artifact.SourceCorpus != c.Manifest.Identity || artifact.OracleSHA256 != readOracleSHA256(t) ||
		artifact.LadderRecipeSHA256 != ladderRecipeSHA {
		t.Fatalf("artifact identity mismatch: %+v", artifact)
	}
	if len(artifact.Rows) != c.Manifest.EntityCount*4 {
		t.Fatalf("rows=%d want=%d", len(artifact.Rows), c.Manifest.EntityCount*4)
	}
	seen := map[[3]string]bool{}
	pass, noArtifact := 0, 0
	for _, row := range artifact.Rows {
		key := [3]string{row.Alpha2, row.Profile, row.Band}
		if seen[key] {
			t.Fatalf("duplicate row %v", key)
		}
		seen[key] = true
		switch row.Status {
		case "pass":
			pass++
			if _, ok := artifact.Candidates[row.GeometryID][row.Selection]; !ok {
				t.Fatalf("row %v selection %q has no stored candidate", key, row.Selection)
			}
		case "no_artifact":
			noArtifact++
			if row.NoArtifactReason == "" {
				t.Fatalf("row %v is no_artifact with no typed reason", key)
			}
		default:
			t.Fatalf("row %v has unknown status %q", key, row.Status)
		}
	}
	if pass != 993 || noArtifact != 3 {
		t.Fatalf("pass=%d no_artifact=%d want pass=993 no_artifact=3", pass, noArtifact)
	}
	want := map[[3]string]string{
		{"SH", "un", "compact"}:        "identity",
		{"SH", "un", "standard"}:       "identity",
		{"SH", "de_facto", "compact"}:  "identity",
		{"SH", "de_facto", "standard"}: "identity",
		{"UM", "un", "standard"}:       "identity",
		{"UM", "de_facto", "standard"}: "identity",
	}
	for key, selection := range want {
		row := findLadderRow(artifact.Rows, key)
		if row == nil || row.Status != "pass" || row.Selection != selection {
			t.Fatalf("%v: got %+v want selection=%s", key, row, selection)
		}
	}
	// Croatia's un card is no longer pinned to a rung. It used to need DEC-009's
	// intermediate 79.63 with two bytes of headroom; the card is now fitted to the
	// drawn silhouette and the compact contribution threshold is 111, so the
	// disputed component falls below it and is recorded as subscale rather than
	// forcing the search onto a knife edge. What must hold is that the card
	// exists, fits, and reaches that outcome through the visibility policy rather
	// than by silencing it.
	if row := findLadderRow(artifact.Rows, [3]string{"HR", "un", "compact"}); row == nil ||
		row.Status != "pass" || row.PathBytes > 2200 {
		t.Fatalf("HR/un/compact: got %+v want a passing card inside the 2200 cap", row)
	} else {
		disputed := false
		for _, component := range row.Visibility {
			if component.SourceOrder == 2 {
				disputed = component.Subscale && component.OmissionReason == "subscale" && !component.Visible
			}
		}
		if !disputed {
			t.Fatalf("HR/un/compact: component 2 must be recorded as subscale with provenance, got %+v", row.Visibility)
		}
	}
	for _, key := range [][3]string{{"UM", "un", "compact"}, {"UM", "de_facto", "compact"}, {"HR", "de_facto", "compact"}} {
		row := findLadderRow(artifact.Rows, key)
		if row == nil || row.Status != "no_artifact" {
			t.Fatalf("%v: got %+v want no_artifact", key, row)
		}
	}
}

func findLadderRow(rows []ladderRow, key [3]string) *ladderRow {
	for i := range rows {
		if rows[i].Alpha2 == key[0] && rows[i].Profile == key[1] && rows[i].Band == key[2] {
			return &rows[i]
		}
	}
	return nil
}

func readOracleSHA256(t *testing.T) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(sourceRoot(), "internal/geometry/lod/silhouette-oracle.v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

// TestLadderRebuildMatchesCommittedArtifact is the T1 determinism gate: a
// clean rebuild from the frozen inputs must reproduce the committed artifact
// byte-for-byte. It runs the full Mapshaper ladder sweep and pure-Go oracle
// search over the whole catalog, so it is slow (a few minutes), matching the
// existing full-catalog builder tests in this package.
func TestLadderRebuildMatchesCommittedArtifact(t *testing.T) {
	c := embeddedCorpus(t)
	artifact, err := buildLadder(c, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	rebuilt, err := marshalLadderArtifact(artifact)
	if err != nil {
		t.Fatal(err)
	}
	committed, err := os.ReadFile(committedLadderArtifactPath())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rebuilt, committed) {
		t.Fatalf("rebuild drift: rebuilt_sha256=%s committed_sha256=%s", digest(rebuilt), digest(committed))
	}
}
