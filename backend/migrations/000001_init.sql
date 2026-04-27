BEGIN;

CREATE SCHEMA IF NOT EXISTS platform;
CREATE SCHEMA IF NOT EXISTS catalog;
CREATE SCHEMA IF NOT EXISTS documents;
CREATE SCHEMA IF NOT EXISTS ledger;

CREATE TABLE IF NOT EXISTS platform.tenants (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'trial' CHECK (status IN ('trial', 'active', 'suspended')),
    plan VARCHAR NOT NULL DEFAULT 'free',
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS platform.tenant_users (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES platform.tenants(id),
    email VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    role VARCHAR NOT NULL DEFAULT 'owner' CHECK (role IN ('owner', 'manager', 'cashier')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS platform.platform_admins (
    id UUID PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS catalog.products (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR NOT NULL,
    is_composite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS catalog.product_options (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    position SMALLINT NOT NULL CHECK (position BETWEEN 1 AND 3),
    UNIQUE (product_id, position)
);

CREATE TABLE IF NOT EXISTS catalog.product_option_values (
    id UUID PRIMARY KEY,
    option_id UUID NOT NULL REFERENCES catalog.product_options(id) ON DELETE CASCADE,
    value VARCHAR NOT NULL,
    position SMALLINT NOT NULL,
    UNIQUE (option_id, position)
);

CREATE TABLE IF NOT EXISTS catalog.product_variants (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    name VARCHAR,
    sku_code VARCHAR,
    barcode VARCHAR,
    unit VARCHAR NOT NULL,
    price DECIMAL(12,4),
    option1 VARCHAR,
    option2 VARCHAR,
    option3 VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS catalog.variant_components (
    id UUID PRIMARY KEY,
    variant_id UUID NOT NULL REFERENCES catalog.product_variants(id) ON DELETE CASCADE,
    component_variant_id UUID NOT NULL REFERENCES catalog.product_variants(id),
    qty DECIMAL(12,3) NOT NULL CHECK (qty > 0),
    CONSTRAINT no_self_reference CHECK (variant_id <> component_variant_id),
    UNIQUE (variant_id, component_variant_id)
);

CREATE TABLE IF NOT EXISTS catalog.warehouses (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    address VARCHAR,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_warehouse_default_per_tenant
    ON catalog.warehouses (tenant_id)
    WHERE is_default = true;

CREATE TABLE IF NOT EXISTS documents.documents (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    type VARCHAR NOT NULL CHECK (type IN ('RECEIPT', 'SALE', 'WRITEOFF', 'INVENTORY', 'TRANSFER')),
    date DATE NOT NULL,
    number VARCHAR,
    status VARCHAR NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'posted', 'cancelled')),
    warehouse_id UUID REFERENCES catalog.warehouses(id),
    source_warehouse_id UUID REFERENCES catalog.warehouses(id),
    target_warehouse_id UUID REFERENCES catalog.warehouses(id),
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS documents.document_items (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents.documents(id) ON DELETE CASCADE,
    variant_id UUID NOT NULL REFERENCES catalog.product_variants(id),
    qty DECIMAL(12,3) NOT NULL CHECK (qty > 0),
    unit_price DECIMAL(12,4) NOT NULL,
    total_amount DECIMAL(14,4) NOT NULL
);

CREATE TABLE IF NOT EXISTS documents.document_item_components (
    id UUID PRIMARY KEY,
    document_item_id UUID NOT NULL REFERENCES documents.document_items(id) ON DELETE CASCADE,
    component_variant_id UUID NOT NULL REFERENCES catalog.product_variants(id),
    qty_per_unit DECIMAL(12,3) NOT NULL CHECK (qty_per_unit > 0),
    qty_total DECIMAL(12,3) NOT NULL CHECK (qty_total > 0)
);

CREATE OR REPLACE FUNCTION catalog.variant_component_guard()
RETURNS trigger AS $$
DECLARE
    variant_is_composite BOOLEAN;
    component_is_composite BOOLEAN;
    variant_tenant UUID;
    component_tenant UUID;
BEGIN
    SELECT p.is_composite, p.tenant_id
    INTO variant_is_composite, variant_tenant
    FROM catalog.product_variants v
    JOIN catalog.products p ON p.id = v.product_id
    WHERE v.id = NEW.variant_id;

    SELECT p.is_composite, p.tenant_id
    INTO component_is_composite, component_tenant
    FROM catalog.product_variants v
    JOIN catalog.products p ON p.id = v.product_id
    WHERE v.id = NEW.component_variant_id;

    IF variant_is_composite IS DISTINCT FROM true THEN
        RAISE EXCEPTION 'variant must belong to composite product';
    END IF;

    IF component_is_composite IS DISTINCT FROM false THEN
        RAISE EXCEPTION 'component variant must be non-composite';
    END IF;

    IF variant_tenant IS DISTINCT FROM component_tenant THEN
        RAISE EXCEPTION 'cross-tenant components are forbidden';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_variant_component_guard ON catalog.variant_components;

CREATE TRIGGER trg_variant_component_guard
BEFORE INSERT OR UPDATE ON catalog.variant_components
FOR EACH ROW EXECUTE FUNCTION catalog.variant_component_guard();

COMMIT;
