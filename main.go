package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	author *author `json:"author"`
}

type author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/books", getBooks).Methods("GET")

}
