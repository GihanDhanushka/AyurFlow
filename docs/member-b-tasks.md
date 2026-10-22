# Member B task list

The current user has confirmed they are **Member B**. Their personal name is still unrecorded. B owns patient records, consultations, care plans, and clinician decisions. C reviews B's work and owns scheduling; B reviews C's scheduling work. This one file tracks B's work from SRS v1.0 and the 1 October 2026 team task plan. Weeks are relative to implementation kickoff. The project-wide board tracks all 64 tasks separately.

## Source review and existing work

- The available SRS v1.0 defines optional NIC/passport, safe identity uniqueness, Unicode search, immutable signed consultations with addenda, one ordered approved pathway, and current-stage clearance at session start. Member B's primary acceptance evidence is AT02–05; B also reviews C's AT06–07 scheduling integration. Clinical clearance does not depend on invoice balance.
- The team plan contains 64 tasks in five two-week sprints. Sprint 1 requires a deployed login → registration → search flow; Sprint 2 requires consultation → approved plan → conflict-safe booking → authorised start. B's feature work depends on A's authentication, record access and audit foundation; C's session work; and F's shared UI.
- Integrated `dev` has Go/Gin and Next.js scaffolds, foundation migrations, health/readiness, Argon2id helpers and a development-only synthetic staff seed. Patient and clinical module folders and screens are placeholders. Protected patient/clinical APIs, feature migrations, browser flows and AT02–05 evidence are absent.
- DEC01–02 clinic rules and DEC08 session/plan details still need reviewers and dated decisions. The SRS and task plan were available for this review; the reduced-scope proposal referenced by them was not available in this workspace. No clinic approval or real-data authorisation is implied by these documents.

## Current work and branch order

`TASK-SET-02` has a recorded provisional development contract on `member-b/confirm-clinic-rules`. Its DEC01–02 clinic review record is in [clinic-rules-review.md](clinic-rules-review.md). The user directed this planning task to finish without waiting for clinic signoff; clinic answers and affected feature UAT remain pending. No patient or clinical feature endpoint has been implemented. Every other B-owned feature below remains Pending; none of AT02–05 has passed.

For every new B-owned TASK ID, start a **new branch** from the latest integrated `dev` named `member-b/<description-of-work>`; use lowercase words and hyphens because Git branch names cannot contain spaces. Keep one TASK ID per branch and put its branch/PR in the board. The names below describe the work and are suggestions until their branches are created. Shared tasks owned by A, C or F use that task owner's branch, with B contributing there. The existing `TASK-SET-02` branch already follows this pattern.

| B-owned task | Status | Branch | Depends on / acceptance gate |
| --- | --- | --- | --- |
| TASK-SET-02 | Done for development planning | `member-b/confirm-clinic-rules` (existing) | DEC01–02 clinic decisions remain gates for affected UAT |
| TASK-PAT-01 | Pending | `member-b/patient-schema` | SET-02/06; agreed fields and uniqueness |
| TASK-PAT-02 | Pending | `member-b/registration-search` | PAT-01, AUTH-03/05; AT02 behavior |
| TASK-PAT-03 | Pending | `member-b/patient-screens` | PAT-02, UI-02/06; registration and lookup |
| TASK-CON-01 | Pending | `member-b/consultation-drafts` | PAT-01, SET-02, AUTH-03; draft without plan |
| TASK-CON-02 | Pending | `member-b/signing-addenda` | CON-01, AUTH-05; AT03 immutability |
| TASK-CON-03 | Pending | `member-b/consultation-screens` | CON-02, UI-02/06; history UI |
| TASK-CON-04 | Pending | `member-b/patient-consultation-tests` | PAT-02/03, CON-01/03; AT02–03 and C review |
| TASK-CLN-01 | Pending | `member-b/pathway-schema` | SET-02/06, PAT-01; approved pathway |
| TASK-CLN-02 | Pending | `member-b/plan-approval` | CLN-01, AUTH-03/05; AT04 |
| TASK-CLN-03 | Pending | `member-b/hold-resume-cancel` | CLN-02, SCH-01, DEC08; AT05 |
| TASK-CLN-04 | Pending | `member-b/session-clearance` | CLN-02/03, SCH-03; start/hold race |
| TASK-CLN-05 | Pending | `member-b/plan-screens` | CLN-02/04, UI-02/06; clinician actions |
| TASK-CLN-06 | Pending | `member-b/clinical-transition-tests` | CLN-03/05, SCH-04; AT04–05 and C review |

