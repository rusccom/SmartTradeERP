package reports

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/httpx"
)

func (r *Repository) FullStockRows(
	ctx context.Context,
	tenantID string,
	query httpx.ListQuery,
) ([]FullStockRow, int, error) {
	total, err := r.countFullStock(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	items, err := r.loadFullStock(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	return r.withFullStockWarehouses(ctx, tenantID, query, items, total)
}

func (r *Repository) countFullStock(ctx context.Context, tenantID string, list httpx.ListQuery) (int, error) {
	query := fullStockCountSQL()
	args := []any{tenantID}
	query, args = appendFullStockFilters(query, args, list)
	row := r.store.Pool.QueryRow(ctx, query, args...)
	total := 0
	err := row.Scan(&total)
	return total, err
}

func (r *Repository) loadFullStock(ctx context.Context, tenantID string, list httpx.ListQuery) ([]FullStockRow, error) {
	query := fullStockSelectSQL()
	args := []any{tenantID}
	query, args = appendFullStockFilters(query, args, list)
	query, args = appendFullStockPaging(query, args, list)
	rows, err := r.store.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFullStockRows(rows)
}

func fullStockCountSQL() string {
	return `SELECT COUNT(*)
        FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        WHERE v.tenant_id=$1 AND p.tenant_id=$1`
}

func fullStockSelectSQL() string {
	return `SELECT p.id::text, p.name, v.id::text, COALESCE(v.name,'Default'),
        COALESCE(v.sku_code,''), COALESCE(v.barcode,''), v.unit,
        COALESCE(latest.running_qty,0), COALESCE(latest.running_avg,0)
        FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        LEFT JOIN LATERAL (` + latestStockSQL() + `) latest ON true
        WHERE v.tenant_id=$1 AND p.tenant_id=$1`
}

func latestStockSQL() string {
	return `SELECT running_qty, running_avg
            FROM ledger.cost_ledger l
            WHERE l.tenant_id=$1 AND l.variant_id=v.id
            ORDER BY l.sequence_num DESC
            LIMIT 1`
}

func appendFullStockFilters(query string, args []any, list httpx.ListQuery) (string, []any) {
	query, args = appendTextFilter(query, args, "p.id::text", list.Filters["product_id"])
	query, args = appendTextFilter(query, args, "v.id::text", list.Filters["variant_id"])
	return appendFullStockSearch(query, args, list.Search)
}

func appendTextFilter(query string, args []any, field string, value string) (string, []any) {
	if value == "" {
		return query, args
	}
	query += ` AND ` + field + `=$` + argPosition(args)
	args = append(args, value)
	return query, args
}

func appendFullStockSearch(query string, args []any, search string) (string, []any) {
	if search == "" {
		return query, args
	}
	query += ` AND ` + fullStockSearchSQL(argPosition(args))
	args = append(args, search)
	return query, args
}

func fullStockSearchSQL(position string) string {
	return `(p.name ILIKE '%' || $` + position + ` || '%'
        OR COALESCE(v.name,'') ILIKE '%' || $` + position + ` || '%'
        OR COALESCE(v.sku_code,'') ILIKE '%' || $` + position + ` || '%'
        OR COALESCE(v.barcode,'') ILIKE '%' || $` + position + ` || '%')`
}

func appendFullStockPaging(query string, args []any, list httpx.ListQuery) (string, []any) {
	query += ` ORDER BY ` + fullStockOrder(list)
	query += ` LIMIT $` + argPosition(args)
	query += ` OFFSET $` + nextArgPosition(args)
	args = append(args, list.PerPage, httpx.Offset(list.Page, list.PerPage))
	return query, args
}

func fullStockOrder(list httpx.ListQuery) string {
	dir := list.SortDir
	switch list.SortBy {
	case "sku_code":
		return "COALESCE(v.sku_code,'') " + dir + ", p.name asc"
	case "global_qty":
		return "COALESCE(latest.running_qty,0) " + dir + ", p.name asc"
	case "avg":
		return "COALESCE(latest.running_avg,0) " + dir + ", p.name asc"
	}
	return "p.name " + dir + ", COALESCE(v.name,'Default') " + dir
}

func scanFullStockRows(rows pgx.Rows) ([]FullStockRow, error) {
	items := make([]FullStockRow, 0)
	for rows.Next() {
		item, err := scanFullStockRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanFullStockRow(rows pgx.Rows) (FullStockRow, error) {
	item := FullStockRow{}
	err := rows.Scan(&item.ProductID, &item.ProductName, &item.VariantID, &item.VariantName,
		&item.SKUCode, &item.Barcode, &item.Unit, &item.GlobalQty, &item.Avg)
	item.Name = fullStockName(item)
	item.StockValue = item.GlobalQty.Mul(item.Avg).Round(4)
	return item, err
}

func fullStockName(item FullStockRow) string {
	if item.VariantName == "" || item.VariantName == "Default" {
		return item.ProductName
	}
	return item.ProductName + " / " + item.VariantName
}

func (r *Repository) withFullStockWarehouses(
	ctx context.Context,
	tenantID string,
	list httpx.ListQuery,
	items []FullStockRow,
	total int,
) ([]FullStockRow, int, error) {
	stocks, err := r.fullStockWarehouses(ctx, tenantID, list, items)
	if err != nil {
		return nil, 0, err
	}
	return attachFullStockWarehouses(items, stocks), total, nil
}

func (r *Repository) fullStockWarehouses(
	ctx context.Context,
	tenantID string,
	list httpx.ListQuery,
	items []FullStockRow,
) (map[string][]FullStockWarehouse, error) {
	ids := fullStockVariantIDs(items)
	if len(ids) == 0 {
		return map[string][]FullStockWarehouse{}, nil
	}
	query, args := fullStockWarehousesQuery(tenantID, ids, list.Filters["warehouse_id"])
	rows, err := r.store.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFullStockWarehouses(rows)
}

func fullStockVariantIDs(items []FullStockRow) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.VariantID)
	}
	return ids
}

func fullStockWarehousesQuery(tenantID string, ids []string, warehouseID string) (string, []any) {
	query := fullStockWarehousesSQL()
	args := []any{tenantID}
	query, args = appendVariantList(query, args, ids)
	query, args = appendTextFilter(query, args, "w.id::text", warehouseID)
	query += ` GROUP BY v.id, w.id, w.name, w.created_at ORDER BY v.id, w.created_at`
	return query, args
}

func fullStockWarehousesSQL() string {
	return `SELECT v.id::text, w.id::text, w.name,
        COALESCE(SUM(CASE WHEN l.type='IN' THEN l.qty ELSE -l.qty END),0)
        FROM catalog.product_variants v
        JOIN catalog.warehouses w ON w.tenant_id=$1
        LEFT JOIN ledger.cost_ledger l ON l.tenant_id=$1
            AND l.variant_id=v.id AND l.warehouse_id=w.id
        WHERE v.tenant_id=$1`
}

func appendVariantList(query string, args []any, ids []string) (string, []any) {
	positions := make([]string, 0, len(ids))
	for _, id := range ids {
		positions = append(positions, "$"+argPosition(args))
		args = append(args, id)
	}
	query += ` AND v.id::text IN (` + strings.Join(positions, ",") + `)`
	return query, args
}

func scanFullStockWarehouses(rows pgx.Rows) (map[string][]FullStockWarehouse, error) {
	result := make(map[string][]FullStockWarehouse)
	for rows.Next() {
		variantID, item, err := scanFullStockWarehouse(rows)
		if err != nil {
			return nil, err
		}
		result[variantID] = append(result[variantID], item)
	}
	return result, rows.Err()
}

func scanFullStockWarehouse(rows pgx.Rows) (string, FullStockWarehouse, error) {
	variantID := ""
	item := FullStockWarehouse{}
	err := rows.Scan(&variantID, &item.WarehouseID, &item.Warehouse, &item.Qty)
	return variantID, item, err
}

func attachFullStockWarehouses(
	items []FullStockRow,
	stocks map[string][]FullStockWarehouse,
) []FullStockRow {
	for index := range items {
		items[index].Warehouses = stocks[items[index].VariantID]
		if items[index].Warehouses == nil {
			items[index].Warehouses = []FullStockWarehouse{}
		}
	}
	return items
}

func argPosition(args []any) string {
	return strconv.Itoa(len(args) + 1)
}

func nextArgPosition(args []any) string {
	return strconv.Itoa(len(args) + 2)
}
