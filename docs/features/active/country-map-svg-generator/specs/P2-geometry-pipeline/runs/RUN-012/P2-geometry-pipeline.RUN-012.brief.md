---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-012"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-012"
status: running
profiles: []
concerns: []
inputs: ["PLAN-010", "CTX-RUN-014"]
---
# RUN-012 — Scale-aware protected-visibility oracle and owner brake

## Authority and accepted inputs

- Accepted PLAN-010 revision:
  `a81da84fc3881068ad0520a3746f3276c18df87a48f10e83d5463fb03fb1db63`.
- Accepted PLAN-010 manifest:
  `fdd0cdfa0c545410f102eb460fef80b57d3b03484a2fc0b1cb423975a7fff59d`.
- Runtime context CTX-RUN-014:
  `f620f2bf164254fec15b8cb5cd5b6b5f8bea1a20cf5ffcbe5beeb901f3a8ad31`.
- Accepted DEC-007:
  `e927675a05f1cee0f4e39532a4f600609fa9560f0e97eca318d7af1ba231e9ef`.
- RUN-010 result:
  `264b1af5262d3c4f88fe9e62f24f29716061067d2e0df87a8571127759914579`.
- RUN-011 result:
  `011ca94ddd060d4f80d1da6894127c6210bf1317a572d861bce0f77ec09f1b32`.
- Product baseline remains PLAN-010's
  `2a7d45953ac857af185a4825ceea7097e0af3a40`, tree
  `a61a0f87122e9094f444fd032dc54a2c04ad3c4b`; later commits are lifecycle
  authority only and do not authorize unrecorded product-byte drift.

## Objective and completion boundary

Execute PLAN-010 T0A only. First reproduce the immutable RUN-010 diagnostic
and RUN-011 Indonesia boundary. Then implement and mutation-test
`protected-visibility/v1`, calibrate one global contribution threshold per
generic scale band, and generate a deterministic representative machine sheet
from the exact candidates that pass every topology, oracle and byte brake.

Stop at the owner checkpoint. Do not start T0B, publish the all-283 ladder,
scan the 996-output catalog, wire public generation, or create a product
commit. Machine PASS permits visual review but does not complete the run:
explicit owner approval must bind both global thresholds and the exact sheet
digest.

## Scope ownership and footprint

The implementer owns only:

- `internal/geometry/cmd/lodbuild/**`;
- `internal/geometry/diagnostic.go`;
- `internal/geometry/gridphase.go`,
  `internal/geometry/gridphase_test.go`;
- `internal/geometry/lod.go`, `internal/geometry/lod_test.go`;
- `internal/geometry/model.go`, `internal/geometry/model_test.go`;
- `internal/geometry/override.go`, `internal/geometry/override_test.go`;
- `internal/geometry/overrides/**`;
- `internal/geometry/preset.go`, `internal/geometry/preset_test.go`,
  `internal/geometry/presets/**`;
- `internal/geometry/projection.go`,
  `internal/geometry/projection_test.go`;
- `internal/geometry/retain.go`, `internal/geometry/retain_test.go`;
- `internal/geometry/silhouette.go`,
  `internal/geometry/silhouette_test.go`;
- `internal/geometry/lod/v1.recipe.json`,
  `internal/geometry/lod/v2.recipe.json`,
  `internal/geometry/lod/silhouette-oracle.v1.json`;
- `internal/geometry/lod/tool/build.mjs`;
- `internal/geometry/testdata/silhouette-oracle/**`.

All other paths are read-only. In particular `.git`, `.mate`, `docs`, `go.mod`,
`go.sum`, `internal/geometry/serialize.go`,
`internal/geometry/serialize_test.go`, public pipeline files, generated ladder
tables, Makefile and T2 integration/mutation suites are forbidden. Existing
scratch is shared: preserve unrelated bytes and never revert another owner's
work. Return without staging, commit or lifecycle mutation.

## Context and doctrine pack

CTX-RUN-014 supplies the current P1 specification and Go implementation
doctrine. PLAN-010 is the exact HOW. Re-own the existing RUN-011 scratch only
after both frozen regressions pass.

Runtime and oracle remain pure Go and offline. Pinned Mapshaper `0.7.44` is
maintainer-only and may be restored from the existing lock solely for candidate
derivation. Never stage `node_modules`, temporary diagnostics or rendered
review artifacts. Keep `go.mod`
`4412d79908d733d715bad8c17334eafbdabb6f96ab8b242d87bdf6637ee9f13b`,
`go.sum`
`139a8de4191f92748a9d3a8c3868f56862e42f1b0d60d8f1dc4a582d273b5daf`,
and the serializer hashes named below unchanged.

## Tasks

1. Re-run RUN-010 twice from identical inputs. Require byte identity, semantic
   digest
   `220aa8cff32ff85a15c62b45eeff79040b66d4e14401f4fe7b02dbf83aa67249`,
   268 parameter/phase rows, 996 catalog rows, two typed failures, frozen
   aggregates and all contradiction predicates false.
