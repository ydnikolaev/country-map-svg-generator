---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "accepted-plan-revision"
id: "PLAN-009"
epic: "country-map-svg-generator"
spec: "P2"
status: accepted
profiles: []
concerns: []
inputs: ["P2", "DEC-006", "PLAN-008", "RUN-009-RESULT", "RESULT-016"]
---
# PLAN-009 — Scale-calibrated, budget-first silhouette ladder

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
- invalidated PLAN-008 revision
  `29084a1805bf957e884558d921581d2601a83e8423e77426d1449004d18bd13a`
  and manifest
  `3a69993c33f56a965902484e032a526f8c8e3c379a056a866b4a314ccc97989e`;
- RUN-009-RESULT
  `5d6f2979cbb5ca210391cb214981662f031fb9589ddf71e201a65bb2c36d2ce9`,
  RESULT-016
  `b44793dcc746e5e55018b49762ab88c29d8bcf063ef8afb70cf6e877ef5dd7bc`,
  and PLAN-008 invalidation receipt
  `c761c0c7425e982a498d32d231e002c8397f782a7ac8db022176053041704f41`.

PLAN-008 remains invalidated. Its aligned projection and selection/render
separation are useful non-authoritative scratch only. T0A must re-own and prove
every reused byte before staging. At planning time `go.mod` is
`4412d79908d733d715bad8c17334eafbdabb6f96ab8b242d87bdf6637ee9f13b`,
`go.sum` is
`139a8de4191f92748a9d3a8c3868f56862e42f1b0d60d8f1dc4a582d273b5daf`,
and the ordered SHA-256 inventory of the 37 current
`internal/geometry/**` files excluding `node_modules` is
`0659b28f949f3f0c8a7ca1021ce0dac9eedab2348035c687709c27c926d73c24`.

The serializer remains edit-frozen at `internal/geometry/serialize.go`
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and `serialize_test.go`
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

## Requirement and acceptance coverage

| Obligations | Tasks | Evidence |
| --- | --- | --- |
| REQ-2–5, REQ-8, VAL-2 | T0A | content-addressed diagnostic, frozen oracle contract, scale-ladder spike and owner threshold/sheet receipt |
| US-1–3, REQ-1–8, INV-1, AC-1–3, VAL-2–3 | T0B | full 996-output gate, budgets, layouts, performance and owner catalog sheet |
| REQ-5, REQ-8 | T1 | deterministic all-283 ladder publication and public wiring |
| US-4, AC-1–4, VAL-1–5 | T2 | parity, mutations, overrides, markers, visual evidence and project ceiling |

## Technical approach

### Generic fitted-scale silhouette ladder

Implement DEC-006 directly. `card` and `hero` remain versioned data presets
which resolve ordinary long-side, padding, quality and byte-policy inputs.
Preset names and entity identifiers are forbidden geometry-selection inputs.

The maintainer builder derives a generic candidate ladder from immutable P1 by
using pinned Mapshaper presimplification and one recorded fine-to-coarse
candidate order. Runtime selection is a pure function of effective fitted
geometry long side, resolved quality and resolved byte policy. `tight` and
arbitrary distortion-free `contain` therefore share the same selection path.

The builder selects the first candidate satisfying all of:

1. non-empty polygonal topology, winding and deterministic `q=0.01`
   canonical coordinates;
2. protected anchors, declared minimum parts and dominant identity geometry;
3. exact P1/projection/recipe/tool/oracle identity and removal provenance;
4. the frozen scale-band silhouette-oracle thresholds;
5. path bytes at or below `2200` for card-scale policy and `7500` for
   hero-scale policy.

The architecture's complete uncompressed SVG maxima remain `2500/8000`.
T0A and T0B use deterministic `300/500` envelope estimates as provisional
cross-spec brakes; P3/P4 later replace estimates with measured delivery-mode
envelopes and may tighten, never expand, the P2 path caps.

