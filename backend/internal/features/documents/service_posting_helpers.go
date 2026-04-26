package documents

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/features/ledger"
)

type compositeReturnInput struct {
	ctx      context.Context
	tenantID string
	doc      Document
	item     postingItem
	shares   map[string]decimal.Decimal
}

func (s *Service) revenueShares(
	ctx context.Context,
	tenantID string,
	item postingItem,
	components []variantComponent,
) (map[string]decimal.Decimal, error) {
	costs := make(map[string]decimal.Decimal)
	totalCost := decimal.Zero
	for _, component := range components {
		cost, err := s.componentCost(ctx, tenantID, item, component)
		if err != nil {
			return nil, err
		}
		costs[component.ComponentVariantID] = cost
		totalCost = totalCost.Add(cost)
	}
	return buildShares(costs, totalCost), nil
}

func (s *Service) componentCost(
	ctx context.Context,
	tenantID string,
	item postingItem,
	component variantComponent,
) (decimal.Decimal, error) {
	qty := item.Qty.Mul(component.QtyPerUnit)
	_, avg, err := s.ledger.GlobalStock(ctx, tenantID, component.ComponentVariantID)
	if err != nil {
		return decimal.Zero, err
	}
	return qty.Mul(avg), nil
}

func buildShares(costs map[string]decimal.Decimal, totalCost decimal.Decimal) map[string]decimal.Decimal {
	shares := make(map[string]decimal.Decimal)
	if totalCost.LessThanOrEqual(decimal.Zero) {
		return equalShares(costs)
	}
	for variantID, cost := range costs {
		shares[variantID] = normalizedShare(cost, totalCost)
	}
	return shares
}

func equalShares(costs map[string]decimal.Decimal) map[string]decimal.Decimal {
	shares := make(map[string]decimal.Decimal)
	if len(costs) == 0 {
		return shares
	}
	share := decimal.NewFromInt(1).Div(decimal.NewFromInt(int64(len(costs)))).Round(8)
	for variantID := range costs {
		shares[variantID] = share
	}
	return shares
}

func normalizedShare(cost, total decimal.Decimal) decimal.Decimal {
	if total.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}
	return cost.Div(total).Round(8)
}

func (s *Service) buildCompositeSaleEntry(
	doc Document,
	tenantID string,
	item postingItem,
	component variantComponent,
	shares map[string]decimal.Decimal,
) ledger.EntryInput {
	qty := item.Qty.Mul(component.QtyPerUnit)
	share := shares[component.ComponentVariantID]
	revenue := item.TotalAmount.Mul(share).Round(4)
	return makeEntry(tenantID, doc.ID, item.ID, component.ComponentVariantID, doc.WarehouseID,
		mustDate(doc.Date), "OUT", "SALE", qty, item.UnitPrice, qty.Mul(item.UnitPrice), &revenue)
}

func (s *Service) buildCompositeReturnEntry(input compositeReturnInput, component variantComponent) (ledger.EntryInput, error) {
	qty := input.item.Qty.Mul(component.QtyPerUnit)
	_, avg, err := s.ledger.GlobalStock(input.ctx, input.tenantID, component.ComponentVariantID)
	if err != nil {
		return ledger.EntryInput{}, err
	}
	share := input.shares[component.ComponentVariantID]
	revenue := input.item.TotalAmount.Mul(share).Round(4).Neg()
	total := qty.Mul(avg).Round(4)
	entry := makeEntry(input.tenantID, input.doc.ID, input.item.ID, component.ComponentVariantID,
		input.doc.WarehouseID, mustDate(input.doc.Date), "IN", "RETURN_IN", qty, avg, total, &revenue)
	return entry, nil
}

func (s *Service) buildDefaultEntry(doc Document, tenantID string, item postingItem) ledger.EntryInput {
	metaType, reason := docMeta(doc.Type)
	revenue := revenueForType(doc.Type, item.TotalAmount)
	return makeEntry(tenantID, doc.ID, item.ID, item.VariantID, doc.WarehouseID,
		mustDate(doc.Date), metaType, reason, item.Qty, item.UnitPrice, item.TotalAmount, revenue)
}

func makeEntry(
	tenantID, documentID, documentItemID, variantID, warehouseID string,
	date time.Time,
	itemType, reason string,
	qty, unitPrice, totalAmount decimal.Decimal,
	revenue *decimal.Decimal,
) ledger.EntryInput {
	return ledger.EntryInput{
		TenantID: tenantID, VariantID: variantID, DocumentID: documentID,
		DocumentItemID: documentItemID, WarehouseID: warehouseID, Date: date,
		Type: itemType, Reason: reason, Qty: qty, UnitPrice: unitPrice,
		TotalAmount: totalAmount, Revenue: revenue,
	}
}

func docMeta(documentType string) (string, string) {
	switch documentType {
	case "RECEIPT":
		return "IN", "PURCHASE"
	case "SALE":
		return "OUT", "SALE"
	case "RETURN":
		return "IN", "RETURN_IN"
	case "WRITEOFF":
		return "OUT", "WRITEOFF"
	}
	return "OUT", "SALE"
}

func revenueForType(documentType string, amount decimal.Decimal) *decimal.Decimal {
	switch documentType {
	case "SALE":
		value := amount
		return &value
	case "RETURN":
		value := amount.Neg()
		return &value
	default:
		return nil
	}
}

func mustDate(raw string) time.Time {
	value, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Now().UTC()
	}
	return value
}
