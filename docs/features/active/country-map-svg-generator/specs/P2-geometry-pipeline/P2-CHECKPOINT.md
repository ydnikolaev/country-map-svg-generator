# P2 checkpoint — resume point for a fresh session

Last updated at commit `36e1043` (T1 complete). Read this first; it lets a new
session resume without the originating chat.

## The one-line status

The 20-run P2 failure is diagnosed, the fix is designed and independently
verified, and the first product artifact (the ladder) is built and verified.
Two things are open: **one owner decision (Croatia de_facto card)** and **one
mate tooling blocker (a stale fence)**. Neither blocks continuing to build.

## What was wrong (the root cause, confirmed)

DEC-006 (accepted) replaced raw full-coastline deviation with a target-scale
raster **silhouette oracle** as the acceptance boundary for derived candidates.
The production selection path in `internal/geometry/lod.go` never received that
change — it still gates on `matchedBoundaryDeviation` (raw deviation) at
`lod.go:245-246` and `518-520`, references the oracle nowhere, and injects
full-detail source components via `restoreRequiredComponents` (`lod.go:239`), the
exact mechanism DEC-006 rejected. The shipped `Generate()` is also source-only
because `publishedLODTable` is nil (`pipeline.go:11`). Twenty runs tried to make
the superseded predicate pass. Under the correct oracle predicate the catalog
goes from **75 failing outputs to 3**.

Three prior coordinator diagnoses (quantization accounting; a `+Inf` mis-count
branch; a predicate-as-objective) were each refuted by independent review or by
measurement. The lesson, recorded in `docs/harness-backlog.yaml`
(`FND-8386048ACA3C`): on this problem, a hypothesis without a full-corpus
measurement is worthless. Do not re-diagnose from a sample.

## Governance state (all in the epic tracker)

- **DEC-009** accepted — identity rung tried last; a band that cannot satisfy
  topology at fitted scale produces a typed no-artifact outcome (generic rule, no
  country literal); intermediate rungs added between 80 and 64.
- **AM-005** verified (`ADV-006` pass, independently recomputed) — the oracle-gate
  fix. Supersedes AM-003/AM-004 in text.