Raw full-coastline symmetric deviation no longer accepts or rejects automatic
derived silhouettes. DEC-005's q-aware deviation guard remains binding for
explicit `source` and custom-quality requests, including typed failure for an
incompatible hard budget.

### `silhouette-oracle/v1`

Freeze one pure-Go, presentation-free oracle before broad candidate generation.
The oracle contract binds:

- exact rasterizer source hash, P1 corpus and projection-v1 identities;
- generic scale-band IDs and maximum effective scales;
- natural viewBox, padding, uniform fit and derived grid dimensions;
- pixel-center sampling, SVG nonzero fill, no antialiasing, stroke, color or
  softening;
- filled-pixel intersection-over-union and reference-to-candidate recall;
- independent protected-anchor and dominant-component coverage;
- one global IoU/recall threshold pair per generic scale band;
- deterministic source polygon/ring order, candidate order and digest tie-break.

Reference and candidate use the same projector, natural bounds and fit. The grid
long side is the maximum effective scale of the band; its other dimension
follows the natural viewBox. Thresholds are calibrated only against the
representative sheet and frozen into the versioned recipe after owner approval.
Entity- and preset-specific thresholds are forbidden.

Candidate search moves only fine-to-coarse and may omit only unprotected
sub-scale components through one versioned visibility policy. It cannot branch
on country/preset, invent or hand-move coordinates, lose protected land, relax
topology, change q/serializer, or accept a byte overflow.

## Tasks and completion conditions

### T0A — diagnostics, oracle and representative ladder brake; no commit

T0A has two serial phases and creates no product commit.

**T0A.1 — reproduce the non-authoritative diagnostics.** Re-own the RUN-009
projection/selection core and maintainer builder by exact hash. Produce one
canonical JSON diagnostic with:

- P1 corpus, baseline, scratch inventory, projection, recipe, Mapshaper,
  lockfile, tool, Go and command identities;
- exact command lines and the following frozen case matrix.

The parameter scans all use P1 entity `AQ`, boundary profile `un`, the resolved
`card` preset, projection v1, the temporary projected v1 table, compact maximum
effective scale `240`, standard maximum `700`, `q=0.01`, and the frozen
serializer. They call selection with the real compact candidate and current
standard/source fallbacks, while path budgets are disabled because the scan
isolates representation:

| Scan | Fixed axes | Ordered variable axis | Required row result |
| --- | --- | --- | --- |
| resolution | weighted Visvalingam `0.7`, planar, keep-shapes, intersection rollback, clean | `64,65,66,68,72,80,96,128,160,192,224,256,320,512,1024` | 64–96 compact raw `1.537690736882369`; 128/160 typed topology failure; 192/224 raw `1.5403736770216907`; 256–1024 raw `2.441111469697744` |
| weighting | resolution `64`, weighted Visvalingam, planar, keep-shapes, intersection rollback, clean | `0,0.1,0.3,0.5,0.7,0.9,1` | raw values `1.559170376003231,1.559170376003231,1.772910882577505,1.5927713853074141,1.537690736882369,1.6726024738845322,1.7575971919760147` |
| coarse phase | resolution `64`, weighting `0.7`, exact restored compact candidate; diagnostic helper bypasses only the production raw-deviation short-circuit | lexicographic `(x,y)` with each axis `0.000..0.009` step `0.001` | exactly 100 rows, no accepted final result; minimum finite final deviation `1.5360371611532861` |
| local phase A | same candidate and helper | lexicographic x `0.0008..0.0012`, y `0.0018..0.0022`, step `0.0001` | exactly 25 rows, no accepted final result |
| local phase B | same candidate and helper | lexicographic x `0.00105..0.00115`, y `0.00215..0.00225`, step `0.00001` | exactly 121 rows, no accepted final result; minimum remains in `[1.536,1.536001]` |

The phase helper performs the same canonicalization, topology, protection and
final-deviation functions as production, but cannot be called by production
selection and cannot change the exact runtime 100-phase schedule.

