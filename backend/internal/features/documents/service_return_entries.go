package documents

import (
	"context"

	"smarterp/backend/internal/features/ledger"
)

func (s *Service) buildReturnEntry(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
) ([]ledger.EntryInput, error) {
	_, avg, err := s.ledger.GlobalStock(ctx, tenantID, item.VariantID)
	if err != nil {
		return nil, err
	}
	revenue := item.TotalAmount.Neg()
	total := item.Qty.Mul(avg).Round(4)
	entry := makeEntry(tenantID, doc.ID, item.ID, item.VariantID, doc.WarehouseID,
		mustDate(doc.Date), "IN", "RETURN_IN", item.Qty, avg, total, &revenue)
	return []ledger.EntryInput{entry}, nil
}
