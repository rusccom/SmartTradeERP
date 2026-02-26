package documents

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
    filters := Filters{Type: r.URL.Query().Get("type")}
    filters.Status = r.URL.Query().Get("status")
    filters.Date = r.URL.Query().Get("date")
    items, total, err := h.service.List(r.Context(), tenantID, filters, page, perPage)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list documents", err.Error())
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
		h.writeDocumentError(w, err, "failed to create document")
		return
	}
	httpx.WriteData(w, http.StatusCreated, map[string]string{"id": id}, nil)
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    item, err := h.service.ByID(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeDocumentError(w, err, "failed to read document")
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
        h.writeDocumentError(w, err, "failed to update document")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Post(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeDocumentError(w, err, "failed to post document")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "posted"}, nil)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Cancel(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeDocumentError(w, err, "failed to cancel document")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "cancelled"}, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    err := h.service.Delete(r.Context(), tenantID, r.PathValue("id"))
    if err != nil {
        h.writeDocumentError(w, err, "failed to delete document")
        return
    }
    httpx.WriteData(w, http.StatusOK, map[string]string{"status": "deleted"}, nil)
}

func (h *Handler) writeDocumentError(w http.ResponseWriter, err error, message string) {
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "document not found", nil)
        return
    }
    if errors.Is(err, ErrPaymentsRequired) {
        httpx.WriteError(w, http.StatusBadRequest, "payments_required", "payments are required for this document type", nil)
        return
    }
    if errors.Is(err, ErrInvalidPaymentMethod) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_payment_method", "payment method must be cash, card or transfer", nil)
        return
    }
    if errors.Is(err, ErrInvalidPaymentAmount) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_payment_amount", "payment amount must be greater than zero", nil)
        return
    }
    if errors.Is(err, ErrPaymentTotalMismatch) {
        httpx.WriteError(w, http.StatusBadRequest, "payment_total_mismatch", "payments total must match document total", nil)
        return
    }
    if errors.Is(err, ErrDraftOnly) {
        httpx.WriteError(w, http.StatusConflict, "draft_only", "operation allowed only for draft", nil)
        return
    }
    if errors.Is(err, ErrPostedOnly) {
        httpx.WriteError(w, http.StatusConflict, "posted_only", "operation allowed only for posted", nil)
        return
    }
    if errors.Is(err, ErrStatusConflict) {
        httpx.WriteError(w, http.StatusConflict, "status_conflict", "invalid document status", nil)
        return
    }
    if errors.Is(err, ErrCompositeWithoutComponents) {
        httpx.WriteError(w, http.StatusConflict, "missing_components", "composite variant has no components", nil)
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}
