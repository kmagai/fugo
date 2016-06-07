package subcmd

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/kmagai/fugo/common"
	"github.com/kmagai/fugo/lib"
)

// Remove is the command name.
type Remove struct {
	Style
}

// Run specifies what this command does.
func (c *Remove) Run(args []string) int {
	stockToRemove := args[0]
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
	}
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	err = pf.RemoveStock(stockToRemove)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}

	fmt.Printf("removed: %s", stockToRemove)
	return common.ExitCodeOK
}

// Synopsis tells what it does.
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
