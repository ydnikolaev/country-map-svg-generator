---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "audit"
id: "AUD-002"
epic: "country-map-svg-generator"
spec: "P1"
status: completed
profiles: []
concerns: []
inputs: ["RUN-002-RESULT"]
---
# AUD-002 — Audit

## Subject scope and exclusions

Audited exact implementation commit
`ed4fdca23965b6aa0c0fc1b8e4c449e79f2d45ee`, parent
`bd23a1423ade0ab3686a89cf250ca37581f1543f`, tree
`1a6bfca82d26027956f92e28c9cd479c277faba5`, against accepted P1,
PLAN-001, DEC-002, DEC-003, ARCH-INV-1, ARCH-INV-7, sealed
AUD-001/ADV-001, RUN-002-RESULT, committed sources/corpus, and recorded
gates. Later coordinator-owned lifecycle commits were excluded from the
product candidate; no product bytes differ between the audited commit and
current HEAD.

## Inputs and evidence

Authority hashes match the durable dispatch: role contract
`ac38fb1d3936459fc1480d239a4d24279d5b600404b6a8d8a07e12f3ee6ed654`,
specification
`a2dd6817331c8da19775f66de49b7aa4526651f06292f8254606c3b3fe3ea3bd`,
plan manifest
`4cd93236c76fbb76ecbff1f83dff79070393718a691925a3f8c5bd2fcf774a6e`,
sealed audit
`da0f7b477e6d044ddb5cedf66b16c41ab8703bc9cfc6c10ccccabc7ba5033be6`,
advisor receipt
`d8e0df1d67b89a48924ffc7443d52bc62225dafd519ee22b8b5620618bc97b2a`,
and run result
`eb64507d5b9c108bb25492c8d2091c3335ca968e73485571783e62f302b99282`.

The commit changes 20 product paths, all inside the accepted P1 footprint. No
forbidden ownership root, source version, dependency, runtime-network,
political-profile, or P1/P2 contract change was found.

Independent hashing reproduced corpus identity
`sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`
from manifest, geometries, receipts, and coverage. The corpus directory has
exactly those four files. All five source receipt digests match committed
bytes, and embedded decoding byte-matches every corpus component.

The corpus contains 249 entities, 283 geometries, 249 explicit UN
resolutions, 34 explicit de-facto resolutions, 215 identical references,
212 capital records over 195 entities, and six protected entities.

## Method and independence

Fresh principal `codex-p1-auditor-2`; no implementation authorship,
remediation, product writes, lifecycle mutation, commit, or audit
materialization. Reviewed authoritative inputs, exact commit scope, source
pins, compiler/decoder/publication code, actual source and corpus data,
generated bundle hashes, focused gates, project ceiling, and obligation
evidence allocation.

Reran focused adversarial tests for every AUD-001 finding, full catalog tests,
`make gen-check`, and `make check`. Independently recomputed the whole-set
identity, source digests, oracle geometry hashes, protected-anchor
containment, reconciliation binding, capital cases, and embedded bytes.
Publication rollback was exercised by inducing real failures after the first
and after the second destination had been replaced.

## Findings

- **F-001 — closed.** Identity covers manifest, geometry, receipts, and
  coverage. Decode requires the exact canonical four-file set and rejects
  receipt/coverage mutations and unexpected files.
- **F-002 — closed.** CN, CY, IL, IN, and RU have exact reviewed UN and
  de-facto geometry identities. All ten referenced geometry IDs independently
  rehash from their coordinates, and every selected entity has mutation teeth.
- **F-003 — closed.** Indonesia's corrected `[110,-7.5]` anchor is contained
  in its 264-part UN geometry. All six protected anchors pass independent
  containment; ocean, wrong-entity, removal, and insufficient-part mutations
  fail.
- **F-004 — closed.** Capital policy supports zero-to-one addition,
  replacement, role assignment, additive selection, and primary-by-ID/name.
  Invalid coordinates, duplicate IDs/roles, and dangling primary references
  fail.
- **F-005 — remains open, high.** Focused gates exist, but almost all
  non-VAL obligation evidence still copies the broad reproducibility result.
  This does not preserve obligation-relevant assertions.
- **F-006 — closed.** A digest-bound reconciliation contains exactly 249
  unique ordered rows matching alpha-2, alpha-3, name, and `matched`
  disposition; missing, duplicate, mismatched, and unbound mutations fail.
