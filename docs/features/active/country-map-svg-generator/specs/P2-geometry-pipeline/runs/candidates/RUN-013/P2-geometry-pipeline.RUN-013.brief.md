---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-013"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-013"
status: draft
profiles: []
concerns: []
inputs: ["PLAN-010", "CTX-RUN-015"]
---
# RUN-013 — Protected-visibility owner-gate closure

## Authority and accepted inputs

- Accepted PLAN-010:
  `a81da84fc3881068ad0520a3746f3276c18df87a48f10e83d5463fb03fb1db63`;
  manifest
  `fdd0cdfa0c545410f102eb460fef80b57d3b03484a2fc0b1cb423975a7fff59d`.
- Runtime context CTX-RUN-015:
  `207d6cc02da8ac5f6610364440ec26b15544a1d728be583ae81a8be6b813ab0b`.
- RUN-012 result:
  `5af7fe748d44bfccaa44599e2e3ed0db070002b7801aecbae56f803b0d0a5b7e`;
  coder result RESULT-020:
  `adcd37664f25d51d67a2ba3e171c5fb37c0eebb57a5a0a91301571afeb66f204`.
- Frozen pre-policy diagnostic authority remains RUN-010/RUN-011 with semantic
  digest
  `220aa8cff32ff85a15c62b45eeff79040b66d4e14401f4fe7b02dbf83aa67249`
  and Indonesia boundary `121 = 20 / 2713 / 2993`,
  `120 = 19 / 2658 / 2938`.
- Re-own the exact RUN-012 scratch and artifact hashes:
  oracle `f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`,
  recipe `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`,
  manifest `b2c467b08fb32603a312918da05ce83da355499fb8cba7f9e73c5a8e916995da`,
  sheet `cf1e568c93646047791b45a329e50a1308683e8c3aae03ee0944654db3718f9f`,
  receipt `fcc1bc307509d55aa281591f1cce0d0e78f73f323a41fc6d78e665a56c2dfd10`.

## Objective and completion boundary

Close the remaining PLAN-010 T0A gap without redoing or broadening T0B:
implement the versioned data-only add-retention group-anchor seam and its
mutation teeth, then rerun only the T0A-owned diagnostic, Indonesia,
protected-visibility and representative A/B gates.

Stop at the owner checkpoint with no product commit. A machine PASS permits
the coordinator to request approval of thresholds `compact=110`,
`standard=900`, oracle floor `IoU=.40 / recall=.42`, and the exact sheet digest.
Do not start T0B or modify inherited production-spike expectations.

## Scope ownership and footprint

Writable:

- `internal/geometry/override.go`, `internal/geometry/override_test.go`;
- `internal/geometry/overrides/**`;
- `internal/geometry/silhouette.go`, `internal/geometry/silhouette_test.go`;
- `internal/geometry/cmd/lodbuild/representative.go`;
- `internal/geometry/cmd/lodbuild/diagnostic_test.go`;
- `internal/geometry/lod/silhouette-oracle.v1.json`;
- `internal/geometry/lod/v2.recipe.json`;
- `internal/geometry/testdata/silhouette-oracle/**`.

Everything else is read-only, including `.git`, `.mate`, `docs`, `go.mod`,
`go.sum`, serializer files, `cmd/lodbuild/main.go`, production pipeline,
generated ladder tables, Makefile and inherited `main_test.go`. Preserve all
shared scratch. Do not stage, commit or mutate lifecycle.

## Context and doctrine pack

CTX-RUN-015 supplies P1 and Go doctrine. PLAN-010 and RUN-013 are the exact
HOW. The group-anchor seam is optional versioned data, initially empty unless a
specific owner-reviewed group is supplied later. It may only add a source
component to the protected identity set; it cannot evict a base member or alter
thresholds, candidate order, budgets, projection, q, serializer, selection
inputs or executable branching by entity.

Current post-policy diagnostic bytes naturally bind changed T0A source hashes;
they are not a replacement for the immutable pre-policy `220aa8...` checkpoint.
Do not relabel the current digest as the frozen input digest.

## Tasks

1. Verify all RUN-012 artifact hashes and unchanged `go.mod`, `go.sum`,
   serializer and serializer-test hashes.
2. Add the versioned group-anchor data schema/parser and an empty default data
   set. Resolve a named group anchor to its containing source component using
   the same projection and lineage path as protected anchors.
3. Make group-anchor membership additive after base identity ranking and anchor
   replacement. Never evict a base member. Record group name, review
   provenance, source order/digest and additive identity reason.
4. Mutation-test stale/invalid review provenance, anchor outside source,
   duplicate/conflicting group names, attempted threshold/budget/order change,
   base-member eviction, executable entity branch and non-additive retention.
5. Rebuild the representative matrix twice. With the empty default group set,
   all five canonical hashes and every selected candidate must remain exactly
   RUN-012-identical.
6. Re-run exact Indonesia boundary and T0A-owned diagnostic tests. Preserve the
   known full-package failures as T0B evidence; do not edit or waive them.
7. Return exact commands, hashes, green selected gates, the two inherited red
   tests as non-selected evidence, and the rendered sheet path. Wait for the
   coordinator to obtain owner approval.

## Validation and evidence

- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry -run
  'Test(Diagnostic|GridPhase|Projection|LOD|Topology|Protection|Preset|Silhouette|Oracle|Selection|Source|GroupAnchor)'
  -count=1`;
- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry/cmd/lodbuild -run '^TestDiagnostic' -count=1`;
- two independent `-representative` builds and byte comparison of oracle,
  recipe, manifest, sheet and receipt;
- one `-indonesia-boundary-out` build with exact 120/121 values;
- compile-only geometry/lodbuild checks and exact source/hash invariants.

`go test ./internal/geometry/cmd/lodbuild -count=1` is explicitly not this
run's passing gate: its inherited `TestLODSpikeFullCorpus` AE budget assertion
and `TestLODAQRUProjectionAlignedCheckpoint` AQ tier assertion are T0B inputs.
Record them unchanged and do not hide their red status.

## Constraints and forbidden actions

Keep the frozen thresholds, contribution units, band grids `240/700`, global
candidate order, `q=.01`, budgets `2200/7500` and `2500/8000`, natural/contain
semantics, pure-Go offline runtime and serializer hashes
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345` /
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

No built-in country exception, automatic owner claim, group-anchor entry
without review provenance, hand-moved point, vertex reinsertion, higher budget,
lower precision, runtime dependency, external process, serializer change or
T0B repair is authorized.

## Stop escalation and resume conditions

Stop if the empty group data changes any RUN-012 candidate or artifact hash; an
override can remove or replace a base member; provenance can be bypassed; any
selected T0A gate fails; Indonesia drifts; or closing the gap requires touching
a forbidden path.

After machine PASS, pause for explicit owner review. Owner approval freezes the
oracle, recipe, thresholds, identity policy and sheet digest. Owner rejection
returns to bounded calibration and cannot waive a machine failure.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
