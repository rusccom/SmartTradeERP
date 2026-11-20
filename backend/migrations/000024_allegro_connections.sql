BEGIN;

CREATE TABLE IF NOT EXISTS platform.allegro_connections (
    tenant_id UUID PRIMARY KEY REFERENCES platform.tenants(id) ON DELETE CASCADE,
    seller_id TEXT NOT NULL,
    method TEXT NOT NULL CHECK (method IN ('token', 'oauth')),
    version UUID NOT NULL,
    login TEXT NOT NULL,
    encrypted_tokens BYTEA NOT NULL,
    expires_at TIMESTAMPTZ,
    verified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    needs_reconnect BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS platform.allegro_authorizations (
    tenant_id UUID PRIMARY KEY REFERENCES platform.tenants(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES platform.tenant_users(id) ON DELETE CASCADE,
    nonce UUID NOT NULL,
    expected_version UUID,
    launched BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE platform.allegro_connections ENABLE ROW LEVEL SECURITY;
ALTER TABLE platform.allegro_connections FORCE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS tenant_isolation ON platform.allegro_connections;
CREATE POLICY tenant_isolation ON platform.allegro_connections
    USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant());

ALTER TABLE platform.allegro_authorizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE platform.allegro_authorizations FORCE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS tenant_isolation ON platform.allegro_authorizations;
CREATE POLICY tenant_isolation ON platform.allegro_authorizations
    USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant());

COMMIT;
