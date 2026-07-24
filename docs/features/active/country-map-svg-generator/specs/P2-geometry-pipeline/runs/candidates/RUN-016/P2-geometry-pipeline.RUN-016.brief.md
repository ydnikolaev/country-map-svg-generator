---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-016"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-016"
status: draft
profiles: []
concerns: []
inputs: ["PLAN-010", "CTX-RUN-018"]
---
# RUN-016 — Full-catalog v2 candidate source and runtime integration

## Authority and accepted inputs

- Accepted PLAN-010:
  `a81da84fc3881068ad0520a3746f3276c18df87a48f10e83d5463fb03fb1db63`;
  manifest
  `fdd0cdfa0c545410f102eb460fef80b57d3b03484a2fc0b1cb423975a7fff59d`.
- Runtime context CTX-RUN-018:
  `907ca2d90c1cf4cd91bfb1897a9160511c4ee8621d93ffda3d643d3935789df8`.
- RUN-015 interrupted result:
  `e990b0db24fe2ce78c7b3ecab1cd81253ef53519a211e75b51dbbf53cba746b4`;
  coder result
  `fcdfa41a754e2fb84d1dd61cfb935d3d4a07995bd7499e2f81af148649b7b5b8`.
- RUN-013 T0A gate
  `792fd75998656c10ecbd7bfdea760685c8f5880d613574b50210aa421339a0ba`
  and accepted owner decision DEC-008
  `d92caff0294b08a46cb7d49dfda8e2309918afe0665a6a01d6a82fa55894d3de`.
- Frozen oracle
  `f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`
  and v2 recipe
  `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`.

## Objective and completion boundary

Return to PLAN-010 T0A only far enough to supply the missing full-catalog
candidate seam proven by T0B. Build a deterministic temporary all-283 v2
compact/standard ladder from the frozen recipe and oracle, consume it in
automatic `GenerateWithLOD` selection with requested-band budgets and
rotation-safe protected-visibility finalization, and make the inherited
full-corpus and AQ/RU checkpoints green.

Rerun the dependent T0A checks and reproduce the five owner-approved
representative artifacts byte-for-byte. Stop at this integration checkpoint.
Do not publish T1 ladder/manifest files, render the owner full-catalog contact
sheet, land the combined T0 product commit, or execute T2.

## Scope ownership and footprint

Writable:

- `internal/geometry/lod.go`;
- `internal/geometry/lod_test.go`;
- `internal/geometry/cmd/lodbuild/main.go`;
- `internal/geometry/cmd/lodbuild/main_test.go`;
- new implementation and test files directly under
  `internal/geometry/cmd/lodbuild/`.

Everything else is read-only, including `.git`, `.mate`, `docs`, `go.mod`,
`go.sum`, `internal/catalog/**`, the oracle, both recipes, group-anchor data,
serializer, representative generator, its approved outputs, all remaining T0B
files and every T1/T2 publication file. Preserve shared scratch. Return without
staging, commit, lifecycle mutation or generated `node_modules`.

## Context and doctrine pack

CTX-RUN-018 supplies the current P1 specification; PLAN-010 and the current Go
profile remain binding. RUN-014 and RUN-015 are the exact integration evidence.

The blocker is exact: the frozen v1 records cannot meet the approved compact
policy for all consumers (`CA/card compact`: path 2286 > 2200 and complete
estimate 2566 > 2500), so runtime changes alone cannot pass. A fallback record
does not inherit its record name's looser policy: the requested/effective-scale
band governs every candidate considered for that request. The maintainer-only
`EvaluateSilhouetteCandidate` uses zero rotation, while runtime geometry,
overrides and markers must remain aligned under the caller's rotation.

Keep global thresholds compact/standard `110/900`, IoU/recall `.40/.42`,
scale boundaries `240/700`, path caps `2200/7500`, complete-file maxima
`2500/8000`, estimate envelope `280`, q `.01`, candidate order, projection,
serializer and empty group-anchor dataset exact.

## Tasks

1. Add one deterministic maintainer builder that executes each frozen v2
   resolution as an all-283 batch, reusing the canonical mapshaper command and
   recipe identity rather than invoking Node per entity.
