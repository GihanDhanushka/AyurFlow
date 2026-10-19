# Agent instructions for AyurFlow

Read [docs/implementation-plan.md](docs/implementation-plan.md), [docs/decisions.md](docs/decisions.md) and the module README before implementing a task. For Member A work, also read [docs/member-a-tasks.md](docs/member-a-tasks.md). The SRS v1.0 and team task plan listed in [docs/source-documents.md](docs/source-documents.md) define requirements and task IDs; explicit SRS rules and reviewed decisions take precedence over illustrative diagrams. Treat undecided items as gates, not approved defaults.

## Repository state and ownership

- This repository is a scaffold. Health/readiness, a migration runner, identity tables and placeholder pages exist; protected business routes and AT01–AT17 evidence do not.
- Member A owns shared Go/API, auth, authorization, audit mechanism and deployment/backup work; B patients/clinical; C scheduling; D inventory/preparation; E billing; F shared frontend/integration/UAT. Every feature owner supplies its own schema, backend, UI, permissions, audit calls and tests. Reviewer pairs: A–F, B–C, D–E.
- Keep one accountable owner and reviewer per task. Coordinate migration numbers, shared files and API changes with affected owners. Do not edit another module's schema to avoid a handoff.

## Build rules

- Use Go/Gin and pgx v5/pgxpool for backend rules and PostgreSQL transactions; Next.js App Router calls `/api/v1` and never opens the database. Add working code in module handlers/services/repositories, with parameterized SQL and forward-only numbered migrations. Never edit an applied migration.
- Enforce authentication, CSRF, role **and record** access in the backend for each protected read/write/report/export. Administrator alone is not clinical permission. Do not expose feature routes before their authorization behavior exists.
- Insert critical audit events within the same transaction as the business mutation. A failed audit write must roll back the mutation. Logs and audit views must omit note bodies, passwords and session secrets.
- Store instants in UTC and display Asia/Colombo. Use integer LKR minor units and fixed-precision decimal grams under the approved rounding rule. Keep confirmed/signed/issued history immutable and retain references when deactivated.
- Preserve the required transaction/race semantics for booking, plan decisions, consumption and receipts. Use operation IDs for consumption and receipt retry safety. A failed mutation leaves the prior committed state intact.
- Use synthetic data only in development, tests, demonstrations and external AI tools. No real patient data or secrets in Git. No public signup, gateways, FEFO, reservations, refunds, recursive recipes, multi-clinic features or offline sync without an approved scope change.

## Delivery and checks

- Before editing, identify the TASK ID, SRS IDs, owner/reviewer, dependencies and acceptance tests. Implement small integrated slices: migration, backend behavior, permissions/audit, UI, OpenAPI and meaningful tests. Update documentation that describes changed behavior.
- Run relevant checks with `./scripts/check.ps1` or the Go and frontend commands in the README. Use real PostgreSQL for concurrency/integrity proof. Record what ran and what could not run; scaffold build checks never count as AT acceptance.
- Keep `docs/api/openapi.yaml` truthful: document only implemented, reviewed endpoints. Record each AT result with setup, steps, expected/actual outcome, tester/date/version/environment/evidence. Failed tests remain open defects.
- Do not assert clinic approval, real-data authorization, backup recovery or pilot readiness from code alone. Review DEC01–08 and the release gates in the implementation plan before marking a task done.
