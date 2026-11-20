package api

import (
	"net/http"

	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/integrations/allegro/oauth"
	"smarterp/backend/internal/integrations/allegro/token"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/config"
	"smarterp/backend/internal/shared/db"
)

func registerAllegro(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, cfg config.Allegro) {
	s := connection.NewService(store, cfg.EncryptionKey)
	login := oauth.New(s, cfg)
	h := connection.Handler{Service: s, OAuthReady: login.Enabled, Renew: login.Refresh, CancelPending: oauth.CancelPending}
	owner := func(route string, fn http.HandlerFunc) {
		mux.Handle(route, auth.Chain(fn, auth.RequireScope(tokens, "client"), auth.RequireRole("owner")))
	}
	handleClient(mux, tokens, "GET "+connection.Path, h.Get)
	owner("PUT "+connection.Path+"/token", token.Handler{Service: s}.Save)
	owner("POST "+connection.Path+"/check", h.Check)
	owner("DELETE "+connection.Path+"/connections/{version}", h.Disconnect)
	owner("POST "+oauth.Path+"/connect", login.Start)
	mux.HandleFunc("POST "+oauth.Path+"/launch", login.Launch)
	mux.HandleFunc("GET "+oauth.Path+"/callback", login.Callback)
}
