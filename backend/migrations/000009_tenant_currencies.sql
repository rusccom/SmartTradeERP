BEGIN;

CREATE TABLE IF NOT EXISTS platform.currencies (
    id UUID PRIMARY KEY,
    code VARCHAR(3) NOT NULL UNIQUE,
    name VARCHAR(80) NOT NULL,
    symbol VARCHAR(8) NOT NULL DEFAULT '',
    decimal_places SMALLINT NOT NULL DEFAULT 2 CHECK (decimal_places BETWEEN 0 AND 4),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (code = upper(code) AND code ~ '^[A-Z]{3}$'),
    CHECK (btrim(name) <> '')
);

INSERT INTO platform.currencies (id, code, name, symbol, decimal_places)
VALUES
    ('00000000-0000-0000-0000-000000000840','USD','US Dollar','$',2),
    ('00000000-0000-0000-0000-000000000978','EUR','Euro','EUR',2),
    ('00000000-0000-0000-0000-000000000933','BYN','Belarusian Ruble','Br',2),
    ('00000000-0000-0000-0000-000000000985','PLN','Polish Zloty','zl',2),
    ('00000000-0000-0000-0000-000000000826','GBP','Pound Sterling','GBP',2),
    ('00000000-0000-0000-0000-000000000398','KZT','Kazakhstani Tenge','KZT',2)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    symbol = EXCLUDED.symbol,
    decimal_places = EXCLUDED.decimal_places,
    is_active = true,
    updated_at = now();

CREATE TABLE IF NOT EXISTS platform.tenant_currencies (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES platform.tenants(id) ON DELETE CASCADE,
    currency_id UUID NOT NULL REFERENCES platform.currencies(id),
    is_base BOOLEAN NOT NULL DEFAULT false,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    display_symbol VARCHAR(8),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, currency_id),
    UNIQUE (tenant_id, id)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_tenant_currencies_base
    ON platform.tenant_currencies (tenant_id)
    WHERE is_base = true;

CREATE INDEX IF NOT EXISTS idx_tenant_currencies_tenant
    ON platform.tenant_currencies (tenant_id);

ALTER TABLE platform.tenant_settings
    ADD COLUMN IF NOT EXISTS base_currency_id UUID REFERENCES platform.currencies(id);

COMMIT;
