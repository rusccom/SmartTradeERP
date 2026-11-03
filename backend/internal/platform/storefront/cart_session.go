package storefront

import (
	"net/http"

	"github.com/google/uuid"
	"smarterp/backend/internal/shared/validation"
)

func cartSession(w http.ResponseWriter, r *http.Request, preview bool) string {
	name := "__Host-sf_cart"
	if preview {
		name = "__Host-sf_draft_cart"
	}
	w.Header().Set("Cache-Control", "no-store")
	cookie, err := r.Cookie(name)
	if err == nil && validation.UUID(cookie.Value) {
		return cookie.Value
	}
	id := uuid.NewString()
	http.SetCookie(w, &http.Cookie{Name: name, Value: id, Path: "/", Secure: true,
		HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: 90 * 24 * 60 * 60})
	return id
}
