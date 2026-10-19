# Member A task list

**Role:** shared Go/API foundation, authentication, role and record access, audit mechanism, deployment and backups. **Review partner:** Member F. Weeks below are relative to implementation kickoff, not calendar dates. Check a box only when its evidence and review are recorded.

Sources: *AyurFlow Team task plan* (1 October 2026), especially pp. 2–5 and 11–14; *AyurFlow SRS* v1.0, especially pp. 4–6, 13, 17–21. See [the implementation plan](implementation-plan.md) for dependencies and release gates.

## Current baseline (repo inspection, 19 October 2026)

- The Go/Gin API, pgx pool, numbered migration runner, `/api/v1/health` and `/api/v1/ready` exist. Next.js has a shared shell and placeholder module pages. Compose and CI/check scripts exist.
- Migrations `000001` and `000002` create staff, role, session and audit tables and role reference rows. They do **not** implement credentials, session behavior, authorization or immutable audit operations.
- OpenAPI documents only health/readiness. No business routes, synthetic seed command, HTTPS pilot deployment, backup automation or restore evidence are present.
- `TASK-SET-05` is **in progress**: each member must still start the application locally and baseline checks must pass. All `TASK-AUTH-*`, pilot operations and AT01–AT17 remain pending. Do not equate scaffold files with passed acceptance tests.

## Weeks 1–2: contracts and first deployed flow

- [ ] **TASK-SET-01 — coordinate ownership and reviews.** Put all 64 task IDs on one board with a named owner, reviewer, status and evidence link; record weekly integration and two-week sprint reviews. A coordinates; every member confirms their tasks. **Done:** no task has an anonymous owner or reviewer.
- [ ] **TASK-SET-04 — settle shared API and security conventions (DEC07).** Review ID formats, UTC instants and Asia/Colombo display, integer LKR minor units, decimal grams policy handoff, pagination, error envelope/status codes, operation-ID retry behavior, cookie/CSRF contract, session idle and absolute expiry, throttling and recovery. Write request/response examples and decisions in `docs/decisions.md` and `docs/api/`. **Done:** all owners review the contract; pending clinic choices stay marked pending.
- [ ] **TASK-SET-05 — finish the scaffold.** Verify Go/Gin, Next.js, migrations and checks on every member's machine; provide repeatable, synthetic-only seeds once identity and patient fields are ready. Resolve actual setup failures. **Done:** each member demonstrates local startup; `scripts/check.ps1` passes; instructions and evidence are recorded. F assists.
- [ ] **TASK-SET-06 — coordinate schema and OpenAPI boundaries.** Collect each module owner's fields, foreign keys, actions, permissions and errors; agree migration numbers and inter-module IDs. A coordinates; each owner writes their own module contract/migration. **Done:** reviewed logical schema, API examples and migration order; no undocumented cross-module assumptions.
- [ ] **TASK-AUTH-01 — finish account/role/session foundation.** Review existing `staff_accounts`, `roles`, `staff_roles`, `staff_sessions` migrations; add needed forward migration(s), Argon2id create/verify and controlled synthetic staff seeding. **Done:** no plaintext or committed default credentials; role assignment and password handling have meaningful tests.
- [ ] **TASK-AUTH-02 — implement sessions.** Add login/logout, opaque random session token stored hashed server-side, expiry, revocation, disabled-account revocation, generic login errors and DEC07 throttling. **Done:** expired, logged-out and disabled sessions fail to act; AT01 cases have evidence.
- [ ] **TASK-AUTH-03 — protect every backend route.** Add Secure/HttpOnly/SameSite cookie behavior, CSRF for state-changing cookie requests, role checks and record checks on reads, writes, reports and exports. Define a narrow authorization interface module owners can use. **Done:** direct API tests reject forbidden requests; administrator-only cannot access consultation text.
- [ ] **TASK-AUTH-05 — transactional audit foundation.** Provide a reusable `audit` writer that accepts the caller's PostgreSQL transaction before patient registration is exposed. Restrict audit retrieval; prevent ordinary-user edit/delete. Require module owners to call it for critical mutations. **Done:** forced audit-write failure rolls back the business change; actor/action/object/time/reason are retained without notes or secrets.
- [ ] **TASK-AUTH-04 — account and role UI with F.** Build the login screen for the Sprint 1 flow, then account create/disable and explicit role assignment screens through week 3; coordinate the admin wireframe (`TASK-UI-05`) and shared controls (`TASK-UI-06`). **Done:** authorized administrators can manage accounts; no public signup; disable revokes sessions.
- [ ] **TASK-OPS-01 — deploy the first slice.** Extend development Compose into a reviewed pilot configuration with HTTPS reverse proxy, restricted database access, separate application/migration credentials and secrets outside Git. Pair with B/F on the login → patient registration → search demonstration. **Done:** deployed Sprint 1 journey works with permissions; F verifies. Pilot host choices remain under DEC06.

