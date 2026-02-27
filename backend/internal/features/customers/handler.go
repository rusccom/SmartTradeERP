package customers

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
	query := httpx.ParseListQuery(r, httpx.SortConfig{
		Allowed:  []string{"name", "created_at"},
		Fallback: "created_at",
	}, nil)
	data, total, err := h.service.List(r.Context(), tenantID, query)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list customers", err.Error())
		return
	}
	meta := &httpx.Meta{Page: query.Page, PerPage: query.PerPage, Total: total}
	httpx.WriteData(w, http.StatusOK, data, meta)
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
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to create customer", err.Error())
		return
	}
	httpx.WriteData(w, http.StatusCreated, map[string]string{"id": id}, nil)
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.FromContext(r.Context())
	item, err := h.service.ByID(r.Context(), tenantID, r.PathValue("id"))
	if err != nil {
		h.writeQueryError(w, err)
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
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to update customer", err.Error())
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

func (h *Handler) writeQueryError(w http.ResponseWriter, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "customer not found", nil)
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to read customer", err.Error())
}

func (h *Handler) writeDeleteError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrIsDefault) {
		httpx.WriteError(w, http.StatusConflict, "is_default", "cannot delete default customer", nil)
		return
	}
	if errors.Is(err, ErrHasDocuments) {
		httpx.WriteError(w, http.StatusConflict, "has_documents", "customer has linked documents", nil)
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to delete customer", err.Error())
}
