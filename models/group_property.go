package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupProperty struct {
	ID          string      `json:"id" bson:"_id"`
	Name        string      `json:"name" bson:"name" example:"Doe"`
	Step        uint32      `json:"step" bson:"step"`
	Type        uint32      `json:"type" bson:"type"`
	Description string      `json:"description" bson:"description"`
	Status      bool        `json:"status" bson:"status"`
	Properties  []*Property `json:"properties" bson:"properties"`
}
type GetAllGroupProperty struct {
	ID            string           `json:"id" bson:"_id"`
	Name          string           `json:"name" bson:"name" example:"Doe"`
	Step          uint32           `json:"step" bson:"step"`
	Type          uint32           `json:"type" bson:"type"`
	Status        bool             `json:"status" bson:"status"`
	Description   string           `json:"description" bson:"description"`
	ReadStatuses  []string         `json:"read_statuses" bson:"read_statuses,omitempty"`
	WriteStatuses []string         `json:"write_statuses" bson:"write_statuses,omitempty"`
	Properties    []*GetProperties `json:"properties" bson:"properties"`
}

type Properties struct {
	Property GetProperty `json:"property" bson:"property"`
	Order    uint32      `json:"order" bson:"order"`
}

type GetProperties struct {
	PropertyID string `json:"property_id" bson:"property_id"`
	Order      uint32 `json:"order" bson:"order"`
}

type GetAllGroupPropertiesResponse struct {
	GroupProperties []*GetAllGroupProperty `json:"group_properties"`
	Count           uint32                 `json:"count"`
}
type GetGroupPropertyByStatusIDResponse struct {
	GroupProperties []*GetGroupPropertyByStatusID `json:"group_properties"`
}

type GetAllGroupPropertiesByTypeResponse struct {
	GroupProperties []*GroupProperty `json:"group_properties" bson:"group_properties"`
	Count           uint32           `json:"count" bson:"count"`
}
type GetGroupPropertyByStatusID struct {
	Id           string          `json:"id" bson:"_id"`
	Name         string          `json:"name" bson:"name"`
	Description  string          `json:"description" bson:"description"`
	Step         uint32          `json:"step" bson:"step"`
	Type         uint32          `json:"type" bson:"type"`
	Status       bool            `json:"status" bson:"status"`
	IsDisable    bool            `json:"is_disable" bson:"is_disable"`
	Properties   []*Property     `json:"properties" bson:"properties"`
	Organization OrganizationGet `json:"organization" bson:"organization"`
}
type CreateGroupProperty struct {
	ID            primitive.ObjectID   `bson:"_id"`
	Step          uint32               `bson:"step"`
	Name          string               `bson:"name"`
	Type          uint32               `bson:"type"`
	Status        bool                 `bson:"status"`
	Description   string               `bson:"description"`
	Organization  OrganizationCreate   `json:"organization" bson:"organization"`
	CreatedAt     time.Time            `bson:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at"`
	Properties    []*CreateProperties  `bson:"properties"`
	ReadStatuses  []primitive.ObjectID `bson:"read_statuses"`
	WriteStatuses []primitive.ObjectID `bson:"write_statuses"`
}

type CreateProperties struct {
	PropertyID primitive.ObjectID `json:"property_id" bson:"property_id"`
	Order      uint32             `json:"order" bson:"order"`
}

type GroupPropertySwag struct {
	Name          string           `json:"name" binding:"required"`
	Step          uint32           `json:"step" binding:"required"`
	Type          uint32           `json:"type" binding:"required"`
	Status        bool             `json:"status" binding:"required"`
	Organization  OrganizationGet  `json:"organization" bson:"organization"`
	Description   string           `json:"description" binding:"required" example:"I have no idea"`
	Properties    []*GetProperties `json:"properties" binding:"required"`
	ReadStatuses  []string         `json:"read_statuses" binding:"required" example:"60dd9c0a4472a2aaa970304e"`
	WriteStatuses []string         `json:"write_statuses" binding:"required" example:"60dd9c1a729317449b1ada03"`
}

type OrganizationGet struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type OrganizationCreate struct {
	Name string             `json:"name"`
	ID   primitive.ObjectID `json:"id"`
}