## Weeks 2–3: administration, audit and identity verification

- [ ] **TASK-AUTH-06 — prove identity and logging behavior.** Execute AT01 with named role accounts, direct forbidden requests, CSRF, expiry, throttling and safe-log checks. F reviews. **Done:** AT01 evidence records setup, steps, actual results, version and reviewer; failures remain open defects.

## Weeks 4–8: integration and reliability

- [ ] **TASK-QA-02 (access reviewer).** Review F's report navigation and D/E's report endpoints for role, record and filter restrictions. **Done:** AT11/AT13 report checks expose only authorized fields.
- [ ] **TASK-QA-04 — audit and permission sweep with F.** Collect each owner's critical actions and test clinical, booking, inventory and billing audit events and direct API boundaries. **Done:** AT13 passes, including failed-audit rollback and preserved history.
- [ ] **TASK-QA-05 (host support).** Supply the declared pilot host/configuration for F's AT14 workload and help diagnose host issues; module owners fix their query paths. **Done:** 10-minute, 5-user run on 500 patients/1,000 sessions/200 batches records per-operation p95 and errors; target is under 2 s for all three named reads.
- [ ] **TASK-OPS-02 — encrypted daily backups and isolated restore (DEC06).** Agree off-host storage, separate key custody, monitoring and restore owner; automate backups and test an isolated restore of a patient, batch, invoice and audit event. **Done:** AT16 evidence records backup age, measured restore duration and inspected records. F/clinic reviewer validates operations.

## Weeks 8–10: clean release and handover

- [ ] **TASK-OPS-03 (A's portion).** Maintain API/setup/migration instructions and the operational outage/re-entry procedure; every owner writes their module guide, F assembles them. **Done:** another member follows the instructions successfully.
- [ ] **TASK-OPS-04 — clean startup and persistence with F.** From an isolated clean environment, apply migrations, load synthetic seeds, restart and verify committed records persist. **Done:** AT17 evidence includes commands, versions, outcomes and reviewed OpenAPI/guides.
- [ ] **TASK-OPS-05 (operations portion).** Record host/funding, backup/key custody, restore owner and support contact under DEC06. Support the team lead/clinic/institution on DEC05 real-data authorization, retention and exports. **Done:** decisions are recorded before pilot/handover; technical completion never implies real-data permission.
- [ ] **TASK-OPS-06 (A's portion).** Teach administration, account recovery and operations in F's training/handover; supply runbook, limitations and support details. **Done:** named recipient, release version, acceptance evidence and responsibilities are recorded.

## Gates and handoffs to track

- **DEC02 (B/clinic):** patient fields and clinical record access are required before finalizing record authorization. **DEC03/04 (D/E/clinic):** units and billing rules must be reflected in shared formats. **DEC05:** required before real patient data. **DEC06:** required before operational handover. **DEC07:** A/team Sprint 1 security/API settings. **DEC08:** plan/session concurrency contract with B/C.
- Give B–E a usable session/authorization/audit contract early. Each feature owner integrates its own migration, UI, permission checks, audit calls and tests; A reviews shared conflicts rather than implementing all module schemas.
- Use `AT01`, `AT13`, `AT16`, `AT17` as A's primary acceptance evidence. Contribute host evidence to `AT14` and permission review to `AT11`. F is the reviewer for A-owned work.
- The Sprint 1 exit is deployed login → registration → search; Sprint 2 is approved plan → safe booking → authorized start; Sprint 3 is the full synthetic patient-to-receipt journey and feature freeze; Sprint 4 is UAT and integrity/restore evidence; Sprint 5 is clean delivery and handover.
