---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "accepted-plan-revision"
id: "PLAN-011"
epic: "country-map-svg-generator"
spec: "P2"
status: accepted
profiles: []
concerns: []
inputs: ["P2"]
---
# PLAN-011 — Lineage-keyed component candidate source

## Accepted inputs and baseline

Implement from commit
`b09883f5417494df313b56fec084511fba733fad`, tree
`89de323f9f14a47d9940b36e0cc61e114d7efd3d`.

Bind:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- invalidated PLAN-010 plan/manifest
  `a81da84fc3881068ad0520a3746f3276c18df87a48f10e83d5463fb03fb1db63` /
  `fdd0cdfa0c545410f102eb460fef80b57d3b03484a2fc0b1cb423975a7fff59d`;
- RUN-017 result
  `ff99d94cbcefd9b10943f3d81ea182d035862ab7d62a302b0a2a5666c6489f07`;
- PLAN-010 invalidation receipt
  `64070d25f6de43180a030e51b390e48e850af907a7c141dbe4e8f19d24d52534`;
- byte-frozen oracle
  `f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`;
- predecessor v2 recipe
  `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`.

PLAN-010 remains invalidated. Its product intent and completed scratch remain
binding except for its incomplete multipart candidate-source assumption. The
current 51-file `internal/geometry` scratch inventory, excluding
`node_modules`, is
`d4d16e5af1b77a02dd4163078558e0a80ebf32c5599ea3132f53e23ab617bb61`
(SHA-256 of `shasum -a 256` lines for C-locale sorted paths). Reconcile it by
exact pre/post inventory; never recreate or broadly delete it. `detail` remains
deferred backlog.

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

Every normative obligation has one manifest owner. Later tasks may consume
earlier evidence without claiming duplicate coverage.

## Technical approach

Preserve centered LAEA projection, natural `tight`, uniform arbitrary
`contain`, q `0.01`, fine-to-coarse candidate order, pure-Go offline runtime,
source/custom semantics, serializer bytes, scale boundaries `240/700`, path
caps `2200/7500`, complete maxima `2500/8000`, contribution thresholds
`110/900`, and global IoU/recall `.40/.42`.

Replace only candidate-source encoding and selection preparation:

1. After exact Go projection, split each source MultiPolygon into temporary
   Polygon features. Each carries a versioned lineage tuple: geometry ID,
   source polygon index, source ring count and canonical projected source
   polygon digest.
2. Submit all component features for one geometry in one FeatureCollection to
   one unchanged pinned Mapshaper 0.7.44 weighted-Visvalingam planar command.
   This keeps `resolution` relative to the country extent while making
   `keep-shapes` operate on explicit component features. Change batching, not
   Mapshaper version, flags, algorithm, weighting or resolution semantics.
3. Require exactly one output feature for every lineage key before policy
   pruning. Reject missing, duplicate, unknown, mutated or cross-geometry keys;
   output order has no authority.
4. Validate each output as one non-empty Polygon. Preserve exterior and
   surviving holes. Match holes to source rings deterministically; record a
   collapsed hole as a ring omission and reject invented, escaped, inverted or
   cross-component holes.
5. For each geometry and band, resolve all entity/profile consumers
   (`identical_to` included) and form one generic mandatory union of source
   indices: every component that is visible, dominant or anchored for any
   consumer. Reassemble mandatory output components in source order. A
   contribution-ranked identity component below the band threshold remains
   optional and, when omitted, must carry `protected_subscale`; other optional
   components require every consumer to classify them subscale/non-identity and
   carry `subscale`. Runtime entity/profile/preset/frame inputs never change
   candidate identity.
6. Candidate selection remains first fine-to-coarse resolution passing every
   consumer under the requested band, topology and exact byte policy. Hash the
   candidate-source version, ordered lineage table, mandatory-union proof,
   per-component output digests and omissions into candidate/recipe/publication
   identities.

The oracle file remains byte-identical. The new v2 recipe may differ from
`da3ff9d3…` only by its versioned candidate-source declaration and derived
identity; resolutions, order, projection, Mapshaper policy and visibility
semantics remain exact.

This explicitly avoids the two invalid extremes: anonymous multipart
Mapshaper loss, and keeping every four-point subscale island until archipelagos
cannot meet budgets.

