package documents

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/features/ledger"
)

func (s *Service) postDocumentTx(ctx context.Context, tx pgx.Tx, tenantID, documentID string) error {
	doc, err := s.repo.PostingDocument(ctx, tx, tenantID, documentID)
	if err != nil {
		return err
	}
	items, err := s.repo.PostingItems(ctx, tx, tenantID, documentID)
	if err != nil {
		return err
	}
	if err := s.postItems(ctx, tx, tenantID, doc, items); err != nil {
		return err
	}
	return s.repo.SetStatus(ctx, tx, tenantID, documentID, "posted")
}

func (s *Service) postItems(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	doc Document,
	items []postingItem,
) error {
	for _, item := range items {
		if err := s.postItem(ctx, tx, tenantID, doc, item); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) postItem(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	doc Document,
	item postingItem,
) error {
	if item.IsComposite {
		return s.postCompositeItem(ctx, tx, tenantID, doc, item)
	}
	return s.postSimpleItem(ctx, tx, tenantID, doc, item)
}

func (s *Service) postSimpleItem(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	doc Document,
	item postingItem,
) error {
	entries, err := s.buildEntriesForSimple(ctx, tenantID, doc, item)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}
	return s.appendEntries(ctx, tx, entries)
}

func (s *Service) postCompositeItem(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	doc Document,
	item postingItem,
) error {
	components, err := s.loadCompositeComponents(ctx, tx, tenantID, item)
	if err != nil {
		return err
	}
	if err := s.repo.SaveItemComponents(ctx, tx, item.ID, components, item.Qty); err != nil {
		return err
	}
	entries, err := s.buildEntriesForComposite(ctx, tenantID, doc, item, components)
	if err != nil {
		return err
	}
	return s.appendEntries(ctx, tx, entries)
}

func (s *Service) loadCompositeComponents(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	item postingItem,
) ([]variantComponent, error) {
	components, err := s.repo.VariantComponents(ctx, tx, tenantID, item.VariantID)
	if err != nil {
		return nil, err
	}
	if len(components) == 0 {
		return nil, ErrCompositeWithoutComponents
	}
	return components, nil
}

func (s *Service) appendEntries(ctx context.Context, tx pgx.Tx, entries []ledger.EntryInput) error {
	for _, entry := range entries {
		if _, err := s.ledger.Append(ctx, tx, entry); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) buildEntriesForSimple(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
) ([]ledger.EntryInput, error) {
	switch doc.Type {
	case "TRANSFER":
		return s.buildTransferEntries(ctx, tenantID, doc, item)
	case "INVENTORY":
		return s.buildInventoryEntry(ctx, tenantID, doc, item)
	case "RETURN":
		return s.buildReturnEntry(ctx, tenantID, doc, item)
	default:
		return []ledger.EntryInput{s.buildDefaultEntry(doc, tenantID, item)}, nil
	}
}

func (s *Service) buildTransferEntries(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
) ([]ledger.EntryInput, error) {
	_, avg, err := s.ledger.GlobalStock(ctx, tenantID, item.VariantID)
	if err != nil {
		return nil, err
	}
	date := mustDate(doc.Date)
	total := item.Qty.Mul(avg).Round(4)
	outEntry := makeEntry(tenantID, doc.ID, item.ID, item.VariantID, doc.SourceWarehouseID,
		date, "OUT", "TRANSFER_OUT", item.Qty, avg, total, nil)
	inEntry := makeEntry(tenantID, doc.ID, item.ID, item.VariantID, doc.TargetWarehouseID,
		date, "IN", "TRANSFER_IN", item.Qty, avg, total, nil)
	return []ledger.EntryInput{outEntry, inEntry}, nil
}

func (s *Service) buildInventoryEntry(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
) ([]ledger.EntryInput, error) {
	diff, err := s.inventoryDiff(ctx, tenantID, doc, item)
	if err != nil {
		return nil, err
	}
	if diff.IsZero() {
		return nil, nil
	}
	_, avg, err := s.ledger.GlobalStock(ctx, tenantID, item.VariantID)
	if err != nil {
		return nil, err
	}
	entry := inventoryEntry(doc, tenantID, item, diff, avg)
	return []ledger.EntryInput{entry}, nil
}

func (s *Service) inventoryDiff(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
) (decimal.Decimal, error) {
	qty, err := s.ledger.WarehouseStock(ctx, tenantID, item.VariantID, doc.WarehouseID)
	if err != nil {
		return decimal.Zero, err
	}
	return item.Qty.Sub(qty), nil
}

func inventoryEntry(doc Document, tenantID string, item postingItem, diff, avg decimal.Decimal) ledger.EntryInput {
	entryType, reason, absQty := inventoryMeta(diff)
	total := absQty.Mul(avg).Round(4)
	return makeEntry(tenantID, doc.ID, item.ID, item.VariantID, doc.WarehouseID,
		mustDate(doc.Date), entryType, reason, absQty, avg, total, nil)
}

func inventoryMeta(diff decimal.Decimal) (string, string, decimal.Decimal) {
	entryType := "IN"
	reason := "SURPLUS"
	if diff.IsPositive() {
		return entryType, reason, diff
	}
	return "OUT", "SHORTAGE", diff.Neg()
}

func (s *Service) buildEntriesForComposite(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	if doc.Type == "SALE" {
		return s.buildCompositeSaleEntries(ctx, tenantID, doc, item, components)
	}
	if doc.Type == "RETURN" {
		return s.buildCompositeReturnEntries(ctx, tenantID, doc, item, components)
	}
	return s.buildCompositeRegularEntries(ctx, tenantID, doc, item, components)
}

func (s *Service) buildCompositeRegularEntries(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	entries := make([]ledger.EntryInput, 0, len(components))
	for _, component := range components {
		derived := deriveComponentItem(item, component)
		rows, err := s.buildEntriesForSimple(ctx, tenantID, doc, derived)
		if err != nil {
			return nil, err
		}
		entries = append(entries, rows...)
	}
	return entries, nil
}

func deriveComponentItem(item postingItem, component variantComponent) postingItem {
	qty := item.Qty.Mul(component.QtyPerUnit)
	return postingItem{
		ID:          item.ID,
		VariantID:   component.ComponentVariantID,
		Qty:         qty,
		UnitPrice:   item.UnitPrice,
		TotalAmount: qty.Mul(item.UnitPrice),
		IsComposite: false,
	}
}

func (s *Service) buildCompositeSaleEntries(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	shares, err := s.revenueShares(ctx, tenantID, item, components)
	if err != nil {
		return nil, err
	}
	entries := make([]ledger.EntryInput, 0, len(components))
	for _, component := range components {
		entry := s.buildCompositeSaleEntry(doc, tenantID, item, component, shares)
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Service) buildCompositeReturnEntries(
	ctx context.Context,
	tenantID string,
	doc Document,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	shares, err := s.revenueShares(ctx, tenantID, item, components)
	if err != nil {
		return nil, err
	}
	input := compositeReturnInput{ctx: ctx, tenantID: tenantID, doc: doc, item: item, shares: shares}
	return s.buildCompositeReturnRows(input, components)
}

func (s *Service) buildCompositeReturnRows(
	input compositeReturnInput,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	entries := make([]ledger.EntryInput, 0, len(components))
	for _, component := range components {
		entry, err := s.buildCompositeReturnEntry(input, component)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
