---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-008"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-008 — Representative geometry owner approval

## Context

PLAN-010 and RUN-013 require a durable product-owner approval after the
representative geometry matrix passes its machine gate and before T0B may scan
the full catalog. The owner reviewed the representative sheet in this task and
explicitly confirmed the proposed visual policy and continuation.

The accepted machine evidence is RUN-013 gate `GATE-458D4D778A57`
(`792fd75998656c10ecbd7bfdea760685c8f5880d613574b50210aa421339a0ba`)
and coder result `RESULT-021`
(`6c34549769804c3a0751954d3db9ee3aa271d0b20511ca34c9456d61389082ee`).
The approval boundary is only the representative T0A freeze; it does not waive
the separate full-catalog machine and owner gates.

## Decision

Approve and freeze the exact RUN-013 representative geometry policy:

- compact protected-component contribution threshold `110`;
- standard protected-component contribution threshold `900`;
- global silhouette oracle minimum `IoU=0.40` and `recall=0.42`;
- oracle digest
  `f2c9cd32e806985be55193e564942e716bb42c72e01d7405c976ddf36666d56f`;
- recipe digest
  `da3ff9d331df37882bb2ed4159e9f14ef2aaeedf2b9fcff9e77c79fed4d751c5`;
- representative manifest digest
  `b2c467b08fb32603a312918da05ce83da355499fb8cba7f9e73c5a8e916995da`;
- representative contact-sheet digest
  `cf1e568c93646047791b45a329e50a1308683e8c3aae03ee0944654db3718f9f`;
- representative machine receipt digest
  `fcc1bc307509d55aa281591f1cce0d0e78f73f323a41fc6d78e665a56c2dfd10`.

The empty `group-anchors.v1` data set is approved as an exact no-op. Any later
non-empty group-anchor entry requires its own owner review and a rerun of the
dependent representative and full-catalog gates.

## Options considered

- Approve the exact machine-passing thresholds and sheet.
- Reject the sheet and return to bounded pre-freeze calibration.
- Request a reviewed group anchor before full-catalog work.

## Rejected alternatives

Recalibration and an immediate group-anchor override are rejected because the
owner accepted the current representative silhouettes and no reviewed example
demonstrated a missing identity-defining group. Machine failures cannot be
waived and were already absent from the selected T0A gate.

## Consequences and tradeoffs

T0B may begin from these exact content identities and must not change the
frozen oracle, recipe, thresholds, contribution model, identity policy,
candidate order, budgets, projection, quantization or serializer. The full
catalog may still expose a failure; that returns to this checkpoint rather than
silently relaxing the frozen policy.

This approval does not publish artifacts, create a T0 product commit, approve
the 996-output catalog sheet or authorize a budget exception.

## Affected artifacts and owners

- P2 / PLAN-010 T0A and RUN-013: representative owner checkpoint satisfied.
- P2 / PLAN-010 T0B: may start from the exact frozen identities above.
- Product owner: approval authority for this visual checkpoint.
- P2 geometry owner: machine integration and full-catalog evidence.
- DEC-007: unchanged and still governs scale-aware protected visibility.

## Validation and revisit trigger

The approval remains valid only while the gate, oracle, recipe, manifest,
sheet, receipt and policy digests above remain exact. Revisit and obtain a new
owner decision if any digest changes, a threshold or oracle floor changes, a
non-empty group-anchor override is proposed, the full catalog reveals an
identity-defining omission, or a T0A-dependent regression fails.

## Supersession

This decision does not supersede an architecture or product-policy decision. It
records the product-owner approval checkpoint required by DEC-007 and PLAN-010.
