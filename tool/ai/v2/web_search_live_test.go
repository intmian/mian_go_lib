//go:build live_ai

package ai

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

// TestLiveOpenAIWebSearch is an opt-in manual connectivity check. Set
// OPENAI_API_KEY before running it; OPENAI_BASE_URL and OPENAI_MODEL are optional.
// The API key is never printed by the test.
func TestLiveOpenAIWebSearch(t *testing.T) {
	token := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if token == "" {
		t.Skip("OPENAI_API_KEY is required for the live web search test")
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = "gpt-5.4"
	}
	reasoningEffort := ReasoningEffort(strings.TrimSpace(os.Getenv("OPENAI_REASONING_EFFORT")))
	if reasoningEffort == "" {
		reasoningEffort = ReasoningEffortMedium
	}
	provider, err := NewOpenAIProvider(
		strings.TrimSpace(os.Getenv("OPENAI_BASE_URL")),
		token,
		ProviderSourceTypeOpenAI,
	)
	if err != nil {
		t.Fatalf("create OpenAI provider: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	resp, err := provider.Chat(ctx, ChatRequest{
		Model: model,
		Messages: []ChatMessage{{
			Role:    ChatRoleUser,
			Content: "deepseek v4 正式版出了吗？",
		}},
		ReasoningEffort: reasoningEffort,
		Tools:           []ChatTool{ChatToolWebSearch},
	})
	if err != nil {
		t.Fatalf("OpenAI web search failed: %v", err)
	}
	if strings.TrimSpace(resp.Text) == "" {
		t.Fatal("OpenAI web search returned empty text")
	}
	if len(resp.Citations) == 0 {
		t.Fatalf("OpenAI web search returned no URL citations; response: %s", resp.Text)
	}

	t.Logf("model: %s", resp.Model)
	t.Logf("answer:\n%s", resp.Text)
	for i, citation := range resp.Citations {
		t.Logf("citation[%d]: %s - %s", i, citation.Title, citation.URL)
	}
}
