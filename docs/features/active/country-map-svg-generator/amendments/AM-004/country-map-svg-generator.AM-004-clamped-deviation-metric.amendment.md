---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "amendment"
id: "AM-004"
epic: "country-map-svg-generator"
status: proposed
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# AM-004 — Amendment

## Trigger and discovery evidence

Independent review `ADV-003` blocked AM-003 with two findings, both verified
against the code before this amendment was written. AM-003 named the right
subsystem and drew a defensible conclusion, but two of the three defects it
listed are wrong, and the guard it produced can be satisfied without touching the
real fault. This amendment corrects the diagnosis and replaces the obligation.

**AM-003's defect #1 was unreachable.** It claimed `matchedBoundaryDeviation`
returns `+Inf` on component-count mismatch at
`internal/geometry/lod.go:622-624`, and attributed 25 infinite rejections across
CL, NO and PH to that branch. The branch exists, but `len(full) != len(candidate)`
cannot occur at any call site. On the tier path the candidate is always the output
of `restoreRequiredComponents`, which appends exactly one polygon per source
polygon (`internal/geometry/lod.go:584-607`), so the lengths are equal by
construction. On the source path the comparison is guarded five lines earlier by
`sameLODStructure` (`internal/geometry/lod.go:436`), which reports a different and
already-named fallback, `source:retention:…:component_or_ring_loss`. Component
loss is therefore already handled by restoration and by the structure check, and
never reaches the metric. The remaining `+Inf` returns at
`internal/geometry/lod.go:640-642` and `667-669` are likewise unreachable once
`sameLODStructure` has proved the sorted multiset of ring counts equal.

The `+Inf` that is reachable is `lodSegmentIndex.distance`
(`internal/geometry/lod.go:728`): when a query point's grid cell contains no
indexed segment, `best` stays `+Inf` and propagates out through
`symmetricRingDeviation`. That is consistent with the affected entities being
archipelagic, but this amendment does not assert it as established — the
attribution is unconfirmed, and settling it is part of the measurement required
below.

**AM-003's magnitudes do not exist as measurements.** `matchedBoundaryDeviation`
is a clamped predicate, not a distance:

- it short-circuits the moment the running maximum crosses the limit —
  `internal/geometry/lod.go:649-651` and `675-677`, and inside
  `symmetricRingDeviation` at `712-714` and `721-723` — so a returned value above
  the limit is the first per-point value to cross it in iteration order, not the
  maximum;
- its spatial index is scoped to the limit — `newLODIndex`
  (`internal/geometry/lod.go:689-703`) registers each segment only in the cells
  its bounding box expanded by `limit` touches, with
  `cell := math.Max(limit, .05)`. That guarantees any segment within `limit` of a
  query point is present, so values **at or below the limit are exact**; it
  guarantees nothing above, where `distance` returns the minimum over whatever
  segments happen to share one cell.

Every figure AM-003 quoted as a measured deviation is above its tolerance and is
therefore an instrument artifact. Specifically withdrawn: "the median excess over
tolerance is 17 times the `q/√2` quantization reserve", "only 10 of 156 deviation
rejections are small enough to be attributable to quantization noise", and the
per-attempt magnitudes quoted for CA/un/hero. The inference those figures
supported — that quantization noise is ruled out and that ladder resolutions are
genuinely too coarse — is withdrawn with them.

**What survives, and why.** The predicate is sound: at or below the limit the
metric is exact, so "this candidate exceeds tolerance" is reliable even when "by
how much" is not. Every conclusion resting only on pass/fail therefore stands:

- 75 of 996 outputs across 27 entities cannot reach a within-budget tier;
- the first rejection at the finest candidate tier is a deviation rejection in 65
  cases, a byte rejection in 8, a topology rejection in 2 — so the byte overrun is
  the last link of the chain, not the first;
- the cheap tiers are rejected, selection falls through to the unsimplified
  `source` tier, and that tier breaks the byte budget;
- the source-tier search does not converge — every attempt in a descending
  sequence continues to fail the predicate.

The corrected root cause is sharper than AM-003's. The search at
`internal/geometry/lod.go:429-445` halves its tolerance through
`nextSourceTolerance` (`internal/geometry/lod.go:543-548`) in order to reduce a
quantity that the instrument refuses to measure whenever it matters. A descent
cannot converge on an objective that reports only "over", plus an arbitrary number
once it is over. That is the defect: not that the metric is wrong, but that a
pass/fail predicate is being used as a minimization objective.

