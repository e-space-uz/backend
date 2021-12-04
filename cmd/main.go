package main

import (
	"github.com/e-space-uz/backend/api"
	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "admin_api_gateway")

	server := api.New(&api.RouterOptions{
		Log:         log,
		Cfg:         cfg,
		Services:    grpcClients,
		RedisClient: redisManager,
		// Kafka:       kafka,
	})
	server.Run(cfg.HttpPort)
}
