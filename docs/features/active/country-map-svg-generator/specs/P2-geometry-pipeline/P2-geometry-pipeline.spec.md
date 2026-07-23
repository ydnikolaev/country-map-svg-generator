---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "specification"
id: "P2"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["DISC-006", "ARCH-001", "DEC-003"]
---
# P2 — Soft-Organic Geometry Pipeline

## Outcome

Every corpus entity can be transformed into deterministic, visually balanced,
presentation-free geometry for `hero` and `card`. The result preserves country
identity at its intended size, softens mechanical edges without fabricating
geography, fits a stable viewBox, and supplies optional projected marker positions.

## User stories and stakeholders

| ID | Story |
| --- | --- |
| US-1 | As a visitor, I see a recognizable, elegant country silhouette rather than noisy GIS detail. |
| US-2 | As a site developer, I receive geometry already fitted for its use scale and do not tune individual SVG paths. |
| US-3 | As a producer, I can override exceptional countries without forking the pipeline. |
| US-4 | As a maintainer, I can compare Go output with a trusted projection/path oracle and detect topology damage. |

## Scope and non-goals

In scope: centered equal-area projection, orientation/fitting, profile-specific
simplification, conservative smoothing, ring normalization, protected-feature
retention, optional marker projection, per-country geometry overrides, structural
validation and golden comparison.

Out of scope: SVG styling/markup, capital selection policy, labels, borders between
subdivisions, pan/zoom maps, geographic analysis and the deferred `detail` profile.

## Context and ground truth

The owner's screenshot fixes two decorative scales: repeated small cards and one
large hero. Small islands and coast detail may be removed when they do not carry
identity at that scale. The desired character is slightly rounded and
hand-finished, not cartographically jagged. D3's azimuthal equal-area projection,
fit and path output are the pinned reference. DEC-003 keeps the Go implementation
narrow and independently testable.

## Requirements and invariants

| ID | Requirement |
| --- | --- |
| REQ-1 | Geometry is centered with a per-entity equal-area projection and fitted into a deterministic profile viewBox with declared padding. |
| REQ-2 | `card` and `hero` use independently configurable tolerances; defaults target their accepted CSS sizes and byte budgets. |
| REQ-3 | Simplification preserves ring validity, winding semantics, non-empty identity geometry and P1 protected features. |
| REQ-4 | Softening is conservative, bounded and deterministic; it cannot move a point beyond a declared profile tolerance or introduce self-intersection. |
| REQ-5 | Tiny unprotected features may be removed by explicit area/visibility policy; removals are inspectable per entity. |
| REQ-6 | Per-country overrides can adjust orientation, padding, tolerance, retained features and marker offset without altering source corpus bytes. |
| REQ-7 | Capital/custom coordinates project through the same transform as geometry and are rejected or warned when implausibly outside the fitted result. |
| REQ-8 | Equal inputs produce byte-equivalent normalized path commands independent of map iteration or host. |
| INV-1 | The geometry result is presentation-free and cannot depend on fill, stroke, CSS or delivery mode. |

## Interfaces data and behavior

Input is one P1 entity/profile record plus a geometry profile and resolved
per-country overrides. Output is CTR-002: normalized path data, deterministic
viewBox, component/ring metadata, removal diagnostics and optional marker points.
It contains no colors, classes, XML or CSS.

Failures distinguish invalid source geometry, projection failure, topology damage,
empty result, protected-feature loss, marker anomaly and exceeded override bounds.

## Dependencies and handoffs

Depends on P1/CTR-001. Hands CTR-002 to P3 and diagnostic metadata to P4. P2 owns
BND-003 exclusively. D3 fixtures are pinned development evidence, not a runtime
dependency.

## Profiles and concerns

- Go: narrow geometry modules and standard Go tests.
- Visual quality: recognizable silhouettes at explicit rendered sizes.
- Performance: point count and path precision are decided before serialization.
- Correctness: geometry simplification can create invalid shapes and must be guarded.

## Validation impact

| ID | Protects | Guard/action | Tier and teeth | Owner/due |
| --- | --- | --- | --- | --- |
| VAL-1 | REQ-1, REQ-8 | compare representative continents, islands, antimeridian and polar entities with pinned D3 fixtures | integration golden; bounded numeric delta and byte-stable rerun | P2 / W2 complete |
| VAL-2 | REQ-3, REQ-4 | feed self-intersection, winding, protected-island and near-collapse mutations | unit+integration; each invalid mutation fails or preserves named feature | P2 / W2 complete |
| VAL-3 | REQ-2, REQ-5 | run hero/card fixtures at intended rasterized sizes | visual+structural; removal report and recognizable approved baselines | P2 / W2 complete |
| VAL-4 | REQ-6 | apply bounded valid and invalid country overrides | unit; valid changes only named dimension, out-of-range fails | P2 / W2 complete |
| VAL-5 | REQ-7 | project inside, edge and implausible-outside markers | integration; same transform and typed anomaly behavior | P2 / W2 complete |

## Acceptance criteria

| ID | Covers | Criterion | Evidence |
| --- | --- | --- | --- |
| AC-1 | US-1, REQ-2, REQ-3, REQ-4, REQ-5 | Approved representative card/hero silhouettes remain recognizable, balanced and free of visible topology artifacts. | raster comparison sheet plus structural receipt |
| AC-2 | US-2, REQ-1, REQ-8 | Every representative output fits its viewBox and repeated runs are byte-equivalent. | golden receipt |
| AC-3 | US-3, REQ-6 | Named exceptional countries are expressible through bounded data overrides, not code forks. | override fixtures |
| AC-4 | US-4, REQ-3, REQ-4 | Mutation tests prove invalid topology and protected-feature loss cannot pass silently. | teeth receipt |

## Rollout rollback and operations

P2 first establishes representative golden fixtures, then P4 exercises the entire
catalog. Geometry-profile defaults are versioned. A changed default requires a new
golden/full-catalog comparison; rollback selects the prior profile/corpus pair.
Diagnostics expose input/output point counts, removed features, projection and
override choices.

## Delivery constraints

Use a maintained Go geometry library for interchange/simplification where it owns
the concern; keep custom projection and SVG-adjacent normalization isolated and
small. Do not introduce Node, Python, GDAL or network runtime requirements. Do not
run generic post-hoc SVG optimization to compensate for bloated geometry. Do not
manually redraw countries in code; exceptional treatment is bounded data.

## Decisions and unresolved questions

DEC-003 is accepted and binding. The full-catalog evidence may tune default
tolerances inside the accepted visual/byte budgets; it may not add `detail` scope
without an amendment.

## Amendments

None.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
