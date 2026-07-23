---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-001"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-001 — Elevated risk policy for the map generator

## Context

The epic introduces a public CLI/config contract, a shared geometry corpus,
performance budgets, cross-platform output, and politically sensitive boundary
profiles. The registered policy must require independent review of architecture,
performance, accessibility, and shared-design decisions without treating this
local, rebuildable generator as irreversible infrastructure.

## Decision

Use the registered `elevated` risk policy. Every specification must provide
deterministic positive and negative validation, the authoritative project gate,
and two blinded independent review lanes before closure. A security, privacy, or
irreversible-data change requires an amendment and reclassification to
`critical`; an external commitment requires `external`.

## Options considered

| Option | Fit |
| --- | --- |
| ordinary | Too weak for the new public CLI/config contract and explicit performance/accessibility concerns |
| elevated | Matches architecture, public-contract, shared-design, accessibility, performance, and distribution triggers |
| critical | Excessive while generation is local, outputs are rebuildable, and no secrets or irreversible state exist |

## Rejected alternatives

`ordinary` is rejected because one review lane cannot independently cover both
geospatial correctness and client/performance behavior. `critical` is rejected
because the accepted scope has no security, privacy, or irreversible
infrastructure boundary.

## Consequences and tradeoffs

Delivery pays for two independent reviews at each spec boundary and cannot
degrade completion when a lane is missing. Deterministic gates remain the inner
loop. Review effort is concentrated on corpus policy, visual balance, SVG
contracts, accessibility, and performance rather than duplicating compiler or
schema checks.

## Affected artifacts and owners

P1 owns corpus provenance and political-profile evidence. P2 owns geometry and
visual-validity evidence. P3 owns CLI, config, output security, and deterministic
generation. P4 owns full-catalog, browser, accessibility, and performance proof.
P5 owns cross-platform packaging. The teamlead owns escalation.

## Validation and revisit trigger

The policy is valid while the tool remains local-first, generated output is
rebuildable, and no secret or irreversible state is introduced. Revisit on an
accepted amendment adding remote execution, uploads, untrusted raw SVG,
credentials, destructive migration, or contractual external publication.

## Supersession

This epic-local decision has no predecessor. A successor must cite the exact
amendment or accepted finding that changed the registered trigger set.
