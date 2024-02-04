package repl

import (
	"bufio"
	"io"
	"os"
)

type REPL struct {
	*bufio.ReadWriter
}

func NewREPL() *REPL {
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	return &REPL{
		ReadWriter: bufio.NewReadWriter(stdin, stdout),
	}
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

func (repl *REPL) Eval(question string) (answer string, err error) {
	return "SOME ANSWER", nil
}

func (repl *REPL) Print(answer string) error {
	if _, err := repl.WriteString("\n🤖\n" + answer + "\n\n"); err != nil {
		return err
	}
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
