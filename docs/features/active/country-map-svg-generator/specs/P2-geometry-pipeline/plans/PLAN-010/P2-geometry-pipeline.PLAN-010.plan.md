---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "accepted-plan-revision"
id: "PLAN-010"
epic: "country-map-svg-generator"
spec: "P2"
status: accepted
profiles: []
concerns: []
inputs: ["P2", "DEC-006", "DEC-007", "PLAN-009", "RUN-010-RESULT", "RESULT-018", "RUN-011-RESULT", "RESULT-019"]
---
# PLAN-010 — Scale-aware protected-visibility silhouette ladder

## Accepted inputs and baseline

Implement P2 from scaffold baseline
`2a7d45953ac857af185a4825ceea7097e0af3a40`, tree
`a61a0f87122e9094f444fd032dc54a2c04ad3c4b`. Lifecycle-only commits after
that baseline do not authorize product-byte drift.

Bind:

- P2 `4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`;
- architecture
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- accepted DEC-006
  `100cfe93e270fbf39b3ed61b78a7a58b8e9dd7f69372acb24a8316d4cc04620f`
  and acceptance receipt
  `01776bab76c5feec35bac816608ab52c6f09af12fabafd1c5fe3660d4514ba05`;
- accepted DEC-007
  `e927675a05f1cee0f4e39532a4f600609fa9560f0e97eca318d7af1ba231e9ef`
  and acceptance receipt
  `3d93aaa6a120e66e6d7d6ddfc76d84ac8b704bfe13c25b326c13ffb0f767c0fc`;
- invalidated PLAN-009 revision
  `a406405b5895244bba918d5c9a1cb913a1144ae08a7785ea441f58cf87baeea8`,
  manifest
  `7fcaf58708caeba7bba7097132c86b898c8bf74d650b9efc91934cf512698e75`,
  and invalidation receipt
  `7a174d2a4fc65cc1317dcc659bda84e1b6ef69a4c54830e407322cea34ac2d15`;
- RUN-010-RESULT
  `264b1af5262d3c4f88fe9e62f24f29716061067d2e0df87a8571127759914579`,
  RESULT-018
  `e637a82c748a648814fe74ae5d8ea8723cba4d3c0758beb98c18df46d9b028c5`,
  VR-VAL-2
  `a501e291f904561787b22aff692f16636e898cd7b9fdccdac16770688635817c`,
  and VERIFY-533A8DF7AE6B
  `b6d913896075cbb5d946fae001f6bf067011dd5bce174b1045370555b9476245`;
- RUN-011-RESULT
  `011ca94ddd060d4f80d1da6894127c6210bf1317a572d861bce0f77ec09f1b32`
  and RESULT-019
  `f1f0f03c518549c5ea1d6f91dbd30755bed7208538dc71a9f0a3cec48f062d06`.

PLAN-009 remains invalidated. Everything in it remains binding successor intent
except the unconditional automatic `minimum_parts` rule superseded by DEC-007.
RUN-010's diagnostics and RUN-011's blocker are authoritative evidence, while
their product bytes remain non-authoritative scratch that T0A must re-own.

At planning time `go.mod` is
`4412d79908d733d715bad8c17334eafbdabb6f96ab8b242d87bdf6637ee9f13b`,
`go.sum` is
`139a8de4191f92748a9d3a8c3868f56862e42f1b0d60d8f1dc4a582d273b5daf`,
and the ordered SHA-256 inventory of the 44 current
`internal/geometry/**` files excluding `node_modules` is
`4132db98460bedc1b898f675e41548a092fd0e5d8dca972e950ffb6b497bd291`.
The inventory is the SHA-256 of `sha256sum` lines for lexically sorted paths.

The current RUN-011 scratch is additionally pinned at:

- `internal/geometry/cmd/lodbuild/diagnostic_test.go`
  `1d780e5442a47ad3fe919b587f5d02ca893c375339077f4caca0c4dbb0aee977`;
- `internal/geometry/cmd/lodbuild/main.go`
  `6d4c9366e06b44cdb0963308c790a436eaed9c3334218d88f146d9eac3fc6aca`;
- `internal/geometry/cmd/lodbuild/representative.go`
  `c7d2d3b9b2cadeee30430db5db34741b60fd4c8bc0ba3b4872f02ff6c304301c`;
