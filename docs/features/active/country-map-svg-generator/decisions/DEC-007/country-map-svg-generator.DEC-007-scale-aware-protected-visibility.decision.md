---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-007"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-007 — Scale-aware protected visibility

## Context

P1 deliberately records protected-feature anchors and source-corpus
`minimum_parts` so an island state, archipelago or microstate cannot silently
disappear. P2 also deliberately removes coast detail and small islands that have
no visual role at a requested scale.

DEC-006 made both rules collide by requiring the raw P1 minimum-part count in
every automatic scale band. RUN-011 proved the collision with content-addressed
evidence. For Indonesia at compact scale, Mapshaper resolution 121 is the
coarsest tested generic candidate retaining 20 parts and emits 2,713 path bytes
with a 2,993-byte complete-file estimate. Resolution 120 retains 19 parts.
Therefore no candidate can satisfy the current 20-part brake and the frozen
2,200/2,500 compact budgets. The stricter provisional IoU/recall pair would
require an even finer candidate.

The owner already fixed the product intent: small card maps are decorative,
recognizable silhouettes; minor islands and coast detail should not be
overworked when they are not visible at that scale. Raising budgets, hiding an
Indonesia branch or weakening all protected checks would contradict that
intent.

## Decision

Keep P1 protected metadata and source bytes unchanged. Interpret protection at
two distinct boundaries:

1. P1 `minimum_parts` remains a source-corpus integrity requirement. Explicit
   `source` and custom-quality requests continue to retain the full declared
   minimum through DEC-005.
2. An automatic derived scale-band candidate uses a versioned,
   presentation-free protected-visibility policy. It must always retain:
   - the dominant identity component;
   - the source component containing every protected anchor;
   - each member of the feature's deterministic identity set that is visible at
     that band's oracle grid.
3. A feature's identity set is the declared `minimum_parts` largest source
   polygon components by reference filled-pixel contribution, with the
   anchor-containing component forced into the set. Ties use source polygon
   order and then content digest. This derives protection from P1 without
   country-specific lists.
4. Visibility uses the same natural projection, bounds, uniform fit,
   pixel-center nonzero rasterizer and maximum effective scale as
   `silhouette-oracle/v1`. One global minimum filled-pixel contribution is
   calibrated and frozen per generic scale band. Preset names, entity IDs and
   contain-frame shape cannot affect it.
5. An identity-set component below the band threshold may be omitted only with
   explicit `protected_subscale` provenance. It is not reclassified as
   unprotected. The oracle recall, dominant coverage, anchor coverage, topology,
   byte caps and owner-approved sheet remain mandatory.
6. Component area is a deterministic default proxy, not semantic truth. If the
   owner sheet proves that it misses a recognizable remote island group, the
   existing bounded data-override seam may add a named protected group anchor.
   Such an anchor forces its containing source component at every band, carries
   review provenance, and reruns the whole representative and catalog gate. It
   cannot change a band threshold, budget, candidate order or code path.
7. Candidate search stays fine-to-coarse. It cannot invent, hand-move or
   reinsert vertices; branch on a country; change q or serializer; raise
   budgets; or bypass the generic silhouette oracle.

This rule changes protection cardinality only for automatic derived
silhouettes. It does not change P1, explicit source/custom behavior, marker
projection, natural/contain layout or public style semantics.

## Options considered

- Preserve all declared parts at every scale and raise the compact budget.
- Preserve all declared parts and optimize or replace the frozen serializer.
- Add an Indonesia-specific retained-part override.
- Treat `minimum_parts` as source-only and drop every non-anchor island
  automatically.
- Derive a scale-aware protected identity set and apply one band-global
  visibility threshold.

## Rejected alternatives

- **Higher budgets:** rejected because 2,713 path bytes already exceed the
  complete 2,500-byte architecture maximum before markup, and finer
  oracle-passing candidates cost substantially more.
- **Serializer work:** rejected because the conflict exists in geometry
  cardinality and visual thresholds, the serializer is edit-frozen, and
  post-hoc optimization would hide the wrong representation boundary.
- **Indonesia branch:** rejected because it does not generalize to other
  archipelagos and would make equal resolved inputs entity-dependent.
- **Anchor-only protection:** rejected because one point is too weak to preserve
  the recognizable multi-island identity of an archipelago.
- **Custom reinsertion:** rejected because it fabricates a special geometry and
  breaks candidate provenance.

## Consequences and tradeoffs

Compact archipelagos may intentionally show fewer parts than the source corpus.
The omission is scale-specific, measurable and inspectable. At larger effective
scales more identity-set components become visible and therefore mandatory.
Source/custom output retains the original minimum count.

The runtime remains pure Go, offline and generic. Arbitrary tight/contain sizes
continue to select from effective fitted scale; `card` and `hero` remain data
presets. A single band policy is harder to calibrate than an entity override,
but it prevents hidden exceptions and gives the owner one coherent visual rule.
Rare owner-proven group anchors remain data, not executable country branches,
and cannot tune quality or budgets.

P2 diagnostics must distinguish `protected_subscale` from ordinary unprotected
omission. Tests and artifacts gain the identity-set membership, per-component
reference pixel contribution, threshold, retained/omitted status and digest
tie-break.

## Affected artifacts and owners

- P2 / CTR-002: protection and removal diagnostics for automatic silhouettes.
- DEC-006: minimum-part retention and visibility-policy clauses only.
- PLAN-009 and RUN-011: invalidated successor inputs; their diagnostic and
  blocker evidence remain authoritative.
- `silhouette-oracle/v1`, v2 recipe and ladder artifacts: record band-global
  protected visibility and identity-set provenance.
- P1 / CTR-001: unchanged owner of source metadata and source-corpus integrity.
- P2 geometry owner: implementation and machine gates.
- Product owner: compact/standard threshold and representative-sheet approval.

## Validation and revisit trigger

Before broad generation:

- reproduce RUN-010 T0A.1 and the RUN-011 Indonesia 120/121 boundary;
- mutation-test dominant loss, anchor loss, identity-set ranking/tie drift,
  visible protected loss, false `protected_subscale` provenance and
  entity/preset/frame branching; mutation-test any group-anchor override so it
  can only add retention and cannot affect thresholds or candidate order;
- calibrate one contribution threshold per generic band over RU, CA, CN, AQ,
  AE, ID, CL, AR, AU, island and microstate fixtures;
- prove every accepted candidate passes topology, global IoU/recall,
  `2200/7500` path caps and `2500/8000` complete-file maxima;
- render retained and omitted protected components visibly in the exact owner
  contact sheet;
- rerun the full 996-output catalog gate before publication.

Revisit if no global band thresholds satisfy budgets and owner recognition, an
anchor/dominant component cannot fit, a supposedly subscale protected component
is visibly identity-defining, larger scales lose a component retained at a
smaller scale, or public callers need a protection control not expressible as
quality/source selection.

## Supersession

DEC-007 supersedes only DEC-006's unconditional automatic
`minimum_parts`-retention rule and its matching protected-visibility language.
DEC-006 remains binding for budget-first generic ladders, global raster
thresholds, scale-only selection, budgets, q, serializer freeze, provenance and
owner gates. DEC-005 remains binding for explicit source/custom requests.
