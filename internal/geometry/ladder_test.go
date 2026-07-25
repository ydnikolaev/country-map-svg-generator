package geometry

import (
	"encoding/json"
	"testing"
)

func embeddedLadder(t *testing.T) *LadderTable {
	t.Helper()
	table, err := EmbeddedLadderTable()
	if err != nil {
		t.Fatal(err)
	}
	return table
}

// TestEmbeddedLadderTableBindsFrozenInputs is the load-time contract: the
// committed artifact must bind the corpus, oracle, ladder recipe and v2 recipe
// this binary actually carries. A build that ships a stale artifact must fail
// here rather than serve geometry judged against different inputs.
func TestEmbeddedLadderTableBindsFrozenInputs(t *testing.T) {
	table := embeddedLadder(t)
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	if table.Kind != "silhouette-ladder/v1" || table.SchemaVersion != 1 {
		t.Fatalf("kind=%q schema_version=%d", table.Kind, table.SchemaVersion)
	}
	if table.CandidateOrder != "fine_to_coarse" {
		t.Fatalf("candidate_order=%q", table.CandidateOrder)
	}
	if table.VisibilityPolicy != oracle.VisibilityPolicy {
		t.Fatalf("visibility_policy=%q oracle=%q", table.VisibilityPolicy, oracle.VisibilityPolicy)
	}
	if table.OracleSHA256 != digestHex(silhouetteOracleJSON) ||
		table.LadderRecipeSHA256 != digestHex(ladderRecipeJSON) ||
		table.BaseRecipeSHA256 != digestHex(ladderBaseRecipeJSON) {
		t.Fatal("committed digests do not bind the embedded oracle/recipes")
	}
}

// TestLadderTableCoverageAndTotals pins the committed catalog outcome: 996
// rows over 498 compact and 498 standard, 993 pass and 3 no_artifact. These
// totals are DEC-009's and DEC-010's result and cannot drift silently.
func TestLadderTableCoverageAndTotals(t *testing.T) {
	table := embeddedLadder(t)
	if len(table.Rows) != 996 {
		t.Fatalf("rows=%d want 996", len(table.Rows))
	}
	byBand := map[string]int{}
	pass, noArtifact := 0, 0
	for _, row := range table.Rows {
		byBand[row.Band]++
		switch row.Status {
		case string(LadderPass):
			pass++
		case string(LadderNoArtifact):
			noArtifact++
		}
	}
	if byBand["compact"] != 498 || byBand["standard"] != 498 {
		t.Fatalf("band coverage: %v", byBand)
	}
	if pass != 993 || noArtifact != 3 {
		t.Fatalf("pass=%d no_artifact=%d want 993/3", pass, noArtifact)
	}
	if got := table.CandidateGeometries(); got != 283 {
		t.Fatalf("candidate geometries=%d want 283", got)
	}
}

// TestLadderKeyIsFunctional proves the index's premise: (geometry, band)
// determines the committed outcome, so the profile need not be carried into
// the lookup key. ParseLadderTable rejects a violation, but a passing parse
// alone would not prove the keys actually collapse — this counts them.
func TestLadderKeyIsFunctional(t *testing.T) {
	table := embeddedLadder(t)
	if len(table.byKey) != 566 {
		t.Fatalf("distinct (geometry, band) keys=%d want 566", len(table.byKey))
	}
	seen := map[ladderKey]LadderRow{}
	for _, row := range table.Rows {
		key := ladderKey{row.GeometryID, row.Band}
		if prior, ok := seen[key]; ok {
			if prior.Status != row.Status || prior.Selection != row.Selection {
				t.Fatalf("%v: %s/%s vs %s/%s", key, prior.Status, prior.Selection, row.Status, row.Selection)
			}
			continue
		}
		seen[key] = row
	}
}

