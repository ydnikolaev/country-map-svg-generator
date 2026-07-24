---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-009"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-009"
status: running
profiles: []
concerns: []
inputs: ["PLAN-008", "CTX-RUN-011"]
---
# RUN-009 — Projection-aligned selection micro brake

## Authority and accepted inputs

- Accepted PLAN-008 revision:
  `29084a1805bf957e884558d921581d2601a83e8423e77426d1449004d18bd13a`.
- Accepted PLAN-008 manifest:
  `3a69993c33f56a965902484e032a526f8c8e3c379a056a866b4a314ccc97989e`.
- Runtime context CTX-RUN-011:
  `1a93fce428febf6a8ff7d00ba37ac647bcade3a2d1bc2653b3bba83a87db9b33`.
- Baseline is PLAN-008's
  `d219c9bae0434df888d6c4763eb0b19c2e260b91`, tree
  `9096ed613412015f3cc86d5f0615ef753b5a399c`; later lifecycle commits
  contain no accepted geometry product bytes.
- RUN-007/008 results and PLAN-007 invalidation bind through PLAN-008.

## Objective and completion boundary

Execute PLAN-008 T0A only. Align runtime/full projection subdivision with the
versioned artifact contract, separate geometry selection from rendering/budget
enforcement, and run the exact RU/AQ production-table micro brake.

Stop at the T0A checkpoint even if it passes. Do not run T0B, create a product
commit, publish artifacts or touch T1/T2.

## Scope ownership and footprint

The implementer may edit only:

- `internal/geometry/cmd/lodbuild/main.go`,
  `internal/geometry/cmd/lodbuild/main_test.go`;
- `internal/geometry/diagnostic.go`;
- `internal/geometry/gridphase.go`,
  `internal/geometry/gridphase_test.go`;
- `internal/geometry/lod.go`, `internal/geometry/lod_test.go`;
- `internal/geometry/model.go`, `internal/geometry/model_test.go`;
- `internal/geometry/preset.go`, `internal/geometry/preset_test.go`,
  `internal/geometry/presets/v1.json`;
- `internal/geometry/projection.go`,
  `internal/geometry/projection_test.go`;
- `internal/geometry/lod/v1.recipe.json`.

All other paths are read-only to the implementer. In particular `go.mod`,
`go.sum`, `internal/geometry/lod/tool/**`, public pipeline files, serializer,
generated artifacts, integration/mutation tests, docs, `.git` and `.mate` are
forbidden. Existing scratch is shared; preserve unrelated bytes and do not
revert another owner's work. Integration target is `main`. Return without
commit or lifecycle mutation.

## Context and doctrine pack

CTX-RUN-011 supplies current implementation doctrine. PLAN-008 supplies the
exact HOW. Runtime remains pure Go and offline; pinned Mapshaper `0.7.44` is
maintainer-only and may be restored solely from the locked tool dependency for
the opt-in temporary rebuild.

## Tasks

1. Define one versioned centered-LAEA projection contract used by runtime,
   builder, recipe, binder and deviation reference: zero build rotation,
   y-down, flatness `0.10`, precision `0.000001`. Automatic/preset execution
   must not use the implicit pre-AutoQuality `0.20`. Expose actual policy in
   deterministic provenance.
2. Split geometry-only selection from rendering. Selection owns tier binding,
   restoration, topology/protection, source tolerance, 100-phase q-grid search
   and final deviation. Selection takes no serializer or byte-budget input.
3. Render the selected immutable representation once. A hard linear byte
   overrun returns the dedicated typed budget error and never re-enters tier,
   tolerance or phase search. `SelectLODProvenance` performs no rendering.
4. Preserve exact phase order and one-time Transform/ViewBox/marker/protected
   composition. Remove the invalid budgeted-path-equals-unlimited assertion.
5. Rebuild temporary compact/standard tables twice and run only the T0A brake.

## Validation and evidence

Record full recipe, lock, tool and both rebuild hashes. Run focused projection,
selection, gridphase, topology, public-contract compile and lodbuild tests.

The micro matrix is:

- forced-source RU/un/card and AQ/un/card+hero, unlimited;
- AQ/card requested compact and AQ/hero requested standard with the rebuilt
  table;
- representative RU production-table card/hero;
- forced-source bounded budget failure;
- deterministic repeats of selection and rendering.

Unlimited cases prove q-grid, polygon topology, protection, component/ring
identity, raw/final deviation, parser round-trip and provenance. Report RU's
phase/attempt count/runtime; `(0.003,0.002)`, 33 and `<4s` is the regression
anchor, and any change requires explicit evidence.

Production table cases must select compact/standard rather than source and be
within `2500/8000`. Forced source budget failure must be the typed budget error
within `10s`, after exactly one valid representation, with no extra
phase/tolerance attempt. Each case is `≤10s`, post-build micro total `≤30s`;
builder time is separate.

## Constraints and forbidden actions

Keep q=`0.01`, ratio `≤0.012`, floor `0.80px`, explicit cap `1.25px`,
resolutions `64/256`, thresholds `240/700`, budgets `2500/8000`,
natural/contain semantics and frozen serializer hashes
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

No country branch, adaptive phase set, tolerance/q/budget relaxation, post-hoc
optimizer, MakeValid, union/buffer, dependency change, CGO, internal JTS import,
fork, new runtime, generated publication or broad matrix is permitted.

## Stop escalation and resume conditions

Stop without product commit if projection authorities differ, compact/standard
still fails representative deviation or budget acceptance, source budget
failure re-enters search, topology/protection fails, frozen bytes change or a
performance ceiling is exceeded.

If T0A passes, return exact evidence and wait for a separate T0B dispatch. If it
fails, preserve the exact projection/tier/deviation/byte matrix and return to
DEC-004's presimplification metadata or additional precomputed-tier decision.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
