package ai

import "testing"

func TestAgentRegistryRejectsDuplicateRegistrations(t *testing.T) {
	resetDefaultRegistryForTest()
	provider := &fakeProvider{}
	if err := AddProvider(1, provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}
	if err := AddProvider(1, provider); err == nil {
		t.Fatal("expected duplicate provider error")
	}

	setting := &BaseAgentSetting{ProviderID: 1, Models: []string{"m1"}}
	if err := AddAgentSetting(AgentIDBase, setting); err != nil {
		t.Fatalf("AddAgentSetting error: %v", err)
	}
	if err := AddAgentSetting(AgentIDBase, setting); err == nil {
		t.Fatal("expected duplicate setting error")
	}
}

func TestNewBaseAgentUsesRegistryAsResolver(t *testing.T) {
	resetDefaultRegistryForTest()
	provider := &fakeProvider{}
	if err := AddProvider(2, provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}

	agent := NewBaseAgent()
	err := agent.InitWithSetting(&BaseAgentSetting{
		ProviderID: 2,
		Models:     []string{"m2"},
	})
	if err != nil {
		t.Fatalf("InitWithSetting error: %v", err)
	}
	setting := agent.setting.GetSetting()
	if setting.ProviderID != 2 {
		t.Fatalf("expected override provider, got %d", setting.ProviderID)
	}
	if setting.Models[0] != "m2" {
		t.Fatalf("expected override model, got %#v", setting.Models)
	}
}

func TestAgentRegistryGetAgentSettingAs(t *testing.T) {
	resetDefaultRegistryForTest()
	setting := &BaseAgentSetting{ProviderID: 1, Models: []string{"m1"}}
	if err := AddAgentSetting(AgentIDBase, setting); err != nil {
		t.Fatalf("AddAgentSetting error: %v", err)
	}

	got, ok := GetAgentSettingAs[*BaseAgentSetting](AgentIDBase)
	if !ok {
		t.Fatal("base agent setting not found")
	}
	got.ProviderID = 2
	got.Models = []string{"m2"}

	gotAgain, ok := GetAgentSettingAs[*BaseAgentSetting](AgentIDBase)
	if !ok {
		t.Fatal("base agent setting not found after update")
	}
	if gotAgain.ProviderID != 2 || gotAgain.Models[0] != "m2" {
		t.Fatalf("unexpected setting: %#v", got)
	}
}

func TestPackageLevelGettersUseDefaultRegistry(t *testing.T) {
	resetDefaultRegistryForTest()
	provider := &fakeProvider{}
	if err := AddProvider(1, provider); err != nil {
		t.Fatalf("AddProvider error: %v", err)
	}
	gotProvider, ok := GetProvider(1)
	if !ok || gotProvider != provider {
		t.Fatal("GetProvider should read from default registry")
	}

	setting := &BaseAgentSetting{ProviderID: 1, Models: []string{"m1"}}
	if err := AddAgentSetting(AgentIDBase, setting); err != nil {
		t.Fatalf("AddAgentSetting error: %v", err)
	}
	gotSetting, ok := GetAgentSetting(AgentIDBase)
	if !ok || gotSetting != setting {
		t.Fatal("GetAgentSetting should read from default registry")
	}
}
