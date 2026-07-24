---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "implementation-plan"
id: "PLAN-012"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["P2"]
---
# PLAN-012 — Component-relative variable-resolution lineage candidates

## Accepted inputs and baseline

Implement from commit
`8bfeeae9f46fd5a93daf7ee2aaf86356cc468891`, tree
`33a25c074df75e9382c8de2099ebbb054f8894f1`.

Bind:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- invalidated PLAN-011 accepted revision/manifest
  `d9dac64bef9f01d32f8ca4a5905b57388f0849e531e566ca2d39a04696d9e799` /
  `75fba90898b5519fb1578a347e13a9ea79b8d1f7e0de2ca5bfb3c45f6e017e1d`;
- RUN-018 result/agent result
  `0dcf517debb6bb0e04c853b9734d3ea38622568ba9cdcdcbe8ca3937abe9f1b0` /
  `af33d1730556bcf2cc4d3349ff1010882141c9554c1d6c94c7e86751da69cae6`;
- PLAN-011 invalidation receipt
  `0fa6cbbe38ba6d12db453bbf1777b3881ce8b6910143bedd057f516eb8fdb420`;
- RUN-017 result
  `ff99d94cbcefd9b10943f3d81ea182d035862ab7d62a302b0a2a5666c6489f07`;
- byte-frozen oracle
  `f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`;
- predecessor v2 recipe
  `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`.

RUN-018 proved that explicit component lineage restores all four SH identities,
but country-relative simplification still reduces mandatory source component 0
to a triangle at every frozen resolution: recall
`0.41904286516503403 < 0.5`, IoU `0.4162379919429811`. Projection, parser and
oracle remain exonerated.

Pinned Mapshaper `0.7.44` local help and implementation expose
`-simplify variable`, evaluating a resolution expression per feature against
layer bounds. Use that framework feature; no custom simplifier is authorized.
Reconcile the existing product scratch by exact pre/post inventory and never
recreate or broadly delete it. `detail` remains deferred backlog.

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

Each normative obligation has exactly one manifest owner.

## Technical approach

Preserve centered LAEA, natural `tight`, arbitrary uniform `contain`, q `0.01`,
fine-to-coarse candidate order, pure-Go offline runtime, source/custom
semantics, serializer bytes, scale boundaries `240/700`, path caps
`2200/7500`, complete maxima `2500/8000`, contribution thresholds `110/900`,
global IoU/recall `.40/.42`, and component recall `.50`.

Replace only candidate-source simplification semantics:

1. Split each projected MultiPolygon into temporary Polygon features. Encode
   `lineage-features/v2` with geometry ID, source polygon index, source ring
   count and canonical projected polygon digest.
2. For each geometry and band, resolve every entity/profile consumer including
   `identical_to`. The generic mandatory union contains any component visible,
   dominant, anchored or group-anchored for any consumer. Runtime entity,
   profile, preset and frame never affect candidate identity.
3. Submit every component for one geometry in one FeatureCollection. Require
   finite positive `country_long_side` and `component_long_side`; invalid values
   fail typed validation before Mapshaper.
4. Invoke pinned weighted-Visvalingam with unchanged weighting `.7`, planar,
   keep-shapes, clean and intersection repair, adding only `variable`.
   Precompute a deterministic numeric `component_resolution` property:

   - mandatory:
     `base_resolution * country_long_side / component_long_side`;
   - optional: `base_resolution`.

   Mapshaper converts resolution against country layer bounds, so the mandatory
   expression yields component interval
   `component_long_side / (2 * base_resolution)`. Optional components retain
   the predecessor country interval.
5. Require exactly one non-empty Polygon for every lineage key before policy
   pruning. Reject missing, duplicate, unknown, mutated or cross-geometry keys.
   Output order has no authority. Match surviving holes by containment/maximum
   overlap, then source ring order/digest; reject invented, promoted, inverted,
   escaped or cross-lineage holes.
6. Reassemble in source order. A contribution-ranked identity below the band
   threshold may be omitted only as `protected_subscale`; other optional parts
   require every consumer to classify them subscale/non-identity and use
   `subscale`. Mandatory parts cannot be pruned.
