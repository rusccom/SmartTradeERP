package oauth

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/shared/config"
)

const Path = connection.Path + "/oauth"
const dashboardPath = "/dashboard/integrations/allegro"

type Handler struct {
	*connection.Service
	cfg config.Allegro
}

func New(s *connection.Service, cfg config.Allegro) *Handler { return &Handler{s, cfg} }

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
	return h.tokenRequest(ctx, form, false)
}

func (h *Handler) Refresh(ctx context.Context, value string) (connection.Credential, error) {
	if !h.Enabled() {
		return connection.Credential{}, connection.Setup
	}
	return h.tokenRequest(ctx, url.Values{"grant_type": {"refresh_token"}, "refresh_token": {value}}, true)
}

func (h *Handler) tokenRequest(ctx context.Context, form url.Values, basic bool) (connection.Credential, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://allegro.pl/auth/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return connection.Credential{}, connection.Remote
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if basic {
		req.SetBasicAuth(h.cfg.ClientID, h.cfg.ClientSecret)
	}
	return h.readCredential(req)
}

func (h *Handler) readCredential(req *http.Request) (connection.Credential, error) {
	var pair struct {
		connection.Tokens
		Seconds int64 `json:"expires_in"`
	}
	if err := h.Client.Request(req, &pair); err != nil {
		return connection.Credential{}, err
	}
	if pair.Access == "" || pair.Refresh == "" || pair.Seconds <= 0 || pair.Seconds > 365*24*3600 {
		return connection.Credential{}, connection.Remote
	}
	expires := time.Now().Add(time.Duration(pair.Seconds) * time.Second)
	return connection.Credential{Tokens: pair.Tokens, ExpiresAt: &expires}, nil
}
