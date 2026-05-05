package v2

import (
	"context"
	"errors"
	"strings"

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
	}, nil
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
	// v2 chat deliberately does not expose external tools/functions yet, so the
	// provider reports no callable tool surface even if the upstream model may
	// support one.
	return nil
}

func (p *OpenAIProvider) SourceType() ProviderSourceType {
	return p.sourceType
}
