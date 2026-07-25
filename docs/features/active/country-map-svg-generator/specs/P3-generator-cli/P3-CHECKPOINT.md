# P3 checkpoint — resume point for a fresh session

Written at T1 rather than at the end, deliberately. P3 delivers ungoverned under
DEC-014, so there is no run result, no audit verdict and no plan manifest to
resume from — this file and the commits are the whole record. Read it first.

## The one-line status

**All five tasks done.** All seven commands are registered: `init`, `validate`,
`explain`, `generate`, `inspect`, `preview`, `version`. The binary generates the
whole catalog offline in 3.3 s, publishes transactionally, and carries a CTR-006
manifest. What remains is not P3 build work — it is the governance reconciliation
debt `#15` and the four findings below, all of which land in P2 or P4.

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

P3's budget is **≤ 90 s added to `make check`**. Final measurement, from one
green `make check` on an idle machine at the close of T5:

| Package | Time |
| --- | --- |
| `cmd/country-map-svg-generator` | 17.9 s |
| `internal/config` | 0.2 s |
| `internal/render` | 35.1 s |
| **P3 total** | **53.2 s** |
| `internal/catalog` | 31.1 s |
| `internal/geometry` | 376.8 s |
| `internal/geometry/cmd/lodbuild` | 296.0 s |

Almost all of `render` is VAL-8 and VAL-3 generating real geometry; the
serializer tests reuse one generated result across the whole style × delivery
matrix precisely so the matrix does not multiply pipeline runs. The levers that
keep P3 inside its budget:

- The ordered product loop (VAL-5) runs a **selected set** of ISO codes chosen for
  shape, never the catalog.
- Full-catalog double-generate is **P4's** E2E per ARCH-001's allocation table.
  P3 proves the contract; P4 proves it at scale.

Also: summarize package-level failures, not just test-level ones. A `go test
-json` package failure carries no `Test` field, so a test-only filter reports a
timed-out package as a clean run. That bit P2 twice.

**Run `make check` detached if the harness keeps killing it.** A backgrounded
invocation was stopped three times in one session — not by a test, by the
harness. `nohup make check > log 2>&1 < /dev/null & disown` survives, and then
the log is read rather than the exit code. Note that macOS has no `setsid`: an
attempt to use it produced a log containing only `nohup: setsid: No such file or
directory`, and the run never started at all. Read the log, not the status.

**And never run `make check` while anything else heavy is on the machine.** The
wall-clock brakes measure the machine, not just the code. A contended run showed
`geometry` at 440 s against 343 s quiet and `catalog` at 45 s against 28 s, and
the `lodbuild` micro-brake failed at 31.3 s against its 30 s limit. That is the
brake working, not flaking — but it means a red brake is only evidence about the
code when the run had the machine to itself. Check `pgrep -f "go test"` first;
`WKI-B05B4B4287A2` is the real fix.

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
| `init`, `validate`, `explain` | T2 — done |
| `generate` | T4 — done |
| `inspect`, `preview` | T5 — done |

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

## T2 — done

`internal/config` (BND-005) carries schema `country-map/v1`, the embedded
presets, the precedence chain and every diagnostic. `init`, `validate` and
`explain` are registered. `internal/render` (BND-006) exists with the one bridge
into geometry; its serializer is T3.

### Four things the tests found that review would not have

**1. `encoding/json` matches field names case-insensitively.**
`{"layout":{"longside":160}}` binds to `LongSide` and decodes cleanly, so
`DisallowUnknownFields` alone accepts it. For a schema frozen at v1 that would
make every casing of every key part of the contract forever, and would make a
genuine typo look like a working document. `internal/config/keys.go` derives the
exact key set from the struct tags by reflection and checks case exactly. The
same traversal answers "what keys exist", which is the discovery surface REQ-1
asks for, so one mechanism serves both.

**2. A mandatory `{iso}` in the filename template made REQ-3's collision check
unreachable.** That requirement was invented here, not specified; with it, two
entities can never resolve to one path, so the check the requirement asks for by
name could not fire. Removed. A template without `{iso}` is a legitimate
single-entity configuration, and for a batch the collision check refuses it with
a better diagnostic — one that names which entities actually collided.

