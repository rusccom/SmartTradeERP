package storefront

import (
	"html/template"

	"smarterp/backend/internal/platform/storefront/catalog"
)

// LayoutVM is the shared chrome every themed page renders: brand, theme tokens,
// SEO head, and navigation. It is identical across themes so all shops share the
// same structure and differ only visually.
type LayoutVM struct {
	Brand    BrandVM
	Theme    ThemeVM
	SEO      SeoVM
	Nav      []NavItemVM
	Preview  bool
	Language string
}

type BrandVM struct {
	TenantName string
	LogoURL    string
	HomeURL    string
}

type ThemeVM struct {
	ID           string
	AssetVersion string
	TokenCSS     template.CSS
}

// SeoVM holds head metadata. JSONLD is a complete <script> element built by the
// core (never by a theme) so structured data cannot be omitted or malformed.
type SeoVM struct {
	Title       string
	Description string
	Canonical   string
	Robots      string
	OGType      string
	OGImage     string
	JSONLD      template.HTML
}

type NavItemVM struct {
	Label string
	URL   string
}

type CartVM struct {
	Layout LayoutVM
}

type HomeVM struct {
	Layout   LayoutVM
	Products []catalog.Card
	Sections []string
}

type ListVM struct {
	Title    string
	Layout   LayoutVM
	Products []catalog.Card
	Page     int
	HasPrev  bool
	HasNext  bool
	PrevURL  string
	NextURL  string
}

type ProductDetailVM struct {
	Selected    *catalog.VariantView
	Layout      LayoutVM
	Name        string
	Description string
	Images      []string
	Price       catalog.Money
	HasPrice    bool
	Variants    []catalog.VariantView
}

type NotFoundVM struct {
	Layout LayoutVM
}
