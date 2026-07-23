---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "specification"
id: "P3"
epic: "country-map-svg-generator"
spec: "P3"
status: draft
profiles: []
concerns: []
inputs: ["DISC-006", "ARCH-001", "DEC-003", "AM-001"]
---
# P3 — Agent-First SVG Generator CLI

## Outcome

A producer or AI agent can initialize, validate, explain and generate optimized SVG
assets locally through one deterministic Go binary. Configuration expresses style,
delivery mode, pins, animation and country exceptions without changing code.
Generated assets are immediately usable either standalone or inline with site-theme
CSS, and every batch carries a machine-readable manifest.

## User stories and stakeholders

| ID | Story |
| --- | --- |
| US-1 | As an AI agent, I can discover commands/schema, run without prompts, receive stable JSON and explain resolved configuration. |
| US-2 | As a producer, I can create reusable branded presets and override one country without duplicating the whole config. |
| US-3 | As a frontend developer, I can theme inline SVG across light/dark modes and add subtle motion without editing generated paths. |
| US-4 | As a build pipeline, I can regenerate a selected or full catalog offline and verify exactly what was emitted. |
| US-5 | As a visitor, I receive safe, accessible SVG with no script/external payload and reduced-motion compliance. |

## Scope and non-goals

In scope: versioned YAML/JSON config, inheritance and overrides, embedded presets
and schemas, CLI commands, human/JSON diagnostics, five base visual styles,
standalone/themed-inline delivery, optional markers, subtle animation hooks, minimal
SVG serialization, generic natural/fixed-frame layout controls, transactional
output and generation manifest.

Out of scope: interactive TUI for MVP, browser editor, hosted API, raster export,
custom JavaScript animation, arbitrary SVG template injection, source/corpus
compilation and public release automation.

## Context and ground truth

The primary operator is an AI agent, so the CLI must be fully headless and
self-describing. Inline SVG can inherit `currentColor` and CSS custom properties;
an SVG loaded through `<img>` cannot consume host-page theme variables, while an
SVG mask can efficiently decorate repeated cards. The owner wants only slight,
optional motion. AM-001 makes natural-aspect `tight` the geometry default and
retains arbitrary dimensions as an explicit distortion-free `contain` frame. The
mate Go CLI profile is active and governs binary-level E2E.

## Requirements and invariants

| ID | Requirement |
| --- | --- |
| REQ-1 | CLI exposes `init`, `validate`, `preview`, `generate`, `explain`, `inspect` and `version`, with complete `--help`, non-interactive flags and stable `--json`. |
| REQ-2 | Config is versioned YAML or JSON and supports embedded presets, single-parent inheritance, global tokens, profile overrides and per-country overrides with deterministic precedence. |
| REQ-3 | Validation rejects unknown fields, cycles, missing presets, invalid ISO/profile/style/mode combinations, unsafe values and output collisions before writing assets. |
| REQ-4 | Base styles are `outline`, `filled`, `bold-soft`, `silhouette` and `ghost`; all resolve through tokens rather than renderer forks. |
| REQ-5 | Presentation tokens cover fill, stroke, stroke width, opacity, line cap/join, marker appearance and optional advanced gradient/pattern/filter references under explicit safe bounds. |
| REQ-6 | `standalone` bakes complete presentation. `themed-inline` emits stable classes/data attributes and CSS variables with usable fallbacks, including `currentColor`. |
| REQ-7 | Marker modes are `none`, `capital`, `all-capitals` and `custom`; default is `none`; custom and registry selection are validated and overrideable per country. |
| REQ-8 | Animation is opt-in, subtle, transform/opacity-only, hook-based, disabled for `card` by default and suppressed by `prefers-reduced-motion`. |
| REQ-9 | SVG has deterministic viewBox/path precision/order and contains no script, event handler, external URL, embedded raster, editor metadata or redundant markup. |
| REQ-10 | Generation can target one ISO, an explicit set or the full catalog; it stages outputs, validates all selected assets, then publishes atomically with CTR-006. |
| REQ-11 | Normal execution is offline and resolves embedded corpus, presets and schemas without Node, Python, GDAL or auxiliary files. |
| REQ-12 | Exit classes and JSON diagnostics distinguish usage/config, data, rendering, validation, budget and filesystem failures and name remediation context. |
| REQ-13 | Config and CLI expose the general P2 layout contract: `tight` accepts a rendered long-side or maximum-box constraint and derives natural viewBox proportions; `contain` accepts any safe finite positive frame, preserves one uniform scale and centers unused space. `card` and `hero` are overrideable presets, not hard-coded size or aspect modes. |
| INV-1 | A successful batch is deterministic and atomically published; a failed batch leaves the previous successful output unchanged. |

## Interfaces data and behavior

Example configuration shape:

```yaml
schema: country-map/v1
extends: site-default
profile: card
layout:
  mode: tight
  longSide: 160
delivery: themed-inline
style: bold-soft
tokens:
  fill: var(--map-fill, currentColor)
  fillOpacity: 0.14
  stroke: var(--map-stroke, currentColor)
countries:
  US:
    profile: hero
    layout: { mode: contain, width: 720, height: 420 }
    marker: { mode: capital }
```

