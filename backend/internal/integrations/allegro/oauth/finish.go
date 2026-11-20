package oauth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/integrations/allegro/connection"
)

func CancelPending(ctx context.Context, tx pgx.Tx, tenant string) error {
	_, err := tx.Exec(ctx, `DELETE FROM platform.allegro_authorizations WHERE tenant_id=$1`, tenant)
	return err
}

func takeRequest(ctx context.Context, tx pgx.Tx, f flow) (string, error) {
	const query = `DELETE FROM platform.allegro_authorizations a USING platform.tenant_users u
		WHERE a.tenant_id=$1 AND a.nonce=$2 AND a.launched AND a.expires_at>now()
		AND u.tenant_id=a.tenant_id AND u.id=a.owner_id AND u.role='owner' AND u.is_active
		RETURNING COALESCE(a.expected_version::text,'')`
	var version string
	err := tx.QueryRow(ctx, query, f.Tenant, f.State).Scan(&version)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", connection.State
	}
	return version, err
}

func (h *Handler) finish(ctx context.Context, f flow, code string) error {
	var failed error
	err := h.DB.WithTenantTx(ctx, f.Tenant, func(tx pgx.Tx) error {
		version, err := takeRequest(ctx, tx, f)
		if err != nil || code == "" {
			return err
		}
		g, err := h.grant(ctx, f, code)
		if err != nil {
			failed = err
			return nil
		}
		g.Version = version
		return h.SaveTx(ctx, tx, f.Tenant, g)
	})
	if err != nil {
		return err
	}
	return failed
}

func (h *Handler) grant(ctx context.Context, f flow, code string) (connection.Grant, error) {
	cred, err := h.exchange(ctx, code, f.Verifier)
	if err != nil {
		return connection.Grant{}, err
	}
	account, err := h.Client.Verify(ctx, cred.Access)
	return connection.Grant{Credential: cred, Seller: account, Method: connection.OAuth}, err
}
