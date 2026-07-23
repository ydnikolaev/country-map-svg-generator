---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "accepted-plan-revision"
id: "PLAN-002"
epic: "country-map-svg-generator"
spec: "P2"
status: accepted
profiles: []
concerns: []
inputs: ["P2", "DEC-003", "DEC-004", "AM-001", "AM-002", "RUN-001-RESULT", "CTR-001"]
---
# PLAN-002 — Precomputed LOD geometry pipeline

## Accepted inputs and baseline

This successor implements P2 at baseline commit
`50f56191c21dc84530f579c4896a59d3a3651476`, tree
`5d812dc559e47c8a33b67e99fa61d0888349473c`. It is bound to:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- DEC-003 plus refining DEC-004
  `3efec49a55fe7289ea528233280671e43db5f114c24586c1780dba98ba2667e5`;
- accepted natural-layout amendments AM-001
  `ab284fd6e2a79c4ef7a8c45ab05fe807d217ce7a8a178be354a44d26ef4ac80a`
  and AM-002
  `ec2203d7132e467c7b1519d765c4790d583c29691311aad208a2ce6092523487`;
- interrupted RUN-001 result
  `2d704d25d76d2646f5cfb8061623ed19bbe43ebb24af5ad60e711c6d9e92ac87`;
- P1 corpus identity
  `sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`,
  manifest SHA-256
  `227374e94910f078623b3c643085d39fa094b4e1ae20a0ffc01cb2d4a19ee616`,
  and geometry SHA-256
  `f3d5d80ac8b9650cb75f68f46593c24476a979855875842bc3274d9d928e7744`.

PLAN-001 is invalidated, not edited. Its uncommitted Go draft is reusable scratch,
not accepted baseline truth; every reused byte remains subject to this plan's
leases and gates.

## Requirement and acceptance coverage

| Obligation | Tasks | Evidence |
| --- | --- | --- |
| US-1, AC-1 | T0, T3, T4 | approved card/hero raster sheet, final path budgets and structural receipt |
| US-2, REQ-1, REQ-2, REQ-8, AC-2 | T2, T4 | natural `tight`, arbitrary `contain`, scale-only LOD selection and deterministic goldens |
| US-3, REQ-6, AC-3 | T2, T3, T4 | bounded override and generic protected-restoration fixtures |
| US-4, REQ-3, REQ-4, AC-4 | T0–T4 | source binding, deviation/topology teeth, rebuild proof and D3 parity |
| REQ-5 | T2, T3, T4 | exact precomputed disappearance, restoration and runtime removal provenance |
| REQ-7 | T2, T4 | unchanged common marker transform and anomaly gates |
| INV-1 | T2, T3 | presentation-free LOD and command result |
| VAL-1 | T2, T4 | pinned D3 parity |
| VAL-2 | T2–T4 | LOD/topology/protected/softening mutation suite |
| VAL-3 | T0, T4 | multi-scale visual, arbitrary-layout and hard-budget gate |
| VAL-4 | T2, T4 | override gate |
| VAL-5 | T2, T4 | marker gate |

## Technical approach

Keep the accepted full-source normalization, centered Lambert azimuthal equal-area
projection, adaptive subdivision, natural-aspect fitting and common marker
transform. Compute the projector, natural bounds, `tight` viewBox or `contain`
transform and effective fitted geometry scale `d` from the full P1 geometry
before any LOD choice.

A maintainer-only builder invokes Mapshaper 0.7.44, pinned by npm integrity
`sha512-3Cx+IABMXt1G28Y8J7oalW5P5VYyt1vHz5FO+KkV51EHnJioo9h9maO+u+4IyCPZ9Mh1hqchgomFmc68GtFwQQ==`
and shasum `e08fc40d50347698ea07d2f82ead42a98b1c62f0`. Each of the 283 P1 geometry
IDs is imported as one whole MultiPolygon. Weighted Visvalingam
(`weighting=0.7`), default intersection rollback, `keep-shapes`, and `clean`
produce two versioned geographic tiers:

- `compact`: initial `resolution=48`, selected initially for `d <= 240`;
- `standard`: initial `resolution=128`, selected initially for `240 < d <= 700`;
- `source`: exact P1 geometry for `d > 700` or any failed guard.

T0 may calibrate compact resolution only within `[32,64]`, standard only within
`[96,256]`, and the two generic scale thresholds only within the representative
reference ranges. It publishes the literal final values in the immutable v1
recipe. Leaving those ranges, adding a tier, inspecting ISO/profile/preset/frame
shape, or changing algorithms requires a new plan revision.

Before tier selection, compute the accepted retention decision from the full
projected source: largest component, every anchor-bearing component, stable
minimum-parts selections, and every component above the resolved area/visibility
threshold. The runtime projects the selected tier with the exact full-source
projector and fit transform, then restores every required source component that
the tier omitted. A missing interior ring may remain absent only when the same
explicit area/visibility policy permits it; otherwise restore the ring or advance
the tier.

