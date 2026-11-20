BEGIN;

CREATE TABLE catalog.vat_countries (
    tenant_id uuid NOT NULL REFERENCES platform.tenants(id),
    country text NOT NULL CHECK (country ~ '^[A-Z]{2}$'),
    PRIMARY KEY (tenant_id, country)
);

CREATE TABLE catalog.product_vat_rates (
    tenant_id uuid NOT NULL,
    product_id uuid NOT NULL,
    country text NOT NULL,
    rate numeric NOT NULL CHECK (rate BETWEEN 0 AND 100 AND rate = round(rate, 2)),
    PRIMARY KEY (tenant_id, product_id, country),
    FOREIGN KEY (tenant_id, product_id) REFERENCES catalog.products (tenant_id, id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id, country) REFERENCES catalog.vat_countries (tenant_id, country)
);

ALTER TABLE catalog.vat_countries ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.vat_countries FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON catalog.vat_countries
    USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant());

ALTER TABLE catalog.product_vat_rates ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.product_vat_rates FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON catalog.product_vat_rates
    USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant());

ALTER TABLE storefront.order_items
    ADD COLUMN vat_country text CHECK (vat_country ~ '^[A-Z]{2}$'),
    ADD COLUMN vat_mode text NOT NULL DEFAULT 'manual' CHECK (vat_mode IN ('manual', 'catalog'));

COMMIT;
