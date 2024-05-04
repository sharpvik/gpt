package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.design/x/clipboard"
)

type (
	History struct {
		io.Writer
		Last            *Entry
		canUseClipboard bool
	}

	Entry struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
)

func NewHistory(file *os.File) (*History, error) {
	history := History{
		Writer: file,
	}
	if err := clipboard.Init(); err == nil {
		history.canUseClipboard = true
	}
	lastEntry, err := history.readLastEntry(file)
	if err != nil {
		return nil, err
	}
	history.Last = lastEntry
	return &history, nil
}

func (h *History) readLastEntry(file *os.File) (*Entry, error) {
	var lastEntry Entry

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read history: %s", err)
	}

	lines := bytes.Split(raw, []byte("\n"))
	//? lines is never empty, even if raw = []byte{}, lines will contain at
	//? least one slice with no bytes inside it, so that's how we check if the
	//? history is actually empty.
	if len(lines[0]) == 0 {
		return &lastEntry, nil
	}

	//? lines[len(lines)-1] is always [] because WriteEntry uses json.Encoder
	//? and it forcefully appends a \n after each write, so we have to use
	//? lines[len(lines)-2] instead.
	lastLine := lines[len(lines)-2]
	if err := json.Unmarshal(lastLine, &lastEntry); err != nil {
		return nil, fmt.Errorf("failed to read last history entry: %s", err)
	}
	return &lastEntry, nil
}

func (h *History) WriteEntry(entry *Entry) error {
	return json.NewEncoder(h).Encode(entry)
}

func (h *History) CopyLastAnswer() error {
	if h.canUseClipboard {
		clipboard.Write(clipboard.FmtText, []byte(h.Last.Answer))
		return nil
	}
	return errors.New("cannot use clipboard on this machine")
}
