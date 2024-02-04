package home

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/sharpvik/gpt/static"
	"github.com/urfave/cli/v2"
)

var (
	homeDir   string //? ~
	configDir string //? ~/.gpt

	historyFilePath string //? ~/.gpt/history.jsonl
	historyFile     *os.File

	openAiApiKeyFilePath string //? ~/.gpt/api.key
	OpenAiApiKey         string
)

func Init(_ *cli.Context) (err error) {
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to find user's home directory: %s", err)
	}

	configDir = path.Join(homeDir, static.ConfigDirName)
	if err := os.MkdirAll(configDir, fs.ModePerm); err != nil {
		return fmt.Errorf("failed to create configuration folder: %s", err)
	}

	historyFilePath = path.Join(configDir, static.HistoryFileName)
	if historyFile, err = os.OpenFile(
		historyFilePath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0600,
	); err != nil {
		return fmt.Errorf("failed to open the history file: %s", err)
	}

	openAiApiKeyFilePath = path.Join(configDir, static.OpenAiApiKeyFileName)
	OpenAiApiKey = ReadApiKey()

	return nil
}

func ReadApiKey() string {
	file, err := os.Open(openAiApiKeyFilePath)
	if err != nil {
		return ""
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func StoreApiKey(key string) error {
	file, err := os.Create(openAiApiKeyFilePath)
	if err != nil {
		return fmt.Errorf("failed to create API key file: %s", err)
	}

	if _, err := file.WriteString(key); err != nil {
		return fmt.Errorf("failed to store API key in file: %s", err)
	}

	return nil
}
