package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// NATSConfig holds the configuration for NATS connection
type NATSConfig struct {
	URL           string        // NATS server URL
	Username      string        // Optional: Username for authentication
	Password      string        // Optional: Password for authentication
	ReconnectWait time.Duration // Time to wait before attempting reconnection
	MaxReconnects int           // Maximum number of reconnection attempts
}

// NATSClient encapsulates the NATS connection and options
type NATSClient struct {
	conn *nats.Conn
}

// NewNATSClient initializes and returns a NATSClient
func NewNATSClient(cfg NATSConfig) (*NATSClient, error) {
	opts := []nats.Option{
		nats.Name("Go NATS Utility"),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("Disconnected due to: %v, will attempt reconnects\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Reconnected to %v\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Printf("Connection closed. Reason: %v\n", nc.LastError())
		}),
	}

	if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, nats.UserInfo(cfg.Username, cfg.Password))
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, err
	}

	return &NATSClient{conn: nc}, nil
}

// Close closes the NATS connection
func (c *NATSClient) Close() {
	if c.conn != nil && !c.conn.IsClosed() {
		c.conn.Close()
	}
}

// Publish publishes a message to a specific subject
func (c *NATSClient) Publish(subject string, data []byte) error {
	return c.conn.Publish(subject, data)
}

// Subscribe subscribes to a subject with a message handler
func (c *NATSClient) Subscribe(subject string, handler func(m *nats.Msg)) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}

// Request sends a request and waits for a reply
func (c *NATSClient) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	return c.conn.Request(subject, data, timeout)
}

// JetStream context (optional, if using JetStream)
func (c *NATSClient) JetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	return c.conn.JetStream(opts...)
}

// Example usage within the utility (optional)
func Example() {
	cfg := NATSConfig{
		URL:           nats.DefaultURL,
		ReconnectWait: 2 * time.Second,
		MaxReconnects: 10,
	}

	client, err := NewNATSClient(cfg)
	if err != nil {
		fmt.Printf("Error connecting to NATS: %v\n", err)
		return
	}
	defer client.Close()

	// Publish example
	err = client.Publish("updates", []byte("Hello NATS"))
	if err != nil {
		fmt.Printf("Error publishing message: %v\n", err)
	}

	// Subscribe example
	_, err = client.Subscribe("updates", func(m *nats.Msg) {
		fmt.Printf("Received message: %s\n", string(m.Data))
	})
	if err != nil {
		fmt.Printf("Error subscribing to subject: %v\n", err)
	}

	// Keep the example running
	select {}
}
