package documents

import (
    "context"
    "time"

    "smarterp/backend/internal/features/ledger"
)

func (s *Service) revenueShares(
    ctx context.Context,
    tenantID string,
    item postingItem,
    components []variantComponent,
) (map[string]float64, error) {
    costs := make(map[string]float64)
    totalCost := 0.0
    for _, component := range components {
        cost, err := s.componentCost(ctx, tenantID, item, component)
        if err != nil {
            return nil, err
        }
        costs[component.ComponentVariantID] = cost
        totalCost += cost
    }
    return buildShares(costs, totalCost), nil
}

func (s *Service) componentCost(
    ctx context.Context,
    tenantID string,
    item postingItem,
    component variantComponent,
) (float64, error) {
    qty := item.Qty * component.QtyPerUnit
    _, avg, err := s.ledger.GlobalStock(ctx, tenantID, component.ComponentVariantID)
    if err != nil {
        return 0, err
    }
    return qty * avg, nil
}

func buildShares(costs map[string]float64, totalCost float64) map[string]float64 {
    shares := make(map[string]float64)
    for variantID, cost := range costs {
        shares[variantID] = normalizedShare(cost, totalCost)
    }
    return shares
}

func normalizedShare(cost, total float64) float64 {
    if total <= 0 {
        return 0
    }
    return cost / total
}

func (s *Service) buildCompositeSaleEntry(
    doc Document,
    tenantID string,
    item postingItem,
    component variantComponent,
    shares map[string]float64,
) ledger.EntryInput {
    qty := item.Qty * component.QtyPerUnit
    share := shares[component.ComponentVariantID]
    revenue := item.TotalAmount * share
    entryRevenue := revenue
    return makeEntry(tenantID, doc.ID, item.ID, component.ComponentVariantID, doc.WarehouseID,
        mustDate(doc.Date), "OUT", "SALE", qty, item.UnitPrice, qty*item.UnitPrice, &entryRevenue)
}

func (s *Service) buildDefaultEntry(doc Document, tenantID string, item postingItem) ledger.EntryInput {
    metaType, reason := docMeta(doc.Type)
    revenue := revenueForType(doc.Type, item.TotalAmount)
    return makeEntry(tenantID, doc.ID, item.ID, item.VariantID, doc.WarehouseID,
        mustDate(doc.Date), metaType, reason, item.Qty, item.UnitPrice, item.TotalAmount, revenue)
}

func makeEntry(
    tenantID string,
    documentID string,
    documentItemID string,
    variantID string,
    warehouseID string,
    date time.Time,
    itemType string,
    reason string,
    qty float64,
    unitPrice float64,
    totalAmount float64,
    revenue *float64,
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
    case "WRITEOFF":
        return "OUT", "WRITEOFF"
    }
    return "OUT", "SALE"
}

func revenueForType(documentType string, amount float64) *float64 {
    if documentType != "SALE" {
        return nil
    }
    value := amount
    return &value
}

func mustDate(raw string) time.Time {
    value, err := time.Parse("2006-01-02", raw)
    if err != nil {
        return time.Now().UTC()
    }
    return value
}
