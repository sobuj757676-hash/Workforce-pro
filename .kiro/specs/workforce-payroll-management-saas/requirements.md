# Requirements Document

## Introduction

This document defines the business and system requirements for an advanced multi-tenant workforce payroll management SaaS. The System digitizes the Monthly Individual Work Record Card and Employee Overtime Request Form while adding explicit authorization states, evidence provenance, disputes, version history, payroll traceability, and safe mobile workflows. The requirements preserve the source artifacts’ monthly dates 1–31, start and end times, worker and supervisor signatures, remarks, normal days, overtime, incentives, advances, total amount, planned overtime dates and times, date-specific pre-work signatures, and staggered rest-day acknowledgement without copying sample personal data.

The System separates the Platform Admin Console from the Tenant Management Admin Panel, protects tenant payroll and salary data from implicit platform access, and applies permission-based roles with configurable maker-checker and segregation-of-duties controls. A published Working/Closed Operating Calendar is the ordinary-work gate. Working-day planned overtime and exceptional Closed-day work use distinct, bounded authorization and worker-consent workflows. Attendance, replay-resistant checkout evidence, supervisor verification, worker acknowledgement, effective-dated compensation, HR review, and immutable payroll finalization create a traceable work-to-payroll chain.

Payroll formulas, eligibility rules, rates, caps, rounding, reports, retention, and statutory interpretations remain configurable and require validation by qualified Singapore PTE, MOM, contractual, payroll, legal, and compliance stakeholders before tenant production use. This document does not prescribe unverified legal rules or select still-open technology, vendor, commercial, browser-baseline, service-level, or numeric performance choices.

## Glossary

- **System**: The complete Workforce Payroll Management SaaS, including all user applications, APIs, domain services, background processing, storage, and integration channels.
- **Tenant**: A customer organization whose identities, configuration, workforce, evidence, compensation, and payroll data form an isolated data boundary.
- **Control_Plane**: Platform-managed tenant lifecycle, provisioning, release, security-default, usage, quota, service-health, and support-case capabilities that exclude tenant payroll payloads by default.
- **Tenant_Data_Plane**: Tenant-scoped workforce, attendance, overtime, work-record, compensation, payroll, evidence, and business-configuration capabilities.
- **Platform_Admin_Console**: The responsive web surface used by authorized Platform_Operators to administer the Control_Plane.
- **Tenant_Management_Admin_Panel**: The responsive web surface used by permissioned Tenant_Administrators to administer one Tenant.
- **Worker_PWA**: The installable and browser-capable progressive web application for a Worker’s own attendance, consent, checkout, Work_Record, dispute, estimate, and Payslip workflows.
- **Supervisor_PWA**: The installable and browser-capable progressive web application for an assigned Supervisor’s roster, Toolbox_Meeting, attendance, exception, and Work_Record workflows.
- **Management_Web**: The responsive web surface for permissioned management, HR, and payroll operations.
- **Principal**: An authenticated human or machine identity.
- **Platform_Operator**: A Principal authorized for specified Control_Plane actions without an implicit Tenant business role.
- **Platform_Support_Agent**: A Platform_Operator authorized to coordinate Support_Cases and request bounded Support_Sessions.
- **Tenant_Administrator**: A Tenant Principal with one or more explicit Tenant administration permissions; the designation does not grant salary access or payroll finalization.
- **Worker**: A Tenant workforce member who may act only on the Worker’s own authorized records.
- **Supervisor**: A Tenant Principal assigned to specified Sites, teams, Workers, dates, or scopes.
- **HR_Reviewer**: A Tenant Principal permitted to review Work_Records and Payroll_Ledger_Entries within an assigned scope.
- **Payroll_Maker**: A Tenant Principal permitted to prepare compensation, adjustment, calculation, or payroll changes.
- **Payroll_Finalizer**: A Tenant Principal permitted to approve readiness and finalize payroll subject to Maker_Checker_Policy.
- **Auditor**: A read-only Principal permitted to inspect authorized reports, evidence, and Audit_Events.
- **Integration_Principal**: A machine Principal with versioned, Tenant-scoped, action-scoped, and data-classification-scoped credentials.
- **Permission**: An explicit authorization to perform one action on one resource class.
- **Resource_Scope**: The Tenant, organization unit, Site, team, Worker, assignment, date, or Payroll_Period boundary attached to a Permission.
- **Role**: A versioned collection of Permissions that can be assigned to a Principal for a Resource_Scope.
- **Maker_Checker_Policy**: A configurable policy that requires one Principal to prepare an action and a distinct qualified Principal to approve the exact submitted version.
- **Segregation_of_Duties**: A policy that prevents prohibited combinations of preparation, approval, access, or finalization actions.
- **Step_Up_Authentication**: A fresh, higher-assurance authentication check required for a sensitive action.
- **Organization_Unit**: A Tenant-configured legal entity, business unit, department, or other controlled organizational scope.
- **Site**: A Tenant-configured work location with an IANA time zone and optional Checkpoints.
- **Checkpoint**: An authorized physical or logical checkout point associated with a Site.
- **Account**: A Tenant-managed login identity and lifecycle record for a human Principal.
- **Worker_Profile**: A protected workforce identity record that may contain a masked employment identifier.
- **Employment_Record**: A versioned record describing a Worker’s employment relationship.
- **Work_Assignment**: An effective-dated relationship among a Worker, Site, team, Supervisor, and permitted work scope.
- **Operating_Calendar**: A versioned monthly schedule that classifies every applicable date as Working or Closed.
- **Calendar_Day**: One date and scope in a published Operating_Calendar.
- **Working**: A Calendar_Day classification that permits ordinary-work workflows subject to all other controls.
- **Closed**: A Calendar_Day classification that blocks ordinary-work workflows and remains unchanged by exceptional work authorization.
- **Toolbox_Meeting**: A Supervisor-created, timed meeting with a versioned roster snapshot that initiates attendance for an eligible work scope.
- **Dual_Attendance_Confirmation**: Matching Supervisor and Worker arrival confirmations for the same Toolbox_Meeting and Worker.
- **Attendance_Exception**: A reason-coded state for lateness, absence, mismatch, missing confirmation, or another attendance anomaly.
- **Overtime_Plan**: A proposed set of named Workers, dates, planned times, Site or location scope, work category, and purpose for overtime.
- **Overtime_Assignment**: One named Worker’s bounded, versioned overtime authorization derived from an approved Overtime_Plan.
- **Working_Day_OT**: Planned overtime on a Working Calendar_Day that is distinct from ordinary work.
- **Closed_Day_OT**: Exceptional authorized overtime on a Closed Calendar_Day.
- **Authorized_Work_Window**: The approved date, start, end, Site, location, Supervisor, Worker, work type, and authorization version that bounds overtime work.
- **Worker_PreWork_Consent**: A Worker’s date-specific decision, recorded before work starts, that binds the Worker to the exact Overtime_Assignment digest.
- **Staggered_Rest_Day_Acknowledgement**: A configurable consent or acknowledgement subtype whose label and legal meaning require Tenant policy validation.
- **QR_Challenge**: A signed, short-lived, purpose-bound, one-time checkout challenge containing Tenant, Site, Checkpoint, issue time, expiry, nonce, and key reference.
- **Checkout_Evidence**: Minimized timestamp, Checkpoint, location, device, session, and challenge evidence submitted for a Work_Record.
- **Work_Record**: A versioned daily record of a Worker’s date, Site, time segments, classifications, remarks, evidence, verification, acknowledgement, and workflow status.
- **Work_Record_Digest**: A cryptographic digest of the canonical contents of one Work_Record version.
- **Supervisor_Verification**: A Supervisor’s decision on an exact Work_Record version, including reason-coded edits where applicable.
- **Worker_Acknowledgement**: A Worker’s signature over one current Work_Record version and Work_Record_Digest.
- **Batch_Envelope**: A Worker-signed collection of ordered per-record digests used to create independently traceable Worker_Acknowledgements.
- **Dispute**: A Worker’s reasoned rejection of one current Work_Record version.
- **Monthly_Work_Record_Card**: A responsive, physical-card-inspired monthly projection with one row or card for each calendar date and recognizable source-artifact fields.
- **Compensation_Profile**: An approved, effective-dated, versioned set of pay basis, currency, compensation values, and rule references for a Worker.
- **Pay_Rule**: A versioned, approved, deterministic rule for normal pay, overtime, allowance, incentive, advance, deduction, adjustment, cap, or rounding.
- **Adjustment**: A versioned, reason-coded payroll input prepared and approved under configured policy.
- **Payroll_Period**: A Tenant-configured interval for calculation, review, readiness, finalization, and Payslip publication.
- **Payroll_Ledger_Entry**: A Worker-level period result containing calculation lines, states, warnings, blockers, review decisions, and provenance.
- **Estimated_Not_Finalized**: A payroll status indicating that displayed amounts can change and are not a finalized payment result.
- **Input_Digest**: A digest of the exact canonical calendar, Work_Record, Compensation_Profile, Pay_Rule, Adjustment, and engine versions used in a calculation.
- **Stale_Review**: A prior HR or payroll review whose reviewed Input_Digest no longer matches current inputs.
- **Readiness_Report**: A versioned result of finalization checks, including blockers, warnings, and the checked Input_Digest.
- **Payroll_Snapshot**: An immutable, versioned final payroll artifact that locks exact Worker results and all input references.
- **Payslip**: A finalized Worker document and machine-readable representation derived from a Payroll_Snapshot.
- **Support_Case**: A tracked request describing a Tenant support purpose without granting Tenant_Data_Plane access.
- **Support_Grant**: A Tenant-approved authorization containing exact purpose, scope, data classes, approver, and expiration for a Support_Session.
- **Support_Session**: A visible, revocable, purpose-bound, scope-limited, time-bound session that attributes Tenant_Data_Plane access to both a Platform_Support_Agent and a Support_Grant.
- **Break_Glass_Session**: A separately governed, short-lived emergency Support_Session requiring elevated approval, alerting, and post-use review.
- **Audit_Event**: An append-only record of an action, authorization decision, version, digest, reason, context, and outcome.
- **Idempotency_Key**: A Tenant-scoped client identifier that makes safe retries converge on one equivalent server effect.
- **Expected_Version**: The aggregate version supplied by a command for optimistic concurrency control.
- **Pending_Offline_Action**: A minimized, encrypted where capability permits, session-bound client action that has no authority until server validation.
- **Capability_Fallback**: A documented browser or process alternative used when installation, push, background synchronization, camera, geolocation, or durable storage is unavailable.
- **Sensitive_Data**: Salary, monetary estimates, Payslips, protected identifiers, signatures, precise location, evidence, authentication secrets, and other classified Tenant data.
- **Published_Object**: A versioned domain object whose content is immutable after publication.
- **Integration_Event**: A versioned, minimized fact emitted for authorized asynchronous processing or external integration.
- **Integrity_Manifest**: A list of export artifacts, source versions, classifications, and digests used to verify completeness and integrity.
- **Qualified_Stakeholder**: A person authorized and competent to validate Singapore PTE, MOM, contract, payroll, legal, privacy, accessibility, security, or operational policy.

