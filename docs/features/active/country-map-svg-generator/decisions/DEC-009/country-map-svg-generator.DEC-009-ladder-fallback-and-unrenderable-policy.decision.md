---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-009"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-009 — Identity fallback rung and the unrenderable-band policy

## Context

AM-005 left two questions to the product owner, both raised by the T0 catalog
feasibility measurement (`readiness/T0-catalog-feasibility.matrix.tsv`, sha256
`711aecae4239e4e7f013bf14a598507b25f6b5b85283f99afb2a0b14aabedcb0`). Under the
DEC-006 oracle predicate with per-entity fine-to-coarse ladder search, 986 of 996
catalog rows pass with the frozen thresholds untouched. Ten fail, across three
entities.

**Where the identity rung belongs.** SH and UM carry components separated by
thousands of kilometres of ocean. Mapshaper sizes its simplification grid from
the feature's own bounding box, so for an ocean-spanning bbox every island falls
below grid resolution at every one of the 29 declared rungs, including the finest.
Both collapse to a single surviving polygon everywhere on the ladder. Evaluated
without simplification they pass cleanly — SH at IoU 1.0 and 637 of 2200 path
bytes at compact, 630 of 7500 at standard; UM at 791 of 7500 at standard. The
ladder needs a candidate that declines to simplify. DEC-006 says selection takes
"the finest candidate that passes every generic predicate", which read literally
places a no-simplification candidate first.

**What to do about a band that cannot be drawn.** UM at compact fails even
without simplification: the geometry does not survive `q=0.01` canonical
quantization at card fitted scale, and the ring collapses. No ladder rung
addresses this. At 90–160 CSS px a dozen atolls scattered across the Pacific and
Caribbean are not a silhouette; they are a few pixels of noise. The owner's
position is that such geometry cannot be drawn meaningfully at that scale and the
product should say so rather than fake it.

The product requirement that frames both: the top ~200 countries must generate
without exception, at high visual quality, maximally optimized, and customizable
through CSS from an agent-driven CLI. Of the three failing entities only Croatia
is in that set, and its gap is nine bytes.

## Decision

**1. The identity candidate is a fallback, tried last.** A no-simplification
candidate is evaluated only after every declared resolution rung has failed to
produce a valid candidate. It is judged by exactly the same frozen gates as any
other candidate — the silhouette oracle, DEC-007 protected visibility, and the
byte caps — and receives no exemption.

It is not a fidelity tier and must not be described as one. It exists because
resolution is bbox-relative and therefore meaningless for a geometry whose bbox
spans an ocean; it is a fallback for a tool limitation, not a quality preference.

Consequence, and the reason for the placement: only rows that currently fail can
reach it. No row that passes today can change its selection, so the 986 passing
selections, the representative artifacts and both owner contact sheets are
unaffected.

**2. A band that cannot satisfy topology has no artifact, by generic rule.** When
no candidate in a band — including the identity fallback — satisfies topology at
the actual fitted scale, that entity/band pair produces no artifact. The manifest
records a typed, machine-readable reason naming the failing gate, and P3 returns
an explainable typed failure for that entity/preset pair.

This is a property of the generic engine, not an exception. No country literal or
branch is introduced; the outcome falls out of the same gates every other entity
runs. It applies automatically to any future corpus entity with the same
character, and AC-3's requirement that exceptional countries be expressible
through data rather than code forks is preserved because there is no per-country
code at all.

**3. Intermediate resolution rungs are added between 80 and 64.** Croatia's
compact candidate serializes to 2209 bytes at resolution 80 against a 2200 cap,
and at the next declared rung, 64, loses a source component whose contribution is
exactly the retention threshold of 110. The declared ladder has nothing between
them. The additions close a nine-byte gap; they change no threshold, tolerance,
weighting, byte cap or budget.

**4. Raw boundary deviation leaves the derived path entirely.** It remains the
gate for an explicit `source` or custom-quality request under DEC-005, which is a
caller-selected mode where coastline closeness is the point. Ordinary catalog
generation does not consult it.

Under this decision the catalog resolves to: SH gains both bands through the
fallback; UM gains its standard band and has no compact artifact; HR gains its
compact band through the intermediate rungs. Every top-200 country generates.

## Options considered

