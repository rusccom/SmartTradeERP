package oauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"smarterp/backend/internal/integrations/allegro/connection"
	"smarterp/backend/internal/shared/auth"
)

const lifetime = 900

type flow struct {
	Tenant   string
	State    string
	Verifier string
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := auth.ClaimsFromContext(ctx)
	ticket, err := h.ticket(ctx, claims.TenantID, claims.UserID)
	connection.Respond(w, map[string]string{"ticket": ticket, "launch_url": h.cfg.APIOrigin + Path + "/launch"}, err)
}

func (h *Handler) ticket(ctx context.Context, tenant, owner string) (string, error) {
	if !h.Enabled() {
		return "", connection.Setup
	}
	const query = `INSERT INTO platform.allegro_authorizations (tenant_id,owner_id,nonce,expires_at,expected_version)
		SELECT $1,$2,$3,now()+interval '15 minutes',c.version FROM (VALUES(1)) AS guard(n)
		LEFT JOIN platform.allegro_connections c ON c.tenant_id=$1 WHERE c.method IS NULL OR c.method='oauth'
		ON CONFLICT (tenant_id) DO UPDATE SET owner_id=EXCLUDED.owner_id,nonce=EXCLUDED.nonce,launched=false,
		expires_at=EXCLUDED.expires_at,expected_version=EXCLUDED.expected_version`
	nonce := uuid.NewString()
	tag, err := h.DB.Exec(ctx, tenant, query, tenant, owner, nonce)
	if err == nil && tag.RowsAffected() != 1 {
		return "", connection.Different
	}
	return tenant + "." + nonce, err
}

func (h *Handler) launch(ctx context.Context, ticket string) (flow, error) {
	tenant, nonce, ok := strings.Cut(ticket, ".")
	if !ok || uuid.Validate(tenant) != nil || uuid.Validate(nonce) != nil {
		return flow{}, connection.State
	}
	f := flow{Tenant: tenant, State: uuid.NewString()}
	verifier := make([]byte, 32)
	if _, err := rand.Read(verifier); err != nil {
		return f, err
	}
	f.Verifier = base64.RawURLEncoding.EncodeToString(verifier)
	const query = `UPDATE platform.allegro_authorizations SET nonce=$3,launched=true
		WHERE tenant_id=$1 AND nonce=$2 AND NOT launched AND expires_at>now()`
	tag, err := h.DB.Exec(ctx, tenant, query, tenant, nonce, f.State)
	if err == nil && tag.RowsAffected() != 1 {
		return f, connection.State
	}
	return f, err
}

func (f flow) challenge() string {
	hash := sha256.Sum256([]byte(f.Verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
