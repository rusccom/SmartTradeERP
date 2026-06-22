BEGIN;

CREATE TABLE IF NOT EXISTS documents.document_sequences (
    tenant_id UUID NOT NULL REFERENCES platform.tenants(id) ON DELETE CASCADE,
    document_type VARCHAR NOT NULL CHECK (
        document_type IN ('RECEIPT', 'SALE', 'WRITEOFF', 'INVENTORY', 'TRANSFER', 'RETURN')
    ),
    year INTEGER NOT NULL CHECK (year BETWEEN 2000 AND 9999),
    next_number INTEGER NOT NULL DEFAULT 1 CHECK (next_number > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, document_type, year)
);

COMMIT;
