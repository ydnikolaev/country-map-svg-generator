# P3 checkpoint — resume point for a fresh session

Written at T1 rather than at the end, deliberately. P3 delivers ungoverned under
DEC-014, so there is no run result, no audit verdict and no plan manifest to
resume from — this file and the commits are the whole record. Read it first.

## The one-line status

**T1 done.** The binary exists, the command surface and its typed failure
taxonomy are in place, and the e2e harness that will hold every later task is
built and proven to bite. Only `version` is registered so far; the other six
commands land with the tasks that implement them.

## Why P3 has no governed tracker

AM-003 and AM-004 sit in amendment state `applied`, which has exactly one
outgoing edge — `verify` — that they can never earn. The impact fence blocks
every spec transition except `cancel` and `invalidate`, so P2 cannot close and P3
cannot start through the lifecycle. **DEC-014** authorizes P3 to ship as ordinary
engineering commits meanwhile and accepts the two dependencies P3 needs.
**`WKI-4062B33B8FEA`** carries the debt, the four source locations of the mate
fix, and the sketch to apply.

Two escapes were considered and refused, both on evidence rather than caution:

| Escape | Why not |
| --- | --- |
| `mate amendment verify AM-003` | Fabricated evidence. The receipt would assert an independent pass that never happened. |
| `mate spec invalidate P2` | Legal and fence-exempt, and it does not lift the fence. `impactFence` (`internal/work/amendment.go:255`) reads only amendment state and `affected_specs`, never the spec's state, so P2 lands in `draft`, stays fenced by the same two amendments, and loses PLAN-013 and its twenty-run history. It exits a state, not the fence. |

`mate decision scaffold/validate/accept` **does** work under the fence at epic
scope — it is how DEC-010 through DEC-015 landed. Use it for anything that needs
a governed record. `mate backlog add` and `mate feedback add` work too.

## Scope, as decided

**DEC-015 moved the crop and component-selection seam out of P3** to a successor
spec, taking DEC-011 item 6's explicit escape. DEC-011 item 5's open sub-decision
(how a cropped variant is judged) went with it. The one obligation this leaves on
P3: **the config schema v1 must reserve the `selection` key namespace** so the
seam arrives later as an additive change rather than as a schema major. That
reservation and its failing-unknown-field assertion are due in T2.

Everything else in the P3 spec stands: 13 requirements, 8 validation obligations,
5 acceptance criteria.

## The `make check` constraint — read before adding any test

`make check` passes `-p 1`, and that is load-bearing: several tests assert
wall-clock brakes that measure CPU contention instead of the code under test when
packages run in parallel. So **every new package adds to the sum, not to the
max**. `geometry` ≈ 322 s and `lodbuild` ≈ 373 s already sit near the default
10-minute per-package timeout; adding a full-catalog sweep to either breaks the
build before it fails an assertion.

P3's budget is **≤ 90 s added to `make check`**. Measured inside a real `-p 1`
run: **2.142 s** for the whole `cmd/country-map-svg-generator` package including
the `GOWORK=off` binary build. For context in the same run, `catalog` 28.7 s,
`geometry` 353.3 s, `lodbuild` 267.6 s. The levers that keep P3 inside its
budget:

- The ordered product loop (VAL-5) runs a **selected set** of ISO codes chosen for
  shape, never the catalog.
- Full-catalog double-generate is **P4's** E2E per ARCH-001's allocation table.
  P3 proves the contract; P4 proves it at scale.

Also: summarize package-level failures, not just test-level ones. A `go test
-json` package failure carries no `Test` field, so a test-only filter reports a
timed-out package as a clean run. That bit P2 twice.

**And never run `make check` through a pipe.** `make check | tail -60` reports the
exit code of `tail`, which is always 0, and buffers the entire stream so the log
is empty until the run ends and truncated afterwards. T1 read a red gate as green
that way — the same failure-summarizing mistake as above, in a new shape. Run
`make check` on its own and read the whole output.

## T1 — done

### The binary and its boundary

`cmd/country-map-svg-generator` owns BND-004 and nothing else. There is no
`internal/cli` package: ARCH-001's boundary table has no such boundary, and
adding one would be a new package crossing an ownership boundary, which its
as-built rule makes an amendment-level change. The command surface is `package
main` with one file per command; `internal/config` (BND-005) and `internal/render`
(BND-006) are the libraries it wires together.

