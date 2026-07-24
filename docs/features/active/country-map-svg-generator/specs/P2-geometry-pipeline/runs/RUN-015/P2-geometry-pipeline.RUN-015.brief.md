---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-015"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-015"
status: running
profiles: []
concerns: []
inputs: ["PLAN-010", "CTX-RUN-017"]
---
# RUN-015 — Frozen protected-visibility runtime integration

## Authority and accepted inputs

- Accepted PLAN-010:
  `a81da84fc3881068ad0520a3746f3276c18df87a48f10e83d5463fb03fb1db63`;
  manifest
  `fdd0cdfa0c545410f102eb460fef80b57d3b03484a2fc0b1cb423975a7fff59d`.
- Runtime context CTX-RUN-017:
  `dcbabf6eea957560dc8e048265d070a68caa6780d460605d8cc4e405f18e6bd0`.
- RUN-014 blocked result:
  `65f91c9a25dfae47a8c89bf54be89aafe80044fe85999712fb40ecba6496dfdf`;
  coder result
  `55445a767c41fe0dab3288a3f7e15c28cc559750a88a55d801e4eb917ba2de7d`.
- RUN-013 T0A gate
  `792fd75998656c10ecbd7bfdea760685c8f5880d613574b50210aa421339a0ba`
  and accepted owner decision DEC-008
  `d92caff0294b08a46cb7d49dfda8e2309918afe0665a6a01d6a82fa55894d3de`.
- Frozen `lod.go` starts at
  `5d03ca1ff5e5cab55f3c9d77a608c102a5be9666f11e5c788d358f131a3b1006`;
  every other T0A and serializer hash remains the RUN-013 inventory.

## Objective and completion boundary

Close only the lease gap proven by RUN-014. Integrate the already frozen
`protected-visibility/v1` oracle into the automatic `GenerateWithLOD` candidate
path, using the exact approved band thresholds and budgets. Make the two
inherited production tests green without changing their assertions.

Rerun all dependent T0A checks and reproduce the five representative artifact
hashes exactly. Stop after this runtime-integration checkpoint. Do not continue
the broader 996-output T0B evidence/sheet, stage, commit, publish T1 or execute
T2.

## Scope ownership and footprint

Writable:

- `internal/geometry/lod.go`;
- `internal/geometry/lod_test.go`.

Everything else is read-only, including `.git`, `.mate`, `docs`, `go.mod`,
`go.sum`, `internal/catalog/**`, `cmd/lodbuild/**`, `silhouette.go` and its
oracle/recipe/testdata, group-anchor code/data, serializer, every T0B file and
all T1/T2 files. Preserve shared scratch. Return without staging, commit or
lifecycle mutation.

## Context and doctrine pack

CTX-RUN-017 supplies the current P1 and Go doctrine. PLAN-010, DEC-007,
DEC-008, RUN-014 and this brief are the exact HOW.

The blocker is exact: direct `GenerateWithLOD` callers cannot be intercepted by
T0B. The old automatic branch computes unconditional P1 `minimumParts`,
restores those components into every candidate, and rejects candidates by
DEC-005 raw deviation. That is correct only for explicit source/custom quality;
it contradicts the approved automatic protected-visibility policy.

Keep oracle `f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`,
recipe `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`,
thresholds `110/900`, IoU/recall `.40/.42`, scale boundaries `240/700`,
path caps `2200/7500`, complete-file maxima `2500/8000`, estimate envelope
`280`, q `.01`, projection, candidate order, serializer and empty group-anchor
data exact.

## Tasks

1. Reproduce the exact RUN-014 failures:
   AE/un/card source fallback at `5910 > 2500`; AQ/un/card compact rejection at
   raw deviation `1.537690736882369 > 1.536`.
2. Separate automatic table-backed selection from explicit source/custom
   finalization. Source requests, a nil table and non-auto quality must retain
   the existing DEC-005 minimum-part, deviation, q-aware source and typed-budget
   behavior exactly.
3. For an automatic table-backed candidate, resolve its exact compact/standard
   band from the embedded oracle. Apply the embedded reviewed group-anchor set
   as additive protected data, then evaluate the candidate with the same
   projection, rasterizer, visibility, topology, fixed-grid and band policy as
   the frozen representative path.
4. Automatic acceptance must require the frozen global IoU/recall pair,
   dominant and anchor coverage, every visible identity component, valid
   `protected_subscale` provenance, topology, `path<=band.PathCap` and
   `path+280<band.CompleteFileMaximum`. It must not consult DEC-005 raw
   deviation or unconditionally reinsert omitted source components.
5. Preserve runtime candidate/fallback semantics and scale-only selection.
   Entity, profile, preset name, contain-frame shape, marker/style inputs and
   caller `MaxPathBytes` cannot change candidate acceptance. The caller hard
   byte limit remains a terminal render error and may not advance tiers.
6. Record protected-visibility metrics, band identity, retained/omitted
   component provenance and group-anchor application in existing
   `LODProvenance` without changing public presentation policy. Fail closed on
   oracle/group data errors, lineage loss or false omission provenance.
7. Add focused mutation teeth for automatic raw-deviation bypass, visible
   protected loss, dominant/anchor loss, threshold/band mismatch, path/file
   overflow, non-additive group behavior, entity/preset/frame branching and
   unchanged source/custom behavior.
8. Run both inherited production tests, the full lod unit suite and the exact
   RUN-013 T0A focused gate. Rebuild representative artifacts twice and require
   the same five hashes and Indonesia 120/121 boundary.
9. Return exact diff, commands, results, artifact hashes and both production
   maxima. Stop; T0B resumes only in a separate run.

## Validation and evidence

- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry/cmd/lodbuild -run
  'Test(LODSpikeFullCorpus|LODAQRUProjectionAlignedCheckpoint)' -count=1`;
- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry -run
  'Test(Diagnostic|GridPhase|Projection|LOD|Topology|Protection|Preset|Silhouette|Oracle|Selection|Source|GroupAnchor)'
  -count=1`;
- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry/cmd/lodbuild -run '^TestDiagnostic' -count=1`;
- two representative builds with exact oracle/recipe/manifest/sheet/receipt
  hashes and exact Indonesia boundary;
- compile-only geometry/lodbuild and complete frozen-file hash inventory.

This is a T0A-dependent integration checkpoint, not the PLAN-010 project
ceiling. No machine failure can be waived by DEC-008.

## Constraints and forbidden actions

No change to oracle, recipe, thresholds, contribution model, identity ranking,
group data, candidate order, projection, q, budgets, serializer or public
presets. No country/profile/preset/frame branch, source metadata stripping,
vertex reinsertion, hand repair, MakeValid, union/buffer, post-hoc SVG
optimization, lower precision, higher budget or runtime external dependency.

Do not weaken or edit the inherited production tests. Do not make automatic
selection depend on the caller hard-byte value. Do not broaden into T0B,
publish the ladder, stage `node_modules`, stage Git paths or commit.

## Stop escalation and resume conditions

Stop if the integration requires any file beyond `lod.go`/`lod_test.go`; source
or custom behavior changes; a frozen artifact hash changes; the representative
sheet changes; a generic automatic candidate cannot pass the approved oracle
and budgets; group-anchor validation cannot remain additive; or either
inherited test remains red after one bounded implementation attempt.

On PASS, return exact content evidence so the coordinator can close this
checkpoint and start a fresh T0B run. If any frozen identity changes, DEC-008
does not authorize continuation and the owner gate must be repeated.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
