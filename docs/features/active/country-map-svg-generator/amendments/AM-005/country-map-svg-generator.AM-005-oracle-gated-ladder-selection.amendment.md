---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "amendment"
id: "AM-005"
epic: "country-map-svg-generator"
status: proposed
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# AM-005 — Amendment

## Trigger and discovery evidence

DEC-006 replaced raw full-coastline symmetric deviation with a target-scale
raster silhouette oracle as the acceptance boundary for scale-calibrated derived
candidates. The production selection path never received that change. Twenty runs
and thirteen plans then tried to make the superseded predicate pass.

Four facts, each verified against the code rather than inferred:

1. `internal/geometry/lod.go` gates every derived tier candidate on
   `matchedBoundaryDeviation` at lines 245-246 and re-gates it after grid-phase
   canonicalization at 518-520. It contains no reference to the silhouette
   oracle, to coverage, or to IoU.
2. The oracle exists, is calibrated, frozen and owner-approved under DEC-008, and
   its own primitive says what the rule is. `EvaluateSilhouetteCandidate`
   (`internal/geometry/silhouette.go:214-217`): "automatic acceptance is
   exclusively oracle and budget driven; DEC-005 raw deviation is intentionally
   not consulted." It is consumed only by the maintainer representative path.
3. `restoreRequiredComponents` is called at `internal/geometry/lod.go:239`,
   immediately before the deviation gate, injecting full-detail source components
   into a simplified candidate. DEC-006's rejected alternatives name exactly this:
   "Restoring every retained source component reverses the accepted visibility
   policy and makes remote islands dominate bytes despite having no visible role
   in a card."
4. The runtime table is both pre-DEC-006 in shape and not wired.
   `internal/geometry/lod/v1.recipe.json` commits two global single resolutions —
   `compact: 64`, `standard: 256` — for all 283 geometries, where DEC-006
   specifies a per-geometry fine-to-coarse search that "selects the finest
   candidate that passes every generic predicate". And
   `internal/geometry/pipeline.go:11` reads `var publishedLODTable *LODTable`,
   nil, so the shipped `Generate()` is source-tier-only.

### The T0 measurement

Before authoring this amendment, the full catalog was swept under the DEC-006
predicate: every geometry × band, fine-to-coarse over the complete frozen 29-rung
v2 ladder, evaluated by `EvaluateSilhouetteCandidate` unmodified, with the frozen
thresholds untouched and no per-country branch. Determinism confirmed by two
independent runs producing byte-identical output.

Bound identities: source corpus
`sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`,
silhouette oracle
`f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`,
v2 recipe `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`.

Result: **986 of 996 rows pass; 10 fail, across 3 entities.** Under the current
production predicate, 75 of 996 fail. The predicate was the defect.

The 10 failures, with their verified mechanisms:

- **SH — 4 rows, both bands, both profiles.** Saint Helena's reference carries 4
  components separated by thousands of kilometres of the South Atlantic.
  Mapshaper's `resolution=N` sizes its simplification grid from the feature's own
  bounding box, so at every one of the 29 rungs — including the finest, 512 —
  every island falls below grid resolution and the geometry collapses to a single
  surviving polygon. This is one structural gap repeated 29 times, not 29
  independent search failures. The **unsimplified** reference geometry passes
  cleanly at both bands: IoU 1.0, 637 of 2200 path bytes at compact, 630 of 7500
  at standard. Saint Helena does not need simplifying; the ladder has no rung that
  declines to simplify.
- **UM — 4 rows.** United States Minor Outlying Islands, 12 Pacific and Caribbean
  atolls, same bbox-relative mechanism. The unsimplified reference passes at
  standard (791 of 7500). At compact it fails independently and earlier, with a
  ring collapse during `q=0.01` canonical quantization — no ladder rung, however
  fine, addresses that. UM/compact is a genuine infeasibility under the current
  quantization contract, and is the one open decision in this amendment.
- **HR — 2 rows, compact only.** Croatia's 25-part archipelago degrades gracefully
  (25→24→…→9 parts, no collapse). At resolution 96 it serializes to 2783 bytes
  against the 2200 path cap; at resolution 80 to **2209 bytes — over by 9**; and at
  the next declared rung, 64, source component 2 is lost with contribution 110
  against a threshold of exactly 110. The frozen ladder declares nothing between
  80 and 64. HR/standard passes; the 7500-byte cap absorbs it.

### Recorded fragility

Byte headroom on the 986 passing rows is real but thin. Compact fills the cap to
72% on average, worst VE/de_facto at **2199 of 2200**. Standard averages 50%,
worst PH/un at 7346 of 7500. A future serializer, projection or precision change
will move those rows first. This is recorded as a known property, not a defect.

## Affected scope and authority

A P2 macro-HOW amendment implementing an already-accepted decision, owned by the
product owner. It affects P2's derived-candidate acceptance path, its build
artifact, and its runtime wiring, plus the P2 validation allocation. P3 is affected
by context and by the provenance fields it receives; no P3 contract text changes.

