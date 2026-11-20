package connection

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Renew func(context.Context, string) (Credential, error)

func (s *Service) Check(ctx context.Context, tenant string, renew Renew) error {
	var result error
	err := s.DB.WithTenantTx(ctx, tenant, func(tx pgx.Tx) error {
		a, err := lockAccount(ctx, tx, tenant)
		if err != nil {
			return err
		}
		result = s.checkSaved(ctx, tx, a, renew)
		return commitCheck(ctx, tx, tenant, result)
	})
	if err != nil {
		return err
	}
	return result
}

func lockAccount(ctx context.Context, tx pgx.Tx, tenant string) (stored, error) {
	a := stored{tenant: tenant}
	const query = `SELECT seller_id,method,expires_at,encrypted_tokens FROM platform.allegro_connections WHERE tenant_id=$1 FOR UPDATE`
	err := tx.QueryRow(ctx, query, tenant).Scan(&a.ID, &a.Method, &a.ExpiresAt, &a.encrypted)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, Missing
	}
	return a, err
}

// Save a rotated refresh token even if a subsequent request to Allegro fails.
func commitCheck(ctx context.Context, tx pgx.Tx, tenant string, result error) error {
	var p *Problem
	if result == nil || !errors.As(result, &p) {
		return result
	}
	if p == Unauthorized || p == Scope || p == Different {
		_, err := tx.Exec(ctx, `UPDATE platform.allegro_connections SET needs_reconnect=true WHERE tenant_id=$1`, tenant)
		return err
	}
	return nil
}

func (s *Service) checkSaved(ctx context.Context, tx pgx.Tx, a stored, renew Renew) error {
	key := binding(a.tenant, a.Method, a.ID)
	var pair Tokens
	if err := s.Vault.Decrypt(key, a.encrypted, &pair); err != nil {
		return err
	}
	if a.Method == Manual {
		return s.verifySaved(ctx, tx, a, pair.Access)
	}
	if a.ExpiresAt == nil || time.Until(*a.ExpiresAt) > time.Minute {
		if err := s.verifySaved(ctx, tx, a, pair.Access); err != Unauthorized {
			return err
		}
	}
	return s.renewAndVerify(ctx, tx, a, renew, pair.Refresh)
}

func (s *Service) renewAndVerify(ctx context.Context, tx pgx.Tx, a stored, renew Renew, refresh string) error {
	fresh, err := renew(ctx, refresh)
	if err != nil {
		return err
	}
	if err = s.rotate(ctx, tx, a, fresh); err != nil {
		return err
	}
	return s.verifySaved(ctx, tx, a, fresh.Access)
}

func (s *Service) rotate(ctx context.Context, tx pgx.Tx, a stored, fresh Credential) error {
	data, err := s.Vault.Encrypt(binding(a.tenant, a.Method, a.ID), fresh.Tokens)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE platform.allegro_connections SET encrypted_tokens=$2,expires_at=$3 WHERE tenant_id=$1`,
		a.tenant, data, fresh.ExpiresAt)
	return err
}

func (s *Service) verifySaved(ctx context.Context, tx pgx.Tx, a stored, token string) error {
	seller, err := s.Client.Verify(ctx, token)
	if err != nil {
		return err
	}
	if seller.ID != a.ID {
		return Different
	}
	_, err = tx.Exec(ctx, `UPDATE platform.allegro_connections SET login=$2,verified_at=now(),needs_reconnect=false
		WHERE tenant_id=$1`, a.tenant, seller.Login)
	return err
}
