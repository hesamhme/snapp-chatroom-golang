package main

import (
	"appchat/internal/infrastructure/nats"
	"appchat/internal/infrastructure/tcp"
	"appchat/internal/interfaces/cli"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "client" {
		log.Println("Starting chat client...")
		cli.StartCLI("localhost:8080")
		return
	}

	// Initialize NATS connection
	natsClient, err := nats.NewNATSConnection("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	// Start the TCP handler to forward messages to NATS
	log.Println("Starting TCP handler on port 8080...")
	tcpHandler := tcp.NewTCPHandler(natsClient)
	tcpHandler.Start("8080")
}
