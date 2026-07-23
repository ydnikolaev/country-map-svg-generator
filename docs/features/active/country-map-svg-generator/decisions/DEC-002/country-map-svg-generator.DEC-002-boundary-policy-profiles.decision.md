---
# SCAFFOLDED by mate from a versioned governed template; AUTHORABLE INSTANCE.
schema_version: "1.2.0"
template_version: "1.2.0"
kind: "decision"
id: "DEC-002"
epic: "country-map-svg-generator"
status: accepted
profiles: []
concerns: []
inputs: ["country-map-svg-generator"]
---
# DEC-002 — Explicit political boundary profiles

## Context

The generator serves a visa site across countries and territories. Any single
unlabelled world boundary dataset silently embeds a political viewpoint.
Natural Earth documents a de-facto default and separate claims, while the owner
requires a UN-aligned default plus an alternative de-facto view.

## Decision

The catalog identity is the current ISO 3166-1 set. `un` is the default boundary
profile and `de-facto` is a named alternative. Every generated manifest records
the profile and exact corpus version. Geometry overrides are explicit,
provenanced, schema-validated, and never inferred from locale or user location.

## Options considered

| Option | Fit |
| --- | --- |
| one unlabelled dataset | Smallest implementation but hides political policy |
| locale-selected boundaries | Flexible but creates implicit behavior and support risk |
| explicit named profiles | Deterministic, auditable, portable, and owner-approved |

## Rejected alternatives

An unlabelled dataset is rejected because consumers cannot identify its policy.
Locale-based switching is rejected because locale is not authority for a
political boundary decision and would make output nondeterministic.

## Consequences and tradeoffs

The corpus carries more than one geometry profile and documentation must state
that maps are decorative rather than legal cartography. Source updates require
profile-specific review. Consumers gain an explicit, stable switch without
forking generated assets.

## Affected artifacts and owners

P1 owns the ISO catalog, source provenance, and profile registry. P2 preserves
profile identity through processing. P3 exposes the profile in config, CLI, and
manifests. P4 documents client selection and verifies representative disputed
areas.

## Validation and revisit trigger

Completion requires profile completeness for every catalog entry, source/version
metadata, and negative tests for missing or implicit profile selection. Revisit
when ISO or either upstream geometry source changes an entity or disputed area,
or when the site owner selects a different default.

## Supersession

No predecessor. A successor must preserve explicit profile identity and cite the
owner authority for changing the default.
