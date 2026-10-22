/**
 * Billing module TypeScript types.
 * All money values are integer LKR minor units (e.g. Rs. 1 500.00 = 150000).
 * Use formatLKR() from @/lib/format for display — never do arithmetic here.
 * Timestamps arrive as UTC RFC 3339 strings; use formatClinicDateTime() for display.
 */

export type InvoiceStatus = "draft" | "issued" | "void";

export type PaymentMethod = "cash" | "direct_deposit";

export interface InvoiceItem {
  id: string;
  invoiceId: string;
  description: string;
  quantity: number;
  unitPriceMinor: number;  // integer LKR minor units
  lineTotalMinor: number;  // always quantity * unitPriceMinor
  createdAt: string;       // UTC ISO string
}

export interface Invoice {
  id: string;
  patientId: string;
  status: InvoiceStatus;
  totalMinor?: number;     // null/undefined on draft
  issuedAt?: string;       // UTC ISO string; null on draft
  issuedBy?: string;       // staff account ID; null on draft
  createdAt: string;
  updatedAt: string;
  items?: InvoiceItem[];
}

export interface Receipt {
  id: string;
  invoiceId: string;
  operationId: string;
  amountMinor: number;     // integer LKR minor units
  paymentMethod: PaymentMethod;
  label?: string;
  issuedAt: string;        // UTC ISO string
  issuedBy: string;        // staff account ID
}

export interface InvoiceBalance {
  invoiceId: string;
  totalMinor: number;
  paidMinor: number;
  outstandingMinor: number;
}

export interface InvoicePage {
  items: Invoice[];
  nextCursor: string | null;
}

export interface OutstandingReportRow {
  invoiceId: string;
  patientId: string;
  issuedAt: string;
  totalMinor: number;
  paidMinor: number;
  outstandingMinor: number;
  ageDays: number;
}

export interface OutstandingReport {
  generatedAt: string;
  rows: OutstandingReportRow[];
}

// ---------------------------------------------------------------------------
// Request shapes (match Go PostReceiptRequest, AddItemRequest, etc.)
// ---------------------------------------------------------------------------

export interface CreateInvoiceRequest {
  patientId: string;
}

export interface AddItemRequest {
  description: string;
  quantity: number;
  unitPriceMinor: number;
}

export interface PostReceiptRequest {
  operationId: string;     // client-generated UUID per intended posting attempt
  amountMinor: number;
  paymentMethod: PaymentMethod;
  label?: string;
}
