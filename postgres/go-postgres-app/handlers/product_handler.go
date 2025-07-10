package handlers

import (
	"encoding/json"
	"go-postgres-app/db"
	"go-postgres-app/models"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var productValidate = validator.New()

func GetProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")

	query := "SELECT id, name, description, price FROM products WHERE 1=1"
	var args []interface{}
	argPos := 1

	if name != "" {
		query += " AND name ILIKE $" + strconv.Itoa(argPos)
		args = append(args, "%"+name+"%")
		argPos++
	}
	if minPrice != "" {
		query += " AND price >= $" + strconv.Itoa(argPos)
		args = append(args, minPrice)
		argPos++
	}
	if maxPrice != "" {
		query += " AND price <= $" + strconv.Itoa(argPos)
		args = append(args, maxPrice)
		argPos++
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
		if err == nil {
			products = append(products, p)
		}
	}

	json.NewEncoder(w).Encode(products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	row := db.DB.QueryRow("SELECT id, name, description, price FROM products WHERE id = $1", id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := productValidate.Struct(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("INSERT INTO products (name, description, price) VALUES ($1, $2, $3)",
		p.Name, p.Description, p.Price)
	if err != nil {
		http.Error(w, "Insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product created"})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := productValidate.Struct(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4",
		p.Name, p.Description, p.Price, id)
	if err != nil {
		http.Error(w, "Update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
}
