package mamapack

import (
	"slices"
	"sort"

	"smarterp/backend/internal/platform/storefront/catalog"
)

type bagCard struct {
	catalog.Card
	Slots []comparisonSlot
}

// BagCollection is the existing public collection URL.
const BagCollection = "gotowe-torby-do-szpitala"

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
	groups := []BagGroup{}
	indices := map[string]int{}
	for _, p := range items {
		category := bagCategory(p.Categories)
		i, ok := indices[category.Slug]
		if !ok {
			i = len(groups)
			indices[category.Slug] = i
			groups = append(groups, BagGroup{ID: category.Slug, Name: category.Name})
		}
		groups[i].Bags = append(groups[i].Bags, bagCard{Card: p, Slots: comparisonSlots(p)})
	}
	for i := range groups {
		groups[i].Rows = compareRows(groups[i].Bags)
	}
	return groups
}

func bagCategory(categories []catalog.Category) catalog.Category {
	var root catalog.Category
	for _, c := range categories {
		if c.ParentSlug == BagCollection {
			return c
		}
		if c.Slug == BagCollection {
			root = c
		}
	}
	return root
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