**3. The specification's own example config did not validate.** Globals set
`mode: tight, longSide: 160`; the `US` override sets `mode: contain, width: 720,
height: 420`. Merged naively, the result carried a contain mode *and* an
inherited long side, and the resolved check refused the documented shape.
`pruneSupersededLayoutFields` drops sizing that belonged to a mode a later layer
replaced — and only sizing set **earlier** than the mode. A layer that sets a
mode and a foreign field together is contradicting itself, and the per-layer
check already refuses that; pruning it would silently accept the mistake.

**4. `init --out config.json` wrote a commented YAML body into a `.json` file.**
Caught by `init` validating what it generates, which is why it does that: a
starter config that does not validate teaches the wrong shape and blames the
author for it. JSON has no comments, so that path now emits plain JSON and
prints the guidance the YAML starter carries inline.

### The layering design, and why it is typed rather than generic

Every field is a pointer or a map so "unset" is distinguishable from "set to the
zero value". Without that, a layer could not tell whether an earlier one had an
opinion, `explain` could not name an origin, and an explicit `fillOpacity: 0` —
an invisible fill under a visible stroke, a real style — would be silently
dropped.

Layers are merged as typed settings rather than as generic maps, so every layer
has already been through key checking and validation by the time it merges. A
generic merge defers both to the end and reports a preset's typo against
whichever document inherited it.

Nested blocks merge field by field. A country override that sets one token must
not erase the rest, and the result of that bug would still look like a valid
config. Lists replace wholesale: appending would make an override unable to
shorten an inherited list, and would make the result depend on how many ancestors
happened to mention it.

`TestMergeCoversEveryField` walks the whole schema rather than a handful of
fields, because every hand-written merge assertion would still pass if the merge
silently skipped a field nobody thought to name.

### Resolution is two-pass, and it has to be

The profile block to apply depends on the resolved profile, and the resolved
profile can be set by the country override — `countries: {US: {profile: hero}}`
is the specification's own example. One pass would choose the block from the
document globals and apply the card block to a hero card. So the profile is
resolved first from every layer that can carry it except the profile blocks, and
only then is the matching block inserted. A profile block may not select a
profile; that is refused, which is what keeps this from being circular.

Each preset in an `extends` chain is named separately in the provenance
(`preset ancestry site-default`), because with a chain of three, "which preset set
this" is the actual question.

Presets carry only what they change. A preset restating a default makes `explain`
attribute the value to the preset, sending an author to edit the wrong place;
`TestPresetsCarryOnlyWhatTheyChange` enforces it.

### One decode path for both formats

goccy turns YAML into a generic value and `encoding/json` turns JSON into one;
from there both go through the same normalization, the same key check and the
same typed decode. Two typed decoders would be two implementations of "what does
this document mean", free to differ exactly where it matters — case sensitivity,
embedded-struct flattening, and the padding union type.

Normalization is where YAML's extra vocabulary is refused with a path: a
non-string mapping key, a `.nan`, a `.inf`. Unchecked, an infinity reaches the
layout as a dimension and a viewBox full of NaN is the first anyone hears of it.

goccy is the parser rather than a minimal one because `FormatError` carries the
line, column and a source excerpt, which is the difference between a fixable
diagnostic and "invalid YAML".

### The vocabulary, enforced rather than documented

`render.Vocabulary` builds the enums from their owners: profiles from
`geometry.Presets()`, boundaries from `geometry.AcceptedBoundaryProfiles`, ISO
codes and capital ids from the corpus. So a new detail preset appears in the CLI
with no edit to the config layer, and no country literal or enum lives here.

`geometry.AcceptedBoundaryProfiles` was added to the geometry adapter in this
task (an existing file — the source-inventory gate only demands registration for
**new** files), with `TestAcceptedBoundaryProfilesMatchesWhatTheAdapterTakes`
keeping the exported list and the switch that implements it from drifting.

`render.GeometryRequest` is the only place a `geometry.Input` is constructed, and
`TestTheTwoAxesAreNotTransposed` is its tooth. Its subtlest rule:
**a layout mode with no sizing must produce an empty geometry layout**, because
`ApplyPreset` fills a layout only when its mode is empty — returning a mode with
no size would suppress the preset's long side and produce a card sized by
nothing.

### The provenance constants moved again, as predicted

