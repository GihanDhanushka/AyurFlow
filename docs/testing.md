# Verification plan and evidence

Starter automated checks cover Go build/vet, health-vs-readiness behavior, frontend ESLint, TypeScript and production build. These are foundation checks, not pilot acceptance.

TASK-AUTH-01 implementation check (19 October 2026, Windows, Go 1.27.1): `go fmt ./...`, `go vet ./...`, `go test -count=1 ./...` and `go build ./cmd/api ./cmd/migrate ./cmd/seedstaff` passed. Tests cover Argon2id round-trip, wrong passwords, unique salts, malformed/oversized hashes and development-only seed input validation. A live PostgreSQL seed/repeat test and F review remain outstanding; AT01 remains pending because login, sessions, CSRF and role/record checks are not implemented.

TASK-AUTH-02 local checks (21 October 2026, Windows): `go test -count=1 ./...`, `go vet ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seedstaff` and integration-test compilation passed. Unit tests cover secret hashing, cookie/origin/CSRF handler behavior and required-policy configuration. `go test -tags=integration ./internal/auth` contains a disposable PostgreSQL lifecycle case for throttling, session replacement, logout, idle/absolute expiry and disabled-account rejection. The database case could not run here because `AUTH_TEST_DATABASE_URL` and PostgreSQL are unavailable. `go test -race` could not build because the installed C compiler reports `64-bit mode not compiled in`. DEC07 and F review remain pending. AT01 is not passed; AUTH-03 role/record access and general CSRF checks remain to be implemented.

The backend GitHub Actions job now starts a disposable PostgreSQL 17 service and runs the AUTH-02 integration test with synthetic credentials. This checks real migration and session SQL on each push/PR. Local execution is still available with a dedicated disposable database and `AUTH_TEST_DATABASE_URL`; do not point it at a shared or production database. CI results are recorded only after the remote job completes.

