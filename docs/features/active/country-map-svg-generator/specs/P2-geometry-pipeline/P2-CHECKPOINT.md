# P2 checkpoint — resume point for a fresh session

Last updated at commit `7ccb43a` (T3 complete). Read this first; it lets a new
session resume without the originating chat.

## The one-line status

**The 20-run P2 failure is over.** The ladder is built (T1), embedded and
row-by-row recomputed (T2), and now actually served by `Generate()` (T3):
measured over the whole catalog through the shipped path, **993 rendered, 3
typed no-artifact, 0 errors, 0 source fallbacks, 0 over the band cap**. Real SVGs
were rendered and looked at — Italy, Japan, Indonesia, Brazil, India, Greece and
the rest are recognizable silhouettes, not source-only fallbacks.

Two findings came out of looking at the output and are captured as backlog work,
not fixed here: **Russia renders as a crude blob** (`WKI-C83B0B9EB4EE`, high) and
**a card frames the full UN claim rather than the mainland** (`WKI-490046152C71`,
owner decision). **Next step is T4**, then T5.

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
- `internal/geometry/cmd/t0sweep`, `t0identity` live in a git worktree at
  `/private/tmp/claude-501/…/b2a645f2-…/scratchpad/t0` (detached HEAD `f03c7fe`).
  It is still on disk and still registered in `git worktree list`, but that path
  is **a dead session's scratchpad** — session-scoped and orphaned, so it can
  vanish on cleanup without warning. Treat `f03c7fe` as the real source of truth
  and rebuild the worktree from it when T5 needs those tools; do not depend on
  the path surviving.

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
- `publishedLODTable` stayed nil through T2 on purpose (T3 flipped it). Wiring
  the ladder into the then-current `lod.go` predicate would have fallen back to
  source on nearly every entity *while reporting success* — a silent green.
- Cost, measured: lazy `sync.Once` parse of ~10 MB takes **707 ms** on first use;
  the linker drops the embed entirely if unreferenced (lodbuild unchanged at
  15.32 MB; a consumer calling the loader grows 15.32 → 24.78 MB). Tracked as
  `WKI-6638BACD6E20`. T3's measurement: a cold single-entity process is **1.53 s**
  end to end for 4 rows, so the load is roughly half of it; the full catalog is
  162 s for 996 rows, i.e. the per-row work dominates once the process is warm.
  Lazy per-geometry decode remains the obvious lever and belongs with P3's CLI,
  where the real invocation shape is known.

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

## T3 — done (commits `4577c04`, `7ccb43a`)

`Generate()` resolves the ladder on first use and serves it. Selection is a
lookup, not a fine-to-coarse walk: exactly one derived candidate is legal per
band, so the pre-DEC-006 substitution of the standard rung inside a card
viewport is gone.

**Three superseded gates no longer run for a ladder candidate**, all kept intact
for the DEC-005 source path:

1. `restoreRequiredComponents` — injected full-detail source components the
   candidate deliberately dropped.
2. The derived raw-deviation gate — replaced by the silhouette oracle (DEC-006).
3. The raw P1 `minimum_parts` count — **DEC-007 item 2 explicitly moved this** to
   the protected-visibility policy, and DEC-007 names Indonesia at compact as the
   case: 20 declared parts cannot coexist with the frozen 2200/2500 budgets at
   any generic resolution. Verified, not assumed: all 6 committed rows below
   their declared minimum (ID both profiles at compact, KI both profiles at both
   bands) carry explicit `protected_subscale` provenance on every omitted
   identity-set member, and KI's two omitted components contribute **exactly 0**
   filled pixels at the compact grid.

**The final-deviation gate at `finalizeGridPhase` is shared by both paths**, so
it became conditional rather than deleted. The T2 checkpoint's plan named it by
line number (`518-520`) and following that literally would have weakened the
source path. Read the call graph, not the line numbers.

