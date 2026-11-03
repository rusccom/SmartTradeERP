package catalog

import "smarterp/backend/internal/integrations/smartbundle"

type Product struct {
	Name, Slug, SEOTitle, SEODescription string
	Composite                            bool
	Categories                           []Category
	Images                               []string
	Variants                             []Variant
}

type Category struct{ Slug, Name, ParentSlug string }

type Variant struct {
	smartbundle.Option
	Images []string
	SKU    string
	Bundle *smartbundle.Bundle
}

type Filter struct {
	Preview                bool
	Slug, Collection, Kind string
	VariantIDs             []string
	Limit, Offset          int
}