## Affected scope and authority

Same scope and authority as AM-003: a P2 macro-HOW amendment to the deviation
acceptance path, owned by the product owner. P2's WHAT is unchanged — no REQ, INV
or AC row is touched, the specification makes no tier label normative and fixes no
byte budget numerically. DEC-004's precomputed-tier architecture and pure-Go
runtime are not reopened. AM-003's classification was independently reviewed and
confirmed correct on this point; only its factual content and its guard are
replaced.

AM-003 remains `applied` and unverified, with `ADV-003` recording why. This
amendment supersedes its diagnosis and its accepted change in full. The VAL-6 row
AM-003 placed in the P2 specification is replaced by the text below rather than
retained.

## Prior accepted truth

AM-003 at `d910385478f3d9eea6c3f3755aa769928da0f7ce0c82c1fedaed68d4835bdffd`,
applied to the P2 specification at
`ad93558b8bd0370fbc8bec5e3f644871e7c52389d343f7959cc4eb8f3c444938`, which carries
VAL-6 in the form AM-003 accepted:

> exercise the boundary-deviation acceptance metric directly: mismatched
> component and ring counts, permuted component order, a perturbation that would
> flip a proximity-based correspondence, and a descending sweep of the
> simplification tolerance across the resolution ladder …

`ADV-003` finds that row defective in three ways: its first clause tests the
unreachable branch, so a synthetic unit test feeding two multipolygons of unequal
length satisfies it in full while production keeps returning `+Inf`; it binds
monotonicity to the metric, which no metric can deliver, because the candidate
geometries in a descending sweep are not a nested family; and it names the
29-resolution build ladder as the search variable when the source-tier search does
not use that ladder at all.

DEC-005's reserve of `q/√2` at `internal/geometry/lod.go:423` and its stated
intent to "search monotonically from the coarsest allowed simplification toward
finer candidates" are unchanged and remain accurate. AM-003 asserted DEC-005 had
already assumed metric monotonicity; DEC-005 describes the search *sequence*, and
this amendment does not attribute more to it than that.

## Accepted change

**1. Measure before deciding.** Before any change to the metric, the ladder or the
budgets, P2 produces one governed measurement of the true deviations: the
short-circuits at `internal/geometry/lod.go:649-651`, `675-677`, `712-714` and
`721-723` disabled and the index unclamped, so that the metric returns an exact
maximum rather than a first crossing. Passing a sufficiently large limit achieves
both at once, since `cell := math.Max(limit, .05)` collapses the index to a scan;
whether that or an explicit unclamped variant is used is an implementation choice.

That measurement reports, for every one of the 75 failing outputs and every
attempt in its search: the true deviation, the tolerance, the candidate tier and
the rejection reason. It settles three questions no current evidence can answer —
which branch actually produces the infinities, what the real excess distribution
is, and whether deviation is in fact monotone in the simplification tolerance once
measured honestly. No change to the metric, the ladder or the budgets is
authorized before this measurement exists as evidence.

**2. Separate the predicate from the objective.** The acceptance path keeps a
clamped predicate — it is correct, and its early exit is what makes it affordable.
What the search consumes must be a real objective: a deviation the implementation
is willing to compute exactly for the candidate under consideration, so that a
descending sequence of tolerances produces a comparable sequence of values. The
two may be one function with an unclamped mode or two functions; that is a HOW
choice for the plan.

**3. Bind the search property to the pipeline, not to the metric.** The obligation
is that the *accepted candidate* improves as the requested tolerance falls, over
the tolerance sequence the source search actually walks — `remaining` halved
through `nextSourceTolerance` — not over the 29-resolution build ladder, which
this search does not use. The property is asserted on the pre-quantization
comparison, because the grid phase adds up to `q/√2` independently of tolerance
(`internal/geometry/lod.go:501`), and it carries an explicit numeric band rather
than a strict float inequality. Where `simplifySharedResolved` substitutes
`simplifyRings` output or the unsimplified polygon on failure
(`internal/geometry/simplify.go:43-50`), that substitution is reported as a named
outcome rather than silently breaking the sequence.

**4. Totality is about the reachable branch.** A finite graded result is required
where `+Inf` can actually arise — the index cell miss at
`internal/geometry/lod.go:728` — and the guard must exercise it through the
production path for an archipelagic entity, not through a synthetic input the
pipeline cannot construct. Protected-feature loss continues to fail through the
protection check under REQ-3, and component loss continues to be handled by
restoration and `sameLODStructure` rather than by the metric.

