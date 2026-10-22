"use client";

/**
 * Billing module placeholder page (TASK-BILL-05).
 *
 * This page will be replaced with real invoice/receipt UI screens once:
 *  - TASK-AUTH-03 (role/record middleware) is integrated by Member A
 *  - TASK-BILL-01 (invoice/receipt schema) migration is applied
 *  - API endpoints are registered in the router
 *
 * The TypeScript types and API client shapes in @/features/billing/types.ts
 * are already defined and match the Go backend domain types.
 *
 * Module: MOD-04 | Owner: Member E | Reviewer: Member D
 * SRS: BIL-01–06 | AT: AT12
 */

import Link from "next/link";

export default function BillingPage() {
  return (
    <section className="max-w-3xl">
      <Link href="/" className="text-sm font-semibold text-green-800">
        ← All work areas
      </Link>

      <p className="mt-10 text-sm font-semibold uppercase tracking-widest text-slate-500">
        MOD-04 · In progress
      </p>
      <h1 className="mt-3 text-4xl font-semibold">Invoices &amp; Receipts</h1>
      <p className="mt-5 text-lg leading-8 text-slate-600">
        Itemized invoices, deposits and recorded receipts for patient billing.
        All amounts are in Sri Lankan Rupees (LKR).
      </p>

      {/* Status notice — remove once routes are live */}
      <div className="mt-8 rounded-2xl border border-amber-200 bg-amber-50 p-7">
        <h2 className="text-xl font-semibold text-amber-900">
          Billing screens are being prepared
        </h2>
        <p className="mt-3 leading-7 text-amber-800">
          The database schema and backend service logic are implemented. Billing
          pages will become available after authentication and authorization
          controls (TASK-AUTH-03) are integrated.
        </p>
      </div>

      {/* Planned screens — shown as a roadmap */}
      <div className="mt-8 space-y-4">
        <h2 className="text-xl font-semibold">Planned screens</h2>
        <ul className="ml-5 list-disc space-y-2 leading-7 text-slate-600">
          <li>
            <strong>Invoice list</strong> — browse draft and issued invoices per
            patient
          </li>
          <li>
            <strong>Create invoice</strong> — open a draft and add line items
            (description, quantity, unit price)
          </li>
          <li>
            <strong>Issue invoice</strong> — lock the invoice and generate the
            final total
          </li>
          <li>
            <strong>Post receipt</strong> — record a cash or direct-deposit
            payment with a deposit label
          </li>
          <li>
            <strong>Receipt history</strong> — view all receipts and the
            outstanding balance for an invoice
          </li>
          <li>
            <strong>Outstanding balance report</strong> — list all issued
            invoices with remaining balances
          </li>
        </ul>
      </div>
    </section>
  );
}
