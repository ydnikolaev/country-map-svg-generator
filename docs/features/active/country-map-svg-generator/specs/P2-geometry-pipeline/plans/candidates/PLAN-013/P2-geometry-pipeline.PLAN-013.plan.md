---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "implementation-plan"
id: "PLAN-013"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["P2"]
---
# PLAN-013 — Deterministic coarser/source-tier selection fallback

## Accepted inputs and baseline

Implement from commit
`eb8689c1e5fb33f5bfc14cf46e9095e31aceea94`, tree
`41d2bf9c57d319b5cf87d80ef7dc39b2aa283bcd`.

That commit tracks the previously untracked `internal/geometry/` tree as an
explicit red baseline, so RUN-020 isolates its write footprint by `git diff`
rather than by mtime and artifact hash as RUN-019 was forced to. The five
failing tests it carries — `TestFullCorpusBothProfiles`,
`TestRepresentativeNaturalRatiosAndArbitraryFrames`,
`TestApprovalDigestAndBudgets`, `TestLODSpikeFullCorpus` and
`TestLODAQRUProjectionAlignedCheckpoint` — are exactly this plan's target.

Bind:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- invalidated PLAN-012 candidate body
  `d8b94e354b96079ba3b8dab443b595b15377c6867629dab19584539992d6752e` and accepted
  manifest input digest
  `94a61a5db90d7beae56210ad500f1be1a0daf0eeb72f47defc3740260d2b0fc9`;
- PLAN-012 invalidation receipt
  `fad48eee40abdfa655c7cd258444e96e7637e189b3d62ee7e4c18f9eed630c3c`.

RUN-019 executed PLAN-012 T0A only and was invalidated: its machine gate
`go test ./internal/geometry/cmd/lodbuild` failed on its own premise. Two frozen
representative entities do not fit the preferred tier under the frozen budgets:

- AE/card selects the compact tier (deviation `1.183 < 1.536`), then exact
  serialization produces `5910` bytes against the frozen `2500` hard maximum;
  the pipeline raised `hard_budget_exceeded` at `internal/geometry/lod.go:338`
  with no coarser rung after a byte-hard failure.
- AQ/card compact deviation `1.5377` exceeds the `1.536` tolerance, so production
  already falls back to the `source` tier and emits valid `1386`-byte output,
  but the over-strict assertion `production tier == compact` fails on the label
  mismatch.

A read-only diagnostic proved PLAN-012's variable-component-resolution mechanism
closes neither gap: AE overhead is per-component at the frozen compact resolution
and AQ's `0.0017` deviation gap is uniform across the whole ladder. Both entities
DO satisfy the frozen byte cap at a coarser or `source` tier (AE `1104` B, AQ
`1386` B, both far below `2500`) while preserving mandatory silhouette (SH)
fidelity. The owner accepts source/coarser-tier fallback. The variable-resolution
surface (a `v2` recipe/ladder, a Mapshaper `-simplify variable` expression, and
silhouette variable bands) is therefore ABANDONED; PLAN-013 keeps the frozen v1
recipe, ladder, serializer and oracle byte-identical.

## Requirement and acceptance coverage

| Obligation | Owning task | Terminal evidence |
| --- | --- | --- |
| US-1, US-2, US-3 | T0B | `VR-US-1`, `VR-US-2`, `VR-US-3` |
| US-4 | T2 | `VR-US-4` |
| REQ-1, REQ-6, INV-1 | T0B | `VR-REQ-1`, `VR-REQ-6`, `VR-INV-1` |
| REQ-2, REQ-3, REQ-4, REQ-5, REQ-8 | T0A | `VR-REQ-2`, `VR-REQ-3`, `VR-REQ-4`, `VR-REQ-5`, `VR-REQ-8` |
| REQ-7 | T2 | `VR-REQ-7` |
| AC-1, AC-2, AC-3, AC-4 | T2 | `VR-AC-1`, `VR-AC-2`, `VR-AC-3`, `VR-AC-4` |
| VAL-1, VAL-4, VAL-5 | T2 | `VR-VAL-1`, `VR-VAL-4`, `VR-VAL-5` |
| VAL-2 | T0A | `VR-VAL-2` |
| VAL-3 | T0B | `VR-VAL-3` |

Each normative obligation has exactly one manifest owner. The byte budgets under
REQ-2 and the byte-equivalent determinism under REQ-8 are the direct target of the
fallback rung; the spec mandates byte budgets and aspect/fidelity, not a specific
tier label, so a `source`-tier selection that meets the frozen byte cap and
mandatory SH fidelity satisfies REQ-2, REQ-3 and AC-2.

