BEGIN;

-- Tenant isolation becomes database-enforced. Every application query runs
-- inside a transaction that sets app.tenant_id (internal/shared/db), and the
-- policies below make other tenants' rows invisible even if a Go query ever
-- forgets its WHERE tenant_id clause. Fail-closed: with no app.tenant_id set,
-- tenant-owned tables appear empty.
--
-- Deliberately NOT under RLS:
--   platform.tenants, platform.tenant_users, platform.platform_admins -- read
--     before a tenant context exists (login, registration, admin console);
--   platform.currencies -- global dictionary shared by all tenants;
--   platform.storefront_domains -- the host->tenant lookup that CREATES the
--     tenant context for public storefront requests;
--   public.schema_migrations -- migration runner bookkeeping.
--
-- FORCE makes policies apply to the table owner too. A future migration that
-- must backfill across tenants should temporarily
-- ALTER TABLE ... DISABLE ROW LEVEL SECURITY inside its own transaction.

CREATE OR REPLACE FUNCTION platform.current_tenant() RETURNS uuid
    LANGUAGE sql STABLE
    AS $fn$ SELECT NULLIF(current_setting('app.tenant_id', true), '')::uuid $fn$;

-- 1. Tables that carry tenant_id directly.

DO $$
DECLARE
    t TEXT;
BEGIN
    FOREACH t IN ARRAY ARRAY[
        'catalog.products', 'catalog.product_variants', 'catalog.warehouses',
        'catalog.customers', 'catalog.categories', 'catalog.product_categories',
        'catalog.product_availability',
        'documents.documents', 'documents.shifts', 'documents.document_sequences',
        'ledger.posting_batches', 'ledger.inventory_movements',
        'ledger.cost_movement_results', 'ledger.stock_balances',
        'ledger.document_item_financials', 'ledger.daily_variant_metrics',
        'platform.tenant_settings', 'platform.tenant_currencies',
        'platform.media_objects',
        'storefront.orders', 'storefront.order_items'
    ] LOOP
        EXECUTE format('ALTER TABLE %s ENABLE ROW LEVEL SECURITY', t);
        EXECUTE format('ALTER TABLE %s FORCE ROW LEVEL SECURITY', t);
        EXECUTE format('DROP POLICY IF EXISTS tenant_isolation ON %s', t);
        EXECUTE format(
            'CREATE POLICY tenant_isolation ON %s
                USING (tenant_id = platform.current_tenant())
                WITH CHECK (tenant_id = platform.current_tenant())', t);
    END LOOP;
END $$;

-- 2. Line tables without tenant_id: scoped through their immediate parent.
--    The parent's own policy already filters by tenant, so visibility (and
--    insertability) of a line row follows visibility of its parent row.

DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT * FROM (VALUES
        ('catalog.product_options',            'catalog.products',         'product_id'),
        ('catalog.product_option_values',      'catalog.product_options',  'option_id'),
        ('catalog.variant_components',         'catalog.product_variants', 'variant_id'),
        ('documents.document_items',           'documents.documents',      'document_id'),
        ('documents.document_item_components', 'documents.document_items', 'document_item_id'),
        ('documents.shift_cash_ops',           'documents.shifts',         'shift_id'),
        ('documents.document_payments',        'documents.documents',      'document_id')
    ) AS v(child, parent, fk) LOOP
        EXECUTE format('ALTER TABLE %s ENABLE ROW LEVEL SECURITY', r.child);
        EXECUTE format('ALTER TABLE %s FORCE ROW LEVEL SECURITY', r.child);
        EXECUTE format('DROP POLICY IF EXISTS tenant_isolation ON %s', r.child);
        EXECUTE format(
            'CREATE POLICY tenant_isolation ON %s
                USING (EXISTS (SELECT 1 FROM %s p WHERE p.id = %s))
                WITH CHECK (EXISTS (SELECT 1 FROM %s p WHERE p.id = %s))',
            r.child, r.parent, r.fk, r.parent, r.fk);
    END LOOP;
END $$;

COMMIT;
