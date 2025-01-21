package cli

import (
	"appchat/internal/domain"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type UserInput struct {
	Username string
	Chatroom string
}

func StartCLI(serverAddress string) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Failed to connect to chat server:", err)
		return
	}
	defer conn.Close()

	userInput := collectUserInput()

	joinMsg := domain.Message{
		Username: userInput.Username,
		Chatroom: userInput.Chatroom,
		Content:  "has joined the chatroom",
	}
	if err := sendMessage(conn, joinMsg); err != nil {
		fmt.Println("Error sending join message:", err)
		return
	}

	go listenForMessages(conn)

	sendUserMessages(conn, userInput)
}

func collectUserInput() UserInput {
	var userInput UserInput
	fmt.Print("Enter your username: ")
	fmt.Scanln(&userInput.Username)
	fmt.Print("Enter chatroom name: ")
	fmt.Scanln(&userInput.Chatroom)
	return userInput
}

func sendUserMessages(conn net.Conn, userInput UserInput) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		if scanner.Scan() {
			message := scanner.Text()
			if message == "exit" {
				break
			}
			msg := domain.Message{
				Username: userInput.Username,
				Chatroom: userInput.Chatroom,
				Content:  message,
			}
			if err := sendMessage(conn, msg); err != nil {
				fmt.Println("Error sending message:", err)
				break
			}
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
		msg := scanner.Text()
		if strings.HasPrefix(msg, "Users in chatroom") {
			fmt.Println(msg) 
		} else {
			var chatMsg domain.Message
			if err := json.Unmarshal([]byte(msg), &chatMsg); err == nil {
				fmt.Printf("[%s] %s: %s\n", chatMsg.Chatroom, chatMsg.Username, chatMsg.Content)
			}
		}
	}
}
