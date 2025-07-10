package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var dataFile = "db.json"

func loadItems() ([]Item, error) {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	var items []Item
	err = json.Unmarshal(file, &items)
	return items, err
}

func saveItems(items []Item) error {
	bytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, bytes, 0644)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	items, _ := loadItems()
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	items, _ := loadItems()
	for _, item := range items {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.NotFound(w, r)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	_ = json.NewDecoder(r.Body).Decode(&newItem)

	items, _ := loadItems()
	items = append(items, newItem)
	saveItems(items)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updated Item
	_ = json.NewDecoder(r.Body).Decode(&updated)

	items, _ := loadItems()
	for i, item := range items {
		if item.ID == params["id"] {
			items[i] = updated
			saveItems(items)
			json.NewEncoder(w).Encode(updated)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	items, _ := loadItems()
	for i, item := range items {
		if item.ID == params["id"] {
			items = append(items[:i], items[i+1:]...)
			saveItems(items)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/items", createItem).Methods("POST")
	router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
