// internal/app/app.go
package emitter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/nexi-intra/koksmat-emit/dependencies"
	"github.com/nexi-intra/koksmat-emit/internal/observability"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt"
)

type EventRecord struct {
	Tenant      string          `json:"tenant"`
	Searchindex string          `json:"searchindex"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Source      string          `json:"source"`
	Tag         string          `json:"tag"`
	Payload     json.RawMessage `json:"payload"`
}

// Define a secret key (ensure this is kept secure in production)
var jwtSecret = []byte("your-very-secure-secret")

// CreateJWT generates a JWT with the app_displayname as a payload claim
func CreateJWT(appDisplayName string) (string, error) {
	// Define the token claims
	claims := jwt.MapClaims{
		"app_displayname": appDisplayName,
		"exp":             time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":             time.Now().Unix(),
	}

	// Create a new token object specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type App struct {
	Obs *observability.Observability
	Mix *dependencies.MicroService
	// Other dependencies can be added here
}

func NewApp(obs *observability.Observability) *App {
	mixClient, err := dependencies.NewMicroserviceConnection()
	if err != nil {
		obs.Error("Failed to connect to MagicMix", zap.Error(err))
		return nil
	}
	return &App{
		Obs: obs,
		Mix: mixClient,
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

func (a *App) SaveWebhook(endpoint string, body string) error {
	a.Obs.Verbose("Saving webhook", zap.String("endpoint", endpoint), zap.String("body", string(body)))

	procedure := "create_event"

	token, err := CreateJWT("koksmat-emit")
	if err != nil {
		a.Obs.Error("Failed to create JWT", zap.Error(err))
		return err
	}
	if !json.Valid([]byte(body)) {
		a.Obs.Error("Invalid JSON", zap.String("body", body))
		return fmt.Errorf("invalid json: %s", body)
	}
	record := EventRecord{
		Tenant:      "",
		Searchindex: "",
		Name:        "webhook",
		Description: "webhook",
		Source:      "koksmat-emit",
		Tag:         endpoint,
		Payload:     json.RawMessage(body),
	}
	payload, err := json.Marshal(record)
	if err != nil {
		a.Obs.Error("Failed to marshal webhook record", zap.Error(err))
		return err
	}

	args := []string{"execute", "mix", procedure, token, body}

	result, err := a.Mix.Request("magic-mix.app", args, string(payload), 5*time.Second)
	if err != nil {
		a.Obs.Error("Failed to save webhook", zap.Error(err))
		return err
	}
	a.Obs.Info("Webhook saved", zap.String("result", *result))

	return nil
}