7. Apply q, topology, winding, per-component/global fidelity and exact
   serialization. Select the first fine-to-coarse resolution passing every
   consumer and aggregate byte cap. Final assembled-country validation, not a
   hidden per-country budget, normalizes component-local detail to the `240/700`
   display bands.

Hash candidate-source version, exact variable-expression identity, numeric
feature table, ordered lineage, mandatory-union proof, per-component outputs
and omissions into recipe/candidate/publication identities. If no frozen
resolution passes, fail; do not relax fidelity, visibility or budgets.

## Tasks and completion conditions

### T0A — Variable-resolution source and representative checkpoint

Implement `lineage-features/v2`, variable resolution, hole validation and
mandatory-union reassembly. T0A creates no product commit.

Run brakes in order:

1. Synthetic two-component geometry proves empirically on pinned `0.7.44` that
   the derived variable expression matches a single-component constant-control
   interval and differs from country-relative constant resolution. Reverse
   input/output feature order and require identical bytes. Constant-mode
   substitution, inverted ratio, zero/non-finite size or expression drift must
   fail typed tests.
2. SH at all 29 frozen resolutions, both consumers and bands: four exact
   lineage outputs before pruning, component 0 recall at least `.50`, nonzero
   contribution, topology and exact budgets. Stop if none passes.
3. CA, ID and KI plus legacy Indonesia `120/121` control and new candidates:
   mandatory-union consumer-order invariance, truthful optional omissions,
   mandatory fidelity and aggregate budgets. Stop on unavoidable overflow or
   country-specific pressure.
4. Mutation teeth: missing/duplicate/unknown/cross-wired lineage; index/digest
   mutation; feature/consumer reordering; mandatory/optional inversion;
   constant variable-expression substitution; hole order/containment/winding/
   collapse provenance; one-byte underreported path/file budgets; country code,
   manual vertex, projection/q/pin/serializer/runtime dependency changes.
5. Rebuild the full representative matrix twice under unchanged oracle and
   budgets. Regenerate recipe, manifest, receipt and contact sheet and obtain
   explicit owner approval of the new sheet digest. Machine failure cannot be
   waived.

### T0B — Full catalog and runtime integration

Start only from T0A machine and owner receipts. T0B may stage but not edit T0A
paths. Generate and validate all 996 ordinary outputs:

- 249 entities × two profiles × card/hero, no gaps or duplicates;
- natural and tiny/huge/portrait/landscape/square contain layouts;
- deterministic lineage, mandatory unions, omissions and output;
- topology, winding, visibility, fidelity and `2200/7500`, `<2500/<8000`
  budgets;
- warm generation at most 120 seconds on YMBPM3;
- normal Go generation performs no Node, Mapshaper, network or external I/O;
- regenerated full-catalog sheet and any bounded group-anchor override receive
  owner approval.

Run clean `GOWORK=off go test ./... -count=1`, public nil/source behavior,
typed source-budget failure and frozen serializer/dependency checks. Only then
create one atomic T0 product commit.

### T1 — Publication

From clean T0, perform two byte-identical maintainer rebuilds and publish the
all-283 ladder, manifest and generated initializer. Bind recipe,
`lineage-features/v2`, variable expression, numeric component provenance,
omissions, oracle, thresholds and tool. Runtime stays embedded, pure Go and
offline.

### T2 — Evidence and P2 closure

