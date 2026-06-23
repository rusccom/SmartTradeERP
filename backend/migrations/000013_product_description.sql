BEGIN;

ALTER TABLE catalog.products
    ADD COLUMN IF NOT EXISTS description_html TEXT;

COMMIT;
