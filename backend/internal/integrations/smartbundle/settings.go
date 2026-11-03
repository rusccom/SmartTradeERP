package smartbundle

import (
	"slices"

	"github.com/shopspring/decimal"
)

type Rule struct {
	Discount decimal.Decimal `json:"discount"`
	Options  []string        `json:"options"`
}

type Settings struct {
	Mode     string          `json:"mode"`
	Discount decimal.Decimal `json:"discount"`
	Slots    map[string]Rule `json:"slots"`
}

type Alternate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EditableSlot struct {
	Key     string      `json:"key"`
	Name    string      `json:"name"`
	Options []Alternate `json:"options"`
}

type Editor struct {
	Settings   Settings       `json:"settings"`
	Components []EditableSlot `json:"components"`
}

func (s Settings) valid(components []EditableSlot) bool {
	if s.Slots == nil || (s.Mode != "fixed" && s.Mode != "dynamic") || !percent(s.Discount) {
		return false
	}
	available := make(map[string][]Alternate, len(components))
	for _, slot := range components {
		available[slot.Key] = slot.Options
	}
	for id, rule := range s.Slots {
		options, ok := available[id]
		if !ok || !validRule(rule, options) {
			return false
		}
	}
	return true
}

func validRule(rule Rule, options []Alternate) bool {
	if !percent(rule.Discount) {
		return false
	}
	for _, id := range rule.Options {
		if !slices.ContainsFunc(options, func(o Alternate) bool { return o.ID == id }) {
			return false
		}
	}
	return true
}

func percent(value decimal.Decimal) bool {
	return !value.IsNegative() && value.LessThanOrEqual(hundred) && value.Equal(value.Round(2))
}