2. Resolve all `249 × {un,de_facto}` catalog consumers, including
   `identical_to`, for every geometry. For compact and standard, select the
   first fine-to-coarse record that passes every consumer of that geometry
   under that band: group anchors, IoU/recall, dominant and anchor coverage,
   visible identity retention or truthful `protected_subscale` omission,
   topology, path cap and `path+280 < CompleteFileMaximum`.
3. Emit an in-memory temporary `LODTable` plus deterministic candidate and
   provenance evidence sufficient for tests. Do not write
   `v2.ladder.json`/`v2.manifest.json`; T1 owns publication.
4. Integrate automatic table-backed runtime selection in `lod.go`. Determine
   policy from requested/effective scale, bind each projected record with the
   runtime projector and rotation, add reviewed group anchors, evaluate
   protected visibility and topology in that aligned space, and finalize via
   the existing deterministic grid-phase path. Do not adopt the transform or
   geometry returned by the zero-rotation maintainer evaluator.
5. For automatic selection, remove unconditional P1 minimum-part restoration
   and DEC-005 raw-deviation rejection. Keep phase protection active with
   minimum one component, revalidate oracle visibility/provenance after
   finalization, and populate existing `LODProvenance` fields truthfully.
6. Preserve explicit source/custom and nil-table behavior exactly. Caller
   `MaxPathBytes` remains a terminal render-only constraint and cannot alter
   candidate selection. Entity, profile, preset name, contain-frame shape,
   marker/style inputs and requested byte cap cannot branch candidate identity.
7. Update only test setup that currently injects frozen v1 records so the
   inherited production assertions consume the temporary v2 ladder. Add
   mutation teeth for all-consumer geometry selection, requested-band fallback,
   rotation/marker alignment, protected visibility loss, false omission
   provenance, group anchors, deterministic rebuild and unchanged source/custom
   behavior. Do not weaken corpus count, budget, determinism, tier, projection,
   grid or runtime assertions.
8. Run the two inherited production checkpoints, full geometry/lodbuild unit
   suites, exact RUN-013 focused gate and two representative rebuilds. Require
   the same oracle, recipe, representative manifest, contact sheet and machine
   receipt hashes. Return exact diff, commands, hashes, maxima, build/runtime
   timings and any residual failure.

## Validation and evidence

- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry/cmd/lodbuild -run
  'Test(LODSpikeFullCorpus|LODAQRUProjectionAlignedCheckpoint)' -count=1`;
- `GOCACHE=/private/tmp/country-map-go-cache GOWORK=off go test
  ./internal/geometry ./internal/geometry/cmd/lodbuild -count=1`;
- the exact RUN-013 focused T0A gate;
- two full-catalog temporary v2 builds with identical table/provenance digests;
- two representative builds with exact frozen identities:
  oracle `f2c9cd32...`, recipe `da3ff9d3...`,
  manifest `b2c467b0...`, sheet `cf1e568c...`,
  machine receipt `fcc1bc30...`;
- all 996 production outputs deterministic and within compact/standard requested
  band limits, with logged maxima and warm runtime no more than 120 seconds;
- compile-only geometry/lodbuild and exact frozen-file hash inventory.

This is a dependent T0A integration checkpoint. No failure can be waived by
DEC-008, and exact representative identity is required to reuse that approval.

## Constraints and forbidden actions

No change to oracle, recipe, thresholds, contribution model, identity ranking,
group data, candidate order, projection, q, budgets, serializer or public
presets. No country/profile/preset/frame-specific candidate branch, source
metadata stripping, vertex reinsertion, hand repair, MakeValid, union/buffer,
post-hoc SVG optimization, lower precision, higher budget or runtime external
dependency.

Mapshaper/Node is maintainer-build-only. Normal runtime remains pure Go and
offline. Do not publish generated ladder files, broaden into the full owner
sheet/T0B closeout, stage shared scratch, or commit.

## Stop escalation and resume conditions

Stop after one bounded implementation attempt if no frozen v2 resolution can
satisfy all consumers of a geometry; selection requires changing an approved
threshold/recipe/oracle; correct runtime integration needs a file outside the
lease; explicit source/custom semantics drift; any representative hash changes;
the temporary ladder is nondeterministic; a runtime rotation cannot preserve
geometry/marker alignment; or either inherited checkpoint remains red.

On PASS, return exact content evidence so the coordinator can independently
verify and close the T0A checkpoint, then resume T0B in a separate run. On a
frozen identity change, DEC-008 is invalid and owner approval must be repeated.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
