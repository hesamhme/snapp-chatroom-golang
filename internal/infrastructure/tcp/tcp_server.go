package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"appchat/internal/domain"
	"appchat/internal/infrastructure/nats"
	"appchat/internal/infrastructure/redis"
	"github.com/sirupsen/logrus"
)

type TCPHandler struct {
	natsClient *nats.NATSClient
	redisClient *redis.RedisClient
	clients    map[string][]net.Conn
	lock       sync.Mutex
}

func NewTCPHandler(natsClient *nats.NATSClient, redisClient *redis.RedisClient) *TCPHandler {
	return &TCPHandler{
		natsClient: natsClient,
		redisClient: redisClient,
		clients:    make(map[string][]net.Conn),
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
		th.lock.Lock()
		th.clients[msg.Chatroom] = append(th.clients[msg.Chatroom], conn)
		th.lock.Unlock()

		logrus.Infof("Adding user %s to chatroom %s", msg.Username, msg.Chatroom)
		if err := th.redisClient.AddUserToChatroom(msg.Chatroom, msg.Username); err != nil {
			logrus.Errorf("Failed to add user to Redis: %v", err)
		}

		logrus.Infof("User %s joined chatroom %s", msg.Username, msg.Chatroom)
		th.natsClient.PublishMessage("chatroom."+msg.Chatroom, string(scanner.Bytes()))
	}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			logrus.Errorf("Error decoding message: %v", err)
			continue
		}

		if msg.Content == "#users" {
			users, err := th.redisClient.GetUsersInChatroom(msg.Chatroom)
			if err != nil {
				logrus.Errorf("Error retrieving users from chatroom %s: %v", msg.Chatroom, err)
				conn.Write([]byte("Error retrieving users\n"))
				continue
			}
			userList := fmt.Sprintf("Users in chatroom %s: %v", msg.Chatroom, users)
			logrus.Infof("Sending user list to %s: %s", msg.Username, userList)
			conn.Write([]byte(userList + "\n"))
		} else {
			logrus.Infof("Received message from %s: %s", msg.Username, msg.Content)
			th.natsClient.PublishMessage("chatroom."+msg.Chatroom, string(scanner.Bytes()))
		}
	}
	conn.Close()
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
