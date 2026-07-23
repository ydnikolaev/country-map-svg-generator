---
name: implement
description: Drive one ready specification through grounded planning, one bounded implementation run, exact verification, independent audit, and lifecycle closure. Use for a single spec; use teamlead for multi-spec waves or parallel orchestration.
user-invocable: true
argument-hint: "<epic> <spec-id> [plan|execute|resume]"
requires_capabilities: [run-shell, read-files, edit-files]
cites: [code/validation.md, agent/verification-honesty.md, agent/commit-hygiene.md]
mate_synced: v1.4.1
---

# implement

Implement exactly one governed spec. Repository lifecycle state is the resume
authority; chat history is not. The skill owns coordination for this one spec,
but independent judgment stays independent.

Invocation mode defines mutation authority:

- `plan`: author and validate the generated plan candidate only; do not edit code,
  accept the plan, or mutate later lifecycle state. The caller accepts or rejects it.
- `execute`: work only inside an already active run and return the real diff plus
  focused evidence. When a teamlead issued the brief, do not commit, issue the
  project gate, finish the run, audit, or complete the spec; the teamlead owns those
  mutations.
- `resume` or no mode: standalone one-spec coordination; follow the complete
  protocol below and own lifecycle closure.

Use `teamlead` instead when the request spans several specs, needs a parallel
wave, or requires coordinating multiple implementers. Do not grow this skill
into a second general orchestrator.

## 1. Orient from machine state

Take `<epic>` and `<spec-id>` from the request. If either is missing and exactly
one actionable spec cannot be established from the repository, stop and ask.

Run:

```text
mate work snapshot <epic> --spec <spec-id> --json
mate spec status <epic> <spec-id> --json
mate work events <epic> --spec <spec-id> --json
mate agent brief status <epic> <spec-id> --json
mate agent result status <epic> <spec-id> --json
```

Then read the exact specification, accepted solution architecture, decisions,
amendments, and required context. Resolve the features home from
`.mate/config.yaml`; do not assume `docs/features`. Use semantic IDs in CLI
mutations, never repository paths or hand-authored receipt YAML.
Read the selected development-pipeline `protocol_version` from project config
and the installed source manifest; never infer it from release names, omit it as
"current", or guess an upgrade. Follow the compatibility error's recovery route
when the pair is unsupported.

Fail closed when:

- lifecycle validation is red, a dependency is incomplete, or an impact fence
  covers this spec;
- required context is missing, stale, superseded, or contradicted by code;
- the selected next action is `planned` rather than `available`;
- the spec cannot be completed in one bounded journey without weakening WHAT.

State routing:

| State | Action |
| --- | --- |
| `draft` | stop; discovery/readiness owns the missing contract |
| `ready`, no accepted plan | continue at planning |
| `ready`, accepted plan | continue at run preparation |
| `in_progress`, active run | resume the exact running brief |
| `in_progress`, passed run | continue at independent audit |
| `completed` | report the existing evidence; do not reimplement |

After orientation, `execute` requires `in_progress` with the exact active run and
jumps to section 4; it must not scaffold or accept another plan/run. `plan`
requires `ready` without an accepted plan and stops at the end of section 2.
Before standalone execution reaches section 3, the parent epic must be
`in_progress`: consume exact `epic.accept-readiness` and `epic.begin` evidence as
needed. In delegated mode this remains teamlead-owned. Never request
`spec.begin` readiness beneath a draft, ready, completed, or cancelled epic; the
backend rejects it and rechecks the same invariant at run start and completion.

Delegated `plan` and `execute` are worker modes: an issued brief is their complete
authority and the caller records their structured terminal result. Standalone
coordination owns brief issuance/result recording, but it still needs a fresh
principal for an independent audit; the same inline principal may not relabel
self-review as audit.

## 2. Ground and plan

Before authoring HOW, verify the spec's paths, APIs, schemas, commands, fixtures,
and test topology against current code. Treat prose as intent and the repository
as ground truth. Record a contradiction; do not silently reinterpret it.

Create the plan only through the lifecycle producer. In delegated `plan` mode,
the caller first issues an immutable implementation-planner brief; consume that
exact brief and the installed role contract instead of reconstructing a prompt.
In standalone mode, issue the same brief before invoking a provider worker:

```text
mate plan scaffold <epic> <spec-id> <PLAN-NNN> --actor <principal> --json
mate agent brief issue <epic> <spec-id> <BRIEF-NNN> \
  --role implementation-planner --provider <provider> \
  --binding-mode <native|conforming-emulation|unavailable> \
  --role-contract <installed-role-path> \
  --context specification=<spec#sha256> --objective <bounded-objective> \
  --validation <command> --expect-revision <revision> --json
```

Fill the generated spec-local plan and manifest. The plan must cover every
applicable `REQ-*`, `INV-*`, `AC-*`, and `VAL-*`, name exact writes and owners,
order tasks, identify the due focused/component/project gates, and state the
deviation boundary. Keep run-scoped verification evidence as pathless logical
slots; the CLI materializes their exact locations for each run. It may choose HOW
but cannot narrow WHAT.

