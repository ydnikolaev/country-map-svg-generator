---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "implementation-plan"
id: "PLAN-008"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["P2", "PLAN-007", "RUN-007-RESULT", "RESULT-014", "RUN-008-RESULT", "RESULT-015"]
---
# PLAN-008 — Projection-aligned LOD selection and bounded rendering

## Accepted inputs and baseline

Implement P2 from scaffold baseline
`d219c9bae0434df888d6c4763eb0b19c2e260b91`, tree
`9096ed613412015f3cc86d5f0615ef753b5a399c`. Lifecycle-only commits after
that baseline do not authorize product-byte drift.

Bind:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- invalidated PLAN-007 revision
  `813d3a6551b9c6b008900ccb344f72dcc18d98bcbc9a72461e1d3dbfbc2adfdd`
  and manifest
  `bdc2f69fb6cc5346e94c45533685d44e11f97b7d156b909da51be1854790338e`;
- RUN-007-RESULT
  `4f52fefcae3079403645b93a87a5e0712d97f3849419b50308dd85253f8b433f`
  and RESULT-014
  `c79b3c753e18fd1f6e3096424db66b06790497eca15bbf03c2b15a12baf1b59f`;
- RUN-008-RESULT
  `6a65475d5e719ef7600b811100e21ae202a65c83dc6a3f2afd49eea3e9df8698`
  and RESULT-015
  `fc15cc94e75c4b67e4d7263032ca224f6d2d50f62ff3fed2de05f3d661a5815f`;
- invalidation receipt
  `624635498aba7459806c92330f43593a584489f3196a50b59ef9e9444078e211`.

PLAN-007 remains invalidated. Reconcile, re-own and freshly prove the existing
non-authoritative geometry scratch; do not treat it as accepted product.
`go.mod` and `go.sum` remain
`4412d79908d733d715bad8c17334eafbdabb6f96ab8b242d87bdf6637ee9f13b`
and
`139a8de4191f92748a9d3a8c3868f56862e42f1b0d60d8f1dc4a582d273b5daf`.
The serializer remains edit-frozen at
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

## Requirement and acceptance coverage

| Obligations | Tasks | Evidence |
| --- | --- | --- |
| REQ-3, REQ-4, REQ-8, VAL-2 | T0A | projection/selection/phase/budget micro gates |
| US-1–3, REQ-1–8, INV-1, AC-1–3, VAL-2–3 | T0B | full matrix, budgets, performance and owner sheet |
| REQ-5, REQ-8 | T1 | deterministic all-283 artifact publication and public wiring |
| US-4, AC-1–4, VAL-1–5 | T2 | parity, mutation, override, marker, visual and project gates |

## Technical approach

### One projection-subdivision authority

The current runtime projects the full reference before resolving `AutoQuality`;
an unset flatness therefore becomes implicit `0.20`. The builder recipe,
artifact metadata and binder require `0.10`. The current AQ selection evidence
then rejects compact/card at raw deviation `1.537690736882369` versus `1.536`,
and standard/hero at `12.205557348488734` versus `7.68`.

Define one versioned projection contract shared by runtime full projection,
maintainer projection, recipe validation, artifact metadata binding and
deviation-reference construction. Version 1 is centered LAEA, zero build
rotation, y-down, subdivision flatness `0.10` and coordinate precision
`0.000001`. Automatic and preset generation use it exactly; the result exposes
the actual policy in deterministic provenance. Do not keep the implicit
pre-policy `0.20`.

Preserve compact/standard resolutions `64/256`, scale thresholds `240/700`,
q=`0.01`, automatic simplification ratio `≤0.012` with `0.80px` floor,
explicit cap `1.25px`, budgets `2500/8000`, natural `tight`, arbitrary
distortion-free `contain`, frozen serializer and pure-Go offline runtime.

### Geometry selection is independent from rendering

Split the runtime boundary:

1. `selectLODRepresentation` performs tier ordering, binding, restoration,
   topology/protection checks, source tolerance search, q-grid phase search,
   final deviation and deterministic provenance. It returns the first valid
   representation.
2. `renderLODRepresentation` receives that immutable selection, tries the
   already accepted bounded softening-to-linear fallback, serializes once with
   the frozen serializer and enforces the hard byte limit.

