package storefront

import (
	"context"
	_ "embed"
	"maps"
	"slices"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/integrations/smartbundle"
	"smarterp/backend/internal/platform/storefront/catalog"
)

//go:embed catalog_products.sql
var catalogSQL string

//go:embed catalog_bundles.sql
var catalogBundlesSQL string

func (repo *Repository) catalog(ctx context.Context, tenantID string, f catalog.Filter) ([]catalog.Product, error) {
	rows, err := repo.store.Query(ctx, tenantID, catalogSQL, tenantID, f.Slug, f.Kind, f.Collection, f.VariantIDs, f.Limit, f.Offset, f.CollectionSlugs)
	if err != nil {
		return nil, err
	}
	products, err := pgx.CollectRows(rows, scanShopProduct)
	if err != nil {
		return nil, err
	}
	return products, repo.loadRecipes(ctx, tenantID, products)
}

func scanShopProduct(row pgx.CollectableRow) (catalog.Product, error) {
	var p catalog.Product
	err := row.Scan(&p.Name, &p.Slug, &p.SEOTitle, &p.SEODescription, &p.Composite, &p.Categories, &p.Images, &p.Variants)
	return p, err
}

func (repo *Repository) loadRecipes(ctx context.Context, tenantID string, products []catalog.Product) error {
	targets := recipeTargets(products)
	if len(targets) == 0 {
		return nil
	}
	ids := slices.Collect(maps.Keys(targets))
	rows, err := repo.store.Query(ctx, tenantID, catalogBundlesSQL, tenantID, ids)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := scanSlot(rows, targets); err != nil {
			return err
		}
	}
	return rows.Err()
}

func indexRecipes(p catalog.Product, targets map[string]*smartbundle.Bundle) {
	for i := range p.Variants {
		recipe := &smartbundle.Bundle{Mode: "fixed"}
		p.Variants[i].Bundle = recipe
		targets[p.Variants[i].ID] = recipe
	}
}

func scanSlot(rows pgx.Rows, targets map[string]*smartbundle.Bundle) error {
	var id, mode string
	var slot smartbundle.Slot
	var recipe smartbundle.Bundle
	if err := rows.Scan(&id, &mode, &recipe.Discount, &slot.Key, &slot.Name, &slot.Qty, &slot.Discount, &slot.Options); err != nil {
		return err
	}
	b := targets[id]
	b.Mode, b.Discount = mode, recipe.Discount
	b.Slots = append(b.Slots, slot)
	return nil
}

func recipeTargets(products []catalog.Product) map[string]*smartbundle.Bundle {
	targets := make(map[string]*smartbundle.Bundle)
	for _, p := range products {
		if p.Composite {
			indexRecipes(p, targets)
		}
	}
	return targets
}
