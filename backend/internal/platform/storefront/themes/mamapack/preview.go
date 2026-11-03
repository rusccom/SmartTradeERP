package mamapack

import (
	"embed"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"smarterp/backend/internal/integrations/smartbundle"
	"smarterp/backend/internal/platform/storefront/catalog"
)

const BagCollection = "gotowe-torby-do-szpitala"
const mamaAssets = "/_assets/themes/mamapack/"

type Guide struct {
	ID, Name string
	Bags     []Link
}
type Link struct{ Label, Slug string }
type sampleItem struct {
	Key, Name, Slug string
	Category        catalog.Category
	Variants        []struct {
		Name  string
		Price decimal.Decimal
	}
}

//go:embed data/*.json
var examples embed.FS

func readThemeData[T any](name string) T {
	var out T
	raw, err := examples.ReadFile("data/" + name + ".json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(raw, &out); err != nil {
		panic(err)
	}
	return out
}

var Guides = sync.OnceValue(func() []Guide { return readThemeData[[]Guide]("bags") })
var samples = sync.OnceValue(buildSamples)

func SampleProducts(f catalog.Filter) []catalog.Product {
	var out []catalog.Product
	for _, p := range samples() {
		if f.Matches(p) {
			out = append(out, p)
		}
	}
	start := min(f.Offset, len(out))
	end := len(out)
	if f.Limit > 0 {
		end = min(end, start+f.Limit)
	}
	return out[start:end]
}

func buildSamples() []catalog.Product {
	byKey := make(map[string]catalog.Product)
	var out []catalog.Product
	for _, item := range readThemeData[[]sampleItem]("preview") {
		p := sampleProduct(item)
		byKey[item.Key] = p
		out = append(out, p)
	}
	for _, group := range Guides() {
		for level, bag := range group.Bags {
			out = append(out, sampleBag(group, bag, level, byKey))
		}
	}
	return out
}

func sampleID(key string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte("mamapack-preview/"+key)).String()
}

func sampleProduct(item sampleItem) catalog.Product {
	p := catalog.Product{Name: item.Name, Slug: item.Slug, Categories: []catalog.Category{item.Category},
		Images: []string{mamaAssets + item.Key + ".webp"}}
	for _, option := range item.Variants {
		v := smartbundle.Option{ID: sampleID(item.Key + option.Name), Name: option.Name, Price: &option.Price, Stock: decimal.NewFromInt(36)}
		p.Variants = append(p.Variants, catalog.Variant{Option: v})
	}
	return p
}

func sampleBag(group Guide, bag Link, level int, items map[string]catalog.Product) catalog.Product {
	recipe := &smartbundle.Bundle{Mode: "dynamic", Discount: decimal.NewFromInt(5)}
	for _, key := range sampleContents(group.ID, level) {
		recipe.Slots = append(recipe.Slots, sampleSlot(items[key], level))
	}
	option := smartbundle.Option{ID: sampleID(bag.Slug), Name: "Default"}
	if level == 0 {
		recipe.Mode = "fixed"
		amount := decimal.NewFromInt(39)
		option.Price = &amount
	}
	return catalog.Product{Name: "Torba " + bag.Label + " · " + group.Name, Slug: bag.Slug, Composite: true,
		Images: []string{mamaAssets + "bag-" + group.ID + ".webp"}, Categories: []catalog.Category{{Slug: BagCollection, Name: "Gotowe torby do szpitala"}},
		Variants: []catalog.Variant{{Option: option, Bundle: recipe}}}
}

func sampleContents(group string, level int) []string {
	keys := map[string][]string{"mama": {"pads", "briefs", "liners", "wipes", "pacifier"},
		"dziecko": {"nappies", "wipes", "pacifier", "pads", "liners"}, "razem": {"pads", "nappies", "wipes", "briefs", "pacifier", "liners"}}
	return keys[group][:min(2+level, len(keys[group]))]
}

func sampleSlot(p catalog.Product, level int) smartbundle.Slot {
	slot := smartbundle.Slot{Key: p.Variants[0].ID, Name: p.Name, Qty: decimal.NewFromInt(1)}
	if level > 1 {
		slot.Qty = decimal.NewFromInt(2)
		slot.Discount = decimal.NewFromInt(10)
	}
	for _, v := range p.Variants {
		slot.Options = append(slot.Options, v.Option)
	}
	return slot
}

func SampleCategories() []catalog.Category {
	var out []catalog.Category
	seen := map[string]bool{}
	for _, p := range SampleProducts(catalog.Filter{Kind: "builder"}) {
		for _, c := range p.Categories {
			if !seen[c.Slug] {
				out = append(out, c)
				seen[c.Slug] = true
			}
		}
	}
	return out
}
