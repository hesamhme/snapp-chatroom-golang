package tcp

import (
	"appchat/internal/domain"
	"appchat/internal/infrastructure/nats"
	"bufio"
	"encoding/json"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type TCPHandler struct {
	natsClient *nats.NATSClient
	clients    map[string][]net.Conn
	lock       sync.Mutex
}

func NewTCPHandler(natsClient *nats.NATSClient) *TCPHandler {
	return &TCPHandler{
		natsClient: natsClient,
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

		logrus.Infof("User %s joined chatroom %s", msg.Username, msg.Chatroom)
		th.natsClient.PublishMessage("chatroom."+msg.Chatroom, string(scanner.Bytes()))
	}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			logrus.Errorf("Error decoding message: %v", err)
			continue
		}
		logrus.Infof("Received message from %s: %s", msg.Username, msg.Content)
		th.natsClient.PublishMessage("chatroom."+msg.Chatroom, string(scanner.Bytes()))
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
