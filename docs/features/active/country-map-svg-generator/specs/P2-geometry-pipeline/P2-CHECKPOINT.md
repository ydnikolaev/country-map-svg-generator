# P2 checkpoint — resume point for a fresh session

Last updated at commit `fd1de11` (T5 complete). Read this first; it lets a new
session resume without the originating chat.

## The one-line status

**P2's product work is done.** The ladder is built (T1), embedded and row-by-row
recomputed (T2), served by `Generate()` (T3), guarded with mutation teeth (T4),
and the full catalog is rendered and reviewed by the owner (T5): **993 rendered,
3 typed no-artifact, 0 errors, 0 source fallbacks, 0 over the band cap**, with
`make check` green at **215 passed / 0 failed**.

Two defects were found by *looking at output*, not by testing, and both are
fixed: Russia rendered as a blob (DEC-012) and France as a speck in an empty card
(DEC-013). Neither was visible to any metric — the oracle's IoU floor is 0.40 and
Russia scored 0.92.

What remains is not product work. **`#15`, the governance reconciliation debt**,
is still blocked by the mate fence below, and P3 inherits the component-selection
seam.

## What was wrong (the root cause, confirmed)

*Historical — fixed in T3. Kept because the lesson outlived the bug.*

DEC-006 (accepted) replaced raw full-coastline deviation with a target-scale
raster **silhouette oracle** as the acceptance boundary for derived candidates.
The production selection path in `internal/geometry/lod.go` never received that
change: it gated on `matchedBoundaryDeviation`, referenced the oracle nowhere,
and injected full-detail source components via `restoreRequiredComponents` — the
exact mechanism DEC-006 rejected. The shipped `Generate()` was source-only
because `publishedLODTable` was nil. Twenty runs tried to make the superseded
predicate pass. Under the correct predicate the catalog went from **75 failing
outputs to 3**.

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
- **DEC-011** accepted (`5be32b7`) — cropping is an explicit seam over the
  selected candidate, owned by P3. Its item 1 (full-claim framing) is superseded
  by DEC-013.
- **DEC-012** accepted (`2241254`) — topology is judged after canonicalization.
- **DEC-013** accepted (`fd1de11`) — the card is fitted to the silhouette it
  draws; compact contribution threshold 111; DEC-008's approval renewed on the
  full catalog.
- **AM-005** verified (`ADV-006` pass, independently recomputed) — the oracle-gate
  fix. Supersedes AM-003/AM-004 in text.
- **AM-003, AM-004** `retracted` on `ADV-007` and `ADV-009`, citing AM-005 as the
  verified successor. **The fence is gone**; P2 is `in_progress` with
  `run.scaffold` available. The section below records why the earlier reading of
  this as unliftable was wrong.

## The mate fence — lifted, and the false blocker that held it up

**Superseded. Kept because the mistake is instructive, not because the finding
stands.** The claim below — that `applied` is a lifecycle sink needing a mate SSOT
fix and a CLI release — was wrong, and it cost this epic a governance detour plus
all of P3 shipping outside the lifecycle under DEC-014.

`retract` was always the terminal edge out of `applied`. The shipped state machine
declares it:

```yaml
retract:
  from: [applied]
  to: retracted
  loop: amendment-retract
```

with `retracted` in both `states` and `terminal`, and named by
`terminal_classes.retraction`. `amendmentFencingStates` is
`{impacting, accepted, applied}` (`internal/work/amendment.go:261`), so a
`retracted` amendment stops fencing **by construction** — exactly the property the
sketched `supersede` edge was designed to obtain, under a name the vocabulary
already had.

**How the error was made, so it is not repeated.** The earlier probe tried
`reject` and `verify`, saw both refuse, and generalized to "exactly one edge out
of `applied`" from a read of `state-machines.yaml:65-70` that did not include the
`retract` entry. Two rules follow:

- **Enumerate a state machine from the registry, never from the verbs you happened
  to try.** Three probes are not a closed set.
- **Distinguish an illegality refusal from an evidence refusal.** They have
  different messages and mean opposite things. The falsification test that settles
  it costs one command per verb:

