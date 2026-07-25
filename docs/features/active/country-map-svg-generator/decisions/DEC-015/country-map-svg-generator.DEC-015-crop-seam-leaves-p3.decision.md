---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-015"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-015 — The crop and component-selection seam leaves P3

## Context

DEC-011 decided that cropping is an explicit seam over the *already selected*
candidate — select from the ladder, then crop, then refit — never an automatic
heuristic, and that the knobs are addressable per request through the identity
vocabulary rather than per country. Its item 6 assigned that work to "P3, or to a
successor spec if P3 is already too large," and its item 5 left an open
sub-decision: the oracle measures IoU and recall against the full visible
reference, and a crop changes what that reference should be, so whether a cropped
variant is re-judged against the cropped reference, inherits its parent's verdict,
or is marked as unjudged provenance must be settled when the seam is built.

DEC-013 then narrowed what remains. The card is now fitted to the silhouette it
draws rather than to the whole territorial claim, so the framing complaint that
motivated `WKI-490046152C71` is answered: France is the hexagon with Corsica,
Chile is tall and narrow, cards under a fifth full went from 20 to 0. What is left
of that work item is component *selection* — which components a request asks for —
not framing.

P3 as specified carries thirteen requirements, eight validation obligations and
the entire CLI: config schema and inheritance, five styles, two delivery modes,
markers, animation, a minimal SVG serializer, transactional publication and a
catalog manifest. P3's spec body states "No unresolved question blocks P3," which
predates DEC-011 and is no longer true if the seam stays.

## Decision

**1. The crop and component-selection seam is not P3 scope.** It moves to a
successor spec, exercising the escape DEC-011 item 6 named. P3 ships the CLI over
the current default: the fitted-to-candidate framing DEC-013 established, with no
crop and no component selection.

**2. DEC-011 item 5 travels with it.** The acceptance rule for a cropped variant
is settled in the successor spec, by the session that builds the seam and can
measure the three candidate rules against real output. It is not settled here and
it is not settled by P3.

**3. P3's configuration schema v1 reserves the `selection` key namespace.** No
other v1 key may be named `selection` or nested beneath it, and an unknown key
there fails validation like any other unknown field. This is the whole cost of the
deferral: the seam arrives later as an additive change within the schema major
rather than as a breaking change requiring a new major and a migration for every
config already written.

**4. The successor spec is created when the impact fence lifts.** Adding a spec to
a fenced epic is itself governed and currently refused. Until then this decision
plus `WKI-490046152C71` are the durable record of the obligation.

## Options considered

1. Build the seam inside P3.
2. **Defer it to a successor spec and reserve the schema namespace (chosen).**
3. Defer it with no reservation, accepting a schema major later.
4. Drop the seam and keep full-claim framing permanently — DEC-011's option 4,
   already rejected there.

## Rejected alternatives

**1 — build it in P3.** It would require settling DEC-011 item 5 before the config
schema can freeze, because component selection has to be expressible in CTR-3.
That is a fidelity-judgement decision the owner would be asked to make from
hypotheticals, before the CLI that renders the comparison exists — the same
mistake DEC-011 item 2 explicitly avoided by deferring the second framing mode
until the T5 contact sheet made it visible. It also makes the largest spec in the
epic larger while P3 is already delivering ungoverned under DEC-014.

**3 — defer without reserving.** Cheap now, expensive exactly once: the seam lands
as `country-map/v2`, and every config a producer or agent has written needs a
migration for a feature they did not ask for. Reserving one key costs a line in
the schema validator.

**4 — drop it.** Already rejected in DEC-011 on the owner's evidence that a
decorative card sometimes needs a different crop.

## Consequences and tradeoffs

- P3 gets smaller and can freeze its config schema at T2 without an open product
  question hanging over it.
- `WKI-490046152C71` stays deferred, with its remaining content restated as
  component selection rather than framing.
- A producer who wants a cropped card cannot have one until the successor spec
  ships. DEC-013 makes that much less painful than it was — the default framing is
  now fitted to what is actually drawn.
- The successor spec inherits a real, bounded surface rather than a vague one: the
  crop order is fixed by DEC-011 item 3, the vocabulary by item 4, and the one
  open question is named in item 5.

## Affected artifacts and owners

- P3 specification — scope reduced; its "no unresolved question blocks P3" line is
  true again as a result of this decision rather than in spite of it.
- P3 configuration schema v1 (CTR-3), `internal/config/**` (BND-005) — must
  reserve `selection`; teamlead.
- `WKI-490046152C71` — product-owner; remains deferred, now scoped to component
  selection.
- DEC-011 — unchanged; this decision exercises its item 6 rather than amending it.
- The successor spec — to be created when the fence lifts.

## Validation and revisit trigger

The reservation is validated by the config schema gate: a v1 document using
`selection` must fail as an unknown field, and that assertion must exist before P3
closes T2. Revisit when the successor spec is created, or earlier if the owner
needs a cropped variant before then — in which case this decision is superseded
rather than quietly widened.

## Supersession

Supersedes nothing. Exercises DEC-011 item 6 and carries its item 5 forward
unresolved. Depends on DEC-013 having made fitted-to-candidate the default, which
is what makes the deferral tolerable.