- **F-007 — closed.** Corpus, report, and embedded output are staged and
  validated before replacement. Failures at the second and third destination
  restored every prior authoritative destination without a mixed output set.

### F-005-R2 — high, lifecycle evidence allocation

Files `VR-US-1..4-P1-RUN-002`, `VR-REQ-1..7-P1-RUN-002`,
`VR-AC-1..4-P1-RUN-002`, and `VR-INV-1-P1-RUN-002` all point to
`SELECT-59FD08536411` and `corpus-reproducibility`. Existing focused receipts
`GATE-EE3514765774` (ISO), `GATE-6F08DD3DA870` (profiles),
`GATE-B475A8CCACA3` (capitals), and `GATE-5D6A3BFF82E9` (protected
features) are allocated only to VAL-1, VAL-2, VAL-4, and VAL-5.

Owner: coordinator/gate evidence allocator. Rematerialize each obligation
result from its relevant focused receipt: ISO for US-1/REQ-1/AC-1; profiles
for US-2/REQ-2/AC-2; capitals and protected features for US-4/REQ-4/REQ-5/
AC-4; appropriate composites for REQ-3 and REQ-6; retain reproducibility for
REQ-7/AC-3/INV-1. No product-code change is required unless the gate system
cannot express composite evidence.

## Acceptance verification

| Obligation | Product result | Governed evidence |
| --- | --- | --- |
| REQ-1 / AC-1 / ARCH-INV-1 | pass: exact 249 unique sorted ISO records and reconciliation teeth | adjust: relevant ISO gate exists but REQ/AC evidence points to broad ceiling |
| REQ-2 / AC-2 | pass: both profiles resolve and exact five-entity oracle matches | adjust: relevant profile gate exists but REQ/AC evidence points to broad ceiling |
| REQ-3 | pass: five complete digest-matching receipts and documented transformations | adjust: evidence is not allocated to source/receipt-relevant assertions |
| REQ-4 / ARCH-INV-7 | pass: zero/single/multiple/add/replace/roles/primary behavior proven | adjust: relevant capital gate exists but REQ evidence points to broad ceiling |
| REQ-5 | pass: six named protected records, contained anchors and minimum parts | adjust: relevant protected gate exists but REQ evidence points to broad ceiling |
| REQ-6 | pass: declared failure classes and repaired invariants fail closed | adjust: composite focused evidence is not allocated |
| REQ-7 / AC-3 / INV-1 | pass: deterministic whole-set identity, create-once publication and reproducibility | pass: reproducibility allocation is relevant |
| AC-4 | pass: normal/multi/zero capital and island/microstate cases proven | adjust: requires capital plus protected evidence, not broad ceiling alone |
| VAL-1 | pass | pass: corpus-iso-integrity |
| VAL-2 | pass | pass: corpus-boundary-profiles |
| VAL-3 | pass | pass: corpus-reproducibility |
| VAL-4 | pass | pass: corpus-capital-registry |
| VAL-5 | pass | pass: corpus-protected-features |

## Test and E2E evidence

- `GOWORK=off go test ./internal/catalog/... -count=1` — pass.
- `make gen-check` — pass; expected corpus identity and counts reproduced.
- `make check` — pass.
- Exact-set identity, receipt/coverage mutations, unexpected-file rejection,
  all five oracle mutations, six-anchor containment, capital add/replace/
  roles/primary, 249-row reconciliation, and source digest probes — pass.
- Publication failure at the report destination after corpus replacement —
  non-zero with complete prior set restored.
- Publication failure at the embedded destination after corpus and report
  replacement — non-zero with complete prior set restored.

## Deviations and follow-ups

No product implementation or authority deviation remains. The only material
deviation is governed evidence allocation F-005-R2. The induced
permission-failure probes may leave non-authoritative staging debris in the
deliberately unwritable temporary directory, but all authoritative
destinations are restored; this does not reopen F-007.

## Verdict and closure eligibility

`adjust`. RUN-002 product implementation satisfies the accepted P1 contract
and closes six of seven prior findings, but P1 is not closure-eligible while
F-005-R2 leaves obligation evidence overbroad.

## Actionable corrections or closeout

Rematerialize the affected US/REQ/AC/INV verification results from existing
obligation-relevant focused gate receipts, using composite evidence where an
obligation spans multiple gates. Validate the corrected lifecycle corpus and
re-evaluate closure. No product remediation or product commit is indicated.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
