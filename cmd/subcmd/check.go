package subcmd

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/kmagai/fugo"
	"github.com/kmagai/fugo/adapter"
	"github.com/kmagai/fugo/cmd/common"
)

// Check is the command name
type Check struct {
	Style
}

// Run specifies what this command does
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

	resource := adapter.NewGoogleAPI()
	stocksInPortfolio, err := fugo.GetStock(resource, portfolio.Stocks)
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

// Synopsis tells what it does
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