Tasks may be prepared while decisions are pending, but they cannot be marked Done until their dependencies, integration evidence and reviewer checks pass. Do not use provisional pathway or field choices as clinic approval.

## Sprint 1 — decisions, contracts, and registration (weeks 1–2)

- [x] **TASK-SET-02 — Record development rules and open clinic gates (DEC01–02).** The synthetic-development contract and unresolved clinic decisions with owners/gates are in [clinic-rules-review.md](clinic-rules-review.md). The user waived waiting for signoff on this planning task; DEC01–02 clinic decisions and affected feature UAT remain pending. Do not treat provisional rules as approved clinic policy.
- [ ] **TASK-SET-06 — Supply B's schema and API contracts.** Agree patient, encounter/addendum, plan/stage, and clinical-clearance IDs and relationships with A and C. Define request/response fields, role and record restrictions, validation, conflicts, and migration numbers in the shared OpenAPI contract. A coordinates the shared contract.
- [ ] **TASK-UI-02 — Design patient and care-plan screens with F.** Show draft versus signed notes, addenda, approval history, current stage, and hold/cancel states. Use F's shared components when available.
- [ ] **TASK-PAT-01 (PAT-01) — Create the patient schema and validation.** Store a server-generated internal ID, name, agreed contact fields, and optional NIC/passport. Enforce uniqueness of a supplied *normalized* identity in PostgreSQL; allow patients without an identity number.
- [ ] **TASK-PAT-02 (PAT-01–02) — Implement registration, search, and demographic updates.** Search by internal ID, name, supplied identity, and agreed contact fields. Bound or paginate results. Enforce backend role/record access and write demographic audit events in the same transaction as changes. Handle simultaneous duplicate registrations safely.
- [ ] **TASK-PAT-03 — Build registration, search, and patient-detail screens.** Show field errors and clear success/failure states; preserve Unicode names. The Sprint 1 exit flow is authorised login → registration → search on the deployed application.

## Sprint 2 — consultations and clinical pathway (weeks 3–4)

- [ ] **TASK-CON-01 (CON-01) — Create consultation/addendum schema and draft actions.** Allow an authorised clinician to create and edit a consultation without a care plan. Enforce record access on reads and writes.
- [ ] **TASK-CON-02 (CON-02–03) — Implement signing and dated addenda.** Signing records clinician and time and locks the original text. Reject stale edits and competing signatures. Append corrections with author, time, and reason; never replace the signed original. Keep required audit writes atomic with the action.
- [ ] **TASK-CON-03 — Build consultation editor and history.** Distinguish draft and signed states and show dated addenda in order. Keep full clinical text unavailable to non-clinical roles.
- [ ] **TASK-CON-04 — Verify patient and consultation behavior.** Complete AT02–03 with synthetic data, direct forbidden API requests, a duplicate-identity race, Unicode search, a signing race, and preserved original/addendum history. C reviews.
- [ ] **TASK-CLN-01 (CLN-01) — Encode the agreed ordered pathway.** Add plan/stage schema with status, current/completed stages, decision actor/time, and the clinic-approved stage labels and checklists. Keep plan status separate from pathway stage.
- [ ] **TASK-CLN-02 (CLN-02) — Implement plan creation, approval, and guarded advancement.** Only authorised clinicians may prescribe or approve. Preserve approver/time; reject missing predecessor approval, stage skipping, and invalid transitions.
- [ ] **TASK-CLN-03 (CLN-03–04) — Implement hold, resume, and cancel with C.** Record clinician decisions and reasons. A hold blocks advancement and *new* session starts without erasing prior approvals. Cancellation is terminal, cancels future unstarted sessions, and retains completed sessions, approvals, and consumption.
- [ ] **TASK-CLN-04 (CLN-05–06) — Provide shared session clearance with C.** At session start, check the latest plan status, current stage, and required approval. Booking alone does not grant clearance; invoice balance is irrelevant. Make hold/cancel/approval versus start races safe with transactional locking or an equivalent conflict rule.
- [ ] **TASK-CLN-05 — Build the care-plan timeline and clinician actions.** Display completed stages and decision history, with usable approval, hold, resume, and cancel actions only for permitted users.
- [ ] **TASK-CLN-06 — Verify clinical transitions with clinic examples.** Complete AT04–05, including unauthorised approval, stage skipping, noncurrent-stage start, retained history, and hold/cancel versus start races. C reviews.

