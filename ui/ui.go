package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sharpvik/gpt/eval"
	"github.com/sharpvik/gpt/llm"
)

type Model struct {
	width  int
	height int

	focus bool
	input textinput.Model

	chatHistory string
	chat        viewport.Model

	gpt4         *llm.GPT4
	history      *eval.History
	historyEntry *eval.Entry
}

func New(gpt4 *llm.GPT4, history *eval.History) *Model {
	focus := true
	input := textinput.New()
	input.Placeholder = "Enter your question here"
	input.Focus()

	localHistory := aiMessage("Hey there! How can I help you today?")
	chat := viewport.New(80, 10)
	chat.MouseWheelEnabled = true
	chat.SetContent(localHistory)

	return &Model{
		focus: focus,
		input: input,

		chatHistory: localHistory,
		chat:        chat,

		gpt4:         gpt4,
		history:      history,
		historyEntry: &eval.Entry{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var gptCmd, inputCmd, chatCmd tea.Cmd

	if m.focus {
		m.input, inputCmd = m.input.Update(msg)
	} else {
		m.chat, chatCmd = m.chat.Update(msg)
	}

	switch msg := msg.(type) {
	case gptMsg:
		m = m.updateWithGptMsg(msg)
	case tea.KeyMsg:
		m, gptCmd = m.updateWithKeyMsg(msg)
	case tea.WindowSizeMsg:
		m = m.updateWithWindowSizeMsg(msg)
	}

	return m, tea.Batch(gptCmd, inputCmd, chatCmd)
}

func (m Model) updateWithGptMsg(msg gptMsg) Model {
	if msg.err != nil {
		return m.updateChatHistory(errorMessage(msg.err))
	}
	answer := msg.answer.Choices[0].Message.Content
	m.historyEntry.Answer = answer
	m.history.WriteEntry(m.historyEntry)
	return m.updateChatHistory(aiMessage(answer))
}

func (m Model) updateWithKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.focus = false
		m.input.Blur()

	case "i":
		m.focus = true
		m.input.Focus()

	case "c":
		m.history.CopyLastAnswer()

	case "enter":
		question := m.input.Value()
		m.input.SetValue("")
		m.historyEntry.Question = question
		m = m.updateChatHistory(humanMessage(question))
		return m, m.askChatGPT(question)
	}
	return m, nil
}

func (m Model) updateWithWindowSizeMsg(msg tea.WindowSizeMsg) Model {
	m.width = msg.Width
	m.height = msg.Height
	m.chat.Height = m.height - 5
	m.chat.Width = m.width - 4
	m.input.Width = m.width - 7
	return m
}

func (m Model) updateChatHistory(message string) Model {
	m.chatHistory += message
	m.chat.
		SetContent(lipgloss.NewStyle().Width(m.chat.Width).
			Render(m.chatHistory))
	m.chat.GotoBottom()
	return m
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.chatStyle().Render(m.chat.View()),
		m.inputStyle().Render(m.input.View()),
	)
}

/* STYLES */

func (m Model) chatStyle() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Border(lipgloss.NormalBorder()).
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
