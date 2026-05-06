package ai

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
	resetDefaultRegistryForTest()
	provider := &fakeProvider{
		failByModel: map[string]error{"m1": errors.New("boom")},
		textByModel: map[string]string{"m2": "answer"},
	}
	if err := AddProvider(1, provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}

	setting := &BaseAgentSetting{
		ProviderID:      1,
		SysPrompt:       "system prompt",
		Models:          []string{"m1", "m2"},
		ReasoningEffort: ReasoningEffortMedium,
	}

	agent := NewBaseAgent()
	if err := agent.InitWithSetting(setting); err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}
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
	resetDefaultRegistryForTest()
	provider := &fakeProvider{
		failByModel: map[string]error{"m1": errors.New("boom")},
	}
	if err := AddProvider(1, provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}
	agent := NewBaseAgent()
	err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 1,
		Models:     []string{"m1"},
	})
	if err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}

	if _, err := agent.Chat(context.Background(), "hello"); err == nil {
		t.Fatal("expected chat error")
	}
	if len(agent.History()) != 0 {
		t.Fatalf("failed chat should not record history: %#v", agent.History())
	}
}

func TestBaseAgentRejectsBlankModels(t *testing.T) {
	agent := NewBaseAgent()
	err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 1,
		Models:     []string{" ", "\t"},
	})
	if err == nil {
		t.Fatal("expected blank models error")
	}
}

func TestBaseAgentResolvesProviderAtChatTime(t *testing.T) {
	resetDefaultRegistryForTest()
	agent := NewBaseAgent()
	err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 3,
		Models:     []string{"m1"},
	})
	if err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}
	if _, err := agent.Chat(context.Background(), "hello"); err == nil {
		t.Fatal("expected missing provider error before registration")
	}

	if err := AddProvider(3, &fakeProvider{}); err != nil {
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

func TestBaseAgentInitReadsRegisteredSetting(t *testing.T) {
	resetDefaultRegistryForTest()
	if err := AddProvider(1, &fakeProvider{}); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}
	if err := AddAgentSetting(AgentIDBase, &BaseAgentSetting{
		ProviderID: 1,
		Models:     []string{"m1"},
	}); err != nil {
		t.Fatalf("AddAgentSetting error: %v", err)
	}

	agent := NewBaseAgent()
	if err := agent.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	text, err := agent.Chat(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Chat error: %v", err)
	}
	if text != "ok" {
		t.Fatalf("unexpected chat text: %q", text)
	}
}

func TestSettingStateGetSettingReturnsStoredReference(t *testing.T) {
	var state BaseAgentSettingState
	if err := state.InitWithSetting(&BaseAgentSetting{
		ProviderID: 1,
		Models:     []string{"m1"},
	}); err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}

	got := state.GetSetting()
	got.ProviderID = 2
	got.Models[0] = "m2"
	gotAgain := state.GetSetting()
	if gotAgain.ProviderID != 2 || gotAgain.Models[0] != "m2" {
		t.Fatalf("GetSetting should return stored setting reference, got %#v", gotAgain)
	}
}

func TestAgentSettingStateInitReadsByAgentID(t *testing.T) {
	resetDefaultRegistryForTest()
	if err := AddAgentSetting(AgentIDBase, &BaseAgentSetting{
		ProviderID: 1,
		Models:     []string{"m1"},
	}); err != nil {
		t.Fatalf("AddAgentSetting error: %v", err)
	}

	agent := NewBaseAgent()
	var state AgentSettingState[*BaseAgentSetting]
	if err := state.Init(agent); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	setting := state.GetSetting()
	if setting.ProviderID != 1 || setting.Models[0] != "m1" {
		t.Fatalf("unexpected setting: %#v", setting)
	}
}

