# Task 2.1 Implementation Plan

## Scope

Implement only Task 2.1 and Requirements 1.1–1.6: authenticated human and Integration Principal context, server-derived effective Tenant context, supplied-Tenant comparison, Tenant-scoped resource lookup semantics, non-disclosing denial, and mutation-free boundary failure. Tenant rollout, production payroll, production operations, release enablement, broad role policy, later storage channels, and later domain workflows remain out of scope and gated.

## Authorization evidence

- Evidence ID: `T21-AUTH-001`.
- Decision/effective date: `2026-07-17`.
- Accountable approver: the user acting as project authority.
- Delegation: the user explicitly authorized the acting Senior Staff Engineer/Architect to select and record a coherent mainstream production-grade exact-version minimum stack for Task 2.1 after confirming readiness to proceed.
- Authorized scope: Task 2.1 and Requirements 1.1–1.6 only.
- Explicit exclusions: Tenant rollout, payroll production, production operations, release activation, and every later-task technology/policy decision.

## Approved execution sequence

1. Record the user-authorized minimum exact-version stack and open `G-IMPLEMENTATION` only for Task 2.1 dependencies.
2. Establish a Go module, production configuration validation, PostgreSQL migration, and Principal persistence contract.
3. Implement verified bearer-token authentication, server-side Principal resolution, effective Tenant context, supplied-Tenant comparison, and reusable Tenant-scoped resource access/service guards.
4. Expose the authenticated context endpoint and reusable HTTP middleware/error contract required by this task, with safe non-disclosing outcomes.
5. Add unit, PostgreSQL integration, HTTP integration, security, and Property 1 tests with at least 100 generated cases.
6. Run formatting, static analysis, race-enabled tests, integration tests, and self-review; update traceability with exact symbols/results; commit only Task 2.1 artifacts while leaving task-status transitions to the orchestrator.

## Non-goals and gate boundaries

- No Tenant provisioning, rollout, support session, role catalog activation, payroll, production operations, PWA, notification, export, or release work.
- No production identity realm, Tenant, or deployment is activated by this implementation.
- No caller-controlled Tenant value is authoritative; cross-Tenant and absent-resource outcomes remain equivalent at the API boundary.
