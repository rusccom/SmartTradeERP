package clientauth

import (
    "errors"
    "net/http"

    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/shared/httpx"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    req := LoginRequest{}
    if err := httpx.DecodeJSON(r, &req); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
    pair, err := h.service.Login(r.Context(), req)
    if err != nil {
        h.writeAuthError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusOK, pair, nil)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    req := RegisterRequest{}
    if err := httpx.DecodeJSON(r, &req); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
    pair, err := h.service.Register(r.Context(), req)
    if err != nil {
        h.writeAuthError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusCreated, pair, nil)
}

func (h *Handler) writeAuthError(w http.ResponseWriter, err error) {
    if errors.Is(err, ErrInvalidCredentials) {
        httpx.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "invalid credentials", nil)
        return
    }
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "invalid credentials", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "auth failed", err.Error())
}
