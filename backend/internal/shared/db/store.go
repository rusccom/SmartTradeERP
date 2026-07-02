package db

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
    Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// ErrNoTenant is returned when a tenant-scoped call is attempted without a
// valid tenant id. Fail-closed: the query never reaches the database.
var ErrNoTenant = errors.New("db: tenant id is missing or invalid")

const setTenantSQL = `SELECT set_config('app.tenant_id', $1, true)`

type Store struct {
    // Pool is raw, tenant-UNSCOPED access. Under RLS (migration 000016) all
    // tenant-owned tables appear empty through it. Only pre-tenant flows may
    // use it: login, registration, platform admin, storefront host lookup,
    // global dictionaries, health checks. Everything else goes through the
    // tenant-scoped Exec/Query/QueryRow/WithTenantTx.
    Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
    return &Store{Pool: pool}
}

func (s *Store) Ping(ctx context.Context) error {
    return s.Pool.Ping(ctx)
}

// WithTx runs fn in a transaction WITHOUT tenant scoping. Reserved for flows
// that legitimately have no tenant yet (registration, platform admin).
// Tenant-owned tables stay invisible inside it until SetTenant is called.
func (s *Store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
    tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)
    if err := fn(tx); err != nil {
        return err
    }
    return tx.Commit(ctx)
}

// WithTenantTx runs fn in a transaction whose app.tenant_id is set, so the
// RLS policies expose exactly this tenant's rows for every statement inside.
func (s *Store) WithTenantTx(ctx context.Context, tenantID string, fn func(tx pgx.Tx) error) error {
    if !validTenant(tenantID) {
        return ErrNoTenant
    }
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        if _, err := tx.Exec(ctx, setTenantSQL, tenantID); err != nil {
            return err
        }
        return fn(tx)
    })
}

// SetTenant scopes an already-open transaction to a tenant. Used by
// registration, which creates the tenant row first and then its child graph.
func SetTenant(ctx context.Context, tx pgx.Tx, tenantID string) error {
    if !validTenant(tenantID) {
        return ErrNoTenant
    }
    _, err := tx.Exec(ctx, setTenantSQL, tenantID)
    return err
}

func validTenant(tenantID string) bool {
    if tenantID == "" {
        return false
    }
    _, err := uuid.Parse(tenantID)
    return err == nil
}
