package api

import (
	"net/http"

	"smarterp/backend/internal/features/ledger"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func Register(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	mux.HandleFunc("GET /health", health)
	ledgerService := ledger.NewService(store)
	registerAdmin(mux, store, tokens)
	registerClientAuth(mux, store, tokens)
	registerCatalog(mux, store, tokens, ledgerService)
	registerOperations(mux, store, tokens, ledgerService)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleClient(
	mux *http.ServeMux,
	tokens *auth.TokenService,
	pattern string,
	handler http.HandlerFunc,
) {
	middleware := auth.RequireScope(tokens, "client")
	mux.Handle(pattern, auth.Chain(http.HandlerFunc(handler), middleware))
}

func handleAdmin(
	mux *http.ServeMux,
	tokens *auth.TokenService,
	pattern string,
	handler http.HandlerFunc,
) {
	middleware := auth.RequireScope(tokens, "admin")
	mux.Handle(pattern, auth.Chain(http.HandlerFunc(handler), middleware))
}
