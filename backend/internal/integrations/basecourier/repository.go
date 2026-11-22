package basecourier

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"smarterp/backend/internal/shared/db"
)

type repository struct{ store *db.Store }

func scanConnection(row pgx.Row, missing error) (*connection, error) {
	var c connection
	err := row.Scan(&c.Email, &c.Version, &c.VerifiedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, missing
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r repository) get(ctx context.Context, tenantID string) (*connection, error) {
	const query = `SELECT email, version, verified_at FROM platform.base_courier_connections WHERE tenant_id=$1`
	return scanConnection(r.store.QueryRow(ctx, tenantID, query, tenantID), nil)
}

func (r repository) save(ctx context.Context, tenantID string, u update, sealed []byte) (*connection, error) {
	const create = `INSERT INTO platform.base_courier_connections (tenant_id, email, encrypted_api_key)
		VALUES ($1, $2, $3) ON CONFLICT (tenant_id) DO NOTHING RETURNING email, version, verified_at`
	const replace = `UPDATE platform.base_courier_connections
		SET email=$2, encrypted_api_key=$3, version=gen_random_uuid(), verified_at=now()
		WHERE tenant_id=$1 AND version=$4 RETURNING email, version, verified_at`
	if u.Version == nil {
		return scanConnection(r.store.QueryRow(ctx, tenantID, create, tenantID, u.Email, sealed), errChanged)
	}
	return scanConnection(r.store.QueryRow(ctx, tenantID, replace, tenantID, u.Email, sealed, *u.Version), errChanged)
}

func (r repository) credential(ctx context.Context, tenantID string, version uuid.UUID) (string, []byte, error) {
	const query = `SELECT email, encrypted_api_key FROM platform.base_courier_connections
		WHERE tenant_id=$1 AND version=$2`
	var email string
	var sealed []byte
	err := r.store.QueryRow(ctx, tenantID, query, tenantID, version).Scan(&email, &sealed)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil, errChanged
	}
	return email, sealed, err
}

func (r repository) markVerified(ctx context.Context, tenantID string, version uuid.UUID) (*connection, error) {
	const query = `UPDATE platform.base_courier_connections SET verified_at=now()
		WHERE tenant_id=$1 AND version=$2 RETURNING email, version, verified_at`
	return scanConnection(r.store.QueryRow(ctx, tenantID, query, tenantID, version), errChanged)
}

func (r repository) remove(ctx context.Context, tenantID string, version uuid.UUID) error {
	const query = `DELETE FROM platform.base_courier_connections WHERE tenant_id=$1 AND version=$2`
	tag, err := r.store.Exec(ctx, tenantID, query, tenantID, version)
	if err == nil && tag.RowsAffected() == 0 {
		return errChanged
	}
	return err
}
