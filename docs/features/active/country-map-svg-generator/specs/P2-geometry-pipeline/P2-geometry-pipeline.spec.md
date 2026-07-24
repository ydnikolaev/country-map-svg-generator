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
inputs: ["DISC-006", "ARCH-001", "DEC-003", "DEC-004", "DEC-005", "DEC-006", "DEC-007", "DEC-008", "DEC-009", "AM-001", "AM-002", "AM-003", "AM-004", "AM-005"]
---
# P2 — Soft-Organic Geometry Pipeline

## Outcome

Every corpus entity can be transformed into deterministic, visually balanced,
presentation-free geometry for `hero` and `card`. The result preserves country
identity at its intended size, softens mechanical edges without fabricating
geography, derives a stable natural-aspect viewBox by default, can opt into a
fixed containing frame without distortion, and supplies optional projected marker
positions.

## User stories and stakeholders

| ID | Story |
| --- | --- |
| US-1 | As a visitor, I see a recognizable, elegant country silhouette rather than noisy GIS detail. |
| US-2 | As a site developer, I receive geometry in the country's natural projected proportions by default, can request any fixed containing frame when layout requires it, and do not tune individual SVG paths. |
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
| REQ-1 | Geometry is centered with a per-entity equal-area projection and uniformly fitted without distortion. Default `tight` layout derives the viewBox aspect ratio from projected geometry plus declared padding; explicit `contain` layout accepts any safe finite `width × height`, preserves geometry proportions, centers it and leaves unused frame area transparent. |
| REQ-2 | `card` and `hero` are configurable named presets over the general engine, not engine modes. Their defaults specify intended rendered long-side scale, padding, tolerances and byte budgets; callers may use arbitrary natural or fixed-frame sizes. |
| REQ-3 | Simplification preserves ring validity, winding semantics, non-empty identity geometry and P1 protected features. |
| REQ-4 | Softening is conservative, bounded and deterministic; it cannot move a point beyond a declared profile tolerance or introduce self-intersection. |
| REQ-5 | Tiny unprotected features may be removed by explicit area/visibility policy; removals are inspectable per entity. |
| REQ-6 | Per-country overrides can adjust orientation, padding, tolerance, retained features and marker offset without altering source corpus bytes. |
| REQ-7 | Capital/custom coordinates project through the same transform as geometry and are rejected or warned when implausibly outside the fitted result. |
| REQ-8 | Equal inputs produce byte-equivalent normalized path commands independent of map iteration or host. |
| INV-1 | The geometry result is presentation-free and cannot depend on fill, stroke, CSS or delivery mode. |

## Interfaces data and behavior

Input is one P1 entity/profile record plus a geometry profile, a layout request,
and resolved per-country overrides. Layout is either `tight`, with a long-side or
maximum-box constraint from which the natural viewBox dimensions are derived, or
`contain`, with an explicit finite positive frame. Both apply one uniform scale
and never stretch geometry. Output is CTR-002: normalized path data, deterministic
viewBox, natural projected aspect ratio, fit transform, component/ring metadata,
resolved layout mode, effective rendered geometry scale, removal diagnostics and
optional marker points. It contains no colors, classes, XML or CSS.

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
| VAL-3 | REQ-1, REQ-2, REQ-5 | run hero/card fixtures at intended rasterized long-side sizes; compare an elongated geometry in multiple `contain` frames | visual+structural; removal report and recognizable approved baselines; resolved quality remains a function of fitted geometry scale rather than the frame's shorter side | P2 / W2 complete |
| VAL-4 | REQ-6 | apply bounded valid and invalid country overrides | unit; valid changes only named dimension, out-of-range fails | P2 / W2 complete |
| VAL-5 | REQ-7 | project inside, edge and implausible-outside markers | integration; same transform and typed anomaly behavior | P2 / W2 complete |
| VAL-6 | REQ-2, REQ-3, REQ-8 | prove derived-candidate acceptance is the frozen silhouette oracle, DEC-007 protected visibility and the byte caps: reintroducing a raw-deviation gate to the derived path, bypassing the oracle, or restoring source components must each redden a mutation; re-verify every committed ladder row offline in pure Go against the frozen thresholds and bound digests, reddening on a perturbed coordinate, changed threshold, stale corpus/oracle/recipe identity, reordered candidates or suppressed omission; keep raw deviation as the gate on the explicit DEC-005 source path only | unit+property+mutation; oracle-gated derived acceptance, committed-artifact re-verification, and explicit-source integrity as three separate obligations | P2 / W2 complete |

Per AM-005 and DEC-006, VAL-6 guards that automatic derived-candidate acceptance
is exclusively the frozen silhouette oracle, DEC-007 protected visibility and the
byte caps, and that the derived path never consults raw full-coastline symmetric
deviation. Raw deviation remains binding for an explicit `source` or
custom-quality request under DEC-005, where skipping final validation fails. Per
DEC-009 the identity candidate is a fallback tried only after every declared
resolution rung fails, and a band whose candidates cannot satisfy topology at the
actual fitted scale emits a typed reason and no artifact rather than a silent gap.

The acceptance metric on the explicit-source path is a clamped predicate: exact at
or below its limit and short-circuiting above, so it is a sound pass/fail test and
not a measurement. Any deviation figure carried in provenance, evidence or a plan
is either at or below its limit or produced by an unclamped path; a first-crossing
value is never reported as a measurement.

Population guards under VAL-1 through VAL-3 report every failing entity with
counts and distribution; a sweep that stops at its first failure does not satisfy
them, because the first failure is a function of iteration order and hides both
the magnitude and every later assertion.

## Acceptance criteria

| ID | Covers | Criterion | Evidence |
| --- | --- | --- | --- |
| AC-1 | US-1, REQ-2, REQ-3, REQ-4, REQ-5 | Approved representative card/hero silhouettes remain recognizable, balanced and free of visible topology artifacts. | raster comparison sheet plus structural receipt |
| AC-2 | US-2, REQ-1, REQ-2, REQ-8 | Every representative `tight` output follows the projected country aspect ratio; every `contain` output fits its arbitrary frame without distortion; quality follows effective fitted scale rather than unused frame space; repeated runs are byte-equivalent. | golden receipt |
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

### AM-001 — Natural-aspect layout and explicit fixed frames

Accepted owner clarification on 2026-07-23: the default SVG viewBox follows the
projected geometry's own proportions (wide Russia, tall Chile/Argentina,
near-square Australia). `card` and `hero` provide scale/quality defaults rather
than fixed aspect ratios. An explicit arbitrary `width × height` remains available
as a containing frame and must preserve proportions. Quality is calibrated from
the actually rendered geometry scale, not the frame's shorter side.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
