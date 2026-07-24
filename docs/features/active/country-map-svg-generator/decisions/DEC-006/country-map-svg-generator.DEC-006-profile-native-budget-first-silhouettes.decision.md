---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-006"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator", "P2", "DEC-004", "DEC-005", "RUN-009-RESULT", "RESULT-016"]
---
# DEC-006 — Scale-calibrated, budget-first silhouettes

## Context

The product needs decorative maps, not analytical GIS layers. The accepted site
composition renders many `card` silhouettes at roughly 90–160 CSS px and a
smaller number of `hero` silhouettes around 500–700 CSS px. The architecture
therefore requires hard complete uncompressed SVG maxima of 2,500 and 8,000
bytes, while P2 explicitly allows unprotected islands and coastline detail to
disappear below the visual scale.

DEC-004 introduced two maintainer-built Mapshaper tiers and DEC-005 added a
q-aware source fallback. Successor plans then treated a fitted-space symmetric
boundary-deviation ratio of `0.012` as the dominant acceptance predicate and
restored required source components before serialization. That rule protects
near-cartographic fidelity, but it is not the user-visible outcome of `card` or
`hero`.

RUN-009 proves the narrower contradiction. After runtime and builder used the
same projection contract, the rebuilt AQ compact candidate still measured
`1.537690736882369px` against a `1.536px` ceiling and the mandated brake stopped
before broad validation.

Coordinator diagnostics on the uncommitted RUN-009 scratch additionally
suggested that resolution, weighting and phase tuning do not solve the
full-catalog byte problem. Those observations are not lifecycle evidence and
are not accepted by this decision as proof. The successor's first no-commit
brake must reproduce them with content-addressed inputs, commands, per-case rows
and aggregate output before any product implementation is authorized.

## Decision

Keep `card` and `hero` as versioned data presets over the general engine. They
only resolve ordinary inputs such as intended long side, padding, quality and
file budget. Neither preset name nor entity identifier may select geometry.

For each corpus geometry and boundary profile, the maintainer build produces a
generic, scale-calibrated silhouette ladder from immutable P1. Candidate
selection is a pure function of the actually fitted geometry long side, resolved
quality and resolved byte policy. `tight` and every arbitrary `contain` frame
therefore use the same engine and respond to effective fitted scale rather than
the frame's shorter side or a preset enum.

The builder may use pinned Mapshaper presimplification and a deterministic
bounded candidate search. It selects the finest candidate that passes every
generic predicate; it never silently overflows the byte policy or invokes an
external tool during ordinary generation.

The hard acceptance boundary for derived silhouettes is:

1. non-empty valid polygon topology, winding and deterministic `q=0.01`
   canonical coordinates;
2. P1 remains authoritative for entity identity, projection center, natural
   aspect input, protected anchors, capital-marker transform and exact removal
   provenance;
3. identity-defining protected land and declared minimum parts remain visible;
   unprotected components and sub-scale detail may be omitted by one versioned
   visibility policy;
4. a versioned target-scale raster silhouette oracle, dominant-component
   coverage and protected-anchor coverage replace raw full-coastline symmetric
   deviation for scale-calibrated derived candidates;
5. P2 path data is capped initially at 2,200 bytes for card-scale requests and
   7,500 bytes for hero-scale requests, reserving deterministic markup headroom;
   P3/P4 separately prove complete uncompressed SVG median, p95 and maxima of
   0.6/1.2/2.5 KB and 1.8/4/8 KB respectively;
6. representative and full-catalog contact sheets are reviewed for
   recognizability, balance and the soft-organic aesthetic before publication;
7. every candidate choice, omission, similarity metric, point count, byte
   count, source identity, recipe identity and tool identity is inspectable.

The raster oracle contract is itself a prerequisite artifact, not an
implementation detail. Version 1 must bind the exact pure-Go rasterizer source
hash, P1/projection identity, reference construction, scale bands, grid
dimensions, fill rule, metrics, thresholds and tie-break:

- reference and candidate use the same exact P1-owned projector, natural
  viewBox, padding and uniform fit;
- the grid long side is the maximum effective scale of the selected generic
  band; the other dimension follows the same natural viewBox;
- occupancy samples the pixel center with SVG's nonzero fill rule and no
  antialiasing, color, stroke or softening;
- the machine metrics are filled-pixel intersection-over-union plus
  reference-to-candidate filled-pixel recall, with protected anchors checked
  independently;
- T0A calibrates one global threshold pair per scale band against the owner
  sheet, records it in the recipe and freezes it before artifact search; entity
  and preset-specific thresholds are forbidden;
- candidates are tried in one recorded fine-to-coarse presimplification order;
  the first passing candidate wins, with source polygon/ring order and content
  digest as the final stable tie-break.

Candidate search may move only from finer to coarser presimplification
thresholds and from visible to sub-scale unprotected components. It cannot
branch on a country or preset identifier, invent or hand-move coordinates,
relax topology, remove a protected feature, change the serializer, or accept a
budget overflow. A versioned data override may preserve a declared identity
feature or adjust presentation after visual review; executable country-specific
branches remain forbidden.

The pure-Go runtime embeds the accepted generic scale ladder, resolves it from
effective fitted scale and byte policy, projects markers, applies uniform
`tight`/`contain` layout transforms and emits presentation-free path commands.
Styling, CSS variables, dark/light themes and small animations remain P3/P4
concerns and do not duplicate geometry.

The 300/500-byte initial envelope reserves are cross-spec brakes, not a claim
that P2 owns SVG markup. P3 must publish deterministic envelope measurements for
each delivery mode and optional pin. If any accepted envelope needs more room,
the corresponding P2 path cap tightens; the architecture's complete-file
maximum never expands silently.

