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

// Add is the command name.
type Add struct {
	Style
}

// Run specifies what this command does.
func (c *Add) Run(args []string) int {
	stockToAdd := args[0]

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

	// TODO: make other APIs available
	var api interfaces.Resourcer = googleFinance.API{}
	newStocks, err := api.GetStocker(stockToAdd)
	fmt.Println(newStocks)

	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}

	// TODO: prompt to confirm
	err = pf.AddStock(stockToAdd)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	return common.ExitCodeOK
}

// Synopsis tells what it does.
func (c *Add) Synopsis() string {
	return fmt.Sprintf("Add stock to your portfolio")
}

// Help text.
func (c *Add) Help() string {
	helpText := `
	fugo add [CODE]
`
	return strings.TrimSpace(helpText)
}