P2's WHAT is unchanged. No REQ, INV or AC row is touched. The specification makes
no tier label normative and fixes no byte budget numerically. Byte budgets
`2200/7500` path and `2500/8000` complete file are unchanged — DEC-006 rejected
raising them explicitly, and nothing in the T0 measurement argues for it.

DEC-004's precomputed-tier architecture and pure-Go offline runtime are not
reopened; this amendment restores the division of labour DEC-004 and DEC-006
already specified. DEC-005 is not reopened: raw deviation remains the binding gate
for an explicit `source` or custom-quality request, which is where DEC-006 left
it. DEC-007's visibility policy is unchanged and is a gate in the new path.
DEC-008's frozen thresholds are unchanged and reused as-is; its approved digests
are re-approved because the ladder recipe changes, per §Accepted change item 5.

This amendment supersedes AM-003 and AM-004 in full. Both are `applied` and
unverified, blocked by `ADV-003` and `ADV-004`; both were authored from diagnoses
that independent review refuted. The VAL-6 row they placed in the P2
specification is replaced.

## Prior accepted truth

- AM-003 `d910385478f3d9eea6c3f3755aa769928da0f7ce0c82c1fedaed68d4835bdffd`,
  blocked by `ADV-003`. Its defect #1 named an unreachable branch; its magnitude
  claims were clamped first-crossing values, not measurements.
- AM-004 `02fd14952364b0ca7ce9bb0a98d20a0c5b670e94f18e639030545aac09aca471`,
  blocked by `ADV-004`. Its "predicate used as a minimization objective" describes
  a mechanism the code does not have: the deviation magnitude enters a boolean and
  a log string, and `nextSourceTolerance`
  (`internal/geometry/lod.go:543-548`) is a function of tolerances only.
- P2 specification currently at
  `e1dbb85ebc5fca0de656adff86acce9d5208b70014f6db659aa0cb01e658a853`, carrying
  AM-004's VAL-6, which `ADV-004` found untestable: its numeric band is
  undeclared, and "the accepted candidate" does not exist in a failing output.
- One kernel of AM-004 survives and is retained verbatim below:
  `matchedBoundaryDeviation` is exact at or below its limit and short-circuits
  above it, so a first-crossing value is never a measurement.

## Accepted change

**1. Derived acceptance is the oracle.** Automatic acceptance of a scale-calibrated
derived candidate is exclusively the frozen silhouette oracle, DEC-007 protected
visibility, and the byte caps. The derived path does not consult raw
full-coastline symmetric deviation. Raw deviation remains binding for an explicit
`source` or custom-quality request under DEC-005, including its typed failure when
a caller combines an incompatible fidelity request with a byte maximum.

**2. Selection moves to maintainer build time.** Candidate search is per geometry
× band, fine-to-coarse over the frozen ladder, first passing candidate wins, with
source polygon/ring order and content digest as the stable tie-break. The result
is a committed ladder artifact carrying, per row: selected resolution, IoU,
recall, visibility and omission provenance, path bytes, complete-file estimate,
and the oracle, recipe and corpus digests it was produced under. The runtime
becomes a lookup that verifies — bind identity, fit, canonicalize, re-validate
topology, protected anchors, containment and the byte cap, serialize — and never
a search that measures. Effective fitted scale selects the band; above the
standard band, and for any explicit non-auto quality, the DEC-005 source path is
used unchanged.

**3. Component restoration is removed from the derived path.** DEC-006 rejected
the mechanism. The 8 byte-cap first-rejections in the pre-T0 population were
restoration-inflated: a candidate built at resolution 64 cannot serialize to
14,015 bytes from Mapshaper output alone.

**4. The global two-tier artifact is superseded.** `v1.recipe.json`'s tier block
is retired; its Mapshaper pin moves into the ladder recipe. "Tier" stops naming a
global geometry set and names this entity's selected candidate for this band. The
band names and their `240`/`700` fitted-scale boundaries survive, because they
carry the byte budgets. `publishedLODTable` is initialized from the committed
artifact.

**5. The ladder gains rungs; thresholds gain nothing.** Two additions, both
following from T0:

- an **identity rung** that declines to simplify, evaluated like any other
  candidate against the same frozen gates. T0 proves SH passes at IoU 1.0 and 637
  bytes and UM/standard at 791 bytes with zero simplification, and that no
  existing rung can reach them because Mapshaper's resolution is bbox-relative;
- **intermediate rungs between 80 and 64**, to close HR/compact's 9-byte gap
  without dropping a component whose contribution sits exactly at the retention
  threshold.

No threshold, tolerance, weighting, byte cap or budget changes. Entity- and
preset-specific thresholds remain forbidden. Because the ladder recipe changes,
DEC-008's approved digests are re-approved on the regenerated sheets; the
thresholds those sheets approved are carried through unchanged.

