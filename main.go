package main

import (
    "log"
    "net/http"

    "go-backend/controllers"
    "go-backend/db"

    "github.com/rs/cors"
)

func main() {
    // Connect to MongoDB
    db.Connect("mongodb://localhost:27017")

    mux := http.NewServeMux()
    mux.HandleFunc("/signup", controllers.Signup)
    mux.HandleFunc("/login", controllers.Login)

    handler := cors.AllowAll().Handler(mux) // allow all origins
    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