Explicit `source` or custom-quality generation remains available for callers
who need closer coastline fidelity. DEC-005's q-aware deviation guard continues
to govern that explicit path, including typed failure when a caller combines an
incompatible fidelity request and byte maximum. Arbitrary dimensions remain
ordinary layout inputs; their effective fitted scale selects generic quality.
The deferred `detail` preset remains the future high-fidelity convenience
surface over the same engine.

## Options considered

| Option | Result |
| --- | --- |
| Continue tuning `resolution`, weighting or grid phase | Rejected as the successor direction; RUN-009 still fails after aligned projection and the first successor brake must reproduce broader diagnostics |
| Reinsert source vertices until every tier passes the `0.012` deviation ceiling | Rejected for named profiles; it spends bytes to recover detail the user explicitly permits removing |
| Increase the 2,500/8,000 maxima | Rejected; repeated-card network and DOM cost is a primary product constraint |
| Accept source fallback whenever topology passes | Rejected as a catalog policy; the diagnostic matrix contains hundreds of source and finer-tier budget overruns |
| Precompute a generic scale-calibrated silhouette ladder with byte and target-scale visual oracles | Selected; presets remain data and arbitrary layouts keep fitted-scale quality |
| Hand-draw all countries | Rejected as the catalog pipeline; it is non-repeatable, though versioned data overrides remain available for reviewed identity fixes |

## Rejected alternatives

Sub-pixel full-coastline deviation is a useful analytical safety measure but the
wrong primary oracle for a 90–160 px decorative silhouette. Restoring every
retained source component reverses the accepted visibility policy and makes
remote islands dominate bytes despite having no visible role in a card.

Adding another dynamic LOD or a runtime repair loop retains the same mistaken
objective and makes the local CLI slower. Post-hoc SVGO, lower output precision
and silent clipping remain rejected because they obscure provenance or can
damage topology. The selected approach changes the derived geometry contract
upstream, where it can be rebuilt, measured and visually approved.

## Consequences and tradeoffs

Ordinary preset generation becomes predictable, fast and offline: the expensive
cartographic work happens once in the maintainer build. The generated SVG
remains style-neutral, naturally proportioned, themeable and
animation-friendly. The accepted generic path can be used repeatedly without
shipping a geometry engine or source corpus to the browser.

Named-profile coastlines are intentionally less literal. Some unprotected
islands and fine inlets disappear, and raster similarity plus visual approval
become first-class product evidence. This is an explicit trade rather than an
unreported fallback. Exact P1 identity and omission provenance keep the result
auditable.

The artifact schema now carries generic scale-band and oracle identities in
addition to boundary profile. A future `detail` preset can resolve a finer band
without changing `card`/`hero` preset semantics. Custom/source generation
remains more expensive and may fail a caller-supplied hard budget honestly.

## Affected artifacts and owners

P1 and CTR-001 remain unchanged and authoritative. P2 owns generic
scale/quality recipes, presimplification selection, visibility policy,
protected-feature enforcement, raster silhouette oracles, omission provenance
and embedded geometry. CTR-002 must expose effective scale, selected generic
band, recipe/oracle versions, removals, similarity evidence and path bytes
without leaking Mapshaper into runtime.

P3 continues to own SVG structure, styling modes, CSS-variable hooks and CLI
configuration. P4 owns full-catalog browser/network proof and producer/client
performance guidance. P5 ships only Go binaries and generated artifacts; Node,
npm and Mapshaper remain maintainer-only.

## Validation and revisit trigger

Before broad implementation, a no-commit diagnostic first reproduces the
coordinator observations against content-addressed source, scratch, recipe,
lockfile and tool hashes. It records exact commands, per-resolution,
per-weighting and phase rows, and all 996 case rows plus aggregate path-byte and
timing digests. A contradiction stops and amends or withdraws this decision
before product bytes are committed.

The same hard brake then freezes `silhouette-oracle/v1` and builds generic
card-scale and hero-scale candidates for RU, CA, CN, AQ, AE, ID, CL, AR, AU and
protected island/microstate fixtures. It records candidate thresholds, omitted
parts, raster metrics, topology, protected coverage, points, path bytes,
complete-SVG envelope estimates, rebuild hashes and runtime. Every
representative path must meet the 2,200/7,500 P2 cap, every estimated complete
file must remain below 2,500/8,000, and the same outputs render into an
owner-reviewable contact sheet.

The full gate covers all 249 entities, both boundary profiles and both preset
inputs through the generic engine. It proves 996 deterministic outputs, P2
path-cap distributions, P3/P4 complete uncompressed SVG median/p95/max targets,
no missing entity, exact removal provenance, fitted-scale selection, stable
natural/contain layout and ordinary offline execution without Node or network.
Mutation tests must fail a lost protected anchor, invalid ring, stale
P1/recipe/oracle identity, changed rasterizer behavior, threshold drift,
nondeterministic candidate order, hidden omission, silhouette-oracle bypass or
budget overflow.

Revisit if a generic profile recipe cannot cover the full catalog, the raster
oracle passes visibly unrecognizable silhouettes, protected microstates cannot
fit the byte limits, the owner rejects the representative sheet, or a later
product genuinely needs analytical coastline fidelity in named profiles.

## Supersession

Refines DEC-004's derived-candidate acceptance from strict full-coastline
deviation to a generic fitted-scale silhouette oracle plus a path sub-budget.
Retains DEC-004's scale-only selection, immutable P1 authority,
maintainer-only Mapshaper, deterministic versioned artifacts, pure-Go offline
runtime and no country-code or preset branches. Refines DEC-005 only by keeping
its q-aware source fallback outside automatic budgeted derived-candidate
acceptance; DEC-005 remains binding for explicit `source` and custom quality.
DEC-003 remains binding for projection, Go rendering and the prohibition on
post-hoc SVG optimization.