- `internal/geometry/lod/silhouette-oracle.v1.json`
  `66e9194b86a8ae79f6ce8570e75dbe970180c02f709f9377c711dfe275f92308`;
- `internal/geometry/lod/v2.recipe.json`
  `f11865346e9e78262dee148ed9d15eff416aeeabeb0d767d604a6e96401c58fc`;
- `internal/geometry/silhouette.go`
  `ff157dfb42cc63c5207b8e37b4380b3a3385d26f78cc15465a3384423d76cec7`;
- `internal/geometry/silhouette_test.go`
  `52ed5f9ca03fa2c54312e0b5aa7ebf851b772cba386300df26f6deb8c1bd7875`.

The serializer remains edit-frozen at `internal/geometry/serialize.go`
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and `serialize_test.go`
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

## Requirement and acceptance coverage

| Obligations | Tasks | Evidence |
| --- | --- | --- |
| REQ-2–5, REQ-8, VAL-2 | T0A | reproduced RUN-010/RUN-011 brakes, frozen oracle and protected-visibility provenance, machine-passing representative sheet and owner approval |
| US-1–3, REQ-1–8, INV-1, AC-1–3, VAL-2–3 | T0B | full 996-output gate, budgets, layouts, performance and owner catalog sheet |
| REQ-5, REQ-8 | T1 | deterministic all-283 ladder publication and public wiring |
| US-4, AC-1–4, VAL-1–5 | T2 | parity, mutations, overrides, markers, visual evidence and project ceiling |

## Technical approach

### Preserve the generic fitted-scale ladder

Retain PLAN-009's DEC-006 implementation: `card` and `hero` are versioned data
presets, while runtime candidate selection is a pure function of effective
fitted geometry long side, resolved quality and resolved byte policy. `tight`
and arbitrary distortion-free `contain` layouts share that path. Preset name,
entity identity, contain-frame shape and presentation inputs are forbidden
selection inputs.

The maintainer builder derives one fine-to-coarse generic ladder from immutable
P1 geometry with pinned Mapshaper presimplification. Every automatic candidate
must pass topology, winding, deterministic `q=0.01` coordinates, oracle
IoU/recall, protected visibility, provenance and path caps of `2200` for
card-scale policy and `7500` for hero-scale policy. Complete uncompressed SVG
maxima remain `2500/8000`; deterministic `300/500` envelope estimates remain
the P2 brake pending P3/P4 measurement. No q, serializer, candidate-order,
budget, runtime, projection, natural/contain, marker or style contract changes
are authorized.

Explicit `source` and custom-quality requests retain P1's full declared
`minimum_parts` through DEC-005, including its q-aware deviation guard and typed
hard-budget failure. Only automatic derived silhouettes use the following
scale-aware rule.

### `protected-visibility/v1`

For each source feature, deterministically form an identity set from the
`minimum_parts` largest source polygon components by reference filled-pixel
contribution. Every anchor-containing component is forced into the set. Ranking
is descending contribution, then source polygon order, then component content
digest. The base set has exactly `min(minimum_parts, source component count)`
members: when a declared anchor is outside the initial top-N, it replaces the
lowest-ranked non-anchor member. A later owner-approved group-anchor override
adds its component without evicting a base member. The dominant identity
component and every protected-anchor component are mandatory at every band
whether or not their contribution is subscale.

Contribution is the integer filled-pixel count obtained by rasterizing each
source component independently with the same projection, natural bounds,
padding, uniform fit, pixel-center sampling, SVG nonzero fill rule and maximum
effective scale grid as `silhouette-oracle/v1`. There is exactly one global
minimum contribution threshold per generic scale band. The recipe freezes its
integer unit, grid identity and value; entity, preset, frame and delivery inputs
cannot influence it.

An identity-set member whose reference contribution meets the band threshold
must retain source-component lineage and a nonzero source-intersecting candidate
pixel contribution on that grid. A member below the threshold may be absent
only as `protected_subscale`; it remains protected in diagnostics and at larger
bands. Larger effective-scale bands may not lose an identity member retained
by a smaller band. Fine-to-coarse search may omit only these subscale members
and ordinary unprotected components. It may not invent, move or reinsert
vertices.

