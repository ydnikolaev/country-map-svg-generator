---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-004"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["DEC-003", "RUN-001-RESULT", "P1/CTR-001", "P2/CTR-002"]
---
# DEC-004 — Precomputed LOD with a pure-Go runtime

## Context

DEC-003 selected a maintained Go simplifier because the owner requires one
offline Go binary for ordinary generation. P2 RUN-001 exercised DEC-003's revisit
trigger: Natural Earth 10m geometry for `CN/un/card` cannot simultaneously keep
the resolved `0.75px` deviation, fixed `0.01px` output grid, valid topology and
the card byte target through the selected runtime RDP path. The observed lower
bounds were 83,827 linear path bytes after fixed-grid topology normalization and
292,423 bytes after finer `0.001px` quantization. Raising the tolerance to
`1.25px` made the case small but violated the resolved quality contract.

A bounded spike with Mapshaper 0.7.44 simplified the same whole China
MultiPolygon deterministically from 12,623 source vertices to 90 vertices at
`resolution=48`, while its documented simplifier repaired intersections.
This decision governs only where cartographic generalization occurs. Projection,
natural-aspect fitting, arbitrary `contain` frames, styling, SVG delivery and the
single-binary installation contract do not change.

## Decision

Use committed, versioned P2-owned levels of detail derived from the immutable P1
corpus as the primary runtime reduction input. A pinned maintainer-only
Mapshaper recipe produces two candidate tiers, `compact` and `standard`; the
unmodified P1 geometry is the lossless `source` tier. The Go runtime selects a
tier only from the effective fitted geometry scale, validates its source binding,
topology, maximum deviation, protected anchors and minimum parts, and advances to
a finer tier or `source` on any failure.

Mapshaper, Node and npm are permitted only in the explicit LOD rebuild workflow.
They are forbidden from the shipped binary, ordinary CLI execution and the
normal offline `make check` gate. Generated LOD bytes, their source identity,
recipe, tool integrity, package lock, thresholds, removals and restoration
provenance are committed and verified in Go. A source, recipe, tool or threshold
change creates a new LOD version rather than mutating accepted bytes.

`card` and `hero` remain versioned data presets over the general engine. Neither
preset names nor country identifiers may select an LOD. Softening remains
optional and bounded; when it would exceed the generic byte budget or fail
flattened-topology validation, the result uses deterministic linear commands and
records the fallback.

## Options considered

| Option | Result |
| --- | --- |
| Continue runtime RDP with tolerance or precision fallback | Rejected by the measured topology, deviation and byte contradictions |
| Add GEOS/CGO or a Node runtime | Rejected because it breaks the portable offline single-binary contract |
| Store Mapshaper presimplification importance/TopoJSON arcs | Deferred; it adds a Go arc decoder, threshold semantics and harder P1 provenance mapping before they are shown necessary |
| Commit two derived LOD tiers plus full-source fallback | Selected as the smallest framework-owned, deterministic runtime seam |
| Hand-edit or branch on exceptional countries | Rejected because it fabricates a non-general policy and cannot scale to the corpus |

## Rejected alternatives

Adaptive `0.001px` precision and fixed-grid repair preserve topology only by
retaining an unusably large path. Automatically raising RDP tolerance changes
the caller's resolved quality contract. Buffer/union cleanup after the failing
RDP leaves thousands of unnecessary vertices. Per-component `keep-shapes`
precomputation is also rejected: China alone has 63 parts, which creates an
artificial minimum vertex floor. Each geometry ID is generalized as one whole
MultiPolygon; only anchor-bearing or required source components may be restored
generically after comparison with P1.

## Consequences and tradeoffs

Ordinary generation becomes faster and its output complexity is bounded before
serialization. The project gains deterministic cartographic generalization
without placing a second runtime on users or colleagues. Arbitrary `tight` and
`contain` dimensions remain flexible because scale selects the tier; large or
unsafe cases fall back to the full source.

The repository gains committed derived data, a pinned Node development toolchain
and a rebuild-only command. The first successor-plan task must measure artifact
size, rebuild time and final Go path bytes before the full corpus is generated.
The LOD corpus is cache-like but governed: it is never authoritative over P1 and
must fail closed on stale identity, invalid topology, excessive deviation or lost
identity features.

## Affected artifacts and owners

P1 and CTR-001 remain unchanged and authoritative. P2 owns the LOD recipe,
artifacts, runtime selection, restoration, provenance and validation additions to
CTR-002. P3 consumes the additive CTR-002 provenance and continues to own SVG
serialization and CLI configuration. P4 proves full-catalog browser performance.
P5 ships only Go binaries and excludes the maintainer toolchain.

The successor P2 plan must bind DEC-003, this decision, both accepted layout
amendments, the RUN-001 result, exact P1 corpus hashes and the verified product
checkpoint. It must not resume RUN-001 under changed HOW.

## Validation and revisit trigger

Before full implementation, a bounded spike must compare `compact`, `standard`
and `source` across the representative corpus at the declared rendered scales.
It stops if any card exceeds 2,500 path bytes, any hero exceeds 8,000 path bytes,
protected restoration or exact removal provenance cannot be proved, visual
approval fails, committed artifact size/runtime is unreasonable, or two clean
rebuilds differ byte-for-byte. Normal `make check` must validate committed
artifacts without Node, npm, Mapshaper or network access.

Revisit if two tiers cannot cover the accepted scales without excessive fallback,
if the full-source tier makes large custom sizes impractical, or if LOD provenance
cannot be mapped exactly to P1 geometry. Presimplification importance metadata is
the next governed option; country branches, silent tolerance increases and
weaker topology gates are not.

## Supersession

Refines DEC-003 at its documented revisit trigger. DEC-003 remains binding for
the Go renderer, offline runtime, immutable corpus and no post-hoc SVG optimizer;
its maintained-Go-simplifier clause is replaced only for the maintainer-time LOD
derivation described here.
