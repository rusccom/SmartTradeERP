package storefront

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/shopspring/decimal"

	"smarterp/backend/internal/platform/storefront/catalog"
)

const listPageSize = 24

func (sh *shop) list(r *http.Request) (ListVM, error) {
	page, collection := pageParam(r), r.PathValue("collection")
	if collection != "" && len(sh.service.registry.theme(sh.brand.ThemeID).manifest.Collections) == 0 {
		return ListVM{}, ErrProductNotFound
	}
	f := catalog.Filter{Collection: collection, Limit: listPageSize + 1, Offset: (page - 1) * listPageSize}
	items, err := sh.catalog(f)
	if err != nil {
		return ListVM{}, err
	}
	title, ok := sh.collectionTitle(collection, items)
	if !ok {
		return ListVM{}, ErrProductNotFound
	}
	return sh.listView(r, items, page, title), nil
}

func (sh *shop) listView(r *http.Request, items []catalog.Product, page int, title string) ListVM {
	seo := seoList(r, sh.brand)
	if title != "" {
		seo.Title = title + " · " + sh.brand.TenantName
	}
	next := len(items) > listPageSize
	if next {
		items = items[:listPageSize]
	}
	return ListVM{Layout: sh.layout(r, seo), Title: title, Products: sh.cards(items), Page: page, HasPrev: page > 1, HasNext: next,
		PrevURL: pageURL(r.URL.Path, page-1), NextURL: pageURL(r.URL.Path, page+1)}
}

func (sh *shop) collectionTitle(slug string, items []catalog.Product) (string, bool) {
	if slug == "" {
		return "", true
	}
	if title, ok := sh.service.registry.theme(sh.brand.ThemeID).manifest.Collections[slug]; ok {
		return title, true
	}
	for _, p := range items {
		for _, c := range p.Categories {
			if c.Slug == slug {
				return c.Name, true
			}
		}
	}
	return "", false
}

func (sh *shop) product(r *http.Request, p catalog.Product) ProductDetailVM {
	variants := sh.variantViews(p)
	vm := ProductDetailVM{Name: p.Name, Description: p.SEODescription, Variants: variants}
	keys := p.Images
	if len(variants) != 0 {
		i := selectedVariant(variants, r.URL.Query().Get("variant"))
		vm.Selected = &variants[i]
		vm.Price, vm.HasPrice = variants[i].Price, variants[i].HasPrice
		if len(p.Variants[i].Images) > 0 {
			keys = p.Variants[i].Images
		}
	}
	vm.Images = sh.service.imageList(keys, mainTransform)
	vm.Layout = sh.layout(r, seoProduct(r, p, firstString(vm.Images), vm.Price, vm.HasPrice))
	return vm
}

func selectedVariant(variants []catalog.VariantView, id string) int {
	for i, v := range variants {
		if v.VariantID == id {
			return i
		}
	}
	return 0
}

func (sh *shop) cards(products []catalog.Product) []catalog.Card {
	cards := make([]catalog.Card, 0, len(products))
	for _, p := range products {
		variants := sh.variantViews(p)
		price, ok := minPrice(variants)
		cards = append(cards, catalog.Card{Categories: p.Categories, Variants: variants, Name: p.Name, Slug: p.Slug, URL: productURL(p.Slug),
			ImageURL: sh.service.image(firstString(p.Images), cardTransform), Price: price, HasPrice: ok})
	}
	return cards
}

func (sh *shop) variantViews(p catalog.Product) []catalog.VariantView {
	list := make([]catalog.VariantView, 0, len(p.Variants))
	for _, v := range p.Variants {
		q, err := resolveVariant(v, nil, sh.brand.Currency.places)
		key := firstString(v.Images)
		if key == "" {
			key = firstString(p.Images)
		}
		view := catalog.VariantView{VariantID: v.ID, Name: v.Name, SKU: v.SKU, ImageURL: sh.service.image(key, cardTransform), HasPrice: err == nil,
			Price: formatMoney(q.Price, sh.brand.Currency), Purchasable: err == nil && availableUnits(q.Components).GreaterThanOrEqual(decimal.NewFromInt(1))}
		if v.Bundle != nil {
			view.Slots = v.Bundle.Slots
		}
		list = append(list, view)
	}
	return list
}

func minPrice(variants []catalog.VariantView) (catalog.Money, bool) {
	var best catalog.Money
	min := decimal.Zero
	for _, v := range variants {
		if !v.HasPrice {
			continue
		}
		value, _ := decimal.NewFromString(v.Price.Amount)
		if best.Amount == "" || value.LessThan(min) {
			best, min = v.Price, value
		}
	}
	return best, best.Amount != ""
}

func (s *Service) imageList(keys []string, transform string) []string {
	var out []string
	for _, key := range keys {
		if u := s.image(key, transform); u != "" {
			out = append(out, u)
		}
	}
	return out
}

func firstString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func productURL(slug string) string {
	if slug == "" {
		return "/products"
	}
	return "/products/" + url.PathEscape(slug)
}

func pageURL(path string, n int) string {
	if n <= 1 {
		return path
	}
	return path + "?page=" + strconv.Itoa(n)
}

func pageParam(r *http.Request) int {
	n, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || n < 1 || n > 1000000 {
		return 1
	}
	return n
}
