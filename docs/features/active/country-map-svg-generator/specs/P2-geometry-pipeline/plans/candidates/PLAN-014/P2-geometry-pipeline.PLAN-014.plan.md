---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "implementation-plan"
id: "PLAN-014"
epic: "country-map-svg-generator"
spec: "P2"
status: draft
profiles: []
concerns: []
inputs: ["P2"]
---
# PLAN-014 — As-built reconciliation of the geometry pipeline

## Accepted inputs and baseline

Reconcile from commit `83600a6049889f481ba59c5d2a02c47747eeec85`, tree
`e01a9a9f2cf0616ea54e173efb0d2162f31e1e13`.

Bind:

- P2 `3a9be62003fb6daa63130a600d5a01511074bb25d57011ec1ea9f568c1a027e0`;
- architecture
  `5325da7358596b56c9dc6f65d3b6af3801174257503a332334a77a9c71433947`;
- PLAN-013 invalidation receipt
  `d4d5f48dc6dac07ed904297874f5de161040a34c5c43e5b08141c9e06f073138`,
  authority `run-result=RUN-020-RESULT`.

This plan is not product work. P2's product work is already built and committed:
the ladder is built and embedded (T1, T2), served by `Generate()` (T3), guarded
with mutation teeth (T4), and the full catalog is rendered and owner-reviewed
(T5). It shipped as ordinary engineering commits because AM-003 and AM-004 were
believed to fence P2 permanently. That reading was wrong, both amendments are
`retracted`, and the fence is lifted — so what is owed is the obligation-to-
evidence mapping that was never written down, together with an honest register of
what the as-built does not satisfy. Debt `#15`.

PLAN-013 described a narrow byte-cap tier-selection fallback from a baseline that
no longer exists and never described what was built; it is invalidated. Part of
its T0A intent did land in the tree — `internal/geometry/lod_test.go` carries
`TestLODHardBudgetAdvancesToFittingTier`, the replacement PLAN-013 named — so the
mapping below reports the tree as it is rather than restating PLAN-013's claims.

Governing decisions: DEC-003 … DEC-013 accepted; AM-001, AM-002 and AM-005
`verified`; AM-003 and AM-004 `retracted`. AM-005 rewrote VAL-6 into three
separate obligations and its item 5 forbids threshold, tolerance, weighting,
byte-cap and budget changes without a further amendment.

**The project gate baseline is unverified at this commit.** `P2-CHECKPOINT.md`
records `make check` green at 215 passed / 0 failed at `fd1de11`; HANDOFF §7
records that the later re-run was never observed to completion and instructs the
next session to treat the baseline as unknown. Only Markdown and governance
artifacts have been committed since. Establishing the gate is this plan's first
task, not an assumption it rests on.

## Requirement and acceptance coverage

| Obligation | Owning task | Terminal evidence |
| --- | --- | --- |
| US-1, US-3 | T2 | `VR-US-1`, `VR-US-3` |
| US-2, US-4 | T3 | `VR-US-2`, `VR-US-4` |
| REQ-3, REQ-5, REQ-6, REQ-7, INV-1 | T2 | `VR-REQ-3`, `VR-REQ-5`, `VR-REQ-6`, `VR-REQ-7`, `VR-INV-1` |
| REQ-1, REQ-2, REQ-4, REQ-8 | T3 | `VR-REQ-1`, `VR-REQ-2`, `VR-REQ-4`, `VR-REQ-8` |
| VAL-4, VAL-6 | T2 | `VR-VAL-4`, `VR-VAL-6` |
| VAL-1, VAL-2, VAL-3, VAL-5 | T3 | `VR-VAL-1`, `VR-VAL-2`, `VR-VAL-3`, `VR-VAL-5` |
| AC-1, AC-3, AC-4 | T2 | `VR-AC-1`, `VR-AC-3`, `VR-AC-4` |
| AC-2 | T3 | `VR-AC-2` |

Each normative obligation has exactly one manifest owner. The split is by verdict,
not by subject: T2 owns every obligation the committed tree satisfies, T3 owns
every obligation it does not, so an obligation cannot be recorded as satisfied by
the same artifact that records the gap.

### As-built evidence, satisfied (T2)

