---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-016"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-016 — The ladder's first-load cost is accepted, and the catalog byte distribution is recorded

## Context

P2 measured the committed detail ladder at **707 ms** to parse on first use and
routed the decision to P3, where the real invocation shape would be known
(`WKI-6638BACD6E20`). It now is. Measured through the built binary, median of
five runs each:

| Invocation | Wall clock |
| --- | --- |
| `version` | 0.596 s |
| `explain --iso FR` | 0.592 s |
| `generate --iso FR` (one asset) | 1.288 s |
| `generate` (whole catalog, 248 assets) | 3.3 s |

Two facts fall out of those numbers and settle the question.

**The discovery commands never pay the cost at all.** `version` and `explain`
come in at 0.59 s, indistinguishable from each other, because the ladder is
loaded lazily by `Generate` and neither command generates. Those are the commands
an agent runs most, and they were already free.

**In a batch the cost is invisible.** The whole catalog is 3.3 s for 248 assets.
The load happens once; per-asset work after the process is warm is roughly 8 ms.

That leaves the single-asset invocation, where the load is about 0.69 s of a
1.29 s run. But **0.596 s of that run is process start plus corpus decode**,
which lazy per-geometry ladder decode cannot touch. The whole prize is taking a
1.29 s single-asset command to roughly 0.7 s.

Separately, the same full-catalog run is the first measured byte distribution
through the shipped path, at the `card` profile with the default configuration:

| | path bytes | file bytes |
| --- | --- | --- |
| median | 1964 | 2165 |
| p95 | 2183 | 2384 |
| max | 2197 | 2397 |

248 assets written, 1 typed no-artifact (UM), 476 106 bytes in total. QAB-2 for
`card` asks for a median of at most 0.6 KB, p95 at most 1.2 KB and max at most
2.5 KB, as uncompressed file sizes.

## Decision

**1. The ladder's first-load cost is accepted. Lazy per-geometry decode is not
implemented, and `WKI-6638BACD6E20` closes.** The measurement above is the
reason: the commands an agent runs repeatedly never pay it, a batch amortizes it
to nothing, and the residual on a single asset is bounded below by a 0.596 s
process floor the optimization cannot reach. Implementing it would add a lazily
decoded index to a hot path in exchange for roughly half a second on the one
invocation shape nobody repeats.

**2. Revisit only on evidence, and the evidence is named.** If a real workflow
appears that invokes the binary once per asset in a loop — a file watcher, a
per-request server, a build step that shells out per country — the arithmetic
changes and this is reopened. Until such a workflow exists, it is speculation.

**3. The measured byte distribution is recorded as the baseline P4 inherits, and
the QAB-2 median and p95 are not treated as delivery blockers today.** ARCH-001
says exceeding a **maximum** blocks delivery; nothing is over the 2500-byte
maximum. It also says a first full-catalog baseline may tighten the budgets and
that architecture budgets are updated only from accepted measured evidence
recorded by P4. This is that evidence arriving early, not a licence to ignore it.

**4. The median overshoot has a named cause, and it is structural rather than
accidental.** The ladder selects, per geometry and band, the **finest**
oracle-passing candidate — P2's own description. Finest-that-fits means every
entity lands just under the cap by construction, so the distribution clusters at
the ceiling instead of spreading below it. QAB-2's median and p95 were written
for a distribution where most countries are simple and cheap. The two are not
reconcilable by tuning; one of them has to change.

**5. That reconciliation is P4's, with the owner.** It is a product judgement
about page weight — 248 cards at roughly 2.1 KB is 476 KB uncompressed — against
fidelity, and it needs a rendered page to judge, exactly as DEC-008 and DEC-013
did. Tracked as `WKI-33B6EC2482B3`.

## Options considered

1. Implement lazy per-geometry decode now.
2. **Accept the load cost on the measurement, and record the byte baseline
   (chosen).**
3. Accept the load cost and say nothing about the byte distribution.
4. Tighten the ladder's selection rule now to hit QAB-2's median.

## Rejected alternatives

**1 — implement it now.** The measurement does not support it. It would buy
about 0.6 s on the least repeated invocation shape while adding state to the
selection path, and P2's own scars are about complexity in that path.

**3 — stay quiet about the bytes.** The whole point of routing this to P3 was to
measure through the shipped path. Measuring and not recording would leave P4 to
rediscover a 3.6x median overshoot, and the manifest already carries the numbers
that make it obvious.

**4 — tighten the selection rule now.** It changes what every asset in the
catalog looks like, on a fidelity judgement no metric can make. DEC-008 and
DEC-013 both established that this class of decision is an owner approval over a
rendered catalog, not a threshold change. Doing it here, in the CLI spec, from a
number, would repeat exactly the mistake P2 spent twenty runs on.

## Consequences and tradeoffs

- A single-asset invocation stays at roughly 1.3 s. For an agent this is one
  round trip; for a producer it is imperceptible against opening the result.
- P4 inherits a concrete baseline rather than an unexamined budget, and the
  manifest reports both file and path bytes per asset so the distribution can be
  recomputed from any run.
- The tension between "finest that fits" and QAB-2's median is now written down
  where the next session sees it, instead of being discovered again from a
  surprising page weight.
- If the reconciliation lands on tightening the ladder, the committed artifact is
  rebuilt and every asset changes — which is why it is P4's with an owner
  approval and not a passing edit.

## Affected artifacts and owners

- `WKI-6638BACD6E20` — closed by item 1; owner teamlead.
- `WKI-33B6EC2482B3` — the QAB-2 reconciliation; owner product-owner, due P4.
- ARCH-001 QAB-2 — unchanged here; P4 updates it from accepted evidence.
- `internal/geometry/ladder.go` — unchanged; no lazy decode is added.
- CTR-006 manifest — carries `bytes` and `path_bytes` per asset, which is what
  makes the baseline recomputable rather than a number in a document.

## Validation and revisit trigger

Item 1 is revisited if a per-asset invocation workflow appears, and the
measurement above is the comparison to redo. Item 3 expires when P4 records its
own full-catalog distribution and the owner accepts a budget change; until then
the numbers here stand as the baseline.

## Supersession

Supersedes nothing. Closes the deferral P2 recorded in `WKI-6638BACD6E20` and
opens `WKI-33B6EC2482B3`. Depends on DEC-008 and DEC-013 for the principle that a
fidelity or threshold change is an owner approval over rendered output.