`goccy/go-yaml` landed with its import, `go.mod` and `go.sum` changed, and
`TestDiagnosticIdentityDriftChecksEveryRecordedLoadBearingIdentity` reddened for
the second time. Re-recorded on the same evidence. This is the gate functioning,
not debt: it fires once per real dependency change. `WKI-1DA58E0FE741` carries
the over-specification.

### VAL-8, and the framing defect it found

`internal/render/layout_test.go` asserts the layout contract against what
geometry **produces**, not against what the configuration was allowed to say.
Shapes are chosen from the corpus by measured aspect, never by ISO literal.

The no-distortion claim is made without reaching inside the pipeline: the same
entity rendered into a 900×300 frame and a 300×900 frame must produce drawn
extents of the *same* aspect. Under a stretch-to-fill bug they would follow the
frames (3.0 and 0.33), so the 2% tolerance separates quantization noise from
distortion by two orders of magnitude.

Two things worth knowing before extending it:

- **`NaturalAspect` is the full projection's aspect, not the drawn one.**
  DEC-013 fits the card to the silhouette it draws, so for an entity with
  excluded components those differ substantially. Assertions compare the frame
  against the measured drawn extent instead.
- **The centring tolerance is one output pixel, not the quantization grid.** The
  drawn extent is the simplified and quantized silhouette, and simplification
  pulls the extreme vertex inward by a different fraction at each edge —
  measured at up to 0.19 px. A real centring failure misses by hundreds.

**It found a live framing defect: `WKI-37F18A2AA6A5`.** An explicit `contain`
frame is fitted *before* visibility removals, so an entity that drops components
draws off-centre — up to **292 px off in a 300 px tall frame**, worst case NC;
FM draws 6 components and removes 14, leaving 3.96 px of slack above and 145.3
below. This is the France-as-a-speck class on the path DEC-013 did not reach: an
explicit frame changes the layout, the committed ladder verdict does not
transfer, and the request takes the DEC-005 source path where the old ordering
still applies. It is `internal/geometry` (BND-003), so it is P2's to fix and
currently fenced.

The centring assertion therefore holds over entities that draw everything they
were fitted for — **verified over 18 sampled entities**, and the test fails if
that population is ever empty rather than passing vacuously. The excluded
population is measured by `TestVisibilityRemovalsPushTheSilhouetteOffCentre`,
which reddens when the gap closes, so a fix cannot be silently outlived by a
green suite. That is P2's own idiom for its superseded characterizations.

### Fixed from T1

`requestedJSON` now skips the value position of known value-taking flags, so
`--config --json` (a missing value) no longer flips error output to envelope mode
while reporting a different mistake. `explain_json_error.txtar` covers it.

Config-shape failures are classified **by type**, not by matching substrings of
an error message. `config.DocumentError` wraps every "the document is wrong"
failure, so the CLI's frozen exit classes cannot be silently reclassified by a
dependency changing its wording.

### Follow-ups inside P3's own code

- **`preview` generates once per style.** `buildPreview` calls `Resolve` →
  `GeometryRequest` → `Generate` inside the style loop, so one entity costs five
  pipeline runs and the twelve-entity ceiling costs sixty. T3's serializer tests
  deliberately generate **once** and re-serialize across the whole matrix, for
  exactly this reason: a style is a serialization concern and nothing about it
  needs a second pipeline run. Hoist the generation out of the loop.
- **The e2e suite has no ordered cross-command loop of its own.** VAL-5's loop
  lives inside `generate_happy.txtar` (init → validate → generate, one tree, each
  step depending on the last), which satisfies it, but the Go profile also asks
  for one *Go* test driving the same spine so programmatic fixture mutation has a
  home. Worth adding when the first mutation is needed rather than before.

### Two things T4 should fix when it adds its first flag

- ~~`valueTakingFlags` is a hand-maintained list~~ — **done in T4.**
  `TestEveryValueTakingFlagIsRegistered` walks the tree in both directions.
- ~~`explain` and `generate` must resolve the preset the same way~~ — **done.**
  Both go through `geometry.ApplyPreset`, and `inspect` does too.

## T3 — done

`internal/render` serializes. Five styles, two delivery modes, CTR-005 hooks,
markers, opt-in animation, the accessibility modes AC-5 asks for, and the
structural gate REQ-9 needs. No command registered — `generate` is T4's.

### A style is data, and the serializer never learns its name

