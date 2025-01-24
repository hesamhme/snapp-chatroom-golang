package tests

import (
	"appchat/internal/application"
	"appchat/internal/domain"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock dependencies for Redis and NATS
type MockNATSClient struct {
	PublishedMessages []string
}

func (m *MockNATSClient) PublishMessage(subject, message string) error {
	m.PublishedMessages = append(m.PublishedMessages, message)
	return nil
}

func (m *MockNATSClient) Subscribe(subject string, handler func(string)) {}

func (m *MockNATSClient) Close() {}

type MockRedisClient struct {
	users map[string][]string
}

func (m *MockRedisClient) AddUserToChatroom(chatroom, username string) error {
	m.users[chatroom] = append(m.users[chatroom], username)
	return nil
}

func (m *MockRedisClient) GetUsersInChatroom(chatroom string) ([]string, error) {
	return m.users[chatroom], nil
}

func (m *MockRedisClient) RemoveUserFromChatroom(chatroom, username string) error {
	return nil
}

func (m *MockRedisClient) AddChatroom(chatroom string) error {
	if m.users == nil {
		m.users = make(map[string][]string)
	}
	m.users[chatroom] = []string{}
	return nil
}

func (m *MockRedisClient) GetChatrooms() ([]string, error) {
	var chatrooms []string
	for chatroom := range m.users {
		chatrooms = append(chatrooms, chatroom)
	}
	return chatrooms, nil
}

func (m *MockRedisClient) RemoveChatroom(chatroom string) error {
	delete(m.users, chatroom)
	return nil
}


// Test function
func TestJoinChatroom(t *testing.T) {
	// Arrange: Set up mock dependencies
	mockRedis := &MockRedisClient{users: make(map[string][]string)}
	mockNats := &MockNATSClient{}

	chatroomUseCase := application.NewChatroomUseCase(mockNats, mockRedis)

	// Test data
	username := "test"
	chatroom := "x"

	// Act: Call JoinChatroom method
	err := chatroomUseCase.JoinChatroom(username, chatroom)

	// Assert: Ensure no error and message sent correctly
	assert.NoError(t, err, "Expected no error when joining chatroom")

	expectedMessage := domain.Message{
		Type:     domain.SystemMessageType,
		Username: username,
		Chatroom: chatroom,
		Content:  "test has joined the chatroom",
	}

	// Convert expected message to JSON string for comparison
	expectedJSON, _ := json.Marshal(expectedMessage)
	actualJSON := mockNats.PublishedMessages[0]

	assert.Equal(t, string(expectedJSON), actualJSON, "Expected output does not match")
}


// Test function for leaving a chatroom
func TestLeaveChatroom(t *testing.T) {
	// Arrange: Set up mock dependencies
	mockRedis := &MockRedisClient{users: make(map[string][]string)}
	mockNats := &MockNATSClient{}

	chatroomUseCase := application.NewChatroomUseCase(mockNats, mockRedis)

	// Test data
	username := "test_user"
	chatroom := "x"

	// Add user to the chatroom first
	_ = mockRedis.AddUserToChatroom(chatroom, username)

	// Act: Call LeaveChatroom method
	err := chatroomUseCase.LeaveChatroom(username, chatroom)

	// Assert: Ensure no error when leaving the chatroom
	assert.NoError(t, err, "Expected no error when leaving chatroom")

	expectedMessage := domain.Message{
		Username: "server",
		Chatroom: chatroom,
		Content:  "test_user has left the chatroom",
	}

	// Convert expected message to JSON string for comparison
	expectedJSON, _ := json.Marshal(expectedMessage)
	actualJSON := mockNats.PublishedMessages[0]

	assert.Equal(t, string(expectedJSON), actualJSON, "Expected output does not match")
}

// Test function for two users communicating in the same chatroom
func TestChatroomMessaging(t *testing.T) {
	// Arrange: Set up mock dependencies
	mockRedis := &MockRedisClient{users: make(map[string][]string)}
	mockNats := &MockNATSClient{}

	chatroomUseCase := application.NewChatroomUseCase(mockNats, mockRedis)

	// Test data
	username1 := "user_test_1"
	username2 := "user_test_2"
	chatroom := "test"

	// Act: Both users join the chatroom
	err := chatroomUseCase.JoinChatroom(username1, chatroom)
	assert.NoError(t, err, "Expected no error when user 1 joins chatroom")

	err = chatroomUseCase.JoinChatroom(username2, chatroom)
	assert.NoError(t, err, "Expected no error when user 2 joins chatroom")

	// User 1 sends a message
	msg1Content := "hi from test one"
	err = chatroomUseCase.SendMessage(username1, chatroom, msg1Content)
	assert.NoError(t, err, "Expected no error when user 1 sends message")

	// User 2 sends a message
	msg2Content := "hi from test two"
	err = chatroomUseCase.SendMessage(username2, chatroom, msg2Content)
	assert.NoError(t, err, "Expected no error when user 2 sends message")

	// Assert: Validate the messages received by both users
	expectedMessages := []domain.Message{
		{
			Username: username1,
			Chatroom: chatroom,
			Content:  msg1Content,
		},
		{
			Username: username2,
			Chatroom: chatroom,
			Content:  msg2Content,
		},
	}

	// Convert expected messages to JSON strings for comparison
	expectedJSON1, _ := json.Marshal(expectedMessages[0])
	expectedJSON2, _ := json.Marshal(expectedMessages[1])

	// Check if both messages are published to NATS
	assert.Contains(t, mockNats.PublishedMessages, string(expectedJSON1), "Expected first message in NATS")
	assert.Contains(t, mockNats.PublishedMessages, string(expectedJSON2), "Expected second message in NATS")
}