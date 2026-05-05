package v2

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type IAgent interface {
	Chat(ctx context.Context, content string) (string, error)
}

type BaseAgentSetting struct {
	ProviderID      ProviderID      `json:"providerID"`
	SysPrompt       string          `json:"sysPrompt"`
	Models          []string        `json:"models"`
	ReasoningEffort ReasoningEffort `json:"reasoningEffort"`
}

type BaseAgent struct {
	core    CoreAgentSetting
	setting BaseAgentSetting
	history []ChatMessage
	mu      sync.Mutex
}

func NewBaseAgentMother(setting BaseAgentSetting) (*AgentMother[*BaseAgent, BaseAgentSetting], error) {
	return NewAgentMother[*BaseAgent, BaseAgentSetting](setting, NewBaseAgent)
}

func NewBaseAgent(core CoreAgentSetting, setting BaseAgentSetting) (*BaseAgent, error) {
	if setting.ProviderID == "" {
		return nil, errors.New("provider id is required")
	}
	if len(setting.Models) == 0 {
		return nil, errors.New("agent models are required")
	}
	return &BaseAgent{
		core:    core,
		setting: setting,
		history: make([]ChatMessage, 0),
	}, nil
}

func (a *BaseAgent) Chat(ctx context.Context, content string) (string, error) {
	if a == nil {
		return "", errors.New("agent is nil")
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return "", errors.New("chat content is required")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	provider, ok := a.core.GetProvider(a.setting.ProviderID)
	if !ok {
		return "", errors.New("provider not registered")
	}
	messages := a.buildMessages(content)

	var errs []error
	for _, model := range a.setting.Models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		resp, err := provider.Chat(ctx, ChatRequest{
			Model:           model,
			Messages:        messages,
			ReasoningEffort: a.setting.ReasoningEffort,
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", model, err))
			continue
		}
		text := strings.TrimSpace(resp.Text)
		if text == "" {
			errs = append(errs, fmt.Errorf("%s: empty response", model))
			continue
		}
		a.history = append(a.history,
			ChatMessage{Role: ChatRoleUser, Content: content},
			ChatMessage{Role: ChatRoleAssistant, Content: text},
		)
		return text, nil
	}

	if len(errs) > 0 {
		return "", errors.Join(errs...)
	}
	return "", errors.New("agent models are required")
}

func (a *BaseAgent) buildMessages(content string) []ChatMessage {
	messages := make([]ChatMessage, 0, len(a.history)+2)
	if strings.TrimSpace(a.setting.SysPrompt) != "" {
		messages = append(messages, ChatMessage{
			Role:    ChatRoleSystem,
			Content: a.setting.SysPrompt,
		})
	}
	messages = append(messages, a.history...)
	messages = append(messages, ChatMessage{
		Role:    ChatRoleUser,
		Content: content,
	})
	return messages
}

func (a *BaseAgent) History() []ChatMessage {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	history := make([]ChatMessage, len(a.history))
	copy(history, a.history)
	return history
}

func (s BaseAgentSetting) ExportJSON() ([]byte, error) {
	return ExportSettingJSON(s)
}

func (s *BaseAgentSetting) ImportJSON(data []byte) error {
	if s == nil {
		return errors.New("base agent setting is nil")
	}
	next, err := ImportSettingJSON(*s, data)
	if err != nil {
		return err
	}
	*s = next
	return nil
}

func (s BaseAgentSetting) ExportJSONDoc() ([]SettingFieldDoc, error) {
	return ExportSettingDoc(s)
}
