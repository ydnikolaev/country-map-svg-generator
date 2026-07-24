# P2 checkpoint — resume point for a fresh session

Last updated at commit `da65ba1` (T2 complete). Read this first; it lets a new
session resume without the originating chat.

## The one-line status

The 20-run P2 failure is diagnosed, the fix is designed and independently
verified, the ladder artifact is built (T1) and now embedded and independently
recomputed row by row (T2). Both prior open items are closed: the Croatia card is
an accepted owner decision (DEC-010), and the mate fence is confirmed unliftable
with the owner's standing instruction to keep building outside the governed
wrapper. **Next step is T3.**

## What was wrong (the root cause, confirmed)

DEC-006 (accepted) replaced raw full-coastline deviation with a target-scale
raster **silhouette oracle** as the acceptance boundary for derived candidates.
The production selection path in `internal/geometry/lod.go` never received that
change — it still gates on `matchedBoundaryDeviation` (raw deviation) at
`lod.go:245-246` and `518-520`, references the oracle nowhere, and injects
full-detail source components via `restoreRequiredComponents` (`lod.go:239`), the
exact mechanism DEC-006 rejected. The shipped `Generate()` is still source-only
because `publishedLODTable` is nil (`pipeline.go:11`). Twenty runs tried to make
the superseded predicate pass. Under the correct oracle predicate the catalog
goes from **75 failing outputs to 3**.

Three prior coordinator diagnoses (quantization accounting; a `+Inf` mis-count
branch; a predicate-as-objective) were each refuted by independent review or by
measurement. The lesson, recorded in `docs/harness-backlog.yaml`
(`FND-8386048ACA3C`): on this problem, a hypothesis without a full-corpus
measurement is worthless. Do not re-diagnose from a sample.

## Governance state

- **DEC-009** accepted — identity rung tried last; a band that cannot satisfy
  topology at fitted scale produces a typed no-artifact outcome (generic rule, no
  country literal); intermediate rungs added between 80 and 64.
- **DEC-010** accepted (`45385fb`) — the owner accepts the typed no-artifact
  outcome for HR/de_facto/compact. See the closed decision below.
- **AM-005** verified (`ADV-006` pass, independently recomputed) — the oracle-gate
  fix. Supersedes AM-003/AM-004 in text.
- **AM-003, AM-004** stuck in `applied`, fencing P2 and P3. **Confirmed
  unliftable at the source**, not merely observed — see the next section.

## The mate fence — verified at source, not inferred

Re-probed under mate **v1.4.3** (the finding `FND-1B950D899CE4` was filed against
v1.4.1, so this is a fresh confirmation, not a citation):

| Probe | Result |
| --- | --- |
| `amendment reject AM-003` | `illegal amendment transition reject from applied` |
| `amendment verify AM-003 --input advisor-receipt=ADV-003#c5b9347f…` | `advisor receipt does not pass and bind the current amendment closure` |
| `spec transition P2 amendment-resolve` | `spec P2 is fenced by impacting amendments AM-003,AM-004` |

Read in the mate SSOT (`~/Developer/projects/mate`):

- `internal/workdocs/assets/bundle-v1/registries/state-machines.yaml:65-70` — the
  amendment machine has **exactly one** edge out of `applied` (`verify`), and
  `reject` is legal only from `[proposed, impacting]`. A failed applied amendment
  is a lifecycle sink.
- `internal/work/amendment.go:283` — the blocking set is
  `{impacting, accepted, applied}`, so `applied` fences by construction.
- The fence is enforced at `internal/work/plan_accept.go:98`,
  `internal/work/run_start.go:95`, `internal/work/run_finish.go:140` and
  `internal/work/transition.go:148`. **Only `cancel` and `invalidate` are exempt.**

So `plan invalidate` is legal (dry-run green: PLAN-013, authority
`run-result=RUN-020-RESULT#b440805079f0242620e9529c75713c8db0b7637dbd8caaf1fae44cc5a4447289`,
spec revision 123 → 124) but `plan supersede`/`accept` and `run start` are not.
There is no governed route to a P2 run. Verifying AM-003 to escape this would be
fabricated evidence and must not be done.

