package smartbundle

import (
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/platform/storefront/catalog"
)

const Key = "smartbundle"

var hundred = decimal.NewFromInt(100)

type Slot struct {
	catalog.Slot
	Discount decimal.Decimal
}

type Bundle struct {
	Mode     string
	Discount decimal.Decimal
	Slots    []Slot
}

func (b Bundle) Choices() []catalog.Slot {
	choices := make([]catalog.Slot, 0, len(b.Slots))
	for _, s := range b.Slots {
		choices = append(choices, s.Slot)
	}
	return choices
}

func (b Bundle) Resolve(parent *decimal.Decimal, picks map[string]string, places int32) (catalog.Quote, error) {
	if len(b.Slots) == 0 || (len(picks) != 0 && len(picks) != len(b.Slots)) {
		return catalog.Quote{}, catalog.ErrSelection
	}
	result := catalog.Quote{}
	for _, slot := range b.Slots {
		option, ok := slot.choice(picks)
		if !ok || (b.Mode == "dynamic" && option.Price == nil) {
			return catalog.Quote{}, catalog.ErrSelection
		}
		result.Components = append(result.Components, slot.component(option))
		if b.Mode == "dynamic" {
			unit := option.Price.Mul(factor(slot.Discount)).Mul(factor(b.Discount)).Round(places)
			result.Price = result.Price.Add(unit.Mul(slot.Qty))
		}
	}
	return b.price(result, parent, places)
}

func (b Bundle) price(result catalog.Quote, parent *decimal.Decimal, places int32) (catalog.Quote, error) {
	switch b.Mode {
	case "fixed":
		if parent == nil {
			return catalog.Quote{}, catalog.ErrSelection
		}
		result.Price = parent.Mul(factor(b.Discount))
	case "dynamic":
	default:
		return catalog.Quote{}, catalog.ErrSelection
	}
	result.Price = result.Price.Round(places)
	return result, nil
}

func (s Slot) choice(picks map[string]string) (catalog.Option, bool) {
	if len(picks) == 0 && len(s.Options) > 0 {
		return s.Options[0], true
	}
	for _, option := range s.Options {
		if option.ID == picks[s.Key] {
			return option, true
		}
	}
	return catalog.Option{}, false
}

func (s Slot) component(option catalog.Option) catalog.Component {
	label := s.Name
	if option.Name != "" && option.Name != "Default" {
		label += " · " + option.Name
	}
	return catalog.Component{VariantID: option.ID, Name: label, Qty: s.Qty, Stock: option.Stock}
}

func factor(discount decimal.Decimal) decimal.Decimal {
	return hundred.Sub(discount).Shift(-2)
}
