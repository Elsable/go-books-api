package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/asdine/storm"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string  `storm:"id"`
	Isbn   string  `storm:"unique"`
	Title  string  `storm:"index"`
	Author *Author `storm:"inline"`
}

type Author struct {
	Firstname string `storm:"index"`
	Lastname  string `storm:"index"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	db, err := storm.Open("books.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	var books []Book
	if err := db.All(&books); err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	return
}

func addBook(w http.ResponseWriter, r *http.Request) {
	db, err := storm.Open("books.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	fmt.Println(book)
	book.ID = uuid.Must(uuid.NewV4()).String()
	if err := db.Save(&book); err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(book)
	return
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	db, err := storm.Open("books.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book Book
	if err := db.One("ID", params["id"], &book); err != nil {
		e := map[string]string{"error": "book not found"}
		json.NewEncoder(w).Encode(e)
		return
	}
	if err := db.DeleteStruct(&book); err != nil {
		e := map[string]string{"error": "book does not want to be removed"}
		json.NewEncoder(w).Encode(e)
		return
	}
	var books []Book
	if err := db.All(&books); err != nil {
		e := map[string]string{"error": "can't fetch remaining books at this time"}
		json.NewEncoder(w).Encode(e)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	db, err := storm.Open("books.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book []Book
	if err := db.One("ID", params["id"], &book); err != nil {
		e := map[string]string{"error": "book not found"}
		json.NewEncoder(w).Encode(e)
		return
	}
	json.NewEncoder(w).Encode(&Book{})
}

func getBook(w http.ResponseWriter, r *http.Request) {
	db, err := storm.Open("books.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	var book []Book
	if err := db.One("ID", params["id"], &book); err != nil {
		e := map[string]string{"error": "book not found"}
		json.NewEncoder(w).Encode(e)
		return
	}
	json.NewEncoder(w).Encode(&Book{})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books", addBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
