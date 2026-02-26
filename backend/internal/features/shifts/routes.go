package shifts

import (
	"net/http"

	"smarterp/backend/internal/shared/auth"
)

func RegisterRoutes(mux *http.ServeMux, service *Service, tokens *auth.TokenService) {
	handler := NewHandler(service)
	chain := []auth.Middleware{auth.RequireScope(tokens, "client")}
	mux.Handle("POST /api/client/shifts/open", auth.Chain(http.HandlerFunc(handler.Open), chain...))
	mux.Handle("GET /api/client/shifts/current", auth.Chain(http.HandlerFunc(handler.Current), chain...))
	mux.Handle("POST /api/client/shifts/cash-op", auth.Chain(http.HandlerFunc(handler.CashOp), chain...))
	mux.Handle("POST /api/client/shifts/close", auth.Chain(http.HandlerFunc(handler.Close), chain...))
	mux.Handle("GET /api/client/shifts/{id}/report", auth.Chain(http.HandlerFunc(handler.Report), chain...))
}
