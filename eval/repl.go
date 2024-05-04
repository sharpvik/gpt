package eval

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sharpvik/gpt/llm"
)

type REPL struct {
	*bufio.ReadWriter
	*History
	gpt4 *llm.GPT4
}

func NewREPL(history *History, gpt4 *llm.GPT4) (*REPL, error) {
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	return &REPL{
		ReadWriter: bufio.NewReadWriter(stdin, stdout),
		History:    history,
		gpt4:       gpt4,
	}, nil
}

func (repl *REPL) Read() (question string, err error) {
	repl.WriteString("\n👾 ")
	repl.Flush()
	if question, err = repl.ReadString(0); err == io.EOF {
		err = nil
	}
	return question, err
}

func (repl *REPL) Print(answer *openai.ChatCompletionStream) (string, error) {
	var buf strings.Builder

	repl.WriteString("\n🤖 ")
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
	stream, err := repl.gpt4.Stream(question)
	if err != nil {
		return err
	}
	answer, err := repl.Print(stream)
	repl.WriteEntry(&Entry{question, answer})
	return err
}