Add D3 projection/layout parity, lineage/component/hole mutations, topology
teeth, overrides, markers, budgets, natural/contain and visual evidence. Run
focused gates, full Go tests, offline `make check` and a fresh independent
audit.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0A | `component-variable-lineage-candidate-source` | `internal/geometry/cmd/lodbuild/**`; `internal/geometry/lod/tool/build.mjs`; `internal/geometry/lod/v2.recipe.json`; `internal/geometry/silhouette.go`; `internal/geometry/silhouette_test.go`; `internal/geometry/testdata/silhouette-oracle/**` |
| T0B | `geometry-foundation-integrator` | `go.mod`; `go.sum`; `internal/geometry/adapter.go`; `adapter_test.go`; `diagnostic.go`; `fit.go`; `fit_test.go`; `gridphase.go`; `gridphase_test.go`; `lod.go`; `lod_test.go`; `lod_catalog_test.go`; `marker.go`; `marker_test.go`; `math.go`; `model.go`; `model_test.go`; `normalize.go`; `normalize_test.go`; `override.go`; `override_test.go`; `overrides/**`; `path.go`; `path_test.go`; `pipeline.go`; `pipeline_test.go`; `preset.go`; `preset_test.go`; `presets/**`; `projection.go`; `projection_test.go`; `retain.go`; `retain_test.go`; `serialize.go`; `serialize_test.go`; `simplify.go`; `simplify_test.go`; `soften.go`; `soften_test.go`; `internal/geometry/lod/silhouette-oracle.v1.json`; `internal/geometry/lod/v1.recipe.json`; `internal/geometry/lod/tool/package.json`; `internal/geometry/lod/tool/package-lock.json`; `internal/geometry/testdata/lod-spike/**` |
| T1 | `silhouette-ladder-publisher` | `internal/geometry/lod/v2.ladder.json`; `internal/geometry/lod/v2.manifest.json`; `internal/geometry/lod_generated.go`; `internal/geometry/lod_build_test.go` |
| T2 | `geometry-verifier` | `Makefile`; `internal/geometry/approval_test.go`; `integration_test.go`; `mutations_test.go`; `oracle_test.go`; `internal/geometry/testdata/approval/**`; `testdata/d3/**`; `testdata/mutations/**` |

Leases do not overlap. Oracle, v1 recipe and serializer stay byte-frozen.
No implementer writes `data/**`, `internal/catalog/**`, public `cmd/**`,
lifecycle documents or another lease. Never stage `node_modules`, temporary
artifacts or unrelated user work.

## Dependencies and execution waves

Critical path: `W0/T0A synthetic → SH → CA/ID/KI → representative machine
gate → owner approval → W1/T0B → W2/T1 → W3/T2`. Concurrency is one. A wave
starts only after predecessor evidence is accepted.

## Validation plan

- `geometry-component-variable-resolution-and-sh-brake`;
- `geometry-lineage-holes-and-mutations`;
- `geometry-protected-visibility-and-scale-ladder-brake`;
- `geometry-ladder-integrated-budget-visual-performance`;
- `geometry-t0-clean-checkout`;
- `geometry-ladder-rebuild-check`;
- `geometry-ladder-public-wiring`;
- `geometry-d3-parity`, `geometry-overrides`, `geometry-markers`;
- `GOWORK=off go test ./... -count=1` and offline `make check`;
- independent terminal audit against PLAN-012 and both owner receipts.

## Deviation and amendment policy

`p2-component-variable-resolution-v1` permits only component feature encoding,
pinned Mapshaper variable resolution with the exact ratio above, deterministic
lineage/hole/mandatory-union validation, and resulting identities.

Any threshold, metric, rasterizer, visibility rule, band boundary, resolution
ladder/order, Mapshaper pin/flags/algorithm/weighting, projection, q, serializer,
budget, runtime dependency or selection input requires another plan. Country
branches, manual vertices, custom simplification, MakeValid, union/buffer
repair, CGO and post-hoc optimization remain forbidden.

## Commit worktree and integration policy

T0A creates no product commit. T0B creates one atomic T0 commit only after both
owner gates and clean-checkout validation. T1 and T2 each create one later
atomic commit after full gates. Stage only exact leased paths. Lifecycle
artifacts remain separate coordinator commits.

## Rollback and recovery

On the first T0A brake failure, restore the exact pre-run scratch inventory and
retain governed evidence naming geometry, band, resolution, lineage, metrics
and bytes. Never relax frozen policy. On T0B failure return to the accepted
representative checkpoint. T1 publishes from a temporary directory only after
two identical rebuilds. Revert completed commits in T2 → T1 → T0 order.

## Completion and handoff

P2 closes with deterministic component-aware all-283 artifacts, 996 ordinary
outputs, both owner-approved sheets, exact protection/omission provenance,
natural/contain layouts, hard budgets, at most 120-second offline generation,
clean T0/T1 checkouts and independent audit.

P3 receives selected band, candidate/oracle/recipe versions, per-component
lineage/protection/resolution provenance, removals, metrics, path bytes,
transform and optional markers. No presentation contract is added.

Estimated remaining effort: 8–18 agent-hours, likely 12, medium confidence.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
