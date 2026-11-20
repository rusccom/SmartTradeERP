BEGIN;

CREATE SCHEMA IF NOT EXISTS sales;
CREATE SCHEMA IF NOT EXISTS region_pl;

ALTER TABLE storefront.orders
    ADD COLUMN source text NOT NULL DEFAULT 'storefront',
    ADD COLUMN revision bigint NOT NULL DEFAULT 1,
    ADD COLUMN entry_hash bytea,
    ADD COLUMN delivery_cost numeric(14,4) CHECK (delivery_cost >= 0),
    ADD COLUMN invoice_requested boolean NOT NULL DEFAULT false,
    ADD COLUMN details jsonb NOT NULL DEFAULT '{}',
    ADD COLUMN billing_customer_id uuid;

ALTER TABLE storefront.orders ADD FOREIGN KEY (tenant_id, billing_customer_id)
    REFERENCES catalog.customers (tenant_id, id);

ALTER TABLE storefront.order_items
    ADD COLUMN position integer,
    ADD COLUMN vat_rate text,
    ADD COLUMN sku_code text,
    ADD COLUMN unit text;
CREATE INDEX IF NOT EXISTS idx_order_items_tenant_order ON storefront.order_items (tenant_id, order_id);

CREATE TABLE sales.sequences (
    tenant_id uuid NOT NULL REFERENCES platform.tenants(id),
    series text NOT NULL,
    year integer NOT NULL,
    last_number bigint NOT NULL CHECK (last_number > 0),
    PRIMARY KEY (tenant_id, series, year)
);
CREATE TABLE sales.order_payments (
    tenant_id uuid NOT NULL,
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    amount numeric(14,4) NOT NULL CHECK (amount > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, id),
    FOREIGN KEY (tenant_id, order_id) REFERENCES storefront.orders (tenant_id, id)
);
CREATE INDEX idx_order_payments_order ON sales.order_payments (tenant_id, order_id);

CREATE TABLE sales.invoices (
    tenant_id uuid NOT NULL,
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    number text NOT NULL,
    format text NOT NULL,
    issued_at timestamptz NOT NULL DEFAULT now(),
    snapshot jsonb NOT NULL,
    PRIMARY KEY (tenant_id, id),
    UNIQUE (tenant_id, order_id),
    UNIQUE (tenant_id, number),
    FOREIGN KEY (tenant_id, order_id) REFERENCES storefront.orders (tenant_id, id)
);

CREATE TABLE region_pl.sellers (
    tenant_id uuid PRIMARY KEY REFERENCES platform.tenants(id),
    name text NOT NULL,
    nip text NOT NULL,
    street text NOT NULL,
    postal_code text NOT NULL,
    city text NOT NULL,
    environment text NOT NULL CHECK (environment IN ('test', 'production')),
    encrypted_token bytea,
    revision bigint NOT NULL DEFAULT 1
);
CREATE TABLE region_pl.invoice_payloads (
    tenant_id uuid NOT NULL,
    invoice_id uuid NOT NULL,
    payload bytea NOT NULL,
    PRIMARY KEY (tenant_id, invoice_id),
    FOREIGN KEY (tenant_id, invoice_id) REFERENCES sales.invoices (tenant_id, id)
);
CREATE TABLE region_pl.invoice_submissions (
    tenant_id uuid NOT NULL,
    invoice_id uuid NOT NULL,
    state text NOT NULL DEFAULT 'ready' CHECK (state IN ('ready','sending','pending','accepted','rejected','uncertain')),
    environment text NOT NULL CHECK (environment IN ('test','production')),
    attempt_id uuid,
    session_ref text,
    invoice_ref text,
    ksef_number text,
    message text,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, invoice_id, environment),
    FOREIGN KEY (tenant_id, invoice_id) REFERENCES region_pl.invoice_payloads (tenant_id, invoice_id),
    CHECK (state <> 'accepted' OR ksef_number IS NOT NULL)
);

DO $rls$
DECLARE relation regclass;
BEGIN
    FOREACH relation IN ARRAY ARRAY[
        'sales.sequences'::regclass, 'sales.order_payments'::regclass, 'sales.invoices'::regclass,
        'region_pl.sellers'::regclass, 'region_pl.invoice_payloads'::regclass, 'region_pl.invoice_submissions'::regclass
    ] LOOP
        EXECUTE format('ALTER TABLE %s ENABLE ROW LEVEL SECURITY', relation);
        EXECUTE format('ALTER TABLE %s FORCE ROW LEVEL SECURITY', relation);
        EXECUTE format('CREATE POLICY tenant_isolation ON %s USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant())', relation);
    END LOOP;
END $rls$;

COMMIT;