`SelectLODProvenance` invokes selection only. `GenerateWithLOD` selects once and
renders once. Remove serializer, path and `MaxPathBytes` from phase/tolerance
search. A selected representation that exceeds its byte limit fails with the
typed budget error; rendering never resumes tier, tolerance or phase search.
Compact advances to standard and standard to source only for representation
failure, not for byte-budget failure.

Preserve the exact generic phase lattice:
`phaseX,phaseY ∈ {0.000,…,0.009}`, lexicographic with `(0,0)` first. Compose
the selected phase exactly once through geometry, full reference, Transform,
ViewBox, protected anchors and markers. RUN-007's RU checkpoint
`(0.003,0.002)`, 33 attempts, under four seconds is an empirical regression
anchor, not a country-specific branch.

Forced-source RU card and AQ card/hero prove unlimited representation validity.
Source is not required to compress Russia's `477784`-byte valid representation
to the compact production budget. A bounded source request instead proves a
typed, fast budget failure after one valid selection, without re-entering
search. Real `2500/8000` acceptance applies to temporary compact/standard
production-table paths. Do not compare a missing budgeted path with an
unlimited path.

Before T1, the public nil table means source-only: unlimited/custom generation
may succeed and a preset may return the documented typed source-budget error.
After T1, the generated table must provide budget-passing production tiers.

## Tasks and completion conditions

### T0A — projection/selection micro brake; no product commit

Re-own only the micro footprint. Align the projection contract, separate
selection from rendering, preserve the exact phase algorithm and provenance,
and add a typed fail-fast budget path. Build temporary compact/standard tables
twice with pinned Mapshaper `0.7.44`; record full artifact, recipe, lock and tool
hashes.

Run:

- forced-source RU/un/card and AQ/un/card+hero with no byte limit;
- AQ/card requested compact and AQ/hero requested standard using rebuilt tables;
- representative RU production table paths;
- forced-source hard-budget failure;
- deterministic repeated selection and rendering.

Unlimited source cases must pass q-grid, topology, protection, component/ring
structure, raw/final deviation, parser round-trip and provenance. RU must
reproduce phase `(0.003,0.002)`, 33 attempts and `<4s`, or return the changed
checkpoint for review.

Production cases must actually select compact/standard—source fallback is not
acceptance—and produce frozen-serializer paths within `2500/8000`. The forced
source budget case must return the typed error within `10s` with no additional
tolerance/phase attempt. The post-build micro brake is `≤30s`; builder time is
reported separately.

Stop without product commit if projection authorities differ, compact/standard
still fails representative deviation or budget, the serializer would need an
edit, topology fails or a performance ceiling is exceeded. Failure returns to
DEC-004's presimplification metadata or additional precomputed-tier decision.

### T0B — full foundation, production matrix and owner approval

Start only from passing T0A evidence. Integrate the remaining foundation and
the byte-identical T0A result. A T0A-owned correction returns to T0A and reruns
downstream evidence.

With temporary production tables, validate all 283 geometry IDs, all 249
entities, both boundary profiles, card/hero, long sides
`90,128,160,240,500,640,700,1024,2048`, softening on/off, natural layouts and
representative `19x19` and `10000x12345` contain frames. Compact/standard
production paths must meet `2500/8000`. Unlimited source fallback proves
geometry; bounded source fallback either meets the limit or fails typed and
bounded.

The warm normal-table batch—249 entities, card and hero, both profiles—must be
`≤120s`. T0A cases remain `≤10s` each and `≤30s` aggregate. Generate the owner
sheet for RU, CL, AR, AU, AQ, CA, ID and island states from budget-passing
production paths.

Only after the full matrix and owner approval, stage the exact T0A/T0B union and
create one atomic T0 commit. In a clean checkout without T1 generated files,
prove `GOWORK=off go test ./... -count=1`, unlimited nil/source success, preset
nil/source typed-budget behavior, no runtime external dependency, and exact
dependency/recipe/builder/serializer hashes.

### T1 — publish compact/standard artifacts

From the clean T0 commit, run two deterministic maintainer rebuilds and publish
all-283 compact, standard and manifest artifacts plus the generated table
initializer and rebuild test. Do not modify T0 files. A clean T1 checkout must
route public generation through the embedded table, select only by effective
scale, meet card/hero budgets and perform no external I/O.

### T2 — close P2 evidence

