// internal/app/app.go
package emitter

import (
	"net/http"

	"github.com/nexi-intra/koksmat-emit/internal/observability"
	"go.uber.org/zap"
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

func (a *App) SaveWebhook(endpoint string, body []byte) {
	a.Obs.Verbose("Saving webhook", zap.String("endpoint", endpoint), zap.String("body", string(body)))
}
