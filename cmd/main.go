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
	stocks_ch := portfolio.Update()

	for i := 0; i < len(portfolio.Stocks); i++ {
		fmt.Println(<-stocks_ch)
	}

}
