BEGIN;

CREATE TABLE IF NOT EXISTS catalog.customers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR NOT NULL,
    phone VARCHAR,
    email VARCHAR,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_customer_default_per_tenant
    ON catalog.customers (tenant_id)
    WHERE is_default = true;

CREATE INDEX IF NOT EXISTS idx_customers_tenant
    ON catalog.customers (tenant_id);

ALTER TABLE documents.documents
    ADD COLUMN IF NOT EXISTS customer_id UUID REFERENCES catalog.customers(id);

COMMIT;
