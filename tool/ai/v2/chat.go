package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

type ChatRole string

const (
	ChatRoleDeveloper ChatRole = "developer"
	ChatRoleSystem    ChatRole = "system"
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
)

type ThinkingType string

const (
	ThinkingTypeUnset    ThinkingType = ""
	ThinkingTypeEnabled  ThinkingType = "enabled"
	ThinkingTypeDisabled ThinkingType = "disabled"
)

type ChatMessage struct {
	Role    ChatRole
	Content string
}

type ChatRequest struct {
	Model           string
	Messages        []ChatMessage
	ReasoningEffort ReasoningEffort
	Thinking        ThinkingType
	MaxOutputTokens int64

	// Tools/functions are intentionally not part of the v2 first slice. Provider
	// methods only handle plain chat plus reasoning/thinking for now, and
	// AvailableTools returns no callable tools until this request surface exists.
	// Tools []ChatTool
	// Functions []FunctionTool
}

type ChatResponse struct {
	Text             string
	ReasoningContent string
	Model            string
	RawJSON          string
}

type ChatStreamEventType string

const (
	ChatStreamEventTextDelta      ChatStreamEventType = "text_delta"
	ChatStreamEventReasoningDelta ChatStreamEventType = "reasoning_delta"
	ChatStreamEventDone           ChatStreamEventType = "done"
)

type ChatStreamEvent struct {
	Type           ChatStreamEventType
	TextDelta      string
	ReasoningDelta string
	Response       *ChatResponse
}

type IChatStream interface {
	Recv() (ChatStreamEvent, error)
	Close() error
}

func (p *OpenAIProvider) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	if err := p.validateChatRequest(req); err != nil {
		return ChatResponse{}, err
	}

	switch p.sourceType {
	case ProviderSourceTypeOpenAI:
		return p.chatOpenAI(ctx, req)
	case ProviderSourceTypeDeepSeek:
		return p.chatDeepSeek(ctx, req)
	default:
		return ChatResponse{}, errors.New("unsupported provider source")
	}
}

func (p *OpenAIProvider) ChatStream(ctx context.Context, req ChatRequest) (IChatStream, error) {
	if err := p.validateChatRequest(req); err != nil {
		return nil, err
	}

	switch p.sourceType {
	case ProviderSourceTypeOpenAI:
		return p.chatOpenAIStream(ctx, req)
	case ProviderSourceTypeDeepSeek:
		return p.chatDeepSeekStream(ctx, req)
	default:
		return nil, errors.New("unsupported provider source")
	}
}

func (p *OpenAIProvider) chatOpenAI(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	params := responses.ResponseNewParams{
		Model: shared.ResponsesModel(req.Model),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toResponseInput(req.Messages),
		},
	}
	if req.MaxOutputTokens > 0 {
		params.MaxOutputTokens = openai.Int(req.MaxOutputTokens)
	}
	if req.ReasoningEffort != "" {
		params.Reasoning = openAIReasoningParam(req.ReasoningEffort)
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return ChatResponse{}, err
	}
	if resp == nil {
		return ChatResponse{}, errors.New("empty response")
	}

	return ChatResponse{
		Text:             resp.OutputText(),
		ReasoningContent: extractOpenAIReasoningSummary(resp.RawJSON()),
		Model:            string(resp.Model),
		RawJSON:          resp.RawJSON(),
	}, nil
}

func (p *OpenAIProvider) chatDeepSeek(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	params := openai.ChatCompletionNewParams{
		Model:    shared.ChatModel(req.Model),
		Messages: toChatCompletionMessages(req.Messages),
	}
	if req.MaxOutputTokens > 0 {
		params.MaxCompletionTokens = openai.Int(req.MaxOutputTokens)
	}
	if req.ReasoningEffort != "" {
		params.ReasoningEffort = shared.ReasoningEffort(req.ReasoningEffort)
	}

	opts := thinkingOptions(req.Thinking)
	resp, err := p.client.Chat.Completions.New(ctx, params, opts...)
	if err != nil {
		return ChatResponse{}, err
	}
	if resp == nil {
		return ChatResponse{}, errors.New("empty response")
	}

	text := ""
	if len(resp.Choices) > 0 {
		text = resp.Choices[0].Message.Content
	}

	return ChatResponse{
		Text:             text,
		ReasoningContent: extractDeepSeekReasoning(resp.RawJSON()),
		Model:            resp.Model,
		RawJSON:          resp.RawJSON(),
	}, nil
}