| Option | Result |
| --- | --- |
| Identity rung tried last, as a fallback for bbox-relative resolution | **Selected**; provably cannot alter any passing row, and matches what the rung actually is |
| Identity rung tried first, as the finest tier per DEC-006's literal "finest wins" | Rejected; see below |
| Identity rung first, after measuring how many of the 986 rows change selection | Rejected as the plan of record; the measurement is real work and the expected benefit is near zero, but it remains available if a later need for maximum fidelity emerges |
| UM compact: typed refusal under a generic unrenderable rule | **Selected** |
| UM compact: group anchor through the DEC-007 override seam | Rejected for now; see below |
| UM compact: revisit the `q=0.01` quantization contract | Rejected; general change, disproportionate blast radius |
| Raise the byte budgets, or relax an oracle threshold | Rejected; DEC-006 rejected both by name and T0 gives no evidence for either — UM fails on topology, not on bytes or similarity |

## Rejected alternatives

**Identity rung first.** DEC-006's "finest candidate that passes" read literally
puts a no-simplification candidate at the head of the search, and any small entity
whose unsimplified geometry fits the cap would then select it — Kiribati passes
today at 38 bytes, Gibraltar at 63, Vatican City at 67, Pitcairn at 68. The
fidelity gain is close to nothing, because weighted Visvalingam at resolution 512
removes almost nothing from a small atoll. The cost is concrete: every one of the
986 selections must be re-measured and both owner sheets re-approved. Poor trade
for a change whose purpose is to rescue two entities the ladder cannot express.

**Group anchor for UM compact.** DEC-007 item 6 provides a versioned data override
for preserving a declared identity feature, and UM's atolls are the shape it was
written for. It is rejected for now not because it is illegitimate but because it
would preserve components that cannot be rendered meaningfully at card scale — it
solves the gate, not the user-visible problem. It remains available if the owner
later wants a card artifact for UM at any cost.

**Revisiting `q=0.01`.** Every one of the 986 passing rows depends on that
quantization contract, along with the byte-equivalence guarantee under REQ-8.
Changing it to rescue one entity/band pair inverts the cost of the fix and the
size of the problem.

## Consequences and tradeoffs

Positive: the top-200 requirement is met in full; no passing row changes; no
contact sheet is re-approved on account of this decision; the engine gains no
country-specific code; the unrenderable rule generalizes to future corpus
changes.

Negative: the catalog is one decorative card short of complete — 995 of 996
ordinary outputs. The "complete ISO catalog" claim in the epic must be stated
precisely: every assigned entity is represented, one entity has no card-scale
artifact, and the reason is recorded and machine-readable rather than silent.

Operational: a consumer requesting UM at card scale receives a typed failure, not
an empty or malformed SVG. P3 and P4 must handle that path deliberately, and P4's
full-catalog site proof must show what a missing decorative card looks like in
place.

Accepted deviation from DEC-006's literal "finest candidate wins": the identity
candidate is ordered last rather than first. DEC-006's intent — maximum quality
subject to the byte policy — is preserved for every declared resolution; the
exception applies only where the ladder cannot express a candidate at all.

## Affected artifacts and owners

- P2 specification: the validation allocation carries the unrenderable outcome as
  a typed result rather than a silent gap. Owner: product owner, through AM-005.
- P2 ladder recipe and build tool: the fallback candidate and the intermediate
  rungs. Owner: geometry maintainer.
- CTR-002 provenance: selected band and resolution, the fallback marker when the
  identity candidate is used, and the typed reason when a band produces no
  artifact. Owner: geometry maintainer.
- P3 CLI: explainable typed failure for an entity/preset pair with no artifact,
  and its JSON shape. Owner: P3.
- P4 site proof: the missing-card path rendered in place. Owner: P4.
- DEC-006 is refined, not superseded: its oracle-first acceptance boundary and its
  rejection of budget increases stand unchanged.

## Validation and revisit trigger

Proof obligations:

- the identity candidate is reachable only after every declared rung fails, and a
  mutation that promotes it ahead of the declared rungs must fail a test;
- adding the identity fallback and the intermediate rungs changes none of the 986
  passing selections recorded in the T0 matrix — this is a direct, mechanical
  comparison and is the gate that proves the placement claim;
- a band with no valid candidate emits a typed reason naming the failing gate, and
  a mutation that emits a silent gap, an empty artifact or a malformed SVG must
  fail;
- no country literal or per-entity branch appears in executable selection code.

Revisit if: the corpus refresh introduces further entities that reach the identity
fallback, which would suggest the declared ladder is mis-parameterized rather than
merely incomplete; the owner requires a card artifact for an entity the policy
declares unrenderable; or a future need for maximum fidelity justifies measuring
identity-first placement across the catalog.

## Supersession

Refines DEC-006 on candidate ordering and on the outcome when no candidate exists.
Resolves the open owner decision recorded in AM-005 §Accepted change item 6. Does
not supersede DEC-004, DEC-005, DEC-007 or DEC-008.
