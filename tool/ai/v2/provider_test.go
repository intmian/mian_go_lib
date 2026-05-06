package ai

import (
	"slices"
	"testing"
)

func TestIsOpenAILikeSource(t *testing.T) {
	if !IsOpenAILikeSource(ProviderSourceTypeOpenAI) {
		t.Fatal("OpenAI should be openai-like")
	}
	if !IsOpenAILikeSource(ProviderSourceTypeDeepSeek) {
		t.Fatal("DeepSeek should be openai-like")
	}
	if IsOpenAILikeSource(ProviderSourceType("unknown")) {
		t.Fatal("unknown provider should not be openai-like")
	}
}

func TestOpenAIProviderCapabilities(t *testing.T) {
	openAIProvider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeOpenAI)
	if err != nil {
		t.Fatalf("NewOpenAIProvider OpenAI error: %v", err)
	}
	if !slices.Contains(openAIProvider.AvailableReasoning(), ReasoningEffortXHigh) {
		t.Fatal("OpenAI reasoning should include xhigh")
	}
	if len(openAIProvider.AvailableTools()) != 0 {
		t.Fatal("OpenAI v2 chat should not report tools before tool request support exists")
	}

	deepSeekProvider, err := NewOpenAIProvider("https://api.deepseek.com", "test-token", ProviderSourceTypeDeepSeek)
	if err != nil {
		t.Fatalf("NewOpenAIProvider DeepSeek error: %v", err)
	}
	if !slices.Contains(deepSeekProvider.AvailableReasoning(), ReasoningEffortMax) {
		t.Fatal("DeepSeek reasoning should include max")
	}
	if len(deepSeekProvider.AvailableTools()) != 0 {
		t.Fatal("DeepSeek v2 chat should not report tools before tool request support exists")
	}
}

func TestNewOpenAIProviderRejectsNonOpenAILikeSource(t *testing.T) {
	if _, err := NewOpenAIProvider("", "test-token", ProviderSourceType("unknown")); err == nil {
		t.Fatal("expected unsupported provider error")
	}
}

func TestValidateChatRequestRejectsInvalidRole(t *testing.T) {
	provider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeOpenAI)
	if err != nil {
		t.Fatalf("NewOpenAIProvider error: %v", err)
	}

	err = provider.validateChatRequest(ChatRequest{
		Model: "gpt-test",
		Messages: []ChatMessage{
			{Role: ChatRole("bad"), Content: "hello"},
		},
	})
	if err == nil {
		t.Fatal("expected invalid role error")
	}
}

func TestValidateChatRequestRejectsUnavailableThinking(t *testing.T) {
	provider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeOpenAI)
	if err != nil {
		t.Fatalf("NewOpenAIProvider error: %v", err)
	}

	err = provider.validateChatRequest(ChatRequest{
		Model:    "gpt-test",
		Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hello"}},
		Thinking: ThinkingTypeEnabled,
	})
	if err == nil {
		t.Fatal("expected OpenAI thinking error")
	}
}
