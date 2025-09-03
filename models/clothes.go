package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cloth struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	ImageURL  string             `bson:"imageUrl" json:"imageUrl"`
	Category  string             `bson:"category" json:"category"` // men or women
	
}