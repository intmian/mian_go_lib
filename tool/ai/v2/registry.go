package ai

import (
	"errors"
	"sync"
)

type providerResolver interface {
	GetProvider(id ProviderID) (IProvider, bool)
}

type agentRegistry struct {
	mu              sync.RWMutex
	id2AgentSetting map[AgentID]IAgentSetting
	id2Provider     map[ProviderID]IProvider
}

func newAgentRegistry() *agentRegistry {
	return &agentRegistry{
		id2AgentSetting: make(map[AgentID]IAgentSetting),
		id2Provider:     make(map[ProviderID]IProvider),
	}
}

func (r *agentRegistry) AddProvider(id ProviderID, provider IProvider) error {
	if r == nil {
		return errors.New("agent registry is nil")
	}
	if id == 0 {
		return errors.New("provider id is required")
	}
	if provider == nil {
		return errors.New("provider is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.id2Provider[id]; ok {
		return errors.New("provider already registered")
	}
	r.id2Provider[id] = provider
	return nil
}

func (r *agentRegistry) GetProvider(id ProviderID) (IProvider, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.id2Provider[id]
	return provider, ok
}

func (r *agentRegistry) AddAgentSetting(agentID AgentID, setting IAgentSetting) error {
	if r == nil {
		return errors.New("agent registry is nil")
	}
	if agentID == 0 {
		return errors.New("agent id is required")
	}
	if setting == nil {
		return errors.New("agent setting is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.id2AgentSetting[agentID]; ok {
		return errors.New("agent setting already registered")
	}
	r.id2AgentSetting[agentID] = setting
	return nil
}

func (r *agentRegistry) GetAgentSetting(agentID AgentID) (IAgentSetting, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	setting, ok := r.id2AgentSetting[agentID]
	return setting, ok
}

func getAgentSettingAs[S IAgentSetting](registry *agentRegistry, agentID AgentID) (S, bool) {
	var empty S
	if registry == nil {
		return empty, false
	}
	setting, ok := registry.GetAgentSetting(agentID)
	if !ok {
		return empty, false
	}
	typed, ok := setting.(S)
	if !ok {
		return empty, false
	}
	return typed, true
}
