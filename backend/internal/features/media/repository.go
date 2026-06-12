package media

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
)

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) List(
	ctx context.Context,
	tenantID string,
	ownerType string,
	ownerID string,
) ([]Item, error) {
	query := `SELECT id::text, owner_type, owner_id::text, object_key, file_name,
        content_type, size_bytes, is_primary, status, created_at::text
        FROM platform.media_objects
        WHERE tenant_id=$1 AND owner_type=$2 AND owner_id=$3 AND status='ready'
        ORDER BY is_primary DESC, sort_order, created_at`
	rows, err := r.store.Pool.Query(ctx, query, tenantID, ownerType, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanItems(rows)
}

func (r *Repository) GetPending(
	ctx context.Context,
	tenantID string,
	ownerType string,
	ownerID string,
	mediaID string,
) (Item, error) {
	query := `SELECT id::text, owner_type, owner_id::text, object_key, file_name,
        content_type, size_bytes, is_primary, status, created_at::text
        FROM platform.media_objects
        WHERE tenant_id=$1 AND owner_type=$2 AND owner_id=$3
        AND id=$4 AND status='pending' AND expires_at > now()`
	row := r.store.Pool.QueryRow(ctx, query, tenantID, ownerType, ownerID, mediaID)
	return scanRow(row)
}

func (r *Repository) CreatePending(
	ctx context.Context,
	tenantID string,
	item Item,
	expiresAt time.Time,
) error {
	query := `INSERT INTO platform.media_objects
        (id, tenant_id, owner_type, owner_id, object_key, file_name,
        content_type, size_bytes, status, expires_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'pending',$9)`
	_, err := r.store.Pool.Exec(ctx, query, item.ID, tenantID, item.OwnerType,
		item.OwnerID, item.ObjectKey, item.FileName, item.ContentType,
		item.SizeBytes, expiresAt)
	return err
}

func (r *Repository) ClearPrimary(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	ownerType string,
	ownerID string,
) error {
	query := `UPDATE platform.media_objects SET is_primary=false, updated_at=now()
        WHERE tenant_id=$1 AND owner_type=$2 AND owner_id=$3 AND is_primary=true`
	_, err := tx.Exec(ctx, query, tenantID, ownerType, ownerID)
	return err
}

func (r *Repository) MarkReady(ctx context.Context, tx pgx.Tx, tenantID string, item Item) error {
	query := `UPDATE platform.media_objects
        SET status='ready', is_primary=true, expires_at=NULL, updated_at=now()
        WHERE tenant_id=$1 AND id=$2 AND status='pending'`
	tag, err := tx.Exec(ctx, query, tenantID, item.ID)
	if err == nil && tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (r *Repository) Delete(ctx context.Context, tenantID string, mediaID string) error {
	query := `DELETE FROM platform.media_objects WHERE tenant_id=$1 AND id=$2`
	_, err := r.store.Pool.Exec(ctx, query, tenantID, mediaID)
	return err
}

func scanItems(rows pgx.Rows) ([]Item, error) {
	items := make([]Item, 0)
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanItem(rows pgx.Rows) (Item, error) {
	item := Item{}
	err := rows.Scan(&item.ID, &item.OwnerType, &item.OwnerID, &item.ObjectKey,
		&item.FileName, &item.ContentType, &item.SizeBytes,
		&item.IsPrimary, &item.Status, &item.CreatedAt)
	return item, err
}

func scanRow(row pgx.Row) (Item, error) {
	item := Item{}
	err := row.Scan(&item.ID, &item.OwnerType, &item.OwnerID, &item.ObjectKey,
		&item.FileName, &item.ContentType, &item.SizeBytes,
		&item.IsPrimary, &item.Status, &item.CreatedAt)
	return item, err
}
