package main

import (
	"fmt"
	"os"

	"github.com/sharpvik/gpt/home"
	"github.com/sharpvik/gpt/repl"
	"github.com/sharpvik/gpt/static"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:     static.Name,
	Usage:    "ChatGPT in your terminal",
	Version:  static.Version,
	Authors:  []*cli.Author{static.Author},
	Before:   home.Init,
	Commands: []*cli.Command{key},
	Action:   evalArgOrLoop,
}

var key = &cli.Command{
	Name:      "key",
	Aliases:   []string{"k"},
	Usage:     "Specify OpenAI API key",
	Args:      true,
	ArgsUsage: "<OPENAI_API_KEY>",
	Action:    storeApiKey,
}

func evalArgOrLoop(ctx *cli.Context) error {
	repl := repl.NewREPL()
	if ctx.Args().Present() {
		question := ctx.Args().First()
		return repl.EvalAndPrint(question)
	}
	return repl.Loop()
}

func storeApiKey(ctx *cli.Context) error {
	key := ctx.Args().First()
	return home.StoreApiKey(key)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		abort(err)
	}
}

func abort(args ...any) {
	fmt.Println(args...)
	os.Exit(1)
}
