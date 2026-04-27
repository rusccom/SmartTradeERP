BEGIN;

CREATE UNIQUE INDEX IF NOT EXISTS ux_tenant_users_tenant_id_id
    ON platform.tenant_users (tenant_id, id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_warehouses_tenant_id_id
    ON catalog.warehouses (tenant_id, id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_customers_tenant_id_id
    ON catalog.customers (tenant_id, id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_documents_tenant_id_id
    ON documents.documents (tenant_id, id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_shifts_tenant_id_id
    ON documents.shifts (tenant_id, id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_document_items_document_id_id
    ON documents.document_items (document_id, id);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM documents.shifts s
        JOIN platform.tenant_users u ON u.id = s.user_id
        WHERE u.tenant_id IS DISTINCT FROM s.tenant_id
    ) THEN
        RAISE EXCEPTION 'cross-tenant shift user references exist';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM documents.shifts s
        JOIN catalog.warehouses w ON w.id = s.warehouse_id
        WHERE w.tenant_id IS DISTINCT FROM s.tenant_id
    ) THEN
        RAISE EXCEPTION 'cross-tenant shift warehouse references exist';
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM documents.document_items i
        JOIN documents.documents d ON d.id = i.document_id
        JOIN catalog.product_variants v ON v.id = i.variant_id
        JOIN catalog.products p ON p.id = v.product_id
        WHERE p.tenant_id IS DISTINCT FROM d.tenant_id
    ) THEN
        RAISE EXCEPTION 'cross-tenant document item variant references exist';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM documents.document_item_components c
        JOIN documents.document_items i ON i.id = c.document_item_id
        JOIN documents.documents d ON d.id = i.document_id
        JOIN catalog.product_variants v ON v.id = c.component_variant_id
        JOIN catalog.products p ON p.id = v.product_id
        WHERE p.tenant_id IS DISTINCT FROM d.tenant_id
    ) THEN
        RAISE EXCEPTION 'cross-tenant document item component references exist';
    END IF;
END $$;

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
    'catalog.products'::regclass,
    'fk_products_tenant',
    'ALTER TABLE catalog.products ADD CONSTRAINT fk_products_tenant
     FOREIGN KEY (tenant_id) REFERENCES platform.tenants (id)'
);

SELECT public.add_constraint_if_missing(
    'catalog.warehouses'::regclass,
    'fk_warehouses_tenant',
    'ALTER TABLE catalog.warehouses ADD CONSTRAINT fk_warehouses_tenant
     FOREIGN KEY (tenant_id) REFERENCES platform.tenants (id)'
);

SELECT public.add_constraint_if_missing(
    'catalog.customers'::regclass,
    'fk_customers_tenant',
    'ALTER TABLE catalog.customers ADD CONSTRAINT fk_customers_tenant
     FOREIGN KEY (tenant_id) REFERENCES platform.tenants (id)'
);

SELECT public.add_constraint_if_missing(
    'documents.documents'::regclass,
    'fk_documents_tenant',
    'ALTER TABLE documents.documents ADD CONSTRAINT fk_documents_tenant
     FOREIGN KEY (tenant_id) REFERENCES platform.tenants (id)'
);

SELECT public.add_constraint_if_missing(
    'documents.shifts'::regclass,
    'fk_shifts_tenant',
    'ALTER TABLE documents.shifts ADD CONSTRAINT fk_shifts_tenant
     FOREIGN KEY (tenant_id) REFERENCES platform.tenants (id)'
);

SELECT public.add_constraint_if_missing(
    'documents.shifts'::regclass,
    'fk_shifts_tenant_user',
    'ALTER TABLE documents.shifts ADD CONSTRAINT fk_shifts_tenant_user
     FOREIGN KEY (tenant_id, user_id) REFERENCES platform.tenant_users (tenant_id, id)'
);

SELECT public.add_constraint_if_missing(
    'documents.shifts'::regclass,
    'fk_shifts_tenant_warehouse',
    'ALTER TABLE documents.shifts ADD CONSTRAINT fk_shifts_tenant_warehouse
     FOREIGN KEY (tenant_id, warehouse_id) REFERENCES catalog.warehouses (tenant_id, id)'
);

SELECT public.add_constraint_if_missing(
    'documents.documents'::regclass,
    'fk_documents_tenant_warehouse',
    'ALTER TABLE documents.documents ADD CONSTRAINT fk_documents_tenant_warehouse
     FOREIGN KEY (tenant_id, warehouse_id) REFERENCES catalog.warehouses (tenant_id, id)'
);

SELECT public.add_constraint_if_missing(
    'documents.documents'::regclass,
    'fk_documents_tenant_source_warehouse',
    'ALTER TABLE documents.documents ADD CONSTRAINT fk_documents_tenant_source_warehouse
     FOREIGN KEY (tenant_id, source_warehouse_id) REFERENCES catalog.warehouses (tenant_id, id)'
);

SELECT public.add_constraint_if_missing(
    'documents.documents'::regclass,
    'fk_documents_tenant_target_warehouse',
    'ALTER TABLE documents.documents ADD CONSTRAINT fk_documents_tenant_target_warehouse
     FOREIGN KEY (tenant_id, target_warehouse_id) REFERENCES catalog.warehouses (tenant_id, id)'
);

SELECT public.add_constraint_if_missing(
    'documents.documents'::regclass,
    'fk_documents_tenant_customer',
    'ALTER TABLE documents.documents ADD CONSTRAINT fk_documents_tenant_customer
     FOREIGN KEY (tenant_id, customer_id) REFERENCES catalog.customers (tenant_id, id)'
);

SELECT public.add_constraint_if_missing(
    'documents.documents'::regclass,
    'fk_documents_tenant_shift',
    'ALTER TABLE documents.documents ADD CONSTRAINT fk_documents_tenant_shift
     FOREIGN KEY (tenant_id, shift_id) REFERENCES documents.shifts (tenant_id, id)'
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

    SELECT p.tenant_id INTO variant_tenant
    FROM catalog.product_variants v
    JOIN catalog.products p ON p.id = v.product_id
    WHERE v.id = NEW.variant_id;

    IF doc_tenant IS NULL OR variant_tenant IS DISTINCT FROM doc_tenant THEN
        RAISE EXCEPTION 'document item variant must belong to document tenant';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_document_item_tenant_guard ON documents.document_items;

CREATE TRIGGER trg_document_item_tenant_guard
BEFORE INSERT OR UPDATE ON documents.document_items
FOR EACH ROW EXECUTE FUNCTION documents.document_item_tenant_guard();

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

    SELECT p.tenant_id INTO component_tenant
    FROM catalog.product_variants v
    JOIN catalog.products p ON p.id = v.product_id
    WHERE v.id = NEW.component_variant_id;

    IF doc_tenant IS NULL OR component_tenant IS DISTINCT FROM doc_tenant THEN
        RAISE EXCEPTION 'item component variant must belong to document tenant';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_document_item_component_tenant_guard ON documents.document_item_components;

CREATE TRIGGER trg_document_item_component_tenant_guard
BEFORE INSERT OR UPDATE ON documents.document_item_components
FOR EACH ROW EXECUTE FUNCTION documents.document_item_component_tenant_guard();

COMMIT;
