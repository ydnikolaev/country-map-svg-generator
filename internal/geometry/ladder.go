package geometry

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

// The committed DEC-006/DEC-009 ladder and the three frozen inputs it binds by
// digest. Embedding the recipes alongside the artifact is what lets the runtime
// verify the binding without touching the filesystem: a build carrying a stale
// artifact fails to load rather than silently serving geometry that was judged
// against a different oracle, corpus or base recipe.
//
//go:embed lod/ladder.artifact.json
var ladderArtifactJSON []byte

//go:embed lod/ladder.recipe.json
var ladderRecipeJSON []byte

//go:embed lod/v2.recipe.json
var ladderBaseRecipeJSON []byte

// LadderOutcome distinguishes the two committed outcomes from the absence of
// any committed row. DEC-009's typed no-artifact result is a decision the build
// recorded, not a gap; collapsing it into "not found" would let a consumer fall
// back to an unjudged path for exactly the cases the oracle refused.
type LadderOutcome string

const (
	// LadderPass — the build selected a candidate that satisfies the oracle.
	LadderPass LadderOutcome = "pass"
	// LadderNoArtifact — no rung, identity included, could satisfy the oracle
	// at this band's fitted scale (DEC-009; DEC-010 for HR/de_facto/compact).
	LadderNoArtifact LadderOutcome = "no_artifact"
	// LadderAbsent — the artifact carries no row for this key at all.
	LadderAbsent LadderOutcome = "absent"
)

// LadderIdentitySelection is the selection value of the identity rung, tried
// last per DEC-009. Any other selection is a formatted Mapshaper resolution.
const LadderIdentitySelection = "identity"

// LadderAttempt is one rejected rung, preserved for inspectability (DEC-006
// item 7). Present only on no-artifact rows.
type LadderAttempt struct {
	Selection string `json:"selection"`
	Reason    string `json:"reason"`
}

// LadderRow mirrors the committed artifact row exactly. Field names and tags
// must stay in lockstep with lodbuild's ladderRow: the artifact is encoded by
// the builder and decoded here, and a silent name drift would drop provenance
// on decode rather than fail.
type LadderRow struct {
	Alpha2, Profile, Band string
	GeometryID            string
	Status                string
	Selection             string `json:",omitempty"`
	IoU, Recall           float64
	PathBytes             int                            `json:",omitempty"`
	CompleteEstimate      int                            `json:",omitempty"`
	Omissions             int                            `json:",omitempty"`
	Points, Parts         int                            `json:",omitempty"`
	Visibility            []ProtectedVisibilityComponent `json:",omitempty"`
	GroupAnchors          []GroupAnchorResolution        `json:",omitempty"`
	NoArtifactReason      string                         `json:",omitempty"`
	Attempts              []LadderAttempt                `json:",omitempty"`
}

// ladderArtifact is the on-disk shape. It is unexported because consumers work
// through LadderTable, which adds the verified indexes.
type ladderArtifact struct {
	SchemaVersion      int                                        `json:"schema_version"`
	Kind               string                                     `json:"kind"`
	SourceCorpus       string                                     `json:"source_corpus"`
	OracleSHA256       string                                     `json:"oracle_sha256"`
	LadderRecipeSHA256 string                                     `json:"ladder_recipe_sha256"`
	BaseRecipeSHA256   string                                     `json:"base_recipe_sha256"`
	CandidateOrder     string                                     `json:"candidate_order"`
	VisibilityPolicy   string                                     `json:"visibility_policy"`
	Candidates         map[string]map[string]ProjectedLODGeometry `json:"candidates"`
	Rows               []LadderRow                                `json:"rows"`
}

// ladderKey is the committed selection's functional key. Profile is absent by
// construction: un and de_facto resolve to different geometry ids whenever they
// differ at all, so (geometry, band) determines the selection. ParseLadderTable
// enforces that rather than assuming it.
type ladderKey struct{ geometryID, band string }

