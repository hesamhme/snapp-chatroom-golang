package main

import (
	"appchat/internal/infrastructure/logger"
	"appchat/internal/infrastructure/nats"
	"appchat/internal/infrastructure/redis"
	"appchat/internal/infrastructure/tcp"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.Init()
	logrus.Info("Chat server starting...")

	natsClient, err := nats.NewNATSConnection("nats://localhost:4222")
	if err != nil {
		logrus.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	redisClient := redis.NewRedisClient()

	logrus.Info("Starting TCP handler on port 8080...")
	tcpHandler := tcp.NewTCPHandler(natsClient, redisClient)
	tcpHandler.Start("8080")
}
