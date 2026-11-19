package catalog

type Product struct {
	Name, Slug, SEOTitle, SEODescription string
	Composite                            bool
	Categories                           []Category
	Images                               []string
	Variants                             []Variant
}

type Category struct{ Slug, Name, ParentSlug string }

type Variant struct {
	Option
	Images []string
	SKU    string
	Bundle Recipe
}

type Filter struct {
	Preview                        bool
	Slug, Collection, ProductType  string
	RequireVariants, ExactCategory bool
	VariantIDs                     []string
	Limit, Offset                  int
}
