package routes

import (
	"go-postgres-app/handlers"
	"net/http"

	"github.com/gorilla/mux"
)
import "go-postgres-app/middleware"


func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// r.HandleFunc("/users", handlers.GetUsers).Methods(http.MethodGet)
	// r.HandleFunc("/users/{id}", handlers.GetUserByID).Methods(http.MethodGet)
	// r.HandleFunc("/users", handlers.CreateUser).Methods(http.MethodPost)
	// r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods(http.MethodPut)
	// r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/auth/register", handlers.Register).Methods("POST")
	r.HandleFunc("/auth/login", handlers.Login).Methods("POST")
	r.Handle("/users", middleware.JWTAuth(http.HandlerFunc(handlers.GetUsers))).Methods("GET")
	r.Handle("/users/{id}", middleware.JWTAuth(http.HandlerFunc(handlers.GetUserByID))).Methods("GET")
	r.Handle("/users", middleware.JWTAuth(http.HandlerFunc(handlers.CreateUser))).Methods("POST")
	r.Handle("/users/{id}", middleware.JWTAuth(http.HandlerFunc(handlers.UpdateUser))).Methods("PUT")
	r.Handle("/users/{id}", middleware.JWTAuth(http.HandlerFunc(handlers.DeleteUser))).Methods("DELETE")
	r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")
	return r
}
