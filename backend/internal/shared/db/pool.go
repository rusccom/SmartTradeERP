package db

import (
	"context"
	"net/url"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolPurpose uint8

const (
	ServicePool PoolPurpose = iota
	MigrationPool
)

func Connect(ctx context.Context, databaseURL string, purpose PoolPurpose) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	configurePool(config, databaseURL, purpose)
	config.AfterConnect = registerTypes
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func configurePool(config *pgxpool.Config, databaseURL string, purpose PoolPurpose) {
	limit := int32(4)
	if purpose == MigrationPool {
		limit = 1
		config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}
	config.MaxConns = min(config.MaxConns, limit)
	if !poolParamSet(databaseURL, "pool_min_conns") {
		config.MinConns = 1
	}
	config.MinConns = min(config.MinConns, config.MaxConns)
	config.MinIdleConns = min(config.MinIdleConns, config.MaxConns)
}

func poolParamSet(databaseURL string, key string) bool {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return false
	}
	_, ok := parsed.Query()[key]
	return ok
}

func registerTypes(_ context.Context, conn *pgx.Conn) error {
	pgxdecimal.Register(conn.TypeMap())
	return nil
}
