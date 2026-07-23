---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "specification"
id: "P5"
epic: "country-map-svg-generator"
spec: "P5"
status: draft
profiles: []
concerns: []
inputs: ["DISC-006", "ARCH-001", "DEC-003"]
---
# P5 — Portable Distribution and Handoff

## Outcome

A colleague can download or receive one archive for a supported platform, verify
it, run the generator offline without installing Go/Node/Python/GDAL, and reproduce
the documented product loop. Public-facing polish is delayed until the owner's
local generator and site path are already usable.

## User stories and stakeholders

| ID | Story |
| --- | --- |
| US-1 | As a colleague, I can unpack one archive and run the CLI immediately on my platform. |
| US-2 | As a maintainer, I can produce traceable checksummed artifacts from one accepted source state. |
| US-3 | As an advanced Go user, I can install the same CLI with `go install`. |
| US-4 | As the owner, distribution work does not delay the local site deliverable. |

## Scope and non-goals

In scope: macOS arm64/amd64, Linux arm64/amd64 and Windows amd64 binaries; archives,
checksums, embedded corpus/schema/presets, version reporting, clean-environment
smoke tests, `go install` compatibility and colleague operating guide.

Out of scope: package managers, signed/notarized installers, hosted update service,
telemetry, public documentation site, marketplace/plugin packaging and automatic
release publication before owner acceptance.

## Context and ground truth

The owner asked for an autonomous cross-platform binary but explicitly prioritized
the local generator and site integration. Go/DEC-003 supports one binary and
embedded immutable data. The repository Go profile requires proving the installed
binary rather than only testing imported packages.

## Requirements and invariants

| ID | Requirement |
| --- | --- |
| REQ-1 | Release matrix contains darwin/arm64, darwin/amd64, linux/arm64, linux/amd64 and windows/amd64 executables. |
| REQ-2 | Each archive contains only the executable, license/source notices, quick-start/operations guide and any explicitly required non-runtime metadata; corpus/schema/presets are embedded. |
| REQ-3 | Checksums cover every archive and `version --json` reports generator, corpus and config-schema identities. |
| REQ-4 | Each platform artifact runs `version`, `validate` and representative generation with network unavailable and no language/geospatial runtime installed. |
| REQ-5 | `go install` builds with `GOWORK=off` from the declared module path and produces behavior compatible with release archives. |
| REQ-6 | Release configuration is deterministic, reviewable and cannot publish from a source state that has not passed the project ceiling. |
| REQ-7 | Operations documentation covers install, upgrade, rollback, corpus identity, checksum verification, troubleshooting and source-refresh separation. |
| INV-1 | A release archive is immutable and self-sufficient for normal offline generation on its named platform. |

## Interfaces data and behavior

CTR-007 is `<name>_<version>_<os>_<arch>.<archive>`, a checksum manifest and
versioned guide. The executable exposes CTR-004 unchanged. Archive installation
does not create global state; users choose config/output paths. Upgrade is binary
replacement plus explicit regeneration when identities change.

## Dependencies and handoffs

Depends on P3. P4 is not a technical dependency and runs in parallel, but public
publication remains sequenced after owner acceptance. Owns BND-009 and BND-010.
Release artifacts consume accepted source, corpus and project-ceiling receipts.

## Profiles and concerns

- Go portability: cross-compilation, `GOWORK=off` and embedded assets.
- Supply-chain integrity: checksums and traceable source identity.
- Operations: explicit upgrade/rollback and no implicit online updater.
- Scope control: distribution cannot enter the P4 critical path.

## Validation impact

| ID | Protects | Guard/action | Tier and teeth | Owner/due |
| --- | --- | --- | --- | --- |
| VAL-1 | REQ-1, REQ-2, REQ-3 | build matrix and inspect archive names/content/checksums/version output | package gate; remove one target/file/checksum field and fail | P5 / W4 complete |
| VAL-2 | REQ-4 | run supported binaries in clean platform environments with network disabled | package E2E; asserts no auxiliary runtime/data lookup | P5 / W4 complete |
| VAL-3 | REQ-5 | install/build module with `GOWORK=off` outside repository workspace | Go-profile E2E; local replacement cannot mask failure | P5 / W4 complete |
| VAL-4 | REQ-6 | attempt packaging with stale/red/missing project-ceiling evidence | release gate; publication plan refuses | P5 / W4 complete |
| VAL-5 | REQ-7 | execute quick-start and rollback examples against an unpacked archive | docs smoke; commands and reported identities match | P5 / W4 complete |

## Acceptance criteria

| ID | Covers | Criterion | Evidence |
| --- | --- | --- | --- |
| AC-1 | US-1, REQ-1, REQ-2, REQ-3, REQ-4 | A clean machine for each target can unpack, verify and generate representative SVG offline. | platform E2E matrix |
| AC-2 | US-2, REQ-3, REQ-6 | Every archive is checksummed and bound to accepted source/corpus/schema identities. | release manifest receipt |
| AC-3 | US-3, REQ-5 | `go install` succeeds under `GOWORK=off` and passes the representative CLI loop. | installation receipt |
| AC-4 | US-4 | P4 can complete without P5 and public publication occurs only after owner-first acceptance. | dependency DAG and release gate |

## Rollout rollback and operations

First retain a local development binary, then create private colleague archives,
then automate/publicize releases after P4 acceptance. Versions are immutable.
Rollback replaces the executable with a previous archive and regenerates from the
recorded config/corpus pair. No automatic update check or telemetry is introduced.

## Delivery constraints

Use the ecosystem release tool at its highest-level supported configuration; avoid
custom packaging scripts unless a proven requirement is unsupported. Do not bundle
Node/Python/GDAL, platform installers, auto-updaters or public-site work into this
spec. Release archives must remain small enough that the embedded corpus, not
packaging overhead, is the dominant data cost.

## Decisions and unresolved questions

DEC-003 governs the single-binary contract. Public release timing remains
owner-controlled, but it does not block implementation or private archive proof.

## Amendments

None.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
