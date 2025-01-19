package domain

type Chatroom struct {
	Name string
}

func NewChatroom(name string) *Chatroom {
	return &Chatroom{Name: name}
}