## Requirements

### Requirement 1: Tenant Isolation and Authorization Boundary

**User Story:** As a Tenant security owner, I want every action and data access isolated by Tenant and scope, so that one Tenant cannot access another Tenant’s information or authority.

#### Acceptance Criteria

1. WHEN a Principal submits a Tenant_Data_Plane request, THE System SHALL derive the effective Tenant from the authenticated server-side context.
2. WHEN a request names a resource, THE System SHALL verify that the resource belongs to the effective Tenant before returning resource data.
3. IF a Principal’s effective Tenant differs from a resource Tenant, THEN THE System SHALL deny the request without disclosing the resource’s existence.
4. IF a Tenant boundary check fails, THEN THE System SHALL leave authoritative domain state unchanged.
5. THE System SHALL apply Tenant scope to database records, cache keys, object keys, search documents, messages, idempotency records, exports, and Integration_Events.
6. WHEN a caller supplies a Tenant identifier, THE System SHALL compare the supplied identifier with the authenticated Tenant context.
7. WHEN the System evaluates a command, THE System SHALL evaluate Permission, Resource_Scope, workflow state, applicable policy version, and Tenant lifecycle status.
8. IF any mandatory authorization condition denies a command, THEN THE System SHALL prevent a later allow condition from overriding the denial.
9. WHEN the same Principal attempts the same action through different user surfaces or APIs, THE System SHALL produce authorization decisions from the same server-side policy inputs.
10. IF a user interface displays an action that the Principal lacks permission to execute, THEN THE System SHALL deny execution at the server boundary.

### Requirement 2: Platform and Tenant Administration Separation

**User Story:** As a SaaS governance owner, I want distinct platform and Tenant administration boundaries, so that platform operations cannot become hidden Tenant payroll authority.

#### Acceptance Criteria

1. THE System SHALL provide a Platform_Admin_Console for Control_Plane administration.
2. THE System SHALL provide a separate Tenant_Management_Admin_Panel for Tenant-owned administration.
3. WHEN a Platform_Operator performs a Control_Plane action, THE System SHALL use a Platform Principal that has no implicit Tenant_Data_Plane role.
4. IF a Platform_Operator requests salary, payroll, Payslip, work-evidence, or protected identity data without a valid Support_Session, THEN THE System SHALL deny the request without returning the protected data.
5. IF a Tenant_Administrator lacks a salary Permission, THEN THE System SHALL omit salary fields from responses available to that Tenant_Administrator.
6. IF a Tenant_Administrator lacks payroll-finalization Permission, THEN THE System SHALL deny payroll finalization regardless of the Tenant_Administrator designation.
7. WHEN an actor uses a surface outside the actor’s primary role, THE System SHALL require an independently granted Role and Resource_Scope.
8. THE System SHALL prevent feature entitlements and feature flags from granting domain Permissions.
9. THE System SHALL prevent feature entitlements and feature flags from bypassing Maker_Checker_Policy, readiness checks, or Payroll_Snapshot immutability.

### Requirement 3: Tenant Lifecycle and Provisioning

**User Story:** As a Platform_Operator, I want auditable and resumable Tenant lifecycle administration, so that customer environments are provisioned and controlled safely.

#### Acceptance Criteria

1. WHEN an authorized Platform_Operator requests Tenant provisioning with an Idempotency_Key, THE System SHALL create or return one Tenant identity for that key and request digest.
2. WHEN a provisioning step completes, THE System SHALL record the step and Tenant placement before executing a later dependent step.
3. IF provisioning is retried after interruption, THEN THE System SHALL resume from recorded completion state without duplicating completed Tenant resources.
4. IF mandatory identity, security, ownership, data-placement, or onboarding checks are incomplete, THEN THE System SHALL prevent the Tenant from entering Active status.
5. WHEN a Tenant lifecycle status changes, THE System SHALL record the prior status, new status, reason, actor, expected version, and Audit_Event.
6. WHILE a Tenant is Suspended, THE System SHALL apply the configured command-specific suspension policy to new sessions, mutations, queues, and retained-record access.
7. WHEN a Tenant is Suspended, THE System SHALL preserve Tenant data according to the approved retention and lifecycle policy.
8. WHERE worker access to already-published Payslips during suspension is enabled by approved policy, THE System SHALL permit only the configured retained-record access.
9. IF a Tenant status transition conflicts with the current Expected_Version, THEN THE System SHALL return a conflict without applying the transition.
10. WHERE plans, billing, entitlements, or quotas are implemented, THE System SHALL version the applicable Control_Plane policy and record the version used for each decision.
11. IF commercial rules for plans, billing, trials, grace periods, entitlements, or quotas remain unapproved, THEN THE System SHALL treat those rules as rollout prerequisites rather than infer values.

### Requirement 4: Tenant-Approved Support and Break-Glass Access

**User Story:** As a Tenant security owner, I want support access to be explicit, bounded, visible, and attributable, so that platform support cannot silently impersonate Tenant users.

#### Acceptance Criteria

1. WHEN a Platform_Support_Agent requests Tenant_Data_Plane access, THE System SHALL require an open Support_Case and a requested purpose, scope, data classification, and duration.
2. WHEN a Tenant approver approves support access, THE System SHALL bind the Support_Grant to the Support_Case, Tenant, exact approved scope, purpose, approver, and expiration.
3. IF a requested Support_Session scope exceeds the Support_Grant, THEN THE System SHALL deny creation of the Support_Session.
4. WHEN a Support_Session starts, THE System SHALL provide a visible support-context indication to authorized Tenant users.
5. WHILE a Support_Session is active, THE System SHALL enforce Tenant domain Permissions, field masking, Maker_Checker_Policy, and immutable-domain rules.
6. WHEN a Support_Session action occurs, THE System SHALL attribute the action to the Platform_Support_Agent, Tenant, Support_Case, Support_Grant, Support_Session, purpose, and outcome.
7. IF a Support_Session is expired, revoked, for a different Tenant, outside purpose, or outside scope, THEN THE System SHALL deny the attempted action.
8. WHEN an authorized Tenant or platform security Principal revokes a Support_Session, THE System SHALL prevent subsequent use of that session.
9. IF salary, Payslip, payroll mutation, or high-sensitivity evidence access is absent from the exact approved scope, THEN THE System SHALL mask or deny that data class.
10. WHEN a Break_Glass_Session is requested, THE System SHALL require the separately approved emergency policy, elevated authorization, bounded duration, and stated emergency reason.
11. WHEN a Break_Glass_Session starts, THE System SHALL alert configured security contacts.
12. WHEN a Break_Glass_Session ends, THE System SHALL require and record a post-use review.
13. IF a Support_Session or Break_Glass_Session attempts to modify a Payroll_Snapshot, THEN THE System SHALL reject the modification.

### Requirement 5: Organization, Site, Account, Worker, and Assignment Administration

**User Story:** As a Tenant_Administrator, I want permission-scoped workforce administration, so that operational structures and responsibilities are current and traceable.

#### Acceptance Criteria

1. WHEN an authorized Tenant_Administrator creates or updates an Organization_Unit, THE System SHALL store the Tenant, type, name, status, version, effective information, and Audit_Event.
2. WHEN an authorized Tenant_Administrator creates or updates a Site, THE System SHALL require an IANA time zone and Tenant scope.
3. WHEN an authorized Tenant_Administrator creates or updates a Checkpoint, THE System SHALL bind the Checkpoint to one Tenant and Site.
4. WHEN an authorized Tenant_Administrator creates or changes an Account, THE System SHALL record lifecycle status, assigned Roles, Resource_Scopes, Expected_Version, and Audit_Event.
5. WHEN an authorized Principal creates or updates a Worker_Profile, THE System SHALL enforce field-level authorization for protected employment identifiers.
6. WHEN the System displays a protected employment identifier to a Principal without full-view Permission, THE System SHALL return the configured masked representation.
7. WHEN an authorized Principal creates an Employment_Record, THE System SHALL record its effective interval, source, version, reason, and approval state.
8. WHEN an authorized Principal creates a Work_Assignment, THE System SHALL bind the Worker, Site, team, Supervisor, effective interval, and Resource_Scope.
9. IF an assignment reference crosses a Tenant boundary, THEN THE System SHALL reject the assignment.
10. WHEN an assignment version changes, THE System SHALL preserve prior versions for historical date resolution.
11. WHEN the System resolves a Worker’s operational scope for a date, THE System SHALL use the approved assignment version effective for that date.
12. IF assignment data is missing or ambiguous for an action that requires assignment authority, THEN THE System SHALL block that action with a stable reason code.

### Requirement 6: Permission-Based Roles, Maker-Checker, and Segregation of Duties

**User Story:** As a Tenant access administrator, I want configurable roles and approval separation, so that sensitive actions follow least privilege and organizational controls.

#### Acceptance Criteria

1. THE System SHALL support versioned role templates and custom Roles composed of explicit Permissions.
2. WHEN a Role is assigned, THE System SHALL bind the Role version to a Principal and Resource_Scope.
3. WHEN a high-risk Role assignment is configured to require approval, THE System SHALL keep the assignment inactive until an authorized checker approves the exact assignment digest.
4. WHEN Maker_Checker_Policy applies to an action, THE System SHALL require the checker to differ from the maker.
5. WHEN a checker approves an action, THE System SHALL bind the approval to the submitted content digest and policy version.
6. IF submitted content changes after approval, THEN THE System SHALL invalidate the prior approval for the changed content.
7. IF a Principal violates a configured Segregation_of_Duties rule, THEN THE System SHALL deny the conflicting action.
8. WHEN an access policy changes, THE System SHALL preserve the policy version used for prior authorization decisions.
9. WHEN a Role, Permission, Resource_Scope, or account status is revoked, THE System SHALL apply the revocation to subsequent server authorization decisions.
10. WHEN configured access-review intervals become due, THE System SHALL present current sensitive grants for authorized review.
11. WHEN a sensitive access grant is approved, rejected, changed, or revoked, THE System SHALL create an Audit_Event.

### Requirement 7: Versioning, Concurrency, and Idempotent Commands

