---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "run-brief"
id: "RUN-010"
epic: "country-map-svg-generator"
spec: "P2"
run: "RUN-010"
status: draft
profiles: []
concerns: []
inputs: ["PLAN-009", "CTX-RUN-012"]
---
# RUN-010 — Content-addressed diagnostic reproduction brake

## Authority and accepted inputs

- Accepted PLAN-009 revision:
  `a406405b5895244bba918d5c9a1cb913a1144ae08a7785ea441f58cf87baeea8`.
- Accepted PLAN-009 manifest:
  `7fcaf58708caeba7bba7097132c86b898c8bf74d650b9efc91934cf512698e75`.
- Runtime context CTX-RUN-012:
  `ee168bae80b3ff88b43c25d715c5e3d2a92d957fe677c0bfe25ca8f8b5418e80`.
- Product baseline remains PLAN-009's
  `2a7d45953ac857af185a4825ceea7097e0af3a40`, tree
  `a61a0f87122e9094f444fd032dc54a2c04ad3c4b`; commits after it are
  lifecycle-only and do not authorize geometry drift.
- DEC-006, PLAN-008/RUN-009 history and the exact invalidation authority are
  bound content-addressably through PLAN-009.

## Objective and completion boundary

Execute PLAN-009 T0A.1 only. Reproduce the previously coordinator-observed
parameter, phase and 996-case diagnostics from immutable inputs and publish
canonical semantic JSON plus separate timing receipts.

Stop at the T0A.1 checkpoint even when every predicate passes. Do not start the
oracle/ladder spike in T0A.2, seek owner approval, make a product commit,
publish generated artifacts or execute T0B/T1/T2.

## Scope ownership and footprint

The implementer may edit only:

- `internal/geometry/cmd/lodbuild/**`;
- `internal/geometry/diagnostic.go`;
- `internal/geometry/gridphase.go`,
  `internal/geometry/gridphase_test.go`;
- `internal/geometry/lod.go`, `internal/geometry/lod_test.go`;
- `internal/geometry/model.go`, `internal/geometry/model_test.go`;
- `internal/geometry/preset.go`, `internal/geometry/preset_test.go`,
  `internal/geometry/presets/**`;
- `internal/geometry/projection.go`,
  `internal/geometry/projection_test.go`;
- `internal/geometry/lod/v1.recipe.json`;
- `internal/geometry/lod/tool/build.mjs`.

All other paths are read-only. In particular `go.mod`, `go.sum`, serializer
files, `silhouette*`, oracle/v2 recipe/testdata, public pipeline files,
generated tables, integration/mutation tests, docs, `.git` and `.mate` are
forbidden. Existing scratch is shared: preserve unrelated bytes and do not
revert another owner's work. Integration target is `main`; return without
commit, staging or lifecycle mutation.

## Context and doctrine pack

CTX-RUN-012 supplies the current implementation doctrine from P1. PLAN-009 is
the exact HOW and is self-contained for the frozen matrix. Runtime remains pure
Go and offline. Pinned Mapshaper `0.7.44` is maintainer-only and may be restored
solely from the existing locked tool dependency for temporary rebuilds; never
stage `node_modules` or diagnostic output.

## Tasks

1. Record exact P1 corpus, baseline, scratch inventory, projection, recipe,
   Mapshaper, lockfile, tool, Go and command identities. Re-own every reused
   RUN-009 scratch byte before executing it.
2. Implement a diagnostic-only helper that shares production
   canonicalization, topology, protection and final-deviation functions but is
   unreachable from production selection and cannot change its 100-phase
   schedule.
3. Run the exact AQ/un/card parameter matrix from PLAN-009 T0A.1:
   resolutions
   `64,65,66,68,72,80,96,128,160,192,224,256,320,512,1024`;
   weightings `0,0.1,0.3,0.5,0.7,0.9,1`; the 100 lexicographic coarse phases;
   local phase A's 25 rows; and local phase B's 121 rows. Preserve every fixed
   axis and expected typed/numeric outcome exactly.
4. Rebuild one temporary v1 table at resolution `64/256`, weighting `0.7`,
   then scan exactly 996 rows in manifest entity order ×
   `un,de_facto` × `card,hero`. Resolve each preset once, clear its token and
   set `MaxPathBytes=0` without changing resolved layout/quality fields.
5. Emit canonical semantic JSON without wall-clock durations in the frozen row
   order. Run twice from identical inputs and require byte-identical output and
   SHA-256. Emit separate timing receipts keyed to that semantic digest.
6. Return the diagnostic paths, digests, exact commands, aggregate values,
   focused test output and explicit PASS/CONTRADICTION verdict. Do not continue
   into T0A.2.

## Validation and evidence

Run focused projection, gridphase, selection, topology, protection, preset and
lodbuild tests plus compile-only checks for the touched packages.

The exact row outcomes and `1e-12` comparison tolerance are PLAN-009
T0A.1 lines 141–193. The 996-case aggregate must reproduce:

- selected: `card/compact=149`, `card/standard=237`,
  `card/source=112`, `hero/standard=380`, `hero/source=118`;
- historical-max overruns: `56,216,76,222,55` in that same key order;
- maximum card path `477784` at `RU/un/source`;
- maximum hero path `675606` at `CA/un/source`.

Evidence must include the two semantic file hashes and byte comparison,
separate timing receipt hashes, exact case/typed-failure counts, minima and
maxima, and whether any parameter-only candidate makes all 996 rows satisfy
the new `2200/7500` P2 caps.

## Constraints and forbidden actions

Keep projection v1, `q=0.01`, compact/standard maximum effective scales
`240/700`, resolution `64/256`, weighting `0.7`, current topology/protection,
natural/contain semantics and frozen serializer hashes
`c9ee9019b9155522acef16cd405ef0a72c93a4d519493a86061617cc2e034345`
and
`0142dea0db810bfa97e4bb23d30cdf039faad860e8e1354a4111b804ee584c8f`.

No country/preset branch, tolerance/q/budget relaxation, production phase
change, oracle/threshold invention, post-hoc optimizer, MakeValid,
union/buffer, CGO, new runtime dependency, generated publication or broad
cleanup is permitted. Diagnostic path budgets remain disabled only as
specified; production behavior must not gain a bypass.

## Stop escalation and resume conditions

Stop without product commit and report CONTRADICTION when any PLAN-009 T0A.1
predicate is true: input/hash/count drift; expected numeric drift above
`1e-12`; typed failure mismatch; a phase accepted at or below `1.536`;
semantic byte/digest drift; aggregate/maxima drift; or a parameter-only
candidate satisfies topology/protection and all new P2 caps.

Also stop on required edits outside the lease, serializer/dependency drift,
production access to the diagnostic helper, or inability to reproduce from
recorded commands. Do not reinterpret a contradiction as calibration freedom.
On PASS, preserve exact content handles and wait for a separate T0A.2 run.

<!-- MATE:extensions — generated by composition from selected profiles and concerns -->
