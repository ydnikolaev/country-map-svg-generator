---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "specification"
id: "P4"
epic: "country-map-svg-generator"
spec: "P4"
status: draft
profiles: []
concerns: []
inputs: ["DISC-006", "ARCH-001", "DEC-001"]
---
# P4 — Full Catalog and Site Integration

## Outcome

The owner receives a complete optimized catalog and can place it into the visa
site with documented patterns that preserve theme behavior, visual quality and
page performance. Full-catalog evidence proves coverage, deterministic output,
SVG safety and byte budgets; browser fixtures prove the hero and repeated-card
use cases in light, dark and reduced-motion modes.

## User stories and stakeholders

| ID | Story |
| --- | --- |
| US-1 | As the site producer, I can generate all countries into one ready-to-copy folder with an inventory and no later optimization step. |
| US-2 | As a frontend developer, I know which delivery pattern to use for a themed hero and many lightweight cards. |
| US-3 | As a visitor, repeated country cards load and render without excessive network, DOM or animation cost. |
| US-4 | As the owner, I can review a balanced visual sheet and approve exceptions without inspecting 200+ raw path files. |
| US-5 | As a maintainer, a regression in coverage, safety, size or browser behavior blocks the catalog. |

## Scope and non-goals

In scope: full ISO catalog generation for hero/card, manifest/statistics, structural
and visual evidence, owner site preset, browser fixtures, producer documentation,
client integration/performance documentation and one copy-ready output layout.

Out of scope: implementing the visa site itself, runtime map generation, a
general-purpose frontend component package, CDN/vendor selection, all-country
preloading, public gallery/marketing and `detail`.

## Context and ground truth

The supplied desktop composition shows one large hero slot and a grid of repeated
cards. Hero art can justify inline SVG for theming and subtle motion. Cards are
decorative, small and numerous; they should avoid persistent animation and
unnecessary inline DOM. Host CSS cannot style internals of an SVG loaded through
`<img>`, while inline SVG and mask/color patterns support site themes differently.
ARCH-001 sets separate size/byte budgets.

## Requirements and invariants

| ID | Requirement |
| --- | --- |
| REQ-1 | The deliverable contains every currently assigned ISO entity for `card` and `hero`, an owner preset/config, manifest and human-readable summary. |
| REQ-2 | Full-catalog validation proves count/identity, XML/SVG structure, safety, viewBox containment, deterministic digests and ARCH-001 byte distributions. |
| REQ-3 | Representative visual evidence covers large, small, fragmented, island, microstate, polar/antimeridian and disputed-boundary cases in both profiles. |
| REQ-4 | Hero guidance defaults to `themed-inline`, stable CSS variables and optional transform/opacity motion with reduced-motion suppression. |
| REQ-5 | Card guidance defaults to external SVG mask/currentColor or standalone `<img>` according to theming need; it forbids persistent per-card animation and bulk-inlining the entire catalog. |
| REQ-6 | Client guidance covers explicit rendered dimensions/aspect handling, caching/fingerprinting, loading only rendered/near-rendered cards, accessibility semantics and fallbacks. |
| REQ-7 | Producer guidance covers profile selection, preset inheritance, per-country exception discipline, batch selection, manifest/budget interpretation and when regeneration is required. |
| REQ-8 | Browser fixtures prove light/dark theme changes, currentColor/custom-property fallback, mask use, no layout shift from missing dimensions and reduced motion. |
| REQ-9 | Every budget exception is explicit in versioned config, names the ISO/reason/metric and remains below an accepted exception ceiling. |
| INV-1 | The delivered site catalog is a complete validated set; individual files are never promoted without their matching manifest. |

## Interfaces data and behavior

Input is the P3 built binary, accepted corpus and owner preset. Output layout:

```text
dist/
  card/{iso2}.svg
  hero/{iso2}.svg
  catalog.json
  README.md
```

