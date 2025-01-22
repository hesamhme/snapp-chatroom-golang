package tcp

import (
	"appchat/internal/application"
	"appchat/internal/domain"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type TCPHandler struct {
	useCase *application.ChatroomUseCase
	clients map[string][]net.Conn
	lock    sync.Mutex
}

func NewTCPHandler(useCase *application.ChatroomUseCase) *TCPHandler {
	return &TCPHandler{
		useCase: useCase,
		clients: make(map[string][]net.Conn),
	}
}

func (th *TCPHandler) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logrus.Fatalf("Error starting TCP handler: %v", err)
	}
	defer listener.Close()

	logrus.Infof("TCP handler started on port %s", port)

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
	th.useCase.JoinChatroom(msg.Username, msg.Chatroom)
}

func (th *TCPHandler) processMessage(conn net.Conn, msg domain.Message) {
	switch msg.Content {
	case "#users":
		users, _ := th.useCase.GetUsers(msg.Chatroom)
		conn.Write([]byte(fmt.Sprintf("Users in chatroom %s: %v\n", msg.Chatroom, users)))
	case "#rooms":
		rooms, _ := th.useCase.GetChatrooms()
		conn.Write([]byte(fmt.Sprintf("Active chatrooms: %v\n", rooms)))
	default:
		th.useCase.SendMessage(msg.Username, msg.Chatroom, msg.Content)
	}
}
