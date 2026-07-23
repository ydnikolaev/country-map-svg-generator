---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "implementation-plan"
id: "PLAN-003"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["P2", "DEC-003", "DEC-004", "AM-001", "AM-002", "RUN-001-RESULT", "RUN-002-RESULT", "CTR-001"]
---
# PLAN-003 — Integrated pre-publication LOD pipeline

## Accepted inputs and baseline

This successor implements P2 at baseline commit
`c2ea11c56c0af83aa25a0cd58188864cb0129275`, tree
`82e476d45052078f87ccc4d4bda06ab0a689cdd7`. It is bound to:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- DEC-003 `f1002bdfbcf6947d11faa339161445614c2ae93ee340c2f97fa90a4036e7261d`
  and DEC-004
  `3efec49a55fe7289ea528233280671e43db5f114c24586c1780dba98ba2667e5`;
- AM-001 `ab284fd6e2a79c4ef7a8c45ab05fe807d217ce7a8a178be354a44d26ef4ac80a`
  and AM-002 `ec2203d7132e467c7b1519d765c4790d583c29691311aad208a2ce6092523487`;
- interrupted RUN-001
  `2d704d25d76d2646f5cfb8061623ed19bbe43ebb24af5ad60e711c6d9e92ac87`;
- interrupted RUN-002
  `d6956cb35cab56e0de9efa3e43ac617d255147d471b7260b8484ba589bceebb8`;
- PLAN-002 invalidation receipt
  `ded9b234d2c832ac0eb2522db5386841415358e0342f66d63a67c69732ed356d`;
- P1 corpus identity
  `sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`,
  manifest SHA-256
  `227374e94910f078623b3c643085d39fa094b4e1ae20a0ffc01cb2d4a19ee616`,
  and geometry SHA-256
  `f3d5d80ac8b9650cb75f68f46593c24476a979855875842bc3274d9d928e7744`.

PLAN-002 is invalidated, not edited. Its T0 builder and the earlier Go pipeline
remain uncommitted scratch. Every reused byte is re-owned and reverified here.
The diagnostic resolutions below are evidence about search bounds, not accepted
calibration values.

## Requirement and acceptance coverage

| Obligation | Tasks | Evidence |
| --- | --- | --- |
| US-1, AC-1 | T0, T2 | final-path budgets, approved card/hero sheet, full-corpus receipt |
| US-2, REQ-1, REQ-2, REQ-8, AC-2 | T0, T2 | natural `tight`, arbitrary `contain`, scale-only selection goldens |
| US-3, REQ-6, AC-3 | T0, T2 | generic protected/minimum-parts restoration and override mutations |
| US-4, REQ-3, REQ-4, AC-4 | T0–T2 | source/tool/artifact binding, topology/deviation teeth and rebuild proof |
| REQ-5 | T0–T2 | exact omission, restoration, removal and fallback provenance |
| REQ-7 | T0, T2 | unchanged common marker transform and anomaly gates |
| INV-1 | T0 | presentation-free geometry and canonical commands |
| VAL-1 | T2 | pinned D3 parity |
| VAL-2 | T0, T2 | LOD/topology/protected/softening mutation suite |
| VAL-3 | T0, T2 | multi-scale visual, arbitrary-layout and hard-budget gates |
| VAL-4 | T2 | bounded override gate |
| VAL-5 | T2 | marker gate |

## Technical approach

Preserve the accepted full-source normalization, centered Lambert azimuthal
equal-area projection, adaptive subdivision, natural-aspect fitting and common
marker transform. The full P1 geometry always determines the projector, natural
bounds, `tight` viewBox or `contain` transform, effective fitted long side,
retention decisions and protected-component binding before a LOD candidate is
read. Selected geographic LOD boundaries use that exact projector and transform.

T0 is an integrated pre-publication brake, not a builder-only spike. Its
production pipeline accepts a temporary injected LOD table in tests while normal
runtime falls back losslessly to P1 until T1 supplies the generated embedded
table. The same production functions perform scale-only tier selection,
component/ring matching, restoration, deviation checks, retention, command
serialization and budget fallback. No acceptance check depends on code scheduled
after T0.

The maintainer builder keeps Mapshaper 0.7.44 pinned by npm integrity
`sha512-3Cx+IABMXt1G28Y8J7oalW5P5VYyt1vHz5FO+KkV51EHnJioo9h9maO+u+4IyCPZ9Mh1hqchgomFmc68GtFwQQ==`
and shasum `e08fc40d50347698ea07d2f82ead42a98b1c62f0`. Each of the 283 geometry IDs
is processed independently as one whole MultiPolygon within one Node process.
Weighted Visvalingam (`weighting=0.7`), default intersection rollback,
`keep-shapes`, and `clean` produce `compact` and `standard`; exact P1 is the
`source` tier.

