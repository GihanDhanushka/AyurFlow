//go:build integration

package billing

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync"
	"testing"
	"time"

	"ayurflow/backend/internal/testutil"

	"github.com/jackc/pgx/v5/pgxpool"
)

const billingTestActor = "d02c1c97-7588-4455-a045-127158e311f6"

func billingFixture(t *testing.T) (*Service, *pgxpool.Pool) {
	t.Helper()
	pool := testutil.Postgres(t, "BILLING_TEST_DATABASE_URL")
	_, err := pool.Exec(t.Context(), `INSERT INTO staff_accounts(id, login_identifier, password_hash)
		VALUES ($1, 'synthetic-billing-actor', '$argon2id$synthetic-non-login-fixture')
		ON CONFLICT (id) DO NOTHING`, billingTestActor)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(t.Context(), `INSERT INTO staff_roles(staff_id, role_code)
		VALUES ($1, 'billing') ON CONFLICT DO NOTHING`, billingTestActor)
	if err != nil {
		t.Fatal(err)
	}
	return NewService(pool), pool
}

func issuedTestInvoice(t *testing.T, service *Service, total int64) Invoice {
	t.Helper()
	// PAT-01 has not supplied its patient FK/fixture yet. These synthetic IDs
	// exercise the current billing schema without inventing a patient schema.
	invoice, err := service.CreateInvoice(t.Context(), billingTestActor, CreateInvoiceRequest{
		PatientID: "30ad6d29-79b0-4028-93c4-614df48ec3ea",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.AddItem(t.Context(), billingTestActor, invoice.ID, AddItemRequest{
		Description: "Synthetic test fee", Quantity: 1, UnitPriceMinor: total,
	})
	if err != nil {
		t.Fatal(err)
	}
	issued, err := service.IssueInvoice(t.Context(), billingTestActor, invoice.ID)
	if err != nil {
		t.Fatal(err)
	}
	return issued
}

func receiptOperation(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	var id string
	if err := pool.QueryRow(t.Context(), `SELECT gen_random_uuid()::text`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func TestReceiptRetryPayloadIntegration(t *testing.T) {
	service, pool := billingFixture(t)
	invoice := issuedTestInvoice(t, service, 100)
	other := issuedTestInvoice(t, service, 100)
	request := PostReceiptRequest{
		OperationID: receiptOperation(t, pool), AmountMinor: 10, PaymentMethod: "cash", Label: "synthetic-reference",
	}
	original, err := service.PostReceipt(t.Context(), billingTestActor, invoice.ID, request)
	if err != nil {
		t.Fatal(err)
	}
	retry := request
	retry.OperationID = strings.ToUpper(retry.OperationID)
	repeated, err := service.PostReceipt(t.Context(), billingTestActor, invoice.ID, retry)
	if err != nil || repeated.ID != original.ID {
		t.Fatalf("identical retry did not return original receipt: %v", err)
	}
	for _, tc := range []struct {
		name    string
		invoice string
		change  func(*PostReceiptRequest)
	}{
		{name: "changed amount", invoice: invoice.ID, change: func(r *PostReceiptRequest) { r.AmountMinor = 20 }},
		{name: "changed method", invoice: invoice.ID, change: func(r *PostReceiptRequest) { r.PaymentMethod = "direct_deposit" }},
		{name: "changed label", invoice: invoice.ID, change: func(r *PostReceiptRequest) { r.Label = "changed-reference" }},
		{name: "different invoice", invoice: other.ID, change: func(*PostReceiptRequest) {}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changed := request
			tc.change(&changed)
			_, err := service.PostReceipt(t.Context(), billingTestActor, tc.invoice, changed)
			if !errors.Is(err, ErrConflict) {
				t.Fatalf("changed retry should conflict: %v", err)
			}
		})
	}
	var receipts, audits int
	if err := pool.QueryRow(t.Context(), `SELECT count(*) FROM receipts`).Scan(&receipts); err != nil {
		t.Fatal(err)
	}
	err = pool.QueryRow(t.Context(), `SELECT count(*) FROM audit_events WHERE action='receipt.posted'`).Scan(&audits)
	if err != nil || receipts != 1 || audits != 1 {
		t.Fatalf("retry duplicated receipt/audit: receipts=%d audits=%d error=%v", receipts, audits, err)
	}
}

func TestConcurrentReceiptOperationsIntegration(t *testing.T) {
	for _, tc := range []struct {
		name          string
		sameOperation bool
		secondAmount  int64
		wantConflicts int
	}{
		{name: "identical simultaneous retry", sameOperation: true, secondAmount: 60},
		{name: "different simultaneous payload", sameOperation: true, secondAmount: 70, wantConflicts: 1},
		{name: "competing payments cannot overpay", secondAmount: 60, wantConflicts: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service, pool := billingFixture(t)
			invoice := issuedTestInvoice(t, service, 100)
			first := PostReceiptRequest{OperationID: receiptOperation(t, pool), AmountMinor: 60, PaymentMethod: "cash"}
			second := first
			second.AmountMinor = tc.secondAmount
			if !tc.sameOperation {
				second.OperationID = receiptOperation(t, pool)
			}
			ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
			var workers sync.WaitGroup
			defer func() { cancel(); workers.Wait() }()
			type outcome struct {
				receipt Receipt
				err     error
			}
			results := make(chan outcome, 2)
			start := make(chan struct{})
			for _, request := range []PostReceiptRequest{first, second} {
				workers.Add(1)
				go func() {
					defer workers.Done()
					<-start
					receipt, err := service.PostReceipt(ctx, billingTestActor, invoice.ID, request)
					results <- outcome{receipt: receipt, err: err}
				}()
			}
			close(start)
			conflicts := 0
			ids := make(map[string]bool)
			for range 2 {
				select {
				case result := <-results:
					if errors.Is(result.err, ErrConflict) {
						conflicts++
						continue
					}
					if result.err != nil {
						t.Fatal(result.err)
					}
					ids[result.receipt.ID] = true
				case <-ctx.Done():
					t.Fatal("concurrent receipt operations did not complete")
				}
			}
			if conflicts != tc.wantConflicts || len(ids) != 1 {
				t.Fatalf("unexpected retry outcomes: conflicts=%d receipts=%d", conflicts, len(ids))
			}
			balance, err := service.GetBalance(t.Context(), billingTestActor, invoice.ID)
			if err != nil || balance.PaidMinor > balance.TotalMinor || balance.PaidMinor == 0 {
				t.Fatalf("invalid committed balance: %+v, %v", balance, err)
			}
		})
	}
}

func TestReceiptOverpaymentCannotOverflowIntegration(t *testing.T) {
	service, pool := billingFixture(t)
	invoice := issuedTestInvoice(t, service, math.MaxInt64)
	_, err := service.PostReceipt(t.Context(), billingTestActor, invoice.ID, PostReceiptRequest{
		OperationID: receiptOperation(t, pool), AmountMinor: 10, PaymentMethod: "cash",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.PostReceipt(t.Context(), billingTestActor, invoice.ID, PostReceiptRequest{
		OperationID: receiptOperation(t, pool), AmountMinor: math.MaxInt64, PaymentMethod: "cash",
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("overflowing overpayment should be rejected: %v", err)
	}
	balance, err := service.GetBalance(t.Context(), billingTestActor, invoice.ID)
	if err != nil || balance.PaidMinor != 10 || balance.OutstandingMinor != math.MaxInt64-10 {
		t.Fatalf("failed overpayment changed committed balance: %+v, %v", balance, err)
	}
}

func TestReceiptAuditFailureRollsBackIntegration(t *testing.T) {
	service, pool := billingFixture(t)
	invoice := issuedTestInvoice(t, service, 100)
	_, err := pool.Exec(t.Context(), `CREATE FUNCTION force_billing_audit_failure() RETURNS trigger
		LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'synthetic audit failure'; END; $$;
		CREATE TRIGGER force_billing_audit_failure BEFORE INSERT ON audit_events
		FOR EACH ROW EXECUTE FUNCTION force_billing_audit_failure()`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.PostReceipt(t.Context(), billingTestActor, invoice.ID, PostReceiptRequest{
		OperationID: receiptOperation(t, pool), AmountMinor: 10, PaymentMethod: "cash",
	})
	if err == nil {
		t.Fatal("forced audit failure was ignored")
	}
	var receipts, receiptAudits int
	if err := pool.QueryRow(t.Context(), `SELECT count(*) FROM receipts`).Scan(&receipts); err != nil {
		t.Fatal(err)
	}
	err = pool.QueryRow(t.Context(), `SELECT count(*) FROM audit_events WHERE action='receipt.posted'`).Scan(&receiptAudits)
	if err != nil || receipts != 0 || receiptAudits != 0 {
		t.Fatal("receipt or audit survived the failed transaction")
	}
}

func TestInvoiceDraftAndIssueLifecycleIntegration(t *testing.T) {
	service, pool := billingFixture(t)
	patientID := "30ad6d29-79b0-4028-93c4-614df48ec3ea"

	// 1. Create draft invoice
	inv, err := service.CreateInvoice(t.Context(), billingTestActor, CreateInvoiceRequest{
		PatientID: patientID,
	})
	if err != nil {
		t.Fatalf("failed to create draft invoice: %v", err)
	}
	if inv.Status != StatusDraft {
		t.Fatalf("expected status draft, got %s", inv.Status)
	}
	if inv.TotalMinor != nil {
		t.Fatalf("expected nil total on draft, got %v", inv.TotalMinor)
	}

	// 2. Add line items with integer LKR minor units
	item1, err := service.AddItem(t.Context(), billingTestActor, inv.ID, AddItemRequest{
		Description: "Ayurvedic Consultation", Quantity: 1, UnitPriceMinor: 250000,
	})
	if err != nil || item1.LineTotalMinor != 250000 {
		t.Fatalf("failed to add item1: item=%+v err=%v", item1, err)
	}

	item2, err := service.AddItem(t.Context(), billingTestActor, inv.ID, AddItemRequest{
		Description: "Shirodhara Oil 500g", Quantity: 2, UnitPriceMinor: 150000,
	})
	if err != nil || item2.LineTotalMinor != 300000 {
		t.Fatalf("failed to add item2: item=%+v err=%v", item2, err)
	}

	tempItem, err := service.AddItem(t.Context(), billingTestActor, inv.ID, AddItemRequest{
		Description: "Mistaken Item", Quantity: 1, UnitPriceMinor: 50000,
	})
	if err != nil {
		t.Fatalf("failed to add tempItem: %v", err)
	}

	// 3. Remove mistakenly added item from draft
	err = service.RemoveItem(t.Context(), billingTestActor, inv.ID, tempItem.ID)
	if err != nil {
		t.Fatalf("failed to remove item: %v", err)
	}

	// Check invoice lines
	fetched, err := service.GetInvoice(t.Context(), billingTestActor, inv.ID)
	if err != nil || len(fetched.Items) != 2 {
		t.Fatalf("expected 2 items after removal, got %d items: %v", len(fetched.Items), err)
	}

	// 4. Atomically issue invoice
	issued, err := service.IssueInvoice(t.Context(), billingTestActor, inv.ID)
	if err != nil {
		t.Fatalf("failed to issue invoice: %v", err)
	}
	if issued.Status != StatusIssued {
		t.Fatalf("expected status issued, got %s", issued.Status)
	}
	if issued.TotalMinor == nil || *issued.TotalMinor != 550000 {
		t.Fatalf("expected total 550000 minor units (Rs. 5,500.00), got %v", issued.TotalMinor)
	}
	if issued.IssuedBy == nil || *issued.IssuedBy != billingTestActor {
		t.Fatalf("expected issued_by %s, got %v", billingTestActor, issued.IssuedBy)
	}
	if issued.IssuedAt == nil {
		t.Fatal("expected non-nil issued_at")
	}

	// Check audit event
	var auditCount int
	err = pool.QueryRow(t.Context(), `SELECT count(*) FROM audit_events WHERE action='invoice.issued' AND object_id=$1`, inv.ID).Scan(&auditCount)
	if err != nil || auditCount != 1 {
		t.Fatalf("expected 1 audit event for invoice.issued, got %d (err: %v)", auditCount, err)
	}

	// 5. Verify immutability: Cannot add or remove items on an issued invoice
	_, err = service.AddItem(t.Context(), billingTestActor, inv.ID, AddItemRequest{
		Description: "Late Addition", Quantity: 1, UnitPriceMinor: 10000,
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("adding item to issued invoice must return ErrConflict, got %v", err)
	}

	err = service.RemoveItem(t.Context(), billingTestActor, inv.ID, item1.ID)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("removing item from issued invoice must return ErrConflict, got %v", err)
	}

	// 6. Cannot void an issued invoice
	err = service.VoidInvoice(t.Context(), billingTestActor, inv.ID)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("voiding issued invoice must return ErrConflict, got %v", err)
	}
}

func TestInvoiceRoleAccessDeniedIntegration(t *testing.T) {
	service, pool := billingFixture(t)
	patientID := "30ad6d29-79b0-4028-93c4-614df48ec3ea"

	// Staff account with therapist role only (no billing role)
	therapistActor := "e11e11e1-1111-4444-a045-127158e311f6"
	_, err := pool.Exec(t.Context(), `INSERT INTO staff_accounts(id, login_identifier, password_hash)
		VALUES ($1, 'synthetic-therapist', '$argon2id$synthetic-non-login-fixture')
		ON CONFLICT (id) DO NOTHING`, therapistActor)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(t.Context(), `INSERT INTO staff_roles(staff_id, role_code)
		VALUES ($1, 'therapist') ON CONFLICT DO NOTHING`, therapistActor)
	if err != nil {
		t.Fatal(err)
	}

	// Attempting to create invoice as therapist must fail
	_, err = service.CreateInvoice(t.Context(), therapistActor, CreateInvoiceRequest{PatientID: patientID})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden for therapist actor, got %v", err)
	}

	// Disabled staff account with billing role
	disabledActor := "d22d22d2-2222-4444-a045-127158e311f6"
	_, err = pool.Exec(t.Context(), `INSERT INTO staff_accounts(id, login_identifier, password_hash, active)
		VALUES ($1, 'synthetic-disabled', '$argon2id$synthetic-non-login-fixture', false)
		ON CONFLICT (id) DO NOTHING`, disabledActor)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(t.Context(), `INSERT INTO staff_roles(staff_id, role_code)
		VALUES ($1, 'billing') ON CONFLICT DO NOTHING`, disabledActor)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.CreateInvoice(t.Context(), disabledActor, CreateInvoiceRequest{PatientID: patientID})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden for disabled actor, got %v", err)
	}
}

func TestVoidDraftInvoiceIntegration(t *testing.T) {
	service, pool := billingFixture(t)
	patientID := "30ad6d29-79b0-4028-93c4-614df48ec3ea"

	inv, err := service.CreateInvoice(t.Context(), billingTestActor, CreateInvoiceRequest{PatientID: patientID})
	if err != nil {
		t.Fatal(err)
	}

	err = service.VoidInvoice(t.Context(), billingTestActor, inv.ID)
	if err != nil {
		t.Fatalf("failed to void draft invoice: %v", err)
	}

	fetched, err := service.GetInvoice(t.Context(), billingTestActor, inv.ID)
	if err != nil || fetched.Status != StatusVoid {
		t.Fatalf("expected status void, got %s (err: %v)", fetched.Status, err)
	}

	var auditCount int
	err = pool.QueryRow(t.Context(), `SELECT count(*) FROM audit_events WHERE action='invoice.voided' AND object_id=$1`, inv.ID).Scan(&auditCount)
	if err != nil || auditCount != 1 {
		t.Fatalf("expected 1 audit event for invoice.voided, got %d", auditCount)
	}

	// Cannot add item or issue void invoice
	_, err = service.AddItem(t.Context(), billingTestActor, inv.ID, AddItemRequest{
		Description: "Item", Quantity: 1, UnitPriceMinor: 1000,
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("adding item to void invoice must return ErrConflict, got %v", err)
	}

	_, err = service.IssueInvoice(t.Context(), billingTestActor, inv.ID)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("issuing void invoice must return ErrConflict, got %v", err)
	}
}
