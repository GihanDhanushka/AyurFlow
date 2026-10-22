package billing

import (
	"errors"
	"time"
)

// ---------------------------------------------------------------------------
// Sentinel errors
// ---------------------------------------------------------------------------

// ErrNotFound is returned when an invoice or receipt does not exist or is
// not visible to the caller.
var ErrNotFound = errors.New("billing: record not found")

// ErrConflict is returned for state transitions that conflict with the
// current record state (e.g. editing an issued invoice, overpayment).
var ErrConflict = errors.New("billing: state conflict")

// ErrForbidden is returned when the caller lacks the required role or
// record-level access. AUTH-03 will supply the shared middleware; this
// sentinel lets service tests verify the rejection path independently.
var ErrForbidden = errors.New("billing: access denied")

// ErrValidation is returned for field-level input problems (negative amount,
// blank description, etc.).
var ErrValidation = errors.New("billing: validation failed")

// ---------------------------------------------------------------------------
// Invoice status
// ---------------------------------------------------------------------------

// InvoiceStatus is the lifecycle state of an invoice.
// Allowed transitions: draft → issued; draft → void.
// An issued invoice is immutable; void is terminal.
type InvoiceStatus string

const (
	StatusDraft  InvoiceStatus = "draft"
	StatusIssued InvoiceStatus = "issued"
	StatusVoid   InvoiceStatus = "void"
)

// ---------------------------------------------------------------------------
// Core types
// All monetary fields use int64 LKR minor units (e.g. Rs. 1 500.00 = 150000).
// Never use float64 for money. Rounding rule is pending DEC04.
// Timestamps are stored in UTC; display in Asia/Colombo is the frontend's job.
// ---------------------------------------------------------------------------

// Invoice is a billing record linked to one patient.
type Invoice struct {
	ID               string        `json:"id"`
	PatientID        string        `json:"patientId"`
	Status           InvoiceStatus `json:"status"`
	TotalMinor       *int64        `json:"totalMinor,omitempty"` // nil on draft
	IssuedAt         *time.Time    `json:"issuedAt,omitempty"`
	IssuedBy         *string       `json:"issuedBy,omitempty"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	Items            []InvoiceItem `json:"items,omitempty"`
}

// InvoiceItem is a single line on an invoice.
// Lines are editable only while the invoice is in draft status.
// line_total_minor = quantity * unit_price_minor (computed; never stored).
type InvoiceItem struct {
	ID              string `json:"id"`
	InvoiceID       string `json:"invoiceId"`
	Description     string `json:"description"`
	Quantity        int    `json:"quantity"`
	UnitPriceMinor  int64  `json:"unitPriceMinor"`
	LineTotalMinor  int64  `json:"lineTotalMinor"` // always Quantity * UnitPriceMinor
	CreatedAt       time.Time `json:"createdAt"`
}

// Receipt records a single payment event against an issued invoice.
// The operation_id provides idempotency for retry-safe posting (AT12).
type Receipt struct {
	ID            string    `json:"id"`
	InvoiceID     string    `json:"invoiceId"`
	OperationID   string    `json:"operationId"`
	AmountMinor   int64     `json:"amountMinor"`
	PaymentMethod string    `json:"paymentMethod"`
	Label         string    `json:"label,omitempty"`
	IssuedAt      time.Time `json:"issuedAt"`
	IssuedBy      string    `json:"issuedBy"`
}

// InvoiceBalance is a derived view of an invoice's payment state.
type InvoiceBalance struct {
	InvoiceID     string `json:"invoiceId"`
	TotalMinor    int64  `json:"totalMinor"`
	PaidMinor     int64  `json:"paidMinor"`
	OutstandingMinor int64 `json:"outstandingMinor"`
}

// ---------------------------------------------------------------------------
// Request shapes (used by service and handler layers)
// ---------------------------------------------------------------------------

// CreateInvoiceRequest contains the minimum fields needed to open a draft.
// PatientID is validated for existence by the service (pending PAT-01 schema).
type CreateInvoiceRequest struct {
	PatientID string `json:"patientId"`
}

// AddItemRequest adds or replaces a line item on a draft invoice.
type AddItemRequest struct {
	Description    string `json:"description"`
	Quantity       int    `json:"quantity"`
	UnitPriceMinor int64  `json:"unitPriceMinor"`
}

// PostReceiptRequest posts a payment against an issued invoice.
// OperationID is client-generated per intended posting; same ID = idempotent return.
type PostReceiptRequest struct {
	OperationID   string `json:"operationId"`
	AmountMinor   int64  `json:"amountMinor"`
	PaymentMethod string `json:"paymentMethod"`
	Label         string `json:"label,omitempty"`
}

// ---------------------------------------------------------------------------
// List / report shapes
// ---------------------------------------------------------------------------

// InvoicePage is the paginated response for invoice list queries.
type InvoicePage struct {
	Items      []Invoice `json:"items"`
	NextCursor *string   `json:"nextCursor"`
}

// OutstandingReportRow is one row in the outstanding balance report (TASK-BILL-06).
type OutstandingReportRow struct {
	InvoiceID        string    `json:"invoiceId"`
	PatientID        string    `json:"patientId"`
	IssuedAt         time.Time `json:"issuedAt"`
	TotalMinor       int64     `json:"totalMinor"`
	PaidMinor        int64     `json:"paidMinor"`
	OutstandingMinor int64     `json:"outstandingMinor"`
	AgeDays          int       `json:"ageDays"`
}

// OutstandingReport wraps the full outstanding balance report.
type OutstandingReport struct {
	GeneratedAt time.Time              `json:"generatedAt"`
	Rows        []OutstandingReportRow `json:"rows"`
}
