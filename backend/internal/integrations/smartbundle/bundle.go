package smartbundle

import (
	"errors"

	"github.com/shopspring/decimal"
)

var ErrSelection = errors.New("invalid or unavailable bundle selection")
var hundred = decimal.NewFromInt(100)

type Option struct {
	ID    string
	Name  string
	Price *decimal.Decimal
	Stock decimal.Decimal
}

type Slot struct {
	Key      string
	Name     string
	Qty      decimal.Decimal
	Discount decimal.Decimal
	Options  []Option
}

type Bundle struct {
	Mode     string
	Discount decimal.Decimal
	Slots    []Slot
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

func (b Bundle) Resolve(parent *decimal.Decimal, picks map[string]string, places int32) (Quote, error) {
	if len(b.Slots) == 0 || (len(picks) != 0 && len(picks) != len(b.Slots)) {
		return Quote{}, ErrSelection
	}
	result := Quote{}
	for _, slot := range b.Slots {
		option, ok := slot.choice(picks)
		if !ok || (b.Mode == "dynamic" && option.Price == nil) {
			return Quote{}, ErrSelection
		}
		result.Components = append(result.Components, slot.component(option))
		if b.Mode == "dynamic" {
			unit := option.Price.Mul(factor(slot.Discount)).Mul(factor(b.Discount)).Round(places)
			result.Price = result.Price.Add(unit.Mul(slot.Qty))
		}
	}
	return b.price(result, parent, places)
}

func (b Bundle) price(result Quote, parent *decimal.Decimal, places int32) (Quote, error) {
	switch b.Mode {
	case "fixed":
		if parent == nil {
			return Quote{}, ErrSelection
		}
		result.Price = parent.Mul(factor(b.Discount))
	case "dynamic":
	default:
		return Quote{}, ErrSelection
	}
	result.Price = result.Price.Round(places)
	return result, nil
}

func (s Slot) choice(picks map[string]string) (Option, bool) {
	if len(picks) == 0 && len(s.Options) > 0 {
		return s.Options[0], true
	}
	for _, option := range s.Options {
		if option.ID == picks[s.Key] {
			return option, true
		}
	}
	return Option{}, false
}

func (s Slot) component(option Option) Component {
	label := s.Name
	if option.Name != "" && option.Name != "Default" {
		label += " · " + option.Name
	}
	return Component{VariantID: option.ID, Name: label, Qty: s.Qty, Stock: option.Stock}
}

func factor(discount decimal.Decimal) decimal.Decimal {
	return hundred.Sub(discount).Shift(-2)
}