| Obligation | Evidence in the committed tree |
| --- | --- |
| US-1 | `readiness/T5-catalog-render.contact-sheet.json` (all 993 drawn silhouettes, owner-reviewed); DEC-013 renewed DEC-008's approval on the full catalog; `internal/geometry/shipped_test.go:TestShippedCatalogServesTheLadder` asserts the silhouette fills a sane share of its frame over every committed row |
| US-3 | `internal/geometry/override_test.go` (all four tests); `internal/geometry/overrides/v1.json`, `overrides/group-anchors.v1.json`; `internal/geometry/pipeline_test.go:TestPublicGeneratePreservesBoundaryValidations` |
| REQ-3 | `internal/geometry/mutations_test.go` (all three); `internal/geometry/simplify_test.go:TestTopologySweepMatchesBrute`, `TestTopologyAQSourceCanonicalProbe`; `internal/geometry/pipeline_test.go:TestProtectedFeatureCannotDisappear`; `internal/geometry/silhouette_test.go:TestProtectedVisibilityRejectsVisibleLossAndAllowsProvenSubscale`, `TestSilhouetteOracleRejectsDominantAndProtectedLoss`; `internal/geometry/ladder_gate_test.go:TestLadderGateRecomputesEveryCommittedRow` recomputes topology and protection for all 996 rows |
| REQ-5 | `internal/geometry/retain.go` and `Result.Removals`; `internal/geometry/silhouette_test.go:TestProtectedVisibilityRejectsVisibleLossAndAllowsProvenSubscale` proves per-component `OmissionReason` provenance; `internal/geometry/ladder_gate_test.go:TestLadderGateRejectsASuppressedOmission`; the removal report is consumed downstream at `internal/render/layout_test.go` |
| REQ-6 | `internal/geometry/override_test.go:TestOverrideFailsClosed`, `TestGroupAnchorIsAdditiveAndCannotEvictBase`, `TestGroupAnchorRejectsStaleReviewAndMissingLineage`, `TestGroupAnchorSchemaMutationTeeth` — the last refuses `threshold`, `path_budget`, `candidate_order`, `entity_branch` and `mode` in override data, which is what keeps an override bounded data rather than a code fork. Source corpus bytes are untouched: `internal/geometry` imports only `internal/catalog` and writes nothing |
| REQ-7 | `internal/geometry/marker_test.go:TestMarkerUsesExactGeometryTransform` (markers pass through the identical projection and fit transform); `internal/geometry/pipeline_test.go:TestMarkerAnomalyIsVisible` (typed anomaly plus a diagnostic for an implausibly-outside marker), `TestPipelineDeterministicPresentationFree` (an inside marker carries no anomaly) |
| INV-1 | `internal/geometry/pipeline_test.go:TestPipelineDeterministicPresentationFree` rejects `fill`, `stroke`, `class` and `<svg` in emitted path data; the whole `internal/geometry` package imports exactly one project package, `internal/catalog`, so BND-003 holds structurally; `internal/geometry/cmd/svgshowcase` renders unmodified `Generate()` output through five CSS treatments and eight palettes |
| VAL-4 | `internal/geometry/override_test.go:TestOverrideFailsClosed` (out-of-range rotation, negative minimum parts, wrong corpus, missing metadata all fail); `TestGroupAnchorIsAdditiveAndCannotEvictBase` (a valid override changes only the named dimension and cannot evict a base member); `TestGroupAnchorRejectsStaleReviewAndMissingLineage`; `TestGroupAnchorSchemaMutationTeeth` |
| VAL-6 | Obligation 1 (oracle-gated derived acceptance): `internal/geometry/silhouette_test.go:TestSilhouetteOracleMutationTeeth` (ten mutations including `candidate_order`, `threshold`, `contribution`, `antialias`, `softening`), `TestSilhouetteOracleContractAndRasterizerIdentity`; `internal/geometry/shipped_test.go:TestShippedGateRejectsALadderlessPipeline` rejects both regression shapes 60 of 60 through a provenance discriminator, so reintroducing a raw-deviation gate or restoring source components on the derived path reddens; `TestShippedCatalogServesTheLadder` records zero source fallbacks over the catalog. Obligation 2 (committed-artifact re-verification): `internal/geometry/ladder_gate_test.go`, the one obligation carrying an explicit marker — all four tests, offline and pure Go, recompute every committed row and redden on a perturbed coordinate, a moved threshold and a suppressed omission, with stale corpus/oracle/recipe identity refused upstream by `internal/geometry/ladder_test.go:TestParseLadderTableRejectsStaleBinding` and `TestEmbeddedLadderTableBindsFrozenInputs`. Obligation 3 (explicit-source integrity): `internal/geometry/lod_test.go:TestSourceFinalizationReservesQuantizationAndGuardsFinalDeviation`, `TestProjectionContractPreservesExplicitSourceFlatness`, `TestSelectionRenderSeparationAndTypedSourceBudget`; `internal/geometry/pipeline_test.go:TestPublicGenerateFallsBackToSourceWhenTheLadderHasNoRow` |
| AC-1 | `readiness/T5-catalog-render.contact-sheet.json`; `internal/geometry/testdata/silhouette-oracle/representative.contact-sheet.svg` and `representative.machine-receipt.json`; DEC-012 and DEC-013 record the two artifacts found by looking and their fixes; `internal/geometry/shipped_test.go:TestShippedCatalogServesTheLadder` is the structural receipt |
| AC-3 | `internal/geometry/override_test.go` (all four); `internal/geometry/overrides/` |
| AC-4 | `internal/geometry/mutations_test.go` (all three); `internal/geometry/pipeline_test.go:TestProtectedFeatureCannotDisappear`; `internal/geometry/silhouette_test.go:TestSilhouetteOracleMutationTeeth`, `TestProtectedVisibilityRejectsVisibleLossAndAllowsProvenSubscale`, `TestSilhouetteOracleRejectsDominantAndProtectedLoss`; `internal/geometry/ladder_gate_test.go` (three teeth); `internal/geometry/shipped_test.go:TestShippedGateRejectsALadderlessPipeline`; `internal/geometry/serialize_test.go:TestSerializeMutationZeroByteFalseSuccessFails`; `internal/geometry/cmd/lodbuild/diagnostic_test.go:TestDiagnosticSourceInventoryIsExactAndBiting`, `TestDiagnosticIdentityDriftChecksEveryRecordedLoadBearingIdentity` |

