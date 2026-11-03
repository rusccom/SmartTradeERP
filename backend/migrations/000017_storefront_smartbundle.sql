BEGIN;

-- ERP owns component quantities. Storefront rules only describe shopper choices and pricing.
CREATE TABLE storefront.smartbundle_rules (
    tenant_id uuid NOT NULL,
    variant_id uuid NOT NULL,
    price_mode text NOT NULL DEFAULT 'fixed' CHECK (price_mode IN ('fixed', 'dynamic')),
    discount numeric(5,2) NOT NULL DEFAULT 0 CHECK (discount BETWEEN 0 AND 100),
    slots jsonb NOT NULL DEFAULT '{}' CHECK (jsonb_typeof(slots) = 'object'),
    PRIMARY KEY (tenant_id, variant_id),
    FOREIGN KEY (tenant_id, variant_id) REFERENCES catalog.product_variants (tenant_id, id) ON DELETE CASCADE
);

ALTER TABLE storefront.smartbundle_rules ENABLE ROW LEVEL SECURITY;
ALTER TABLE storefront.smartbundle_rules FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON storefront.smartbundle_rules
    USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant());

ALTER TABLE storefront.order_items
    ADD COLUMN bundle_components jsonb NOT NULL DEFAULT '[]';

COMMIT;
