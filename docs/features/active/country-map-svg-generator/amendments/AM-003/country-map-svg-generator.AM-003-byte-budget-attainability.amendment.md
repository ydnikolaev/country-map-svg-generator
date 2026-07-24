---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "amendment"
id: "AM-003"
epic: "country-map-svg-generator"
status: proposed
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# AM-003 — Amendment

## Trigger and discovery evidence

P2 reached twenty implementation runs, thirteen plans and twelve invalidations
without a single passing run. Every plan from PLAN-001 to PLAN-013 framed the
defect the same way — LOD output exceeds the frozen byte budget — and each
successor attacked that symptom with a different mechanism.

RUN-020 executed the PLAN-013 byte-cap fallback rung. The rung works: failing
outputs fell from 225 to 75 of 996, and AE/card and AQ/card, the two entities
PLAN-013 was written around, now select a within-budget tier. The rung is
committed at `d052f9b8dacff41719ddbc19431955c3bcbac690`. RUN-020 is interrupted
and RESULT-029 records the block.

A bounded architecture review of the remaining 75 found that the byte overrun is
the last link of a four-link chain, and that no plan had examined the first three.
Measured over the full corpus rather than the first failing entity:

| First rejection reason at the finest candidate tier | Outputs |
| --- | --- |
| `raw_deviation` | 65 |
| `hard_budget` | 8 |
| `topology` | 2 |

Only 8 of 75 begin as a byte problem. The chain is: the cheap, small-byte
candidate tiers are rejected by the deviation gate; selection therefore falls
through to the unsimplified `source` tier; and `source` is what breaks the byte
budget, by up to a factor of forty (RU/un/card is 101,605 bytes against a 2,500
maximum).

Three independent defects in the acceptance metric drive the deviation
rejections. `matchedBoundaryDeviation` at `internal/geometry/lod.go:621`:

1. Returns `+Inf` whenever the component count differs
   (`internal/geometry/lod.go:622-624`), so losing one small island is
   indistinguishable from arbitrary distortion. 25 rejections across CL, NO and
   PH are infinite.
2. Matches polygons greedily by centroid proximity and requires exact per-polygon
   ring-count equality (`internal/geometry/lod.go:630-638`), so the metric is not
   a continuous function of the geometry: as simplification changes, centroids
   move, the greedy assignment flips, and the maximum jumps. The source-tier
   search at `internal/geometry/lod.go:429-445` descends strictly through
   `nextSourceTolerance` (`internal/geometry/lod.go:543`) expecting deviation to
   fall with it. Of the 13 outputs whose search made two or more attempts, 11 are
   non-monotone — tightening the tolerance left the measured deviation unchanged
   or made it worse. CA/un/hero: attempt `7.672929` measured `7.697673`, attempt
   `3.836464` measured `8.010662`, attempt `1.918232` measured `8.010662`. A
   descent search cannot converge on a metric that does not respond to its own
   control variable.
3. Where the excess is finite it is generally large, not marginal: the median
   excess over tolerance is 17 times the `q/√2` quantization reserve. Only 10 of
   156 deviation rejections are small enough to be attributable to quantization
   noise. The ladder resolutions genuinely do not meet the tolerance for these
   entities.

Two supporting observations. First, `TestLODSpikeFullCorpus` aborts on the first
over-budget entity, so for twenty runs it reported exactly one name and every
plan scoped itself to that name; PLAN-013 asserts as established fact that "RU
card/hero hold their tiers", which the corpus has never produced and which the
test never reached. Second, three of the four hard-coded tier expectations in
`TestLODAQRUProjectionAlignedCheckpoint` assert selections the corpus does not
make, for the same masked reason.

## Affected scope and authority

This is a P2 macro-HOW amendment to the acceptance metric and the ladder
calibration accepted under DEC-005. The product owner is the decision authority.

It does not change P2 WHAT. The specification never makes a tier label normative
and never fixes byte budgets numerically: REQ-2 states only that preset defaults
specify byte budgets. REQ-3 (ring validity, winding, non-empty identity, P1
protected features), AC-1 (recognizable, balanced, artifact-free silhouettes) and
AC-2 (aspect, containment, byte-equivalent reruns) are unchanged and remain the
obligations any replacement metric must serve.

It does not change the layered architecture. DEC-004 rejected runtime RDP on
measured topology, deviation and byte contradictions and chose precomputed tiers
with a pure-Go runtime; that choice is unaffected and is not reopened. The
compact/standard/source ladder stays. What changes is how a candidate is judged
against the tolerance, and how the ladder resolutions are calibrated once that
judgement is sound.

P3, P4 and P5 are unaffected in outcome and boundary. P3 continues to receive
selected tier, requested tier, fallback reasons, deviation and fidelity metrics,
path bytes, transform and optional markers; the meaning of the deviation metric it
receives changes, its shape does not.

## Prior accepted truth

DEC-005 (`q-aware source fallback`) accepted a bounded source-tier fallback that
reserves the maximum `q = 0.01` grid displacement (`q/√2`) from the resolved
symmetric visual tolerance and "searches monotonically from the coarsest allowed
simplification toward finer candidates without ever exceeding the remaining
tolerance". The implementation reserves the headroom correctly at
`internal/geometry/lod.go:423`. The monotonic search the decision assumes is not
achievable against the metric as implemented, and that assumption was never
tested.

