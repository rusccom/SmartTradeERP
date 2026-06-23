BEGIN;

-- 1. Public visibility. Fail-closed: nothing is reachable by the public
--    storefront until it is explicitly published.
ALTER TABLE catalog.products
    ADD COLUMN IF NOT EXISTS is_published BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE catalog.product_variants
    ADD COLUMN IF NOT EXISTS is_published BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS ix_products_published
    ON catalog.products (tenant_id)
    WHERE is_published;

-- 2. Browseable taxonomy. Storefront navigation and category landing pages.
CREATE TABLE IF NOT EXISTS catalog.categories (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    parent_id UUID REFERENCES catalog.categories(id) ON DELETE SET NULL,
    name VARCHAR(160) NOT NULL,
    slug VARCHAR(200) NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_categories_slug_per_tenant
    ON catalog.categories (tenant_id, slug)
    WHERE slug <> '';

CREATE TABLE IF NOT EXISTS catalog.product_categories (
    tenant_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES catalog.categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_product_categories_category
    ON catalog.product_categories (category_id);

-- 3. Denormalized availability projection so public reads never touch the
--    ledger live. Refreshed on ledger writes (Phase 3); checkout re-validates.
CREATE TABLE IF NOT EXISTS catalog.product_availability (
    variant_id UUID PRIMARY KEY REFERENCES catalog.product_variants(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    available_qty DECIMAL(14,4) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- 4. Host -> tenant map. The storefront service resolves the request Host
--    against this table on every request (fail-closed on miss). Installed
--    themes are the embedded theme.json manifests (source of truth), so there
--    is no themes table: adding a theme is dropping a folder, validated in Go.
CREATE TABLE IF NOT EXISTS platform.storefront_domains (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES platform.tenants(id) ON DELETE CASCADE,
    host VARCHAR(255) NOT NULL UNIQUE,
    kind VARCHAR(16) NOT NULL DEFAULT 'subdomain' CHECK (kind IN ('subdomain', 'custom')),
    status VARCHAR(16) NOT NULL DEFAULT 'active' CHECK (status IN ('pending', 'verifying', 'active', 'disabled')),
    verify_token VARCHAR(120),
    tls_status VARCHAR(16) NOT NULL DEFAULT 'pending' CHECK (tls_status IN ('pending', 'active', 'error')),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_storefront_domains_tenant
    ON platform.storefront_domains (tenant_id);

-- 5. Per-tenant theme selection and customization. One row per tenant, so it
--    extends the existing platform.tenant_settings. draft_* hold unpublished
--    edits; publish copies draft_* -> published_* atomically. The theme id is
--    validated in Go against the installed themes, not by a DB foreign key.
ALTER TABLE platform.tenant_settings
    ADD COLUMN IF NOT EXISTS storefront_theme_id VARCHAR(64) NOT NULL DEFAULT 'classic',
    ADD COLUMN IF NOT EXISTS published_tokens JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS draft_tokens JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS draft_theme_id VARCHAR(64),
    ADD COLUMN IF NOT EXISTS sections JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS draft_sections JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS logo_media_id UUID;

COMMIT;
