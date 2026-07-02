package reports

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
)

type Repository struct {
	store *db.Store
}

const topProductsSQL = `SELECT p.id::text, p.name, COALESCE(SUM(f.gross_profit),0) AS profit
        FROM ledger.document_item_financials f
        JOIN documents.document_items i ON i.id=f.document_item_id AND i.document_id=f.document_id
        JOIN documents.documents d ON d.id=f.document_id
        JOIN catalog.product_variants v ON v.id=i.variant_id
        JOIN catalog.products p ON p.id=v.product_id
        WHERE d.status='posted'
            AND d.type IN ('SALE','RETURN') AND d.date BETWEEN $1 AND $2
        GROUP BY p.id, p.name
        ORDER BY profit DESC
        LIMIT 20`

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) StockRows(ctx context.Context, tenantID, warehouseID string) ([]StockRow, error) {
	query, args := stockRowsQuery(warehouseID)
	rows, err := r.store.Query(ctx, tenantID, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStock(rows)
}

func stockRowsQuery(warehouseID string) (string, []any) {
	query := stockRowsSelect()
	args := []any{}
	query, args = appendWarehouseStockJoin(query, args, warehouseID)
	query += stockRowsGroup()
	return query, args
}

func stockRowsSelect() string {
	return `SELECT v.id::text, COALESCE(v.name, 'Default'),
        COALESCE(SUM(sb.qty),0),
        ` + stockAverageSQL() + `
        FROM catalog.product_variants v
        LEFT JOIN ledger.stock_balances sb ON sb.variant_id=v.id`
}

func stockAverageSQL() string {
	return `COALESCE((
            SELECT r.running_avg_cost
            FROM ledger.cost_movement_results r
            WHERE r.variant_id=v.id
            ORDER BY r.sequence_num DESC
            LIMIT 1
        ),0)`
}

func stockRowsGroup() string {
	return ` GROUP BY v.id, v.name
        ORDER BY COALESCE(v.name, 'Default')`
}

func appendWarehouseStockJoin(query string, args []any, warehouseID string) (string, []any) {
	if warehouseID == "" {
		return query, args
	}
	query += ` AND sb.warehouse_id=$` + strconv.Itoa(len(args)+1)
	args = append(args, warehouseID)
	return query, args
}

func scanStock(rows pgx.Rows) ([]StockRow, error) {
	items := make([]StockRow, 0)
	for rows.Next() {
		item := StockRow{}
		if err := rows.Scan(&item.VariantID, &item.Name, &item.Qty, &item.Avg); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) TopProducts(ctx context.Context, tenantID string, fromDate, toDate time.Time) ([]TopProduct, error) {
	rows, err := r.store.Query(ctx, tenantID, topProductsSQL, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTopProducts(rows)
}

func scanTopProducts(rows pgx.Rows) ([]TopProduct, error) {
	items := make([]TopProduct, 0)
	for rows.Next() {
		item := TopProduct{}
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Profit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
