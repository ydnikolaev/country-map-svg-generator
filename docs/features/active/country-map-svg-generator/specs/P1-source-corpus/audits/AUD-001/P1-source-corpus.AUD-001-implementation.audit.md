---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "audit"
id: "AUD-001"
epic: "country-map-svg-generator"
spec: "P1"
status: completed
profiles: []
concerns: []
inputs: ["RUN-001-RESULT"]
---
# AUD-001 — Audit

## Subject scope and exclusions

Audited exact implementation commit
`6f5d5fdd73c4bfbdc13202afec7432e2ae24c974`, parent
`668023213ad90529bad519da1b7a1fe5e42dc617`, tree
`9de827a0a57d7d4f8eed1fac978f9d0582648652`, against accepted P1,
PLAN-001, DEC-002, DEC-003, ARCH-INV-1, ARCH-INV-7, RUN-001-RESULT,
committed sources/corpus, and recorded gates. Coordinator-owned uncommitted
lifecycle material and unrelated worktree state were excluded from the product
candidate.

## Inputs and evidence

Authority hashes match BRIEF-005: role contract
`ac38fb1d3936459fc1480d239a4d24279d5b600404b6a8d8a07e12f3ee6ed654`,
specification
`a2dd6817331c8da19775f66de49b7aa4526651f06292f8254606c3b3fe3ea3bd`,
plan manifest
`4cd93236c76fbb76ecbff1f83dff79070393718a691925a3f8c5bd2fcf774a6e`,
and run result
`de1aacd394ec4a00fdfedd6aae532c53663ce60f00dbddb23d795c0a59a8ce90`.
RESULT-004 exactly enumerates all 34 commit paths and hashes. Local source receipt
hashes pass. Independent Natural Earth v5.1.2 downloads byte-match all three
committed upstream files. The local ISO snapshot has 249 unique sorted
alpha-2/alpha-3 rows, but no retained raw OBP export or independent review
attestation allows semantic row-by-row reconstruction.

The corpus reports identity
`sha256:d52557c62d119da3c447acf39817a6eff6fe8fbc7d70db9ebe6aa9dd49c16c4f`,
249 entities, 283 geometries, 249 UN resolutions, 34 explicit de-facto
resolutions, 215 identical references, 212 capital records over 195 entities, and
six protected entities.

## Method and independence

Fresh principal `codex-p1-auditor-1`; no authorship, remediation, product writes,
lifecycle mutation, or audit materialization. Reviewed authority, exact
commit/parent/diff, source receipts, compiler/decoder/publication code, corpus
outputs, all tests, gate receipts/logs, and evidence allocation. Reran focused
tests, `make gen-check`, and `make check`; performed independent upstream-hash,
content-identity, and protected-anchor probes. Provider enforcement degradations
are disclosed in the installed role contract.

## Findings

- **F-001 — blocker.** Published corpus bytes are not content-addressed. Identity
  covers only manifest and geometry, excluding receipts and coverage; validation
  ignores both. Different corpus bytes can retain the same identity. Owner:
  catalog publisher. Bind every authoritative component into identity, validate
  receipts/coverage on decode, regenerate, and add receipt/coverage mutations.
- **F-002 — blocker.** AC-2 lacks an exact approved profile oracle. Five entities
  are checked only for non-identical IDs and the sole sabotage collapses CN;
  arbitrary incorrect but distinct geometry passes. Owner: corpus/profile policy.
  Commit exact expected profile geometry identities (or equivalent semantic
  oracle) for the accepted pinned source profiles and mutate every fixture.
- **F-003 — blocker.** Indonesia's protected anchor `[118,-2]` lies in none of its
  264 UN geometry parts; validation checks only name, coordinate range, minimum
  parts, and hard-coded presence. Owner: protected policy. Correct the anchor,
  validate containment against resolved geometry, and add ocean/misbound sabotage.
- **F-004 — major.** Capital overrides can only select an existing
  `primary_name`; they cannot replace or augment capital records or roles as the
  accepted plan requires, including zero-capital cases. Owner: catalog compiler.
  Implement additions/replacements/roles/primary semantics and positive/negative
  tests.
