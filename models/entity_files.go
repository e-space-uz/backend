package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EntityFiles struct {
	ID       string `json:"id" bson:"_id"`
	FileName string `json:"name" bson:"name"`
	Url      string `json:"url" bson:"url"`
	Comment  string `json:"comment" bson:"comment"`
	User     string `json:"user" bson:"user"`
}

type CreateEntityFiles struct {
	ID        primitive.ObjectID `bson:"_id"`
	FileName  string             `bson:"name"`
	Url       string             `bson:"url"`
	Comment   string             `bson:"comment"`
	User      string             `bson:"user"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type EntityFilesSwag struct {
	Comment string `json:"comment"`
}
