package documents

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) WarehouseExists(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.warehouses WHERE id=$1
    )`
	return exists(ctx, tx, query, id)
}

func (r *Repository) CustomerExists(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.customers WHERE id=$1
    )`
	return exists(ctx, tx, query, id)
}

func (r *Repository) OpenShiftExists(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, error) {
	query := `SELECT status FROM documents.shifts
        WHERE id=$1
        FOR UPDATE`
	status := ""
	err := tx.QueryRow(ctx, query, id).Scan(&status)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return status == "open", err
}

func (r *Repository) VariantComposite(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, bool, error) {
	query := `SELECT COALESCE(bool_or(p.is_composite), false), COUNT(*) > 0
        FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        WHERE v.id=$1`
	composite := false
	exists := false
	row := tx.QueryRow(ctx, query, id)
	err := row.Scan(&composite, &exists)
	return composite, exists, err
}

func exists(ctx context.Context, tx pgx.Tx, query string, id string) (bool, error) {
	value := false
	row := tx.QueryRow(ctx, query, id)
	err := row.Scan(&value)
	return value, err
}
