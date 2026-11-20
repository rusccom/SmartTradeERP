package connection

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Renew func(context.Context, string) (Credential, error)

type PrivateRenewal struct {
	TenantID string
	SellerID string
	ClientID string
	Renew    Renew
}

func (p PrivateRenewal) matches(a stored) bool {
	return p.Renew != nil && a.privateClient != "" && p.TenantID == a.tenant && p.SellerID == a.ID && p.ClientID == a.privateClient
}

type lockedCheck func(pgx.Tx, stored) error

func (s *Service) runCheck(ctx context.Context, tenant string, check lockedCheck) error {
	var result error
	err := s.DB.WithTenantTx(ctx, tenant, func(tx pgx.Tx) error {
		a, err := lockAccount(ctx, tx, tenant)
		if err != nil {
			return err
		}
		result = check(tx, a)
		return commitCheck(ctx, tx, tenant, result)
	})
	if err != nil {
		return err
	}
	return result
}

func lockAccount(ctx context.Context, tx pgx.Tx, tenant string) (stored, error) {
	a := stored{tenant: tenant}
	const query = `SELECT seller_id,method,expires_at,encrypted_tokens,COALESCE(private_client_id,''),refreshed_at,needs_reconnect
		FROM platform.allegro_connections WHERE tenant_id=$1 FOR UPDATE`
	err := tx.QueryRow(ctx, query, tenant).Scan(&a.ID, &a.Method, &a.ExpiresAt, &a.encrypted,
		&a.privateClient, &a.RefreshedAt, &a.NeedsReconnect)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, Missing
	}
	return a, err
}

// Keep a rotated credential even when the subsequent order-access check fails.
func commitCheck(ctx context.Context, tx pgx.Tx, tenant string, result error) error {
	var p *Problem
	if result == nil || !errors.As(result, &p) {
		return result
	}
	if p == Unauthorized || p == Scope || p == Different || p == ReadOnly {
		_, err := tx.Exec(ctx, `UPDATE platform.allegro_connections SET needs_reconnect=true WHERE tenant_id=$1`, tenant)
		return err
	}
	return nil
}

func (s *Service) Check(ctx context.Context, tenant string, commercial Renew, private PrivateRenewal) error {
	return s.runCheck(ctx, tenant, func(tx pgx.Tx, a stored) error {
		renew := commercial
		if a.privateClient != "" {
			if !private.matches(a) {
				return Setup
			}
			renew = private.Renew
		}
		return s.checkSaved(ctx, tx, a, renew)
	})
}

func (s *Service) checkSaved(ctx context.Context, tx pgx.Tx, a stored, renew Renew) error {
	var pair Tokens
	if err := s.Vault.Decrypt(binding(a.tenant, a.Method, a.ID), a.encrypted, &pair); err != nil {
		return err
	}
	if a.Method == Manual && a.privateClient == "" {
		return s.verifySaved(ctx, tx, a, pair.Access)
	}
	if a.ExpiresAt == nil || time.Until(*a.ExpiresAt) > time.Minute {
		if err := s.verifySaved(ctx, tx, a, pair.Access); err != Unauthorized {
			return err
		}
	}
	return s.renewAndVerify(ctx, tx, a, renew, pair.Refresh)
}

func (s *Service) MaintainPrivate(ctx context.Context, private PrivateRenewal) error {
	return s.runCheck(ctx, private.TenantID, func(tx pgx.Tx, a stored) error {
		if !private.matches(a) || a.NeedsReconnect {
			return nil
		}
		if a.RefreshedAt != nil && a.ExpiresAt != nil && time.Until(*a.ExpiresAt) > 5*time.Minute {
			return nil
		}
		var pair Tokens
		if err := s.Vault.Decrypt(binding(a.tenant, a.Method, a.ID), a.encrypted, &pair); err != nil {
			return err
		}
		return s.renewAndVerify(ctx, tx, a, private.Renew, pair.Refresh)
	})
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
	_, err = tx.Exec(ctx, `UPDATE platform.allegro_connections SET encrypted_tokens=$2,expires_at=$3,refreshed_at=now()
		WHERE tenant_id=$1`, a.tenant, data, fresh.ExpiresAt)
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
