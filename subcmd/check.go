package subcmd

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/kmagai/fugo/common"
	"github.com/kmagai/fugo/lib"
	"github.com/kmagai/fugo/lib/interfaces"
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
	pf := fugo.NewPortfolio(usr.HomeDir + fugo.Fugorc)
	err = pf.GetPortfolio()

	if err != nil {
		fmt.Println(err)
		pf, err = pf.SetDefaultPortfolio()
		if err != nil {
			fmt.Println(err)
			return common.ExitCodeError
		}
	}
	var api interfaces.Resourcer = googleFinance.API{}
	stocks, err := api.GetStockers(pf.Codes)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	fmt.Println(stocks)

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