- **F-005 — major.** Evidence allocation overclaims obligation coverage by binding
  broad project output to unrelated acceptance claims. Owner: coordinator/gate
  registry. Reissue focused evidence after remediation and preserve
  obligation-relevant assertions.
- **F-006 — major.** The ISO receipt authenticates the transformed snapshot, not a
  reconstructable raw authority acquisition or row-level review attestation.
  Owner: source maintainer. Retain a reviewed reconciliation artifact mapping all
  local rows to the observed authority snapshot where direct export retention is
  unavailable.
- **F-007 — moderate.** Publication is not atomic across corpus, report, and
  embedded output: corpus publication completes before later writes. Owner: corpus
  publisher. Stage and validate the complete output set before publication.

## Acceptance verification

| Obligation | Result | Evidence or gap |
| --- | --- | --- |
| REQ-1 / ARCH-INV-1 | partial | 249 sorted unique uppercase alpha-2 records with alpha-3/name and mutation teeth pass; authoritative row provenance is not independently reconstructable. |
| REQ-2 | partial | Both profiles resolve for all 249; exact accepted dispute semantics lack an oracle. |
| REQ-3 | partial | Four complete receipts and Natural Earth hashes pass; ISO transformed-source lineage remains weak. |
| REQ-4 / ARCH-INV-7 | partial | Zero/single/multi roles and primary selection work; replace/augment override semantics are absent. |
| REQ-5 | fail | Six protected records exist, but Indonesia's anchor is outside its geometry. |
| REQ-6 | partial | Declared failures are exercised; decoded receipt/coverage tampering is not rejected. |
| REQ-7 / INV-1 | fail | Repeated builds are byte-identical, but not every published corpus byte participates in identity. |
| AC-1 | partial | Count/uniqueness proven; authoritative ISO row reconciliation lacks audit evidence. |
| AC-2 | fail | No exact approved profile fixture report or oracle. |
| AC-3 | fail | Reproducibility passes; full corpus content addressing does not. |
| AC-4 | fail | Capital representative cases exist, but protected metadata and override behavior are incomplete. |
| VAL-1 | pass | Missing and duplicate mutations fail with invariant diagnostics. |
| VAL-2 | partial | Missing/collapsed profile fails; incorrect-but-distinct geometry is not detected. |
| VAL-3 | partial | Double compilation/gen-check passes; receipt/coverage identity mutations are absent. |
| VAL-4 | partial | Existing cases pass; override role/record augmentation is absent. |
| VAL-5 | fail | Removal is detected, but invalid/ocean-bound anchors pass. |

## Test and E2E evidence

Independent rerun: focused catalog tests pass; `make gen-check` passes; `make
check` passes. Natural Earth remote/local hashes pass. Runtime boundary passes:
standard library only, embedded corpus, no runtime network/source-file/subprocess
acquisition. Green tests are insufficient for closure because their oracles omit
F-001 through F-004.

## Deviations and follow-ups

No changed-path, ownership-boundary, source-version, dependency, runtime-network,
or commit-parent deviation was found. Substantive deviations are incomplete
content identity/validation, absent profile oracle, invalid protected anchor,
narrower capital override contract, weak ISO review provenance, and non-atomic
complete publication. Each requires correction in a successor same-plan run; none
is audit-remediated here.

## Verdict and closure eligibility

`block`. RUN-001 is not closure-eligible. REQ-5, REQ-7, INV-1, AC-2, AC-3,
AC-4, VAL-2, VAL-3, and VAL-5 lack sufficient proof or are contradicted by direct
probes.

## Actionable corrections or closeout

1. Correct corpus identity/decoder validation and regenerate under a new identity.
2. Materialize exact accepted dispute expectations and mutation oracles.
3. Correct protected anchors and enforce geometry containment.
4. Complete capital override semantics and tests.
5. Make publication atomic across corpus/report/embed.
6. Strengthen ISO row-level reconciliation provenance.
7. Rerun focused gates, `make gen-check`, and `make check`, regenerate honest
   evidence, and dispatch a fresh independent auditor.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
