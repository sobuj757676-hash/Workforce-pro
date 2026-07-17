# Access Governance Baseline

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**
>
> This marking applies independently to every policy in this document: Permission catalog, Role templates, custom-role constraints, risk tiers, Maker_Checker_Policy, Segregation_of_Duties, step-up authentication, access review, support access, and Break_Glass_Session governance.

## 1. Status, scope, and authority

| Field | Value |
|---|---|
| Policy identifier | `ACCESS-GOVERNANCE-BASELINE` |
| Version | `1.0.0-draft` |
| Status | **DRAFT – NOT APPROVED – NOT FOR PRODUCTION** |
| Task | 1.3 |
| Controlling gates | `XG-02` and `G-TENANT-ROLLOUT` remain closed |
| Payroll production gate | `G-PAYROLL-PRODUCTION` remains closed under `XG-01` |
| Approval needed | Explicit Tenant Security Owner and Platform Security approval of this exact version and digest before production rollout |

This document is a complete, review-ready baseline produced after explicit direction to draft an unapproved baseline using industry-standard RBAC, Maker-Checker, Segregation of Duties, Break Glass, least-privilege, and audit-logging practices. It does not claim Singapore MOM, legal, contract, payroll, privacy, compliance, Qualified_Stakeholder, Tenant, or security approval. Approval must create a new immutable policy version with approver identities, scope, effective time, canonical digest, and evidence reference; it must not relabel this draft version in place.

The requirements are authoritative, followed by the design. This baseline defines access governance only. It does not define salary amounts, payroll formulas, statutory rates, thresholds, standard days or hours, rounding, legal interpretations, or production payroll eligibility. No role title, user interface, feature entitlement, plan, installation, support case, or approval request is authority by itself.

## 2. Normative authorization semantics

### 2.1 Decision inputs and order

Every server authorization decision MUST evaluate these inputs from current authoritative state:

1. authenticated Principal, Principal kind, Account status, session status, and effective Tenant derived by the server;
2. platform Control_Plane versus Tenant_Data_Plane boundary;
3. Tenant lifecycle policy for the command;
4. active Role assignment versions and explicit Permission grants;
5. intersection of every applicable Resource_Scope;
6. resource Tenant ownership and field/data classification;
7. workflow state, resource version, submitted digest, and policy version;
8. Segregation_of_Duties rules and actor participation in the resource lineage;
9. Maker_Checker_Policy approval count, checker qualification, scope, digest, and policy version;
10. required Step_Up_Authentication level and freshness;
11. support or Break_Glass_Session purpose, scope, data classes, status, and expiry when applicable;
12. immutable-domain and non-downgradable integrity rules.

Evaluation outcomes are `DENY`, `REQUIRE_STEP_UP`, `REQUIRE_APPROVAL`, or `ALLOW`. Mandatory denial has precedence. A later allow, Role, entitlement, feature flag, support grant, break-glass approval, or surface MUST NOT override an earlier mandatory denial. Denial causes no authoritative domain mutation and does not disclose cross-Tenant resource existence.

### 2.2 Effective grants

An effective grant is the intersection of:

`active account ∩ active server session ∩ active Role-version assignment ∩ Permission ∩ assignment Resource_Scope ∩ resource scope ∩ Tenant lifecycle allowance ∩ workflow policy ∩ current policy version`.

For a Support_Session, the intersection additionally includes:

`Platform_Support_Agent eligibility ∩ Support_Grant purpose ∩ Support_Grant actions ∩ Support_Grant data classes ∩ Support_Grant Resource_Scope ∩ Support_Session expiry ∩ Tenant support policy`.

A Support_Grant does not create a Tenant Role or standing Tenant identity. A Break_Glass_Session does not create standing access and cannot grant a Permission prohibited in Section 10.

### 2.3 Scope model

The only scope dimensions are `Tenant`, `Organization_Unit`, `Site`, `Team`, `Worker`, `Work_Assignment`, `Local_Date_Interval`, and `Payroll_Period`. Scope containment is conjunctive across dimensions. Missing scope never means all scope. `OWN` is resolved from the authenticated Principal-to-Worker relationship and cannot be supplied by the caller. `ASSIGNED` is resolved from the approved effective Work_Assignment version for the action date.

Roles provide Permissions; assignments provide Resource_Scope. A Role template never supplies global scope. Platform Permissions cannot be placed in Tenant Roles, Tenant Permissions cannot be placed in Platform Roles, and one assignment cannot cross a Tenant boundary.

## 3. Risk tiers and authentication policy

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

### 3.1 Risk tiers

| Tier | Definition | Default decision controls |
|---|---|---|
| `R1_STANDARD` | Own/assigned-scope operational reads and low-impact actions that do not expose Sensitive_Data or approve another actor's work. | Current authenticated session; explicit Permission and scope. |
| `R2_SENSITIVE` | Tenant administration, operational mutation, protected operational evidence, or broad non-monetary reads. | MFA-authenticated session; reason and Audit_Event for mutation. |
| `R3_HIGH` | Salary/protected-identity/Payslip read, sensitive export, support session use, high-risk access administration, or approval of a high-impact operational action. | `STEP_UP_1`; immutable decision record; maker-checker where listed in Section 8. |
| `R4_CRITICAL` | Payroll finalization, salary approval, access-policy activation, Break_Glass_Session, Tenant suspension/termination, or assignment of an `R4_CRITICAL` Permission. | `STEP_UP_2`; exact-digest approval; critical audit atomicity; maker-checker count in Section 8. |

Risk is the maximum of Permission risk, resource data classification, action risk, and workflow-state risk. Tenant configuration may raise but MUST NOT lower the tier defined here. There are no monetary amount thresholds in this baseline: any salary change, Adjustment, salary-bearing export, or payroll finalization is governed regardless of amount. This avoids inventing payroll or legal thresholds.

### 3.2 Step-up levels

| Level | Exact baseline requirement |
|---|---|
| `SESSION` | Current non-revoked authenticated session. |
| `STEP_UP_1` | Successful MFA with primary-authentication revalidation no more than 15 minutes before the action; bound to Principal, Tenant, session, and action family. |
| `STEP_UP_2` | Successful phishing-resistant MFA with primary-authentication revalidation no more than 5 minutes before the action; one-time authorization bound to Principal, Tenant, action, resource, current version, and submitted digest. |

Step-up expires immediately on logout, session/device/account revocation, Role or scope loss, Tenant suspension, security reset, or digest/version change. Step-up is actor-specific and cannot be delegated or reused by a support session. An implementation may use a stronger approved method without changing this policy version.

## 4. Permission catalog

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

### 4.1 Catalog rules

Permission identifiers are stable, case-sensitive, and deny-by-default. Wildcards and prefix matching are forbidden. `Max scope` is a ceiling, not an automatic grant. `TENANT` means an assignment may be narrowed to any Tenant scope dimension; `ASSIGNED` and `OWN` cannot be broadened. Every mutation requires idempotency, Expected_Version where mutable, stable safe errors, and Audit_Event behavior required by the controlling requirements.

### 4.2 Platform Control_Plane catalog

