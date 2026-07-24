---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-011"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-011"
status: draft
profiles: []
concerns: []
inputs: ["PLAN-009", "CTX-RUN-013"]
---
# RUN-011 — Silhouette oracle and representative scale-ladder brake

## Authority and accepted inputs

- Accepted PLAN-009 revision:
  `a406405b5895244bba918d5c9a1cb913a1144ae08a7785ea441f58cf87baeea8`.
- Accepted PLAN-009 manifest:
  `7fcaf58708caeba7bba7097132c86b898c8bf74d650b9efc91934cf512698e75`.
- Runtime context CTX-RUN-013:
  `ee974b14dc4bba4164644ba8d435d7897bdb9798a26d647a1cc0a3fdb437c723`.
- RUN-010 immutable checkpoint:
  run result `RUN-010-RESULT`
  `264b1af5262d3c4f88fe9e62f24f29716061067d2e0df87a8571127759914579`,
  agent result `RESULT-018`
  `e637a82c748a648814fe74ae5d8ea8723cba4d3c0758beb98c18df46d9b028c5`,
  gate receipt `GATE-533A8DF7AE6B`
  `b5b92cb89caffc1fb9dc4d81131c7ac58d008069a001f7f235fee1e8ad849ed7`,
  and verification `VERIFY-533A8DF7AE6B`
  `b6d913896075cbb5d946fae001f6bf067011dd5bce174b1045370555b9476245`.
- T0A.1 canonical semantic digest:
  `220aa8cff32ff85a15c62b45eeff79040b66d4e14401f4fe7b02dbf83aa67249`;
  exact executing-Go-source inventory digest:
  `9974b3cf551d4d8c162dceafef6f42a1ed2a44361070684bbb00dfc6541b6e9e`.
- Product baseline remains PLAN-009's
  `2a7d45953ac857af185a4825ceea7097e0af3a40`, tree
  `a61a0f87122e9094f444fd032dc54a2c04ad3c4b`; lifecycle commits do not
  authorize unrecorded geometry drift.

## Objective and completion boundary

Execute PLAN-009 T0A.2 only. Implement and mutation-test
`silhouette-oracle/v1`, derive a deterministic fine-to-coarse candidate ladder
for the representative matrix, calibrate one global IoU/recall threshold pair
per generic scale band, prove budget and selection-policy invariants, and render
one contact sheet from the exact machine-accepted candidates.

Stop at the T0A.2 owner checkpoint. Do not start T0B, scan or publish the full
996-output catalog, create a product commit, wire public generation, or publish
the all-283 ladder. A passing run requires both machine evidence and explicit
owner approval of the two global threshold pairs and the representative sheet.

## Scope ownership and footprint

The implementer may edit only:

- `internal/geometry/cmd/lodbuild/**`;
- `internal/geometry/diagnostic.go`;
- `internal/geometry/gridphase.go`,
  `internal/geometry/gridphase_test.go`;
- `internal/geometry/lod.go`, `internal/geometry/lod_test.go`;
- `internal/geometry/model.go`, `internal/geometry/model_test.go`;
- `internal/geometry/preset.go`, `internal/geometry/preset_test.go`,
  `internal/geometry/presets/**`;
- `internal/geometry/projection.go`,
  `internal/geometry/projection_test.go`;
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
tables, Makefile and integration/mutation suites outside the lease are
forbidden. Existing scratch is shared: preserve unrelated bytes and do not
revert another owner's work. Integration target is `main`; return without
commit, staging or lifecycle mutation.

## Context and doctrine pack

CTX-RUN-013 supplies the P1 and Go implementation doctrine. PLAN-009 is the
exact HOW for T0A.2. Re-run the content-addressed T0A.1 diagnostic before
building on it and require semantic digest
`220aa8cff32ff85a15c62b45eeff79040b66d4e14401f4fe7b02dbf83aa67249`.

Runtime and oracle remain pure Go and offline. Pinned Mapshaper `0.7.44` is
maintainer-only and may be restored from the existing lock solely for candidate
derivation. Never stage `node_modules`, temporary diagnostics, or timing
receipts. Keep `go.mod`
`4412d79908d733d715bad8c17334eafbdabb6f96ab8b242d87bdf6637ee9f13b`,
`go.sum`
`139a8de4191f92748a9d3a8c3868f56862e42f1b0d60d8f1dc4a582d273b5daf`,
and the serializer hashes named below unchanged.

## Tasks

1. Re-run T0A.1 twice from identical inputs; require byte identity, the exact
   semantic digest, 268 parameter/phase rows, 996 catalog rows, two typed
   failures, the frozen aggregates, and all six contradiction predicates
   false.
