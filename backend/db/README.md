# Database foundation

Run `go run ./cmd/migrate` from `backend/` with `DATABASE_URL` set. Forward migrations are embedded in the executable, applied in filename order, recorded in `schema_migrations`, and protected by a transaction-level advisory lock. The pending batch runs in one transaction. Run migrations explicitly before starting the API.

Use six-digit sequential filenames. `000001`–`000004` are allocated to identity, roles, throttling and audit integrity. The next available number on this feature branch is `000005`; recheck the latest `dev` before assigning it. Coordinate numbers across owners and never edit an applied migration. Review backup/restore before destructive changes; there is deliberately no automatic down migration.

The schema contains identity, role, server-session, throttle and audit storage. Optional session behavior and the transactional audit writer are implemented; audit UPDATE/DELETE/TRUNCATE operations are rejected by database triggers. The audit reader requires an explicit authorization adapter; AUTH-03 still supplies shared role/record protection. No credentials are seeded automatically; `cmd/seedstaff` requires explicit synthetic development input. `btree_gist` is enabled for future scheduling exclusion constraints.

Clinical, patient, scheduling, stock and invoice schemas remain pending agreed field rules. Use `timestamptz` for event times, PostgreSQL `date` for expiry, integer minor units for LKR, and `numeric` with an agreed decimal scale for grams. Enforce half-open booking intervals and uniqueness for retries in their module migrations.
