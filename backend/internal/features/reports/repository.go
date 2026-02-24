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
    query := `SELECT variant_id::text, MAX(name),
        COALESCE(MAX(running_qty),0), COALESCE(MAX(running_avg),0)
        FROM (
            SELECT l.variant_id, COALESCE(v.name, 'Default') AS name, l.running_qty, l.running_avg,
                ROW_NUMBER() OVER (PARTITION BY l.variant_id ORDER BY l.sequence_num DESC) AS rn
            FROM ledger.cost_ledger l
            JOIN catalog.product_variants v ON v.id=l.variant_id
            JOIN catalog.products p ON p.id=v.product_id
            WHERE l.tenant_id=$1 AND p.tenant_id=$1`
    args := []any{tenantID}
    query, args = appendWarehouse(query, args, warehouseID)
    query += `
        ) x WHERE rn=1
        GROUP BY variant_id
        ORDER BY MAX(name)`
    rows, err := r.store.Pool.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanStock(rows)
}

func appendWarehouse(query string, args []any, warehouseID string) (string, []any) {
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
        WHERE l.tenant_id=$1 AND l.date BETWEEN $2 AND $3 AND l.type='OUT'
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
