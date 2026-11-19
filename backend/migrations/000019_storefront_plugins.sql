BEGIN;

CREATE TABLE platform.plugins (
    key text PRIMARY KEY,
    name text NOT NULL
);
INSERT INTO platform.plugins (key, name) VALUES ('smartbundle', 'SmartBundle');

CREATE TABLE platform.tenant_plugins (
    tenant_id uuid NOT NULL REFERENCES platform.tenants(id) ON DELETE CASCADE,
    plugin_key text NOT NULL REFERENCES platform.plugins(key),
    enabled boolean NOT NULL DEFAULT false,
    PRIMARY KEY (tenant_id, plugin_key)
);
ALTER TABLE platform.tenant_plugins ENABLE ROW LEVEL SECURITY;
ALTER TABLE platform.tenant_plugins FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON platform.tenant_plugins
    USING (tenant_id = platform.current_tenant())
    WITH CHECK (tenant_id = platform.current_tenant());

COMMIT;
