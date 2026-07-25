# Handoff — country-map-svg-generator, epic remainder

Written 2026-07-25 for a fresh session. Everything here was verified against
source, the CLI, or the committed artifact — citations are given so you do not
re-derive. Read `specs/P2-geometry-pipeline/P2-CHECKPOINT.md` and
`specs/P3-generator-cli/P3-CHECKPOINT.md` too; both were corrected in this session
and are now trustworthy on the fence question.

Repo root: `/Users/yuranikolaev/Developer/projects/country-map-svg-generator`,
branch `main`, HEAD `a44fda4` at the time of writing.

---

## 1. What changed in this session

**The impact fence is lifted.** Commit `a44fda4` — `fix(epic): lift the impact
fence that was never actually stuck`.

AM-003 and AM-004 sat in `applied` for all of P2's T2–T5 and all of P3 because two
checkpoints recorded that `applied` has exactly one outgoing edge (`verify`), which
neither could ever earn. **That was false.** `retract` is the second edge and was
always there.

| Action | Result |
| --- | --- |
| `mate advisor record … ADV-007 --amendment AM-003 --verdict block` | blocking receipt bound to the then-current closure |
| `mate amendment retract … AM-003` | `applied` → `retracted` |
| `mate advisor record … ADV-009 --amendment AM-004 --verdict block` | (ADV-008 became an orphan — see trap 5) |
| `mate amendment retract … AM-004` | `applied` → `retracted` |
| `mate backlog dispose WKI-4062B33B8FEA reject` | premise falsified |

Fence lifted, verified via `mate work snapshot`: no spec carries
`impact_blocked_by`, and P2 moved from the phantom `amendment.resolve` to
`run.scaffold` (available). `mate status` green.

**Do not re-litigate this.** Evidence, if you need it:

- shipped state machine (`~/Developer/projects/mate/internal/workdocs/assets/bundle-v1/registries/state-machines.yaml`,
  machine `amendment`): `retract: {from: [applied], to: retracted, loop: amendment-retract}`;
  `retracted` is in `states`, in `terminal`, and named by `terminal_classes.retraction`.
- `amendmentFencingStates = {impacting, accepted, applied}` — `internal/work/amendment.go:261`.
  `retracted` is absent ⇒ stops fencing by construction.
- `amendmentTerminalStates = {verified, rejected, retracted}` — same file, `:270`.
- Falsification test on the installed CLI v1.4.3: `accept`, `apply`, `reject` each
  refuse with `illegal amendment transition X from applied`; `retract` passes
  legality and refuses only on evidence. **Legality is checked before evidence**, so
  an evidence refusal proves the transition is legal.

Why `retract` and not `verify`: `verify` would have been fabricated evidence —
ADV-003/ADV-004 failed both amendments on independent review. AM-005 (`verified`)
states they "were authored from diagnoses that independent review refuted" and
supersedes both in full, so `retract` is the disposition already declared correct.

---

## 2. Governed state right now

```
epic country-map-svg-generator  in_progress  rev 44
P1  completed
P2  in_progress   next=run.scaffold           available    (PLAN-013 accepted, RUN-020 interrupted)
P3  draft         next=spec.accept-readiness  blocked by dep P2
P4  draft         next=spec.accept-readiness  blocked by dep P3
P5  draft         next=spec.accept-readiness  blocked by dep P3
AM-001 verified · AM-002 verified · AM-003 retracted · AM-004 retracted · AM-005 verified
```

`mate work snapshot country-map-svg-generator --json` is **the authority on
legality**. Re-derive next actions from it; never trust a prose list including this
one.

**Reality vs trackers:** P3's tracker says `draft, revision 1`, but **all of P3's
code is written, committed and green** — 7 commands (`init`, `validate`, `explain`,
`generate`, `inspect`, `preview`, `version`), T1–T5 done, whole catalog generated
offline in 3.3 s. It shipped outside the lifecycle under DEC-014 because of the
fence. That gap is debt `#15`.

---

## 3. Measured facts you should not recompute

From the committed ladder artifact (`internal/geometry/lod/ladder.artifact.json`,
996 rows, 993 pass / 3 no_artifact):

| Band | n | path median | p95 | max | IoU floor | IoU median |
| --- | --- | --- | --- | --- | --- | --- |
| compact (card) | 495 | 1962 B | 2184 B | 2200 B | 0.40 | **0.995** |
| standard (hero) | 498 | 4140 B | 7022 B | 7436 B | 0.40 | **0.999** |

Complete-file bytes = path + 280. QAB-2 (ARCH-001 line 93) requires card
median ≤600 / p95 ≤1200 / max ≤2500 and hero median ≤1800 / p95 ≤4000 / max ≤8000.
**Median and p95 fail in both bands; max passes in both** (card max 2480 of 2500 —
20 bytes of headroom).

