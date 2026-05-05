package v2

import "errors"

type IAgentMother interface {
	CreateAnyAgent(core CoreAgentSetting) (IAgent, error)
	CreateAnyAgentWithSetting(core CoreAgentSetting, setting any) (IAgent, error)
	SetAnyAgentSetting(setting any) error
	GetAnyAgentSetting() any
}

type AgentMother[A IAgent, S any] struct {
	setting S
	builder func(core CoreAgentSetting, setting S) (A, error)
}

func NewAgentMother[A IAgent, S any](setting S, builder func(core CoreAgentSetting, setting S) (A, error)) (*AgentMother[A, S], error) {
	if builder == nil {
		return nil, errors.New("agent builder is nil")
	}
	return &AgentMother[A, S]{
		setting: setting,
		builder: builder,
	}, nil
}

func (m *AgentMother[A, S]) CreateAgent(core CoreAgentSetting) (A, error) {
	return m.CreateAgentWithSetting(core, m.setting)
}

func (m *AgentMother[A, S]) CreateAgentWithSetting(core CoreAgentSetting, setting S) (A, error) {
	var empty A
	if m == nil {
		return empty, errors.New("agent mother is nil")
	}
	if m.builder == nil {
		return empty, errors.New("agent builder is nil")
	}
	return m.builder(core, setting)
}

func (m *AgentMother[A, S]) SetAgentSetting(setting S) error {
	if m == nil {
		return errors.New("agent mother is nil")
	}
	m.setting = setting
	return nil
}

func (m *AgentMother[A, S]) GetAgentSetting() S {
	if m == nil {
		var empty S
		return empty
	}
	return m.setting
}

func (m *AgentMother[A, S]) CreateAnyAgent(core CoreAgentSetting) (IAgent, error) {
	return m.CreateAgent(core)
}

func (m *AgentMother[A, S]) CreateAnyAgentWithSetting(core CoreAgentSetting, setting any) (IAgent, error) {
	typedSetting, ok := setting.(S)
	if !ok {
		return nil, errors.New("agent setting type mismatch")
	}
	return m.CreateAgentWithSetting(core, typedSetting)
}

func (m *AgentMother[A, S]) SetAnyAgentSetting(setting any) error {
	typedSetting, ok := setting.(S)
	if !ok {
		return errors.New("agent setting type mismatch")
	}
	return m.SetAgentSetting(typedSetting)
}

func (m *AgentMother[A, S]) GetAnyAgentSetting() any {
	return m.GetAgentSetting()
}
