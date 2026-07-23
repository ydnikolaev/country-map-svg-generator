<!-- SYNCED FROM mate@v1.4.1 · mode:sync · DO NOT EDIT — change upstream then `mate pull` -->
<!-- SSOT SOURCE (mate repo). Consumers receive a provenance-stamped copy via `mate pull` — edit HERE, never the synced copies. -->

# Validation — Go CLI profile

> **Thesis.** A Go CLI's real contract is the **binary a stranger installs**, and
> Go ships the canonical machinery to test exactly that. So on this stack the CORE
> validation doctrine's *"validator tech = a test in the project's runner"* seam
> resolves to one shape: drive the **built artifact** (not the library) through its
> behaviors — breadth as many isolated `testscript` scenarios, depth as one
> sequential loop for cross-command state. Each principle below is Go-ecosystem
> canon (the standard `go test` runner, `os/exec`, `GOWORK`, and
> `rogpeppe/go-internal/testscript` — the harness Go itself, `cue`, and `goreleaser`
> use), so it holds on any Go CLI regardless of its domain.

## 0. When to apply (and when not)

Apply to any consumer whose facet vector carries `lang: go` **and** `kind: cli` —
a Go program whose surface is subcommands a user runs. The four principles are the
Go instantiation of the validation doctrine's validator-tech seam; they assume no
particular CLI framework, output format, or domain.

Do **not** stretch this profile to a Go *library* (there is no binary to
black-box — a library's e2e is its exported API exercised under `go test`), nor to
a long-running *service* (whose e2e is a running process and a client, a different
facet). This profile carries only what every Go **CLI** shares.

## 1. Black-box the built binary, not the library

**The e2e suite drives the compiled binary through `os/exec`; the unit suite
drives the library.** They answer different questions: units prove a package's
logic, e2e proves the artifact a stranger downloads actually works in a directory
that never heard of it. Build the binary once (a package `TestMain`), run it in a
fresh `t.TempDir()`, and assert on the three things a user sees — stdout, stderr,
exit code — plus the files it writes. A test that imports the library to "e2e" it
is a unit test wearing an e2e label: it can pass while the wired-up binary is
broken.

## 2. Build the artifact as a fresh checkout would — `GOWORK=off`

**Build the e2e binary with the Go workspace disabled, so a green e2e can never
rest on a `go.work` that a consumer's `go install` will not have.** A workspace
silently injects local module replacements; a stranger running `go install
<module>` gets none of them. Building under `GOWORK=off` reproduces that stranger's
checkout — a module graph that only resolves against published dependencies — so a
break that a workspace would paper over fails the test instead of shipping. This is
the make-impossible rung of the CORE prevention hierarchy applied to the build
environment: the one build that proves portability is the one with no local help.

## 3. Cover breadth with `testscript` isolated scenarios

**Use `rogpeppe/go-internal/testscript`: one `.txtar` script per behaviour, each in
its own hermetic work directory.** This is where the per-subcommand matrix lives —
one small script apiece for the happy path (output shape), argument validation, the
error paths, and each flag. Hand testscript the built binary on `PATH` and invoke
it explicitly (`exec <cli> …` for success, `! exec <cli> …` for an expected
non-zero exit); assert with `stdout`/`stderr`/`exists`/`cmp`, and stage fixtures as
txtar file sections so each script is self-contained and readable standalone.
Isolated scenarios are cheap to add and independently debuggable — the natural home
for the long tail of edge cases a single linear test would bury.

## 4. Keep one sequential loop for cross-command state

**One ordered Go e2e test drives a single working tree through the whole product
loop, sharing state command→command.** Per-scenario isolation is exactly wrong for
the bug where a command breaks the state its *successor* needs: the loop (install →
author → check → hand-off, whatever the product's spine is) catches it because the
steps run in order against one tree. It is also the home for **programmatic fixture
mutation** — bump a version, rename a contract field, marshal a verdict — that reads
worse as many pre-staged txtar variants. Profile-and-loop are complementary: the
isolated matrix for breadth, the one loop for state flow. Neither replaces the other.

## 5. Hermetic by construction

**No e2e run reaches the network or touches the developer's real machine state.**
testscript hands each script a clean allowlist env; the setup pins the CLI's
machine-state/home variable to a per-run temp directory (so a scenario can neither
read nor pollute the real one) and sets whatever switches mute a passive
update/version check. A network-bound subcommand is exercised **only at the
CLI-parser layer** — an unknown flag that fails before any remote call — never a
real round-trip; its live behaviour belongs to a unit test with a stubbed client.
A hermetic suite is the precondition for the CORE doctrine's determinism (§5): a
test that can reach the network is not deterministic.

## 6. Instantiation seam

What a Go CLI project swaps into this frame — everything here is a *value*, and a
value never belongs in the principles above it:

- **The binary's import path and name** — what `TestMain` builds and what scripts
  invoke; the project's, not the profile's.
- **The machine-state/home env var and the network-mute switches** — the CLI reads
  its own; the test setup pins the first to a per-run temp dir and sets the rest.
  Their *names* are the project's value; "pin state per-run, mute the network" is
  the principle.
- **The subcommand set the matrix must cover** — enumerate it from the project's own
  command-surface SSOT (the same list a docs page or a freshness gate reads), so
  "cover every command" is *provably* complete and a newly added subcommand shows up
  as a missing script, not a silent gap. Driving the checklist off a hand-kept list
  is how the matrix rots.
- **The product loop's spine** — the ordered sequence §4's single test walks; the
  domain names the steps, the profile only mandates that one ordered test exists.

## 7. Anti-patterns

| ❌ | ✅ |
|---|---|
| "e2e" test that imports the library and calls its functions | `os/exec` the built binary; assert stdout/stderr/exit + files written |
| Building the e2e binary under the ambient workspace | `GOWORK=off` build, so a broken fresh checkout fails instead of hiding |
| One giant linear test carrying every edge case | Breadth as isolated `testscript` scenarios; keep the linear test for state flow |
| Only isolated scenarios, no ordered loop | One sequential test for the command→command state a per-script sandbox can't see |
| An e2e that hits the network (a real update check, a live API) | Pin state to a temp dir, mute the check; exercise network verbs at the parser layer only |
| "cover everything" against a hand-maintained command list | Drive the matrix checklist off the command-surface SSOT; a new command is a visible gap |

## 8. Porting checklist

- [ ] The e2e package builds the binary once with `GOWORK=off` and drives it via `os/exec` — not the library.
- [ ] A `testscript` suite exists with one `.txtar` per behaviour: happy path, arg validation, error paths, flags — one family per subcommand.
- [ ] Exactly one ordered Go test walks the product loop against a single tree, sharing state command→command.
- [ ] Setup is hermetic: the machine-state/home var is pinned to a per-run temp dir, the network-mute switches are set, and no scenario reaches the network.
- [ ] Network-bound subcommands are covered only at the CLI-parser layer (unknown flag), never a live call.
- [ ] The scenario checklist is enumerated from the command-surface SSOT, so a new subcommand surfaces as a missing script.

---

## Cross-links

This profile **refines** the CORE **validation** doctrine
(`doctrine/code/validation.md` in the SSOT; `../../doctrine/code/validation.md`
relative to this file once pulled into a consumer's `.mate/`) — it is that
doctrine's Go-CLI instantiation of the *validator-tech* instantiation seam, and it
serves the **coverage / test-honesty** dimension of that doctrine's validatable-
dimension taxonomy (tests that actually exercise the artifact, not a proxy for it).
A sibling profile on another stack refines the same parent for that stack's test
canon.
