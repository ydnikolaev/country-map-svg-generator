---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "specification"
id: "P1"
epic: "country-map-svg-generator"
spec: "P1"
status: draft
profiles: []
concerns: []
inputs: ["DISC-006", "ARCH-001", "DEC-002", "DEC-003"]
---
# P1 — Versioned Country Geometry Corpus

## Outcome

A maintainer can compile pinned public sources into one immutable, versioned corpus
that maps every currently assigned ISO 3166-1 entity to canonical geometry,
boundary-profile variants, protected-feature metadata and capital records. Normal
generation can trust the corpus without consulting the network or repairing data.

## User stories and stakeholders

| ID | Story |
| --- | --- |
| US-1 | As a producer, I need every supported country addressed by ISO alpha-2 so site data and generated assets join predictably. |
| US-2 | As an owner, I need an explicit boundary posture so maps do not silently communicate an unintended political position. |
| US-3 | As a maintainer, I need source versions, licenses and transformations recorded so a corpus can be reproduced and audited. |
| US-4 | As a generator, I need capitals and protected islands/microstates encoded locally so rendering stays offline and visually complete. |

## Scope and non-goals

In scope: ISO catalog reconciliation, source receipts, admin geometry variants,
capital registry with role metadata, protected-feature metadata, compilation
validation, immutable bundle and comparison report.

Out of scope: runtime downloads, projection, simplification, SVG styling, CLI
commands, automatic scheduled refresh, arbitrary subnational data.

## Context and ground truth

ISO 3166-1 is catalog authority. Natural Earth is public domain and provides
admin-0 geometry and populated places, but represents de facto boundaries by
default; DEC-002 therefore requires explicit `un` and `de-facto` profiles.
Natural Earth capital records can contain multiple capital roles. The greenfield
repository has no legacy corpus to migrate.

## Requirements and invariants

| ID | Requirement |
| --- | --- |
| REQ-1 | The manifest contains exactly one record per currently assigned ISO 3166-1 alpha-2 entity and records alpha-3 and display-name metadata. |
| REQ-2 | Each record resolves geometry for both `un` and `de-facto`, either as an explicit variant or a declared identical reference. |
| REQ-3 | Source receipts name origin, upstream version/date, license, acquisition digest and transformation identity. |
| REQ-4 | Capital records support zero, one or multiple points, named roles, primary display choice and explicit per-country override. |
| REQ-5 | Protected-feature metadata identifies identity-defining islands or microstate geometry that downstream simplification may not silently erase. |
| REQ-6 | Compilation rejects missing/duplicate ISO mappings, invalid coordinates, empty required geometry, dangling references and undocumented source inputs. |
| REQ-7 | Corpus bytes and manifest ordering are deterministic and content-addressed. |
| INV-1 | An accepted corpus is immutable: refresh creates a new identity and never rewrites existing bytes. |

`ARCH-INV-1` and `ARCH-INV-7` from ARCH-001 are non-negotiable.

## Interfaces data and behavior

Input is an explicit maintainer source set plus mapping/override records. Output is
CTR-001: a corpus directory or embedded bundle containing a versioned manifest,
canonical geometry records, boundary variants, capitals, protected-feature
metadata and receipts.

Entity keys are uppercase ISO alpha-2. Coordinates remain source-geographic until
P2. A corpus version is immutable; a refresh produces a new version and a
machine-readable added/removed/changed report. Errors identify source, entity,
field and violated invariant.

## Dependencies and handoffs

Upstream sources are external and pinned by receipt. P1 has no epic-spec
dependency. It hands CTR-001 to P2 and exposes corpus identity to P3/P4/P5.
Owner approval is required for political overrides; maintainer approval is required
for source-version changes.

## Profiles and concerns

- Go: compilation and validation code follows repository Go conventions.
- Data provenance: public-domain status does not remove version/digest obligations.
- Geopolitical representation: DEC-002 governs profile semantics.
- Determinism: the corpus is a build input and must be reproducible.

## Validation impact

| ID | Protects | Guard/action | Tier and teeth | Owner/due |
| --- | --- | --- | --- | --- |
| VAL-1 | REQ-1, REQ-6 | compile the authoritative ISO fixture with one entity missing and one duplicated | integration; both mutations fail non-zero with entity diagnostics | P1 / W1 complete |
| VAL-2 | REQ-2 | remove one boundary-profile resolution | integration; missing profile fails closed | P1 / W1 complete |
| VAL-3 | REQ-3, REQ-7 | compile twice from identical pinned inputs | integration; byte/digest equality plus missing-receipt mutation failure | P1 / W1 complete |
| VAL-4 | REQ-4 | exercise zero/single/multi-capital and override fixtures | unit+integration; roles round-trip and invalid primary/coordinate fails | P1 / W1 complete |
| VAL-5 | REQ-5 | remove protected metadata from named microstate/island fixtures | integration; coverage gate names uncovered entity/features | P1 / W1 complete |

## Acceptance criteria

| ID | Covers | Criterion | Evidence |
| --- | --- | --- | --- |
| AC-1 | US-1, REQ-1 | Manifest count equals the pinned authoritative ISO set with no missing or duplicate alpha-2 key. | corpus gate receipt |
| AC-2 | US-2, REQ-2 | Every entity resolves both accepted boundary profiles and selected dispute fixtures match owner-approved expectations. | profile fixture report |
| AC-3 | US-3, REQ-3, REQ-7 | Two clean compilations are byte-identical and every input has a valid receipt. | reproducibility receipt |
| AC-4 | US-4, REQ-4, REQ-5 | Capital and protected-feature fixtures cover normal, multi-capital, island and microstate cases. | registry coverage report |

## Rollout rollback and operations

The first accepted corpus is version 1. Later refreshes publish side-by-side with a
diff report. The generator pins one corpus identity. Rollback selects the prior
bundle; no in-place data migration exists. Compiler output reports entity counts,
profile variants, capital coverage, protected-feature coverage and digest.

## Delivery constraints

BND-001 and BND-002 are exclusive to P1. Compilation may use maintainer-only
source adapters, but runtime generation cannot. Do not hand-edit compiled bytes,
infer ISO identity from display names, silently synthesize missing boundaries, or
couple corpus records to SVG/CSS presentation. Prefer a compact representation only
after correctness and deterministic decoding are proven.

## Decisions and unresolved questions

DEC-002 and DEC-003 are accepted and binding. No unresolved authority blocks P1.
Source-version selection is an implementation-plan input and must be pinned before
the run begins.

## Amendments

None.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