// LadderTable is the verified, indexed committed ladder. It is data only: it
// answers "what did the build select for this geometry at this band, and what
// was the judged outcome", and holds no selection policy of its own.
type LadderTable struct {
	SchemaVersion      int
	Kind               string
	SourceCorpus       string
	OracleSHA256       string
	LadderRecipeSHA256 string
	BaseRecipeSHA256   string
	CandidateOrder     string
	VisibilityPolicy   string
	Rows               []LadderRow

	byKey      map[ladderKey]int
	candidates map[string]map[string]ProjectedLODGeometry
}

// LadderSelection is one committed pass outcome: the winning rung and the exact
// projected geometry the oracle judged.
type LadderSelection struct {
	Row    LadderRow
	Record ProjectedLODGeometry
}

// Lookup resolves one (geometry, band) key. On LadderPass the selection carries
// the winning candidate; on LadderNoArtifact the row carries the typed reason
// and the full per-rung rejection log; on LadderAbsent both are zero.
func (t *LadderTable) Lookup(geometryID, bandID string) (LadderSelection, LadderOutcome) {
	if t == nil {
		return LadderSelection{}, LadderAbsent
	}
	index, ok := t.byKey[ladderKey{geometryID, bandID}]
	if !ok {
		return LadderSelection{}, LadderAbsent
	}
	row := t.Rows[index]
	if row.Status != string(LadderPass) {
		return LadderSelection{Row: row}, LadderNoArtifact
	}
	return LadderSelection{Row: row, Record: t.candidates[row.GeometryID][row.Selection]}, LadderPass
}

// Candidate returns the stored geometry for an exact (geometry, selection)
// pair. Selection keys are shared across rows, so the same simplified geometry
// is stored once even when both profiles and both bands choose it.
func (t *LadderTable) Candidate(geometryID, selection string) (ProjectedLODGeometry, bool) {
	if t == nil {
		return ProjectedLODGeometry{}, false
	}
	record, ok := t.candidates[geometryID][selection]
	return record, ok
}

// CandidateGeometries reports how many distinct geometries carry candidates.
func (t *LadderTable) CandidateGeometries() int {
	if t == nil {
		return 0
	}
	return len(t.candidates)
}

