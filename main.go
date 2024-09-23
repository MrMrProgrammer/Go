package main

import (
	"encoding/json"
	"net/http"
	"log"
	"strconv"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	books   = []Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan"},
		{ID: 2, Title: "Clean Code", Author: "Robert C. Martin"},
		{ID: 3, Title: "The Pragmatic Programmer", Author: "Andy Hunt"},
	}
	mu      sync.Mutex
	nextID  = 4
)

func bookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w, r)
	case http.MethodPost:
		createBook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	book.ID = nextID
	nextID++
	books = append(books, book)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/books/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.ID == id {
			if err := json.NewDecoder(r.Body).Decode(&books[i]); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			books[i].ID = id // ensure ID remains the same
			json.NewEncoder(w).Encode(books[i])
			return
		}
	}

	http.NotFound(w, r)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/books/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/books", bookHandler)
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			updateBook(w, r)
		case http.MethodDelete:
			deleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
