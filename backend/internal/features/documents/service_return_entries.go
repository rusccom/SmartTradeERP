package documents

import "smarterp/backend/internal/features/ledger"

func (s *Service) buildReturnEntry(
	run postingRun,
	item postingItem,
) ([]ledger.EntryInput, error) {
	_, avg, err := s.ledger.GlobalStockTx(run.ctx, run.tx, run.tenantID, item.VariantID)
	if err != nil {
		return nil, err
	}
	revenue := item.TotalAmount.Neg()
	total := item.Qty.Mul(avg).Round(4)
	entry := makeEntry(run.tenantID, run.doc.ID, item.ID, item.VariantID, run.doc.WarehouseID,
		mustDate(run.doc.Date), "IN", "RETURN_IN", item.Qty, avg, total, &revenue)
	return []ledger.EntryInput{entry}, nil
}
