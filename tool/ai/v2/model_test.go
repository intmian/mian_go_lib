package ai

import (
	"testing"
)

func TestModelConfigValidatesCallProtocolAgainstType(t *testing.T) {
	tests := []struct {
		name      string
		modelType ModelType
		protocol  ModelCallProtocol
		wantErr   bool
	}{
		{name: "text inherits", modelType: ModelTypeText},
		{name: "openai text", modelType: ModelTypeText, protocol: ModelCallProtocolOpenAIText},
		{name: "text rejects stt", modelType: ModelTypeText, protocol: ModelCallProtocolOpenAISTT, wantErr: true},
		{name: "stt inherits", modelType: ModelTypeSTT},
		{name: "openai stt", modelType: ModelTypeSTT, protocol: ModelCallProtocolOpenAISTT},
		{name: "dashscope qwen3 asr", modelType: ModelTypeSTT, protocol: ModelCallProtocolDashScopeQwen3ASR},
		{name: "dashscope fun asr", modelType: ModelTypeSTT, protocol: ModelCallProtocolDashScopeFunASR},
		{name: "stt rejects text", modelType: ModelTypeSTT, protocol: ModelCallProtocolOpenAIText, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := (ModelConfig{ID: "model", Name: "upstream", Type: tt.modelType, CallProtocol: tt.protocol}).Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpenAIProviderRegistersModelAndUsesUpstreamName(t *testing.T) {
	provider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeOpenAI)
	if err != nil {
		t.Fatalf("NewOpenAIProvider: %v", err)
	}
	if err := provider.RegisterModel(ModelConfig{
		ID:        "main",
		Name:      "upstream-model",
		Type:      ModelTypeText,
		Reasoning: []ReasoningEffort{ReasoningEffortLow},
		Tools:     []ChatTool{ChatToolWebSearch},
	}); err != nil {
		t.Fatalf("RegisterModel: %v", err)
	}

	req, err := provider.resolveRegisteredModel(ChatRequest{
		Model:           "main",
		Messages:        []ChatMessage{{Role: ChatRoleUser, Content: "hello"}},
		ReasoningEffort: ReasoningEffortLow,
		Tools:           []ChatTool{ChatToolWebSearch},
	})
	if err != nil {
		t.Fatalf("resolveRegisteredModel: %v", err)
	}
	if req.Model != "upstream-model" {
		t.Fatalf("resolved model = %q", req.Model)
	}
}

func TestOpenAIProviderRejectsModelCapabilityMismatch(t *testing.T) {
	provider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeDeepSeek)
	if err != nil {
		t.Fatalf("NewOpenAIProvider: %v", err)
	}
	if err := provider.RegisterModel(ModelConfig{
		ID:    "deepseek",
		Name:  "deepseek-chat",
		Type:  ModelTypeText,
		Tools: []ChatTool{ChatToolWebSearch},
	}); err == nil {
		t.Fatal("expected unsupported tool error")
	}
}

func TestRegisteredModelCapabilitiesAreEnforced(t *testing.T) {
	provider, err := NewOpenAIProvider("", "test-token", ProviderSourceTypeOpenAI)
	if err != nil {
		t.Fatalf("NewOpenAIProvider: %v", err)
	}
	if err := provider.RegisterModel(ModelConfig{
		ID:        "limited",
		Name:      "upstream-model",
		Type:      ModelTypeText,
		Reasoning: []ReasoningEffort{ReasoningEffortLow},
	}); err != nil {
		t.Fatalf("RegisterModel: %v", err)
	}
	if _, err := provider.resolveRegisteredModel(ChatRequest{
		Model:           "limited",
		ReasoningEffort: ReasoningEffortHigh,
	}); err == nil {
		t.Fatal("expected unsupported reasoning error")
	}
	if _, err := provider.resolveRegisteredModel(ChatRequest{Model: "missing"}); err == nil {
		t.Fatal("expected missing registered model error")
	}
}
