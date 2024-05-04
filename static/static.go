package static

import "github.com/urfave/cli/v2"

const (
	/* General information */

	Name    = "gpt"
	Version = "v0.2.2"

	/* File names */

	ConfigDirName        = "." + Name
	HistoryFileName      = "history.jsonl"
	OpenAiApiKeyFileName = "api.key"
)

var Author = &cli.Author{
	Name:  "Viktor A. Rozenko Voitenko",
	Email: "sharp.vik@gmail.com",
}
