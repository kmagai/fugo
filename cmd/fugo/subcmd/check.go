package subcmd

import (
	"fmt"
	"strings"

	"github.com/kmagai/fugo"
)

type Check struct {
	Style
}

func (c *Check) Run(args []string) int {
	portfolio := fugo.GetPortfolio()
	portfolio = portfolio.Update()
	portfolio.PrintStocks()

	return 0
}

func (c *Check) Synopsis() string {
	return fmt.Sprint("Check stock data in your portfolio")
}

func (c *Check) Help() string {
	helpText := `
`
	return strings.TrimSpace(helpText)
}
