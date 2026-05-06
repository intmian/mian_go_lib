package ai

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

type AgentProviderBinding struct {
	resolver   providerResolver
	providerID ProviderID
	mu         sync.RWMutex
}

func NewAgentProviderBinding(resolver providerResolver) AgentProviderBinding {
	return AgentProviderBinding{resolver: resolver}
}

func (c *AgentProviderBinding) InitProvider(providerID ProviderID) error {
	if providerID == 0 {
		return errors.New("provider id is required")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.providerID = providerID
	return nil
}

func (c *AgentProviderBinding) GetProvider() (IProvider, error) {
	if c == nil {
		return nil, errors.New("agent provider binding is nil")
	}
	c.mu.RLock()
	providerID := c.providerID
	resolver := c.resolver
	c.mu.RUnlock()
	if providerID == 0 {
		return nil, errors.New("agent provider is not initialized")
	}
	if resolver == nil {
		return nil, errors.New("provider resolver is nil")
	}
	provider, ok := resolver.GetProvider(providerID)
	if !ok {
		return nil, errors.New("provider not registered")
	}
	return provider, nil
}

type AgentSettingState[S IAgentSetting] struct {
	setting     S
	initialized bool
	mu          sync.RWMutex
}

type BaseAgentSettingState = AgentSettingState[*BaseAgentSetting]

func (c *AgentSettingState[S]) Init(agent AgentIDProvider) error {
	if agent == nil {
		return errors.New("agent is nil")
	}
	setting, ok := GetAgentSettingAs[S](agent.GetID())
	if !ok {
		return errors.New("agent setting not registered")
	}
	return c.InitWithSetting(setting)
}

func (c *AgentSettingState[S]) InitWithSetting(setting S) error {
	if isNilSetting(setting) {
		return errors.New("agent setting is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setting = setting
	c.initialized = true
	return nil
}

func (c *AgentSettingState[S]) GetSetting() S {
	var empty S
	if c == nil {
		return empty
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.initialized {
		return empty
	}
	return c.setting
}

func (c *AgentSettingState[S]) settingSnapshot() (S, error) {
	var empty S
	if c == nil {
		return empty, errors.New("agent setting state is nil")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.initialized {
		return empty, errors.New("agent setting is not initialized")
	}
	return c.setting, nil
}

type AgentMessageHistory struct {
	history []ChatMessage
	mu      sync.RWMutex
}

func NewAgentMessageHistory() AgentMessageHistory {
	return AgentMessageHistory{history: make([]ChatMessage, 0)}
}

func (c *AgentMessageHistory) AppendExchange(userContent, assistantContent string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.history = append(c.history,
		ChatMessage{Role: ChatRoleUser, Content: userContent},
		ChatMessage{Role: ChatRoleAssistant, Content: assistantContent},
	)
}

func (c *AgentMessageHistory) History() []ChatMessage {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	history := make([]ChatMessage, len(c.history))
	copy(history, c.history)
	return history
}

func (c *AgentMessageHistory) BuildMessages(sysPrompt, content string) []ChatMessage {
	history := c.History()
	messages := make([]ChatMessage, 0, len(history)+2)
	if strings.TrimSpace(sysPrompt) != "" {
		messages = append(messages, ChatMessage{
			Role:    ChatRoleSystem,
			Content: sysPrompt,
		})
	}
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{
		Role:    ChatRoleUser,
		Content: content,
	})
	return messages
}

func isNilSetting[S IAgentSetting](setting S) bool {
	value := reflect.ValueOf(setting)
	if !value.IsValid() {
		return true
	}
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
