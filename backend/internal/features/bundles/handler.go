package bundles

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
	items, total, err := h.service.List(r.Context(), tenantID, page, perPage)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list bundles", err.Error())
		return
	}
	meta := &httpx.Meta{Page: page, PerPage: perPage, Total: total}
	httpx.WriteData(w, http.StatusOK, items, meta)
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.FromContext(r.Context())
	item, err := h.service.ByID(r.Context(), tenantID, r.PathValue("id"))
	if err != nil {
		h.writeError(w, err, "failed to read bundle")
		return
	}
	httpx.WriteData(w, http.StatusOK, item, nil)
}

func (h *Handler) Components(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.FromContext(r.Context())
	items, err := h.service.Components(r.Context(), tenantID, r.PathValue("id"))
	if err != nil {
		h.writeError(w, err, "failed to read components")
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
	h.updateComponents(w, r, payload.Components)
}

func (h *Handler) updateComponents(w http.ResponseWriter, r *http.Request, items []Component) {
	tenantID := tenant.FromContext(r.Context())
	err := h.service.SetComponents(r.Context(), tenantID, r.PathValue("id"), items)
	if err != nil {
		h.writeError(w, err, "failed to update components")
		return
	}
	httpx.WriteData(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
}

func (h *Handler) writeError(w http.ResponseWriter, err error, message string) {
	if writeRequestError(w, err) {
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "bundle not found", nil)
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}

func writeRequestError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, validation.ErrInvalidData) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_components", "invalid components", nil)
		return true
	}
	if errors.Is(err, ErrInvalidComponentState) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_component_state", "invalid component state", nil)
		return true
	}
	if errors.Is(err, ErrMissingComponents) {
		httpx.WriteError(w, http.StatusConflict, "missing_components", "bundle has no components", nil)
		return true
	}
	return false
}
