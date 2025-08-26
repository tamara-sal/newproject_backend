package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-backend/db"
	"go-backend/models"

	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
        return
    }
	user.Email = strings.ToLower(user.Email) // convert email to lowercase


    hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashed)

    collection := db.Client.Database("uni_notes").Collection("users")
    _, err := collection.InsertOne(r.Context(), user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

func Login(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
	credentials.Email = strings.ToLower(credentials.Email)

    json.NewDecoder(r.Body).Decode(&credentials)

    collection := db.Client.Database("uni_notes").Collection("users")
    var user models.User
    err := collection.FindOne(r.Context(), map[string]string{"email": credentials.Email}).Decode(&user)

    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
        return
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"user": map[string]string{"name": user.Name, "email": user.Email}})
}