**User Story:** As an operator, I want safe retries and conflict detection, so that network failures and concurrent edits do not duplicate or overwrite business actions.

#### Acceptance Criteria

1. WHEN a mutation is submitted, THE System SHALL require a Tenant-scoped Idempotency_Key.
2. WHEN an Idempotency_Key is first accepted, THE System SHALL store the request digest and result reference for the configured idempotency interval.
3. WHEN the same Idempotency_Key and request digest are retried, THE System SHALL return the existing equivalent result.
4. IF the same Idempotency_Key is reused with a different request digest, THEN THE System SHALL reject the request.
5. WHEN a mutable aggregate command is submitted, THE System SHALL compare the supplied Expected_Version with the current aggregate version.
6. IF an Expected_Version does not match, THEN THE System SHALL reject the write and return the current version plus a safe changed-fields summary.
7. IF a conflict affects content requiring review or signature, THEN THE System SHALL require review of the current version before resubmission.
8. WHEN a Published_Object changes within an allowed workflow, THE System SHALL create a new version rather than overwrite the published content.
9. WHEN the System canonicalizes content for a digest, THE System SHALL use deterministic field ordering and representation.
10. WHEN an API error is returned, THE System SHALL include a stable code, human-safe message, and correlation identifier without cross-Tenant disclosure.

### Requirement 8: Monthly Operating Calendar

**User Story:** As a Management Calendar Maker, I want every date classified and approved, so that ordinary work eligibility is explicit before operations begin.

#### Acceptance Criteria

1. WHEN an authorized calendar maker creates a monthly Operating_Calendar, THE System SHALL create one Calendar_Day for every date in that month and scope.
2. WHEN a Calendar_Day is configured, THE System SHALL require exactly one Working or Closed classification.
3. WHEN a calendar maker submits an Operating_Calendar, THE System SHALL bind submission to the exact calendar version and digest.
4. WHEN Maker_Checker_Policy is enabled for calendar publication, THE System SHALL require an eligible checker other than the maker.
5. WHEN a checker publishes an Operating_Calendar, THE System SHALL make the published calendar version immutable.
6. IF any date lacks a valid classification, THEN THE System SHALL block calendar publication.
7. IF calendar scope, dates, or classifications change after submission, THEN THE System SHALL require resubmission and approval of the new digest.
8. WHEN a later calendar version is proposed, THE System SHALL calculate impact on existing meetings, overtime, Work_Records, reviews, and payroll inputs.
9. IF the configured impact policy classifies a calendar change as blocking, THEN THE System SHALL prevent publication until the blocker is resolved or an authorized policy-defined override is approved.
10. WHEN a later calendar version is published, THE System SHALL preserve prior published versions for historical resolution.
11. WHEN the System resolves work eligibility, THE System SHALL use the published Calendar_Day version applicable to the Site, scope, and date.
12. IF no published Calendar_Day applies to an operational action, THEN THE System SHALL block the action with a stable missing-calendar reason.

### Requirement 9: Closed-Day Hard Gate

**User Story:** As a workforce governance owner, I want Closed dates to block ordinary work, so that exceptional work cannot be recorded through convenience paths.

#### Acceptance Criteria

1. WHEN an ordinary Toolbox_Meeting is requested for a Closed Calendar_Day, THE System SHALL reject the request with a Closed-calendar reason.
2. WHEN ordinary attendance is requested for a Closed Calendar_Day, THE System SHALL reject the request with a Closed-calendar reason.
3. WHEN an ordinary Work_Record is requested for a Closed Calendar_Day, THE System SHALL reject the request with a Closed-calendar reason.
4. WHEN ordinary overtime entry is requested for a Closed Calendar_Day without Closed_Day_OT authorization, THE System SHALL reject the request with a Closed-calendar reason.
5. WHEN Closed_Day_OT is authorized, THE System SHALL preserve the Calendar_Day classification as Closed.
6. WHEN a Work_Record results from valid Closed_Day_OT, THE System SHALL label the record “Closed — Authorized OT Work”.
7. IF a Closed-day action falls outside the applicable Authorized_Work_Window, THEN THE System SHALL reject the action without changing the Calendar_Day.
8. WHEN a Closed-day rejection occurs, THE System SHALL present the authorized proposal or exception path available to the current Principal.

### Requirement 10: Supervisor-Created Toolbox Meetings

**User Story:** As a Supervisor, I want to create scoped toolbox meetings, so that attendance begins from an accountable roster and calendar context.

#### Acceptance Criteria

1. WHEN an assigned Supervisor creates an ordinary Toolbox_Meeting, THE System SHALL require a Working Calendar_Day, Site, date, start time, and authorized assignment scope.
2. WHEN a Toolbox_Meeting is created, THE System SHALL capture the current roster as a versioned roster snapshot.
3. WHEN a Toolbox_Meeting is created, THE System SHALL record the applicable calendar and assignment version references.
4. IF a Supervisor is outside the Site, team, date, or Worker scope, THEN THE System SHALL deny meeting creation for that scope.
5. WHEN an authorized Supervisor creates a Closed-day special Toolbox_Meeting, THE System SHALL require an approved Closed_Day_OT authorization covering the meeting scope and time.
6. IF a special Toolbox_Meeting exceeds the authorized Worker, Supervisor, Site, location, date, or time scope, THEN THE System SHALL reject the meeting.
7. WHEN a meeting is activated, THE System SHALL make attendance actions available only to Workers in the roster snapshot or explicitly authorized exception scope.
8. WHEN a meeting is changed after attendance begins, THE System SHALL create a new auditable version and protect existing attendance provenance.
9. WHEN a Toolbox_Meeting reaches its configured end state, THE System SHALL prevent new ordinary attendance confirmations except through an authorized exception workflow.

### Requirement 11: Dual Attendance Confirmation

**User Story:** As a Worker, I want my arrival confirmed by both my Supervisor and me, so that an ordinary Work_Record begins from mutual attendance evidence.

#### Acceptance Criteria

1. WHEN a Supervisor confirms a Worker’s arrival, THE System SHALL bind the confirmation to the Supervisor, Worker, Toolbox_Meeting, timestamp, and meeting version.
2. WHEN a Worker confirms arrival, THE System SHALL bind the confirmation to the authenticated Worker, attendance entry, timestamp, and meeting version.
3. WHEN matching valid Supervisor and Worker confirmations exist, THE System SHALL create one draft Work_Record for the meeting-Worker pair.
4. WHEN creation of the same draft is retried, THE System SHALL return the existing draft rather than create a duplicate.
5. IF the Worker is not in the authorized roster or exception scope, THEN THE System SHALL reject the Worker confirmation.
6. IF the two confirmations refer to different Workers or meetings, THEN THE System SHALL create an Attendance_Exception rather than a draft Work_Record.
7. IF either required confirmation is absent, THEN THE System SHALL keep attendance in a non-work-record state.
8. WHEN a lateness, absence, or arrival-time mismatch is reported, THE System SHALL create a reason-coded Attendance_Exception.
9. WHEN an Attendance_Exception is resolved, THE System SHALL record the resolver, reason, evidence, before-and-after state, and Audit_Event.
10. IF an attendance action is Pending offline, THEN THE System SHALL display Pending status until server validation completes.

### Requirement 12: Working-Day Planned Overtime

**User Story:** As an OT Proposer, I want Working-day overtime planned and approved separately from ordinary work, so that overtime is intentional and traceable.

#### Acceptance Criteria

1. WHEN an OT Proposer creates a Working_Day_OT plan, THE System SHALL require named Workers, a Working date, planned start and end times, Site or location scope, work category, purpose, and proposed Supervisor scope.
2. WHEN a Working_Day_OT plan is submitted, THE System SHALL validate assignment eligibility and time overlap against current approved data.
3. WHEN Maker_Checker_Policy applies, THE System SHALL require an eligible OT authorizer other than the proposer.
4. WHEN an OT authorizer approves a plan, THE System SHALL create bounded Overtime_Assignments for the named Workers.
5. IF a plan changes after authorization, THEN THE System SHALL require a new authorization digest and new Worker_PreWork_Consent for affected assignments.
6. WHEN planned overtime is displayed to a Worker, THE System SHALL show exact date, planned times, Site or location, work category, consent deadline, and authorization state.
7. WHEN actual Working_Day_OT is recorded, THE System SHALL link the Work_Record to the applicable Overtime_Assignment.
8. IF overtime hours lack an applicable authorized assignment or approved exception, THEN THE System SHALL prevent normal overtime inclusion in payroll readiness.
9. THE System SHALL prevent hours beyond a configured threshold from becoming authorized overtime solely because the threshold was exceeded.
10. WHEN actual overtime differs from planned overtime, THE System SHALL present the variance for Supervisor and payroll review.

### Requirement 13: Closed-Day Named and Bounded Authorization

**User Story:** As an OT Authorizer, I want Closed-day work limited to named people and an exact scope, so that exceptional work remains controlled without changing the calendar.

#### Acceptance Criteria

1. WHEN a Closed_Day_OT proposal is created, THE System SHALL require named Workers, named or scoped Supervisors, date, bounded start and end times, Site, bounded location scope, purpose, and work category.
2. WHEN a Closed_Day_OT proposal is submitted, THE System SHALL verify that the applicable published Calendar_Day is Closed.
3. WHEN a Closed_Day_OT proposal is approved, THE System SHALL create an immutable authorization digest.
4. WHEN a Closed_Day_OT proposal is approved, THE System SHALL create a separate consent request for each named Worker.
5. IF a Worker is not named in the current authorization, THEN THE System SHALL reject creation of that Worker’s Closed-day Work_Record.
6. IF the acting Supervisor is outside the current authorization, THEN THE System SHALL reject the Closed-day attendance or Work_Record action.
7. IF actual work time falls outside the Authorized_Work_Window, THEN THE System SHALL reject normal Closed_Day_OT treatment and route the variance to exception review.
8. IF actual work location falls outside the authorized location scope, THEN THE System SHALL reject normal Closed_Day_OT treatment and route the variance to exception review.
9. WHEN Closed_Day_OT authorization is superseded before work, THE System SHALL invalidate unsigned consent requests tied to the prior digest.
10. WHEN Closed_Day_OT authorization is displayed, THE System SHALL retain the Closed calendar classification and authorized-work label.

### Requirement 14: Date-Specific Worker Pre-Work Consent

**User Story:** As a Worker, I want overtime consent tied to the exact date and assignment before work begins, so that my decision cannot be reused for different work.

#### Acceptance Criteria

