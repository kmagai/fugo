package common

import (
	"math"

	"github.com/fatih/color"
	"github.com/kmagai/fugo"
)

var (
	defaultColor     = color.New(color.FgBlack)
	headColor        = color.New(color.FgGreen, color.Underline)
	priceColor       = color.New(color.FgBlack, color.BgHiWhite, color.Bold)
	priceChangeColor = color.New(color.FgBlack, color.Bold)
)

// ShowPortfolio prints stock data in portfolio
func ShowPortfolio(portfolio *fugo.Portfolio) {
	if len(portfolio.Stocks) > 0 {
		for _, stock := range portfolio.Stocks {
			printStock(stock)
			defaultColor.Println()
			defaultColor.Println()
		}
	} else {
		defaultColor.Println("No stock data found")
		defaultColor.Println("You may need to add stock to your portfolio")
		defaultColor.Println("add it like `fugo add [CODE]`")
	}
}

func printStock(stock fugo.Stock) {
	headColor.Printf(stock.Name+" (%s)\n", stock.Code)
	const layout = "2006-01-02 15:04:05"
	defaultColor.Println(stock.UpdatedAt.Format(layout))
	priceColor.Print(stock.Price)
	if stock.Change > 0 {
		printPriceUp(stock)
	} else {
		printPriceDown(stock)
	}
}

func printPriceUp(stock fugo.Stock) {
	priceChangeColor.Print(color.RedString("↑") + " ")
	priceChangeColor.Print(roundAt(stock.Change, 2))
	priceChangeColor.Print(" (")
	priceChangeColor.Print(roundAt(stock.ChangePercent, 2))
	priceChangeColor.Println("%)")
}

func printPriceDown(stock fugo.Stock) {
	priceChangeColor.Print(color.BlueString("↓") + " ")
	priceChangeColor.Print(roundAt(stock.Change, 2))
	priceChangeColor.Print("(")
	priceChangeColor.Print(roundAt(stock.ChangePercent, 2))
	priceChangeColor.Println("%)")
}

func roundAt(f float64, roundAt int) float64 {
	shift := math.Pow(10, float64(roundAt))
	return round(f*shift) / shift
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
