package api

import (
	"net/http"

	"smarterp/backend/internal/features/customers"
	"smarterp/backend/internal/features/ledger"
	"smarterp/backend/internal/features/products"
	"smarterp/backend/internal/features/variants"
	"smarterp/backend/internal/features/warehouses"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func registerCatalog(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	registerProducts(mux, store, tokens, ledgerService)
	registerVariants(mux, store, tokens, ledgerService)
	registerWarehouses(mux, store, tokens, ledgerService)
	registerCustomers(mux, store, tokens)
}

func registerProducts(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	repo := products.NewRepository(store)
	service := products.NewService(store, repo, ledgerService)
	handler := products.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/products", handler.List)
	handleClient(mux, tokens, "POST /api/client/products", handler.Create)
	handleClient(mux, tokens, "GET /api/client/products/{id}", handler.ByID)
	handleClient(mux, tokens, "PUT /api/client/products/{id}", handler.Update)
	handleClient(mux, tokens, "DELETE /api/client/products/{id}", handler.Delete)
}

func registerVariants(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	repo := variants.NewRepository(store)
	service := variants.NewService(store, repo, ledgerService)
	handler := variants.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/variants", handler.List)
	handleClient(mux, tokens, "POST /api/client/variants", handler.Create)
	handleClient(mux, tokens, "GET /api/client/variants/{id}", handler.ByID)
	handleClient(mux, tokens, "PUT /api/client/variants/{id}", handler.Update)
	handleClient(mux, tokens, "DELETE /api/client/variants/{id}", handler.Delete)
	handleClient(mux, tokens, "GET /api/client/variants/{id}/components", handler.Components)
	handleClient(mux, tokens, "PUT /api/client/variants/{id}/components", handler.SetComponents)
	handleClient(mux, tokens, "GET /api/client/variants/{id}/stock", handler.Stock)
}

func registerWarehouses(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	repo := warehouses.NewRepository(store)
	service := warehouses.NewService(repo, ledgerService)
	handler := warehouses.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/warehouses", handler.List)
	handleClient(mux, tokens, "POST /api/client/warehouses", handler.Create)
	handleClient(mux, tokens, "PUT /api/client/warehouses/{id}", handler.Update)
	handleClient(mux, tokens, "DELETE /api/client/warehouses/{id}", handler.Delete)
}

func registerCustomers(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	repo := customers.NewRepository(store)
	service := customers.NewService(repo)
	handler := customers.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/customers", handler.List)
	handleClient(mux, tokens, "POST /api/client/customers", handler.Create)
	handleClient(mux, tokens, "GET /api/client/customers/{id}", handler.ByID)
	handleClient(mux, tokens, "PUT /api/client/customers/{id}", handler.Update)
	handleClient(mux, tokens, "DELETE /api/client/customers/{id}", handler.Delete)
}
