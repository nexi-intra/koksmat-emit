package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookGitHub(t *testing.T) {
	tests := []struct {
		name     string
		input    GitHubWebhookInput
		expected GitHubWebhookOutput
	}{
		{
			name: "created action",
			input: GitHubWebhookInput{
				Action: "created",
				Repository: struct {
					Name  string `json:"name"`
					Owner struct {
						Login string `json:"login"`
					} `json:"owner"`
				}{
					Name: "test-repo",
				},
			},
			expected: GitHubWebhookOutput{
				Message: "A new issue or PR was opened in test-repo",
				Status:  "success",
			},
		},
		{
			name: "opened action",
			input: GitHubWebhookInput{
				Action: "opened",
				Repository: struct {
					Name  string `json:"name"`
					Owner struct {
						Login string `json:"login"`
					} `json:"owner"`
				}{
					Name: "test-repo",
				},
			},
			expected: GitHubWebhookOutput{
				Message: "A new issue or PR was opened in test-repo",
				Status:  "success",
			},
		},
		{
			name: "closed action",
			input: GitHubWebhookInput{
				Action: "closed",
				Repository: struct {
					Name  string `json:"name"`
					Owner struct {
						Login string `json:"login"`
					} `json:"owner"`
				}{
					Name: "test-repo",
				},
			},
			expected: GitHubWebhookOutput{
				Message: "An issue or PR was closed in test-repo",
				Status:  "success",
			},
		},
		{
			name: "unhandled action",
			input: GitHubWebhookInput{
				Action: "unhandled",
				Repository: struct {
					Name  string `json:"name"`
					Owner struct {
						Login string `json:"login"`
					} `json:"owner"`
				}{
					Name: "test-repo",
				},
			},
			expected: GitHubWebhookOutput{
				Message: "Action not handled: unhandled",
				Status:  "ignored",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interactor := webhook_GitHub()
			var output GitHubWebhookOutput
			err := interactor.Interact(context.Background(), tt.input, &output)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, output)
		})
	}
}
