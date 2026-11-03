package storefront

import (
	"errors"
	"net/http"

	"smarterp/backend/internal/platform/storefront/catalog"
)

type Handler struct {
	service  *Service
	registry *Registry
	preview  *previewAuth
}

type shopHandler func(http.ResponseWriter, *http.Request, *shop)

func (h *Handler) withShop(next shopHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Cookie")
		preview, ok := true, true
		if h.service.repo != nil {
			preview, ok = h.preview.session(w, r)
		}
		if !ok {
			return
		}
		sh, err := h.service.open(r, preview)
		if err != nil {
			http.Error(w, "Store unavailable", 500)
			return
		}
		cachePage(w, r, sh.preview)
		next(w, r, sh)
	}
}

func cachePage(w http.ResponseWriter, r *http.Request, preview bool) {
	if preview {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	}
	if preview || r.Method != http.MethodGet {
		w.Header().Set("Cache-Control", "no-store")
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60, s-maxage=300")
	w.Header().Set("Cache-Tag", normalizeHost(r.Host))
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request, sh *shop) {
	products, err := sh.catalog(catalog.Filter{Limit: 12})
	vm := HomeVM{Layout: sh.layout(r, seoHome(r, sh.brand)), Sections: homeSections(sh.brand.Sections), Products: sh.cards(products)}
	h.page(w, vm.Layout, "home", vm, err)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request, sh *shop) {
	vm, err := sh.list(r)
	if errors.Is(err, ErrProductNotFound) {
		h.notFound(w, r, sh)
		return
	}
	h.page(w, vm.Layout, "list", vm, err)
}

func (h *Handler) product(w http.ResponseWriter, r *http.Request, sh *shop) {
	items, err := sh.catalog(catalog.Filter{Slug: r.PathValue("slug"), Limit: 1})
	if err != nil {
		http.Error(w, "Catalog unavailable", 500)
		return
	}
	if len(items) == 0 {
		h.notFound(w, r, sh)
		return
	}
	vm := sh.product(r, items[0])
	h.page(w, vm.Layout, "product", vm, nil)
}

func (h *Handler) cart(w http.ResponseWriter, r *http.Request, sh *shop) {
	w.Header().Set("Cache-Control", "no-store")
	vm := CartVM{Layout: sh.layout(r, seoCart(r, sh.brand))}
	h.page(w, vm.Layout, "cart", vm, nil)
}

func (h *Handler) page(w http.ResponseWriter, l LayoutVM, page string, vm any, err error) {
	if err != nil {
		http.Error(w, "Page unavailable", 500)
		return
	}
	h.registry.Render(w, http.StatusOK, l.Theme.ID, page, vm)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request, sh *shop) {
	vm := NotFoundVM{Layout: sh.layout(r, seoNotFound(r, sh.brand))}
	h.registry.Render(w, http.StatusNotFound, vm.Layout.Theme.ID, "404", vm)
}

func writeNotFound(w http.ResponseWriter) { http.Error(w, "Store not found", http.StatusNotFound) }
