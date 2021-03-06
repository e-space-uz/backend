package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entity struct {
	ID             string          `json:"id" bson:"_id"`
	EntitySoato    string          `json:"entity_soato" bson:"entity_soato"`
	Address        string          `json:"address" bson:"address"`
	RevertComment  string          `json:"revert_comment" bson:"revert_comment"`
	EntityNumber   string          `json:"entity_number" bson:"entity_number"`
	EntityTypeCode uint64          `json:"entity_type_code" bson:"entity_type_code"`
	Version        uint64          `json:"version" bson:"version"`
	Organizations  map[string]bool `json:"organizations" bson:"organizations"`
	Status         string          `json:"status" bson:"status"`
	City           *City           `json:"city" bson:"city"`
	Region         *Region         `json:"region" bson:"region"`
	District       *District       `json:"district" bson:"district"`
	StaffIds       []string        `json:"staff_ids" bson:"staff_ids"`
	EntityGallery  []string        `json:"entity_gallery" bson:"entity_gallery"`
	// EntityDrafts   []*GetAllEntityDraft `json:"entity_drafts" bson:"entity_drafts"`
	EntityFiles    []*EntityFiles       `json:"entity_files" bson:"entity_files"`
	EntityProperty []*GetEntityProperty `json:"entity_properties" bson:"entity_properties"`
	CreatedAt      primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt      primitive.DateTime   `json:"updated_at" bson:"updated_at"`
}

type GetAllEntities struct {
	ID               string             `json:"id" bson:"_id"`
	Address          string             `json:"address" bson:"address"`
	EntitySoato      string             `json:"entity_soato" bson:"entity_soato"`
	EntityNumber     string             `json:"entity_number" bson:"entity_number"`
	EntityTypeCode   uint64             `json:"entity_type_code" bson:"entity_type_code"`
	Version          uint64             `json:"version" bson:"version"`
	Status           string             `json:"status" bson:"status"`
	EntityProperties []*EntityProperty  `json:"entity_properties" bson:"entity_properties"`
	City             *City              `json:"city" bson:"city"`
	Region           *Region            `json:"region" bson:"region"`
	District         *District          `json:"district" bson:"district"`
	EntityFiles      []string           `json:"entity_files" bson:"entity_files"`
	EntityGallery    []string           `json:"entity_gallery" bson:"entity_gallery"`
	CreatedAt        primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt        primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type CreateUpdateEntity struct {
	ID                 primitive.ObjectID      `bson:"_id"`
	Address            string                  `bson:"address"`
	EntitySoato        string                  `bson:"entity_soato"`
	EntityNumber       string                  `bson:"entity_number"`
	EntityTypeCode     uint64                  `bson:"entity_type_code"`
	Version            uint64                  `bson:"version"`
	City               *City                   `bson:"city"`
	Region             *Region                 `bson:"region"`
	District           *District               `bson:"district"`
	Status             string                  `bson:"status"`
	StaffIds           []primitive.ObjectID    `bson:"staff_ids"`
	EntityFiles        []primitive.ObjectID    `bson:"entity_files"`
	EntityDrafts       []primitive.ObjectID    `bson:"entity_drafts"`
	EntityGallery      []string                `bson:"entity_gallery"`
	EntityProperties   []*CreateEntityProperty `bson:"entity_properties"`
	CreatedAt          time.Time               `bson:"created_at"`
	UpdatedAt          time.Time               `bson:"updated_at"`
	EntityStatusUpdate time.Time               `bson:"entity_status_update"`
	DeletedAt          time.Time               `bson:"deleted_at"`
}

type EntityProperty struct {
	PropertyID string `json:"property_id" bson:"property_id"`
	Value      string `json:"value" bson:"value"`
}
type GetEntityProperty struct {
	Property *GetProperty `json:"property" bson:"property"`
	Value    string       `json:"value" bson:"value"`
}
type CreateEntityProperty struct {
	PropertyID primitive.ObjectID `json:"property_id" bson:"property_id"`
	Value      string             `json:"value" bson:"value"`
}
type UpdateEntityStatus struct {
	EntityID string `json:"entity_id"`
	StatusID string `json:"status"`
}

type CreateUpdateEntitySwag struct {
	EntityTypeCode   uint64            `json:"entity_type_code" binding:"required"  example:"1"`
	Address          string            `json:"address" binding:"required"`
	City             *City             `json:"city" binding:"required"`
	Region           *Region           `json:"region" binding:"required"`
	District         *District         `json:"district" binding:"required"`
	EntityFiles      []string          `json:"entity_files" binding:"required"`
	EntityGallery    []string          `json:"entity_gallery" binding:"required"`
	EntityProperties []*EntityProperty `json:"entity_properties" binding:"required"`
}
type UpdateWithActionIDEntitySwag struct {
	ActionDescription string            `json:"action_description" binding:"required"`
	StatusID          string            `json:"status" binding:"required"`
	ActionType        string            `json:"action_type" binding:"required"`
	Deadline          int               `json:"deadline" binding:"required"`
	EntityGallery     []string          `json:"entity_gallery" binding:"required"`
	EntityFiles       []string          `json:"entity_files" binding:"required"`
	EntityProperties  []*EntityProperty `json:"entity_properties" binding:"required"`
	Organizations     map[string]bool   `json:"organizations"`
}

type UpdateEntityStatusSwag struct {
	ActionID string `json:"action_id" binding:"required"`
	EntityID string `json:"entity_id" binding:"required"`
}
type UpdateEntityPropertySwag struct {
	EntityFile     []string         `json:"entity_file" binding:"required"`
	EntityProperty []EntityProperty `json:"entity_properties" binding:"required"`
	Organizations  map[string]bool  `json:"organizations"`
	ActionID       string           `json:"action_id" binding:"required"`
}

type GetAllEntitiesRequest struct {
	EntitySoato  string `json:"entity_soato"`
	CityID       string `json:"city_id"`
	RegionID     string `json:"region_id"`
	EntityNumber string `json:"entity_number"`
	Page         uint32 `json:"page"`
	Limit        uint32 `json:"limit"`
}
