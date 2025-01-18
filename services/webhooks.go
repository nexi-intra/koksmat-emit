package services

import (
	"github.com/nexi-intra/koksmat-emit/dependencies"
	"github.com/nexi-intra/koksmat-emit/internal/emitter"
	"go.uber.org/zap"
)

type WebhooksService struct {
	mixService *dependencies.MicroService
	app        *emitter.App
}

func NewWebhooksService(app *emitter.App) *WebhooksService {
	mixService, err := dependencies.NewMicroserviceConnection()
	if err != nil {
		app.Obs.Error("Failed to connect to NATS", zap.Error(err))
		return nil
	}
	return &WebhooksService{
		mixService: mixService,
		app:        app,
	}
}

func (s *WebhooksService) Start() {
	s.app.Obs.Info("Starting webhooks service")
}

func (s *WebhooksService) Stop() {
	s.app.Obs.Info("Stopping webhooks service")
}

func (s *WebhooksService) Close() string {
	s.mixService.Close()
	return "webhooks"
}
