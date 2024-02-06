package ui

import (
	_ "embed"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
)

//go:embed ui.go
var mockContent []byte

type gptMsg struct {
	answer openai.ChatCompletionResponse
	err    error
}

func (m Model) askChatGPT(question string) tea.Cmd {
	return func() tea.Msg {
		answer, err := m.gpt4.Ask(question)
		return gptMsg{
			answer: answer,
			err:    err,
		}
	}
}

func (m Model) mockAskChatGPT(question string) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.DefaultClient.Get("https://www.google.com")
		if err != nil {
			return gptMsg{
				err: err,
			}
		}
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return gptMsg{
				err: err,
			}
		}
		return gptMsg{
			answer: openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: string(content),
						},
					},
				},
			},
		}
	}
}