`ladderVerdictApplies` is one named predicate on purpose — it decides whether a
committed oracle verdict may be trusted at all. A verdict earned at rotation 0,
under the v1 projection contract, at automatic quality, with the caller's byte
ceiling cleared does not transfer to a request that changes any of those; such a
request takes the source path.

DEC-009's outcome is now the typed `ErrNoArtifact` with an `IsNoArtifact`
predicate. Falling back to source there would emit the exact silhouette the
oracle refused (18190 bytes for Croatia's de_facto card against a 2200 cap).

**One defect found by rendering, not by testing.** Softening is a runtime
representation choice the oracle never judged, and on small compact geometries it
inflates the path several-fold: Grenada's 575 committed bytes became **2230**,
clearing the 2200 band cap while still fitting the preset's 2500 complete-file
maximum. Eight rows did this (GD, VI at compact; LU, TC at standard). The band
cap now bounds whatever representation is finally emitted.

**Measured over the whole catalog through the shipped path:** 993 rendered, 3
typed no_artifact, **0 errors, 0 source fallbacks, 0 over the band cap**, 754
byte-identical to the committed row. The remaining 239 differ because softening
produced a smaller path that still fits.

Guarded by `TestShippedCatalogServesTheLadder` (~172 s) plus
`TestShippedPathHoldsTheBandCapAgainstSoftening`. This is the gate that would
have caught the original 20-run failure: T2's gate proves the artifact is
internally sound, T3's proves the artifact is what `Generate()` actually serves,
and the twenty runs failed exactly in the gap between those two claims.

`internal/geometry/cmd/svgproof` runs the same sweep and writes real SVG files
plus an HTML contact sheet. It is not in `make check` — writing ~1000 files is
not a gate — but it is how a human looks.

## `make check` — 207 passed, 3 failed

Baseline was 174/5 before T2, 205/5 after T2. Now **207 passed, 3 failed**: T3
fixed two of the five and added none.

- `geometry`: `TestRepresentativeNaturalRatiosAndArbitraryFrames` — the Chile
  framing question, now an owner decision (`WKI-490046152C71`). **Do not "fix"
  this by editing the assertion.**
- `lodbuild`: `TestLODAQRUProjectionAlignedCheckpoint`, `TestLODSpikeFullCorpus` —
  pre-existing, T4's work.

Full `make check` now costs roughly 7 minutes: the ~175 s Node ladder-rebuild
determinism gate plus the ~172 s pure-Go shipped-catalog gate. Both earn it.

Note `TestDiagnosticSourceInventoryIsExactAndBiting` is an exact source-file
inventory gate: **any new file under `internal/geometry` or
`internal/geometry/cmd/lodbuild` must be registered** in its
`postDiagnosticAdditions` list or the gate reddens. That is by design;
`internal/geometry/ladder.go` was registered in `3de2ea5`. Files under other
`cmd/` subdirectories (such as `svgproof`) are outside its globs.

## OPEN — two findings the T3 render surfaced

Both are captured as durable backlog rows; neither was fixed in T3.

**`WKI-C83B0B9EB4EE` (high) — Russia renders as a crude blob.** RU/un/standard
selects rung **24**: 508 bytes of a 7500 cap, 45 points, 1 part, IoU 0.9204.
Italy at the same band spends 6829. Russia is the only anomaly of its kind —
exactly 14 committed rows select a rung at or below 64 and the other 12 (CA, GR,
AX, DK, FK, HK) sit at 68–98% of budget, i.e. coarse because they hit the cap.
Russia is coarse while leaving 93% of its hero budget unspent, so the
fine-to-coarse search rejected every rung from 512 down to 32 for a reason the
artifact does not record (per-attempt logs are stored only on no_artifact rows).
The oracle cannot see it: minimum IoU is 0.40. Suspected antimeridian topology
damage, **unverified**. Diagnosing needs a Mapshaper sweep for RU alone with
per-attempt logging — T1/T5 work.

**`WKI-490046152C71` (owner) — mainland or full claim?** The viewBox is fitted to
the whole projected source before any candidate is chosen, so it reserves space
for components the silhouette may deliberately not draw. Chile spans 43.0° of
longitude via Easter Island, Salas y Gómez and Juan Fernández against 13.5° for
the mainland, so it renders against one edge of a near-square box; the United
States does the same via Alaska. **Pre-existing, not a T3 regression** — the AU
viewBox is identical on the ladder and source paths (126.4×144 both) — visible
now only because these entities render at all. Decide before T5's contact sheet.

## CLOSED — Croatia de_facto card (DEC-010, `45385fb`)

Decided by the owner this session: **accept the typed no-artifact outcome.** No
override, no threshold change, no country literal. Croatia renders in 3 of its 4
profile × band slots. The committed rejection log is the evidence: rung 79.63 →
2205 bytes against the 2200 cap, rung 64 → source component 2 lost (contribution
111 against the frozen 110 threshold); the two crossings coincide, so no rung
bridges them. Option 2 (a DEC-007 group anchor) stays pre-authorized as a *data*
change if the site later shows the missing card is a material defect — that would
not reopen DEC-010.

## Next step — T4 (do this first in the new session)

Three failing tests remain. Two are genuinely stale and one is not.

1. **`lodbuild` — `TestLODAQRUProjectionAlignedCheckpoint` and
   `TestLODSpikeFullCorpus`.** Pre-existing since the T0A.1 spike. The exact
   reason was read after T3, not assumed: the checkpoint test builds its **own**
   v1 tier table (`stage=production_table`) with no `Ladder`, so it takes the
   legacy branch and is rejected by the superseded predicate —
   `AQ/hero production tier mismatch ... Fallbacks:[standard:raw_deviation:12.2]`.
   That is the DEC-006 problem surviving inside the spike harness rather than in
   production, which T3 fixed. So these are genuinely untouched by T3's ladder
   path, but they are not "unrelated to the ladder" either: they exercise exactly
   the predicate the ladder replaced.

   The T4 decision is therefore whether the spike harness moves to the oracle
   predicate or is retired as historical evidence — not whether to move its
   numbers. `DiagnosticGridPhases` deliberately still reproduces the superseded
   path and now says so in its doc comment; the same question applies to it.
2. **`geometry` — `TestRepresentativeNaturalRatiosAndArbitraryFrames`.** Blocked
   on the owner's framing decision (`WKI-490046152C71`). **Do not edit the
   assertion to make it pass.** It encodes mainland framing; production frames
   the full claim. One of the two has to change, and which one is not a test
   question.
3. **VAL-6 mutation teeth for the runtime path — partly done, not yet
   automated.** The shipped gate was hand-mutated at the end of T3 and bites on
   both regression vectors, with the right message:

   - forcing `ladderVerdictApplies` to `return false` →
     `AD/de_facto/compact: served tier "source", not the committed band — the
     ladder is not reaching the shipped path`;
   - regressing `publishedLODTable` to nil → the same failure, plus
     `TestPublishedTableIsTheCommittedLadder: the published table is not
     ladder-backed`.

   That is a hand-verification, which by the project's own rule is a gate not yet
   written. T4 should make it automatic — the cheap shape is a table-driven test
   that runs the shipped assertions against a deliberately ladder-less table and
   requires them to fail.

Then T5: full-catalog SVG proof and the owner contact sheet
(`internal/geometry/cmd/svgproof -out <dir>` already produces both). Resolve
`WKI-C83B0B9EB4EE` (Russia) before that sheet goes to the owner — a top-200
country rendering as a blob is the first thing anyone will notice.

## Hard constraints (from accepted decisions — do not violate)

Byte budgets 2200/7500 path, 2500/8000 file — not to be raised. No country literal
or per-entity branch. Pure-Go offline runtime, no CGO/Node/network at generation.
P1 corpus bytes immutable. Byte-identical determinism (REQ-8). Presentation-free
geometry (INV-1) — CSS/theming/customization is P3, and INV-1 is precisely what
makes it possible. q=0.01 quantization untouched. Thresholds frozen (DEC-008).
