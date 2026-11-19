package plugins

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
	"smarterp/backend/internal/shared/validation"
)

type Handler struct{ repo *Repository }

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Client(w http.ResponseWriter, r *http.Request) {
	h.list(w, r, tenant.FromContext(r.Context()), auth.ClaimsFromContext(r.Context()).Role == "owner")
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	h.list(w, r, r.PathValue("id"), true)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request, id string, manage bool) {
	w.Header().Set("Cache-Control", "private, no-store")
	if !validation.UUID(id) {
		respond(w, validation.ErrInvalidData)
		return
	}
	states, err := h.repo.List(r.Context(), id)
	if err != nil {
		respond(w, err)
		return
	}
	httpx.WriteData(w, 200, map[string]any{"plugins": states, "can_manage": manage}, nil)
}

func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Enabled *bool `json:"enabled"`
	}
	if httpx.DecodeJSON(r, &input) != nil || input.Enabled == nil {
		respond(w, validation.ErrInvalidData)
		return
	}
	respond(w, h.repo.Enable(r.Context(), tenant.FromContext(r.Context()), r.PathValue("key"), *input.Enabled))
}

func (h *Handler) Grant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input struct {
		Granted *bool `json:"granted"`
	}
	if !validation.UUID(id) || httpx.DecodeJSON(r, &input) != nil || input.Granted == nil {
		respond(w, validation.ErrInvalidData)
		return
	}
	respond(w, h.repo.Grant(r.Context(), id, r.PathValue("key"), *input.Granted))
}

func respond(w http.ResponseWriter, err error) {
	switch {
	case err == nil:
		httpx.WriteData(w, 200, map[string]bool{"saved": true}, nil)
	case errors.Is(err, ErrUnavailable):
		httpx.WriteError(w, 403, "plugin_unavailable", "Plugin access is required", nil)
	case errors.Is(err, pgx.ErrNoRows):
		httpx.WriteError(w, 404, "not_found", "Organization or plugin not found", nil)
	case errors.Is(err, validation.ErrInvalidData):
		httpx.WriteError(w, 400, "invalid_data", "Invalid plugin settings", nil)
	default:
		httpx.WriteError(w, 500, "internal_error", "Cannot access plugin settings", nil)
	}
}
