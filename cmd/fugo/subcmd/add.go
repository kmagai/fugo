package subcmd

import (
	"fmt"
	"strings"

	"github.com/kmagai/fugo"
)

type Add struct {
	Style
}

func (c *Add) Run(args []string) int {
	stockToAdd := args[0]
	portfolio, err := fugo.GetPortfolio()
	if err != nil {
		return 1
	}
	added, err := portfolio.AddStock(stockToAdd)
	if err != nil || added == nil {
		fmt.Println("Failed to add the stock")
		return 1
	}
	// TODO: need better printing
	fmt.Printf("Successfully added: %s", added)
	return 0
}

func (c *Add) Synopsis() string {
	return fmt.Sprintf("Add stock to your portfolio")
}

func (c *Add) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