func (p *OpenAIProvider) chatOpenAIStream(ctx context.Context, req ChatRequest) (IChatStream, error) {
	params := responses.ResponseNewParams{
		Model: shared.ResponsesModel(req.Model),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toResponseInput(req.Messages),
		},
	}
	if req.MaxOutputTokens > 0 {
		params.MaxOutputTokens = openai.Int(req.MaxOutputTokens)
	}
	if req.ReasoningEffort != "" {
		params.Reasoning = openAIReasoningParam(req.ReasoningEffort)
	}

	return &openAIChatStream{
		stream: p.client.Responses.NewStreaming(ctx, params),
	}, nil
}

func (p *OpenAIProvider) chatDeepSeekStream(ctx context.Context, req ChatRequest) (IChatStream, error) {
	params := openai.ChatCompletionNewParams{
		Model:    shared.ChatModel(req.Model),
		Messages: toChatCompletionMessages(req.Messages),
	}
	if req.MaxOutputTokens > 0 {
		params.MaxCompletionTokens = openai.Int(req.MaxOutputTokens)
	}
	if req.ReasoningEffort != "" {
		params.ReasoningEffort = shared.ReasoningEffort(req.ReasoningEffort)
	}

	return &deepSeekChatStream{
		stream: p.client.Chat.Completions.NewStreaming(ctx, params, thinkingOptions(req.Thinking)...),
	}, nil
}

func toResponseInput(messages []ChatMessage) responses.ResponseInputParam {
	input := make(responses.ResponseInputParam, 0, len(messages))
	for _, message := range messages {
		input = append(input, responses.ResponseInputItemUnionParam{
			OfMessage: &responses.EasyInputMessageParam{
				Role: responses.EasyInputMessageRole(message.Role),
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(message.Content),
				},
			},
		})
	}
	return input
}

func toChatCompletionMessages(messages []ChatMessage) []openai.ChatCompletionMessageParamUnion {
	params := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, message := range messages {
		switch message.Role {
		case ChatRoleDeveloper:
			params = append(params, openai.DeveloperMessage(message.Content))
		case ChatRoleSystem:
			params = append(params, openai.SystemMessage(message.Content))
		case ChatRoleAssistant:
			params = append(params, openai.AssistantMessage(message.Content))
		default:
			params = append(params, openai.UserMessage(message.Content))
		}
	}
	return params
}

func (p *OpenAIProvider) validateChatRequest(req ChatRequest) error {
	if strings.TrimSpace(req.Model) == "" {
		return errors.New("chat model is required")
	}
	if len(req.Messages) == 0 {
		return errors.New("chat messages are required")
	}
	for i, message := range req.Messages {
		if !isValidChatRole(message.Role) {
			return fmt.Errorf("chat messages[%d] has invalid role %q", i, message.Role)
		}
	}
	if req.ReasoningEffort != "" && !slices.Contains(p.AvailableReasoning(), req.ReasoningEffort) {
		return fmt.Errorf("reasoning effort %q is not available for %s", req.ReasoningEffort, p.sourceType)
	}
	if !isValidThinkingType(req.Thinking) {
		return fmt.Errorf("thinking type %q is invalid", req.Thinking)
	}
	if p.sourceType != ProviderSourceTypeDeepSeek && req.Thinking != ThinkingTypeUnset {
		return fmt.Errorf("thinking is not available for %s", p.sourceType)
	}
	return nil
}

func isValidChatRole(role ChatRole) bool {
	switch role {
	case ChatRoleDeveloper, ChatRoleSystem, ChatRoleAssistant, ChatRoleUser:
		return true
	default:
		return false
	}
}

func isValidThinkingType(thinking ThinkingType) bool {
	switch thinking {
	case ThinkingTypeUnset, ThinkingTypeEnabled, ThinkingTypeDisabled:
		return true
	default:
		return false
	}
}

func openAIReasoningParam(effort ReasoningEffort) shared.ReasoningParam {
	param := shared.ReasoningParam{
		Effort: shared.ReasoningEffort(effort),
	}
	if effort != ReasoningEffortNone {
		param.Summary = shared.ReasoningSummaryAuto
	}
	return param
}

func thinkingOptions(thinking ThinkingType) []option.RequestOption {
	if thinking == ThinkingTypeUnset {
		return nil
	}
	return []option.RequestOption{
		option.WithJSONSet("thinking", map[string]string{"type": string(thinking)}),
	}
}

type openAIChatStream struct {
	stream    *ssestream.Stream[responses.ResponseStreamEventUnion]
	text      strings.Builder
	reasoning strings.Builder
}

