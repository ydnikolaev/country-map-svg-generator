---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-020"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-020"
status: draft
profiles: []
concerns: []
inputs: ["PLAN-013", "CTX-RUN-022"]
---
# RUN-020 — Byte-cap tier-selection fallback and assertion correction

## Authority and accepted inputs

Execute accepted PLAN-013
`0efd5be8f545f28edd11262d39cd0f96bf765f2698db93fb932ec36187295dbe`
with manifest
`00e8d0ffd3e5f84d4ae434f2c46e0fa9d87250e177474d283cfc8eae3bed9938`,
from baseline commit `eb8689c1e5fb33f5bfc14cf46e9095e31aceea94`, tree
`41d2bf9c57d319b5cf87d80ef7dc39b2aa283bcd`.

Bind RUN-019 result
`c1ac9c1a2e867a5be8784881342852cf44b9ec7981f20a0169fea1ff032fe7db`,
PLAN-012 invalidation receipt
`fad48eee40abdfa655c7cd258444e96e7637e189b3d62ee7e4c18f9eed630c3c`,
frozen silhouette oracle
`f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`
and CTX-RUN-022
`baa8ed04a3825f41a86f4d0f20d58b45ac48b4257ec52196fd945481f6cbfd1e`.

Unlike every prior P2 run, the baseline is tracked in git. Isolate your write
footprint with `git status` and `git diff` against `eb8689c1`, not by mtime or
artifact hash.

## Objective and completion boundary

Execute PLAN-013 T0A only.

Fold exact serialization and the frozen per-tier byte cap into LOD candidate
evaluation in `internal/geometry/lod.go`, so a candidate whose exactly
serialized path exceeds the frozen complete-file maximum for the active preset
is treated exactly like a deviation or topology failure: record a truthful
`<tier>:hard_budget` fallback reason and advance to the next coarser candidate.
The first candidate satisfying deviation tolerance, topology, mandatory
silhouette fidelity AND the frozen byte cap becomes `SelectedTier`. Correct the
over-strict tier assertions that encode the abandoned contract. Pass the
synthetic, silhouette and representative brakes, rebuild the representative
matrix twice, and stop with the new contact sheet for owner review.

T0A creates no product commit. Do not enter T0B, do not generate the 996-output
catalog, do not wire runtime publication, do not stage or commit product bytes.
T0A is complete only when every machine gate passes and the exact new
sheet/receipt hashes are available for owner approval.

## Scope ownership and footprint

Exclusive writes, derived from the accepted manifest `tasks[T0A].writes`:

- `internal/geometry/cmd/lodbuild/**`;
- `internal/geometry/lod.go`;
- `internal/geometry/lod_test.go`.

Off-limits, non-exhaustive: `go.mod`, `go.sum`, `.gitignore`, `data/**`,
`internal/catalog/**`, `internal/geometry/silhouette.go`,
`internal/geometry/silhouette_test.go`,
`internal/geometry/testdata/silhouette-oracle/**`,
`internal/geometry/lod/**` (recipes, ladder, build tool),
`internal/geometry/lod_catalog_test.go`, `internal/geometry/lod_build_test.go`,
`internal/geometry/approval_test.go`, `internal/geometry/integration_test.go`,
`internal/geometry/mutations_test.go`, `internal/geometry/oracle_test.go`,
`internal/geometry/pipeline_test.go`, `Makefile`, public `cmd/**`, every
lifecycle document under `docs/features/**`, and any other task's lease. Never
stage `node_modules` or temporary artifacts.

## Context and doctrine pack

Load CTX-RUN-022 before implementation. It carries current P1 source-corpus
authority. PLAN-013 and this brief carry the exact geometry policy and runtime
fences. The user-visible intent remains lightweight, soft-organic,
aspect-ratio-correct country SVGs for arbitrary viewports; `detail` remains
backlog.

## Tasks

1. Read `internal/geometry/lod.go` end to end before editing. The candidate
   order is built at lod.go:186-192 (`compact → standard → source`); the loop
   already advances past a candidate failing topology (lod.go:239-243) or
   deviation tolerance (lod.go:244-247); the frozen byte cap is enforced only
   later in `renderLODRepresentation`, which raises a terminal
   `ErrBudget`/`max_path_bytes` at lod.go:337 with no coarser rung. Confirm this
   reading empirically before changing anything.
2. Add the deterministic byte-cap fallback rung. Byte evaluation must use the
   same exact serialization as production output and the same frozen per-preset
   complete-file maximum. Selection must stay a pure function of geometry and
   the frozen ladder — byte-identical on rerun, independent of iteration order,
   host, and of the render-time budget parameter.
3. Keep mandatory silhouette fidelity gating every candidate, including coarser
   and `source` tiers, through the existing per-component and matched-boundary
   deviation checks and the finalize path (`finalizeSourceLOD`,
   `finalizeGridPhases`). A tier that fits the byte cap but loses mandatory
   fidelity is rejected, never selected.
4. Retain the terminal `ErrBudget`/`max_path_bytes` at lod.go:337. It must still
   fire when no candidate fits — for example a `table == nil` forced-source
   input whose source geometry genuinely exceeds the cap.