Two of these carry a declared bound rather than an unqualified claim, and T2
records the bound with the evidence. A declared bound is not an open finding and
does not move the obligation to T3: it names work that cannot be done at all under
a constraint the spec itself accepts — an offline gate cannot re-derive geometry
the artifact does not store, and an override mechanism cannot demonstrate a named
exceptional country while none is required — whereas a finding names work that is
owed and unstarted. The distinction is why the split below reads cleanly as
satisfied against open: both bounded cells are fully claimable at the strength
stated, and the bound travels with the claim so no reader has to infer it.

- **VAL-6 obligation 2.** The offline gate does not recompute the *rejected* rungs
  of a `no_artifact` row: those geometries are not stored in the artifact, so
  re-deriving them needs the Mapshaper sweep. That stays the determinism gate's
  job (`internal/geometry/cmd/lodbuild/ladder_test.go:TestLadderRebuildMatchesCommittedArtifact`),
  which is not offline. The boundary is already stated in `P2-CHECKPOINT.md` §T2
  and is carried forward rather than quietly dropped.
- **AC-3.** The override mechanism and its fail-closed teeth are proven on
  fixtures. No shipped entity currently exercises it: `overrides/v1.json` is `[]`
  and `TestGroupAnchorSchemaMutationTeeth` asserts the embedded group-anchor set is
  empty. That is consistent with P2's own scope — exceptional treatment is bounded
  data, and none is presently required — but the criterion reads "named
  exceptional countries", so the absence of a named instance is stated, not
  implied.

### As-built evidence, declared open (T3)