## Technical approach

Preserve every frozen quantity: centered LAEA projection, natural `tight` and
arbitrary uniform `contain` layout, q `0.01`, the 29-resolution ladder
`[512…8]`, weighted Visvalingam `0.7`, the `1.536` deviation tolerance and its
derivation, path caps `2200/7500`, complete-file maxima `2500/8000`, the
silhouette oracle, the v1 recipe, the source corpus and the serializer. No
budget, threshold or tolerance is recalibrated; no country literal, custom
simplifier, projection or corpus change, CGO, or runtime Node/Mapshaper
dependency is introduced.

Change only production tier selection in `internal/geometry/lod.go`. Today the
fine-to-coarse candidate order is `compact → standard → source` (lod.go:172-192);
the loop already advances past a candidate that fails topology or exceeds the
deviation tolerance (lod.go:239-247), which is why AQ/card already reaches
`source`. The defect is that byte-budget satisfaction is NOT part of selection: a
candidate that passes deviation and topology is committed as `SelectedTier`, and
the frozen byte cap is enforced only afterward in `renderLODRepresentation`
(lod.go:313-341), which raises a terminal `ErrBudget`/`max_path_bytes` at
lod.go:337-338 with no coarser rung. AE/card commits compact, then hard-errors.

Add one deterministic fallback rung:

1. Fold exact serialization and the frozen per-tier byte cap into candidate
   evaluation. A candidate whose exactly serialized path exceeds the frozen
   complete-file maximum for the active preset is treated exactly like a
   deviation or topology failure: record a truthful `<tier>:hard_budget` fallback
   reason and advance to the next coarser candidate. The first candidate that
   satisfies deviation tolerance, topology, mandatory SH fidelity AND the frozen
   byte cap becomes `SelectedTier`. Selection stays a pure function of geometry
   and the frozen ladder, independent of iteration order or host, and byte
   evaluation is done identically on every rerun.
2. Mandatory SH fidelity continues to gate every candidate, including the coarser
   and `source` tiers, through the existing per-component and matched-boundary
   deviation checks and the finalize path (`finalizeSourceLOD`,
   `finalizeGridPhases`); a tier that fits the byte cap but loses mandatory
   fidelity is rejected, never selected.
3. The terminal `ErrBudget`/`max_path_bytes` at lod.go:337-338 is RETAINED, but
   is reached only when no candidate fits — i.e. when even `source` over-caps and
   there is no coarser rung left (for example a `table == nil` forced-source
   input whose source geometry genuinely exceeds the cap). The pipeline never
   hard-errors on over-cap bytes while a coarser-or-source tier fits; it falls
   back deterministically instead.

Correct the over-strict assertions that encode the abandoned "commit the
preferred tier then hard-error" contract, scoped narrowly:

- `internal/geometry/cmd/lodbuild/main_test.go` production-table cases
  (main_test.go:225-244) assert a hard-coded `tier == compact/standard` per
  entity. Correct the expectation to the byte-plus-fidelity-satisfying selected
  tier, accepting a legitimate `source` selection for AQ/card while keeping the
  frozen byte-cap and mandatory-fidelity guards.
- `internal/geometry/lod_test.go` `TestLODHardBudgetDoesNotAdvanceSelectedTier`
  (lod_test.go:84-123) asserts that a compact over-cap keeps `SelectedTier ==
  compact` and raises a terminal budget error. Its invariant is exactly the bug.
  Replace it with a fallback assertion: a compact over-cap that a coarser tier
  fits advances `SelectedTier` to the fitting tier, emits valid within-cap bytes,
  and remains a pure, byte-identical function of selection (selection stays
  independent of the render budget on rerun).

Leave the genuinely-terminal tests unchanged, because they exercise the retained
hard error: `internal/geometry/lod_test.go`
`TestSelectionRenderSeparationAndTypedSourceBudget` (lod_test.go:223-282) and the
`main_test.go` forced-source budget case (main_test.go:258-282), both with
`table == nil`, and `internal/geometry/pipeline_test.go:67-83`
`TestPublicGenerateNilTableMatchesIsolatedSourceInjection`, which observes a
correct `source/source` selection when no ladder is present.

Hash the selection-order and byte-cap contract, the retained terminal condition,
the per-entity selected tiers and the fallback reasons into the recipe/candidate
provenance. If no candidate fits the frozen cap while preserving mandatory
fidelity, fail; do not relax fidelity, visibility or budgets.

