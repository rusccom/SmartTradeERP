package db

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
    Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Store struct {
    Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
    return &Store{Pool: pool}
}

func (s *Store) Ping(ctx context.Context) error {
    return s.Pool.Ping(ctx)
}

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
