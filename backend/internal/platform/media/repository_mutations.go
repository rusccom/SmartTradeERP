package media

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) LockOwner(ctx context.Context, tx pgx.Tx, tenantID, ownerType, ownerID string) error {
	_, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, tenantID+"|"+ownerType+"|"+ownerID)
	return err
}

func (r *Repository) HasReady(ctx context.Context, tx pgx.Tx, tenantID, ownerType, ownerID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM platform.media_objects
        WHERE tenant_id=$1 AND owner_type=$2 AND owner_id=$3 AND status='ready')`
	exists := false
	err := tx.QueryRow(ctx, query, tenantID, ownerType, ownerID).Scan(&exists)
	return exists, err
}

func (r *Repository) DeleteReady(
	ctx context.Context,
	tx pgx.Tx,
	tenantID, ownerType, ownerID, mediaID string,
) (string, bool, error) {
	query := `DELETE FROM platform.media_objects
        WHERE tenant_id=$1 AND owner_type=$2 AND owner_id=$3 AND id=$4 AND status='ready'
        RETURNING object_key, is_primary`
	objectKey := ""
	isPrimary := false
	err := tx.QueryRow(ctx, query, tenantID, ownerType, ownerID, mediaID).Scan(&objectKey, &isPrimary)
	return objectKey, isPrimary, err
}

func (r *Repository) PromoteNextPrimary(ctx context.Context, tx pgx.Tx, tenantID, ownerType, ownerID string) error {
	query := `UPDATE platform.media_objects SET is_primary=true, updated_at=now()
        WHERE tenant_id=$1 AND id=(
            SELECT id FROM platform.media_objects
            WHERE tenant_id=$1 AND owner_type=$2 AND owner_id=$3 AND status='ready'
            ORDER BY sort_order, created_at LIMIT 1)`
	_, err := tx.Exec(ctx, query, tenantID, ownerType, ownerID)
	return err
}

func (r *Repository) SetPrimaryByID(
	ctx context.Context,
	tx pgx.Tx,
	tenantID, ownerType, ownerID, mediaID string,
) error {
	query := `UPDATE platform.media_objects SET is_primary=true, updated_at=now()
        WHERE tenant_id=$1 AND owner_type=$2 AND owner_id=$3 AND id=$4 AND status='ready'`
	tag, err := tx.Exec(ctx, query, tenantID, ownerType, ownerID, mediaID)
	if err == nil && tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (r *Repository) PrimaryByOwners(
	ctx context.Context,
	tenantID, ownerType string,
	ownerIDs []string,
) ([]PrimaryRef, error) {
	query := `SELECT owner_id::text, object_key FROM platform.media_objects
        WHERE tenant_id=$1 AND owner_type=$2 AND owner_id::text = ANY($3) AND is_primary AND status='ready'`
	rows, err := r.store.Pool.Query(ctx, query, tenantID, ownerType, ownerIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	refs := make([]PrimaryRef, 0)
	for rows.Next() {
		ref := PrimaryRef{}
		if err := rows.Scan(&ref.OwnerID, &ref.ObjectKey); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, rows.Err()
}
