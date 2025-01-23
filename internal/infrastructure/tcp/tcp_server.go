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
	th.useCase.JoinChatroom(msg.Username, msg.Chatroom)
}

func (th *TCPHandler) handleUserExit(conn net.Conn, msg domain.Message) {
    th.lock.Lock()
    defer th.lock.Unlock()

    if err := th.useCase.LeaveChatroom(msg.Username, msg.Chatroom); err != nil {
        logrus.Errorf("Error leaving chatroom: %v", err)
        return
    }

    exitMessage := fmt.Sprintf("%s has left the chatroom", msg.Username)
    for _, clientConn := range th.clients[msg.Chatroom] {
        clientConn.Write([]byte(exitMessage + "\n"))
    }
    logrus.Infof("User %s left chatroom %s", msg.Username, msg.Chatroom)
}


func (th *TCPHandler) processMessage(conn net.Conn, msg domain.Message) {
    switch msg.Content {
    case "#users":
        users, _ := th.useCase.GetUsers(msg.Chatroom)
        conn.Write([]byte(fmt.Sprintf("Users in chatroom %s: %v\n", msg.Chatroom, users)))
    case "#rooms":
        rooms, _ := th.useCase.GetChatrooms()
        conn.Write([]byte(fmt.Sprintf("Active chatrooms: %v\n", rooms)))
    case "has left the chatroom":
        th.handleUserExit(conn, msg)
    default:
        // Process chat message
        th.useCase.SendMessage(msg.Username, msg.Chatroom, msg.Content)
    }
}


func (th *TCPHandler) subscribeToNATS() {
    go func() {
        th.useCase.SubscribeToMessages(func(msg domain.Message) {
            th.lock.Lock()
            defer th.lock.Unlock()
            
            messageContent := fmt.Sprintf("[%s] %s: %s", msg.Chatroom, msg.Username, msg.Content)
            for _, conn := range th.clients[msg.Chatroom] {
                _, err := conn.Write([]byte(messageContent + "\n"))
                if err != nil {
                    logrus.Errorf("Error sending message to client: %v", err)
                }
            }
            logrus.Infof("Broadcasted message in chatroom %s: %s", msg.Chatroom, msg.Content)
        })
    }()
}

