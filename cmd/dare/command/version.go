package command

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/cli"
)

// VersionCommand is a Command implementation prints the version.
type VersionCommand struct {
	Name              string
	Revision          string
	Version           string
	VersionPrerelease string
	Ui                cli.Ui
}

var _ cli.Command = &VersionCommand{}

// Help ...
func (c *VersionCommand) Help() string {
	return fmt.Sprintf("Prints %s version", c.Name)
}

// Run ...
func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "%s v%s", c.Name, c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, ".%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	c.Ui.Output(versionString.String())
	return 0
}

// Synopsis ...
func (c *VersionCommand) Synopsis() (s string) {
	s = fmt.Sprintf("Prints %s version", c.Name)
	return s
}