| Obligation | Nearest as-built evidence | What is not satisfied |
| --- | --- | --- |
| US-2 | `internal/geometry/model_test.go:TestTightNaturalAspectAndContainNeverStretch`, `TestLayoutValidationArbitraryFrames`; `internal/render/layout_test.go:TestTightDerivesNaturalProportions`, `TestContainHonoursTheFrameExactly` | "can request any fixed containing frame" is served only at the profile's own long side; any other size takes the DEC-005 source path and exceeds the byte ceiling — OF-2 |
| US-4 | `internal/geometry/silhouette_test.go`; `internal/geometry/mutations_test.go`; `internal/geometry/simplify_test.go:TestTopologySweepMatchesBrute` | topology-damage detection is satisfied; comparison against the *pinned D3 projection/path oracle* the spec names as ground truth is not — no D3 fixture tree exists anywhere in the repository — OF-4 |
| REQ-1 | `internal/geometry/model_test.go:TestTightNaturalAspectAndContainNeverStretch`; `internal/render/layout_test.go:TestContainNeverDistorts`, `TestContainNeverLetsTheSilhouetteEscapeItsFrame`, `TestContainCentresTheUnusedSpace` | the equal-area projection, uniform fit, natural `tight` aspect, proportion preservation and containment all hold. "Centers it" holds only for entities that draw everything they were fitted for — OF-1 |
| REQ-2 | `internal/geometry/preset_test.go` (all three); `internal/geometry/presets/v1.json`; `internal/geometry/integration_test.go:TestApprovalDigestAndBudgets`; `internal/geometry/shipped_test.go:TestShippedPathHoldsTheBandCapAgainstSoftening` | presets-over-an-engine, scale/padding/tolerance defaults and the frozen byte maxima hold. "Callers may use arbitrary natural or fixed-frame sizes" does not hold within budget — OF-2. The QAB-2 distribution finding is recorded against ARCH-001, not against this requirement — OF-3 |
| REQ-4 | `internal/geometry/soften.go:validateSoftenedGeometry` re-validates topology on the sampled curve and measures matched-boundary deviation against the declared tolerance; `internal/geometry/lod.go` emits typed `softening_topology_fallback` and `softening_tolerance_fallback` diagnostics; `internal/geometry/model_test.go:TestAutoQualityUsesFittedGeometryScale` pins the softening coefficient; `internal/geometry/shipped_test.go:TestShippedPathHoldsTheBandCapAgainstSoftening` | the bound is enforced in production code and exercised transitively over the catalog, but no test drives a tolerance-violating or self-intersecting soften and asserts it is refused. There is no `soften_test.go` — OF-8 |
| REQ-8 | `internal/geometry/pipeline_test.go:TestPipelineDeterministicPresentationFree`; `internal/geometry/serialize_test.go` (six tests); `internal/geometry/ladder_gate_test.go:TestLadderGateRecomputesEveryCommittedRow` reproduces all 993 pass rows byte-exactly; `internal/geometry/cmd/lodbuild/ladder_test.go:TestLadderRebuildMatchesCommittedArtifact` rebuilds the artifact byte-for-byte; `internal/geometry/cmd/lodbuild/main_test.go:TestLODSpikeFullCorpus` asserts per-case path equality on rerun | map-iteration independence and repeated-run byte equivalence hold. Host independence has no evidence: every determinism check has only ever run on one host and one architecture — OF-10 |
| VAL-1 | `internal/geometry/integration_test.go:TestApprovalDigestAndBudgets` covers 18 representative entities across continents, islands, antimeridian (FJ, KI, RU) and polar (AQ, RU) at both profiles and both presets; byte-stable rerun is held by the REQ-8 evidence above | the pinned D3 fixtures the guard names do not exist — OF-4. The approval digest is computed and logged, never asserted against a committed value, so representative output can change without reddening — OF-5. The sweep stops at its first failure — OF-6 |
| VAL-2 | `internal/geometry/mutations_test.go` covers self-intersection (bow tie), winding/containment (escaped hole) and near-collapse (degenerate ring); `internal/geometry/pipeline_test.go:TestProtectedFeatureCannotDisappear` covers the protected island; `internal/geometry/ladder_gate_test.go:TestLadderGateRejectsAPerturbedCoordinate` | all four named mutation classes have a tooth and each invalid mutation fails. The integration leg is `TestFullCorpusBothProfiles`, which stops at its first failure and therefore does not satisfy the population-guard clause — OF-6 |
| VAL-3 | `internal/geometry/integration_test.go:TestRepresentativeNaturalRatiosAndArbitraryFrames`; `internal/geometry/model_test.go:TestAutoQualityUsesFittedGeometryScale`; `internal/geometry/preset_test.go:TestAutoQualityVersionedRelativePolicy`; `internal/render/layout_test.go:TestContainNeverDistorts`; the T5 contact sheet and DEC-013's owner approval | the natural-ratio leg asserts RU, CL and AU. The multiple-`contain`-frames leg runs four frames but asserts nothing beyond the absence of an error — its trailing comparison is an empty block — OF-7. Nothing rasterizes at the intended long-side sizes inside `make check`; "recognizable approved baselines" rests on owner review of the T5 sheet, which is evidence but not a gate. The sweep stops at its first failure — OF-6 |
| VAL-5 | `internal/geometry/marker_test.go:TestMarkerUsesExactGeometryTransform` (same transform); `internal/geometry/pipeline_test.go:TestMarkerAnomalyIsVisible` (implausible-outside, typed anomaly plus diagnostic), `TestPipelineDeterministicPresentationFree` (inside, no anomaly) | the edge case the guard names between inside and implausible-outside has no fixture — OF-9 |
| AC-2 | the REQ-1, REQ-2 and REQ-8 evidence above | the criterion asks for a golden receipt. `TestApprovalDigestAndBudgets` logs a digest instead of asserting one — OF-5. "Every `contain` output" excludes the entities whose drawn silhouette sits off centre — OF-1. "Quality follows effective fitted scale rather than unused frame space" is stated in `AutoQuality` and pinned as a formula, but the across-frames comparison that would prove it is the assertion-free loop — OF-7 |

### Open findings register

