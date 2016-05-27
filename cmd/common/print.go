package common

import (
	"math"

	"fmt"
	"github.com/kmagai/fugo"
	"github.com/olekukonko/tablewriter"
	"os"
)

// ShowPortfolio prints stock data in portfolio
func ShowPortfolio(portfolio *fugo.Portfolio) {
	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CODE", "NAME", "PRICE", "CHANGE", "LAST_UPDATE"})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	const layout = "2006-01-02 15:04:05"
	for _, stock := range portfolio.Stocks {
		table.Append([]string{stock.Code, stock.Name, fmt.Sprint(stock.Price), fmt.Sprint(roundAt(stock.Change, 2)) + " (" + fmt.Sprint(roundAt(stock.ChangePercent, 2)) + "%)", stock.UpdatedAt.Format(layout)})
	}
	table.Render()
	fmt.Println()
}

func roundAt(f float64, roundAt int) float64 {
	shift := math.Pow(10, float64(roundAt))
	return round(f*shift) / shift
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
