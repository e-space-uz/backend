package config

import (
	"os"

	"github.com/spf13/cast"
)

const (
	EntityCollection        = "EntityCollection"
	EntityDraftCollection   = "EntityDraftCollection"
	EntityFilesCollection   = "EntityFilesCollection"
	PropertyCollection      = "PropertyCollection"
	GroupPropertyCollection = "GroupPropertyCollection"
	UserCollection          = "UserCollection"
	CityCollection             = "CityCollection"
	RegionCollection           = "RegionCollection"
	DistrictCollection         = "DistrictCollection"
	
)

type Config struct {
	Environment string
	LogLevel    string
	HttpPort    string

	MongoHost     string
	MongoPort     int
	MongoUser     string
	MongoPassword string
	MongoDatabase string
}

func Load() Config {

	cfg := Config{}

	cfg.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	cfg.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

	cfg.HttpPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8000"))

	cfg.MongoHost = cast.ToString(getOrReturnDefault("MONGO_HOST", "localhost"))
	cfg.MongoPort = cast.ToInt(getOrReturnDefault("MONGO_PORT", 27017))
	cfg.MongoUser = cast.ToString(getOrReturnDefault("MONGO_USER", "espace"))
	cfg.MongoPassword = cast.ToString(getOrReturnDefault("MONGO_PASSWORD", "mongodb"))
	cfg.MongoDatabase = cast.ToString(getOrReturnDefault("MONGO_DATABASE", "espace"))

	return cfg
}
func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