Validate source/recipe/tool/artifact identities, polygon topology, winding,
non-empty geometry, fixed `0.01px` command quantization and a component/ring
matched symmetric maximum boundary-deviation guard. Classify permitted missing
components/rings before comparison. Densify both full and selected fitted
boundaries to at most `0.05px` and require
`max(full→LOD, LOD→full) <= resolved simplification tolerance`. Failure advances
`compact → standard → source`; it never increases tolerance. Every precomputed
disappearance, restored component/ring, runtime removal and fallback names exact
source polygon/ring provenance in additive CTR-002 LOD metadata.

Softening is presentation-free geometry and remains optional, conservative and
topology-checked. A generic resolved `MaxPathBytes` field carries the hard command
budget: `card` defaults to hard `2,500`/advisory `2,200`, `hero` to hard
`8,000`/advisory `7,500`, bounded caller overrides are allowed, and zero means no
ceiling for an unpresetted request. Preset names never reach the fallback logic.
Flatten curves below `min(0.05px, tolerance/8)`. If softened commands exceed
`MaxPathBytes` or fail topology, emit linear commands with
`softening_budget_fallback` or
`softening_topology_fallback`. P3 may add round stroke joins as style, but P2
never depends on CSS for correctness.

Node, npm and Mapshaper exist only under the LOD rebuild command. Shipped Go
binaries, ordinary CLI generation and `make check` neither execute nor require
them. No SVGO or generic post-hoc SVG optimizer is introduced.

## Tasks and completion conditions

1. **T0 — Calibrate and accept the LOD spike.** Pin the maintainer toolchain,
   finish the production-capable builder, and temporarily build compact/standard
   candidates for all 283 geometry IDs. Run structural, protection,
   determinism and hard-budget checks over all 249 entities, both profiles and
   both presets; keep human visual approval representative. Compare final Go
   paths at `d=90,128,160,240,500,640,700,1024,2048` with
   softening on/off, visual digests, topology, deviation, protected restoration,
   artifact size, generation time and two clean rebuilds. Done only when card
   hard maximum is `2,500` path bytes, hero hard maximum is `8,000`, every
   protected/minimum-parts case is generic, and the owner-approval sheet is
   accepted. Stop rather than implement T1 if any invariant fails. T0 output
   freezes literal v1 recipe values and hashes.
2. **T1 — Publish canonical LOD artifacts.** Read the T0-frozen builder/tool/recipe
   and publish compact/standard artifacts for all 283 geometry IDs plus one manifest.
   Bind exact P1, recipe, package-lock, Mapshaper integrity/git head, tier hashes,
   coverage counts and component/ring disappearance provenance. Two clean rebuilds
   must be byte-identical. Committed artifacts are validated without invoking
   Node.
3. **T2 — Integrate scale-driven runtime selection.** Reconcile the reusable
   RUN-001 projection/layout/marker draft, replace runtime RDP with LOD lookup,
   scale selection, component/ring-matched symmetric deviation and topology
   guards, restoration and lossless full-source fallback. Additive CTR-002
   provenance names schema/algorithm/LOD
   versions and fallback diagnostics. Natural viewBoxes and arbitrary frames
   remain regression-compatible.
4. **T3 — Finish retention, commands and bounded softening.** Complete exact
   removal/restoration accounting, fixed-grid canonical commands, override
   bounds, optional markers and topology-checked selective softening with
   deterministic linear fallback. No preset or country branch is permitted.
5. **T4 — Allocate complete evidence.** Add full-corpus both-profile coverage,
   D3 fixtures, natural-ratio and arbitrary-frame goldens, scale-only selection
   mutations, protected restoration, budget/visual approval, deterministic
   rebuild validation and the offline project ceiling. Done with green focused
   gates, green `make check`, approved visual evidence and an independently
   auditable product commit.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0 | `lod-spike` | `internal/geometry/cmd/lodbuild/**`, `internal/geometry/lod/tool/**`, `internal/geometry/lod/v1.recipe.json`, `internal/geometry/testdata/lod-spike/**` |
| T1 | `lod-artifact` | `internal/geometry/lod/v1.manifest.json`, `internal/geometry/lod/v1.compact.json`, `internal/geometry/lod/v1.standard.json`, `internal/geometry/lod_build_test.go` |
| T2 | `lod-runtime` | `go.mod`, `go.sum`, `internal/geometry/{adapter.go,adapter_test.go,diagnostic.go,fit.go,fit_test.go,lod.go,lod_test.go,marker.go,marker_test.go,math.go,model.go,model_test.go,normalize.go,normalize_test.go,override.go,override_test.go,pipeline.go,pipeline_test.go,preset.go,preset_test.go,projection.go,projection_test.go,simplify.go}`, `internal/geometry/overrides/**`, `internal/geometry/presets/**` |
| T3 | `lod-path` | `internal/geometry/{path.go,path_test.go,retain.go,retain_test.go,soften.go,soften_test.go}` |
| T4 | `lod-verifier` | `Makefile`, `internal/geometry/{approval_test.go,integration_test.go,mutations_test.go,oracle_test.go}`, `internal/geometry/testdata/{approval,d3,mutations}/**` |

