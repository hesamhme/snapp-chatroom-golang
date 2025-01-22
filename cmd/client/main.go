package main

import (
	"appchat/internal/infrastructure/logger"
	"appchat/internal/interfaces/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.Init()
	logrus.Info("Chat client starting...")

	serverAddress := "localhost:8080"
	cli.StartCLI(serverAddress)
}
