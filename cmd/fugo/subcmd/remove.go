package subcmd

import (
	"fmt"
	"strings"

	"github.com/kmagai/fugo"
)

type Remove struct {
	Style
}

// Run to remove a stock specified with stock code
func (c *Remove) Run(args []string) int {
	stockToRemove := args[0]
	portfolio := fugo.GetPortfolio()
	removed, err := portfolio.RemoveStock(stockToRemove)
	if err != nil {
		fmt.Println("Couldn't find the stock in your portfolio")
		return 1
	}
	// TODO: need better printing
	fmt.Printf("removed: %s", removed)
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
