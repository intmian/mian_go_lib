package ai

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type ModelID string

type ModelType string

type ModelCallProtocol string

const (
	ModelTypeText ModelType = "text"
	ModelTypeSTT  ModelType = "stt"

	// ModelTypeChat is kept as a source-compatible alias. New configuration
	// should persist ModelTypeText.
	ModelTypeChat = ModelTypeText
)

const (
	ModelCallProtocolOpenAIText        ModelCallProtocol = "OpenAIText"
	ModelCallProtocolDeepSeekText      ModelCallProtocol = "DeepSeekText"
	ModelCallProtocolOpenAISTT         ModelCallProtocol = "OpenAISTT"
	ModelCallProtocolDashScopeQwen3ASR ModelCallProtocol = "DashScopeQwen3ASR"
	ModelCallProtocolDashScopeFunASR   ModelCallProtocol = "DashScopeFunASR"
)

type ModelConfig struct {
	ID           ModelID           `json:"id"`
	Name         string            `json:"name"`
	Type         ModelType         `json:"type"`
	CallProtocol ModelCallProtocol `json:"callProtocol,omitempty"`
	Reasoning    []ReasoningEffort `json:"reasoning,omitempty"`
	Tools        []ChatTool        `json:"tools,omitempty"`
}

func (m ModelConfig) Validate() error {
	if strings.TrimSpace(string(m.ID)) == "" {
		return errors.New("model id is required")
	}
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("model name is required")
	}
	switch m.Type {
	case ModelTypeText, ModelTypeSTT:
	default:
		return fmt.Errorf("model type %q is invalid", m.Type)
	}
	if err := validateModelCallProtocol(m.Type, m.CallProtocol); err != nil {
		return err
	}
	for i, effort := range m.Reasoning {
		if !isValidReasoningEffort(effort) {
			return fmt.Errorf("model reasoning[%d] %q is invalid", i, effort)
		}
		if slices.Contains(m.Reasoning[:i], effort) {
			return fmt.Errorf("model reasoning[%d] %q is duplicated", i, effort)
		}
	}
	for i, tool := range m.Tools {
		if !isValidChatTool(tool) {
			return fmt.Errorf("model tools[%d] %q is invalid", i, tool)
		}
		if slices.Contains(m.Tools[:i], tool) {
			return fmt.Errorf("model tools[%d] %q is duplicated", i, tool)
		}
	}
	if m.Type != ModelTypeText && (len(m.Reasoning) > 0 || len(m.Tools) > 0) {
		return errors.New("non-text model cannot configure reasoning or chat tools")
	}
	return nil
}

func validateModelCallProtocol(modelType ModelType, protocol ModelCallProtocol) error {
	if protocol == "" {
		return nil
	}
	switch modelType {
	case ModelTypeText:
		switch protocol {
		case ModelCallProtocolOpenAIText, ModelCallProtocolDeepSeekText:
			return nil
		default:
			return fmt.Errorf("text model call protocol %q is invalid", protocol)
		}
	case ModelTypeSTT:
		switch protocol {
		case ModelCallProtocolOpenAISTT, ModelCallProtocolDashScopeQwen3ASR, ModelCallProtocolDashScopeFunASR:
			return nil
		default:
			return fmt.Errorf("stt model call protocol %q is invalid", protocol)
		}
	default:
		return fmt.Errorf("model type %q is invalid", modelType)
	}
	return nil
}

func isValidReasoningEffort(effort ReasoningEffort) bool {
	switch effort {
	case ReasoningEffortNone,
		ReasoningEffortMinimal,
		ReasoningEffortLow,
		ReasoningEffortMedium,
		ReasoningEffortHigh,
		ReasoningEffortXHigh,
		ReasoningEffortMax:
		return true
	default:
		return false
	}
}

func isValidChatTool(tool ChatTool) bool {
	switch tool {
	case ChatToolWebSearch:
		return true
	default:
		return false
	}
}