Validate and seal it:

```text
mate plan validate <epic> <spec-id> <PLAN-NNN> --json
mate plan seal <epic> <spec-id> <PLAN-NNN> --actor <principal> --json
mate readiness issue <epic> <spec-id> --for plan.accept --plan <PLAN-NNN> \
  --provider <adapter-provider> --session-id <session> --json -- <due preflight argv>
mate plan accept <epic> <spec-id> <PLAN-NNN> --expect-revision <revision> \
  --readiness <REV-ID#sha256> --agent-result <RESULT-NNN#sha256> \
  --actor <principal> --json
```

Use deterministic validation first. For ordinary in-envelope work, self-review
the plan once. Require one fresh independent plan review only when risk policy
classifies a public contract, schema, registry, lifecycle guard, migration,
security/data boundary, or consequential architecture change. One remediation
is the limit; repeated structural disagreement stops for an owner decision.

In `plan` mode, seal the validated candidate, then return the candidate and
manifest handles, validation result, ground-truth contradictions, risk
classification, and unresolved decisions here. The caller records the terminal
planner result against the issued brief and owns readiness/acceptance. Do not
issue readiness, accept the plan, or continue into implementation. In standalone
mode, record `--outcome passed --output plan-manifest=<PLAN-NNN#sha256>` and pass
that exact result to `plan accept`; never accept from the prose summary alone.

## 3. Prepare and start one run

Ensure every required implementation source is current in the context registry,
then compile the exact run pack. Add or supersede sources through `mate context`;
never edit the registry directly.

```text
mate context list <epic> --json
mate run scaffold <epic> <spec-id> <RUN-NNN> --actor <principal> --json
mate context compile <epic> <spec-id> <RUN-NNN> <CTX-ID> --actor <principal> --json
mate context verify <epic> <CTX-ID> --json
mate run validate <epic> <spec-id> <RUN-NNN> --json
mate readiness issue <epic> <spec-id> --for spec.begin --run <RUN-NNN> \
  --provider <adapter-provider> --session-id <session> --json -- <due preflight argv>
mate run start <epic> <spec-id> <RUN-NNN> --expect-revision <revision> \
  --context <CTX-ID#sha256> --readiness <REV-ID#sha256> --actor <principal> --json
```

The scaffold must exist before context compilation because the compiler binds
the pack to that run target; an untouched scaffold correctly reports
`complete=false`. Compile and verify the pack, then fill the generated brief
before `run validate`. Its inputs are exactly the accepted plan and current
context-pack IDs. Validation resolves both bindings and fails before readiness
can execute its gate. Readiness is bound separately by the atomic `run start`
transaction. Preserve the returned brief handle and digest.
On resume, reconstruct the same handle from `work events`: the `run.start`
successor digest is the immutable running brief digest.

After `run start`, standalone coordination issues a coder brief with `--run
<RUN-NNN>`, exact plan/context/readiness/subject handles, manifest-derived write
roots, forbidden roots, and assigned validation. Delegated `execute` consumes the
already-issued brief. Follow its recorded `named-role` or
`generic-serial-installed-contract` route exactly; fallback reads the installed
contract named by the brief and never substitutes remembered prompt text.

If validation returns `accepted_plan_unexecutable`, no run exists to interrupt.
In standalone coordination, execute the returned exact plan-invalidation command
and replan through a successor; in delegated mode, stop and return that command
to the teamlead. The current invalid manifest is failure-only authority: never
edit accepted bytes or manufacture a run/audit to obtain another route.

## 4. Execute inside the accepted envelope

Implement in the order and ownership envelope of the accepted plan:

1. add or update the smallest failing test where the behavior is testable;
2. make the smallest coherent production change;
3. run focused checks early, then the registered component/stack depths;
4. update docs, validators, fixtures, generated projections, and migration bytes
   exactly where their obligations are due;
5. route useful non-blocking work through `mate followup route`; never expand the
   current spec because an adjacent improvement is attractive.

At the execution handoff (or, standalone, before finishing the run), capture
each grounded **non-blocking harness** discovery once with
`mate feedback add --scope project|mate --severity high|medium|low`: use
`project` when the remedy belongs only to this consumer and `mate` only when
the portable harness should change. Assign severity from verified impact, not
from how strongly the agent prefers the improvement.
Attach repository evidence and retain the returned finding handle in the run
handoff. Feedback capture must not expand or delay the accepted run, replace a
required follow-up/deviation, or grant authority to queue, decide, or promote.

Do not weaken a test, validator, AC, or gate to obtain green. Do not hand-edit a
generated projection. Use project commands and `make` targets; `make check` is
the authoritative ceiling, not the inner loop.

