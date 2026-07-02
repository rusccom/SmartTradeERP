package db

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Tenant-scoped single-query access. Each call sends one pgx batch --
// set_config('app.tenant_id', ...) followed by the query -- which Postgres
// executes in one implicit transaction: a single round trip, and the setting
// is transaction-local so it can never leak to another pooled connection.

func (s *Store) Exec(ctx context.Context, tenantID, sql string, args ...any) (pgconn.CommandTag, error) {
    if !validTenant(tenantID) {
        return pgconn.CommandTag{}, ErrNoTenant
    }
    var tag pgconn.CommandTag
    err := s.Pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
        results := conn.SendBatch(ctx, scopedBatch(tenantID, sql, args))
        defer results.Close()
        if _, err := results.Exec(); err != nil {
            return err
        }
        var execErr error
        tag, execErr = results.Exec()
        return execErr
    })
    return tag, err
}

func (s *Store) Query(ctx context.Context, tenantID, sql string, args ...any) (pgx.Rows, error) {
    if !validTenant(tenantID) {
        return nil, ErrNoTenant
    }
    conn, err := s.Pool.Acquire(ctx)
    if err != nil {
        return nil, err
    }
    rows, err := scopedQuery(ctx, conn, tenantID, sql, args)
    if err != nil {
        conn.Release()
        return nil, err
    }
    return rows, nil
}

func (s *Store) QueryRow(ctx context.Context, tenantID, sql string, args ...any) pgx.Row {
    return &scopedRow{store: s, ctx: ctx, tenantID: tenantID, sql: sql, args: args}
}

func scopedQuery(ctx context.Context, conn *pgxpool.Conn, tenantID, sql string, args []any) (pgx.Rows, error) {
    results := conn.SendBatch(ctx, scopedBatch(tenantID, sql, args))
    if _, err := results.Exec(); err != nil {
        results.Close()
        return nil, err
    }
    rows, err := results.Query()
    if err != nil {
        results.Close()
        return nil, err
    }
    return &scopedRows{Rows: rows, results: results, conn: conn}, nil
}

func scopedBatch(tenantID, sql string, args []any) *pgx.Batch {
    batch := &pgx.Batch{}
    batch.Queue(setTenantSQL, tenantID)
    batch.Queue(sql, args...)
    return batch
}

// scopedRows keeps the pooled connection alive until the caller closes the
// rows, then releases batch results and connection in order. Close is safe to
// call more than once.
type scopedRows struct {
    pgx.Rows
    results pgx.BatchResults
    conn    *pgxpool.Conn
}

func (r *scopedRows) Close() {
    r.Rows.Close()
    r.results.Close()
    r.conn.Release()
}

// scopedRow defers execution until Scan, mirroring pgx.Row semantics
// (pgx.ErrNoRows when the query matches nothing).
type scopedRow struct {
    store    *Store
    ctx      context.Context
    tenantID string
    sql      string
    args     []any
}

func (r *scopedRow) Scan(dest ...any) error {
    rows, err := r.store.Query(r.ctx, r.tenantID, r.sql, r.args...)
    if err != nil {
        return err
    }
    defer rows.Close()
    return scanFirst(rows, dest)
}

func scanFirst(rows pgx.Rows, dest []any) error {
    if !rows.Next() {
        if err := rows.Err(); err != nil {
            return err
        }
        return pgx.ErrNoRows
    }
    if err := rows.Scan(dest...); err != nil {
        return err
    }
    rows.Close()
    return rows.Err()
}