1. WHEN a Worker reviews an Overtime_Assignment, THE System SHALL display the exact Worker, date, time window, Site or location, work category, applicable terms, and authorization digest.
2. WHEN a Worker consents, THE System SHALL bind Worker_PreWork_Consent to the authenticated Worker and exact Overtime_Assignment digest.
3. WHEN a Worker consents, THE System SHALL record the decision timestamp before the Authorized_Work_Window start.
4. IF the Authorized_Work_Window has started before consent is accepted, THEN THE System SHALL mark the consent request Expired and prevent synthetic backdating.
5. WHEN a Worker declines an Overtime_Assignment, THE System SHALL record the decline without treating the Worker as consented.
6. IF an authorization digest changes, THEN THE System SHALL require a new Worker_PreWork_Consent for the changed assignment.
7. IF a Worker_PreWork_Consent is missing, late, declined, expired, or tied to another digest, THEN THE System SHALL block normal authorized overtime inclusion.
8. WHERE Staggered_Rest_Day_Acknowledgement is configured, THE System SHALL record the configured label, subtype, policy version, date-specific decision, and assignment digest.
9. WHERE Staggered_Rest_Day_Acknowledgement is configured, THE System SHALL avoid assigning legal meaning beyond the Qualified_Stakeholder-approved Tenant policy.
10. IF consent is captured offline under an approved policy, THEN THE System SHALL keep the consent Pending until the server validates identity, deadline, authorization, digest, Tenant status, and scope.

### Requirement 15: Replay-Resistant Site QR Checkout

**User Story:** As a Tenant operations owner, I want checkout evidence resistant to replay and scope substitution, so that recorded end times have trustworthy provenance.

#### Acceptance Criteria

1. WHEN a Checkpoint issues a QR_Challenge, THE System SHALL bind the challenge to Tenant, Site, Checkpoint, checkout purpose, issue time, expiry, random nonce, and signing-key reference.
2. WHEN a Worker submits checkout, THE System SHALL require an authenticated Worker with an open eligible Work_Record.
3. WHEN checkout is submitted, THE System SHALL verify the QR_Challenge signature, purpose, freshness, Tenant, Site, Checkpoint, Work_Record, session, and Worker binding.
4. WHEN location evidence is required by approved policy, THE System SHALL evaluate the submitted evidence against the configured Checkpoint or location scope.
5. WHEN a valid nonce is accepted, THE System SHALL consume the nonce atomically.
6. IF a nonce was previously consumed, THEN THE System SHALL reject the submission as replay.
7. IF a QR_Challenge is expired, has an invalid signature, or has the wrong purpose, THEN THE System SHALL reject checkout without setting an end time.
8. IF checkout scope does not match the current Work_Record, THEN THE System SHALL reject checkout without attaching evidence to another record.
9. WHEN checkout is accepted, THE System SHALL attach minimized Checkout_Evidence to the correct Work_Record and use a trusted receipt time according to approved policy.
10. IF location evidence fails policy, THEN THE System SHALL route checkout to an explicit exception state rather than fabricate acceptance.
11. WHEN key rotation occurs, THE System SHALL accept only keys in the configured active or grace state.
12. WHEN replay, repeated device sharing, impossible travel, excessive failure, location mismatch, or abnormal timing is detected, THE System SHALL create a security signal without automatically treating the signal as proof of misconduct.
13. THE System SHALL prevent a static reusable QR code from being accepted as equivalent to a valid one-time QR_Challenge.

### Requirement 16: Checkout and Work Evidence Exceptions

**User Story:** As a Worker or Supervisor, I want explicit exception paths for evidence failures, so that legitimate work can be reviewed without weakening normal controls.

#### Acceptance Criteria

1. WHEN a Worker has no supported phone or camera, THE System SHALL offer the configured supervised checkout exception path.
2. WHEN a checkout is missed, THE System SHALL keep the final end time unresolved until an authorized reason-coded correction or exception is approved.
3. WHEN a wrong-location event occurs, THE System SHALL require review of reassignment or supporting evidence before acceptance.
4. WHEN a Worker is reassigned, THE System SHALL require a current assignment version or an approved assignment exception.
5. WHEN emergency work is reported, THE System SHALL route the report through the configured elevated retrospective review workflow.
6. IF emergency work lacks evidence required by approved policy, THEN THE System SHALL prevent the work from being labeled ordinary authorized work.
7. WHEN an exception is requested, THE System SHALL record requester, reason code, optional controlled text, supporting evidence references, and requested effect.
8. WHEN an exception is approved or rejected, THE System SHALL record approver, decision, before-and-after effect, payroll impact, and Audit_Event.
9. IF an exception would exceed the approver’s Resource_Scope, THEN THE System SHALL deny the exception decision.
10. WHEN an exception changes an acknowledged or reviewed source, THE System SHALL create a new Work_Record version and trigger affected re-review.

### Requirement 17: Supervisor Work Record Review and Editing

**User Story:** As a Supervisor, I want to review and reason-code Work_Record changes, so that verified time is accurate and every edit is attributable.

#### Acceptance Criteria

1. WHEN checkout or approved exception evidence completes a draft Work_Record, THE System SHALL route the record to Supervisor review.
2. WHEN a Supervisor reviews a Work_Record, THE System SHALL display current date, Site, time segments, breaks, classifications, remarks, evidence, authorization, consent, and version.
3. WHEN a Supervisor verifies an unchanged Work_Record, THE System SHALL bind Supervisor_Verification to the exact Work_Record version and digest.
4. WHEN a Supervisor changes time, break, classification, or remark, THE System SHALL require a reason code.
5. WHEN a Supervisor changes a Work_Record, THE System SHALL record before-and-after values, optional controlled reason text, actor, timestamp, and evidence references.
6. WHEN a Supervisor change is accepted, THE System SHALL create the required new Work_Record version.
7. IF a Supervisor is outside the Worker, Site, team, date, or authorization scope, THEN THE System SHALL deny verification or editing.
8. IF an Expected_Version conflict occurs during verification, THEN THE System SHALL reject the change and require review of the current version.
9. WHEN a verified record is ready for Worker review, THE System SHALL change the state to Pending Worker Acknowledgement.
10. WHEN an acknowledged or HR-reviewed Work_Record is later corrected, THE System SHALL preserve the prior signed or reviewed version as historical evidence.
11. WHEN an acknowledged or HR-reviewed Work_Record is later corrected, THE System SHALL mark the new version unacknowledged and affected payroll review stale.
12. IF time segments overlap another Work_Record for the same Worker without an approved concurrent-assignment policy, THEN THE System SHALL reject verification.
13. IF a time segment has a non-positive duration, negative break, or negative payable duration, THEN THE System SHALL reject verification.

### Requirement 18: Single Worker Acknowledgement and Dispute

**User Story:** As a Worker, I want to acknowledge or dispute the exact daily record I reviewed, so that my signature cannot apply to changed information.

#### Acceptance Criteria

1. WHEN a Worker opens a record for acknowledgement, THE System SHALL display the current version, Work_Record_Digest, date, start, end, breaks, classifications, remarks, and visible evidence summary.
2. WHEN a Worker acknowledges a record, THE System SHALL bind Worker_Acknowledgement to the authenticated Worker, current version, current Work_Record_Digest, and signing timestamp.
3. IF the submitted version or digest differs from the current record, THEN THE System SHALL reject acknowledgement as changed since review.
4. IF the record is not owned by the authenticated Worker, THEN THE System SHALL deny acknowledgement without disclosing protected details.
5. IF the record is not in Pending Worker Acknowledgement state, THEN THE System SHALL reject acknowledgement with an ineligible-state reason.
6. WHEN a Worker disputes a record, THE System SHALL bind the Dispute to the current version and require a reason.
7. WHEN a Dispute is accepted, THE System SHALL route the record to the configured Supervisor or HR resolution state.
8. WHILE a Dispute remains unresolved, THE System SHALL exclude the record from eligible batch acknowledgement.
9. WHILE a Dispute remains a configured readiness blocker, THE System SHALL prevent finalization of the affected scope.
10. WHEN a disputed record is superseded by an authorized correction, THE System SHALL require Worker review of the new version.
11. WHEN a record changes after Worker_Acknowledgement, THE System SHALL preserve the prior acknowledgement only as evidence of the prior version.

### Requirement 19: Safe Batch and Monthly Worker Confirmation

**User Story:** As a Worker, I want to confirm multiple eligible daily records safely, so that monthly confirmation reduces effort without becoming a blanket signature.

#### Acceptance Criteria

1. WHEN the System presents batch candidates, THE System SHALL list each date, current version, digest indicator, eligibility state, and exclusion reason.
2. WHEN a Worker prepares a batch, THE System SHALL require explicit selection of each included Work_Record.
3. WHEN a Batch_Envelope is signed, THE System SHALL bind the envelope to the authenticated Worker and ordered selected record digests.
4. WHEN a batch is submitted, THE System SHALL revalidate ownership, state, current version, current digest, unresolved Dispute, and blocking exception for each selected item.
5. WHEN a selected item remains eligible, THE System SHALL create one independently traceable Worker_Acknowledgement for that item.
6. IF a selected item changed since review, THEN THE System SHALL exclude the item with a changed-since-review reason.
7. IF a selected item is disputed, blocked, ineligible, or owned by another Worker, THEN THE System SHALL exclude the item with a specific safe reason.
8. WHEN a batch contains eligible and ineligible items, THE System SHALL return separate acknowledged and excluded item lists.
9. WHEN a batch retry repeats the same valid request, THE System SHALL avoid duplicate Worker_Acknowledgements.
10. WHEN the mobile batch view detects a changed item, THE System SHALL deselect the item and mark the item “Changed since review”.
11. WHEN a Worker reaches final batch signing, THE System SHALL present a per-date summary of every selected item.
12. THE System SHALL prevent a monthly signature from applying to unresolved, disputed, unselected, ineligible, or changed records.

### Requirement 20: Monthly Work Record Card Experience

**User Story:** As a Worker or authorized reviewer, I want a familiar monthly Work Record view, so that I can reconcile digital records with the physical-card workflow.

#### Acceptance Criteria

