package subcmd

import (
	"fmt"
	"strings"

	"github.com/kmagai/fugo"
	"github.com/kmagai/fugo/cmd/fugo/common"
)

type Add struct {
	Style
}

func (c *Add) Run(args []string) int {
	stockToAdd := args[0]
	portfolio, err := fugo.GetPortfolio()
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}

	added, err := portfolio.AddStock(stockToAdd)
	if err != nil {
		fmt.Println(err)
		return common.ExitCodeError
	}
	// TODO: need better printing
	fmt.Printf("Successfully added: %s", added)
	return common.ExitCodeOK
}

func (c *Add) Synopsis() string {
	return fmt.Sprintf("Add stock to your portfolio")
}

func (c *Add) Help() string {
	helpText := `
	fugo add [CODE]
`
	return strings.TrimSpace(helpText)
}
