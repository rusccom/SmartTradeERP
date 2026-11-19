BEGIN;

CREATE TABLE IF NOT EXISTS platform.inpost_connections (
    tenant_id UUID PRIMARY KEY REFERENCES platform.tenants(id) ON DELETE CASCADE,
    environment TEXT NOT NULL CHECK (environment IN ('production', 'sandbox')),
    organization_id BIGINT NOT NULL CHECK (organization_id > 0),
    encrypted_token BYTEA NOT NULL,
    organization_name TEXT NOT NULL,
    services TEXT[] NOT NULL,
    verified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revision BIGINT NOT NULL DEFAULT 1
);

ALTER TABLE platform.inpost_connections ENABLE ROW LEVEL SECURITY;
ALTER TABLE platform.inpost_connections FORCE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS tenant_isolation ON platform.inpost_connections;
CREATE POLICY tenant_isolation ON platform.inpost_connections
    USING (tenant_id = platform.current_tenant())
    WITH CHECK (tenant_id = platform.current_tenant());

COMMIT;
