---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-005"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["DEC-004", "RUN-001-RESULT", "RUN-002-RESULT", "RUN-004-RESULT", "RUN-005-RESULT", "P1/CTR-001", "P2/CTR-002"]
---
# DEC-005 — Bounded q-aware source fallback

## Context

DEC-004 selected two maintainer-precomputed LOD tiers plus unmodified P1 as the
lossless runtime `source` tier. It rejected continued runtime RDP after the
original implementation could not jointly satisfy topology, deviation and byte
budgets. It also declared a revisit trigger when full source makes large custom
sizes impractical or two tiers cannot cover accepted scales.

RUN-004 and RUN-005 reached that trigger without weakening any guard. The
compact and standard artifacts fail AQ's fitted-space deviation at the accepted
`0.012` ceiling because they are simplified in a different coordinate space.
After both candidates fail, the exact fitted source is valid before
canonicalization but pointwise `q=0.01px` rounding creates proper crossings.
AQ's 23,994-point fitted source passes topology; its quantized probe contains
three crossings within the first 500 segments and seven within the first 1,000.
RU card independently fails after quantization at ring `0/0`, segments
`7212/7214`.

The contradiction is now narrower than DEC-004's original runtime-RDP failure.
The source authority, visual tolerance and fixed grid are all known. What is
missing is a generic way to materialize a final q-safe representation from exact
retained P1 when precomputed tiers cannot be used.

## Decision

Keep exact immutable P1 as the sole authority and input of runtime tier
`source`. Permit the pure-Go runtime to finalize that fallback through a bounded,
q-aware reduction using the existing validated geometry simplifier. This is
allowed only after compact and standard are unavailable or fail, or when
effective fitted scale selects source.

The source finalizer:

1. begins from exact projected, fitted and retained P1 geometry;
2. reserves the maximum `q=0.01px` grid displacement (`q/sqrt(2)`) from the
   resolved symmetric visual tolerance;
3. searches monotonically from the coarsest allowed simplification toward finer
   candidates without ever exceeding the remaining tolerance;
4. accepts a candidate only after exact `q=0.01` canonicalization, final topology
   and winding validation, protected/minimum-part retention, and final symmetric
   full↔canonical deviation within the caller's resolved tolerance;
5. serializes only the accepted canonical geometry and applies the unchanged
   generic hard byte limit;
6. fails typed when no candidate satisfies every guard.

`source` provenance continues to identify the authoritative tier, not a claim
that final SVG commands contain every P1 vertex. Add explicit finalization
provenance: exact P1 identity, `source_reduced`, attempted/selected tolerance,
q, final deviation, restorations, output points and fallbacks. A caller and an
auditor can distinguish an unreduced exact-source attempt from its bounded final
representation.

This decision does not authorize a general runtime repair layer. `MakeValid`,
buffer/union, polygon clipping, ad-hoc vertex movement, relaxed crossing rules,
country-specific branches, lower precision, higher visual tolerance, silent
budget overflow and post-hoc SVG optimization remain forbidden. Compact and
standard stay maintainer-precomputed and primary.

## Options considered

| Option | Result |
| --- | --- |
| Keep pointwise-rounded exact source only | Rejected: verified AQ and RU inputs become topologically invalid at mandatory q |
| Add a third precomputed tier or presimplification hierarchy | Deferred: DEC-004 names it as the next option, but the measured source-finalization defect is narrower |
| Implement custom snap-rounding/untangling | Rejected: it is a new geometry repair engine and can change polygon identity without a visual-error proof |
| Use bounded q-aware existing Go simplification only for source finalization | Selected: reuses the validated pure-Go path, keeps exact P1 authority and subjects final bytes to every existing guard |
| Raise tolerance, change q or relax budgets | Rejected: each changes a user-visible or performance contract and does not prove topology |

## Rejected alternatives

The original failed runtime RDP loop simplified first and tried to recover by
changing tolerance or precision. The selected finalizer has a different
acceptance boundary: q, topology, protection and final symmetric deviation are
one indivisible predicate, and the search can only move finer. No result exists
until the final representation passes.

Mapshaper `-clean`, `MakeValid` and polygon union are inappropriate at runtime:
they can split, merge or fill regions, making provenance and protected-feature
identity harder to prove. Exact source with more decimal places violates the
canonical SVG grid. A third tier adds artifact schema and calibration before the
fallback path itself has been tested.

## Consequences and tradeoffs

Normal card and hero generation still uses compact or standard and retains the
precomputed performance benefit. Source fallback pays additional CPU only when
selected or reached after a failed candidate. Runtime remains one offline Go
binary with no Node, npm, Mapshaper, network, CGO or new dependency.

Final source commands are visually bounded rather than vertex-lossless.
Explicit finalization provenance makes that tradeoff inspectable. Exact P1
continues to determine projection, fit, retention, protected anchors, marker
transform, deviation authority and artifact binding.

The finalizer may still exceed a byte budget or exhaust its valid search. That is
a typed product limitation and a new measured planning input, never permission
to weaken the guard.

## Affected artifacts and owners

P1, CTR-001 and the corpus bytes remain unchanged and authoritative. P2 owns the
source-finalization algorithm, provenance, tests and additive CTR-002 fields.
P2's projected compact/standard artifact work is independent but must use the
same final-representation guard. P3 consumes canonical commands and provenance;
P4 measures browser performance; P5 continues to ship only Go binaries.

Every successor P2 plan binds DEC-004, this decision, both accepted layout
amendments, P1 RUN-001 and RUN-002 results, exact P1 manifest/geometry hashes,
and the current verified product checkpoint.

## Validation and revisit trigger

Before publication, force source for AQ card/hero and RU card. Run first without
a byte limit to isolate representation validity, then with `2500/8000`. Record
attempted and selected tolerance, q, final deviation, topology, protection,
points, bytes, parser round-trip and runtime. Mutations must fail if search moves
coarser, skips final q validation, reports zero/unmeasured deviation, loses a
protected part, changes q, uses a country branch or accepts an invalid crossing.

The all-corpus matrix covers both profiles, card/hero reference sizes, `1024`,
`2048`, `19x19` and `10000x12345`, with deterministic repeated output and
ordinary offline `make check`.

Revisit if source finalization cannot pass the unlimited representation brake,
if it is selected frequently enough to violate generation targets, if it cannot
stay within hard bytes at accepted scales, or if its provenance cannot map
exactly to P1. The next governed options are Mapshaper presimplification
importance metadata or an additional precomputed tier.

## Supersession

Refines DEC-004 only where it calls source lossless and rejects all runtime RDP.
DEC-004 remains binding for two primary precomputed tiers, scale-only selection,
full-source authority, maintainer-only Mapshaper, pure-Go offline runtime,
generic fallbacks and hard validation. DEC-003 remains binding except where
already refined by DEC-004 and this decision.
