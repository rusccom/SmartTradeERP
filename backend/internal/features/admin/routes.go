package admin

import (
    "net/http"

    "smarterp/backend/internal/shared/auth"
)

func RegisterRoutes(mux *http.ServeMux, service *Service, tokens *auth.TokenService) {
    handler := NewHandler(service)
    mux.HandleFunc("POST /api/admin/auth/login", handler.Login)
    registerProtected(mux, handler, tokens)
}

func registerProtected(mux *http.ServeMux, handler *Handler, tokens *auth.TokenService) {
    middleware := []auth.Middleware{auth.RequireScope(tokens, "admin"), auth.RequireRole("owner")}
    mux.Handle("GET /api/admin/tenants", auth.Chain(http.HandlerFunc(handler.ListTenants), middleware...))
    mux.Handle("GET /api/admin/tenants/{id}", auth.Chain(http.HandlerFunc(handler.TenantByID), middleware...))
    mux.Handle("GET /api/admin/stats", auth.Chain(http.HandlerFunc(handler.Stats), middleware...))
}
