package variants

import (
    "net/http"

    "smarterp/backend/internal/shared/auth"
)

func RegisterRoutes(mux *http.ServeMux, service *Service, tokens *auth.TokenService) {
    handler := NewHandler(service)
    chain := []auth.Middleware{auth.RequireScope(tokens, "client")}
    mux.Handle("GET /api/client/variants", auth.Chain(http.HandlerFunc(handler.List), chain...))
    mux.Handle("POST /api/client/variants", auth.Chain(http.HandlerFunc(handler.Create), chain...))
    mux.Handle("GET /api/client/variants/{id}", auth.Chain(http.HandlerFunc(handler.ByID), chain...))
    mux.Handle("PUT /api/client/variants/{id}", auth.Chain(http.HandlerFunc(handler.Update), chain...))
    mux.Handle("DELETE /api/client/variants/{id}", auth.Chain(http.HandlerFunc(handler.Delete), chain...))
    mux.Handle("GET /api/client/variants/{id}/components", auth.Chain(http.HandlerFunc(handler.Components), chain...))
    mux.Handle("PUT /api/client/variants/{id}/components", auth.Chain(http.HandlerFunc(handler.SetComponents), chain...))
    mux.Handle("GET /api/client/variants/{id}/stock", auth.Chain(http.HandlerFunc(handler.Stock), chain...))
}