The catalog scan uses manifest entity order, then boundary profiles
`un,de_facto`, then presets `card,hero`: exactly 996 rows. For each row,
`InputFromCatalog` resolves the preset once; the diagnostic then clears the
preset token and sets `MaxPathBytes=0` without changing any resolved layout or
quality field. It generates against one temporary v1 table rebuilt at
resolution `64/256`, weighting `0.7`, and records effective scale,
requested/selected candidate, topology, protection, omissions, raw/final
diagnostic deviation, points, path bytes, comparisons to both the historical
`2500/8000` maxima and new `2200/7500` P2 caps, fallbacks and deterministic
output digest.

The historical aggregate reproduction must equal:

- selected counts:
  `card/compact=149`, `card/standard=237`, `card/source=112`,
  `hero/standard=380`, `hero/source=118`;
- historical-max overruns:
  `card/compact=56`, `card/standard=216`, `card/source=76`,
  `hero/standard=222`, `hero/source=55`;
- maximum card path `477784` at `RU/un/source`;
- maximum hero path `675606` at `CA/un/source`.

Canonical semantic JSON excludes wall-clock durations and contains rows in the
orders above. Run it twice from identical inputs; bytes and SHA-256 must match.
Each run writes a separate timing receipt keyed to the semantic digest; timing
is reported but is not required to be byte-identical.

Stop before T0A.2 and withdraw or amend DEC-006 when any machine predicate is
true: an input/hash/case count differs; an expected numeric row differs by more
than `1e-12`; an expected typed failure identity differs; a phase scan accepts
at or below `1.536`; either semantic output/digest drifts; selected counts,
historical-overrun counts or maxima differ; or any parameter-only candidate
both satisfies current topology/protection and makes every 996 row fit the new
P2 cap. Do not reinterpret a contradiction as calibration freedom.

**T0A.2 — freeze oracle v1 and spike the ladder.** Only after T0A.1 passes:

1. implement and mutation-test the pure-Go rasterizer and contract parser;
2. build deterministic fine-to-coarse presimplification candidates for RU, CA,
   CN, AQ, AE, ID, CL, AR, AU and protected island/microstate fixtures;
3. calibrate one global IoU/recall threshold pair per generic scale band;
4. apply one global visibility policy with exact omission provenance;
5. exercise multiple fitted scales across and immediately around every band
   boundary, including natural and portrait/landscape/square contain frames;
6. render the representative contact sheet from the exact machine-accepted
   candidates;
7. obtain owner approval of both threshold pairs and the sheet;
8. freeze oracle, recipe, visibility policy and sheet digest before T0B.

Every representative output must pass topology, protection, oracle and
determinism, stay within `2200/7500` path bytes and have a deterministic
complete-file estimate below `2500/8000`. Runtime must select by effective scale,
quality and byte policy only. Repeating preset/entity names with equal resolved
inputs must select equal candidate identities.

T0A stops with no product commit if no single global threshold pair per band
satisfies machine and owner evidence, a protected microstate cannot fit, a
candidate is visibly unrecognizable, complete-envelope reserve is insufficient,
selection inspects preset/entity/frame shape, rebuilds drift, runtime needs an
external process, or serializer changes are requested.

### T0B — full-catalog integration and owner gate

Start only from passing T0A.1/T0A.2 evidence and explicit owner approval.
Integrate the remaining geometry foundation without modifying T0A-owned oracle,
recipe, search or selection files. A required correction returns to T0A and
reruns both dependent receipts.

Generate all 996 ordinary outputs through the generic engine. Prove:

- all 249 entities, both profiles and both preset inputs, with no missing or
  duplicate output;
- natural layout and representative tiny, huge, portrait, landscape and square
  contain frames with uniform scale and fitted-scale selection;
- deterministic candidate identity, topology, winding, protected coverage,
  removals, oracle metrics, points, path bytes and repeated output;
- path-cap distributions and maxima `≤2200/7500`;
- deterministic complete-file estimates `<2500/8000`;
- normal warm 996-output production `≤120s` on recorded YMBPM3;
- ordinary Go generation performs no Node, Mapshaper, network or external-file
  access;
