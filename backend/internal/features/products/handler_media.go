package products

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	mediafeature "smarterp/backend/internal/features/media"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
)

func (h *Handler) ListMedia(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.FromContext(r.Context())
	items, err := h.service.ListMedia(r.Context(), tenantID, r.PathValue("id"))
	if err != nil {
		h.writeMediaError(w, err)
		return
	}
	httpx.WriteData(w, http.StatusOK, items, nil)
}

func (h *Handler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	input := mediafeature.DirectUploadInput{}
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "bad_request", "invalid payload", err.Error())
		return
	}
	tenantID := tenant.FromContext(r.Context())
	upload, err := h.service.CreateMediaUpload(r.Context(), tenantID, r.PathValue("id"), input)
	if err != nil {
		h.writeMediaError(w, err)
		return
	}
	httpx.WriteData(w, http.StatusCreated, upload, nil)
}

func (h *Handler) CompleteMediaUpload(w http.ResponseWriter, r *http.Request) {
	tenantID := tenant.FromContext(r.Context())
	item, err := h.service.CompleteMediaUpload(
		r.Context(), tenantID, r.PathValue("id"), r.PathValue("mediaID"),
	)
	if err != nil {
		h.writeMediaError(w, err)
		return
	}
	httpx.WriteData(w, http.StatusOK, item, nil)
}

func (h *Handler) writeMediaError(w http.ResponseWriter, err error) {
	if errors.Is(err, mediafeature.ErrInvalidMedia) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_media", "invalid product media", nil)
		return
	}
	if errors.Is(err, mediafeature.ErrStorageNotConfigured) {
		httpx.WriteError(w, http.StatusServiceUnavailable, "media_storage_not_configured", "media storage is not configured", nil)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "product not found", nil)
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "media_error", "failed to process product media", err.Error())
}
