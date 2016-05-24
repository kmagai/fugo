package subcmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kmagai/fugo/cmd/common"
)

type VersionCommand struct {
	Style

	Name     string
	Version  string
	Revision string
}

func (c *VersionCommand) Run(args []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "%s version %s", c.Name, c.Version)
	if c.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	c.Ui.Output(versionString.String())
	return common.ExitCodeOK
}

func (c *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Print %s version and quit", c.Name)
}

func (c *VersionCommand) Help() string {
	helpText := `
	fugo help
`
	return strings.TrimSpace(helpText)
}