2. Implement a versioned pure-Go `silhouette-oracle/v1` contract and parser.
   Bind rasterizer source, P1/projection/recipe/tool identities, scale-band
   boundaries, natural fit, grid derivation, pixel-center nonzero fill, no
   antialias/stroke/color/softening, IoU, reference recall, protected-anchor
   coverage, dominant-component coverage, candidate order and digest tie-break.
3. Add mutation teeth for stale identity, fill-rule/sample/grid drift, threshold
   drift, lost protected land, dominant-component loss, candidate reordering,
   nondeterministic tie-break, hidden omission, byte overflow and oracle bypass.
4. Build deterministic fine-to-coarse presimplification candidates for
   `RU,CA,CN,AQ,AE,ID,CL,AR,AU` plus protected island and microstate fixtures.
   Record exact candidate parameters, digests, removed components and protected
   anchors. Use one versioned visibility policy; no entity or preset branch.
5. Calibrate one global IoU/recall threshold pair for the compact band
   (maximum effective scale `240`) and one for the standard band (maximum
   effective scale `700`). Search fine-to-coarse and select the first candidate
   passing topology, protection, oracle, visibility and path-budget checks.
6. Exercise natural output and portrait, landscape and square `contain` frames
   at multiple fitted scales on both sides of every band boundary. Prove that
   equal effective scale, requested quality and byte policy select the same
   candidate regardless of entity token, preset token or contain-frame shape.
7. Require every representative automatic output to fit path caps
   `2200/7500` and deterministic complete-file estimates `2500/8000`, retain
   protected land, preserve non-empty polygon topology and repeat byte-for-byte.
   Keep explicit `source` and custom-quality behavior on the DEC-005 q-aware
   deviation guard with typed hard-budget failure.
8. Render a standalone representative contact sheet and manifest from the
   exact accepted candidates. Show country/profile, effective scale, selected
   band/candidate, path bytes, complete estimate, IoU, recall and omissions.
   Produce a machine receipt binding sheet bytes, manifest and thresholds.
9. Return exact commands, changed-path hashes, rebuild digests, focused test
   output, machine verdict and sheet path. Wait for the owner to approve or
   reject both threshold pairs and the visual sheet before run completion.

## Validation and evidence

Required gates at the project ceiling:

- `geometry-diagnostic-reproduction`: exact T0A.1 input and semantic digest;
- `geometry-silhouette-oracle-v1`: contract/parser/rasterizer identity,
  natural grid, metrics, global thresholds, protection and mutations;
- `geometry-scale-ladder-spike`: fine-to-coarse order, stable tie-break,
  visibility provenance, boundary behavior, budgets, envelope estimates,
  determinism and contact-sheet identity;
- `geometry-selection-policy`: selection independence from preset/entity/frame
  shape for equal resolved inputs;
- `geometry-source-contract`: unchanged explicit source/custom DEC-005 guard.

Run focused tests for projection, gridphase, topology, protection, preset,
oracle, silhouette, selection and lodbuild plus compile-only checks for touched
packages. Build the representative artifacts twice and require byte-identical
oracle, recipe, machine receipt, sheet manifest and sheet.

Evidence must bind exact input hashes; Go identity; rasterizer source; candidate
matrix/order; thresholds; scale probes; selected identities; topology and
protection results; omissions; points; path/file estimates; output digests;
mutation failures; test commands and exit codes. Owner evidence must quote the
approved oracle/recipe/sheet digests, not a chat-only description.

## Constraints and forbidden actions

Keep projection v1, `q=0.01`, compact/standard maximum effective scales
`240/700`, path caps `2200/7500`, complete-file maxima `2500/8000`, natural and
uniform-contain semantics, and frozen serializer hashes
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

No country/preset threshold or candidate branch, raw-deviation acceptance for
automatic silhouettes, custom vertex reinsertion, hand-moved coordinate,
runtime repair, phase refinement, lower precision, higher budgets, post-hoc
SVG optimizer, MakeValid, union/buffer, CGO, new runtime dependency, external
runtime process, serializer change, generated publication or broad cleanup is
permitted.

Threshold calibration may move only before owner freeze and must remain one
global pair per generic band. Candidate tuning is limited to the recorded
fine-to-coarse set. Every omission requires deterministic versioned provenance.

## Stop escalation and resume conditions

Stop without product commit and return machine failure when the T0A.1 digest
drifts; a global threshold pair cannot satisfy the representative matrix; a
protected microstate cannot fit; a candidate is visibly unrecognizable; reserve
is insufficient; selection inspects preset/entity/frame shape; an automatic
candidate bypasses the oracle; rebuild bytes drift; runtime needs an external
process; or serializer/dependency edits are requested.

Pause for explicit owner review after machine PASS. Owner rejection of the
thresholds or sheet returns to bounded pre-freeze calibration in a successor
run; it does not authorize T0B or relaxed budgets. Owner approval freezes
oracle, recipe, visibility policy and sheet digest, after which any change to
those surfaces requires replan and a new version.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
