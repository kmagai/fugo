package subcmd

import (
	"strings"
)

type Remove struct {
	Style
}

func (c *Remove) Run(args []string) int {
	println("to be implemented")

	return 0
}

func (c *Remove) Synopsis() string {
	return "Sub-command"
}

func (c *Remove) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
