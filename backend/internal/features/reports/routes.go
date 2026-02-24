package reports

import (
    "net/http"

    "smarterp/backend/internal/shared/auth"
)

func RegisterRoutes(mux *http.ServeMux, service *Service, tokens *auth.TokenService) {
    handler := NewHandler(service)
    chain := []auth.Middleware{auth.RequireScope(tokens, "client")}
    mux.Handle("GET /api/client/reports/profit", auth.Chain(http.HandlerFunc(handler.Profit), chain...))
    mux.Handle("GET /api/client/reports/stock", auth.Chain(http.HandlerFunc(handler.Stock), chain...))
    mux.Handle("GET /api/client/reports/top-products", auth.Chain(http.HandlerFunc(handler.TopProducts), chain...))
    mux.Handle("GET /api/client/reports/movements", auth.Chain(http.HandlerFunc(handler.Movements), chain...))
}
