---
name: discover
description: >-
  Turn rough product intent into a resumable, grounded discovery checkpoint and
  a validated governed epic corpus: outcome, boundaries, macro architecture,
  self-contained specs, dependency DAG, validation allocation, and curated
  planning context. Use before teamlead when a feature or epic still needs
  interviewing, codebase grounding, decomposition, or readiness validation.
user-invocable: true
argument-hint: "<feature or epic idea>"
requires_capabilities: [run-shell, read-files, edit-files]
allowed-tools: Bash Read Edit Write
model_tier: reasoning
model: opus
cites: [code/validation.md, code/registries.md, code/structure.md, agent/verification-honesty.md]
mate_synced: v1.4.1
---

# discover

Convert rough intent into the governed inputs that `teamlead` can execute. Ask
one focused question at a time and persist accepted knowledge as immutable
discovery checkpoints. Define WHAT and cross-spec macro-HOW; never write a
task-level implementation plan or begin implementation.

The repository lifecycle is the control plane. The skill supplies judgment and
authors generated documents; `mate` commands own scaffolding, validation,
digests, state transitions, and machine authority. Provider-native read-only
fan-out is an optional acceleration. A serial path must always produce the same
semantic result.

## 1. Orient and choose the feature identity

Read `.mate/config.yaml` first. Resolve configured homes from it, then use the
feature registry's registered epic roots and tracker locators when authoring
generated files. Never assume `docs/features`, reconstruct a path from an
example, or search for the first same-named file.

Inspect repository state before asking the operator to repeat facts:

```text
mate epic list --json
mate epic status <epic> --json
mate discovery status <epic> --json
mate work snapshot <epic> --json
mate work events <epic> --json
```

Use only the commands applicable to an existing epic. Also inspect related
active and archived work, backlog items, accepted decisions, code and tests in
the likely footprint, external requests, and recent relevant changes. Ground
claims in current repository bytes. Report unavailable or contradictory truth;
do not replace it with an assumption.

If existing work overlaps the idea, show the overlap and ask whether to extend
it or create a separate epic. Once the operator chooses a new identity, create
it only through:

```text
mate epic new <epic> --title <title> --actor <principal> --json
```

## 2. Run the resumable interview

Ask exactly one unresolved product question at a time. Skip anything already
answered by the request, an accepted decision, or verified ground truth. Cover,
as applicable:

1. user-visible outcome and measurable success;
2. in-scope and out-of-scope boundaries;
3. actors, authority, and ownership;
4. data, API, UI, CLI, events, and public contracts;
5. macro architecture and cross-spec contracts;
6. real dependencies, external requests, and rollout constraints;
7. risk policy and irreversible choices;
8. testing depth, exact E2E boundary, and validator changes;
9. context each future spec needs for planning and execution.

After the first accepted answer and at least one current repository source,
start discovery with a strict JSON input. Use a short-lived input file; the
published checkpoint, not that file or chat history, is the resume authority.

```text
mate discovery start <epic> DISC-001 --input <input.json> \
  --actor <principal> --json
```

The input carries all accepted answers, digest-bound ground-truth sources,
accepted decision handles, open questions, and a precise resume anchor. Answers
use `owner`, `ground-truth`, or `inference` honestly. An inference never closes
a question that needs owner authority.

After material answers or grounding changes, reconstruct the complete next
input from the current checkpoint and publish a successor with the exact
revision:

```text
mate discovery checkpoint <epic> <DISC-NNN> --input <input.json> \
  --expect-revision <revision> --actor <principal> --json
```

If the session must stop, checkpoint it with `--event interrupt`. A fresh
session resumes from `mate discovery status` and the current immutable
checkpoint using `mate discovery resume`; it does not replay the full interview.
On CAS or unknown-success errors, inspect status and events before retrying with
the same operation identity.

## 3. Author the governed corpus

Create and edit authorable documents through production scaffolds. Resolve the
returned locator when one is provided; otherwise resolve the generated document
from the configured feature home and the feature registry's exact epic/spec
record. Never hand-edit trackers, registries, receipts, or checkpoint YAML.

```text
mate epic architecture scaffold <epic> --actor <principal> --json
mate spec add <epic> <spec-id> <slug> --title <title> \
  [--depends-on <upstream>] --actor <principal> --json
```

