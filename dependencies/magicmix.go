package dependencies

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	natsutil "github.com/nexi-intra/koksmat-emit/dependencies/nats"
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

func (c *MicroService) Request(subject string, req interface{}, resp interface{}, timeout time.Duration) error {
	// Safety check: resp must be a pointer, or json.Unmarshal will fail
	if reflect.ValueOf(resp).Kind() != reflect.Ptr {
		return errors.New("resp argument must be a pointer")
	}

	// Marshal the request into bytes (JSON in this example)
	reqData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	log.Println("Sending request", subject, string(reqData))
	// Send the request to NATS
	msg, err := c.client.Request(subject, reqData, timeout)
	if err != nil {
		return fmt.Errorf("NATS request failed: %w", err)
	}
	log.Println("Received response", string(msg.Data))
	// Unmarshal the response into 'resp'
	if err := json.Unmarshal(msg.Data, resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
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

	// Prepare request
	req := MyRequest{Message: "Hello Service", Args: []string{"query", "mix", "select 1"}}

	// Prepare a variable to hold the response
	// Note that we must pass a pointer: &resp
	var resp MyResponse

	// Make the request
	err = service.Request("magic-mix.app", req, &resp, 5*time.Second)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}

	// We now have an unmarshaled response in 'resp'
	fmt.Printf("Received reply: %+v\n", resp)

	// Send a request

}
