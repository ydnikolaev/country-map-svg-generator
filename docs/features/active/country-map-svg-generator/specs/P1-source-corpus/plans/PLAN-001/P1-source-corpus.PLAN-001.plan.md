---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "accepted-plan-revision"
id: "PLAN-001"
epic: "country-map-svg-generator"
spec: "P1"
status: accepted
profiles: []
concerns: []
inputs: ["P1"]
---
# PLAN-001 — Implementation plan

## Accepted inputs and baseline

- Normative specification: `P1-source-corpus.spec.md` at SHA-256
  `a2dd6817331c8da19775f66de49b7aa4526651f06292f8254606c3b3fe3ea3bd`.
- Epic architecture: `country-map-svg-generator.solution-architecture.md` at
  SHA-256
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`.
- Accepted decisions: DEC-002 fixes `un` as the default boundary posture with
  configurable `de-facto`; DEC-003 fixes a Go runtime over an immutable local
  corpus.
- Repository baseline: commit `c7116841e99ae099b75b3c9eabdfef1805c7543f`,
  tree `7749f5e628feecd38dff900c621e2eea095f0376`.
- Source pins for the first corpus are fixed before execution:

  - ISO 3166-1 assigned-set snapshot observed `2026-07-23`; exactly 249 assigned
    alpha-2 entities, with alpha-3 and English display names recorded in a
    reviewed local snapshot. ISO OBP is the authority; the acquisition receipt
    records the observation date and snapshot digest.
  - Natural Earth `ne_10m_admin_0_countries_iso` version `5.1.1` supplies the
    `un` profile.
  - Natural Earth `ne_10m_admin_0_countries` version `5.1.1` supplies the
    `de-facto` profile.
  - Natural Earth `ne_10m_populated_places_simple` version `5.1.2` supplies the
    capital candidates.

  All acquired source bytes are committed below `data/sources/` with receipts;
  normal generation never downloads them.
- P1 exclusively owns BND-001 (`internal/catalog/**`) and BND-002 (`data/**`).
  The greenfield module roots `go.mod`, `go.sum`, and `Makefile` are established
  here as shared foundations; later specs extend but do not replace them.
- No accepted amendment or unresolved authority changes this plan.

## Requirement and acceptance coverage

| Contract | Tasks | Required evidence |
| --- | --- | --- |
| REQ-1, AC-1, ARCH-INV-1 | T1, T2, T4 | Exactly 249 sorted unique uppercase alpha-2 records, each with alpha-3 and display name; missing and duplicate mutations fail with entity diagnostics. |
| REQ-2, AC-2 | T2, T4 | Every entity resolves `un` and `de-facto` through an explicit geometry object or an explicit identical reference; removed-profile and disputed-fixture mutations fail closed. |
| REQ-3 | T1, T2 | Every source input is named by a receipt containing origin, upstream version/date, license, acquisition SHA-256, and transformation identity; undocumented input fails. |
| REQ-4, ARCH-INV-7 | T2, T4 | Zero/single/multi-capital records, roles, primary choice, and explicit overrides round-trip; invalid coordinates and dangling primary IDs fail. |
| REQ-5 | T2, T4 | Protected anchor/part metadata covers named island and microstate fixtures and is rejected when required coverage is removed. |
| REQ-6 | T2, T4 | Missing/duplicate ISO, empty geometry, invalid coordinate, dangling reference, missing receipt, and uncovered protected-feature fixtures all fail non-zero with source/entity/field/invariant context. |
| REQ-7, AC-3 | T3, T4 | Two clean compiles emit byte-identical manifest, geometry, embedded bundle, report, and content digests. |
| INV-1 | T3 | Corpus version directories are create-once/content-addressed; changed bytes require a new identity and never overwrite v1. |
| AC-4 | T2, T4 | Representative normal, multi-capital, island, archipelago, and microstate fixtures produce a registry coverage report. |
| VAL-1–5 | T4 | Focused unit/integration mutation tests and the registered project ceiling provide the exact evidence slots declared in the manifest. |

## Technical approach

Create a small standard-library-first Go package under `internal/catalog` with
closed models for entities, geographic polygon/multipolygon geometry, boundary
profile references, capital roles, protected anchors, receipts, and corpus
identity. Coordinates remain WGS84 longitude/latitude. The runtime-facing decoder
has no acquisition logic and returns immutable values keyed by uppercase alpha-2.

The maintainer compiler lives under `internal/catalog/cmd/corpuscompile`. It reads
only explicit files below a supplied source root: the reviewed ISO snapshot, the
two pinned Natural Earth GeoJSON themes, capital GeoJSON, mapping/override records,
protected-feature policy, and receipts. The compiler never infers ISO identity
from a display name. Natural Earth identifiers are joined only through reviewed
mapping fields; countries absent from an upstream theme require a documented
explicit mapping or explicit geometry source, never synthesized geometry.

The `un` and `de-facto` views are separate manifest contracts. When their canonical
geometry bytes are equal the second profile stores an explicit `identical_to`
reference; otherwise each stores its own geometry ID. Disputed fixtures lock the
owner-approved expectations without embedding presentation policy.

Capital ingestion keeps every admin-0 capital candidate with a stable local ID,
role, name, coordinates, and primary flag. A reviewed override file can replace
or augment source roles and primary choice, including countries with zero or
multiple capitals. Protected-feature policy records named anchor points and
minimum retained part counts for representative microstates, islands, and
archipelagos; P2 consumes these constraints rather than guessing from area.

Compilation canonicalizes strings and floating-point precision, sorts every map
as a slice before encoding, validates all cross-references, and writes into a
temporary directory. The accepted output is deterministic JSON plus a
deterministically compressed embedded bundle generated beneath
`internal/catalog`; compression metadata is fixed. The corpus identity is SHA-256
over canonical manifest and geometry bytes. Publication renames the complete
temporary version only when validation passes and refuses to overwrite an existing
version with different bytes.

## Tasks and completion conditions

1. **T1 — Pin source set and define the catalog contract.** Establish the Go
   module, commit the reviewed 249-entity ISO snapshot, Natural
   Earth 5.1.1/5.1.2 source bytes, license text, receipts, mappings, boundary
   policy, capital overrides, and protected-feature policy. Implement strict
   source/model decoding and receipt verification. Done when every acquisition
   digest matches committed bytes, the ISO fixture is exactly 249 unique sorted
   records, all source inputs have receipts, and malformed source/model fixtures
   fail with typed path/entity/field diagnostics.
2. **T2 — Compile and reconcile geometry, profiles, capitals, and protections.**
   Implement the maintainer-only GeoJSON adapters and compiler. Join only reviewed
   identifiers, normalize geographic geometry deterministically, resolve both
   boundary profiles, apply capital overrides, and bind protected anchors/part
   counts. Done when all 249 entities have non-empty resolvable geometry in both
   profiles, capital fixtures cover zero/single/multiple/override cases, protected
   fixtures cover microstates/islands, and all REQ-6 mutations fail closed.
3. **T3 — Publish immutable deterministic corpus v1.** Implement canonical
   encoding, content identity, comparison report, transactional create-once
   publication, decoder, and generated embedded bundle. Compile v1 twice into
   separate temporary roots and compare every byte/digest. Done when the committed
   `data/corpus/v1` manifest, geometry, coverage, receipts index, and embedded Go
   bundle decode to the same identity; overwrite-with-different-bytes fails; and a
   refresh helper produces added/removed/changed output without mutating v1.
4. **T4 — Add validation teeth and project ceiling.** Add unit, integration, and
   mutation fixtures for VAL-1–5; exercise the real committed source set and
   generated bundle through `go test`; wire deterministic generate/check targets
   into `Makefile`. Done when every named sabotage reds for the expected invariant,
   the unmodified full corpus passes, repeated compilation remains byte-identical,
   and `make check` is green with no network access.

## Ownership and write footprint

- T1 (`catalog-contract`) owns:
  `go.mod`, `go.sum`,
  `internal/catalog/{model.go,model_test.go,source.go,source_test.go}`,
  `internal/catalog/testdata/source/**`, and
  `data/sources/**`.
- T2 (`catalog-compiler`) owns:
  `internal/catalog/{compile.go,compile_test.go,geojson.go,geojson_test.go,validate.go,validate_test.go}`,
  `internal/catalog/cmd/corpuscompile/{main.go,main_test.go}`, and reviewed
  records below `data/policy/**` plus `internal/catalog/testdata/compile/**`.
- T3 (`corpus-publisher`) owns:
  `internal/catalog/{canonical.go,canonical_test.go,compare.go,compare_test.go,embedded.go,embedded_gen.go}`,
  `data/corpus/v1/**`, and `data/reports/v1/**`.
- T4 (`corpus-verifier`) owns focused additions under
  `internal/catalog/{fullcatalog_test.go,mutations_test.go}`,
  `internal/catalog/testdata/mutations/**`, and `Makefile`.
- Generated lifecycle, plan, run, audit, readiness, and evidence artifacts remain
  coordinator-owned. P1 does not write `internal/geometry/**`, `internal/config/**`,
  `internal/render/**`, `cmd/**`, `docs/guides/**`, `.goreleaser.yaml`, or
  `docs/operations/**`.

## Dependencies and execution waves

- **W1:** T1, concurrency 1. Freeze model and exact source bytes.
- **W2:** T2, concurrency 1 after T1. Freeze reconciliation semantics.
- **W3:** T3, concurrency 1 after T2. Publish the immutable corpus.
- **W4:** T4, concurrency 1 after T3. Prove mutation teeth and project ceiling.

The spec is intentionally serial. Source/model identity flows into every later
byte, and concurrent writers to the compiler or generated corpus would weaken the
reproducibility claim.

## Validation plan

- T1: `GOWORK=off go test ./internal/catalog -run
  'Test(Source|Model|ISO|Receipt)' -count=1`; missing/duplicate ISO, malformed
  alpha codes, wrong acquisition digest, duplicate source ID, and undocumented
  input mutations must fail.
- T2: `GOWORK=off go test ./internal/catalog -run
  'Test(Compile|Profiles|Capitals|Protected|Validation)' -count=1`; removed profile,
  invalid/dangling geometry, invalid coordinates, invalid primary capital,
  missing protected anchor, and display-name-only joins must fail.
- T3: run the built compiler twice with distinct temporary output roots and compare
  recursive SHA-256 manifests; decode the committed and embedded bundles and assert
  the same corpus identity. Attempting to republish changed v1 bytes must fail
  without altering the prior directory.
- T4: `GOWORK=off go test ./... -count=1`, then `make gen-check` and `make check`.
  Tests use committed fixtures and must not consult the network, user home, locale,
  or system geodata.
- VAL evidence mapping:

  - VAL-1 → `corpus-iso-integrity`, integration.
  - VAL-2 → `corpus-boundary-profiles`, integration.
  - VAL-3 → `corpus-reproducibility`, integration.
  - VAL-4 → `corpus-capital-registry`, component+integration.
  - VAL-5 → `corpus-protected-features`, integration.

## Deviation and amendment policy

`p1-corpus-local-how-v1`: helper naming, file splits inside declared roots,
canonical encoder implementation, fixture layout, and diagnostic wording are local
HOW. A new path inside BND-001/BND-002 or a different deterministic encoding with
the same contract requires a reviewed plan revision before mutation.

Stop and return to the teamlead for any change to the 249-entity catalog authority,
source versions, boundary-profile meaning/default, political override, corpus
identity rules, protected-feature semantics, capital-role contract, runtime
download policy, ownership boundary, or P1/P2 interface. Those are WHAT or
macro-HOW and require an accepted amendment/decision. Do not silently fill missing
geometry, drop an ISO entity, collapse profile meaning, or move acquisition into
normal generation.

## Commit worktree and integration policy

Execute T1–T4 serially from the accepted baseline. Each task produces one
conventional atomic commit containing only its declared lease; committed source
receipts travel with source bytes, and generated corpus bytes travel with the
compiler change that produced them. The coordinator verifies the actual diff and
parent, runs the focused task gate, and runs the combined gate before the next
wave. Implementers do not commit lifecycle artifacts or stage unrelated files.

## Rollback and recovery

Compilation stages into a new temporary directory and publishes only after all
validators pass. On failure, discard only that explicit temporary directory and
leave the accepted corpus untouched. An interrupted task resumes from its atomic
commit and the committed source receipts. Rollback selects the previous immutable
corpus identity; it never edits v1 in place.

A changed source digest, ISO snapshot, accepted decision, specification,
architecture, or baseline invalidates the plan before redispatch. Never recover by
hand-editing compiled bytes, bypassing receipt checks, weakening mutation fixtures,
or synthesizing a placeholder polygon.

## Completion and handoff

P1 completes only with an accepted PLAN-001, one governed implementation run,
exactly 249 ISO entities, both boundary profiles per entity, deterministic corpus
v1 and embedded bundle, valid source receipts, capital/protected-feature coverage,
all VAL-1–5 evidence, one fresh independent audit, and green `make check`.

The P2 handoff is CTR-1: exact corpus identity, decoder API, geographic geometry,
profile resolution, capital roles/primary selection, protected anchors/minimum
parts, coverage report, and rollback identity. P2 may transform geometry but may
not reinterpret ISO identity, political profile semantics, or source provenance.

Estimated effort is 6–10 agent-hours, likely 8, medium confidence. Source
reconciliation for small territories and multi-capital overrides is the bounded
uncertainty; it must be resolved with explicit reviewed records rather than hidden
compiler heuristics.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
