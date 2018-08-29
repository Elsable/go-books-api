package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asdine/storm"
	uuid "github.com/satori/go.uuid"

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

var db storm.DB

func openDB() error {
	db, err := storm.Open("books.db")
	return err
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	if err := openDB(); err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	var books []book
	if err := db.All(&books); err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	return
}

func addBook(w http.ResponseWriter, r *http.Request) {
	if err := openDB(); err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var book book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = uuid.Must(uuid.NewV4()).String()
	if err := db.Save(&book); err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(book)
	return
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	err := db.DeleteStruct(&user)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}

	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/addbook", addBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

}
