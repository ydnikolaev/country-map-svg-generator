---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "audit"
id: "AUD-003"
epic: "country-map-svg-generator"
spec: "P1"
status: completed
profiles: []
concerns: []
inputs: ["RUN-002-RESULT"]
---
# AUD-003 — Audit

## Subject scope and exclusions

Audited exact implementation commit
`ed4fdca23965b6aa0c0fc1b8e4c449e79f2d45ee`, parent
`bd23a1423ade0ab3686a89cf250ca37581f1543f`, tree
`1a6bfca82d26027956f92e28c9cd479c277faba5`, against accepted P1,
PLAN-001, sealed AUD-001 and AUD-002, RUN-002-RESULT, the complete RUN-002
focused/project gate corpus, and harness finding `FND-ED003909B82A`.

Later lifecycle-only commits were excluded from the product subject. Product
roots do not differ between the audited commit and current HEAD. Audit
materialization, lifecycle mutation, sealing, and commit are coordinator-owned
and excluded from this read-only review.

## Inputs and evidence

Authority hashes independently match: specification
`a2dd6817331c8da19775f66de49b7aa4526651f06292f8254606c3b3fe3ea3bd`;
PLAN-001 body
`4e39f29252ac88caf43144e625c1cbff18111720f7cb1a3fffde2d3e7ed82a27`;
PLAN-001 manifest
`4cd93236c76fbb76ecbff1f83dff79070393718a691925a3f8c5bd2fcf774a6e`;
AUD-001
`da0f7b477e6d044ddb5cedf66b16c41ab8703bc9cfc6c10ccccabc7ba5033be6`;
AUD-002
`5c3fa16baba13eeb15fbaf9e810ea6574294bd318c14f6eae2c67a7a677673b7`;
RUN-002-RESULT
`eb64507d5b9c108bb25492c8d2091c3335ca968e73485571783e62f302b99282`;
and RESULT-007
`0c87a88f4bfc4821824ceac0a822f2ebc786360f492da4d299e0bf3636b02ada`.

RUN-002 contains five passing receipts against the same brief:

- `GATE-EE3514765774`, `corpus-iso-integrity`;
- `GATE-6F08DD3DA870`, `corpus-boundary-profiles`;
- `GATE-B475A8CCACA3`, `corpus-capital-registry`;
- `GATE-5D6A3BFF82E9`, `corpus-protected-features`;
- `GATE-59FD08536411`, `corpus-reproducibility`.

Their logs prove missing/duplicate ISO rejection, receipt completeness, 249-row
reconciliation, exact five-entity boundary oracles, capital add/replace/roles/
primary behavior, protected-anchor containment and sabotage, and the full
project ceiling. AUD-002 adds independent whole-set identity/source hashing,
embedded-byte comparison, geometry-oracle hashing, six-anchor containment,
capital cases, and real rollback failures after the first and second
authoritative destination replacements.

## Method and independence

Fresh principal `codex-p1-auditor-3`; no implementation, prior-audit
authorship, remediation, product write, lifecycle mutation, or commit. Reviewed
the accepted spec and plan, exact revision and footprint, sealed audits, run
result, every focused gate receipt/selection/result/log, representative
generated obligation evidence, publication code/tests, and the harness finding.

Fresh verification:

- `make check` — pass;
- `make gen-check` — pass, reproducing identity
  `sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`,
  249 entities, 283 geometries, 212 capitals, and six protected entities;
- focused identity/create-once/staging/complete-publication tests — pass,
  including receipt, coverage, manifest, geometry, and extra-file mutations;
- product-root diff and worktree checks — clean.

Provider limitations for native model selection and capability enforcement are
disclosed; semantic independence and read-only behavior were observed.

## Findings

- **F-005-R2 — harness-owned, non-product follow-up.** Generated non-VAL
  evidence allocation remains overbroad. It is not waived or represented as
  corrected. Its terminal disposition is `FND-ED003909B82A`, scope `mate`,
  severity `high`, area `evidence-execution`, release `v1.4.1`, state `new`.
  Upstream Mate must support focused/composite non-VAL bindings without
  mutating accepted artifacts.

