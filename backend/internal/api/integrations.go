package api

import (
	"net/http"

	"smarterp/backend/internal/integrations/basecourier"
	"smarterp/backend/internal/integrations/inpost"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/config"
	"smarterp/backend/internal/shared/db"
)

func registerIntegrations(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, cfg config.Config) {
	registerPlugins(mux, store, tokens)
	registerInPost(mux, store, tokens, cfg.InPostEncryptionKey)
	registerBaseCourier(mux, store, tokens, cfg.BaseCourierEncryptionKey)
	registerAllegro(mux, store, tokens, cfg.Allegro)
}

func registerBaseCourier(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, key string) {
	const path = "/api/client/integrations/base-courier"
	const account = path + "/connections/{version}"
	handler := basecourier.NewHandler(store, key)
	owner := func(route string, fn http.HandlerFunc) {
		mux.Handle(route, auth.Chain(fn, auth.RequireScope(tokens, "client"), auth.RequireRole("owner")))
	}
	handleClient(mux, tokens, "GET "+path, handler.Get)
	owner("PUT "+path, handler.Save)
	owner("POST "+account+"/check", handler.Check)
	owner("DELETE "+account, handler.Disconnect)
}

func registerInPost(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, key string) {
	const path = "/api/client/integrations/inpost"
	handler := inpost.NewHandler(store, key)
	handleClient(mux, tokens, "GET "+path, handler.Get)
	mux.Handle("PUT "+path, auth.Chain(http.HandlerFunc(handler.Save),
		auth.RequireScope(tokens, "client"), auth.RequireRole("owner")))
}
