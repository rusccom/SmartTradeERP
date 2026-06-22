BEGIN;

ALTER TABLE catalog.products
    ADD COLUMN IF NOT EXISTS slug VARCHAR(200) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS seo_title VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS seo_description VARCHAR(320) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS ux_products_slug_per_tenant
    ON catalog.products (tenant_id, slug)
    WHERE slug <> '';

COMMIT;
