package storefront

import (
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"smarterp/backend/internal/platform/storefront/catalog"
)

type customPageVM struct {
	Layout               LayoutVM
	Products             []catalog.Card
	Categories           []catalog.Category
	Selected, Prev, Next string
}

func (p PageDefinition) available(plugins []string) bool {
	return p.Plugin == "" || slices.Contains(plugins, p.Plugin)
}

func (h *Handler) custom(w http.ResponseWriter, r *http.Request, sh *shop) {
	p, ok := h.registry.theme(sh.brand.ThemeID).manifest.Pages[r.URL.Path]
	if !ok || !p.available(sh.brand.Plugins) {
		h.notFound(w, r, sh)
		return
	}
	vm, err := sh.customPage(r, p)
	h.page(w, vm.Layout, p.Template, vm, err)
}

func (sh *shop) customPage(r *http.Request, p PageDefinition) (customPageVM, error) {
	seo := seoHome(r, sh.brand)
	seo.Title = p.Title + " · " + sh.brand.TenantName
	vm := customPageVM{Layout: sh.layout(r, seo)}
	filter := p.Catalog.filter(r)
	filter.Preview = sh.preview
	products, err := sh.catalog(filter)
	if err != nil {
		return vm, err
	}
	vm.Products = sh.cards(p.Catalog.paginate(r, products, &vm))
	if p.Catalog.CategoryParam != "" {
		vm.Selected = filter.Collection
		vm.Categories, err = sh.service.repo.categories(sh.ctx, sh.tenantID, filter)
	}
	return vm, err
}

func (c PageCatalog) filter(r *http.Request) catalog.Filter {
	f := catalog.Filter{ProductType: c.ProductType, Collection: c.Collection, RequireVariants: c.RequireVariants}
	if c.CategoryParam != "" {
		f.Collection = r.URL.Query().Get(c.CategoryParam)
		f.ExactCategory = true
	}
	if c.PageSize > 0 {
		f.Offset = (pageParam(r) - 1) * c.PageSize
		f.Limit = c.PageSize + 1
	}
	return f
}

func (c PageCatalog) paginate(r *http.Request, items []catalog.Product, vm *customPageVM) []catalog.Product {
	if c.PageSize <= 0 {
		return items
	}
	page := pageParam(r)
	if page > 1 {
		vm.Prev = c.pageURL(r, page-1)
	}
	if len(items) > c.PageSize {
		vm.Next = c.pageURL(r, page+1)
		return items[:c.PageSize]
	}
	return items
}

func (c PageCatalog) pageURL(r *http.Request, page int) string {
	q := url.Values{}
	if value := r.URL.Query().Get(c.CategoryParam); c.CategoryParam != "" && value != "" {
		q.Set(c.CategoryParam, value)
	}
	if page > 1 {
		q.Set("page", strconv.Itoa(page))
	}
	return (&url.URL{Path: r.URL.Path, RawQuery: q.Encode()}).String()
}
