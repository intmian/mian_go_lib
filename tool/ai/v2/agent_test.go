package v2

import (
	"context"
	"errors"
	"io"
	"testing"
)

type fakeProvider struct {
	calls       []ChatRequest
	failByModel map[string]error
	textByModel map[string]string
}

func (p *fakeProvider) AvailableModels(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (p *fakeProvider) AvailableReasoning() []ReasoningEffort {
	return []ReasoningEffort{ReasoningEffortNone, ReasoningEffortMedium}
}

func (p *fakeProvider) AvailableTools() []string {
	return nil
}

func (p *fakeProvider) SourceType() ProviderSourceType {
	return ProviderSourceTypeOpenAI
}

func (p *fakeProvider) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	p.calls = append(p.calls, req)
	if err := p.failByModel[req.Model]; err != nil {
		return ChatResponse{}, err
	}
	if text := p.textByModel[req.Model]; text != "" {
		return ChatResponse{Text: text, Model: req.Model}, nil
	}
	return ChatResponse{Text: "ok", Model: req.Model}, nil
}

func (p *fakeProvider) ChatStream(ctx context.Context, req ChatRequest) (IChatStream, error) {
	return &fakeChatStream{}, nil
}

type fakeChatStream struct{}

func (s *fakeChatStream) Recv() (ChatStreamEvent, error) {
	return ChatStreamEvent{}, io.EOF
}

func (s *fakeChatStream) Close() error {
	return nil
}

func TestBaseAgentFallsBackModelsAndRecordsHistory(t *testing.T) {
	provider := &fakeProvider{
		failByModel: map[string]error{"m1": errors.New("boom")},
		textByModel: map[string]string{"m2": "answer"},
	}
	registry := NewAgentRegistry()
	if err := registry.AddProvider("main", provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}

	mother, err := NewBaseAgentMother(BaseAgentSetting{
		ProviderID:      "main",
		SysPrompt:       "system prompt",
		Models:          []string{"m1", "m2"},
		ReasoningEffort: ReasoningEffortMedium,
	})
	if err != nil {
		t.Fatalf("NewBaseAgentMother error: %v", err)
	}
	if err := registry.AddAgentMother(AgentIDBase, mother); err != nil {
		t.Fatalf("AddAgentMother error: %v", err)
	}

	motherAny, ok := registry.GetAgentMother(AgentIDBase)
	if !ok {
		t.Fatal("agent mother not registered")
	}
	agentAny, err := motherAny.CreateAnyAgent(CoreAgentSetting{ProviderResolver: registry})
	if err != nil {
		t.Fatalf("CreateAgent error: %v", err)
	}
	agent := agentAny.(*BaseAgent)
	text, err := agent.Chat(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Chat error: %v", err)
	}
	if text != "answer" {
		t.Fatalf("unexpected chat text: %q", text)
	}
	if len(provider.calls) != 2 {
		t.Fatalf("expected 2 model calls, got %d", len(provider.calls))
	}
	if provider.calls[0].Model != "m1" || provider.calls[1].Model != "m2" {
		t.Fatalf("unexpected model fallback order: %#v", provider.calls)
	}
	if len(agent.History()) != 2 {
		t.Fatalf("expected user and assistant history, got %d", len(agent.History()))
	}

	_, err = agent.Chat(context.Background(), "next")
	if err != nil {
		t.Fatalf("second Chat error: %v", err)
	}
	last := provider.calls[len(provider.calls)-1]
	if len(last.Messages) != 4 {
		t.Fatalf("expected system + previous exchange + current user, got %d", len(last.Messages))
	}
	if last.Messages[1].Role != ChatRoleUser || last.Messages[1].Content != "hello" {
		t.Fatalf("previous user message missing from history: %#v", last.Messages)
	}
	if last.Messages[2].Role != ChatRoleAssistant || last.Messages[2].Content != "answer" {
		t.Fatalf("previous assistant message missing from history: %#v", last.Messages)
	}
}

func TestBaseAgentDoesNotRecordFailedChat(t *testing.T) {
	provider := &fakeProvider{
		failByModel: map[string]error{"m1": errors.New("boom")},
	}
	agent, err := NewBaseAgent(CoreAgentSetting{
		ProviderResolver: providerResolverFunc(func(id ProviderID) (IProvider, bool) {
			return provider, true
		}),
	}, BaseAgentSetting{
		ProviderID: "main",
		Models:     []string{"m1"},
	})
	if err != nil {
		t.Fatalf("NewBaseAgent error: %v", err)
	}

	if _, err := agent.Chat(context.Background(), "hello"); err == nil {
		t.Fatal("expected chat error")
	}
	if len(agent.History()) != 0 {
		t.Fatalf("failed chat should not record history: %#v", agent.History())
	}
}

func TestBaseAgentResolvesProviderAtChatTime(t *testing.T) {
	registry := NewAgentRegistry()
	agent, err := NewBaseAgent(CoreAgentSetting{ProviderResolver: registry}, BaseAgentSetting{
		ProviderID: "late",
		Models:     []string{"m1"},
	})
	if err != nil {
		t.Fatalf("NewBaseAgent error: %v", err)
	}
	if _, err := agent.Chat(context.Background(), "hello"); err == nil {
		t.Fatal("expected missing provider error before registration")
	}

	if err := registry.AddProvider("late", &fakeProvider{}); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}
	text, err := agent.Chat(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Chat after provider registration error: %v", err)
	}
	if text != "ok" {
		t.Fatalf("unexpected chat text: %q", text)
	}
}

type providerResolverFunc func(id ProviderID) (IProvider, bool)

func (f providerResolverFunc) GetProvider(id ProviderID) (IProvider, bool) {
	return f(id)
}