## Tasks and completion conditions

### T0A — Component-aware source and representative checkpoint

Re-own the existing candidate-builder scratch. Implement
`lineage-features/v1`, deterministic mandatory-union selection/reassembly and
typed lineage failures.

First run a hard micro-brake over SH plus high-part/worst-budget witnesses
CA, ID and KI:

- SH at every frozen resolution must preserve component 0 lineage and nonzero
  contribution for both consumers;
- mandatory-union construction must be identical regardless of entity/profile
  iteration or Mapshaper output order;
- optional subscale islands must remain omittable with truthful provenance;
- at least one frozen candidate per requested band must pass existing budgets.

Stop immediately if component-aware source cannot satisfy that brake without a
policy change.

Add mutation teeth for missing, duplicate, unknown and cross-geometry lineage;
source index/digest mutation; output reordering; exterior/hole reordering;
missing, promoted, inverted or escaped holes; per-consumer candidate branching;
and false Mapshaper-loss/subscale provenance.

Then rerun:

- RUN-017 twice with identical semantic evidence;
- Indonesia 120/121 boundary and every existing component/topology/protection/q
  regression;
- full representative matrix with unchanged thresholds and budgets;
- two deterministic representative builds.

The candidate-source version changes recipe/candidate identities, so regenerate
the representative manifest, receipt and contact sheet and obtain explicit
owner approval. T0A creates no product commit. Machine failure cannot be waived
by the owner.

### T0B — Full catalog and runtime integration

Start only from accepted T0A machine and owner receipts. Generate all 996
ordinary outputs and prove:

- 249 entities × two profiles × card/hero, no missing or duplicates;
- natural plus tiny/huge/portrait/landscape/square contain layouts;
- deterministic lineage, mandatory unions, omissions, paths and repeats;
- topology, winding, visibility and IoU/recall;
- path maxima at most `2200/7500` and complete estimates below `2500/8000`;
- warm 996-output generation at most 120 seconds on YMBPM3;
- normal Go generation performs no Node, Mapshaper, network or external I/O;
- owner approval of the regenerated full-catalog sheet and any bounded reviewed
  group-anchor override.

Integrate requested-band, rotation-safe protected visibility into
`GenerateWithLOD`; explicit source/custom and caller hard-byte semantics remain
unchanged. Run clean `GOWORK=off go test ./... -count=1`, public nil/source
behavior, typed source-budget failure and frozen serializer/dependency checks.
Only then stage the exact T0A/T0B union and create one atomic T0 product commit.

### T1 — Publication

From clean T0, perform two byte-identical maintainer rebuilds and publish the
all-283 v2 ladder, manifest and generated initializer. Bind P1, projection,
oracle, recipe, `lineage-features/v1`, mandatory unions, ordered component
provenance, thresholds, tool, candidate order and omission table. Public
runtime remains embedded, pure Go, fitted-scale-only and offline.

### T2 — Evidence and P2 closure

Add D3 projection/layout parity, component/hole mutations, topology teeth,
overrides, markers, budgets, natural/contain and visual evidence. Run focused
gates, `GOWORK=off go test ./... -count=1`, offline `make check` and a fresh
independent audit.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0A | `component-lineage-candidate-source` | `internal/geometry/cmd/lodbuild/**`; `internal/geometry/lod/tool/build.mjs`; `internal/geometry/lod/v2.recipe.json`; `internal/geometry/silhouette.go`; `internal/geometry/silhouette_test.go`; `internal/geometry/testdata/silhouette-oracle/**` |
| T0B | `geometry-foundation-integrator` | `go.mod`; `go.sum`; `internal/geometry/adapter.go`; `adapter_test.go`; `diagnostic.go`; `fit.go`; `fit_test.go`; `gridphase.go`; `gridphase_test.go`; `lod.go`; `lod_test.go`; `lod_catalog_test.go`; `marker.go`; `marker_test.go`; `math.go`; `model.go`; `model_test.go`; `normalize.go`; `normalize_test.go`; `override.go`; `override_test.go`; `overrides/**`; `path.go`; `path_test.go`; `pipeline.go`; `pipeline_test.go`; `preset.go`; `preset_test.go`; `presets/**`; `projection.go`; `projection_test.go`; `retain.go`; `retain_test.go`; `serialize.go`; `serialize_test.go`; `simplify.go`; `simplify_test.go`; `soften.go`; `soften_test.go`; `internal/geometry/lod/silhouette-oracle.v1.json`; `internal/geometry/lod/v1.recipe.json`; `internal/geometry/lod/tool/package.json`; `internal/geometry/lod/tool/package-lock.json`; `internal/geometry/testdata/lod-spike/**` |
| T1 | `silhouette-ladder-publisher` | `internal/geometry/lod/v2.ladder.json`; `internal/geometry/lod/v2.manifest.json`; `internal/geometry/lod_generated.go`; `internal/geometry/lod_build_test.go` |
| T2 | `geometry-verifier` | `Makefile`; `internal/geometry/approval_test.go`; `integration_test.go`; `mutations_test.go`; `oracle_test.go`; `internal/geometry/testdata/approval/**`; `testdata/d3/**`; `testdata/mutations/**` |