| Probe on `AM-003` (`applied`, CLI v1.4.3) | Result |
| --- | --- |
| `amendment accept` | `illegal amendment transition accept from applied` |
| `amendment apply` | `illegal amendment transition apply from applied` |
| `amendment reject` | `illegal amendment transition reject from applied` |
| `amendment retract` | `requires exactly one blocking advisor-receipt` — **legality passed** |

Legality is checked before evidence, so a refusal that names evidence is proof the
transition is legal.

**What the retract needed.** One blocking advisor receipt bound to the *current*
epic closure, and no readiness receipt. `ADV-003`/`ADV-004` are blocking and
`result: fail`, but bind an older tracker digest, so `mate advisor record` issued
`ADV-007` and `ADV-009` restating the same reviewers' verdict against the current
closure. This is a re-binding of an accepted independent verdict, not a fresh
review: AM-005 — `verified` — already states that both "were authored from
diagnoses that independent review refuted". A `verify` receipt would have been
fabrication; a `block` receipt authorizing `retract` is the disposition AM-005
already declared correct.

**Sequencing trap.** Each retract bumps the epic revision, which invalidates any
receipt recorded before it. Record the receipt for one amendment, retract it, then
record the next. `ADV-008` was recorded too early and is a live orphan artifact —
immutable, unused, and it cannot be rebound (`already exists outside this
operation`).

### Historical: the reading that was wrong

Kept verbatim below the line so the correction above has something to correct.

Re-probed under mate **v1.4.3** (the finding `FND-1B950D899CE4` was filed against
v1.4.1, so this is a fresh confirmation, not a citation):

Re-probed under mate **v1.4.3** (the finding `FND-1B950D899CE4` was filed against
v1.4.1, so this is a fresh confirmation, not a citation):

| Probe | Result |
| --- | --- |
| `amendment reject AM-003` | `illegal amendment transition reject from applied` |
| `amendment verify AM-003 --input advisor-receipt=ADV-003#c5b9347f…` | `advisor receipt does not pass and bind the current amendment closure` |
| `spec transition P2 amendment-resolve` | `spec P2 is fenced by impacting amendments AM-003,AM-004` |

Read in the mate SSOT (`~/Developer/projects/mate`):

- ~~`internal/workdocs/assets/bundle-v1/registries/state-machines.yaml:65-70` — the
  amendment machine has **exactly one** edge out of `applied` (`verify`), and
  `reject` is legal only from `[proposed, impacting]`. A failed applied amendment
  is a lifecycle sink.~~ **False.** There are two edges out of `applied`:
  `verify` → `verified` and `retract` → `retracted`. The rest of the sentence is
  right — `reject` is indeed limited to `[proposed, impacting]`, which is what made
  the wrong generalization feel confirmed.
- `internal/work/amendment.go:283` — the blocking set is
  `{impacting, accepted, applied}`, so `applied` fences by construction.
- The fence is enforced at `internal/work/plan_accept.go:98`,
  `internal/work/run_start.go:95`, `internal/work/run_finish.go:140` and
  `internal/work/transition.go:148`. **Only `cancel` and `invalidate` are exempt.**

So `plan invalidate` is legal (dry-run green: PLAN-013, authority
`run-result=RUN-020-RESULT#b440805079f0242620e9529c75713c8db0b7637dbd8caaf1fae44cc5a4447289`,
spec revision 123 → 124) but `plan supersede`/`accept` and `run start` are not.
~~There is no governed route to a P2 run.~~ **False, per the correction above** —
`retract` was the route. What remains true: verifying AM-003 to escape the fence
would have been fabricated evidence and was correctly refused.

**Owner instruction (that session): continue T2+ as ordinary engineering commits
outside the governed run wrapper**, with the reconciliation debt tracked. That
instruction was sound given what was believed at the time; it is spent now, and
the debt it created is `#15`.

### Dead: the mate fix, sketched from the SSOT

**Do not build this.** `WKI-4062B33B8FEA` is `rejected` on a falsified premise:
`retract` already provides this edge, so the sketch below would have added a
second, redundant terminal disposition to a shared, fleet-wide state machine. Kept
only to document what was nearly shipped upstream on a misreading.

