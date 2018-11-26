package handlers

import (
	"log"
	"mime"
	"net/http"

	"github.com/arunsworld/go-app/assets"

	"github.com/gorilla/mux"
)

const staticPrefix = "/static/"

// SetupStatic sets up the static routes
func SetupStatic(m *mux.Router) {
	mime.AddExtensionType(".map", "text/plain")
	mime.AddExtensionType(".woff2", "font/woff2")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".ttf", "font/ttf")

	staticFS := assets.StaticFiles()
	m.PathPrefix(staticPrefix).Handler(http.StripPrefix(staticPrefix, http.FileServer(staticFS.FS())))
}

// IndexHandler serves /
func IndexHandler() http.HandlerFunc {
	templatesFS := assets.TemplatesFiles()
	index, err := templatesFS.ReadFile("/index.html")
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	}
}
