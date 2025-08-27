package controllers

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "go-backend/db"
    "go-backend/models"

    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // change this to a secure key

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

    user.Email = strings.ToLower(user.Email)
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

    json.NewDecoder(r.Body).Decode(&credentials)
    credentials.Email = strings.ToLower(credentials.Email)

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

    // Create JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": user.Email,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    })
    tokenString, _ := token.SignedString(jwtSecret)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "user":  map[string]string{"name": user.Name, "email": user.Email},
        "token": tokenString,
    })
}
