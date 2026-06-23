package storefront

import (
	"net/http"
	"strings"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/tenant"
)

type Resolver struct {
	service *Service
}

func NewResolver(service *Service) *Resolver {
	return &Resolver{service: service}
}

// Middleware resolves the request Host to a tenant and stores it on the
// context. It fails closed (404) when no active storefront matches and never
// falls back to a default tenant. Only r.Host is trusted: Cloudflare forwards
// the original requested host, so no client-set header can spoof the tenant.
func (rs *Resolver) Middleware() auth.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resolved, err := rs.service.ResolveHost(r.Context(), normalizeHost(r.Host))
			if err != nil {
				writeNotFound(w)
				return
			}
			ctx := tenant.WithTenantID(r.Context(), resolved.TenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// normalizeHost lower-cases the host and strips any port and trailing dot so
// lookups match the stored host value regardless of request formatting.
func normalizeHost(raw string) string {
	host := strings.ToLower(strings.TrimSpace(raw))
	if i := strings.LastIndex(host, ":"); i != -1 {
		host = host[:i]
	}
	return strings.TrimSuffix(host, ".")
}