`internal/config/styles/v1.json` holds the five styles as **token defaults**, and
they enter the precedence chain as their own layer immediately above the embedded
defaults. That placement is what makes REQ-4's "resolve through tokens rather
than renderer forks" true rather than aspirational, and it is why the style
layer lives in `config` and not in `render`: with the defaults in the serializer,
`explain` could not name where a token came from, which is the one guarantee T2
was built around.

`TestEveryStyleSerializesThroughTokens` asserts the negative — the emitted markup
never contains a style name — because a serializer that grew a `switch` on style
would otherwise still pass. `TestStylesProduceDistinctOutput` keeps two styles
from collapsing into each other; a style that renders identically to another is a
name that means nothing and nothing else would notice.

**Consequence for the presets:** with styles carrying tokens, the shipped presets
shrank to what they actually change. `standalone-default` also stopped extending
`site-default` — it overrode everything its parent set, so the inheritance did
nothing. They are siblings now, and `presetChainIn` gained a seam so ancestry and
cycle detection are exercised against synthetic sets rather than against shipped
data that has no ancestry to walk.

`TestPresetsCarryOnlyWhatTheyChange` compares a preset against what it **selects
or inherits** — its style's tokens and its ancestry — and deliberately *not*
against the embedded defaults. A preset is a named public contract, so restating
a default pins it: without the restatement, changing a default later would
silently change what the preset means for every config that extends it.

### The two delivery modes are genuinely different assets

`standalone` bakes presentation and must need no host CSS (ARCH-INV-6). An
author's `var(--brand, rebeccapurple)` therefore collapses to its fallback, and a
`var(--brand)` with no fallback is **refused** with a diagnostic naming
`themed-inline` — emitting it would produce a file that renders as nothing, which
reads as a broken asset rather than as a configuration mistake. Standalone also
carries no `country-map` classes: hooks nothing can target are exactly the
redundant markup REQ-9 forbids.

`themed-inline` emits every paint as `var(--country-map-<token>, <resolved>)`, so
the host sets the property and the resolved value is the fallback — REQ-6's
"usable fallbacks, including `currentColor`".

Animation is refused with `standalone` at validation. It is hook-based by REQ-8,
so the host stylesheet owns the motion; a standalone asset has no host stylesheet,
which would mean motion nothing can deliver and — worse — nothing could switch
off for `prefers-reduced-motion`.

Accessibility defaults to `decorative` (`aria-hidden` plus `focusable="false"`,
because some browsers put an inline SVG in the tab order and that is a keyboard
trap on a decoration). The accepted composition is cards beside text that already
names the country; announcing each one makes the page worse. `labelled` is the
opt-in.

### REQ-9 made concrete

"No redundant markup" is the requirement nobody ever fails, so it was written
down as specific absences and asserted: no `<g>`, no `<metadata>`, no comments,
no `<desc>`, no default-valued attributes, no stroke attributes on a style with
no stroke, and no `fill-opacity` on a fill that is `none`.

**Two of those came from looking at real output rather than from the tests.** The
tight viewBox read `0 0 144 122.9026241596183` — thirteen decimals on a document
whose path is quantized to q=0.01, which is bytes in every asset for a distinction
no renderer can draw. `number()` now rounds to the grid the geometry already uses.
And `outline` emitted a `fill-opacity` for a fill of `none`.

`Structure` parses with `encoding/xml` rather than searching text: a search for
`<script` is defeated by any of the ways markup can spell the same thing, and a
parser sees the element regardless. Elements are an **allowlist** — a denylist is
wrong the first time either list changes.

VAL-4 mutates **real emitted output**, not a hand-written fixture: a gate that
only sees markup written to be caught proves it can catch that markup, not that it
guards what the tool produces. Fifteen mutations, each a way a defect or a hostile
value would actually get in. Its counterpart asserts the gate still accepts every
legitimate style × delivery combination — a gate nobody can satisfy gets deleted,
and then nothing is guarded.

### The escaper is written out, and why

`xml.EscapeText` is specified for *character data*. Its behaviour on the quote
characters is adjacent to what an attribute value needs rather than identical to
it, and "adjacent" is not a contract to rest markup safety on. `writeEscaped`
handles the five XML entities plus the three whitespace characters that must not
survive raw in an attribute, and no more.

