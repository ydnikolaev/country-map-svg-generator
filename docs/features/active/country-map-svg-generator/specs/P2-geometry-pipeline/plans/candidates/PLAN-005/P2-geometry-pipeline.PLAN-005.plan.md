---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "implementation-plan"
id: "PLAN-005"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["P2", "PLAN-004", "RUN-004-RESULT", "DEC-003", "DEC-004", "AM-001", "AM-002", "CTR-001"]
---
# PLAN-005 — Narrow one-percent visual calibration successor

## Accepted inputs and baseline

This successor implements P2 at baseline commit
`89b62ac3621e410e99ebc9d3e3e3a08b8b2537aa`, tree
`6109cde3d1147ef3f89768bcef4e91377f43f6c4`. It is bound to:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- accepted PLAN-004 revision
  `48791fbd35f22b030817852734410e108f86abcfb52c6bd456a9d645633efb6d`;
- interrupted RUN-004
  `4297343e6af06c66e7a1de374c0ac3e518a5179e244eef2faf80d4f95bca5d1e`;
- PLAN-004 invalidation receipt
  `39fdbe88fc23698d23504573b2a4f6be823248536249ca329c17a919228cf409`;
- DEC-003 `f1002bdfbcf6947d11faa339161445614c2ae93ee340c2f97fa90a4036e7261d`,
  DEC-004 `3efec49a55fe7289ea528233280671e43db5f114c24586c1780dba98ba2667e5`,
  AM-001 `ab284fd6e2a79c4ef7a8c45ab05fe807d217ce7a8a178be354a44d26ef4ac80a`
  and AM-002 `ec2203d7132e467c7b1519d765c4790d583c29691311aad208a2ce6092523487`;
- P1 corpus identity
  `sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`.

PLAN-004 is invalidated, not edited. Its data-driven policy scratch is re-owned
and reverified here. RUN-004 exhausted ratios `0.0075`, `0.009` and `0.010`.
At the prior maximum, AQ card's first over-limit witness was
`1.2804229983137465px` against `1.28px`; AQ hero's was
`6.498361937568921px` against `6.4px`. Both source fallbacks then failed the
unchanged `0.01px` topology guard. No path bytes, broad matrix or visual sheet
exist, so none may be cited as acceptance evidence.

## Requirement and acceptance coverage

| Obligation | Tasks | Evidence |
| --- | --- | --- |
| US-1, AC-1 | T0, T2 | hard budgets and owner-approved card/hero sheet |
| US-2, REQ-1, REQ-2, REQ-8, AC-2 | T0, T2 | natural `tight`, arbitrary `contain`, scale-only quality goldens |
| US-3, REQ-6, AC-3 | T0, T2 | bounded data override and restoration mutations |
| US-4, REQ-3, REQ-4, AC-4 | T0–T2 | topology, deviation, softening, artifact and serializer teeth |
| REQ-5 | T0–T2 | omission, restoration, removal and fallback provenance |
| REQ-7 | T0, T2 | common-transform marker tests |
| INV-1 | T0 | presentation-free output |
| VAL-1 | T2 | pinned D3 parity |
| VAL-2 | T0, T2 | topology/protected/softening mutation suite |
| VAL-3 | T0, T2 | all-scale visual and budget matrix |
| VAL-4, VAL-5 | T2 | override and marker gates |

## Technical approach

All PLAN-004 production boundaries remain: full P1 geometry chooses the centered
equal-area projector, fit, natural viewBox, effective fitted scale, retention,
protected anchors and marker transform. The runtime is pure Go and offline.
Mapshaper 0.7.44 remains maintainer-only and pinned to the same integrity,
shasum and weighted-Visvalingam settings. Compact, standard and source remain the
only tiers; selection may inspect only effective fitted scale and table
availability, never ISO, profile, preset name or frame shape.

The successor changes one root assumption. Simplification error is a visual
silhouette budget, so its automatic default is scale-relative in fitted CSS-pixel
space rather than tightening toward an arbitrarily small absolute value as the
map grows. The only accepted formula family is:

`resolved = max(floor_px, relative_ratio * effective_fitted_long_side)`.

