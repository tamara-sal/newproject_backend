package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Note struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
    Title       string             `bson:"title" json:"title"`
    Description string             `bson:"description" json:"description"` 
    AuthorEmail string             `bson:"authorEmail" json:"authorEmail"`
}