AUTH-02 PostgreSQL follow-up: the [workflow for merged PR #12](https://github.com/thimira20011/AyurFlow/actions/runs/39788040138), commit `812a93f`, completed successfully on 22 October 2026. This proves the existing synthetic session lifecycle case against PostgreSQL; it does not close DEC07 or full AT01.

## TASK-AUTH-05 audit mechanism checks (22 October 2026)

- **Scope/version:** AUD-01–02 and shared PAT-02 support on `codex/auth-05-audit`, implementation commit `2831d76`; owner A, reviewer F. This is shared-mechanism evidence, not a complete AT13 result.
- **Tester/environment:** Codex automation; local Windows/Go checks, plus GitHub Actions Ubuntu/Go 1.27.1/PostgreSQL 17.11 for database/race checks.
- **Setup:** disposable PostgreSQL with synthetic staff and a generated isolated schema per fixture; migration `000004` is applied through the real migration runner. Synthetic business mutations exist only in test schemas.
- **Steps:** commit a mutation with its audit event; force an audit INSERT failure; submit invalid metadata after a mutation; roll back a transaction containing a successful audit write; attempt UPDATE/DELETE/TRUNCATE; bypass Go validation with invalid SQL metadata; test scoped pagination, denied objects and disabled accounts.
- **Expected:** business/audit rows commit together or both roll back; failed audit leaves no committable transaction; stored metadata includes actor/action/object/time/reason and changed-field names; history changes and free-text metadata fail; reads require explicit permission and target access.
- **Actual:** local `go test -count=1 ./...`, `go vet ./...`, command builds and integration-test compilation passed. The [branch workflow](https://github.com/thimira20011/AyurFlow/actions/runs/39791162419) passed both backend and frontend. Its [backend job logs](https://github.com/thimira20011/AyurFlow/actions/runs/39791162419/job/119867241400) confirm all four PostgreSQL audit cases (`TestWriteAtomicityIntegration`, `TestAuditImmutabilityIntegration`, `TestAuditMetadataConstraintsIntegration`, `TestReaderScopeAndPermissionsIntegration`) ran and passed with `-race`, alongside auth regression tests. The forced audit failure left no business mutation; history mutation and unstructured metadata attempts were rejected; scoped pagination and disabled-account denial passed. Local PostgreSQL and Windows race execution remain unavailable, with the remote Linux job supplying those checks.
- **Remaining acceptance:** F review, AUTH-03 adapter/HTTP permissions, module-specific mutation/audit calls and direct API checks. AT13 remains Pending.

Scaffold verification on 2 October 2026: Go formatting, vet, tests and both command builds passed; frontend lint, typecheck and production build passed; JSON and YAML syntax checks passed. Docker is not installed on the verification machine, so container startup, SQL execution and restart persistence remain unverified. No AT01-AT17 result is marked passed.

| Test | Scope | Proposed owner | Status |
| --- | --- | --- | --- |
| AT01 | Accounts, sessions, CSRF, roles and record restrictions | A | Pending |
| AT02 | Optional identity, duplicate race, Unicode search, demographic audit | B | Pending |
| AT03 | Draft/sign locking race and preserved dated addenda | B | Pending |
| AT04 | Ordered plan and clinician approval; current-stage clearance | B/C | Pending |
| AT05 | Hold/resume/cancel, retained history and hold/start race | B/C | Pending |
| AT06 | Availability, adjacent/overlapping booking races, atomic rescheduling | C | Pending |
| AT07 | Assigned therapist delivery and cancellation boundaries | C | Pending |
| AT08 | Batch receipt/block and expiry/eligibility | D | Pending |
| AT09 | Recipe scaling, immutable snapshot, no preparation deduction | D | Pending |
| AT10 | Consumption retries, races, insufficiency and complete rollback | D | Pending |
| AT11 | Bidirectional actual traceability and report permissions | D/F | Pending |
| AT12 | Locked invoice, deposit, balance, receipt retries/races and retention | E | Pending |
| AT13 | Reference/report access, critical audit and failed-audit rollback | A/all | Pending |
| AT14 | 10 minutes, 5 users, 500 patients, 1000 sessions, 200 batches; p95 <2s on three queries | F | Pending |
| AT15 | 1366×768 and 1024×768, keyboard, labels, Unicode, errors | F/all | Pending |
| AT16 | Outage re-entry and encrypted isolated restore of four record categories | A/F | Pending |
| AT17 | Clean startup, migrations, seed, restart persistence, docs and handover | A/F | Pending |

For each result record setup, steps, expected/actual behavior, tester, date, application version, environment and evidence location. Go PostgreSQL tests live beside module code as `*_integration_test.go` under the `integration` build tag; CI discovers all packages. Store browser scenarios in `tests/e2e`. Record unresolved defects rather than marking them passed.

## Dev integration repair (23 October 2026)

- **Tasks/scope:** SET-05 CI maintenance; BILL-01 receipt integrity support (BIL-01–06, partial AT12); AUTH-05/AUD-01–02 regression support (partial AT13). Owner A, reviewer F; E/D review for billing changes. Branch `codex/dev-ci-integration-fixes`, based on `dev` at `a9eeb93`.
- **Tester/environment:** Codex automation, Windows, Go 1.27.3, synthetic fixtures only.
- **Setup/steps:** restore a consistent frontend lockfile and clean-install it; run frontend lint, generated-route type checking and production build; run Go unit tests, vet, all-package builds, tagged integration compilation and vulnerability scanning. Add real-migration isolation/repeat checks plus changed-payload, concurrent retry/payment, overflow and failed-audit rollback tests against actual billing services.
- **Expected/actual local:** all listed local checks pass; `govulncheck` reports zero reachable vulnerabilities. PostgreSQL and race execution require the Linux workflow because local PostgreSQL/Docker and a supported Windows race compiler are unavailable. Remote execution result pending the branch workflow.
- **Limits:** this evidence does not close AT12/AT13, patient FK, AUTH-03, DEC04/07, clinic review or operational release gates. See [review findings](dev-integration-review.md) and [incoming test instructions](ci-and-test-guide.md).
