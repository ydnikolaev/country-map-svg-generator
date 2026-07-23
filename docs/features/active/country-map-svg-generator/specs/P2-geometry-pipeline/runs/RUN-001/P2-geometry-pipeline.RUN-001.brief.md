---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-001"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-001"
status: running
profiles: []
concerns: []
inputs: ["PLAN-001", "CTX-RUN-003"]
---
# RUN-001 — Run brief

## Authority and accepted inputs

- Accepted plan revision: `PLAN-001` at SHA-256
  `3c7147cd86a0be2c65236fbcc9939ba5c83961b5541c2bc12745385328c4b67e`.
- Accepted plan manifest: `PLAN-001` at SHA-256
  `340693ac032d454609087b3656caa01b0d7580684d01790dbce3543ebf6c4070`.
- Run context pack: `CTX-RUN-003` at SHA-256
  `6c69281d77c8acd3b77e1ac76e4edf81e2f3d1c79a60a7135c861c4a96560029`.
- Accepted-plan baseline is commit
  `80c50352615b5d55a6efdc2fdffc5d9db3c9fa88`, tree
  `53948381a787202e9edf5e200749743e6220325c`; implementation integrates from
  current lifecycle-only descendant
  `4ed1100bcce7dbe040350b15e863cc9b1a556893`.
- Verified AM-001 and AM-002 are binding. Run readiness is intentionally not
  embedded in this candidate; `mate run start` binds it atomically.

## Objective and completion boundary

Implement and verify CTR-002: deterministic presentation-free soft-organic
geometry for every P1 entity/profile. Default `tight` output follows each
country's natural projected aspect ratio; explicit `contain` accepts arbitrary
safe frames and never stretches geometry. The run includes fitted-scale quality,
versioned card/hero convenience presets, protected retention, inspectable removal,
bounded smoothing, canonical commands, and optional markers.

Finish only when D3 oracle, topology/mutation, natural-ratio, arbitrary-frame,
fitted-scale, override, marker, full-corpus, approval-sheet, and offline project
gates pass. This run does not implement SVG/XML styling, CLI/config parsing,
browser delivery, public packaging, or the deferred `detail` profile.

## Scope ownership and footprint

Writable product paths are exactly the accepted PLAN-001 union:

- `go.mod`, `go.sum`, `Makefile`;
- `internal/geometry/**`.

`internal/catalog/**` is read-only input. Forbidden product roots are `data/**`,
`internal/config/**`, `internal/render/**`, `cmd/**`, `docs/guides/**`,
`docs/operations/**`, `.goreleaser.yaml`, `.git/**`, and `.mate/**`.

The integration target is current `main`. The implementer returns atomic product
commits or one bounded final product commit and never commits lifecycle evidence.

## Context and doctrine pack

`CTX-RUN-003` supplies the current implementation context selected by the epic
registry. The accepted plan and P2 specification additionally bind natural-aspect
layout, arbitrary distortion-free containment, fitted-geometry-scale quality,
P1 protected metadata, and CTR-002 diagnostics/provenance.

The repository Go profile is active. Normal generation and `make check` use
`GOWORK=off`, are offline, and do not require Node, Python, GDAL, GEOS, PROJ,
CGO, user-home data, locale, or system geodata. Node is permitted only for
explicit maintainer refresh of committed D3 fixtures.

## Tasks

1. T1: establish CTR-002 models, typed diagnostics, generic `tight`/`contain`
   layout requests, fitted-scale auto quality, preset data, dependencies, and the
   P1 adapter.
2. T2: implement normalization, centered Lambert azimuthal equal-area projection,
   adaptive subdivision, natural tight bounds, uniform centered contain fitting,
   and the shared marker transform.
3. T3: implement protected retention, inspectable removal/collapse accounting,
   validated simplification, bounded softening, marker anomalies, versioned
   overrides, canonical commands, and deterministic pipeline output.
4. T4: add pinned D3 Geo oracle fixtures and Go parity tests without a normal Node
   or network dependency.
5. T5: add full-corpus, sabotage, natural-ratio, arbitrary-frame, paired-frame
   fitted-scale, visual approval, byte-budget, and project-ceiling evidence.

T1 → T2 is serial. T3 and T4 may proceed independently after T2 only if their
leases do not overlap. T5 follows both.

## Validation and evidence

- `geometry-d3-parity`: projection, rotation, subdivision, bounds, fitting, and
  marker samples stay within declared oracle deltas.
- `geometry-topology-teeth`: bow ties, escaped holes, collapsed rings, smoothing
  crossings, protected removal, implicit disappearance, and unchecked
  simplification all red.
- `geometry-visual-and-budget`: `RU`, `CL`/`AR`, and `AU` prove wide/tall/near
  square natural ratios; tiny/huge portrait/landscape/square frames prove uniform
  centered containment; paired frames prove fitted-scale quality; approval
  digest, diagnostics, and budgets have teeth.
- `geometry-overrides`: invalid orientation, padding, tolerance, anchor, marker,
  applicability, or unknown field fails closed.
- `geometry-markers`: a divergent transform, non-finite/outside point, or
  suppressed anomaly fails.
- Focused commands and ceiling:
  `GOWORK=off go test ./internal/geometry -count=1`,
  `GOWORK=off go test ./... -count=1`, and `make check`.

The run result names exact commands, outcomes, changed paths, commits, deviations,
approval evidence, and feedback handles. Narrow green tests do not replace the
project ceiling.

## Constraints and forbidden actions

- Never stretch geometry or expose independent X/Y geometry scales.
- Never hard-code `card`/`hero` or country switches in the engine; presets and
  exceptions are versioned data resolving through the generic contract.
- Never synthesize, bridge, buffer, union, or manually redraw geography.
- Never let simplification/softening bypass final topology validation, protected
  retention, or explicit disappearance accounting.
- Do not mutate P1 corpus bytes or reinterpret boundary/capital policy.
- Do not add runtime Node/Python/GDAL/GEOS/PROJ/network dependencies.
- Do not weaken a validator, oracle tolerance, mutation fixture, or approval
  digest to make a failure green.

## Stop escalation and resume conditions

Helper naming, numerically equivalent projection implementation within oracle
bounds, file splits inside `internal/geometry`, and fixture layout are local HOW
and must be reported. Stop before mutation for any new write root, dependency
substitution, changed formula/preset/oracle tolerance/point ceiling/override
field, hard-coded preset branch, projection-policy change, protected-policy
change, fabricated geography, `detail` scope, forbidden runtime, or P3/P4 concern.

After a few speculative fixes for the same failure, reduce to a minimal fixture
and read the owning library documentation. On interruption, resume from accepted
PLAN-001, this brief, `CTX-RUN-003`, current product commit/diff, and last exact
validation output; never reconstruct authority from chat history.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
