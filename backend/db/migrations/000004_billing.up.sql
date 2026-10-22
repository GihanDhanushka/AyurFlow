-- Migration 000004: billing module schema
-- Owner: Member E  Reviewer: Member D
-- Depends on: 000001 (staff_accounts, audit_events)
-- NOTE: References patients(id) — apply only after Member B's patient migration exists.
--       In development, the FK is deferred; co-ordinate migration order with Member A.

-- invoices -----------------------------------------------------------------
CREATE TABLE invoices (
    id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    -- patient_id references patients(id) once Member B's migration is applied.
    -- Recorded as text here so the billing schema compiles independently;
    -- the FK will be added in a follow-up migration after PAT-01 lands.
    patient_id       uuid        NOT NULL,
    status           text        NOT NULL DEFAULT 'draft'
                                 CHECK (status IN ('draft', 'issued', 'void')),
    -- total_minor_units is stored after issue to avoid repeated summing.
    -- NULL on a draft; set atomically during issue.
    total_minor_units bigint,
    issued_at        timestamptz,
    issued_by        uuid        REFERENCES staff_accounts(id),
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT invoices_total_requires_issue
        CHECK (
            (status = 'draft'  AND total_minor_units IS NULL AND issued_at IS NULL AND issued_by IS NULL) OR
            (status = 'issued' AND total_minor_units IS NOT NULL AND issued_at IS NOT NULL AND issued_by IS NOT NULL) OR
            (status = 'void'   AND total_minor_units IS NULL AND issued_at IS NULL)
        )
);

CREATE INDEX invoices_patient_idx  ON invoices(patient_id);
CREATE INDEX invoices_status_created_idx ON invoices(status, created_at DESC);

-- invoice_items ------------------------------------------------------------
-- Lines are editable only while the invoice is in 'draft' status.
-- After issue the invoice and all its lines are immutable.
CREATE TABLE invoice_items (
    id              uuid    PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id      uuid    NOT NULL REFERENCES invoices(id),
    description     text    NOT NULL CHECK (length(btrim(description)) > 0),
    quantity        integer NOT NULL CHECK (quantity > 0),
    -- unit_price_minor: LKR in minor units (e.g. Rs. 500.00 = 50000)
    unit_price_minor bigint  NOT NULL CHECK (unit_price_minor >= 0),
    -- line_total_minor is always computed; never stored separately.
    -- Calculated as quantity * unit_price_minor in application layer.
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX invoice_items_invoice_idx ON invoice_items(invoice_id);

-- receipts -----------------------------------------------------------------
CREATE TABLE receipts (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id      uuid        NOT NULL REFERENCES invoices(id),
    -- operation_id is client-generated per intended posting attempt.
    -- UNIQUE enforces idempotency: same operation_id = same receipt (no duplicate).
    operation_id    uuid        NOT NULL UNIQUE,
    amount_minor    bigint      NOT NULL CHECK (amount_minor > 0),
    payment_method  text        NOT NULL CHECK (payment_method IN ('cash', 'direct_deposit')),
    -- label is used for deposit references, bank slip numbers, etc.
    label           text,
    issued_at       timestamptz NOT NULL DEFAULT now(),
    issued_by       uuid        NOT NULL REFERENCES staff_accounts(id)
);

CREATE INDEX receipts_invoice_idx    ON receipts(invoice_id);
CREATE INDEX receipts_operation_idx  ON receipts(operation_id);
