package shifts

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/auth"
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

func (h *Handler) Open(w http.ResponseWriter, r *http.Request) {
	req := OpenRequest{}
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
		return
	}
	tenantID, userID := tenantAndUser(r)
	id, err := h.service.Open(r.Context(), tenantID, userID, req)
	if err != nil {
		h.writeShiftError(w, err, "failed to open shift")
		return
	}
	httpx.WriteData(w, http.StatusCreated, map[string]string{"id": id}, nil)
}

func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	tenantID, userID := tenantAndUser(r)
	item, err := h.service.Current(r.Context(), tenantID, userID)
	if err != nil {
		h.writeShiftError(w, err, "failed to load current shift")
		return
	}
	httpx.WriteData(w, http.StatusOK, item, nil)
}

func (h *Handler) CashOp(w http.ResponseWriter, r *http.Request) {
	req := CashOpRequest{}
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
		return
	}
	tenantID, userID := tenantAndUser(r)
	err := h.service.CashOp(r.Context(), tenantID, userID, req)
	if err != nil {
		h.writeShiftError(w, err, "failed to apply cash operation")
		return
	}
	httpx.WriteData(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
}

func (h *Handler) Close(w http.ResponseWriter, r *http.Request) {
	tenantID, userID := tenantAndUser(r)
	report, err := h.service.Close(r.Context(), tenantID, userID)
	if err != nil {
		h.writeShiftError(w, err, "failed to close shift")
		return
	}
	httpx.WriteData(w, http.StatusOK, report, nil)
}

func (h *Handler) Report(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.FromContext(r.Context())
	report, err := h.service.Report(r.Context(), tenantID, r.PathValue("id"))
	if err != nil {
		h.writeShiftError(w, err, "failed to load shift report")
		return
	}
	httpx.WriteData(w, http.StatusOK, report, nil)
}

func tenantAndUser(r *http.Request) (string, string) {
	tenantID := tenant.FromContext(r.Context())
	userID := auth.ClaimsFromContext(r.Context()).UserID
	return tenantID, userID
}

func (h *Handler) writeShiftError(w http.ResponseWriter, err error, message string) {
	if writeShiftRequestError(w, err) {
		return
	}
	if writeShiftStateError(w, err) {
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}

func writeShiftRequestError(w http.ResponseWriter, err error) bool {
	if writeShiftIdentityError(w, err) {
		return true
	}
	return writeShiftCashError(w, err)
}

func writeShiftIdentityError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, validation.ErrInvalidData) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid shift data", nil)
		return true
	}
	if errors.Is(err, ErrInvalidShiftReference) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_reference", "invalid shift reference", nil)
		return true
	}
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "shift not found", nil)
		return true
	}
	return false
}

func writeShiftCashError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, ErrInvalidCashOpType) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_cash_op_type", ErrInvalidCashOpType.Error(), nil)
		return true
	}
	if errors.Is(err, ErrInvalidAmount) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_amount", ErrInvalidAmount.Error(), nil)
		return true
	}
	return false
}

func writeShiftStateError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, ErrShiftAlreadyOpen) {
		httpx.WriteError(w, http.StatusConflict, "shift_already_open", ErrShiftAlreadyOpen.Error(), nil)
		return true
	}
	if errors.Is(err, ErrNoOpenShift) {
		httpx.WriteError(w, http.StatusConflict, "no_open_shift", ErrNoOpenShift.Error(), nil)
		return true
	}
	if errors.Is(err, ErrShiftAlreadyClosed) {
		httpx.WriteError(w, http.StatusConflict, "shift_already_closed", ErrShiftAlreadyClosed.Error(), nil)
		return true
	}
	return false
}
