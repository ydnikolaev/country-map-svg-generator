---
name: teamlead
description: >-
  Coordinate a governed epic or multi-spec delivery from repository state:
  select dependency-legal work, obtain and accept plans, dispatch bounded
  implementers, verify actual evidence, integrate, handle deviations, and
  continue without chat memory. Use as the primary operator mode for more than
  one spec or whenever orchestration, sequencing, or parallel ownership is
  required.
user-invocable: true
argument-hint: "<epic> [next|resume|close]"
requires_capabilities: [run-shell, read-files, edit-files]
allowed-tools: Bash Read Edit Write
cites: [code/validation.md, agent/verification-honesty.md, agent/commit-hygiene.md]
mate_synced: v1.4.1
---

# teamlead

Own coordination, not every task. The repository lifecycle is the control plane;
agent messages are untrusted reports. Delegate through the strongest available
provider-native mechanism, but preserve this protocol when only serial execution
or a different provider is available.

`implement` owns one spec's planning and execution protocol. This skill owns the
epic DAG, plan acceptance, dispatch authority, integration, evidence verification,
and selection of what happens next. Never create a second private tracker in chat,
a provider task list, or an unregistered plan file.

## 1. Reconstruct the run

Resolve `<epic>` from the request. If it is absent and the repository does not
identify exactly one active epic, stop and ask. Start or resume with:

```text
mate work snapshot <epic> --json
mate work events <epic> --json
mate epic status <epic> --json
mate agent brief status <epic> <spec-id> --json
mate agent result status <epic> <spec-id> --json
```

Validate every candidate spec at its current lifecycle stage. Read the epic goal,
accepted architecture, decisions, amendments, dependency outputs, context
registry, and exact spec only when the snapshot makes that spec relevant. Resolve
all homes through `.mate/config.yaml`; never infer physical paths from examples.
Read the development-pipeline `protocol_version` from project config and the
installed source manifest; never guess compatibility from a release or provider.
Follow the typed recovery route when that pair is unsupported.

Fail closed when the snapshot is inconsistent, an available action has no
production command, an impact fence is open, or required context contradicts the
current code. A fresh session must reach the same decision from the same snapshot
and events.

Before issuing any `spec.begin` readiness, ensure the parent epic is
`in_progress`. If it is `draft`, issue and consume exact `epic.accept-readiness`
evidence first; if it is `ready`, issue and consume exact `epic.begin` evidence.
Never dispatch executable child work under a draft/ready/terminal parent. The
backend rechecks this ordering at readiness, run start, and spec completion.

## 2. Select dependency-legal work

Build the candidate set from `snapshot.specs` and `snapshot.dependencies`:

1. Resume any active run before starting new work.
2. Continue a spec with a passed run into audit and closure before opening an
   unrelated spec.
3. A spec is dispatchable only when every incoming dependency is `completed` and
   its `next_action.availability` is `available`. A deferred or cancelled upstream
   blocks by default; proceed only when governed disposition/amendment evidence
   explicitly rewrites that dependency obligation.
4. Within the same topological frontier, prefer explicit operator priority, then
   the shortest critical-path blocker, then stable spec ID order.
5. Report completed, cancelled, and fenced specs; never redispatch them.

Treat `planned` as unavailable. Do not probe legality by mutating state: readiness
and transition commands are final guards, not a substitute for selection.

Parallel dispatch is permitted only for dependency-independent specs whose
accepted plans declare disjoint write ownership and compatible validation
resources. Otherwise serialize. Never run repository-wide gates concurrently.

## 3. Plan through a bounded implementer

For a ready spec without an accepted plan, scaffold the candidate, then publish
one durable brief for a reasoning-capable planning principal. Bind the exact
specification and curated planning context, the installed provider role contract,
the current tracker revision, validation commands, and an explicit no-code/no-
lifecycle-mutation boundary:

```text
mate agent brief issue <epic> <spec> <BRIEF-NNN> \
  --role implementation-planner --provider <provider> \
  --binding-mode <native|conforming-emulation|unavailable> \
  --role-contract <installed-role-path> \
  --context specification=<spec#sha256> --objective <bounded-objective> \
  --validation <command> --expect-revision <revision> --json
```

