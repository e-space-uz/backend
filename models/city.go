package models

type City struct {
	ID     string `json:"id" bson:"_id"`
	RuName string `json:"ru_name" bson:"ru_name"`
	Name   string `json:"name" bson:"name"`
	Soato  uint32 `json:"soato" bson:"soato"`
	Code   uint32 `json:"code" bson:"code"`
}
