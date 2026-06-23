BEGIN;

CREATE SCHEMA IF NOT EXISTS storefront;

-- Customer orders placed through a tenant's public storefront. These are NOT
-- ERP documents: the merchant reviews and fulfils them (which posts a SALE to
-- the ledger) from the dashboard. Prices are stamped server-side at checkout.
CREATE TABLE IF NOT EXISTS storefront.orders (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    number VARCHAR(40) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'confirmed', 'cancelled', 'fulfilled')),
    customer_name VARCHAR(160) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(40) NOT NULL DEFAULT '',
    shipping_address TEXT NOT NULL DEFAULT '',
    currency_code VARCHAR(3) NOT NULL DEFAULT '',
    total_amount DECIMAL(14,4) NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_storefront_orders_number
    ON storefront.orders (tenant_id, number);
CREATE INDEX IF NOT EXISTS idx_storefront_orders_tenant
    ON storefront.orders (tenant_id, created_at DESC);

CREATE TABLE IF NOT EXISTS storefront.order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES storefront.orders(id) ON DELETE CASCADE,
    variant_id UUID NOT NULL REFERENCES catalog.product_variants(id),
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(255) NOT NULL DEFAULT '',
    unit_price DECIMAL(12,4) NOT NULL,
    qty DECIMAL(12,3) NOT NULL CHECK (qty > 0),
    total_amount DECIMAL(14,4) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_storefront_order_items_order
    ON storefront.order_items (order_id);

COMMIT;