1. WHEN a Monthly_Work_Record_Card is displayed, THE System SHALL present one row or card for every calendar date in the selected month.
2. WHEN a month contains fewer than 31 dates, THE System SHALL limit active date rows or cards to the valid dates in that month.
3. WHEN a daily record exists, THE System SHALL present date, start time, end time, Worker acknowledgement status, Supervisor verification status, and remark.
4. WHEN permitted calculation data exists, THE System SHALL present normal days or hours, categorized overtime, incentive, advance, and total amount projections using configured terminology.
5. WHEN monetary data is not finalized, THE System SHALL label the monetary data Estimated/Not Finalized.
6. WHEN a Calendar_Day is Closed, THE System SHALL keep the date visible and visually distinguish the Closed state using text or icon in addition to color.
7. WHEN valid Closed_Day_OT exists, THE System SHALL display “Closed — Authorized OT Work”.
8. WHEN monthly totals are displayed, THE System SHALL distinguish verified, pending, and disputed hours.
9. WHEN evidence drill-down is authorized, THE System SHALL expose meeting, attendance, checkout, edit, authorization, consent, and Audit_Event references applicable to the record.
10. IF evidence drill-down is not authorized, THEN THE System SHALL omit protected evidence fields.
11. WHEN the Monthly_Work_Record_Card is used on mobile, THE System SHALL provide date cards with status, start, end, duration, and the current primary action.
12. WHEN the Monthly_Work_Record_Card is used on larger screens, THE System SHALL support a monthly grid or table with expandable decision-critical details.
13. WHEN a user navigates the monthly card by keyboard or assistive technology, THE System SHALL expose date, calendar state, record state, and available action through semantic labels.
14. THE System SHALL preserve recognizable source-artifact concepts without requiring a pixel-perfect reproduction of the paper forms.
15. THE System SHALL use generic placeholders rather than sample personal names, identifiers, or signatures in seeded data and documentation artifacts.

### Requirement 21: Effective-Dated Compensation and Pay Rules

**User Story:** As a Payroll_Maker, I want compensation and calculation policies versioned by effective date, so that historical and current payroll resolve deterministically.

#### Acceptance Criteria

1. WHEN a Payroll_Maker drafts a Compensation_Profile, THE System SHALL record Worker, effective interval, pay basis, currency, compensation values, rule references, reason, and version.
2. WHEN a Compensation_Profile is submitted for approval, THE System SHALL bind submission to the exact content digest.
3. WHEN Maker_Checker_Policy applies, THE System SHALL require an eligible checker other than the Payroll_Maker.
4. WHEN a Compensation_Profile is approved, THE System SHALL make the approved version immutable.
5. IF approved Compensation_Profile intervals overlap ambiguously for the same Worker and pay basis, THEN THE System SHALL reject approval.
6. WHEN a calculation resolves compensation for a work date, THE System SHALL select the approved version effective for that date.
7. IF no unambiguous approved Compensation_Profile covers a required date, THEN THE System SHALL block the affected calculation.
8. WHEN a Pay_Rule is created or changed, THE System SHALL version the rule, effective interval, units, conditions, tables, rounding mode, caps, references, and approval state.
9. IF a Pay_Rule contains a cycle, ambiguous currency, unsupported reference, or non-deterministic operation, THEN THE System SHALL reject approval.
10. THE System SHALL calculate monetary values using fixed-precision decimal arithmetic and explicit currency.
11. THE System SHALL prevent binary floating-point arithmetic from being authoritative for payroll monetary values.
12. WHEN compensation or Pay_Rules change, THE System SHALL preserve prior approved versions and the reason for change.
13. WHEN an Adjustment is created, THE System SHALL require a reason code, effective Payroll_Period, amount or approved formula, maker, and version.
14. WHEN an Adjustment requires approval, THE System SHALL exclude the Adjustment from finalization until an eligible checker approves the exact digest.

### Requirement 22: Configurable and Validated Payroll Policy

**User Story:** As a payroll policy owner, I want formulas and legal assumptions configured and qualified, so that the System does not invent statutory or contractual treatment.

#### Acceptance Criteria

1. THE System SHALL represent normal pay, overtime categories, allowances, incentives, advances, deductions, adjustments, caps, eligibility, and rounding through versioned approved configuration.
2. THE System SHALL avoid hard-coded unverified statutory rates, thresholds, standard days, standard hours, or contract interpretations.
3. WHEN a Tenant configures a payroll policy, THE System SHALL record the policy owner, source reference, effective date, approval evidence, and version.
4. WHEN a Tenant requests production payroll enablement, THE System SHALL require recorded Qualified_Stakeholder validation of Singapore PTE, MOM, contract-specific, payroll, and legal policy applicable to that Tenant.
5. IF required policy validation is missing or expired under configured governance, THEN THE System SHALL block production payroll finalization.
6. WHEN a statutory or contractual policy changes, THE System SHALL require a new effective-dated version and approval rather than mutate historical policy.
7. WHEN the System explains a calculated line, THE System SHALL identify the configured rule version and source inputs rather than claim independent legal advice.
8. WHERE statutory reporting adapters are enabled, THE System SHALL require a Qualified_Stakeholder-approved contract and validation suite before production submission.
9. IF a policy choice remains undecided, THEN THE System SHALL expose the choice as an implementation or Tenant-rollout prerequisite without assigning a default legal meaning.

### Requirement 23: Estimated Payroll Calculation and Provenance

**User Story:** As an HR_Reviewer, I want current explainable payroll estimates, so that I can identify issues before finalization.

#### Acceptance Criteria

1. WHEN an eligible Work_Record, Compensation_Profile, Pay_Rule, Calendar_Day, or Adjustment changes, THE System SHALL initiate recalculation of the affected Worker and Payroll_Period without waiting for period finalization.
2. WHEN a calculation starts, THE System SHALL resolve exact applicable calendar, Work_Record, compensation, rule, adjustment, currency, and engine versions.
3. WHEN a calculation succeeds, THE System SHALL store line-level quantities, units, rates or multipliers, amounts, source references, rule version, and explanation key.
4. WHEN a calculation succeeds, THE System SHALL label the result Estimated/Not Finalized.
5. WHEN the same canonical inputs and engine version are calculated repeatedly, THE System SHALL return the same result digest and calculation lines.
6. WHEN repeated calculation uses the same calculation Idempotency_Key, THE System SHALL create at most one successful result for that key.
7. WHEN Work_Record segments are calculated, THE System SHALL prevent duplicate counting across mutually exclusive payable classifications.
8. WHERE an approved rule decomposes a Work_Record into classifications, THE System SHALL use non-overlapping segments for the decomposition.
9. WHEN the System aggregates gross pay, THE System SHALL sum configured positive basic, normal, overtime, incentive, allowance, and positive Adjustment lines under the approved rounding sequence.
10. WHEN the System aggregates net pay, THE System SHALL subtract configured advance, deduction, and negative Adjustment lines from gross under the approved rounding sequence.
11. IF currency is ambiguous or inconsistent, THEN THE System SHALL block the affected calculation.
12. IF a mandatory source is missing, invalid, duplicated, overlapping, disputed, unverified, or unauthorized, THEN THE System SHALL produce a warning or blocker according to approved non-downgradable and Tenant policy.
13. WHEN a calculation is stored, THE System SHALL store the Input_Digest, result digest, calculated time, and exact version provenance.
14. WHEN a user drills into a calculation line with permission, THE System SHALL present the source records, rule version, and calculation explanation.
15. THE System SHALL prevent a payroll estimate from being represented as a finalized Payroll_Snapshot or Payslip.

### Requirement 24: Daily and Monthly HR Review

**User Story:** As an HR_Reviewer, I want daily and monthly review with source-change awareness, so that approvals apply only to current evidence and calculations.

#### Acceptance Criteria

1. WHEN an authorized HR_Reviewer opens a daily record, THE System SHALL present current workflow state, evidence, exceptions, versions, and payroll inclusion status.
2. WHEN an authorized HR_Reviewer opens a monthly ledger, THE System SHALL list Workers with calculation, readiness, review, exception, and source-change states.
3. WHEN an HR_Reviewer approves a Payroll_Ledger_Entry, THE System SHALL bind the decision to the current Input_Digest.
4. IF any referenced source change produces a different Input_Digest, THEN THE System SHALL mark the affected review Stale_Review or Reopened.
5. WHEN a review becomes stale, THE System SHALL identify the affected Worker, Payroll_Period, changed source classes, and recalculation status.
6. WHEN a source change affects only a subset of Workers or periods, THE System SHALL limit recalculation and re-review to the dependency-affected scope.
7. IF recalculation fails, THEN THE System SHALL keep the affected entry blocked and provide a safe diagnostic and correlation identifier.
8. WHEN recalculation completes after a source change, THE System SHALL require review of the new Input_Digest before finalization.
9. WHEN an approved Adjustment changes an estimate, THE System SHALL display before-and-after impact to authorized HR or payroll users.
10. IF an HR_Reviewer lacks salary-view Permission, THEN THE System SHALL omit monetary values while presenting only separately authorized operational review data.
11. WHEN review, rejection, reopening, or stale-state transition occurs, THE System SHALL create an Audit_Event.

### Requirement 25: Payroll Readiness and Finalization

**User Story:** As a Payroll_Finalizer, I want readiness checks and atomic finalization, so that finalized payroll is complete, current, and reproducible.

#### Acceptance Criteria

1. WHEN readiness is requested, THE System SHALL evaluate published calendar coverage, required Work_Record verification, unresolved Disputes, compensation coverage, rule evaluation, overlap, duplicate work, overtime authorization, Worker_PreWork_Consent, review freshness, Adjustment approval, currency consistency, Maker_Checker_Policy, and input integrity.
2. WHEN readiness completes, THE System SHALL produce a Readiness_Report with blockers, warnings, checked Input_Digest, policy versions, and timestamp.
3. THE System SHALL prevent Tenant configuration from downgrading cross-Tenant mismatch, structural corruption, stale Input_Digest, or failed snapshot integrity from blocker severity.
4. IF any readiness blocker remains, THEN THE System SHALL reject finalization and return the blocker list.
5. WHEN finalization is requested, THE System SHALL require Payroll_Finalizer Permission, applicable Resource_Scope, Step_Up_Authentication, and Maker_Checker_Policy satisfaction.
6. WHEN finalization begins, THE System SHALL acquire a lock scoped to the Payroll_Period or approved finalization cohort.
7. WHEN finalization reloads inputs, THE System SHALL compare the current Input_Digest with the Readiness_Report Input_Digest.
8. IF the current Input_Digest differs from the checked Input_Digest, THEN THE System SHALL roll back the finalization attempt and require new readiness and review.
9. IF any included Worker result lacks review at the current Input_Digest, THEN THE System SHALL block finalization.
10. WHEN finalization succeeds, THE System SHALL create one immutable, versioned Payroll_Snapshot containing every included Worker result and locked input reference.
11. WHEN finalization succeeds, THE System SHALL record finalizer, finalization time, approval evidence, Readiness_Report digest, and snapshot digest.
12. IF any step fails during snapshot or Payslip creation, THEN THE System SHALL roll back every new snapshot and Payslip artifact from that attempt.
13. WHEN concurrent finalization attempts target the same scope, THE System SHALL permit at most one attempt to create the applicable snapshot version.
14. WHEN a Payroll_Snapshot is finalized, THE System SHALL prevent mutation of its content and locked references.
15. WHEN finalization completes or fails, THE System SHALL release the finalization lock without unlocking immutable snapshot inputs.