| Permission | Action/resource | Data boundary | Max scope | Risk | Mandatory constraints |
|---|---|---|---|---|---|
| `platform.tenant.view` | View Tenant lifecycle/onboarding metadata | Control plane only | Platform portfolio | R2_SENSITIVE | No Tenant payroll payload. |
| `platform.tenant.provision` | Request/resume Tenant provisioning | Control plane only | Platform portfolio | R3_HIGH | Idempotent; activation gates remain mandatory. |
| `platform.tenant.status.submit` | Submit exact Tenant lifecycle transition and reason | Control plane only | Platform portfolio | R3_HIGH | Step-up 1; maker identity/from-version/digest recorded. |
| `platform.tenant.status.change` | Approve and execute restrict, suspend, reactivate, or termination transition | Control plane only | Platform portfolio | R4_CRITICAL | Step-up 2; exact transition digest; maker-checker. |
| `platform.security.defaults.view` | View platform security defaults | Control plane only | Platform | R2_SENSITIVE | No Tenant secret values. |
| `platform.security.defaults.change` | Submit security-default policy version | Control plane only | Platform | R4_CRITICAL | Step-up 2; two-checker activation. |
| `platform.release.view` | View feature/release state | Control plane only | Platform | R1_STANDARD | Entitlements never grant domain Permission. |
| `platform.release.change` | Change controlled feature release | Control plane only | Platform | R3_HIGH | Step-up 1; maker-checker; cannot bypass domain rules. |
| `platform.support.case.view` | View assigned Support_Case metadata | Control plane only | Assigned cases | R2_SENSITIVE | Purpose metadata only; no Tenant payload. |
| `platform.support.case.manage` | Create/update assigned Support_Cases | Control plane only | Assigned cases | R2_SENSITIVE | A case grants no Tenant access. |
| `platform.support.session.request` | Request exact Support_Grant scope | Control plane request | Assigned cases | R3_HIGH | Step-up 1; purpose/actions/data classes/duration required. |
| `platform.support.session.start` | Start session from approved grant | Broker metadata | Approved grant | R3_HIGH | Step-up 1; agent eligibility and grant intersection. |
| `platform.support.session.revoke` | Revoke a Support_Session | Broker metadata | Assigned cases | R3_HIGH | Immediate subsequent denial. |
| `platform.support.operational.use` | Use approved operational support actions | Tenant data via broker | Grant intersection | R3_HIGH | Active Support_Session; Tenant domain controls remain. |
| `platform.support.protected_id.use` | View approved full protected identifiers | Tenant data via broker | Grant intersection | R4_CRITICAL | Step-up 2; exact data class; access audit. |
| `platform.support.evidence.use` | View approved high-sensitivity evidence | Tenant data via broker | Grant intersection | R4_CRITICAL | Step-up 2; exact evidence class; access audit. |
| `platform.support.salary.use` | View approved salary data | Tenant data via broker | Grant intersection | R4_CRITICAL | Step-up 2; exact salary-read scope; read only. |
| `platform.support.payslip.use` | View approved Payslip content | Tenant data via broker | Grant intersection | R4_CRITICAL | Step-up 2; exact Payslip scope; read only. |
| `platform.support.payroll_mutation.use` | Perform explicitly approved mutable payroll support action | Tenant data via broker | Grant intersection | R4_CRITICAL | Step-up 2; domain Permission/policy; never snapshot mutation/finalization. |
| `platform.break_glass.request` | Request Break_Glass_Session | Broker metadata | One Tenant/case | R4_CRITICAL | Section 10 governs. |
| `platform.break_glass.approve.security` | Security approval of exact request | Broker metadata | One Tenant/case | R4_CRITICAL | Distinct from requester and other approver. |
| `platform.break_glass.approve.incident` | Incident-command approval of exact request | Broker metadata | One Tenant/case | R4_CRITICAL | Distinct from requester and security approver. |
| `platform.break_glass.start` | Start fully approved session | Broker metadata | Approved request | R4_CRITICAL | Step-up 2; agent must be named in request. |
| `platform.break_glass.review` | Complete independent post-use review | Broker/audit metadata | One session | R4_CRITICAL | Reviewer did not request, approve, or use session. |
| `platform.audit.view` | View platform Audit_Events | Platform audit | Assigned platform scope | R2_SENSITIVE | No implicit Tenant audit access. |
| `platform.audit.export` | Export platform Audit_Events | Platform audit | Assigned platform scope | R3_HIGH | Step-up 1; protected delivery and audit. |

### 4.3 Tenant access and administration catalog

| Permission | Action/resource | Max scope | Risk | Mandatory constraints |
|---|---|---|---|---|
| `tenant.access.role.view` | View Role templates/versions and grants | Tenant | R2_SENSITIVE | Sensitive grants masked unless separately authorized by this Permission's scope. |
| `tenant.access.role.draft` | Draft custom Role version | Tenant | R3_HIGH | Catalog IDs only; Section 6 constraints. |
| `tenant.access.role.submit` | Submit Role version/assignment digest | Tenant | R3_HIGH | Maker record required. |
| `tenant.access.role.approve` | Approve non-critical Role/assignment | Tenant | R3_HIGH | Distinct checker; exact digest/scope. |
| `tenant.access.role.approve_critical` | Approve Role/assignment containing R4_CRITICAL | Tenant | R4_CRITICAL | One of two distinct critical checkers. |
| `tenant.access.role.activate` | Activate fully approved Role version/assignment | Tenant | R4_CRITICAL | System-only completion after required approvals. |
| `tenant.access.role.revoke` | Revoke Role/assignment | Tenant | R3_HIGH | Step-up 1; subsequent decisions use revocation. |
| `tenant.access.policy.view` | View access/SoD/risk policy versions | Tenant | R2_SENSITIVE | Historical versions retained. |
| `tenant.access.policy.draft` | Draft stricter Tenant access policy | Tenant | R4_CRITICAL | Cannot lower this baseline. |
| `tenant.access.policy.approve` | Approve exact Tenant policy version | Tenant | R4_CRITICAL | Two distinct checkers; exact digest. |
| `tenant.access.review.execute` | Review current sensitive grants | Tenant | R3_HIGH | Cannot approve own grant; decision audited. |
| `tenant.account.view` | View Account lifecycle and non-secret grants | Tenant | R2_SENSITIVE | Field authorization still applies. |
| `tenant.account.manage` | Create/update/disable Account | Tenant | R3_HIGH | Role activation is separate. |
| `tenant.session.revoke` | Revoke Account sessions | Tenant | R3_HIGH | Immediate subsequent denial/purge signal. |
| `tenant.device.revoke` | Revoke registered device/installation | Tenant | R3_HIGH | Immediate session/local-state invalidation signal. |

### 4.4 Workforce, calendar, attendance, overtime, and Work_Record catalog

| Permission | Action/resource | Max scope | Risk | Mandatory constraints |
|---|---|---|---|---|
| `tenant.organization.view` | View Organization_Unit/Site/Checkpoint | Tenant | R1_STANDARD | Tenant-scoped. |
| `tenant.organization.manage` | Create/update Organization_Unit/Site | Tenant | R2_SENSITIVE | Versioned and audited. |
| `tenant.checkpoint.manage` | Create/update Checkpoint metadata | Site | R3_HIGH | Signing-key management remains separate. |
| `tenant.worker.view` | View non-protected Worker_Profile fields | Tenant | R2_SENSITIVE | Salary/protected IDs excluded. |
| `tenant.worker.protected_id.view` | View full protected employment identifier | Worker | R3_HIGH | Step-up 1; access audit. |
| `tenant.worker.manage` | Create/update Worker_Profile/Employment_Record | Tenant | R3_HIGH | Does not grant protected-ID or salary view. |
| `tenant.assignment.view` | View Work_Assignment | Tenant | R1_STANDARD | Effective-date scope. |
| `tenant.assignment.manage` | Draft/submit Work_Assignment version | Tenant | R2_SENSITIVE | Cross-Tenant references denied. |
| `tenant.assignment.approve` | Approve Work_Assignment version | Tenant | R3_HIGH | Distinct checker where policy applies. |
| `tenant.calendar.view` | View calendars and publication state | Tenant | R1_STANDARD | Published/draft visibility follows scope. |
| `tenant.calendar.draft` | Create/edit calendar draft | Tenant | R2_SENSITIVE | Every valid date classified. |
| `tenant.calendar.submit` | Submit exact calendar digest | Tenant | R2_SENSITIVE | Maker bound to version/digest. |
| `tenant.calendar.publish` | Approve/publish calendar | Tenant | R3_HIGH | Distinct checker; immutable publication. |
| `tenant.calendar.impact_override` | Approve policy-defined blocking impact override | Tenant | R4_CRITICAL | Cannot override missing classification or Tenant mismatch. |
| `tenant.meeting.view` | View assigned Toolbox_Meeting/roster | Assigned | R1_STANDARD | Minimized assigned roster only. |
| `tenant.meeting.create` | Create ordinary Working-day meeting | Assigned | R2_SENSITIVE | Calendar/assignment guard. |
| `tenant.meeting.create_closed_special` | Create meeting under Closed_Day_OT authorization | Assigned | R3_HIGH | Exact authorization window/scope. |
| `tenant.attendance.supervisor_confirm` | Confirm assigned Worker arrival | Assigned | R2_SENSITIVE | Meeting roster/version bound. |
| `tenant.attendance.worker_confirm_own` | Confirm own arrival | Own | R1_STANDARD | Authenticated Worker only. |
| `tenant.attendance.exception.request` | Request attendance exception | Own or Assigned | R2_SENSITIVE | Requested effect and evidence required. |
| `tenant.attendance.exception.decide` | Approve/reject attendance exception | Assigned | R3_HIGH | Cannot decide own request; scope/evidence audit. |
| `tenant.overtime.view` | View OT plan/assignment | Own or Assigned | R1_STANDARD | Exact scope; terms by authorization. |
| `tenant.overtime.propose` | Draft/submit OT plan | Assigned | R2_SENSITIVE | Named/bounded scope. |
| `tenant.overtime.authorize` | Approve OT plan | Assigned | R3_HIGH | Distinct checker; exact digest. |
| `tenant.overtime.consent_own` | Consent/decline own assignment | Own | R2_SENSITIVE | Before start; exact digest; cannot backdate. |
| `tenant.overtime.exception.decide` | Decide OT variance/emergency exception | Assigned | R3_HIGH | Cannot fabricate consent/ordinary authorization. |
| `tenant.work_record.view_own` | View own Work_Record | Own | R1_STANDARD | Evidence fields separately authorized. |
| `tenant.work_record.view_assigned` | View assigned Work_Record | Assigned | R2_SENSITIVE | Salary omitted. |
| `tenant.work_record.evidence.view` | View operational evidence summary | Own or Assigned | R2_SENSITIVE | High-sensitivity payload excluded. |
| `tenant.work_record.evidence.view_sensitive` | View high-sensitivity evidence | Assigned | R3_HIGH | Step-up 1; access audit. |
| `tenant.work_record.checkout_own` | Submit own checkout evidence | Own | R2_SENSITIVE | QR/session/record binding. |
| `tenant.work_record.verify` | Verify assigned record unchanged | Assigned | R2_SENSITIVE | Exact current digest. |
| `tenant.work_record.edit` | Reason-coded edit of assigned record | Assigned | R3_HIGH | New version; prior signatures historical only. |
| `tenant.work_record.acknowledge_own` | Acknowledge own current record | Own | R2_SENSITIVE | Exact version/digest. |
| `tenant.work_record.dispute_own` | Dispute own current record | Own | R2_SENSITIVE | Reason required. |
| `tenant.work_record.exception.request` | Request evidence/work exception | Own or Assigned | R2_SENSITIVE | Requested effect/evidence required. |
| `tenant.work_record.exception.decide` | Decide evidence/work exception | Assigned | R3_HIGH | Cannot decide own request; payroll impact recorded. |

