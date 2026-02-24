package warehouses

import (
    "net/http"

    "smarterp/backend/internal/shared/auth"
)

func RegisterRoutes(mux *http.ServeMux, service *Service, tokens *auth.TokenService) {
    handler := NewHandler(service)
    chain := []auth.Middleware{auth.RequireScope(tokens, "client")}
    mux.Handle("GET /api/client/warehouses", auth.Chain(http.HandlerFunc(handler.List), chain...))
    mux.Handle("POST /api/client/warehouses", auth.Chain(http.HandlerFunc(handler.Create), chain...))
    mux.Handle("PUT /api/client/warehouses/{id}", auth.Chain(http.HandlerFunc(handler.Update), chain...))
    mux.Handle("DELETE /api/client/warehouses/{id}", auth.Chain(http.HandlerFunc(handler.Delete), chain...))
}
