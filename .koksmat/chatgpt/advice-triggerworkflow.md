---
title: Trigger GitHub Workflow
chatgpt-prompt: like a go function that can trigger a github workflow
---

> like a go function that can trigger a github workflow

---

Certainly! You can trigger a GitHub Actions workflow programmatically using Go by leveraging GitHub's [Workflow Dispatch API](https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28#create-a-workflow-dispatch-event). Below is a comprehensive example of how to achieve this using Go's `net/http` package. Additionally, I'll show you how to do it using the [go-github](https://github.com/google/go-github) library for a more streamlined approach.

---

## Prerequisites

1. **Personal Access Token (PAT):**

   - Ensure you have a GitHub PAT with the `repo` and `workflow` scopes. You can create one [here](https://github.com/settings/tokens).

2. **Workflow Configuration:**

   - Your GitHub Actions workflow YAML file should be configured to allow manual triggers. Ensure it includes the `workflow_dispatch` trigger.

   ```yaml
   # .github/workflows/your_workflow.yml
   name: Your Workflow

   on:
     workflow_dispatch:
       inputs:
         # Define any inputs here if needed
         example_input:
           description: "An example input"
           required: false
           default: "default value"

   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         - name: Checkout repository
           uses: actions/checkout@v2
         # Add more steps as needed
   ```

---

## Method 1: Using `net/http`

Here's how you can trigger a GitHub Actions workflow using Go's standard `net/http` package.

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// WorkflowDispatchPayload represents the JSON payload for the workflow dispatch event
type WorkflowDispatchPayload struct {
	Ref    string                 `json:"ref"`              // The git reference for the workflow (e.g., "main")
	Inputs map[string]interface{} `json:"inputs,omitempty"` // Optional inputs for the workflow
}

func triggerGitHubWorkflow(owner, repo, workflowID, ref string, inputs map[string]interface{}, token string) error {
	// Construct the API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches", owner, repo, workflowID)

	// Create the payload
	payload := WorkflowDispatchPayload{
		Ref:    ref,
		Inputs: inputs,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set necessary headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to trigger workflow, status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Example usage
	owner := "your-github-username"
	repo := "your-repo-name"
	workflowID := "your_workflow.yml" // Can also be the workflow ID number
	ref := "main"                      // The branch or tag to run the workflow on

	// Optional inputs if your workflow requires them
	inputs := map[string]interface{}{
		"example_input": "Hello from Go!",
	}

	// Your GitHub Personal Access Token (ensure it's kept secure)
	token := os.Getenv("GITHUB_TOKEN") // It's recommended to set this as an environment variable

	err := triggerGitHubWorkflow(owner, repo, workflowID, ref, inputs, token)
	if err != nil {
		fmt.Printf("Error triggering workflow: %v\n", err)
	} else {
		fmt.Println("Workflow triggered successfully!")
	}
}
```

### Explanation:

1. **WorkflowDispatchPayload Struct:**

   - Represents the JSON payload required by the GitHub API. It includes the `ref` (branch/tag) and optional `inputs`.

2. **triggerGitHubWorkflow Function:**

   - Constructs the API endpoint URL based on the repository details.
   - Marshals the payload into JSON.
   - Sets the necessary HTTP headers, including the `Authorization` header with the PAT.
   - Sends a POST request to trigger the workflow.
   - Checks the response status code to ensure the request was successful (`204 No Content`).

3. **main Function:**
   - Demonstrates how to use the `triggerGitHubWorkflow` function.
   - Retrieves the GitHub PAT from an environment variable for security.
   - Handles success and error messages accordingly.

---

## Method 2: Using the `go-github` Library

The [`go-github`](https://github.com/google/go-github) library by Google provides a more convenient way to interact with GitHub's API. Here's how you can use it to trigger a workflow.

### Install the Library

First, install the `go-github` library and its dependencies:

```bash
go get github.com/google/go-github/v50/github
go get golang.org/x/oauth2
```

### Example Code

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func triggerGitHubWorkflow(owner, repo, workflowID, ref string, inputs map[string]interface{}, token string) error {
	// Create an OAuth2 authenticated client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Prepare the dispatch request
	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    ref,
		Inputs: inputs,
	}

	// Trigger the workflow
	_, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, workflowID, event)
	if err != nil {
		return fmt.Errorf("failed to trigger workflow: %v", err)
	}

	return nil
}

func main() {
	// Example usage
	owner := "your-github-username"
	repo := "your-repo-name"
	workflowID := "your_workflow.yml" // Can also be the workflow ID number
	ref := "main"                      // The branch or tag to run the workflow on

	// Optional inputs if your workflow requires them
	inputs := map[string]interface{}{
		"example_input": "Hello from go-github!",
	}

	// Your GitHub Personal Access Token (ensure it's kept secure)
	token := os.Getenv("GITHUB_TOKEN") // It's recommended to set this as an environment variable

	err := triggerGitHubWorkflow(owner, repo, workflowID, ref, inputs, token)
	if err != nil {
		fmt.Printf("Error triggering workflow: %v\n", err)
	} else {
		fmt.Println("Workflow triggered successfully!")
	}
}
```

### Explanation:

1. **OAuth2 Authentication:**

   - Uses the `oauth2` package to authenticate requests with your PAT.

2. **go-github Client:**

   - Creates a new GitHub client using the authenticated HTTP client.

3. **CreateWorkflowDispatchEventRequest:**

   - Prepares the payload with the `ref` and optional `inputs`.

4. **Triggering the Workflow:**

   - Uses `CreateWorkflowDispatchEventByFileName` to trigger the workflow by its filename (e.g., `your_workflow.yml`). Alternatively, you can use `CreateWorkflowDispatchEvent` with the workflow ID.

5. **Handling Responses:**
   - The `go-github` library handles the HTTP response. If there's an error, it will be returned; otherwise, the workflow is successfully triggered.

---

## Security Considerations

- **Protect Your PAT:**

  - Never hard-code your GitHub PAT in your source code. Use environment variables or secure secret management systems.

- **Least Privilege:**

  - Ensure your PAT has only the necessary scopes (`repo` and `workflow` for triggering workflows).

- **Error Handling:**
  - Implement robust error handling in production code to handle various failure scenarios gracefully.

---

## Conclusion

Triggering GitHub Actions workflows programmatically using Go can be efficiently achieved either by making direct HTTP requests with `net/http` or by leveraging the `go-github` library for a more abstracted approach. Choose the method that best fits your project's needs and complexity.

If you have any further questions or need assistance with specific aspects of the implementation, feel free to ask!
