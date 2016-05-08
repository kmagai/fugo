/*
Fugo(富豪), meaning a rich man in Japanese, is a stock portfolio manager for japanese stock market.
*/
package main

import (
	"fmt"

	"github.com/kmagai/fugo"
)

func main() {
	portfolio := fugo.GetPortfolio()
	portfolio = portfolio.Update()
	for _, stock := range portfolio.Stocks {
		// TODO: show formatted
		fmt.Println(stock)
	}
}
