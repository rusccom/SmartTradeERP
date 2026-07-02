BEGIN;

-- Restores the composite-FK rigor of 000004/000006 that migrations 000012 and
-- 000014 dropped: every tenant-owned row must reference same-tenant parents,
-- enforced by the database, not only by Go code.

-- 1. storefront.orders: tenant must exist; (tenant_id, id) key for composite
--    references from line rows.

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_storefront_orders_tenant') THEN
        ALTER TABLE storefront.orders
            ADD CONSTRAINT fk_storefront_orders_tenant
            FOREIGN KEY (tenant_id) REFERENCES platform.tenants(id);
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS ux_storefront_orders_tenant_id
    ON storefront.orders (tenant_id, id);

-- 2. storefront.order_items: add tenant_id and pair every reference, so a
--    line can only point at an order and a variant of the same tenant.

ALTER TABLE storefront.order_items
    ADD COLUMN IF NOT EXISTS tenant_id UUID;

UPDATE storefront.order_items oi
    SET tenant_id = o.tenant_id
    FROM storefront.orders o
    WHERE o.id = oi.order_id AND oi.tenant_id IS NULL;

ALTER TABLE storefront.order_items
    ALTER COLUMN tenant_id SET NOT NULL;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_order_items_tenant_order') THEN
        ALTER TABLE storefront.order_items
            ADD CONSTRAINT fk_order_items_tenant_order
            FOREIGN KEY (tenant_id, order_id)
            REFERENCES storefront.orders (tenant_id, id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_order_items_tenant_variant') THEN
        ALTER TABLE storefront.order_items
            ADD CONSTRAINT fk_order_items_tenant_variant
            FOREIGN KEY (tenant_id, variant_id)
            REFERENCES catalog.product_variants (tenant_id, id);
    END IF;
END $$;

-- 3. catalog.categories: tenant FK, (tenant_id, id) key, and a same-tenant
--    guard on parent_id. The original single-column parent FK keeps its
--    ON DELETE SET NULL behaviour; the composite one passes on NULL parent_id
--    (MATCH SIMPLE) and blocks cross-tenant parents otherwise.

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_categories_tenant') THEN
        ALTER TABLE catalog.categories
            ADD CONSTRAINT fk_categories_tenant
            FOREIGN KEY (tenant_id) REFERENCES platform.tenants(id);
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS ux_categories_tenant_id
    ON catalog.categories (tenant_id, id);

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_categories_tenant_parent') THEN
        ALTER TABLE catalog.categories
            ADD CONSTRAINT fk_categories_tenant_parent
            FOREIGN KEY (tenant_id, parent_id)
            REFERENCES catalog.categories (tenant_id, id);
    END IF;
END $$;

-- 4. catalog.product_categories: tenant FK plus same-tenant pairing to both
--    products and categories. Original CASCADE FKs stay for delete behaviour.

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_product_categories_tenant') THEN
        ALTER TABLE catalog.product_categories
            ADD CONSTRAINT fk_product_categories_tenant
            FOREIGN KEY (tenant_id) REFERENCES platform.tenants(id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_product_categories_tenant_product') THEN
        ALTER TABLE catalog.product_categories
            ADD CONSTRAINT fk_product_categories_tenant_product
            FOREIGN KEY (tenant_id, product_id)
            REFERENCES catalog.products (tenant_id, id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_product_categories_tenant_category') THEN
        ALTER TABLE catalog.product_categories
            ADD CONSTRAINT fk_product_categories_tenant_category
            FOREIGN KEY (tenant_id, category_id)
            REFERENCES catalog.categories (tenant_id, id);
    END IF;
END $$;

-- 5. catalog.product_availability: tenant FK and same-tenant variant pairing.

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_product_availability_tenant') THEN
        ALTER TABLE catalog.product_availability
            ADD CONSTRAINT fk_product_availability_tenant
            FOREIGN KEY (tenant_id) REFERENCES platform.tenants(id);
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_product_availability_tenant_variant') THEN
        ALTER TABLE catalog.product_availability
            ADD CONSTRAINT fk_product_availability_tenant_variant
            FOREIGN KEY (tenant_id, variant_id)
            REFERENCES catalog.product_variants (tenant_id, id);
    END IF;
END $$;

COMMIT;
