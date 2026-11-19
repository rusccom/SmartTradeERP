package storefront

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/platform/storefront/catalog"
)

//go:embed catalog_bundles.sql
var baseRecipesSQL string

type fixedRecipe struct {
	components  []catalog.Component
	unavailable bool
}

func (b *fixedRecipe) Choices() []catalog.Slot { return nil }

func (b *fixedRecipe) Resolve(price *decimal.Decimal, picks map[string]string, places int32) (catalog.Quote, error) {
	if price == nil || len(picks) > 0 || len(b.components) == 0 || b.unavailable {
		return catalog.Quote{}, catalog.ErrSelection
	}
	return catalog.Quote{Price: price.Round(places), Components: b.components}, nil
}

func (r *Repository) fixedRecipes(ctx context.Context, req catalog.RecipeRequest) (map[string]catalog.Recipe, error) {
	rows, err := r.store.Query(ctx, req.TenantID, baseRecipesSQL, req.TenantID, req.VariantIDs, req.Preview)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]catalog.Recipe, len(req.VariantIDs))
	for _, id := range req.VariantIDs {
		result[id] = &fixedRecipe{}
	}
	for rows.Next() {
		if err := readComponent(rows, result); err != nil {
			return nil, err
		}
	}
	return result, rows.Err()
}

func readComponent(row pgx.Row, recipes map[string]catalog.Recipe) error {
	var id, variant string
	var ready bool
	var c catalog.Component
	if err := row.Scan(&id, &c.VariantID, &c.Name, &variant, &c.Qty, &c.Stock, &ready); err != nil {
		return err
	}
	if variant != "" && variant != "Default" {
		c.Name += " · " + variant
	}
	b := recipes[id].(*fixedRecipe)
	b.components = append(b.components, c)
	b.unavailable = b.unavailable || !ready
	return nil
}
