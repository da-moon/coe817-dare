package command

import (
	"encoding/hex"
	"testing"

	"github.com/mitchellh/cli"
)

func TestKeygenCommand(t *testing.T) {
	ui := new(cli.MockUi)
	c := &KeygenCommand{Ui: ui}
	code := c.Run(nil)
	if code != 0 {
		t.Fatalf("bad: %d", code)
	}

	output := ui.OutputWriter.String()
	result, err := hex.DecodeString(output)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(result) != 26 {
		t.Fatalf("bad: %#v", result)
	}
}
