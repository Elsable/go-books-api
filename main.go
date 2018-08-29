package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asdine/storm"

	"github.com/gorilla/mux"
)

type book struct {
	ID     string  `storm:"id"`
	Isbn   string  `storm:"unique"`
	Title  string  `storm:"index"`
	Author *author `storm:"inline"`
}

type author struct {
	Firstname string `storm:"index"`
	Lastname  string `storm:"index"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	db, err := storm.Open("books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var books []book
	if err := db.All(&books); err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/books", getBooks).Methods("GET")

}