No implementer writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle docs,
or another task's lease. The coordinator owns plans, runs, evidence and audit.

## Dependencies and execution waves

Critical path is T0 → T1 → T2 → T3 → T4. T0 is a brake, not a prototype that
silently rolls into production, and T1 dispatch requires its recorded passing
receipt. A failed T0 creates no product commit and terminates RUN-002 blocked.
T2 may reuse RUN-001 scratch only after T1's artifact contract is frozen. No
parallel mutation is useful until T2; T3 and D3-fixture preparation inside T4
may then run in separate non-overlapping leases.

## Validation plan

- `geometry-lod-spike`: bounded T0 command; emits exact final-path bytes, visual
  digest, topology/deviation/protection, artifact-size/runtime and rebuild
  determinism. The 2,500/8,000 maxima are hard; preset 2,200/7,500 remain advisory.
- `geometry-lod-rebuild-check`: two temporary rebuilds from exact P1 and pinned
  package lock, then byte comparison with committed artifacts.
- `geometry-lod-runtime`:
  `GOWORK=off go test ./internal/geometry -run
  'TestLOD(Artifact|SourceBinding|Selection|Deviation|Topology|Protected|MinimumParts|Fallback|Determinism)' -count=1`.
- `geometry-topology-teeth`:
  `GOWORK=off go test ./internal/geometry -run
  'Test(Mutation|Topology|Protected|Soften|Budget)' -count=1`.
- `geometry-d3-parity`, `geometry-overrides`, and `geometry-markers` retain their
  PLAN-001 focused meanings.
- `make check` validates committed LOD artifacts, all Go tests, full corpus,
  approval sheet and budgets offline. It must fail if a normal gate invokes Node,
  npm, Mapshaper, network, SVGO or an uncommitted source.

Teeth mutate P1/corpus/recipe/tool/tier hashes; remove a geometry ID; corrupt a
ring; inspect ISO/preset/frame shape during selection; suppress a fallback; lose
a protected/minimum part; exceed deviation; make rebuild bytes drift; or accept a
softened over-budget/self-intersecting path. RU, CL/AR and AU prove natural aspect;
portrait, landscape, square, `19x19` and `10000x12345` prove uniform containment.

## Deviation and amendment policy

`p2-lod-local-how-v1` allows helper/file splits within leased roots, equivalent
hashing/encoding details and T0 calibration inside the declared ranges. A new
tier, dependency, tool version, algorithm, calibration range, budget maximum,
selection input, `MaxPathBytes` bounds, provenance field family, source root or
runtime external process requires plan revision. A change to natural/contain
semantics, protected policy, P1 authority, presentation-free CTR-002, or
single-binary/offline runtime requires an amendment or successor decision. Never
recover with country branches, silent tolerance/precision escalation, manual
geography or weaker gates.

## Commit worktree and integration policy

Each accepted task lands as one conventional atomic product commit after its
focused gate. Stage only exact leased paths. RUN-001 scratch has no authorship
authority: before the first product commit, reconcile every reused file against
this plan and remove probe-only code. The coordinator verifies expected parent,
actual path scope and gate receipt before the next task. Lifecycle evidence is
never included in a product commit. T0 may land only if its brake passes; otherwise
finish the run blocked with no production artifact commit.

## Rollback and recovery

The source tier is always the semantic rollback. A bad compact tier advances to
standard; a bad standard tier advances to exact P1 and records why. Artifact
rebuilds use explicit temporary directories and publish only after complete hash,
coverage and topology validation. A P1, recipe, Mapshaper, Node/npm or threshold
change creates a new LOD version and full comparison; accepted v1 bytes are never
regenerated in place. Revert product commits in reverse task order and select the
prior algorithm/LOD manifest pair. Failure returns to decision/plan review, never
to tolerance escalation or manual redrawing.

## Completion and handoff

P2 closes only after T0 acceptance, deterministic all-283 LOD artifacts, both
profiles for all 249 entities, exact provenance/restoration, green natural and
arbitrary layouts, fixed-grid commands, hard budgets, owner-approved visual sheet,
offline `make check`, exact run result and fresh independent audit.

P3 receives additive CTR-002 LOD provenance plus the unchanged general
layout/quality API: natural `tight`, distortion-free arbitrary `contain`,
data-driven presets, canonical presentation-free commands, stable diagnostics,
metrics, removals and optional transformed capital markers. P3 serializes and
styles; it does not reinterpret LOD policy or optimize bloated geometry after the
fact.

Estimated effort is 14–24 agent-hours, likely 18, medium confidence. T0 bounds the
remaining uncertainty before artifact and runtime work.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
