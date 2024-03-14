package ai

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	cl     *openai.Client
	model  string
	renshe string
}

const DefaultRenshe = "You are ChatGPT, a large language model trained by OpenAI."

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
		o.model = "gpt-3.5-turbo"
	} else {
		o.model = "gpt-4"
	}
}

func (o *OpenAI) Chat(content string) (string, error) {
	resp, err := o.cl.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: o.model,
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
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
