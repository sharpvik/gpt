package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sharpvik/gpt/llm"
)

type Model struct {
	width   int
	height  int
	history string

	input textinput.Model
	chat  viewport.Model

	gpt4 *llm.GPT4
}

func New(gpt4 *llm.GPT4) *Model {
	input := textinput.New()
	input.Placeholder = "Enter your question here"
	input.Focus()

	chat := viewport.New(80, 10)
	chat.MouseWheelEnabled = true

	return &Model{
		input: input,
		chat:  chat,
		gpt4:  gpt4,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var gptCmd tea.Cmd

	switch msg := msg.(type) {
	case gptResponse:
		if msg.err != nil {
			m.history += fmt.Sprintf(" 🚨  Error fetching response from ChatGPT: %s\n\n", msg.err)
		} else {
			m.history += " 🤖  " + msg.answer.Choices[0].Message.Content + "\n\n"
		}
		m.chat.SetContent(m.history)
		m.chat.GotoBottom()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			question := m.input.Value()
			m.history += " 👾  " + question + "\n\n"
			gptCmd = m.askChatGPT(question)
			m.input.SetValue("")
			m.chat.SetContent(m.history)
			m.chat.GotoBottom()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.chat.Height = m.height - 4
		return m, nil
	}

	var inputCmd, chatCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	m.chat, chatCmd = m.chat.Update(msg)

	return m, tea.Batch(gptCmd, inputCmd, chatCmd)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Bottom,
		m.chatStyle().Render(m.chat.View()),
		m.inputStyle().Render(m.input.View()),
	)
}

/* STYLES */

func (m Model) chatStyle() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Border(lipgloss.NormalBorder(), true, true, false, true).
		BorderForeground(lipgloss.Color("#333333")).
		Padding(0, 1).
		Width(m.width - 2)
}

func (m Model) inputStyle() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666666")).
		Padding(0, 1).
		Width(m.width - 2)
}
