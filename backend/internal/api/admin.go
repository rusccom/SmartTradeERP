package api

import (
	"net/http"

	"smarterp/backend/internal/features/admin"
	"smarterp/backend/internal/features/clientauth"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
)

func registerAdmin(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	repo := admin.NewRepository(store)
	service := admin.NewService(repo, tokens)
	handler := admin.NewHandler(service)
	mux.HandleFunc("POST /api/admin/auth/login", handler.Login)
	handleAdmin(mux, tokens, "GET /api/admin/tenants", handler.ListTenants)
	handleAdmin(mux, tokens, "GET /api/admin/tenants/{id}", handler.TenantByID)
	handleAdmin(mux, tokens, "GET /api/admin/stats", handler.Stats)
}

func registerClientAuth(mux *http.ServeMux, store *db.Store, tokens *auth.TokenService) {
	repo := clientauth.NewRepository()
	service := clientauth.NewService(store, repo, tokens)
	handler := clientauth.NewHandler(service)
	mux.HandleFunc("POST /api/client/auth/login", handler.Login)
	mux.HandleFunc("POST /api/client/auth/register", handler.Register)
}
