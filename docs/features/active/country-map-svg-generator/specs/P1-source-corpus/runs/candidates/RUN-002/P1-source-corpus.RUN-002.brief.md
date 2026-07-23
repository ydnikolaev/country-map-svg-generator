---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-002"
epic: "country-map-svg-generator"
spec: "P1"
run: "RUN-002"
status: draft
profiles: []
concerns: []
inputs: ["PLAN-001", "CTX-RUN-002"]
---
# RUN-002 — Run brief

## Authority and accepted inputs

- Accepted plan revision: `PLAN-001` at SHA-256
  `4e39f29252ac88caf43144e625c1cbff18111720f7cb1a3fffde2d3e7ed82a27`.
- Accepted plan manifest: `PLAN-001` at SHA-256
  `4cd93236c76fbb76ecbff1f83dff79070393718a691925a3f8c5bd2fcf774a6e`.
- Run context pack: `CTX-RUN-002` at SHA-256
  `5f5d7fcb829c078920e9dd3a49ec1161492822b0018c7229a31fb51da160f657`.
- Blocking audit: `AUD-001` at SHA-256
  `da0f7b477e6d044ddb5cedf66b16c41ab8703bc9cfc6c10ccccabc7ba5033be6`.
- Same-plan remediation authority: `ADV-001` at SHA-256
  `d8e0df1d67b89a48924ffc7443d52bc62225dafd519ee22b8b5620618bc97b2a`.
- Product baseline: commit
  `6f5d5fdd73c4bfbdc13202afec7432e2ae24c974`, tree
  `9de827a0a57d7d4f8eed1fac978f9d0582648652`.

DEC-002 and DEC-003 remain binding. This run may repair only the implementation
and evidence gaps identified by AUD-001; it may not revise the accepted WHAT or
macro-HOW. Run readiness is bound atomically by `mate run start`.

## Objective and completion boundary

Remediate all seven AUD-001 findings without changing the accepted P1 architecture:
make the published corpus identity cover every authoritative artifact, enforce an
exact reviewed boundary-profile oracle, validate protected-anchor containment,
support additive and replacement capital overrides, provide a reconstructable
249-row ISO reconciliation attestation, publish the complete output set atomically,
and replace overbroad evidence with focused mutation teeth.

Finish only when each finding has a concrete code/data/test disposition, all
focused gates pass against adversarial mutations, two clean compilations are
byte-identical, the committed and embedded bundles agree, and `make check` is
green offline.

## Scope ownership and footprint

Writable product paths remain exactly the PLAN-001 footprint:

- `go.mod`, `go.sum`, `Makefile`;
- `data/sources/**`, `data/policy/**`, `data/corpus/v1/**`,
  `data/reports/v1/**`;
- `internal/catalog/**`, including compiler, decoder, validation, tests, and the
  generated embedded bundle.

Forbidden product roots: `internal/geometry/**`, `internal/config/**`,
`internal/render/**`, `cmd/**`, `docs/guides/**`, `docs/operations/**`,
`.goreleaser.yaml`, `.git/**`, `.mate/**`, and all governed lifecycle documents.

The integration target is the current `main` line. The implementer owns one atomic
remediation commit and must preserve unrelated worktree changes.

## Context and doctrine pack

Use `CTX-RUN-002`, the accepted plan, and AUD-001 as the complete authority set.
The pinned sources and 249-entity contract remain unchanged. The Go profile is
active; all build and black-box checks use `GOWORK=off` and must remain hermetic:
no network, user-home state, locale dependence, Node, Python, GDAL, or ambient
workspace data.

The audit's adversarial probes are required inputs, not optional suggestions. When
an audit finding admits several local implementations, choose the smallest design
that preserves deterministic canonical output and fail-closed validation.

## Tasks

1. F-001/F-007: define a documented canonical identity envelope over every
   authoritative published component without self-reference; make decode/validate
   verify receipts and coverage; stage and atomically replace the entire corpus,
   report, and embedded-output set.
2. F-002/F-003: add a reviewed exact per-profile oracle for selected
   politically-sensitive fixtures; validate expected profile digests/IDs; verify
   every protected anchor is contained by its entity geometry and correct the
   invalid Indonesia anchor from authoritative geometry.
3. F-004: extend capital override policy to add, replace, assign roles and primary
   status, including entities with zero discovered capitals; retain strict
   coordinate/reference validation and deterministic ordering.
4. F-006: add a row-level, digest-bound ISO reconciliation attestation covering
   all 249 assigned entities and validate it against the snapshot. Do not claim a
   raw ISO export that is not present.
5. F-005: add focused positive and negative fixtures for every repaired invariant,
   regenerate the corpus, and run the complete focused and project ceilings.

Execute serially in dependency order. Do not weaken an existing invariant while
adding the new oracle.

## Validation and evidence

- Identity/receipt/coverage teeth must fail after changing any authoritative
  published byte while leaving the rest untouched, and after deleting, adding, or
  mis-digesting a receipt or coverage row.
- Publication must prove that a failed staged compile cannot leave a mixed
  old/new output set; a complete successful compile must replace all outputs.
- Boundary-profile tests must compare selected `un` and `de-facto` output to exact
  reviewed expectations, not merely prove that both are present or different.
- Protected-feature tests must include containment success plus ocean and
  wrong-entity mutations; Indonesia's accepted anchor must lie inside its geometry.
- Capital tests must cover zero-to-one addition, replacement, role assignment,
  primary selection, invalid coordinates, duplicates, and dangling references.
- ISO attestation tests must prove exactly 249 row-level reconciliations and fail
  on missing, duplicate, mismatched, or unbound rows.
- Reproducibility must compare two clean compilations and committed/embedded
  decoding at both byte and identity levels.
- Focused commands:
  `GOWORK=off go test ./internal/catalog/... -count=1`,
  `make gen-check`, and `make check`.

The result must map F-001 through F-007 to exact changed paths and test evidence,
name all commands/outcomes and the product commit, and disclose any residual
uncertainty. Passing only the pre-audit tests is insufficient.

## Constraints and forbidden actions

- Preserve the 249 assigned-entity set, source pins, `un`/`de-facto` semantics,
  corpus v1/P2 contract, offline runtime, and standard-library-only preference.
- Do not synthesize geometry, infer ISO identity from names, silently discard
  entities, or manufacture source provenance.
- Do not hand-edit compiled corpus or embedded bytes; modify authoritative source,
  policy, mapping, or compiler and regenerate.
- Do not change the accepted plan, source versions, political policy, ownership
  roots, public API scope, or future P2 rendering behavior.
- Do not weaken validators, fixtures, or audit expectations to make the gate green.
- Do not commit lifecycle evidence or unrelated files.

## Stop escalation and resume conditions

Stop before mutation if remediation requires a source-version change, new political
decision, changed corpus/P2 interface, new write root, runtime acquisition, or a
claim of authority unsupported by checked-in evidence. Route such changes back to
the teamlead for amendment or plan supersession.

Local schema extensions, compiler/helper splits, deterministic staging mechanics,
and diagnostic improvements inside the accepted footprint are allowed and must be
reported. After several speculative failures on one invariant, reduce to a minimal
fixture and inspect the owning code/source contract. Resume only from the accepted
plan, CTX-RUN-002, AUD-001/ADV-001, this brief, and the current product commit/diff.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