RUN-002 proved that the prior calibration range was not a valid all-corpus
bound: resolution 48 yielded AQ card 2,613 bytes, 32 yielded CA card 4,184,
128 yielded CA hero 10,761, and 96 yielded CA hero 8,250. A diagnostic resolution
12 let CA card advance but was never visually or structurally accepted.
Accordingly T0 searches only global literal values within compact `[8,64]` and
standard `[64,256]`, plus global scale thresholds compact `[160,320]` and
standard `[500,900]`. It begins at `24/96` and `240/700`. The final recipe freezes
one value per field only after every machine brake and the representative visual
sheet pass. Leaving these bounds, adding a tier, or inspecting ISO, profile,
preset or frame shape requires another plan revision.

Topology validation distinguishes a proper crossing from a shared source vertex.
Projection/subdivision carries stable source-vertex identity; non-adjacent
segments may meet only when the full-source ring has the same shared vertex.
Numeric endpoint equivalence is bounded by machine precision and canonical
quantization and cannot excuse a proper crossing, changed winding, escaped hole,
collapsed ring or new self-touch. AQ is the mandatory regression for the
approximately `1e-17` projection drift observed in RUN-002.

Retention is computed from full projected source: largest component, every
anchor-bearing component, stable minimum-parts selections and every component
above the resolved visibility threshold. Any required component/ring missing
from a candidate is restored from full source before validation. ID is the
mandatory generic protected-restoration regression; restoration never branches
on country identity.

Candidate boundaries are component/ring matched and densely compared in fitted
space. Require symmetric
`max(full→candidate, candidate→full) <= resolved simplification tolerance`.
Failure advances compact → standard → source and records exact provenance.

Canonical SVG path serialization is part of T0's byte brake. From fixed `0.01px`
commands, a deterministic generic serializer chooses the shorter valid encoding
from absolute and relative forms, elides repeated command letters where SVG
grammar permits, removes redundant separators and leading/trailing zeroes, and
never changes command meaning. A strict in-repo parser round-trips the emitted
`d` bytes back to the canonical commands; adversarial negative-zero, sign,
decimal and subpath fixtures guard minification. No SVGO or post-hoc optimizer
exists.

Generic `MaxPathBytes` remains hard: card defaults to 2,500 bytes, hero to 8,000;
2,200/7,500 are advisory. Zero is unlimited only for unpresetted requests.
Optional softening is topology checked and falls back to deterministic linear
commands when softening is invalid, over budget or exceeds the resolved profile
tolerance. Before serialization, densely compare every softened fitted-space
curve to its canonical linear source at no more than `0.05px` spacing and require
maximum symmetric displacement no greater than the resolved tolerance. Excess
displacement emits `softening_tolerance_fallback` and uses linear commands; a
failure to validate the linear fallback fails closed. If even the valid linear
source fallback exceeds a nonzero hard maximum, generation fails with an exact
diagnostic; it never weakens geometry silently.

Node, npm and Mapshaper are maintainer-only. Shipped Go binaries, normal
generation and `make check` neither execute nor require them. `detail` remains
deferred.

## Tasks and completion conditions

1. **T0 — Implement and accept the integrated LOD brake.** Finish the pinned
   builder and recipe together with the production projector/fit, injected LOD
   table seam, full-source retention/restoration, scale-only selection,
   topology/deviation guards, fixed-grid commands, minified serializer, generic
   budgets and softening fallback. Temporarily build all 283 IDs for each
   calibration candidate and exercise the exact production path across all 249
   entities, both profiles, card/hero, reference long sides
   `90,128,160,240,500,640,700,1024,2048`, softening on/off and arbitrary
   `19x19`/`10000x12345` frames. Emit exact maxima, fallbacks, restored source
   rings, artifact bytes, runtime, hashes, two clean rebuilds and a representative
   RU/CL/AR/AU/AQ/CA/ID plus island-state visual sheet. Done only when every
   structural/protection/deviation/determinism brake passes, final path maxima are
   ≤2,500/8,000 and the owner accepts the sheet. Otherwise stop without a product
   commit or T1.
2. **T1 — Publish and embed the frozen LOD corpus.** Rebuild once from the
   T0-frozen tool and literal recipe, publish compact/standard artifacts plus one
   manifest and generated Go embedding adapter, and validate exact P1, recipe,
   package-lock, Mapshaper, tier hashes, all-283 coverage and component/ring
   omission provenance without Node. Two clean maintainer rebuilds must reproduce
   the committed bytes. Ordinary Go tests use only the embedded artifacts.
3. **T2 — Allocate complete evidence and close the product gate.** Add the
   full-corpus, D3, natural-ratio, arbitrary-frame, scale-only selection,
   protected-restoration, budget, serializer, visual-approval, override and marker
   evidence. Run focused gates and offline `make check`; produce an auditable
   product commit and exact P2 verification inputs.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0 | `lod-integrated-brake` | `go.mod`, `go.sum`, `internal/geometry/cmd/lodbuild/**`, `internal/geometry/lod/tool/**`, `internal/geometry/lod/v1.recipe.json`, `internal/geometry/testdata/lod-spike/**`, `internal/geometry/{adapter.go,adapter_test.go,diagnostic.go,fit.go,fit_test.go,lod.go,lod_test.go,marker.go,marker_test.go,math.go,model.go,model_test.go,normalize.go,normalize_test.go,override.go,override_test.go,path.go,path_test.go,pipeline.go,pipeline_test.go,preset.go,preset_test.go,projection.go,projection_test.go,retain.go,retain_test.go,serialize.go,serialize_test.go,simplify.go,simplify_test.go,soften.go,soften_test.go}`, `internal/geometry/overrides/**`, `internal/geometry/presets/**` |
