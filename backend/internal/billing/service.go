package billing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Service enforces authorization, workflow rules and transaction boundaries
// for all billing operations. All protected mutations write an audit event
// in the same transaction; a failed audit write rolls back the mutation.
//
// AUTH-03 (role and record checks) is not yet implemented. Every method that
// requires a caller identity accepts actorID and roles parameters so that the
// permission logic can be added in place once AUTH-03 lands. Methods marked
// with "AUTH-03 TODO" must not be exposed on protected routes until that work
// is integrated.
//
// DEC04 TODO: rounding rule, allowed payment methods and tariff examples must
// be confirmed with the clinic billing staff before AT12 acceptance.
type Service struct {
	pool *pgxpool.Pool
	repo *Repository
}

// NewService creates a billing Service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool, repo: NewRepository(pool)}
}

// ---------------------------------------------------------------------------
// Invoice lifecycle
// ---------------------------------------------------------------------------

// CreateInvoice opens a new draft invoice for the given patient.
//
// AUTH-03 TODO: check caller has billing role; check patient record access.
func (s *Service) CreateInvoice(ctx context.Context, actorID string, req CreateInvoiceRequest) (Invoice, error) {
	if strings.TrimSpace(req.PatientID) == "" {
		return Invoice{}, fmt.Errorf("%w: patientId is required", ErrValidation)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Invoice{}, err
	}
	defer tx.Rollback(ctx)

	inv, err := s.repo.CreateInvoice(ctx, tx, req.PatientID)
	if err != nil {
		return Invoice{}, err
	}

	if err := s.repo.RecordAuditEvent(ctx, tx,
		actorID, "invoice.created", "invoice", inv.ID, ""); err != nil {
		return Invoice{}, err // rolls back the invoice insert
	}

	if err := tx.Commit(ctx); err != nil {
		return Invoice{}, err
	}
	return inv, nil
}

// GetInvoice returns a single invoice with its line items.
//
// AUTH-03 TODO: check caller has billing role; check patient record access.
func (s *Service) GetInvoice(ctx context.Context, _ string, invoiceID string) (Invoice, error) {
	inv, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return Invoice{}, err
	}
	inv.Items, err = s.repo.GetInvoiceItems(ctx, invoiceID)
	return inv, err
}

// ListInvoicesForPatient returns a paginated list of invoices for one patient.
//
// AUTH-03 TODO: check caller has billing role; check patient record access.
func (s *Service) ListInvoicesForPatient(ctx context.Context, _ string, patientID string, limit int, afterID string) (InvoicePage, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	// Fetch limit+1 to detect whether there is a next page.
	items, err := s.repo.ListInvoicesForPatient(ctx, patientID, limit+1, afterID)
	if err != nil {
		return InvoicePage{}, err
	}
	var nextCursor *string
	if len(items) > limit {
		items = items[:limit]
		last := items[len(items)-1].ID
		nextCursor = &last
	}
	return InvoicePage{Items: items, NextCursor: nextCursor}, nil
}

// AddItem adds a line item to a draft invoice.
//
// AUTH-03 TODO: check caller has billing role.
func (s *Service) AddItem(ctx context.Context, actorID, invoiceID string, req AddItemRequest) (InvoiceItem, error) {
	if err := validateAddItem(req); err != nil {
		return InvoiceItem{}, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return InvoiceItem{}, err
	}
	defer tx.Rollback(ctx)

	inv, err := s.repo.GetInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		return InvoiceItem{}, err
	}
	if inv.Status != StatusDraft {
		return InvoiceItem{}, fmt.Errorf("%w: invoice is %s, only draft invoices accept new items", ErrConflict, inv.Status)
	}

	item, err := s.repo.AddInvoiceItem(ctx, tx, invoiceID, req.Description, req.Quantity, req.UnitPriceMinor)
	if err != nil {
		return InvoiceItem{}, err
	}

	if err := s.repo.RecordAuditEvent(ctx, tx,
		actorID, "invoice.item_added", "invoice", invoiceID, ""); err != nil {
		return InvoiceItem{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return InvoiceItem{}, err
	}
	return item, nil
}

