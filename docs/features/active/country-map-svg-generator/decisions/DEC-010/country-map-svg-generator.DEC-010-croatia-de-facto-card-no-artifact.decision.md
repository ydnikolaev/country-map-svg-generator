---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-010"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-010 — Croatia de_facto card renders no artifact

## Context

DEC-009 established the generic rule for a band that cannot be drawn: when no
rung on the resolution ladder — including the identity rung tried last — can
satisfy topology, protected visibility and the frozen byte caps at the actual
fitted scale, the pipeline emits a typed no-artifact outcome rather than a
degraded silhouette. It also fixed the ladder itself: identity last, plus one
intermediate rung at 79.63 between the declared 80 and 64.

DEC-009 was accepted on the T0 catalog feasibility measurement, which reported
Croatia's gap as nine bytes and expected the 79.63 rung to close it. Building the
committed ladder artifact (`internal/geometry/lod/ladder.artifact.json`) showed
that those figures belonged to Croatia's **un** profile. The **de_facto** profile
is a different geometry (`geo-ab96f1c1…`) and it is harder. Its committed
per-attempt rejection log is unambiguous:

| Rung | Outcome |
| --- | --- |
| 80 | `bytes=2215 (cap 2200)` |
| 79.63 | `bytes=2205 (cap 2200)` |
| 64 | `protected_feature_loss entity=HR field=visible_component: source component 2 contribution=111 threshold=110 was lost` |
| 48 … 8 | protected component loss, progressively worse |
| identity | `bytes=18190 (cap 2200)` |

The byte-cap crossing and the protected-component-loss crossing happen at the
same simplification step, verified to `.001` resolution precision. There is no
rung at any granularity that bridges them: anything fine enough to keep the
disputed-border component (contribution 111, one unit over the 110 retention
threshold) is over the 2200-byte card path cap, and anything coarse enough to fit
the cap has already dropped that component.

Croatia is inside the top ~200 set the product requirement names, which is why
this case is escalated to an explicit owner decision instead of being absorbed
silently by DEC-009's generic rule. The same mechanism already produces
no-artifact for UM/un/compact and UM/de_facto/compact, which are not in that set.

## Decision

**Croatia's de_facto profile at the compact band emits the DEC-009 typed
no-artifact outcome.** No override, no threshold change, no byte-budget change,
no country literal and no per-entity branch anywhere in the pipeline. The outcome
is produced by the generic rule alone; DEC-010 records that the owner reviewed
this specific instance of it and accepted the result.

Croatia therefore renders in three of its four profile × band slots: un/card,
un/hero and de_facto/hero. The consuming site treats the missing de_facto card
the same way it treats any absent decorative card.

The committed final catalog state is **993 pass / 3 no_artifact** over 996 rows.

## Options considered

1. **Accept the typed no-artifact outcome (chosen).** Costs one decorative card
   on one profile of one country. Requires no code, no data, no threshold move
   and no new mechanism — DEC-009's rule already produces it. Keeps the pipeline
   free of entity-specific knowledge, which INV-1 and the DEC-006 predicate both
   depend on.
2. **Group anchor via the DEC-007 override seam.** A data override preserving the
   disputed-border component would keep the card. It is data rather than code, so
   it does not violate the no-country-literal constraint, but it requires owner
   visual review of the resulting silhouette and adds a permanent per-entity data
   row that must be re-validated on every corpus bump.
3. **Relax the card byte budget or the 110 contribution threshold.** Either would
   let some rung pass.

## Rejected alternatives

**Option 2** was rejected for now, not on principle. It remains available through
the existing DEC-007 seam and needs no further decision authority to adopt later;
the owner declined to spend visual-review effort on a single decorative card when
the honest absence costs nothing structural. If the site later shows that the
missing card is visible enough to matter, adopting option 2 is a data change
under DEC-007, not a reopening of DEC-010.

**Option 3** was rejected outright. The 2200/7500 path and 2500/8000 file budgets
are frozen by DEC-006 and are the reason the product's optimization claim is
meaningful; raising one to rescue one card inverts that. The 110 contribution
threshold is frozen by DEC-008 and a change would require re-running its owner
approval across the whole catalog, not just Croatia. Moving a global threshold to
fix a single entity is precisely the per-entity special case this epic has
refused everywhere else, wearing a generic costume.

## Consequences and tradeoffs

- One top-200 country is missing one decorative card. This is a visible product
  gap and is accepted knowingly rather than hidden.
- The pipeline stays entity-agnostic. No selection code, threshold, budget or
  corpus byte changes as a result of this decision.
- The no-artifact path now has a top-200 exemplar, so its typed outcome must be
  handled properly by every downstream consumer rather than treated as an
  unreachable edge — P3's CLI and P4's site both have to render absence
  deliberately.
- Reversal is cheap and pre-authorized: option 2 through the DEC-007 seam.

## Affected artifacts and owners

- P2 (`geometry-pipeline`) — owns the ladder artifact and the no-artifact typed
  outcome; carries the committed rejection log that evidences this decision.
- P3 (`generator-cli`) — must surface the typed no-artifact outcome as a first
  class result, not an error, for HR/de_facto/compact.
- P4 (`site-delivery`) — must render a missing decorative card without breaking
  layout for a top-200 entity.
- `internal/geometry/lod/ladder.artifact.json` — the evidence; rows for
  HR/de_facto/compact, UM/un/compact and UM/de_facto/compact carry
  `Status: no_artifact` with the full per-rung rejection log.
- Accountable owner: product owner (accepting principal).

## Validation and revisit trigger

Validation is the committed artifact itself: the HR/de_facto/compact row must
remain `no_artifact` with the rejection log above, and the catalog totals must
remain 993 pass / 3 no_artifact. The VAL-6 offline ladder gate recomputes every
committed row against the frozen thresholds, so a silent change to this outcome
reddens the project ceiling.

Revisit when any of: the P1 corpus is bumped and Croatia's de_facto geometry
changes; the DEC-006 byte budgets or the DEC-008 contribution threshold are
reopened for an unrelated reason; or the site shows the missing card is a
material product defect, in which case option 2 is adopted through DEC-007
without reopening this record.

## Supersession

None. DEC-010 refines the application of DEC-009 to one measured instance and
supersedes no predecessor. DEC-009 remains the generic rule; DEC-010 is the owner
acceptance of its outcome for a top-200 entity.
