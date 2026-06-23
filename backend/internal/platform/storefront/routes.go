package storefront

import (
	"context"
	"net/http"
	"time"

	"smarterp/backend/internal/shared/auth"
	"smarterp/backend/internal/shared/db"
	"smarterp/backend/internal/shared/httpx"
)

const healthDBTimeout = 2 * time.Second

// Register wires the public, tenant-by-Host storefront routes onto mux. Pages
// resolve the tenant from the request Host; /health and theme assets do not.
func Register(mux *http.ServeMux, store *db.Store, opts Options) error {
	registry, err := NewRegistry()
	if err != nil {
		return err
	}
	repo := NewRepository(store)
	service := NewService(repo, registry, opts)
	resolver := NewResolver(service)
	handler := NewHandler(service, registry, newPreviewAuth(opts.JWTSecret))
	commerce := NewCommerceHandler(NewCommerce(store, repo))
	mux.HandleFunc("GET /health", health(store))
	mux.HandleFunc("GET /_assets/themes/{theme}/theme.css", handler.ThemeCSS)
	mux.HandleFunc("GET /_assets/storefront.js", handler.AppJS)
	mw := resolver.Middleware()
	registerPages(mux, handler, mw)
	registerCommerce(mux, commerce, mw)
	return nil
}

func registerCommerce(mux *http.ServeMux, handler *CommerceHandler, mw auth.Middleware) {
	mux.Handle("POST /api/storefront/cart", auth.Chain(http.HandlerFunc(handler.Cart), mw))
	mux.Handle("POST /api/storefront/checkout", auth.Chain(http.HandlerFunc(handler.Checkout), mw))
}

func registerPages(mux *http.ServeMux, handler *Handler, mw auth.Middleware) {
	mux.Handle("GET /{$}", auth.Chain(http.HandlerFunc(handler.Home), mw))
	mux.Handle("GET /products", auth.Chain(http.HandlerFunc(handler.ProductList), mw))
	mux.Handle("GET /products/{slug}", auth.Chain(http.HandlerFunc(handler.ProductDetail), mw))
	mux.Handle("GET /cart", auth.Chain(http.HandlerFunc(handler.Cart), mw))
	mux.Handle("GET /robots.txt", auth.Chain(http.HandlerFunc(handler.Robots), mw))
	mux.Handle("GET /sitemap.xml", auth.Chain(http.HandlerFunc(handler.Sitemap), mw))
	mux.Handle("/", auth.Chain(http.HandlerFunc(handler.NotFound), mw))
}

func health(store *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), healthDBTimeout)
		defer cancel()
		if err := store.Ping(ctx); err != nil {
			httpx.WriteError(w, http.StatusServiceUnavailable, "db_unavailable", "database unavailable", nil)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