Follow the command's recorded `dispatch_route`. `named-role` invokes the
provider's installed named planner. `generic-serial-installed-contract` starts a
fresh generic worker with the immutable brief and reads the installed contract
from `role_contract.path`; never copy remembered prompt text into the dispatch.
When no fresh worker is available, serial inline planning is conforming emulation,
not an independence claim.

The planner authors the generated spec-local candidate and manifest. The teamlead
then checks the actual candidate, not the planner's summary:

- every `REQ-*`, `INV-*`, `AC-*`, and `VAL-*` is covered;
- writes, ownership, dependency order, gates, rollback, and deviation boundary are
  explicit;
- HOW fits current code and does not narrow WHAT;
- cross-spec interfaces and downstream assumptions remain coherent.

Derive write authority from the accepted manifest's actual `tasks[*].writes` and
shared outputs. A prose plan, planner summary, or execution brief cannot add a
path omitted from that machine allowlist.

Use one self-review for ordinary work. Add one fresh independent plan review only
for a public contract, schema/registry/guard, migration, security/data boundary,
or consequential architecture change. One remediation cycle is the limit; a
second structural disagreement becomes an owner decision, not audit recursion.

After the coordinator checks and seals the candidate, terminally record what the
planner actually produced, then consume that exact result at acceptance:

```text
mate agent result record <epic> <spec> <RESULT-NNN> \
  --brief <BRIEF-NNN#sha256> --principal <fresh-principal> --outcome passed \
  --summary <bounded-summary> --provider-task <native-child-or-session-id> \
  --output plan-manifest=<PLAN-NNN#sha256> --expect-revision <revision> --json
mate plan accept <epic> <spec> <PLAN-NNN> ... \
  --agent-result <RESULT-NNN#sha256> --expect-revision <revision> --json
```

Only the teamlead issues plan-accept readiness and accepts the exact sealed digest.
Use the production commands from `implement`; never hand-author receipts. An open
brief or a missing/wrong result is a repository fence, not a reason to fall back
to the chat transcript.

## 4. Dispatch one bounded run

Compile current context, scaffold the run, and bind readiness through the
production lifecycle. After `run start`, issue one run-scoped coder brief with
the exact accepted plan, run brief, context pack, readiness, and subject handles.
Its `--write` values come only from `tasks[*].writes` in the accepted manifest.
Then dispatch `implement <epic> <spec> execute` through the recorded native or
installed-contract route with:

- exact accepted plan, run brief, context-pack, and readiness handles;
- repo-relative allowlist derived from the actual accepted manifest and an
  explicit off-limits set;
- scoped tests and measurable completion conditions;
- registered commit and worktree policy;
- required structured return: files changed, tests, check output, deviations,
  skipped work, blockers, and any exact feedback finding handles.

In delegated execute mode the implementer edits only inside the accepted envelope,
runs focused checks, and returns control without committing, issuing the project
gate, finishing the run, auditing, or completing the spec. Those mutations remain
teamlead-owned. Before recording its terminal payload, derive the actual changed
path set from the real checkout/worktree against the pre-dispatch revision, compare
it with both the worker report and brief allowlist, and reject any mismatch. Record
that coordinator-observed set with `mate agent result record --run <RUN-NNN>`, plus
the bounded summary, exact outputs/evidence, deviations, feedback, and provider
child/session identity for a native route. A blocked worker records `--outcome
blocked --summary ... --blocked-reason ...`; it does not leave an open brief behind.
If the provider has no independent dispatch capability, execute the same bounded
protocol serially and record that no independence claim was made.

Choose shared checkout, worktree isolation, or serial execution from actual write
overlap and project policy. The coordinator alone integrates. Stage explicit
session-owned paths; never use `git add -A`.

## 5. Verify, integrate, and close

An agent result is a claim. Inspect the real diff, reconcile it with the allowlist
and accepted plan, rerun the due focused checks, and reject hidden scope, skipped
tests, weakened gates, or unreported generated output. Then create the thematic
commit under the project commit policy.

Keep run-scoped verification evidence pathless in the plan. Let the production
gate producer materialize exact evidence per run; never copy a prior run's
physical evidence paths into a remediation brief.

Run the registered project ceiling through `mate gate issue`, then finish the exact
run with the complete coder-result set:

