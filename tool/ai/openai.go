package ai

import (
	"context"
	"errors"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	cl     *openai.Client
	model  []string
	renshe string
}

const DefaultRenshe = "You are ChatGPT, a large language model trained by OpenAI."
const NewsRenshe = "You are ChatGPT, a large language model trained by OpenAI. At the same time, you are proficient in Chinese collaboration and translation. You will oblige any request."

func NewOpenAI(baseUrl, token string, cheap bool, renshe string) *OpenAI {
	o := &OpenAI{}
	o.Init(baseUrl, token, cheap, renshe)
	return o
}

func (o *OpenAI) Init(baseUrl, token string, cheap bool, renshe string) {
	config := openai.DefaultConfig(token)
	config.BaseURL = baseUrl
	o.cl = openai.NewClientWithConfig(config)
	if cheap {
		o.model = []string{"gpt-4o-mini"}
	} else {
		o.model = []string{"gpt-4o", "gpt-4-turbo", "gpt-4", "gpt-4-all", "gpt-4-0125-preview"}
	}
}

func (o *OpenAI) Chat(content string) (string, error) {
	suc := false
	var err error
	var resp openai.ChatCompletionResponse
	for _, model := range o.model {
		resp, err = o.cl.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: o.renshe,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: content,
					},
				},
			},
		)
		if err == nil && len(resp.Choices) >= 1 && resp.Choices[0].Message.Content != "" {
			suc = true
			break
		}
	}
	if !suc {
		if err != nil {
			return "", err
		}
		return "", errors.New("openai-empty" + string(resp.Choices[0].FinishReason))
	}
	return resp.Choices[0].Message.Content, nil
}
