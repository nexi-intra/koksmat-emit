package api

import (
	"context"

	"github.com/nexi-intra/koksmat-emit/internal/emitter"
	"github.com/swaggest/usecase"
	"go.uber.org/zap"
)

// GitHubWebhookInput defines the expected payload from GitHub webhook.
type GitHubWebhookInput struct {
	Action     string `json:"action"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

// GitHubWebhookOutput defines the response for the webhook handler.
type GitHubWebhookOutput struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// webhook_GitHub creates a new usecase.Interactor to handle GitHub webhook events.
// It processes the webhook payload and responds based on the action specified in the input.
//
// Supported actions:
// - "created": Responds with a message indicating a new issue or PR was opened.
// - "opened": Responds with a message indicating a new issue or PR was opened.
// - "closed": Responds with a message indicating an issue or PR was closed.
// - Default: Responds with a message indicating the action is not handled.
//
// The interactor sets the title, description, and tags for the use case.
func webhook_GitHub(app *emitter.App) usecase.Interactor {
	// Create a new interactor for the webhook.

	u := usecase.NewInteractor(func(ctx context.Context, input GitHubWebhookInput, output *GitHubWebhookOutput) error {
		app.Obs.Info("Hook", zap.String("action", input.Action))

		// Example logic: respond based on the action.
		switch input.Action {
		case "created":
			*output = GitHubWebhookOutput{
				Message: "A new issue or PR was opened in " + input.Repository.Name,
				Status:  "success",
			}
		case "opened":
			*output = GitHubWebhookOutput{
				Message: "A new issue or PR was opened in " + input.Repository.Name,
				Status:  "success",
			}
			return nil
		case "closed":
			*output = GitHubWebhookOutput{
				Message: "An issue or PR was closed in " + input.Repository.Name,
				Status:  "success",
			}
			return nil
		default:
			*output = GitHubWebhookOutput{
				Message: "Action not handled: " + input.Action,
				Status:  "ignored",
			}
			return nil
		}
		return nil
	})

	// Describe the usecase.
	u.SetTitle("GitHub Webhook Handler")
	u.SetDescription("Handles POST requests from GitHub webhooks.")
	u.SetTags("GitHub")
	return u
}

/*
TODO: Adding HMAC validation ensures that the webhook payload is from GitHub and has not been tampered with. GitHub sends the payload with a `X-Hub-Signature-256` header, which contains a HMAC hash of the payload signed with your secret key. Here's how you can implement it in your Go program:

### Steps to Add HMAC Validation:

1. **Read the Secret Key:**
   Define a shared secret that GitHub and your server use for signing and validating payloads.

2. **Compute HMAC Hash:**
   Compute the HMAC-SHA256 hash of the payload using the secret key.

3. **Compare the Hash:**
   Compare the computed hash with the one provided in the `X-Hub-Signature-256` header.

4. **Reject Invalid Requests:**
   If the hashes don't match, reject the request with a 403 Forbidden response.

### Updated Example with HMAC Validation:

```go
package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/swaggest/usecase"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/nethttp"
)

// GitHubWebhookInput defines the expected payload from GitHub webhook.
type GitHubWebhookInput struct {
	Action     string `json:"action"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

// GitHubWebhookOutput defines the response for the webhook handler.
type GitHubWebhookOutput struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Shared secret for HMAC validation (configure this securely).
const secret = "your_shared_secret"

// validateHMAC validates the HMAC-SHA256 signature.
func validateHMAC(payload []byte, signature string) error {
	// Remove the "sha256=" prefix from the signature header.
	expectedSig := strings.TrimPrefix(signature, "sha256=")

	// Compute the HMAC-SHA256 of the payload.
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	computedSig := hex.EncodeToString(mac.Sum(nil))

	// Compare the computed signature with the expected signature.
	if !hmac.Equal([]byte(computedSig), []byte(expectedSig)) {
		return errors.New("invalid signature")
	}

	return nil
}

// HandleGitHubWebhook processes GitHub webhook events.
func HandleGitHubWebhook(ctx context.Context, input GitHubWebhookInput) (GitHubWebhookOutput, error) {
	// Example logic: respond based on the action.
	switch input.Action {
	case "opened":
		return GitHubWebhookOutput{
			Message: "A new issue or PR was opened in " + input.Repository.Name,
			Status:  "success",
		}, nil
	case "closed":
		return GitHubWebhookOutput{
			Message: "An issue or PR was closed in " + input.Repository.Name,
			Status:  "success",
		}, nil
	default:
		return GitHubWebhookOutput{
			Message: "Action not handled: " + input.Action,
			Status:  "ignored",
		}, nil
	}
}

func main() {
	// Create a new interactor for the webhook.
	webhookInteractor := usecase.NewInteractor(func(ctx context.Context, input GitHubWebhookInput) (GitHubWebhookOutput, error) {
		return HandleGitHubWebhook(ctx, input)
	})

	// Describe the usecase.
	webhookInteractor.SetTitle("GitHub Webhook Handler with HMAC Validation")
	webhookInteractor.SetDescription("Handles POST requests from GitHub webhooks and validates HMAC signatures.")
	webhookInteractor.SetTags("webhook", "GitHub")
	webhookInteractor.Input = GitHubWebhookInput{}
	webhookInteractor.Output = GitHubWebhookOutput{}

	// Create a custom REST handler for the interactor.
	handler := nethttp.NewHandler(webhookInteractor)
	handler.PrepareRequest = func(r *http.Request) (context.Context, interface{}, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, nil, rest.ErrRequestBody{Err: err}
		}

		// Restore the body for further processing by the interactor.
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		// Validate the HMAC signature.
		signature := r.Header.Get("X-Hub-Signature-256")
		if signature == "" {
			return nil, nil, rest.ErrRequestBody{Err: errors.New("missing signature header")}
		}

		if err := validateHMAC(body, signature); err != nil {
			return nil, nil, rest.ErrRequestBody{Err: err}
		}

		// Parse the JSON payload into the input struct.
		var input GitHubWebhookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, nil, rest.ErrRequestBody{Err: err}
		}

		return r.Context(), input, nil
	}

	handler.Method = http.MethodPost
	handler.Path = "/webhook"

	// Create the router and register the handler.
	router := http.NewServeMux()
	router.Handle(handler.Path, handler)

	// Start the HTTP server.
	http.ListenAndServe(":8080", router)
}
```

### Key Notes:
1. **Secret Management:**
   Use environment variables or a secure vault to store the `secret` instead of hardcoding it.

2. **Error Responses:**
   Respond with `403 Forbidden` if the HMAC validation fails. Adjust the `PrepareRequest` method to handle errors accordingly.

3. **Payload Replay Prevention:**
   Consider implementing replay protection by checking for unique IDs or timestamps in the payload.

4. **Testing:**
   Use a tool like `ngrok` to expose your local server to GitHub for testing webhook events.

With this setup, your Go program securely validates GitHub webhook payloads using HMAC.
*/
