package documents

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) WarehouseExists(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.warehouses WHERE tenant_id=$1 AND id=$2
    )`
	return exists(ctx, tx, query, tenantID, id)
}

func (r *Repository) CustomerExists(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.customers WHERE tenant_id=$1 AND id=$2
    )`
	return exists(ctx, tx, query, tenantID, id)
}

func (r *Repository) ShiftExists(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM documents.shifts WHERE tenant_id=$1 AND id=$2
    )`
	return exists(ctx, tx, query, tenantID, id)
}

func (r *Repository) OpenShiftExists(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM documents.shifts WHERE tenant_id=$1 AND id=$2 AND status='open'
    )`
	return exists(ctx, tx, query, tenantID, id)
}

func (r *Repository) VariantComposite(ctx context.Context, tx pgx.Tx, tenantID string, id string) (bool, bool, error) {
	query := `SELECT COALESCE(bool_or(p.is_composite), false), COUNT(*) > 0
        FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        WHERE v.tenant_id=$1 AND p.tenant_id=$1 AND v.id=$2`
	composite := false
	exists := false
	row := tx.QueryRow(ctx, query, tenantID, id)
	err := row.Scan(&composite, &exists)
	return composite, exists, err
}

func exists(ctx context.Context, tx pgx.Tx, query string, tenantID string, id string) (bool, error) {
	value := false
	row := tx.QueryRow(ctx, query, tenantID, id)
	err := row.Scan(&value)
	return value, err
}
