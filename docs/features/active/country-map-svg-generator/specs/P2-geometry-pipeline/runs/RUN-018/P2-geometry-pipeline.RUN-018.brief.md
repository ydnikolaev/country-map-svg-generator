---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-018"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-018"
status: running
profiles: []
concerns: []
inputs: ["PLAN-011", "CTX-RUN-020"]
---
# RUN-018 — Component-lineage source and representative checkpoint

## Authority and accepted inputs

- Accepted PLAN-011 plan
  `d9dac64bef9f01d32f8ca4a5905b57388f0849e531e566ca2d39a04696d9e799`;
  manifest
  `75fba90898b5519fb1578a347e13a9ea79b8d1f7e0de2ca5bfb3c45f6e017e1d`.
- Runtime context CTX-RUN-020:
  `f67d818d1b5b900067a568c8d5ceef92c48f4d029146ddc7bf7af6fb5b360403`.
- RUN-017 diagnostic result
  `ff99d94cbcefd9b10943f3d81ea182d035862ab7d62a302b0a2a5666c6489f07`;
  semantic digest
  `94b756629056505eb85d0af568a44d478173b79b0aa8bf3613567c7f87179c54`.
- Frozen oracle
  `f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`;
  predecessor recipe
  `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`.
- Current 51-file geometry scratch inventory
  `d4d16e5af1b77a02dd4163078558e0a80ebf32c5599ea3132f53e23ab617bb61`.

## Objective and completion boundary

Execute PLAN-011 T0A only. Replace anonymous multipart candidate input with
`lineage-features/v1`, prove the hard SH/CA/ID/KI micro-brake, then rebuild and
machine-validate the representative matrix under unchanged policy. Produce the
new deterministic contact sheet and stop for owner approval.

Do not integrate runtime `GenerateWithLOD`, run the 996-output T0B catalog,
publish T1 artifacts, create a product commit or seek to reuse DEC-008.

## Scope ownership and footprint

Writable:

- `internal/geometry/cmd/lodbuild/**`;
- `internal/geometry/lod/tool/build.mjs`;
- `internal/geometry/lod/v2.recipe.json`;
- `internal/geometry/silhouette.go`;
- `internal/geometry/silhouette_test.go`;
- `internal/geometry/testdata/silhouette-oracle/**`.

Everything else is read-only, including `.git`, `.mate`, `docs`, Go module
files, catalog/data, serializer, runtime LOD, v1 recipe, oracle JSON, group
anchors, presets/overrides, T0B/T1/T2 files. Preserve shared scratch; no staging,
commit, lifecycle mutation or `node_modules` output.

## Context and doctrine pack

CTX-RUN-020 supplies P1 authority. PLAN-011 is the exact HOW.

RUN-017 proved current Mapshaper candidate generation turns SH's four projected
components into one four-point Polygon at every frozen resolution. Parser,
projection and oracle are correct. The fix changes batching and candidate-source
identity only; Mapshaper 0.7.44, command flags, resolutions, weighting,
projection, candidate order, thresholds, metrics, q and budgets remain exact.

Contribution-ranked components below the threshold are optional and omitted as
`protected_subscale`; do not accidentally make every identity-ranked island
mandatory.

## Tasks

1. Reproduce pre-run SH evidence and exact frozen-file/scratch inventory.
2. Extend v2 recipe schema with the versioned candidate-source declaration.
   Its new digest may change only because of that declaration.
3. Encode every projected source polygon as a lineage-keyed temporary Polygon
   feature: geometry ID, source index, source ring count and canonical projected
   digest. Submit all components of one geometry together to one pinned
   Mapshaper command so resolution remains country-relative.
4. Parse output by lineage key, not feature order. Require exactly one non-empty
   output per input before policy pruning. Validate exterior/hole lineage and
   fail typed on missing, duplicate, unknown, mutated, cross-geometry, inverted,
   escaped or invented data.
5. Resolve all entity/profile consumers of each geometry. Per band, form one
   generic mandatory union of components visible, dominant or anchored for any
   consumer. An identity-ranked subscale component remains optional with
   `protected_subscale`; other optional components require all-consumer
   `subscale`. Reassemble in source order with exact omission provenance.
6. Add teeth for output reorder, consumer iteration reorder, false provenance,
   mandatory pruning, optional retention until budget overflow and all lineage/
   hole failure classes. No entity/profile/preset/frame runtime branch.
7. Run the micro-brake first:
   - SH both consumers at all 29 resolutions;
   - CA, ID and KI compact/standard worst-multipart cases;
   - two identical builds;
   - at least one candidate per band passes every consumer, topology and exact
     `2200/7500` path plus `2500/8000` complete limits.
   Stop without representative rebuild if it fails.
8. On micro-brake PASS, rerun Indonesia 120/121 and full T0A focused tests.
   Rebuild representative artifacts twice under the unchanged oracle and
   thresholds; require machine PASS and byte identity.
9. Return exact diff, commands, timings, new recipe/manifest/sheet/receipt
   hashes, micro-brake maxima, representative rows and contact-sheet path. Stop
   for coordinator verification and owner review.

## Validation and evidence

- two runs of
  `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry/cmd/lodbuild -run
  'Test(SHCandidateLineageDiagnostic|ComponentLineage|MultipartMicroBrake)'
  -count=1 -v`;
- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry/cmd/lodbuild -count=1`;
- focused `./internal/geometry` silhouette/component/topology/protection tests;
- two representative builds with byte-identical new manifest/sheet/receipt;
- exact old/new recipe semantic diff and frozen oracle/serializer/hash inventory;
- compile-only all geometry packages.

Owner approval is not a substitute for any machine predicate.

## Constraints and forbidden actions

No threshold, metric, rasterizer, visibility rule, band boundary, resolution,
Mapshaper version/flag/algorithm/weighting, candidate order, projection, q,
serializer, budget, runtime dependency or selection-input change. No country
branch, manual vertices, MakeValid, union/buffer repair, CGO, post-hoc SVG
optimization or false omission provenance.

Node/Mapshaper remains maintainer-only. Normal runtime stays pure Go/offline.

## Stop escalation and resume conditions

Stop after one bounded attempt if exact lineage cannot be preserved; SH remains
lost; CA/ID/KI mandatory unions cannot fit existing budgets; hole matching is
ambiguous; any policy value must change; oracle or serializer drifts; output is
nondeterministic; representative machine gate fails; or a writable path outside
the lease is required.

On PASS, return machine evidence and the new sheet only. The coordinator must
independently verify, record owner approval and close the checkpoint before a
separate T0B run.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
