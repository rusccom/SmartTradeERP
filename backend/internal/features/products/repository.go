package products

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

func (r *Repository) List(ctx context.Context, tenantID, isComposite string, page, perPage int) ([]Product, int, error) {
    total, err := r.count(ctx, tenantID, isComposite)
    if err != nil {
        return nil, 0, err
    }
    data, err := r.load(ctx, tenantID, isComposite, page, perPage)
    if err != nil {
        return nil, 0, err
    }
    return data, total, nil
}

func (r *Repository) count(ctx context.Context, tenantID, isComposite string) (int, error) {
    query := `SELECT COUNT(*) FROM catalog.products WHERE tenant_id=$1`
    args := []any{tenantID}
    if isComposite != "" {
        query += ` AND is_composite=$2`
        args = append(args, isComposite == "true")
    }
    row := r.store.Pool.QueryRow(ctx, query, args...)
    total := 0
    return total, row.Scan(&total)
}

func (r *Repository) load(ctx context.Context, tenantID, isComposite string, page, perPage int) ([]Product, error) {
    offset := (page - 1) * perPage
    query := `SELECT id::text, name, is_composite, created_at::text, updated_at::text
        FROM catalog.products WHERE tenant_id=$1`
    args := []any{tenantID}
    query, args = addCompositeFilter(query, args, isComposite)
    query += ` ORDER BY created_at DESC LIMIT $` + position(len(args)+1)
    query += ` OFFSET $` + position(len(args)+2)
    args = append(args, perPage, offset)
    rows, err := r.store.Pool.Query(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanProducts(rows)
}

func addCompositeFilter(query string, args []any, isComposite string) (string, []any) {
    if isComposite == "" {
        return query, args
    }
    query += ` AND is_composite=$` + position(len(args)+1)
    args = append(args, isComposite == "true")
    return query, args
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

func (r *Repository) CreateDefaultVariant(ctx context.Context, tx pgx.Tx, variantID, productID string, req CreateRequest) error {
    query := `INSERT INTO catalog.product_variants
        (id, product_id, name, sku_code, barcode, unit, price)
        VALUES ($1,$2,'Default',$3,$4,$5,$6)`
    _, err := tx.Exec(ctx, query, variantID, productID, req.SKUCode, req.Barcode, req.Unit, req.Price)
    return err
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

func (r *Repository) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    query := `UPDATE catalog.products
        SET name=$3, is_composite=$4, updated_at=now()
        WHERE tenant_id=$1 AND id=$2`
    _, err := r.store.Pool.Exec(ctx, query, tenantID, id, req.Name, req.IsComposite)
    return err
}

func (r *Repository) Delete(ctx context.Context, tenantID, id string) error {
    query := `DELETE FROM catalog.products WHERE tenant_id=$1 AND id=$2`
    _, err := r.store.Pool.Exec(ctx, query, tenantID, id)
    return err
}