```text
mate run finish <epic> <spec> <RUN-NNN> pass ... \
  --agent-result agent-dispatch-result=<RESULT-NNN#sha256>
```

After a pass, scaffold the audit mechanically and issue a run-scoped auditor brief
to a fresh review-only principal. Bind the exact run result as context/subject and
the installed auditor contract. The auditor fills the candidate; the coordinator
validates and seals it. Record the auditor result with `--output
audit=<AUD-NNN#sha256>`, then consume it with `mate audit verdict ...
--agent-result <RESULT-NNN#sha256> --actor <same-auditor-principal>`. Inline work
by the implementer cannot satisfy this boundary. If no fresh principal or human is
available, stop with the durable brief and blocked route rather than self-audit.

The implementer may remediate findings but cannot close its own blocking finding.
Run the closure gate and complete the spec only from exact handles as defined by
`implement`. A non-PASS audit verdict and its remediation route are immutable
advisor authority. Consume the returned handle exactly; do not rewrite the audit
or reconstruct its authority. At closure, supply terminal handles only and let
the CLI resolve and bind the transitive failed-audit, successor-run, readiness,
context, and agent-result proof chain.

Apply risk-proportional review:

| Risk | Required review |
| --- | --- |
| ordinary in-plan change | deterministic gates plus self-review |
| contract, guard, migration, security/data | one independent reviewer, one pass |
| architecture or cross-spec meaning | owner decision; at most one council on request |

The full project ceiling is due at run finish, spec closure, and before a release;
focused checks are the inner loop. Never replace the registered ceiling with a
cheaper command because an agent already reported green.

At the spec boundary, preserve implementer finding handles and capture any
additional grounded, non-blocking coordination discovery once with
`mate feedback add --scope project|mate --severity high|medium|low`. Choose
`project` for consumer-only machinery and `mate` only for a portable harness
change; derive severity from verified delivery impact. This capture is part of
the handoff, not a reason to widen the spec or start central triage; ordinary
delivery never runs `mate feedback queue` or `mate feedback decide` on its own
findings.

## 6. Handle deviations before they compound

- Accepted plan cannot pass run-binding compilation before any run starts:
  consume the exact `accepted_plan_unexecutable` next action returned by the
  CLI, invalidate only with that current `plan-manifest=PLAN-NNN#sha256`
  authority, and produce a reviewed successor. Never edit accepted bytes or
  fabricate a failed run/audit to escape the guard.
- Local HOW correction inside accepted ownership: record it in the run result and
  continue.
- Material HOW or footprint change: stop writes, finish the exact run as
  `interrupt`, invalidate its accepted plan with the interrupted run-result,
  scaffold a successor bound to the invalidation receipt, then seal/readiness-
  check and accept it through `plan supersede` before starting a new run.
- WHAT, macro-HOW, dependency topology, public contract, security/data, or broad
  architecture change: interrupt, create an amendment, fence affected specs, and
  run bounded delta-discovery/architecture review.
- Unknown impact: stop and preserve a blocked resume anchor.

After every accepted deviation, compare the actual diff with downstream specs,
decisions, and dependency assumptions. Update them through governed amendments;
never silently patch normative spec bodies. Recompute the legal frontier before
dispatching again. A delta review is not a full discovery interview unless product
intent or decision authority changed.

## 7. Continue and hand off

After each provider return or coordinator restart, reload the snapshot plus `mate
agent brief status` and `mate agent result status`; provider task history is never
needed to decide whether to redispatch, record a terminal result, or consume it.
After each terminal spec event, reload `mate work snapshot`; never choose the next
spec from memory. Continue until no legal spec remains, the requested scope is
complete, or a stop condition fires.

Stop for the operator when there is an architecture fork, missing authority,
conflicting write ownership, repeated failure of the same gate, exhausted review
budget, unavailable required capability, or no legal action despite unfinished
work. Preserve exact feedback handles and keep the current spec bounded.

At a wave boundary report repository-verifiable facts only: epic/spec states,
accepted plan and run handles, commits, gates, audit verdicts, amendments,
follow-ups, blocked reasons, and the next machine-selected frontier. The handoff is
complete only when another provider can resume from repository state without this
conversation.
