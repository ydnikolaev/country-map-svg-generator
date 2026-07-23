---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-003"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-003 — Go renderer over a versioned geometry corpus

## Context

The owner prefers Go and requires a single offline binary for colleagues. D3
offers the strongest reference behavior for projection and SVG paths, but a
Node runtime would violate the accepted installation contract. Upstream geometry
updates are rare; style and package generation are frequent.

## Decision

Ship a Go CLI that renders from an embedded, immutable, versioned geometry
corpus. Separate source adapters and corpus compilation from ordinary asset
generation. Use a maintained Go geometry library for GeoJSON and simplification;
keep centered equal-area projection and minimal SVG serialization as small
isolated modules verified against pinned D3 golden fixtures. Do not require
Node, Python, GDAL, or network access at generation time.

## Options considered

| Option | Fit |
| --- | --- |
| Node/D3/SVGO runtime | Fastest prototype and canonical APIs, but violates single-binary/offline distribution |
| Go plus embedded JavaScript engine | Preserves D3 behavior but adds runtime complexity and a large hidden compatibility surface |
| Go plus versioned corpus and narrow verified geometry modules | Meets distribution, determinism, performance, and maintainability constraints |

## Rejected alternatives

Node is rejected as a user runtime dependency. An embedded JavaScript engine is
rejected because the accepted operations need only a narrow subset of D3 and do
not justify a second runtime. Generic post-hoc SVG optimization is rejected:
the renderer emits minimal deterministic markup directly.

## Consequences and tradeoffs

The shipping binary is simple to install and generation is deterministic. The
project owns a small projection/serialization surface and must protect it with
golden fixtures, topology/validity checks, and full-catalog budgets. Upstream
refresh is a separate maintainer path and cannot block normal users.

## Affected artifacts and owners

P1 owns source adapters and corpus identity. P2 owns projection,
simplification, smoothing, validation, and D3 parity. P3 owns Go CLI and SVG
serialization contracts. P4 owns full-catalog/browser proof. P5 owns binary
matrix and handoff.

## Validation and revisit trigger

Revisit if the Go pipeline cannot meet the accepted full-catalog visual or byte
budgets without reproducing a broad geospatial framework, or if a supported
single-binary runtime exposes canonical D3/SVGO behavior with lower complexity.
Until the full-catalog gate passes, performance budgets remain targets.

## Supersession

No predecessor. A successor must preserve offline deterministic generation or
obtain new owner authority.
