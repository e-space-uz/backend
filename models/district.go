package models

type District struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	RuName     string `json:"ru_name" bson:"ru_name"`
	Code       uint32 `json:"code" bson:"code"`
	ExternalID uint32 `json:"external_id" bson:"external_id"`
	Soato      uint32 `json:"soato" bson:"soato"`
	City       City   `json:"city" bson:"city"`
	Region     Region `json:"region" bson:"region"`
}