**Owner instruction (this session): continue T2+ as ordinary engineering commits
outside the governed run wrapper**, with the reconciliation debt tracked. The
mate fix is the owner's separate work.

Governed operations that *do* still work under the fence and should keep being
used: `mate decision scaffold/validate/accept` at epic scope (verified by
dry-run, then used for DEC-010).

## Evidence artifacts (committed, recomputable)

- `readiness/T0-catalog-feasibility.matrix.tsv` (sha `711aecae…`) — 996 rows under
  the oracle predicate: 986 pass / 10 fail. The ground truth.
- `readiness/T0-identity-fallback.evaluation.tsv` (sha `48beef14…`) — identity
  candidate evaluation; recomputes byte-identical.
- `internal/geometry/cmd/t0sweep`, `t0identity` live in the T0 git worktree at
  `/private/tmp/claude-501/…/b2a645f2-…/scratchpad/t0` (detached HEAD `f03c7fe`),
  **verified still present** as of this session. If it is ever gone, rebuild from
  `f03c7fe`.

## T1 — done and independently verified (commit `36e1043`)

Built the DEC-006 ladder the 20 runs never built:
- `internal/geometry/lod/ladder.artifact.json` (sha `f0646ed7…`, 10.4 MB) — per
  geometry × band, the finest oracle-passing candidate, 996 rows, 283 candidate
  geometries carrying 475 dedup candidates, full provenance and per-attempt
  rejection log.
- `internal/geometry/lod/ladder.recipe.json` (sha `050fa04e…`) — binds v2 recipe
  `da3ff9d3…`, oracle `f2c9cd32…`, corpus, Mapshaper pin; `resolutions` is
  `[]float64`; one intermediate rung `79.63`; `identity_fallback: tried_last`.
- `internal/geometry/cmd/lodbuild/ladder.go` + `ladder_test.go` — the builder and
  its gates, incl. `TestLadderRebuildMatchesCommittedArtifact` (full rebuild,
  ~175 s, Node required, byte-identical — the determinism gate).

**Final catalog state: 993 pass, 3 no_artifact.** From 75 failing to 3.

## T2 — done (commits `3de2ea5`, `da65ba1`)

**T2.1 — `internal/geometry/ladder.go`.** The artifact, the ladder recipe and the
v2 recipe are embedded; `LadderTable` is a verified, indexed, **data-only** view.

- Load fails closed on a stale binding: corpus identity, oracle digest, ladder
  recipe digest, v2 recipe digest, `candidate_order`, `visibility_policy`,
  schema/kind. Also on per-row incoherence (a pass row with no stored candidate,
  a no-artifact row with no rejection log, an unknown band or status).
- `(geometry, band)` is proven a functional key: 996 rows → **566 keys, zero
  conflicts**. Profile is not needed in the key because the two profiles resolve
  to different geometry ids whenever they differ at all. `ParseLadderTable`
  enforces this rather than assuming it.
- `Lookup` returns three outcomes: `LadderPass`, `LadderNoArtifact` (DEC-009's
  typed absence, a recorded decision) and `LadderAbsent` (no row). Collapsing the
  first two would let a consumer fall back to an unjudged path for exactly the
  cases the oracle refused.
- **`publishedLODTable` stays nil on purpose**, pinned by
  `TestGenerateStaysSourceOnlyUntilSelectionIsRewritten`. Wiring the ladder into
  the current `lod.go` predicate would fall back to source on nearly every entity
  *while reporting success* — a silent green. The flip belongs with T3.
- Cost, measured: lazy `sync.Once` parse of ~10 MB takes **707 ms** on first use;
  the linker drops the embed entirely if unreferenced (lodbuild unchanged at
  15.32 MB; a consumer calling the loader grows 15.32 → 24.78 MB). **The 707 ms
  per-process cost is a real T3 input** for a CLI that generates one SVG.

**T2.2 — `internal/geometry/ladder_gate_test.go` (VAL-6 obligation 2).** Pure-Go,
offline, inside `make check` (~28 s). Recomputes topology, protected visibility,
IoU/recall, path bytes, complete-file estimate, points, parts, omissions and
every visibility component for **all 996 rows** and compares against the
committed claims. All 993 pass rows reproduce exactly; the 3 no-artifact rows are
rechecked at the identity rung and still fail. The predicate is restated in the
gate rather than shared with the builder — a gate that calls the same helper as
the thing it checks proves only self-consistency.

