# Billing module

Proposed owner: E. Reviewer: D. SRS: BIL-01–06. AT: AT12.

Draft/issued invoices, integer LKR minor units, direct deposits/receipts,
balance protection and idempotent receipt operations.

## Status (TASK-BILL-01 — in progress)

| File | Status | Notes |
|---|---|---|
| `domain.go` | ✅ Implemented | Invoice, InvoiceItem, Receipt, Balance, Report types; all money as int64 LKR minor units |
| `repository.go` | ✅ Implemented | Parameterized pgx queries; FOR UPDATE locking; audit event helper |
| `service.go` | ✅ Implemented | Full invoice lifecycle, idempotent receipt posting, balance derivation |
| `handler.go` | ✅ Implemented | All route handlers; `RegisterRoutes` defined but **not wired** (AUTH-03 gate) |
| `billing_test.go` | ✅ Implemented | Unit tests for validation logic; integration stub documented |
| `000004_billing.up.sql` | ✅ Written | invoices, invoice_items, receipts with constraints |
| HTTP routes registered | ❌ Blocked | Requires TASK-AUTH-03 (Member A) — role/record middleware |
| Patient FK enforced | ❌ Blocked | Requires TASK-PAT-01 (Member B) — patients table |
| AT12 evidence | ❌ Pending | Requires live PostgreSQL + AUTH-03 integration |

## Financial rules (enforced throughout this module)

- Money: **int64 LKR minor units** everywhere (Rs. 1 500.00 = 150000). No float64.
- Instants: stored as UTC, displayed as Asia/Colombo by the frontend.
- Rounding rule: **pending DEC04** clinic billing review.
- Receipt `operation_id`: server UNIQUE constraint — idempotency gate for AT12.

## Transaction rules

Every mutation (create draft, add item, issue, void, post receipt) writes an
`audit_events` row **in the same PostgreSQL transaction**. If the audit insert
fails, the entire transaction rolls back. No orphan business state can exist
without an audit record.

## Authorization gate (AUTH-03 TODO markers in each method)

All service methods accept `actorID` and are ready for role/record injection.
The handler calls `RegisterRoutes` only after AUTH-03 middleware exists.
Do **not** call `RegisterRoutes` from `router.go` until AUTH-03 is reviewed.
Administrator role alone is not billing access.

## Adding working behavior

When AUTH-03 (Member A) lands:
1. Replace `"TODO-auth03"` actorID stubs in `handler.go` with the middleware
   session value.
2. Add role check: billing staff only (`billing_staff` or equivalent role code).
3. Add record check: caller must be authorized for the referenced patient.
4. Call `billing.RegisterRoutes(api, handler)` from `transport/http/router.go`.
5. Apply migration `000004_billing.up.sql` against the dev database (after
   Member B's patient migration is applied first).
6. Update `docs/api/openapi.yaml` with the billing paths.
7. Run `go test -count=1 ./internal/billing/...` and record results.

## Branch

`member-e/task-bill-01-invoice-schema`