func digestHex(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

// ParseLadderTable decodes and fully verifies an artifact against the corpus
// and oracle this binary actually carries. Every failure here is a build that
// must not ship: a mismatched digest means the committed selections were judged
// against inputs that are no longer the ones in use.
func ParseLadderTable(raw []byte, corpusIdentity string, oracle SilhouetteOracle) (*LadderTable, error) {
	var artifact ladderArtifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return nil, fmt.Errorf("ladder artifact: %w", err)
	}
	if artifact.SchemaVersion != 1 || artifact.Kind != "silhouette-ladder/v1" {
		return nil, fmt.Errorf("ladder artifact: unsupported schema_version=%d kind=%q", artifact.SchemaVersion, artifact.Kind)
	}
	if artifact.SourceCorpus != corpusIdentity {
		return nil, fmt.Errorf("ladder artifact: source_corpus %s does not bind the embedded corpus %s", artifact.SourceCorpus, corpusIdentity)
	}
	if want := digestHex(silhouetteOracleJSON); artifact.OracleSHA256 != want {
		return nil, fmt.Errorf("ladder artifact: oracle_sha256 %s does not bind the embedded oracle %s", artifact.OracleSHA256, want)
	}
	if want := digestHex(ladderRecipeJSON); artifact.LadderRecipeSHA256 != want {
		return nil, fmt.Errorf("ladder artifact: ladder_recipe_sha256 %s does not bind the embedded ladder recipe %s", artifact.LadderRecipeSHA256, want)
	}
	if want := digestHex(ladderBaseRecipeJSON); artifact.BaseRecipeSHA256 != want {
		return nil, fmt.Errorf("ladder artifact: base_recipe_sha256 %s does not bind the embedded v2 recipe %s", artifact.BaseRecipeSHA256, want)
	}
	if artifact.CandidateOrder != "fine_to_coarse" {
		return nil, fmt.Errorf("ladder artifact: candidate_order=%q", artifact.CandidateOrder)
	}
	if artifact.VisibilityPolicy != oracle.VisibilityPolicy {
		return nil, fmt.Errorf("ladder artifact: visibility_policy %q does not bind the oracle policy %q", artifact.VisibilityPolicy, oracle.VisibilityPolicy)
	}
	if len(artifact.Rows) == 0 {
		return nil, fmt.Errorf("ladder artifact: no rows")
	}

	bands := map[string]bool{}
	for _, band := range oracle.Bands {
		bands[band.ID] = true
	}
	table := &LadderTable{
		SchemaVersion: artifact.SchemaVersion, Kind: artifact.Kind,
		SourceCorpus: artifact.SourceCorpus, OracleSHA256: artifact.OracleSHA256,
		LadderRecipeSHA256: artifact.LadderRecipeSHA256, BaseRecipeSHA256: artifact.BaseRecipeSHA256,
		CandidateOrder: artifact.CandidateOrder, VisibilityPolicy: artifact.VisibilityPolicy,
		Rows:       artifact.Rows,
		byKey:      make(map[ladderKey]int, len(artifact.Rows)),
		candidates: artifact.Candidates,
	}
	for i, row := range artifact.Rows {
		if !bands[row.Band] {
			return nil, fmt.Errorf("ladder artifact: row %s/%s carries unknown band %q", row.Alpha2, row.Profile, row.Band)
		}
		switch row.Status {
		case string(LadderPass):
			if row.Selection == "" {
				return nil, fmt.Errorf("ladder artifact: %s/%s/%s passes with no selection", row.Alpha2, row.Profile, row.Band)
			}
			if _, ok := artifact.Candidates[row.GeometryID][row.Selection]; !ok {
				return nil, fmt.Errorf("ladder artifact: %s/%s/%s selects %s with no stored candidate", row.Alpha2, row.Profile, row.Band, row.Selection)
			}
		case string(LadderNoArtifact):
			if row.NoArtifactReason == "" || len(row.Attempts) == 0 {
				return nil, fmt.Errorf("ladder artifact: %s/%s/%s is no_artifact without a reason and rejection log", row.Alpha2, row.Profile, row.Band)
			}
		default:
			return nil, fmt.Errorf("ladder artifact: %s/%s/%s carries unknown status %q", row.Alpha2, row.Profile, row.Band, row.Status)
		}
		key := ladderKey{row.GeometryID, row.Band}
		if previous, clash := table.byKey[key]; clash {
			prior := artifact.Rows[previous]
			if prior.Status != row.Status || prior.Selection != row.Selection {
				return nil, fmt.Errorf(
					"ladder artifact: (%s, %s) resolves to both %s/%s and %s/%s — (geometry, band) is not a functional key",
					row.GeometryID, row.Band, prior.Status, prior.Selection, row.Status, row.Selection)
			}
			continue
		}
		table.byKey[key] = i
	}
	return table, nil
}

var (
	embeddedLadderOnce  sync.Once
	embeddedLadderTable *LadderTable
	embeddedLadderErr   error
)

// EmbeddedLadderTable parses the committed artifact once, on first use. The
// artifact is ~10 MB, so the cost is deliberately kept off package init: a
// caller that never consults the ladder never pays for it.
func EmbeddedLadderTable() (*LadderTable, error) {
	embeddedLadderOnce.Do(func() {
		corpus, err := catalog.Embedded()
		if err != nil {
			embeddedLadderErr = err
			return
		}
		oracle, err := EmbeddedSilhouetteOracle()
		if err != nil {
			embeddedLadderErr = err
			return
		}
		embeddedLadderTable, embeddedLadderErr = ParseLadderTable(ladderArtifactJSON, corpus.Manifest.Identity, oracle)
	})
	return embeddedLadderTable, embeddedLadderErr
}

// EmbeddedLadderArtifactJSON returns the raw committed bytes. Verification
// tooling recomputes from these rather than re-reading the working tree, so it
// checks what the binary would actually serve.
func EmbeddedLadderArtifactJSON() []byte { return ladderArtifactJSON }
