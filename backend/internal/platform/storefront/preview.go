package storefront

import (
	"net/http"
	"time"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/tenant"
)

const previewScope = "storefront_preview"

// previewAuth verifies the ?preview=<token> query against the shared JWT secret
// so a store owner can view their unpublished draft on real catalog data.
type previewAuth struct {
	tokens *auth.TokenService
}

func newPreviewAuth(secret string) *previewAuth {
	if secret == "" {
		return &previewAuth{}
	}
	return &previewAuth{tokens: auth.NewTokenService(secret, time.Hour)}
}

// check reports whether the request carries a valid preview token for the
// tenant resolved from the Host.
func (p *previewAuth) check(r *http.Request) bool {
	if p.tokens == nil {
		return false
	}
	raw := r.URL.Query().Get("preview")
	if raw == "" {
		return false
	}
	claims, err := p.tokens.Parse(raw)
	if err != nil {
		return false
	}
	return claims.Scope == previewScope && claims.TokenType == "preview" &&
		claims.TenantID == tenant.FromContext(r.Context())
}