### The failure taxonomy is the contract, not an afterthought

`exit.go` holds CTR-004's typed exit classes and the `--json` envelope. The
numeric codes are frozen at v1 and appended to, never renumbered — a build
pipeline branches on them:

| Code | Class | What it means |
| --- | --- | --- |
| 0 | ok | |
| 1 | usage | flags, arguments, command selection — fix the invocation |
| 2 | config | the document parsed but violates the schema |
| 3 | data | corpus side: unknown ISO, missing geometry, corpus mismatch |
| 4 | render | geometry or serialization failed on otherwise legal input |
| 5 | validation | emitted output failed its own structural gate |
| 6 | budget | a byte budget was exceeded |
| 7 | filesystem | staging, publication, path resolution |

`budget` is separate from `validation` on purpose: it is the one failure a caller
fixes by asking for less detail rather than by fixing a defect.

Two details that are load-bearing rather than stylistic:

- **The top-level fallback maps to `usage`, not to an internal class.** Command
  bodies always return a `*CLIError`; anything else reaching the top came from
  cobra's own parsing and routing, which is a caller mistake. Mapping the fallback
  to `render` would report a mistyped command as a rendering failure.
- **The error path reads the raw arguments for `--json`, not the parsed flag.**
  An unknown flag fails *before* the parsed value exists, and that is exactly the
  case where an agent still needs the envelope.

### The e2e harness

`main_test.go` builds the real binary once per package run with `GOWORK=off` and
puts it on `PATH` for testscript. It deliberately does **not** use
`testscript.RunMain`: that re-executes the test binary, which `go test` compiled
with whatever workspace was ambient, and the Go CLI validation profile asks for
the opposite on both counts — black-box the built artifact, and build it the way
a stranger's checkout would so a green suite cannot rest on a `go.work` a
consumer will not have. The cost is one `go build` per package run.

`RequireExplicitExec` is on, so a script can never accidentally exercise a host
tool instead of the CLI.

### VAL-1's gate, and why it is not vacuous

`surface_test.go` enumerates runnable command leaves **from cobra's own tree**,
not from a maintained list — a gate reading its own list of commands proves only
that the list matches itself. For each leaf it requires four scripts:
`<leaf>_help`, `_happy`, `_invalid`, `_json_error`. It checks the other direction
too, so a renamed command cannot leave its old scripts behind still green.

The enumeration has its own tooth. `TestLeafEnumerationCatchesANewCommand` runs
it against a synthetic tree carrying every shape it must tell apart: a leaf, a
parent that is not a leaf, a nested leaf, a hidden command, `completion`, and a
grouping command with no `Run`. Without it, both coverage gates stay green when
the enumeration returns nothing.

Both gates were verified to redden by mutation, not by inspection: removing
`version_happy.txtar` fails the coverage gate, and adding an orphan `.txtar`
fails the reverse gate.

### Command registration order

Registering a command is what makes the VAL-1 gate demand its scripts, so
commands are registered by the task that implements them, not stubbed up front. A
stub would satisfy the gate with scripts asserting that the command does nothing.

| Command | Registered by |
| --- | --- |
| `version` | T1 — done |
| `init`, `validate`, `explain` | T2 |
| `generate` | T4 |
| `inspect`, `preview` | T5 |

### Dependencies added, and the gate that noticed

`spf13/cobra v1.10.2` and `rogpeppe/go-internal v1.15.0` (testscript), both
accepted by DEC-014. `go mod tidy` also pulled `golang.org/x/tools` (txtar, a
testscript dependency) and dropped `golang.org/x/image`, which nothing imports
any more — verified by grep across all `.go` files including tests. Both
dependencies are pure Go and offline; neither weakens ARCH-INV-2.

**That reddened `TestDiagnosticIdentityDriftChecksEveryRecordedLoadBearingIdentity`
in `lodbuild`, and the response is worth reading before the next dependency.**
`loadDiagnosticIdentities` records `go.mod` and `go.sum` as file digests and
compares them against frozen constants. Those two fields are the coarsest in the
record: they cover the whole module graph, including dependencies nothing under
`internal/geometry` imports. Every other live field is scoped to something the
diagnostic actually consumes.

