package main

import (
	"log"
	"net/http"

	"go-backend/controllers"
	"go-backend/db"

	"github.com/rs/cors"
)

func main() {
	db.Connect("mongodb://localhost:27017") // connect to MongoDB

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/signup", controllers.Signup)
	mux.HandleFunc("/login", controllers.Login)

	// Notes CRUD
	mux.HandleFunc("/note", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetNotes(w, r)
		case http.MethodPost:
			controllers.CreateNote(w, r)
		case http.MethodPut:
			controllers.UpdateNote(w, r)
		case http.MethodDelete:
			controllers.DeleteNote(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	handler := cors.AllowAll().Handler(mux)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
