package application

import (
	"appchat/internal/infrastructure/nats"
	"fmt"
)

type ChatroomUseCase struct {
	natsClient *nats.NATSClient
}

func NewChatroomUseCase(natsClient *nats.NATSClient) *ChatroomUseCase {
	return &ChatroomUseCase{natsClient: natsClient}
}

func (uc *ChatroomUseCase) JoinChatroom(username, chatroomName string) {
	joinMessage := fmt.Sprintf("%s has joined the chatroom: %s", username, chatroomName)
	uc.natsClient.PublishMessage(chatroomName, joinMessage)
}

func (uc *ChatroomUseCase) SendMessage(chatroomName, username, message string) {
	fullMessage := fmt.Sprintf("[%s] %s: %s", chatroomName, username, message)
	uc.natsClient.PublishMessage(chatroomName, fullMessage)
}

func (uc *ChatroomUseCase) SubscribeToMessages(chatroomName string) {
	uc.natsClient.Subscribe(chatroomName, func(msg string) {
		fmt.Println(msg)
	})
}
