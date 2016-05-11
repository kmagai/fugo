// fugo, meaning a rich man in Japanese, is yet another stock portfolio manager.
package main

import "github.com/kmagai/fugo"

func main() {
	portfolio := fugo.GetPortfolio()
	portfolio = portfolio.Update()
	portfolio.PrintStocks()
}
