package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arunsworld/go-app/handlers"
	"github.com/unrolled/secure"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	mux := mux.NewRouter()

	indexHandler := handlers.IndexHandler()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/form", indexHandler)
	mux.HandleFunc("/chat", indexHandler)
	mux.HandleFunc("/qr", indexHandler)
	mux.HandleFunc("/chatws", handlers.ChatWebSocketHandler)

	api := mux.PathPrefix("/api").Subrouter()
	api.HandleFunc("/choices", handlers.ChoicesHandler).Methods("GET")
	api.HandleFunc("/form-submit", handlers.FormHandler).Methods("POST")
	api.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")

	handlers.SetupStatic(mux)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "80"
	}
	fmt.Printf("Serving on port %s...\n", port)
	serve(secureMux(mux), fmt.Sprintf(":%s", port))
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
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}
	log.Fatal(srv.ListenAndServe())
}
