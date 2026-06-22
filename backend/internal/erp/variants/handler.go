package variants

import (
    "errors"
    "net/http"

    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/shared/httpx"
    "smarterp/backend/internal/shared/tenant"
    "smarterp/backend/internal/shared/validation"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    page, perPage := httpx.ParsePagination(r)
    productID := r.URL.Query().Get("product_id")
    items, total, err := h.service.List(r.Context(), tenantID, productID, page, perPage)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list variants", err.Error())
        return
    }
    meta := &httpx.Meta{Page: page, PerPage: perPage, Total: total}
    httpx.WriteData(w, http.StatusOK, items, meta)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    req := CreateRequest{}
    if err := httpx.DecodeJSON(r, &req); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
    tenantID := tenant.FromContext(r.Context())
    id, err := h.service.Create(r.Context(), tenantID, req)
    if err != nil {
        h.writeMutationError(w, err, "failed to create variant")
        return
    }
    httpx.WriteData(w, http.StatusCreated, map[string]string{"id": id}, nil)
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    item, err := h.service.ByID(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeNotFoundAwareError(w, err, "failed to read variant")
        return
    }
    httpx.WriteData(w, http.StatusOK, item, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    req := UpdateRequest{}
    if err := httpx.DecodeJSON(r, &req); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Update(r.Context(), tenantID, r.PathValue("id"), req)
    if err != nil {
        h.writeMutationError(w, err, "failed to update variant")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Delete(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeDeleteError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "deleted"}, nil)
}

func (h *Handler) Stock(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    stock, err := h.service.Stock(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load stock", err.Error())
        return
    }
    httpx.WriteData(w, http.StatusOK, stock, nil)
}

func (h *Handler) writeNotFoundAwareError(w http.ResponseWriter, err error, message string) {
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "variant not found", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}

func (h *Handler) writeDeleteError(w http.ResponseWriter, err error) {
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "variant not found", nil)
        return
    }
    if errors.Is(err, ErrHasMovements) {
        httpx.WriteError(w, http.StatusConflict, "has_movements", "variant has inventory movements", nil)
        return
    }
    if errors.Is(err, ErrLastVariant) {
        httpx.WriteError(w, http.StatusConflict, "last_variant", "product must keep one variant", nil)
        return
    }
    if errors.Is(err, ErrUsedInBundle) {
        httpx.WriteError(w, http.StatusConflict, "used_in_bundle", "variant is used in bundle", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to delete variant", err.Error())
}

func (h *Handler) writeMutationError(w http.ResponseWriter, err error, message string) {
    if errors.Is(err, validation.ErrInvalidData) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid variant data", nil)
        return
    }
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "variant not found", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}
