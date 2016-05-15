package subcmd

import (
	"fmt"
	"strings"

	"github.com/kmagai/fugo"
	"github.com/kmagai/fugo/cmd/fugo/common"
)

type Remove struct {
	Style
}

// Run to remove a stock specified with stock code
func (c *Remove) Run(args []string) int {
	stockToRemove := args[0]
	portfolio, err := fugo.GetPortfolio()
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	removed, err := portfolio.RemoveStock(stockToRemove)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	// TODO: need better printing
	fmt.Printf("removed: %s", removed)
	return common.ExitCodeOK
}

func (c *Remove) Synopsis() string {
	return "Remove stock from your portfolio"
}

func (c *Remove) Help() string {
	helpText := `
	fugo remove [CODE]
`
	return strings.TrimSpace(helpText)
}
