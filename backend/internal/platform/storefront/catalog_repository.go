package storefront

import (
	"context"
	_ "embed"
	"fmt"
	"maps"
	"slices"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/platform/storefront/catalog"
)

//go:embed catalog_products.sql
var catalogSQL string

func (r *Repository) catalog(ctx context.Context, tenantID string, f catalog.Filter, loader catalog.RecipeLoader) ([]catalog.Product, error) {
	rows, err := r.store.Query(ctx, tenantID, catalogSQL, tenantID, f.Slug, f.ProductType, f.Collection, f.VariantIDs,
		f.Limit, f.Offset, f.Preview, f.RequireVariants, f.ExactCategory)
	if err != nil {
		return nil, err
	}
	products, err := pgx.CollectRows(rows, scanShopProduct)
	if err != nil {
		return nil, err
	}
	return products, applyRecipes(ctx, tenantID, products, f.Preview, loader)
}

func scanShopProduct(row pgx.CollectableRow) (catalog.Product, error) {
	var p catalog.Product
	err := row.Scan(&p.Name, &p.Slug, &p.SEOTitle, &p.SEODescription, &p.Composite, &p.Categories, &p.Images, &p.Variants)
	return p, err
}

func applyRecipes(ctx context.Context, tenantID string, products []catalog.Product, preview bool, loader catalog.RecipeLoader) error {
	targets := recipeVariants(products)
	if len(targets) == 0 {
		return nil
	}
	req := catalog.RecipeRequest{TenantID: tenantID, VariantIDs: slices.Collect(maps.Keys(targets)), Preview: preview}
	recipes, err := loader(ctx, req)
	if err != nil {
		return err
	}
	for id, v := range targets {
		recipe := recipes[id]
		if recipe == nil {
			return fmt.Errorf("missing catalog recipe for composite variant %s", id)
		}
		v.Bundle = recipe
	}
	return nil
}

func recipeVariants(products []catalog.Product) map[string]*catalog.Variant {
	out := map[string]*catalog.Variant{}
	for _, p := range products {
		if !p.Composite {
			continue
		}
		for i := range p.Variants {
			out[p.Variants[i].ID] = &p.Variants[i]
		}
	}
	return out
}