### 4.5 Compensation, payroll, Payslip, export, audit, and integration catalog

| Permission | Action/resource | Max scope | Risk | Mandatory constraints |
|---|---|---|---|---|
| `tenant.salary.view` | View compensation, salary, monetary estimates, and salary-bearing lines | Worker/Payroll_Period | R3_HIGH | Separate from Worker identity and all edit/approval/finalization Permissions. |
| `tenant.salary.edit` | Draft/edit Compensation_Profile or salary-bearing configuration | Worker | R3_HIGH | Does not grant view, approval, finalization, or export. |
| `tenant.salary.submit` | Submit exact compensation/salary digest | Worker | R3_HIGH | Maker identity/version recorded. |
| `tenant.salary.approve` | Approve exact compensation/salary digest | Worker | R4_CRITICAL | Distinct from maker; Step-up 2. |
| `tenant.adjustment.prepare` | Draft/submit reason-coded Adjustment | Worker/Payroll_Period | R3_HIGH | Any amount; no monetary threshold. |
| `tenant.adjustment.approve` | Approve exact Adjustment digest | Worker/Payroll_Period | R4_CRITICAL | Distinct from maker; Step-up 2. |
| `tenant.pay_rule.view` | View approved/draft rule definitions | Tenant | R2_SENSITIVE | Does not assert legal correctness. |
| `tenant.pay_rule.draft` | Draft constrained Pay_Rule version | Tenant | R3_HIGH | No arbitrary code; production gate remains closed. |
| `tenant.pay_rule.approve` | Approve exact Pay_Rule version | Tenant | R4_CRITICAL | Distinct checker; does not open payroll production gate. |
| `tenant.payroll.operational.view` | View status, hours, blockers without money | Worker/Payroll_Period | R2_SENSITIVE | Monetary fields omitted. |
| `tenant.payroll.calculate` | Request calculation/recalculation | Worker/Payroll_Period | R2_SENSITIVE | Output remains Estimated/Not Finalized. |
| `tenant.payroll.ledger.view` | View salary-bearing payroll ledger | Worker/Payroll_Period | R3_HIGH | Requires this explicit Permission; salary data access audited. |
| `tenant.payroll.review` | Approve/reject current Input_Digest | Worker/Payroll_Period | R3_HIGH | Current digest; stale change reopens. |
| `tenant.payroll.readiness.run` | Run Readiness_Report | Payroll_Period | R3_HIGH | Cannot downgrade mandatory blockers. |
| `tenant.payroll.finalization.submit` | Submit readiness/current digest for finalization | Payroll_Period | R4_CRITICAL | Maker identity/digest recorded; no snapshot created. |
| `tenant.payroll.finalize` | Atomically create immutable Payroll_Snapshot | Payroll_Period | R4_CRITICAL | Separate Permission; Step-up 2; distinct checker/finalizer; all gates/readiness apply. |
| `tenant.payslip.view_own` | View own published Payslip | Own | R3_HIGH | Fresh Step-up 1; no general cache. |
| `tenant.payslip.view` | View another Worker's Payslip | Worker/Payroll_Period | R4_CRITICAL | Step-up 2; separately authorized; access audit. |
| `tenant.export.work_record` | Export operational monthly work card | Assigned/Payroll_Period | R2_SENSITIVE | No salary/evidence payload unless separately permitted. |
| `tenant.export.overtime` | Export OT package | Assigned/Payroll_Period | R2_SENSITIVE | Preserve authorization/consent states. |
| `tenant.export.payroll` | Export salary-bearing payroll ledger | Worker/Payroll_Period | R3_HIGH | Separate export Permission; Step-up 1; maker-checker. |
| `tenant.export.payslip` | Export Payslip package | Worker/Payroll_Period | R4_CRITICAL | Step-up 2; maker-checker; protected delivery. |
| `tenant.export.evidence` | Export protected evidence package | Assigned/Payroll_Period | R4_CRITICAL | Step-up 2; maker-checker; Integrity_Manifest. |
| `tenant.audit.view` | Query authorized Audit_Events | Tenant | R3_HIGH | Row/field/time/classification scope. |
| `tenant.audit.export` | Export authorized Audit_Events | Tenant | R4_CRITICAL | Step-up 2; maker-checker; Integrity_Manifest. |
| `tenant.integration.view` | View Integration_Principal metadata | Tenant | R2_SENSITIVE | Secrets never returned. |
| `tenant.integration.manage` | Create/rotate/revoke Integration_Principal scope | Tenant | R4_CRITICAL | Step-up 2; two-checker activation. |
| `tenant.integration.invoke` | Invoke one explicit machine action | Contract scope | R2_SENSITIVE–R4_CRITICAL | Risk equals mapped action; no human-only own/sign/approve actions. |

## 5. Default Role templates

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

### 5.1 Template semantics

Templates are immutable, versioned collections of Permission IDs. Assignment is always Tenant/resource scoped and never implied by a title. A Principal receives no authority merely because the Principal is described as an administrator, Worker, Supervisor, reviewer, maker, checker, finalizer, auditor, operator, or support agent. Templates with no listed Permission have no grants. No default template combines maker and checker Permissions for the same governed action.

### 5.2 Platform templates

| Template | Explicit Permission grants | Required assignment scope | Prohibited additions to this template version |
|---|---|---|---|
| `PLATFORM_TENANT_LIFECYCLE_OPERATOR` | `platform.tenant.view`, `platform.tenant.provision`, `platform.tenant.status.submit` | Platform portfolio | Tenant_Data_Plane and status-checker Permissions. |
| `PLATFORM_TENANT_LIFECYCLE_CHECKER` | `platform.tenant.view`, `platform.tenant.status.change` | Platform portfolio | Provisioning-maker and Tenant_Data_Plane Permissions. |
| `PLATFORM_SUPPORT_COORDINATOR` | `platform.support.case.view`, `platform.support.case.manage`, `platform.support.session.request` | Assigned cases | Support-session use, Tenant Roles, break-glass approval. |
| `PLATFORM_SUPPORT_AGENT_OPERATIONAL` | `platform.support.case.view`, `platform.support.session.start`, `platform.support.session.revoke`, `platform.support.operational.use` | Assigned cases plus grant intersection | Sensitive support-use Permissions unless separately assigned through a different approved Role version. |
| `PLATFORM_SUPPORT_AGENT_SENSITIVE_READ` | `platform.support.protected_id.use`, `platform.support.evidence.use`, `platform.support.salary.use`, `platform.support.payslip.use` | Assigned cases plus grant intersection | Payroll mutation, finalization, export, role/policy administration. |
| `PLATFORM_SECURITY_APPROVER` | `platform.security.defaults.view`, `platform.security.defaults.change`, `platform.break_glass.approve.security`, `platform.support.session.revoke`, `platform.audit.view` | Platform | Break-glass request/use/review. |
| `PLATFORM_INCIDENT_APPROVER` | `platform.break_glass.approve.incident`, `platform.support.session.revoke`, `platform.audit.view` | Platform | Break-glass request/use/review/security approval. |
| `PLATFORM_BREAK_GLASS_OPERATOR` | `platform.break_glass.request`, `platform.break_glass.start` | Named incident/Tenant | Both approval Permissions and post-use review. |
| `PLATFORM_BREAK_GLASS_REVIEWER` | `platform.break_glass.review`, `platform.audit.view` | Named completed sessions | Request, approval, start, or support-use Permissions. |
| `PLATFORM_RELEASE_MAKER` | `platform.release.view`, `platform.release.change` | Platform/release scope | Security-default or Tenant_Data_Plane Permissions. |
| `PLATFORM_AUDITOR` | `platform.audit.view`, `platform.audit.export`, `platform.tenant.view` | Assigned platform scope | All mutation and support-use Permissions. |

