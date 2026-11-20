package ai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

// ProviderOption 调整 OpenAIProvider 的底层请求行为。
type ProviderOption func(*providerOptions)

type providerOptions struct {
	httpClient *http.Client
	maxRetries *int
}

// WithHTTPClient 指定发请求用的 HTTP 客户端（默认 http.DefaultClient，连接超时 30 秒、等待响应不限时）。
func WithHTTPClient(client *http.Client) ProviderOption {
	return func(o *providerOptions) { o.httpClient = client }
}

// WithMaxRetries 设置 SDK 内部的自动重试次数（默认 2 次，连接失败也会重试）。传 0 表示由调用方自己决定是否重试。
func WithMaxRetries(retries int) ProviderOption {
	return func(o *providerOptions) {
		if retries < 0 {
			retries = 0
		}
		o.maxRetries = &retries
	}
}

// HTTPStatusFromError 取出上游返回的 HTTP 状态码；连接失败、超时等没有收到响应的错误返回 false。
func HTTPStatusFromError(err error) (int, bool) {
	var apiErr *openai.Error
	if errors.As(err, &apiErr) && apiErr.StatusCode > 0 {
		return apiErr.StatusCode, true
	}
	return 0, false
}

func NewOpenAIProvider(baseURL, token string, sourceType ProviderSourceType, providerOpts ...ProviderOption) (*OpenAIProvider, error) {
	if !IsOpenAILikeSource(sourceType) {
		return nil, errors.New("unsupported openai-like provider source")
	}

	opts := []option.RequestOption{option.WithAPIKey(token)}
	if strings.TrimSpace(baseURL) != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	var settings providerOptions
	for _, apply := range providerOpts {
		if apply != nil {
			apply(&settings)
		}
	}
	if settings.httpClient != nil {
		opts = append(opts, option.WithHTTPClient(settings.httpClient))
	}
	if settings.maxRetries != nil {
		opts = append(opts, option.WithMaxRetries(*settings.maxRetries))
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

// AvailableModelInfo 是上游 /models 里的一个模型。各家扩展字段不同，Raw 保留原始 JSON 由调用方解析。
type AvailableModelInfo struct {
	ID  string
	Raw string
}

func (p *OpenAIProvider) AvailableModelInfos(ctx context.Context) ([]AvailableModelInfo, error) {
	pager := p.client.Models.ListAutoPaging(ctx)
	models := make([]AvailableModelInfo, 0)
	for pager.Next() {
		model := pager.Current()
		if model.ID == "" {
			continue
		}
		models = append(models, AvailableModelInfo{ID: model.ID, Raw: model.RawJSON()})
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
			ReasoningEffortMax, // GPT-6 系列支持
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
