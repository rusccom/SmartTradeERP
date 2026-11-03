package storefront

import (
	"io/fs"
	"net/http"
	"path"
)

func (h *Handler) assets(mux *http.ServeMux) {
	mux.HandleFunc("GET /_assets/themes/{theme}/theme.css", h.ThemeCSS)
	mux.HandleFunc("GET /_assets/themes/{theme}/{file}", func(w http.ResponseWriter, r *http.Request) {
		theme := r.PathValue("theme")
		if !h.registry.Has(theme) {
			http.NotFound(w, r)
			return
		}
		serveAsset(w, r, path.Join("themes", theme, "assets"), r.PathValue("file"))
	})
	mux.HandleFunc("GET /_assets/{file}", func(w http.ResponseWriter, r *http.Request) {
		serveAsset(w, r, "assets", r.PathValue("file"))
	})
}

func serveAsset(w http.ResponseWriter, r *http.Request, dir, file string) {
	if !fs.ValidPath(dir) || !fs.ValidPath(file) || path.Base(file) != file || file == "." {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
	http.ServeFileFS(w, r, themesFS, path.Join(dir, file))
}
