package billing

// Package billing implements the invoices, receipts, balances and billing
// reports module (Member E, SRS BIL-01–06, AT12).
//
// Status: TASK-BILL-01 in progress.
// Domain types, repository queries, service logic and unit tests are
// implemented. The HTTP handler and route registration are present but
// routes are NOT registered in the router until AUTH-03 (role/record checks)
// is integrated (TASK-AUTH-03, Member A).
//
// Financial rules (always enforced):
//   - All monetary values use int64 LKR minor units; no float64 anywhere.
//   - Instants are stored in UTC; the frontend displays Asia/Colombo.
//   - Rounding rule is pending DEC04 (clinic billing review).
//
// Transaction rules:
//   - Every mutation (issue, void, post receipt, add item) writes its audit
//     event in the same PostgreSQL transaction. A failed audit write causes
//     the mutation to roll back automatically.
//   - Receipt posting uses operation_id for idempotency (AT12 retry safety).
//   - Concurrent receipt posts serialize on a FOR UPDATE lock on the invoice row.
//
// Authorization (to be enforced when AUTH-03 lands):
//   - Every service method accepts actorID and will check role + record access.
//   - Administrator role alone does not confer billing access.
//   - Do not register /api/v1/billing/... routes before AUTH-03 is reviewed.
package billing