// RemoveItem removes a line item from a draft invoice.
//
// AUTH-03 TODO: check caller has billing role.
func (s *Service) RemoveItem(ctx context.Context, actorID, invoiceID, itemID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	inv, err := s.repo.GetInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		return err
	}
	if inv.Status != StatusDraft {
		return fmt.Errorf("%w: cannot remove items from a %s invoice", ErrConflict, inv.Status)
	}

	if err := s.repo.RemoveInvoiceItem(ctx, tx, itemID, invoiceID); err != nil {
		return err
	}

	if err := s.repo.RecordAuditEvent(ctx, tx,
		actorID, "invoice.item_removed", "invoice", invoiceID, ""); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// IssueInvoice atomically locks the invoice, computes the total from its lines,
// and transitions it to 'issued'. An issued invoice is immutable.
//
// AUTH-03 TODO: check caller has billing role.
func (s *Service) IssueInvoice(ctx context.Context, actorID, invoiceID string) (Invoice, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Invoice{}, err
	}
	defer tx.Rollback(ctx)

	// Lock the row to serialize concurrent issue attempts.
	inv, err := s.repo.GetInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		return Invoice{}, err
	}
	if inv.Status != StatusDraft {
		return Invoice{}, fmt.Errorf("%w: invoice is already %s", ErrConflict, inv.Status)
	}

	totalMinor, err := s.repo.SumInvoiceItems(ctx, tx, invoiceID)
	if err != nil {
		return Invoice{}, err
	}
	if totalMinor == 0 {
		return Invoice{}, fmt.Errorf("%w: cannot issue an invoice with no items or a zero total", ErrValidation)
	}

	if err := s.repo.IssueInvoice(ctx, tx, invoiceID, totalMinor, actorID); err != nil {
		return Invoice{}, err
	}

	// Audit event in the same transaction; failure rolls back the issue.
	if err := s.repo.RecordAuditEvent(ctx, tx,
		actorID, "invoice.issued", "invoice", invoiceID, ""); err != nil {
		return Invoice{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Invoice{}, err
	}

	// Re-fetch the fully updated invoice with lines.
	return s.GetInvoice(ctx, actorID, invoiceID)
}

// VoidInvoice transitions a draft invoice to 'void'. Fails if any receipts exist.
//
// AUTH-03 TODO: check caller has billing role.
func (s *Service) VoidInvoice(ctx context.Context, actorID, invoiceID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	inv, err := s.repo.GetInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		return err
	}
	if inv.Status != StatusDraft {
		return fmt.Errorf("%w: can only void a draft invoice, current status is %s", ErrConflict, inv.Status)
	}

	// Guard: no receipts should exist on a draft, but check defensively.
	paid, err := s.repo.SumReceipts(ctx, tx, invoiceID)
	if err != nil {
		return err
	}
	if paid > 0 {
		return fmt.Errorf("%w: invoice has receipts and cannot be voided", ErrConflict)
	}

	if err := s.repo.VoidInvoice(ctx, tx, invoiceID); err != nil {
		return err
	}

	if err := s.repo.RecordAuditEvent(ctx, tx,
		actorID, "invoice.voided", "invoice", invoiceID, ""); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Receipts and balance
// ---------------------------------------------------------------------------