For every source component, the oracle, recipe, ladder and CTR-002 diagnostics
record source order and digest, contribution, identity-set rank/reason,
dominant/anchor flags, band threshold, visible/subscale classification,
candidate lineage and pixel contribution, retained/omitted status, omission
reason and tie-break keys. Missing lineage, a visible protected omission or a
false `protected_subscale` label fails closed.

The existing bounded data-override seam may add an owner-approved named group
anchor only. It forces the containing source component at every band and
records review provenance. Tests prove that it can only add retention and can
change neither thresholds, candidate order, budgets nor executable control
flow. There is no country-specific quality or protection branch.

The oracle remains pure Go, offline and presentation-free: no antialiasing,
stroke, fill color, CSS, delivery mode or softening participates. Global
per-band IoU/recall pairs, dominant coverage, anchor coverage, topology,
budgets and owner review remain independent mandatory gates.

## Tasks and completion conditions

### T0A — regression, protected-visibility oracle and representative brake; no commit

T0A is serial, re-owns the pinned scratch and creates no product commit.

**T0A.1 — reproduce the frozen evidence.**

1. Re-run the RUN-010 diagnostic twice from its pinned inputs. Both canonical
   semantic outputs must contain the exact 268 parameter/phase rows and 996
   catalog rows, produce SHA-256
   `220aa8cff32ff85a15c62b45eeff79040b66d4e14401f4fe7b02dbf83aa67249`,
   retain the recorded counts/maxima, and leave every contradiction predicate
   false.
2. Re-run the RUN-011 Indonesia compact boundary under a diagnostic-only
   unconditional-minimum brake. Resolution `121` must be the coarsest tested
   candidate retaining `20` parts and emit exactly `2713` path bytes and `2993`
   estimated complete-file bytes; resolution `120` must retain `19` parts.
   This frozen regression cannot be called by production selection.
3. Stop on any input, row, digest, typed failure, count, byte or boundary drift.
   Do not reinterpret drift as threshold-calibration freedom.

**T0A.2 — implement and freeze `protected-visibility/v1`.** Only after T0A.1:

1. implement and mutation-test component contribution, identity-set ranking,
   lineage and provenance using the existing pure-Go oracle;
2. calibrate one contribution threshold for each generic band over RU, CA, CN,
   AQ, AE, ID, CL, AR, AU and island/microstate fixtures;
3. search only fine-to-coarse and prove every candidate topology, global
   IoU/recall, dominant/anchor coverage, `2200/7500` path caps, `2500/8000`
   estimates, monotonic band retention and deterministic rebuilds;
4. mutation-test dominant loss, anchor loss, rank/tie drift, visible protected
   loss, false `protected_subscale`, entity/preset/frame branching and an
   add-only group-anchor override;
5. exercise multiple fitted scales around every band boundary in natural and
   portrait/landscape/square contain frames;
6. machine-generate the representative sheet from the exact accepted
   candidates, visibly labeling retained and omitted protected components,
   contribution, threshold and omission provenance;
7. require the whole representative machine gate to pass before asking for
   owner review; then obtain owner approval of both global thresholds and the
   exact sheet digest;
8. freeze oracle, recipe, thresholds, identity-set provenance and sheet digest.

T0A stops with no product commit if the Indonesia regression drifts, no single
global threshold per band satisfies machine budgets and owner recognition, a
dominant/anchor component cannot fit, a supposedly subscale member remains
identity-defining, retention is non-monotonic, the machine sheet fails, the
owner declines it, selection inspects forbidden inputs, rebuilds drift, runtime
needs an external process, or serializer/budget/q changes are requested. Owner
approval cannot waive a machine failure. T0B cannot start before both the
machine-passing sheet receipt and owner approval exist.

### T0B — full-catalog integration and owner gate

Start only from the exact passing T0A receipts and owner-approved frozen sheet.
Integrate the remaining foundation without modifying T0A-owned diagnostics,
oracle, recipe, thresholds, provenance, search or selection. Any required
change returns to T0A and reruns its dependent evidence.

Generate all 996 ordinary outputs and prove:

- all 249 entities, both boundary profiles and both preset inputs, with no
  missing or duplicate output;
- natural layout and tiny, huge, portrait, landscape and square contain frames
  with uniform scale and fitted-scale selection;
