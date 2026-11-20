package token

import (
	"context"
	"net/http"
	"strings"
	"time"
	"unicode"

	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/shared/httpx"
	"smarterp/backend/internal/shared/tenant"
)

type Handler struct{ Service *connection.Service }

type input struct {
	Access  string `json:"access_token"`
	Version string `json:"version"`
}

func (h Handler) Save(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
	var in input
	if httpx.DecodeJSON(r, &in) != nil {
		connection.Respond(w, nil, connection.Invalid)
		return
	}
	in.Access = strings.TrimSpace(in.Access)
	if in.Access == "" || len(in.Access) > 16384 || strings.IndexFunc(in.Access, invalidChar) >= 0 {
		connection.Respond(w, nil, connection.Invalid)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	connection.Respond(w, nil, h.save(ctx, tenant.FromContext(ctx), in))
}

func invalidChar(r rune) bool { return r > 126 || unicode.IsSpace(r) || unicode.IsControl(r) }

func (h Handler) save(ctx context.Context, tenant string, in input) error {
	if h.Service.Vault.AEAD == nil {
		return connection.Setup
	}
	account, err := h.Service.Client.Verify(ctx, in.Access)
	if err != nil {
		return err
	}
	cred := connection.Credential{Tokens: connection.Tokens{Access: in.Access}, ExpiresAt: expiry(in.Access)}
	return h.Service.Save(ctx, tenant, connection.Grant{Seller: account, Credential: cred, Method: connection.Manual, Version: in.Version})
}
