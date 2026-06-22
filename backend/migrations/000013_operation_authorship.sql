BEGIN;

ALTER TABLE documents.documents
    ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES platform.tenant_users(id);

ALTER TABLE ledger.posting_batches
    ADD COLUMN IF NOT EXISTS posted_by UUID REFERENCES platform.tenant_users(id);

COMMIT;
