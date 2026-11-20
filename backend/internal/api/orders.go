package api

import (
	"net/http"

	"smarterp/backend/internal/erp/orders"
	"smarterp/backend/internal/platform/media"
	"smarterp/backend/internal/regions/pl"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func registerOrders(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService, key string, images *media.Service) {
	h := orders.Handler{Store: store, Media: images}
	handleClient(mux, tokens, "GET /api/client/orders", h.List)
	handleClient(mux, tokens, "POST /api/client/orders", h.Save)
	handleClient(mux, tokens, "POST /api/client/orders/quote", h.Quote)
	handleClient(mux, tokens, "GET /api/client/orders/{id}", h.Get)
	handleClient(mux, tokens, "PUT /api/client/orders/{id}", h.Save)
	handleClient(mux, tokens, "POST /api/client/orders/{id}/quote", h.Quote)
	handleClient(mux, tokens, "POST /api/client/orders/{id}/payments", h.Payment)
	registerPolishInvoices(mux, tokens, pl.New(store, key))
}

func registerPolishInvoices(mux *http.ServeMux, tokens *auth.TokenService, h *pl.Handler) {
	handleClient(mux, tokens, "GET /api/client/orders/invoice-formats", h.Formats)
	handleClient(mux, tokens, "POST /api/client/orders/{id}/invoice", h.Issue)
	handleClient(mux, tokens, "GET /api/client/regions/pl/settings", h.Settings)
	mux.Handle("PUT /api/client/regions/pl/settings", auth.Chain(http.HandlerFunc(h.SaveSettings), auth.RequireScope(tokens, "client"), auth.RequireRole("owner")))
	handleClient(mux, tokens, "GET /api/client/regions/pl/invoices/{id}/ksef", h.Delivery)
	handleClient(mux, tokens, "POST /api/client/regions/pl/invoices/{id}/ksef", h.Send)
	handleClient(mux, tokens, "POST /api/client/regions/pl/invoices/{id}/ksef/check", h.Check)
}
