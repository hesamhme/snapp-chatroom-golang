package main

import (
	"appchat/internal/application"
	"appchat/internal/infrastructure/nats"
	"appchat/internal/interfaces/cli"
	"log"
)

func main() {
	// Initialize NATS connection
	natsConn, err := nats.NewNATSConnection("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()

	// Initialize chatroom use case
	chatroomUseCase := application.NewChatroomUseCase(natsConn)

	// Start CLI
	cli.StartCLI(chatroomUseCase)
}
