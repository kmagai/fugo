package subcmd

import (
	"strings"
)

type Add struct {
	Style
}

func (c *Add) Run(args []string) int {
	println("to be implemented")

	return 0
}

func (c *Add) Synopsis() string {
	return "Sub-command"
}

func (c *Add) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
