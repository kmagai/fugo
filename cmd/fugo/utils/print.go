package utils

import (
	"fmt"
	"math"

	"github.com/fatih/color"
	"github.com/kmagai/fugo"
)

func PrintStocks(portfolio *fugo.Portfolio) {
	if len(portfolio.Stocks) > 0 {
		for _, stock := range portfolio.Stocks {
			fmt.Println("------------------------------------")
			headColor := color.New(color.FgGreen, color.Underline)
			priceColor := color.New(color.FgBlack, color.BgHiWhite, color.Bold)
			priceChangeColor := color.New(color.FgBlack, color.Bold)
			headColor.Printf(stock.Name+" (%s)\n", stock.Code)
			const layout = "2006-01-02 15:04:05"
			fmt.Println(stock.UpdatedAt.Format(layout))
			priceColor.Print(stock.Price)

			if stock.Change > 0 {
				priceChangeColor.Print(color.RedString("↑") + " ")
				priceChangeColor.Print(roundAt(stock.Change, 2))
				priceChangeColor.Print(" (")
				priceChangeColor.Print(roundAt(stock.ChangePercent, 2))
				priceChangeColor.Println("%)")
			} else {
				priceChangeColor.Print(color.BlueString("↓") + " ")
				priceChangeColor.Print(roundAt(stock.Change, 2))
				priceChangeColor.Print("(")
				priceChangeColor.Print(roundAt(stock.ChangePercent, 2))
				priceChangeColor.Println("%)")
			}
		}
		fmt.Println("------------------------------------")
	} else {
		fmt.Println("Nothing to print!")
		fmt.Println("TODO: SHOW HELP MESSAGE HERE")
	}
}

func roundAt(f float64, roundAt int) float64 {
	shift := math.Pow(10, float64(roundAt))
	return round(f*shift) / shift
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
