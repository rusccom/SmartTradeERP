BEGIN;

CREATE TABLE storefront.carts (
    tenant_id uuid NOT NULL REFERENCES platform.tenants(id) ON DELETE CASCADE,
    id uuid NOT NULL,
    preview boolean NOT NULL,
    version bigint NOT NULL DEFAULT 1,
    items jsonb NOT NULL CHECK (jsonb_typeof(items) = 'array'),
    PRIMARY KEY (tenant_id, id, preview)
);

ALTER TABLE storefront.carts ENABLE ROW LEVEL SECURITY;
ALTER TABLE storefront.carts FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON storefront.carts
    USING (tenant_id = platform.current_tenant()) WITH CHECK (tenant_id = platform.current_tenant());

COMMIT;
