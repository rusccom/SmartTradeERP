BEGIN;

CREATE TABLE IF NOT EXISTS documents.document_sequences (
    tenant_id UUID NOT NULL REFERENCES platform.tenants(id) ON DELETE CASCADE,
    document_type VARCHAR NOT NULL CHECK (
        document_type IN ('RECEIPT', 'SALE', 'WRITEOFF', 'INVENTORY', 'TRANSFER', 'RETURN')
    ),
    year INTEGER NOT NULL CHECK (year BETWEEN 2000 AND 9999),
    next_number INTEGER NOT NULL DEFAULT 1 CHECK (next_number > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, document_type, year)
);

WITH parsed AS (
    SELECT tenant_id, type, parts[1]::integer AS year, parts[2]::integer AS sequence_num
    FROM (
        SELECT tenant_id, type, regexp_match(number, '^' || type || '-([0-9]{4})-([0-9]{6})$') AS parts
        FROM documents.documents
        WHERE number IS NOT NULL AND btrim(number) <> ''
    ) numbered
    WHERE parts IS NOT NULL
)
INSERT INTO documents.document_sequences (tenant_id, document_type, year, next_number)
SELECT tenant_id, type, year, MAX(sequence_num) + 1
FROM parsed
GROUP BY tenant_id, type, year
ON CONFLICT (tenant_id, document_type, year)
DO UPDATE SET
    next_number = GREATEST(documents.document_sequences.next_number, EXCLUDED.next_number),
    updated_at = now();

WITH pending AS (
    SELECT id, tenant_id, type, EXTRACT(YEAR FROM date)::integer AS year,
        (ROW_NUMBER() OVER (
            PARTITION BY tenant_id, type, EXTRACT(YEAR FROM date)::integer
            ORDER BY date, created_at, id
        ))::integer AS row_num
    FROM documents.documents
    WHERE number IS NULL OR btrim(number) = ''
),
assigned AS (
    SELECT pending.id, pending.tenant_id, pending.type, pending.year,
        pending.type || '-' || pending.year::text || '-' ||
            lpad((COALESCE(s.next_number, 1) + pending.row_num - 1)::text, 6, '0') AS number,
        COALESCE(s.next_number, 1) + pending.row_num AS next_number
    FROM pending
    LEFT JOIN documents.document_sequences s
        ON s.tenant_id = pending.tenant_id
        AND s.document_type = pending.type
        AND s.year = pending.year
),
updated AS (
    UPDATE documents.documents d
    SET number = assigned.number
    FROM assigned
    WHERE d.id = assigned.id
    RETURNING d.id
)
INSERT INTO documents.document_sequences (tenant_id, document_type, year, next_number)
SELECT assigned.tenant_id, assigned.type, assigned.year, MAX(assigned.next_number)
FROM assigned
JOIN updated ON updated.id = assigned.id
GROUP BY assigned.tenant_id, assigned.type, assigned.year
ON CONFLICT (tenant_id, document_type, year)
DO UPDATE SET
    next_number = GREATEST(documents.document_sequences.next_number, EXCLUDED.next_number),
    updated_at = now();

ALTER TABLE documents.documents
    ALTER COLUMN number SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'documents.documents'::regclass
            AND conname = 'chk_documents_number_not_blank'
    ) THEN
        ALTER TABLE documents.documents
            ADD CONSTRAINT chk_documents_number_not_blank CHECK (btrim(number) <> '');
    END IF;
END;
$$ LANGUAGE plpgsql;

COMMIT;