Where the excess sits, compact band:

| Group | n | median bytes |
| --- | --- | --- |
| selected the finest rung (`512`) — small/simple entities | 165 | **522 B** — already compliant |
| fell to a coarser rung — large/complex entities | 328 | **2066 B** (min 1423, max 2200) |

Mechanism: `evaluateLadderCaseExplained`
(`internal/geometry/cmd/lodbuild/ladder.go:337-354`) walks rungs **fine-to-coarse**
and returns the **first** that passes IoU + bytes ≤ `band.PathCap` + file max +
topology + protection + dominant-component. Finest-that-fits under a 2200 cap lands
everything just under the cap by construction. IoU 0.995 against a 0.40 floor means
the ladder is buying near-pixel fidelity nobody asked for.

**Root diagnosis (use this framing, it is better than "page weight vs fidelity"):**
DEC-006 accepted two incompatible rules in one document — item 5 says path data is
"capped **initially** at 2,200 bytes … P3/P4 separately prove … median 0.6 KB",
while the same decision says "selects the **finest** candidate that passes every
generic predicate" and "the first passing candidate wins". Finest-that-fits-under-2200
cannot yield a 600 B median. The word **"initially"** is the lever DEC-006 left for
the number. DEC-006 also forbids reversing the search: "Candidate search may move
**only from finer to coarser**" — so lowering the cap is in-bounds, reversing the
order is not.

Page-weight reality check, so nobody over-dramatizes: 20 cards × 2165 B ≈ 43 KB
uncompressed, ~16 KB gzipped. The waste is invisible precision (226 points median
on a 128 px card), not bandwidth.

Nothing was lost by retracting AM-003 despite its slug `byte-budget-attainability`
— **checked, not inferred**: AM-003 is about outputs *exceeding* the cap (RU/un/card
at 101,605 B vs a 2,500 max) via source-tier fallthrough, not about median
clustering. Its one relevant rule ("any ladder adjustment requires a rebuild and a
fresh owner-approved contact sheet") survives in AM-005 item 5.

---

## 4. Blockers

| # | Blocker | Owner | Notes |
| --- | --- | --- | --- |
| 1 | ~~Impact fence~~ | — | ✅ lifted this session |
| 2 | Debt `#15` — P2/P3 have no governed receipts | agent | unblocked; first work item |
| 3 | `WKI-37F18A2AA6A5` — explicit `contain` frame fitted **before** visibility removals; up to 292 px off centre in a 300 px frame (worst NC; FM draws 6, removes 14, 3.96 px slack above vs 145.3 below) | agent, in P2 | **precondition of P4 VAL-7** |
| 4 | `WKI-C35A01E965DC` — only the profile's own long side is served from the ladder; any other size falls back to source and blows the ceiling (Greenland card: 2400 B at longSide 128, 59892 B at 160 — and 160 is a documented reference size **and** the P3 spec's own example) | agent, in P2 | **precondition of P4 VAL-2 + VAL-7**; effort=large |
| 5 | `WKI-33B6EC2482B3` — QAB-2 median 3.3× over | **owner** + rebuild | blocks P4 AC-1/AC-3 |
| 6 | No browser driver for P4's VAL-4/5/7 | decision, then agent | nothing in `go.mod`; precedent for a maintainer-only non-Go toolchain exists (`internal/geometry/lod/tool/node_modules`, pinned mapshaper). Candidates: chromedp (pure Go) or playwright-go |
| 7 | Contact-sheet acceptance (P4 AC-4) | **owner** | DEC-006 item 6 + P4's own text call this a planned acceptance boundary, not an open question. What the owner decides is *which byte exceptions are named and kept* (REQ-9) — and that depends on #5 |
| 8 | `WKI-B05B4B4287A2` — wall-clock brakes force `-p 1` | agent | promoted to **P4 precondition**: current serial gate ~757 s, P4 adds full-catalog double-generate + byte distributions + browser fixtures |
| 9 | P5 VAL-2 Windows amd64 clean-machine proof | **owner/infra** | no host on darwin/arm64. Linux dockerable, darwin local. Options: CI runner, borrowed machine, or narrow REQ-1 by decision |
| 10 | No successor spec for the crop / component-selection seam (DEC-015, `WKI-490046152C71` deferred) | **owner** | not in the DAG (only P1–P5). Either a P6 or an explicit non-goal |

**Why #3 and #4 are P4 preconditions, not deferred debt** — P4's own text:

- REQ-10: "tiny/huge portrait/landscape/square `contain` assets use one uniform
  scale, **remain centered** and never distort geography."
- VAL-7: "tiny, huge, portrait, landscape and square containing frames … intrinsic
  ratio, CSS sizing, uniform transform, **centering**, containment and
  no-distortion mutations **have teeth**."

