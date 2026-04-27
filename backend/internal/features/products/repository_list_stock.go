package products

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

type variantWarehouseRow struct {
	VariantID string
	Item      ProductWarehouseItem
}

func (r *Repository) attachVariantWarehouses(
	ctx context.Context,
	tenantID string,
	variants []productVariantRow,
	filter ProductStockFilter,
) ([]productVariantRow, error) {
	stocks, err := r.loadVariantWarehouses(ctx, tenantID, variants, filter)
	if err != nil {
		return nil, err
	}
	for index := range variants {
		variants[index].Item.Warehouses = stocks[variants[index].Item.ID]
		if variants[index].Item.Warehouses == nil {
			variants[index].Item.Warehouses = []ProductWarehouseItem{}
		}
	}
	return variants, nil
}

func (r *Repository) loadVariantWarehouses(
	ctx context.Context,
	tenantID string,
	variants []productVariantRow,
	filter ProductStockFilter,
) (map[string][]ProductWarehouseItem, error) {
	if len(variants) == 0 {
		return map[string][]ProductWarehouseItem{}, nil
	}
	query, args := variantWarehousesQuery(tenantID, variants, filter)
	rows, err := r.store.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVariantWarehouses(rows)
}

func variantWarehousesQuery(
	tenantID string,
	variants []productVariantRow,
	filter ProductStockFilter,
) (string, []any) {
	query := variantWarehousesSQL()
	args := []any{tenantID}
	query, args = appendVariantIDList(query, args, variants)
	query, args = appendWarehouseScope(query, args, filter)
	query += ` GROUP BY v.id, w.id, w.name, w.created_at, sb.qty ORDER BY v.id, w.created_at`
	return query, args
}

func variantWarehousesSQL() string {
	return `SELECT v.id::text, w.id::text, w.name,
        COALESCE(sb.qty,0)
        FROM catalog.product_variants v
        JOIN catalog.warehouses w ON w.tenant_id=$1
        LEFT JOIN ledger.stock_balances sb ON sb.tenant_id=$1
            AND sb.variant_id=v.id AND sb.warehouse_id=w.id
        WHERE v.tenant_id=$1`
}

func appendVariantIDList(
	query string,
	args []any,
	variants []productVariantRow,
) (string, []any) {
	positions := make([]string, 0, len(variants))
	for _, item := range variants {
		positions = append(positions, "$"+position(len(args)+1))
		args = append(args, item.Item.ID)
	}
	query += ` AND v.id::text IN (` + strings.Join(positions, ",") + `)`
	return query, args
}

func scanVariantWarehouses(rows pgx.Rows) (map[string][]ProductWarehouseItem, error) {
	result := make(map[string][]ProductWarehouseItem)
	for rows.Next() {
		variantID, item, err := scanVariantWarehouse(rows)
		if err != nil {
			return nil, err
		}
		result[variantID] = append(result[variantID], item)
	}
	return result, rows.Err()
}

func scanVariantWarehouse(rows pgx.Rows) (string, ProductWarehouseItem, error) {
	variantID := ""
	item := ProductWarehouseItem{}
	err := rows.Scan(&variantID, &item.WarehouseID, &item.Warehouse, &item.Qty)
	return variantID, item, err
}
