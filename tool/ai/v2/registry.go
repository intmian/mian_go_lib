package v2

import (
	"errors"
	"sync"
)

type ProviderResolver interface {
	GetProvider(id ProviderID) (IProvider, bool)
}

type CoreAgentSetting struct {
	ProviderResolver ProviderResolver
}

func (c CoreAgentSetting) GetProvider(id ProviderID) (IProvider, bool) {
	if c.ProviderResolver == nil {
		return nil, false
	}
	return c.ProviderResolver.GetProvider(id)
}

type AgentRegistry struct {
	mu             sync.RWMutex
	id2AgentMother map[AgentID]IAgentMother
	id2Provider    map[ProviderID]IProvider
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		id2AgentMother: make(map[AgentID]IAgentMother),
		id2Provider:    make(map[ProviderID]IProvider),
	}
}

func (r *AgentRegistry) AddProvider(id ProviderID, provider IProvider) error {
	if r == nil {
		return errors.New("agent registry is nil")
	}
	if id == "" {
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

func (r *AgentRegistry) GetProvider(id ProviderID) (IProvider, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.id2Provider[id]
	return provider, ok
}

func (r *AgentRegistry) AddAgentMother(agentID AgentID, mother IAgentMother) error {
	if r == nil {
		return errors.New("agent registry is nil")
	}
	if agentID == 0 {
		return errors.New("agent id is required")
	}
	if mother == nil {
		return errors.New("agent mother is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.id2AgentMother[agentID]; ok {
		return errors.New("agent mother already registered")
	}
	r.id2AgentMother[agentID] = mother
	return nil
}

func (r *AgentRegistry) GetAgentMother(agentID AgentID) (IAgentMother, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	mother, ok := r.id2AgentMother[agentID]
	return mother, ok
}
