package storefront

import "html/template"

// LayoutVM is the shared chrome every themed page renders: brand, theme tokens,
// SEO head, and navigation. It is identical across themes so all shops share the
// same structure and differ only visually.
type LayoutVM struct {
	Brand BrandVM
	Theme ThemeVM
	SEO   SeoVM
	Nav   []NavItemVM
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

// MoneyVM keeps a display string for the theme plus a raw amount and ISO code
// for structured data, so themes never format currency themselves.
type MoneyVM struct {
	Display string
	Amount  string
	Code    string
}

type ProductCardVM struct {
	Name     string
	URL      string
	ImageURL string
	Price    MoneyVM
	HasPrice bool
}

type ProductVariantVM struct {
	VariantID string
	Name      string
	SKU       string
	Price     MoneyVM
	HasPrice  bool
}

type CartVM struct {
	Layout LayoutVM
}

type HomeVM struct {
	Layout   LayoutVM
	Products []ProductCardVM
	Sections []string
}

type ListVM struct {
	Layout   LayoutVM
	Products []ProductCardVM
	Page     int
	HasPrev  bool
	HasNext  bool
	PrevURL  string
	NextURL  string
}

type ProductDetailVM struct {
	Layout      LayoutVM
	Name        string
	Description string
	Images      []string
	Price       MoneyVM
	HasPrice    bool
	Variants    []ProductVariantVM
}

type NotFoundVM struct {
	Layout LayoutVM
}
