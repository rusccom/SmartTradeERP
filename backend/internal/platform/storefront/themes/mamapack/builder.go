package mamapack

import "smarterp/backend/internal/platform/storefront/catalog"

type BuilderGroup struct {
	Name, ID string
	Products []catalog.Card
}

func BuilderGroups(items []catalog.Card, selected string) []BuilderGroup {
	groups := []BuilderGroup{}
	indices := map[string]int{}
	for _, p := range items {
		if len(p.Variants) == 0 {
			continue
		}
		category := builderCategory(p, selected)
		index, ok := indices[category.Slug]
		if !ok {
			index = len(groups)
			indices[category.Slug] = index
			groups = append(groups, BuilderGroup{Name: category.Name, ID: category.Slug})
		}
		groups[index].Products = append(groups[index].Products, p)
	}
	return groups
}

func builderCategory(p catalog.Card, selected string) catalog.Category {
	for _, c := range p.Categories {
		if c.Slug == selected {
			return c
		}
	}
	if len(p.Categories) > 0 {
		return p.Categories[0]
	}
	return catalog.Category{Slug: "pozostale", Name: "Pozostałe produkty"}
}
