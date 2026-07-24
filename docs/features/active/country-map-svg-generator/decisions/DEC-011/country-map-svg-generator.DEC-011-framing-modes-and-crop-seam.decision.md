---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-011"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-011 — Framing stays full-claim; cropping is an explicit seam, not a heuristic

## Context

T3 made the whole catalog render for the first time, and rendering it exposed a
framing question that had been invisible while these entities failed on the byte
budget. The viewBox is fitted to the whole projected source before any candidate
is chosen, so it reserves space for components the silhouette may deliberately
not draw. Chile's UN geometry spans 43.0° of longitude because of Easter Island
(-109.45), Salas y Gómez (-105.47) and the Juan Fernández group (-80.7, -80.1);
the mainland alone spans 13.5°. Rendered, Chile sits against one edge of a
near-square box. The United States shows the same effect through Alaska and
Hawaii. This is pre-existing behaviour and not a T3 regression — the Australia
viewBox is identical on the ladder and source paths (126.4×144 both).

The owner's requirement is not a cleverer default. It is an iteration loop: hand
an agent the binary, ask for the United States, look at it on the site, ask for a
different crop, look again. That is a request for **explicit, addressable
parameters**, and it changes what has to be decided. An automatic
"mainland detection" rule would have to be guessed now; explicit knobs let the
choice be made per request, with evidence on screen.

Three candidate automatic rules were considered and each breaks somewhere. A pure
distance threshold dismembers archipelago states — Indonesia is ~5000 km wide, so
Papua would be dropped. A "distance AND low pixel contribution" rule keeps
archipelagos intact and removes Easter Island, but keeps Alaska, so it does not
produce the tight United States card the owner had in mind. A "smallest box
holding 98% of silhouette mass" rule behaves the same way on Alaska. None of the
three is obviously right, and choosing between them before seeing the full
catalog would be guesswork.

**Why cropping is not merely a missing flag.** `Override` today carries rotation,
padding, quality, retention anchors, minimum parts and marker offsets — no
component selection and no frame. More importantly, the committed ladder is keyed
to the exact source geometry, so a cropped request is a geometry the ladder never
saw. Measured: dropping Alaska and Hawaii from the United States and calling
`Generate()` falls through to the DEC-005 source path and fails with
`hard_budget_exceeded` at **40885 bytes against 2500** for the card and 42868
against 8000 for the hero; cropped Chile fails at 21198 against 2500. Simplifying
on the fly is not available either — that is Mapshaper, and a pure-Go offline
runtime with no Node and no network at generation is a frozen constraint.

**What does work, measured.** Cropping the *already selected candidate* rather
than the source stays comfortably inside the frozen budgets, in pure Go:

| Case | Parts | Bytes | Band cap |
| --- | --- | --- | --- |
| US card, cropped | 7 → 4 | 1817 | 2200 |
| US hero, cropped | 17 → 9 | 3137 | 7500 |
| Chile card, cropped | 9 → 5 | 1675 | 2200 |

The candidate is already simplified and already oracle-judged, so cropping it is
cheap and cannot inflate the path. The cost is that a cropped frame displays
geometry simplified for the *previous* fitted scale: at the ~1.2× zoom measured
above this is invisible, and it grows with the crop-zoom factor.

## Decision

**1. The default framing stays the full claim.** Every committed row keeps
today's behaviour: the viewBox is fitted to the whole projected source. Nothing
about the current 993-row catalog changes.

**2. A second named framing mode is a product requirement, but its definition is
deferred until the T5 contact sheet exists.** The owner will decide what
"central" means with all 996 cards on screen rather than from three hypothetical
rules. Deferring costs nothing: no code written now would be wasted, because the
seam below is needed for either answer.

**3. Cropping is implemented as an explicit seam over the selected candidate,
never as an automatic heuristic in the default path.** The order is fixed: select
from the ladder first, then crop, then refit. Cropping before selection is
forbidden — it leaves the ladder and blows the frozen budgets, as measured above.

**4. The knobs are addressable per request, not per country.** Component
selection is expressed through the existing identity vocabulary (DEC-007 source
order, contribution rank, digest), never through a country literal or a
per-entity branch. An agent asks for the same generic thing for every entity.

**5. The acceptance rule for a cropped variant is an open sub-decision, named
here so it is not skipped.** The oracle measures IoU and recall against the full
visible reference; a crop changes what the reference should be. Whether a cropped
variant is re-judged against the cropped reference, carries its parent's verdict,
or is marked as unjudged provenance must be settled when the seam is built.

**6. This work is not part of P2.** P2 owns the geometry pipeline and its ladder,
which are complete. The crop seam plus its CLI surface belong to P3, or to a
successor spec if P3 is already too large.

## Options considered

1. **Pick an automatic mainland rule now** — distance, distance-and-contribution,
   or silhouette-mass share.
2. **Explicit knobs, full-claim default, second mode deferred (chosen).**
3. **Precompute ladder rows for both framings** — double the artifact, both modes
   fully oracle-judged.
4. **Do nothing** — accept full-claim framing permanently.

## Rejected alternatives

**Option 1** was rejected because all three candidate rules were measured against
real cases and each fails on one: distance dismembers Indonesia, and both
contribution-based rules keep Alaska, which is the case that prompted the
question. Choosing one now would freeze a guess into a frozen constraint.

**Option 3** was rejected as premature, not wrong. Precomputing a second framing
would double the 10.4 MB artifact and the ~175 s determinism gate for a mode
whose definition is not yet decided. It stays available once option 2's deferred
definition lands, and is the better answer if the second mode turns out to be
used for the whole catalog rather than per request.

**Option 4** was rejected by the owner: the full-claim frame leaves top-200
countries visibly off-centre in a mostly empty box, which is a product defect for
a decorative card even if it is geographically honest.

## Consequences and tradeoffs

- The current catalog is unaffected; T4 and T5 proceed on today's framing.
- `TestRepresentativeNaturalRatiosAndArbitraryFrames` asserts `ratios["CL"] <
  0.75`, which encodes mainland framing that was never decided and is not current
  behaviour. Under this decision it may be realigned to the full-claim contract,
  with the deferred second mode referenced so the question is not lost.
- The crop seam carries a real quality cost that grows with zoom, and it must
  report that rather than hide it — a cropped result is displaying geometry
  simplified for a wider frame.
- P3 grows a surface it did not have. If that makes P3 too large, a successor
  spec is the honest split.

## Affected artifacts and owners

- P2 (`geometry-pipeline`) — unchanged by this decision; its ladder and framing
  behaviour stay as committed.
- P3 (`generator-cli`) — gains the crop seam and its addressable knobs, or hands
  them to a successor spec.
- P4 (`site-delivery`) — consumes whichever framing the owner selects per card.
- `WKI-490046152C71` — the tracking row for the deferred second-mode definition.
- Accountable owner: product owner (accepting principal).

## Validation and revisit trigger

Validation for part 1 is that the committed catalog totals and viewBoxes do not
move: 993 rendered, 3 typed no_artifact, and the guarded band caps hold. Parts 3
to 5 are validated when the seam is built, and must include a case proving that
cropping before selection is rejected rather than silently falling back to source.

Revisit at the T5 contact sheet, which is the explicit trigger for deciding the
second mode's definition. Revisit sooner only if the crop seam's measured quality
cost at realistic zoom turns out to be unacceptable, which would make option 3
the answer instead.

## Supersession

None. DEC-011 does not alter DEC-006's oracle boundary, DEC-007's
protected-visibility policy, DEC-009's ladder rules or DEC-010's Croatia outcome;
it adds a framing seam above them and defers one sub-decision to T5.
