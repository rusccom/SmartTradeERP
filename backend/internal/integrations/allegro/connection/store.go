package connection

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/db"
)

type Service struct {
	DB     *db.Store
	Client Client
	Vault  Vault
}

func NewService(store *db.Store, key string) *Service {
	return &Service{store, NewClient(), NewVault(key)}
}

func (s *Service) Status(ctx context.Context, tenant string) (*Account, error) {
	const query = `SELECT seller_id,login,method,version,verified_at,expires_at,needs_reconnect
		FROM platform.allegro_connections WHERE tenant_id=$1`
	var a Account
	err := s.DB.QueryRow(ctx, tenant, query, tenant).Scan(&a.ID, &a.Login, &a.Method, &a.Version, &a.VerifiedAt, &a.ExpiresAt, &a.NeedsReconnect)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &a, err
}

func (s *Service) Save(ctx context.Context, tenant string, g Grant) error {
	return s.DB.WithTenantTx(ctx, tenant, func(tx pgx.Tx) error { return s.SaveTx(ctx, tx, tenant, g) })
}

func (s *Service) SaveTx(ctx context.Context, tx pgx.Tx, tenant string, g Grant) error {
	encrypted, err := s.Vault.Encrypt(binding(tenant, g.Method, g.ID), g.Tokens)
	if err != nil {
		return err
	}
	const query = `INSERT INTO platform.allegro_connections AS c
		(tenant_id,seller_id,login,method,encrypted_tokens,expires_at,version) SELECT $1,$2,$3,$4,$5,$6,$7
		WHERE $8='' OR EXISTS (SELECT FROM platform.allegro_connections WHERE tenant_id=$1 AND version::text=$8)
		ON CONFLICT (tenant_id) DO UPDATE SET login=EXCLUDED.login,encrypted_tokens=EXCLUDED.encrypted_tokens,
		expires_at=EXCLUDED.expires_at,version=EXCLUDED.version,verified_at=now(),needs_reconnect=false
		WHERE c.version::text=$8 AND c.seller_id=EXCLUDED.seller_id AND c.method=EXCLUDED.method`
	tag, err := tx.Exec(ctx, query, tenant, g.ID, g.Login, g.Method, encrypted, g.ExpiresAt, uuid.NewString(), g.Version)
	if err == nil && tag.RowsAffected() != 1 {
		return conflict(ctx, tx, tenant, g)
	}
	return err
}

func conflict(ctx context.Context, tx pgx.Tx, tenant string, g Grant) error {
	var other bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM platform.allegro_connections
		WHERE tenant_id=$1 AND (method<>$2 OR seller_id<>$3))`, tenant, g.Method, g.ID).Scan(&other)
	if err != nil {
		return err
	}
	if other {
		return Different
	}
	return Changed
}

type CancelPending func(context.Context, pgx.Tx, string) error

func (s *Service) Disconnect(ctx context.Context, tenant, version string, cancelPending CancelPending) error {
	if uuid.Validate(version) != nil {
		return Changed
	}
	return s.DB.WithTenantTx(ctx, tenant, func(tx pgx.Tx) error {
		if err := cancelPending(ctx, tx, tenant); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `DELETE FROM platform.allegro_connections WHERE tenant_id=$1 AND version=$2`, tenant, version)
		if err == nil && tag.RowsAffected() == 0 {
			return Changed
		}
		return err
	})
}