No product, specification, authority, ownership, source-version, runtime,
rollback, or regression finding remains.

### Closure-boundary decision

F-005-R2 is genuine: US/REQ/AC/INV results bind to
`SELECT-59FD08536411` and the broad `corpus-reproducibility` result.
It does not leave an unproved product or specification obligation:

1. Immutable focused gate receipts contain the relevant assertions and teeth.
2. AUD-002 independently exercised and correlated those assertions to the
   exact commit.
3. AUD-003 rechecked authority, revision, receipts, ceiling, identity, and
   publication behavior.
4. PLAN-001's five declared VAL bindings are correctly focused; the defect is
   in generated non-VAL coverage.
5. Mate v1.4.1 `compileDeclaredEvidence` forces non-validation coverage to the
   single project-tier gate; accepted evidence cannot be honestly rebound here.
6. The limitation and owner are durably recorded as `FND-ED003909B82A`.

Changing product code would therefore be unrelated remediation, editing
accepted evidence would break immutability, and editing the managed consumer
would violate harness ownership. No executable P1 correction remains.

## Acceptance verification

| Obligation | Result | Composite evidence |
| --- | --- | --- |
| US-1 / REQ-1 / AC-1 / ARCH-INV-1 | pass | ISO receipt plus exact 249-row reconciliation and mutations |
| US-2 / REQ-2 / AC-2 | pass | Profile receipt plus five exact reviewed two-profile geometry oracles |
| US-3 / REQ-3 | pass | ISO/receipt assertions, project ceiling, and independent five-source hashing |
| US-4 / REQ-4 / ARCH-INV-7 | pass | Capital receipt plus add/replace/roles/primary positive and negative cases |
| REQ-5 | pass | Protected receipt plus six contained anchors and ocean/wrong-entity/part sabotage |
| REQ-6 | pass | Composite ISO, profile, capital, protected, identity, and ceiling failures |
| REQ-7 / INV-1 / AC-3 | pass | Reproducibility, fresh gen-check, whole-set identity, create-once, rollback |
| AC-4 | pass | Capital plus protected receipts and independent representative-case review |
| VAL-1 | pass | `corpus-iso-integrity` |
| VAL-2 | pass | `corpus-boundary-profiles` |
| VAL-3 | pass | `corpus-reproducibility` |
| VAL-4 | pass | `corpus-capital-registry` |
| VAL-5 | pass | `corpus-protected-features` |

All product requirements, invariants, acceptance criteria, validation
obligations, rollback behavior, ownership boundaries, and completion conditions
are satisfied.

## Test and E2E evidence

- Five focused/project RUN-002 receipts and their exact logs — pass.
- `make check` — fresh pass.
- `make gen-check` — fresh pass with the accepted identity and counts.
- Identity and publication-focused adversarial tests — fresh pass.
- Product roots unchanged since the subject commit.
- AUD-002 second/third-destination failure injection — non-zero failures with
  every prior authoritative output restored.

## Deviations and follow-ups

The only deviation is F-005-R2, terminally captured as the harness finding
`FND-ED003909B82A`. It remains open as upstream harness work and must not be
silently discarded or represented as product debt.

## Verdict and closure eligibility

`pass`.

Exact product commit
`ed4fdca23965b6aa0c0fc1b8e4c449e79f2d45ee` satisfies accepted P1 and
PLAN-001. P1 is product/spec closure-eligible. F-005-R2 continues only as the
harness-owned non-product follow-up `FND-ED003909B82A`.

## Actionable corrections or closeout

1. Materialize, validate, seal, and record this AUD-003 pass against unchanged
   RUN-002.
2. Close P1 without another implementation run or product commit.
3. Retain `FND-ED003909B82A` in the harness backlog and route its correction
   upstream; do not rewrite RUN-002 evidence.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
