package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

type creds struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Datetime time.Time `json:"datetime"`
	Choice   string    `json:"choice"`
}

// ChoicesHandler returns an array of choices
func ChoicesHandler(w http.ResponseWriter, r *http.Request) {
	result := []map[string]string{
		{"id": "1", "text": "Assamese"},
		{"id": "2", "text": "Gujarati"},
		{"id": "3", "text": "Kannada"},
		{"id": "4", "text": "Kashmiri"},
	}
	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

// FormHandler handles input from a form
func FormHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	c := &creds{}
	err := decoder.Decode(c)
	if err != nil {
		log.Println("Invalid JSON being sent to FormHandler. Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Printf("%+v\n", *c)
	fmt.Fprintln(w, "OK")
}

// UploadHandler deals with file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	f, header, err := r.FormFile("file")
	if err != nil {
		log.Println("Bad request being sent to UploadHandler. Error:", err)
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
	newFilename := u.String() + header.Filename
	targetFile, err := os.Create("/tmp/uploads/" + newFilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()
	io.Copy(targetFile, f)
	fmt.Fprintln(w, newFilename)
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
