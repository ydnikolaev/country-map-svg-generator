<!-- SYNCED FROM mate@v1.4.1 · mode:sync · DO NOT EDIT — change upstream then `mate pull` -->
<!-- SSOT SOURCE (mate repo). Consumers receive a provenance-stamped copy via `mate pull` — edit HERE, never the synced copies. -->

# Framework-first — Go profile

> **Thesis.** Interactive terminal UI is a solved concern with a maintained
> canonical library, so on the Go stack "framework-first" (the CORE doctrine)
> means: when a CLI needs to ask a human something — a choice, a multi-select, a
> confirm, a typed value — reach for the fleet's chosen forms library, never a
> hand-rolled `bufio` prompt loop and never the unmaintained option. This profile
> names that one canon today; it is the Go home the CORE doctrine's §5 seam points
> at, and it grows as more Go concerns earn a named primitive.

## 0. When to apply (and when not)

Apply to any consumer whose facet vector carries `lang: go` and whose CLI has an
**interactive surface** — a command a person runs and answers. The principle is
the Go instantiation of the CORE doctrine's §1 (search the canon before you
invent): terminal forms are canon, not domain logic.

Do **not** stretch it to a CLI with **no** interactive surface — most CLIs are
pure flags/args/stdin and need no forms library at all; taking the dependency
"to be consistent" is exactly the debt the CORE doctrine §3 warns against. And do
**not** read this as "every prompt is interactive": a value an agent or CI passes
belongs on a flag or env var (the headless path), never behind a form. The library
is for the **human** path; the machine path stays non-interactive by construction.

## 1. Interactive prompts use the fleet's forms library, never a hand-rolled loop

**When a Go CLI needs interactive input, build it with `charm.land/huh` (Bubble
Tea–based forms) — not a hand-rolled `fmt.Scanln`/`bufio.Scanner` prompt loop, and
not `AlecAivazis/survey` (unmaintained as of 2026).** huh is the ecosystem's
canonical, maintained mechanism for terminal forms: single- and multi-select,
confirm, typed input with validation, and a first-class accessible mode for screen
readers — all the edge cases (raw-mode restore, resize, Ctrl-C, paste, width) a
hand-rolled loop re-solves slightly wrong. A prompt loop written by hand is the
custom-code-is-debt anti-pattern (CORE §3) made concrete: it is terminal-handling
you now own and must maintain, to re-implement a solved concern.

Two guard rails make this bite without over-reaching:

- **The library is TTY-only; gate it.** Detect a non-interactive stdin
  (`term.IsTerminal`) *before* constructing the form and take the headless path or
  error cleanly — a form must never block on input that will never arrive (the
  guarantee that keeps an agent's invocation from hanging). huh's accessible mode
  is for screen readers; it is **not** the headless path (it is still interactive).
- **Interactive is the default surface, headless is the opt-out.** A human runs the
  bare command and answers the form; an agent passes `--headless` (or the flags/env
  it reads) and the form is skipped entirely. The form only ever *pre-fills* from
  flags — it is a front end over the same values, never a second source of truth.
- **A field with a vocabulary is a huh select/multiselect, never an input.** The
  Go mechanic for the cli-ux doctrine's "select, don't type": a `huh.NewSelect`/
  `NewMultiSelect` reads its options from an enumerable registry (`registry/facets.yaml`
  is the reference case), so a new value appears in the picker with no code edit;
  `huh.NewInput` is reserved for the genuinely open (a name, a path). Resolve a value
  the tool can discover (a source path from an env var) before prompting, and hint the
  near-certain default — the *what to ask* is the cli-ux doctrine, huh is the *how*.

## 2. Instantiation seam

What a Go CLI project swaps into this frame — each is a *value*, and a value never
belongs in the principle above it:

- **The forms library import path and version** — the concrete `charm.land/huh/v2`
  (or the pinned major) the project depends on; the *name of the fleet's chosen
  library* is this profile's fill, the principle ("use it, don't hand-roll") is the
  CORE doctrine's.
- **The TTY-detection primitive** — the `IsTerminal` the project calls and the fd
  it checks; "gate the form on an interactive stdin" is the principle.
- **The headless opt-out's shape** — the flag/env a project exposes for its
  machine path (`--headless`, a `CI` probe, a `--yes`), and which fields it seeds
  from flags. The principle is "a machine path exists and never enters the form";
  its spelling is the project's.

## 3. Anti-patterns

| ❌ | ✅ |
|---|---|
| Hand-rolled `bufio.Scanner`/`fmt.Scanln` prompt loop for a menu or confirm | Build a `huh` form (Select/MultiSelect/Confirm/Input) |
| `AlecAivazis/survey` in new code | `charm.land/huh` — the maintained canon |
| Running the form unconditionally, then hanging on a piped stdin | Gate on `term.IsTerminal`; take the headless path or error cleanly first |
| Repurposing huh's accessible mode as the "non-interactive" path | Accessible mode is still interactive; headless = no form call at all |
| Taking the forms dependency in a flags-only CLI "for consistency" | No interactive surface → no forms library (CORE §3: custom/extra code is debt) |

## 4. Porting checklist

- [ ] Every interactive prompt is a `huh` field, not a hand-rolled read loop; no `survey` import remains.
- [ ] The form is gated on `term.IsTerminal` (or equivalent) so a non-TTY stdin errors cleanly and never blocks.
- [ ] A headless opt-out (flag/env) skips the form entirely and reads flags/env only — the agent/CI contract.
- [ ] Form fields pre-seed from any flags passed, so the interactive and headless paths agree on the same values.
- [ ] A flags-only CLI takes **no** forms dependency (the principle is scoped to an interactive surface).

---

## Cross-links

This profile **refines** the CORE **framework-first** doctrine
(`doctrine/code/framework-first.md` in the SSOT; `../../doctrine/code/framework-first.md`
relative to this file once pulled into a consumer's `.mate/`) — it is that doctrine's
Go instantiation of the §5 seam (the concrete "search the canon" answer for one Go
concern: terminal forms). §1's guard rails are the Go *mechanics* of the
**cli-ux** doctrine (`doctrine/code/cli-ux.md`) — that doctrine says *what* a good
interactive command asks (resolve-don't-ask, select-don't-type, both front doors); this
profile says the Go *how* (huh). It sits beside `profiles/go/validation.md` (the Go-CLI
validation instantiation); a sibling profile on another stack refines the same parent
for that stack's canon.
