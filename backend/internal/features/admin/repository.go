package admin

import (
    "context"

    "smarterp/backend/internal/shared/db"
)

type Repository struct {
    store *db.Store
}

func NewRepository(store *db.Store) *Repository {
    return &Repository{store: store}
}

func (r *Repository) FindAdminByEmail(ctx context.Context, email string) (adminUser, error) {
    query := `SELECT id::text, password_hash
        FROM platform.platform_admins
        WHERE email=$1`
    row := r.store.Pool.QueryRow(ctx, query, email)
    user := adminUser{}
    err := row.Scan(&user.ID, &user.PasswordHash)
    return user, err
}

func (r *Repository) ListTenants(ctx context.Context, page, perPage int) ([]Tenant, int, error) {
    total, err := r.countTenants(ctx)
    if err != nil {
        return nil, 0, err
    }
    tenants, err := r.loadTenants(ctx, page, perPage)
    if err != nil {
        return nil, 0, err
    }
    return tenants, total, nil
}

func (r *Repository) countTenants(ctx context.Context) (int, error) {
    row := r.store.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM platform.tenants`)
    total := 0
    err := row.Scan(&total)
    return total, err
}

func (r *Repository) loadTenants(ctx context.Context, page, perPage int) ([]Tenant, error) {
    offset := (page - 1) * perPage
    query := `SELECT id::text, name, status, plan, created_at::text
        FROM platform.tenants
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2`
    rows, err := r.store.Pool.Query(ctx, query, perPage, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanTenants(rows)
}

func scanTenants(rows interface{ Next() bool; Scan(...any) error; Err() error }) ([]Tenant, error) {
    items := make([]Tenant, 0)
    for rows.Next() {
        item := Tenant{}
        err := rows.Scan(&item.ID, &item.Name, &item.Status, &item.Plan, &item.CreatedAt)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func (r *Repository) GetTenantByID(ctx context.Context, id string) (Tenant, error) {
    query := `SELECT id::text, name, status, plan, created_at::text
        FROM platform.tenants
        WHERE id=$1`
    row := r.store.Pool.QueryRow(ctx, query, id)
    item := Tenant{}
    err := row.Scan(&item.ID, &item.Name, &item.Status, &item.Plan, &item.CreatedAt)
    return item, err
}

func (r *Repository) Stats(ctx context.Context) (Stats, error) {
    query := `SELECT
        COUNT(*) AS total,
        COUNT(*) FILTER (WHERE status='active') AS active,
        COUNT(*) FILTER (WHERE created_at >= now() - interval '30 day') AS recent
        FROM platform.tenants`
    row := r.store.Pool.QueryRow(ctx, query)
    stats := Stats{}
    err := row.Scan(&stats.TotalTenants, &stats.ActiveTenants, &stats.NewLast30Days)
    return stats, err
}
