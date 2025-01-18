/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/nexi-intra/koksmat-emit/api"
	"github.com/nexi-intra/koksmat-emit/internal/emitter"
	"github.com/nexi-intra/koksmat-emit/internal/observability"
	"github.com/spf13/cobra"

	"context"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the services.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

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
		app := emitter.NewApp(obs)

		// Setup HTTP handlers
		mux := app.Routes()

		// Setup /metrics endpoint
		mux.(*http.ServeMux).Handle("/metrics", obs.MetricsHandler)

		// Start the server
		server := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}

		// Graceful shutdown
		go func() {
			api.Start(":4321", app)
			obs.Info("Starting server on :8080")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				obs.Error("Server failed", zap.Error(err))
			}

		}()

		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		obs.Info("Shutting down server...")

		if err := server.Shutdown(context.Background()); err != nil {
			obs.Error("Server shutdown failed", zap.Error(err))
		}

		obs.Info("Server exited gracefully")

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
