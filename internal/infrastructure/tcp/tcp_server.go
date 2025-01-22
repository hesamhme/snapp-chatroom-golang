package tcp

import (
	"appchat/internal/domain"
	"appchat/internal/infrastructure/nats"
	"appchat/internal/infrastructure/redis"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type TCPHandler struct {
	natsClient  *nats.NATSClient
	redisClient *redis.RedisClient
	clients     map[string][]net.Conn
	lock        sync.Mutex
}

func NewTCPHandler(natsClient *nats.NATSClient, redisClient *redis.RedisClient) *TCPHandler {
	return &TCPHandler{
		natsClient:  natsClient,
		redisClient: redisClient,
		clients:     make(map[string][]net.Conn),
	}
}

func (th *TCPHandler) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logrus.Fatalf("Error starting TCP handler: %v", err)
	}
	defer listener.Close()

	logrus.Infof("TCP handler started on port %s", port)

	go th.subscribeToNATS()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Errorf("Error accepting connection: %v", err)
			continue
		}
		logrus.Infof("New client connected: %s", conn.RemoteAddr())
		go th.handleConnection(conn)
	}
}

func (th *TCPHandler) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	var msg domain.Message

	if scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			logrus.Errorf("Error decoding message: %v", err)
			return
		}
		th.registerClient(conn, msg)
	}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			logrus.Errorf("Error decoding message: %v", err)
			continue
		}
		th.processMessage(conn, msg)
	}
	conn.Close()
}

func (th *TCPHandler) registerClient(conn net.Conn, msg domain.Message) {
	th.lock.Lock()
	defer th.lock.Unlock()
	th.clients[msg.Chatroom] = append(th.clients[msg.Chatroom], conn)
	th.redisClient.AddUserToChatroom(msg.Chatroom, msg.Username)
	th.redisClient.AddChatroom(msg.Chatroom)
	logrus.Infof("User %s joined chatroom %s", msg.Username, msg.Chatroom)
	th.natsClient.PublishMessage("chatroom."+msg.Chatroom, msg.Content)
}

func (th *TCPHandler) processMessage(conn net.Conn, msg domain.Message) {
	switch msg.Content {
	case "#users":
		th.sendUserList(conn, msg)
	case "#rooms":
		th.sendRoomList(conn, msg)
	default:
		logrus.Infof("Received message from %s: %s", msg.Username, msg.Content)
		th.natsClient.PublishMessage("chatroom."+msg.Chatroom, string(encodeMessage(msg)))
	}
}

func (th *TCPHandler) sendUserList(conn net.Conn, msg domain.Message) {
	users, err := th.redisClient.GetUsersInChatroom(msg.Chatroom)
	if err != nil {
		conn.Write([]byte("Error retrieving users\n"))
		return
	}
	response := fmt.Sprintf("Users in chatroom %s: %v", msg.Chatroom, users)
	conn.Write([]byte(response + "\n"))
}

func (th *TCPHandler) sendRoomList(conn net.Conn, msg domain.Message) {
	chatrooms, err := th.redisClient.GetChatrooms()
	if err != nil {
		conn.Write([]byte("Error retrieving chatrooms\n"))
		return
	}
	response := fmt.Sprintf("Active chatrooms: %v", chatrooms)
	conn.Write([]byte(response + "\n"))
}

func (th *TCPHandler) subscribeToNATS() {
	th.natsClient.Subscribe("chatroom.*", func(msg string) {
		var incomingMsg domain.Message
		if err := json.Unmarshal([]byte(msg), &incomingMsg); err == nil {
			th.lock.Lock()
			for _, conn := range th.clients[incomingMsg.Chatroom] {
				_, err := conn.Write(append([]byte(msg), '\n'))
				if err != nil {
					logrus.Errorf("Error sending message to client: %v", err)
				}
			}
			th.lock.Unlock()
			logrus.Infof("Broadcasting message to chatroom %s: %s", incomingMsg.Chatroom, incomingMsg.Content)
		}
	})
}

func encodeMessage(msg domain.Message) []byte {
	data, _ := json.Marshal(msg)
	return data
}