### 5.3 Tenant templates

| Template | Explicit Permission grants | Required assignment scope | Explicit exclusions |
|---|---|---|---|
| `TENANT_ACCESS_MAKER` | `tenant.access.role.view`, `tenant.access.role.draft`, `tenant.access.role.submit`, `tenant.access.policy.view`, `tenant.account.view`, `tenant.account.manage`, `tenant.session.revoke`, `tenant.device.revoke` | Tenant or narrower | Role/policy approval, salary, finalization. |
| `TENANT_ACCESS_CHECKER` | `tenant.access.role.view`, `tenant.access.role.approve`, `tenant.access.review.execute`, `tenant.access.policy.view` | Tenant or narrower | Role drafting/submission and critical approval. |
| `TENANT_ACCESS_CRITICAL_CHECKER` | `tenant.access.role.view`, `tenant.access.role.approve_critical`, `tenant.access.policy.view`, `tenant.access.policy.approve`, `tenant.access.review.execute` | Tenant | Role/policy drafting/submission; salary/finalization. |
| `TENANT_WORKFORCE_ADMIN` | `tenant.organization.view`, `tenant.organization.manage`, `tenant.checkpoint.manage`, `tenant.worker.view`, `tenant.worker.manage`, `tenant.assignment.view`, `tenant.assignment.manage` | Tenant or narrower | Protected-ID view, salary, assignment approval. |
| `TENANT_PROTECTED_ID_READER` | `tenant.worker.view`, `tenant.worker.protected_id.view` | Worker/Organization_Unit | Worker mutation, salary. |
| `TENANT_ASSIGNMENT_CHECKER` | `tenant.assignment.view`, `tenant.assignment.approve` | Tenant or narrower | Assignment management, salary. |
| `CALENDAR_MAKER` | `tenant.calendar.view`, `tenant.calendar.draft`, `tenant.calendar.submit` | Calendar scope | Calendar publish/impact override. |
| `CALENDAR_CHECKER` | `tenant.calendar.view`, `tenant.calendar.publish` | Calendar scope | Calendar draft/submit. |
| `CALENDAR_IMPACT_CHECKER` | `tenant.calendar.view`, `tenant.calendar.impact_override` | Calendar scope | Calendar draft/submit; cannot override structural blockers. |
| `OT_PROPOSER` | `tenant.overtime.view`, `tenant.overtime.propose` | Assigned scope | OT authorization. |
| `OT_AUTHORIZER` | `tenant.overtime.view`, `tenant.overtime.authorize` | Assigned scope | OT proposal on overlapping scope. |
| `WORKER_SELF_SERVICE` | `tenant.attendance.worker_confirm_own`, `tenant.overtime.view`, `tenant.overtime.consent_own`, `tenant.work_record.view_own`, `tenant.work_record.evidence.view`, `tenant.work_record.checkout_own`, `tenant.work_record.acknowledge_own`, `tenant.work_record.dispute_own`, `tenant.work_record.exception.request`, `tenant.payslip.view_own` | Own | Any non-own scope; salary estimate visibility still follows separate Tenant policy and response control. |
| `SUPERVISOR_OPERATIONS` | `tenant.meeting.view`, `tenant.meeting.create`, `tenant.meeting.create_closed_special`, `tenant.attendance.supervisor_confirm`, `tenant.attendance.exception.request`, `tenant.overtime.view`, `tenant.work_record.view_assigned`, `tenant.work_record.evidence.view`, `tenant.work_record.verify`, `tenant.work_record.edit`, `tenant.work_record.exception.request` | Assigned | Exception decisions, OT authorization, salary. |
| `OPERATIONS_EXCEPTION_CHECKER` | `tenant.attendance.exception.decide`, `tenant.overtime.exception.decide`, `tenant.work_record.exception.decide`, `tenant.work_record.view_assigned`, `tenant.work_record.evidence.view` | Assigned | Requester actions on overlapping resources; salary. |
| `HR_OPERATIONAL_REVIEWER` | `tenant.worker.view`, `tenant.assignment.view`, `tenant.work_record.view_assigned`, `tenant.work_record.evidence.view`, `tenant.payroll.operational.view` | Assigned | Salary, ledger money, finalization. |
| `PAYROLL_LEDGER_REVIEWER` | `tenant.salary.view`, `tenant.payroll.operational.view`, `tenant.payroll.ledger.view`, `tenant.payroll.review`, `tenant.payroll.readiness.run` | Worker/Payroll_Period | Salary edit/approval, adjustment, finalization, export. |
| `COMPENSATION_MAKER` | `tenant.salary.view`, `tenant.salary.edit`, `tenant.salary.submit`, `tenant.pay_rule.view`, `tenant.pay_rule.draft`, `tenant.adjustment.prepare`, `tenant.payroll.calculate` | Worker/Payroll_Period or Tenant rule scope | Salary/Pay_Rule/Adjustment approval and finalization. |
| `COMPENSATION_CHECKER` | `tenant.salary.view`, `tenant.salary.approve`, `tenant.pay_rule.view`, `tenant.pay_rule.approve`, `tenant.adjustment.approve` | Worker/Payroll_Period or Tenant rule scope | Maker Permissions and finalization. |
| `PAYROLL_FINALIZATION_MAKER` | `tenant.salary.view`, `tenant.payroll.ledger.view`, `tenant.payroll.review`, `tenant.payroll.readiness.run`, `tenant.payroll.finalization.submit` | Payroll_Period | `tenant.payroll.finalize`, compensation/adjustment approval. |
| `PAYROLL_FINALIZER` | `tenant.salary.view`, `tenant.payroll.ledger.view`, `tenant.payroll.readiness.run`, `tenant.payroll.finalize` | Payroll_Period | Finalization submission and source-making Permissions. |
| `OPERATIONAL_EXPORTER` | `tenant.export.work_record`, `tenant.export.overtime` | Assigned/Payroll_Period | Salary-bearing, Payslip, evidence, audit export. |
| `PAYROLL_EXPORT_MAKER` | `tenant.salary.view`, `tenant.payroll.ledger.view`, `tenant.export.payroll`, `tenant.export.payslip` | Worker/Payroll_Period | Export approval is a workflow decision by a distinct authorized checker; no finalization. |
| `PAYROLL_EXPORT_CHECKER` | `tenant.salary.view`, `tenant.payroll.ledger.view`, `tenant.export.payroll`, `tenant.export.payslip` | Worker/Payroll_Period | Cannot approve own export request; no finalization. |
| `EVIDENCE_EXPORT_MAKER` | `tenant.work_record.evidence.view_sensitive`, `tenant.export.evidence` | Assigned/Payroll_Period | Salary/finalization. |
| `EVIDENCE_EXPORT_CHECKER` | `tenant.work_record.evidence.view_sensitive`, `tenant.export.evidence` | Assigned/Payroll_Period | Cannot approve own export request; salary/finalization. |
| `TENANT_AUDITOR` | `tenant.audit.view`, `tenant.organization.view`, `tenant.assignment.view`, `tenant.calendar.view`, `tenant.payroll.operational.view` | Tenant or narrower | All mutations, salary, protected evidence, export. |
| `TENANT_SENSITIVE_AUDITOR` | `tenant.audit.view`, `tenant.audit.export`, `tenant.work_record.evidence.view_sensitive` | Tenant or narrower | All mutations, salary unless separately assigned, finalization. |
| `INTEGRATION_ADMIN` | `tenant.integration.view`, `tenant.integration.manage` | Tenant | `tenant.integration.invoke`, human signatures/approvals. |

