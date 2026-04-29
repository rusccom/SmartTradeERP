package bundles

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

type Repository struct {
	store *db.Store
}

func NewRepository(store *db.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) List(ctx context.Context, tenantID string, page, perPage int) ([]Bundle, int, error) {
	total, err := r.count(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	data, err := r.load(ctx, tenantID, page, perPage)
	return data, total, err
}

func (r *Repository) count(ctx context.Context, tenantID string) (int, error) {
	query := `SELECT COUNT(*)
        FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        WHERE v.tenant_id=$1 AND p.tenant_id=$1 AND p.is_composite=true`
	total := 0
	err := r.store.Pool.QueryRow(ctx, query, tenantID).Scan(&total)
	return total, err
}

func (r *Repository) load(ctx context.Context, tenantID string, page, perPage int) ([]Bundle, error) {
	query := bundleSelect() + ` ORDER BY p.created_at DESC, v.created_at DESC LIMIT $2 OFFSET $3`
	offset := httpx.Offset(page, perPage)
	rows, err := r.store.Pool.Query(ctx, query, tenantID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBundles(rows)
}

func (r *Repository) ByVariantID(ctx context.Context, tenantID, variantID string) (Bundle, error) {
	query := bundleSelect() + ` AND v.id::text=$2`
	row := r.store.Pool.QueryRow(ctx, query, tenantID, variantID)
	return scanBundle(row)
}

func bundleSelect() string {
	return `SELECT p.id::text, v.id::text, p.name, COALESCE(v.name,'Default'),
        COALESCE(v.sku_code,''), v.unit, COALESCE(v.price,0)
        FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        WHERE v.tenant_id=$1 AND p.tenant_id=$1 AND p.is_composite=true`
}

func scanBundles(rows pgx.Rows) ([]Bundle, error) {
	items := make([]Bundle, 0)
	for rows.Next() {
		item, err := scanBundle(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanBundle(row pgx.Row) (Bundle, error) {
	item := Bundle{}
	err := row.Scan(&item.ProductID, &item.VariantID, &item.ProductName,
		&item.VariantName, &item.SKUCode, &item.Unit, &item.Price)
	return item, err
}

func (r *Repository) Components(ctx context.Context, tenantID, variantID string) ([]Component, error) {
	rows, err := r.store.Pool.Query(ctx, componentsSQL(), tenantID, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanComponents(rows)
}

func (r *Repository) ComponentsTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID string,
	variantID string,
) ([]Component, error) {
	rows, err := tx.Query(ctx, componentsSQL(), tenantID, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanComponents(rows)
}

func componentsSQL() string {
	return `SELECT vc.component_variant_id::text, vc.qty,
        p.name, COALESCE(component.name,'Default'), component.unit
        FROM catalog.variant_components vc
        JOIN catalog.product_variants parent ON parent.id=vc.variant_id
        JOIN catalog.product_variants component ON component.id=vc.component_variant_id
        JOIN catalog.products p ON p.id=component.product_id
        WHERE parent.tenant_id=$1 AND component.tenant_id=$1 AND vc.variant_id=$2
        ORDER BY vc.id`
}

func scanComponents(rows pgx.Rows) ([]Component, error) {
	items := make([]Component, 0)
	for rows.Next() {
		item := Component{}
		err := rows.Scan(&item.ComponentVariantID, &item.Qty,
			&item.ProductName, &item.VariantName, &item.Unit)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) VariantComposite(ctx context.Context, tenantID, variantID string) (bool, error) {
	query := `SELECT p.is_composite
        FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        WHERE v.tenant_id=$1 AND v.id=$2`
	value := false
	err := r.store.Pool.QueryRow(ctx, query, tenantID, variantID).Scan(&value)
	return value, err
}

func (r *Repository) ComponentUsable(ctx context.Context, tenantID, variantID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.product_variants v
        JOIN catalog.products p ON p.id=v.product_id
        WHERE v.tenant_id=$1 AND v.id=$2 AND p.is_composite=false
    )`
	exists := false
	err := r.store.Pool.QueryRow(ctx, query, tenantID, variantID).Scan(&exists)
	return exists, err
}

func (r *Repository) ProductHasComponents(ctx context.Context, tenantID, productID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.variant_components c
        JOIN catalog.product_variants v ON v.id=c.variant_id
        WHERE v.tenant_id=$1 AND v.product_id=$2
    )`
	exists := false
	err := r.store.Pool.QueryRow(ctx, query, tenantID, productID).Scan(&exists)
	return exists, err
}

func (r *Repository) ProductUsedAsComponent(ctx context.Context, tenantID, productID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.variant_components c
        JOIN catalog.product_variants v ON v.id=c.component_variant_id
        WHERE v.tenant_id=$1 AND v.product_id=$2
    )`
	exists := false
	err := r.store.Pool.QueryRow(ctx, query, tenantID, productID).Scan(&exists)
	return exists, err
}

func (r *Repository) VariantHasComponents(ctx context.Context, tenantID, variantID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.variant_components c
        JOIN catalog.product_variants v ON v.id=c.variant_id
        WHERE v.tenant_id=$1 AND c.variant_id=$2
    )`
	exists := false
	err := r.store.Pool.QueryRow(ctx, query, tenantID, variantID).Scan(&exists)
	return exists, err
}

func (r *Repository) VariantUsedAsComponent(ctx context.Context, tenantID, variantID string) (bool, error) {
	query := `SELECT EXISTS(
        SELECT 1 FROM catalog.variant_components c
        JOIN catalog.product_variants v ON v.id=c.component_variant_id
        WHERE v.tenant_id=$1 AND c.component_variant_id=$2
    )`
	exists := false
	err := r.store.Pool.QueryRow(ctx, query, tenantID, variantID).Scan(&exists)
	return exists, err
}

func (r *Repository) ReplaceComponents(ctx context.Context, tx pgx.Tx, variantID string, items []Component) error {
	if err := deleteComponents(ctx, tx, variantID); err != nil {
		return err
	}
	for _, item := range items {
		row := componentInsert{ID: uuid.NewString(), VariantID: variantID, Item: item}
		if err := insertComponent(ctx, tx, row); err != nil {
			return err
		}
	}
	return nil
}

func deleteComponents(ctx context.Context, tx pgx.Tx, variantID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM catalog.variant_components WHERE variant_id=$1`, variantID)
	return err
}

func insertComponent(ctx context.Context, tx pgx.Tx, row componentInsert) error {
	query := `INSERT INTO catalog.variant_components (id, variant_id, component_variant_id, qty)
        VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(ctx, query, row.ID, row.VariantID, row.Item.ComponentVariantID, row.Item.Qty)
	return err
}

func (r *Repository) SaveSnapshot(ctx context.Context, tx pgx.Tx, input SnapshotInput) error {
	for _, item := range input.Components {
		if err := insertSnapshot(ctx, tx, input, item); err != nil {
			return err
		}
	}
	return nil
}

func insertSnapshot(ctx context.Context, tx pgx.Tx, input SnapshotInput, item Component) error {
	query := `INSERT INTO documents.document_item_components
        (id, document_item_id, component_variant_id, qty_per_unit, qty_total)
        VALUES ($1,$2,$3,$4,$5)`
	total := item.Qty.Mul(input.DocumentQty)
	_, err := tx.Exec(ctx, query, uuid.NewString(), input.DocumentItemID,
		item.ComponentVariantID, item.Qty, total)
	return err
}
