BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TABLE IF EXISTS ledger.cost_ledger CASCADE;
DROP FUNCTION IF EXISTS ledger.cost_ledger_tenant_guard() CASCADE;

CREATE TABLE IF NOT EXISTS platform.tenant_settings (
    tenant_id UUID PRIMARY KEY REFERENCES platform.tenants (id) ON DELETE CASCADE,
    allow_negative_stock BOOLEAN NOT NULL DEFAULT false,
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ledger.posting_batches (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    document_id UUID NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'superseded', 'reversed')),
    effective_date DATE NOT NULL,
    posted_at TIMESTAMP NOT NULL DEFAULT now(),
    supersedes_batch_id UUID,
    reason VARCHAR NOT NULL DEFAULT 'posting',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, id),
    FOREIGN KEY (tenant_id, document_id)
        REFERENCES documents.documents (tenant_id, id),
    FOREIGN KEY (supersedes_batch_id)
        REFERENCES ledger.posting_batches (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_posting_batches_active_document
    ON ledger.posting_batches (tenant_id, document_id)
    WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_posting_batches_document
    ON ledger.posting_batches (tenant_id, document_id, status);

CREATE TABLE IF NOT EXISTS ledger.inventory_movements (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    posting_batch_id UUID NOT NULL,
    document_id UUID NOT NULL,
    document_item_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    movement_date DATE NOT NULL,
    direction VARCHAR(3) NOT NULL CHECK (direction IN ('IN', 'OUT')),
    reason VARCHAR NOT NULL CHECK (reason IN (
        'PURCHASE', 'SALE', 'WRITEOFF', 'SHORTAGE', 'SURPLUS',
        'RETURN_IN', 'RETURN_OUT', 'TRANSFER_IN', 'TRANSFER_OUT'
    )),
    qty DECIMAL(12,3) NOT NULL CHECK (qty > 0),
    unit_price DECIMAL(12,4) NOT NULL CHECK (unit_price >= 0),
    total_amount DECIMAL(14,4) NOT NULL CHECK (total_amount >= 0),
    revenue_amount DECIMAL(14,4),
    posting_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (tenant_id, posting_batch_id)
        REFERENCES ledger.posting_batches (tenant_id, id),
    FOREIGN KEY (tenant_id, document_id)
        REFERENCES documents.documents (tenant_id, id),
    FOREIGN KEY (tenant_id, variant_id)
        REFERENCES catalog.product_variants (tenant_id, id),
    FOREIGN KEY (tenant_id, warehouse_id)
        REFERENCES catalog.warehouses (tenant_id, id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_batch
    ON ledger.inventory_movements (tenant_id, posting_batch_id);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_variant_date
    ON ledger.inventory_movements (tenant_id, variant_id, movement_date, posting_order, created_at, id);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_document_item
    ON ledger.inventory_movements (tenant_id, document_item_id);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_warehouse
    ON ledger.inventory_movements (tenant_id, warehouse_id, variant_id);

CREATE TABLE IF NOT EXISTS ledger.cost_movement_results (
    movement_id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    sequence_num BIGINT NOT NULL,
    qty_delta DECIMAL(12,3) NOT NULL,
    unit_cost DECIMAL(12,4) NOT NULL CHECK (unit_cost >= 0),
    movement_cost DECIMAL(14,4) NOT NULL CHECK (movement_cost >= 0),
    revenue_amount DECIMAL(14,4),
    cogs_amount DECIMAL(14,4),
    gross_profit DECIMAL(14,4),
    running_qty DECIMAL(12,3) NOT NULL,
    running_avg_cost DECIMAL(12,4) NOT NULL CHECK (running_avg_cost >= 0),
    running_inventory_value DECIMAL(14,4) NOT NULL,
    calculation_version UUID NOT NULL,
    calculated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, variant_id, sequence_num),
    FOREIGN KEY (movement_id)
        REFERENCES ledger.inventory_movements (id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id, variant_id)
        REFERENCES catalog.product_variants (tenant_id, id)
);

CREATE INDEX IF NOT EXISTS idx_cost_results_variant_sequence
    ON ledger.cost_movement_results (tenant_id, variant_id, sequence_num DESC);

CREATE TABLE IF NOT EXISTS ledger.stock_balances (
    tenant_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    qty DECIMAL(12,3) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, variant_id, warehouse_id),
    FOREIGN KEY (tenant_id, variant_id)
        REFERENCES catalog.product_variants (tenant_id, id),
    FOREIGN KEY (tenant_id, warehouse_id)
        REFERENCES catalog.warehouses (tenant_id, id)
);

CREATE TABLE IF NOT EXISTS ledger.document_item_financials (
    tenant_id UUID NOT NULL,
    document_id UUID NOT NULL,
    document_item_id UUID NOT NULL,
    revenue_amount DECIMAL(14,4) NOT NULL DEFAULT 0,
    cogs_amount DECIMAL(14,4) NOT NULL DEFAULT 0,
    gross_profit DECIMAL(14,4) NOT NULL DEFAULT 0,
    calculation_version UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, document_item_id),
    FOREIGN KEY (tenant_id, document_id)
        REFERENCES documents.documents (tenant_id, id),
    FOREIGN KEY (document_id, document_item_id)
        REFERENCES documents.document_items (document_id, id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_document_item_financials_document
    ON ledger.document_item_financials (tenant_id, document_id);

CREATE TABLE IF NOT EXISTS ledger.daily_variant_metrics (
    tenant_id UUID NOT NULL,
    date DATE NOT NULL,
    variant_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    qty_in DECIMAL(12,3) NOT NULL DEFAULT 0,
    qty_out DECIMAL(12,3) NOT NULL DEFAULT 0,
    revenue_amount DECIMAL(14,4) NOT NULL DEFAULT 0,
    cogs_amount DECIMAL(14,4) NOT NULL DEFAULT 0,
    gross_profit DECIMAL(14,4) NOT NULL DEFAULT 0,
    ending_qty DECIMAL(12,3) NOT NULL DEFAULT 0,
    ending_avg_cost DECIMAL(12,4) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, date, variant_id, warehouse_id),
    FOREIGN KEY (tenant_id, variant_id)
        REFERENCES catalog.product_variants (tenant_id, id),
    FOREIGN KEY (tenant_id, warehouse_id)
        REFERENCES catalog.warehouses (tenant_id, id)
);

COMMIT;