No default Role grants `tenant.payslip.view`, `platform.support.payroll_mutation.use`, or `tenant.export.evidence` together with salary/finalization authority. Such access requires a custom Role or separate assignment that passes Sections 6–9.

## 6. Custom Role and assignment constraints

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

1. A custom Role MUST reference only Permission IDs in Section 4 and one boundary (`platform` or one Tenant).
2. Wildcards, implicit grants, inherited title semantics, unknown Permission IDs, and client-defined Permissions are invalid.
3. A Role version is immutable after submission. Any grant change creates a new version and digest; existing assignments remain bound to their assigned version until separately changed.
4. Every assignment states Principal, Role version, Resource_Scope, effective start, optional expiry, maker, policy version, and digest. Platform support eligibility Roles also state assigned Support_Case scope.
5. `OWN` Permissions can only be assigned to a Principal with an unambiguous current Worker relationship and cannot be broadened.
6. `ASSIGNED` Permissions resolve from approved date-effective Work_Assignment data; missing or ambiguous assignment denies the action.
7. Human-only consent, acknowledgement, dispute, approval, review, and finalization Permissions cannot be assigned to Integration_Principals.
8. Integration_Principals receive only explicit contract-mapped `tenant.integration.invoke` actions and cannot receive Role/policy administration, salary approval, payroll finalization, support, or break-glass Permissions.
9. A Principal cannot grant a Permission or scope beyond the grant ceiling explicitly delegated to that Principal. Access administration Permission alone is not a delegation ceiling.
10. Any Role or assignment containing an `R3_HIGH` Permission is a high-risk assignment. It requires one distinct qualified checker.
11. Any Role or assignment containing an `R4_CRITICAL` Permission requires two distinct qualified checkers, at least one holding `tenant.access.role.approve_critical` (or the corresponding platform security authority), and system activation only after both exact-digest approvals.
12. Any change in Principal, Role version, Permission set, Resource_Scope, effective interval, expiry, or policy version changes the assignment digest and invalidates prior pending approvals.
13. A checker cannot approve the checker's own Role, assignment, scope expansion, access-review decision, or any assignment benefiting an Account controlled by the same authenticated human identity.
14. Conflicting maker/checker Permissions may exist only when assigned scopes are provably disjoint. Overlap, missing scope, or ambiguous date coverage denies activation/action.
15. Salary view, salary edit, salary approval, payroll finalization, each sensitive export, and each support data class remain separate grants. None is a prerequisite that automatically grants another.
16. Tenant_Administrator designation, platform employment, support assignment, feature entitlement, or dashboard visibility grants no catalog Permission.
17. Custom Roles cannot weaken immutable Payroll_Snapshot rules, mandatory Tenant isolation, deny precedence, audit atomicity, step-up minima, Maker_Checker_Policy, or Section 10 break-glass prohibitions.
18. Revocation, Account disablement, session/device revocation, scope expiry, or Tenant lifecycle denial applies to the next server decision; cached authorization is not authoritative.
19. Sensitive grants are reviewable at least every 90 days in this proposed baseline and immediately after manager/assignment change, security incident, or extended leave. `R4_CRITICAL` grants are reviewable every 30 days. Review does not renew an expired grant.
20. Role assignment duration defaults to no more than 365 days for `R3_HIGH` and 90 days for `R4_CRITICAL`; shorter assignment or Tenant policy wins. Renewal is a new exact-digest approval decision.

## 7. Segregation_of_Duties rules

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

### 7.1 Identity and overlap semantics

Distinctness is evaluated by authenticated human identity, not Account, session, device, role title, API credential, or support-session identity. Multiple Accounts controlled by one human are one actor. An Integration_Principal cannot satisfy a human checker requirement. Scope overlap exists when the maker and checker scopes could both cover the same Tenant/resource/date/Payroll_Period; ambiguity is treated as overlap.

### 7.2 Mandatory rules

| Rule | Maker/requester action | Checker/decider action | Enforcement |
|---|---|---|---|
| `SOD-01-ROLE` | Draft/submit Role or assignment | Approve/activate same digest | Same actor denied; R4_CRITICAL requires two distinct checkers. |
| `SOD-02-POLICY` | Draft access/security policy | Approve/activate same version | Maker and both checkers distinct. |
| `SOD-03-CALENDAR` | Draft/submit calendar | Publish same version | Same actor denied for overlapping scope. |
| `SOD-04-OT` | Propose/submit OT plan | Authorize same plan | Same actor denied for overlapping scope. |
| `SOD-05-ASSIGNMENT` | Draft/submit Work_Assignment | Approve same assignment | Same actor denied where approval policy applies. |
| `SOD-06-EXCEPTION` | Request attendance/OT/work exception or benefit from requested effect | Decide same exception | Requester, affected Worker, and decider must not be same human; decider scope must cover effect. |
| `SOD-07-SALARY` | Edit/submit Compensation_Profile or salary-bearing configuration | Approve same digest | Same actor denied regardless of amount. |
| `SOD-08-PAY_RULE` | Draft/submit Pay_Rule | Approve same version | Same actor denied; approval does not satisfy compliance gate. |
| `SOD-09-ADJUSTMENT` | Prepare/submit Adjustment | Approve same digest | Same actor denied regardless of amount/sign. |
| `SOD-10-REVIEW` | Create/change a payroll source used by a ledger entry | Review that entry at resulting Input_Digest | Source maker cannot be sole reviewer; independent reviewer required. |
| `SOD-11-FINALIZE` | Submit finalization/readiness digest or make included compensation/Adjustment in current cycle | Finalize snapshot | Finalizer must differ from submitter and cannot finalize a result containing the finalizer's unindependently approved source change. |
| `SOD-12-EXPORT` | Request sensitive payroll/Payslip/evidence/audit export | Approve export manifest | Same actor denied; approval bound to exact query, fields, recipients, classification, expiry, and digest. |
| `SOD-13-SUPPORT` | Platform agent requests Support_Grant | Tenant approver grants access | Boundaries and humans distinct; agent cannot approve Tenant grant. |
| `SOD-14-BREAK_GLASS` | Request/use Break_Glass_Session | Security approval, incident approval, post-use review | Requester/user, two approvers, and reviewer separated as Section 10 states. |
| `SOD-15-AUDIT` | Perform sensitive action | Alter/delete its Audit_Event | Audit alteration/deletion Permission does not exist. |

A SoD denial is not bypassed by adding Roles. If a Principal holds both sides of a rule, the Principal may perform at most the side allowed by the resource's recorded lineage; the conflicting action is denied. Disjoint scopes may permit both Permissions, but scope is re-evaluated for every action.

## 8. Maker_Checker_Policy

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

### 8.1 Approval record

Every approval is a first-class immutable record containing Tenant/boundary, policy ID and version, action, resource ID/version, Resource_Scope, canonical submitted digest, maker, checker, checker Permission and scope at decision time, decision, reason, timestamp, step-up evidence reference, correlation ID, and Audit_Event reference. Changed content, scope, actor, version, policy, export recipient, or data classification creates a different digest and invalidates approval for the changed action.

### 8.2 Exact baseline matrix

