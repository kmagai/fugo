package subcmd

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/kmagai/fugo"
	"github.com/kmagai/fugo/cmd/fugo/common"
)

type Check struct {
	Style
}

func (c *Check) Run(args []string) int {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	portfolio := &fugo.Portfolio{}
	portfolio.Path = usr.HomeDir + fugo.Fugorc
	portfolio, err = portfolio.GetPortfolio()
	if err != nil {
		fmt.Println(err)
		portfolio, err = portfolio.SetDefaultPortfolio()
	}
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
