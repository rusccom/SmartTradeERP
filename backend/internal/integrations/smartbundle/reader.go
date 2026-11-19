package smartbundle

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"

	"smarterp/backend/internal/platform/storefront/catalog"
	"smarterp/backend/internal/shared/db"
)

//go:embed recipes.sql
var recipesSQL string

func Loader(store *db.Store) catalog.RecipePlugin {
	return catalog.RecipePlugin{Key: Key, Load: func(ctx context.Context, req catalog.RecipeRequest) (map[string]catalog.Recipe, error) {
		rows, err := store.Query(ctx, req.TenantID, recipesSQL, req.TenantID, req.VariantIDs, req.Preview)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return readRecipes(rows, req.VariantIDs)
	}}
}

func readRecipes(rows pgx.Rows, ids []string) (map[string]catalog.Recipe, error) {
	result := make(map[string]catalog.Recipe, len(ids))
	for _, id := range ids {
		result[id] = &Bundle{Mode: "fixed"}
	}
	for rows.Next() {
		if err := readSlot(rows, result); err != nil {
			return nil, err
		}
	}
	return result, rows.Err()
}

func readSlot(row pgx.Row, recipes map[string]catalog.Recipe) error {
	var id string
	var b Bundle
	var s Slot
	if err := row.Scan(&id, &b.Mode, &b.Discount, &s.Key, &s.Name, &s.Qty, &s.Discount, &s.Options); err != nil {
		return err
	}
	recipe := recipes[id].(*Bundle)
	recipe.Mode, recipe.Discount = b.Mode, b.Discount
	recipe.Slots = append(recipe.Slots, s)
	return nil
}
