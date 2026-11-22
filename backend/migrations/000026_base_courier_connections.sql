BEGIN;

CREATE TABLE platform.base_courier_connections (
    tenant_id UUID PRIMARY KEY REFERENCES platform.tenants(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    encrypted_api_key BYTEA NOT NULL,
    version UUID NOT NULL DEFAULT gen_random_uuid(),
    verified_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE platform.base_courier_connections ENABLE ROW LEVEL SECURITY;
ALTER TABLE platform.base_courier_connections FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON platform.base_courier_connections
    USING (tenant_id = platform.current_tenant())
    WITH CHECK (tenant_id = platform.current_tenant());

COMMIT;
