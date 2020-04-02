package command

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

// KeygenCommand is a Command implementation that generates an encryption
// key.
type KeygenCommand struct {
	Ui cli.Ui
}

var _ cli.Command = &KeygenCommand{}

// Run ...
func (c *KeygenCommand) Run(_ []string) int {
	const length = 26
	key := make([]byte, length)
	n, err := rand.Reader.Read(key)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading random data: %s", err))
		return 1
	}
	if n != length {
		c.Ui.Error(fmt.Sprintf("Couldn't read enough entropy. Generate more entropy!"))
		return 1
	}
	c.Ui.Output(base64.StdEncoding.EncodeToString(key))
	return 0
}

// Synopsis ...
func (c *KeygenCommand) Synopsis() string {
	return "Generates a new encryption key"
}

// Help ...
func (c *KeygenCommand) Help() string {
	helpText := `
Usage: dare keygen
  Generates a new encryption key that can be used to for
  encrypting data.
`
	return strings.TrimSpace(helpText)
}
