package subcmd

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/kmagai/fugo"
	"github.com/kmagai/fugo/cmd/common"
)

// Remove is the command name
type Remove struct {
	Style
}

// Run specifies what this command does
func (c *Remove) Run(args []string) int {
	stockToRemove := args[0]
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
	removed, err := portfolio.RemoveStock(stockToRemove)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	// TODO: need better printing
	fmt.Printf("removed: %p", removed)
	return common.ExitCodeOK
}

// Synopsis tells what it does
func (c *Remove) Synopsis() string {
	return "Remove stock from your portfolio"
}

// Help text
func (c *Remove) Help() string {
	helpText := `
	fugo remove [CODE]
`
	return strings.TrimSpace(helpText)
}
