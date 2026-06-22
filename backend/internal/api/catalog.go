package api

import (
	"net/http"

	"smarterp/backend/internal/erp/bundles"
	"smarterp/backend/internal/directory/customers"
	"smarterp/backend/internal/erp/ledger"
	mediafeature "smarterp/backend/internal/platform/media"
	"smarterp/backend/internal/erp/products"
	"smarterp/backend/internal/erp/variants"
	"smarterp/backend/internal/erp/warehouses"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

type catalogDeps struct {
	store   *db.Store
	tokens  *auth.TokenService
	ledger  *ledger.Service
	bundles *bundles.Service
	media   *mediafeature.Service
}

func registerCatalog(
	mux *http.ServeMux,
	deps catalogDeps,
) {
	registerProducts(mux, deps)
	registerVariants(mux, deps.store, deps.tokens, deps.ledger, deps.bundles)
	registerBundles(mux, deps.tokens, deps.bundles)
	registerWarehouses(mux, deps.store, deps.tokens, deps.ledger)
	registerCustomers(mux, deps.store, deps.tokens)
}

func registerProducts(
	mux *http.ServeMux,
	deps catalogDeps,
) {
	repo := products.NewRepository(deps.store)
	service := products.NewService(deps.store, repo, deps.ledger, deps.bundles)
	service.SetMediaService(deps.media)
	handler := products.NewHandler(service)
	handleClient(mux, deps.tokens, "GET /api/client/products", handler.List)
	handleClient(mux, deps.tokens, "POST /api/client/products", handler.Create)
	handleClient(mux, deps.tokens, "GET /api/client/products/{id}", handler.ByID)
	handleClient(mux, deps.tokens, "PUT /api/client/products/{id}", handler.Update)
	handleClient(mux, deps.tokens, "DELETE /api/client/products/{id}", handler.Delete)
	handleClient(mux, deps.tokens, "GET /api/client/products/{id}/media", handler.ListMedia)
	handleClient(mux, deps.tokens, "POST /api/client/products/{id}/media", handler.UploadMedia)
	handleClient(mux, deps.tokens, "POST /api/client/products/{id}/media/{mediaID}/complete", handler.CompleteMediaUpload)
}

func registerVariants(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
	bundleService *bundles.Service,
) {
	repo := variants.NewRepository(store)
	service := variants.NewService(store, repo, ledgerService, bundleService)
	handler := variants.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/variants", handler.List)
	handleClient(mux, tokens, "POST /api/client/variants", handler.Create)
	handleClient(mux, tokens, "GET /api/client/variants/{id}", handler.ByID)
	handleClient(mux, tokens, "PUT /api/client/variants/{id}", handler.Update)
	handleClient(mux, tokens, "DELETE /api/client/variants/{id}", handler.Delete)
	handleClient(mux, tokens, "GET /api/client/variants/{id}/stock", handler.Stock)
}

func registerBundles(
	mux *http.ServeMux,
	tokens *auth.TokenService,
	service *bundles.Service,
) {
	handler := bundles.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/bundles", handler.List)
	handleClient(mux, tokens, "GET /api/client/bundles/{id}", handler.ByID)
	handleClient(mux, tokens, "GET /api/client/bundles/{id}/components", handler.Components)
	handleClient(mux, tokens, "PUT /api/client/bundles/{id}/components", handler.SetComponents)
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
