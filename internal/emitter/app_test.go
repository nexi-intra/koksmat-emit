// internal/app/app.go
package emitter

import (
	"fmt"
	"os"
	"testing"

	"github.com/nexi-intra/koksmat-emit/services"

	"github.com/nexi-intra/koksmat-emit/internal/observability"
	"go.uber.org/zap"
)

func TestApp_SaveWebhook(t *testing.T) {

	type fields struct {
		Obs *observability.Observability
		Mix *services.MicroService
	}
	type args struct {
		endpoint string
		body     string
	}
	// Initialize Observability
	obs, err := observability.NewObservability()
	if err != nil {
		fmt.Printf("Failed to initialize observability: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := obs.Shutdown(); err != nil {
			obs.Error("Error during shutdown", zap.Error(err))
		}
	}()

	// Initialize Application
	app := NewApp(obs)

	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "Test SaveWebhook",

			args: args{endpoint: "http://localhost:8080", body: `
		{
		"tenant": "",
		"searchindex": "",
		"name": "name",
		"description": "",
		"source": "\u007B\u007D",
		"tag": ""
		
		}

			
			`},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := app.SaveWebhook(tt.args.endpoint, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("App.SaveWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