Read at `~/Developer/projects/mate` so a future session does not re-derive it.
The change is small and mostly data:

| Where | Change |
| --- | --- |
| `internal/workdocs/assets/bundle-v1/registries/state-machines.yaml:61-70` | add `supersede: {from: [applied], to: superseded, loop: amendment-supersede}`, add `superseded` to `states` and `terminal` |
| `internal/workdocs/assets/bundle-v1/registries/loops.yaml` | add the `amendment-supersede` loop entry, modelled on `amendment-reject` |
| `internal/cli/amendment.go:14` | add `newAmendmentEvidenceTransitionCmd("supersede")` — the factory is already generic |
| `internal/work/amendment.go:283` | nothing: the blocking set is `{impacting, accepted, applied}`, so a `superseded` amendment stops fencing by construction |

Evidence probably needs no new code. `validateAmendmentTransitionEvidence`
(`internal/work/amendment_transition.go:167`) falls through to a default branch
requiring one readiness receipt whose `AuthorizedNext` is `amendment.<event>`, so
`mate readiness issue --for amendment.supersede --amendment AM-003` should
authorize it as-is. Confirm rather than assume — the schema may also need
`superseded` added wherever amendment states are validated.

**Do this in its own session, in the mate repo.** It is a fleet-wide change to
the tool that governs every project, with its own doctrines, gates and release
process, and it ends in `mate fleet pull` back into this consumer. It is the
natural first step of the next phase rather than the last step of this one.

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

Guarded by `TestShippedCatalogServesTheLadder` (~197 s) plus
`TestShippedPathHoldsTheBandCapAgainstSoftening`. This is the gate that would
have caught the original 20-run failure: T2's gate proves the artifact is
internally sound, T3's proves the artifact is what `Generate()` actually serves,
and the twenty runs failed exactly in the gap between those two claims.

**T3 shipped that gate with a hole, and T4's tooth found it** — see below. Read
the T4 section before trusting any claim in this one.

`internal/geometry/cmd/svgproof` runs the same sweep and writes real SVG files
plus an HTML contact sheet. It is not in `make check` — writing ~1000 files is
not a gate — but it is how a human looks.

## T4 — done

**1. The shipped gate had a hole, and closing it is T4's most important change.**
`TestShippedGateRejectsALadderlessPipeline` drives the gate's own assertions
against two deliberately broken pipelines. It caught the first shape
(`publishedLODTable` regressed to nil) on every sampled row — but on the second
(the ladder loading yet selection not consulting it, so the superseded tier walk
runs) **40 of 60 rows still satisfied every assertion**. The legacy path reads
the same stored candidates, so it reports the same band and the same byte
budget; it merely restores source components and re-judges by raw deviation
first, producing geometry no oracle saw.

The fix is a provenance discriminator: a ladder candidate computes no raw
deviation and restores nothing, so `RawDeviation != 0 || len(Restored) != 0`
means the superseded path ran. With it, both shapes are rejected 60 of 60. **T3's
commit message and this checkpoint both claimed that gate "would have caught the
20-run failure" before this was verified** — it would have caught one shape of it.

**2. `TestRepresentativeNaturalRatiosAndArbitraryFrames` realigned under
DEC-011.** The old assertion demanded `ratios["CL"] < .75`, i.e. mainland
framing, which was never a decided contract — it was unreachable while Chile
failed on the byte budget, so nothing ever tested it. DEC-011 fixes the default
as full-claim, so the test now asserts the full-claim band and states that a
value below .75 would mean the deferred second mode appeared by accident.

**3. The two `lodbuild` failures are restated, not deleted.** Both asserted that
the superseded pre-DEC-006 path succeeds — the v1 tier artifact selected by raw
boundary deviation — which is the exact claim DEC-006 was accepted for refuting.
They are now characterizations of that path:

