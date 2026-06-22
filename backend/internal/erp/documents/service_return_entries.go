package documents

import (
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/erp/ledger"
)

func (s *Service) buildReturnEntry(
	run postingRun,
	item postingItem,
) ([]ledger.EntryInput, error) {
	revenue := item.TotalAmount.Neg()
	entry := makeEntry(run.tenantID, run.doc.ID, item.ID, item.VariantID, run.doc.WarehouseID,
		mustDate(run.doc.Date), "IN", "RETURN_IN", item.Qty, decimal.Zero, decimal.Zero, &revenue)
	return []ledger.EntryInput{entry}, nil
}
