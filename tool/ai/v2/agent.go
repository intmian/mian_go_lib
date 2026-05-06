package ai

type AgentIDProvider interface {
	GetID() AgentID
}

type IAgent[S IAgentSetting] interface {
	AgentIDProvider
	Init() error
	InitWithSetting(setting S) error
}
