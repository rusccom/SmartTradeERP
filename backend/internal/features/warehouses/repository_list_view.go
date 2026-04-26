package warehouses

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type warehouseStockRow struct {
	WarehouseID string
	Item        WarehouseStockItem
}

func (r *Repository) ListWithIncludes(
	ctx context.Context,
	tenantID string,
	include WarehouseListInclude,
) ([]WarehouseListItem, error) {
	items, err := r.loadListItems(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if !include.Stock {
		return items, nil
	}
	return r.attachWarehouseStock(ctx, tenantID, items)
}

func (r *Repository) loadListItems(ctx context.Context, tenantID string) ([]WarehouseListItem, error) {
	query := `SELECT id::text, name, COALESCE(address,''), is_default, is_active, created_at::text
        FROM catalog.warehouses
        WHERE tenant_id=$1
        ORDER BY created_at`
	rows, err := r.store.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWarehouseListItems(rows)
}

func scanWarehouseListItems(rows pgx.Rows) ([]WarehouseListItem, error) {
	items := make([]WarehouseListItem, 0)
	for rows.Next() {
		item := WarehouseListItem{Stock: []WarehouseStockItem{}}
		err := rows.Scan(&item.ID, &item.Name, &item.Address, &item.IsDefault, &item.IsActive, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) attachWarehouseStock(
	ctx context.Context,
	tenantID string,
	items []WarehouseListItem,
) ([]WarehouseListItem, error) {
	stocks, err := r.loadWarehouseStockRows(ctx, tenantID, items)
	if err != nil {
		return nil, err
	}
	return attachWarehouseStockRows(items, stocks), nil
}

func (r *Repository) loadWarehouseStockRows(
	ctx context.Context,
	tenantID string,
	items []WarehouseListItem,
) ([]warehouseStockRow, error) {
	if len(items) == 0 {
		return []warehouseStockRow{}, nil
	}
	query, args := warehouseStockRowsQuery(tenantID, items)
	rows, err := r.store.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWarehouseStockRows(rows)
}

func warehouseStockRowsQuery(tenantID string, items []WarehouseListItem) (string, []any) {
	query := warehouseStockRowsSQL()
	args := []any{tenantID}
	query, args = appendWarehouseIDList(query, args, items)
	query += warehouseStockRowsGroup()
	return query, args
}

func warehouseStockRowsSQL() string {
	return `SELECT w.id::text, p.id::text, p.name, v.id::text, COALESCE(v.name,'Default'),
        COALESCE(v.sku_code,''), COALESCE(v.barcode,''), v.unit,
        COALESCE(SUM(CASE WHEN l.type='IN' THEN l.qty WHEN l.type='OUT' THEN -l.qty ELSE 0 END),0),
        COALESCE(latest.running_avg,0)
        FROM catalog.warehouses w
        JOIN catalog.product_variants v ON v.tenant_id=$1
        JOIN catalog.products p ON p.id=v.product_id AND p.tenant_id=$1
        LEFT JOIN ledger.cost_ledger l ON l.tenant_id=$1
            AND l.warehouse_id=w.id AND l.variant_id=v.id
        LEFT JOIN LATERAL (` + latestWarehouseStockSQL() + `) latest ON true
        WHERE w.tenant_id=$1`
}

func latestWarehouseStockSQL() string {
	return `SELECT running_avg
            FROM ledger.cost_ledger la
            WHERE la.tenant_id=$1 AND la.variant_id=v.id
            ORDER BY la.sequence_num DESC
            LIMIT 1`
}

func appendWarehouseIDList(
	query string,
	args []any,
	items []WarehouseListItem,
) (string, []any) {
	positions := make([]string, 0, len(items))
	for _, item := range items {
		positions = append(positions, "$"+warehouseArgPosition(args))
		args = append(args, item.ID)
	}
	query += ` AND w.id::text IN (` + strings.Join(positions, ",") + `)`
	return query, args
}

func warehouseStockRowsGroup() string {
	return ` GROUP BY w.id, p.id, p.name, v.id, v.name, v.sku_code, v.barcode,
        v.unit, latest.running_avg
        HAVING COALESCE(SUM(CASE WHEN l.type='IN' THEN l.qty WHEN l.type='OUT' THEN -l.qty ELSE 0 END),0) <> 0
        ORDER BY w.id, p.name, COALESCE(v.name,'Default')`
}

func scanWarehouseStockRows(rows pgx.Rows) ([]warehouseStockRow, error) {
	items := make([]warehouseStockRow, 0)
	for rows.Next() {
		item, err := scanWarehouseStockRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanWarehouseStockRow(rows pgx.Rows) (warehouseStockRow, error) {
	row := warehouseStockRow{Item: WarehouseStockItem{}}
	err := rows.Scan(&row.WarehouseID, &row.Item.ProductID, &row.Item.ProductName,
		&row.Item.VariantID, &row.Item.VariantName, &row.Item.SKUCode, &row.Item.Barcode,
		&row.Item.Unit, &row.Item.Qty, &row.Item.AvgCost)
	row.Item.Name = warehouseStockName(row.Item)
	row.Item.StockValue = row.Item.Qty.Mul(row.Item.AvgCost).Round(4)
	return row, err
}

func warehouseStockName(item WarehouseStockItem) string {
	if item.VariantName == "" || item.VariantName == "Default" {
		return item.ProductName
	}
	return item.ProductName + " / " + item.VariantName
}

func attachWarehouseStockRows(
	items []WarehouseListItem,
	stocks []warehouseStockRow,
) []WarehouseListItem {
	index := warehouseIndex(items)
	for _, row := range stocks {
		position, ok := index[row.WarehouseID]
		if ok {
			items[position] = appendWarehouseStock(items[position], row.Item)
		}
	}
	return items
}

func warehouseIndex(items []WarehouseListItem) map[string]int {
	result := make(map[string]int, len(items))
	for index, item := range items {
		result[item.ID] = index
	}
	return result
}

func appendWarehouseStock(item WarehouseListItem, stock WarehouseStockItem) WarehouseListItem {
	item.StockValue = item.StockValue.Add(stock.StockValue).Round(4)
	item.Stock = append(item.Stock, stock)
	return item
}

func warehouseArgPosition(args []any) string {
	return strconv.Itoa(len(args) + 1)
}