- deterministic candidate identity, topology, winding, identity-set
  provenance, visible protected retention, `protected_subscale` omissions,
  oracle metrics, points, path bytes and repeated output;
- path distributions and maxima `≤2200/7500`, with deterministic complete-file
  estimates `<2500/8000`;
- normal warm 996-output production `≤120s` on recorded YMBPM3;
- ordinary Go generation uses no Node, Mapshaper, network or external files;
- owner approval of the full-catalog sheet and every versioned group-anchor
  data override.

No executable country/preset branch, hidden omission or oracle bypass is
allowed. Only after every machine gate and owner approval may T0B stage the
exact T0A/T0B union and create one atomic T0 commit. A clean checkout without
T1 artifacts must pass `GOWORK=off go test ./... -count=1`, public nil/source
behavior, typed source-budget failure, and frozen serializer/dependency hashes.

### T1 — publish the generic ladder

From the clean T0 commit, run two deterministic maintainer rebuilds and publish
the all-283 ladder artifact, manifest and generated table initializer. Bind P1,
projection, recipe, oracle, band thresholds, identity-set/component provenance,
tool, candidate order and omissions. T1 may not modify T0 files.

In a clean T1 checkout, public generation uses the embedded ladder, selects only
from effective scale/quality/byte policy, meets `2200/7500`, preserves
natural/contain and source/custom semantics, and performs no external I/O.

### T2 — allocate complete P2 evidence

Add D3 projection/layout parity, oracle, candidate-order and protected-visibility
mutations, topology teeth, overrides, markers, budgets, natural/contain and
visual evidence. Run focused gates, `GOWORK=off go test ./... -count=1`,
offline `make check` and a fresh independent audit.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0A | `protected-visibility-oracle` | `internal/geometry/cmd/lodbuild/**`, `internal/geometry/{diagnostic.go,gridphase.go,gridphase_test.go,lod.go,lod_test.go,model.go,model_test.go,override.go,override_test.go,preset.go,preset_test.go,projection.go,projection_test.go,retain.go,retain_test.go,silhouette.go,silhouette_test.go}`, `internal/geometry/{overrides,presets}/**`, `internal/geometry/lod/{v1.recipe.json,v2.recipe.json,silhouette-oracle.v1.json}`, `internal/geometry/lod/tool/build.mjs`, `internal/geometry/testdata/silhouette-oracle/**` |
| T0B | `geometry-foundation-integrator` | `go.mod`, `go.sum`, remaining foundation, public seam, full-catalog brake, owner fixtures, unchanged tool lock/package metadata and frozen serializer staging |
| T1 | `silhouette-ladder-publisher` | `internal/geometry/lod/v2.ladder.json`, `internal/geometry/lod/v2.manifest.json`, `internal/geometry/lod_generated.go`, `internal/geometry/lod_build_test.go` |
| T2 | `geometry-verifier` | `Makefile`, approval/integration/mutation/oracle tests and `internal/geometry/testdata/{approval,d3,mutations}/**` |

T0A has no commit authority. T0B may stage proven T0A bytes but cannot rewrite
them. Serializer files are integration-owned and edit-frozen at their hashes
above. Current `integration_test.go` and `mutations_test.go` remain T2-owned.

No implementer writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle
documents, plans, trackers, readiness or another task's lease. Node modules and
temporary diagnostics, sheets and artifact output are never staged.

## Dependencies and execution waves

Critical path is `T0A.1 → T0A.2 machine gate → representative owner approval →
T0B → T1 → T2`, concurrency one.

T0A.2 requires both exact frozen regressions. Owner review starts only after a
machine-passing representative sheet. T0B requires the frozen global thresholds,
sheet digest and explicit owner approval. T1 requires complete T0B evidence,
full-catalog owner approval, one T0 commit and its clean-checkout receipt. T2
requires reproducible T1 publication.

## Validation plan

- `geometry-diagnostic-reproduction`: exact RUN-010 semantic digest, rows,
  ordering, counts, maxima and contradiction predicates across two runs.
- `geometry-indonesia-120-121-regression`: exact legacy minimum-part boundary
  and `2713/2993` resolution-121 bytes without production reachability.
- `geometry-protected-visibility-v1`: common rasterizer/fit, per-component
  contribution, identity-set rank/tie, global band thresholds, lineage,
  monotonic retention and exact `protected_subscale` provenance.