**6. UM/compact is an open owner decision.** T0 establishes that UM has no compact
candidate under any ladder rung, and none even without simplification, because the
geometry does not survive `q=0.01` canonical quantization at card fitted scale.
This amendment does not choose. The options, and what each costs:

- **Typed refusal.** UM has no `card` silhouette; P3 returns a typed, explainable
  failure for that entity/preset pair and the catalog documents the gap. Cheapest,
  honest, and leaves one hole in a "complete catalog" claim.
- **Group anchor via the existing override seam.** DEC-007 item 6 already provides
  a versioned data override for preserving a declared identity feature; UM's atolls
  are exactly the shape it was written for. Keeps the catalog complete, requires
  owner visual review, and is data rather than a code branch — so it does not
  violate AC-3.
- **Revisit the quantization contract for sub-scale multi-component geometries.**
  General, not per-country, and therefore the most invasive: it touches `q=0.01`,
  which every other row depends on, and would re-open determinism and byte
  guarantees across the catalog.

Raising the byte budgets and relaxing an oracle threshold are not options here.
DEC-006 rejected both by name, and the T0 measurement gives no evidence for
either — UM/compact fails on topology, not on bytes or on similarity.

**7. Reporting honesty is normative.** A deviation figure carried in provenance,
evidence or a plan is either at or below its limit, where the metric is exact, or
is produced by an unclamped path. A first-crossing value is never reported as a
measurement. Population gates report every failing entity with counts and
distribution; a sweep that stops at its first failure does not satisfy them.

## Downstream impact map

- The P2 specification's VAL-6 row is replaced per the validation delta below.
- P2's context registry advances so the next plan is authored against this
  amendment and the T0 matrix, not against AM-003's or AM-004's framing.
- CTR-002 carries selected band, selected resolution, IoU, recall, omission and
  visibility provenance, and the oracle and recipe digests. P3 consumes these for
  explainability; its contract text needs no change.
- The representative artifacts and both owner contact sheets are regenerated,
  because selection outcomes change across the corpus. DEC-008 re-approval is
  required before P2 closes.
- The silhouette oracle itself is not rebuilt. It is built, frozen, owner-approved
  and already correctly exercised by the representative path.
- P4's byte-distribution evidence inherits the T0 headroom figures as its
  starting expectation, including the VE and PH near-cap rows.

## Validation delta

VAL-6 is replaced by three obligations:

- **Oracle-gated derived acceptance.** The derived path's acceptance is the oracle,
  DEC-007 visibility and the byte caps. A mutation test reddens both on bypassing
  the oracle and on reintroducing a raw-deviation gate to the derived path. A
  second reddens on reintroducing component restoration.
- **Committed-artifact re-verification.** Every committed ladder row is
  re-verified offline, in pure Go, without Node or network, against the frozen
  thresholds and bound digests: topology, protected visibility, IoU, recall and
  byte caps recomputed and required to agree with the committed provenance.
  Mutations that must redden: a perturbed coordinate, a changed threshold, a stale
  corpus, oracle or recipe identity, reordered candidates, and a suppressed
  omission.
- **Explicit-source integrity.** The DEC-005 source path keeps raw deviation as
  its gate, and skipping its final validation fails.

Population guards under VAL-1 through VAL-3 report every failing entity with
counts and distribution. `t.Fatalf` inside a corpus loop does not satisfy them;
`internal/geometry/cmd/lodbuild/main_test.go:88` is the instance that concealed
the scale of this defect for twenty runs.

VAL-6's owner/due column reads `P2 / W2 complete`, matching its siblings.

## Migration rollout and rollback

There are no runtime instances and no published artifacts. Selection outcomes
change across the corpus, so both owner sheets are regenerated and re-approved
before P2 closes.

Rollback restores the P2 specification to its pre-AM-003 bytes
`4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`, the
`v1.recipe.json` tier block, and the current `lod.go` selection path. The RUN-020
byte-cap rung at `d052f9b8dacff41719ddbc19431955c3bcbac690` is retained under
rollback as a defensive terminal check, though item 2 moves byte evaluation to
build time and the rung stops being load-bearing.

## Application and status

This immutable proposal records the intended change. Current lifecycle state and
application authority are carried by the epic tracker and transition receipts, not
by mutating this body after publication.

Five portable harness findings from this episode are recorded in the project
feedback queue. Two of them describe how this took twenty runs:
`FND-1B3065946FF8`, no lifecycle brake on repeated non-convergence — every
invalidate/plan/run/interrupt cycle was individually legal, so nothing ever forced
a review of the shared diagnosis; and `FND-A50ADF64186A`, a population gate that
aborts on first failure — which is why every plan saw one failing country instead
of 75. The others: `FND-8386048ACA3C`, a verification brief that carries its own
conclusion is confirmed rather than testing it; `FND-CC7A8AAB239A`, an applied
amendment that fails independent review has no corrective transition, which is why
AM-003 and AM-004 remain `applied`; and `FND-D35B92BE8504`, no native receipt for
a required product-owner approval checkpoint.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