The two constants were re-recorded rather than the predicate weakened, on
evidence rather than convenience: exactly one test failed in the whole suite, the
`internal/geometry` package stayed green at 353 s including the shipped-catalog
gate, and `lodbuild`'s full 30-resolution × 249-entity mapshaper ladder rebuild
reproduced byte-identically under the new graph. The diagnostic's outputs are
guarded independently of the module hashes by the frozen expectations in
`evaluatePredicates` (`expectedResolution`, `expectedWeights`, `expectedSelected`,
`MaximumCard`/`MaximumHero`) — had the dependency change altered any computed
value, those would have fired on their own.

Expect this to fire once more, in T2, when `goccy/go-yaml` lands **with its
import**. Do not pre-add it: `go get` on a module nothing imports is dropped by
the next `go mod tidy`, and the requirement moves between the direct and indirect
blocks when the import appears, so the hash changes then anyway. Re-record with
the same evidence. `WKI-1DA58E0FE741` carries the underlying over-specification;
narrowing the predicate weakens an existing gate and belongs to P2's
reconciliation, not to P3.

Two things about that record that save a re-derivation:

- `BaselineCommit`, `BaselineTree`, `ScratchBefore` and `Command` are inert.
  `loadDiagnosticIdentities` assigns them from the same constants the drift check
  compares against, so they can never drift.
- `diagnosticSourceInventoryDrift` is a **self-consistency** check, not a
  frozen-value one: it verifies the path list, the hash well-formedness and the
  aggregate over the list it just built. So editing `diagnostic.go` — which is
  itself in `diagnosticSourcePaths` — does not create a fixed-point problem.

## The vocabulary, settled before T2 freezes the schema

Two orthogonal axes, and the epic's documents and the Go packages spell them
differently. Getting this wrong is a schema major, so it was settled first.

| Axis | Values | Config key (CTR-3) | Go field in `geometry.Input` |
| --- | --- | --- | --- |
| detail / byte budget | `card`, `hero` | `profile` | **`Preset`** |
| boundary posture (DEC-002) | `un`, `de_facto` | `boundary` | **`Profile`** |
| named inheritable bundle | any | `extends` | — (config only) |

The Go names invert both terms relative to the documents. The documents are
authoritative for the config, and the decisive citation is **P4 REQ-7** —
"producer guidance covers *profile selection*, *preset inheritance*" — which
makes `profile` and `preset` distinct concepts, with `preset` meaning the
`extends` bundle. **P4 VAL-2** corroborates: it measures "median/p95/max **per
profile**" against QAB-2, whose columns are `card` and `hero`. So spelling
card/hero as `preset` in the config would have collided with `extends`.

Three rules follow, each with a test due in T2:

1. **Both enums are read from their source, never written as literals.** `profile`
   from `geometry.Presets()`; `boundary` from geometry's accepted set — **not**
   from `catalog.Manifest.Profiles`, see below.
2. **One adapter maps config → `geometry.Input`, and a test fails if the two are
   transposed.** A document with `profile: hero, boundary: de_facto` must produce
   `Input.Preset == "hero"` and `Input.Profile == "de_facto"`. Transposed, it
   would render the wrong boundary posture and still satisfy every byte budget —
   a DEC-002-class defect invisible to every metric.
3. **Budgets are never restated in config.** They come from
   `geometry.Presets()[name].MaxPathBytes`.

### The `de-facto` / `de_facto` split — a live trap

The corpus and the runtime disagree on the spelling:

| Where | Spelling |
| --- | --- |
| `catalog/compile.go:125` — `Manifest.Profiles` | `["un", "de-facto"]` — **hyphen** |
| `catalog/validate.go:219` | accepts the hyphen |
| `catalog/model.go:27` — JSON tag | `de_facto` — underscore |
| `geometry/adapter.go:23`, `geometry/pipeline.go:64` | `de_facto` — underscore |

Nothing translates between them. So sourcing the `boundary` enum from
`Manifest.Profiles` would accept `de-facto` and hand it to geometry.

Two guards already catch that, and one gap remains — stated at the strength it
was actually traced, not at the strength that would make the point louder:

