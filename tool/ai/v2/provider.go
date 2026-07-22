package ai

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ProviderSourceType string

const (
	ProviderSourceTypeOpenAI   ProviderSourceType = "OpenAI"
	ProviderSourceTypeDeepSeek ProviderSourceType = "DeepSeek"
)

type ReasoningEffort string

const (
	ReasoningEffortNone    ReasoningEffort = "none"
	ReasoningEffortMinimal ReasoningEffort = "minimal"
	ReasoningEffortLow     ReasoningEffort = "low"
	ReasoningEffortMedium  ReasoningEffort = "medium"
	ReasoningEffortHigh    ReasoningEffort = "high"
	ReasoningEffortXHigh   ReasoningEffort = "xhigh"
	ReasoningEffortMax     ReasoningEffort = "max"
)

type IProvider interface {
	AvailableModels(ctx context.Context) ([]string, error)
	AvailableReasoning() []ReasoningEffort
	AvailableTools() []string
	SourceType() ProviderSourceType
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
	ChatStream(ctx context.Context, req ChatRequest) (IChatStream, error)
}

type IModelProvider interface {
	IProvider
	RegisterModel(model ModelConfig) error
	GetModel(id ModelID) (ModelConfig, bool)
	RegisteredModels() []ModelConfig
}

func IsOpenAILikeSource(sourceType ProviderSourceType) bool {
	switch sourceType {
	case ProviderSourceTypeOpenAI, ProviderSourceTypeDeepSeek:
		return true
	default:
		return false
	}
}

type OpenAIProvider struct {
	client     openai.Client
	sourceType ProviderSourceType
	modelsMu   sync.RWMutex
	models     map[ModelID]ModelConfig
}

func NewOpenAIProvider(baseURL, token string, sourceType ProviderSourceType) (*OpenAIProvider, error) {
	if !IsOpenAILikeSource(sourceType) {
		return nil, errors.New("unsupported openai-like provider source")
	}

	opts := []option.RequestOption{option.WithAPIKey(token)}
	if strings.TrimSpace(baseURL) != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	return &OpenAIProvider{
		client:     openai.NewClient(opts...),
		sourceType: sourceType,
		models:     make(map[ModelID]ModelConfig),
	}, nil
}

func (p *OpenAIProvider) RegisterModel(model ModelConfig) error {
	if p == nil {
		return errors.New("provider is nil")
	}
	model.ID = ModelID(strings.TrimSpace(string(model.ID)))
	model.Name = strings.TrimSpace(model.Name)
	if err := model.Validate(); err != nil {
		return err
	}
	for _, effort := range model.Reasoning {
		if !slices.Contains(p.AvailableReasoning(), effort) {
			return fmt.Errorf("reasoning effort %q is not supported by %s adapter", effort, p.sourceType)
		}
	}
	for _, tool := range model.Tools {
		if !slices.Contains(p.AvailableTools(), string(tool)) {
			return fmt.Errorf("chat tool %q is not supported by %s adapter", tool, p.sourceType)
		}
	}

	p.modelsMu.Lock()
	defer p.modelsMu.Unlock()
	if p.models == nil {
		p.models = make(map[ModelID]ModelConfig)
	}
	if _, ok := p.models[model.ID]; ok {
		return fmt.Errorf("model %q already registered", model.ID)
	}
	model.Reasoning = append([]ReasoningEffort(nil), model.Reasoning...)
	model.Tools = append([]ChatTool(nil), model.Tools...)
	p.models[model.ID] = model
	return nil
}

func (p *OpenAIProvider) GetModel(id ModelID) (ModelConfig, bool) {
	if p == nil {
		return ModelConfig{}, false
	}
	p.modelsMu.RLock()
	defer p.modelsMu.RUnlock()
	model, ok := p.models[id]
	if !ok {
		return ModelConfig{}, false
	}
	model.Reasoning = append([]ReasoningEffort(nil), model.Reasoning...)
	model.Tools = append([]ChatTool(nil), model.Tools...)
	return model, true
}

func (p *OpenAIProvider) RegisteredModels() []ModelConfig {
	if p == nil {
		return nil
	}
	p.modelsMu.RLock()
	defer p.modelsMu.RUnlock()
	models := make([]ModelConfig, 0, len(p.models))
	for _, model := range p.models {
		model.Reasoning = append([]ReasoningEffort(nil), model.Reasoning...)
		model.Tools = append([]ChatTool(nil), model.Tools...)
		models = append(models, model)
	}
	sort.Slice(models, func(i, j int) bool { return models[i].ID < models[j].ID })
	return models
}

func (p *OpenAIProvider) AvailableModels(ctx context.Context) ([]string, error) {
	pager := p.client.Models.ListAutoPaging(ctx)
	models := make([]string, 0)
	for pager.Next() {
		model := pager.Current()
		if model.ID == "" {
			continue
		}
		models = append(models, model.ID)
	}
	if err := pager.Err(); err != nil {
		return nil, err
	}
	return models, nil
}

func (p *OpenAIProvider) AvailableReasoning() []ReasoningEffort {
	switch p.sourceType {
	case ProviderSourceTypeOpenAI:
		return []ReasoningEffort{
			ReasoningEffortNone,
			ReasoningEffortMinimal,
			ReasoningEffortLow,
			ReasoningEffortMedium,
			ReasoningEffortHigh,
			ReasoningEffortXHigh,
		}
	case ProviderSourceTypeDeepSeek:
		return []ReasoningEffort{
			ReasoningEffortHigh,
			ReasoningEffortMax,
		}
	default:
		return nil
	}
}

func (p *OpenAIProvider) AvailableTools() []string {
	if p.sourceType == ProviderSourceTypeOpenAI {
		return []string{string(ChatToolWebSearch)}
	}
	return nil
}

func (p *OpenAIProvider) SourceType() ProviderSourceType {
	return p.sourceType
}
