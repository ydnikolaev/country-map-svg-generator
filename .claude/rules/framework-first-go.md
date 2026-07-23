<!-- SYNCED FROM mate@v1.4.1 · mode:sync · DO NOT EDIT — change upstream then `mate pull` -->
<!-- cites: [profiles/go/framework-first.md] -->

# Framework-first (Go) — reach for the maintained primitive

Always-on in a Go consumer. When a Go CLI needs to ask a human something, reach for the ecosystem's canonical primitive over a hand-rolled equivalent:

- **Interactive prompts are `charm.land/huh` forms**, not a hand-rolled `bufio`/`fmt.Scanln` loop and not `AlecAivazis/survey` (unmaintained) — huh owns single/multi-select, confirm, typed input with validation, and screen-reader accessibility; a hand-rolled loop re-solves raw-mode/resize/Ctrl-C/paste slightly wrong and is yours to maintain forever.
- **Gate the form on a TTY** (`term.IsTerminal`) *before* you build it — a non-interactive stdin takes the headless path or errors cleanly; a form must never block on input that will never arrive (the guarantee that keeps an agent's invocation from hanging). huh's accessible mode is for screen readers, **not** the headless path.
- **Interactive is the human default; headless is the agent opt-out** — expose a `--headless` (or flags/env) path that skips the form entirely, and let the form only *pre-fill* from flags. The form is a front end over the same values, never a second source of truth.
- **No interactive surface ⇒ no forms dependency** — a flags-only CLI takes neither huh nor a prompt loop; adding one "for consistency" is the custom-code-is-debt anti-pattern.

Full treatment: the Go profile (`.mate/profiles/go/`), which refines the CORE framework-first doctrine; the concrete library version, the TTY-detection primitive, and the headless opt-out's spelling are project values named there.