| Governed action | Risk | Required distinct approvals | Required checker qualification | Activation condition |
|---|---:|---:|---|---|
| Non-critical Role/assignment containing any R3_HIGH Permission | R3_HIGH | 1 | `tenant.access.role.approve` or platform equivalent | Exact assignment digest approved; checker scope contains assignment. |
| Role/assignment containing any R4_CRITICAL Permission | R4_CRITICAL | 2 | Two critical checkers; at least one security-designated | Both approve same digest; neither maker/beneficiary. |
| Access-policy or platform security-default activation | R4_CRITICAL | 2 | Two critical/security checkers | Both approve same policy digest. |
| Tenant status restrict/suspend/reactivate/termination transition | R4_CRITICAL | 1 | Independent lifecycle checker | Exact transition/from-version/reason digest. |
| Calendar publication | R3_HIGH | 1 | `tenant.calendar.publish` | Exact submitted calendar digest and impact result. |
| Blocking calendar impact override | R4_CRITICAL | 2 | Calendar impact checker plus Tenant security/operations checker | Exact override scope/reason/digest; non-downgradable blockers remain. |
| Work_Assignment approval when enabled | R3_HIGH | 1 | `tenant.assignment.approve` | Exact assignment digest. |
| Working_Day_OT or Closed_Day_OT authorization | R3_HIGH | 1 | `tenant.overtime.authorize` | Exact plan digest; worker consent remains separately required. |
| Attendance/OT/Work_Record exception decision | R3_HIGH | 1 | Corresponding exception decide Permission | Exact requested effect/evidence digest. |
| Compensation_Profile or salary change | R4_CRITICAL | 1 | `tenant.salary.approve` | Exact salary digest; any amount. |
| Pay_Rule version | R4_CRITICAL | 1 | `tenant.pay_rule.approve` | Exact constrained-rule digest; production compliance gate unaffected. |
| Adjustment | R4_CRITICAL | 1 | `tenant.adjustment.approve` | Exact digest; any amount/sign/formula. |
| Payroll ledger review | R3_HIGH | 1 reviewer decision | `tenant.payroll.review` | Current Input_Digest; source maker cannot be sole reviewer. |
| Payroll finalization | R4_CRITICAL | 1 distinct finalizer/checker | `tenant.payroll.finalize` | Submitted/readiness/current Input_Digests equal; all included current reviews; Step-up 2; lock acquired. |
| Payroll/Payslip export | R3_HIGH/R4_CRITICAL | 1 | Separate authorized checker with equivalent export Permission and containing scope | Exact manifest/query/fields/recipient/expiry digest. |
| Evidence/audit export | R4_CRITICAL | 1 | Separate authorized checker with same export class and containing scope | Exact manifest/query/fields/recipient/expiry digest. |
| Support_Grant | R3_HIGH/R4_CRITICAL by data class | 1 Tenant approver | Tenant approver authorized for every requested action/data class | Exact case/purpose/scope/actions/data classes/expiry digest. |
| Break_Glass_Session | R4_CRITICAL | 2 | Security approver and incident approver | Section 10 conditions and same request digest. |

Rejection requires a controlled reason and returns the resource to the policy-defined draft/rejected state without mutating approved content. Approval never authorizes beyond the checker's current Resource_Scope and never survives checker/Role/scope revocation for an action not yet activated.

## 9. Step-up action catalog

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

`STEP_UP_1` is mandatory for: high-risk Role assignment approval/revocation; access review decision; full protected-identifier view; high-sensitivity evidence view; salary/ledger view; Payroll/Payslip export request or approval at R3; support session request/start/use; own Payslip view; release change; and platform audit export.

`STEP_UP_2` is mandatory for: any R4 Role/policy assignment approval or activation; salary approval; Adjustment approval; Pay_Rule approval; payroll finalization submission and finalization; another Worker's Payslip view; Payslip/evidence/audit export; Integration_Principal creation/rotation; Tenant suspension/termination transition; platform security-default activation; calendar blocking-impact override; all Break_Glass request/approval/start actions; and support access to protected ID, sensitive evidence, salary, Payslip, or payroll mutation.

Where an action appears at both levels, `STEP_UP_2` wins. Viewing a button, holding a Role, or having completed a prior unrelated step-up never satisfies action-bound step-up.

## 10. Break_Glass_Session governance

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

### 10.1 Preconditions and request

A Break_Glass_Session is available only for a declared active security or service incident where delay would create material harm and ordinary Support_Grant timing is inadequate. It requires:

- one open Support_Case linked to one Tenant and one active incident record;
- a named Platform_Break_Glass_Operator;
- emergency reason code plus controlled rationale;
- exact requested actions, Resource_Scope, data classes, and expected outcome;
- requested duration of 30 minutes or less;
- canonical request digest and `STEP_UP_2` by the requester;
- two approvals of the same digest: one `platform.break_glass.approve.security` and one `platform.break_glass.approve.incident`.

Requester/operator, security approver, and incident approver are three distinct humans. An approver cannot be the operator's current delegated session identity or approve an assignment that benefits the approver.

### 10.2 Permitted and prohibited authority

The session grants only the intersection of the named operator's existing support-use eligibility and the exact approved actions/scope/data classes. It is read-only unless an explicit non-payroll operational mutation is named and independently allowed by Tenant domain policy.

Break-glass can never grant or perform:

- Payroll_Snapshot or published Payslip mutation;
- payroll finalization or readiness bypass;
- salary, compensation, Pay_Rule, or Adjustment edit/approval;
- Role, Permission, access-policy, Support_Grant, or Tenant security-policy change;
- export of salary, Payslip, evidence, payroll, or audit packages;
- audit alteration/deletion;
- Tenant boundary substitution;
- worker consent, acknowledgement, Supervisor verification, HR review, or checker action on behalf of another human;
- disabling field masking, Maker_Checker_Policy, Segregation_of_Duties, immutable-domain, or audit controls.

Salary, Payslip, protected-identifier, or high-sensitivity evidence **read** is denied unless that exact data class and resource scope appears in both approvals and the operator already has the corresponding support-use eligibility Permission. Every such object access requires `STEP_UP_2` and an access Audit_Event.

### 10.3 Session lifecycle

- The session starts only after both approvals, while request/case/incident/approver authority remain current.
- Effective Tenant is derived from the approved request; caller-supplied substitution is denied without existence disclosure.
- The session expires at the earliest of 30 minutes after start, approved expiry, case closure, incident closure, grant/session revocation, operator session expiry, Role/scope loss, Tenant policy denial, or security reset.
- Extension is forbidden. Continued need requires a new request, digest, Step-up 2, and two new approvals.
- Start immediately alerts configured Platform security, incident response, and Tenant security contacts with purpose, operator, scope, start, and expiry but no protected payload.
- A visible support-context indicator is active for authorized Tenant users for the entire session.
- Every attempted action, including denial, records dual attribution to operator, Tenant, Support_Case, incident, request, approvals, session, purpose, resource/version, policy version, and outcome.
- Revocation is immediate for subsequent requests and emits session/device/client purge signals where applicable.

### 10.4 Post-use review

Within 24 hours after end, an independent `PLATFORM_BREAK_GLASS_REVIEWER` who did not request, approve, or use the session reviews the exact action log, scope adherence, accessed data classes, mutations, denials, alerts, end condition, and stated outcome. The Tenant Security Owner receives the review package and may record a response. Review result is `COMPLIANT`, `EXCESS_SCOPE_ATTEMPT_BLOCKED`, or `POLICY_BREACH`; findings require a reason and linked incident/remediation record. An overdue review triggers a security alert and blocks the operator from starting another Break_Glass_Session until review completion, except a different operator may act under a separately approved new incident request.

## 11. Access review, revocation, and audit

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

Current R3_HIGH/R4_CRITICAL grants, support eligibility, custom Roles, expiring assignments, dormant Accounts, Integration_Principals, and SoD conflicts are included in access review. Decisions are approve-current, narrow, revoke, or reject-renewal; broadening always uses a new assignment workflow. Reviewers cannot review their own access.

Role, Permission, scope, Account, session, device, Tenant status, Support_Session, or Break_Glass revocation applies to every subsequent server decision. Sensitive mutations fail closed if the required authorization Audit_Event cannot be committed atomically or safely quarantined under an approved atomicity policy.

Audit records for this policy include actor boundary, Tenant, Principal, action, resource/version, scope, decision, policy version, maker/checker identities, exact digest, reason, step-up level/evidence reference, support/break-glass context, correlation ID, timestamp, and outcome. Logs and metrics exclude salary amounts, raw identifiers, exact locations, push endpoints, and evidence payloads.

## 12. Configuration persistence and Admin Panel contract

> **DRAFT – NOT APPROVED – NOT FOR PRODUCTION**

Every value in Sections 3–11 is configuration data. Implementations MUST NOT hardcode Permission grants, Role contents, risk mappings, approval counts, SoD pairs, step-up freshness, support limits, Break_Glass duration, review intervals, assignment lifetimes, or prohibited-action lists in application code. The values in this document are proposed seed rows for `ACCESS-GOVERNANCE-BASELINE` version `1.0.0-draft`; the authorization engine resolves the current approved version from authoritative configuration tables.

### 12.1 Configuration tables

Names are logical and technology-neutral. Every table is Tenant-scoped unless `boundary = PLATFORM`; every primary/foreign reference includes the same boundary and Tenant key. All versioned rows include `policy_set_id`, `version`, `status`, `effective_interval`, `created_by`, `created_at`, `reason`, and `content_digest` even where abbreviated below.

