BEGIN;

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

DROP FUNCTION public.add_constraint_if_missing(regclass, text, text);

COMMIT;
