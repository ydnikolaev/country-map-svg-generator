---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "solution-architecture"
id: "ARCH-001"
epic: "country-map-svg-generator"
status: draft
profiles: []
concerns: []
inputs: ["DISC-006", "DEC-001", "DEC-002", "DEC-003"]
---
# ARCH-001 — Solution architecture

## Outcomes and constraints

Produce a complete, deterministic catalog of ready-to-use country SVGs from one
local command. The first delivery serves the owner's visa site: a large hero map
and many small card decorations. The generator must remain useful to colleagues
through portable configuration, stable machine output, and no runtime dependency
other than one binary.

The runtime is offline. Source acquisition and corpus refresh are maintainer
operations, not part of normal generation. The MVP optimizes the `hero` and `card`
profiles; a higher-detail profile is explicitly deferred.

## Ground truth and probes

- The repository began as a mate-governed greenfield scaffold with no product code.
- ISO 3166-1 is the catalog authority; the output set is every currently assigned
  ISO 3166-1 entity, not an informal "countries" list.
- Natural Earth is public domain and supplies admin geometry and populated-place
  data, but its default boundary representation is de facto. Boundary posture is
  therefore an explicit profile governed by DEC-002.
- The accepted site composition contains a roughly 500–700 CSS px hero illustration
  and repeated 90–160 CSS px card illustrations, with occasional larger accents.
- The current mate consumer is `lang: [go], kind: cli`. Its Go profile requires
  black-box testing of the built binary, a `GOWORK=off` build, hermetic testscript
  breadth, and one sequential product-loop test.
- D3's azimuthal equal-area projection and SVG path output are the comparison
  oracle, not a runtime dependency. The chosen implementation is governed by
  DEC-003.

## Architecture invariants

| ID | Invariant | Owner | Evidence |
| --- | --- | --- | --- |
| ARCH-INV-1 | Every emitted country is keyed by stable ISO alpha-2; catalog completeness and uniqueness fail closed. | P1 | corpus manifest gate |
| ARCH-INV-2 | Normal generation performs no network access and needs no Node, Python, GDAL, or external data files. | P3 | built-binary E2E in an isolated environment |
| ARCH-INV-3 | Geometry, presentation, and delivery profiles remain independent contracts. | P2/P3 | package-boundary tests and config schema |
| ARCH-INV-4 | Equal inputs, corpus version, and generator version produce byte-identical output. | P2/P3 | golden and repeat-generation digest tests |
| ARCH-INV-5 | `card` output protects identity-defining land while removing detail below its visual scale; simplification never silently drops an ISO entity. | P2 | full-catalog structural and visual review |
| ARCH-INV-6 | CSS theme hooks and animation hooks are stable only in `themed-inline`; standalone assets remain visually complete without host CSS. | P3/P4 | browser fixtures in light/dark/reduced-motion modes |
| ARCH-INV-7 | Pins are opt-in and derive from a versioned registry with explicit multi-capital roles and overrides. | P1/P3 | registry and rendering tests |
| ARCH-INV-8 | Public distribution work cannot delay the first local generator and site integration. | P4/P5 | wave DAG and acceptance ordering |

## Components and ownership boundaries

| Boundary | Owner | Scope | Responsibility | Forbidden coupling |
| --- | --- | --- | --- | --- |
| BND-001 | P1 | `internal/catalog/**` | catalog model, ISO identity, capitals, boundary-profile metadata | renderer or CLI policy |
| BND-002 | P1 | `data/**` | pinned source receipts and immutable compiled corpus | runtime source downloads |
| BND-003 | P2 | `internal/geometry/**` | projection, fitting, simplification, smoothing, feature retention | CSS/style serialization |
| BND-004 | P3 | `cmd/**` | executable and command surface | direct geodata parsing |
| BND-005 | P3 | `internal/config/**` | schema, presets, inheritance, overrides, diagnostics | geometry mutation |
| BND-006 | P3 | `internal/render/**` | minimal SVG serialization and style/pin contracts | source acquisition |
| BND-007 | P4 | `docs/guides/**` | producer and client integration/performance guidance | release automation |
| BND-008 | P4 | `test/e2e/site/**` | full-catalog browser and site-use fixtures | corpus compilation |
| BND-009 | P5 | `.goreleaser.yaml` | portable artifact matrix | product behavior |
| BND-010 | P5 | `docs/operations/**` | installation, handoff, refresh and release operations | client integration truth |

## Contracts and interfaces

| ID | Contract | Producer | Consumers | Semantics |
| --- | --- | --- | --- | --- |
| CTR-1 | corpus bundle | P1 | P2 | versioned manifest plus canonical geometry/capital records; content digest is identity |
| CTR-2 | geometry result | P2 | P3 | normalized paths, viewBox and optional projected markers; presentation-free and deterministic |
| CTR-3 | config | P3 | agents/producers | versioned YAML or JSON; presets inherit; global and per-country overrides; unknown fields fail |
| CTR-4 | CLI | P3 | agents/producers/P4 | `init`, `validate`, `preview`, `generate`, `explain`, `inspect`, `version`; stable exit classes and `--json` |
| CTR-5 | SVG | P3 | browsers/P4 | `standalone` or `themed-inline`; stable classes, data attributes, CSS custom properties and accessible metadata options |
| CTR-6 | catalog manifest | P3 | site/build pipelines | one machine-readable index with ISO id, files, profile, byte count, dimensions, pin state and digests |
| CTR-7 | release archive | P5 | colleagues | one platform binary, checksums, schemas, presets and operating guide |

