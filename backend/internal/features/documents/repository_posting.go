package documents

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

func (r *Repository) PostingDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string) (Document, error) {
	query := `SELECT id::text, type, date::text, COALESCE(number,''), status,
        COALESCE(warehouse_id::text,''), COALESCE(source_warehouse_id::text,''),
        COALESCE(target_warehouse_id::text,''), COALESCE(note,'')
        FROM documents.documents
        WHERE tenant_id=$1 AND id=$2`
	row := tx.QueryRow(ctx, query, tenantID, documentID)
	doc := Document{}
	err := row.Scan(&doc.ID, &doc.Type, &doc.Date, &doc.Number, &doc.Status,
		&doc.WarehouseID, &doc.SourceWarehouseID, &doc.TargetWarehouseID, &doc.Note)
	return doc, err
}

func (r *Repository) PostingItems(ctx context.Context, tx pgx.Tx, tenantID, documentID string) ([]postingItem, error) {
	query := `SELECT i.id::text, i.variant_id::text, i.qty, i.unit_price, i.total_amount, p.is_composite
        FROM documents.document_items i
        JOIN documents.documents d ON d.id=i.document_id
        JOIN catalog.product_variants v ON v.id=i.variant_id
        JOIN catalog.products p ON p.id=v.product_id
        WHERE d.tenant_id=$1 AND d.id=$2`
	rows, err := tx.Query(ctx, query, tenantID, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostingItems(rows)
}

func scanPostingItems(rows pgx.Rows) ([]postingItem, error) {
	items := make([]postingItem, 0)
	for rows.Next() {
		item := postingItem{}
		err := rows.Scan(&item.ID, &item.VariantID, &item.Qty, &item.UnitPrice, &item.TotalAmount, &item.IsComposite)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) VariantComponents(ctx context.Context, tx pgx.Tx, variantID string) ([]variantComponent, error) {
	query := `SELECT component_variant_id::text, qty
        FROM catalog.variant_components
        WHERE variant_id=$1`
	rows, err := tx.Query(ctx, query, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVariantComponents(rows)
}

func scanVariantComponents(rows pgx.Rows) ([]variantComponent, error) {
	items := make([]variantComponent, 0)
	for rows.Next() {
		item := variantComponent{}
		if err := rows.Scan(&item.ComponentVariantID, &item.QtyPerUnit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) SaveItemComponents(
	ctx context.Context,
	tx pgx.Tx,
	documentItemID string,
	components []variantComponent,
	documentQty decimal.Decimal,
) error {
	for _, item := range components {
		qtyTotal := item.QtyPerUnit.Mul(documentQty)
		if err := r.insertItemComponent(ctx, tx, documentItemID, item, qtyTotal); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) insertItemComponent(
	ctx context.Context,
	tx pgx.Tx,
	documentItemID string,
	item variantComponent,
	qtyTotal decimal.Decimal,
) error {
	query := `INSERT INTO documents.document_item_components
        (id, document_item_id, component_variant_id, qty_per_unit, qty_total)
        VALUES ($1,$2,$3,$4,$5)`
	_, err := tx.Exec(ctx, query, uuid.NewString(), documentItemID, item.ComponentVariantID, item.QtyPerUnit, qtyTotal)
	return err
}

func (r *Repository) DeleteItemComponentsByDocument(ctx context.Context, tx pgx.Tx, documentID string) error {
	query := `DELETE FROM documents.document_item_components dic
        USING documents.document_items di
        WHERE dic.document_item_id=di.id AND di.document_id=$1`
	_, err := tx.Exec(ctx, query, documentID)
	return err
}
