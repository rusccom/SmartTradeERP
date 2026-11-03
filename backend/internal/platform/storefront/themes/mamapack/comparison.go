package mamapack

import (
	"slices"
	"sort"

	"smarterp/backend/internal/platform/storefront/catalog"
)

type bagCard struct {
	Label string
	catalog.Card
	Slots []comparisonSlot
}
type comparisonSlot struct{ Key, Name, Qty string }
type comparisonRow struct {
	Name   string
	Values []string
}
type BagGroup struct {
	ID, Name string
	Bags     []bagCard
	Rows     []comparisonRow
}

func CompareBags(items []catalog.Card) []BagGroup {
	bySlug := map[string]bagCard{}
	for _, p := range items {
		bySlug[p.Slug] = bagCard{Card: p, Slots: comparisonSlots(p)}
	}
	groups := []BagGroup{}
	for _, guide := range Guides() {
		g := assembleGroup(guide, bySlug)
		if len(g.Bags) > 0 {
			groups = append(groups, g)
		}
	}
	return appendOtherBags(groups, items, bySlug)
}

func appendOtherBags(groups []BagGroup, items []catalog.Card, bySlug map[string]bagCard) []BagGroup {
	extra := Guide{ID: "pozostale", Name: "Pozostałe torby"}
	for _, p := range items {
		if catalog.HasCategory(p.Categories, BagCollection) {
			extra.Bags = append(extra.Bags, Link{Label: p.Name, Slug: p.Slug})
		}
	}
	if g := assembleGroup(extra, bySlug); len(g.Bags) > 0 {
		groups = append(groups, g)
	}
	return groups
}

func comparisonSlots(p catalog.Card) []comparisonSlot {
	if len(p.Variants) == 0 {
		return nil
	}
	var rows []comparisonSlot
	for _, slot := range p.Variants[0].Slots {
		rows = append(rows, comparisonSlot{Key: slot.Key, Name: slot.Name, Qty: slot.Qty.String()})
	}
	return rows
}

func assembleGroup(guide Guide, cards map[string]bagCard) BagGroup {
	g := BagGroup{ID: guide.ID, Name: guide.Name}
	for _, link := range guide.Bags {
		if card, ok := cards[link.Slug]; ok {
			card.Label = link.Label
			g.Bags = append(g.Bags, card)
			delete(cards, link.Slug)
		}
	}
	g.Rows = compareRows(g.Bags)
	return g
}

func compareRows(bags []bagCard) []comparisonRow {
	rows := map[string]*comparisonRow{}
	for i, bag := range bags {
		for _, slot := range bag.Slots {
			putComparison(rows, slot, i, len(bags))
		}
	}
	out := make([]comparisonRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, *row)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func putComparison(rows map[string]*comparisonRow, slot comparisonSlot, index, size int) {
	row := rows[slot.Key]
	if row == nil {
		row = &comparisonRow{Name: slot.Name, Values: slices.Repeat([]string{"—"}, size)}
		rows[slot.Key] = row
	}
	row.Values[index] = slot.Qty
}

func BagSlugs() []string {
	var slugs []string
	for _, guide := range Guides() {
		for _, bag := range guide.Bags {
			slugs = append(slugs, bag.Slug)
		}
	}
	return slugs
}
