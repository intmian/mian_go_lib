package ai

var defaultRegistry = newAgentRegistry()

func AddProvider(id ProviderID, provider IProvider) error {
	return defaultRegistry.AddProvider(id, provider)
}

func GetProvider(id ProviderID) (IProvider, bool) {
	return defaultRegistry.GetProvider(id)
}

func AddAgentSetting(agentID AgentID, setting IAgentSetting) error {
	return defaultRegistry.AddAgentSetting(agentID, setting)
}

func GetAgentSetting(agentID AgentID) (IAgentSetting, bool) {
	return defaultRegistry.GetAgentSetting(agentID)
}

func GetAgentSettingAs[S IAgentSetting](agentID AgentID) (S, bool) {
	return getAgentSettingAs[S](defaultRegistry, agentID)
}

func resetDefaultRegistryForTest() {
	defaultRegistry = newAgentRegistry()
}
