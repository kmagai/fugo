package subcmd

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/kmagai/fugo"
	"github.com/kmagai/fugo/cmd/common"
)

// Add is the command name
type Add struct {
	Style
}

// Run specifies what this command does
func (c *Add) Run(args []string) int {
	stockToAdd := args[0]

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

	added, err := portfolio.AddStock(stockToAdd)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	fmt.Printf("Successfully added: %p", added)
	return common.ExitCodeOK
}

// Synopsis tells what it does
func (c *Add) Synopsis() string {
	return fmt.Sprintf("Add stock to your portfolio")
}

// Help text
func (c *Add) Help() string {
	helpText := `
	fugo add [CODE]
`
	return strings.TrimSpace(helpText)
}