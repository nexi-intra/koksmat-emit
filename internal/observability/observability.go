// internal/observability/observability.go
package observability

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Observability encapsulates logging and metrics functionalities.
type Observability struct {
	Logger          *zap.Logger
	MetricsRegistry *prometheus.Registry
	HttpRequests    *prometheus.CounterVec
	MetricsHandler  http.Handler
}

// Config holds the configuration for Observability.
type Config struct {
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	LogOutputPaths string `mapstructure:"LOG_OUTPUT_PATHS"`
	MetricsPort    string `mapstructure:"METRICS_PORT"`
	ServiceName    string `mapstructure:"SERVICE_NAME"`
}

// NewObservability initializes logging and metrics based on environment variables.
func NewObservability() (*Observability, error) {
	// Initialize Viper to read from environment variables.
	viper.AutomaticEnv()
	// Do not set an environment variable prefix or key replacer to preserve capitalization.

	// Define environment variable bindings explicitly.
	// This ensures that each config field maps to its exact environment variable.
	if err := viper.BindEnv("LOG_LEVEL", "LOG_LEVEL"); err != nil {
		return nil, fmt.Errorf("error binding LOG_LEVEL: %w", err)
	}
	if err := viper.BindEnv("LOG_OUTPUT_PATHS", "LOG_OUTPUT_PATHS"); err != nil {
		return nil, fmt.Errorf("error binding LOG_OUTPUT_PATHS: %w", err)
	}
	if err := viper.BindEnv("METRICS_PORT", "METRICS_PORT"); err != nil {
		return nil, fmt.Errorf("error binding METRICS_PORT: %w", err)
	}
	if err := viper.BindEnv("SERVICE_NAME", "SERVICE_NAME"); err != nil {
		return nil, fmt.Errorf("error binding SERVICE_NAME: %w", err)
	}

	// Optionally, set default values.
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_OUTPUT_PATHS", "stdout")
	viper.SetDefault("METRICS_PORT", "9090")
	viper.SetDefault("SERVICE_NAME", "my-go-service")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Initialize Logger.
	logger, err := initLogger(cfg.LogLevel, cfg.LogOutputPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize Metrics.
	metricsRegistry := prometheus.NewRegistry()

	// Initialize HTTP Requests Counter.
	httpRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)
	metricsRegistry.MustRegister(httpRequests)

	// Initialize Metrics Handler.
	metricsHandler := promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{})

	return &Observability{
		Logger:          logger,
		MetricsRegistry: metricsRegistry,
		HttpRequests:    httpRequests,
		MetricsHandler:  metricsHandler,
	}, nil
}

// initLogger sets up the Zap logger based on the provided log level and output paths.
func initLogger(level, outputPath string) (*zap.Logger, error) {
	var zapLevel zapcore.Level
	if _, err := zapcore.ParseLevel(level); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json", // Use "console" for human-readable logs.
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{outputPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("unable to build logger: %w", err)
	}

	return logger, nil
}

// Shutdown gracefully shuts down the logger.
func (o *Observability) Shutdown() error {
	return o.Logger.Sync()
}

// Verbose logs a message at the DEBUG level.
func (o *Observability) Verbose(msg string, fields ...zap.Field) {
	o.Logger.Debug(msg, fields...)
}

// Info logs a message at the INFO level.
func (o *Observability) Info(msg string, fields ...zap.Field) {
	o.Logger.Info(msg, fields...)
}

// Warning logs a message at the WARN level.
func (o *Observability) Warning(msg string, fields ...zap.Field) {
	o.Logger.Warn(msg, fields...)
}

// Error logs a message at the ERROR level.
func (o *Observability) Error(msg string, fields ...zap.Field) {
	o.Logger.Error(msg, fields...)
}

// InstrumentedHandler wraps an HTTP handler with observability (metrics).
func (o *Observability) InstrumentedHandler(path string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Increment Prometheus counter directly using the reference.
		o.HttpRequests.WithLabelValues(path).Inc()

		// Log the request.
		o.Info("Handling request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)

		// Call the actual handler.
		handlerFunc(w, r)
	}
}