// TestLadderLookupOutcomes covers the three outcomes, including the typed
// absence DEC-010 accepted for Croatia. A no-artifact row must be
// distinguishable from a key the artifact never carried.
func TestLadderLookupOutcomes(t *testing.T) {
	table := embeddedLadder(t)
	rowFor := func(alpha2, profile, band string) LadderRow {
		t.Helper()
		for _, row := range table.Rows {
			if row.Alpha2 == alpha2 && row.Profile == profile && row.Band == band {
				return row
			}
		}
		t.Fatalf("no row for %s/%s/%s", alpha2, profile, band)
		return LadderRow{}
	}

	// DEC-010: Croatia's de_facto card is a decided no-artifact, and it must
	// carry the rejection log that evidenced the decision.
	hr := rowFor("HR", "de_facto", "compact")
	selection, outcome := table.Lookup(hr.GeometryID, "compact")
	if outcome != LadderNoArtifact {
		t.Fatalf("HR/de_facto/compact outcome=%s want %s", outcome, LadderNoArtifact)
	}
	if selection.Row.NoArtifactReason == "" || len(selection.Row.Attempts) == 0 {
		t.Fatal("HR/de_facto/compact carries no typed reason or rejection log")
	}
	if selection.Record.Geometry.ID != "" {
		t.Fatal("a no-artifact outcome must carry no candidate geometry")
	}

	// The same Croatia geometry passes at the hero band — the outcome is per
	// band, not per entity.
	if _, outcome := table.Lookup(hr.GeometryID, "standard"); outcome != LadderPass {
		t.Fatalf("HR/de_facto/standard outcome=%s want %s", outcome, LadderPass)
	}

	// DEC-009's identity fallback, tried last.
	sh := rowFor("SH", "un", "compact")
	selection, outcome = table.Lookup(sh.GeometryID, "compact")
	if outcome != LadderPass || selection.Row.Selection != LadderIdentitySelection {
		t.Fatalf("SH/un/compact outcome=%s selection=%q want pass/identity", outcome, selection.Row.Selection)
	}

	// Croatia keeps its un card. It used to need DEC-009's intermediate 79.63
	// rung with two bytes of headroom; fitting the card to the drawn silhouette
	// moved the byte count, and raising the compact contribution threshold by one
	// let the disputed component fall below it, so a plain declared rung now
	// carries it with room to spare.
	hrUN := rowFor("HR", "un", "compact")
	if selection, outcome := table.Lookup(hrUN.GeometryID, "compact"); outcome != LadderPass || selection.Row.PathBytes > 2200 {
		t.Fatalf("HR/un/compact outcome=%s selection=%q bytes=%d want a passing card inside the 2200 cap",
			outcome, selection.Row.Selection, selection.Row.PathBytes)
	}

	if _, outcome := table.Lookup("geo-does-not-exist", "compact"); outcome != LadderAbsent {
		t.Fatalf("unknown geometry outcome=%s want %s", outcome, LadderAbsent)
	}
	if _, outcome := table.Lookup(sh.GeometryID, "no-such-band"); outcome != LadderAbsent {
		t.Fatalf("unknown band outcome=%s want %s", outcome, LadderAbsent)
	}
}

// TestLadderPassRowsResolveCandidates proves every committed pass row can
// actually be served: its stored candidate exists and carries geometry.
func TestLadderPassRowsResolveCandidates(t *testing.T) {
	table := embeddedLadder(t)
	for _, row := range table.Rows {
		if row.Status != string(LadderPass) {
			continue
		}
		record, ok := table.Candidate(row.GeometryID, row.Selection)
		if !ok {
			t.Fatalf("%s/%s/%s: no candidate for selection %s", row.Alpha2, row.Profile, row.Band, row.Selection)
		}
		if record.Geometry.ID != row.GeometryID || len(record.Geometry.Coordinates) == 0 {
			t.Fatalf("%s/%s/%s: candidate geometry id=%q parts=%d",
				row.Alpha2, row.Profile, row.Band, record.Geometry.ID, len(record.Geometry.Coordinates))
		}
		if record.Projection.Version != "v1" {
			t.Fatalf("%s/%s/%s: candidate projection version=%q", row.Alpha2, row.Profile, row.Band, record.Projection.Version)
		}
	}
}

// TestParseLadderTableRejectsStaleBinding is the mutation tooth for the load
// gate: every frozen digest the artifact binds must be load-blocking.
func TestParseLadderTableRejectsStaleBinding(t *testing.T) {
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	corpusIdentity := embeddedLadder(t).SourceCorpus

	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(ladderArtifactJSON, &decoded); err != nil {
		t.Fatal(err)
	}
	mutate := func(t *testing.T, field string, value any) []byte {
		t.Helper()
		clone := make(map[string]json.RawMessage, len(decoded))
		for k, v := range decoded {
			clone[k] = v
		}
		raw, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		clone[field] = raw
		out, err := json.Marshal(clone)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}

	const stale = "0000000000000000000000000000000000000000000000000000000000000000"
	for _, tc := range []struct {
		name  string
		field string
		value any
	}{
		{"stale corpus", "source_corpus", "sha256:" + stale},
		{"stale oracle", "oracle_sha256", stale},
		{"stale ladder recipe", "ladder_recipe_sha256", stale},
		{"stale base recipe", "base_recipe_sha256", stale},
		{"reordered candidates", "candidate_order", "coarse_to_fine"},
		{"foreign visibility policy", "visibility_policy", "protected-visibility/v2"},
		{"unsupported schema", "schema_version", 2},
		{"foreign kind", "kind", "silhouette-ladder/v2"},
		{"suppressed rows", "rows", []LadderRow{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseLadderTable(mutate(t, tc.field, tc.value), corpusIdentity, oracle); err == nil {
				t.Fatalf("%s loaded without error", tc.name)
			}
		})
	}
}