- **`geometry.InputFromCatalog(corpus, alpha2, profile, preset)` fails closed.**
  Its `switch profile` has a `default: return ... unknown profile %q`, so a
  hyphenated value is a typed error. It also fails closed on an unknown entity and
  a missing geometry, and it fills `Markers` from the entity's capitals.
- **`Generate()` fails closed too, indirectly.** It calls `validatePublicInput`,
  which cross-checks `Geometry.ID` against the id the profile references
  (`pipeline.go:62-73`); a hyphenated profile resolves `expected` to the UN
  geometry, so a de-facto geometry mismatches and errors.
- **The gap:** `pipeline.go:64` is a bare `if in.Profile == "de_facto"`. An `Input`
  built by hand with the hyphen *and* the UN geometry (or an empty `Geometry.ID`)
  passes the cross-check and silently renders the `un` posture.

**Therefore the CLI resolves geometry only through `geometry.InputFromCatalog` and
never constructs `geometry.Input.Profile` itself.** That makes the gap
unreachable. Add a test that a hyphenated `boundary` is rejected at the config
boundary rather than reaching geometry at all.

### CTR-6 carries both axes

ARCH-001 line 81 defines the manifest as "ISO id, files, **profile**, byte count,
dimensions, pin state and digests" — written before the two axes were
disentangled. One `profile` field cannot distinguish two assets generated at
different boundary postures. The T4 manifest carries **both** `profile` and
`boundary`. No governed decision is needed: ARCH-001 says compatibility is
additive within a major and only *removing* a field requires one.

## Known limitation to fix when it starts to matter

`requestedJSON` in `root.go` scans the raw arguments for a literal `--json` on the
error path, because a flag that fails to parse never reaches the parsed value.
Once T2 adds value-taking flags, `--config --json` (a missing value) or any flag
whose value is the string `--json` will flip error output to envelope mode. Fix it
when the first value-taking flag lands: stop scanning at the value position of a
known value-taking flag.

## Remaining tasks

- **T2 — `internal/config`.** Schema `country-map/v1` in YAML and JSON, embedded
  presets, single-parent inheritance, global tokens, profile and per-country
  overrides, documented precedence. Unknown fields, cycles, missing presets,
  invalid combinations and output collisions fail before anything is written.
  `explain` emits resolved value plus origin layer. REQ-13's layout contract
  belongs here, not later — it is the same surface freeze, and `card`/`hero` must
  come out as overrideable presets rather than modes. Reserve `selection`
  (DEC-015). VAL-2, VAL-8.
- **T3 — `internal/render`.** Five styles through tokens, two delivery modes,
  CTR-005 hooks, markers, opt-in animation. VAL-3, VAL-4.
- **T4 — `generate`.** Staging, atomic publication, CTR-006 manifest. VAL-5,
  VAL-6.
- **T5 — `inspect`, `preview`,** the QAB-2 byte baseline for P4, and
  `WKI-6638BACD6E20`: measure a real single-SVG invocation through the built
  binary — P2 measured 707 ms of ladder load inside a 1.53 s cold process — then
  accept the cost or implement lazy per-geometry decode.

## What P3 inherits from P2

- `geometry.Generate()` serves the committed ladder. Selection is a lookup, not a
  fine-to-coarse walk.
- `geometry.ErrNoArtifact` / `IsNoArtifact` is DEC-009's typed absence and must
  stay a first-class CLI result. Falling back to source there emits the exact
  silhouette the oracle refused.
- **Do not add a file under `internal/geometry`.**
  `TestDiagnosticSourceInventoryIsExactAndBiting` is an exact source inventory and
  reddens on any unregistered file there. `internal/render` consumes the package;
  it does not extend it.
- Byte budgets 2200/7500 path and 2500/8000 file are frozen. No country literal
  or per-entity branch anywhere. Pure-Go offline runtime. Presentation-free
  geometry (INV-1) is precisely what makes P3's theming possible.

## The habit that found both P2 defects

Russia rendered as a blob and France as a speck. Neither was visible to any
metric — the oracle's IoU floor is 0.40 and Russia scored 0.92. Both were found
by rendering the output and looking at it. `internal/geometry/cmd/svgshowcase`
and `cmd/svgproof` already do this for the geometry layer; extend the habit to
the CLI's own output before calling a task done. A green suite is not a rendered
card.
