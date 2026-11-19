package api

import (
	"net/http"

	"smarterp/backend/internal/integrations/inpost"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func registerInPost(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, key string) {
	const path = "/api/client/integrations/inpost"
	handler := inpost.NewHandler(store, key)
	handleClient(mux, tokens, "GET "+path, handler.Get)
	mux.Handle("PUT "+path, auth.Chain(http.HandlerFunc(handler.Save),
		auth.RequireScope(tokens, "client"), auth.RequireRole("owner")))
}
