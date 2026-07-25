---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-013"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-013 — The card is fitted to the silhouette it draws

## Context

DEC-011 fixed the default framing as the full territorial claim and deferred the
definition of a second, tighter mode to the T5 contact sheet. That sheet now
exists, and it settled the question in a way the hypothetical options could not:
the problem was never that a second mode was missing, it was that the first mode
was wrong.

France made it unmistakable. Its UN geometry carries eleven components spanning
119 degrees of longitude; a card draws two of them, metropolitan France and
Corsica. Everything else is below the band's visibility threshold and
deliberately not drawn. Yet the layout was fitted to the whole projection, so the
card reserved space for territory the oracle had already excluded and the drawn
silhouette occupied **11.9% of the frame width** — France as a speck in an empty
box. Twenty cards were in that state.

The geometry was never wrong. Fidelity is already judged against the visible
reference — that is what `CompareSilhouettes` rasterizes — so the layout fit was
the one remaining place still reasoning about the whole claim.

Two measurements decided the shape of the fix rather than the reasoning behind
it. Fitting to the **visible reference** is wrong: a component can sit in the
candidate while below the visibility threshold, so a frame sized to the visible
set leaves it outside the viewBox, containment fails, and the search walks to a
coarser rung until the candidate loses it — on France that traded Corsica and 936
bytes for a single 35-point blob. Fitting to the **candidate** is right, and it
is what this decision adopts.

## Decision

**1. The layout is fitted to the candidate, in the oracle and in the runtime
alike.** The card frames the geometry it draws.

**2. The two sides must agree, and that is not a detail.** The band is still
chosen from the full-source fit, which keeps the rule non-circular — the build
assigns bands by preset and a card request lands in compact either way — and the
ladder candidate is then refitted, canonicalized and rendered under its own
transform and quality. Fitting only the build would commit one frame and serve
another, which is the failure class this spec exists to end.

**3. The compact protected-component contribution threshold moves from 110 to
111.** Croatia is the reason, and the reason is exact: under the tighter frame
its `un` card missed the 2200-byte cap by two bytes at rung 79.63. An exhaustive
sweep of eleven additional rungs, down to 0.05 granularity, found no rung that
both fits the cap and keeps the disputed component — it survives at 79.6 and dies
at 79.55. At 111 that component, contributing exactly 110, falls below the
threshold and is recorded as `subscale` with provenance, which is DEC-007's
mechanism working rather than a check being silenced.

**4. DEC-008's approval is renewed on the full catalog rather than the
representative sheet.** DEC-008 froze an oracle digest that has since changed
twice — once when DEC-012 edited the rasterizer, once here — so its stated
validity condition was already broken independently of this decision. DEC-008 is
not a calibration procedure; its chosen option was to *approve the exact
machine-passing thresholds and sheet*, so renewing it means the owner looking at
the rendered catalog and accepting it. That is what happened.

No byte budget, IoU floor, recall floor, rung set, quantization or corpus byte
changes.

## Options considered

1. **Fit the card to the candidate, and move the compact threshold by one
   (chosen).**
2. Fit to the candidate and accept that Croatia loses its `un` card, leaving it
   with no compact card at all.
3. Add further intermediate rungs to close Croatia's two-byte gap.
4. Leave the framing as the full claim and build a separate cropping mode later.

## Rejected alternatives

**Option 2** was rejected by the owner. DEC-010 accepted losing one Croatia card
after weighing it; losing both, silently, as a side effect of a framing change is
a different decision and a worse one.

**Option 3** was tested, not argued. Eleven extra rungs between 79.63 and 64,
including fractional steps at 0.05 granularity, all lose the protected component.
The window is genuinely closed; no rung set closes it.

**Option 4** was rejected because the sheet showed the premise was wrong. A
second mode would have been a workaround for a default that framed cards for
geometry it does not draw, and it would have left the twenty sparse cards sparse
until that mode shipped.

## Consequences and tradeoffs

Measured against the previous artifact:

| | Before | After |
| --- | --- | --- |
| Totals | 993 pass / 3 no_artifact | **993 / 3, unchanged** |
| Cards lost or gained | — | **none / none** |
| Rungs changed | — | 12 of 996 |
| Byte delta | — | median 0 |
| Cards filling under a fifth of their frame | **20** | **0** |

- France is the hexagon with Corsica. Chile is tall and narrow. Croatia keeps its
  card. The Netherlands gains parts rather than losing them.
- **Two rows drop one component each**: `AU/de_facto/compact` from 5 parts to 4,
  and `NZ/un/compact` from 6 to 5. Both were reviewed and accepted.
- `TestRepresentativeNaturalRatiosAndArbitraryFrames` returns to its original
  assertion. It always said a Chile card is tall and narrow; that was briefly
  rewritten to demand the opposite while the frame followed the whole claim. The
  original expectation was right and the pipeline was wrong — Chile measures 0.29
  rather than 0.98.
- The whole binding chain was rebuilt: rasterizer digest into the oracle, oracle
  into the ladder recipe, recipe into the artifact and every stored candidate.
- DEC-011's deferred second framing mode is **no longer needed as a framing
  mode**. What remains of it is component *selection* — "draw this country
  without that island" — which is a different capability and stays routed to P3.

## Affected artifacts and owners

- `internal/geometry/silhouette.go` — the oracle's fit.
- `internal/geometry/lod.go` — the runtime's matching fit, quality and render.
- `internal/geometry/lod/silhouette-oracle.v1.json` — threshold 111, rebound
  rasterizer digest; `ladder.recipe.json` and `ladder.artifact.json` rebuilt.
- `readiness/T5-catalog-render.contact-sheet.json` — the reviewed evidence.
- DEC-008 — approval renewed on this catalog; DEC-010 — unchanged, Croatia's
  `de_facto` card remains a typed no-artifact; DEC-011 — its framing default is
  superseded by item 1 here, its crop seam is not.
- Accountable owner: product owner (accepting principal).

## Validation and revisit trigger

Validated by the full gate set against the rebuilt artifact: 215 passed, 0
failed, including the offline VAL-6 recomputation of every committed row, the
Node determinism rebuild proving the artifact reproduces byte-identically, and
the shipped-catalog gate proving `Generate()` serves the ladder inside every band
cap. Rendered evidence: 993 rendered, 3 typed no-artifact, 0 errors, 0 source
fallbacks, 0 over the band cap.

Revisit if the corpus or the Mapshaper pin changes, since contribution is
measured in the frame and the frame now depends on the candidate; if any further
digest in the chain changes, which renews the DEC-008 obligation again; or if the
P3 component-selection seam lands, since a selected subset is a new frame and
must be judged rather than inherited.

## Supersession

Supersedes DEC-011 item 1 only — the default framing. DEC-011's remaining items
stand: cropping is an explicit seam over the selected candidate, never an
automatic heuristic; selection is expressed through the identity vocabulary and
never a country literal; and the acceptance rule for a cropped variant is still
an open sub-decision. DEC-006, DEC-007, DEC-009, DEC-010 and DEC-012 are
unchanged.
