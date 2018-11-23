package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/uuid"
)

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// FormHandler handles input from a form
func FormHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	c := &creds{}
	err := decoder.Decode(c)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(*c)
	fmt.Fprintln(w, "OK")
}

// UploadHandler deals with file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	f, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	fmt.Printf("Got file: %s with size: %d.\n", header.Filename, header.Size)
	u, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	targetFile, err := os.Create(u.String() + header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()
	io.Copy(targetFile, f)
	fmt.Fprintln(w, "OK")
}

// UploadIntrospectionHandler deals with file uploads for introspection purposes
func UploadIntrospectionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	fmt.Println("Values in the POST: ")
	for k := range r.MultipartForm.Value {
		fmt.Println(k)
	}
	fmt.Println("Files in the POST: ")
	for k := range r.MultipartForm.File {
		fmt.Println(k)
	}
	fmt.Fprintln(w, "OK")
}
