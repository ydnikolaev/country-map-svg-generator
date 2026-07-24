---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-012"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-012 — Topology is judged after canonicalization, on both sides

## Context

T5's contact-sheet preparation was blocked by `WKI-C83B0B9EB4EE`: Russia
rendered as a crude single-part blob spending 508 of its 7500 hero bytes, while
Italy at the same band spent 6829. Russia was the only anomaly of its kind — of
the 14 committed rows selecting a rung at or below 64, the other 12 sat at 68–98%
of their byte budget, coarse because they had hit the cap.

The committed artifact could not explain it: per-attempt rejection logs are
stored only on `no_artifact` rows, so a row that passes at a coarse rung records
nothing about why the finer rungs lost. `lodbuild -ladder-explain` was added to
replay the search for one entity and print every rejected rung. It reports what
the build did rather than a reimplementation, and it is exact — `build.mjs`
simplifies one feature at a time, so sweeping a single entity reproduces the
catalog sweep's geometry per rung.

Three measurements followed, each checked rather than inferred:

1. **Only the `un` profile is affected.** RU/de_facto is topologically clean at
   every rung and settles at 80 with 5 parts and IoU 0.9617.
2. **The source is not at fault.** Both profiles validate clean at 209 and 214
   parts, in projection and as the identity record. There is no antimeridian
   seam: the longest segment in the largest ring is unremarkable against its
   neighbours. The self-crossing is introduced by simplification, and it is the
   *same* crossing at every rung from 512 down to 32.
3. **The pipeline's own canonicalization removes it.** The defect measures
   ~1e-4 in fitted units against an output grid of q=0.01 — four orders of
   magnitude coarser. Replaying rungs 512, 128 and 64 through canonicalization:
   topology broken before, clean after, on-grid and contained, at 39, 9 and 3
   parts.

So the acceptance predicate rejected candidates for a defect that the pipeline
erases two steps later. It was validating topology at a precision no emitted path
ever has.

## Decision

**Topology is judged on the canonical geometry, not on the pre-canonical fitted
coordinates — in the oracle and in the runtime alike.**

1. `EvaluateSilhouetteCandidate` no longer rejects a candidate on a
   pre-canonical topology check. The authoritative check remains inside the
   phase loop, on the canonical geometry.
2. The runtime's derived-candidate path applies the same boundary for a ladder
   candidate. `canonicalizeShared` validates topology on its output, so the
   post-canonical guarantee is unchanged.
3. **The two sides must agree.** A ladder selected under one predicate and served
   under another is the failure class this spec exists to end; the runtime change
   is not optional trimming, it is half of the decision. Fixing only the build
   left Russia committed at rung 80 and falling back to source at 101605 bytes
   against a 2500 ceiling.
4. The DEC-005 source path keeps its pre-canonical check unchanged. It is not
   judged by the oracle and makes no such promise.

No threshold, byte budget, tolerance, rung set or quantization value changes.
`q=0.01` is untouched; this decision changes *where* topology is judged, not what
counts as valid.

## Options considered

1. **Judge topology after canonicalization on both sides (chosen).**
2. **Add a tolerance to the intersection test.** Treat a crossing below some
   epsilon as absent.
3. **Leave it.** Accept Russia at rung 24.

## Rejected alternatives

**Option 2** was rejected because it introduces a new frozen constant that has to
be justified, calibrated and defended for every entity, and it would weaken the
topology guarantee everywhere in exchange for a case that needs no weakening at
all. The crossing found is a *strict* crossing, not a near-touch, so a tolerance
would have to be large enough to swallow real defects. Judging at the output's
own precision is exact rather than approximate, and needs no new number.

**Option 3** was rejected by the owner: a top-200 country rendering as a blob at
7% of its byte budget is not shippable, and it would have been the first thing
visible on the T5 contact sheet.

## Consequences and tradeoffs

The catalog effect is surgical, and this is the load-bearing measurement:
**exactly 2 of 996 rows changed**, both of them Russia's `un` profile.

| Row | Before | After |
| --- | --- | --- |
| RU/un/compact | rung 24, 437 B, 1 part, IoU 0.9212 | rung 80, 1759 B, 4 parts, IoU 0.9613 |
| RU/un/standard | rung 24, 508 B, 1 part, IoU 0.9204 | rung 224, 6577 B, 16 parts, IoU 0.9722 |

Everything else is unchanged, which independently confirms the diagnosis that
Russia was the sole victim of the misplaced check.

- **Catalog totals hold**: 993 pass, 3 no_artifact, 6 identity rungs, 283
  candidate geometries. DEC-009's ladder rules and DEC-010's Croatia outcome keep
  their evidence exactly.
- **The artifact and its whole binding chain were rebuilt.** Editing
  `silhouette.go` changed the rasterizer digest the oracle records, which changed
  the oracle digest the ladder recipe binds, which changed the recipe digest the
  artifact and every stored candidate carry. That cascade is the provenance chain
  working as designed, and it forces a rebuild rather than allowing a silent
  drift.
- Rendered through the shipped path afterwards: 993 rendered, 3 typed
  no_artifact, 0 errors, 0 source fallbacks, 0 over the band cap, 756
  byte-identical to the committed row.
- A candidate whose topology is genuinely broken now fails later — inside the
  phase loop rather than before it — so it costs more phases before rejection.
  Measured build time is unchanged in practice.

## Affected artifacts and owners

- `internal/geometry/silhouette.go` — the oracle's acceptance predicate.
- `internal/geometry/lod.go` — the runtime's matching boundary for a ladder
  candidate.
- `internal/geometry/lod/silhouette-oracle.v1.json`,
  `ladder.recipe.json`, `ladder.artifact.json` — the rebuilt binding chain.
- `internal/geometry/cmd/lodbuild` — `-ladder-explain`, the diagnostic that made
  this findable.
- P2 (`geometry-pipeline`) — owns all of the above.
- `WKI-C83B0B9EB4EE` — the release blocker this closes.
- Accountable owner: product owner (accepting principal).

## Validation and revisit trigger

Validated by the existing gates re-run against the rebuilt artifact, all green:
the offline VAL-6 recomputation of every committed row, the Node determinism
rebuild proving the artifact reproduces byte-identically, and the shipped-catalog
gate proving `Generate()` serves the ladder within every band cap. The
rasterizer-identity test binds `silhouette.go` to the oracle, so any future edit
to the predicate forces the same rebuild rather than allowing drift.

Revisit if a future corpus or Mapshaper pin produces a candidate whose
pre-canonical damage is *not* erased by canonicalization and which therefore
survives to the post-canonical check — that is the case this decision routes
correctly by design, and it should be observed rather than assumed impossible.

## Supersession

None. DEC-012 relocates where DEC-006's topology obligation is evaluated; it does
not alter the oracle's metrics, DEC-007's protected-visibility policy, DEC-009's
ladder rules, DEC-010's Croatia outcome or DEC-011's framing decision.