#3 breaks centering. #4 breaks arbitrary-size containment and the byte gate. Neither
can be worked around inside P4.

---

## 5. Task list, in dependency order

The key economy: **#3, #4 and #5 all require a ladder rebuild. Do them in ONE
rebuild, not three.** A rebuild is 30 resolutions × 249 entities through mapshaper.

| # | Task | Blocked by | Notes |
| --- | --- | --- | --- |
| T1 | Establish a green `make check` baseline | — | **not done — see §7.** Must precede any geometry change, or a later red gate is unattributable |
| T2 | Reconcile P2's as-built: `mate plan invalidate` PLAN-013 (authority = RUN-020's result), successor plan bound to committed T1–T5, then a run recording it | T1 | first half of debt `#15`; no backdated receipts |
| T3 | Fix `WKI-37F18A2AA6A5` — fit the `contain` frame **after** visibility removals | T2 | `internal/geometry` (BND-003). `TestVisibilityRemovalsPushTheSilhouetteOffCentre` in `internal/render/layout_test.go` reddens when the gap closes — that is deliberate, it is P2's idiom for a superseded characterization |
| T4 | Fix `WKI-C35A01E965DC` — serve the ladder at any requested long side, not only the preset's own | T2 | effort=large; touches the ladder artifact + the DEC-005 fallback path |
| T5 | Add an all-rungs diagnostic to lodbuild: per entity × band, bytes/IoU/visibility at **every** rung | T2 | this is what makes #5 answerable offline for any candidate cap — see trap 3 for why a naive sweep does not work |
| T6 | One ladder rebuild covering T3 + T4 + the cap chosen in T8 | T3, T4, T5, T8 | changes all 248 assets; DEC-008 digests need re-approval per AM-005 item 5 |
| T7 | `WKI-B05B4B4287A2` — replace lodbuild's wall-clock assertions with work counters, drop `-p 1`, confirm green | T2 | P4 precondition |
| T8 | **OWNER SESSION** — present rendered catalog at several candidate caps + the T5 table; owner picks the byte threshold (#5) and accepts the contact sheet with its named exceptions (#7) | T5 | **write to the user here.** Two decisions, in this order: threshold first, exception list second |
| T9 | Author AM-006 for the byte-budget change | T8 | **an amendment, not a bare decision** — see trap 4 |
| T10 | Complete P2 (audit, `spec.complete`) | T6, T7, T9 | |
| T11 | Reconcile and complete P3: `spec.accept-readiness`, `begin`, plan bound to what T1–T5 built, run, audit, `spec.complete` | T10 | second half of debt `#15` |
| T12 | Decide the browser driver (#6) and stand up P4's browser lane | T11 | **may need the user** if the choice is contentious |
| T13 | Implement P4 — 10 REQ, 7 VAL, 5 AC, `dist/` layout, producer + client docs | T11, T12 | |
| T14 | Implement P5 — 7 REQ, 5 VAL; runs in parallel with P4, gated only on P3 | T11 | |
| T15 | **OWNER** — Windows clean-machine route for P5 VAL-2 (#9) | T14 | **write to the user** |
| T16 | **OWNER** — P6 or non-goal for the crop/selection seam (#10) | T13 | **write to the user** |

Rough size, judgment not calibration (basis: P2 consumed 13 plans and 20 runs):
~10–16 sessions, ~40 % of it P4.

---

## 6. Traps — read before touching anything

1. **`make check` discipline.** Run it detached (`nohup make check > log 2>&1
   < /dev/null & disown`) and **read the log, not the exit code** — a backgrounded
   run was killed by the harness three times in one session, and once more in this
   one (§7). macOS has no `setsid`. **Never pipe it** (`make check | tail` reports
   `tail`'s exit code, always 0 — a red gate was read as green that way). **Never
   run it while anything heavy is on the machine**: the wall-clock brakes measure
   the machine, not the code (`geometry` 440 s contended vs 343 s quiet). Check
   `pgrep -f "go test"` first. Also summarize **package**-level failures, not just
   test-level ones — a `go test -json` package failure carries no `Test` field, so a
   test-only filter reports a timed-out package as a clean run. That bit P2 twice.

2. **Do not add a file under `internal/geometry`** without registering it.
   `TestDiagnosticSourceInventoryIsExactAndBiting` is an exact source inventory and
   reddens on any unregistered file there. Extending `diagnostic.go` in place is
   safe — `diagnosticSourceInventoryDrift` is a self-consistency check, so editing a
   file that is itself in the list creates no fixed-point problem.

3. **A byte-cap "sweep" is NOT N runs of lodbuild.** Verified: lodbuild has no cap
   flag (`cmd/lodbuild/main.go:94-102`), and 2200/7500 are frozen in
   `ParseSilhouetteOracle`'s `want` list (`internal/geometry/silhouette.go:95-98`),
   so each candidate cap would mean editing a frozen contract. And `-ladder-explain`
   cannot help: the search returns on the first passing rung (`ladder.go:351`), so
   rejected rungs are always **finer** than the winner — bytes at **coarser** rungs
   were never computed. Hence T5: one diagnostic pass over all rungs, then every
   candidate cap is answerable offline.

4. **Changing a byte cap is amendment-class.** AM-005 item 5 states "No threshold,
   tolerance, weighting, **byte cap** or budget changes", and its precedent is
   explicit: "Because the ladder recipe changes, DEC-008's approved digests are
   re-approved on the regenerated sheets." So T9 is an amendment (AM-006) plus a DEC
   for the owner's number, not a bare DEC.

5. **Governed-mutation sequencing.** Every amendment transition bumps the epic
   revision, which invalidates any advisor receipt recorded before it. Record the
   receipt for one amendment, retract it, **then** record the next. ADV-008 was
   recorded too early and is a live orphan — immutable, unused, and unrebindable
   (`already exists outside this operation`). Leave it; it is a wart, not a defect,
   and it is documented in `a44fda4`.

6. **SSOT for the byte budget is undeclared, and that is the root of the
   contradiction.** The number lives in four places: QAB-2 in ARCH-001
   (600/1200/2500), `presets/v1.json` `advisory_bytes: 2200`, `path_cap` in
   `lod/silhouette-oracle.v1.json`, and the frozen `want` list in `silhouette.go`.
   **Declare which one is authoritative before T8**, or the contradiction returns on
   the next cycle. Recommendation: the oracle's `path_cap`, with QAB-2 and
   `advisory_bytes` derived from it.

7. **Render and look before calling anything done.** Both P2 defects — Russia as a
   blob, France as a speck — were invisible to every metric (oracle IoU floor 0.40,
   Russia scored 0.92) and were found by looking at output.
   `internal/geometry/cmd/svgshowcase` and `cmd/svgproof` exist for this. A green
   suite is not a rendered card.

8. **Hard constraints, do not violate.** Byte budgets 2200/7500 path and 2500/8000
   file are frozen as *maxima* (lowering the selection target is the open question;
   raising the maxima was explicitly rejected by DEC-006). No country literal or
   per-entity branch anywhere. Pure-Go offline runtime, no CGO/Node/network at
   generation time. P1 corpus bytes immutable. Byte-identical determinism (REQ-8).
   Presentation-free geometry (INV-1) — that invariant is exactly what makes P3's
   theming possible. `q=0.01` quantization untouched.

9. **P3's own small follow-up:** `buildPreview` calls `Generate` inside the style
   loop, so one entity costs five pipeline runs and the twelve-entity ceiling costs
   sixty. Hoist it. The T3 serializer tests generate **once** and re-serialize across
   the matrix for exactly this reason.

---

## 7. The gate baseline is NOT established — start here

A `make check` was started detached in this session and **never reported a
completion**. The background watcher was torn down with no exit record, so the run
may have been killed mid-flight or may have finished unobserved. **Treat the
baseline as unknown, not as green.**

What was observed before the record was lost:

```
ok  cmd/country-map-svg-generator   22.101s
ok  internal/catalog                36.003s
ok  internal/config                  0.355s
internal/geometry (~377 s) and internal/geometry/cmd/lodbuild (~296 s) — NOT observed
```

The two unobserved packages are the slow, brake-carrying ones — precisely where a
red result would matter. Nothing in this session touched Go code (governance
artifacts and Markdown only), so there is no *reason* to expect red; that is an
argument, not evidence.

**First action for whoever picks this up:** re-run per trap 1, detached, on an idle
machine, and read the log to completion. Expected total ~757 s serial. Package times
from P3's close: cmd 17.9 · config 0.2 · render 35.1 · catalog 31.1 · geometry 376.8
· lodbuild 296.0. Do not begin T2 until the whole log is read.

---

## 8. When to write to the user

Four points, and only these:

- **T8** — the byte threshold and the contact-sheet acceptance. Needs rendered
  output in front of them; cannot be delegated or inferred.
- **T12** — if the browser-driver choice turns out to be contentious rather than
  obvious.
- **T15** — the Windows clean-machine route.
- **T16** — P6 or non-goal for the crop/selection seam.

Everything else on the list is agent work.

---

## 9. Housekeeping, unblocked and trivial

- `docs/status.md` is empty ("nothing yet" in both sections) despite being the
  designated first-orientation read of any session. Fill it.
- mate artifact pin is stale: CLI v1.4.3 installed, artifact pinned v1.4.1, latest
  v1.4.3. `mate version status` reports `update-available`. Not required by anything
  on this list.