Author the epic and solution architecture with outcome, boundaries,
stakeholders, ground truth, macro architecture, cross-spec contracts, rollout,
risk, external ownership, and the specification map. Cut each spec as one
bounded implementation journey with complete user stories, requirements,
invariants, acceptance criteria, validation obligations, dependencies, and
non-goals. Do not prescribe task-level HOW.

The spec dependency graph is the only wave authority. Add an edge only for a
real data, interface, authority, or write-order dependency. Derive parallel
frontiers and critical path from that DAG, then record the projection and owner
in the architecture; never maintain a second editable wave tracker. In
`Delivery waves and ownership`, use canonical `W1..Wn` earliest-wave rows and
one `Critical path: P1 -> P3` line; when maximum paths tie, choose the
lexicographically smallest complete spec-ID sequence. In `Components and
ownership boundaries`, use a `Boundary | Owner | Scope` table with unique
`BND-NNN` identities, exactly one registered spec owner per row, at least one
row per spec, and non-overlapping repository-path scopes across owners.

Allocate unit, component, integration, E2E, and validator work at the lowest
boundary that proves the risk. Every validation obligation names an owner and
the lifecycle transition where it becomes due. Put E2E only where a cross-boundary
behavior cannot be proven more cheaply.

Record consequential choices as governed decisions. Discovery completion
requires exactly one accepted decision whose slug is
`risk-policy-<registered-policy>`; select a policy registered by the installed
pipeline rather than inventing one. Scaffold, author, validate, and accept the
decision through `mate decision` commands.

Register current planning knowledge only through context CAS operations:

```text
mate context list <epic> --json
mate context add <epic> <resource-id> --kind <kind> --title <title> \
  --source <registered-path> --owner <owner> --required-at planning \
  --expect-sha <registry-sha> --actor <principal> --json
```

Register exactly one current canonical solution architecture for `planning`.
Use `planning:<spec-id>` for spec-specific sources and plain `planning` only for
knowledge every spec genuinely needs. Supersede or stale sources through the
CLI; do not edit the context registry. Keep context curated and bounded rather
than accumulating every related document.

If the corpus needs a governed artifact for which the installed CLI exposes no
legal producer, stop with that exact capability gap. Do not fabricate machine
authority to make discovery green.

## 4. Complete only through the corpus gate

Before completion, ensure that no blocking product or authority question is
open. Completion cannot silently refresh knowledge: if a ground-truth source
was added, removed, changed, or was not `current` in the latest checkpoint,
publish and inspect one ordinary nonterminal checkpoint first. Request the
terminal checkpoint only when its ground-truth set, paths, and digests are
identical to that immediate predecessor:

```text
mate discovery checkpoint <epic> <DISC-NNN> --event complete \
  --input <input.json> --expect-revision <revision> \
  --actor <principal> --json
```

The command is the corpus authority. It must reject untouched slots, invalid
document identity, missing architecture, incomplete specs, DAG/wave mismatch,
unowned validation or external work, unknown risk policy, contradictory or
stale truth, and invalid planning context. Never weaken a document or remove a
blocking question merely to obtain green.

On refusal, preserve the active checkpoint and follow the named recovery route:
ground the missing fact, ask the one required owner question, repair the
authorable document, or use the proper producer. Re-run completion once. A
second structural contradiction is an architecture or owner decision, not an
invitation to enter a review loop.

## 5. Validate and hand off

After a completed checkpoint, run the production document and state checks:

```text
mate epic validate <epic> --stage draft --json
mate spec validate <epic> <spec-id> --stage draft --json
mate work snapshot <epic> --json
mate work events <epic> --json
```

Report the exact completed checkpoint handle, epic and spec identities, derived
parallel frontier and critical path, accepted risk decision, current context
registry digest, validation results, and any blocker. Recommend `teamlead
<epic>` as the next step. Do not accept readiness, create implementation plans,
start runs, or claim that a provider transcript is lifecycle evidence.

## Stop conditions

Stop and preserve a resume anchor when product authority is missing, current
code contradicts a premise, an architecture fork changes WHAT or macro-HOW, an
external owner is absent, a required producer/capability is unavailable, or the
same deterministic gate fails twice without new evidence. Never hide the stop
by guessing, flattening the DAG, batching questions, or moving work into chat.