| ID | Finding | Bears on | Item | Class |
| --- | --- | --- | --- | --- |
| OF-1 | An explicit `contain` frame is fitted before visibility removals, so the drawn silhouette sits off centre — up to 292 px in a 300 px frame. DEC-013 fixed this for a committed ladder candidate; an explicit frame changes the layout, the committed verdict does not transfer, and the request takes the source path where the fit still precedes removals. `fitGeometry` does centre what it is given; the defect is in the caller | REQ-1, US-2, AC-2 | `WKI-37F18A2AA6A5` | code fix in BND-003, no threshold change |
| OF-2 | Only the profile's own long side is served from the ladder; any other requested size falls back to source and exceeds the byte ceiling | US-2, REQ-2, AC-2 | `WKI-C35A01E965DC` | touches the ladder artifact ⇒ rebuild ⇒ **amendment-class** under AM-005 item 5 |
| OF-3 | Finest-that-fits clusters every asset at the byte cap, well over ARCH-001 QAB-2's median in both bands — from the committed ladder artifact, a card path median of 1962 B against QAB-2's 600 B and a hero median of 4140 B against 1800 B. QAB-2 is an architecture quality attribute, not a P2 obligation id; P2's own text names only byte budgets and the frozen maxima, and per the committed ladder artifact the maxima pass in both bands. Recorded against ARCH-001, noted against REQ-2's budget language | ARCH-001 QAB-2; noted at REQ-2 | `WKI-33B6EC2482B3` | **owner decision** on the threshold, then **amendment-class** |
| OF-4 | No pinned D3 projection/path oracle fixtures exist. The spec's ground truth names D3's azimuthal equal-area projection, fit and path output as the pinned reference and VAL-1 requires comparison against them; there is no `testdata/d3` tree in the repository. What exists is the raster *silhouette* oracle, which is a different instrument answering a different question | US-4, VAL-1 | none — new | fixture plus parity test |
| OF-5 | The representative approval digest is computed and logged, not asserted against a committed value, so a change in representative output does not redden | VAL-1, AC-2 | none — new | test-only |
| OF-6 | The population-guard clause is unmet. P2's corpus sweeps stop at their first failure, and the spec states that a sweep which does so does not satisfy VAL-1 through VAL-3. `TestFullCorpusBothProfiles` and `TestApprovalDigestAndBudgets` abort; `TestShippedCatalogServesTheLadder` aborts on the first row error. None reports counts or a distribution of failures. The downstream layout tests in `internal/render` accumulate instead, which is the shape P2's own sweeps need. `TestFullCorpusBothProfiles` is additionally card-only | VAL-1, VAL-2, VAL-3 | none — new | test-only |
| OF-7 | VAL-3's multiple-`contain`-frames leg is assertion-free: the loop runs four frames and only checks that no error was returned, and the comparison that would have checked quality across frames is an empty block | VAL-3, AC-2, REQ-1 | none — new | test-only |
| OF-8 | REQ-4's softening bound has no direct test. The production path validates topology on the sampled curve and measures deviation against the declared tolerance, and falls back with a typed diagnostic — but nothing drives a violating soften and asserts refusal | REQ-4 | none — new | test-only; a new file under `internal/geometry` must be registered in `TestDiagnosticSourceInventoryIsExactAndBiting` |
| OF-9 | VAL-5's edge-marker case has no fixture; only inside and implausible-outside are covered | VAL-5 | none — new | test-only |
| OF-10 | REQ-8's host independence has no evidence; every determinism check runs on one host and architecture. The route to closing it is P5 VAL-2's clean-machine proof, which has no available host | REQ-8 | none — new | **owner/infra** |
| OF-11 | Wall-clock brakes force `make check` to run `-p 1`, which is why the project ceiling costs ~14 minutes serially. This bears on the validation apparatus, not on any product obligation | none | `WKI-B05B4B4287A2` | test-only |
| OF-12 | The diagnostic identity record hashes the whole module graph, so any dependency change reddens it. Its recorded next action names this reconciliation as the place to decide | none | `WKI-1DA58E0FE741` | decision, then test-only |

`WKI-490046152C71` (mainland versus full claim) is `deferred` and DEC-011 and
DEC-013 route what remains of it — component *selection* rather than framing — to
P3. It is recorded as routed, not as a P2 open finding.

## Technical approach

Write down what is true and name what is not. The run this plan governs produces
four Markdown artifacts under the spec's own directory and changes no Go source,
no test, no data file and no committed artifact.

Three properties make this a reconciliation rather than a receipt:

1. **Every claim is bound to a file and a test function in the tree at the plan
   baseline.** A claim that names no executable evidence is recorded as an open
   finding, not as coverage. The mapping above was established by reading the
   tests; only VAL-6 obligation 2 carried a marker, and nothing was inferred from
   a test's name.
2. **The project ceiling is measured before the mapping is trusted.** A test
   function is not evidence for an obligation until it is known to pass. The
   baseline is unverified at this commit and establishing it is T1, ahead of
   everything else.
3. **Nothing is claimed at a strength that was not checked.** Where the as-built
   satisfies part of an obligation, the satisfied part is stated and the rest goes
   to the register with the work item or the successor route that owns it. Where a
   gate carries a boundary — VAL-6 obligation 2's unrecomputed rejected rungs,
   AC-3's absent shipped instance — the boundary travels with the claim.

The register is written against the spec's own obligation text, not against a
convenient reading of it. Two mappings were checked rather than assumed and came
out differently from the obvious guess. The off-centre `contain` frame is
characterized in `internal/render/layout_test.go`, which is P3's package, and
attributed there to **REQ-13**, a P3 requirement; the same file records that
`fitGeometry` does centre the geometry it is given. So the honest statement under
P2's REQ-1 is that centring is implemented and holds wherever nothing is removed
after the fit, and that the explicit-`contain` path did not inherit DEC-013's
fit-to-what-is-drawn rule. And the QAB-2 median overshoot is a finding against an
ARCH-001 quality attribute; P2's own text names no median or p95, and the frozen
maxima it does name pass in both bands, so forcing that item onto a P2 requirement
as a violation would be as dishonest as hiding it.

