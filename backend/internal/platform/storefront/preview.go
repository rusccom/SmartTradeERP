package storefront

import (
	"net/http"
	"time"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/tenant"
)

const previewScope = "storefront_preview"
const previewCookie = "sf_preview"

type previewAuth struct{ tokens *auth.TokenService }

func newPreviewAuth(secret string) *previewAuth {
	if secret == "" {
		return &previewAuth{}
	}
	return &previewAuth{tokens: auth.NewTokenService(secret, time.Hour)}
}

func (p *previewAuth) valid(raw string, r *http.Request) bool {
	if raw == "" || p.tokens == nil {
		return false
	}
	claims, err := p.tokens.Parse(raw)
	return err == nil && claims.Scope == previewScope && claims.TokenType == "preview" &&
		claims.TenantID == tenant.FromContext(r.Context())
}

func (p *previewAuth) session(w http.ResponseWriter, r *http.Request) (preview, proceed bool) {
	if r.URL.Query().Has("preview") {
		return false, p.enter(w, r)
	}
	cookie, err := r.Cookie(previewCookie)
	preview = err == nil && p.valid(cookie.Value, r)
	if !preview && r.Header.Get("X-Storefront-Preview") == "1" {
		w.Header().Set("Cache-Control", "no-store")
		http.Error(w, "Preview expired", http.StatusForbidden)
		return false, false
	}
	return preview, true
}

func (p *previewAuth) enter(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	token := r.URL.Query().Get("preview")
	if r.Method != http.MethodGet || !p.valid(token, r) {
		http.Error(w, "Invalid or expired preview", http.StatusForbidden)
		return false
	}
	secure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	http.SetCookie(w, &http.Cookie{Name: previewCookie, Value: token, Path: "/", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode})
	target := *r.URL
	query := target.Query()
	query.Del("preview")
	target.RawQuery = query.Encode()
	http.Redirect(w, r, target.RequestURI(), http.StatusSeeOther)
	return false
}
