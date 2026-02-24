package variants

import (
    "errors"
    "net/http"

    "github.com/jackc/pgx/v5"

    "smarterp/backend/internal/shared/httpx"
    "smarterp/backend/internal/shared/tenant"
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
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to create variant", err.Error())
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
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to update variant", err.Error())
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

func (h *Handler) Components(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    items, err := h.service.Components(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeNotFoundAwareError(w, err, "failed to read components")
        return
    }
    httpx.WriteData(w, http.StatusOK, items, nil)
}

func (h *Handler) SetComponents(w http.ResponseWriter, r *http.Request) {
    payload := struct {
        Components []Component `json:"components"`
    }{}
    if err := httpx.DecodeJSON(r, &payload); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
    tenantID := tenant.FromContext(r.Context())
    err := h.service.SetComponents(r.Context(), tenantID, r.PathValue("id"), payload.Components)
    if err != nil {
        h.writeComponentsError(w, err)
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
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
    if errors.Is(err, ErrHasMovements) {
        httpx.WriteError(w, http.StatusConflict, "has_movements", "variant has ledger movements", nil)
        return
    }
    if errors.Is(err, ErrLastVariant) {
        httpx.WriteError(w, http.StatusConflict, "last_variant", "product must keep one variant", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to delete variant", err.Error())
}

func (h *Handler) writeComponentsError(w http.ResponseWriter, err error) {
    if errors.Is(err, ErrInvalidComponentState) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_components", "invalid component state", nil)
        return
    }
    h.writeNotFoundAwareError(w, err, "failed to update components")
}
