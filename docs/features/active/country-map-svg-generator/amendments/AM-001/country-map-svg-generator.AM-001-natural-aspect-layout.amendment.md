---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "amendment"
id: "AM-001"
epic: "country-map-svg-generator"
status: proposed
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# AM-001 — Amendment

## Trigger and discovery evidence

During P2 planning on 2026-07-23, before any geometry implementation began, the
owner clarified that an SVG canvas should normally follow the visual proportions
of the country: Russia is wide, Chile and Argentina are tall, and Australia is
near-square. A fixed square `card`/`hero` canvas would add unnecessary transparent
area and conceal an incorrect coupling between use-scale presets and aspect ratio.

The same discussion retained the previously accepted flexibility requirement:
agents and producers must still be able to request arbitrary dimensions such as
`19x19` or `10000x12345`. Those dimensions describe an explicit layout frame,
not permission to stretch the geography.

## Affected scope and authority

This is a P2 WHAT/macro-HOW clarification to CTR-002 layout semantics. The product
owner is the decision authority. It affects P2 REQ-1, REQ-2, US-2, AC-2, interface
text, and validation; it invalidates the existing P2 readiness evidence and the
unaccepted PLAN-001 proposal.

P3 must expose the resulting generic layout options without hard-coded country or
preset branches. P4 must document natural aspect handling and fixed-frame
containment. Their outcomes and boundaries remain unchanged, so they require
fresh planning context rather than a separate specification amendment now.

## Prior accepted truth

P2 specification SHA-256
`9cfdf6bff8a5b49559fbbcfd0679517891b27087c7d587b2fbd51646760aa74f`
required a deterministic profile viewBox with padding and described `card` and
`hero` as separate tolerance profiles, but did not state whether the viewBox aspect
was fixed or geometry-derived.

Readiness `REV-29FB24300D07#3448feebbd87c5bb0e8dca2e3cfc3ee23ccd56035a6d40099c1569c4aabd6914`
authorized planning against that byte set. Planner brief
`BRIEF-001#48241f40004011bfc560c1b9508e8c891d33e12df56a684b784e63c6f0a93ba6`
was closed as cancelled before plan acceptance because the proposal assumed a
fixed width and height.

## Accepted change

P2 defines two general layout modes:

- `tight` is the default. It uniformly scales projected geometry to a rendered
  long-side or maximum-box constraint, derives the second canvas dimension from
  the projected geometry, and adds declared padding. The output viewBox therefore
  follows the country's natural projected aspect ratio.
- `contain` is explicit. It accepts any safe finite positive width and height,
  uniformly fits and centers the geometry inside the padded frame, and leaves
  unused area transparent. It never stretches or independently scales axes.

`card` and `hero` become versioned data presets over this API. Their defaults use
`tight` and define intended rendered long-side scale, padding, tolerances, and byte
budgets—not a fixed aspect ratio. Callers may override every field or choose
`contain`.

Quality calibration uses the long side of the actually fitted geometry in target
pixels. It does not use the frame's shorter side, which would over-simplify
naturally elongated countries.

## Downstream impact map

- **P2:** invalidate prior readiness; validate amended spec bytes; issue fresh
  readiness; produce and accept a new plan before implementation.
- **CTR-002/P3:** carry layout mode, natural projected aspect ratio, resolved
  viewBox, uniform fit transform, and effective rendered scale. P3 CLI/config
  planning must expose natural size and explicit frame controls through the same
  engine contract.
- **P4:** full-catalog and browser evidence must include wide, tall, and
  near-square natural viewBoxes plus arbitrary `contain` frames. Client guidance
  must distinguish intrinsic SVG aspect ratio from CSS slot dimensions.
- **P1/P5:** no contract or implementation impact.

## Validation delta

P2 adds deterministic tests that:

- wide (`RU`), tall (`CL`/`AR`), and near-square (`AU`) `tight` outputs derive
  materially different viewBox ratios from projected bounds;
- changing only target long-side scale preserves aspect ratio;
- `contain` accepts tiny, huge, portrait, landscape, and square frames, preserves
  one uniform scale, centers unused axes, and never distorts geometry;
- `card` and `hero` resolve through the generic layout contract and no engine enum
  or country switch hard-codes their behavior;
- quality thresholds respond to fitted geometry scale rather than the frame's
  shorter dimension.

## Migration rollout and rollback

No runtime or generated asset exists, so migration is documentation and planning
only. Apply the amendment to P2 before implementation, replace the candidate plan,
and retain the cancelled planner result as historical evidence. Rollback would
reject this amendment and restore the prior specification bytes before any P2
run; after implementation, changing aspect semantics requires a new amendment and
golden/full-catalog comparison.

## Application and status

Pending governed proposal, impact fencing, owner-authority acceptance, application
to the P2 specification, and fresh readiness/plan evidence. No implementation run
may begin while this amendment is unapplied.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