## Integration, review, and release support

- [ ] **TASK-SCH-04 — Support C's session-delivery integration.** Agree the clinical-clearance call and the behavior of a session already in progress when a plan is held or cancelled. The assigned therapist may record delivery, but cannot approve a stage or prescribe a plan.
- [ ] **TASK-SCH-06 — Review C's scheduling tests.** Check AT06–07 cases that touch plan/stage links, clearance, cancellation, and concurrent session start. C owns scheduling and its resource constraints.
- [ ] **TASK-QA-01/03/04/06 — Supply B's examples and fix B-module integration defects.** B's [synthetic AT02–05 case contribution](qa/member-b-examples.md) and [fixture data](../tests/fixtures/member-b-cases.json) are drafted on `member-b/clinical-test-examples`; F/A review and execution evidence remain pending. Continue supporting the patient-to-receipt walkthrough, audit/access review, clinic UAT, and defect resolution.
- [ ] **TASK-UI-07 / AT15 — Check B's screens.** Verify Unicode round trip, keyboard focus, labels, validation messages, and 1366 × 768 and 1024 × 768 layouts with F.
- [ ] **TASK-OPS-03 — Write B-module usage guidance.** Document registration, search, consultation signing/addenda, care-plan decisions, and how staff handle failures; have another member verify the instructions.

## Completion rules for every B-owned feature

- [ ] Add the relevant numbered migration, Go handler/service/repository behavior, parameterized SQL, and OpenAPI request/response/error contract.
- [ ] Enforce backend authentication, role **and record** access, CSRF for mutations, validation, and an audit event in the same transaction as each critical change. A supplies the shared auth/audit mechanisms.
- [ ] Add the feature screen or documented interface and meaningful automated tests, including real PostgreSQL tests for uniqueness and state-transition races.
- [ ] Record AT02–05 setup, steps, expected and actual result, tester, date, environment, version, and evidence. Obtain C's review and verify the feature in the integrated app before marking it complete.

## Decisions and dependencies to track

- **DEC02:** Clinic-approved fields, normalization, clinician record access, one pathway, and completion rules are required before B's affected features are accepted.
- **DEC08:** Confirm resumption, terminal cancellation, in-progress session completion during hold, and concurrency semantics with C and the clinic reviewer.
- **Member A:** Authentication, session/CSRF controls, backend permission framework, audit mechanism, and shared API conventions must be integrated before protected routes are released.
- **Member C:** Session schema, booking, therapist delivery, and plan/session transition contract must agree with B's clinical model.
- **Member F:** Shared UI patterns and end-to-end integration support B's screens; B still owns their implementation and tests.

Source documents: SRS v1.0 and the team task plan listed in [source-documents.md](source-documents.md). [decisions.md](decisions.md) and [testing.md](testing.md) track project-wide decisions and acceptance evidence. The reduced-scope proposal is referenced by those documents but was not available in this workspace during this review.
