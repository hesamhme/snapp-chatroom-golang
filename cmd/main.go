package main

// import (
// 	"appchat/internal/infrastructure/nats"
// 	"appchat/internal/infrastructure/redis"
// 	"appchat/internal/infrastructure/tcp"
// 	"appchat/internal/interfaces/cli"
// 	"os"

// 	"github.com/sirupsen/logrus"
// )

// func initLogger() {
// 	logFile, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		logrus.Fatalf("Failed to open log file: %v", err)
// 	}

// 	logrus.SetOutput(logFile)
// 	logrus.SetFormatter(&logrus.JSONFormatter{})
// 	logrus.SetLevel(logrus.InfoLevel)
// }

// func main() {
// 	initLogger()
// 	logrus.Info("Application starting...")

// 	if len(os.Args) > 1 && os.Args[1] == "client" {
// 		logrus.Info("Starting chat client...")
// 		cli.StartCLI("localhost:8080")
// 		return
// 	}

// 	natsClient, err := nats.NewNATSConnection("nats://localhost:4222")
// 	if err != nil {
// 		logrus.Fatalf("Failed to connect to NATS: %v", err)
// 	}
// 	defer natsClient.Close()

// 	redisClient := redis.NewRedisClient()

// 	logrus.Info("Starting TCP handler on port 8080...")
// 	tcpHandler := tcp.NewTCPHandler(natsClient, redisClient)
// 	tcpHandler.Start("8080")
// }