Add D3 parity, topology mutation teeth, override, marker, natural/contain,
scale-only, phase, protection, restoration, budget and visual evidence. Run
focused gates, `GOWORK=off go test ./... -count=1`, offline `make check` and a
fresh independent audit.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0A | `lod-selection-micro` | `internal/geometry/cmd/lodbuild/{main.go,main_test.go}`, `internal/geometry/{diagnostic.go,gridphase.go,gridphase_test.go,lod.go,lod_test.go,model.go,model_test.go,preset.go,preset_test.go,projection.go,projection_test.go}`, `internal/geometry/presets/v1.json`, `internal/geometry/lod/v1.recipe.json` |
| T0B | `geometry-foundation-integrator` | `go.mod`, `go.sum`, remaining foundation, public seam, `internal/geometry/lod/tool/**`, catalog brake, owner fixtures, and serializer staging |
| T1 | `lod-publisher` | `internal/geometry/lod/v1.{compact,standard,manifest}.json`, `internal/geometry/lod_generated.go`, `internal/geometry/lod_build_test.go` |
| T2 | `lod-verifier` | `Makefile`, parity/integration/mutation/oracle/approval tests and `internal/geometry/testdata/{approval,d3,mutations}/**` |

T0A has no commit authority. T0B may stage proven T0A bytes but cannot rewrite
them. Serializer files are T0B integration-owned and edit-frozen at their bound
hashes. Current `integration_test.go` and `mutations_test.go` remain T2-owned.
No task writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle documents
or another task's lease.

## Dependencies and execution waves

Critical path is `T0A → T0B → T1 → T2`, concurrency one. Source validity alone
does not release T0B: T0A requires real production-table deviation and budget
acceptance. T1 requires full T0 evidence, owner approval and its clean commit;
T2 requires reproducible T1 publication.

## Validation plan

- `geometry-projection-contract`: runtime, builder, recipe, binder and deviation
  reference share exactly one versioned `0.10` policy.
- `geometry-selection-render-separation`: search invokes no serializer/budget
  predicate; selection-only and production share the same representation.
- `geometry-source-budget-fast-fail`: over-budget source returns the typed error
  without another phase/tolerance attempt.
- `geometry-gridphase-contract`: exact 100-phase order, one-time transform,
  q-grid, containment, protection, deviation and failure isolation.
- `geometry-fixedgrid-ru-aq`: source validity, RU phase anchor and actual
  compact/standard production acceptance.
- `geometry-gridphase-performance`: `10s/30s/120s` ceilings.
- `geometry-t0-clean-checkout`, `geometry-lod-public-wiring` and
  `geometry-lod-rebuild-check`.
- Existing D3, topology, override, marker, serializer and project gates retain
  their accepted meanings.

## Deviation and amendment policy

Local helper/test splits inside a task lease are allowed. Changes to q,
tolerance, phase set/order, `64/256`, `240/700`, `2500/8000`, natural/contain,
serializer bytes, projection policy values, P1 authority, runtime dependency
posture or performance ceilings require replan.

If aligned temporary artifacts still fail representative deviation or budget,
stop T0A; fail-fast behavior is not production feasibility. Return to DEC-004's
presimplification metadata or additional precomputed-tier options.

No country branch, adaptive phase set, higher tolerance, finer q, post-hoc
optimizer, MakeValid, union/buffer, CGO, internal JTS import, fork or additional
runtime is permitted.

## Commit worktree and integration policy

T0A creates no commit. After T0B and owner approval, stage only the exact
T0A/T0B paths for one conventional T0 product commit. T1 and T2 each create one
later atomic commit after complete gates. Recheck every released task from an
isolated clean checkout and preserve unrelated/user scratch.

## Rollback and recovery

Before T0 completion, remove only experiment-owned additions; never broadly
delete predecessor scratch. Resume from exact projection, representation,
artifact and timing receipts. A selected representation's budget failure is
terminal and typed; phase/source-tolerance exhaustion is terminal and typed.
Revert committed work in T2 → T1 → T0 order. Reverting T1 leaves the T0
nil/source seam compiling with documented preset budget behavior.

## Completion and handoff

P2 closes only with aligned artifacts, real production budgets, all-283 artifact
coverage, all-249/both-profile runtime coverage, owner-approved visuals,
deterministic natural/contain output, phase/restoration provenance, frozen
serializer proof, `10s/30s/120s` receipts, pure-Go offline clean checkouts and
independent audit.

Estimated effort: 8–18 agent-hours, likely 12, low confidence. T0A is a
2–4-hour hard brake; failure prevents broad implementation.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
