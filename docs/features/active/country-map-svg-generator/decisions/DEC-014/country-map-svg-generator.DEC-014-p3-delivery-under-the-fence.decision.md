---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-014"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-014 — P3 ships under the fence, and takes two named dependencies

## Context

AM-003 and AM-004 sit in amendment state `applied`. That state has exactly one
outgoing edge, `verify`, which requires a passing advisor receipt they can never
earn — both failed independent review (ADV-003, ADV-004) and the corrected
successor AM-005 was authored, applied and verified in their place. The impact
fence blocks every spec transition except `cancel` and `invalidate`, so P2 cannot
close and P3 cannot start through the governed lifecycle. Re-verified at mate CLI
v1.4.3 against the SSOT, not cited from the v1.4.1 filing: `state-machines.yaml`
amendment machine, `transition.go:144` exemption set, `amendment.go:255` fence
computation. Tracked as `WKI-4062B33B8FEA`; filed upstream as `FND-1B950D899CE4`
and `FND-CC7A8AAB239A`.

`spec invalidate P2` is a legal, fence-exempt transition, and it is not an escape.
`impactFence` keys only on amendment state and `affected_specs` and never reads
the spec's state, so P2 would land in `draft`, remain fenced by the same two
amendments, refuse `spec.accept-readiness` with the identical error, and lose its
accepted PLAN-013 and twenty-run history — the exact record reconciliation will
consume. Verifying AM-003/AM-004 to escape would be fabricated evidence.

Separately, ARCH-001's as-built reconciliation rule requires an amendment or a
superseding decision before a new runtime dependency can close out. Amendments are
fenced; decisions are not. P3 needs two.

## Decision

**1. P3 ships as ordinary engineering commits outside the governed run wrapper
while the fence stands.** This extends to P3 the route the owner authorized for
P2 from T1 onward. It is a delivery route, not a lowering of the bar: every task
still ends at a green `make check`, and the reconciliation debt is tracked at
`WKI-4062B33B8FEA` for discharge once mate grows a terminal edge out of `applied`.

**2. No lifecycle mutation is taken to escape the fence.** P2 stays `in_progress`
with PLAN-013 accepted and its run history intact. `spec invalidate` and `plan
invalidate` are both refused for the reason above — they exit a state, not the
fence, and they destroy the reconciliation input.

**3. `spf13/cobra` is accepted as P3's command-surface dependency.** It is the Go
CLI canon and the fleet's own precedent — mate is built on it. VAL-1 requires
enumerating runnable command leaves and proving each one covered; cobra's command
tree *is* that enumeration, so hand-rolling the router would also mean
hand-rolling the registry the gate reads.

**4. `goccy/go-yaml` is accepted as P3's configuration-parsing dependency.**
REQ-3 requires rejecting unknown fields with an actionable diagnostic. Verified
against the library's current documentation rather than recalled:
`yaml.DisallowUnknownField()` is a decoder option, `FormatError` and
`Path.AnnotateSource` render line-column errors with a source excerpt, and the
decoder honours both `yaml` and `json` struct tags, so one struct set serves both
config formats. JSON input is decoded through `encoding/json` with
`DisallowUnknownFields`; whether the YAML parser accepts JSON as a YAML 1.2 subset
was **not** verified and is not relied on.

**5. Both dependencies are runtime, offline and pure Go.** Neither weakens
ARCH-INV-2 (no network, no Node/Python/GDAL) or the P3 delivery constraint that
normal generation resolves embedded corpus, presets and schemas.

## Options considered

1. Wait for the mate fix before starting P3.
2. Escape the fence with a legal but destructive `spec invalidate`.
3. Verify AM-003/AM-004 to clear the fence.
4. **Ship P3 ungoverned with the debt tracked, and record the dependencies as a
   decision (chosen).**
5. Hand-roll the command router and the YAML parser to avoid both dependencies.

## Rejected alternatives

**1 — wait.** The fix is a fleet-wide change to the tool that governs every
project, with its own doctrines, gates and release process, and it ends in `mate
fleet pull` back into this consumer. It is real work on someone else's schedule.
P2 already demonstrated that product delivery under this fence is possible and
honest when the debt is written down.

**2 — invalidate.** Refuted by reading the fence, not by opinion: it does not lift.
It costs PLAN-013 and the run history and buys a different blocked state.

**3 — verify.** Fabricated evidence. The receipt would assert an independent pass
that never happened, in the one record whose whole purpose is to prove one did.

**5 — hand-roll both.** The framework-first rule names this exact shape as debt:
custom code re-solving a concern the canon owns. A hand-rolled subcommand router
re-implements help rendering, flag inheritance and the leaf registry VAL-1 needs;
a hand-rolled YAML reader re-implements unknown-field detection and error
positioning, which is precisely the diagnostic quality REQ-12 asks for.

## Consequences and tradeoffs

- P3's evidence lives in commits, tests and `P3-CHECKPOINT.md` rather than in run
  results and audit verdicts. The checkpoint is therefore load-bearing and is
  written at T1, not at the end.
- The reconciliation debt now spans two specs instead of one. Discharging it means
  reconstructing P2's and P3's delivery against real receipts after the fence
  lifts, which is more work the longer this runs.
- No independent audit gate stands between implementation and closure. The
  compensating controls are the deterministic gates, `make check` green at every
  task boundary, and looking at rendered output — the two P2 defects that mattered
  (Russia as a blob, France as a speck) were both invisible to every metric and
  found by rendering.
- Two dependencies enter the module graph and must build under `GOWORK=off`, since
  that is the build VAL-7 gates and the one a stranger's `go install` reproduces.

## Affected artifacts and owners

- P3 specification and every P3 commit — teamlead.
- ARCH-001 contracts CTR-3 through CTR-6 and boundaries BND-004 through BND-006 —
  unchanged by this decision; it authorizes the route, not a scope change.
- `WKI-4062B33B8FEA` — product-owner; carries the debt and the mate fix sketch.
- `FND-1B950D899CE4`, `FND-CC7A8AAB239A` — routed to the mate fleet queue.
- P2's lifecycle state — explicitly frozen, owned by the reconciliation.

## Validation and revisit trigger

Revisit when mate grows a terminal transition out of `applied` and the fence
lifts: at that point item 1 expires by its own terms and the reconciliation runs.
Revisit item 3 or 4 if either dependency stops being maintained, requires CGO,
reaches the network at generation time, or fails to build under `GOWORK=off`.
This decision does not expire on a date and does not authorize any further
ungoverned scope beyond P3.

## Supersession

Supersedes nothing. Extends to P3 the delivery route the owner authorized for P2
in the P2 checkpoint's fence section; that instruction was session-scoped and this
record replaces it with a governed one. Does not touch DEC-011 or DEC-013.
