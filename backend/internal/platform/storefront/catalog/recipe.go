package catalog

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

var ErrSelection = errors.New("invalid or unavailable selection")

type Option struct {
	ID    string
	Name  string
	Price *decimal.Decimal
	Stock decimal.Decimal
}

type Slot struct {
	Key     string
	Name    string
	Qty     decimal.Decimal
	Options []Option
}

type Component struct {
	Stock     decimal.Decimal `json:"-"`
	VariantID string          `json:"variant_id"`
	Name      string          `json:"name"`
	Qty       decimal.Decimal `json:"qty"`
}

type Quote struct {
	Price      decimal.Decimal
	Components []Component
}

type Recipe interface {
	Resolve(*decimal.Decimal, map[string]string, int32) (Quote, error)
	Choices() []Slot
}

type RecipeRequest struct {
	TenantID   string
	VariantIDs []string
	Preview    bool
}

type RecipeLoader func(context.Context, RecipeRequest) (map[string]Recipe, error)

type RecipePlugin struct {
	Key  string
	Load RecipeLoader
}
