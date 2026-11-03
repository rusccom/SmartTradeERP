package catalog

import "smarterp/backend/internal/integrations/smartbundle"

// Money keeps a display string for the theme plus a raw amount and ISO code
// for structured data, so themes never format currency themselves.
type Money struct {
	Display string
	Amount  string
	Code    string
}

type Card struct {
	Categories []Category
	Variants   []VariantView
	Slug       string
	Name       string
	URL        string
	ImageURL   string
	Price      Money
	HasPrice   bool
}

type VariantView struct {
	Purchasable bool
	Slots       []smartbundle.Slot
	VariantID   string
	Name        string
	SKU         string
	Price       Money
	HasPrice    bool
}
