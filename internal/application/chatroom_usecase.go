package application

import (
	"appchat/internal/domain"
	"appchat/internal/infrastructure/nats"
	"appchat/internal/infrastructure/redis"
	"encoding/json"
	"fmt"
)

type ChatroomUseCase struct {
	natsClient  *nats.NATSClient
	redisClient *redis.RedisClient
}

// Factory function to create a new use case
func NewChatroomUseCase(natsClient *nats.NATSClient, redisClient *redis.RedisClient) *ChatroomUseCase {
	return &ChatroomUseCase{
		natsClient:  natsClient,
		redisClient: redisClient,
	}
}

// Handles user joining the chatroom
func (c *ChatroomUseCase) JoinChatroom(username, chatroom string) error {
	if err := c.redisClient.AddUserToChatroom(chatroom, username); err != nil {
		return err
	}
	if err := c.redisClient.AddChatroom(chatroom); err != nil {
		return err
	}

	message := domain.Message{
		Username: username,
		Chatroom: chatroom,
		Content:  fmt.Sprintf("%s has joined the chatroom", username),
	}

	return c.natsClient.PublishMessage("chatroom."+chatroom, string(encodeMessage(message)))
}

// Handles user leaving the chatroom
func (c *ChatroomUseCase) LeaveChatroom(username, chatroom string) error {
	if err := c.redisClient.RemoveUserFromChatroom(chatroom, username); err != nil {
		return err
	}

	users, _ := c.redisClient.GetUsersInChatroom(chatroom)
	if len(users) == 0 {
		c.redisClient.RemoveChatroom(chatroom)
	}

	message := domain.Message{
		Username: "server",
		Chatroom: chatroom,
		Content:  fmt.Sprintf("%s has left the chatroom", username),
	}

	return c.natsClient.PublishMessage("chatroom."+chatroom, string(encodeMessage(message)))
}

// Retrieves active users in the chatroom
func (c *ChatroomUseCase) GetUsers(chatroom string) ([]string, error) {
	return c.redisClient.GetUsersInChatroom(chatroom)
}

// Retrieves all active chatrooms
func (c *ChatroomUseCase) GetChatrooms() ([]string, error) {
	return c.redisClient.GetChatrooms()
}

// Sends a message to the chatroom
func (c *ChatroomUseCase) SendMessage(username, chatroom, content string) error {
	message := domain.Message{
		Username: username,
		Chatroom: chatroom,
		Content:  content,
	}
	return c.natsClient.PublishMessage("chatroom."+chatroom, string(encodeMessage(message)))
}

func encodeMessage(msg domain.Message) []byte {
	data, _ := json.Marshal(msg)
	return data
}
