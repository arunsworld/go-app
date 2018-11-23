package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arunsworld/go-app/handlers"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/unrolled/secure"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	box := packr.New("templates", "./assets/templates/")
	staticBox := packr.New("static", "./assets/static/")

	mux := mux.NewRouter()

	mux.HandleFunc("/", handlers.IndexHandler(box))
	mux.HandleFunc("/form", handlers.IndexHandler(box))

	api := mux.PathPrefix("/api").Subrouter()
	api.HandleFunc("/form-submit", handlers.FormHandler).Methods("POST")
	api.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")

	handlers.SetupStatic(mux, staticBox)

	fmt.Println("Serving on port 9095...")
	serve(secureMux(mux), ":9095")
}

func secureMux(mux *mux.Router) http.Handler {
	c := cors.New(cors.Options{})

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
	})

	handler := c.Handler(mux)
	handler = secureMiddleware.Handler(handler)

	handler = gziphandler.GzipHandler(handler)

	return handler
}

func serve(handler http.Handler, address string) {
	srv := http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	log.Fatal(srv.ListenAndServe())
}