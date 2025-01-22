package main

import (
	"appchat/internal/application"
	"appchat/internal/infrastructure/logger"
	"appchat/internal/infrastructure/nats"
	"appchat/internal/infrastructure/redis"
	"appchat/internal/infrastructure/tcp"

	"github.com/sirupsen/logrus"
)

func main() {
	logger.Init()
	logrus.Info("Chat server starting...")

	natsClient, _ := nats.NewNATSConnection("nats://localhost:4222")
	redisClient := redis.NewRedisClient()

	chatroomUseCase := application.NewChatroomUseCase(natsClient, redisClient)

	tcpHandler := tcp.NewTCPHandler(chatroomUseCase)
	tcpHandler.Start("8080")
}
