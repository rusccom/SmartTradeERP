BEGIN;

ALTER TABLE platform.allegro_connections
    ADD COLUMN private_client_id TEXT CHECK (private_client_id IS NULL OR method = 'token'),
    ADD COLUMN refreshed_at TIMESTAMPTZ;

CREATE TABLE platform.allegro_device_requests (
    tenant_id UUID PRIMARY KEY REFERENCES platform.tenants(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES platform.tenant_users(id) ON DELETE CASCADE,
    request_id UUID NOT NULL,
    connection_version TEXT NOT NULL,
    encrypted_code BYTEA NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    poll_after TIMESTAMPTZ NOT NULL,
    poll_seconds INTEGER NOT NULL CHECK (poll_seconds > 0)
);

ALTER TABLE platform.allegro_device_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE platform.allegro_device_requests FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON platform.allegro_device_requests
    USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant());

COMMIT;