No deviation figure is carried into this plan. The acceptance metric on the
explicit-source path is a clamped predicate — exact at or below its limit and
short-circuiting above — so a first-crossing value from it is not a measurement,
and the spec forbids reporting one as though it were. The predecessor plan quoted
two such figures; they are not repeated here.

## Tasks and completion conditions

### T1 — Establish the project ceiling at the plan baseline

Run the registered project gate, `make check` (`go test ./... -count=1 -p 1`), on
an idle machine at commit `83600a6`, and write the receipt.

Discipline, because this gate has produced a false red and a false green in this
repository before:

- run it detached and read the log to completion; never read the exit code of a
  backgrounded run and never pipe it, because the pipeline reports the last
  command's status;
- check `pgrep -f "go test"` first — the brakes measure the machine, not the code,
  and a contended run has been observed at 440 s against 343 s quiet;
- summarize **package**-level failures as well as test-level ones. A `go test
  -json` package failure carries no `Test` field, so a test-only filter reports a
  timed-out package as a clean run. That has happened twice here;
- expect roughly 757 s serially, with `internal/geometry` and
  `internal/geometry/cmd/lodbuild` the slow, brake-carrying packages and both near
  the default ten-minute per-package timeout.

Done when the receipt records the invocation, the host, the wall-clock total, the
per-package result including any package-level failure, and the pass/fail counts,
and states plainly whether the ceiling is green. **A red or unread gate stops the
run**: a red result is reported as a blocker rather than worked around, because
an unattributable red at this baseline would make every coverage claim below it
unsound.

### T2 — Map every satisfied obligation to its as-built evidence

Start only from a green T1 receipt. For each of US-1, US-3, REQ-3, REQ-5, REQ-6,
REQ-7, INV-1, VAL-4, VAL-6, AC-1, AC-3 and AC-4, write the obligation, the exact
file and test function (or committed artifact) that proves it, what the evidence
actually asserts, and the strength of the claim.

Run each satisfied obligation's attestation gate — `p2-asbuilt-override-attestation`
for VAL-4 and `p2-asbuilt-oracle-gate-attestation` for VAL-6 — and record the
observed function set from its `-v` output against the set the matrix names.

Done when every one of the twelve has at least one named file-and-function whose
presence is confirmed in the tree at the baseline and, where a gate covers it,
observed passing by name; when VAL-6 is recorded as
AM-005's three separate obligations, each with its own evidence; when VAL-6
obligation 2's unrecomputed-rejected-rungs boundary and AC-3's absent shipped
instance are recorded with their claims; and when no cell names a test that T1 did
not observe passing.

### T3 — Register every open finding against the obligation it bears on

Start only from T2. For each of US-2, US-4, REQ-1, REQ-2, REQ-4, REQ-8, VAL-1,
VAL-2, VAL-3, VAL-5 and AC-2, record the nearest as-built evidence, the precise
part of the obligation text it does not reach, the finding OF-1 … OF-12 that names
the gap, the backlog item where one exists, and the class of work that would close
it.

Run the four remaining attestation gates — `p2-asbuilt-golden-parity-attestation`,
`p2-asbuilt-topology-protection-attestation`, `p2-asbuilt-layout-scale-attestation`
and `p2-asbuilt-marker-attestation` — and check each gate's condition 3: every
clause of the VAL's own guard text is matched either to evidence in T2's matrix or
to a finding here.

Done when all eleven are recorded; when every finding names either a backlog item
or is marked explicitly as newly discovered by this reconciliation; when each
finding is classified as test-only, code-fix, amendment-class or owner-gated; and
when no finding is described as closed, mitigated or scheduled by this plan.
**Nothing in this register is fixed by this run.** OF-2 and OF-3 are recorded as
amendment-class and owner-gated respectively and are explicitly out of scope.

### T4 — Write the reconciliation record and update the checkpoint

Start only from T2 and T3. Write the reconciliation record: what P2 built, which
commits carry it, why it shipped outside the governed wrapper, that the fence is
lifted, that PLAN-013 is invalidated and this plan succeeds it, the counts from
T2 and T3, and the terminal state — P2 remains `in_progress` with declared open
work, and this plan does not close it.

Update `P2-CHECKPOINT.md` so the resume record points at the three artifacts
rather than at the debt. Correct nothing else in it: the checkpoint's historical
sections, including the superseded fence reading it deliberately preserves, stay
as they are.

Done when the record exists, when the checkpoint references all three artifacts,
when the count of satisfied and open obligations sums to 23, and when the record
states in its own words that this is a reconciliation of the as-built and not a
claim of completion.

## Ownership and write footprint