Compatibility is additive within a schema major. Removing a CLI field, CSS hook,
manifest field, or preset token requires a new schema major or an accepted
superseding decision.

## Quality attribute budgets

| ID | Attribute | `card` | `hero` | Gate |
| --- | --- | --- | --- | --- |
| QAB-1 | Intended CSS size | 90–160 px, up to 240 px accent | 500–700 px | browser fixtures |
| QAB-2 | SVG bytes | median ≤0.6 KB; p95 ≤1.2 KB; max ≤2.5 KB | median ≤1.8 KB; p95 ≤4 KB; max ≤8 KB | full-catalog manifest gate |
| QAB-3 | Persistent animation | none | opt-in, transform/opacity only | reduced-motion/browser gate |
| QAB-4 | Capital pin | off | optional | config/render tests |
| QAB-5 | Runtime network | none | none | hermetic E2E |
| QAB-6 | Catalog count | all currently assigned ISO 3166-1 entities | same | corpus gate |

Budgets are uncompressed file sizes. A first full-catalog baseline may tighten
them; exceeding a maximum blocks delivery unless an explicit per-country exception
names the reason and remains within an accepted exception ceiling.

## Data state and consistency

The corpus is immutable and content-addressed. It contains source version/license
receipts, ISO mappings, geometry variants required by boundary profiles, protected
feature metadata, and capital records. A refresh creates a new corpus version and
comparison report; it never mutates an existing version.

Generated output is disposable and reproducible. A generation writes to a staging
directory and replaces the selected output set only after validation succeeds.
Configuration is user-owned and never rewritten implicitly after `init`.

## Failure recovery and rollback

- Invalid source/corpus data fails during corpus compilation, before it can enter
  normal generation.
- Invalid config, unknown ISO codes, impossible pin overrides, output collisions,
  and budget violations return typed diagnostics and non-zero status.
- Catalog generation is transactional at the output-directory boundary; a failed
  run leaves the last successful set usable.
- Rollback means select the prior binary/corpus/config and regenerate; output
  manifests expose all three identities.
- A visually disputed boundary or protected-feature regression rolls back the
  corpus version, not presentation code.

## Security and privacy

The product processes public geospatial data and local configuration; it stores no
personal data or credentials. Normal generation has no network capability. Paths
are resolved beneath explicit input/output roots, generated names derive from
validated ISO identifiers, and SVG output forbids scripts, event attributes,
external references, embedded raster data, and unsafe XML constructs.

The owner accepts political-boundary policy and any named exceptions. Maintainers
accept source refresh receipts. P4 proves that client-side use does not introduce
unsafe inline markup.

## Observability and operations

Human output is concise and actionable; `--json` emits stable diagnostics for
agents. Every run reports generator version, corpus version, config digest, profile,
country counts, byte statistics, warnings and output manifest path. `explain`
resolves inheritance and per-country decisions without writing assets. `inspect`
reports geometry/marker facts for one ISO code.

## Delivery waves and ownership

| Wave | Owner | Exit |
| --- | --- | --- |
| W1 | P1 | trusted, complete, versioned corpus contract |
| W2 | P2 | deterministic soft-organic geometry for both profiles |
| W3 | P3 | usable local agent-first generator |
| W4 | P4, P5 | site-ready full catalog and, independently, colleague distribution |

Critical path: P1 -> P2 -> P3 -> P4

P5 may run beside P4 and cannot block the
owner's first production use.

## E2E allocation

| Journey | Lowest proving boundary | Due |
| --- | --- | --- |
| Compile pinned source into a complete corpus and reject a missing/duplicate ISO entity | P1 integration | W1 |
| Render representative geometry and match pinned D3/path fixtures | P2 integration | W2 |
| Run `init → validate → generate → inspect` through the built binary in one working tree | P3 E2E | W3 |
| Generate the full catalog twice and prove digests, count, structure and byte budgets | P4 E2E | W4 |
| Render representative inline/mask assets in light, dark, reduced-motion and card-density fixtures | P4 browser E2E | W4 |
| Install each packaged binary in a clean environment and execute the product loop offline | P5 package E2E | W4 |

## Rollout migration and compatibility

The owner-first rollout is: representative fixtures, local full catalog, site
integration, then cross-platform packaging. `hero` and `card` are the only MVP
detail profiles. Public repository polish, hosted services, plugin surfaces and
automatic upstream refresh are deferred until the local workflow is accepted.

Existing configuration remains valid across additive releases. Corpus upgrades are
explicit. Generated assets may be replaced only with a new manifest and successful
full-catalog gate.

## As-built reconciliation

Implementation plans must cite boundary IDs. Any new package crossing an ownership
boundary, new runtime dependency, changed budget, CSS hook, CLI contract, catalog
authority or boundary posture requires an amendment or superseding decision before
closeout. P4 records measured full-catalog distributions; architecture budgets are
updated only from that accepted evidence.

## Decisions and unresolved forks

- DEC-001: elevated risk policy governs architecture, performance, accessibility,
  public-contract and geopolitical decisions.
- DEC-002: `un` is the default boundary profile and `de-facto` is configurable.
- DEC-003: Go renderer over a versioned immutable geometry corpus.

No unresolved product fork blocks planning. `detail` output, automatic upstream
refresh and public distribution polish are deferred backlog, not hidden scope.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
