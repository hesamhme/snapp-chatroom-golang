package domain


type MessageType string

const (
    ChatMessageType  MessageType = "chat"
    CommandMessageType MessageType = "command"
    ResponseMessageType MessageType = "response"
)


type Message struct {
    Type     MessageType `json:"type"`
    Username string      `json:"username"`
    Chatroom string      `json:"chatroom"`
    Content  string      `json:"content"`
}
