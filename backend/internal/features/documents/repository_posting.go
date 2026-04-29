package documents

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) PostingDocument(ctx context.Context, tx pgx.Tx, tenantID, documentID string) (Document, error) {
	query := `SELECT id::text, type, date::text, COALESCE(number,''), status,
        COALESCE(warehouse_id::text,''), COALESCE(source_warehouse_id::text,''),
        COALESCE(target_warehouse_id::text,''), COALESCE(shift_id::text,''), COALESCE(note,'')
        FROM documents.documents
        WHERE tenant_id=$1 AND id=$2`
	row := tx.QueryRow(ctx, query, tenantID, documentID)
	doc := Document{}
	err := row.Scan(&doc.ID, &doc.Type, &doc.Date, &doc.Number, &doc.Status,
		&doc.WarehouseID, &doc.SourceWarehouseID, &doc.TargetWarehouseID, &doc.ShiftID, &doc.Note)
	return doc, err
}

func (r *Repository) PostingItems(ctx context.Context, tx pgx.Tx, tenantID, documentID string) ([]postingItem, error) {
	query := `SELECT i.id::text, i.variant_id::text, i.qty, i.unit_price, i.total_amount, p.is_composite
        FROM documents.document_items i
        JOIN documents.documents d ON d.id=i.document_id
        JOIN catalog.product_variants v ON v.id=i.variant_id
        JOIN catalog.products p ON p.id=v.product_id
        WHERE d.tenant_id=$1 AND v.tenant_id=$1 AND p.tenant_id=$1 AND d.id=$2`
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
