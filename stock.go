package fugo

import (
	"fmt"
	"time"
)

// TODO: ちゃんとする
// 現在、東京証券取引所のみ対応
const googleFinanceAPI = "http://www.google.com/finance/info?infotype=infoquoteall&q=%s:%s"

type Stock struct {
	Market       string
	Code         string // Depending on the Market (f.g. Tokyo Stock Exchange provides number for Code)
	Name         string
	Price        int
	ClosingPrice int
	UpdatedAt    time.Time
}

func (stock *Stock) SetPrice(price int) {
	stock.Price = price
	stock.UpdatedAt = time.Now()
}

func (stock *Stock) FetchURL() string {
	return fmt.Sprintf(googleFinanceAPI, stock.Market, stock.Code)
}