func TestChatOnceDoesNotRecordHistory(t *testing.T) {
	resetDefaultRegistryForTest()
	if err := AddProvider(1, &fakeProvider{}); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}

	agent := NewBaseAgent()
	if err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 1,
		SysPrompt:  "translate",
		Models:     []string{"m1"},
	}); err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}

	text, err := agent.ChatOnce(context.Background(), "hello")
	if err != nil {
		t.Fatalf("ChatOnce error: %v", err)
	}
	if text != "ok" {
		t.Fatalf("unexpected text: %q", text)
	}
	if len(agent.History()) != 0 {
		t.Fatalf("ChatOnce should not record history, got %d messages", len(agent.History()))
	}

	// Chat should not see ChatOnce history.
	text, err = agent.Chat(context.Background(), "world")
	if err != nil {
		t.Fatalf("Chat error: %v", err)
	}
	if len(agent.History()) != 2 {
		t.Fatalf("Chat should only record its own exchange, got %d", len(agent.History()))
	}
}

func TestChatOnceIncludesSysPrompt(t *testing.T) {
	resetDefaultRegistryForTest()
	provider := &fakeProvider{}
	if err := AddProvider(1, provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}

	agent := NewBaseAgent()
	if err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 1,
		SysPrompt:  "translate to english",
		Models:     []string{"m1"},
	}); err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}

	if _, err := agent.ChatOnce(context.Background(), "hello"); err != nil {
		t.Fatalf("ChatOnce error: %v", err)
	}

	req := provider.calls[len(provider.calls)-1]
	if len(req.Messages) != 2 {
		t.Fatalf("expected sys + user, got %d messages", len(req.Messages))
	}
	if req.Messages[0].Role != ChatRoleSystem || req.Messages[0].Content != "translate to english" {
		t.Fatalf("sys prompt missing: %#v", req.Messages[0])
	}
	if req.Messages[1].Role != ChatRoleUser || req.Messages[1].Content != "hello" {
		t.Fatalf("user message wrong: %#v", req.Messages[1])
	}
}

func TestChatOnceFailsWithoutLeakingToHistory(t *testing.T) {
	resetDefaultRegistryForTest()
	provider := &fakeProvider{
		failByModel: map[string]error{"m1": errors.New("boom")},
	}
	if err := AddProvider(1, provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}

	agent := NewBaseAgent()
	if err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 1,
		Models:     []string{"m1"},
	}); err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}

	if _, err := agent.ChatOnce(context.Background(), "hello"); err == nil {
		t.Fatal("expected ChatOnce error")
	}
	if len(agent.History()) != 0 {
		t.Fatal("failed ChatOnce should not leak to history")
	}
}

func TestClearHistory(t *testing.T) {
	resetDefaultRegistryForTest()
	if err := AddProvider(1, &fakeProvider{}); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}

	agent := NewBaseAgent()
	if err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 1,
		Models:     []string{"m1"},
	}); err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}

	if _, err := agent.Chat(context.Background(), "first"); err != nil {
		t.Fatalf("Chat error: %v", err)
	}
	if len(agent.History()) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(agent.History()))
	}

	agent.ClearHistory()
	if len(agent.History()) != 0 {
		t.Fatalf("ClearHistory should empty history, got %d", len(agent.History()))
	}

	// Next Chat should start fresh (no history).
	if _, err := agent.Chat(context.Background(), "second"); err != nil {
		t.Fatalf("Chat error: %v", err)
	}
	if len(agent.History()) != 2 {
		t.Fatalf("after clear, chat should record its own exchange only, got %d", len(agent.History()))
	}
}

func TestHistoryClearOnNil(t *testing.T) {
	var h *AgentMessageHistory
	h.Clear() // must not panic
}

func TestClearHistoryOnNilAgent(t *testing.T) {
	var a *BaseAgent
	a.ClearHistory() // must not panic
}

func TestChatOnceOnNilAgent(t *testing.T) {
	var a *BaseAgent
	if _, err := a.ChatOnce(context.Background(), "hello"); err == nil {
		t.Fatal("expected error")
	}
}

var _ IAgent[*BaseAgentSetting] = (*BaseAgent)(nil)