Byte budgets `2200/7500` and `2500/8000` are unchanged and no per-country
exception is introduced; country literals and branches remain forbidden. Ladder
recalibration is not authorized by this amendment at all — AM-003 opened it as a
dependent second step on the strength of the withdrawn distribution, and it
returns to being an open question the step-one measurement must answer first.

No estimate is offered for how many of the 75 outputs any of this recovers.

## Downstream impact map

- The P2 specification's VAL-6 row is rewritten as stated in the validation delta
  below, replacing the AM-003 text.
- P2's `inputs` frontmatter carries AM-001, AM-002, AM-003, AM-004, DEC-004 and
  DEC-005. `ADV-003` finding F5 established that it currently carries none of
  them: `inputs: ["DISC-006", "ARCH-001", "DEC-003"]`. AM-002's precedent for this
  is P3's frontmatter, which does carry `AM-001`.
- P3's `inputs` frontmatter carries AM-003 and AM-004. No P3 spec body change is
  required, consistent with AM-003's impact map; the obligation is context, not
  contract.
- The P2 context registry advances so the next plan is authored against this
  amendment and the step-one measurement, not against AM-003's framing.
- No change to the silhouette oracle, the v1 recipe or the ladder is authorized
  here, so neither owner contact sheet is invalidated by this amendment itself.

## Validation delta

VAL-6 is replaced by:

- **Reachable totality.** Drive an archipelagic entity through the production path
  and assert the metric returns a finite graded value where the index cell miss at
  `internal/geometry/lod.go:728` would otherwise yield `+Inf`. A synthetic input
  with mismatched component counts does not satisfy this, because the pipeline
  cannot produce one.
- **Order stability.** Identical inputs give identical values under permuted
  component order, and correspondence does not change under a perturbation that
  would flip the greedy centroid assignment at
  `internal/geometry/lod.go:630-638`. This preserves REQ-8.
- **Search progress.** Over the tolerance sequence the source search actually
  walks, the accepted candidate's pre-quantization deviation does not increase
  beyond a declared numeric band as the requested tolerance falls; a
  `simplifySharedResolved` substitution is reported as a named outcome rather than
  silently breaking the sequence.
- **Honest magnitudes.** Any deviation figure carried in provenance, evidence or a
  plan is either at or below its limit, where the clamped metric is exact, or is
  produced by the unclamped path. A first-crossing value is never reported as a
  measurement.

Population guards under VAL-1 through VAL-3 continue to report every failing
entity with counts and distribution; a sweep that stops at its first failure does
not satisfy them. That clause of AM-003 is retained unchanged — `ADV-003` did not
dispute it, and `t.Fatalf` inside the entity loop at
`internal/geometry/cmd/lodbuild/main_test.go:88` is confirmed.

VAL-6's owner/due column reads `P2 / W2 complete`, matching its five siblings.
AM-003's `W0` was inconsistent with the convention across all five spec files.

## Migration rollout and rollback

There are no runtime instances and no published artifacts. The step-one
measurement changes nothing observable and can be discarded. Any subsequent metric
change alters selection outcomes across the corpus, so every representative and
full-catalog sheet is regenerated and re-approved by the owner before P2 closes.

Rollback restores the P2 specification to its pre-AM-003 bytes
`4fb4535e74ee072d9db6ecc6e1ee92e667ac9c78f35a8adbf342c7b93d1c2741`. The PLAN-013
rung committed at `d052f9b8dacff41719ddbc19431955c3bcbac690` is independent of
both amendments and is retained under rollback.

## Application and status

This immutable proposal records the intended change. Current lifecycle state and
application authority are carried by the epic tracker and transition receipts, not
by mutating this body after publication.

Four portable harness findings from this episode are recorded in the project
feedback queue: `FND-1B3065946FF8` (no lifecycle brake on repeated
non-convergence), `FND-8386048ACA3C` (a verification brief that carries its own
conclusion is confirmed rather than testing it), `FND-A50ADF64186A` (population
gates that abort on first failure hide the scale), and `FND-CC7A8AAB239A` (an
applied amendment that fails independent verification has no corrective
transition). The first two are visible in this amendment's own history: AM-003 was
the fourteenth artifact to attack the same symptom, and its diagnosis was
confirmed by a scout whose brief carried the conclusion under test.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
