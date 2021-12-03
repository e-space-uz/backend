package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct {
	ID               string            `json:"id" bson:"_id"`
	Name             string            `json:"name" bson:"name"`
	Label            string            `json:"label" bson:"label"`
	Placeholder      string            `json:"placeholder" bson:"placeholder"`
	Type             string            `json:"type" bson:"type"`
	Validation       string            `json:"validation" bson:"validation"`
	Description      string            `json:"description" bson:"description"`
	CollectionName   string            `json:"collection_name" bson:"collection_name"`
	Status           bool              `json:"status" bson:"status"`
	IsRequired       bool              `json:"is_required" bson:"is_required"`
	WithConfirmation bool              `json:"with_confirmation" bson:"with_confirmation"`
	PropertyOptions  []*PropertyOption `json:"property_options" bson:"property_options"`
}
type GetProperty struct {
	ID               string            `json:"id" bson:"_id"`
	Name             string            `json:"name" bson:"name" example:"Doe"`
	Label            string            `json:"label" bson:"label" example:"Doe"`
	Placeholder      string            `json:"placeholder" bson:"placeholder" example:"Doe"`
	Type             string            `json:"type" bson:"type" example:"radio"`
	Validation       string            `json:"validation" bson:"validation"`
	Description      string            `json:"description" bson:"description"`
	CollectionName   string            `json:"collection_name" bson:"collection_name"`
	Status           bool              `json:"status" bson:"status"`
	IsRequired       bool              `json:"is_required" bson:"is_required" example:"false"`
	WithConfirmation bool              `json:"with_confirmation" bson:"with_confirmation"`
	PropertyOptions  []*PropertyOption `json:"property_options" bson:"property_options"`
}
type GetAllPropertiesResponse struct {
	Properties []*Property `json:"properties"`
	Count      uint32      `json:"count"`
}

type CreateProperty struct {
	ID               primitive.ObjectID `bson:"_id"`
	Name             string             `bson:"name"`
	Type             string             `bson:"type"`
	Label            string             `bson:"label"`
	Placeholder      string             `bson:"placeholder"`
	Validation       string             `bson:"validation"`
	Description      string             `bson:"description"`
	CollectionName   string             `json:"collection_name" bson:"collection_name"`
	Status           bool               `bson:"status"`
	IsRequired       bool               `bson:"is_required"`
	WithConfirmation bool               `bson:"with_confirmation"`
	PropertyOptions  []*PropertyOption  `bson:"property_options"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}

type PropertyOption struct {
	Name  string `json:"name" binding:"required" example:"Option"`
	Value string `json:"value" binding:"required"`
}

type PropertySwag struct {
	Name             string            `json:"name" binding:"required"`
	Type             string            `json:"type" binding:"required" example:"radio"`
	Label            string            `json:"label" binding:"required" example:"Doe"`
	Placeholder      string            `json:"placeholder" binding:"required" example:"Doe"`
	Validation       string            `json:"validation" binding:"required" example:"This is required"`
	Description      string            `json:"description" binding:"required" example:"This is description for property"`
	CollectionName   string            `json:"collection_name" bson:"collection_name"`
	IsRequired       bool              `json:"is_required" binding:"required" example:"false"`
	Status           bool              `json:"status" binding:"required" example:"true"`
	WithConfirmation bool              `json:"with_confirmation" binding:"required" example:"true"`
	PropertyOptions  []*PropertyOption `json:"property_options" binding:"required"`
}
