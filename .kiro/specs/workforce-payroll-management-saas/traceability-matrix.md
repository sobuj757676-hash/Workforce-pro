# Requirements → Design → Tasks → Code → Tests Traceability Matrix

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**
>
> This artifact records specification traceability only. It does not approve a policy, implementation, Tenant rollout, payroll production use, or Singapore MOM, legal, contract, payroll, privacy, or compliance position.

## 1. Purpose, authority, and gate state

This matrix is the authoritative traceability artifact for Task 1.4. It links every acceptance-criteria set in [requirements.md](requirements.md) to controlling sections in [design.md](design.md), implementation and verification work in [tasks.md](tasks.md), current artifacts, and future Code → Tests evidence. Requirements remain authoritative, followed by design. A matrix entry cannot change either source.

The Task 1.2 disposition is preserved: qualified Singapore PTE, MOM, contract, payroll, legal, privacy, and compliance validation is an **External Dependency; Manual Compliance Validation Before Production; Not Required For Development**. No such approval is claimed here, and `G-PAYROLL-PRODUCTION` remains closed.

The Task 1.3 baseline in [access-governance-baseline.md](access-governance-baseline.md) remains **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**. It may guide sequential development only. `G-TENANT-ROLLOUT` remains closed. Every policy value must be persisted in versioned configuration tables, managed through the applicable Platform Admin Console or Tenant Management Admin Panel workflow, and never hardcoded. Production activation requires explicit business/security approval of the exact policy digest.

## 2. Traceability status and update contract

### 2.1 Status vocabulary

| Status | Meaning |
|---|---|
| `SPECIFIED` | Requirement, design, and planned task/test links exist; no implementation or passing-test claim is made. |
| `DRAFT_GOVERNANCE` | A review-ready governance artifact exists but remains **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**. |
| `IMPLEMENTED_UNVERIFIED` | Concrete code symbols and commit evidence exist, but required verification is incomplete or failing. |
| `VERIFIED` | Concrete code and tests exist, required test classes pass, and evidence is recorded. |
| `BLOCKED_BY_GATE` | Implementation or rollout cannot proceed because a decision-register gate is closed. |
| `EXTERNAL_VALIDATION_REQUIRED` | Qualified manual validation is required and cannot be replaced by code or automated tests. |

All Requirement 1–43 rows are currently `SPECIFIED` unless the evidence column explicitly records a draft governance artifact. No application code or automated-test completion is claimed by Task 1.4.

### 2.2 Mandatory Code → Tests update record

When an implementation task changes code or tests, update the applicable requirement row and add one evidence record in Section 5. Each record must contain all of these values; omission is a traceability failure:

1. acceptance criteria in `R<requirement>.<criterion>` form;
2. implementing task ID and verification task ID;
3. repository-relative production path plus symbol, route, schema, migration, configuration table, or UI component;
4. repository-relative test path plus stable test/property/state-machine identifier;
5. test class: unit, example, edge, integration, security, state-machine, accessibility, performance, property, or manual validation;
6. result and durable evidence: commit, CI run, report, or Qualified_Stakeholder evidence reference;
7. policy/configuration impact, including Admin Panel surface and configuration-table names where policy values are involved;
8. gate impact and confirmation that no closed gate was treated as open.

Use `**Validates: Requirements X.Y**` in property-based tests. Property tests must also identify feature name, design property number and statement, run at least 100 generated iterations unless an approved cost rationale is recorded, shrink failures, and retain minimized counterexamples as regression fixtures. A task checkbox or feature flag is not test evidence or policy approval.

### 2.3 Change-control rules

- Update the matrix in the same commit as the code/tests that change coverage.
- Never mark `VERIFIED` from task state alone; record executable evidence.
- A changed acceptance criterion, design invariant, task scope, implementation symbol, test identifier, or policy version requires impact review of all linked rows.
- Policy-bearing implementation must resolve versioned configuration from authoritative tables and expose authorized Admin Panel management. Missing, invalid, ambiguous, unapproved, or inactive policy configuration fails closed; application defaults cannot authorize.
- `G-TENANT-ROLLOUT` and `G-PAYROLL-PRODUCTION` remain closed until the [implementation decision register](implementation-decision-register.md) contains the required evidence. This matrix cannot open either gate.
- Post-final correction, off-cycle payroll, and revised Payslip workflows remain out of scope.

## 3. Complete Requirements → Design → Task → Verification matrix

