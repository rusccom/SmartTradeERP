package storefront

import (
	"context"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/platform/storefront/catalog"
)

const categoriesSQL = `SELECT c.slug,c.name FROM catalog.categories c
	WHERE c.tenant_id=$1 AND (c.is_published OR $2) AND c.slug<>'' AND EXISTS (
		SELECT FROM catalog.product_categories pc JOIN catalog.products p ON p.id=pc.product_id AND p.tenant_id=pc.tenant_id
		WHERE pc.tenant_id=c.tenant_id AND pc.category_id=c.id AND (p.is_published OR $2)
		AND ($3::text='' OR p.is_composite=($3='composite'))
		AND (NOT $4 OR EXISTS(SELECT FROM catalog.product_variants v WHERE v.product_id=p.id AND v.tenant_id=p.tenant_id
			AND (v.is_published OR $2))))
	ORDER BY c.sort_order,c.name,c.slug`

func (r *Repository) categories(ctx context.Context, tenantID string, f catalog.Filter) ([]catalog.Category, error) {
	rows, err := r.store.Query(ctx, tenantID, categoriesSQL, tenantID, f.Preview, f.ProductType, f.RequireVariants)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (catalog.Category, error) {
		var c catalog.Category
		err := row.Scan(&c.Slug, &c.Name)
		return c, err
	})
}