### Requirement 26: Finalized Payslips and Worker Monetary Visibility

**User Story:** As a Worker, I want clear estimates and reliable finalized Payslips, so that I can distinguish provisional information from final payroll.

#### Acceptance Criteria

1. WHEN a Payroll_Snapshot is finalized successfully, THE System SHALL create one versioned Payslip for each included Worker result.
2. WHEN Payslips are published, THE System SHALL make each Payslip visible only to its Worker and separately authorized Principals.
3. WHEN a Worker opens a finalized Payslip, THE System SHALL require authentication and the configured fresh-authentication level.
4. WHEN a Worker views a finalized Payslip, THE System SHALL display the Payroll_Period, snapshot version, finalized status, line breakdown, currency, and publication state.
5. WHERE Tenant policy permits Worker monetary estimates, THE System SHALL display current estimates with a prominent Estimated/Not Finalized label.
6. WHERE Tenant policy hides Worker monetary estimates, THE System SHALL omit current monetary estimate amounts from Worker responses.
7. WHEN Worker monetary estimates are hidden, THE System SHALL continue to display separately authorized hours and Work_Record statuses.
8. THE System SHALL make a published finalized Payslip available to the applicable Worker regardless of the Tenant’s current estimate-visibility setting.
9. WHEN a Payslip file is delivered, THE System SHALL use authenticated retrieval or a short-lived protected reference.
10. THE System SHALL exclude Payslip documents from general-purpose service-worker caches.

### Requirement 27: Post-Finalization Scope Boundary

**User Story:** As a product owner, I want post-finalization behavior explicitly bounded, so that implementation does not invent unsupported correction workflows.

#### Acceptance Criteria

1. THE System SHALL treat finalized Payroll_Snapshots and published Payslips as immutable in this feature version.
2. IF a user requests a post-final payroll correction workflow, THEN THE System SHALL return an explicit out-of-scope response without modifying the Payroll_Snapshot.
3. IF a user requests an off-cycle payroll workflow, THEN THE System SHALL return an explicit out-of-scope response without creating an off-cycle payroll artifact.
4. IF a user requests a revised Payslip workflow, THEN THE System SHALL return an explicit out-of-scope response without replacing a published Payslip.
5. THE System SHALL permit data-model and integration extension references for future correction, off-cycle, or revised-Payslip capabilities without exposing an executable workflow.
6. WHEN future-scope metadata references a prior snapshot, THE System SHALL prevent that metadata from changing the prior snapshot digest or content.

### Requirement 28: Worker and Supervisor Browser-Capable Installable PWAs

**User Story:** As a Worker or Supervisor, I want role-focused applications that work installed or in a browser, so that installation is optional and core workflows remain available.

#### Acceptance Criteria

1. THE System SHALL provide installable Worker_PWA and Supervisor_PWA role shells in supported capable browsers.
2. THE System SHALL provide the same authorized core workflow semantics through normal supported browser tabs without installation.
3. WHEN a role shell is launched, THE System SHALL use a role-specific name, icon, start destination, navigation, and cache scope.
4. WHEN a Principal is authorized for multiple roles, THE System SHALL require an explicit server-authorized role switch.
5. THE System SHALL prevent installation of one role shell from granting Permission to another role.
6. WHEN installation capability is unavailable or installation is declined, THE System SHALL keep the browser workflow reachable.
7. WHEN an install prompt is displayed, THE System SHALL make the prompt dismissible and explain the effect of installation.
8. THE System SHALL prevent an install prompt from blocking an operational workflow.
9. WHEN an application update is available, THE System SHALL report whether activation is safe, deferred, requires reload, or is security-required.
10. WHEN a non-critical update arrives during an in-progress sensitive form, THE System SHALL defer disruptive activation until a safe boundary.
11. WHEN a security-critical update or revocation requires session invalidation, THE System SHALL require fresh authentication before sensitive work resumes.
12. IF the supported browser, device, operating-system, camera, storage, or assistive-technology baseline remains undecided, THEN THE System SHALL treat the baseline as an implementation sign-off prerequisite.
13. THE System SHALL keep Management_Web, Tenant_Management_Admin_Panel, and Platform_Admin_Console responsive without assuming those surfaces are installable.

### Requirement 29: Pending Offline Actions and Reconciliation

**User Story:** As a mobile Worker or Supervisor, I want safe offline handling, so that temporary connectivity loss does not create false authoritative records.

#### Acceptance Criteria

1. WHEN an approved offline-capable action is captured without network access, THE System SHALL store a Pending_Offline_Action with local action identifier, Tenant, Principal, session binding, action digest, Expected_Version, creation time, expiry, and Pending status.
2. WHEN local encryption capability is available, THE System SHALL encrypt the minimized Pending_Offline_Action payload with session-bound key material.
3. WHILE a Pending_Offline_Action lacks server acceptance, THE System SHALL keep authoritative domain state unchanged for that Pending_Offline_Action.
4. WHILE a Pending_Offline_Action lacks server acceptance, THE System SHALL display the action as Pending rather than accepted, verified, acknowledged, approved, or finalized.
5. WHEN connectivity, launch, focus, or an authorized manual retry occurs, THE System SHALL attempt foreground reconciliation of eligible Pending_Offline_Actions.
6. WHERE browser background synchronization is available, THE System SHALL treat background retry as opportunistic rather than required for correctness.
7. WHEN reconciliation starts, THE System SHALL revalidate authentication, Tenant status, Role, Resource_Scope, session binding, action digest, Expected_Version, workflow state, and applicable deadline.
8. WHEN a reconciled action is accepted, THE System SHALL apply at most one authoritative mutation for the local action identifier.
9. IF a reconciled action is rejected, expired, or conflicts, THEN THE System SHALL display a specific safe outcome and required recovery path.
10. IF a Pending_Offline_Action belongs to a changed session binding, THEN THE System SHALL reject reconciliation and purge the sensitive payload.
11. IF a Pending_Offline_Action is expired, THEN THE System SHALL reject reconciliation and purge the sensitive payload.
12. WHEN reconciliation reaches a terminal accepted or rejected outcome, THE System SHALL delete the sensitive local payload.
13. WHEN reconciliation produces a conflict, THE System SHALL retain only the minimum information required for explicit user resolution.
14. THE System SHALL prevent last-write-wins resolution of an offline conflict involving attendance, consent, checkout, acknowledgement, approval, compensation, or payroll.
15. IF offline depth, action allowlist, queue size, or queue expiry remains undecided, THEN THE System SHALL treat those values as implementation and rollout prerequisites.

### Requirement 30: Progressive Capability Fallback

**User Story:** As a Worker or Supervisor, I want safe alternatives when device capabilities are missing, so that unavailable features do not break or weaken the workflow.

#### Acceptance Criteria

1. IF push capability is unavailable or denied, THEN THE System SHALL provide an in-app task inbox or badge.
2. WHERE approved fallback notification channels are configured, THE System SHALL use the configured channel when push delivery is unavailable.
3. IF background synchronization is unavailable, THEN THE System SHALL provide foreground, online-event, and manual retry with visible queue status.
4. IF camera capability is unavailable or denied, THEN THE System SHALL offer an accessible system picker where policy permits or a supervised checkout exception.
5. IF geolocation is unavailable, denied, or fails policy, THEN THE System SHALL explain the purpose and route the action through available Checkpoint evidence or exception review.
6. IF durable local storage is unavailable, THEN THE System SHALL require online submission for actions that cannot be stored safely.
7. IF the network is unavailable, THEN THE System SHALL capture only policy-permitted minimized Pending_Offline_Actions.
8. IF an offline action is not permitted by policy, THEN THE System SHALL present an online-required state and available exception process.
9. WHEN a capability permission is requested, THE System SHALL explain the purpose at the time of use and provide denial recovery guidance.
10. IF a capability fallback is used, THEN THE System SHALL preserve Tenant isolation, consent, attendance, checkout, acknowledgement, evidence, and payroll invariants.
11. THE System SHALL prevent a Capability_Fallback from silently marking required evidence complete.

### Requirement 31: Sensitive Client Data and Session Clearing

**User Story:** As a user on a personal or shared device, I want sensitive local data cleared on access termination, so that prior sessions cannot leak payroll or evidence.

#### Acceptance Criteria

1. THE System SHALL exclude salary, monetary estimates, Payslip files, high-sensitivity evidence, authentication secrets, and complete Worker rosters from general-purpose client caches.
2. WHEN the System stores policy-permitted local Tenant data, THE System SHALL minimize the data and bind the data to Tenant, Principal, role, session, and expiry.
3. WHEN logout occurs, THE System SHALL stop background retry and invalidate the server session.
4. WHEN logout, session revocation, Role removal, Tenant suspension, remote wipe, device revocation, or security reset occurs, THE System SHALL delete session keys and sensitive session-bound local payloads.
5. WHEN a termination signal occurs, THE System SHALL clear role-specific read models containing protected Tenant data.
6. WHEN a termination signal occurs, THE System SHALL revoke or render unusable the associated push subscription and session binding.
7. IF local payload deletion cannot be confirmed, THEN THE System SHALL destroy the key material required to decrypt the payload.
8. WHEN termination-signal processing completes, THE System SHALL require fresh authentication and authorization for subsequent sensitive access.
9. WHERE policy permits static app-shell caching after logout, THE System SHALL retain only assets that reveal no Tenant or user data.
10. WHEN fresh authentication expires for a sensitive action, THE System SHALL require Step_Up_Authentication before the action proceeds.

### Requirement 32: Responsive Product Experience and State Communication

**User Story:** As a user of any role, I want a coherent responsive experience, so that I can complete authorized work on mobile, tablet, or desktop without losing critical context.

#### Acceptance Criteria

