package warehouses

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
    includes := warehouseIncludes(r)
    if includes.Stock {
        h.listWithIncludes(w, r, tenantID, includes)
        return
    }
    items, err := h.service.List(r.Context(), tenantID)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list warehouses", err.Error())
        return
    }
    httpx.WriteData(w, http.StatusOK, items, nil)
}

func (h *Handler) listWithIncludes(
	w http.ResponseWriter,
	r *http.Request,
	tenantID string,
	includes WarehouseListInclude,
) {
	items, err := h.service.ListWithIncludes(r.Context(), tenantID, includes)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list warehouses", err.Error())
		return
	}
	httpx.WriteData(w, http.StatusOK, items, nil)
}

func warehouseIncludes(r *http.Request) WarehouseListInclude {
	values := httpx.ParseIncludes(r)
	return WarehouseListInclude{Stock: httpx.HasInclude(values, "stock")}
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
        h.writeWarehouseError(w, err, "failed to create warehouse")
        return
    }
    httpx.WriteData(w, http.StatusCreated, map[string]string{"id": id}, nil)
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
        h.writeWarehouseError(w, err, "failed to update warehouse")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Delete(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeWarehouseError(w, err, "failed to delete warehouse")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "deleted"}, nil)
}

func (h *Handler) writeWarehouseError(w http.ResponseWriter, err error, message string) {
    if errors.Is(err, validation.ErrInvalidData) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid warehouse data", nil)
        return
    }
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "warehouse not found", nil)
        return
    }
    if errors.Is(err, ErrHasMovements) {
        httpx.WriteError(w, http.StatusConflict, "has_movements", "warehouse has ledger movements", nil)
        return
    }
    if errors.Is(err, ErrMustKeepDefault) {
        httpx.WriteError(w, http.StatusConflict, "default_required", "tenant must keep default warehouse", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}
