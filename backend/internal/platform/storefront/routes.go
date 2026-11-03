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

// LocalPreview runs the normal shop handlers against the theme's sample catalog.
func LocalPreview() (http.Handler, error) {
	h, err := newStoreHandler(nil, Options{})
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	h.mount(mux, func(next http.Handler) http.Handler { return next })
	return mux, nil
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
		"POST /api/storefront/cart": h.quote, "POST /api/storefront/checkout": h.checkout,
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