1. THE System SHALL use a shared versioned design system for semantic color, typography, spacing, density, focus, status, motion, breakpoints, and data visualization.
2. WHEN a surface renders on mobile, THE System SHALL use touch-capable controls and progressive disclosure without removing decision-critical data.
3. WHEN a surface renders on tablet, THE System SHALL adapt eligible workflows to split panes, grids, cards, or detail views according to available space.
4. WHEN a surface renders on desktop, THE System SHALL support dense comparison, side-by-side evidence, keyboard operation, and complex configuration.
5. WHEN a complex management table reflows on mobile, THE System SHALL preserve blockers, approval context, and required actions.
6. WHEN data is loading, empty, erroneous, partial, denied, offline, Pending, conflicting, stale, successful, maintained, restricted, or suspended, THE System SHALL present a distinguishable state and available next action.
7. THE System SHALL use optimistic authoritative status only for reversible low-risk presentation state.
8. WHEN attendance, consent, checkout, acknowledgement, approval, compensation, or payroll is awaiting server acceptance, THE System SHALL display Pending status.
9. WHEN a form has safe non-sensitive draft content, THE System SHALL preserve the content according to approved autosave policy and present unsaved or conflict recovery.
10. WHEN an action is destructive or high risk, THE System SHALL present the affected scope and require explicit confirmation or approval according to policy.
11. WHEN charts present operational or payroll information, THE System SHALL provide an equivalent accessible table or textual summary.
12. WHEN route or data loading is delayed, THE System SHALL provide non-deceptive progress feedback while retaining stale or Pending labels where applicable.
13. THE System SHALL compose dashboard tiles, totals, actions, alerts, and links from the Principal’s current Permissions and Resource_Scopes.
14. IF a Principal lacks access to underlying data, THEN THE System SHALL omit the corresponding dashboard data and deep link.

### Requirement 33: Accessibility and Localization Readiness

**User Story:** As a user with diverse access and language needs, I want accessible and localization-ready workflows, so that critical workforce and payroll tasks remain perceivable and operable.

#### Acceptance Criteria

1. THE System SHALL provide semantic landmarks, logical focus order, visible focus, keyboard operation, and associated form errors for critical workflows.
2. THE System SHALL express status through text or icon in addition to color.
3. WHEN content is zoomed or reflowed within the approved accessibility baseline, THE System SHALL preserve content and functionality.
4. WHEN reduced-motion preference is active, THE System SHALL reduce non-essential motion.
5. WHEN touch input is used, THE System SHALL use touch targets that meet the approved accessibility target.
6. WHEN authentication is required, THE System SHALL provide an accessible authentication flow within the approved conformance target.
7. THE System SHALL use message keys, plural rules, and locale-aware date, number, and currency formatting.
8. WHEN time is displayed, THE System SHALL present the relevant Site time zone or an unambiguous localized zone indication.
9. THE System SHALL support text expansion and bidirectional-layout readiness without embedding required text in imagery.
10. IF a translation is unavailable, THEN THE System SHALL use the approved fallback locale and identify untranslated content for remediation.
11. WHEN accessibility release-gate enablement is requested, THE System SHALL require approval of the target conformance level, assistive-technology matrix, audit cadence, and exception governance.
12. WHERE WCAG 2.2 AA is approved as the conformance target, THE System SHALL apply WCAG 2.2 AA checks as release criteria.
13. IF the initial language set or right-to-left requirement remains undecided, THEN THE System SHALL treat the decision as a rollout prerequisite rather than infer a language scope.

### Requirement 34: Notifications and Action Deadlines

**User Story:** As a user with pending work, I want reliable, privacy-minimized notifications, so that I can act without treating delivery as domain authority.

#### Acceptance Criteria

1. WHEN Supervisor attendance is called, THE System SHALL create a Worker task to confirm arrival.
2. WHEN an Overtime_Assignment is authorized, THE System SHALL create a task for the named Worker to consent or decline.
3. WHEN a consent deadline approaches according to configured policy, THE System SHALL notify the Worker and configured escalation recipient.
4. WHEN checkout is expected or missed, THE System SHALL create the configured Worker and Supervisor tasks.
5. WHEN a Work_Record is verified, THE System SHALL create a Worker acknowledgement task.
6. WHEN a reviewed record changes, THE System SHALL notify the affected Worker or HR_Reviewer according to current workflow state.
7. WHEN a Dispute is created, THE System SHALL notify the assigned Supervisor or HR resolution role.
8. WHEN payroll review becomes stale or readiness is blocked, THE System SHALL notify the responsible authorized role.
9. WHEN a Payslip is published, THE System SHALL notify the applicable Worker without including sensitive payroll amounts in a general notification payload.
10. THE System SHALL keep domain deadlines and states authoritative independently of notification delivery.
11. IF notification delivery fails, THEN THE System SHALL retain the underlying task and apply configured idempotent retry or fallback.
12. WHEN notification templates change, THE System SHALL version Tenant, locale, channel, content, and approval state.
13. THE System SHALL minimize personal and payroll data in notification payloads.
14. IF channel, escalation timing, or fallback priority remains undecided, THEN THE System SHALL treat the choice as Tenant configuration or rollout prerequisite.

### Requirement 35: Security and Privacy Controls

**User Story:** As a security and privacy owner, I want layered protection for workforce and payroll information, so that sensitive operations and evidence remain confidential and accountable.

#### Acceptance Criteria

1. THE System SHALL enforce deny-by-default access and least privilege for human and Integration_Principals.
2. WHEN payroll finalization, high-risk Role assignment, sensitive export, or Break_Glass_Session occurs, THE System SHALL require MFA or Step_Up_Authentication according to approved risk policy.
3. THE System SHALL encrypt data in transit and at rest using approved implementation standards.
4. WHEN encryption or signing keys rotate, THE System SHALL preserve controlled validation of approved active or grace versions and retire disallowed versions.
5. THE System SHALL apply field-level authorization to salary, protected identifiers, signatures, location, and evidence.
6. WHEN protected objects are downloaded, THE System SHALL use authenticated access or short-lived signed references.
7. THE System SHALL prevent protected evidence and Payslip objects from being publicly readable.
8. WHEN files are uploaded, THE System SHALL validate content type and scan for malicious content before trusted use.
9. THE System SHALL apply rate limiting, anti-automation, secure session management, and request-forgery protections appropriate to each interface.
10. WHEN location or device evidence is collected, THE System SHALL disclose the approved purpose and collect no greater precision or duration than the approved policy requires.
11. THE System SHALL prevent sample-image personal names, identifiers, and signatures from being copied into seeds, tests, demonstrations, or generated examples.
12. WHEN platform health or usage telemetry is produced, THE System SHALL use aggregate or opaque Tenant references and exclude payroll amounts and raw evidence.
13. THE System SHALL prevent feature-release controls from overriding domain authorization or payroll invariants.
14. WHEN security scans detect a release-blocking issue under approved policy, THE System SHALL prevent the affected release until remediation or authorized time-bound exception.

### Requirement 36: Auditability and Integrity

**User Story:** As an Auditor, I want append-only, integrity-verifiable history, so that each sensitive decision and signature can be reconstructed.

#### Acceptance Criteria

1. WHEN an auditable action occurs, THE System SHALL append an Audit_Event containing actor, actor boundary, Tenant, action, resource, version, timestamp, correlation, and outcome.
2. WHEN data changes, THE System SHALL include before-and-after digests and reason in the Audit_Event.
3. WHEN authorization is evaluated for a sensitive action, THE System SHALL record the decision and policy version.
4. WHEN approval, consent, acknowledgement, review, or finalization occurs, THE System SHALL record the exact signed or approved digest.
5. WHEN a Support_Session action occurs, THE System SHALL record Support_Case, Support_Grant, Support_Session, Platform_Support_Agent, Tenant, purpose, scope, and outcome.
6. WHEN salary or Payslip content is accessed, THE System SHALL record the access according to approved audit policy.
7. WHEN Tenant lifecycle, entitlement, quota, feature release, app installation, push subscription, offline reconciliation, or client purge state changes, THE System SHALL create an Audit_Event.
8. THE System SHALL make Audit_Events append-only to application workflows.
9. WHEN audit integrity anchoring is configured, THE System SHALL chain or periodically anchor high-value Audit_Events with verifiable hashes.
10. IF an audit pipeline write required for a sensitive mutation fails, THEN THE System SHALL block or safely quarantine the mutation according to approved atomicity policy.
11. WHEN an authorized Auditor queries history, THE System SHALL enforce Tenant, row, field, time, and data-classification scope.
12. THE System SHALL represent approval, consent, acknowledgement, and review as first-class domain records rather than rely only on generic Audit_Events.

### Requirement 37: Observability, Reliability, and Recovery

**User Story:** As a service operator, I want privacy-safe observability and recoverable processing, so that failures are detectable without leaking sensitive data.

#### Acceptance Criteria

1. THE System SHALL measure calendar, attendance, consent, checkout, Work_Record, payroll, notification, Tenant provisioning, support, PWA update, offline reconciliation, authorization, and finalization health.
2. THE System SHALL exclude Worker identifiers, salary amounts, exact locations, push endpoints, and raw protected identifiers from general logs and metric labels.
3. WHEN a business action crosses asynchronous components, THE System SHALL propagate correlation and causation identifiers.
4. WHEN an Integration_Event is delivered more than once, THE System SHALL make the consumer idempotent by event identifier.
5. WHEN a transactional state change emits an Integration_Event, THE System SHALL preserve state-and-event atomicity through an approved transactional delivery mechanism.
6. WHEN asynchronous processing fails beyond configured retry policy, THE System SHALL retain a diagnosable recovery item without silently losing the domain action.
7. WHEN persistent calculation failure, queue lag, replay spikes, finalization conflicts, audit failure, Tenant-isolation guard activation, unauthorized platform access, support overreach, app-update failure, sensitive-purge failure, or evidence-store failure crosses an approved threshold, THE System SHALL alert the configured operational owner.
8. WHEN recovery tooling is used, THE System SHALL enforce Tenant scope and create an Audit_Event.
9. THE System SHALL support point-in-time recovery and versioned immutable-object recovery under approved recovery policy.
10. WHEN restore procedures are executed, THE System SHALL verify restored Tenant isolation, object integrity, and Payroll_Snapshot immutability.
11. IF availability, recovery time, recovery point, or throughput targets remain undecided, THEN THE System SHALL treat numeric targets as implementation sign-off prerequisites.

### Requirement 38: Performance and Scalability Gates

**User Story:** As a product and operations owner, I want measurable performance gates, so that responsive behavior and month-end processing are validated on representative conditions.

#### Acceptance Criteria

