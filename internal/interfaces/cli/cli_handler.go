package cli

import (
	"appchat/internal/application"
	"bufio"
	"fmt"
	"os"
	"strings"
)
// this is temp version
func StartCLI(chatroomUseCase *application.ChatroomUseCase) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Chatroom!")
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter chatroom name: ")
	chatroomName, _ := reader.ReadString('\n')
	chatroomName = strings.TrimSpace(chatroomName)

	chatroomUseCase.JoinChatroom(username, chatroomName)
	go chatroomUseCase.SubscribeToMessages(chatroomName)

	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			break
		}

		chatroomUseCase.SendMessage(chatroomName, username, message)
	}
}
