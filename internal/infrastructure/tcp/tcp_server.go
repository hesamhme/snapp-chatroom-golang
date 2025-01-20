package tcp

import (
	"appchat/internal/domain"
	"appchat/internal/infrastructure/nats"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
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
		fmt.Println("Error starting TCP handler:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP handler started on port", port)

	go th.subscribeToNATS()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go th.handleConnection(conn)
	}
}

func (th *TCPHandler) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	var msg domain.Message

	if scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Println("Error decoding message:", err)
			return
		}
		th.lock.Lock()
		th.clients[msg.Chatroom] = append(th.clients[msg.Chatroom], conn)
		th.lock.Unlock()

		th.natsClient.PublishMessage("chatroom."+msg.Chatroom, string(scanner.Bytes()))
	}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Println("Error decoding message:", err)
			continue
		}
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
					fmt.Println("Error sending message to client:", err)
				}
			}
			th.lock.Unlock()
		}
	})
}
