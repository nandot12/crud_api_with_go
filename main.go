package main

import (
	"log"
	"net/http"
	"go-mysql-api/db"
	"go-mysql-api/handler"
	"github.com/gorilla/mux"
)

func main() {
	// Koneksi ke database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Router
	r := mux.NewRouter()

	// Endpoint
	r.HandleFunc("/users", handlers.GetUsers(database)).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser(database)).Methods("POST") 
	r.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser(database)).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", handlers.DeleteUser(database)).Methods("DELETE") 

	// Jalankan server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
