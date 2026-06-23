package api

import (
	"context"
	"net/http"
	"time"

	"smarterp/backend/internal/erp/bundles"
	"smarterp/backend/internal/erp/ledger"
	mediafeature "smarterp/backend/internal/platform/media"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/storage"
)

const healthDBTimeout = 2 * time.Second

func Register(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	mediaStore storage.ObjectStore,
	mediaURL string,
) error {
	mux.HandleFunc("GET /health", health(store))
	ledgerService := ledger.NewService(store)
	bundleService := newBundleService(store)
	mediaService := newMediaService(store, mediaStore)
	registerAdmin(mux, store, tokens)
	registerClientAuth(mux, store, tokens)
	registerCatalog(mux, catalogDeps{
		store:    store,
		tokens:   tokens,
		ledger:   ledgerService,
		bundles:  bundleService,
		media:    mediaService,
		mediaURL: mediaURL,
	})
	registerOperations(mux, store, tokens, ledgerService, bundleService)
	registerSettings(mux, store, tokens)
	return registerStorefrontAdmin(mux, store, tokens)
}

func newBundleService(store *db.Store) *bundles.Service {
	repo := bundles.NewRepository(store)
	return bundles.NewService(store, repo)
}

func newMediaService(store *db.Store, objectStore storage.ObjectStore) *mediafeature.Service {
	repo := mediafeature.NewRepository(store)
	return mediafeature.NewService(store, repo, objectStore)
}

func health(store *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), healthDBTimeout)
		defer cancel()
		if err := store.Ping(ctx); err != nil {
			httpx.WriteError(w, http.StatusServiceUnavailable, "db_unavailable", "database unavailable", nil)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleClient(
	mux *http.ServeMux,
	tokens *auth.TokenService,
	pattern string,
	handler http.HandlerFunc,
) {
	middleware := auth.RequireScope(tokens, "client")
	mux.Handle(pattern, auth.Chain(http.HandlerFunc(handler), middleware))
}

func handleAdmin(
	mux *http.ServeMux,
	tokens *auth.TokenService,
	pattern string,
	handler http.HandlerFunc,
) {
	middleware := auth.RequireScope(tokens, "admin")
	mux.Handle(pattern, auth.Chain(http.HandlerFunc(handler), middleware))
}