The floor is frozen at `0.80px`. T0 searches only `relative_ratio` in
`[0.010,0.012]`, beginning at `0.011` and moving monotonically by at most
`0.001`. This is a narrow successor around the measured approximately one-percent
AQ boundary, not permission for open-ended tolerance escalation. There is no
country, preset or frame branch and no hidden absolute ceiling. The exact ratio
freezes only after all-corpus structural/byte brakes and owner visual approval.
Values outside this range require another plan revision.

Every automatic quality coefficient—flatness, simplification floor/ratio,
softening, minimum area and quantization—is stored once in the versioned
`presets/v1.json` quality policy. `AutoQuality` consumes that policy for any
fitted size; card and hero only supply default rendered scale, padding and byte
budgets. Explicit `Input.Quality` and bounded per-country overrides remain the
general customization seams. Automatic resolved simplification is validated
against the versioned relative policy; explicit non-auto quality retains the
existing absolute `1.25px` safety bound and is never silently promoted to the
coarser automatic value. Invalid embedded policy is a deterministic package
configuration failure guarded by tests.

The LOD recipe remains at the RUN-004 endpoint `64/256` and thresholds
`240/700`. T0 may tune thresholds only through the already accepted compact
`[160,320]` and standard `[500,900]` ranges. Above the standard threshold the
source tier is attempted; a topology-invalid or over-budget fixed-grid source
fails typed and cannot weaken canonical precision. The all-scale matrix,
including `1024`, `2048`, `19x19` and `10000x12345`, decides whether the bounded
threshold range is viable. Failure interrupts; it does not add a tier, increase
recipe resolution, change `0.01px`, or reintroduce runtime RDP repair.

Candidate component/ring matching, required-component restoration, robust shared
source-vertex topology, symmetric deviation, fixed-grid canonicalization,
strict path-parser round trip, shortest valid absolute/relative serialization,
softening displacement/topology fallback and hard byte enforcement remain
unchanged. A candidate advances compact → standard → source on any guard failure.
Generic `MaxPathBytes` stays hard: card `2500`, hero `8000`;
`2200/7500` are advisory and zero is unlimited only without a preset.

## Tasks and completion conditions

1. **T0 — Calibrate and pass the integrated production brake.** Reconcile all
   RUN-004 scratch, retain the versioned quality policy, calibrate only the
   bounded relative ratio and rerun the exact production path over
   all 249 entities, both profiles, card/hero, long sides
   `90,128,160,240,500,640,700,1024,2048`, softening on/off and arbitrary
   `19x19`/`10000x12345` frames. Emit maxima, selected tiers, deviations,
   resolved tolerances, fallbacks, restorations, final path bytes, runtime,
   artifact hashes and two clean builder rebuilds. Produce the representative
   RU/CL/AR/AU/AQ/CA/ID plus island-state card/hero sheet. Done only when every
   machine brake passes, path maxima are ≤`2500/8000`, and the owner approves the
   sheet. Otherwise interrupt without a T0 product commit or T1 publication.
2. **T1 — Publish and embed the frozen LOD corpus.** Rebuild from the T0-frozen
   tool and literal recipe; publish compact/standard artifacts, manifest and Go
   embedding adapter. Prove exact P1, recipe, lockfile, Mapshaper and tier hashes,
   all-283 coverage, omission provenance, offline runtime and two byte-identical
   clean maintainer rebuilds.
3. **T2 — Complete evidence and the project ceiling.** Add full-corpus, D3,
   natural-ratio, arbitrary-frame, scale-only policy, protected restoration,
   budget, serializer, visual-approval, override and marker evidence. Run focused
   gates and offline `make check`; create the auditable P2 product commit and
   closure inputs.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0 | `lod-integrated-brake` | `go.mod`, `go.sum`, `internal/geometry/cmd/lodbuild/**`, `internal/geometry/lod/tool/**`, `internal/geometry/lod/v1.recipe.json`, `internal/geometry/testdata/lod-spike/**`, `internal/geometry/{adapter.go,adapter_test.go,diagnostic.go,fit.go,fit_test.go,lod.go,lod_test.go,marker.go,marker_test.go,math.go,model.go,model_test.go,normalize.go,normalize_test.go,override.go,override_test.go,path.go,path_test.go,pipeline.go,pipeline_test.go,preset.go,preset_test.go,projection.go,projection_test.go,retain.go,retain_test.go,serialize.go,serialize_test.go,simplify.go,simplify_test.go,soften.go,soften_test.go}`, `internal/geometry/overrides/**`, `internal/geometry/presets/**` |
