package inpost

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
)

type repository struct{ store *db.Store }

func (r repository) get(ctx context.Context, tenantID string) (*connection, error) {
	const query = `SELECT environment, organization_id, organization_name, services, verified_at, revision
		FROM platform.inpost_connections WHERE tenant_id = $1`
	var c connection
	err := r.store.QueryRow(ctx, tenantID, query, tenantID).Scan(
		&c.Environment, &c.OrganizationID, &c.OrganizationName, &c.Services, &c.VerifiedAt, &c.Revision)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &c, err
}

func (r repository) token(ctx context.Context, tenantID string, u update) ([]byte, error) {
	const query = `SELECT encrypted_token FROM platform.inpost_connections
		WHERE tenant_id = $1 AND environment = $2 AND organization_id = $3 AND revision = $4`
	var encrypted []byte
	err := r.store.QueryRow(ctx, tenantID, query, tenantID, u.Environment, u.OrganizationID, u.Revision).Scan(&encrypted)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errToken
	}
	return encrypted, err
}

const saveConnection = `INSERT INTO platform.inpost_connections AS saved
	(tenant_id, environment, organization_id, encrypted_token, organization_name, services)
	SELECT $1, $2, $3, $4, $5, $6 WHERE $7::bigint = 0 OR EXISTS (
		SELECT 1 FROM platform.inpost_connections WHERE tenant_id = $1 AND revision = $7)
	ON CONFLICT (tenant_id) DO UPDATE SET environment = EXCLUDED.environment,
		organization_id = EXCLUDED.organization_id, encrypted_token = EXCLUDED.encrypted_token,
		organization_name = EXCLUDED.organization_name, services = EXCLUDED.services,
		verified_at = now(), revision = saved.revision + 1 WHERE saved.revision = $7
	RETURNING verified_at, revision`

func (r repository) save(ctx context.Context, tenantID string, c *connection, encrypted []byte) error {
	row := r.store.QueryRow(ctx, tenantID, saveConnection, tenantID, c.Environment, c.OrganizationID,
		encrypted, c.OrganizationName, c.Services, c.Revision)
	err := row.Scan(&c.VerifiedAt, &c.Revision)
	if errors.Is(err, pgx.ErrNoRows) {
		return errConflict
	}
	return err
}
