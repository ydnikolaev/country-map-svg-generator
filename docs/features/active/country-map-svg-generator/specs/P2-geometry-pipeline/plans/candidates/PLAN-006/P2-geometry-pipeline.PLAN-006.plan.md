---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "implementation-plan"
id: "PLAN-006"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["P2", "PLAN-005", "RUN-004-RESULT", "RUN-005-RESULT", "RESULT-012", "DEC-003", "DEC-004", "AM-001", "AM-002", "CTR-001"]
---
# PLAN-006 — Projected-coordinate LOD and fixed-grid source finalization

## Accepted inputs and baseline

This successor implements P2 at baseline commit
`7ab35331e3e30291769fd178d7b2f0a9aa91c3df`, tree
`50b5d30f4ee904d6c53c92a8610a988f2526c5d3`. It is bound to:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- invalidated PLAN-005 revision
  `4c1137e5413869375de6a9e6b98fc6a350873da6e02a1f13a032409b80b2a847`;
- interrupted RUN-004
  `4297343e6af06c66e7a1de374c0ac3e518a5179e244eef2faf80d4f95bca5d1e`;
- interrupted RUN-005
  `8bb67d3580618337cc73323ccfeeaa6c59f4037026f36cfbecc5dbd7cc38582f`;
- RESULT-012
  `bc14eb9e4b209501a25e6e09a2ad622311d9b53d5e8dacfb565812be21f654a7`;
- PLAN-005 invalidation receipt
  `81671b152b43cfa80634840de6d80587927eff415819b4b45c0a09aadf01113d`;
- DEC-003 `f1002bdfbcf6947d11faa339161445614c2ae93ee340c2f97fa90a4036e7261d`,
  DEC-004 `3efec49a55fe7289ea528233280671e43db5f114c24586c1780dba98ba2667e5`,
  AM-001 `ab284fd6e2a79c4ef7a8c45ab05fe807d217ce7a8a178be354a44d26ef4ac80a`
  and AM-002 `ec2203d7132e467c7b1519d765c4790d583c29691311aad208a2ce6092523487`;
- P1 corpus identity
  `sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`.

PLAN-005 remains invalidated. Its product scratch is re-owned only after
reconciliation against this plan.

RUN-005 proved that the load-bearing defect is a representation mismatch, not
insufficient visual tolerance. At the accepted ceiling `0.012`, AQ card compact
deviation was `1.583086150019346px` against `1.536px`; AQ hero standard was
`8.163153181359744px` against `7.68px`. Compact and standard are currently
simplified in geographic longitude/latitude space, while their guard is measured
after a per-geometry centered LAEA projection and pixel fit.

The exact retained P1 fallback is valid before fixed-grid conversion, but
pointwise `0.01px` rounding can create proper crossings. The committed scratch
tests establish this independently: AQ has 23,994 valid fitted source points,
while its canonical probe observes three crossings in the first 500 segments
and seven in the first 1,000. RU card supplies a second non-polar witness.

## Requirement and acceptance coverage

| Obligation | Tasks | Evidence |
| --- | --- | --- |
| US-1, AC-1 | T0, T2 | final budgets and owner-approved card/hero sheet |
| US-2, REQ-1, REQ-2, REQ-8, AC-2 | T0, T2 | natural `tight`, arbitrary `contain`, scale-only and deterministic goldens |
| US-3, REQ-6, AC-3 | T0, T2 | bounded overrides and generic restoration |
| US-4, REQ-3, REQ-4, AC-4 | T0–T2 | final-representation topology, deviation and mutation teeth |
| REQ-5 | T0–T2 | omission, restoration, removal and fallback provenance |
| REQ-7 | T0, T2 | unchanged full-source marker transform |
| INV-1 | T0 | presentation-free result and canonical commands |
| VAL-1 | T2 | pinned D3 parity |
| VAL-2 | T0, T2 | topology, protection, source-finalization and softening mutations |
| VAL-3 | T0, T2 | full scale/layout/visual/budget matrix |
| VAL-4, VAL-5 | T2 | override and marker gates |