- **AM-003, AM-004** stuck in `applied`, blocked by `ADV-003`/`ADV-004`. They
  fence P2 and **cannot be retired** — mate has no terminal transition out of
  `applied` for a failed amendment (`FND-1B950D899CE4`, high). This blocks the
  governed lifecycle (`spec P2 amendment-resolve` refuses), which is why T1 was
  built as ordinary engineering commits, not governed runs. **This is the mate
  fix the owner said they would do separately.** Until then, product work
  continues outside the governed run wrapper and is reconciled later (task #15).

## Evidence artifacts (committed, recomputable)

- `readiness/T0-catalog-feasibility.matrix.tsv` (sha `711aecae…`) — 996 rows under
  the oracle predicate: 986 pass / 10 fail. The ground truth.
- `readiness/T0-identity-fallback.evaluation.tsv` (sha `48beef14…`) — identity
  candidate evaluation; recomputes byte-identical.
- `internal/geometry/cmd/t0sweep`, `t0identity` live in the T0 git worktree at
  `/private/tmp/.../scratchpad/t0` (detached HEAD `f03c7fe`), not in main. They are
  the proven reference implementations.

## T1 — done and independently verified (commit `36e1043`)

Built the DEC-006 ladder the 20 runs never built:
- `internal/geometry/lod/ladder.artifact.json` (sha `f0646ed7…`, 10.4 MB) — per
  geometry × band, the finest oracle-passing candidate, 996 rows, 475 dedup
  candidate geometries, full provenance and per-attempt rejection log.
- `internal/geometry/lod/ladder.recipe.json` (sha `050fa04e…`) — binds v2 recipe
  `da3ff9d3…`, oracle `f2c9cd32…`, corpus, Mapshaper pin; `resolutions` is now
  `[]float64` (Mapshaper accepts fractional resolution); one intermediate rung
  `79.63`; `identity_fallback: tried_last`.
- `internal/geometry/cmd/lodbuild/ladder.go` + `ladder_test.go` — the builder and
  its gates, incl. `TestLadderRebuildMatchesCommittedArtifact` (full rebuild,
  ~3 min, byte-identical — the determinism gate).

**Verified by the coordinator, not taken on faith:**
- 986/986 T0-passing rows reproduce with zero drift (independent diff script).
  This is DEC-009's placement proof: identity-last + the 79.63 rung changed no
  existing selection.
- The 10 T0 failures resolve exactly as designed: SH×4 and UM/standard×2 pass via
  identity; HR/un/compact passes via rung 79.63 at 2199 bytes; UM/compact×2 and
  HR/de_facto/compact are typed `no_artifact`.
- Determinism: three byte-identical rebuilds.
- `go build ./...`, `go vet`, `gofmt` clean; fast ladder tests green.

**Final catalog state: 993 pass, 3 no_artifact.** From 75 failing to 3.

## OPEN OWNER DECISION #1 — Croatia de_facto card (decide before T5)

HR (Croatia) is top-200. T1 surfaced that AM-005's HR figures were the **un**
profile's; **de_facto is a different geometry** and worse. Its byte-cap crossing
and its protected-component-loss crossing happen at the *same* simplification
step (verified to .001 resolution precision): rung 79.63 → 2205 bytes (still over
the 2200 cap), rung 64 → the disputed-border component (contribution 111 vs
threshold 110) is lost. No rung at any granularity bridges this without relaxing
the byte cap or the 110 threshold — both out of scope by AM-005/DEC-006.

So Croatia renders in 3 of its 4 profile×band slots: un/card ✓, un/hero ✓,
de_facto/hero ✓, **de_facto/card ✗**. This is the same generic no-artifact
mechanism as UM/compact — but UM is not top-200 and Croatia is, so the owner
should decide consciously. Options, none requiring a country literal:
1. Accept the typed no-artifact for HR/de_facto/card (the site shows whatever it
   uses for a missing decorative card). Cheapest, honest.
2. Group anchor via the DEC-007 override seam — a data override preserving the
   disputed component; keeps the card, needs owner visual review, is data not
   code.
3. Revisit the card byte budget or the 110 contribution threshold — broad,
   rejected by DEC-006 for the budget; a threshold change re-opens DEC-008
   approval. Least preferred.

DEC-009 already established the generic rule; this is only whether HR/de_facto
warrants option 2's override. Recommend option 1 now, option 2 only if the owner
wants that specific card.

## Next step — T2 (do this first in the new session)

Embed the ladder artifact and wire the runtime. Precisely:
1. `go:embed` `ladder.artifact.json`; initialize `publishedLODTable` (or a new
   `LadderTable`) from it in `pipeline.go` so shipped `Generate()` stops being
   source-only. The artifact's `Candidates` map stores `ProjectedLODGeometry` per
   (geometry_id, resolution|"identity"), the same shape as the existing v1 tier
   artifacts, so it embeds the same way.
2. Add a pure-Go, offline (no Node, no network) `make check` gate that recomputes
   topology, protected visibility, IoU/recall and byte caps for **every** committed
   row against the frozen thresholds and bound digests, and reddens on: perturbed
   coordinate, changed threshold, stale corpus/oracle/recipe identity, reordered
   candidates, suppressed omission. This is VAL-6 obligation 2 and it is what makes
   the oracle un-bypassable without paying its cost per CLI run.

Then T3 (rewrite `lod.go` selection to lookup+verify; delete
`restoreRequiredComponents` and the two derived-path deviation gates; route
explicit source/>700 scale to the unchanged DEC-005 path; **emit real SVGs for a
sample of top-200 countries as the empirical proof**), T4 (realign the 5 failing
tests + VAL-6 mutation teeth), T5 (full-catalog SVG proof + owner contact sheet).

Task list: #11 T2, #12 T3, #13 T4, #14 T5, #15 governance reconcile (deferred).

## Hard constraints (from accepted decisions — do not violate)

Byte budgets 2200/7500 path, 2500/8000 file — not to be raised. No country literal
or per-entity branch. Pure-Go offline runtime, no CGO/Node/network at generation.
P1 corpus bytes immutable. Byte-identical determinism (REQ-8). Presentation-free
geometry (INV-1) — CSS/theming/customization is P3, and INV-1 is precisely what
makes it possible. q=0.01 quantization untouched. Thresholds frozen (DEC-008).
