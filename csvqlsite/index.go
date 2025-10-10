package csvqlsite

import (
	"embed"
	"net/http"
)

//go:embed index.html index.js index.css
var assets embed.FS

var Handler = http.FileServerFS(assets)

func Route(mux *http.ServeMux) {
	mux.Handle("/", Handler)
}
