BEGIN;

CREATE UNIQUE INDEX IF NOT EXISTS ux_documents_number_per_tenant_type
    ON documents.documents (tenant_id, type, number)
    WHERE number IS NOT NULL AND btrim(number) <> '';

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
    'platform.tenants'::regclass,
    'chk_tenants_name_not_blank',
    'ALTER TABLE platform.tenants ADD CONSTRAINT chk_tenants_name_not_blank
     CHECK (btrim(name) <> '''')'
);

SELECT public.add_constraint_if_missing(
    'catalog.products'::regclass,
    'chk_products_name_not_blank',
    'ALTER TABLE catalog.products ADD CONSTRAINT chk_products_name_not_blank
     CHECK (btrim(name) <> '''')'
);

SELECT public.add_constraint_if_missing(
    'catalog.product_variants'::regclass,
    'chk_product_variants_price_non_negative',
    'ALTER TABLE catalog.product_variants ADD CONSTRAINT chk_product_variants_price_non_negative
     CHECK (price IS NULL OR price >= 0)'
);

SELECT public.add_constraint_if_missing(
    'catalog.warehouses'::regclass,
    'chk_warehouses_name_not_blank',
    'ALTER TABLE catalog.warehouses ADD CONSTRAINT chk_warehouses_name_not_blank
     CHECK (btrim(name) <> '''')'
);

SELECT public.add_constraint_if_missing(
    'catalog.customers'::regclass,
    'chk_customers_name_not_blank',
    'ALTER TABLE catalog.customers ADD CONSTRAINT chk_customers_name_not_blank
     CHECK (btrim(name) <> '''')'
);

SELECT public.add_constraint_if_missing(
    'documents.documents'::regclass,
    'chk_documents_warehouse_shape',
    'ALTER TABLE documents.documents ADD CONSTRAINT chk_documents_warehouse_shape
     CHECK (
        (
            type = ''TRANSFER''
            AND warehouse_id IS NULL
            AND source_warehouse_id IS NOT NULL
            AND target_warehouse_id IS NOT NULL
            AND source_warehouse_id <> target_warehouse_id
        )
        OR
        (
            type <> ''TRANSFER''
            AND warehouse_id IS NOT NULL
            AND source_warehouse_id IS NULL
            AND target_warehouse_id IS NULL
        )
     )'
);

SELECT public.add_constraint_if_missing(
    'documents.document_items'::regclass,
    'chk_document_items_amounts',
    'ALTER TABLE documents.document_items ADD CONSTRAINT chk_document_items_amounts
     CHECK (unit_price >= 0 AND total_amount = round(qty * unit_price, 4))'
);

SELECT public.add_constraint_if_missing(
    'documents.shifts'::regclass,
    'chk_shifts_cash_non_negative',
    'ALTER TABLE documents.shifts ADD CONSTRAINT chk_shifts_cash_non_negative
     CHECK (opening_cash >= 0 AND (closing_cash IS NULL OR closing_cash >= 0))'
);

SELECT public.add_constraint_if_missing(
    'ledger.cost_ledger'::regclass,
    'chk_cost_ledger_amounts',
    'ALTER TABLE ledger.cost_ledger ADD CONSTRAINT chk_cost_ledger_amounts
     CHECK (
        unit_price >= 0
        AND total_amount >= 0
        AND running_avg >= 0
        AND (cogs IS NULL OR cogs >= 0)
     )'
);

DROP FUNCTION public.add_constraint_if_missing(regclass, text, text);

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

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_document_posting_guard ON documents.documents;

CREATE TRIGGER trg_document_posting_guard
BEFORE INSERT OR UPDATE OF status ON documents.documents
FOR EACH ROW EXECUTE FUNCTION documents.document_posting_guard();

COMMIT;
