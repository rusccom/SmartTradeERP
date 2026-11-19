package inpost

import (
	"net/http"
	"strings"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
)

type Handler struct{ service service }

func NewHandler(store *db.Store, key string) *Handler {
	return &Handler{service{repo: repository{store}, client: newShipXClient(), cipher: newTokenCipher(key)}}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-store")
	ctx := r.Context()
	c, err := h.service.repo.get(ctx, tenant.FromContext(ctx))
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.WriteData(w, http.StatusOK, map[string]any{
		"can_manage": auth.ClaimsFromContext(ctx).Role == "owner", "connection": c,
	}, nil)
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-store")
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var u update
	if err := httpx.DecodeJSON(r, &u); err != nil {
		writeError(w, errInvalid)
		return
	}
	u.Token = strings.TrimSpace(u.Token)
	c, err := h.service.save(r.Context(), tenant.FromContext(r.Context()), u)
	if err != nil {
		writeError(w, err)
		return
	}
	httpx.WriteData(w, http.StatusOK, c, nil)
}

func writeError(w http.ResponseWriter, err error) {
	switch err {
	case errInvalid, errToken:
		httpx.WriteError(w, 400, "invalid_settings", "Check the environment, organization ID and ShipX token.", nil)
	case errCredentials:
		httpx.WriteError(w, 422, "inpost_access_denied", "InPost could not verify this organization with the supplied token and environment.", nil)
	case errCarrier:
		httpx.WriteError(w, 502, "inpost_unavailable", "Cannot verify the organization with InPost right now. Try again.", nil)
	case errThrottled:
		httpx.WriteError(w, 429, "inpost_rate_limit", "InPost has limited verification requests. Try again later.", nil)
	case errVault:
		httpx.WriteError(w, 503, "credential_storage_unavailable", "InPost credential storage is unavailable.", nil)
	case errConflict:
		httpx.WriteError(w, 409, "settings_changed", "These settings changed. Reload the page before saving.", nil)
	default:
		httpx.WriteError(w, 500, "internal_error", "Cannot access InPost settings.", nil)
	}
}
