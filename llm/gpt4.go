package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type GPT4 struct {
	*openai.Client
}

func NewGPT4(token string) *GPT4 {
	return &GPT4{
		Client: openai.NewClient(token),
	}
}

func (g *GPT4) Ask(question string) (openai.ChatCompletionResponse, error) {
	return g.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT4,
			Temperature: 0.8,
			N:           1,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		})
}

func (g *GPT4) Stream(question string) (*openai.ChatCompletionStream, error) {
	return g.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT4,
			Temperature: 0.8,
			N:           1,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		})
}
