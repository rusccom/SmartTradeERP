package storefront

import (
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/integrations/smartbundle"
)

// A shopper chooses a variant, quantity and optional component variants. The server prices them.
type cartItemInput struct {
	Picks     map[string]string `json:"picks,omitempty"`
	VariantID string            `json:"variant_id"`
	Qty       decimal.Decimal   `json:"qty"`
}

// cartLine is a server-priced line returned to the storefront cart.
type cartLine struct {
	Picks       map[string]string       `json:"picks,omitempty"`
	Components  []smartbundle.Component `json:"components,omitempty"`
	VariantID   string                  `json:"variant_id"`
	ProductName string                  `json:"product_name"`
	VariantName string                  `json:"variant_name"`
	Qty         string                  `json:"qty"`
	UnitPrice   string                  `json:"unit_price"`
	LineTotal   string                  `json:"line_total"`
	Available   string                  `json:"available"`
	Purchasable bool                    `json:"purchasable"`
}

type cartResult struct {
	Lines    []cartLine `json:"lines"`
	Total    string     `json:"total"`
	Currency string     `json:"currency"`
}

type checkoutCustomer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type checkoutRequest struct {
	Items    []cartItemInput  `json:"items"`
	Customer checkoutCustomer `json:"customer"`
	Note     string           `json:"note"`
}

type orderResult struct {
	Number   string `json:"number"`
	Total    string `json:"total"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}
