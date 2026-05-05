package v2

import "testing"

func TestAgentRegistryRejectsDuplicateRegistrations(t *testing.T) {
	registry := NewAgentRegistry()
	provider := &fakeProvider{}
	if err := registry.AddProvider("main", provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}
	if err := registry.AddProvider("main", provider); err == nil {
		t.Fatal("expected duplicate provider error")
	}

	mother, err := NewBaseAgentMother(BaseAgentSetting{ProviderID: "main", Models: []string{"m1"}})
	if err != nil {
		t.Fatalf("NewBaseAgentMother error: %v", err)
	}
	if err := registry.AddAgentMother(AgentIDBase, mother); err != nil {
		t.Fatalf("AddAgentMother error: %v", err)
	}
	if err := registry.AddAgentMother(AgentIDBase, mother); err == nil {
		t.Fatal("expected duplicate mother error")
	}
}

func TestAgentMotherCreateAgentWithSettingUsesRegistryAsResolver(t *testing.T) {
	registry := NewAgentRegistry()
	provider := &fakeProvider{}
	if err := registry.AddProvider("override", provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}
	mother, err := NewBaseAgentMother(BaseAgentSetting{ProviderID: "default", Models: []string{"m1"}})
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
	agentAny, err := motherAny.CreateAnyAgentWithSetting(CoreAgentSetting{ProviderResolver: registry}, BaseAgentSetting{
		ProviderID: "override",
		Models:     []string{"m2"},
	})
	if err != nil {
		t.Fatalf("CreateAnyAgentWithSetting error: %v", err)
	}
	agent := agentAny.(*BaseAgent)
	if agent.setting.ProviderID != "override" {
		t.Fatalf("expected override provider, got %q", agent.setting.ProviderID)
	}
	if agent.setting.Models[0] != "m2" {
		t.Fatalf("expected override model, got %#v", agent.setting.Models)
	}
}

func TestAgentRegistrySetAndGetAgentSetting(t *testing.T) {
	registry := NewAgentRegistry()
	mother, err := NewBaseAgentMother(BaseAgentSetting{ProviderID: "old", Models: []string{"m1"}})
	if err != nil {
		t.Fatalf("NewBaseAgentMother error: %v", err)
	}
	if err := registry.AddAgentMother(AgentIDBase, mother); err != nil {
		t.Fatalf("AddAgentMother error: %v", err)
	}

	next := BaseAgentSetting{ProviderID: "new", Models: []string{"m2"}}
	motherAny, ok := registry.GetAgentMother(AgentIDBase)
	if !ok {
		t.Fatal("agent mother not registered")
	}
	if err := motherAny.SetAnyAgentSetting(next); err != nil {
		t.Fatalf("SetAgentSetting error: %v", err)
	}
	got := motherAny.GetAnyAgentSetting().(BaseAgentSetting)
	if got.ProviderID != "new" || got.Models[0] != "m2" {
		t.Fatalf("unexpected setting: %#v", got)
	}
}

func TestGlobalOnlyMaintainsRegistries(t *testing.T) {
	registry := NewAgentRegistry()
	id := RegistryID("test-global")
	if err := AddRegistry(id, registry); err != nil {
		t.Fatalf("AddRegistry error: %v", err)
	}
	got, ok := GetRegistry(id)
	if !ok {
		t.Fatal("registry not found")
	}
	if got != registry {
		t.Fatal("GetRegistry returned a different registry")
	}
	if err := AddRegistry(id, registry); err == nil {
		t.Fatal("expected duplicate registry error")
	}
	if defaultRegistry, ok := GetRegistry(RegistryIDDefault); !ok || defaultRegistry != DefaultRegistry {
		t.Fatal("default registry should be registered globally")
	}
}
