package storefront

import "net/http"

func (h *Handler) Robots(w http.ResponseWriter, r *http.Request) {
	body := "User-agent: *\nAllow: /\nSitemap: " + absoluteURL(r, "/sitemap.xml") + "\n"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write([]byte(body))
}

func (h *Handler) Sitemap(w http.ResponseWriter, r *http.Request) {
	body, err := h.service.BuildSitemap(r)
	if err != nil {
		h.serverError(w)
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")
	_, _ = w.Write(body)
}

func (h *Handler) ThemeCSS(w http.ResponseWriter, r *http.Request) {
	css, _, ok := h.registry.CSS(r.PathValue("theme"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	_, _ = w.Write(css)
}
