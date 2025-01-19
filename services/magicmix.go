package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	natsutil "github.com/nexi-intra/koksmat-emit/services/nats"
	"github.com/spf13/viper"
)

func connect(url string) (*natsutil.NATSClient, error) {
	cfg := natsutil.NATSConfig{
		URL:           url,
		ReconnectWait: 2 * time.Second, // Wait 2 seconds before reconnect
		MaxReconnects: 10,              // Attempt to reconnect 10 times
	}
	return natsutil.NewNATSClient(cfg)

}

// MicroService encapsulates the NATS connection and options
type MicroService struct {
	client *natsutil.NATSClient
}

func NewMicroserviceConnection() (*MicroService, error) {
	client, err := connect(viper.GetString("NATS_URL"))
	if err != nil {
		return nil, err
	}
	return &MicroService{client: client}, nil
}

// Close closes the NATS connection
func (c *MicroService) Close() {
	c.client.Close()
}

func (c *MicroService) Request(subject string, args []string, body string, timeout time.Duration) (*string, error) {
	// Safety check: resp must be a pointer, or json.Unmarshal will fail

	type MyRequest struct {
		Args    []string `json:"args"`
		Body    string   `json:"body"`
		Channel string   `json:"channel"`
	}

	req := MyRequest{Args: args, Body: body, Channel: "noma2"}

	// Marshal the request into bytes (JSON in this example)
	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Println("Sending request", subject, string(reqData))
	// Send the request to NATS
	msg, err := c.client.Request(subject, reqData, timeout)
	if err != nil {
		return nil, fmt.Errorf("NATS request failed: %w", err)
	}
	responseData := string(msg.Data)
	log.Println("Received response", responseData)

	return &responseData, nil
}

type MyRequest struct {
	Args    []string `json:"args"`
	Message string   `json:"message"`
}

type MyResponse struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

func sample() {

	service, err := NewMicroserviceConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MagicMix: %v", err)
	}
	defer service.Close()

	// Make the request
	result, err := service.Request("magic-mix.app", []string{"query", "mix", "select 1"}, "", 5*time.Second)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}

	// We now have an unmarshaled response in 'resp'
	fmt.Printf("Received reply: %+v\n", result)

	// Send a request

}