The client contract is generated files plus CTR-006; generator internals are not
shipped to the site. The integration guide presents copyable inline-hero,
mask-card and standalone-image examples and states their theme/DOM/network
tradeoffs.

## Dependencies and handoffs

Depends on P3. Owns BND-007 and BND-008. It consumes CTR-004–CTR-006 and closes
the owner-value critical path. P5 may execute concurrently but cannot gate P4.
Owner approval is required for the visual sheet and any budget exception.

## Profiles and concerns

- Browser performance: repeated-card network, paint, DOM and animation cost.
- Accessibility: decorative versus meaningful image semantics and reduced motion.
- Visual regression: the aesthetic contract cannot be proven by XML alone.
- Documentation: separate producer and client paths prevent expensive misuse.
- Elevated risk: performance/accessibility/public CSS behavior follow DEC-001.

## Validation impact

| ID | Protects | Guard/action | Tier and teeth | Owner/due |
| --- | --- | --- | --- | --- |
| VAL-1 | REQ-1, REQ-2 | generate the entire catalog twice from clean dirs and compare manifest/digests/counts | E2E; missing/extra/non-deterministic asset fails | P4 / W4 complete |
| VAL-2 | REQ-2, REQ-9 | calculate uncompressed median/p95/max per profile and inject an oversized asset/undeclared exception | full-catalog gate; thresholds and exception schema have teeth | P4 / W4 complete |
| VAL-3 | REQ-3 | render a stratified contact sheet at actual card/hero CSS sizes | visual review; owner accepts sheet and named exceptions | P4 / W4 complete |
| VAL-4 | REQ-4, REQ-5, REQ-6, REQ-8 | exercise inline hero, mask card and image card in light/dark/reduced-motion fixtures | browser E2E; theme, fallback, dimensions and motion assertions | P4 / W4 complete |
| VAL-5 | REQ-5, REQ-6 | render a representative card grid and assert no all-catalog inline payload or persistent animation | browser/performance guard; fixture mutation fails | P4 / W4 complete |
| VAL-6 | REQ-7 | execute documented producer commands against shipped preset | docs smoke; examples and expected manifest facts match | P4 / W4 complete |

## Acceptance criteria

| ID | Covers | Criterion | Evidence |
| --- | --- | --- | --- |
| AC-1 | US-1, REQ-1, REQ-2 | Copy-ready folder contains the complete card/hero catalog and every asset passes structural/budget gates without post-processing. | full-catalog receipt |
| AC-2 | US-2, REQ-4, REQ-5, REQ-6, REQ-7, REQ-8 | Documented inline hero, mask card and image card examples work in browser fixtures across themes and reduced motion. | browser E2E |
| AC-3 | US-3, REQ-5, REQ-6 | Card grid follows bounded loading/DOM/animation rules and ships profile-budgeted assets. | performance fixture and catalog stats |
| AC-4 | US-4, REQ-3, REQ-9 | Owner accepts the contact sheet and every remaining exception is named, bounded and documented. | owner acceptance record |
| AC-5 | US-5, all | Mutation tests demonstrate coverage, size, unsafe markup and browser regressions turn the gate red. | gate teeth receipt |

## Rollout rollback and operations

Start with a representative site fixture, then generate/review the full catalog,
then copy it into the owner site. The manifest supports cache fingerprinting and
rollback to a prior complete folder. Catalog replacement is atomic at the site's
asset-version boundary. Documentation records measured distributions rather than
unverified performance claims.

## Delivery constraints

Do not add a site framework dependency to the generator. Do not inline 200+ maps
into an initial page, preload the full catalog, animate every card, or rely on
client-side SVG optimization. Preserve user-site autonomy: examples use plain
HTML/CSS and explain adaptation rather than mandating a component library.

## Decisions and unresolved questions

DEC-001 governs acceptance of performance/accessibility exceptions. The owner must
approve the representative visual sheet; this is a planned acceptance boundary,
not an unresolved product question. `detail` remains backlog.

## Amendments

None.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