5. Correct only the over-strict assertions, and state for each why the change is
   not a weakened gate:
   - `internal/geometry/cmd/lodbuild/main_test.go:225-244` production-table cases
     assert a hard-coded `tier == compact/standard` per entity. Expect the
     byte-plus-fidelity-satisfying selected tier instead, accepting a legitimate
     `source` selection for AQ/card, while keeping the frozen byte-cap and
     mandatory-fidelity guards.
   - `internal/geometry/lod_test.go:84-123`
     `TestLODHardBudgetDoesNotAdvanceSelectedTier` asserts that a compact
     over-cap keeps `SelectedTier == compact` and raises a terminal budget
     error. That invariant is exactly the defect. Replace it with a fallback
     assertion: a compact over-cap that a coarser tier fits advances
     `SelectedTier` to the fitting tier, emits valid within-cap bytes, and
     preserves the rerun-equality (selection purity) assertion that test already
     carries.
   Leave unchanged: `lod_test.go:223-282`
   `TestSelectionRenderSeparationAndTypedSourceBudget`, the
   `main_test.go:258-282` forced-source budget case, and `pipeline_test.go:67-83`
   `TestPublicGenerateNilTableMatchesIsolatedSourceInjection`. They exercise the
   retained terminal error and correct `source/source` selection.
6. Run the brakes in order and stop at the first failure:
   1. Synthetic teeth: an over-cap candidate advances and never becomes
      `SelectedTier`; the first byte-cap-and-fidelity-satisfying candidate is
      selected; selection is byte-identical on rerun and independent of the
      render budget; a forced-source over-cap still raises the retained terminal
      error. Order perturbation, budget-dependent selection, a fidelity-losing
      but byte-fitting tier, and deletion of the terminal error must each fail a
      typed test.
   2. Silhouette at all 29 frozen resolutions, both consumers and bands:
      mandatory fidelity gates every selected tier including coarser and
      `source`; frozen oracle and v1 recipe byte-identical.
   3. Representative builds AE/card, AQ/card, RU/card, RU/hero on the frozen
      ladder: AE/card falls back deterministically to a within-cap tier
      (`<= 2500` B), AQ/card selects `source` (`1386` B) under the
      `1.5377 > 1.536` deviation, RU card/hero hold their tiers, and every
      selected output is within the frozen `2200/7500` path and `2500/8000`
      complete-file budgets with mandatory fidelity preserved.
   4. Mutation teeth: a one-byte-underreported cap, a mutated candidate order, a
      budget-dependent selection, a suppressed fallback reason, a
      fidelity-losing `source` selection, a deleted terminal error, and any
      frozen budget, tolerance, ladder, q, projection, serializer, oracle,
      recipe, pin or runtime-dependency change must fail.
   5. Rebuild the representative matrix twice under the unchanged oracle, recipe
      and budgets; require byte-identical output across both rebuilds.
      Regenerate the contact sheet and stop for owner approval of the new sheet
      digest.

## Validation and evidence

Required evidence:

- `go test ./internal/geometry/cmd/lodbuild -count=1` full output, with
  `TestLODSpikeFullCorpus` and `TestLODAQRUProjectionAlignedCheckpoint` passing;
- `go test ./internal/geometry/... -count=1` full output. The five tests red at
  baseline (`TestFullCorpusBothProfiles`,
  `TestRepresentativeNaturalRatiosAndArbitraryFrames`,
  `TestApprovalDigestAndBudgets`, `TestLODSpikeFullCorpus`,
  `TestLODAQRUProjectionAlignedCheckpoint`) must be green, or — where a test is
  outside the T0A lease and still red — named explicitly with its reason;
- focused synthetic and fallback mutation tests, each shown to fail without the
  rung;
- per-entity requested tier, selected tier, truthful fallback reasons,
  deviation, fidelity metrics and exact byte counts for AE/card, AQ/card,
  RU/card and RU/hero;
- frozen oracle, v1 recipe and serializer hashes shown unchanged;
- two representative builds with byte-identical manifest, receipt and sheet,
  plus the new contact sheet path and digest;
- `git status --porcelain` and `git diff --name-only eb8689c1` proving the
  changed-path set equals the lease.

## Constraints and forbidden actions

Keep centered LAEA, q `.01`, the 29-resolution ladder `[512…8]` and its order,
weighted Visvalingam `.7`, the `1.536` deviation tolerance and its derivation,
path caps `2200/7500`, complete-file maxima `2500/8000`, the silhouette oracle,
the v1 recipe, the source corpus, the serializer, and the pure-Go offline
runtime.

Permitted deviation is exactly `p2-tier-fallback-source-selection-v1`: folding
exact serialization and the frozen per-tier byte cap into fine-to-coarse tier
selection, plus the narrowly-scoped assertion corrections named above.

Forbidden: any budget, threshold, deviation tolerance or its derivation, metric,
band boundary, ladder or order, weighting, projection, q, serializer, oracle,
recipe, tool pin, runtime dependency or selection-input change; country branches
or literals; manual vertices; custom simplification; a variable-resolution or v2
recipe/ladder; MakeValid; union/buffer repair; CGO; post-hoc SVG optimization.
Owner approval cannot waive a machine failure, and a machine failure cannot be
worked around by relaxing an assertion outside the two named corrections.

## Stop escalation and resume conditions

Stop immediately, leave the tree inspectable, and report if:

- no candidate tier fits the frozen cap while preserving mandatory fidelity for
  some entity/preset;
- the fallback would require relaxing fidelity, visibility or a budget;
- a correction beyond the two named over-strict assertions appears necessary;
- outputs drift between identical runs;
- a change outside the lease appears necessary; or
- the representative machine gate fails.

Return exact entity, preset, requested tier, selected tier, fallback reason,
deviation, fidelity and byte evidence, plus the narrowest next decision. Do not
commit, do not issue the project gate, do not finish the run, do not audit, do
not complete the spec — those are coordinator-owned. Resume only after the
coordinator records machine evidence and the owner explicitly approves the new
contact sheet.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
