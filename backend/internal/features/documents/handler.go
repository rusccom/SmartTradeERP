package documents

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
	query := httpx.ParseListQuery(r, httpx.SortConfig{
		Allowed: []string{"date", "number", "total_cost"},
		Fallback: "date",
	}, []string{"type", "status", "date"})
	items, total, err := h.service.List(r.Context(), tenantID, query)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list documents", err.Error())
		return
	}
	meta := &httpx.Meta{Page: query.Page, PerPage: query.PerPage, Total: total}
	httpx.WriteData(w, http.StatusOK, items, meta)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    req := CreateRequest{}
    if err := httpx.DecodeJSON(r, &req); err != nil {
        httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
        return
    }
	tenantID := tenant.FromContext(r.Context())
	result, err := h.service.Create(r.Context(), tenantID, req)
	if err != nil {
		h.writeDocumentError(w, err, "failed to create document")
		return
	}
	httpx.WriteData(w, http.StatusCreated, result, nil)
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
    if writeDocumentRequestError(w, err) {
        return
    }
    if writeDocumentStateError(w, err) {
        return
    }
    if writeDocumentPostingError(w, err) {
        return
    }
    httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}

func writeDocumentRequestError(w http.ResponseWriter, err error) bool {
    if writeDocumentIdentityError(w, err) {
        return true
    }
    return writeDocumentPaymentError(w, err)
}

func writeDocumentIdentityError(w http.ResponseWriter, err error) bool {
    if errors.Is(err, validation.ErrInvalidData) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid document data", nil)
        return true
    }
    if errors.Is(err, ErrInvalidDocumentReference) {
        httpx.WriteError(w, http.StatusBadRequest, "invalid_reference", "invalid document reference", nil)
        return true
    }
    if errors.Is(err, ErrDocumentNumberConflict) {
        httpx.WriteError(w, http.StatusConflict, "number_conflict", "document number already exists", nil)
        return true
    }
    if errors.Is(err, pgx.ErrNoRows) {
        httpx.WriteError(w, http.StatusNotFound, "not_found", "document not found", nil)
        return true
    }
    return false
}

func writeDocumentPaymentError(w http.ResponseWriter, err error) bool {
	if writePaymentPresenceError(w, err) {
		return true
	}
	return writePaymentValueError(w, err)
}

func writePaymentPresenceError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, ErrPaymentsRequired) {
		httpx.WriteError(w, http.StatusBadRequest, "payments_required", "payments are required for this document type", nil)
		return true
	}
	if errors.Is(err, ErrPaymentsNotAllowed) {
		httpx.WriteError(w, http.StatusBadRequest, "payments_not_allowed", "payments are not allowed for this document type", nil)
		return true
	}
	if errors.Is(err, ErrInvalidPaymentMethod) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_payment_method", "payment method must be cash, card or transfer", nil)
		return true
	}
	return false
}

func writePaymentValueError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, ErrInvalidPaymentAmount) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_payment_amount", "payment amount must be greater than zero", nil)
		return true
	}
	if errors.Is(err, ErrPaymentTotalMismatch) {
		httpx.WriteError(w, http.StatusBadRequest, "payment_total_mismatch", "payments total must match document total", nil)
		return true
	}
	return false
}

func writeDocumentStateError(w http.ResponseWriter, err error) bool {
    if errors.Is(err, ErrDraftOnly) {
        httpx.WriteError(w, http.StatusConflict, "draft_only", "operation allowed only for draft", nil)
        return true
    }
    if errors.Is(err, ErrPostedOnly) {
        httpx.WriteError(w, http.StatusConflict, "posted_only", "operation allowed only for posted", nil)
        return true
    }
    if errors.Is(err, ErrStatusConflict) {
        httpx.WriteError(w, http.StatusConflict, "status_conflict", "invalid document status", nil)
        return true
    }
    return false
}

func writeDocumentPostingError(w http.ResponseWriter, err error) bool {
    if errors.Is(err, ErrCompositeWithoutComponents) {
        httpx.WriteError(w, http.StatusConflict, "missing_components", "composite variant has no components", nil)
        return true
    }
    return false
}
