package admin

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
        h.writeLoginError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusOK, pair, nil)
}

func (h *Handler) writeLoginError(w http.ResponseWriter, err error) {
    if errors.Is(err, ErrInvalidCredentials) {
        httpx.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "invalid credentials", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "login failed", err.Error())
}

func (h *Handler) ListTenants(w http.ResponseWriter, r *http.Request) {
    page, perPage := httpx.ParsePagination(r)
    data, total, err := h.service.ListTenants(r.Context(), page, perPage)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load tenants", err.Error())
        return
    }
    meta := &httpx.Meta{Page: page, PerPage: perPage, Total: total}
    httpx.WriteData(w, http.StatusOK, data, meta)
}

func (h *Handler) TenantByID(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    item, err := h.service.TenantByID(r.Context(), id)
    if err != nil {
        h.writeTenantError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusOK, item, nil)
}

func (h *Handler) writeTenantError(w http.ResponseWriter, err error) {
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "tenant not found", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load tenant", err.Error())
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
    stats, err := h.service.PlatformStats(r.Context())
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load stats", err.Error())
        return
    }
    httpx.WriteData(w, http.StatusOK, stats, nil)
}
