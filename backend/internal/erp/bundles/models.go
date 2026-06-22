package bundles

import "github.com/shopspring/decimal"

type Bundle struct {
	ProductID   string          `json:"product_id"`
	VariantID   string          `json:"variant_id"`
	ProductName string          `json:"product_name"`
	VariantName string          `json:"variant_name"`
	SKUCode     string          `json:"sku_code"`
	Unit        string          `json:"unit"`
	Price       decimal.Decimal `json:"price"`
	Components  []Component     `json:"components,omitempty"`
}

type Component struct {
	ComponentVariantID string          `json:"component_variant_id"`
	Qty                decimal.Decimal `json:"qty"`
	ProductName        string          `json:"product_name,omitempty"`
	VariantName        string          `json:"variant_name,omitempty"`
	Unit               string          `json:"unit,omitempty"`
}

type SnapshotInput struct {
	DocumentItemID string
	DocumentQty    decimal.Decimal
	Components     []Component
}

type componentInsert struct {
	ID        string
	VariantID string
	Item      Component
}
