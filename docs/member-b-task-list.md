# Member B task list

Member B owns patient records, consultations, care plans, and clinician decisions. Member C reviews B's work and owns scheduling; Member B reviews C's scheduling work. This checklist follows the SRS v1.0 and team task plan dated 1 October 2026. All items are pending in the current scaffold. Weeks are relative to the team's implementation kickoff.

## Sprint 1 — decisions, contracts, and registration (weeks 1–2)

- [ ] **TASK-SET-02 — Confirm clinic rules (DEC01–02).** With the clinic reviewer, record patient/contact and consultation fields, identity normalization, clinician record access, the one pathway, stage/session checklists, and completion conditions. Record unresolved decisions with owners and gates; do not treat provisional rules as approved.
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
- [ ] **TASK-QA-01/03/04/06 — Supply B's examples and fix B-module integration defects.** Add synthetic patient/clinical cases to the requirements checklist; support the patient-to-receipt walkthrough, audit/access review, clinic UAT, and defect resolution.
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

Source documents: `C:/Users/ASUS/Downloads/AyurFlow_Software_Requirements_Specification.pdf` and `C:/Users/ASUS/Downloads/AyurFlow_Team_Task_Plan.pdf`. The repository's `docs/decisions.md` and `docs/testing.md` track project-wide decisions and acceptance evidence.