| Task | Owner | Exclusive writes |
| --- | --- | --- |
| T1 | `p2-project-ceiling-recorder` | `docs/features/active/country-map-svg-generator/specs/P2-geometry-pipeline/evidence/P2-project-ceiling.gate-receipt.md` |
| T2 | `p2-as-built-coverage-author` | `docs/features/active/country-map-svg-generator/specs/P2-geometry-pipeline/evidence/P2-as-built.coverage-matrix.md` |
| T3 | `p2-open-findings-registrar` | `docs/features/active/country-map-svg-generator/specs/P2-geometry-pipeline/evidence/P2-as-built.open-findings.md` |
| T4 | `p2-reconciliation-recorder` | `docs/features/active/country-map-svg-generator/specs/P2-geometry-pipeline/evidence/P2-as-built.reconciliation-record.md`; `docs/features/active/country-map-svg-generator/specs/P2-geometry-pipeline/P2-CHECKPOINT.md` |

Leases do not overlap. Every write is spec-local, under
`docs/features/active/country-map-svg-generator/specs/P2-geometry-pipeline/`.

**No task writes Go source, a test, a data file, a committed artifact, a
Makefile, or anything under `internal/`, `cmd/`, `data/`, `.git` or `.mate`.** A
diff touching any of those is a lease violation, not a deviation: this run cannot
change the thing it is describing without invalidating the description. If an
implementer concludes a code write is unavoidable, that is a blocker to escalate,
not a task to add.

Lifecycle artifacts written by `mate` itself — the run directory, transition
receipts, the tracker — are the coordinator's, not an implementer's, and are not
in any lease above.

## Dependencies and execution waves

Critical path: `W0/T1 project ceiling → W1/T2 satisfied coverage → W2/T3 open
findings → W3/T4 reconciliation record and checkpoint`. Concurrency is one.

The ordering is load-bearing rather than conventional. T2 cannot honestly claim a
test proves an obligation until T1 has observed that test pass, and T3's register
is the complement of T2's mapping over the same 23 obligations — writing it first
would let a gap be recorded as coverage. T4 depends on both because its counts
must sum.

## Validation plan

The registered project ceiling is `make check`, that is
`go test ./... -count=1 -p 1`, roughly 14 minutes serially. It is declared here as
this plan's project tier and is executed exactly once, by T1, at the plan
baseline. It is **unverified at commit `83600a6`**: `P2-CHECKPOINT.md` records it
green at `fd1de11`, HANDOFF §7 records that the later re-run was never observed to
completion, and only Markdown has been committed since. Establishing it is T1's
completion condition, not this plan's premise.

Six per-obligation attestation gates sit under it. Each is a real argv over the
committed tree, not a claim about it: a `go test -run` selector naming exactly the
test functions that VAL's coverage cell cites, executed read-only against the
package that holds them.

| Gate | Obligation | Tier | Argv shape |
| --- | --- | --- | --- |
| `p2-asbuilt-golden-parity-attestation` | VAL-1 | integration | `go test -count=1 -v -run` over `TestApprovalDigestAndBudgets`, `TestChinaCanonicalQuantization`, `TestSerializeRoundTrip`, `TestLadderGateRecomputesEveryCommittedRow` |
| `p2-asbuilt-topology-protection-attestation` | VAL-2 | integration | over `TestMutationTopologyBowTieFails`, `TestMutationEscapedHoleFails`, `TestMutationCollapsedRingFails`, `TestProtectedFeatureCannotDisappear`, `TestTopologySweepMatchesBrute`, `TestLadderGateRejectsAPerturbedCoordinate` |
| `p2-asbuilt-layout-scale-attestation` | VAL-3 | project | over `TestRepresentativeNaturalRatiosAndArbitraryFrames`, `TestTightNaturalAspectAndContainNeverStretch`, `TestLayoutValidationArbitraryFrames`, `TestAutoQualityUsesFittedGeometryScale`, `TestAutoQualityVersionedRelativePolicy`, plus `internal/render`'s `TestContainNeverDistorts`, `TestContainNeverLetsTheSilhouetteEscapeItsFrame`, `TestContainCentresTheUnusedSpace` |
| `p2-asbuilt-override-attestation` | VAL-4 | component | over the four `internal/geometry/override_test.go` functions |
| `p2-asbuilt-marker-attestation` | VAL-5 | integration | over `TestMarkerUsesExactGeometryTransform`, `TestMarkerAnomalyIsVisible`, `TestPipelineDeterministicPresentationFree` |
| `p2-asbuilt-oracle-gate-attestation` | VAL-6 | project | over the five `silhouette_test.go`, four `ladder_gate_test.go`, three `shipped_test.go` and three explicit-source `lod_test.go` functions named in the coverage cell |

Each gate has three failing conditions, and they are what make it a gate rather
than a restatement:

