package ai

import (
	"testing"

	copilot "github.com/github/copilot-sdk/go"
)

func TestAi(t *testing.T) {
	client := copilot.NewClient(nil)
	if err := client.Start(); err != nil {
		t.Fatal(err)
	}
	defer client.Stop()

	t.Log("Copilot client started successfully")
	session, err := client.CreateSession(&copilot.SessionConfig{Model: "gpt-4.1"})
	if err != nil {
		t.Fatal(err)
	}
	response, err := session.SendAndWait(copilot.MessageOptions{Prompt: "What is 2 + 2?"}, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Response: %s", *response.Data.Content)
}
