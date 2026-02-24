package clientauth

import (
    "context"

    "github.com/jackc/pgx/v5"
)

type Repository struct{}

func NewRepository() *Repository {
    return &Repository{}
}

func (r *Repository) FindByEmail(ctx context.Context, dbtx interface {
    QueryRow(context.Context, string, ...any) pgx.Row
}, email string) (userRecord, error) {
    query := `SELECT id::text, tenant_id::text, role, password_hash
        FROM platform.tenant_users
        WHERE email=$1 AND is_active=true`
    row := dbtx.QueryRow(ctx, query, email)
    item := userRecord{}
    err := row.Scan(&item.ID, &item.TenantID, &item.Role, &item.PasswordHash)
    return item, err
}