Leases do not overlap. T0B may stage verified T0A bytes but cannot edit them.
The oracle data and v1 recipe stay byte-frozen in the T0B lease; any change
returns to T0A and owner review. Serializer files remain frozen at
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

No implementer writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle
documents or another lease. Never stage `node_modules`, temporary artifacts or
unrelated user work.

## Dependencies and execution waves

Critical path is `W0/T0A micro-brake → representative machine gate → owner
approval → W1/T0B → W2/T1 → W3/T2`; concurrency is one. A wave begins only
after its predecessor evidence is accepted.

## Validation plan

- `geometry-component-lineage-and-sh-brake`: lineage, mandatory union,
  component/hole mutations, SH/CA/ID/KI and Indonesia;
- `geometry-regression-protected-visibility-and-scale-ladder-brake`:
  representative policy and deterministic artifacts;
- `geometry-ladder-integrated-budget-visual-performance`: 996 outputs,
  layouts, budgets, runtime and full-catalog owner sheet;
- `geometry-t0-clean-checkout`;
- `geometry-ladder-rebuild-check`;
- `geometry-ladder-public-wiring`;
- `geometry-d3-parity`, `geometry-overrides`, `geometry-markers`;
- offline project ceiling.

Teeth must fail if components are anonymously dropped, output order controls
identity, optional parts are kept until budgets fail, mandatory parts are
pruned, a consumer changes candidate identity, lineage/provenance lies, or a
runtime external dependency appears.

## Deviation and amendment policy

`p2-component-lineage-candidate-source-v1` permits only temporary component
feature encoding, deterministic mandatory-union validation/reassembly and the
resulting candidate/hash/fixture identities.

Any threshold, metric, rasterizer, visibility rule, band boundary, resolution,
Mapshaper pin/flags/algorithm/weighting, candidate order, projection, q,
serializer, budget, runtime dependency or selection-input change requires a new
plan. Country branches, manual vertices, MakeValid, union/buffer repair, CGO and
post-hoc optimization remain forbidden.

## Commit worktree and integration policy

T0A creates no product commit. T0B creates one atomic T0 commit only after both
owner gates and clean-checkout validation. T1 and T2 each create one later
atomic commit after their full gates. Stage only exact leased paths. Lifecycle
artifacts are separate coordinator commits.

## Rollback and recovery

On T0A failure retain prior scratch and record exact
geometry/band/resolution/lineage evidence; never relax policy. On T0B failure
return to the accepted representative checkpoint. T1 publishes from a temporary
directory only after two identical rebuilds. Revert completed product commits
in T2 → T1 → T0 order.

## Completion and handoff

P2 closes with deterministic component-aware all-283 artifacts, 996 ordinary
outputs, both owner-approved sheets, exact protection/omission provenance,
natural/contain layouts, budgets, at most 120-second offline generation, clean
T0/T1 checkouts and independent audit.

P3 receives existing selected-band, candidate/oracle/recipe versions,
per-component lineage/protection provenance, removals, metrics, path bytes,
transform and optional markers. No presentation contract is added.

Estimated remaining effort: 7–16 agent-hours, likely 11, medium confidence.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
