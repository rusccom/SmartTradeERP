package app

import (
	"context"
	"net/http"

	"smarterp/backend/internal/api"
	"smarterp/backend/internal/integrations/smartbundle"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/config"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/storage"
)

func Build(ctx context.Context, cfg config.Config) (*http.Server, func(), error) {
	pool, err := db.Connect(ctx, cfg.DatabaseURL, db.ServicePool)
	if err != nil {
		return nil, nil, err
	}
	handler, err := apiRouter(db.NewStore(pool), cfg)
	if err != nil {
		pool.Close()
		return nil, nil, err
	}
	return &http.Server{Addr: cfg.HTTPAddr, Handler: handler}, pool.Close, nil
}

func apiRouter(store *db.Store, cfg config.Config) (http.Handler, error) {
	mediaStore, err := storage.NewR2Store(cfg.R2)
	if err != nil {
		return nil, err
	}
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTTL)
	mux := http.NewServeMux()
	if err := api.Register(mux, store, tokens, mediaStore, cfg); err != nil {
		return nil, err
	}
	smartbundle.RegisterSettings(mux, store, tokens)
	return cors(mux), nil
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
