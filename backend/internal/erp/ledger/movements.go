package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) Append(ctx context.Context, tx pgx.Tx, batchID string, input EntryInput) (string, error) {
	id := uuid.NewString()
	query := `INSERT INTO ledger.inventory_movements
        (id, tenant_id, posting_batch_id, document_id, document_item_id,
         variant_id, warehouse_id, movement_date, direction, reason, qty,
         unit_price, total_amount, revenue_amount, posting_order)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,
            (SELECT COALESCE(MAX(posting_order), 0) + 1
             FROM ledger.inventory_movements
             WHERE tenant_id=$2 AND posting_batch_id=$3))`
	_, err := tx.Exec(ctx, query, id, input.TenantID, batchID, input.DocumentID,
		input.DocumentItemID, input.VariantID, input.WarehouseID, input.Date,
		input.Type, input.Reason, input.Qty, input.UnitPrice, input.TotalAmount,
		input.Revenue)
	return id, err
}

func AffectedFromEntries(entries []EntryInput) []VariantSequence {
	seen := map[string]bool{}
	affected := make([]VariantSequence, 0)
	for _, entry := range entries {
		if seen[entry.VariantID] {
			continue
		}
		seen[entry.VariantID] = true
		affected = append(affected, VariantSequence{VariantID: entry.VariantID})
	}
	return affected
}

func MergeAffected(items ...[]VariantSequence) []VariantSequence {
	seen := map[string]bool{}
	merged := make([]VariantSequence, 0)
	for _, group := range items {
		merged = appendAffectedGroup(merged, seen, group)
	}
	return merged
}

func appendAffectedGroup(
	merged []VariantSequence,
	seen map[string]bool,
	group []VariantSequence,
) []VariantSequence {
	for _, item := range group {
		if item.VariantID == "" || seen[item.VariantID] {
			continue
		}
		seen[item.VariantID] = true
		merged = append(merged, item)
	}
	return merged
}
