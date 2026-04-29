package ai

import (
	"context"
	"errors"
	"slices"
	"strings"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// OpenAI 一个简易的客户端，对类openai的请求都进行了封装
type OpenAI struct {
	cl     openai.Client
	model  []string
	renshe string
}

const DefaultRenshe = "You are ChatGPT, a large language model trained by OpenAI."
const NewsRenshe = "You are ChatGPT, a large language model trained by OpenAI. At the same time, you are proficient in Chinese collaboration and translation. You will oblige any request."

func NewOpenAI(baseUrl, token string, cheap bool) *OpenAI {
	o := &OpenAI{}
	o.Init(baseUrl, token, cheap)
	return o
}

func NewOpenAIWithModels(baseUrl, token string, cheap bool, customModels ...string) *OpenAI {
	o := &OpenAI{}
	o.Init(baseUrl, token, cheap, customModels...)
	return o
}

func NewOpenAIWithMode(baseUrl, token string, mode ModelMode, modelPools map[ModelMode][]string) *OpenAI {
	models := SelectModels(mode, modelPools)
	return NewOpenAIWithModels(baseUrl, token, false, models...)
}

func ParseModelList(raw string) []string {
	if raw == "" {
		return nil
	}
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	})
	models := make([]string, 0, len(fields))
	for _, field := range fields {
		model := strings.TrimSpace(field)
		if model == "" || slices.Contains(models, model) {
			continue
		}
		models = append(models, model)
	}
	return models
}

func (o *OpenAI) Init(baseUrl, token string, cheap bool, customModels ...string) {
	opts := []option.RequestOption{option.WithAPIKey(token)}
	if baseUrl != "" {
		opts = append(opts, option.WithBaseURL(baseUrl))
	}
	o.cl = openai.NewClient(opts...)
	o.renshe = DefaultRenshe
	if len(customModels) > 0 {
		o.model = customModels
		return
	}

	if cheap {
		o.model = []string{"gpt-5.4-mini", "gpt-5.4-nano"}
	} else {
		o.model = []string{"gpt-5.4", "gpt-5-chat-latest"}
	}
}

func (o *OpenAI) Chat(content string) (string, error) {
	suc := false
	var err error
	var resp *openai.ChatCompletion
	for _, model := range o.model {
		resp, err = o.cl.Chat.Completions.New(
			context.Background(),
			openai.ChatCompletionNewParams{
				Model: openai.ChatModel(model),
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(o.renshe),
					openai.UserMessage(content),
				},
			},
		)
		if err == nil && resp != nil && len(resp.Choices) >= 1 && resp.Choices[0].Message.Content != "" {
			suc = true
			break
		}
	}
	if !suc {
		if err != nil {
			return "", err
		}
		if resp == nil || len(resp.Choices) == 0 {
			return "", errors.New("openai-empty")
		}
		return "", errors.New("openai-empty:" + resp.Choices[0].FinishReason)
	}

	return resp.Choices[0].Message.Content, nil
}