| T1 | `lod-publisher` | `internal/geometry/lod/v1.compact.json`, `internal/geometry/lod/v1.standard.json`, `internal/geometry/lod/v1.manifest.json`, `internal/geometry/lod_generated.go`, `internal/geometry/lod_build_test.go` |
| T2 | `lod-verifier` | `Makefile`, `internal/geometry/{approval_test.go,integration_test.go,mutations_test.go,oracle_test.go}`, `internal/geometry/testdata/{approval,d3,mutations}/**` |

No implementer writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle docs
or another task's lease. The coordinator owns plans, runs, evidence and audit.

## Dependencies and execution waves

Critical path is W1/T0 → W2/T1 → W3/T2, concurrency one. T1 requires the exact
passing T0 machine receipt and explicit owner visual approval. T2 requires the
committed artifact reproducibility proof. Repository-wide gates are serial.

## Validation plan

- `geometry-quality-policy`: invalid/missing coefficients fail; explicit quality
  wins; equal fitted scale resolves equal defaults across preset/frame/ISO;
  reference sizes resolve inside the frozen formula and policy version.
- `geometry-lod-integrated-brake`: all-corpus matrix over final `d` bytes,
  topology/deviation/protection/restoration, visual digest, runtime, artifact
  sizes and two rebuilds.
- `geometry-svg-serializer`: strict command round trip, browser grammar,
  deterministic shortest encoding, negative-zero/sign/decimal mutation teeth.
- `geometry-lod-runtime`: source/artifact binding, scale-only selection,
  candidate advancement, cache parity/mutation and exact provenance.
- `geometry-topology-teeth`: shared source endpoint versus proper crossing,
  winding, holes, protected parts, fixed-grid failure and softened fallback.
- `geometry-d3-parity`, `geometry-overrides`, `geometry-markers` and offline
  `make check` retain their normative P2 meanings.

Teeth mutate policy/recipe/tool/tier/P1 hashes; introduce preset/ISO/frame
selection branches; exceed the tolerance formula; corrupt topology; lose a
protected part; accept an invalid `0.01px` source; drift serialized bytes; or
silently exceed a hard byte maximum.

## Deviation and amendment policy

`p2-lod-scale-relative-how-v2` permits helper/file splits inside the exact lease
and calibration only inside the declared global formula and threshold ranges.
A new tier, dependency, Mapshaper setting, recipe resolution, threshold range,
formula family, quantization, hard budget, selection input, serializer grammar
or runtime external process requires plan revision. Natural/contain semantics,
P1 authority, protected policy, presentation-free CTR-002 or single-binary
runtime changes require an amendment. Never recover with country branches,
manual geography, weakened topology or silent budget overflow.

## Commit worktree and integration policy

Each task lands as one conventional atomic product commit after its focused gate.
T0 commits only after the complete machine brake and owner visual approval.
Stage exact leased paths only; lifecycle evidence is never part of product
commits. RUN-004 scratch has no authority until reconciled and reverified here.

## Rollback and recovery

The previous versioned quality policy and LOD artifact set remain selectable.
A bad compact candidate advances to standard; a bad standard advances to source;
an invalid fixed-grid source fails typed. T1 stages all artifacts and publishes
only after complete validation. Revert T2 → T1 → T0 and select the prior
policy/corpus pair. Any exhausted search range returns to plan review.

## Completion and handoff

P2 closes only after owner-approved visuals, deterministic all-283 artifacts,
both-profile all-249 coverage, exact restoration/provenance, natural and arbitrary
layouts, canonical paths within hard budgets, offline `make check`, an exact run
result and fresh independent audit.

P3 receives CTR-002 with natural `tight`, distortion-free `contain`,
versioned data-driven defaults, explicit quality/override seams, canonical
presentation-free commands, stable diagnostics, metrics and optional transformed
capital markers. `detail` remains deferred.

Estimated remaining effort is 6–14 agent-hours, likely 10, medium confidence.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
