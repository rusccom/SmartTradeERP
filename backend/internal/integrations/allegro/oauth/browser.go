package oauth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"smarterp/backend/internal/integrations/allegro/connection"
)

const flowCookie = "__Host-smarttrade-allegro"
const cookieBinding = "allegro:browser:v1"

func browserHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func setFlowCookie(w http.ResponseWriter, value string, age int) {
	http.SetCookie(w, &http.Cookie{Name: flowCookie, Value: value, Path: "/", MaxAge: age,
		Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}

func (h *Handler) Launch(w http.ResponseWriter, r *http.Request) {
	browserHeaders(w)
	if !h.Enabled() || r.Header.Get("Origin") != h.cfg.DashboardOrigin {
		h.redirect(w, r, connection.State)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	f, err := h.launch(r.Context(), r.PostFormValue("ticket"))
	if err != nil {
		h.redirect(w, r, err)
		return
	}
	encoded, err := h.Vault.Encrypt(cookieBinding, f)
	if err != nil {
		h.redirect(w, r, err)
		return
	}
	setFlowCookie(w, base64.RawURLEncoding.EncodeToString(encoded), lifetime)
	http.Redirect(w, r, h.authorizeURL(f.State, f.challenge()), http.StatusSeeOther)
}

func (h *Handler) readFlow(r *http.Request) (flow, error) {
	var f flow
	cookie, err := r.Cookie(flowCookie)
	if err != nil {
		return f, connection.State
	}
	raw, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil || h.Vault.Decrypt(cookieBinding, raw, &f) != nil {
		return f, connection.State
	}
	if subtle.ConstantTimeCompare([]byte(r.URL.Query().Get("state")), []byte(f.State)) != 1 {
		return f, connection.State
	}
	return f, nil
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	browserHeaders(w)
	f, err := h.readFlow(r)
	if err != nil {
		h.redirect(w, r, err)
		return
	}
	setFlowCookie(w, "", -1)
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	h.redirect(w, r, h.complete(ctx, r.URL.Query(), f))
}

func (h *Handler) complete(ctx context.Context, values url.Values, f flow) error {
	code := values.Get("code")
	if values.Get("error") != "" || code == "" {
		if err := h.finish(ctx, f, ""); err != nil {
			return err
		}
		if values.Get("error") == "access_denied" {
			return connection.Denied
		}
		return connection.State
	}
	if len(code) > 4096 {
		return connection.State
	}
	return h.finish(ctx, f, code)
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request, err error) {
	if !validOrigin(h.cfg.DashboardOrigin) {
		connection.Respond(w, nil, connection.Setup)
		return
	}
	result := "connected"
	if err != nil {
		result = connection.PublicError(err).Code
	}
	http.Redirect(w, r, h.cfg.DashboardOrigin+dashboardPath+"?result="+result, http.StatusSeeOther)
}
