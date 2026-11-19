package smartbundle

import (
	"net/http"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func RegisterSettings(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	handler := NewSettingsHandler(store)
	client := auth.RequireScope(tokens, "client")
	mux.Handle("GET /api/client/bundles/{id}/storefront", client(http.HandlerFunc(handler.Get)))
	mux.Handle("PUT /api/client/bundles/{id}/storefront", client(http.HandlerFunc(handler.Put)))
}
