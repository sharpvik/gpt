package repl

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/sharpvik/gpt/home"
)

type REPL struct {
	*bufio.ReadWriter
	gpt *openai.Client
}

func NewREPL() (*REPL, error) {
	if home.OpenAiApiKey == "" {
		return nil, errors.New("supply an API key using `gpt key <OPENAI_API_KEY>`")
	}
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	return &REPL{
		ReadWriter: bufio.NewReadWriter(stdin, stdout),
		gpt:        openai.NewClient(home.OpenAiApiKey),
	}, nil
}

func (repl *REPL) Read() (question string, err error) {
	if _, err := repl.WriteString("\n👾\n"); err != nil {
		return "", err
	}
	if err := repl.Flush(); err != nil {
		return "", err
	}
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

func (repl *REPL) Print(answer *openai.ChatCompletionStream) error {
	repl.WriteString("\n🤖\n")
	for {
		response, err := answer.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			repl.WriteString("\n🚨 Stream error: " + err.Error() + "\n")
			break
		}
		repl.WriteString(response.Choices[0].Delta.Content)
		repl.Flush()
	}
	repl.WriteString("\n\n")
	return repl.Flush()
}

func (repl *REPL) Loop() error {
	for {
		question, err := repl.Read()
		if err != nil {
			return err
		}
		answer, err := repl.Eval(question)
		if err != nil {
			return err
		}
		if err := repl.Print(answer); err != nil {
			return err
		}
	}
}

func (repl *REPL) EvalAndPrint(question string) error {
	answer, err := repl.Eval(question)
	if err != nil {
		return err
	}
	return repl.Print(answer)
}
