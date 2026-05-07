package currencies

import (
	"errors"
	"net/http"

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
	items, total, err := h.service.List(r.Context(), tenantID, page, perPage)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list currencies", err.Error())
		return
	}
	meta := &httpx.Meta{Page: page, PerPage: perPage, Total: total}
	httpx.WriteData(w, http.StatusOK, items, meta)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	page, perPage := httpx.ParsePagination(r)
	items, total, err := h.service.Options(r.Context(), page, perPage)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list currency options", err.Error())
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
		h.writeCurrencyError(w, err, "failed to create currency")
		return
	}
	httpx.WriteData(w, http.StatusCreated, map[string]string{"id": id}, nil)
}

func (h *Handler) SetBase(w http.ResponseWriter, r *http.Request) {
	req := BaseRequest{}
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
		return
	}
	tenantID := tenant.FromContext(r.Context())
	if err := h.service.SetBase(r.Context(), tenantID, req); err != nil {
		h.writeCurrencyError(w, err, "failed to set base currency")
		return
	}
	httpx.WriteData(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
}

func (h *Handler) writeCurrencyError(w http.ResponseWriter, err error, message string) {
	if errors.Is(err, validation.ErrInvalidData) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid currency data", nil)
		return
	}
	if errors.Is(err, ErrCurrencyExists) {
		httpx.WriteError(w, http.StatusConflict, "currency_exists", "currency already exists", nil)
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}
