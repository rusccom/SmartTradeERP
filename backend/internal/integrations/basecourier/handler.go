package basecourier

import (
	"net/http"

	"github.com/google/uuid"
	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
)

type Handler struct{ service service }

func NewHandler(store *db.Store, key string) *Handler {
	return &Handler{service{repo: repository{store}, client: newProfileClient(), vault: newVault(key)}}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c, err := h.service.repo.get(ctx, tenant.FromContext(ctx))
	respond(w, map[string]any{"can_manage": auth.ClaimsFromContext(ctx).Role == "owner", "connection": c}, err)
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var u update
	if err := httpx.DecodeJSON(r, &u); err != nil {
		respond(w, nil, errInvalid)
		return
	}
	ctx := r.Context()
	c, err := h.service.save(ctx, tenant.FromContext(ctx), u)
	respond(w, c, err)
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	version, err := uuid.Parse(r.PathValue("version"))
	if err != nil {
		respond(w, nil, errVersion)
		return
	}
	ctx := r.Context()
	c, err := h.service.check(ctx, tenant.FromContext(ctx), version)
	respond(w, c, err)
}

func (h *Handler) Disconnect(w http.ResponseWriter, r *http.Request) {
	version, err := uuid.Parse(r.PathValue("version"))
	if err != nil {
		respond(w, nil, errVersion)
		return
	}
	ctx := r.Context()
	if err = h.service.repo.remove(ctx, tenant.FromContext(ctx), version); err != nil {
		respond(w, nil, err)
		return
	}
	w.Header().Set("Cache-Control", "private, no-store")
	w.WriteHeader(http.StatusNoContent)
}