2. Re-run the diagnostic-only unconditional-minimum Indonesia boundary.
   Resolution `121` must be the coarsest candidate retaining 20 parts and emit
   exactly `2713` path bytes and `2993` estimated file bytes; resolution `120`
   must retain 19 parts. Production selection must not call this brake.
3. Implement versioned pure-Go `protected-visibility/v1`. Rank source
   components by independent reference-grid filled-pixel contribution,
   source order and content digest. Build exactly
   `min(minimum_parts, source component count)` base identity members, replacing
   the lowest non-anchor when an anchor falls outside the initial top N.
4. Record for every source component its source order/digest, contribution,
   identity rank/reason, dominant/anchor flags, band threshold,
   visible/subscale class, candidate lineage/contribution, retained state,
   omission reason and tie-break keys. Missing lineage, false
   `protected_subscale`, visible protected loss, dominant loss or anchor loss
   fails closed.
5. Add a versioned data-only group-anchor override seam. It may add retention
   after owner evidence; it may not evict a base member or change thresholds,
   budgets, candidate order or executable control flow.
6. Calibrate one integer contribution threshold for compact and one for
   standard over RU, CA, CN, AQ, AE, ID, CL, AR, AU plus protected-island and
   microstate fixtures. Search fine-to-coarse only. Larger bands may not lose
   an identity member retained by smaller bands.
7. Require every representative automatic output to pass topology, global
   IoU/recall, dominant/anchor coverage, path caps `2200/7500`, deterministic
   complete-file estimates `2500/8000`, and byte-identical rebuilds. Explicit
   `source` and custom quality retain P1 `minimum_parts` and the DEC-005
   q-aware typed budget failure.
8. Exercise natural output and portrait, landscape and square `contain` frames
   on both sides of every band boundary. Equal effective scale, quality and
   byte policy must select the same candidate regardless of preset, entity or
   frame shape.
9. Mutation-test dominant/anchor/rank/tie drift, visible protected loss, false
   subscale provenance, non-monotonic retention, entity/preset/frame branching,
   additive-override violations, oracle bypass and byte overflow.
10. Generate a standalone representative sheet and machine manifest from the
    exact passing candidates. Label country/profile, effective scale, selected
    candidate, bytes, IoU/recall, retained/omitted identity components,
    contribution, threshold and omission provenance. Rebuild twice and bind
    exact sheet, manifest, oracle and recipe digests.
11. Return exact commands, hashes, focused test output and machine verdict.
    Pause for owner approval of both thresholds and the exact sheet digest.

## Validation and evidence

Required gates:

- `geometry-diagnostic-reproduction`;
- `geometry-indonesia-120-121-regression`;
- `geometry-protected-visibility-v1`;
- `geometry-scale-ladder-spike`;
- `geometry-selection-policy`;
- `geometry-source-contract`.

Run focused projection, gridphase, topology, protection, preset, oracle,
silhouette, selection and lodbuild tests plus compile-only checks for touched
packages. Rebuild all representative artifacts twice.

Evidence must bind exact input hashes, Go/rasterizer identity, candidate
matrix/order, thresholds, scale probes, selected identities, topology and
protection results, omissions, points, path/file estimates, output digests,
mutation failures, commands and exit codes. The owner receipt must quote the
approved oracle, recipe and sheet digests. Chat-only approval is insufficient.

## Constraints and forbidden actions

Keep projection v1, `q=0.01`, compact/standard maximum effective scales
`240/700`, path caps `2200/7500`, complete-file maxima `2500/8000`, natural and
uniform-contain semantics, and frozen serializer hashes
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

No country/preset threshold or candidate branch, vertex reinsertion,
hand-moved coordinate, runtime repair, phase refinement, lower precision,
higher budget, post-hoc SVG optimizer, MakeValid, union/buffer, CGO, runtime
dependency, external runtime process, serializer change, generated
publication or broad cleanup is permitted.

Threshold calibration may move only before owner freeze and remains one global
integer threshold per generic band. The bounded override is data-only and
additive. Every omission requires deterministic versioned provenance.

## Stop escalation and resume conditions

Stop without product commit if either frozen regression drifts; no global
threshold per band satisfies all machine gates; a dominant or anchor component
cannot fit; a subscale member remains visually identity-defining; retention is
non-monotonic; selection observes forbidden inputs; rebuilds drift; the runtime
needs an external process; or serializer, q, budget or dependency changes are
requested.

Pause only after machine PASS and expose the exact sheet to the owner. Owner
rejection returns to bounded pre-freeze calibration; it cannot waive a machine
failure or authorize T0B. Owner approval freezes oracle, recipe, thresholds,
identity provenance and sheet digest. Any later change requires replan and a
new version.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
