BEGIN;

ALTER TABLE platform.tenant_settings
    DROP COLUMN IF EXISTS published_tokens,
    DROP COLUMN IF EXISTS draft_tokens,
    DROP COLUMN IF EXISTS sections,
    DROP COLUMN IF EXISTS draft_sections,
    ALTER COLUMN storefront_theme_id SET DEFAULT '';

COMMIT;
