package subcmd

import (
	"fmt"
	"strings"

	"github.com/kmagai/fugo"
	"github.com/kmagai/fugo/cmd/fugo/common"
)

type Check struct {
	Style
}

func (c *Check) Run(args []string) int {
	portfolio, err := fugo.GetPortfolio()
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	portfolio, err = portfolio.Update()
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	common.ShowPortfolio(portfolio)

	return common.ExitCodeOK
}

func (c *Check) Synopsis() string {
	return fmt.Sprint("Check stock data in your portfolio")
}

func (c *Check) Help() string {
	helpText := `
	fugo check
	or
	fugo
`
	return strings.TrimSpace(helpText)
}
