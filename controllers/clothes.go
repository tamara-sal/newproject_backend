package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"go-backend/db"
	"go-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getClothesCollection() *mongo.Collection {
	return db.Client.Database("uni_notes").Collection("clothes")
}

// Add a cloth
func AddCloth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cloth models.Cloth
	json.NewDecoder(r.Body).Decode(&cloth)

	if cloth.Category != "men" && cloth.Category != "women" {
		http.Error(w, "Category must be 'men' or 'women'", http.StatusBadRequest)
		return
	}

	res, _ := getClothesCollection().InsertOne(context.Background(), cloth)
	cloth.ID = res.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(cloth)
}

// Get clothes by category
func GetClothesByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	category := r.URL.Query().Get("category")

	filter := bson.M{"category": category}
	cursor, _ := getClothesCollection().Find(context.Background(), filter)
	defer cursor.Close(context.Background())

	var clothes []models.Cloth
	cursor.All(context.Background(), &clothes)
	json.NewEncoder(w).Encode(clothes)
}
