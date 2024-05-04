package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sharpvik/gpt/eval"
	"github.com/sharpvik/gpt/home"
	"github.com/sharpvik/gpt/llm"
	"github.com/sharpvik/gpt/static"
	"github.com/sharpvik/gpt/ui"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:     static.Name,
	Usage:    "ChatGPT in your terminal",
	Version:  static.Version,
	Authors:  []*cli.Author{static.Author},
	Before:   home.Init,
	Commands: []*cli.Command{key, repl, copy},
	Action:   quickAnswer,
}

var key = &cli.Command{
	Name:      "key",
	Aliases:   []string{"k"},
	Usage:     "Specify OpenAI API key",
	Args:      true,
	ArgsUsage: "<OPENAI_API_KEY>",
	Action:    storeApiKey,
}

var repl = &cli.Command{
	Name:    "repl",
	Aliases: []string{"r"},
	Usage:   "Boot up the REPL",
	Before:  checkApiKey,
	Action:  REPL,
}

var copy = &cli.Command{
	Name:    "copy",
	Aliases: []string{"c"},
	Usage:   "Copy last response",
	Action:  copyLastAnswer,
}

func checkApiKey(_ *cli.Context) error {
	if home.OpenAiApiKey == "" {
		return errors.New("supply an API key using `gpt key <OPENAI_API_KEY>`")
	}
	return nil
}

func quickAnswer(ctx *cli.Context) error {
	if err := checkApiKey(ctx); err != nil {
		return err
	}

	if !ctx.Args().Present() {
		return cli.ShowAppHelp(ctx)
	}

	history, err := eval.NewHistory(home.HistoryFile)
	if err != nil {
		return err
	}

	repl, err := eval.NewREPL(history, llm.NewGPT4(home.OpenAiApiKey))
	if err != nil {
		return err
	}
	question := ctx.Args().First()
	return repl.EvalAndPrint(question)
}

func REPL(ctx *cli.Context) error {
	history, err := eval.NewHistory(home.HistoryFile)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(
		ui.New(llm.NewGPT4(home.OpenAiApiKey), history),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	).Run()
	return err
}

func storeApiKey(ctx *cli.Context) error {
	key := ctx.Args().First()
	return home.StoreApiKey(key)
}

func copyLastAnswer(ctx *cli.Context) error {
	history, err := eval.NewHistory(home.HistoryFile)
	if err != nil {
		return err
	}
	return history.CopyLastAnswer()
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
