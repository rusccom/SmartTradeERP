package smartbundle

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
	"smarterp/backend/internal/shared/validation"
)

type SettingsHandler struct{ store *db.Store }

func NewSettingsHandler(store *db.Store) *SettingsHandler { return &SettingsHandler{store: store} }

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) { h.respond(w, r, nil) }

func (h *SettingsHandler) Put(w http.ResponseWriter, r *http.Request) {
	var settings Settings
	if err := httpx.DecodeJSON(r, &settings); err != nil {
		httpx.WriteError(w, 400, "bad_request", "Invalid settings", nil)
		return
	}
	h.respond(w, r, &settings)
}

func (h *SettingsHandler) respond(w http.ResponseWriter, r *http.Request, s *Settings) {
	result, err := h.editor(r.Context(), tenant.FromContext(r.Context()), r.PathValue("id"), s)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		httpx.WriteError(w, 404, "not_found", "Bundle not found", nil)
	case errors.Is(err, validation.ErrInvalidData):
		httpx.WriteError(w, 400, "invalid_data", "Invalid bundle rules", nil)
	case err != nil:
		httpx.WriteError(w, 500, "internal_error", "Cannot load or save bundle settings", nil)
	default:
		httpx.WriteData(w, 200, result, nil)
	}
}