- owner approval of the full-catalog sheet and any versioned data override.

No executable country/preset branch or oracle bypass is allowed. T0B completes
only after all machine gates and owner approval. Then stage the exact T0A/T0B
union and create one atomic T0 commit. In a clean checkout without T1 generated
artifacts, prove `GOWORK=off go test ./... -count=1`, public nil/source behavior,
typed source-budget failure and frozen serializer/dependency hashes.

### T1 — publish the generic ladder

From the clean T0 commit, run two deterministic maintainer rebuilds and publish
the all-283 ladder artifact, manifest and generated table initializer. Bind P1,
projection, recipe, oracle, visibility, tool, candidate order, omissions and
per-band thresholds. T1 may not modify T0 files.

In a clean T1 checkout, public generation must use the embedded ladder, select
only from effective scale/quality/byte policy, meet `2200/7500`, preserve
natural/contain semantics and perform no external I/O.

### T2 — allocate complete P2 evidence

Add D3 projection/layout parity, oracle and candidate-order mutations,
topology/protection teeth, override, marker, budget, natural/contain and visual
evidence. Run focused gates, `GOWORK=off go test ./... -count=1`, offline
`make check` and a fresh independent audit.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T0A | `silhouette-oracle-foundation` | `internal/geometry/cmd/lodbuild/**`, `internal/geometry/{diagnostic.go,gridphase.go,gridphase_test.go,lod.go,lod_test.go,model.go,model_test.go,preset.go,preset_test.go,projection.go,projection_test.go,silhouette.go,silhouette_test.go}`, `internal/geometry/presets/**`, `internal/geometry/lod/{v1.recipe.json,v2.recipe.json,silhouette-oracle.v1.json}`, `internal/geometry/lod/tool/build.mjs`, `internal/geometry/testdata/silhouette-oracle/**` |
| T0B | `geometry-foundation-integrator` | `go.mod`, `go.sum`, remaining foundation, public seam, full-catalog brake, owner fixtures, overrides, unchanged tool lock/package metadata and frozen serializer staging |
| T1 | `silhouette-ladder-publisher` | `internal/geometry/lod/v2.ladder.json`, `internal/geometry/lod/v2.manifest.json`, `internal/geometry/lod_generated.go`, `internal/geometry/lod_build_test.go` |
| T2 | `geometry-verifier` | `Makefile`, approval/integration/mutation/oracle tests and `internal/geometry/testdata/{approval,d3,mutations}/**` |

T0A has no commit authority. T0B may stage proven T0A bytes but cannot rewrite
them. Serializer files are integration-owned and edit-frozen at their hashes
above. Current `integration_test.go` and `mutations_test.go` remain T2-owned.

No implementer writes `data/**`, `internal/catalog/**`, `cmd/**`, lifecycle
documents, plans, trackers, readiness or another task's lease. Node modules and
temporary diagnostic/artifact output are never staged.

## Dependencies and execution waves

Critical path is `T0A.1 → T0A.2 → owner approval → T0B → T1 → T2`, concurrency
one.

T0A.2 cannot start on assumed coordinator observations; it requires the exact
content-addressed reproduction. T0B cannot start on machine oracle success
alone; it requires frozen global thresholds and representative owner approval.
T1 requires complete T0B evidence, full-catalog owner approval, one T0 commit
and its clean-checkout receipt. T2 requires reproducible T1 publication.

## Validation plan

- `geometry-diagnostic-reproduction`: content-addressed inputs, exact parameter,
  phase and 996-case rows, canonical ordering, two identical outputs and
  aggregate digest.
- `geometry-silhouette-oracle-v1`: rasterizer source identity, common
  projection/fit, natural grid, pixel-center nonzero fill, IoU/recall,
  protected coverage, global thresholds and mutation teeth.
- `geometry-scale-ladder-spike`: fine-to-coarse order, stable tie-break,
  visibility provenance, scale-boundary behavior, `2200/7500`, envelope
  estimates, determinism and representative contact sheet.