| Configuration table | Required policy values | Integrity and authorization constraints |
|---|---|---|
| `access_policy_set` | Boundary, Tenant/cohort scope, semantic version, status (`DRAFT`, `PENDING_APPROVAL`, `APPROVED`, `ACTIVE`, `SUPERSEDED`, `REJECTED`), effective interval, baseline parent, digest | One active applicable version per boundary/Tenant/time; draft cannot authorize; production activation requires explicit approval evidence and open rollout gate. |
| `permission_catalog_entry` | Permission ID, action, resource class, data classification, max scope dimensions, risk tier, human/machine eligibility, constraints | Permission ID unique within catalog version; unknown/wildcard ID invalid; published row immutable. |
| `role_template_version` | Template ID, boundary, version, purpose, status, assignment-scope requirements | No title-derived grants; published version immutable. |
| `role_template_permission` | Role template version, Permission ID | Same-boundary foreign keys; exact catalog version; no implicit inheritance. |
| `custom_role_version` | Tenant, Role ID/version, purpose, maker, status, digest | Tenant-only; published version immutable; validation against approved catalog/guardrails. |
| `custom_role_permission` | Custom Role version, Permission ID | Exact approved catalog reference; no wildcard or cross-boundary grant. |
| `role_assignment_policy` | Risk tier, approval count, checker qualification, max duration, review interval | Tenant may strengthen only; no self/beneficiary approval. |
| `role_assignment` | Principal, Role version, Resource_Scope reference, effective interval, maker, status, digest | Same Tenant; no scope beyond delegation ceiling; inactive until all approvals. |
| `risk_tier_definition` | Tier ID/order, default Step-up level, audit class, assignment/review controls | Tier ordering immutable within active version; Tenant override cannot lower baseline. |
| `action_risk_binding` | Permission/action/resource/workflow/data-class condition, resulting tier | Most restrictive matching tier wins; ambiguity fails closed. |
| `maker_checker_rule` | Action/resource condition, required approvals, checker Permissions, distinctness/beneficiary rules, digest fields, activation rule | Counts are positive; exact policy/version/digest binding; Tenant can increase only. |
| `segregation_of_duties_rule` | Rule ID, maker action set, checker action set, scope-overlap semantics, enforcement (`ASSIGNMENT_DENY`, `ACTION_DENY`), exceptions | Mandatory rules cannot be disabled; no exception may bypass Tenant isolation, immutable objects, or self-approval ban. |
| `step_up_policy` | Level, allowed assurance methods, freshness duration, one-time flag, binding fields, invalidation events | Approved identity capabilities only; stronger method/fresher result wins. |
| `support_access_policy` | Allowed support actions/data classes, max duration, Tenant approver qualification, masking, visibility/notification requirements | Support grant remains intersectional; cannot create Tenant Role or bypass domain controls. |
| `break_glass_policy` | Eligibility, incident/case prerequisites, max duration, approver count/types, distinctness, allowed/prohibited actions/data classes, alert recipients, review deadline, overdue consequence | Mandatory prohibitions in Section 10 cannot be removed by Tenant override; extension disabled. |
| `access_review_policy` | Risk tier, review interval, max assignment duration, triggering events, reviewer qualification, outcomes | Reviewer cannot review own access; expiry/renewal creates new decision. |
| `policy_approval_record` | Policy/assignment/action version, digest, decision, maker, checker, checker Permission/scope, Step-up reference, reason, timestamp | Append-only first-class record; changed digest invalidates applicability. |
| `policy_activation_record` | Approved policy version/digest, activation scope/time, rollout-gate evidence, activator, prior version | Activation denied while `G-TENANT-ROLLOUT` is closed; supersession preserves history. |

Resource scopes are stored in a normalized `resource_scope` plus dimension rows for Tenant, Organization_Unit, Site, Team, Worker, Work_Assignment, Local_Date_Interval, and Payroll_Period. A serialized client scope is never authoritative.

### 12.2 Admin Panel configurability

All policy values are configurable through an authenticated, accessible Admin Panel workflow backed by versioned server APIs and the tables above:

- the **Platform Admin Console** manages platform catalog guardrails, platform Role templates, platform risk/step-up minima, support eligibility, and Break_Glass policy;
- the **Tenant Management Admin Panel** manages Tenant custom Roles, assignments, stricter risk/step-up mappings, Maker_Checker rules, SoD rules, access-review schedules, support approval policy, and authorized contacts within platform guardrails;
- every screen provides current/effective/historical version, field-level diff, impacted Roles/assignments/actions, validation errors, maker/checker state, digest, effective interval, approval evidence state, and rollback-by-new-version;
- draft, validate, submit, approve/reject, activate, supersede, and revoke are separate server-authorized commands; showing a control grants no authority;
- no direct table edit or client-only configuration can activate policy;
- activation is blocked until explicit business approval evidence for the exact digest is recorded and `G-TENANT-ROLLOUT` is opened by the authorized governance process;
- all changes require idempotency, Expected_Version, safe conflict handling, Step-up where required, first-class approval records, and append-only Audit_Events;
- configuration export/import preserves policy IDs, versions, status, digest, dependencies, and classification but cannot import an `ACTIVE` status or approval evidence as authority.

### 12.3 Configuration resolution and safe failure

At request time, the server resolves the active policy by boundary, Tenant, action time, policy scope, and version. It snapshots the resolved policy version into the authorization decision and approval record. Missing, duplicate, ambiguous, expired, not-yet-effective, unsigned, unapproved, or digest-invalid configuration fails closed with a stable safe code and no mutation. Cache entries are keyed by Tenant/boundary and policy version; activation/revocation invalidates affected caches before a later request may rely on the change.

The draft seed may be loaded only into a non-production `DRAFT` state. While `G-TENANT-ROLLOUT` is closed, the system MUST reject transition of this baseline or derived Tenant policy to `ACTIVE` in a production environment.

## 13. Conformance tests and acceptance criteria

These are normative tests for the future implementation. `PASS` means the observable result exactly matches the expected result and no prohibited data/state change occurs. Test data uses generic identifiers only.