1. WHEN implementation performance sign-off is requested, THE System SHALL require approved numeric budgets for startup, interaction, API classes, offline reconciliation, data views, background work, payroll throughput, availability, and recovery.
2. WHEN performance budgets are approved, THE System SHALL associate each budget with percentile, workload, device class, network class, data volume, and measurement method.
3. WHEN performance tests run, THE System SHALL include representative low-, mid-, and high-tier supported devices and constrained networks.
4. WHEN load tests run, THE System SHALL include morning attendance, shift-end checkout, mass foreground reconciliation, notification bursts, month-end recalculation, bulk review, export, provisioning, support audit, and finalization workloads.
5. WHEN one Tenant creates sustained load, THE System SHALL apply approved fair scheduling and quota controls to limit noisy-neighbor impact.
6. WHEN source data changes, THE System SHALL recalculate only dependency-affected Workers and Payroll_Periods unless a full recalculation is explicitly required.
7. WHEN concurrent identical calculations occur, THE System SHALL deduplicate work by Tenant, Worker, Payroll_Period, Input_Digest, and engine version.
8. WHEN large exports are requested, THE System SHALL process the export asynchronously or stream the export according to the approved size policy.
9. WHEN sensitive response caching would violate data classification, THE System SHALL retain the no-store policy even if caching would improve a synthetic performance score.
10. IF a measured release exceeds an approved performance budget, THEN THE System SHALL block release or require an explicit owner, rationale, expiration, and remediation plan for a time-bound exception.

### Requirement 39: Export and Integration Extension Points

**User Story:** As an authorized Tenant operator, I want traceable exports and scoped integrations, so that data can move to approved business systems without losing status or provenance.

#### Acceptance Criteria

1. WHEN an authorized Principal exports a monthly work card, THE System SHALL include recognizable dates, start and end times, verification, acknowledgement, remarks, statuses, and permitted totals.
2. WHEN an authorized Principal exports an overtime package, THE System SHALL include planned dates and times, named Worker references, authorization scope, consent state, applicable staggered-rest-day acknowledgement, and digests.
3. WHEN an authorized Principal exports a payroll ledger, THE System SHALL include permitted calculation lines, statuses, warnings, blockers, source versions, and result digest.
4. WHEN an authorized Principal exports a Payslip, THE System SHALL include the applicable snapshot and Payslip version references.
5. WHEN an authorized Principal exports an audit or evidence package, THE System SHALL include an Integrity_Manifest.
6. WHEN any export is generated, THE System SHALL include Tenant, period or scope, generation time, report version, source versions, classification, and integrity data.
7. THE System SHALL preserve Pending, disputed, stale, excluded, and blocked states in exports rather than flatten those states into accepted facts.
8. WHEN a sensitive export is requested, THE System SHALL enforce separate export Permission, Resource_Scope, Step_Up_Authentication where configured, and Audit_Event creation.
9. WHEN a protected export is delivered, THE System SHALL encrypt or protect the export with time-limited access according to approved policy.
10. WHEN an Integration_Principal calls an API, THE System SHALL enforce Tenant, action, resource, data classification, credential version, and rotation state.
11. WHEN inbound integration data is retried, THE System SHALL apply replay protection and idempotency.
12. WHEN an integration record fails processing, THE System SHALL produce a reconciliation result and recoverable failure state.
13. WHEN Integration_Events contain sensitive facts, THE System SHALL use minimized references or approved encrypted restricted channels.
14. IF an HRIS, accounting, banking, signing, notification, checkpoint, or statutory contract remains undecided, THEN THE System SHALL treat the contract and vendor selection as a non-binding implementation prerequisite.

### Requirement 40: Retention, Data Lifecycle, and Privacy Rights

**User Story:** As a privacy and records owner, I want classified retention and deletion controls, so that each data type follows approved legal and contractual obligations.

#### Acceptance Criteria

1. THE System SHALL classify identity, employment, Work_Record, signature, consent, location, device, payroll, Payslip, Audit_Event, notification, and idempotency data into retention classes.
2. WHEN a Tenant retention policy is approved, THE System SHALL version the duration, trigger, archive, deletion, anonymization, hold, and jurisdiction rules for each class.
3. WHEN a retention action becomes due, THE System SHALL evaluate legal hold, contractual retention, Payroll_Snapshot integrity, and approved deletion constraints before execution.
4. IF deletion would break a locked Payroll_Snapshot reference or required audit integrity, THEN THE System SHALL apply the approved retention or anonymization alternative rather than corrupt the snapshot.
5. WHEN a retention or deletion action executes, THE System SHALL record scope, policy version, actor or scheduler, result, and Audit_Event.
6. WHEN a permitted data-subject or Tenant request is processed, THE System SHALL limit the response to the requester’s authorized Tenant and data scope.
7. IF exact retention durations, data residency, lawful-deletion constraints, or anchoring policy remain undecided, THEN THE System SHALL block production rollout of the affected data class until Qualified_Stakeholder approval is recorded.
8. WHEN transient idempotency or local-action data exceeds approved retention, THE System SHALL delete the transient payload while preserving only the minimum permitted reconciliation or audit evidence.

### Requirement 41: Historical Paper Migration and Reconciliation

**User Story:** As a Tenant migration owner, I want paper records imported as clearly identified evidence, so that historical data is useful without fabricating digital events.

#### Acceptance Criteria

1. WHEN a historical paper card or overtime form is imported, THE System SHALL store the source period, importer, verifier, quality status, field mapping, protected-file reference, and file digest.
2. WHEN transcribed data is created from a historical artifact, THE System SHALL label the data as migrated and human-verified or awaiting verification.
3. WHERE OCR assistance is enabled, THE System SHALL keep OCR output unverified until an authorized human verifies the transcription.
4. THE System SHALL prevent historical import from fabricating digital Dual_Attendance_Confirmation, QR_Challenge, Checkout_Evidence, Worker_PreWork_Consent, Supervisor_Verification, or Worker_Acknowledgement events.
5. WHEN a completed source-artifact example informs field mapping, THE System SHALL use generic placeholders and exclude sample personal data.
6. WHEN payroll go-live approval is requested, THE System SHALL require reconciliation of Worker/day totals, overtime categories, incentives, allowances, advances, deductions, gross, and net against approved legacy results.
7. WHEN a reconciliation difference is identified, THE System SHALL require a reason-coded disposition.
8. THE System SHALL keep shadow payroll results Estimated/Not Finalized and prevent shadow results from publishing Payslips.
9. IF import volume, OCR policy, dual-verification policy, or historical payroll scope remains undecided, THEN THE System SHALL treat those decisions as migration prerequisites.

### Requirement 42: Implementation and Rollout Decision Gates

**User Story:** As an implementation owner, I want open decisions tracked as explicit gates, so that delivery teams do not mistake proposals for confirmed business policy.

#### Acceptance Criteria

1. WHEN implementation-start approval is requested for a dependency, THE System SHALL require identification of approved and still-open programming language, framework, database, cloud, event infrastructure, identity provider, notification provider, signing provider, report renderer, and property-testing library choices applicable to that dependency.
2. IF a technology or vendor choice remains open, THEN THE System SHALL require technology-neutral contracts and a recorded implementation prerequisite.
3. WHEN Tenant rollout approval is requested, THE System SHALL require approval of the exact Permission catalog, default Roles, risk tiers, approval counts, thresholds, and Break_Glass_Session policy.
4. WHEN attendance rollout approval is requested, THE System SHALL require approval of meeting timing, grace periods, late or absence treatment, and confirmation-time disagreement policy.
5. WHEN QR checkout rollout approval is requested, THE System SHALL require approval of Checkpoint model, challenge lifetime, device and location evidence, offline tolerance, privacy notice, and retention.
6. WHEN Worker acknowledgement rollout approval is requested, THE System SHALL require approval of acknowledgement cadence, batch size, cutoff, and signing-authentication level.
7. WHEN payroll rollout approval is requested, THE System SHALL require approval of Payroll_Period cutoff, time-zone boundary, approval calendar, publication timing, finalization scope, formula configuration, blocker policy, and Worker estimate visibility.
8. WHEN production-operations approval is requested, THE System SHALL require approval of availability, latency, throughput, recovery, support, browser, device, accessibility, localization, notification, retention, and data-residency targets.
9. IF branding, custom domains, administration-surface installability, PWA distribution, or app-store wrapping remains undecided, THEN THE System SHALL record the decision as non-binding and avoid inferring product scope.
10. THE System SHALL exclude post-final correction, off-cycle payroll, and revised Payslip workflows from implementation tasks for this feature version.
11. WHEN an implementation task or automated test is created, THE System SHALL require traceability to one or more acceptance criteria in this document.

### Requirement 43: Correctness and Testability

**User Story:** As an engineering quality owner, I want automated tests derived from observable requirements, so that work-to-payroll invariants remain valid across broad inputs and failure orderings.

#### Acceptance Criteria

1. WHEN automated tests are designed, THE System SHALL require the test suite to combine example-based, edge-case, integration, security, state-machine, accessibility, performance, and property-based tests according to behavior type.
2. WHEN a property-based test is implemented, THE System SHALL require the test suite to run at least 100 generated iterations unless a documented test-cost rationale approves a different count.
3. WHEN a property-based test is implemented, THE System SHALL require the test suite to identify the feature name, property number, property statement, and validated requirement references.
4. WHEN property generators are implemented, THE System SHALL require the test suite to include Tenants, calendars, 28–31 day months, time zones, daylight-saving transitions where applicable, assignments, authorization windows, record histories, compensation intervals, monetary values, retries, and concurrency orderings relevant to the property.
5. WHEN a generated property test fails, THE System SHALL require the test suite to retain the minimized counterexample as a regression fixture.
6. WHEN state-machine tests run, THE System SHALL require the test suite to generate valid and invalid calendar, overtime, Work_Record, batch acknowledgement, offline, review, and payroll transitions.
7. WHEN idempotency tests run, THE System SHALL require the test suite to verify retry convergence and rejection of key reuse with a different digest.
8. WHEN concurrency tests run, THE System SHALL require the test suite to verify Expected_Version conflicts, digest revalidation, finalization locking, and absence of partial finalization.
9. WHEN Tenant-isolation security tests run, THE System SHALL require the test suite to substitute cross-Tenant identifiers across APIs, objects, caches, queues, support scope, exports, and offline actions.
10. WHEN privacy tests run, THE System SHALL require the test suite to verify absence of salary, identifiers, exact locations, push endpoints, and evidence payloads from unauthorized responses, general logs, metric labels, notifications, and caches.
11. WHEN PWA tests run, THE System SHALL require the test suite to verify installed and browser parity, update recovery, capability denial, offline non-finality, reconciliation outcomes, and sensitive-data purge.
12. WHEN accessibility tests run, THE System SHALL require the test suite to combine automated checks with manual assistive-technology review for critical workflows under the approved accessibility target.
13. WHEN finalization fault-injection tests run, THE System SHALL require the test suite to verify that an attempt produces either a complete Payroll_Snapshot with all Payslips or no new snapshot and Payslip artifacts.
14. WHEN export round-trip or parser functionality is introduced by a selected integration contract, THE System SHALL require the test suite to include grammar validation, a corresponding pretty-printer or serializer, and parse-print-parse equivalence for valid supported values.
