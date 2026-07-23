---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "accepted-plan-revision"
id: "PLAN-001"
epic: "country-map-svg-generator"
spec: "P2"
status: accepted
profiles: []
concerns: []
inputs: ["P2"]
---
# PLAN-001 — Implementation plan

## Accepted inputs and baseline

- Normative specification: `P2-geometry-pipeline.spec.md` at SHA-256
  `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`.
- Epic architecture: `country-map-svg-generator.solution-architecture.md` at
  SHA-256
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`.
- Accepted architecture decision DEC-003 fixes a Go runtime over the immutable
  local corpus. P2 consumes P1 contract CTR-001 and produces CTR-002.
- Repository baseline: commit
  `80c50352615b5d55a6efdc2fdffc5d9db3c9fa88`, tree
  `53948381a787202e9edf5e200749743e6220325c`.
- P2 exclusively owns BND-003 (`internal/geometry/**`). It may extend the shared
  Go module and project `Makefile`, but does not own SVG/XML styling, CLI/config
  parsing, source-corpus mutation, public packaging, or browser integration.
- Normal execution remains pure Go, local, deterministic, and offline. Node is
  allowed only to regenerate committed D3 oracle fixtures; Python, GDAL, GEOS,
  PROJ, and network access are not runtime or project-gate dependencies.
- `card` and `hero` are named data presets, never engine modes. The default
  `tight` layout derives a natural-aspect viewBox from projected geometry; the
  explicit `contain` layout accepts any finite positive frame and valid padding,
  including `19x19` and `10000x12345`, without stretching. Detail-mode, public
  distribution, P3 markup/CLI, and P4 browser proof remain deferred.

## Requirement and acceptance coverage

| Contract | Tasks | Required evidence |
| --- | --- | --- |
| US-1 | T3, T5 | Every P1 entity/profile produces a valid deterministic geometry result or a typed actionable error. |
| US-2 | T1–T3 | One layout-agnostic contract produces natural-aspect `tight` output and handles arbitrary portrait, landscape, square, tiny, and very large `contain` frames without preset branches or distortion. |
| US-3 | T1, T3, T5 | Auto quality, explicit tolerances, presets, and versioned overrides are inspectable and reproducible. |
| US-4 | T2, T4, T5 | Projection and fitting match committed D3 oracle fixtures within declared tolerances. |
| REQ-1 | T1, T2, T5 | Natural-aspect `tight` viewBox, distortion-free `contain` fitting, finite frame/padding validation, positive drawable area, typed overflow/point-ceiling failures, and advisory extreme-size diagnostics. |
| REQ-2 | T1, T3, T5 | `card`/`hero` resolve from versioned data through the same generic layout/quality contract; every field remains overrideable. |
| REQ-3 | T3, T5 | Simplification preserves valid rings, winding, non-empty identity geometry, and all P1 protected features. |
| REQ-4 | T3, T5 | Softening is bounded by resolved target-pixel tolerance and cannot introduce invalid topology. |
| REQ-5 | T3, T5 | Every removed or collapsed component/ring is explicitly diagnosed with provenance and reason; implicit library omission cannot bypass accounting. |
| REQ-6 | T1, T3, T5 | Versioned per-country data overrides are bounded, inspectable, applicability-scoped, and require no code fork. |
| REQ-7 | T2, T3, T5 | Geometry and optional capital/custom markers use the exact same projection and fit transform; markers do not influence fit. |
| REQ-8 | T2, T3, T5 | Canonical ordering, quantization, commands, diagnostics, and markers produce byte-equivalent output for equal input. |
| INV-1 | T1, T3, T5 | CTR-002 is presentation-free and records all inputs, transforms, removals, provenance, metrics, diagnostics, and versions. |
| AC-1 | T3, T5 | Approved representative raster evidence remains recognizable and visually balanced while structural evidence proves topology and protection. |
| AC-2 | T1, T2, T4, T5 | `tight` natural aspect, arbitrary distortion-free `contain`, fitted-scale quality, D3 parity, and repeated byte equivalence are proven. |
| AC-3 | T1, T3, T5 | Named exceptional countries are expressed through bounded versioned data overrides without engine code forks. |
| AC-4 | T3, T5 | Topology, smoothing, disappearance-accounting, and protected-feature mutation teeth all fail closed. |
| VAL-1–5 | T3–T5 | Focused oracle, mutation, visual/budget, override, and marker gates plus the project ceiling provide evidence. |

## Technical approach

Create a pure-Go `internal/geometry` package with a closed input/output contract.
Input contains a P1 entity/profile geometry, a layout request, an auto or explicit
quality policy, a versioned geometry override, and optional geographic markers.
`tight` accepts a rendered long-side target or maximum width/height box, derives
the other viewBox dimension from projected geometry, and crops to geometry plus
padding. `contain` accepts explicit finite positive width and height, uniformly
fits and centers geometry, and preserves transparent unused frame area. Both use
four nonnegative padding insets and one uniform scale; neither can stretch a
country. CTR-002 contains schema/algorithm/corpus identities, layout mode, natural
projected aspect ratio, resolved preset and quality, normalized presentation-free
commands, bounds and fit transform, ring/component provenance, markers, anomalies,
point/path metrics, removals and reasons, protected-feature retention, overrides,
and stable diagnostics.

No dimension or aspect-ratio whitelist exists. Padding must leave a positive
drawable extent; non-finite values, unsafe arithmetic, invalid projection, or more
than 2,000,000 projected points fail with typed errors. Safe extremes pass with
diagnostics. The effective scale `d` is the long side of the fitted geometry in
target pixels, not the shorter frame dimension; `d < 32` warns about low
resolution and `d > 2,048` warns about generation cost. Dimensions are
optimization hints in target pixel units; the resulting vector remains scalable.

Auto quality derives only from that effective rendered geometry scale `d`:

- flatness: `clamp(0.10 + 16/d, 0.10, 0.30)` px;
- simplification: `clamp(0.25 + 64/d, 0.25, 1.25)` px;
- softening: `clamp(0.12 + 24/d, 0.12, 0.35)` px;
- minimum component area: `clamp(0.50 + 112/d, 0.50, 1.50)` px²;
- canonical quantization: `0.01` target-pixel units.

Versioned `presets/v1.json` supplies convenience defaults only: `card` uses
`tight`, a 128 px fitted-geometry long side, 8 px padding added around its natural
projected bounds, and auto quality (reference
checks at 90/128/160/240 rendered long-side pixels, advisory 2,200 path bytes);
`hero` uses `tight`, a 640 px fitted-geometry long side and 32 px padding added
around its natural projected bounds (checks at
500/640/700, advisory 7,500 bytes). A caller can override every resolved field or
choose `contain`; P3 may later load additional named presets without a Go fork.

Normalize geographic rings, compute a spherical centroid, center a Lambert
azimuthal equal-area projection per entity, adaptively subdivide great-circle
segments, reflect Y for SVG coordinates, then uniformly fit the raw projected
geometry. In `tight`, derive the final viewBox from fitted geometry plus padding;
in `contain`, fit into the padded frame and center unused axes. Optional markers
use the exact same transform and never affect fitting. A test-only oracle mode
disables reduction/softening and compares committed fixtures generated by D3 Geo.

After raw fitting, work in target-pixel units. Keep the largest component, every
anchor-bearing component, and the configured minimum part count. Validate explicit
retention anchors. Remove only remaining unprotected components below the resolved
area threshold and record each removal. Use the validated RDP implementation from
`github.com/peterstace/simplefeatures/geom` v0.59.0; halve tolerance on topology
damage down to zero and then fail typed. Never use an unchecked simplifier.

Softening converts eligible corners to bounded quadratic segments. Entry/exit
distances are capped by the resolved tolerance, one quarter of adjacent edge
lengths, and local clearance. Collinear, short, or low-clearance corners are
skipped. Flatten below `min(0.05px, tolerance/8)` for complete-geometry validation;
fall back from corner to ring to unchanged geometry. No buffer, union, bridging,
manual redraw, or invented geography is permitted.

Canonicalization fixes winding, rotates each quantized ring to its lexicographically
smallest start, and stably sorts holes and polygons. Commands use uppercase
`M/L/Q/Z`, fixed decimals, no negative zero, and no redundant segments. Diagnostics
and markers are stably sorted, so identical inputs produce identical bytes.

Versioned `overrides/v1.json` may adjust projection center/orientation, padding,
quality tolerances, retention anchors/minimum parts, and marker offsets. Every
entry names ISO, corpus applicability, reason, and version. Overrides cannot
replace geometry, fabricate a component, or weaken protected-feature constraints.

Production dependencies are the Go standard library plus
`github.com/peterstace/simplefeatures/geom` v0.59.0. The deterministic visual
evidence rasterizer `golang.org/x/image/vector` v0.44.0 is test-only. D3 Geo 3.1.1
is a development oracle whose generated fixtures are committed; Node is never
needed for normal generation or `make check`.

## Tasks and completion conditions

1. **T1 — Freeze the layout-agnostic contract.** Add the selected Go dependencies,
   CTR-002 models, typed diagnostics, preset schema/data, auto quality, `tight` and
   `contain` validation, and P1 adapter. Done when presets resolve through the same
   generic contract as explicit `19x19`, portrait, landscape, and `10000x12345`
   frames; the contract exposes no independent X/Y geometry scale; unsafe values
   fail typed; and safe extremes warn rather than being rejected by a whitelist.
2. **T2 — Implement normalization, projection, and fitting.** Add canonical ring
   normalization, spherical centroid, Lambert azimuthal equal-area projection,
   adaptive subdivision, natural-aspect tight bounds, padding-aware uniform
   contain fit, and the shared marker transform. Done when projection invariants,
   antimeridian/polar fixtures, padding, and target-pixel coordinates pass
   deterministic tests; `RU`, `CL`/`AR`, and `AU` derive materially wide, tall,
   and near-square `tight` viewBoxes; changing only target long side preserves
   aspect ratio; and tiny, huge, portrait, landscape, and square `contain` frames
   use one uniform scale with centered unused axes.
3. **T3 — Implement reduction, retention, softening, markers, overrides, and
   canonical commands.** Complete the production pipeline with topology validation
   and rich provenance/metrics. Compare pre/post simplification component and ring
   provenance; reject or explicitly diagnose every disappearance, including
   collapsed interior rings, so an implicit library omission cannot bypass REQ-5
   or protected-feature checks. Done when the full P1 corpus succeeds for both
   profiles, protected features survive, every removal/collapse has a stable
   reason, explicit overrides validate, markers remain bounded, and repeated calls
   are byte-identical.
4. **T4 — Commit the D3 oracle.** Add a pinned, maintainer-only D3 fixture
   generator and committed oracle outputs. Done when the Go projection matches the
   fixture suite within declared numerical tolerances and `make check` needs no
   Node installation or network.
5. **T5 — Add mutation, integration, budget, and visual approval evidence.** Add
   full-corpus integration tests, sabotage tests, deterministic approval-sheet
   rendering, and wire the ceiling into `Makefile`. Fitted-scale mutations must
   prove that different frame shapes with the same fitted geometry scale resolve
   identical quality, that a changed fitted scale in the same frame changes
   quality, and that using the frame's shorter side fails. Done when all named
   mutations red, the clean corpus greens, arbitrary layouts and preset budgets
   are visible, and the approval sheet is explicitly reviewed without silently
   refreshing baselines.

## Ownership and write footprint

- T1 (`geometry-contract`) owns `go.mod`, `go.sum`,
  `internal/geometry/{model.go,model_test.go,diagnostic.go,preset.go,preset_test.go,adapter.go,adapter_test.go}`,
  and `internal/geometry/presets/v1.json`.
- T2 (`projection-kernel`) owns
  `internal/geometry/{normalize.go,normalize_test.go,projection.go,projection_test.go,fit.go,fit_test.go}`
  and `internal/geometry/testdata/projection/**`.
- T3 (`geometry-pipeline`) owns
  `internal/geometry/{simplify.go,simplify_test.go,retain.go,retain_test.go,soften.go,soften_test.go,marker.go,marker_test.go,override.go,override_test.go,path.go,path_test.go,pipeline.go,pipeline_test.go}`,
  `internal/geometry/overrides/v1.json`, and
  `internal/geometry/testdata/pipeline/**`.
- T4 (`d3-oracle`) owns `internal/geometry/oracle_test.go` and
  `internal/geometry/testdata/d3/{package.json,package-lock.json,generate.mjs,README.md,fixtures/**}`.
- T5 (`geometry-verifier`) owns
  `internal/geometry/{integration_test.go,mutations_test.go,approval_test.go}`,
  `internal/geometry/testdata/{mutations/**,approval/**}`, and `Makefile`.
- Lifecycle, plan, run, audit, readiness, and evidence artifacts remain
  coordinator-owned. No task writes `cmd/**`, `internal/config/**`,
  `internal/render/**`, public docs, or release packaging.

## Dependencies and execution waves

- **W1:** T1, concurrency 1.
- **W2:** T2 after T1, concurrency 1.
- **W3:** T3 and T4 after T2, concurrency 2; their leases do not overlap.
- **W4:** T5 after T3 and T4, concurrency 1.

Critical path: T1 → T2 → T3 → T5. T4 can run alongside T3 after the projection
contract is frozen.

## Validation plan

- `geometry-d3-parity`:
  `GOWORK=off go test ./internal/geometry -run TestD3Oracle -count=1`.
  Teeth include wrong rotation sign, insufficient subdivision, padding drift, and
  winding changes.
- `geometry-topology-teeth`:
  `GOWORK=off go test ./internal/geometry -run
  'Test(Mutation|Topology|Protected|Soften)' -count=1`.
  Teeth include bow ties, escaped holes, ring collapse, curve crossing, protected
  removal, and unchecked simplification.
- `geometry-visual-and-budget` (project evidence slot): `make check`. It includes
  deterministic approval-sheet, diagnostics, arbitrary viewport, byte-budget,
  full-corpus, and all focused Go tests. Teeth include approval digest drift,
  missing removal reasons, hard-coded card/hero engine branches, fixed-aspect
  `tight` output, stretched `contain` output, arbitrary-size rejection, and
  suppressed extreme-size warnings. Additional teeth keep fitted geometry scale
  constant across different frame shapes, change fitted scale inside one frame,
  and replace fitted scale with the frame's shorter side: the first must preserve
  quality, the second must change it, and the third implementation mutation must
  fail. Component or interior-ring disappearance without an explicit stable
  diagnostic must also red.
- `geometry-overrides`:
  `GOWORK=off go test ./internal/geometry -run TestOverride -count=1`.
  Invalid orientation, padding, tolerance, anchor, marker, applicability, and
  unknown fields must fail closed.
- `geometry-markers`:
  `GOWORK=off go test ./internal/geometry -run TestMarker -count=1`.
  A separate marker transform, non-finite/outside result, and suppressed
  near-boundary anomaly must red.

The approval sheet covers normal (`BR`, `FR`, `US`), elongated (`CL`),
archipelagos (`ID`, `JP`), antimeridian (`FJ`, `KI`, `RU`), polar (`AQ`, `RU`),
microstates (`MC`, `SM`, `VA`), small island (`NR`), and disputed-profile
fixtures (`CN`, `CY`, `IL`, `IN`, `RU`) in both profiles. Columns show raw
projection, final custom viewport, card reference sizes, and hero reference
sizes; annotations show tolerances, point counts, exact path bytes, removals,
protected parts, diagnostics, and overrides. Tests may detect a changed approval
digest but may not approve it.

Structural annotations explicitly compare `RU`, `CL`/`AR`, and `AU` natural
viewBox ratios, target-scale-only aspect invariance, and centered `contain`
transforms. Paired `contain` frames with equal fitted geometry scale but different
unused space must resolve identical quality; a changed fitted scale must change
quality.

## Deviation and amendment policy

`p2-layout-geometry-local-how-v1`: helper names, file splits inside the declared roots,
numerically equivalent implementation details, and fixture layout are local HOW.
A new root, dependency, formula, preset field, oracle tolerance, point ceiling, or
override field requires a reviewed plan revision before mutation.

Stop and return to the teamlead for a hard-coded preset branch or size whitelist,
a `detail` preset, projection-policy change, protected-feature-policy change,
fabricated or manually redrawn geography, weaker validation, runtime use of a
forbidden tool/network, or movement of P3 CLI/render concerns into P2.

## Commit worktree and integration policy

Execute one conventional atomic commit per accepted task lease. Before each wave,
the coordinator verifies the expected parent and lease; after each task, verifies
the actual diff, parent, focused gate, and combined gate. Parallel W3 tasks must
not write each other's paths. Implementers stage only their declared session
files and do not commit lifecycle artifacts. Integration uses compare-and-swap
semantics against the accepted baseline and prior accepted task commits.

## Rollback and recovery

Every task is recoverable at its atomic commit. Fixture generation stages into an
explicit temporary directory and only reviewed bytes enter the tree. On failure,
discard only that temporary output and retain the last green commit. A changed P1
corpus identity, specification, architecture, selected dependency, or baseline
invalidates the plan before redispatch. Never recover by weakening topology
checks, hand-editing oracle output, dropping protected parts, or special-casing a
country in the engine.

## Completion and handoff

P2 completes only with one accepted PLAN-001, a governed implementation run, all
five validation obligations, deterministic full-corpus output for both profiles,
D3 and arbitrary-viewport evidence, mutation teeth, explicit human approval of
the visual sheet, green `make check`, and one fresh independent audit.

The P3 handoff is CTR-002: a general layout/quality API, natural-aspect `tight`
output, distortion-free arbitrary `contain` frames, versioned convenience presets,
custom dimensions and overrides, canonical presentation-free commands, stable
diagnostics/metrics/provenance, and optional transformed capital markers. P3 may
serialize and style the commands but may not add geography heuristics or
reinterpret P1/P2 policy.

Estimated effort is 10–20 agent-hours, likely 14, medium confidence. The bounded
uncertainty is topology-safe softening across the full corpus and calibration of
oracle/visual tolerances; failures become explicit versioned overrides or a plan
revision, never hidden country branches.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