## Technical approach

Full P1 geometry remains authoritative. It determines the centered equal-area
projector, natural bounds, `tight` or `contain` fit, effective fitted long side,
retained components, protected anchors and marker transform.

Change compact and standard artifact coordinates. During the maintainer build,
normalize each of the 283 P1 geometries, derive its centered projector with zero
planar rotation, adaptively subdivide and project it through the exact Go LAEA
implementation, then pass those Cartesian coordinates to pinned Mapshaper 0.7.44
with explicit `planar`, weighted Visvalingam `0.7`, intersection rollback,
`keep-shapes`, `clean`, and the unchanged resolutions `64/256`.

This aligns the builder with Mapshaper's documented semantics: `resolution=`
targets the intended display extent, longitude/latitude data is simplified on a
sphere, and projected data is simplified in a Cartesian plane
([Mapshaper simplification guide](https://mapshaper.org/docs/guides/simplification.html)).
The runtime acceptance guard is already Cartesian fitted-pixel space.

Each artifact record stores and hashes its coordinate-space contract:
`centered_laea`, center longitude/latitude, zero build rotation, y-axis
convention, projection flatness `0.10`, coordinate precision, corpus identity,
recipe identity and geometry ID. Runtime recomputes the center from exact P1 and
rejects mismatched metadata. A per-country rotation override applies the same
planar rotation to full projected source and projected LOD candidate. Full
source still owns fit; the projected candidate receives that exact uniform scale
and translation.

Runtime selection remains based only on effective fitted scale and artifact
availability. Compact, standard and source remain the only tiers at thresholds
`240/700`. ISO, profile, preset and frame shape remain forbidden selection
inputs.

All tiers pass through one final-representation guard:

1. restore required full-source components/rings and validate raw topology;
2. for compact/standard, reject raw fitted-space symmetric deviation above the
   resolved tolerance;
3. canonicalize at exactly `q=0.01`;
4. validate final topology, winding, holes, non-empty identity and protected
   retention;
5. measure symmetric deviation from full fitted source to final canonical
   representation and back; record only this final maximum as authoritative;
6. serialize and enforce the caller's generic hard byte limit.

Selection-only diagnostics execute steps 1–5. They may skip command
serialization but may not report a tier that production would reject at
`q=0.01`.

Source fallback remains provenance tier `source`, derived from exact retained
P1, but it no longer pointwise-rounds the unsimplified ring. It materializes a
deterministic fitted-space reduction through the existing validated Go
simplification path. Make the tolerance search q-aware: accept a candidate only
when `q=0.01` canonicalization, final topology, protected retention and final
symmetric deviation pass. Begin at the coarsest permitted simplification after
reserving the derived maximum grid displacement `q/sqrt(2)`, then move
monotonically toward finer representations. Never exceed the resolved visual
tolerance; typed exhaustion is terminal.

Do not use `MakeValid`, polygon union, ad-hoc vertex movement, relaxed
shared-endpoint rules, country-specific repair, smaller q, larger visual
tolerance or post-hoc SVG optimization. These operations could hide changed
identity instead of proving it.

The versioned quality policy keeps floor `0.80px` and
`relative_ratio=0.012` as the hard experiment ceiling. T0 cannot increase it.
Owner review may reduce it only with a full machine rerun. Explicit quality
retains its existing `1.25px` bound.

Preserve `q=0.01`, `64/256`, `240/700`, compact/standard/source provenance,
card `2500` and hero `8000` hard budgets, `2200/7500` advisory thresholds,
natural-aspect `tight`, distortion-free `contain`, and deterministic offline Go
runtime.

## Tasks and completion conditions

1. **T0 — Repair representation and pass the integrated production brake.**

   First run a small representation experiment with serializer bytes frozen:

   - rebuild `64/256` artifacts in centered unrotated LAEA planar coordinates;
   - run AQ/un card and hero through compact/standard fallback and final
     `q=0.01` validation;
   - force source for AQ card/hero and RU card through the q-aware reducer;
   - run each case first as an equivalently resolved unpresetted input with
     `MaxPathBytes=0`, then with normal `2500/8000`;
   - emit requested/selected tier, coordinate-space identity, raw/final
     deviations, resolved/source tolerance, q, restorations, fallbacks, output
     points, path bytes, parser round-trip and runtime.

   Outcomes are exclusive. If an unlimited case fails topology, protection or
   final deviation, representation repair failed: stop without serializer
   changes or a product commit. If unlimited passes but a budgeted case fails,
   representation is proven and the next blocker is exact serializer/byte
   evidence: stop with measured gaps and do not edit the serializer. Continue
   only if AQ card/hero and forced-source witnesses produce final round-tripping
   paths within `2500/8000`.

   Then run the production path over all 249 entities, both profiles, card/hero,
   long sides `90,128,160,240,500,640,700,1024,2048`, softening on/off, and
   arbitrary `19x19` and `10000x12345` contain frames. Emit maxima, selected
   tiers, final deviations, source reductions, fallbacks, restorations, path
   bytes, runtime, artifact hashes and two clean builder rebuilds. Produce the
   RU/CL/AR/AU/AQ/CA/ID plus island-state card/hero visual sheet.

   T0 completes only when all machine brakes pass, maxima are within
   `2500/8000`, and the owner approves the sheet. Otherwise interrupt without a
   T0 product commit or T1 publication.

2. **T1 — Publish and embed the frozen projected LOD corpus.** Rebuild from the
   T0-frozen tool and recipe. Publish compact/standard projected-coordinate
   artifacts, manifest and Go embedding adapter. Prove exact P1,
   coordinate-space metadata, projection policy, recipe, lockfile, Mapshaper and
   tier hashes, all-283 coverage, omission provenance, offline runtime, and two
   byte-identical clean maintainer rebuilds.
3. **T2 — Complete P2 evidence and the project ceiling.** Add full-corpus, D3,
   natural-ratio, arbitrary-frame, scale-only policy, projected-artifact parity,
   source-finalization, protected restoration, budget, serializer,
   visual-approval, override and marker evidence. Run focused gates and offline
   `make check`, then create the auditable P2 product commit and closure inputs.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0 | `lod-representation-brake` | `go.mod`, `go.sum`, `internal/geometry/cmd/lodbuild/**`, `internal/geometry/lod/tool/**`, `internal/geometry/lod/v1.recipe.json`, `internal/geometry/testdata/lod-spike/**`, `internal/geometry/{adapter.go,adapter_test.go,diagnostic.go,fit.go,fit_test.go,lod.go,lod_test.go,marker.go,marker_test.go,math.go,model.go,model_test.go,normalize.go,normalize_test.go,override.go,override_test.go,path.go,path_test.go,pipeline.go,pipeline_test.go,preset.go,preset_test.go,projection.go,projection_test.go,retain.go,retain_test.go,simplify.go,simplify_test.go,soften.go,soften_test.go}`, `internal/geometry/overrides/**`, `internal/geometry/presets/**` |
| T1 | `lod-publisher` | `internal/geometry/lod/v1.compact.json`, `internal/geometry/lod/v1.standard.json`, `internal/geometry/lod/v1.manifest.json`, `internal/geometry/lod_generated.go`, `internal/geometry/lod_build_test.go` |
| T2 | `lod-verifier` | `Makefile`, `internal/geometry/{approval_test.go,integration_test.go,mutations_test.go,oracle_test.go}`, `internal/geometry/testdata/{approval,d3,mutations}/**` |

T0 may reconcile `go.mod` and `go.sum`, but no new runtime dependency is
authorized. `internal/geometry/serialize.go` and `serialize_test.go` are
read-only T0 controls: the small experiment may diagnose a byte failure but may
not alter serializer behavior. No implementer writes `data/**`,
`internal/catalog/**`, `cmd/**`, lifecycle documents or another task's lease.
The coordinator owns lifecycle, evidence and audit artifacts.

## Dependencies and execution waves

Critical path is W1/T0 → W2/T1 → W3/T2, concurrency one. The AQ/RU
representation experiment is a hard checkpoint before the all-corpus matrix.
T1 requires exact passing final paths, the complete T0 receipt and explicit owner
visual approval. T2 requires committed artifact reproducibility.

## Validation plan

- `geometry-projected-lod-builder`: exact Go LAEA metadata, explicit Mapshaper
  planar mode, `64/256`, 283 IDs, deterministic rebuild and mutation teeth for
  center/flatness/axis/corpus drift.
- `geometry-final-representation`: selection-only/production parity, `q=0.01`
  topology, final symmetric deviation, protected restoration and exact
  compact/standard/source provenance.
- `geometry-source-finalization`: AQ/RU q-crossing witnesses, monotonic bounded
  source reduction, no MakeValid, no visual-limit escape and typed exhaustion.
- `geometry-lod-integrated-brake`: all-corpus/all-scale matrix over final `d`,
  layouts, topology, protection, paths, budgets and visual digest.
- `geometry-svg-serializer`: unchanged parser round-trip, deterministic shortest
  valid absolute/relative encoding and adversarial grammar teeth.
- `geometry-lod-runtime`: source/artifact binding, scale-only selection,
  fallback/cache parity and deterministic offline operation.
- `geometry-topology-teeth`: proper crossing versus genuine shared source
  endpoint, winding, holes, fixed-grid damage and softened fallback.
- `geometry-d3-parity`, `geometry-overrides`, `geometry-markers` and offline
  `make check` retain their P2 meanings.

Teeth fail if geographic coordinates are accepted as projected LOD, planar mode
is removed, metadata drifts, selection returns before q validation, source
reports zero deviation without measuring its final form, q changes, a crossing
is excused, tolerance exceeds the resolved value, selection inspects
ISO/preset/frame, a protected part is lost, final bytes drift or a hard budget
is silently exceeded.

## Deviation and amendment policy

`p2-lod-projected-coordinate-how-v3` permits helper/file splits inside the exact
T0 lease, the projected-coordinate artifact schema above, and q-aware reuse of
the existing source simplifier. It does not permit visual tolerance above
`0.012`, q other than `0.01`, different tier resolutions or thresholds, another
tier/dependency, MakeValid/union repair, changed serializer grammar, changed hard
budgets, country branches or runtime external processes.

Natural/contain semantics, P1 authority, protected-feature policy,
presentation-free CTR-002 or the single-binary offline contract require an
amendment. Unlimited-path success followed by budget failure is measured
successor input, not permission to expand T0 into speculative serializer work.

## Commit worktree and integration policy

Each task lands as one conventional atomic product commit after its complete
focused gate. Stage exact leased paths only. T0 commits only after the full
machine brake and owner visual approval; partial AQ/RU evidence remains
uncommitted scratch plus a recorded blocked result. Lifecycle evidence never
enters a product commit.

## Rollback and recovery

Compact failure advances to standard; standard advances to validated q-aware
source; invalid or over-budget source fails typed. T1 stages all artifacts and
publishes atomically only after validation. Rollback selects the previous
policy/corpus pair and reverts T2 → T1 → T0. Maintainer `node_modules` is absent
from product commits; normal generation and ordinary checks remain Node-free.

## Completion and handoff

P2 closes only after owner-approved visuals, deterministic all-283 projected
artifacts, both-profile all-249 coverage, exact restoration/provenance, natural
and arbitrary layouts, final paths within hard budgets, offline `make check`,
exact run evidence and fresh independent audit.

P3 receives CTR-002 with natural `tight`, distortion-free `contain`, versioned
data-driven defaults, explicit quality/override seams, canonical
presentation-free commands, stable diagnostics and optional capital marker
coordinates. `detail` remains deferred.

Estimated remaining effort is 8–18 agent-hours, likely 13, medium confidence.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