- `geometry-selection-policy`: equal effective scale/quality/byte policy yields
  equal candidate identity independent of preset, entity and contain-frame
  shape.
- `geometry-source-contract`: explicit source/custom quality retains DEC-005
  q-aware deviation and typed budget failure.
- `geometry-ladder-integrated`: all 996 ordinary outputs, layouts, metrics,
  removals, budget distributions, `≤120s` warm production and full-catalog
  owner sheet.
- `geometry-t0-clean-checkout`, `geometry-ladder-public-wiring`,
  `geometry-ladder-rebuild-check`, `geometry-d3-parity`, `geometry-overrides`,
  `geometry-markers` and offline project ceiling.

Teeth fail on stale P1/projection/recipe/oracle/tool identity, changed rasterizer
behavior, antialiasing/stroke/color/softening in the oracle, threshold drift,
nondeterministic candidate order, hidden omission, lost protected land, invalid
ring, oracle bypass, byte overflow, preset/entity branch, frame-shape quality,
external runtime process, serializer drift or complete-envelope overflow.

## Deviation and amendment policy

`p2-budget-first-silhouette-ladder-v1` permits helper/test splits within an
existing lease, bounded tuning of the recorded fine-to-coarse presimplification
candidate set during T0A.2, and threshold calibration before owner freeze.

A contradiction in T0A.1 requires DEC-006 withdrawal or amendment before product
work. After freeze, any threshold, rasterizer, metric, fill rule, visibility
policy, band boundary, candidate order, `2200/7500` cap, projection, q,
serializer, runtime dependency or selection-input change requires replan and a
new oracle/recipe version. Architecture maxima `2500/8000` can only tighten
without explicit architecture amendment.

Custom vertex reinsertion, country/preset branches, runtime repair, phase
refinement, higher path/file budgets, lower precision, post-hoc SVG optimization,
MakeValid, union/buffer, CGO and hand-moved coordinates are forbidden.

## Commit worktree and integration policy

T0A creates no commit. After T0B and owner approval, stage only the exact
T0A/T0B paths and create one conventional T0 product commit. T1 and T2 each
create one later atomic commit after complete gates.

Before integration verify expected parent, actual path scope, scratch and
dependency hashes, oracle/recipe/tool identities, frozen serializer and gate
receipts. Lifecycle/evidence records never enter product commits. Recheck T0 and
T1 from isolated clean checkouts. Preserve unrelated/user scratch and never use
blanket staging.

## Rollback and recovery

Before T0 completion, rollback removes only experiment-owned additions; preserved
RUN-009 scratch is not broadly deleted. Resume from the exact diagnostic digest,
oracle/recipe identities, owner receipt and candidate matrix rather than chat
history.

A T0A.1 contradiction withdraws/amends DEC-006. A T0A.2 failure returns to
decision review without broad generation. A T0B failure returns to the frozen
owner/oracle checkpoint and does not relax thresholds or budgets.

T1 publishes through a temporary directory and only after two identical
rebuilds. Revert completed commits in T2 → T1 → T0 order. Reverting T1 leaves
the T0 nil/source seam compiling with honest typed budget behavior.

## Completion and handoff

P2 closes only with reproduced diagnostics, frozen `silhouette-oracle/v1`,
owner-approved representative and full-catalog sheets, deterministic all-283
ladder artifacts, 996 ordinary outputs, exact omissions/protection, natural and
contain layouts, `2200/7500` paths, complete-file evidence beneath `2500/8000`,
offline `≤120s` production, clean T0/T1 checkouts and independent audit.

CTR-002 hands P3 effective scale, selected generic band, recipe/oracle versions,
removals, similarity evidence, path bytes, transform and optional markers. P3
must measure actual complete SVG envelopes per delivery mode and pin state; it
may tighten P2 caps but cannot silently expand architecture maxima.

Estimated effort is 10–22 agent-hours, likely 15, low confidence. T0A.1 is a
1–2-hour contradiction brake; T0A.2 is a 3–6-hour machine-and-owner brake before
the full-catalog cost.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
