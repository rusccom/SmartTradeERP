package widget

import (
	_ "embed"
	"net/http"
)

//go:embed configurator.js
var configurator []byte

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /_assets/plugins/smartbundle/configurator.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		_, _ = w.Write(configurator)
	})
}
