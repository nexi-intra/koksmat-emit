---
title: Observability
chatgpt-prompt: should I make a new instance of the observability in each module or should I go for a global var
chatgpt-model: o1-mini
---

> should I make a new instance of the observability in each module or should I go for a global var

---

Great question! Deciding between using a **global variable** or **creating new instances** of your `Observability` struct in each module is crucial for the maintainability, scalability, and testability of your Go application. Let’s delve into the pros and cons of each approach and explore the best practices to guide your decision.

---

## Table of Contents

1. [Understanding the Options](#1-understanding-the-options)
2. [Global Variable Approach](#2-global-variable-approach)
   - [Pros](#pros)
   - [Cons](#cons)
3. [Dependency Injection Approach](#3-dependency-injection-approach)
   - [Pros](#pros-1)
   - [Cons](#cons-1)
4. [Best Practice Recommendation](#4-best-practice-recommendation)
5. [Implementing Dependency Injection](#5-implementing-dependency-injection)
   - [Step-by-Step Guide](#step-by-step-guide)
   - [Example Code](#example-code)
6. [Advanced Patterns](#6-advanced-patterns)
   - [Using Interfaces](#using-interfaces)
   - [Singleton Pattern (With Caution)](#singleton-pattern-with-caution)
7. [Conclusion](#7-conclusion)

---

## 1. Understanding the Options

When structuring your application, especially concerning cross-cutting concerns like logging and metrics, you typically face two main architectural choices:

- **Global Variable**: Declare a single, globally accessible instance of the `Observability` struct that all modules can access.
- **Dependency Injection**: Pass an instance of the `Observability` struct to each module or component that requires it, typically through function parameters or struct fields.

Each approach has its own set of advantages and trade-offs, which we’ll explore next.

---

## 2. Global Variable Approach

### Pros

1. **Simplicity**: Easy to implement. Any part of the application can access the `Observability` instance without needing to pass it around.
2. **Convenience**: Reduces the need to modify function signatures or struct definitions to accommodate the `Observability` dependency.

### Cons

1. **Tight Coupling**: Modules become tightly coupled to the global `Observability` instance, making them less modular and harder to reuse.
2. **Testing Difficulties**: Harder to mock or replace the `Observability` instance during testing, leading to potential challenges in unit testing.
3. **Global State Risks**: Increases the risk of unintended side effects since any part of the code can modify the global instance.
4. **Concurrency Issues**: If not handled carefully, global variables can lead to race conditions in concurrent applications.

### Example of Global Variable

```go
// internal/observability/global.go
package observability

var Global *Observability

func Initialize() (*Observability, error) {
    obs, err := NewObservability()
    if err != nil {
        return nil, err
    }
    Global = obs
    return obs, nil
}
```

```go
// main.go
package main

import (
    "myapp/internal/observability"
)

func main() {
    obs, err := observability.Initialize()
    if err != nil {
        // handle error
    }
    defer obs.Shutdown()

    // Access globally
    observability.Global.Info("Server started")
}
```

---

## 3. Dependency Injection Approach

### Pros

1. **Loose Coupling**: Modules depend on abstractions rather than concrete implementations, enhancing modularity and reusability.
2. **Enhanced Testability**: Easily replace or mock the `Observability` instance during testing, allowing for more effective unit tests.
3. **Clear Dependencies**: Makes the dependencies of each module explicit, improving code readability and maintainability.
4. **Concurrency Safety**: Reduces risks associated with global state in concurrent environments.

### Cons

1. **Increased Boilerplate**: Requires passing the `Observability` instance through multiple layers of function calls or struct initializations.
2. **Setup Complexity**: Initial setup can be more involved, especially in larger applications with many dependencies.

### Example of Dependency Injection

```go
// internal/observability/observability.go
package observability

type Observability struct {
    Logger          *zap.Logger
    MetricsRegistry *prometheus.Registry
    HttpRequests    *prometheus.CounterVec
    MetricsHandler  http.Handler
    // Other fields...
}

// NewObservability initializes and returns a new Observability instance.
func NewObservability() (*Observability, error) {
    // Initialization logic...
}

// Shutdown gracefully shuts down the observability components.
func (o *Observability) Shutdown() error {
    return o.Logger.Sync()
}
```

```go
// main.go
package main

import (
    "myapp/internal/observability"
)

type App struct {
    Obs *observability.Observability
    // Other dependencies...
}

func NewApp(obs *observability.Observability) *App {
    return &App{
        Obs: obs,
        // Initialize other dependencies...
    }
}

func main() {
    obs, err := observability.NewObservability()
    if err != nil {
        // Handle error
    }
    defer obs.Shutdown()

    app := NewApp(obs)

    // Use app.Obs in your handlers and modules
}
```

---

## 4. Best Practice Recommendation

**Favor Dependency Injection over Global Variables.**

While global variables offer simplicity, dependency injection provides significant advantages in terms of modularity, testability, and maintainability. In larger and more complex applications, dependency injection helps manage dependencies more effectively, reducing tight coupling and enhancing code quality.

**Key Reasons to Prefer Dependency Injection:**

- **Testability**: Easily inject mocks or stubs for testing.
- **Maintainability**: Clear and explicit dependencies make the codebase easier to understand and maintain.
- **Scalability**: Better suited for scaling applications as it manages dependencies systematically.

---

## 5. Implementing Dependency Injection

Let’s walk through a step-by-step guide to implement dependency injection for your `Observability` struct across different modules.

### Step-by-Step Guide

1. **Initialize `Observability` in `main.go`**: Create a single instance during application startup.
2. **Pass `Observability` to Handlers and Modules**: Inject the instance where needed.
3. **Use Contexts for Request-Level Dependencies** (Optional): Pass `Observability` via `context.Context` if it's required in deeper call stacks.

### Example Code

#### a. Defining the `Observability` Struct

```go
// internal/observability/observability.go
package observability

import (
    "fmt"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "github.com/spf13/viper"
)

type Observability struct {
    Logger          *zap.Logger
    MetricsRegistry *prometheus.Registry
    HttpRequests    *prometheus.CounterVec
    MetricsHandler  http.Handler
}

type Config struct {
    LogLevel       string `mapstructure:"LOG_LEVEL"`
    LogOutputPaths string `mapstructure:"LOG_OUTPUT_PATHS"`
    MetricsPort    string `mapstructure:"METRICS_PORT"`
    ServiceName    string `mapstructure:"SERVICE_NAME"`
}

func NewObservability() (*Observability, error) {
    viper.AutomaticEnv()

    // Explicitly bind environment variables
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

    // Set defaults
    viper.SetDefault("LOG_LEVEL", "info")
    viper.SetDefault("LOG_OUTPUT_PATHS", "stdout")
    viper.SetDefault("METRICS_PORT", "9090")
    viper.SetDefault("SERVICE_NAME", "my-go-service")

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unable to decode into struct: %w", err)
    }

    // Initialize Logger
    logger, err := initLogger(cfg.LogLevel, cfg.LogOutputPaths)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize logger: %w", err)
    }

    // Initialize Metrics
    metricsRegistry := prometheus.NewRegistry()

    httpRequests := prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"path"},
    )
    metricsRegistry.MustRegister(httpRequests)

    metricsHandler := promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{})

    return &Observability{
        Logger:          logger,
        MetricsRegistry: metricsRegistry,
        HttpRequests:    httpRequests,
        MetricsHandler:  metricsHandler,
    }, nil
}

func initLogger(level, outputPath string) (*zap.Logger, error) {
    var zapLevel zapcore.Level
    if err := zapcore.ParseLevel(level, &zapLevel); err != nil {
        return nil, fmt.Errorf("invalid log level: %w", err)
    }

    cfg := zap.Config{
        Level:            zap.NewAtomicLevelAt(zapLevel),
        Development:      false,
        Encoding:         "json",
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

func (o *Observability) Shutdown() error {
    return o.Logger.Sync()
}

func (o *Observability) Verbose(msg string, fields ...zap.Field) {
    o.Logger.Debug(msg, fields...)
}

func (o *Observability) Info(msg string, fields ...zap.Field) {
    o.Logger.Info(msg, fields...)
}

func (o *Observability) Warning(msg string, fields ...zap.Field) {
    o.Logger.Warn(msg, fields...)
}

func (o *Observability) Error(msg string, fields ...zap.Field) {
    o.Logger.Error(msg, fields...)
}

func (o *Observability) InstrumentedHandler(path string, handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        o.HttpRequests.WithLabelValues(path).Inc()
        o.Info("Handling request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.String("remote_addr", r.RemoteAddr),
        )
        handlerFunc(w, r)
    }
}
```

#### b. Defining Your Application Structure

```go
// internal/app/app.go
package app

import (
    "net/http"

    "myapp/internal/observability"
)

type App struct {
    Obs *observability.Observability
    // Other dependencies can be added here
}

func NewApp(obs *observability.Observability) *App {
    return &App{
        Obs: obs,
        // Initialize other dependencies here
    }
}

func (a *App) Routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/hello", a.Obs.InstrumentedHandler("/hello", a.HelloHandler))
    mux.HandleFunc("/verbose", a.Obs.InstrumentedHandler("/verbose", a.VerboseHandler))
    mux.HandleFunc("/health", a.Obs.InstrumentedHandler("/health", a.HealthHandler))

    return mux
}

func (a *App) HelloHandler(w http.ResponseWriter, r *http.Request) {
    a.Obs.Info("Processing /hello request")
    w.Write([]byte("Hello, World!"))
}

func (a *App) VerboseHandler(w http.ResponseWriter, r *http.Request) {
    a.Obs.Verbose("Verbose endpoint hit")
    w.Write([]byte("This is a verbose message"))
}

func (a *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
    a.Obs.Info("Health check endpoint hit")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```

#### c. Initializing in `main.go`

```go
// main.go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "myapp/internal/app"
    "myapp/internal/observability"

    "go.uber.org/zap"
)

func main() {
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
    application := app.NewApp(obs)

    // Setup HTTP handlers
    mux := application.Routes()

    // Setup /metrics endpoint
    mux.(*http.ServeMux).Handle("/metrics", obs.MetricsHandler)

    // Start the server
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Graceful shutdown
    go func() {
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
}
```

---

## 6. Advanced Patterns

### a. Using Interfaces for Observability

To further enhance testability and flexibility, you can define an interface for your `Observability` methods. This allows you to inject mock implementations during testing.

```go
// internal/observability/observer.go
package observability

import "go.uber.org/zap"

type Observer interface {
    Verbose(msg string, fields ...zap.Field)
    Info(msg string, fields ...zap.Field)
    Warning(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
}
```

Update the `Observability` struct to implement this interface:

```go
// internal/observability/observability.go

type Observability struct {
    // Existing fields...
}

func (o *Observability) Verbose(msg string, fields ...zap.Field) {
    o.Logger.Debug(msg, fields...)
}

// Implement other methods similarly...
```

In your application modules, depend on the `Observer` interface instead of the concrete `Observability` struct:

```go
// internal/app/app.go
package app

import (
    "net/http"

    "myapp/internal/observability"
)

type App struct {
    Obs observability.Observer
    // Other dependencies...
}

func NewApp(obs observability.Observer) *App {
    return &App{
        Obs: obs,
        // Initialize other dependencies...
    }
}
```

During testing, you can create a mock implementation of the `Observer` interface:

```go
// internal/observability/mock_observer.go
package observability

import "go.uber.org/zap"

type MockObserver struct {
    // Add fields to track calls if needed
}

func (m *MockObserver) Verbose(msg string, fields ...zap.Field) {}
func (m *MockObserver) Info(msg string, fields ...zap.Field) {}
func (m *MockObserver) Warning(msg string, fields ...zap.Field) {}
func (m *MockObserver) Error(msg string, fields ...zap.Field) {}
```

### b. Singleton Pattern (With Caution)

While the singleton pattern can enforce a single instance, it introduces global state and tight coupling, which are generally discouraged in Go. If you choose to use it, be aware of the trade-offs.

```go
// internal/observability/singleton.go
package observability

import (
    "sync"
)

var (
    instance *Observability
    once     sync.Once
)

func GetInstance() (*Observability, error) {
    var err error
    once.Do(func() {
        instance, err = NewObservability()
    })
    return instance, err
}
```

**Caution**: This approach reintroduces some of the downsides of global variables, such as testing difficulties and tight coupling.

---

## 7. Conclusion

Choosing between **global variables** and **dependency injection** for managing your `Observability` instance hinges on the specific needs and scale of your application. However, adhering to best practices in Go development generally favors **dependency injection** due to its numerous advantages:

- **Enhanced Testability**: Easily mock dependencies during testing without relying on global state.
- **Loose Coupling**: Promotes modularity, making your codebase easier to maintain and extend.
- **Clear Dependencies**: Makes the dependencies of each module explicit, improving code readability.

While the **global variable** approach offers simplicity, it sacrifices the modularity and testability that become increasingly important as your application grows. By leveraging dependency injection, you establish a more robust and scalable architecture that aligns with Go’s emphasis on simplicity and clarity.

**Final Recommendation**: **Use dependency injection** to pass a single instance of the `Observability` struct to modules that require it. This approach balances maintainability, scalability, and testability, ensuring your application remains robust as it evolves.

---

**Happy Coding! 🚀**

If you have further questions or need more detailed examples, feel free to ask!
