package auth

import (
    "context"
    "net/http"
    "strings"

    "smarterp/backend/internal/shared/tenant"
)

type key string

type Middleware func(http.Handler) http.Handler

const claimsKey key = "claims"

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
    wrapped := handler
    for i := len(middlewares) - 1; i >= 0; i-- {
        wrapped = middlewares[i](wrapped)
    }
    return wrapped
}

func RequireScope(tokens *TokenService, scope string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, err := parseBearer(r, tokens)
            if err != nil || claims.Scope != scope || claims.TokenType != "access" {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            ctx := WithClaims(r.Context(), claims)
            if claims.TenantID != "" {
                ctx = tenant.WithTenantID(ctx, claims.TenantID)
            }
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func RequireRole(roles ...string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := ClaimsFromContext(r.Context())
            if !isRoleAllowed(claims.Role, roles) {
                w.WriteHeader(http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func WithClaims(ctx context.Context, claims *Claims) context.Context {
    return context.WithValue(ctx, claimsKey, claims)
}

func ClaimsFromContext(ctx context.Context) *Claims {
    value := ctx.Value(claimsKey)
    claims, _ := value.(*Claims)
    if claims == nil {
        return &Claims{}
    }
    return claims
}

func parseBearer(r *http.Request, tokens *TokenService) (*Claims, error) {
    raw := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
    return tokens.Parse(raw)
}

func isRoleAllowed(role string, allowed []string) bool {
    for _, item := range allowed {
        if role == item {
            return true
        }
    }
    return false
}
