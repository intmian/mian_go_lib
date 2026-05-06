package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type BaseAgentSetting struct {
	ProviderID      ProviderID      `json:"providerID"`
	SysPrompt       string          `json:"sysPrompt"`
	Models          []string        `json:"models"`
	ReasoningEffort ReasoningEffort `json:"reasoningEffort"`
}

type BaseAgent struct {
	provider AgentProviderBinding
	setting  BaseAgentSettingState
	history  AgentMessageHistory
	mu       sync.Mutex
}

func NewBaseAgent() *BaseAgent {
	return &BaseAgent{
		provider: NewAgentProviderBinding(defaultRegistry),
		history:  NewAgentMessageHistory(),
	}
}

func (a *BaseAgent) GetID() AgentID {
	return AgentIDBase
}

func (a *BaseAgent) Init() error {
	if a == nil {
		return errors.New("agent is nil")
	}
	setting, ok := GetAgentSettingAs[*BaseAgentSetting](a.GetID())
	if !ok {
		return errors.New("base agent setting not registered")
	}
	return a.InitWithSetting(setting)
}

func (a *BaseAgent) InitWithSetting(setting *BaseAgentSetting) error {
	if a == nil {
		return errors.New("agent is nil")
	}
	if setting == nil {
		return errors.New("base agent setting is nil")
	}
	if setting.ProviderID == 0 {
		return errors.New("provider id is required")
	}
	if len(setting.Models) == 0 {
		return errors.New("agent models are required")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.setting.InitWithSetting(setting); err != nil {
		return err
	}
	return a.provider.InitProvider(setting.ProviderID)
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

	setting, err := a.setting.settingSnapshot()
	if err != nil {
		return "", err
	}
	provider, err := a.provider.GetProvider()
	if err != nil {
		return "", err
	}
	messages := a.history.BuildMessages(setting.SysPrompt, content)

	var errs []error
	for _, model := range setting.Models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		resp, err := provider.Chat(ctx, ChatRequest{
			Model:           model,
			Messages:        messages,
			ReasoningEffort: setting.ReasoningEffort,
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
		a.history.AppendExchange(content, text)
		return text, nil
	}

	if len(errs) > 0 {
		return "", errors.Join(errs...)
	}
	return "", errors.New("agent models are required")
}

func (a *BaseAgent) History() []ChatMessage {
	if a == nil {
		return nil
	}
	return a.history.History()
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
