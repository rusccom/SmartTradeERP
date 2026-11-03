package catalog

import (
	"slices"

	"smarterp/backend/internal/integrations/smartbundle"
)

type Product struct {
	Name, Slug, SEOTitle, SEODescription string
	Composite                            bool
	Categories                           []Category
	Images                               []string
	Variants                             []Variant
}

type Category struct{ Slug, Name string }

type Variant struct {
	smartbundle.Option
	SKU    string
	Bundle *smartbundle.Bundle
}

type Filter struct {
	Slug, Collection, Kind string
	VariantIDs             []string
	CollectionSlugs        []string
	Limit, Offset          int
}

func (p Product) Belongs(slug string) bool {
	return HasCategory(p.Categories, slug)
}

func HasCategory(categories []Category, slug string) bool {
	return slices.ContainsFunc(categories, func(c Category) bool { return c.Slug == slug })
}

func (f Filter) Matches(p Product) bool {
	if f.Slug != "" && f.Slug != p.Slug || f.Collection != "" && !p.Belongs(f.Collection) && !slices.Contains(f.CollectionSlugs, p.Slug) {
		return false
	}
	if f.Kind == "builder" && p.Composite || f.Kind == "bags" && !p.Composite {
		return false
	}
	if f.VariantIDs == nil {
		return true
	}
	return slices.ContainsFunc(p.Variants, func(v Variant) bool { return slices.Contains(f.VariantIDs, v.ID) })
}
