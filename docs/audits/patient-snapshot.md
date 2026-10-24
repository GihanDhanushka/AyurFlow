# Patient snapshot audit repair

**Task:** QA-04 integrity review, PAT-01–02 support. Owner A; reviewer F with B/C patient handoff. Branch `codex/audit-patient-snapshot`, based on AUTH-03 commit `d592199`. Synthetic data only; no clinic approval or AT02/AT13 acceptance claim.

## Approved inventory and repair

`readPatient` fetched demographics/version and optional identity in separate SELECTs. In PostgreSQL's default READ COMMITTED isolation, a concurrent update could commit between them, producing an old name/version with a new identity. The approved behavioral repair uses one parameterized LEFT JOIN and nullable identity fields. Callers, fields, permission checks, transaction boundaries and schema remain unchanged.

This removes one SQL statement per patient read as a consequence of making the snapshot consistent. No latency or throughput improvement is claimed; broader search batching/index work is deferred until measured. Migration files are unchanged. The database guide now correctly reserves the next number after the existing `000005_patients.up.sql` as `000006`, subject to checking latest `dev`.

## Regression setup and expected result

- Use the existing real-migration PostgreSQL helper with a disposable isolated schema and synthetic staff/patient fixtures.
- A test transaction/row wrapper completes the first Scan, then invokes the real audited `Service.Update` on another connection before the read continues. No sleeps or scheduling assumptions are required.
- The expected regression failure with the old two-query implementation is the original name/version combined with the newly committed identity. The single statement must retain the complete original snapshot while a subsequent read sees the complete updated snapshot. Local PostgreSQL execution is unavailable, so this before/after interleaving is not claimed as locally executed evidence.
- Additional cases verify that absent/cleared optional identity and phone remain null and a nonexistent patient still returns `ErrNotFound`.

## Verification

Local Windows checks passed: `go test -count=1 ./internal/patients`, `go vet ./internal/patients`, `go build ./...`, tagged integration-test compilation (`go test -tags=integration -run '^$' ./internal/patients`), and `git diff --check`. Local PostgreSQL/Docker is unavailable; tagged compilation is not execution evidence. PostgreSQL/race CI remains pending; root will record the feature branch workflow result after push. Remaining AT02/AT13, B/C/F review and clinic decision gates remain open.
