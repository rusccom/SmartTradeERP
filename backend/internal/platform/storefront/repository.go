package storefront

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
)

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

// FindTenantByHost looks up the tenant that owns an active storefront host.
// Returns ErrTenantNotFound when no active mapping exists.
func (r *Repository) FindTenantByHost(ctx context.Context, host string) (ResolvedTenant, error) {
	const query = `SELECT d.tenant_id, t.status
        FROM platform.storefront_domains d
        JOIN platform.tenants t ON t.id = d.tenant_id
        WHERE d.host = $1 AND d.status = 'active'`
	row := r.store.Pool.QueryRow(ctx, query, host)
	resolved := ResolvedTenant{}
	if err := row.Scan(&resolved.TenantID, &resolved.Status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ResolvedTenant{}, ErrTenantNotFound
		}
		return ResolvedTenant{}, err
	}
	return resolved, nil
}
