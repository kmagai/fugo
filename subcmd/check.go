package subcmd

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/kmagai/fugo/common"
	"github.com/kmagai/fugo/lib"
	"github.com/kmagai/googleFinance"
)

// Check is the command name.
type Check struct {
	Style
}

// Run specifies what this command does.
func (c *Check) Run(args []string) int {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	portfolio := fugo.NewPortfolio(usr.HomeDir + fugo.Fugorc)
	portfolio, err = portfolio.GetPortfolio()

	if err != nil {
		fmt.Println(err)
		portfolio, err = portfolio.SetDefaultPortfolio()
		if err != nil {
			fmt.Println(err)
			return common.ExitCodeError
		}
	}
	source := googleFinance.API{}
	var codes []string
	for _, stock := range portfolio.Stocks {
		codes = append(codes, stock.Code)
	}
	stocksInPortfolio, err := source.GetStocks(codes)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}

	portfolio, err = portfolio.Update(stocksInPortfolio)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	common.ShowPortfolio(portfolio)

	return common.ExitCodeOK
}

// Synopsis tells what it does.
func (c *Check) Synopsis() string {
	return fmt.Sprint("Check stock data in your portfolio")
}

// Help text
func (c *Check) Help() string {
	helpText := `
	fugo check
	or
	fugo
`
	return strings.TrimSpace(helpText)
}
