package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DraftEntity struct {
	ID             string    `json:"id" bson:"_id"`
	EntityNumber   string    `json:"entity_number" bson:"entity_number"`
	KadastrNumber  string    `json:"kadastr_number" bson:"kadastr_number"`
	Address        string    `json:"address" bson:"address"`
	EntityTypeCode uint64    `json:"entity_type_code" bson:"entity_type_code"`
	City           *City     `json:"city" bson:"city"`
	Region         *Region   `json:"region" bson:"region"`
	District       *District `json:"district" bson:"district"`
	Status         string    `json:"status" bson:"status"`
	StaffIds       []string  `json:"staff_ids" bson:"staff_ids"`
	EntityGallery  []string  `json:"entity_gallery" bson:"entity_gallery"`
}
type EntityDraft struct {
	ID                string               `json:"id" bson:"_id"`
	EntityDraftSoato  string               `json:"entity_draft_soato" bson:"entity_draft_soato"`
	Comment           string               `json:"comment" bson:"comment"`
	EntityDraftNumber string               `json:"entity_draft_number" bson:"entity_draft_number"`
	City              *City                `json:"city" bson:"city"`
	Region            *Region              `json:"region" bson:"region"`
	District          *District            `json:"district" bson:"district"`
	Status            string               `json:"status" bson:"status"`
	Entity            *DraftEntity         `json:"entity" bson:"entity"`
	EntityGallery     []string             `json:"entity_gallery" bson:"entity_gallery"`
	EntityProperty    []*GetEntityProperty `json:"entity_properties" bson:"entity_properties"`
	CreatedAt         primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt         primitive.DateTime   `json:"updated_at" bson:"updated_at"`
}
type GetAllEntityDraft struct {
	ID                string             `json:"id" bson:"_id"`
	EntityDraftNumber string             `json:"entity_draft_number" bson:"entity_draft_number"`
	EntityDraftSoato  string             `json:"entity_draft_soato" bson:"entity_draft_soato"`
	Comment           string             `json:"comment" bson:"comment"`
	Entity            *DraftEntity       `json:"entity" bson:"entity"`
	City              *City              `json:"city" bson:"city"`
	Region            *Region            `json:"region" bson:"region"`
	District          *District          `json:"district" bson:"district"`
	Status            string             `json:"status" bson:"status"`
	EntityProperty    []*EntityProperty  `json:"entity_properties" bson:"entity_properties"`
	CreatedAt         primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt         primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
type UpdateEntityDraftStatus struct {
	EntityDraftID string `json:"entity_draft_id"`
	StatusID      string `json:"status_id"`
}

type CreateEntityDraft struct {
	ID                primitive.ObjectID      `bson:"_id"`
	EntityID          primitive.ObjectID      `bson:"entity_id"`
	EntityDraftNumber string                  `bson:"entity_draft_number" example:"T123"`
	EntityDraftSoato  string                  `bson:"entity_draft_soato"`
	Comment           string                  `bson:"comment"`
	City              City                    `bson:"city"`
	Region            Region                  `bson:"region"`
	District          District                `bson:"district"`
	StatusId          primitive.ObjectID      `bson:"status_id"`
	EntityGallery     []string                `bson:"entity_gallery"`
	EntityProperties  []*CreateEntityProperty `bson:"entity_properties"`
	CreatedAt         time.Time               `bson:"created_at"`
	UpdatedAt         time.Time               `bson:"updated_at"`
	DeletedAt         time.Time               `bson:"deleted_at"`
}

// swagger requests
type EntityDraftSwag struct {
	City             City              `json:"city"`
	Region           Region            `json:"region"`
	District         District          `json:"district"`
	Comment          string            `json:"comment" bson:"comment"`
	EntityGallery    []string          `json:"entity_gallery"`
	EntityProperties []*EntityProperty `json:"entity_properties"`
	EntityID         string            `json:"entity_id"`
}

type ConfirmEntityDraftSwag struct {
	StatusID string `json:"status_id" binding:"required"`
	EntityID string `json:"entity_id" binding:"required"`
	Comment  string `json:"comment" binding:"required"`
}

type GetAllEntityDraftsRequestsSwag struct {
	Date              string `json:"date"`
	CityID            string `json:"city_id"`
	RegionID          string `json:"region_id"`
	StatusID          string `json:"status_id"`
	FromDate          string `json:"from_date" example:"2021-11-21"`
	ToDate            string `json:"to_date" example:"2021-11-21"`
	EntityDraftNumber string `json:"entity_draft_number"`
	Page              uint32 `json:"page"`
	Limit             uint32 `json:"limit"`
}
type UpdateEntityDraftPropertySwag struct {
	EntityProperty []EntityProperty `json:"entity_properties" binding:"required"`
}
