BEGIN;

ALTER TABLE documents.documents
    DROP CONSTRAINT IF EXISTS documents_documents_type_check;

ALTER TABLE documents.documents
    ADD CONSTRAINT documents_documents_type_check
        CHECK (type IN ('RECEIPT', 'SALE', 'WRITEOFF', 'INVENTORY', 'TRANSFER', 'RETURN'));

CREATE TABLE IF NOT EXISTS documents.shifts (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES platform.tenant_users(id),
    warehouse_id UUID NOT NULL REFERENCES catalog.warehouses(id),
    opened_at TIMESTAMP NOT NULL DEFAULT now(),
    closed_at TIMESTAMP,
    opening_cash DECIMAL(14,4) NOT NULL DEFAULT 0,
    closing_cash DECIMAL(14,4),
    status VARCHAR NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed')),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_shifts_open_per_user
    ON documents.shifts (tenant_id, user_id)
    WHERE status = 'open';

CREATE INDEX IF NOT EXISTS idx_shifts_tenant_user_status
    ON documents.shifts (tenant_id, user_id, status);

ALTER TABLE documents.documents
    ADD COLUMN IF NOT EXISTS shift_id UUID REFERENCES documents.shifts(id);

CREATE TABLE IF NOT EXISTS documents.shift_cash_ops (
    id UUID PRIMARY KEY,
    shift_id UUID NOT NULL REFERENCES documents.shifts(id) ON DELETE CASCADE,
    type VARCHAR NOT NULL CHECK (type IN ('cash_in', 'cash_out')),
    amount DECIMAL(14,4) NOT NULL CHECK (amount > 0),
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_shift_cash_ops_shift
    ON documents.shift_cash_ops (shift_id);

CREATE TABLE IF NOT EXISTS documents.document_payments (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents.documents(id) ON DELETE CASCADE,
    method VARCHAR NOT NULL CHECK (method IN ('cash', 'card', 'transfer')),
    amount DECIMAL(14,4) NOT NULL CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_document_payments_document
    ON documents.document_payments (document_id);

COMMIT;