| Test ID | Scenario | Expected result | Requirement/design trace |
|---|---|---|---|
| `AG-001` | Tenant Administrator has workforce permissions but lacks `tenant.salary.view`. | Salary fields omitted; no inference through error/count/export. | R2.5; Design §§4.1, 4.5 |
| `AG-002` | Tenant Administrator lacks `tenant.payroll.finalize`. | Finalization denied regardless of title/surface/entitlement. | R2.6–2.9; Design §§4.2, 4.4 |
| `AG-003` | Platform operator requests salary without active exact Support_Session. | Denied without protected payload or resource-existence disclosure; no mutation. | R2.3–2.4, R4.1–4.9; Design §§4.3–4.5 |
| `AG-004` | Same actor accesses same command through two surfaces. | Same server authorization decision from same policy inputs. | R2.7–2.9; Design §§4.2, 4.4 |
| `AG-005` | Role contains unknown or wildcard Permission. | Role validation rejected. | R6.1; Design §4.4 |
| `AG-006` | Assignment crosses Tenant or exceeds delegation ceiling. | Denied without activation; safe audit. | R6.2, R35.1; Design §§4.4, 7.3 |
| `AG-007` | R3 Role assignment has no checker. | Assignment remains inactive. | R6.3–6.5; Design §4.4 |
| `AG-008` | R4 assignment has only one checker or checker is maker/beneficiary. | Assignment remains inactive. | R6.3–6.7; Design §§4.4, 23 P16 |
| `AG-009` | Approved assignment scope or Permission set changes. | Digest changes; prior approvals invalid. | R6.5–6.6; Design §§7.3, 13 |
| `AG-010` | Revoked Role/scope/session is reused from cache. | Next server decision denies. | R6.9; Design §§4.4, 7.3 |
| `AG-011` | Calendar maker attempts to publish own submitted digest. | Denied by `SOD-03-CALENDAR`. | R6.4–6.7; Design §§4.1, 8.1 |
| `AG-012` | OT proposer attempts to authorize own plan. | Denied by `SOD-04-OT`. | R6.4–6.7; Design §§4.1, 8.3 |
| `AG-013` | Salary maker attempts to approve own salary change, including zero amount. | Denied; no amount threshold exemption. | R6.4–6.7, R35.1–35.2; Design §§4.5, 23 P16 |
| `AG-014` | Finalizer is same human as finalization submitter. | Denied; snapshot not created. | R6.4–6.7, R35.2; Design §§9.5, 23 P14/P16 |
| `AG-015` | Checker approves digest then maker changes content. | Prior approval invalid; resubmission required. | R6.5–6.6; Design §§7.3, 13 |
| `AG-016` | Feature flag/entitlement enabled but Permission absent. | Domain action denied. | R2.8–2.9; Design §§4.4, 23 P27 |
| `AG-017` | Support session scope exceeds Support_Grant. | Session/action denied; attempt fully attributed. | R4.2–4.9; Design §§4.3, 23 P20 |
| `AG-018` | Support session has salary data class but agent lacks `platform.support.salary.use`. | Salary masked/denied. | R4.5, R4.9; Design §§4.3–4.5 |
| `AG-019` | Support agent has salary eligibility but grant omits salary. | Salary masked/denied. | R4.5, R4.9; Design §§4.3–4.5 |
| `AG-020` | Support or break-glass attempts Payroll_Snapshot mutation. | Hard denial; snapshot digest/content unchanged. | R4.13; Design §§4.3, 13.3 |
| `AG-021` | R3 action uses stale (>15 minute) or wrong-family step-up. | `REQUIRE_STEP_UP`; no mutation. | R35.2; Design §20.1 |
| `AG-022` | R4 action uses stale (>5 minute), reusable, or wrong-digest step-up. | `REQUIRE_STEP_UP`; no mutation. | R35.2; Design §20.1 |
| `AG-023` | Break-glass request has one approval, same-human approvers, >30 minutes, or changed digest. | Start denied. | R4.10–4.11, R35.2; Design §§4.3, 20.1 |
| `AG-024` | Valid break-glass starts. | Alerts sent, visible indicator active, expiry bounded, dual attribution/audit recorded. | R4.10–4.12; Design §§4.3, 20.4 |
| `AG-025` | Break-glass attempts prohibited role, payroll, export, consent, or audit operation. | Hard denial even if a broad client request was approved. | R4.5, R4.13, R35.1; Design §§4.3–4.5 |
| `AG-026` | Break-glass reaches expiry/revocation. | Every subsequent use denied; no grace reuse. | R4.7–4.8, R4.10–4.12; Design §4.3 |
| `AG-027` | Post-use review is overdue. | Security alert; same operator blocked from new start until review. | R4.12; Design §§4.3, 21.1–21.2 |
| `AG-028` | Custom Role assigns human approval to Integration_Principal. | Role assignment rejected. | R6.1–6.2, R35.1; Design §§4.1, 7.3 |
| `AG-029` | Principal holds maker/checker grants on provably disjoint scopes. | Actions only allowed within the non-overlapping scope; overlap/ambiguity denied. | R6.7; Design §4.4 |
| `AG-030` | Salary view, edit, approve, finalize, export, and support are assigned independently. | Each action succeeds only with its exact Permission and controls; no grant implies another. | R2.5–2.7, R4.9, R35.1–35.2; Design §4.5 |
| `AG-031` | Sensitive export query/fields/recipient changes after approval. | Digest mismatch; export blocked pending new approval. | R6.5–6.6, R35.2; Design §§4.4–4.5 |
| `AG-032` | Role/access policy version is revoked after prior decisions. | Historical decisions retain old policy reference; new actions use current state and deny as applicable. | R6.8–6.11; Design §§4.4, 13 |
| `AG-033` | Required sensitive-action audit write fails. | Mutation blocked or safely quarantined; no unaudited success. | R6.11; Design §§20.4, 24.5 |
| `AG-034` | Caller supplies a different Tenant in support, export, or role assignment. | Denied without existence detail; no state change. | R4.2–4.9, R35.1; Design §§4.3–4.4 |
| `AG-035` | Payroll policy/compliance validation is absent. | This access policy does not open `G-PAYROLL-PRODUCTION`; final production enablement remains blocked. | R35.1–35.2 plus R22.4–22.9 boundary; Register XG-01 |
| `AG-036` | Draft seed is loaded while `G-TENANT-ROLLOUT` is closed. | Policy remains `DRAFT`; production `ACTIVE` transition is denied and audited. | R6.8–6.11, R42.3; Design §§4.4, 28.3 |
| `AG-037` | Tenant Admin attempts to lower a baseline risk, approval count, Step-up level, SoD rule, or Break Glass prohibition. | Validation rejects the draft or activation; active policy unchanged. | R6.3–6.8, R35.1–35.2; Design §§4.4, 20.1 |
| `AG-038` | Caller, import, or direct data change claims `ACTIVE` without exact approval/activation records. | Resolver treats policy as unauthorized, fails closed, and alerts/audits integrity failure. | R6.8–6.11, R35.1; Design §§4.4, 20.4 |
| `AG-039` | Two applicable policy versions overlap ambiguously or cache holds a superseded version. | Authorization fails closed or reloads the single current version before decision; no stale allow. | R6.8–6.9; Design §§4.4, 13 |
| `AG-040` | Admin Panel hides/shows a policy control without corresponding server Permission. | Server decision is unchanged; unauthorized command denied. | R2.7–2.9, R35.1; Design §§4.2, 4.4 |
| `AG-041` | A policy value is absent from configuration and application code supplies a built-in fallback. | Conformance fails; authorization must fail closed rather than use hardcoded policy. | R6.8–6.9, R35.1; Design §§4.4, 7.3 |

### 13.1 Structural validation assertions

A conforming policy artifact and implementation MUST prove:

- every Permission ID is unique and belongs to exactly one boundary;
- every default Role grant references a catalog Permission;
- no default Role includes a prohibited maker/checker pair;
- every R3_HIGH/R4_CRITICAL Permission is covered by the step-up and/or explicit exception rules in this document;
- the six salary/support boundaries (`tenant.salary.view`, `tenant.salary.edit`, `tenant.salary.approve`, `tenant.payroll.finalize`, sensitive export Permissions, and support-use Permissions) remain distinct;
- every governed action records exact digest and policy version;
- Break_Glass duration, approver count, prohibited actions, expiry, alerting, and review are determinate;
- every policy value is represented in a versioned configuration table and manageable through the applicable Admin Panel workflow;
- no hardcoded fallback can authorize when configuration is missing, invalid, ambiguous, unapproved, or inactive;
- production activation is impossible while `G-TENANT-ROLLOUT` is closed;
- no Permission or Role opens `G-PAYROLL-PRODUCTION` or assigns legal meaning to payroll policy.

## 14. Requirement → Design → Artifact → Test traceability

| Requirement | Design source | Artifact sections | Tests |
|---|---|---|---|
| 2.5 | §§4.1, 4.5 | §§4.5, 5.3, 6 | AG-001, AG-030 |
| 2.6 | §§4.1–4.2, 4.5 | §§4.5, 5.3 | AG-002, AG-014, AG-030 |
| 2.7 | §§4.2, 4.4 | §§2, 5.1 | AG-004, AG-016 |
| 4.1–4.3 | §4.3 | §§2.2, 4.2, 8.2 | AG-003, AG-017 |
| 4.4–4.6 | §§4.3, 20.4 | §§2.2, 10.3, 11 | AG-018, AG-024 |
| 4.7–4.9 | §§4.3–4.5 | §§2.2, 4.2, 10.2–10.3 | AG-017–AG-020, AG-026 |
| 4.10–4.12 | §§4.3, 20.1, 20.4 | §10 | AG-023–AG-027 |
| 4.13 | §§4.3, 13.3 | §§6, 10.2 | AG-020, AG-025 |
| 6.1–6.2 | §§4.4, 7.2–7.3 | §§4–6 | AG-005–AG-006, AG-028 |
| 6.3 | §§4.4, 10 | §§6, 8 | AG-007–AG-009 |
| 6.4–6.7 | §§4.1, 4.4, 23 Property 16 | §§7–8 | AG-011–AG-015, AG-029, AG-031 |
| 6.8 | §§4.4, 13 | §§5.1, 6, 8.1, 11 | AG-032 |
| 6.9–6.10 | §§4.4, 7.3 | §§6, 11 | AG-010, AG-032 |
| 6.11 | §20.4 | §§8.1, 11 | AG-009, AG-033 |
| 35.1 | §§4.4–4.5, 20.1 | §§2–7, 10 | AG-002–AG-006, AG-016–AG-020, AG-025, AG-028, AG-034 |
| 35.2 | §§4.3, 20.1 | §§3, 8–10 | AG-008, AG-013–AG-015, AG-021–AG-024, AG-031, AG-037 |
| 42.3 | §28 item 3 | §§1, 3–12, 15 | AG-030, AG-035–AG-041 |

## 15. Approval and change control

To approve this baseline, authorized approvers must create a successor policy record containing the canonical digest of the exact approved content, approval identities and scopes, evidence references, effective date, Tenant/cohort scope, review date, and any stricter overrides. Approval must also update `XG-02`; it must not alter `XG-01` or open `G-PAYROLL-PRODUCTION`.

Any change to a Permission, risk tier, Role template, custom-role constraint, approval count, SoD rule, step-up level/freshness, Break_Glass limit, prohibited action, or review requirement creates a new semantic version and requires impact review, migration/rollback plan, updated tests, and new approval. Historical authorization and approval records retain the exact policy version used.
