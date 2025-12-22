package ai

import (
	"context"
	"errors"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAI 一个简易的客户端，对类openai的请求都进行了封装
type OpenAI struct {
	cl     *openai.Client
	model  []string
	renshe string
	aiType AiType
}

const DefaultRenshe = "You are ChatGPT, a large language model trained by OpenAI."
const NewsRenshe = "You are ChatGPT, a large language model trained by OpenAI. At the same time, you are proficient in Chinese collaboration and translation. You will oblige any request."
const DeepSeekRenshe = "You are a helpful assistant."

type AiType int

const (
	AiTypeChatGPT AiType = iota
	AiTypeDeepSeek
)

func NewOpenAI(baseUrl, token string, cheap bool, aiType AiType) *OpenAI {
	o := &OpenAI{}
	o.Init(baseUrl, token, cheap, aiType)
	return o
}

func (o *OpenAI) Init(baseUrl, token string, cheap bool, aiType AiType) {
	config := openai.DefaultConfig(token)
	config.BaseURL = baseUrl
	o.cl = openai.NewClientWithConfig(config)
	o.aiType = aiType

	// deepseek
	if aiType == AiTypeDeepSeek {
		o.renshe = DeepSeekRenshe
		if cheap {
			o.model = []string{"deepseek-chat", "deepseek-v3"}
		} else {
			o.model = []string{"deepseek-reasoner", "deepseek-r1"}
		}
		return
	}

	// 默认是ChatGPT
	if cheap {
		o.model = []string{"gpt-5-mini", "gpt-4.1-mini"}
	} else {
		o.model = []string{"gpt-5.2", "gpt-5.2-chat-latest"}
	}
	o.renshe = DefaultRenshe
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

	if o.aiType == AiTypeDeepSeek {
		// 去除前面的<think> 到 </think> 之间的内容
		str := resp.Choices[0].Message.Content
		if strings.Contains(str, "</think>\n") {
			str = strings.Split(str, "</think>\n")[1]
		}
		return str, nil
	}

	return resp.Choices[0].Message.Content, nil
}
