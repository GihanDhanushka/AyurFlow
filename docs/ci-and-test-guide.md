# CI and test instructions for incoming tasks

Member A maintains shared CI; each module owner writes its behavior tests and its reviewer checks the evidence. Passing CI is not AT01–AT17 acceptance or clinic approval.

## Required checks

| Job | Checks |
| --- | --- |
| `backend` | Parse/format all Go files, verify tidy module files, vet, all unit tests with race detection and shuffled order, static analysis, build every package/command. |
| `integration` | Start disposable PostgreSQL 17, migrate real schemas, execute all Go packages with the `integration` tag and race detector, reject skipped tests or an empty integration suite, retain JSON results. |
| `security` | Run a pinned `govulncheck` against reachable backend code. |
| `frontend` | Clean install from the committed lockfile, verify it remains unchanged, lint, TypeScript and production build. This does not claim browser/E2E coverage. |

The workflow runs on pushes and PRs. A syntax error, dependency mismatch, database failure or test failure must fail its job; fix the cause rather than removing the gate. Coverage/results artifacts are diagnostic evidence, not a coverage target or a claim that unimplemented workflows are tested.

## Add backend tests with each feature

1. Put ordinary tests in `backend/internal/<module>/*_test.go`. CI discovers them through `go test ./...`; do not add per-module paths to the workflow.
2. Put real database tests in `*_integration_test.go`, starting with `//go:build integration`. Name the top-level test `Test...Integration` so the results gate can identify executed integration cases.
3. Use `testutil.Postgres(t)` from `ayurflow/backend/internal/testutil`. It applies the real migrations and returns a pool confined to a generated isolated schema. Each independent case should create its own fixture; concurrent requests for the same business record must intentionally share one fixture.
4. Set `TEST_DATABASE_URL` to an explicitly disposable local database when running integration tests. CI supplies it automatically. Missing configuration fails in CI; local tests may skip with a clear message. Auth/audit/billing legacy variable names remain supported for local workflows.
5. Do not add always-skipped placeholders to a required suite. Tests for an unavailable external service belong to a separately configured suite. Do not use external clinic records, real credentials or production databases.
6. Exercise observable behavior: authorized success, forbidden direct access, failure preserving prior state, race outcomes, and retries where required. A test that merely recalculates the implementation's formula is insufficient.
7. Record the task/SRS IDs, setup, steps, expected/actual result, version, environment and evidence. Keep incomplete/failed AT cases open.

```go
//go:build integration

package yourmodule

import (
    "testing"
    "ayurflow/backend/internal/testutil"
)

func TestYourRequiredBehaviorIntegration(t *testing.T) {
    pool := testutil.Postgres(t)
    // Seed synthetic prerequisites, call the real service and assert its outcome.
    // Never create an alternate mock schema to bypass the real migrations.
    _ = pool
}
```

This snippet demonstrates fixture setup, not a completed behavior test. Before committing, replace the comment with real service calls and assertions. The helper migrates a shared extension baseline and an isolated schema; cleanup removes only that generated schema. An explicit disposable database with schema/migration privileges is required.

## Upcoming task coverage

| Tasks | Required behavior evidence when implemented |
| --- | --- |
| AUTH-03 | Direct unauthenticated/forbidden requests, CSRF rejection, role **and** record checks, administrator-only denial of clinical text. |
| AUTH-04/06 | Controlled account/role changes, disable revoking sessions, audit rollback and complete AT01 evidence. |
| PAT/CON/CLN | Identity uniqueness races, signing immutability, addenda/history, ordered transitions and hold/start races. |
| SCH | Room/therapist overlap conflicts, adjacent slots, atomic reschedule and assigned-therapist clearance. |
| INV/PRE | No preparation deduction, actual consumption rollback, insufficient stock, operation retries and traceability. |
| BILL | Changed receipt payload conflicts, identical/concurrent retry safety, overpayment/overflow rejection, audit rollback, patient reference/access and issued-history immutability. Current tests cover only the implemented subset. |
| OPS | Fresh migrations, repeat migrations retaining data, restart persistence, backup/isolated restore and deployment checks. CI schema tests alone do not prove operational recovery. |

Browser tests should accompany functional screens once their real backend paths exist. Add the browser test command/job in that feature PR and exercise the actual API; do not infer login or billing behavior from placeholder pages or a successful Next.js build.

## Local commands

From `backend`: `go test ./...`, `go vet ./...`, `go build ./...`; with a disposable `TEST_DATABASE_URL`, run `go test -tags=integration -shuffle=on -count=1 ./...`. On a supported C toolchain add `-race`; the Windows toolchain on the original verification machine cannot build 64-bit cgo, so Linux CI supplies that check.

From `frontend`: use `npm ci` before lint/typecheck/build to detect lockfile problems. Use the Node version pinned in the workflow. `scripts/check.ps1` runs local checks and database tests when `TEST_DATABASE_URL` is configured.