## Tasks and completion conditions

### T0A — Fallback rung and assertion correction, representative checkpoint

Implement the deterministic byte-cap fallback rung in
`internal/geometry/lod.go`, fold exact serialization and the frozen per-tier byte
cap into candidate evaluation, and correct the over-strict tier assertions in
`internal/geometry/lod_test.go` and `internal/geometry/cmd/lodbuild`. T0A creates
no product commit.

Run brakes in order:

1. Synthetic teeth: a candidate whose exactly serialized bytes exceed the frozen
   cap advances to the next coarser candidate and never becomes `SelectedTier`;
   the first byte-cap-and-fidelity-satisfying candidate is selected; selection is
   a pure function of geometry and the frozen ladder, byte-identical on rerun and
   independent of the render budget. A forced-source input whose source over-caps
   still raises the retained terminal `ErrBudget`/`max_path_bytes`. Order
   perturbation, budget-dependent selection, a fidelity-losing but byte-fitting
   tier, and deletion of the terminal error must fail typed tests.
2. SH at all 29 frozen resolutions, both consumers and bands: mandatory SH
   fidelity gates every selected tier including coarser and `source`; the frozen
   oracle and recipe are byte-identical. Stop if any selected tier loses
   mandatory fidelity.
3. Representative builds AE/card, AQ/card, RU/card and RU/hero on the frozen
   ladder: AE/card falls back deterministically to a within-cap tier
   (`≤ 2500` B), AQ/card selects `source` (`1386` B) under the `1.5377 > 1.536`
   deviation, RU card/hero hold their tiers, every selected output is within the
   frozen `2200/7500` path and `2500/8000` complete-file budgets and preserves
   mandatory fidelity. Stop on any over-cap without a fitting coarser rung.
4. Mutation teeth: a one-byte-underreported cap, a mutated candidate order, a
   budget-dependent selection, a suppressed fallback reason, a fidelity-losing
   `source` selection, a deleted terminal error, and any frozen budget, tolerance,
   ladder, q, projection, serializer, oracle, recipe, pin or runtime-dependency
   change must fail.
5. Rebuild the representative matrix twice under the unchanged oracle, recipe and
   budgets and require byte-identical output across both rebuilds. Regenerate the
   contact sheet and obtain explicit owner approval of the new sheet digest before
   any downstream wave. Machine failure cannot be waived.

### T0B — Full-catalog fallback validation and clean T0 commit

Start only from T0A machine and owner receipts. T0B may stage but not edit T0A
paths. Generate and validate all 996 ordinary outputs:

- 249 entities × two profiles × card/hero, no gaps or duplicates;
- natural `tight` and tiny/huge/portrait/landscape/square `contain` layouts;
- deterministic tier selection with truthful fallback reasons and no over-cap
  hard error while a coarser-or-source tier fits;
- topology, winding, visibility, mandatory fidelity and the frozen `2200/7500`,
  `<2500/<8000` budgets on every selected tier;
- warm generation at most 120 seconds on YMBPM3;
- normal Go generation performs no Node, Mapshaper, network or external I/O;
- regenerated full-catalog sheet receives owner approval.

Run clean `GOWORK=off go test ./... -count=1`, public nil/source behavior, the
retained typed forced-source budget failure and frozen serializer/oracle/recipe
checks. Only then create one atomic T0 product commit carrying the T0A rung and
assertion corrections plus this catalog gate.

### T1 — Publication and frozen-ladder verification

From clean T0, perform two byte-identical maintainer rebuilds and verify that the
frozen v1 recipe, ladder, generated initializer and oracle are byte-identical and
unchanged by PLAN-013 — the fallback is runtime selection, not a ladder change.
Confirm the public embedded offline selection now exercises the fallback and
emits within-cap output for every entity. Runtime stays embedded, pure Go and
offline.

### T2 — Evidence and P2 closure

