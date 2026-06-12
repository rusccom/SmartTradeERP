BEGIN;

CREATE TABLE IF NOT EXISTS platform.media_objects (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES platform.tenants (id),
    owner_type VARCHAR(80) NOT NULL,
    owner_id UUID NOT NULL,
    object_key TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL DEFAULT '',
    content_type VARCHAR(80) NOT NULL,
    size_bytes BIGINT NOT NULL CHECK (size_bytes > 0),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    status VARCHAR(20) NOT NULL DEFAULT 'ready' CHECK (status IN ('pending', 'ready')),
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, object_key),
    UNIQUE (tenant_id, id)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_media_objects_primary
    ON platform.media_objects (tenant_id, owner_type, owner_id)
    WHERE is_primary = true AND status = 'ready';

CREATE INDEX IF NOT EXISTS idx_media_objects_owner
    ON platform.media_objects (tenant_id, owner_type, owner_id, status, sort_order, created_at);

CREATE INDEX IF NOT EXISTS idx_media_objects_pending_expiry
    ON platform.media_objects (status, expires_at)
    WHERE status = 'pending';

COMMIT;
