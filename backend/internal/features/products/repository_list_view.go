package products

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

type productVariantRow struct {
	ProductID string
	Item      ProductVariantItem
}

type productIncludeRequest struct {
	tenantID string
	query    ProductListQuery
	include  ProductListInclude
	items    []ProductListItem
	total    int
}

func (r *Repository) ListWithIncludes(
	ctx context.Context,
	tenantID string,
	query ProductListQuery,
	include ProductListInclude,
) ([]ProductListItem, int, error) {
	total, err := r.count(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	items, err := r.loadListItems(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	req := productIncludeRequest{tenantID: tenantID, query: query, include: include, items: items, total: total}
	return r.attachProductIncludes(ctx, req)
}

func (r *Repository) loadListItems(
	ctx context.Context,
	tenantID string,
	query ProductListQuery,
) ([]ProductListItem, error) {
	sqlQuery := `SELECT p.id::text, p.name, p.is_composite, p.created_at::text, p.updated_at::text
        FROM catalog.products p WHERE p.tenant_id=$1`
	args := []any{tenantID}
	sqlQuery, args = appendListFilters(sqlQuery, args, query)
	sqlQuery, args = appendSortAndPaging(sqlQuery, args, query)
	rows, err := r.store.Pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProductListItems(rows)
}

func scanProductListItems(rows pgx.Rows) ([]ProductListItem, error) {
	items := make([]ProductListItem, 0)
	for rows.Next() {
		item := ProductListItem{}
		err := rows.Scan(&item.ID, &item.Name, &item.IsComposite, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		item.Variants = []ProductVariantItem{}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) attachProductIncludes(
	ctx context.Context,
	req productIncludeRequest,
) ([]ProductListItem, int, error) {
	if !req.include.Variants || len(req.items) == 0 {
		return req.items, req.total, nil
	}
	variants, err := r.loadProductVariantRows(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if req.include.Warehouses {
		variants, err = r.attachVariantWarehouses(ctx, req.tenantID, variants, req.query.Stock)
	}
	return attachProductVariants(req.items, variants), req.total, err
}

func (r *Repository) loadProductVariantRows(
	ctx context.Context,
	req productIncludeRequest,
) ([]productVariantRow, error) {
	query := productVariantRowsSQL()
	args := []any{req.tenantID}
	query, args = appendProductIDList(query, args, req.items)
	query, args = appendVariantStockFilter(query, args, req.query.Stock)
	query += ` ORDER BY v.created_at`
	rows, err := r.store.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProductVariantRows(rows, req.include)
}

func productVariantRowsSQL() string {
	return `SELECT v.product_id::text, v.id::text, COALESCE(v.name,'Default'),
        COALESCE(v.sku_code,''), COALESCE(v.barcode,''), v.unit, COALESCE(v.price,0),
        COALESCE(latest.running_qty,0), COALESCE(latest.running_avg,0)
        FROM catalog.product_variants v
        LEFT JOIN LATERAL (` + latestVariantStockSQL() + `) latest ON true
        WHERE v.tenant_id=$1`
}

func latestVariantStockSQL() string {
	return `SELECT running_qty, running_avg
            FROM ledger.cost_ledger l
            WHERE l.tenant_id=$1 AND l.variant_id=v.id
            ORDER BY l.sequence_num DESC
            LIMIT 1`
}

func appendProductIDList(
	query string,
	args []any,
	items []ProductListItem,
) (string, []any) {
	positions := make([]string, 0, len(items))
	for _, item := range items {
		positions = append(positions, "$"+position(len(args)+1))
		args = append(args, item.ID)
	}
	query += ` AND v.product_id::text IN (` + strings.Join(positions, ",") + `)`
	return query, args
}

func scanProductVariantRows(rows pgx.Rows, include ProductListInclude) ([]productVariantRow, error) {
	items := make([]productVariantRow, 0)
	for rows.Next() {
		item, err := scanProductVariantRow(rows, include)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanProductVariantRow(rows pgx.Rows, include ProductListInclude) (productVariantRow, error) {
	row := productVariantRow{Item: ProductVariantItem{Warehouses: []ProductWarehouseItem{}}}
	err := rows.Scan(&row.ProductID, &row.Item.ID, &row.Item.Name, &row.Item.SKUCode,
		&row.Item.Barcode, &row.Item.Unit, &row.Item.Price, &row.Item.GlobalQty, &row.Item.AvgCost)
	if !include.Stock {
		row.Item.GlobalQty = row.Item.GlobalQty.Sub(row.Item.GlobalQty)
		row.Item.AvgCost = row.Item.AvgCost.Sub(row.Item.AvgCost)
	}
	row.Item.StockValue = row.Item.GlobalQty.Mul(row.Item.AvgCost).Round(4)
	return row, err
}

func attachProductVariants(items []ProductListItem, variants []productVariantRow) []ProductListItem {
	index := productIndex(items)
	for _, row := range variants {
		position, ok := index[row.ProductID]
		if ok {
			items[position] = appendProductVariant(items[position], row.Item)
		}
	}
	return items
}

func productIndex(items []ProductListItem) map[string]int {
	result := make(map[string]int, len(items))
	for index, item := range items {
		result[item.ID] = index
	}
	return result
}

func appendProductVariant(item ProductListItem, variant ProductVariantItem) ProductListItem {
	item.GlobalQty = item.GlobalQty.Add(variant.GlobalQty)
	item.StockValue = item.StockValue.Add(variant.StockValue).Round(4)
	item.Variants = append(item.Variants, variant)
	return item
}
