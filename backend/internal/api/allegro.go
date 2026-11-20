package api

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/integrations/allegro/oauth"
	"smarterp/backend/internal/integrations/allegro/personal"
	"smarterp/backend/internal/integrations/allegro/token"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/config"
	"smarterp/backend/internal/shared/db"
)

func registerAllegro(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, cfg config.Allegro) {
	s := connection.NewService(store, cfg.EncryptionKey)
	login := oauth.New(s, cfg)
	private := personal.New(s, cfg.Private)
	h := connection.Handler{Service: s, OAuthReady: login.Enabled, Renew: login.Refresh, Private: private.Renewal(), CancelPending: cancelAllegro}
	owner := func(route string, fn http.HandlerFunc) {
		mux.Handle(route, auth.Chain(fn, auth.RequireScope(tokens, "client"), auth.RequireRole("owner")))
	}
	handleClient(mux, tokens, "GET "+connection.Path, h.Get)
	owner("PUT "+connection.Path+"/token", token.Handler{Service: s}.Save)
	owner("PUT "+personal.Path+"/refresh-token", private.Import)
	owner("POST "+personal.Path+"/device", private.Start)
	owner("POST "+personal.Path+"/device/check", private.Poll)
	owner("POST "+connection.Path+"/check", h.Check)
	owner("DELETE "+connection.Path+"/connections/{version}", h.Disconnect)
	owner("POST "+oauth.Path+"/connect", login.Start)
	mux.HandleFunc("POST "+oauth.Path+"/launch", login.Launch)
	mux.HandleFunc("GET "+oauth.Path+"/callback", login.Callback)
}

func cancelAllegro(ctx context.Context, tx pgx.Tx, tenant string) error {
	if err := personal.Cancel(ctx, tx, tenant); err != nil {
		return err
	}
	return oauth.CancelPending(ctx, tx, tenant)
}
