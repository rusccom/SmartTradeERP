package api

import (
	"net/http"

	"smarterp/backend/internal/platform/plugins"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func registerPlugins(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	handler := plugins.NewHandler(plugins.NewRepository(store))
	handleAdmin(mux, tokens, "GET /api/admin/tenants/{id}/plugins", handler.Admin)
	handleAdmin(mux, tokens, "PUT /api/admin/tenants/{id}/plugins/{key}", handler.Grant)
	handleClient(mux, tokens, "GET /api/client/plugins", handler.Client)
	mux.Handle("PUT /api/client/plugins/{key}", auth.Chain(http.HandlerFunc(handler.Activate),
		auth.RequireScope(tokens, "client"), auth.RequireRole("owner")))
}
