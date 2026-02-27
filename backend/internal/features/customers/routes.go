package customers

import (
	"net/http"

	"smarterp/backend/internal/shared/auth"
)

func RegisterRoutes(mux *http.ServeMux, service *Service, tokens *auth.TokenService) {
	handler := NewHandler(service)
	chain := []auth.Middleware{auth.RequireScope(tokens, "client")}
	mux.Handle("GET /api/client/customers", auth.Chain(http.HandlerFunc(handler.List), chain...))
	mux.Handle("POST /api/client/customers", auth.Chain(http.HandlerFunc(handler.Create), chain...))
	mux.Handle("GET /api/client/customers/{id}", auth.Chain(http.HandlerFunc(handler.ByID), chain...))
	mux.Handle("PUT /api/client/customers/{id}", auth.Chain(http.HandlerFunc(handler.Update), chain...))
	mux.Handle("DELETE /api/client/customers/{id}", auth.Chain(http.HandlerFunc(handler.Delete), chain...))
}
