package subcmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kmagai/fugo/common"
)

// VersionCommand is the command name
type VersionCommand struct {
	Style

	Name     string
	Version  string
	Revision string
}

// Run specifies what this command does
func (c *VersionCommand) Run(args []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "%s version %s", c.Name, c.Version)
	if c.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	c.UI.Output(versionString.String())
	return common.ExitCodeOK
}

// Synopsis tells what it does
func (c *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Print %s version and quit", c.Name)
}

// Help text
func (c *VersionCommand) Help() string {
	helpText := `
	fugo help
`
	return strings.TrimSpace(helpText)
}
