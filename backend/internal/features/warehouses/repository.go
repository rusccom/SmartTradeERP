package warehouses

import (
    "context"

    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/shared/db"
)

type Repository struct {
    store *db.Store
}

func NewRepository(store *db.Store) *Repository {
    return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID string) ([]Warehouse, error) {
    query := `SELECT id::text, name, COALESCE(address,''), is_default, is_active, created_at::text
        FROM catalog.warehouses
        WHERE tenant_id=$1
        ORDER BY created_at`
    rows, err := r.store.Pool.Query(ctx, query, tenantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scan(rows)
}

func scan(rows pgx.Rows) ([]Warehouse, error) {
    items := make([]Warehouse, 0)
    for rows.Next() {
        item := Warehouse{}
        err := rows.Scan(&item.ID, &item.Name, &item.Address, &item.IsDefault, &item.IsActive, &item.CreatedAt)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func (r *Repository) Create(ctx context.Context, tenantID, id string, req CreateRequest) error {
    query := `INSERT INTO catalog.warehouses
        (id, tenant_id, name, address, is_default, is_active)
        VALUES ($1,$2,$3,$4,$5,$6)`
    _, err := r.store.Pool.Exec(ctx, query, id, tenantID, req.Name, req.Address, req.IsDefault, req.IsActive)
    return err
}

func (r *Repository) Update(ctx context.Context, tenantID, id string, req UpdateRequest) error {
    query := `UPDATE catalog.warehouses
        SET name=$3, address=$4, is_default=$5, is_active=$6
        WHERE tenant_id=$1 AND id=$2`
    _, err := r.store.Pool.Exec(ctx, query, tenantID, id, req.Name, req.Address, req.IsDefault, req.IsActive)
    return err
}

func (r *Repository) Delete(ctx context.Context, tenantID, id string) error {
    _, err := r.store.Pool.Exec(ctx, `DELETE FROM catalog.warehouses WHERE tenant_id=$1 AND id=$2`, tenantID, id)
    return err
}

func (r *Repository) CountDefaults(ctx context.Context, tenantID string) (int, error) {
    row := r.store.Pool.QueryRow(ctx,
        `SELECT COUNT(*) FROM catalog.warehouses WHERE tenant_id=$1 AND is_default=true`, tenantID)
    count := 0
    return count, row.Scan(&count)
}

func (r *Repository) IsDefault(ctx context.Context, tenantID, id string) (bool, error) {
    row := r.store.Pool.QueryRow(ctx,
        `SELECT is_default FROM catalog.warehouses WHERE tenant_id=$1 AND id=$2`, tenantID, id)
    isDefault := false
    return isDefault, row.Scan(&isDefault)
}