Mutation teeth (each logs which mechanism caught it, so none passes vacuously):
perturbed coordinate → off-grid topology damage; four tightened thresholds →
every sampled row flips to rejected; dropped component → omission/part drift.
Stale corpus/oracle/recipe is load-blocking upstream in `ParseLadderTable`.

**Coverage boundary, stated honestly:** the gate does *not* recompute the
rejected rungs of a no-artifact row — those geometries are not stored, so they
need the Mapshaper sweep. That stays the determinism gate's job.

## `make check` baseline — 5 known failures, unchanged

Before T2: **174 passed, 5 failed**. After T2: **205 passed, 5 failed** — the
same five, no regressions. They are T4's work:

- `geometry`: `TestFullCorpusBothProfiles`, `TestApprovalDigestAndBudgets`,
  `TestRepresentativeNaturalRatiosAndArbitraryFrames` — all
  `hard_budget_exceeded`, because `Generate()` is still source-only.
- `lodbuild`: `TestLODAQRUProjectionAlignedCheckpoint`, `TestLODSpikeFullCorpus`.

Note `TestDiagnosticSourceInventoryIsExactAndBiting` is an exact source-file
inventory gate: **any new file under `internal/geometry(/cmd/lodbuild)` must be
registered** in its `postDiagnosticAdditions` list or the gate reddens. That is
by design; `internal/geometry/ladder.go` was registered in `3de2ea5`.

## CLOSED — Croatia de_facto card (DEC-010, `45385fb`)

Decided by the owner this session: **accept the typed no-artifact outcome.** No
override, no threshold change, no country literal. Croatia renders in 3 of its 4
profile × band slots. The committed rejection log is the evidence: rung 79.63 →
2205 bytes against the 2200 cap, rung 64 → source component 2 lost (contribution
111 against the frozen 110 threshold); the two crossings coincide, so no rung
bridges them. Option 2 (a DEC-007 group anchor) stays pre-authorized as a *data*
change if the site later shows the missing card is a material defect — that would
not reopen DEC-010.

## Next step — T3 (do this first in the new session)

Rewrite `lod.go` selection to lookup-and-verify against `LadderTable`, and flip
`publishedLODTable` in the **same** commit:

1. Consult `LadderTable.Lookup(geometryID, band)` for the fitted band. On
   `LadderPass`, bind the stored candidate and verify it rather than re-deriving
   it. On `LadderNoArtifact`, return the typed no-artifact outcome (a first-class
   result, not an error — DEC-010 makes HR/de_facto/compact a top-200 exemplar).
   On `LadderAbsent`, route to the unchanged DEC-005 source path.
2. Delete `restoreRequiredComponents` and the two derived-path deviation gates
   (`lod.go:245-246`, `518-520`). They are the superseded predicate.
3. Route explicit source / >700 effective scale to the unchanged DEC-005 path.
4. Weigh the 707 ms first-load cost for a one-SVG CLI invocation. Per-geometry
   lazy decode is the obvious lever if it matters; decide deliberately.
5. **Emit real SVGs for a sample of top-200 countries as the empirical proof.**
   Tests passing is not the same as the product working.

Then T4 (realign the 5 failing tests + VAL-6 mutation teeth for the runtime path),
T5 (full-catalog SVG proof + owner contact sheet).

Task list: #1 DEC-010 (done), #2 T2.1 (done), #3 T2.2 (done), #4 this checkpoint,
#5 governance reconcile (deferred, blocked on the mate fix).

## Hard constraints (from accepted decisions — do not violate)

Byte budgets 2200/7500 path, 2500/8000 file — not to be raised. No country literal
or per-entity branch. Pure-Go offline runtime, no CGO/Node/network at generation.
P1 corpus bytes immutable. Byte-identical determinism (REQ-8). Presentation-free
geometry (INV-1) — CSS/theming/customization is P3, and INV-1 is precisely what
makes it possible. q=0.01 quantization untouched. Thresholds frozen (DEC-008).
