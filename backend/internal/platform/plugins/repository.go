package plugins

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
)

var ErrUnavailable = errors.New("plugin access is required")

type State struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Granted bool   `json:"granted"`
	Enabled bool   `json:"enabled"`
}

type Repository struct{ store *db.Store }

func NewRepository(store *db.Store) *Repository { return &Repository{store: store} }

func (r *Repository) List(ctx context.Context, id string) ([]State, error) {
	const query = `SELECT p.key,p.name,t.plugin_key IS NOT NULL,COALESCE(t.enabled,false)
		FROM platform.plugins p LEFT JOIN platform.tenant_plugins t ON t.plugin_key=p.key AND t.tenant_id=$1
		ORDER BY p.name,p.key`
	rows, err := r.store.Query(ctx, id, query, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[State])
}

func (r *Repository) Grant(ctx context.Context, id, key string, value bool) error {
	const grant = `INSERT INTO platform.tenant_plugins (tenant_id,plugin_key)
		VALUES ($1,$2) ON CONFLICT DO NOTHING`
	const revoke = `DELETE FROM platform.tenant_plugins WHERE tenant_id=$1 AND plugin_key=$2`
	return r.store.WithTenantTx(ctx, id, func(tx pgx.Tx) error {
		var exists bool
		err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT FROM platform.tenants WHERE id=$1)
			AND EXISTS(SELECT FROM platform.plugins WHERE key=$2)`, id, key).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return pgx.ErrNoRows
		}
		query := revoke
		if value {
			query = grant
		}
		_, err = tx.Exec(ctx, query, id, key)
		return err
	})
}

func (r *Repository) Enable(ctx context.Context, id, key string, value bool) error {
	const query = `UPDATE platform.tenant_plugins SET enabled=$3 WHERE tenant_id=$1 AND plugin_key=$2`
	return r.store.WithTenantTx(ctx, id, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, query, id, key, value)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrUnavailable
		}
		return err
	})
}

// Enabled locks an existing grant for the duration of a write that uses it.
func Enabled(ctx context.Context, tx db.DBTX, id, key string) (bool, error) {
	var enabled bool
	err := tx.QueryRow(ctx, `SELECT enabled FROM platform.tenant_plugins WHERE tenant_id=$1 AND plugin_key=$2 FOR SHARE`, id, key).Scan(&enabled)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return enabled, err
}

func RequireEnabled(ctx context.Context, tx db.DBTX, id, key string) error {
	on, err := Enabled(ctx, tx, id, key)
	if err == nil && !on {
		return ErrUnavailable
	}
	return err
}