func (s *openAIChatStream) Recv() (ChatStreamEvent, error) {
	for s.stream.Next() {
		event := s.stream.Current()
		switch v := event.AsAny().(type) {
		case responses.ResponseTextDeltaEvent:
			if v.Delta == "" {
				continue
			}
			s.text.WriteString(v.Delta)
			return ChatStreamEvent{Type: ChatStreamEventTextDelta, TextDelta: v.Delta}, nil
		case responses.ResponseReasoningTextDeltaEvent:
			if v.Delta == "" {
				continue
			}
			s.reasoning.WriteString(v.Delta)
			return ChatStreamEvent{Type: ChatStreamEventReasoningDelta, ReasoningDelta: v.Delta}, nil
		case responses.ResponseReasoningSummaryTextDeltaEvent:
			if v.Delta == "" {
				continue
			}
			s.reasoning.WriteString(v.Delta)
			return ChatStreamEvent{Type: ChatStreamEventReasoningDelta, ReasoningDelta: v.Delta}, nil
		case responses.ResponseCompletedEvent:
			text := v.Response.OutputText()
			if text == "" {
				text = s.text.String()
			}
			resp := ChatResponse{
				Text:             text,
				ReasoningContent: firstNonEmpty(s.reasoning.String(), extractOpenAIReasoningSummary(v.Response.RawJSON())),
				Model:            string(v.Response.Model),
				RawJSON:          v.Response.RawJSON(),
			}
			return ChatStreamEvent{Type: ChatStreamEventDone, Response: &resp}, nil
		}
	}
	if err := s.stream.Err(); err != nil {
		return ChatStreamEvent{}, err
	}
	return ChatStreamEvent{}, io.EOF
}

func (s *openAIChatStream) Close() error {
	return s.stream.Close()
}

type deepSeekChatStream struct {
	stream    *ssestream.Stream[openai.ChatCompletionChunk]
	text      strings.Builder
	reasoning strings.Builder
	pending   []ChatStreamEvent
}

func (s *deepSeekChatStream) Recv() (ChatStreamEvent, error) {
	if event, ok := s.popPending(); ok {
		return event, nil
	}

	for s.stream.Next() {
		chunk := s.stream.Current()
		if len(chunk.Choices) == 0 {
			continue
		}
		choice := chunk.Choices[0]
		reasoningDelta := extractDeepSeekStreamReasoning(chunk.RawJSON())
		if reasoningDelta != "" {
			s.reasoning.WriteString(reasoningDelta)
			s.pending = append(s.pending, ChatStreamEvent{Type: ChatStreamEventReasoningDelta, ReasoningDelta: reasoningDelta})
		}
		if choice.Delta.Content != "" {
			s.text.WriteString(choice.Delta.Content)
			s.pending = append(s.pending, ChatStreamEvent{Type: ChatStreamEventTextDelta, TextDelta: choice.Delta.Content})
		}
		if choice.FinishReason != "" {
			resp := ChatResponse{
				Text:             s.text.String(),
				ReasoningContent: s.reasoning.String(),
				Model:            chunk.Model,
				RawJSON:          chunk.RawJSON(),
			}
			s.pending = append(s.pending, ChatStreamEvent{Type: ChatStreamEventDone, Response: &resp})
		}
		if event, ok := s.popPending(); ok {
			return event, nil
		}
	}
	if err := s.stream.Err(); err != nil {
		return ChatStreamEvent{}, err
	}
	return ChatStreamEvent{}, io.EOF
}

func (s *deepSeekChatStream) Close() error {
	return s.stream.Close()
}

func (s *deepSeekChatStream) popPending() (ChatStreamEvent, bool) {
	if len(s.pending) == 0 {
		return ChatStreamEvent{}, false
	}
	event := s.pending[0]
	s.pending = s.pending[1:]
	return event, true
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func extractOpenAIReasoningSummary(raw string) string {
	var resp struct {
		Output []struct {
			Type    string `json:"type"`
			Summary []struct {
				Text string `json:"text"`
			} `json:"summary"`
		} `json:"output"`
	}
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		return ""
	}

	var summary strings.Builder
	for _, item := range resp.Output {
		if item.Type != "reasoning" {
			continue
		}
		for _, part := range item.Summary {
			summary.WriteString(part.Text)
		}
	}
	return summary.String()
}

func extractDeepSeekReasoning(raw string) string {
	var resp struct {
		Choices []struct {
			Message struct {
				ReasoningContent string `json:"reasoning_content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(raw), &resp); err != nil || len(resp.Choices) == 0 {
		return ""
	}
	return resp.Choices[0].Message.ReasoningContent
}

func extractDeepSeekStreamReasoning(raw string) string {
	var resp struct {
		Choices []struct {
			Delta struct {
				ReasoningContent string `json:"reasoning_content"`
			} `json:"delta"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(raw), &resp); err != nil || len(resp.Choices) == 0 {
		return ""
	}
	return resp.Choices[0].Delta.ReasoningContent
}
