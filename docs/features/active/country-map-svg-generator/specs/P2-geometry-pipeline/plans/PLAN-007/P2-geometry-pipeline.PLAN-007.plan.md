---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "accepted-plan-revision"
id: "PLAN-007"
epic: "country-map-svg-generator"
spec: "P2"
status: accepted
profiles: []
concerns: []
inputs: ["P2", "PLAN-006", "RUN-006-RESULT", "RESULT-013", "DEC-003", "DEC-004", "DEC-005", "AM-001", "AM-002", "CTR-001"]
---
# PLAN-007 — Bounded-phase topology-safe fixed-grid finalization

## Accepted inputs and baseline

This successor plans P2 from baseline commit
`0306704d328b487a4e558553d00646b2ab7d112b`, tree
`f412b9f5fb8356d990d313d342cafcd8d0175575`. It is bound to:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`
  and architecture
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- invalidated PLAN-006
  `5df964761f2591b49776ca7e93209e2c1dbb45a22e3e6771f0fad2cd7db196d3`
  and manifest
  `b62f645fb2ae84a3a8ff07a670299eee1764bedd0c8d5cc5d1264331881ff190`;
- interrupted RUN-006
  `3d7bf4a812ff64af4509f08ebd9a6186f29e705967afb2add78a9c20e7bdbd6e`,
  RESULT-013
  `f30f3134fb21e9364d2283ddde2aee1889fac6d5abb203733b3b6343cb627a78`,
  and PLAN-006 invalidation receipt
  `f2704508c447a1b38dacc94d318d75e2d3b119c9a070408fbb067efc1e7f9ccb`;
- DEC-003 `f1002bdfbcf6947d11faa339161445614c2ae93ee340c2f97fa90a4036e7261d`,
  DEC-004 `3efec49a55fe7289ea528233280671e43db5f114c24586c1780dba98ba2667e5`,
  and DEC-005
  `77d1b69ad5e49ecf8324af4250de8f39f24c3b9547d6001d69c09dc4c06f2bc3`;
- AM-001 `ab284fd6e2a79c4ef7a8c45ab05fe807d217ce7a8a178be354a44d26ef4ac80a`
  and AM-002
  `ec2203d7132e467c7b1519d765c4790d583c29691311aad208a2ce6092523487`;
- P1 RUN-001
  `de1aacd394ec4a00fdfedd6aae532c53663ce60f00dbddb23d795c0a59a8ce90`
  and RUN-002
  `eb64507d5b9c108bb25492c8d2091c3335ca968e73485571783e62f302b99282`;
- verified P1 product checkpoint
  `ed4fdca23965b6aa0c0fc1b8e4c449e79f2d45ee`, tree
  `1a6bfca82d26027956f92e28c9cd479c277faba5`, corpus identity
  `sha256:9d56b4d205eb2f5979e0d1f82d84222946b9f4e12e8a2b7ec09204e2225f2d95`,
  manifest SHA-256
  `227374e94910f078623b3c643085d39fa094b4e1ae20a0ffc01cb2d4a19ee616`,
  and geometry SHA-256
  `f3d5d80ac8b9650cb75f68f46593c24476a979855875842bc3274d9d928e7744`.

PLAN-006 remains invalidated. Its uncommitted product scratch is preserved and
has no authority except where this plan explicitly re-leases and eventually
stages it. At planning time `go.mod` is
`4412d79908d733d715bad8c17334eafbdabb6f96ab8b242d87bdf6637ee9f13b`,
`go.sum` is
`139a8de4191f92748a9d3a8c3868f56862e42f1b0d60d8f1dc4a582d273b5daf`,
and the ordered SHA-256 inventory of current `internal/geometry/**` files
excluding `node_modules` is
`9218d841dedfb9c2a1a76087f0eb94a7015fdc43a38dda843f1224f9be4648fc`.
RESULT-013 supplies the exact per-file hashes for its nine changed paths and
proved the centered, zero-rotation LAEA planar artifact/runtime binding and AQ
unlimited representation path. T0 must reconcile those hashes, the complete
pre-dispatch inventory and focused behavior before staging; it does not pretend
the scratch belongs to baseline commit `0306704`.

The new load-bearing fact is narrower. Forced-source RU/card at `q=0.01` failed
final topology at all 14 monotonically halved tolerances from
`1.5289289321881343` through `0.00018663683254249686`. The valid pre-grid
geometry therefore cannot be made representable merely by selecting a finer
line-simplification tolerance. Independent point rounding is the failing
operation: it can collapse near edges onto the same grid cells and create
crossings even when the pre-grid polygon and simplification are valid.

The serializer is frozen at `internal/geometry/serialize.go`
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and `serialize_test.go`
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.
PLAN-007 neither edits it nor moves byte responsibility after serialization.

### Framework search outcome and bounded successor

The current pure-Go dependency
`github.com/peterstace/simplefeatures/geom` v0.59.0 exposes
[`SnapToGrid`](https://pkg.go.dev/github.com/peterstace/simplefeatures/geom#MultiPolygon.SnapToGrid),
whose own contract says a valid polygon may become invalid. That is equivalent
to the failed operation, not a solution. The same module contains a JTS
snap-rounding precision reducer which promises valid, fully noded polygonal
output, but it is under `internal/jtsport` and cannot be imported. Its public
overlay API is robust but does not expose a fixed precision model. The canonical
JTS
[`GeometryPrecisionReducer`](https://locationtech.github.io/jts/javadoc/org/locationtech/jts/precision/GeometryPrecisionReducer.html)
defines a possible governed fallback, but JTS/CGO/new runtime dependencies are
outside the accepted contract.

The accepted authority does not require the fitted coordinate system or viewBox
to start at absolute zero. P2 and AM-001 require deterministic natural dimensions
or the requested contain dimensions, uniform scale, centering and a common
geometry/marker transform. DEC-005 requires `q=0.01` and reserves the
phase-independent maximum rounding displacement `q/sqrt(2)`; it does not fix the
grid's phase relative to the projected geometry.

PLAN-007 therefore tries a smaller framework-free operation first: select one
deterministic global grid phase and compose it exactly once into the candidate's
fit translation. That one resulting translation is the authority for fitted
candidate and full-reference coordinates, Result.Transform, shifted ViewBox,
protected checks and later marker projection before the unchanged zero-origin
`q=0.01` point quantizer and frozen serializer run. This changes no edge
relative to another edge before quantization and performs no per-vertex repair.
In rendered viewBox coordinates the phase cancels, so the only visual
displacement remains ordinary nearest-grid rounding bounded by `q/sqrt(2)`.

If the bounded RU/AQ phase experiment exhausts without a valid representation,
the plan stops. Only then is the unsupported public precision-reducer gap a
resumable governance/dependency blocker. A fork, vendored internal JTS copy,
`replace`, CGO/GEOS, generic union/buffer/MakeValid or project-local snap-rounder
remains forbidden.

## Requirement and acceptance coverage

| Obligation | Tasks | Evidence |
| --- | --- | --- |
| US-1, AC-1 | T0, T2 | RU/AQ experiment, final budgets and owner-approved card/hero sheet |
| US-2, REQ-1, REQ-2, REQ-8, AC-2 | T0, T2 | natural `tight`, arbitrary `contain`, scale-only selection, q-grid and deterministic goldens |
| US-3, REQ-6, AC-3 | T0, T2 | bounded overrides and generic protected/minimum-part retention after phase finalization |
| US-4, REQ-3, REQ-4, AC-4 | T0, T2 | bounded phase contract, final-representation topology/deviation and mutation teeth |
| REQ-5 | T0–T2 | exact omission, restoration, removal and phase/fallback provenance |
| REQ-7 | T0, T2 | unchanged full-source marker transform and anomaly evidence |
| INV-1 | T0 | presentation-free canonical geometry and frozen serializer |
| VAL-1 | T2 | pinned D3 parity |
| VAL-2 | T0, T2 | RU/AQ, fixed-grid phase, protection, softening and mutation teeth |
| VAL-3 | T0, T2 | all-corpus scale/layout/visual/budget matrix |
| VAL-4 | T2 | valid and invalid bounded override fixtures |
| VAL-5 | T2 | inside, edge and outside marker fixtures |

## Technical approach

Full P1 remains the sole authority. Preserve PLAN-006's proven centered LAEA
projection, zero build rotation, projected compact/standard artifacts, full-source
natural bounds, `tight`/`contain` fit, effective fitted scale, retained
components, protected anchors and marker transform. Runtime selection remains
only a function of effective fitted scale and guard results: compact resolution
`64` through scale `240`, standard resolution `256` through scale `700`, then
source. Country, boundary profile, preset and frame shape are forbidden selection
inputs.

Replace only the fixed zero-phase assumption. For each already valid
candidate-specific reduction/restoration result, search the finite phase set
`phaseX,phaseY ∈ {0.000,0.001,...,0.009}` in lexicographic
`(phaseX,phaseY)` order, with `(0,0)` first. The set contains exactly 100
generic phases inside one `q=0.01` cell and is independent of entity, profile,
preset, frame and tier.

Phase selection occurs after one tier candidate has completed reduction,
restoration and its raw fitted-space guard, but before canonicalization, final
deviation or any final-representation acceptance. For a phase `(dx,dy)`, derive
one phase-composed transform from the unshifted fit transform by setting its
translation to `(base.TranslateX+dx, base.TranslateY+dy)`. Materialize both the
candidate and the full projected reference used by the final guard from their
unshifted forms through that resulting transform exactly once. Do not translate
an already fitted geometry and then apply the phase-composed transform.

Use the same resulting translation for protected-anchor projection during the
phase guard, and derive one shifted ViewBox by adding `(dx,dy)` once to both
bounds. A failed phase or tier discards all of this candidate-local state. Only
after a tier passes topology, protection, final deviation and byte guards does
its selected phase-composed Transform and ViewBox enter Result; optional markers
are then projected once through that accepted Transform and checked against that
accepted ViewBox.

The serialized path therefore still contains coordinates on `0.01*Z`; only the
accepted ViewBox/Transform records the deterministic sub-q origin. Width,
height, natural aspect, effective scale, padding, containment and marker
position relative to geometry are invariant.

Accept the first lexicographic phase only when all of these predicates pass
together:

1. input is the valid projected/fitted/restored polygon candidate;
2. output is non-empty and valid under P2's topology/winding/hole guard;
3. every output coordinate lies on `q=0.01`;
4. every protected anchor remains covered and each declared `minimum_parts`
   obligation remains satisfied;
5. no phase changes component/ring structure; any earlier simplification or
   visibility removal retains its existing exact source provenance;
6. the canonical result remains contained by the shifted viewBox with the same
   padding/centering contract;
7. final full-source↔grid-result symmetric deviation, measured in the common
   shifted coordinate frame or equivalently after subtracting the phase, is
   finite and at most the
   unchanged resolved tolerance (`relative_ratio <= 0.012`, floor `0.80px`,
   explicit quality at most `1.25px`);
8. parser round-trip is exact and the frozen serializer's final bytes fit the
   generic `MaxPathBytes` (`2500` card, `8000` hero).

The existing monotonic source simplifier remains the outer search. The bounded
phase set is the inner search for each candidate, so the complete forced-source
bound is 14 known tolerance attempts × 100 phases, with early acceptance of the
first valid pair. The 14-run failure proves tolerance alone is not a topology
cure; the phase experiment tests whether a different global lattice origin can
represent the same valid geometry. Phase exhaustion is typed and terminal.

Keep Mapshaper 0.7.44, Node and npm maintainer-only for LOD rebuilding. Ordinary
generation and `make check` remain offline pure Go. Preserve centered-LAEA
artifact metadata, `64/256`, `240/700`, `q=0.01`, `0.012`, hard
`2500/8000`, advisory `2200/7500`, frozen serializer grammar, natural-aspect
`tight`, distortion-free arbitrary `contain`, and deferred `detail`.

### Public runtime wiring before publication

T0 owns a nil-safe published-table seam in `pipeline.go`/`lod.go`. Public
`Generate` always routes through `GenerateWithLOD` using the immutable published
table when present and `nil` otherwise. `nil` is a first-class source-only mode,
not an error and not a test bypass, so a clean checkout of the eventual T0
commit compiles and generates from P1 before any T1 artifact exists. Focused
tests inject a table directly into `GenerateWithLOD`; they never mutate global
published state.

T1 owns only the canonical compact/standard/manifest bytes,
`lod_generated.go`, and the artifact rebuild test. The generated adapter embeds
and validates the published table and initializes the T0 seam during package
initialization. It may not modify any T0 file. A clean checkout of T1 must prove
that public `Generate`, not only `GenerateWithLOD`, selects published tiers
offline. The clean rebuild executes the builder, tool and recipe already
committed by T0 and compares temporary output to the T1 artifacts.

## Tasks and completion conditions

1. **T0 — Prove and integrate bounded fixed-grid phase finalization.**

   First re-own the complete PLAN-006 T0 foundation plus the frozen serializer
   and new grid-phase files. Reconcile every existing byte against RESULT-013,
   the planning inventory and focused proof; remove probe-only behavior and wire
   the nil/source public `Generate` seam. The serializer is integration-owned
   but edit-frozen at `serialize.go`
   `c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
   and `serialize_test.go`
   `0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.
   The coder brief forbids changing those bytes; after their focused proof the
   coordinator stages them with the rest of T0.

   T0 then has two serial brake phases and produces no commit between them.

   **T0a, RU/AQ hard experiment:** add only the bounded phase helper,
   provenance fields and focused tests. Run forced-source RU/un/card and
   AQ/un/card+hero first with
   `MaxPathBytes=0`, then with `2500/8000`. Run AQ compact/standard regression
   through the proven centered-LAEA binding. For every case emit input validity,
   selected tier, attempted/selected tolerance and phase, q-grid proof,
   viewBox/transform/marker shift, component/ring identity,
   protected/minimum-parts results, raw and final
   symmetric deviations, points, frozen-serializer bytes, parser round-trip,
   determinism, phases attempted, per-phase elapsed time and generation runtime.

   T0a passes only if RU and AQ are valid on `q=0.01`, within the unchanged
   tolerance, deterministic, provenance-complete and within their hard budgets;
   AQ's proven projected-artifact behavior may not regress. If the bounded
   phase set exhausts, any output moves off-grid, relative geometry/marker/frame
   placement changes, a required part is lost, deviation/budget is exceeded, or
   serializer changes are needed, stop with no product commit. Do not add a
   dependency, union, buffer, MakeValid, custom noding or relaxed guard inside
   this plan.

   RUN-006 required runtime evidence but stopped at the representation failure
   before budget and broad stages, so it supplies no acceptable production
   latency receipt. On the recorded `YMBPM3` (`darwin/arm64`, Go 1.26) profile,
   T0a must stop if any single RU/AQ generation exceeds `10s` wall time or the
   complete post-build RU/AQ micro-brake exceeds `30s` aggregate wall time.
   Builder time is reported separately and cannot be used to dilute these
   runtime ceilings.

   **T0b, integrated production brake:** only after T0a passes, run all 283
   geometry IDs and all 249 entities under both boundary profiles, card/hero,
   long sides `90,128,160,240,500,640,700,1024,2048`, softening on/off, natural
   layout and `19x19`/`10000x12345` contain frames. Exercise compact, standard
   and forced source through the exact runtime path. Record selected tier,
   fixed-grid phase provenance, all component/ring changes, final deviation,
   protection, bytes, runtime, artifact hashes and two deterministic reruns.
   Produce the RU/CL/AR/AU/AQ/CA/ID and island-state visual sheet.

   Separately run the normal warm production batch—249 entities, card and hero,
   both boundary profiles, ordinary scale selection, no forced-source or
   multiscale maintenance expansion—and record total/per-output timing. It must
   complete within `120s` wall time on `YMBPM3`. The exhaustive forced-source,
   multiscale and mutation matrix is maintainer evidence with separately
   reported timing; its results may not substitute for or average away the
   normal-production ceiling.

   T0 completes only when all machine brakes pass and the owner approves the
   sheet. Otherwise stop without a T0 commit. After acceptance, stage the entire
   re-owned foundation and create one clean compiling T0 commit. In an isolated
   clean checkout of that commit, with no T1 artifacts or generated adapter,
   prove `GOWORK=off go test ./... -count=1` and public `Generate` nil/source
   behavior. T1 and T2 remain blocked until that clean-checkout receipt passes.

2. **T1 — Publish the frozen projected LOD corpus.** Only after a complete T0
   receipt, owner approval and clean T0 checkout, rebuild from the T0-committed
   builder/tool/recipe
   plus the accepted phase contract. Publish compact/standard artifacts,
   manifest and generated embed adapter for all 283 IDs. Prove exact P1, centered-LAEA
   metadata, recipe/lockfile/Mapshaper hashes, omission provenance, offline
   runtime and two byte-identical clean rebuilds. In a clean checkout of the T1
   commit, prove public `Generate` selects injected published tiers and remains
   offline. Publication may not modify T0 behavior or accepted artifacts in
   place.

3. **T2 — Allocate complete P2 evidence and project ceiling.** Add D3 parity,
   full-corpus, natural/aspect, arbitrary-frame, scale-only, fixed-grid phase,
   protected restoration, frozen serializer, budget, visual, override and
   marker evidence. Run all focused gates and offline `make check`; return exact
   product commit/tree and evidence for independent audit.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0 | `gridphase-foundation` | `go.mod`, `go.sum`, `internal/geometry/cmd/lodbuild/**`, `internal/geometry/lod/tool/**`, `internal/geometry/lod/v1.recipe.json`, `internal/geometry/testdata/lod-spike/**`, `internal/geometry/{adapter.go,adapter_test.go,diagnostic.go,fit.go,fit_test.go,gridphase.go,gridphase_test.go,lod.go,lod_test.go,marker.go,marker_test.go,math.go,model.go,model_test.go,normalize.go,normalize_test.go,override.go,override_test.go,path.go,path_test.go,pipeline.go,pipeline_test.go,preset.go,preset_test.go,projection.go,projection_test.go,retain.go,retain_test.go,serialize.go,serialize_test.go,simplify.go,simplify_test.go,soften.go,soften_test.go}`, `internal/geometry/overrides/**`, `internal/geometry/presets/**` |
| T1 | `lod-publisher` | `internal/geometry/lod/v1.compact.json`, `internal/geometry/lod/v1.standard.json`, `internal/geometry/lod/v1.manifest.json`, `internal/geometry/lod_generated.go`, `internal/geometry/lod_build_test.go` |
| T2 | `lod-verifier` | `Makefile`, `internal/geometry/approval_test.go`, `internal/geometry/integration_test.go`, `internal/geometry/mutations_test.go`, `internal/geometry/oracle_test.go`, `internal/geometry/testdata/approval/**`, `internal/geometry/testdata/d3/**`, `internal/geometry/testdata/mutations/**` |

T0 re-owns the complete uncommitted PLAN-006 prerequisite set because the
baseline contains none of that product foundation. Re-ownership permits
reconciliation, proof and eventual staging; it does not authorize needless
rewrites. No new runtime dependency is permitted. The serializer files are in
the T0 integration footprint only so the clean foundation commit contains the
verified implementation. They remain edit-frozen at the hashes above: the coder
does not write them, while the coordinator may stage their verified bytes.

No implementer writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle
documents, accepted plans, trackers, readiness, or another task's lease. The
coordinator owns the successor decision, lifecycle, evidence and audit records.

## Dependencies and execution waves

Critical path is T0a → T0b → T1 → T2, concurrency one.

- W0/T0a is the first executable step and has no product dependency.
- W1/T0b cannot start from representative success alone; it requires the exact
  passing RU/AQ receipt.
- W2/T1 cannot start until all-corpus T0, owner visual approval, one complete T0
  foundation commit and its nil/source clean-checkout receipt pass.
- W3/T2 cannot start until canonical artifact publication, public wiring and
  clean-checkout rebuild are reproducible.

There is no safe parallel product mutation. Test-fixture preparation may be
read-only before its task, but it cannot create leased bytes early.

## Validation plan

- `geometry-gridphase-contract`: proves exactly 100 lexicographic phases,
  `(0,0)` first, path coordinates on `0.01*Z`, unchanged viewBox dimensions and
  relative geometry/marker placement, maximum visual rounding displacement
  `q/sqrt(2)`, one-time transform composition, candidate-local failed-phase
  isolation, deterministic selection and typed exhaustion.
- `geometry-t0-clean-checkout`: from the exact T0 commit in an isolated clean
  worktree with no T1 files, runs `GOWORK=off go test ./... -count=1`, proves
  public `Generate` routes through the nil-safe seam to source, and verifies the
  staged builder/tool/recipe and frozen serializer hashes.
- `geometry-lod-public-wiring`: focused injection proves
  `GenerateWithLOD(nil)` source behavior and injected-table selection; a clean
  T1 checkout proves public `Generate` observes the generated embedded table,
  selects compact/standard by scale and performs no external I/O.
- `geometry-gridphase-performance`: on recorded `YMBPM3`, records attempted
  phase count and per-phase time, rejects a single RU/AQ generation over `10s`
  or post-build micro-brake over `30s`, and rejects the ordinary warm
  249-entity card+hero both-profile production batch over `120s`. Maintenance
  matrix timing is reported in a separate field.
- `geometry-fixedgrid-ru-aq`: forced-source RU card and AQ card/hero unlimited
  then budgeted, plus AQ compact/standard regression; asserts q-grid, polygonal
  validity, selected phase, exact provenance, protection, final deviation,
  round-trip, bytes and determinism.
- `geometry-final-representation`: all tiers use one final predicate and cannot
  report success before phase canonicalization, final topology, protection,
  symmetric deviation and serialization budget.
- `geometry-lod-integrated-brake`: all-283/all-249/both-profile scale/layout
  matrix with `64/256`, `240/700`, `q=0.01`, `0.012`, `2500/8000`, softening
  states, source forcing, artifact/runtime hashes and owner visual digest.
- `geometry-lod-rebuild-check`: two clean maintainer rebuilds from pinned
  Mapshaper 0.7.44 and lockfile using the T0-committed builder/tool/recipe equal
  each other and T1 committed artifacts.
- `geometry-svg-serializer`: frozen hash, parser round-trip and deterministic
  minified bytes; this gate may diagnose but not rewrite the serializer.
- `geometry-topology-teeth`: RU crossing witness, proper crossings, shared
  endpoints, winding, holes, grid collapse/merge/split, protected anchors,
  minimum parts and typed phase exhaustion.
- `geometry-lod-runtime`: centered-LAEA metadata binding, scale-only selection,
  fallback/cache parity and offline deterministic execution.
- `geometry-d3-parity`, `geometry-overrides`, `geometry-markers`, focused
  `GOWORK=off go test ./internal/geometry -count=1`,
  `GOWORK=off go test ./... -count=1`, and offline `make check` retain their
  accepted P2 meanings.

Teeth fail if the phase set/order changes; the phase schedule or branch logic
inspects country/profile/preset/frame/tier; phase is applied both to fitted
coordinates and again through Transform; failed phase/tier state leaks into
Result; geometry, full reference, protected projection, ViewBox or markers use
different translations; width/height/aspect/centering changes; path coordinates
move off-grid; provenance is invented; a protected or minimum part disappears;
final deviation is measured before phase or is zero/unmeasured/over tolerance;
per-phase/production timing is absent, averaged with the maintenance matrix or
exceeds `10s`/`30s`/`120s`; q, tier, threshold, budget or serializer changes; a
dependency/Node/CGO enters runtime; the T0 clean checkout requires T1 files;
public `Generate` bypasses the seam or ignores the T1 table; generated T1 code
modifies T0 files; a rebuild uses uncommitted builder/tool/recipe bytes; or
T1/T2 starts after only the RU/AQ sample.

## Deviation and amendment policy

`p2-bounded-gridphase-how-v1` permits only the exact 10×10 lexicographic global
phase set, common geometry/viewBox/transform/marker translation, additive
attempted/selected phase provenance, and helper/test splits inside T0's lease.

A larger/finer/adaptive phase set, different phase order, dependency change,
fork, `replace`, internal import, vendored JTS, custom noder, MakeValid,
union/buffer repair, CGO, changed projection/artifact binding, q, tolerance,
`64/256`, `240/700`, budget, serializer grammar, protected policy, selection
input, machine profile, or relaxation/reclassification of the
`10s`/`30s`/`120s` performance ceilings requires replan; a
change to natural/contain semantics, P1 authority, presentation-free CTR-002 or
offline single-binary runtime requires an amendment/decision.

Any T0a failure returns the exact exhausted tolerance×phase matrix to decision
review. It does not widen the phase lattice or try another algorithm.

## Commit worktree and integration policy

Each completed task lands as one conventional atomic product commit after its
complete gate. Stage exact leased paths only. T0a produces no standalone commit;
T0 commits the entire geometry foundation only after T0b and owner approval.
The coordinator stages the edit-frozen serializer only after hash/proof checks.
Preserve unrelated scratch and reconcile every reused PLAN-006 byte by exact
pre-existing hash, baseline diff and focused behavior.

Before integration, the coordinator checks expected parent, actual path scope,
dependency sum, frozen serializer hash and gate receipt. Lifecycle/evidence
records never enter product commits. No task may amend an earlier task's files;
discovery of such a need returns to that owner and reruns all downstream gates.
The T0 commit is not released to T1 until its exact commit/tree is checked out
cleanly and the nil/source receipt passes. The T1 commit receives the same
clean-checkout proof for published public routing and deterministic rebuild.

## Rollback and recovery

Before T0a passes, rollback is deletion of only experiment-owned uncommitted
bytes; existing PLAN-006 scratch stays untouched. A blocked experiment resumes
from the exact RU/AQ tolerance×phase receipt—not chat history.

At runtime, invalid compact advances to standard; invalid standard advances to
phase-finalized source; phase-exhausted, invalid, over-deviation or over-budget
source fails typed. There is no second grid or repair fallback.

T1 stages artifacts in an explicit temporary directory and publishes atomically
only after complete checks. Revert completed product commits in T2 → T1 → T0
order and select the prior algorithm/LOD pair. Any accepted artifact or
dependency change creates a new version; accepted bytes are never regenerated
in place.

Reverting T1 removes only artifacts, rebuild test and generated initializer;
the T0 public seam remains compiling and deterministically falls back to
nil/source. Reverting T0 then removes the complete geometry foundation back to
baseline `0306704`, rather than leaving unstaged prerequisite fragments.

## Completion and handoff

P2 closes only after the passing RU/AQ phase experiment, owner-approved visuals,
deterministic all-283 projected artifacts,
both-profile all-249 coverage, exact phase/restoration/removal provenance,
natural and arbitrary layouts, final paths within `2500/8000`, offline
`make check`, passing `10s`/`30s`/`120s` performance receipts on `YMBPM3`,
nil/source and published-table clean-checkout receipts, immutable run evidence
and fresh independent audit.

P3 receives CTR-002 with natural `tight`, distortion-free `contain`, versioned
data-driven presets, canonical `q=0.01` presentation-free commands, deterministic
viewBox origin and phase provenance, stable diagnostics, removals and optional
transformed capital markers. P3 does not repair or optimize paths. `detail`
remains deferred.

Estimated effort is 7–16 agent-hours, likely 11, medium-low confidence. The
RU/AQ experiment is intentionally 1–3 hours and stops before the all-corpus cost
if the bounded phase set cannot represent the witnesses.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
