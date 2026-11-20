package api

import (
	"net/http"

	"smarterp/backend/internal/platform/currencies"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/countries"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

func registerSettings(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	repo := currencies.NewRepository(store)
	service := currencies.NewService(store, repo)
	handler := currencies.NewHandler(service)
	handleClient(mux, tokens, "GET /api/client/currencies", handler.List)
	handleClient(mux, tokens, "PUT /api/client/currencies/base", handler.SetBase)
	handleClient(mux, tokens, "GET /api/client/currency-options", handler.Options)
	handleClient(mux, tokens, "GET /api/client/country-options", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteData(w, http.StatusOK, countries.All(), nil)
	})
}