Add D3 projection/layout parity, topology and protected-feature mutation teeth,
overrides, markers, byte-budget and natural/contain visual evidence, and evidence
that the tier-fallback selection is deterministic and within cap. Run focused
gates, full Go tests, offline `make check` and a fresh independent audit against
PLAN-013 and both owner receipts.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0A | `tier-fallback-selection-and-assertion-correction` | `internal/geometry/lod.go`; `internal/geometry/lod_test.go`; `internal/geometry/cmd/lodbuild/**` |
| T0B | `geometry-foundation-integrator` | `internal/geometry/lod_catalog_test.go` |
| T1 | `silhouette-ladder-publisher` | `internal/geometry/lod_build_test.go` |
| T2 | `geometry-verifier` | `Makefile`; `internal/geometry/approval_test.go`; `internal/geometry/integration_test.go`; `internal/geometry/mutations_test.go`; `internal/geometry/oracle_test.go`; `internal/geometry/testdata/approval/**`; `internal/geometry/testdata/d3/**`; `internal/geometry/testdata/mutations/**` |

Leases do not overlap. `internal/geometry/lod.go` and its unit gate
`internal/geometry/lod_test.go` are owned solely by T0A, co-located with the
fallback rung; the `internal/geometry/cmd/lodbuild/**` lease carries the corrected
`main_test.go` production-table assertions. Oracle, v1 recipe, ladder and
serializer stay byte-frozen and are leased by no task. No implementer writes
`data/**`, `internal/catalog/**`, public `cmd/**`, silhouette or ladder-generation
sources, lifecycle documents or another lease. Never stage `node_modules`,
temporary artifacts or unrelated user work.

## Dependencies and execution waves

Critical path: `W0/T0A synthetic → SH → representative machine gate → owner
approval → W1/T0B full catalog + clean commit → W2/T1 frozen-ladder verification
→ W3/T2 evidence and audit`. Concurrency is one. A wave starts only after
predecessor evidence is accepted.

## Validation plan

- `geometry-tier-fallback-and-sh-brake`;
- `geometry-fallback-mutations-and-terminal-budget`;
- `geometry-fallback-integrated-budget-visual-performance`;
- `geometry-t0-clean-checkout`;
- `geometry-frozen-ladder-rebuild-check`;
- `geometry-frozen-ladder-public-wiring`;
- `geometry-d3-parity`, `geometry-overrides`, `geometry-markers`;
- `GOWORK=off go test ./... -count=1` and offline `make check`;
- independent terminal audit against PLAN-013 and both owner receipts.

## Deviation and amendment policy

`p2-tier-fallback-source-selection-v1` permits only folding exact serialization
and the frozen per-tier byte cap into fine-to-coarse tier selection, so a
preferred tier that hard-fails the frozen byte cap or the deviation tolerance is
skipped for the first coarser-or-source tier that fits the frozen byte cap and
preserves mandatory SH fidelity, plus the narrowly-scoped correction of the
over-strict `tier == compact/standard` assertions. The terminal
`ErrBudget`/`max_path_bytes` is retained for the no-fitting-tier case.

Any budget, threshold, deviation tolerance or its derivation, metric, band
boundary, resolution ladder or order, weighting, projection, q, serializer,
oracle, recipe, tool pin, runtime dependency or selection input requires another
plan. Country branches, country literals, manual vertices, custom simplification,
a variable-resolution or v2 recipe/ladder, MakeValid, union/buffer repair, CGO
and post-hoc SVG optimization remain forbidden.

## Commit worktree and integration policy

T0A creates no product commit. T0B creates one atomic T0 commit only after both
owner gates and clean-checkout validation, carrying the T0A rung and assertion
corrections plus the catalog gate. T1 and T2 each create one later atomic commit
after full gates. Stage only exact leased paths. Lifecycle artifacts remain
separate coordinator commits.

## Rollback and recovery

On the first T0A brake failure, restore the exact pre-run scratch inventory and
retain governed evidence naming entity, preset, requested tier, selected tier,
fallback reasons, deviation, fidelity and bytes. Never relax frozen policy. On
T0B failure return to the accepted representative checkpoint. T1 verifies from a
temporary directory only after two byte-identical rebuilds. Revert completed
commits in T2 → T1 → T0 order.

## Completion and handoff

P2 closes with deterministic tier-fallback selection across all-283 artifacts and
996 ordinary outputs, both owner-approved sheets, every selected tier within the
frozen `2200/7500` and `2500/8000` budgets with mandatory SH fidelity preserved,
truthful fallback provenance, natural/contain layouts, the retained terminal
budget error only when no tier fits, at most 120-second offline generation, clean
T0/T1 checkouts, a byte-identical frozen v1 ladder and independent audit.

P3 receives the selected tier, requested tier, fallback reasons, deviation and
fidelity metrics, path bytes, transform and optional markers. No presentation
contract is added.

Estimated remaining effort: 4–12 agent-hours, likely 7, medium confidence.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
