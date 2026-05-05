package v2

import (
	"errors"
	"sync"
)

type RegistryID string

const RegistryIDDefault RegistryID = "default"

var DefaultRegistry = NewAgentRegistry()

var globalRegistry = struct {
	mu          sync.RWMutex
	id2Registry map[RegistryID]*AgentRegistry
}{
	id2Registry: map[RegistryID]*AgentRegistry{
		RegistryIDDefault: DefaultRegistry,
	},
}

func AddRegistry(id RegistryID, registry *AgentRegistry) error {
	if id == "" {
		return errors.New("registry id is required")
	}
	if registry == nil {
		return errors.New("agent registry is nil")
	}

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	if _, ok := globalRegistry.id2Registry[id]; ok {
		return errors.New("agent registry already registered")
	}
	globalRegistry.id2Registry[id] = registry
	return nil
}

func GetRegistry(id RegistryID) (*AgentRegistry, bool) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	registry, ok := globalRegistry.id2Registry[id]
	return registry, ok
}
