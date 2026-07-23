<!-- MATE:mate-reflexes START -->
Always-on reflexes (mate-managed; full treatment in the mate doctrine corpus under `.mate/doctrine/`):

# Validation discipline — verify before you conclude, guard what you verified

Always-on. On every claim you build on and every invariant you touch:

- **Verify before you conclude.** A claim about a system you don't own — a library, a provider, a version, a flag — rests on that system's own current source or an empirical check, never on memory; verify at the point of assertion, before the conclusion is load-bearing.
- **Claim only the strength you checked.** "Proven / verified / confirmed" must match what you actually ran; name degradation and uncertainty instead of rounding up, and when correcting a wrong claim re-anchor on the source rather than swinging to the opposite one.
- **Guard what you verified once (the maniac loop).** A hand-verification you did once is a gate you haven't written yet — new entity, SSOT, boundary, or convention means propose the gate, don't wait to be asked.

Full treatment: the validation and verification-honesty doctrines (`.mate/doctrine/`).

# Framework-first — use the canon, custom code is debt

Always-on. When you reach for custom code — a plugin, wrapper, low-level flag, or hack — to make something behave the way a framework or library might already own:

- **Search the canon first, in order.** Framework docs → library docs → repo precedent → only then custom. "I didn't know it could do that" is the normal result of searching, not an embarrassment.
- **Configure at the highest-level knob that owns the concern.** Don't drop to a lower layer (raw runtime flag, bundler option) to force behavior the higher one governs — that coupling breaks silently on upgrade.
- **Treat custom code as debt to retire, not precedent to extend.** An existing wrapper never justifies the next one; removing custom code beats adding it.
- **Round-N stop.** After a few failed speculative fixes, stop guessing — switch to a minimal reproduction and the canonical doc read end-to-end.

Full treatment: the framework-first doctrine (`.mate/doctrine/`); the concrete search order, config layering, escape hatches, and Round-N threshold for a given stack live in its `.mate/profiles/<stack>/` refinement.

# Commit hygiene — one intent, this session's, in the project's convention

Always-on (you commit outside an explicit `/commit` too, so the reflex can't wait to be summoned). On every commit:

- **Stage only what this session touched.** Never `git add -A` / `.` / `*` — a blanket add sweeps in another terminal's, worktree's, or manual edit's files, mis-attributing work you can't vouch for. List the paths you changed, reconcile against `git status`, leave the rest.
- **Write the message in the project's convention.** `type(scope): subject` — imperative, bounded, machine-readable so changelog / semver / blame keep working; the body explains *why*, not what. The concrete scope, type, and trailer vocabulary is a project value — it lives in the project's own commit-convention rule, not here.
- **Keep it atomic.** One logical change per commit, so it reverts, bisects, and reviews as one decision; split a multi-intent tree at the intent boundary.

Full treatment: the commit-hygiene doctrine (`.mate/doctrine/`); the concrete scope, type, and trailer vocabulary lives in the project's own commit-convention rule.

# Harness discipline — whose surface is this?

Always-on. Before writing to any harness artifact (skill, rule, gate, doctrine, config) — including the ambient "keep this for next time":

- **Never edit a managed copy.** A file with a mate provenance stamp (or in `.mate/lock.json`) is read-only here — the fix is authored upstream in the mate SSOT, released, and pulled back; editing it locally forks the source of truth and the next pull eats it.
- **Classify before placing; project-only is the default.** A project value or project-specific artifact stays home (native name, own file, or `.mate/config.yaml` — the seam synced skills already read); only a change with a nameable second consumer promotes up (`mate promote` → `/mate-promote` classifies core/profile/operator — or refuses).
- **Principle into the shared body, value into the config.** A shared artifact never carries one project's paths, commands, thresholds, or ids — and a fix to a *shared* artifact never stays local.
- **A provider fact is read, written down, and branched on by surface — never recalled.** Designing anything that touches the agent runtime starts at *its current docs*, not memory; what you verify goes into the adapter's surface registry with the build you checked it against; and the code asks that surface ("has a rules home?"), never the provider's name.

Full treatment: the harness-layers doctrine (`.mate/doctrine/`).

# Framework-first (Go) — reach for the maintained primitive

Always-on in a Go consumer. When a Go CLI needs to ask a human something, reach for the ecosystem's canonical primitive over a hand-rolled equivalent:

- **Interactive prompts are `charm.land/huh` forms**, not a hand-rolled `bufio`/`fmt.Scanln` loop and not `AlecAivazis/survey` (unmaintained) — huh owns single/multi-select, confirm, typed input with validation, and screen-reader accessibility; a hand-rolled loop re-solves raw-mode/resize/Ctrl-C/paste slightly wrong and is yours to maintain forever.
- **Gate the form on a TTY** (`term.IsTerminal`) *before* you build it — a non-interactive stdin takes the headless path or errors cleanly; a form must never block on input that will never arrive (the guarantee that keeps an agent's invocation from hanging). huh's accessible mode is for screen readers, **not** the headless path.
- **Interactive is the human default; headless is the agent opt-out** — expose a `--headless` (or flags/env) path that skips the form entirely, and let the form only *pre-fill* from flags. The form is a front end over the same values, never a second source of truth.
- **No interactive surface ⇒ no forms dependency** — a flags-only CLI takes neither huh nor a prompt loop; adding one "for consistency" is the custom-code-is-debt anti-pattern.

Full treatment: the Go profile (`.mate/profiles/go/`), which refines the CORE framework-first doctrine; the concrete library version, the TTY-detection primitive, and the headless opt-out's spelling are project values named there.
<!-- MATE:mate-reflexes END -->
