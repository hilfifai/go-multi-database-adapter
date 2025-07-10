package main

import (
	"go-postgres-app/config"
	"go-postgres-app/db"
	"go-postgres-app/routes"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()
	db.Init()

	r := routes.RegisterRoutes()

	log.Println("🚀 Server started at :8181")
	http.ListenAndServe(":8181", r)
}
