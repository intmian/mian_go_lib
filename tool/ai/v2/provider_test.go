package ai

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/openai/openai-go/v3/responses"
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
	if !slices.Contains(openAIProvider.AvailableTools(), string(ChatToolWebSearch)) {
		t.Fatal("OpenAI v2 chat should report web_search")
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

func TestOpenAIResponseParamsAddsWebSearch(t *testing.T) {
	params := openAIResponseParams(ChatRequest{
		Model:        "gpt-test",
		Messages:     []ChatMessage{{Role: ChatRoleUser, Content: "latest news"}},
		MaxToolCalls: 2,
		Tools:        []ChatTool{ChatToolWebSearch},
	})

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal response params: %v", err)
	}
	var payload struct {
		MaxToolCalls int64 `json:"max_tool_calls"`
		Tools        []struct {
			Type string `json:"type"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal response params: %v", err)
	}
	if payload.MaxToolCalls != 2 {
		t.Fatalf("unexpected max tool calls: %d", payload.MaxToolCalls)
	}
	if len(payload.Tools) != 1 || payload.Tools[0].Type != string(ChatToolWebSearch) {
		t.Fatalf("unexpected tools payload: %s", data)
	}
}

func TestOpenAIResponseParamsCompatibilityWebSearchKeepsReasoning(t *testing.T) {
	params := openAIResponseParams(ChatRequest{
		Model:           "gpt-test",
		Messages:        []ChatMessage{{Role: ChatRoleUser, Content: "latest news"}},
		ReasoningEffort: ReasoningEffortLow,
		Tools:           []ChatTool{ChatToolWebSearch},
	})

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal response params: %v", err)
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal response params: %v", err)
	}
	for _, field := range []string{"max_tool_calls", "max_output_tokens"} {
		if _, ok := payload[field]; ok {
			t.Fatalf("compatibility web search request should omit %q: %s", field, data)
		}
	}
	var reasoning struct {
		Effort string `json:"effort"`
	}
	if err := json.Unmarshal(payload["reasoning"], &reasoning); err != nil {
		t.Fatalf("unmarshal reasoning params: %v", err)
	}
	if reasoning.Effort != string(ReasoningEffortLow) {
		t.Fatalf("compatibility web search request should keep reasoning: %s", data)
	}
}

func TestExtractOpenAIURLCitations(t *testing.T) {
	resp := &responses.Response{
		Output: []responses.ResponseOutputItemUnion{{
			Type: "message",
			Content: []responses.ResponseOutputMessageContentUnion{{
				Type: "output_text",
				Annotations: []responses.ResponseOutputTextAnnotationUnion{{
					Type:       "url_citation",
					StartIndex: 3,
					EndIndex:   8,
					Title:      "Example",
					URL:        "https://example.com",
				}},
			}},
		}},
	}

	citations := extractOpenAIURLCitations(resp)
	if len(citations) != 1 {
		t.Fatalf("expected one citation, got %#v", citations)
	}
	got := citations[0]
	if got.StartIndex != 3 || got.EndIndex != 8 || got.Title != "Example" || got.URL != "https://example.com" {
		t.Fatalf("unexpected citation: %#v", got)
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

func TestValidateChatRequestRejectsWebSearchForDeepSeek(t *testing.T) {
	provider, err := NewOpenAIProvider("https://api.deepseek.com", "test-token", ProviderSourceTypeDeepSeek)
	if err != nil {
		t.Fatalf("NewOpenAIProvider error: %v", err)
	}

	err = provider.validateChatRequest(ChatRequest{
		Model:    "deepseek-chat",
		Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hello"}},
		Tools:    []ChatTool{ChatToolWebSearch},
	})
	if err == nil {
		t.Fatal("expected DeepSeek web_search validation error")
	}
}

func TestValidateChatRequestRejectsDuplicateTools(t *testing.T) {
	provider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeOpenAI)
	if err != nil {
		t.Fatalf("NewOpenAIProvider error: %v", err)
	}

	err = provider.validateChatRequest(ChatRequest{
		Model:    "gpt-test",
		Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hello"}},
		Tools:    []ChatTool{ChatToolWebSearch, ChatToolWebSearch},
	})
	if err == nil {
		t.Fatal("expected duplicate tool validation error")
	}
}

func TestValidateChatRequestRejectsMaxToolCallsWithoutTools(t *testing.T) {
	provider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeOpenAI)
	if err != nil {
		t.Fatalf("NewOpenAIProvider error: %v", err)
	}

	err = provider.validateChatRequest(ChatRequest{
		Model:        "gpt-test",
		Messages:     []ChatMessage{{Role: ChatRoleUser, Content: "hello"}},
		MaxToolCalls: 1,
	})
	if err == nil {
		t.Fatal("expected max tool calls validation error")
	}
}
