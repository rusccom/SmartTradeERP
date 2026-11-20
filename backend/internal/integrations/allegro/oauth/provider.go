package oauth

import (
	"context"
	"net/url"

	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/shared/config"
)

const Path = connection.Path + "/oauth"
const dashboardPath = "/dashboard/integrations/allegro"

type Handler struct {
	*connection.Service
	cfg config.Allegro
	app connection.Application
}

func New(s *connection.Service, cfg config.Allegro) *Handler {
	return &Handler{s, cfg, connection.Application{Client: s.Client, ID: cfg.ClientID, Secret: cfg.ClientSecret}}
}

func validOrigin(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host != "" && u.User == nil && u.String() == "https://"+u.Host
}

func (h *Handler) Enabled() bool {
	return h.cfg.ClientID != "" && h.cfg.ClientSecret != "" && h.Vault.AEAD != nil &&
		validOrigin(h.cfg.DashboardOrigin) && validOrigin(h.cfg.APIOrigin)
}

func (h *Handler) callbackURL() string { return h.cfg.APIOrigin + Path + "/callback" }

func (h *Handler) authorizeURL(state, challenge string) string {
	query := url.Values{"response_type": {"code"}, "client_id": {h.cfg.ClientID}, "redirect_uri": {h.callbackURL()},
		"scope": {connection.Scopes}, "state": {state}, "prompt": {"confirm"}, "code_challenge_method": {"S256"}, "code_challenge": {challenge}}
	return "https://allegro.pl/auth/oauth/authorize?" + query.Encode()
}

func (h *Handler) exchange(ctx context.Context, code, verifier string) (connection.Credential, error) {
	form := url.Values{"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {h.callbackURL()},
		"client_id": {h.cfg.ClientID}, "code_verifier": {verifier}}
	return h.app.Grant(ctx, form, false)
}

func (h *Handler) Refresh(ctx context.Context, value string) (connection.Credential, error) {
	if !h.Enabled() {
		return connection.Credential{}, connection.Setup
	}
	return h.app.Refresh(ctx, value)
}
