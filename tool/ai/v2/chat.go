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

type ChatTool string

const (
	ChatToolWebSearch ChatTool = "web_search"
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
	MaxToolCalls    int64
	Tools           []ChatTool
}

type ChatCitation struct {
	StartIndex int64
	EndIndex   int64
	Title      string
	URL        string
}

type ChatResponse struct {
	Text             string
	ReasoningContent string
	Model            string
	Citations        []ChatCitation
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
	var err error
	req, err = p.resolveRegisteredModel(req)
	if err != nil {
		return ChatResponse{}, err
	}
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
	var err error
	req, err = p.resolveRegisteredModel(req)
	if err != nil {
		return nil, err
	}
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

func (p *OpenAIProvider) resolveRegisteredModel(req ChatRequest) (ChatRequest, error) {
	registered := p.RegisteredModels()
	if len(registered) == 0 {
		return req, nil
	}
	modelID := ModelID(strings.TrimSpace(req.Model))
	model, ok := p.GetModel(modelID)
	if !ok {
		return ChatRequest{}, fmt.Errorf("model %q is not registered", modelID)
	}
	if model.Type != ModelTypeText {
		return ChatRequest{}, fmt.Errorf("model %q is not a text model", modelID)
	}
	if req.ReasoningEffort != "" && !slices.Contains(model.Reasoning, req.ReasoningEffort) {
		return ChatRequest{}, fmt.Errorf("reasoning effort %q is not available for model %q", req.ReasoningEffort, modelID)
	}
	for i, tool := range req.Tools {
		if !slices.Contains(model.Tools, tool) {
			return ChatRequest{}, fmt.Errorf("chat tools[%d] %q is not available for model %q", i, tool, modelID)
		}
	}
	req.Model = model.Name
	return req, nil
}

func (p *OpenAIProvider) chatOpenAI(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	params := openAIResponseParams(req)
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
		Citations:        extractOpenAIURLCitations(resp),
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
	params := openAIResponseParams(req)
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

func openAIResponseParams(req ChatRequest) responses.ResponseNewParams {
	params := responses.ResponseNewParams{
		Model: shared.ResponsesModel(req.Model),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toResponseInput(req.Messages),
		},
	}
	if req.MaxOutputTokens > 0 {
		params.MaxOutputTokens = openai.Int(req.MaxOutputTokens)
	}
	if req.MaxToolCalls > 0 {
		params.MaxToolCalls = openai.Int(req.MaxToolCalls)
	}
	if req.ReasoningEffort != "" {
		params.Reasoning = openAIReasoningParam(req.ReasoningEffort)
	}
	for _, tool := range req.Tools {
		switch tool {
		case ChatToolWebSearch:
			params.Tools = append(params.Tools, responses.ToolParamOfWebSearch(responses.WebSearchToolTypeWebSearch))
		}
	}
	return params
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
	if req.MaxToolCalls < 0 {
		return errors.New("max tool calls cannot be negative")
	}
	if req.MaxToolCalls > 0 && len(req.Tools) == 0 {
		return errors.New("max tool calls requires at least one chat tool")
	}
	availableTools := p.AvailableTools()
	seenTools := make(map[ChatTool]struct{}, len(req.Tools))
	for i, tool := range req.Tools {
		if !slices.Contains(availableTools, string(tool)) {
			return fmt.Errorf("chat tools[%d] %q is not available for %s", i, tool, p.sourceType)
		}
		if _, ok := seenTools[tool]; ok {
			return fmt.Errorf("chat tools[%d] %q is duplicated", i, tool)
		}
		seenTools[tool] = struct{}{}
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
				Citations:        extractOpenAIURLCitations(&v.Response),
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

func extractOpenAIURLCitations(resp *responses.Response) []ChatCitation {
	if resp == nil {
		return nil
	}

	var citations []ChatCitation
	for _, item := range resp.Output {
		if item.Type != "message" {
			continue
		}
		for _, content := range item.Content {
			if content.Type != "output_text" {
				continue
			}
			for _, annotation := range content.Annotations {
				if annotation.Type != "url_citation" {
					continue
				}
				citations = append(citations, ChatCitation{
					StartIndex: annotation.StartIndex,
					EndIndex:   annotation.EndIndex,
					Title:      annotation.Title,
					URL:        annotation.URL,
				})
			}
		}
	}
	return citations
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
