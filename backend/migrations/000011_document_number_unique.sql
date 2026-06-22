BEGIN;

-- Clean cut: operational data is reset, catalog/master data is kept.
TRUNCATE TABLE
    ledger.cost_movement_results,
    ledger.inventory_movements,
    ledger.posting_batches,
    ledger.stock_balances,
    ledger.daily_variant_metrics,
    ledger.document_item_financials,
    documents.document_item_components,
    documents.document_payments,
    documents.document_items,
    documents.documents,
    documents.shift_cash_ops,
    documents.shifts,
    documents.document_sequences;

CREATE UNIQUE INDEX IF NOT EXISTS ux_documents_tenant_number
    ON documents.documents (tenant_id, number);

COMMIT;
