package main

import (
	"context"
	"fmt"

	"github.com/e-space-uz/backend/api"
	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/pkg/logger"
	"github.com/e-space-uz/backend/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "backend")

	credential := options.Credential{
		Username: cfg.MongoUser,
		Password: cfg.MongoPassword,
	}
	mongoString := fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort)

	mongoConn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoString).SetAuth(credential))
	if err != nil {
		log.Error("error to connect to mongo database", logger.Error(err))
	}
	defer mongoConn.Disconnect(context.Background())
	if err := mongoConn.Ping(context.Background(), nil); err != nil {
		log.Error("Cannot connect to database error ->", logger.Error(err))
		panic(err)
	}
	connDB := mongoConn.Database(cfg.MongoDatabase)
	log.Info("Connected to MongoDB", logger.Any("database: ", connDB.Name()))

	strg := storage.NewStorageMongo(connDB)

	server := api.New(&api.RouterOptions{
		Log:     log,
		Cfg:     cfg,
		Storage: strg,
	})
	server.Run(cfg.HttpPort)
}