It is **defence in depth, not the primary guard**: the config layer already
refuses markup characters and external references in every author-supplied token,
which is where a bad value gets a diagnostic naming the key.

## T4 — done

`generate` is registered. The batch is built and structurally validated **entirely
in memory** before anything touches the output directory, then published
atomically with backup-and-rollback. That ordering is INV-1: a batch that wrote
as it went would leave a half-regenerated catalog whose manifest agrees with it,
and nothing downstream could tell that from a complete one. The manifest is
staged and published last, so a consumer never reads an index describing files
that are not there yet.

DEC-009's typed no-artifact is a **recorded skip**, not a failure, and the human
output names the skipped entities rather than counting them — a number is easy to
read past, and a catalog quietly short by three is the failure this command is
shaped around.

### The finding that matters most in P3 so far

**`WKI-C35A01E965DC`: only the profile's own long side is served from the ladder.**
Greenland at profile `card` renders in 2400 bytes at `longSide: 128` and fails at
59892 bytes against the 2500 ceiling at `longSide: 160`. 128 is the card preset's
own long side; **160 is one of that same preset's documented reference sizes**
(`presets/v1.json` lists 90, 128, 160, 240) **and the value the P3 specification's
own example configuration uses**. Any other size falls back to the DEC-005
full-detail source path.

REQ-13 makes the long side a first-class knob and QAB-1 documents card use from
90 to 240 px, so most of the sizes an author is invited to ask for are sizes that
leave the ladder. A configuration that works on a small entity fails on a large
one with no warning at authoring time.

P3 mitigates, it does not fix — the ladder is `internal/geometry` (BND-003), P2
territory, fenced:

- `init` no longer writes a `longSide`, with the reason in the file.
- Geometry's own budget refusal is caught and reclassified. Its message —
  "linear path bytes 59892 exceed hard maximum 2500" — is true and says nothing
  about what to change; the CLI now names the profile's served size, explains the
  fallback, and offers the three real options.

The precise mechanism (effective scale outside every band's ceiling, or simply no
ladder row at that scale) was **not** traced. What is verified is the cliff and
its two endpoints.

### A second manifest defect found by looking

The manifest recorded `"height": "135.36669468232515"` while the asset's viewBox
said `135.37`. A consumer laying out a grid from the manifest would reserve a box
the file does not fill. `render.Number` is now exported and used for both, so the
manifest cannot quote a dimension the asset does not carry.

### `valueTakingFlags` is now gated

`TestEveryValueTakingFlagIsRegistered` walks the command tree in both directions.
It was the one maintained list in the package, and `--out` arriving in this task
is exactly the case that would have gone stale.

## T5 — done

`inspect` reports what the pipeline **did**, not what was asked for: tier,
viewBox, path bytes against the budget, points, components drawn *and removed*,
capitals, and geometry's own diagnostics. That distinction is where every
surprise in this epic has lived — Greenland draws 15 components and removes 112,
which `inspect` says and nothing else did. A typed no-artifact is reported as the
recorded decision it is. Nothing is written.

`preview` writes **one self-contained HTML page** with every style inlined. The
spec is explicit that it must not become a hidden web server: no localhost, no
port, no process left behind. The page inlines exactly what `generate` would
write, so what is looked at is the artifact rather than a rendering of it. A
catalog-wide selection is refused at twelve entities — a page nobody can read is
not a preview.

### DEC-016: two decisions, both on measurement

**The ladder's first-load cost is accepted; `WKI-6638BACD6E20` is closed.**
Through the built binary, median of five: `version` 0.596 s, `explain` 0.592 s,
`generate --iso FR` 1.288 s, whole catalog **3.3 s for 248 assets**. The
discovery commands never pay it — the ladder is loaded lazily by `Generate` and
neither generates. A batch amortizes it to nothing. On a single asset the load is
0.69 s of a 1.29 s run **against a 0.596 s process floor lazy decode cannot
touch**. Reopen if a workflow appears that invokes the binary once per asset.

**The first measured byte distribution does not meet QAB-2.** 248 assets, 1 typed
no-artifact (UM), 476 106 bytes total:

| | path bytes | file bytes | QAB-2 (card) |
| --- | --- | --- | --- |
| median | 1964 | 2165 | ≤ 600 |
| p95 | 2183 | 2384 | ≤ 1200 |
| max | 2197 | 2397 | ≤ 2500 ✓ |

