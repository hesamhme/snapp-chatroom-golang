package domain

type Message struct {
	Username string `json:"username"`
	Chatroom string `json:"chatroom"`
	Content  string `json:"content"`
}