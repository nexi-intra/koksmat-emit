package emit

import (
	"bytes"
	"log"
	"testing"
)

func TestTriggerGitHubActions(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil)
	}()

	err := TriggerGitHubActions()

	if err != nil {
		t.Errorf("TriggerGitHubActions() = %v; want nil", err)
	}
}