// PostReceipt posts a payment receipt against an issued invoice.
// Idempotency: if operationID already exists, the existing receipt is returned
// without writing a duplicate (AT12 retry safety).
// Concurrency: the invoice row is locked FOR UPDATE before the overpayment
// check, serializing concurrent receipt posts on the same invoice.
//
// AUTH-03 TODO: check caller has billing role.
func (s *Service) PostReceipt(ctx context.Context, actorID, invoiceID string, req PostReceiptRequest) (Receipt, error) {
	if err := validatePostReceipt(req); err != nil {
		return Receipt{}, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Receipt{}, err
	}
	defer tx.Rollback(ctx)

	// Idempotency check: return existing receipt if operationID already committed.
	existing, err := s.repo.FindReceiptByOperationID(ctx, tx, req.OperationID)
	if err == nil {
		// Already committed — idempotent return. Roll back the read-only tx cleanly.
		tx.Rollback(ctx)
		return existing, nil
	}
	if err != ErrNotFound {
		return Receipt{}, err
	}

	// Lock the invoice row to serialize balance checks.
	inv, err := s.repo.GetInvoiceForUpdate(ctx, tx, invoiceID)
	if err != nil {
		return Receipt{}, err
	}
	if inv.Status != StatusIssued {
		return Receipt{}, fmt.Errorf("%w: can only post receipts against an issued invoice (status: %s)", ErrConflict, inv.Status)
	}

	paid, err := s.repo.SumReceipts(ctx, tx, invoiceID)
	if err != nil {
		return Receipt{}, err
	}
	if paid+req.AmountMinor > *inv.TotalMinor {
		return Receipt{}, fmt.Errorf("%w: receipt amount would exceed the invoice total (outstanding: %d minor units)", ErrConflict, *inv.TotalMinor-paid)
	}

	rec, err := s.repo.InsertReceipt(ctx, tx, invoiceID, req.OperationID, req.AmountMinor, req.PaymentMethod, req.Label, actorID)
	if err != nil {
		return Receipt{}, err
	}

	// Audit in the same transaction; failure rolls back the receipt insert.
	if err := s.repo.RecordAuditEvent(ctx, tx,
		actorID, "receipt.posted", "receipt", rec.ID, ""); err != nil {
		return Receipt{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Receipt{}, err
	}
	return rec, nil
}

// GetBalance returns the derived balance for an issued invoice.
//
// AUTH-03 TODO: check caller has billing role.
func (s *Service) GetBalance(ctx context.Context, _ string, invoiceID string) (InvoiceBalance, error) {
	inv, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return InvoiceBalance{}, err
	}
	if inv.Status != StatusIssued {
		return InvoiceBalance{}, fmt.Errorf("%w: balance is only available for issued invoices", ErrConflict)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return InvoiceBalance{}, err
	}
	defer tx.Rollback(ctx)

	paid, err := s.repo.SumReceipts(ctx, tx, invoiceID)
	if err != nil {
		return InvoiceBalance{}, err
	}
	tx.Rollback(ctx)

	return InvoiceBalance{
		InvoiceID:        invoiceID,
		TotalMinor:       *inv.TotalMinor,
		PaidMinor:        paid,
		OutstandingMinor: *inv.TotalMinor - paid,
	}, nil
}

// ListReceipts returns all receipts for an invoice.
//
// AUTH-03 TODO: check caller has billing role.
func (s *Service) ListReceipts(ctx context.Context, _ string, invoiceID string) ([]Receipt, error) {
	// Confirm invoice exists before querying receipts.
	if _, err := s.repo.GetInvoice(ctx, invoiceID); err != nil {
		return nil, err
	}
	return s.repo.ListReceipts(ctx, invoiceID)
}

// ---------------------------------------------------------------------------
// Reports (TASK-BILL-06)
// ---------------------------------------------------------------------------

// OutstandingReport returns all issued invoices with a remaining balance.
//
// AUTH-03 TODO: check caller has billing or admin role.
func (s *Service) OutstandingReport(ctx context.Context, _ string) (OutstandingReport, error) {
	rows, err := s.repo.OutstandingInvoices(ctx)
	if err != nil {
		return OutstandingReport{}, err
	}
	return OutstandingReport{
		GeneratedAt: time.Now().UTC(),
		Rows:        rows,
	}, nil
}

// ---------------------------------------------------------------------------
// Validation helpers
// ---------------------------------------------------------------------------

func validateAddItem(req AddItemRequest) error {
	var errs []string
	if strings.TrimSpace(req.Description) == "" {
		errs = append(errs, "description is required")
	}
	if req.Quantity <= 0 {
		errs = append(errs, "quantity must be positive")
	}
	if req.UnitPriceMinor < 0 {
		errs = append(errs, "unitPriceMinor must be non-negative")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrValidation, strings.Join(errs, "; "))
	}
	return nil
}

func validatePostReceipt(req PostReceiptRequest) error {
	var errs []string
	if strings.TrimSpace(req.OperationID) == "" {
		errs = append(errs, "operationId is required")
	}
	if req.AmountMinor <= 0 {
		errs = append(errs, "amountMinor must be positive")
	}
	if req.PaymentMethod != "cash" && req.PaymentMethod != "direct_deposit" {
		errs = append(errs, "paymentMethod must be 'cash' or 'direct_deposit'")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrValidation, strings.Join(errs, "; "))
	}
	return nil
}
