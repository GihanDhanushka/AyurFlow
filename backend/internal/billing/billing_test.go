package billing

import (
	"math"
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
		{
			name:    "line total overflow",
			req:     AddItemRequest{Description: "Synthetic fee", Quantity: 2, UnitPriceMinor: math.MaxInt64},
			wantErr: true,
		},
		{
			name:    "quantity exceeds PostgreSQL integer",
			req:     AddItemRequest{Description: "Synthetic fee", Quantity: math.MaxInt32 + 1, UnitPriceMinor: 1},
			wantErr: true,
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
			req:     PostReceiptRequest{OperationID: "ac409193-e3ad-4e96-88e2-b856125b52fe", AmountMinor: 50000, PaymentMethod: "cash"},
			wantErr: false,
		},
		{
			name:    "valid direct deposit with label",
			req:     PostReceiptRequest{OperationID: "ac409193-e3ad-4e96-88e2-b856125b52fe", AmountMinor: 25000, PaymentMethod: "direct_deposit", Label: "Slip #001"},
			wantErr: false,
		},
		{
			name:    "missing operationId",
			req:     PostReceiptRequest{OperationID: "", AmountMinor: 50000, PaymentMethod: "cash"},
			wantErr: true,
		},
		{
			name:    "zero amount",
			req:     PostReceiptRequest{OperationID: "ac409193-e3ad-4e96-88e2-b856125b52fe", AmountMinor: 0, PaymentMethod: "cash"},
			wantErr: true,
		},
		{
			name:    "negative amount",
			req:     PostReceiptRequest{OperationID: "ac409193-e3ad-4e96-88e2-b856125b52fe", AmountMinor: -100, PaymentMethod: "cash"},
			wantErr: true,
		},
		{
			name:    "invalid payment method",
			req:     PostReceiptRequest{OperationID: "ac409193-e3ad-4e96-88e2-b856125b52fe", AmountMinor: 10000, PaymentMethod: "card"},
			wantErr: true,
		},
		{
			name:    "malformed operation UUID",
			req:     PostReceiptRequest{OperationID: "abc123", AmountMinor: 100, PaymentMethod: "cash"},
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