- `TestLODSpikeFullCorpus` keeps its strict determinism check and requires the
  over-budget set to be non-empty, so DEC-006's motivating evidence cannot be
  silently outlived by a green test. It is now **sampled, not a full sweep**:
  every 8th entity plus BD, ~32 of 249. The reason is concrete — once the
  over-budget outcome stopped being fatal the sweep completed, which meant
  falling back to the full-detail source path for most of the catalog and
  pushing the single `lodbuild` package past the **10-minute go test timeout**.
  Note precisely what the non-empty guard verifies: BD is force-included because
  it does not fit, so the guard is "BD still does not fit".
- `TestLODAQRUProjectionAlignedCheckpoint` drops only its tier expectation and
  keeps the runtime bound, binding errors, coordinate space, grid-phase guards
  and determinism. A deeper pre-existing bug surfaced once the tier assertion
  stopped failing first: it compared a full-generate provenance against a
  selection-only one and required `DeepEqual`, but `ParserRoundTrip` is a
  render-stage flag that `SelectLODProvenance` never sets. Normalized, with the
  reason recorded in the test.

## `make check` — green

Baseline was 174 passed / 5 failed before T2, 205/5 after T2, 207/3 after T3.
**After T4 it is green for the first time in this epic.**

**215 passed, 0 failed, 0 package failures**, measured serially: `catalog` +
`cmd` 109/0, `geometry` 106/0.

Three things whoever runs it next needs to know:

- **`make check` now passes `-p 1`, and that is load-bearing.** Several tests
  assert wall-clock brakes (a 30 s post-build micro-brake, a 10 s per-generation
  bound) that exist to catch a runaway search — the failure mode behind the
  twenty interrupted runs. With packages in parallel they measure CPU contention
  instead: `lodbuild`'s micro-brake was observed at **34.5 s alongside the 381 s
  geometry package** and comfortably under 30 s on its own. Raising the
  thresholds until they stop flaking would leave the brakes shaped like guards
  while guarding nothing. Cost: check is now roughly the sum of the packages
  (~12 min) rather than the slowest. Tracked as `WKI-B05B4B4287A2` — the real fix
  is to express the brakes in work done rather than seconds.
- **Summarize package-level failures, not just test-level ones.** A `go test
  -json` package failure carries no `Test` field, so a test-only filter reports
  a timed-out package as a clean run. That happened twice during T4 — once
  `lodbuild` died at 600.4 s and the summary said "26 passed, 0 failed".
- Both heavy packages sit near the default 10-minute per-package timeout:
  `geometry` ≈ 322 s (the 197 s shipped-catalog gate plus the 48 s tooth plus the
  28 s offline ladder gate) and `lodbuild` ≈ 373 s (the 175 s Node rebuild plus
  the 75 s sampled characterization). Adding another full-catalog sweep to either
  package will break the build before it fails an assertion.

Note `TestDiagnosticSourceInventoryIsExactAndBiting` is an exact source-file
inventory gate: **any new file under `internal/geometry` or
`internal/geometry/cmd/lodbuild` must be registered** in its
`postDiagnosticAdditions` list or the gate reddens. That is by design;
`internal/geometry/ladder.go` was registered in `3de2ea5`. Files under other
`cmd/` subdirectories (such as `svgproof`) are outside its globs.

## T5 — done (commits `d84c16c` … `f522095`, `fd1de11`)

**The catalog renders and was reviewed.** `internal/geometry/cmd/svgproof -out
<dir>` writes every SVG plus one self-contained contact sheet with all 993
silhouettes inlined. `contact-sheet.json` is committed as readiness evidence at
`readiness/T5-catalog-render.contact-sheet.json`; the ~5.5 MB of SVG files and
the 3.1 MB page are build outputs and deliberately not committed.

The sheet carries a **frame-coverage** column — what share of the viewBox the
drawn silhouette occupies — because answering "which cards need a look" by eye
over 996 cells is not answering it. That column is what found the framing defect.

`internal/geometry/cmd/svgshowcase -out <file>` renders the capability page: one
country through five CSS treatments, ten more across different geographies, eight
saturated palettes, a gradient, a two-tone, and a reframing pair. Every silhouette
is unmodified `Generate()` output — the page proves INV-1 rather than asserting it.

## Two framing defects, both found by looking

