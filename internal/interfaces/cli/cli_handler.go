package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"appchat/internal/domain"
)



func StartCLI(serverAddress string) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Failed to connect to chat server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Chatroom!")
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter chatroom name: ")
	chatroom, _ := reader.ReadString('\n')
	chatroom = strings.TrimSpace(chatroom)

	// Join the chatroom
	joinMsg := domain.Message{Username: username, Chatroom: chatroom, Content: "has joined the chatroom"}
	if err := sendMessage(conn, joinMsg); err != nil {
		fmt.Println("Error sending join message:", err)
		return
	}

	// Start listening for messages
	go listenForMessages(conn)

	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			break
		}

		msg := domain.Message{Username: username, Chatroom: chatroom, Content: message}
		if err := sendMessage(conn, msg); err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}

func sendMessage(conn net.Conn, msg domain.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(append(data, '\n'))
	return err
}

func listenForMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var msg domain.Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err == nil {
			fmt.Printf("[%s] %s: %s\n", msg.Chatroom, msg.Username, msg.Content)
		}
	}
}