Precedence is embedded defaults → named preset ancestry → document globals →
profile block → country override → allowed CLI selection flags. `explain` emits
the resolved value and origin at every layer.

CTR-004 owns command names, exit classes and JSON envelopes. CTR-005 owns SVG
hooks such as root `.country-map`, `[data-country]`, `.country-map__shape`,
`.country-map__marker`, plus versioned `--country-map-*` custom properties.
CTR-006 records selected inputs, versions, bytes and digests. CSS hooks are a
public compatibility surface.

## Dependencies and handoffs

Depends on P2/CTR-002 and the P1 corpus identity. Supplies generated assets,
CTR-004–CTR-006 to P4 and the binary contract to P5. Owns BND-004 through BND-006.
No client-site framework dependency is permitted.

## Profiles and concerns

- Go CLI: built-binary E2E, `GOWORK=off`, testscript breadth and one ordered loop.
- Agent UX: flags/JSON are complete; no prompt may be required.
- Accessibility: safe titles/descriptions, reduced motion and no semantic reliance
  on decorative maps.
- Performance: minimal direct serialization, bounded advanced effects and no
  optimizer runtime.
- Public contract: config schema, JSON, manifest and CSS hooks are versioned.

## Validation impact

| ID | Protects | Guard/action | Tier and teeth | Owner/due |
| --- | --- | --- | --- | --- |
| VAL-1 | REQ-1, REQ-12 | enumerate runnable command leaves and require testscript coverage for help, happy path, invalid args and JSON errors | Go-profile E2E; new uncovered command reddens gate | P3 / W3 complete |
| VAL-2 | REQ-2, REQ-3 | exercise precedence, unknown fields, cycles and invalid combinations | unit+testscript; mutations fail before output | P3 / W3 complete |
| VAL-3 | REQ-4, REQ-5, REQ-6, REQ-7, REQ-8 | golden each style/delivery/marker mode and inspect hook/token resolution | integration; missing style/mode/hook or unsafe advanced reference fails | P3 / W3 complete |
| VAL-4 | REQ-9 | parse emitted XML/SVG and inject forbidden nodes/attributes/references | structural gate; every forbidden mutation fails | P3 / W3 complete |
| VAL-5 | REQ-10, REQ-11 | run `init → validate → generate → inspect` against one shared temp tree with network unavailable | Go-profile ordered binary E2E; asserts outputs, manifest, exit and no network | P3 / W3 complete |
| VAL-6 | REQ-10 | force one selected-country failure during batch publication | E2E; prior output remains byte-identical and staging is diagnosed | P3 / W3 complete |
| VAL-7 | all | build tested executable with `GOWORK=off` | build gate; ambient workspace cannot satisfy dependencies | P3 / W3 complete |
| VAL-8 | REQ-2, REQ-3, REQ-13 | generate wide, tall and near-square countries with natural long-side sizing, then tiny/huge portrait/landscape/square fixed frames | config+CLI E2E; exact viewBox/layout diagnostics, uniform-scale and no-distortion assertions; invalid dimensions fail before output | P3 / W3 complete |

## Acceptance criteria

| ID | Covers | Criterion | Evidence |
| --- | --- | --- | --- |
| AC-1 | US-1, REQ-1, REQ-2, REQ-3, REQ-12 | An agent can initialize, validate, explain and generate using JSON/flags only; invalid input is typed and actionable. | command-matrix E2E |
| AC-2 | US-2, REQ-2, REQ-4, REQ-5, REQ-6, REQ-7, REQ-13 | One portable preset expresses the owner's brand; a country override can select natural long-side sizing or any safe containing frame; only declared dimensions change and geometry is never stretched. | config/golden fixtures |
| AC-3 | US-3, REQ-6, REQ-8 | Inline output reacts to host CSS tokens and respects reduced motion without path edits or JS. | P3 hook fixture; final browser proof P4 |
| AC-4 | US-4, REQ-9, REQ-10, REQ-11 | Built binary generates selected assets and a deterministic manifest twice while offline. | ordered E2E digest receipt |
| AC-5 | US-5, REQ-8, REQ-9 | Structural validator rejects unsafe SVG and accessible/decorative metadata modes serialize predictably. | mutation receipt |

## Rollout rollback and operations

W3 delivers the first local binary and representative outputs. Schemas and CSS
hooks start at v1; additive fields require safe defaults, breaking changes require
a major. Output manifests make regeneration/rollback explicit. The CLI prints no
telemetry and reads no remote version service. `preview` is local and must not
become a hidden web server requirement.

## Delivery constraints

Use the standard library and maintained Go packages before custom plumbing. The CLI
is flags/JSON-first; because no interactive surface is required, do not take a
`huh` dependency for MVP. Do not embed arbitrary user SVG/CSS/JS, post-process with
SVGO, shell out to external runtimes, or duplicate geometry per style. Keep the
serializer deliberately small and deterministic.

## Decisions and unresolved questions

DEC-003 is accepted and binding. Stable public surfaces listed here require
DEC-001 elevated review to break. No unresolved question blocks P3.

## Amendments

AM-001 is accepted and binding: natural-aspect `tight` is the default layout;
arbitrary width-by-height uses explicit distortion-free `contain`.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
