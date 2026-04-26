package api

import (
	"net/http"

	"smarterp/backend/internal/features/documents"
	"smarterp/backend/internal/features/ledger"
	"smarterp/backend/internal/features/reports"
	"smarterp/backend/internal/features/shifts"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func registerOperations(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	registerDocuments(mux, store, tokens, ledgerService)
	registerShifts(mux, store, tokens)
	registerReports(mux, store, tokens, ledgerService)
}

func registerDocuments(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	repo := documents.NewRepository(store)
	service := documents.NewService(store, repo, ledgerService)
	handler := documents.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/documents", handler.List)
	handleClient(mux, tokens, "POST /api/client/documents", handler.Create)
	handleClient(mux, tokens, "GET /api/client/documents/{id}", handler.ByID)
	handleClient(mux, tokens, "PUT /api/client/documents/{id}", handler.Update)
	handleClient(mux, tokens, "POST /api/client/documents/{id}/post", handler.Post)
	handleClient(mux, tokens, "POST /api/client/documents/{id}/cancel", handler.Cancel)
	handleClient(mux, tokens, "DELETE /api/client/documents/{id}", handler.Delete)
}

func registerShifts(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	repo := shifts.NewRepository(store)
	service := shifts.NewService(store, repo)
	handler := shifts.NewHandler(service)
	handleClient(mux, tokens, "POST /api/client/shifts/open", handler.Open)
	handleClient(mux, tokens, "GET /api/client/shifts/current", handler.Current)
	handleClient(mux, tokens, "POST /api/client/shifts/cash-op", handler.CashOp)
	handleClient(mux, tokens, "POST /api/client/shifts/close", handler.Close)
	handleClient(mux, tokens, "GET /api/client/shifts/{id}/report", handler.Report)
}

func registerReports(
	mux *http.ServeMux,
	store *db.Store,
	tokens *auth.TokenService,
	ledgerService *ledger.Service,
) {
	repo := reports.NewRepository(store)
	service := reports.NewService(repo, ledgerService)
	handler := reports.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/reports/profit", handler.Profit)
	handleClient(mux, tokens, "GET /api/client/reports/stock", handler.Stock)
	handleClient(mux, tokens, "GET /api/client/reports/top-products", handler.TopProducts)
	handleClient(mux, tokens, "GET /api/client/reports/movements", handler.Movements)
}