- `geometry-scale-ladder-spike`: stable fine-to-coarse order, global
  IoU/recall, `2200/7500`, complete-file estimates, determinism and
  machine-passing owner sheet.
- `geometry-selection-policy`: equal effective scale/quality/byte policy yields
  equal candidate identity independent of preset, entity and contain frame.
- `geometry-source-contract`: source/custom retains full P1 `minimum_parts`,
  DEC-005 q-aware deviation and typed budget failure.
- `geometry-ladder-integrated`: all 996 ordinary outputs, layouts, component
  provenance, budgets, `≤120s` production and full-catalog owner sheet.
- `geometry-t0-clean-checkout`, `geometry-ladder-public-wiring`,
  `geometry-ladder-rebuild-check`, `geometry-d3-parity`, `geometry-overrides`,
  `geometry-markers` and offline project ceiling.

Teeth fail on stale P1/projection/recipe/oracle/tool identity, changed
rasterizer behavior, antialiasing/stroke/color/CSS/softening in the oracle,
threshold drift, rank/tie drift, missing lineage, false omission provenance,
dominant/anchor/visible-identity loss, non-monotonic band retention,
nondeterministic candidate order, invalid ring, oracle bypass, byte overflow,
preset/entity/frame branch, external runtime process, serializer drift or
complete-envelope overflow.

## Deviation and amendment policy

`p2-scale-aware-protected-visibility-ladder-v1` permits helper/test splits inside
an existing lease, bounded tuning of the recorded fine-to-coarse candidate set,
and global threshold calibration only before owner freeze.

Any frozen-regression contradiction returns to decision review. After freeze,
any threshold, identity-set algorithm, contribution unit, rasterizer, metric,
fill rule, visibility policy, band boundary, candidate order, `2200/7500` cap,
projection, q, serializer, runtime dependency or selection-input change
requires replan and a new oracle/recipe version. Architecture maxima
`2500/8000` can only tighten without an architecture amendment.

Custom vertex reinsertion, country/preset branches, runtime repair, higher
budgets, lower precision, post-hoc SVG optimization, MakeValid, union/buffer,
CGO and hand-moved coordinates are forbidden. An owner group anchor must remain
bounded versioned data and add retention only.

## Commit worktree and integration policy

T0A creates no commit. After T0B and both owner gates, stage only the exact
T0A/T0B paths and create one conventional T0 product commit. T1 and T2 each
create one later atomic commit after their complete gates.

Before integration verify expected parent, path scope, scratch and dependency
hashes, oracle/recipe/threshold/tool identities, frozen serializer and receipts.
Lifecycle/evidence records never enter product commits. Recheck T0 and T1 from
isolated clean checkouts. Preserve unrelated/user scratch and never use blanket
staging.

## Rollback and recovery

Before T0 completion, rollback removes only experiment-owned additions; the
pinned RUN-010/RUN-011 scratch is not broadly deleted. Resume from exact
diagnostic and Indonesia regression digests, oracle/recipe identities, threshold
matrix, machine sheet digest and owner receipt rather than chat history.

A regression drift returns to decision review. A T0A.2 failure stops before
broad generation and commits nothing. A T0B failure returns to the frozen owner
checkpoint and does not relax thresholds or budgets.

T1 publishes through a temporary directory only after two identical rebuilds.
Revert completed commits in T2 → T1 → T0 order. Reverting T1 leaves the T0
nil/source seam compiling with honest typed budget behavior.

## Completion and handoff

P2 closes only with both frozen regressions, `protected-visibility/v1`,
owner-approved representative and full-catalog sheets, global band thresholds,
identity-set provenance, deterministic all-283 ladder artifacts, 996 ordinary
outputs, natural/contain layouts, `2200/7500` paths, complete-file evidence
beneath `2500/8000`, offline `≤120s` production, clean T0/T1 checkouts and
independent audit.

CTR-002 hands P3 effective scale, selected generic band, recipe/oracle/policy
versions, per-component protected-visibility provenance, removals, similarity
evidence, path bytes, transform and optional markers. P3 may tighten P2 caps
from measured delivery envelopes but cannot expand architecture maxima.

Estimated effort is 10–22 agent-hours, likely 15, low confidence. The exact
regression and machine-plus-owner representative brake precede all broad
catalog work.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