**Russia rendered as a blob** (`WKI-C83B0B9EB4EE`, closed by DEC-012). The
acceptance predicate judged topology on **pre-canonical** fitted coordinates — a
precision no emitted path has. Simplification left a self-crossing of ~1e-4
against a q=0.01 output grid, and canonicalization erased it two steps later.
RU/un was rejected at every rung 512 down to 32 and settled on 24: one part, 508
of 7500 hero bytes. Fixed on both the oracle and the runtime side; fixing only
the build left Russia committed at rung 80 and falling back to source at 101605
bytes. Exactly **2 of 996 rows changed**, both Russia's — which confirmed the
diagnosis.

**France rendered as a speck** (DEC-013). The layout was fitted to the whole
territorial claim while fidelity was already judged against the visible
reference, so the card framed for components the oracle had excluded. France drew
2 of 11 components and filled 11.9% of its frame. Fitted to the candidate
instead: France is the hexagon with Corsica, Chile is tall and narrow, cards
under a fifth full went **20 → 0**. Croatia's `un` card cost two bytes, so the
compact contribution threshold moved 110 → 111 and its disputed component is now
recorded as `subscale` with provenance.

`lodbuild -ladder-explain <ISO>` replays the rung search for one entity and prints
every rejection. It is what made both diagnoses possible: the artifact stores
per-attempt logs only on `no_artifact` rows, so a row passing at a coarse rung
recorded nothing about why the finer rungs lost.

## Governance state as of T5

- **DEC-010** — Croatia's `de_facto` card stays a typed no-artifact. Unchanged.
- **DEC-011** — item 1 (full-claim framing) superseded by DEC-013. The rest
  stands: cropping is an explicit seam over the selected candidate, expressed
  through the identity vocabulary, never a country literal; the acceptance rule
  for a cropped variant is still open. Routed to P3.
- **DEC-012** — topology judged after canonicalization.
- **DEC-013** — the card is fitted to what it draws; compact threshold 111;
  **DEC-008's approval renewed on the full catalog**. DEC-008 froze an oracle
  digest that changed twice, so its validity condition was already broken; note
  it is an owner approval, not a calibration procedure, so any future threshold
  or digest change means showing the owner a rendered catalog again.

## Backlog

Run `mate backlog list`. As of T5: `WKI-C83B0B9EB4EE` **done**;
`WKI-490046152C71` deferred and now largely answered by DEC-013 — what remains of
it is component *selection*, not framing; `WKI-B05B4B4287A2` (wall-clock brakes
force `-p 1`) and `WKI-6638BACD6E20` (707 ms first load) captured.

## Next step — P2 closure, then P3

P2's product work is done: the pipeline serves the committed ladder, the catalog
renders, and the evidence is committed and reviewed. Two things remain.

1. **The governance reconciliation debt, `#15`.** Everything from T1 onward
   shipped as ordinary engineering commits while AM-003/AM-004 were believed to
   fence P2 permanently (`FND-1B950D899CE4`). **The fence is now lifted** — both
   are `retracted` and `run.scaffold` is available on P2 — so the debt is
   reconcilable today and nothing upstream is waited on. The honest shape is a
   governed successor plan bound to the as-built followed by a run recording it,
   not a backdated receipt: `mate plan invalidate` on PLAN-013 is legal and its
   authority is RUN-020's result. Re-derive the ordering from
   `mate work snapshot country-map-svg-generator --json`, which is the authority on
   legality.
2. **P3.** It inherits a real surface: the typed no-artifact outcome as a
   first-class result, the component-selection seam DEC-011 and DEC-013 both
   route to it, and the 707 ms first-load decision that belongs where the
   invocation shape is known.

## Hard constraints (from accepted decisions — do not violate)

Byte budgets 2200/7500 path, 2500/8000 file — not to be raised. No country literal
or per-entity branch. Pure-Go offline runtime, no CGO/Node/network at generation.
P1 corpus bytes immutable. Byte-identical determinism (REQ-8). Presentation-free
geometry (INV-1) — CSS/theming/customization is P3, and INV-1 is precisely what
makes it possible. q=0.01 quantization untouched. Thresholds frozen (DEC-008).