The frozen quantities set by DEC-005 and its neighbours — centered LAEA
projection, `q = 0.01`, the 29-resolution ladder `[512…8]`, weighted Visvalingam
`0.7`, the `1.536` / `7.68` deviation tolerances, path caps `2200/7500`,
complete-file maxima `2500/8000`, the serializer, the silhouette oracle and the
v1 recipe — were all treated as immutable by PLAN-001 through PLAN-013.

PLAN-013 is accepted at
`0efd5be8f545f28edd11262d39cd0f96bf765f2698db93fb932ec36187295dbe` and remains
accepted; it is not invalidated by this amendment, because its rung is correct
work that a successor builds on.

## Accepted change

P2 replaces the boundary-deviation acceptance metric with one that is well-posed
for the search that consumes it. The replacement must satisfy three properties,
each of which is a testable obligation rather than an implementation choice:

1. **Total.** It returns a finite value for any two geometries being compared.
   Component or ring count mismatch is reported as a named, graded structural
   result, not as `+Inf`. Loss of a protected or mandatory feature remains a hard
   rejection under REQ-3, decided by the protection check rather than by an
   infinite distance.
2. **Stable.** Correspondence between source and candidate parts does not depend
   on a greedy centroid ordering or on exact ring-count equality. Equal inputs
   produce an equal value independent of iteration order, preserving REQ-8.
3. **Monotone in the search variable.** Reducing the simplification tolerance
   does not increase the measured deviation. This is the property DEC-005 already
   assumed and is what makes the source-tier descent able to terminate.

The byte budgets `2200/7500` and `2500/8000` are NOT changed by this amendment,
and no per-country exception is introduced — country literals and branches remain
forbidden. The expectation is that once the cheap tiers stop being falsely
rejected there is no reason to fall through to `source`, and the byte pressure
resolves as a consequence rather than as a target.

Ladder resolution calibration is reopened as a second, dependent step: once the
metric is sound, the 29-resolution ladder is re-measured against the corrected
tolerance and adjusted if entities still cannot reach a within-budget tier at
mandatory fidelity. Any such adjustment is a ladder change requiring a rebuild and
a fresh owner-approved contact sheet.

No estimate is offered here for how many of the 75 outputs this recovers. An
estimate produced before the metric is corrected and re-measured would be a guess;
the measurement is part of the work.

## Downstream impact map

- P2 planning context advances to include this amendment. The next plan is
  authored against the corrected diagnosis, not as a successor to PLAN-013's
  framing.
- The silhouette oracle and the v1 recipe are unchanged by the metric replacement
  itself. If step two adjusts ladder resolutions, both are rebuilt and both owner
  contact sheets are regenerated and re-approved.
- `TestLODAQRUProjectionAlignedCheckpoint` carries three tier expectations the
  corpus does not produce (AQ/card, AQ/hero and RU/card). They are corrected
  against measured selection as part of the work, not preserved.
- `TestLODSpikeFullCorpus` and any other population gate must report the full
  failure set with counts and distribution instead of aborting on the first
  entity. Fail-fast on a corpus sweep is what concealed the scale of this defect
  for twenty runs.
- P3 receives the same provenance fields with a changed metric meaning; its
  contract text needs no change, but its planning context must carry this
  amendment so its assertions are written against the new semantics.
- P4 full-catalog and browser evidence is unaffected in shape.

## Validation delta

- A teeth test proving the metric is total: a component-count mismatch produces a
  finite graded result, and a protected-feature loss still fails through the
  protection check.
- A teeth test proving stability: identical inputs give identical values under
  permuted component order, and correspondence does not change under a centroid
  perturbation that would flip a greedy assignment.
- A property test proving monotonicity across the frozen ladder: for a
  representative sample spanning continents, archipelagos, fjord coasts and polar
  entities, deviation is non-increasing as the simplification tolerance falls.
  This is the guard that would have caught the defect, and it is the one gate this
  amendment most requires.
- The full-corpus sweep reports every failing output with tier, fallback reasons,
  deviation and bytes, and fails once with the count and distribution.
- The 27 entities currently unreachable — BD BR BS CA CI CL DE DK EG FR GB GL GR
  GW HR ID IE IS IT KR NO PH RU SN TR VE VN — are named in the evidence as the
  regression population, with RU/un's ring self-intersection at both table
  resolutions tracked separately as a topology defect rather than a deviation one.

## Migration rollout and rollback

There are no runtime instances and no published artifacts. The corrected metric
changes selection outcomes across the corpus, so every representative and
full-catalog sheet is regenerated and re-approved by the owner before P2 closes.

Rollback restores the metric and, if changed, the ladder to their pre-amendment
bytes; the PLAN-013 rung at `d052f9b` is independent of this amendment and is
retained under rollback.

## Application and status

This immutable proposal records the intended change. Current lifecycle state and
application authority are carried by the epic tracker and transition receipts, not
by mutating this body after publication.

Three portable harness findings were raised from this episode and are recorded in
the project feedback queue rather than here: the absence of a lifecycle brake on
repeated non-convergence (`FND-1B3065946FF8`), verification briefs that state the
conclusion under test and are confirmed rather than testing it
(`FND-8386048ACA3C`), and population gates that abort on first failure
(`FND-A50ADF64186A`).

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
