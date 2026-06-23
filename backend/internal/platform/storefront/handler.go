package storefront

import (
	"errors"
	"net/http"
)

type Handler struct {
	service  *Service
	registry *Registry
	preview  *previewAuth
}

func NewHandler(service *Service, registry *Registry, preview *previewAuth) *Handler {
	return &Handler{service: service, registry: registry, preview: preview}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	preview := h.preview.check(r)
	vm, err := h.service.BuildHome(r, preview)
	if err != nil {
		h.serverError(w)
		return
	}
	h.setCache(w, r, preview)
	h.registry.Render(w, http.StatusOK, vm.Layout.Theme.ID, "home", vm)
}

func (h *Handler) ProductList(w http.ResponseWriter, r *http.Request) {
	preview := h.preview.check(r)
	vm, err := h.service.BuildList(r, preview)
	if err != nil {
		h.serverError(w)
		return
	}
	h.setCache(w, r, preview)
	h.registry.Render(w, http.StatusOK, vm.Layout.Theme.ID, "list", vm)
}

func (h *Handler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	preview := h.preview.check(r)
	vm, err := h.service.BuildProduct(r, r.PathValue("slug"), preview)
	if err != nil {
		h.handleProductError(w, r, err)
		return
	}
	h.setCache(w, r, preview)
	h.registry.Render(w, http.StatusOK, vm.Layout.Theme.ID, "product", vm)
}

func (h *Handler) handleProductError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, ErrProductNotFound) {
		h.renderNotFound(w, r)
		return
	}
	h.serverError(w)
}

func (h *Handler) Cart(w http.ResponseWriter, r *http.Request) {
	vm, err := h.service.BuildCart(r, h.preview.check(r))
	if err != nil {
		h.serverError(w)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	h.registry.Render(w, http.StatusOK, vm.Layout.Theme.ID, "cart", vm)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.renderNotFound(w, r)
}

func (h *Handler) renderNotFound(w http.ResponseWriter, r *http.Request) {
	preview := h.preview.check(r)
	vm, err := h.service.BuildNotFound(r, preview)
	if err != nil {
		writeNotFound(w)
		return
	}
	h.registry.Render(w, http.StatusNotFound, vm.Layout.Theme.ID, "404", vm)
}

func (h *Handler) serverError(w http.ResponseWriter) {
	http.Error(w, "internal error", http.StatusInternalServerError)
}

// setCache lets the edge cache published pages but never previews (a draft is
// per-owner and must not be shared or served stale).
func (h *Handler) setCache(w http.ResponseWriter, r *http.Request, preview bool) {
	if preview {
		w.Header().Set("Cache-Control", "no-store")
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60, s-maxage=300, stale-while-revalidate=86400")
	w.Header().Set("Cache-Tag", normalizeHost(r.Host))
}

func writeNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(notFoundHTML))
}

const notFoundHTML = `<!doctype html><html lang="en"><head><meta charset="utf-8">` +
	`<title>Not found</title></head><body><h1>404 — store not found</h1></body></html>`
