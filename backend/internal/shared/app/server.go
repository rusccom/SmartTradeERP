package app

import (
	"context"
	"net/http"

	"smarterp/backend/internal/features/admin"
	"smarterp/backend/internal/features/clientauth"
	"smarterp/backend/internal/features/documents"
	"smarterp/backend/internal/features/ledger"
	"smarterp/backend/internal/features/products"
	"smarterp/backend/internal/features/reports"
	"smarterp/backend/internal/features/variants"
	"smarterp/backend/internal/features/warehouses"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/config"
	"smarterp/backend/internal/shared/db"
)

func Build(ctx context.Context, cfg config.Config) (*http.Server, func(), error) {
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}
	store := db.NewStore(pool)
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTTL)
	mux := http.NewServeMux()
	registerRoutes(mux, store, tokens)
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: cors(mux)}
	cleanup := func() { pool.Close() }
	return server, cleanup, nil
}

func registerRoutes(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	ledgerService := ledger.NewService(store)
	registerAdmin(mux, store, tokens)
	registerClientAuth(mux, store, tokens)
	registerCatalog(mux, store, tokens, ledgerService)
	registerDocumentsAndReports(mux, store, tokens, ledgerService)
}

func registerAdmin(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	adminRepo := admin.NewRepository(store)
	adminService := admin.NewService(adminRepo, tokens)
	admin.RegisterRoutes(mux, adminService, tokens)
}

func registerClientAuth(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	authRepo := clientauth.NewRepository()
	authService := clientauth.NewService(store, authRepo, tokens)
	clientauth.RegisterRoutes(mux, authService)
}

func registerCatalog(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, ledgerService *ledger.Service) {
	productRepo := products.NewRepository(store)
	productService := products.NewService(store, productRepo, ledgerService)
	products.RegisterRoutes(mux, productService, tokens)
	variantRepo := variants.NewRepository(store)
	variantService := variants.NewService(store, variantRepo, ledgerService)
	variants.RegisterRoutes(mux, variantService, tokens)
	warehouseRepo := warehouses.NewRepository(store)
	warehouseService := warehouses.NewService(warehouseRepo, ledgerService)
	warehouses.RegisterRoutes(mux, warehouseService, tokens)
}

func registerDocumentsAndReports(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	documentsRepo := documents.NewRepository(store)
	documentsService := documents.NewService(store, documentsRepo, ledgerService)
	documents.RegisterRoutes(mux, documentsService, tokens)
	reportsRepo := reports.NewRepository(store)
	reportsService := reports.NewService(reportsRepo, ledgerService)
	reports.RegisterRoutes(mux, reportsService, tokens)
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
