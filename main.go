package main

import (
	"log"
	"net/http"

	"mustika/config"
	"mustika/routes"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDatabase()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
