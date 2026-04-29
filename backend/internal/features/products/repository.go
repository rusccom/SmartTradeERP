package products

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/search"
)

type Repository struct {
    store *db.Store
}

func NewRepository(store *db.Store) *Repository {
    return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID string, query ProductListQuery) ([]Product, int, error) {
	total, err := r.count(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	data, err := r.load(ctx, tenantID, query)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (r *Repository) count(ctx context.Context, tenantID string, query ProductListQuery) (int, error) {
	sqlQuery := `SELECT COUNT(*) FROM catalog.products p WHERE p.tenant_id=$1`
	args := []any{tenantID}
	sqlQuery, args = appendListFilters(sqlQuery, args, query)
	row := r.store.Pool.QueryRow(ctx, sqlQuery, args...)
	total := 0
	err := row.Scan(&total)
	return total, err
}

func (r *Repository) load(ctx context.Context, tenantID string, query ProductListQuery) ([]Product, error) {
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
	return scanProducts(rows)
}

func appendListFilters(query string, args []any, productQuery ProductListQuery) (string, []any) {
	listQuery := productQuery.List
	query = addProductScopeFilter(query)
	query, args = search.AppendProductSearch(query, args, listQuery.Search)
	return appendProductStockFilter(query, args, productQuery.Stock)
}

func addProductScopeFilter(query string) string {
	return query + ` AND p.is_composite=false`
}

func appendSortAndPaging(query string, args []any, productQuery ProductListQuery) (string, []any) {
	sortBy, sortDir := readSort(productQuery.List)
	query += ` ORDER BY ` + sortBy + ` ` + sortDir
	query += ` LIMIT $` + position(len(args)+1)
	query += ` OFFSET $` + position(len(args)+2)
	listQuery := productQuery.List
	args = append(args, listQuery.PerPage, httpx.Offset(listQuery.Page, listQuery.PerPage))
	return query, args
}

func readSort(query httpx.ListQuery) (string, string) {
	sortBy := query.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDir := query.SortDir
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}
	return sortBy, sortDir
}

func position(value int) string {
	return strconv.Itoa(value)
}

func scanProducts(rows pgx.Rows) ([]Product, error) {
	items := make([]Product, 0)
	for rows.Next() {
		item := Product{}
		err := rows.Scan(&item.ID, &item.Name, &item.IsComposite, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) Create(ctx context.Context, tx pgx.Tx, tenantID, productID string, req CreateRequest) error {
    query := `INSERT INTO catalog.products (id, tenant_id, name, is_composite)
        VALUES ($1,$2,$3,$4)`
    _, err := tx.Exec(ctx, query, productID, tenantID, req.Name, req.IsComposite)
    return err
}

func (r *Repository) CreateDefaultVariant(ctx context.Context, tx pgx.Tx, input createProductTx) error {
    query := `INSERT INTO catalog.product_variants
        (id, tenant_id, product_id, name, sku_code, barcode, unit, price)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
    _, err := tx.Exec(ctx, query, input.variantID, input.tenantID, input.productID,
        readDefaultVariantName(input.req),
        input.req.SKUCode, input.req.Barcode, input.req.Unit, input.req.Price)
    return err
}

func readDefaultVariantName(req CreateRequest) string {
    if req.VariantName != "" {
        return req.VariantName
    }
    return "Default"
}

func (r *Repository) GetByID(ctx context.Context, tenantID, id string) (Product, error) {
    query := `SELECT id::text, name, is_composite, created_at::text, updated_at::text
        FROM catalog.products
        WHERE tenant_id=$1 AND id=$2`
    row := r.store.Pool.QueryRow(ctx, query, tenantID, id)
    item := Product{}
    err := row.Scan(&item.ID, &item.Name, &item.IsComposite, &item.CreatedAt, &item.UpdatedAt)
    return item, err
}

func (r *Repository) CompositeFlag(ctx context.Context, tenantID, id string) (bool, error) {
	query := `SELECT is_composite FROM catalog.products WHERE tenant_id=$1 AND id=$2`
	value := false
	err := r.store.Pool.QueryRow(ctx, query, tenantID, id).Scan(&value)
	return value, err
}

func (r *Repository) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    query := `UPDATE catalog.products
        SET name=$3, is_composite=$4, updated_at=now()
        WHERE tenant_id=$1 AND id=$2`
    tag, err := r.store.Pool.Exec(ctx, query, tenantID, id, req.Name, req.IsComposite)
    if err == nil && tag.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }
    return err
}

func (r *Repository) Delete(ctx context.Context, tenantID, id string) error {
    query := `DELETE FROM catalog.products WHERE tenant_id=$1 AND id=$2`
    tag, err := r.store.Pool.Exec(ctx, query, tenantID, id)
    if err == nil && tag.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }
    return err
}