| T1 | `lod-publisher` | `internal/geometry/lod/v1.compact.json`, `internal/geometry/lod/v1.standard.json`, `internal/geometry/lod/v1.manifest.json`, `internal/geometry/lod_generated.go`, `internal/geometry/lod_build_test.go` |
| T2 | `lod-verifier` | `Makefile`, `internal/geometry/{approval_test.go,integration_test.go,mutations_test.go,oracle_test.go}`, `internal/geometry/testdata/{approval,d3,mutations}/**` |

No implementer writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle docs
or another task's lease. The coordinator owns plans, runs, evidence and audit.
The leases are disjoint; generated artifact and adapter names are reserved to T1.

## Dependencies and execution waves

Critical path is T0 → T1 → T2. T0 owns every production behavior required by its
acceptance brake and tests temporary candidates through injection, so it has no
dependency on publication. T1 requires the exact passing T0 receipt plus owner
visual approval. T2 requires committed artifact reproducibility. Repository-wide
gates are serial.

## Validation plan

- `geometry-lod-integrated-brake`: T0 all-corpus matrix, exact final `d` bytes,
  topology/deviation/protection/restoration, visual digest, artifact size/runtime
  and two clean rebuilds.
- `geometry-svg-serializer`: canonical-command round trip, browser-valid grammar,
  absolute/relative shortest-choice determinism and mutation teeth.
- `geometry-lod-rebuild-check`: two temporary maintainer rebuilds followed by
  byte comparison with T1 artifacts; never part of ordinary `make check`.
- `geometry-lod-runtime`: artifact/source binding, scale-only selection,
  restoration, fallback and determinism.
- `geometry-topology-teeth`: shared-source endpoint versus proper-crossing,
  winding, holes, protected parts and softened fallback mutations.
- `geometry-d3-parity`, `geometry-overrides` and `geometry-markers` retain their
  normative P2 meanings.
- `make check` validates committed artifacts and every Go test offline; it fails
  if normal validation invokes Node, npm, Mapshaper, network, SVGO or an
  uncommitted source.

Teeth mutate P1/recipe/tool/tier hashes; remove a geometry ID; inspect
ISO/profile/preset/frame shape in selection; corrupt a ring; collapse or fake a
shared endpoint; lose a protected/minimum part; exceed deviation; drift rebuild
bytes; misparse minified commands; or accept a softened over-budget,
over-tolerance or intersecting path.

## Deviation and amendment policy

`p2-lod-integrated-how-v1` permits helper/file splits inside one leased task,
equivalent encoding/hash details and calibration inside the declared global
ranges. A new tier, dependency, Mapshaper version/algorithm, search range, hard
budget, selection input, serializer grammar family, source root or runtime
external process requires plan revision. Natural/contain semantics, protected
policy, P1 authority, presentation-free CTR-002 or single-binary/offline runtime
changes require an amendment or successor decision. Never recover with country
branches, manual geography, weakened topology or silent budget overflow.

## Commit worktree and integration policy

Each task lands as one conventional atomic product commit after its focused gate.
T0 commits only after its complete machine brake and owner visual acceptance.
Stage exact leased paths only. RUN-001/RUN-002 scratch has no authorship authority
until reconciled against this plan. Lifecycle evidence is never included in a
product commit.

## Rollback and recovery

Source is always the semantic fallback. A bad compact candidate advances to
standard; a bad standard candidate advances to exact P1 with provenance. T1
stages all generated outputs and publishes only after full validation. P1,
recipe, tool or threshold changes create a new LOD version; accepted v1 bytes are
never regenerated in place. Revert T2 → T1 → T0 and select the prior
algorithm/LOD pair. Failure returns to plan review, never tolerance escalation or
manual redrawing.

## Completion and handoff

P2 closes only after accepted T0 visuals, deterministic all-283 artifacts,
both-profile all-249 coverage, exact restoration/provenance, natural and arbitrary
layouts, minified canonical paths within hard budgets, offline `make check`, an
exact run result and fresh independent audit.

P3 receives additive CTR-002 LOD provenance plus the general layout/quality API:
natural `tight`, distortion-free arbitrary `contain`, data-driven presets,
canonical presentation-free commands, stable diagnostics, metrics, removals and
optional transformed capital markers. P3 serializes/stylizes the SVG; it does not
optimize bloated geometry after the fact.

Estimated effort is 12–20 agent-hours, likely 16, medium confidence. The integrated
T0 brake now bounds calibration, restoration, serializer and visual uncertainty
before artifact publication.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
