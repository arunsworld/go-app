package handlers

import (
	"log"
	"mime"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

const staticPrefix = "/static/"

// SetupStatic sets up the static routes
func SetupStatic(m *mux.Router, staticBox *packr.Box) {
	mime.AddExtensionType(".map", "text/plain")

	m.PathPrefix(staticPrefix).Handler(http.StripPrefix(staticPrefix, http.FileServer(staticBox)))
}

// IndexHandler serves /
func IndexHandler(box *packr.Box) http.HandlerFunc {
	index, err := box.Find("index.html")
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	}
}
