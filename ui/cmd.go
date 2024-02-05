package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
)

type gptResponse struct {
	answer openai.ChatCompletionResponse
	err    error
}

func (m Model) askChatGPT(question string) tea.Cmd {
	return func() tea.Msg {
		answer, err := m.gpt4.Ask(question)
		return gptResponse{
			answer: answer,
			err:    err,
		}
	}
}
