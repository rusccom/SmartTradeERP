package clientauth

import "net/http"

func RegisterRoutes(mux *http.ServeMux, service *Service) {
    handler := NewHandler(service)
    mux.HandleFunc("POST /api/client/auth/login", handler.Login)
    mux.HandleFunc("POST /api/client/auth/register", handler.Register)
}
