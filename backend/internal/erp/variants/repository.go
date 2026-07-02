package variants

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
)

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID, productID string, page, perPage int) ([]Variant, int, error) {
    total, err := r.count(ctx, tenantID, productID)
    if err != nil {
        return nil, 0, err
    }
    data, err := r.load(ctx, tenantID, productID, page, perPage)
    if err != nil {
        return nil, 0, err
    }
    return data, total, nil
}

func (r *Repository) count(ctx context.Context, tenantID, productID string) (int, error) {
    query := `SELECT COUNT(*)
        FROM catalog.product_variants v`
    args := []any{}
    query, args = addProductFilter(query, args, productID)
    row := r.store.QueryRow(ctx, tenantID, query, args...)
    total := 0
    err := row.Scan(&total)
    return total, err
}

func (r *Repository) load(ctx context.Context, tenantID, productID string, page, perPage int) ([]Variant, error) {
    query := `SELECT v.id::text, v.product_id::text, COALESCE(v.name,''),
        COALESCE(v.sku_code,''), COALESCE(v.barcode,''), v.unit, COALESCE(v.price,0),
        COALESCE(v.option1,''), COALESCE(v.option2,''), COALESCE(v.option3,'')
        FROM catalog.product_variants v`
    args := []any{}
    query, args = addProductFilter(query, args, productID)
    query, args = addLimitOffset(query, args, page, perPage)
    rows, err := r.store.Query(ctx, tenantID, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanVariants(rows)
}

func addProductFilter(query string, args []any, productID string) (string, []any) {
    if productID == "" {
        return query, args
    }
    query += ` WHERE v.product_id=$` + strconv.Itoa(len(args)+1)
    args = append(args, productID)
    return query, args
}

func addLimitOffset(query string, args []any, page, perPage int) (string, []any) {
    query += ` ORDER BY v.created_at DESC`
    query += ` LIMIT $` + strconv.Itoa(len(args)+1)
    query += ` OFFSET $` + strconv.Itoa(len(args)+2)
    offset := (page - 1) * perPage
    args = append(args, perPage, offset)
    return query, args
}

func scanVariants(rows pgx.Rows) ([]Variant, error) {
    items := make([]Variant, 0)
    for rows.Next() {
        item := Variant{}
        err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.SKUCode,
            &item.Barcode, &item.Unit, &item.Price, &item.Option1, &item.Option2, &item.Option3)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func (r *Repository) Create(ctx context.Context, tenantID, variantID string, req CreateRequest) error {
    query := `INSERT INTO catalog.product_variants
        (id, tenant_id, product_id, name, sku_code, barcode, unit, price, option1, option2, option3)
        SELECT $1, p.tenant_id, p.id, $2, $3, $4, $5, $6, $7, $8, $9
        FROM catalog.products p
        WHERE p.id=$10`
    tag, err := r.store.Exec(ctx, tenantID, query, variantID, req.Name, req.SKUCode, req.Barcode,
        req.Unit, req.Price, req.Option1, req.Option2, req.Option3, req.ProductID)
    if err == nil && tag.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }
    return err
}

func (r *Repository) ByID(ctx context.Context, tenantID, id string) (Variant, error) {
    query := `SELECT v.id::text, v.product_id::text, COALESCE(v.name,''),
        COALESCE(v.sku_code,''), COALESCE(v.barcode,''), v.unit, COALESCE(v.price,0),
        COALESCE(v.option1,''), COALESCE(v.option2,''), COALESCE(v.option3,'')
        FROM catalog.product_variants v
        WHERE v.id=$1`
    row := r.store.QueryRow(ctx, tenantID, query, id)
    item := Variant{}
    err := row.Scan(&item.ID, &item.ProductID, &item.Name, &item.SKUCode,
        &item.Barcode, &item.Unit, &item.Price, &item.Option1, &item.Option2, &item.Option3)
    return item, err
}

func (r *Repository) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    query := `UPDATE catalog.product_variants v
        SET name=$2, sku_code=$3, barcode=$4, unit=$5, price=$6,
            option1=$7, option2=$8, option3=$9
        WHERE v.id=$1`
    tag, err := r.store.Exec(ctx, tenantID, query, id, req.Name, req.SKUCode,
        req.Barcode, req.Unit, req.Price, req.Option1, req.Option2, req.Option3)
    if err == nil && tag.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }
    return err
}

func (r *Repository) Delete(ctx context.Context, tenantID, id string) error {
    query := `DELETE FROM catalog.product_variants
        WHERE id=$1`
    tag, err := r.store.Exec(ctx, tenantID, query, id)
    if err == nil && tag.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }
    return err
}

func (r *Repository) ProductVariantCount(ctx context.Context, tenantID, variantID string) (int, error) {
    query := `SELECT COUNT(*)
        FROM catalog.product_variants v
        JOIN catalog.product_variants source ON source.product_id=v.product_id
        WHERE source.id=$1`
    row := r.store.QueryRow(ctx, tenantID, query, variantID)
    count := 0
	err := row.Scan(&count)
	return count, err
}

func (r *Repository) WarehouseStock(ctx context.Context, tenantID, variantID string) ([]StockByWarehouse, error) {
    query := `SELECT w.id::text, w.name,
        COALESCE(sb.qty,0)
        FROM catalog.warehouses w
        LEFT JOIN ledger.stock_balances sb
            ON sb.warehouse_id=w.id AND sb.variant_id=$1
        GROUP BY w.id, w.name, sb.qty
        ORDER BY w.created_at`
    rows, err := r.store.Query(ctx, tenantID, query, variantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanStockByWarehouse(rows)
}

func scanStockByWarehouse(rows pgx.Rows) ([]StockByWarehouse, error) {
    items := make([]StockByWarehouse, 0)
    for rows.Next() {
        item := StockByWarehouse{}
        err := rows.Scan(&item.WarehouseID, &item.Warehouse, &item.Qty)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}
