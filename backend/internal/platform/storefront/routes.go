package storefront

import (
	"context"
	"net/http"
	"time"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

func Register(mux *http.ServeMux, store *db.Store, opts Options) error {
	h, err := newStoreHandler(NewRepository(store), opts)
	if err != nil {
		return err
	}
	mux.HandleFunc("GET /health", health(store))
	h.mount(mux, NewResolver(h.service).Middleware())
	return nil
}

func newStoreHandler(repo *Repository, opts Options) (*Handler, error) {
	registry, err := NewRegistry()
	if err != nil {
		return nil, err
	}
	return &Handler{service: NewService(repo, registry, opts), registry: registry, preview: newPreviewAuth(opts.JWTSecret)}, nil
}

func (h *Handler) mount(mux *http.ServeMux, mw auth.Middleware) {
	h.assets(mux)
	pages := map[string]shopHandler{
		"GET /{$}": h.home, "GET /products": h.list, "GET /products/{slug}": h.product,
		"GET /collections/{collection}": h.list, "GET /pages/{page}": h.custom, "GET /cart": h.cart,
		"GET /api/storefront/cart": h.currentCart, "PUT /api/storefront/cart": h.updateCart,
		"POST /api/storefront/quote": h.quote, "POST /api/storefront/checkout": h.checkout,
		"/": h.notFound,
	}
	for pattern, page := range pages {
		mux.Handle(pattern, mw(h.withShop(page)))
	}
	mux.Handle("GET /robots.txt", mw(http.HandlerFunc(h.Robots)))
	mux.Handle("GET /sitemap.xml", mw(http.HandlerFunc(h.Sitemap)))
}

func health(store *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := store.Ping(ctx); err != nil {
			httpx.WriteError(w, 503, "db_unavailable", "database unavailable", err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