// TestParseLadderTableRejectsIncoherentRows covers the per-row invariants: a
// pass row that cannot be served, a no-artifact row that hides why, and a
// (geometry, band) key that resolves two ways.
func TestParseLadderTableRejectsIncoherentRows(t *testing.T) {
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	table := embeddedLadder(t)

	candidates := map[string]map[string]ProjectedLODGeometry{"geo-a": {"64": {}}}
	newArtifact := func(rows []LadderRow) []byte {
		t.Helper()
		raw, err := json.Marshal(ladderArtifact{
			SchemaVersion: 1, Kind: "silhouette-ladder/v1",
			SourceCorpus: table.SourceCorpus, OracleSHA256: table.OracleSHA256,
			LadderRecipeSHA256: table.LadderRecipeSHA256, BaseRecipeSHA256: table.BaseRecipeSHA256,
			CandidateOrder: "fine_to_coarse", VisibilityPolicy: table.VisibilityPolicy,
			Candidates: candidates, Rows: rows,
		})
		if err != nil {
			t.Fatal(err)
		}
		return raw
	}

	for _, tc := range []struct {
		name string
		rows []LadderRow
	}{
		{"pass without a stored candidate", []LadderRow{
			{Alpha2: "AA", Profile: "un", Band: "compact", GeometryID: "geo-a", Status: "pass", Selection: "32"},
		}},
		{"pass without a selection", []LadderRow{
			{Alpha2: "AA", Profile: "un", Band: "compact", GeometryID: "geo-a", Status: "pass"},
		}},
		{"no_artifact without a rejection log", []LadderRow{
			{Alpha2: "AA", Profile: "un", Band: "compact", GeometryID: "geo-a", Status: "no_artifact", NoArtifactReason: "why"},
		}},
		{"unknown status", []LadderRow{
			{Alpha2: "AA", Profile: "un", Band: "compact", GeometryID: "geo-a", Status: "degraded", Selection: "64"},
		}},
		{"unknown band", []LadderRow{
			{Alpha2: "AA", Profile: "un", Band: "billboard", GeometryID: "geo-a", Status: "pass", Selection: "64"},
		}},
		{"non-functional (geometry, band) key", []LadderRow{
			{Alpha2: "AA", Profile: "un", Band: "compact", GeometryID: "geo-a", Status: "pass", Selection: "64"},
			{Alpha2: "AA", Profile: "de_facto", Band: "compact", GeometryID: "geo-a", Status: "no_artifact", NoArtifactReason: "why", Attempts: []LadderAttempt{{Selection: "64", Reason: "r"}}},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseLadderTable(newArtifact(tc.rows), table.SourceCorpus, oracle); err == nil {
				t.Fatalf("%s loaded without error", tc.name)
			}
		})
	}
}

// TestPublishedTableIsTheCommittedLadder pins the T3 flip: the shipped table is
// the committed ladder, its band ceilings come from the oracle rather than a
// constant, and every passing row is reachable through the runtime maps.
func TestPublishedTableIsTheCommittedLadder(t *testing.T) {
	table, err := publishedLODTable()
	if err != nil {
		t.Fatal(err)
	}
	if table == nil || table.Ladder == nil {
		t.Fatal("the published table is not ladder-backed")
	}
	ladder := embeddedLadder(t)
	if table.RecipeSHA256 != ladder.LadderRecipeSHA256 {
		t.Fatalf("published recipe %q does not bind the ladder recipe %q", table.RecipeSHA256, ladder.LadderRecipeSHA256)
	}
	oracle, err := EmbeddedSilhouetteOracle()
	if err != nil {
		t.Fatal(err)
	}
	for _, band := range oracle.Bands {
		got := table.StandardMaximumScale
		if band.ID == "compact" {
			got = table.CompactMaximumScale
		}
		if got != band.MaximumEffectiveScale {
			t.Fatalf("%s ceiling=%g want the oracle's %g", band.ID, got, band.MaximumEffectiveScale)
		}
	}

	compact, standard := 0, 0
	for _, row := range ladder.Rows {
		if row.Status != string(LadderPass) {
			continue
		}
		switch row.Band {
		case "compact":
			if _, ok := table.Compact[row.GeometryID]; !ok {
				t.Fatalf("%s/%s/compact is not reachable in the published table", row.Alpha2, row.Profile)
			}
			compact++
		case "standard":
			if _, ok := table.Standard[row.GeometryID]; !ok {
				t.Fatalf("%s/%s/standard is not reachable in the published table", row.Alpha2, row.Profile)
			}
			standard++
		}
	}
	if compact == 0 || standard == 0 {
		t.Fatalf("published table coverage: compact=%d standard=%d", compact, standard)
	}
}
