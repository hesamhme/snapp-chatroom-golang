package nats

import (
	"github.com/nats-io/nats.go"
)

type NATSClientInterface interface {
    PublishMessage(subject, message string) error
    Subscribe(subject string, handler func(string))
    Close()
}


type NATSClient struct {
	conn *nats.Conn
}

func NewNATSConnection(url string) (*NATSClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NATSClient{conn: conn}, nil
}

func (nc *NATSClient) PublishMessage(subject, message string) error {
    return nc.conn.Publish(subject, []byte(message))
}


func (nc *NATSClient) Subscribe(subject string, handler func(string)) {
	nc.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(string(m.Data))
	})
}

// Close closes the NATS connection
func (nc *NATSClient) Close() {
	nc.conn.Close()
}