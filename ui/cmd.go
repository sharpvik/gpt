package ui

import (
	_ "embed"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
)

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
