package connection

import (
	"context"
	"net/http"
	"time"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/tenant"
)

type Handler struct {
	*Service
	OAuthReady    func() bool
	Renew         Renew
	CancelPending CancelPending
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c, err := h.Status(ctx, tenant.FromContext(ctx))
	Respond(w, map[string]any{"connection": c, "token_available": h.Vault.AEAD != nil, "oauth_available": h.OAuthReady(),
		"can_manage": auth.ClaimsFromContext(ctx).Role == "owner"}, err)
}

func (h Handler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	Respond(w, nil, h.Service.Check(ctx, tenant.FromContext(ctx), h.Renew))
}

func (h Handler) Disconnect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	Respond(w, nil, h.Service.Disconnect(ctx, tenant.FromContext(ctx), r.PathValue("version"), h.CancelPending))
}
