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

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) StockRows(ctx context.Context, tenantID, warehouseID string) ([]StockRow, error) {
	query, args := stockRowsQuery(tenantID, warehouseID)
	rows, err := r.store.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStock(rows)
}

func stockRowsQuery(tenantID, warehouseID string) (string, []any) {
	query := stockRowsSelect()
	args := []any{tenantID}
	query, args = appendWarehouseStockJoin(query, args, warehouseID)
	query += stockRowsGroup()
	return query, args
}

func stockRowsSelect() string {
	return `SELECT v.id::text, COALESCE(v.name, 'Default'),
        COALESCE(SUM(CASE WHEN l.type='IN' THEN l.qty ELSE -l.qty END),0),
        ` + stockAverageSQL() + `
        FROM catalog.product_variants v
        LEFT JOIN ledger.cost_ledger l ON l.tenant_id=$1 AND l.variant_id=v.id`
}

func stockAverageSQL() string {
	return `COALESCE((
            SELECT la.running_avg
            FROM ledger.cost_ledger la
            WHERE la.tenant_id=$1 AND la.variant_id=v.id
            ORDER BY la.sequence_num DESC
            LIMIT 1
        ),0)`
}

func stockRowsGroup() string {
	return ` WHERE v.tenant_id=$1
        GROUP BY v.id, v.name
        ORDER BY COALESCE(v.name, 'Default')`
}

func appendWarehouseStockJoin(query string, args []any, warehouseID string) (string, []any) {
	if warehouseID == "" {
		return query, args
	}
	query += ` AND l.warehouse_id=$` + strconv.Itoa(len(args)+1)
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
	query := `SELECT p.id::text, p.name, COALESCE(SUM(l.profit),0) AS profit
        FROM ledger.cost_ledger l
        JOIN catalog.product_variants v ON v.id=l.variant_id
        JOIN catalog.products p ON p.id=v.product_id
        WHERE l.tenant_id=$1 AND v.tenant_id=$1 AND l.date BETWEEN $2 AND $3 AND l.type='OUT'
        GROUP BY p.id, p.name
        ORDER BY profit DESC
        LIMIT 20`
	rows, err := r.store.Pool.Query(ctx, query, tenantID, fromDate, toDate)
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
