# Auth module

Proposed owner: A. Reviewer: F. SRS: IAM-01-05.

Staff accounts, Argon2id passwords, opaque sessions, expiry/revocation, throttling, CSRF and backend role/record checks.

Status: Argon2id hash/verify is implemented in `password.go` for TASK-AUTH-01. New hashes use a random 16-byte salt, 64 MiB memory, 3 passes and 4 threads; encoded parameters are checked and bounded before verification. The 64 MiB/3-pass/4-thread choice follows the second option in the [Go Argon2id documentation](https://pkg.go.dev/golang.org/x/crypto/argon2). A development-only `cmd/seedstaff` command creates explicit synthetic accounts after migrations. Login, sessions, throttling, CSRF, authorization and account-management routes remain pending; no protected route is registered.

When adding those features, services own backend permission checks and transaction boundaries; handlers validate inputs; repositories use parameterized pgx queries. Update migrations, OpenAPI, frontend screens and relevant acceptance tests together. Never log passwords or hash strings.
