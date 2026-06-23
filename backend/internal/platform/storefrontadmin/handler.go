package storefrontadmin

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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.Get(r.Context(), tenant.FromContext(r.Context()))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to load storefront settings", err.Error())
		return
	}
	httpx.WriteData(w, http.StatusOK, settings, nil)
}

func (h *Handler) Themes(w http.ResponseWriter, r *http.Request) {
	httpx.WriteData(w, http.StatusOK, h.service.Themes(), nil)
}

func (h *Handler) SaveDraft(w http.ResponseWriter, r *http.Request) {
	req := DraftRequest{}
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
		return
	}
	if err := h.service.SaveDraft(r.Context(), tenant.FromContext(r.Context()), req); err != nil {
		h.writeError(w, err, "failed to save storefront draft")
		return
	}
	httpx.WriteData(w, http.StatusOK, map[string]string{"status": "saved"}, nil)
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Publish(r.Context(), tenant.FromContext(r.Context())); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to publish storefront", err.Error())
		return
	}
	httpx.WriteData(w, http.StatusOK, map[string]string{"status": "published"}, nil)
}

func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	preview, err := h.service.Preview(r.Context(), tenant.FromContext(r.Context()))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to create preview", err.Error())
		return
	}
	httpx.WriteData(w, http.StatusOK, preview, nil)
}

func (h *Handler) writeError(w http.ResponseWriter, err error, message string) {
	if errors.Is(err, validation.ErrInvalidData) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_data", "invalid storefront data", nil)
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "internal_error", message, err.Error())
}
