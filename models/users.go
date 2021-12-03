package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          string             `json:"id" bson:"_id"`
	UserType    string             `json:"user_type" bson:"user_type"`
	Login       string             `json:"login" bson:"login"`
	Password    string             `json:"password" bson:"password"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
