package app

import (
	"context"
	"net/http"
	"time"

	"smarterp/backend/internal/platform/storefront"
	"smarterp/backend/internal/shared/config"
	"smarterp/backend/internal/shared/db"
)

// BuildStorefront constructs the public storefront HTTP server. It runs as a
// separate process from the ERP API with its own database pool and request
// timeouts so a crawler storm on a shop cannot starve the authenticated API.
func BuildStorefront(ctx context.Context, cfg config.Config) (*http.Server, func(), error) {
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}
	store := db.NewStore(pool)
	mux := http.NewServeMux()
	opts := storefront.Options{MediaBaseURL: cfg.R2.PublicBaseURL, JWTSecret: cfg.JWTSecret}
	if err := storefront.Register(mux, store, opts); err != nil {
		pool.Close()
		return nil, nil, err
	}
	cleanup := func() { pool.Close() }
	return storefrontServer(cfg.HTTPAddr, mux), cleanup, nil
}

func storefrontServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
