package db

import (
	"context"
	"net/url"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	configurePool(config, databaseURL)
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

func configurePool(config *pgxpool.Config, databaseURL string) {
	if !poolParamSet(databaseURL, "pool_max_conns") {
		config.MaxConns = 20
	}
	if !poolParamSet(databaseURL, "pool_min_conns") {
		config.MinConns = 1
	}
	if !poolParamSet(databaseURL, "pool_max_conn_lifetime") {
		config.MaxConnLifetime = time.Hour
	}
	if !poolParamSet(databaseURL, "pool_max_conn_idle_time") {
		config.MaxConnIdleTime = 30 * time.Minute
	}
	if !poolParamSet(databaseURL, "pool_health_check_period") {
		config.HealthCheckPeriod = time.Minute
	}
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
