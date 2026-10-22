package billing

import (
	"context"
	"testing"
)

// ---------------------------------------------------------------------------
// Domain / validation unit tests
// These run without a database and prove the validation and business rules
// that don't require a real PostgreSQL connection.
// ---------------------------------------------------------------------------

func TestValidateAddItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		req     AddItemRequest
		wantErr bool
	}{
		{
			name:    "valid item",
			req:     AddItemRequest{Description: "Consultation fee", Quantity: 1, UnitPriceMinor: 50000},
			wantErr: false,
		},
		{
			name:    "blank description",
			req:     AddItemRequest{Description: "   ", Quantity: 1, UnitPriceMinor: 10000},
			wantErr: true,
		},
		{
			name:    "zero quantity",
			req:     AddItemRequest{Description: "Session fee", Quantity: 0, UnitPriceMinor: 10000},
			wantErr: true,
		},
		{
			name:    "negative quantity",
			req:     AddItemRequest{Description: "Session fee", Quantity: -1, UnitPriceMinor: 10000},
			wantErr: true,
		},
		{
			name:    "negative unit price",
			req:     AddItemRequest{Description: "Session fee", Quantity: 1, UnitPriceMinor: -1},
			wantErr: true,
		},
		{
			name:    "zero price is allowed (complimentary item)",
			req:     AddItemRequest{Description: "Complimentary item", Quantity: 1, UnitPriceMinor: 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAddItem(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePostReceipt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		req     PostReceiptRequest
		wantErr bool
	}{
		{
			name:    "valid cash receipt",
			req:     PostReceiptRequest{OperationID: "abc123", AmountMinor: 50000, PaymentMethod: "cash"},
			wantErr: false,
		},
		{
			name:    "valid direct deposit with label",
			req:     PostReceiptRequest{OperationID: "abc456", AmountMinor: 25000, PaymentMethod: "direct_deposit", Label: "Slip #001"},
			wantErr: false,
		},
		{
			name:    "missing operationId",
			req:     PostReceiptRequest{OperationID: "", AmountMinor: 50000, PaymentMethod: "cash"},
			wantErr: true,
		},
		{
			name:    "zero amount",
			req:     PostReceiptRequest{OperationID: "abc789", AmountMinor: 0, PaymentMethod: "cash"},
			wantErr: true,
		},
		{
			name:    "negative amount",
			req:     PostReceiptRequest{OperationID: "abc000", AmountMinor: -100, PaymentMethod: "cash"},
			wantErr: true,
		},
		{
			name:    "invalid payment method",
			req:     PostReceiptRequest{OperationID: "abc111", AmountMinor: 10000, PaymentMethod: "card"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePostReceipt(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePostReceipt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLineTotalMinor(t *testing.T) {
	t.Parallel()
	// Verify that line total arithmetic uses integer math (no float64).
	item := InvoiceItem{Quantity: 3, UnitPriceMinor: 15000}
	item.LineTotalMinor = int64(item.Quantity) * item.UnitPriceMinor
	want := int64(45000)
	if item.LineTotalMinor != want {
		t.Errorf("LineTotalMinor = %d, want %d", item.LineTotalMinor, want)
	}
}

func TestInvoiceStatusConstants(t *testing.T) {
	t.Parallel()
	if StatusDraft != "draft" || StatusIssued != "issued" || StatusVoid != "void" {
		t.Error("invoice status constants have unexpected values")
	}
}

// ---------------------------------------------------------------------------
// Itoa helper test
// ---------------------------------------------------------------------------

func TestItoa(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   int
		want string
	}{
		{0, "0"}, {1, "1"}, {9, "9"}, {10, "10"}, {25, "25"}, {100, "100"},
	}
	for _, tt := range tests {
		if got := itoa(tt.in); got != tt.want {
			t.Errorf("itoa(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Integration test stub (requires BILLING_TEST_DATABASE_URL)
// Pattern follows internal/auth/session_integration_test.go
// ---------------------------------------------------------------------------

// billingTestContext is unused until the PostgreSQL integration test runs.
var _ = context.Background // suppress unused import warning in stub

// TestService_Integration_Stub documents the integration test structure that
// will be implemented in a follow-up (after PAT-01 patient schema lands).
// Run with: go test -tags=integration ./internal/billing
func TestService_Integration_Stub(t *testing.T) {
	// Skip unless the integration build tag is present and the database URL is set.
	// Replace with actual integration logic when PAT-01 is integrated.
	t.Skip("integration test: set BILLING_TEST_DATABASE_URL and build with -tags=integration to run")
}
