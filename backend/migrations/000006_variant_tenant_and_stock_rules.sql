BEGIN;

ALTER TABLE catalog.product_variants
    ADD COLUMN IF NOT EXISTS tenant_id UUID;

UPDATE catalog.product_variants v
SET tenant_id = p.tenant_id
FROM catalog.products p
WHERE p.id = v.product_id AND v.tenant_id IS NULL;

ALTER TABLE catalog.product_variants
    ALTER COLUMN tenant_id SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_products_tenant_id_id
    ON catalog.products (tenant_id, id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_product_variants_tenant_id_id
    ON catalog.product_variants (tenant_id, id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_product_variants_sku_per_tenant
    ON catalog.product_variants (tenant_id, lower(btrim(sku_code)))
    WHERE sku_code IS NOT NULL AND btrim(sku_code) <> '';

CREATE UNIQUE INDEX IF NOT EXISTS ux_product_variants_barcode_per_tenant
    ON catalog.product_variants (tenant_id, lower(btrim(barcode)))
    WHERE barcode IS NOT NULL AND btrim(barcode) <> '';

CREATE INDEX IF NOT EXISTS idx_cl_tenant_warehouse_variant
    ON ledger.cost_ledger (tenant_id, warehouse_id, variant_id);

ALTER TABLE ledger.cost_ledger
    DROP CONSTRAINT IF EXISTS chk_cost_ledger_amounts;

ALTER TABLE ledger.cost_ledger
    ADD CONSTRAINT chk_cost_ledger_amounts
    CHECK (
        unit_price >= 0
        AND total_amount >= 0
        AND running_avg >= 0
        AND (cogs IS NULL OR cogs >= 0)
    );

CREATE OR REPLACE FUNCTION public.add_constraint_if_missing(target_table regclass, constraint_name text, ddl text)
RETURNS void AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = target_table AND conname = constraint_name
    ) THEN
        EXECUTE ddl;
    END IF;
END;
$$ LANGUAGE plpgsql;

SELECT public.add_constraint_if_missing(
    'catalog.product_variants'::regclass,
    'fk_product_variants_tenant',
    'ALTER TABLE catalog.product_variants ADD CONSTRAINT fk_product_variants_tenant
     FOREIGN KEY (tenant_id) REFERENCES platform.tenants (id)'
);

SELECT public.add_constraint_if_missing(
    'catalog.product_variants'::regclass,
    'fk_product_variants_tenant_product',
    'ALTER TABLE catalog.product_variants ADD CONSTRAINT fk_product_variants_tenant_product
     FOREIGN KEY (tenant_id, product_id) REFERENCES catalog.products (tenant_id, id)'
);

SELECT public.add_constraint_if_missing(
    'ledger.cost_ledger'::regclass,
    'fk_cost_ledger_tenant_variant',
    'ALTER TABLE ledger.cost_ledger ADD CONSTRAINT fk_cost_ledger_tenant_variant
     FOREIGN KEY (tenant_id, variant_id) REFERENCES catalog.product_variants (tenant_id, id)'
);

DROP FUNCTION public.add_constraint_if_missing(regclass, text, text);

CREATE OR REPLACE FUNCTION documents.document_item_tenant_guard()
RETURNS trigger AS $$
DECLARE
    doc_tenant UUID;
    variant_tenant UUID;
BEGIN
    SELECT tenant_id INTO doc_tenant
    FROM documents.documents
    WHERE id = NEW.document_id;

    SELECT tenant_id INTO variant_tenant
    FROM catalog.product_variants
    WHERE id = NEW.variant_id;

    IF doc_tenant IS NULL OR variant_tenant IS DISTINCT FROM doc_tenant THEN
        RAISE EXCEPTION 'document item variant must belong to document tenant';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION documents.document_item_component_tenant_guard()
RETURNS trigger AS $$
DECLARE
    doc_tenant UUID;
    component_tenant UUID;
BEGIN
    SELECT d.tenant_id INTO doc_tenant
    FROM documents.document_items i
    JOIN documents.documents d ON d.id = i.document_id
    WHERE i.id = NEW.document_item_id;

    SELECT tenant_id INTO component_tenant
    FROM catalog.product_variants
    WHERE id = NEW.component_variant_id;

    IF doc_tenant IS NULL OR component_tenant IS DISTINCT FROM doc_tenant THEN
        RAISE EXCEPTION 'item component variant must belong to document tenant';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ledger.cost_ledger_tenant_guard()
RETURNS trigger AS $$
DECLARE
    item_document UUID;
    variant_tenant UUID;
BEGIN
    SELECT i.document_id INTO item_document
    FROM documents.document_items i
    WHERE i.id = NEW.document_item_id;

    SELECT tenant_id INTO variant_tenant
    FROM catalog.product_variants
    WHERE id = NEW.variant_id;

    IF item_document IS DISTINCT FROM NEW.document_id THEN
        RAISE EXCEPTION 'ledger item must belong to ledger document';
    END IF;

    IF variant_tenant IS DISTINCT FROM NEW.tenant_id THEN
        RAISE EXCEPTION 'ledger variant must belong to ledger tenant';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION catalog.variant_component_guard()
RETURNS trigger AS $$
DECLARE
    variant_is_composite BOOLEAN;
    component_is_composite BOOLEAN;
    variant_tenant UUID;
    component_tenant UUID;
BEGIN
    SELECT p.is_composite, v.tenant_id
    INTO variant_is_composite, variant_tenant
    FROM catalog.product_variants v
    JOIN catalog.products p ON p.id = v.product_id
    WHERE v.id = NEW.variant_id;

    SELECT p.is_composite, v.tenant_id
    INTO component_is_composite, component_tenant
    FROM catalog.product_variants v
    JOIN catalog.products p ON p.id = v.product_id
    WHERE v.id = NEW.component_variant_id;

    IF variant_is_composite IS DISTINCT FROM true THEN
        RAISE EXCEPTION 'variant must belong to composite product';
    END IF;

    IF component_is_composite IS DISTINCT FROM false THEN
        RAISE EXCEPTION 'component variant must be non-composite';
    END IF;

    IF variant_tenant IS DISTINCT FROM component_tenant THEN
        RAISE EXCEPTION 'cross-tenant components are forbidden';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION documents.document_posting_guard()
RETURNS trigger AS $$
DECLARE
    item_count INTEGER;
    item_total NUMERIC;
    payment_total NUMERIC;
BEGIN
    IF NEW.status <> 'posted' THEN
        RETURN NEW;
    END IF;

    SELECT COUNT(*), COALESCE(SUM(total_amount), 0)
    INTO item_count, item_total
    FROM documents.document_items
    WHERE document_id = NEW.id;

    IF item_count = 0 THEN
        RAISE EXCEPTION 'posted document must have items';
    END IF;

    SELECT COALESCE(SUM(amount), 0)
    INTO payment_total
    FROM documents.document_payments
    WHERE document_id = NEW.id;

    IF NEW.type IN ('SALE', 'RETURN') AND payment_total <> item_total THEN
        RAISE EXCEPTION 'sale and return payments must match item total';
    END IF;

    IF NEW.type NOT IN ('SALE', 'RETURN') AND payment_total <> 0 THEN
        RAISE EXCEPTION 'payments are allowed only for sales and returns';
    END IF;

    IF NEW.type IN ('SALE', 'RETURN') AND NEW.shift_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1
            FROM documents.shifts
            WHERE tenant_id = NEW.tenant_id AND id = NEW.shift_id AND status = 'open'
        ) THEN
            RAISE EXCEPTION 'sale and return documents require an open shift';
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMIT;