The **maximum is satisfied**, and ARCH-001 says exceeding a maximum is what
blocks delivery. The median is 3.6× over, and the cause is structural: the ladder
selects the **finest** oracle-passing candidate per band, so finest-that-fits
puts every entity just under the cap by construction. One of the two has to
change; it is a page-weight against fidelity judgement, which DEC-008 and DEC-013
established as an owner approval over a rendered page. **`WKI-33B6EC2482B3`, due
P4.**

## Two gaps a spec re-read found after the suite was already green

A method note, not just a fix: every test passed and `make check` was green
before either was noticed. They were found by reading VAL-5 and VAL-6 word by
word against what the scripts actually did.

- **VAL-5 asks for `init → validate → generate → inspect`, and the loop stopped
  at `generate`.** The last step is the point of an ordered loop: a command that
  breaks the state its successor needs is invisible to isolated scenarios, which
  is why the Go profile asks for one sequential test alongside the matrix.
- **VAL-6 asks for a failure *during batch publication*, and every forced failure
  happened before it.** Configuration and render refusals never reach `Publish`,
  so the transactional guarantee — the reason the whole batch is built in memory
  first — had never once been exercised. The output directory is now made
  read-only mid-scenario, and the script asserts both that the previous output is
  byte-identical and that the diagnostic *says so*.

REQ-11's offline guarantee is stated at its verified strength. The load-bearing
argument is structural: testscript hands the script an empty work directory and
the loop succeeds with no corpus, preset or schema file on disk. The proxy
variables are a net rather than a proof — they catch a client that honours the
environment and would miss a raw dial.

## Where each acceptance criterion's evidence lives

Swept against the spec's own table so whoever closes P3 does not derive it again.

| AC | Evidence |
| --- | --- |
| AC-1 | Four testscript scenarios per registered leaf including the JSON error path, gated by `TestCommandSurfaceIsFullyCovered` in both directions |
| AC-2 | `site-default` is the portable preset; `TestSwitchingModeDiscardsTheOtherModesSizing` for a country override selecting either sizing mode with only declared dimensions changing; `TestContainNeverDistorts` for "geometry is never stretched" |
| AC-3 | `TestThemedInlineEmitsTheContractedHooks`; `TestAnimationIsOptInAndHookBased` asserts no `<animate>`, no `@keyframes`, no `<style>`, which is what leaves reduced motion for the host to honour. **The browser proof is P4's by the spec's own allocation**, and `preview` produces the page it loads |
| AC-4 | `generate_happy.txtar` generates twice and `cmp`s both asset and manifest, under proxies pointing at a closed port |
| AC-5 | VAL-4's fifteen mutations against real emitted output, plus `TestAccessibilityModesSerializePredictably` |

The one soft spot is AC-3's reduced-motion clause, and the spec assigns its final
proof to P4 rather than to P3.

## What P3 leaves behind

Five findings, none of them P3's to fix — all `internal/geometry` (P2, fenced) or
P4:

| Item | What |
| --- | --- |
| `WKI-4062B33B8FEA` | The amendment fence itself, with the mate fix sketched from the SSOT |
| `WKI-1DA58E0FE741` | The diagnostic identity record hashes the whole module graph |
| `WKI-37F18A2AA6A5` | A `contain` frame is fitted **before** visibility removals — up to 292 px off centre in a 300 px frame |
| `WKI-C35A01E965DC` | Only the profile's own long side is served from the ladder; any other size falls back to source and blows the ceiling |
| `WKI-33B6EC2482B3` | Finest-that-fits clusters every asset at the cap, 3.6× over QAB-2's median |

And the governance debt `#15`: everything in P3 shipped as ordinary engineering
commits under DEC-014, so **P2 and P3 both** need reconciling against real
receipts once mate grows a terminal edge out of `applied`.

### Follow-ups inside P3's own code

- **`preview` generates once per style.** `buildPreview` calls `Generate` inside
  the style loop, so one entity costs five pipeline runs and the twelve-entity
  ceiling costs sixty. The T3 serializer tests deliberately generate **once** and
  re-serialize across the matrix, for exactly this reason. Hoist it.
- **`valueTakingFlags`** is gated in both directions by
  `TestEveryValueTakingFlagIsRegistered` — done in T4, listed here because it is
  the one maintained list left in the package.

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
