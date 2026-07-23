---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "epic-overview"
id: "country-map-svg-generator"
epic: "country-map-svg-generator"
status: draft
profiles: []
concerns: []
inputs: ["DISC-006", "ARCH-001", "DEC-001", "DEC-002", "DEC-003"]
---
# country-map-svg-generator — Country Map SVG Generator

## Outcome

A local, agent-friendly CLI generates an optimized, visually coherent SVG map
catalog for every currently assigned ISO 3166-1 entity. The catalog is ready to
place into the owner's visa site without a post-optimization step, supports
themeable inline hero artwork and lightweight decorative card artwork, and can add
capital markers when requested.

Success means the complete catalog passes deterministic, structural, visual and
byte-budget gates and the owner can integrate it using the documented producer and
client paths.

## Stakeholders and users

- Product owner/site producer: selects the site style, approves visual balance,
  political-boundary posture and first production integration.
- AI agents: primary CLI operators; require stable JSON, explainability and
  non-interactive completeness.
- Colleagues: later local-tool consumers on supported desktop/server platforms.
- Site developers and browsers: consume standalone, inline-themed and mask-based
  assets without shipping generator code.
- Maintainers: refresh upstream sources and release corpus/binary versions.

## Scope and non-goals

In scope:

- complete ISO catalog with `un` and `de-facto` boundary profiles;
- versioned offline corpus and capital registry;
- centered equal-area, soft-organic `hero` and `card` geometry;
- portable style tokens, presets, inheritance and per-country overrides;
- standalone and themed-inline SVG contracts;
- optional capital/custom markers;
- deterministic optimized generation, manifest and agent-first CLI;
- full-catalog site proof and producer/client performance documentation;
- portable binary handoff after the local site path works.

Non-goals for this epic:

- hosted service, GUI, web editor or public SaaS;
- runtime geodata downloads or automatic upstream refresh;
- cartographic navigation/analytics or geographic precision at arbitrary zoom;
- country flags, labels or general-purpose map composition;
- a `detail` profile, public marketplace/plugin packaging and repository marketing.

## Context and ground truth

Discovery checkpoint DISC-006 captures the accepted product contract. ARCH-001 is
the canonical architecture. ISO 3166-1 is catalog authority; Natural Earth is the
candidate public-domain source for admin geometry and capital points; D3 is the
golden projection/path oracle; the owner's supplied site screenshot fixes the
hero/card usage scales. DEC-001 through DEC-003 hold elevated risk, boundary and
runtime architecture authority.

The repository is greenfield for product code. The mate consumer now carries the
Go CLI profile (`lang: go`, `kind: cli`).

## Specification map

| Spec | Outcome | Depends on | Owner |
| --- | --- | --- | --- |
| P1 | trusted versioned country/capital corpus | — | P1 |
| P2 | deterministic soft-organic projected geometry | P1 | P2 |
| P3 | local agent-first generator and SVG/config contracts | P2 | P3 |
| P4 | complete site-ready catalog and integration proof | P3 | P4 |
| P5 | portable colleague handoff | P3 | P5 |

## Architecture and decisions

[ARCH-001](architecture/country-map-svg-generator.solution-architecture.md) is
the one canonical solution architecture. DEC-001 selects elevated decision policy;
DEC-002 governs boundary profiles; DEC-003 governs the Go/corpus separation.

## Dependencies and external requests

- Upstream: pinned ISO entity data, public-domain geometry/capital data and
  explicit source/license/version receipts.
- Build dependencies: Go module dependencies selected by implementation plans;
  D3 fixtures are development oracle data, never runtime code.
- Downstream: the visa site consumes generated SVG plus its manifest; colleagues
  consume release archives only after P5.
- Harness prerequisite is complete: the repository has the mate Go CLI profile.

## Validation and E2E strategy

Each spec owns requirement-level unit/integration teeth. P3 follows the Go profile:
build with `GOWORK=off`, drive the binary via `os/exec`, cover command breadth with
hermetic `testscript` scenarios, and keep one ordered product-loop test. P4 owns
the full-catalog and browser/site boundary because only it can prove combined
geometry, rendering, budget and client behavior. P5 owns clean-machine package
proof. The project ceiling remains red until every due E2E allocation passes.

## Delivery waves and forecast

| Wave | Specs | Forecast | Confidence |
| --- | --- | --- | --- |
| W1 | P1 | medium | medium; source reconciliation is the main uncertainty |
| W2 | P2 | medium–large | medium; full-catalog visual behavior needs evidence |
| W3 | P3 | medium | high after P2 contracts stabilize |
| W4 | P4 and P5 | P4 medium, P5 small–medium | P4 medium, P5 high |

Critical path is `P1 → P2 → P3 → P4`. The first usable local generator exits W3;
site-ready output exits P4. Public/release polish in P5 is deliberately off the
owner-value critical path.

## Risks and open questions

| Risk | Treatment | Authority/status |
| --- | --- | --- |
| disputed borders and country identity | explicit profile plus corpus fixtures | DEC-002 accepted |
| islands/microstates disappear under simplification | protected-feature metadata and full-catalog gates | P1/P2 |
| Go implementation drifts from projection/path oracle | pinned D3 golden fixtures and narrow owned modules | DEC-003 accepted |
| repeated cards cost too much network/DOM work | card budgets, mask guidance, lazy/use-site policy | P4 |
| flexible styling bloats every SVG | geometry/style separation and profile-specific serialization | P3 |
| public packaging distracts from site delivery | parallel P5 after P3; P4 is critical path | accepted rollout |

There are no open questions requiring owner authority. `detail` remains a tracked
backlog item after hero/card evidence exists.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
