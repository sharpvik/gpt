package repl

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type REPL struct {
	*bufio.ReadWriter
	*History
	gpt *openai.Client
}

func NewREPL(historyFile *os.File, apiKey string) (*REPL, error) {
	if apiKey == "" {
		return nil, errors.New("supply an API key using `gpt key <OPENAI_API_KEY>`")
	}

	history, err := NewHistory(historyFile)
	if err != nil {
		return nil, err
	}

	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	return &REPL{
		ReadWriter: bufio.NewReadWriter(stdin, stdout),
		History:    history,
		gpt:        openai.NewClient(apiKey),
	}, nil
}

func (repl *REPL) Read() (question string, err error) {
	repl.WriteString("\n👾\n")
	repl.Flush()
	if question, err = repl.ReadString(0); err == io.EOF {
		err = nil
	}
	return question, err
}

func (repl *REPL) Eval(question string) (answer *openai.ChatCompletionStream, err error) {
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4,
		Temperature: 0.8,
		N:           1,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	}
	return repl.gpt.CreateChatCompletionStream(context.Background(), req)
}

func (repl *REPL) Print(answer *openai.ChatCompletionStream) (string, error) {
	var buf strings.Builder

	repl.WriteString("\n🤖\n")
	for {
		response, err := answer.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return buf.String(), fmt.Errorf("stream error: %s", err)
		}
		delta := response.Choices[0].Delta.Content
		buf.WriteString(delta)
		repl.WriteString(delta)
		repl.Flush()
	}
	repl.WriteString("\n\n")

	return buf.String(), repl.Flush()
}

func (repl *REPL) Loop() error {
	for {
		question, err := repl.Read()
		if err != nil {
			return err
		}
		if err := repl.EvalAndPrint(question); err != nil {
			return err
		}
	}
}

func (repl *REPL) EvalAndPrint(question string) error {
	stream, err := repl.Eval(question)
	if err != nil {
		return err
	}
	answer, err := repl.Print(stream)
	repl.WriteEntry(&Entry{question, answer})
	return err
}
