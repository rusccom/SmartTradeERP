package documents

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/features/ledger"
)

func (s *Service) postItems(
	run postingRun,
	items []postingItem,
) ([]ledger.VariantSequence, error) {
	affected := make([]ledger.VariantSequence, 0)
	for _, item := range items {
		itemAffected, err := s.postItem(run, item)
		if err != nil {
			return nil, err
		}
		affected = ledger.MergeAffected(affected, itemAffected)
	}
	return affected, nil
}

func (s *Service) postItem(
	run postingRun,
	item postingItem,
) ([]ledger.VariantSequence, error) {
	if item.IsComposite {
		return s.postCompositeItem(run, item)
	}
	return s.postSimpleItem(run, item)
}

func (s *Service) postSimpleItem(
	run postingRun,
	item postingItem,
) ([]ledger.VariantSequence, error) {
	entries, err := s.buildEntriesForSimple(run, item)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	return s.appendEntries(run, entries)
}

func (s *Service) postCompositeItem(
	run postingRun,
	item postingItem,
) ([]ledger.VariantSequence, error) {
	components, err := s.loadCompositeComponents(run.ctx, run.tx, run.tenantID, item)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveItemComponents(run.ctx, run.tx, item.ID, components, item.Qty); err != nil {
		return nil, err
	}
	entries, err := s.buildEntriesForComposite(run, item, components)
	if err != nil {
		return nil, err
	}
	return s.appendEntries(run, entries)
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

func (s *Service) buildEntriesForSimple(
	run postingRun,
	item postingItem,
) ([]ledger.EntryInput, error) {
	switch run.doc.Type {
	case "TRANSFER":
		return s.buildTransferEntries(run, item)
	case "INVENTORY":
		return s.buildInventoryEntry(run, item)
	case "RETURN":
		return s.buildReturnEntry(run, item)
	default:
		return []ledger.EntryInput{s.buildDefaultEntry(run.doc, run.tenantID, item)}, nil
	}
}

func (s *Service) buildTransferEntries(
	run postingRun,
	item postingItem,
) ([]ledger.EntryInput, error) {
	_, avg, err := s.ledger.GlobalStockTx(run.ctx, run.tx, run.tenantID, item.VariantID)
	if err != nil {
		return nil, err
	}
	date := mustDate(run.doc.Date)
	total := item.Qty.Mul(avg).Round(4)
	outEntry := makeEntry(run.tenantID, run.doc.ID, item.ID, item.VariantID, run.doc.SourceWarehouseID,
		date, "OUT", "TRANSFER_OUT", item.Qty, avg, total, nil)
	inEntry := makeEntry(run.tenantID, run.doc.ID, item.ID, item.VariantID, run.doc.TargetWarehouseID,
		date, "IN", "TRANSFER_IN", item.Qty, avg, total, nil)
	return []ledger.EntryInput{outEntry, inEntry}, nil
}

func (s *Service) buildInventoryEntry(
	run postingRun,
	item postingItem,
) ([]ledger.EntryInput, error) {
	diff, err := s.inventoryDiff(run, item)
	if err != nil {
		return nil, err
	}
	if diff.IsZero() {
		return nil, nil
	}
	_, avg, err := s.ledger.GlobalStockTx(run.ctx, run.tx, run.tenantID, item.VariantID)
	if err != nil {
		return nil, err
	}
	entry := inventoryEntry(run.doc, run.tenantID, item, diff, avg)
	return []ledger.EntryInput{entry}, nil
}

func (s *Service) inventoryDiff(
	run postingRun,
	item postingItem,
) (decimal.Decimal, error) {
	qty, err := s.ledger.WarehouseStockTx(run.ctx, run.tx, run.tenantID, item.VariantID, run.doc.WarehouseID)
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
	run postingRun,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	if run.doc.Type == "SALE" {
		return s.buildCompositeSaleEntries(run, item, components)
	}
	if run.doc.Type == "RETURN" {
		return s.buildCompositeReturnEntries(run, item, components)
	}
	return s.buildCompositeRegularEntries(run, item, components)
}

func (s *Service) buildCompositeRegularEntries(
	run postingRun,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	entries := make([]ledger.EntryInput, 0, len(components))
	for _, component := range components {
		derived := deriveComponentItem(item, component)
		rows, err := s.buildEntriesForSimple(run, derived)
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
	run postingRun,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	shares, err := s.revenueShares(run, item, components)
	if err != nil {
		return nil, err
	}
	entries := make([]ledger.EntryInput, 0, len(components))
	for _, component := range components {
		entry := s.buildCompositeSaleEntry(run, item, component, shares)
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Service) buildCompositeReturnEntries(
	run postingRun,
	item postingItem,
	components []variantComponent,
) ([]ledger.EntryInput, error) {
	shares, err := s.revenueShares(run, item, components)
	if err != nil {
		return nil, err
	}
	input := compositeReturnInput{run: run, item: item, shares: shares}
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
