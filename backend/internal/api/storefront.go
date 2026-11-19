package api

import (
	"net/http"

	"smarterp/backend/internal/platform/storefront/themes"
	"smarterp/backend/internal/platform/storefrontadmin"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

// registerStorefrontAdmin wires the tenant-facing storefront settings API. It loads theme descriptors for the picker.
func registerStorefrontAdmin(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, shopsBase string) error {
	registry, err := themes.Load()
	if err != nil {
		return err
	}
	service := storefrontadmin.NewService(store, storefrontadmin.NewRepository(store), registry, tokens, shopsBase)
	handler := storefrontadmin.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/storefront/settings", handler.Get)
	handleClient(mux, tokens, "GET /api/client/storefront/themes", handler.Themes)
	handleClient(mux, tokens, "PUT /api/client/storefront/draft", handler.SaveDraft)
	handleClient(mux, tokens, "POST /api/client/storefront/publish", handler.Publish)
	handleClient(mux, tokens, "GET /api/client/storefront/preview", handler.Preview)
	return nil
}