In teamlead-delegated `execute` mode, stop after focused verification and return
files changed, tests added, exact commands/output, deviations, skipped work, and
blockers. The coordinator records those fields with `mate agent result record
--run <RUN-NNN>` and performs sections 5–7. It derives changed paths from the real
checkout/worktree against the pre-dispatch revision and rejects disagreement with
either the worker report or brief allowlist. Standalone coordination performs the
same derivation before recording its result. Every result carries a bounded
`--summary`; native results also carry the provider child/session identity. Every
issued brief must reach one terminal result, including `blocked`, `cancelled`, or
`no-result`; otherwise the lifecycle fence remains intentionally closed.

### Deviation protocol

- Local HOW correction inside the declared footprint: record it and continue.
- Material HOW or footprint change: stop mutation, finish the exact run as
  `interrupt`, invalidate the accepted plan with that immutable interrupted
  run-result, scaffold a successor that binds the returned invalidation receipt,
  then seal/readiness-check and accept it through `plan supersede`. Only then
  scaffold and start a new run.
- WHAT, macro-HOW, dependency topology, public contract, security/data, or wide
  architecture change: interrupt, create/route the amendment, and return to
  bounded delta-discovery or `teamlead`.
- Unknown or unsafe side effects: stop immediately and report a blocked resume
  anchor. Never let an uncertain run authorize sibling work.

This is delta review, not a fresh operator interview, unless the original
product premise or decision authority actually changed.

## 5. Verify and finish the run

Before a passing finish, use the project's commit policy to produce a clean,
reviewable implementation revision. Then run the exact due ceiling through the
production evidence producer:

```text
mate gate issue <epic> --spec <spec-id> --run <RUN-NNN> \
  --subject run-brief=<RUN-NNN#sha256> --gate <gate-id> --actor <principal> \
  --json -- make check
mate run finish <epic> <spec-id> <RUN-NNN> pass --expect-revision <revision> \
  --gate <GATE-ID#sha256> \
  --agent-result agent-dispatch-result=<RESULT-NNN#sha256> \
  --actor <principal> --json
```

Use `fail`, `interrupt`, or `cancel` honestly when applicable; only `pass` accepts
a gate. Never fabricate selection, attempt, result, gate, run-result, or
transition YAML. A red gate leaves the run active for correction; an
infrastructure failure is reported separately from a product failure.

## 6. Independent audit and spec closure

After a passed run, scaffold the audit, then issue a run-scoped auditor brief to
a fresh review-only principal with the spec, accepted plan, amendments, exact run
result, gate evidence, changed code, and installed auditor contract.
The implementer may supply evidence and remediate findings but cannot declare
its own blocking findings closed.

Create the audit frame mechanically:

```text
mate audit scaffold <epic> <spec-id> <AUD-*> \
  --run-result <RUN-RESULT-ID#sha256> --actor <auditor> --json
```

The command returns a mutable `locator`. The independent auditor fills every
semantic section and emits actionable findings or an explicit closure verdict.
Then:

```text
mate audit validate <epic> <spec-id> <AUD-*> --json
mate audit seal <epic> <spec-id> <AUD-*> --expect-sha <candidate-sha256> \
  --actor <auditor> --json
mate agent result record <epic> <spec-id> <RESULT-NNN> --run <RUN-NNN> \
  --brief <BRIEF-NNN#sha256> --principal <auditor> --outcome passed \
  --summary <bounded-summary> --provider-task <native-child-or-session-id> \
  --output audit=<AUD-ID#sha256> --expect-revision <revision> --json
mate audit verdict <epic> <spec-id> --audit <AUD-ID#sha256> --verdict pass \
  --agent-result <RESULT-NNN#sha256> --expect-revision <revision> \
  --actor <auditor> --json
mate gate issue <epic> --spec <spec-id> --subject audit=<AUD-ID#sha256> \
  --gate <gate-id> --actor <principal> --json -- make check
mate spec transition <epic> <spec-id> complete --expect-revision <revision> \
  --input run-result=<RUN-RESULT-ID#sha256> --input audit=<AUD-ID#sha256> \
  --input gate-receipt=<GATE-ID#sha256> --actor <principal> --json
```

A failed audit is immutable evidence. Record `--verdict adjust` with each
material `--finding <ID>=<severity>` and one explicit `--remediation-route`; this
atomically registers the findings and returns the immutable advisor authority for
that route. Consume that exact handle for same-plan remediation, plan
invalidation, or amendment routing; never edit, replace, or fabricate it. After
the authorized successor run, create a new audit and record its independent PASS
verdict, which closes the open remediation findings. Never edit a failed audit
into green.

At completion, provide only the terminal handles shown above. The closure CLI
resolves and binds the full transitive proof chain, including failed-audit
authority and the successor run-start inputs; agents must not discover paths or
fabricate extra proof references. Commit closure artifacts only after the final
project gate is green.

## 7. Handoff

Report only repository-verifiable facts: spec state and revision, accepted plan,
runs and outcomes, implementation commit(s), gate and audit verdicts, amendments
or routed follow-ups, and the next machine-selected action. A fresh session must
be able to resume from `work snapshot` and `work events` without this chat.
