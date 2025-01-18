package dependencies

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func triggerGitHubWorkflow(owner, repo, workflowID, ref string, inputs map[string]interface{}, token string) error {
	// Your GitHub Personal Access Token (ensure it's kept secure)

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

func sampleGithub() {
	// Example usage
	owner := "your-github-username"
	repo := "your-repo-name"
	workflowID := "your_workflow.yml" // Can also be the workflow ID number
	ref := "main"                     // The branch or tag to run the workflow on

	// Optional inputs if your workflow requires them
	inputs := map[string]interface{}{
		"example_input": "Hello from go-github!",
	}

	// Your GitHub Personal Access Token (ensure it's kept secure)
	token := viper.GetString("GITHUB_PAT") // It's recommended to set this as an environment variable

	err := triggerGitHubWorkflow(owner, repo, workflowID, ref, inputs, token)
	if err != nil {
		fmt.Printf("Error triggering workflow: %v\n", err)
	} else {
		fmt.Println("Workflow triggered successfully!")
	}
}