Design references use the numbered headings in [design.md](design.md). Task references use exact IDs in [tasks.md](tasks.md). “Planned” means no code/test evidence is claimed yet.

### 3.1 Platform, tenancy, access, and shared command foundations

| Requirement and acceptance criteria | Controlling design | Implementation tasks | Verification tasks / properties | Current evidence and status |
|---|---|---|---|---|
| [R1 Tenant Isolation and Authorization Boundary](requirements.md#requirement-1-tenant-isolation-and-authorization-boundary), AC 1.1–1.10 | §§4.4, 5.3, 7.3, 20.1–20.2; Property 1 | 2.1–2.4, 3.2, 25.7 | 3.5, 29.1 (P1/P21), 30.5 | `SPECIFIED`; planned, no code/test evidence. |
| [R2 Platform and Tenant Administration Separation](requirements.md#requirement-2-platform-and-tenant-administration-separation), AC 2.1–2.9 | §§4.2–4.5, 5.1, 15.2, 20.2; Properties 19/21/27 | 3.2, 4.1, 6.1, 23.3 | 3.5, 29.1, 29.7, 30.5 | `DRAFT_GOVERNANCE`: access baseline §§2, 4–6, 12–14; application evidence planned. |
| [R3 Tenant Lifecycle and Provisioning](requirements.md#requirement-3-tenant-lifecycle-and-provisioning), AC 3.1–3.11 | §§5.2, 7.2, 9.5, 21, 23 Property 26 | 4.2–4.4 | 4.5, 29.7 (P26/P27), 30.5, 30.10 | `SPECIFIED`; commercial policy remains blocked by register `DEP-04`. |
| [R4 Tenant-Approved Support and Break-Glass Access](requirements.md#requirement-4-tenant-approved-support-and-break-glass-access), AC 4.1–4.13 | §§4.3–4.5, 5.1, 9.5, 20.2–20.4; Properties 19/20 | 5.1–5.4 | 5.5, 29.1 (P19/P20), 30.5 | `DRAFT_GOVERNANCE`: access baseline §§2.2, 4.2, 7–12; rollout remains closed. |
| [R5 Organization, Site, Account, Worker, and Assignment Administration](requirements.md#requirement-5-organization-site-account-worker-and-assignment-administration), AC 5.1–5.12 | §§4.1, 6, 7.2–7.3, 10, 13.1 | 6.2–6.4 | 6.5 | `SPECIFIED`; planned, no code/test evidence. |
| [R6 Permission-Based Roles, Maker-Checker, and Segregation of Duties](requirements.md#requirement-6-permission-based-roles-maker-checker-and-segregation-of-duties), AC 6.1–6.11 | §§4.1, 4.4–4.5, 10, 20, 23 Property 16 | 3.1–3.4 | 3.5, 29.6 (P16) | `DRAFT_GOVERNANCE`: complete Task 1.3 draft baseline; every value configuration/Admin Panel driven; not approved. |
| [R7 Versioning, Concurrency, and Idempotent Commands](requirements.md#requirement-7-versioning-concurrency-and-idempotent-commands), AC 7.1–7.10 | §§7.2–7.3, 10.1, 13, 17.2–17.3, 19; Properties 5/11/26 | 2.3–2.5 | 6.5, 8.7, 17.7, 19.7, 29.3, 29.5, 29.8–29.9 | `SPECIFIED`; planned, no code/test evidence. |

### 3.2 Calendar, attendance, overtime, consent, and evidence

| Requirement and acceptance criteria | Controlling design | Implementation tasks | Verification tasks / properties | Current evidence and status |
|---|---|---|---|---|
| [R8 Monthly Operating Calendar](requirements.md#requirement-8-monthly-operating-calendar), AC 8.1–8.12 | §§3, 7, 8.1, 10, 13, 19; Properties 2/3 | 7.1–7.3 | 7.6, 29.2, 29.9, 30.1 | `SPECIFIED`; planned, no code/test evidence. |
| [R9 Closed-Day Hard Gate](requirements.md#requirement-9-closed-day-hard-gate), AC 9.1–9.8 | §§3, 8.1, 9.3, 19; Properties 2/3 | 7.4–7.5 | 7.6, 29.2 (P2/P3), 30.2 | `SPECIFIED`; planned, no code/test evidence. |
| [R10 Supervisor-Created Toolbox Meetings](requirements.md#requirement-10-supervisor-created-toolbox-meetings), AC 10.1–10.9 | §§8.2, 9.1, 9.3, 10, 19 | 8.1–8.3 | 8.7, 30.1–30.2 | `SPECIFIED`; planned, no code/test evidence. |
| [R11 Dual Attendance Confirmation](requirements.md#requirement-11-dual-attendance-confirmation), AC 11.1–11.10 | §§7, 8.2, 9.1, 17, 23 Property 5 | 8.4–8.6 | 8.7, 29.3 (P5), 30.1 | `SPECIFIED`; planned, no code/test evidence. |
| [R12 Working-Day Planned Overtime](requirements.md#requirement-12-working-day-planned-overtime), AC 12.1–12.10 | §§8.3, 9.2, 10, 19 | 9.1–9.3, 9.5 | 9.6, 30.2 | `SPECIFIED`; planned, no code/test evidence. |
| [R13 Closed-Day Named and Bounded Authorization](requirements.md#requirement-13-closed-day-named-and-bounded-authorization), AC 13.1–13.10 | §§8.3, 9.3, 19; Properties 3/4 | 9.1–9.4 | 9.6, 29.2 (P3/P4), 30.2 | `SPECIFIED`; planned, no code/test evidence. |
| [R14 Date-Specific Worker Pre-Work Consent](requirements.md#requirement-14-date-specific-worker-pre-work-consent), AC 14.1–14.10 | §§8.3, 9.2–9.3, 15.5, 17; Property 4 | 10.1–10.5 | 10.6, 29.2 (P4), 30.2 | `SPECIFIED`; legal meaning remains unassigned; labels/subtypes must be versioned configuration. |
| [R15 Replay-Resistant Site QR Checkout](requirements.md#requirement-15-replay-resistant-site-qr-checkout), AC 15.1–15.13 | §§7.2–7.3, 9.4, 18–20; Property 6 | 11.1–11.5 | 11.6, 25.7, 29.3 (P6), 30.1 | `SPECIFIED`; QR policy values remain rollout prerequisites/configuration. |
| [R16 Checkout and Work Evidence Exceptions](requirements.md#requirement-16-checkout-and-work-evidence-exceptions), AC 16.1–16.10 | §§18–19, 24.4–24.5 | 12.1–12.4 | 12.5, 30.2 | `SPECIFIED`; planned, no code/test evidence. |

### 3.3 Work records and worker acknowledgement

| Requirement and acceptance criteria | Controlling design | Implementation tasks | Verification tasks / properties | Current evidence and status |
|---|---|---|---|---|
| [R17 Supervisor Work Record Review and Editing](requirements.md#requirement-17-supervisor-work-record-review-and-editing), AC 17.1–17.13 | §§7.2–7.3, 8.2, 9.4, 13.2, 19; Property 9 | 13.1–13.6 | 13.7, 29.4 (P9), 30.3 | `SPECIFIED`; planned, no code/test evidence. |
| [R18 Single Worker Acknowledgement and Dispute](requirements.md#requirement-18-single-worker-acknowledgement-and-dispute), AC 18.1–18.11 | §§7.2, 8.2, 9.4, 13.2; Properties 7/9 | 14.1–14.3 | 14.7, 29.4 (P7/P9), 30.3 | `SPECIFIED`; planned, no code/test evidence. |
| [R19 Safe Batch and Monthly Worker Confirmation](requirements.md#requirement-19-safe-batch-and-monthly-worker-confirmation), AC 19.1–19.12 | §§2, 9.4, 15.4; Properties 7/8 | 14.4–14.6 | 14.7, 29.4 (P7/P8), 30.3 | `SPECIFIED`; planned, no code/test evidence. |
| [R20 Monthly Work Record Card Experience](requirements.md#requirement-20-monthly-work-record-card-experience), AC 20.1–20.15 | §§2, 15.1–15.4, 24.6 | 15.1–15.5 | 15.6, 30.1 | `SPECIFIED`; planned, no code/test evidence. |

### 3.4 Compensation, calculation, review, finalization, and Payslips

| Requirement and acceptance criteria | Controlling design | Implementation tasks | Verification tasks / properties | Current evidence and status |
|---|---|---|---|---|
| [R21 Effective-Dated Compensation and Pay Rules](requirements.md#requirement-21-effective-dated-compensation-and-pay-rules), AC 21.1–21.14 | §§7.2, 8.4, 10, 13.1, 14; Properties 10/18 | 16.1–16.5 | 16.6, 17.7, 29.5 (P10/P18) | `SPECIFIED`; formulas/values must be versioned configuration, not hardcoded. |
| [R22 Configurable and Validated Payroll Policy](requirements.md#requirement-22-configurable-and-validated-payroll-policy), AC 22.1–22.9 | §§1.2–1.3, 14, 22, 26, 28–29 | 16.3–16.5, 17.3 | 16.6, 17.7, 30.7, 30.9 | `EXTERNAL_VALIDATION_REQUIRED`; register `XG-01`; `G-PAYROLL-PRODUCTION` closed. |
| [R23 Estimated Payroll Calculation and Provenance](requirements.md#requirement-23-estimated-payroll-calculation-and-provenance), AC 23.1–23.15 | §§7.2, 8.4, 9.5, 13–14; Properties 11/12/18 | 17.1–17.6 | 17.7, 29.5 (P11/P12/P18), 30.1 | `SPECIFIED`; planned, no code/test evidence. |
| [R24 Daily and Monthly HR Review](requirements.md#requirement-24-daily-and-monthly-hr-review), AC 24.1–24.11 | §§8.4, 9.5, 13, 15.6; Property 13 | 18.1–18.5 | 18.6, 29.5 (P13), 30.3 | `SPECIFIED`; planned, no code/test evidence. |
| [R25 Payroll Readiness and Finalization](requirements.md#requirement-25-payroll-readiness-and-finalization), AC 25.1–25.15 | §§8.4, 9.5, 13.3, 14.3, 19; Properties 14/15/16 | 19.1–19.6 | 19.7, 29.6 (P14/P15/P16), 30.1, 30.6 | `BLOCKED_BY_GATE` for production; development contracts planned; both rollout/payroll gates remain closed. |
| [R26 Finalized Payslips and Worker Monetary Visibility](requirements.md#requirement-26-finalized-payslips-and-worker-monetary-visibility), AC 26.1–26.10 | §§7.2, 9.5, 13.3, 15.6, 17, 20; Property 17 | 20.1–20.3 | 20.6, 29.7 (P17), 30.1, 30.6 | `BLOCKED_BY_GATE` for production publication; development contracts planned. |
| [R27 Post-Finalization Scope Boundary](requirements.md#requirement-27-post-finalization-scope-boundary), AC 27.1–27.6 | §§1.3, 13.3, 19, 29 | 20.4–20.5 | 20.6, 30.6 | `SPECIFIED`; excluded workflows remain non-executable. |

### 3.5 PWA, offline, responsive experience, accessibility, and notifications

| Requirement and acceptance criteria | Controlling design | Implementation tasks | Verification tasks / properties | Current evidence and status |
|---|---|---|---|---|
| [R28 Worker and Supervisor Browser-Capable Installable PWAs](requirements.md#requirement-28-worker-and-supervisor-browser-capable-installable-pwas), AC 28.1–28.13 | §§4.2, 5.1, 7, 8.5, 15.7–15.9, 17; Properties 21/24 | 21.1–21.5 | 21.6, 29.7 (P24), 30.4 | `SPECIFIED`; distribution/support matrix gates remain closed. |
| [R29 Pending Offline Actions and Reconciliation](requirements.md#requirement-29-pending-offline-actions-and-reconciliation), AC 29.1–29.15 | §§7.2, 8.5, 15.8, 17; Property 22 | 22.1–22.4 | 22.7, 29.7 (P22), 29.9, 30.4 | `SPECIFIED`; offline policy values must be versioned configuration/Admin Panel driven. |
| [R30 Progressive Capability Fallback](requirements.md#requirement-30-progressive-capability-fallback), AC 30.1–30.11 | §§15.8, 17, 19; Property 25 | 22.5, 24.2 | 22.7, 29.7 (P25), 30.4 | `SPECIFIED`; planned, no code/test evidence. |
| [R31 Sensitive Client Data and Session Clearing](requirements.md#requirement-31-sensitive-client-data-and-session-clearing), AC 31.1–31.10 | §§7.2–7.3, 15.8, 17, 20; Property 23 | 20.2, 22.2, 22.6, 24.5 | 20.6, 22.7, 29.7 (P23), 30.4 | `SPECIFIED`; planned, no code/test evidence. |
| [R32 Responsive Product Experience and State Communication](requirements.md#requirement-32-responsive-product-experience-and-state-communication), AC 32.1–32.14 | §§15.1–15.4, 15.9, 24.6 | 15.2–15.5, 23.1–23.4 | 15.6, 23.6, 30.1 | `SPECIFIED`; planned, no code/test evidence. |
| [R33 Accessibility and Localization Readiness](requirements.md#requirement-33-accessibility-and-localization-readiness), AC 33.1–33.13 | §§15.1–15.4, 24.6, 28 | 15.5, 23.4–23.5, 28.5 | 15.6, 23.6, 30.8–30.9 | `BLOCKED_BY_GATE` for rollout target/matrix; readiness implementation planned. |
| [R34 Notifications and Action Deadlines](requirements.md#requirement-34-notifications-and-action-deadlines), AC 34.1–34.14 | §§12, 16, 21, 24 | 24.1–24.5 | 24.6, 30.4 | `SPECIFIED`; channel/timing/template values must be versioned configuration/Admin Panel driven. |

### 3.6 Security, audit, operations, export, retention, migration, and quality gates

| Requirement and acceptance criteria | Controlling design | Implementation tasks | Verification tasks / properties | Current evidence and status |
|---|---|---|---|---|
| [R35 Security and Privacy Controls](requirements.md#requirement-35-security-and-privacy-controls), AC 35.1–35.14 | §§4.5, 7.3, 15.8, 18, 20, 24.5 | 25.1–25.3 | 25.7, 28.5, 30.9 | `DRAFT_GOVERNANCE` for access controls only; no security/privacy approval claimed. |
| [R36 Auditability and Integrity](requirements.md#requirement-36-auditability-and-integrity), AC 36.1–36.12 | §§4.3–4.4, 13, 20.4, 22.2, 24.5 | 25.4–25.6 | 25.7, 27.7, 30.5–30.6 | `SPECIFIED`; draft baseline defines proposed audit fields; implementation planned. |
| [R37 Observability, Reliability, and Recovery](requirements.md#requirement-37-observability-reliability-and-recovery), AC 37.1–37.11 | §§5, 12, 21, 24.4, 24.7 | 26.1–26.5 | 26.6, 28.6, 30.8–30.9 | `BLOCKED_BY_GATE` for numeric targets/providers; neutral contracts specified. |
| [R38 Performance and Scalability Gates](requirements.md#requirement-38-performance-and-scalability-gates), AC 38.1–38.10 | §§21, 24.7, 27, 28 | 28.1–28.4 | 28.6, 30.8 | `BLOCKED_BY_GATE` for numeric budgets; workload dimensions specified. |
| [R39 Export and Integration Extension Points](requirements.md#requirement-39-export-and-integration-extension-points), AC 39.1–39.14 | §§10–12, 20, 22, 24 | 27.1–27.2 | 27.7 | `BLOCKED_BY_GATE` for contracts/vendors; neutral export/integration contracts specified. |
| [R40 Retention, Data Lifecycle, and Privacy Rights](requirements.md#requirement-40-retention-data-lifecycle-and-privacy-rights), AC 40.1–40.8 | §§20.3, 22.1, 24 | 25.3, 27.3 | 27.7, 30.9 | `BLOCKED_BY_GATE` for durations/residency/legal constraints; no approval claimed. |
| [R41 Historical Paper Migration and Reconciliation](requirements.md#requirement-41-historical-paper-migration-and-reconciliation), AC 41.1–41.9 | §§2, 22, 25, 24 | 27.4–27.6 | 27.7, 30.7 | `EXTERNAL_VALIDATION_REQUIRED` for go-live reconciliation; no digital event fabrication permitted. |
| [R42 Implementation and Rollout Decision Gates](requirements.md#requirement-42-implementation-and-rollout-decision-gates), AC 42.1–42.11 | §§1.1, 26–28 | 1.1–1.4; gate-dependent tasks 28.1–28.5 | 28.6, 30.7–30.10 | `DRAFT_GOVERNANCE`: decision register, access baseline, and this matrix exist; all named gates retain recorded state. |
| [R43 Correctness and Testability](requirements.md#requirement-43-correctness-and-testability), AC 43.1–43.14 | §§23–24 | Test-bearing subtasks throughout 3–28; 29.1–29.9 | 29.1–29.9, 30.1–30.10 | `SPECIFIED`; test taxonomy/property metadata/iteration/update rules controlled here; no test execution claimed. |

## 4. Correctness-property crosswalk

All 27 named design properties have planned implementation in Task 29; this table prevents a property from becoming detached from its requirements or test task.

| Design properties | Requirement coverage | Property-test task |
|---|---|---|
| P1 Tenant Isolation; P19 Platform-Admin Isolation; P20 Support Boundedness; P21 Surface Equivalence | R1, R2, R4 | 29.1 |
| P2 Closed-Day Prohibition; P3 Closed Preservation; P4 Bounded Closed-Day Authorization | R9, R13, R14 | 29.2 |
| P5 Dual Confirmation Uniqueness; P6 QR Replay Resistance | R11, R15 | 29.3 |
| P7 No Unsafe Batch Signature; P8 Per-Record Traceability; P9 Acknowledgement Invalidation | R17, R18, R19 | 29.4 |
| P10 Effective Compensation; P11 Calculation Idempotency; P12 No Duplicate Counting; P13 Stale Review; P18 Fixed-Precision Conservation | R21, R23, R24 | 29.5 |
| P14 Finalization Atomicity; P15 Provenance Closure; P16 Maker-Checker Separation | R6, R25 | 29.6 |
| P17 Estimate Visibility; P22 Offline Non-Finality; P23 Safe Purge; P24 Browser/Install Parity; P25 Fallback Safety; P26 Provisioning Convergence; P27 Entitlement Non-Authority | R2, R3, R26, R28–R31 | 29.7 |
| Generator iterations, shrinking, deterministic replay, minimized regression fixtures | R43.2–43.5 | 29.8 |
| Calendar, overtime, Work_Record, batch, offline, review, and finalization state machines | R43.6–43.8 plus linked domain requirements | 29.9 |

## 5. Code → Tests evidence register

This register contains actual evidence only. Future tasks append records; they do not replace historical entries. “Not applicable” below is a reasoned value, not a placeholder: Task 1.4 creates specification governance and no application code.

| Evidence ID | Acceptance criteria | Task | Production/code evidence | Test/validation evidence | Result | Policy/config and gate effect |
|---|---|---|---|---|---|---|
| `TR-001` | R42.11 | 1.4 | `traceability-matrix.md`, complete R1–R43 matrix and update contract | Structural validation checks all 43 requirement rows, 27 properties, task references, relative links, closed-gate text, and draft markings | Passed at Task 1.4 completion | No policy activated; `G-TENANT-ROLLOUT` and `G-PAYROLL-PRODUCTION` remain closed. |
| `TR-002` | R43.1–R43.14 | 1.4 | `traceability-matrix.md` §§2, 4–5 | Manual cross-review against requirements §43, design §§23–24, and tasks 29.1–29.9 | Passed at Task 1.4 completion | No testing library selected and no test run claimed; `TECH-03` remains a blocking prerequisite. |
| `TR-003` | R2.5–2.9, R4.1–4.13, R6.1–6.11, R35.1–35.2, R42.3 | 1.3 carried baseline | `access-governance-baseline.md` and register `XG-02`; no application code | Baseline conformance cases `AG-001`–`AG-043` are specified, not executed | `DRAFT_GOVERNANCE`, not verified implementation | Every value is configuration-table/Admin Panel driven; explicit approval required; rollout gate closed. |
| `TR-004` | R22.1–22.9 | 1.2 disposition | Register `XG-01`; no application code or policy approval | Qualified manual validation not yet performed | `EXTERNAL_VALIDATION_REQUIRED` | No Singapore MOM/legal/payroll/privacy/contract/compliance approval claimed; payroll production gate closed. |

## 6. Coverage and self-review record

Task 1.4 completion requires all of the following to remain true:

- Requirement coverage: 43 of 43 top-level requirements represented, each with its full acceptance-criteria range.
- Property coverage: 27 of 27 named design properties mapped to Tasks 29.1–29.7; property harness and state-machine controls mapped to 29.8–29.9.
- Planned task coverage: every requirement has at least one implementation or governance task and one verification/evidence path.
- Link integrity: every local Markdown link resolves to an existing specification file; linked requirement anchors are present.
- No application implementation was started, and Task 2.1 has no Task 1.4 evidence record.
- No unverified statutory, contractual, payroll, legal, privacy, or compliance value or approval was introduced.
- Every policy-bearing row requires versioned configuration tables and applicable Admin Panel management; no hardcoded policy default is authorized.
- `G-TENANT-ROLLOUT` and `G-PAYROLL-PRODUCTION` remain closed.
