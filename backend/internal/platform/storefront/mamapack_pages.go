package storefront

import (
	"net/http"

	"smarterp/backend/internal/platform/storefront/catalog"
	"smarterp/backend/internal/platform/storefront/themes/mamapack"
)

type builderVM struct {
	Layout               LayoutVM
	Groups               []mamapack.BuilderGroup
	Categories           []catalog.Category
	Selected, Prev, Next string
}
type bagsVM struct {
	Layout LayoutVM
	Groups []mamapack.BagGroup
}

func (h *Handler) custom(w http.ResponseWriter, r *http.Request, sh *shop) {
	page := h.registry.theme(sh.brand.ThemeID).manifest.Pages[r.URL.Path]
	switch page {
	case "builder":
		vm, err := sh.builder(r)
		h.page(w, vm.Layout, page, vm, err)
	case "bags":
		items, err := sh.catalog(catalog.Filter{Kind: "bags"})
		vm := bagsVM{Layout: sh.pageLayout(r, "Gotowe torby do szpitala dla mamy i noworodka"), Groups: mamapack.CompareBags(sh.cards(items))}
		h.page(w, vm.Layout, page, vm, err)
	default:
		h.notFound(w, r, sh)
	}
}

func (sh *shop) pageLayout(r *http.Request, title string) LayoutVM {
	seo := seoHome(r, sh.brand)
	seo.Title = title + " · " + sh.brand.TenantName
	return sh.layout(r, seo)
}
