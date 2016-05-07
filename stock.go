package fugo

import (
	"fmt"
	"strconv"
	"time"
)

// TODO: ちゃんとする
const googleFinanceAPI = "http://www.google.com/finance/getprices?p=1d&f=d,h,o,l,c,v&i=300&x=TYO&q=%s"

type Stock struct {
	Code      int
	Name      string
	Price     int
	UpdatedAt time.Time
}

func (stock *Stock) SetPrice(price int) {
	stock.Price = price
	stock.UpdatedAt = time.Now()
}

func (stock *Stock) FetchURL() string {
	return fmt.Sprintf(googleFinanceAPI, strconv.Itoa(stock.Code))
}