1. **A named function does not exist or does not pass.** A `-run` selector that
   matches a function which was renamed or deleted is the usual way an as-built
   claim rots; the receipt records the observed function names from `-v` output
   and fails when the set does not equal the set the coverage cell named. A
   selector that silently matches nothing fails here rather than passing green.
2. **A coverage cell names evidence no gate argv covers.** Every function cited in
   T2's matrix must appear in at least one gate's selector. A function may serve
   two guards — the ladder gate's teeth carry both VAL-2's perturbation class and
   VAL-6's obligation 2 — and being named twice is coverage, not duplication; being
   named nowhere is the failure.
3. **A clause of the guard's own text is neither claimed nor registered.** Every
   sentence of the VAL row in the spec is either matched to evidence in T2 or
   matched to a finding in T3. A clause that appears in neither fails the gate.

Condition 3 is the one that bites this particular run. VAL-1's gate does not pass
because the D3 fixtures are absent — it passes because their absence is registered
as OF-4 and their absence is therefore visible in the artifact rather than in
nobody's notes. Had OF-4 gone unwritten, condition 3 fails.

Plan validation is `mate plan validate country-map-svg-generator P2 PLAN-014
--json`. No plan-authoring step runs the project ceiling; only T1 does, and only
on an idle machine.

## Deviation and amendment policy

`p2-as-built-reconciliation-documentation-only-v1` permits exactly one thing:
writing spec-local governance and evidence documents that describe the committed
tree at baseline `83600a6`, plus the checkpoint update that points at them.

Everything else escalates. Specifically forbidden without another plan:

- any edit to Go source or tests, including adding the test that would close
  OF-4 … OF-9;
- any ladder rebuild, any change to the committed ladder artifact, recipe, oracle,
  serializer or corpus;
- any byte-cap, threshold, tolerance, weighting, budget, band, resolution, `q`,
  projection or tool-pin change — AM-005 item 5 makes these amendment-class, and
  the byte threshold is additionally a pending owner decision;
- any change to the spec, its amendments, its decisions or its acceptance
  criteria. This run reconciles the as-built against P2 as written; it does not
  edit P2 to fit the as-built.

An obligation that cannot be honestly claimed without one of the above is recorded
as an open finding and escalated. It is never softened into coverage, and the
obligation text is never reinterpreted to make the as-built fit.

## Commit worktree and integration policy

Each task produces one atomic commit staging only its own leased paths, `docs`
scope, conventional message, no `git add -A`. T4 may stage `P2-CHECKPOINT.md`
alongside its record because the two are one intent.

No product commit is created; the tree's Go source, tests, data and artifacts are
byte-identical before and after. A diff that shows otherwise fails the run.
Lifecycle artifacts remain separate coordinator commits.

## Rollback and recovery

Every write is a new Markdown file except the checkpoint, so rollback is deleting
the files and reverting the commits in T4 → T3 → T2 → T1 order. Nothing built,
embedded, generated or approved is touched, so no rebuild, re-approval or
regeneration is ever needed to recover.

`P2-CHECKPOINT.md` is the one pre-existing file: T4 appends and repoints, and
restores it from `83600a6` on any failure. Its historical and superseded sections
are preserved verbatim — they exist so a documented mistake keeps something to
correct.

If T1's ceiling is red, stop and report the failing package and test with the log.
Do not proceed to T2, do not attribute the failure, and do not fix it here: this
plan's lease does not include the code that would have to change, and a
reconciliation authored over a red gate is worth nothing. If a later task finds an
obligation whose evidence the ceiling did not actually exercise, move that
obligation from T2's mapping to T3's register rather than weakening the claim.

## Completion and handoff

P2's reconciliation closes with three artifacts and a checkpoint that points at
them: a measured project-ceiling receipt at the plan baseline, a mapping of every
satisfied obligation to the file and test function in the committed tree that
proves it, and a register of every obligation the as-built does not satisfy with
the work item or successor route that owns each. Twenty-three obligations, each
appearing exactly once, in exactly one of the two.

**P2 does not complete here.** The terminal state of this plan is a spec that is
`in_progress` with its open work declared instead of undocumented — which is what
debt `#15` actually asked for, and strictly more than a backdated receipt would
have given. Twelve findings are recorded. OF-1 and OF-4 … OF-9 are ordinary
successor work inside BND-003. OF-2 and OF-3 are amendment-class, OF-3 gated on an
owner decision on the byte threshold that no agent may make. OF-10 needs a host
that does not exist yet.

P3 receives the same CTR-002 surface it already consumes — no contract changes —
plus the two findings that reach its territory: the off-centre explicit-`contain`
frame it already characterizes in its own package, and the component-selection
seam DEC-011 and DEC-013 route to it.

Estimated remaining effort: 3–8 agent-hours, likely 5, medium confidence. The
dominant costs are the serial project ceiling and reading roughly twenty test
files closely enough that no cell in either table is asserted from a test's name.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
