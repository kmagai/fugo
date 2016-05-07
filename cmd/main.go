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

	// TODO: 多分最初に前回のデータをどこかに保存しておいてそれを表示→アップデートというのが親切か
	// TODO: とりあえず出力
	for i := 0; i < len(portfolio.Stocks); i++ {
		fmt.Println(<-stocks_ch)
	}

}
