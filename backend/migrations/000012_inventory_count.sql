BEGIN;

ALTER TABLE documents.document_items
    DROP CONSTRAINT IF EXISTS document_items_qty_check;
ALTER TABLE documents.document_items
    ADD CONSTRAINT document_items_qty_check CHECK (qty >= 0);

ALTER TABLE ledger.inventory_movements
    DROP CONSTRAINT IF EXISTS inventory_movements_direction_check;
ALTER TABLE ledger.inventory_movements
    ADD CONSTRAINT inventory_movements_direction_check
        CHECK (direction IN ('IN', 'OUT', 'SET'));

ALTER TABLE ledger.inventory_movements
    DROP CONSTRAINT IF EXISTS inventory_movements_qty_check;
ALTER TABLE ledger.inventory_movements
    ADD CONSTRAINT inventory_movements_qty_check
        CHECK (qty > 0 OR (direction = 'SET' AND qty >= 0));

ALTER TABLE ledger.inventory_movements
    DROP CONSTRAINT IF EXISTS inventory_movements_reason_check;
ALTER TABLE ledger.inventory_movements
    ADD CONSTRAINT inventory_movements_reason_check
        CHECK (reason IN (
            'PURCHASE', 'SALE', 'WRITEOFF', 'COUNT',
            'RETURN_IN', 'TRANSFER_IN', 'TRANSFER_OUT'
        ));

COMMIT;
