package api

import (
	"net/http"

	"smarterp/backend/internal/platform/storefront"
	"smarterp/backend/internal/platform/storefrontadmin"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

// registerStorefrontAdmin wires the tenant-facing storefront settings API. It
// loads the embedded theme registry once for the picker and token validation.
func registerStorefrontAdmin(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) error {
	registry, err := storefront.NewRegistry()
	if err != nil {
		return err
	}
	service := storefrontadmin.NewService(store, storefrontadmin.NewRepository(store), registry, tokens)
	handler := storefrontadmin.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/storefront/settings", handler.Get)
	handleClient(mux, tokens, "GET /api/client/storefront/themes", handler.Themes)
	handleClient(mux, tokens, "PUT /api/client/storefront/draft", handler.SaveDraft)
	handleClient(mux, tokens, "POST /api/client/storefront/publish", handler.Publish)
	handleClient(mux, tokens, "GET /api/client/storefront/preview", handler.Preview)
	return nil
}
