---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-001"
epic: "country-map-svg-generator"
spec: "P1"
run: "RUN-001"
status: running
profiles: []
concerns: []
inputs: ["PLAN-001", "CTX-RUN-001"]
---
# RUN-001 — Run brief

## Authority and accepted inputs

- Accepted plan revision: `PLAN-001` at SHA-256
  `4e39f29252ac88caf43144e625c1cbff18111720f7cb1a3fffde2d3e7ed82a27`.
- Accepted plan manifest: `PLAN-001` at SHA-256
  `4cd93236c76fbb76ecbff1f83dff79070393718a691925a3f8c5bd2fcf774a6e`.
- Run context pack: `CTX-RUN-001` at SHA-256
  `57a6a6f4e08fec174bf0c3bad4f47121c13cb4d63c8757f5f59a6647f48a82fd`.
- Baseline inherited from the accepted plan: repository
  `country-map-svg-generator`, commit
  `c7116841e99ae099b75b3c9eabdfef1805c7543f`, tree
  `7749f5e628feecd38dff900c621e2eea095f0376`.
- DEC-002 and DEC-003 remain binding. Run readiness is intentionally not embedded
  in this candidate; `mate run start` binds it atomically.

## Objective and completion boundary

Implement and verify CTR-1: an immutable deterministic local corpus containing
exactly the current 249 ISO 3166-1 assigned entities, both `un` and `de-facto`
geometry profile resolutions, source receipts, capital roles/overrides, and
protected-feature metadata. Finish only when the committed corpus and embedded
bundle decode to the same identity, all VAL-1–5 mutation teeth pass, and
`make check` is green without network access.

This run does not implement projection, simplification, SVG rendering, CLI product
commands, site delivery, or public release automation.

## Scope ownership and footprint

Writable product paths are exactly the union accepted by PLAN-001:

- `go.mod`, `go.sum`, `Makefile`;
- `data/sources/**`, `data/policy/**`, `data/corpus/v1/**`,
  `data/reports/v1/**`;
- `internal/catalog/**`, including its maintainer compiler, focused tests, and
  generated embedded bundle.

Forbidden product roots: `internal/geometry/**`, `internal/config/**`,
`internal/render/**`, `cmd/**`, `docs/guides/**`, `docs/operations/**`,
`.goreleaser.yaml`, `.git/**`, and `.mate/**`.

The integration target is the current `main` line. The implementer returns one
reviewable atomic product commit and may not commit lifecycle evidence.

## Context and doctrine pack

`CTX-RUN-001` contains the current P1 specification. The accepted plan separately
binds the current solution architecture and source pins: ISO assigned-set snapshot
observed `2026-07-23`, Natural Earth admin-0 and ISO POV `5.1.1`, and populated
places simple `5.1.2`.

The repository Go profile is active. Build and black-box checks use `GOWORK=off`;
tests are hermetic and must not depend on network, user home, locale, machine
geodata, Node, Python, GDAL, or an ambient workspace.

## Tasks

1. T1: establish the Go module and strict source/catalog models; acquire and pin
   exact reviewed source bytes and receipts; prove the 249-entity ISO invariant.
2. T2: implement explicit identifier reconciliation, GeoJSON decoding, both
   boundary profiles, capital roles/overrides, protected anchors/parts, and
   fail-closed validation.
3. T3: implement canonical deterministic encoding, content identity, comparison,
   create-once transactional publication, decoder, corpus v1, and embedded bundle.
4. T4: add the full-corpus and mutation teeth for VAL-1–5, wire `gen-check` and
   `check`, and run the complete project ceiling.

Execute serially. Each task consumes the preceding contract and cannot weaken it.

## Validation and evidence

- `corpus-iso-integrity`: missing and duplicate entity mutations fail; the real
  manifest has exactly 249 unique uppercase alpha-2 keys.
- `corpus-boundary-profiles`: removed-profile, empty/dangling geometry, and selected
  dispute mutations fail.
- `corpus-reproducibility`: two clean compilations and embedded/committed decoding
  produce identical bytes and identities; missing receipts fail.
- `corpus-capital-registry`: zero/single/multiple/override cases round-trip; invalid
  coordinate and primary references fail.
- `corpus-protected-features`: named microstate/island coverage passes and removed
  coverage fails with the entity/feature.
- Focused commands:
  `GOWORK=off go test ./internal/catalog/... -count=1`,
  `make gen-check`, and `make check`.

The run result must name the exact commands, outcomes, changed paths, commit, any
deviations, and feedback handles. Green narrow tests do not replace `make check`.

## Constraints and forbidden actions

- Do not infer ISO identity from names, synthesize missing polygons, silently drop
  an entity, or collapse the political meaning of either profile.
- Do not hand-edit compiled corpus or generated embedded bytes; change the source,
  mapping, policy, or compiler and regenerate.
- Do not add runtime downloads or runtime dependencies on source files, Node,
  Python, GDAL, or system geodata.
- Do not mix presentation, projection, simplification, SVG, CSS, or product CLI
  policy into `internal/catalog`.
- Do not change accepted source versions, catalog authority, profile semantics,
  corpus identity, ownership boundaries, or P1/P2 interface inside the run.
- Do not weaken a validator or mutation fixture to obtain green.

## Stop escalation and resume conditions

Local helper/file splits and diagnostic wording inside the accepted footprint are
allowed and must be reported. Stop before mutation for any source-version change,
political override, entity-count mismatch that requires changing authority,
missing real geometry that would require synthesis, new write root, runtime
acquisition, or changed P1/P2 contract. Route WHAT or macro-HOW changes through the
teamlead and amendment flow.

After a few speculative fixes for the same failure, reduce to a minimal fixture and
read the owning library/source contract. On interruption, resume from the accepted
plan, this run brief, current product commit/diff, and last exact validation output;
never reconstruct authority from chat history.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
