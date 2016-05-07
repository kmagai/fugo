package fugo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"
)

const fugorc = `/.fugorc`

type Portfolio struct {
	Stocks []Stock
}

func GetPortfolio() *Portfolio {
	portfolio := &Portfolio{}
	dat, err := ioutil.ReadFile(portfolio.fileName())
	if err != nil {
		portfolio := portfolio.defaultPortfolio()
		portfolio.saveToFile()
	} else {
		json.Unmarshal(dat, portfolio)
	}

	return portfolio
}

// StockPriceUpdate updates stock price using Google Finance API
func (portfolio *Portfolio) Update() <-chan string {
	stocksChan := make(chan string)
	for _, stock := range portfolio.Stocks {
		go func(stock Stock) {
			fmt.Println(stock.FetchURL())
			res, err := http.Get(stock.FetchURL())
			if err != nil {
				fmt.Println("Failed to Fetch: " + err.Error())
			}
			dat, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Failed to Fetch: " + err.Error())
				return
			}

			stocksChan <- string(dat[:])
		}(stock)
	}
	return stocksChan
}

func (portfolio *Portfolio) fileName() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir + fugorc
}

func (portfolio *Portfolio) saveToFile() {
	dat, err := json.Marshal(portfolio)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(portfolio.fileName(), dat, 0644)
	return
}

// defaultPortfolio stock's are selected and ordered by market capitalization
func (portfolio *Portfolio) defaultPortfolio() *Portfolio {
	portfolio.Stocks = []Stock{
		{Market: "INDEXNIKKEI", Code: "NI225"}, // 日経平均
		{Market: "TYO", Code: "7203"},          // トヨタ自動車(株)
		{Market: "TYO", Code: "9437"},          // (株)NTTドコモ
		{Market: "TYO", Code: "9432"},          // 日本電信電話(株)
		{Market: "TYO", Code: "2914"},          // JT
		{Market: "TYO", Code: "9433"},          // KDDI(株)
		{Market: "TYO", Code: "8306"},          // (株)三菱UFJフィナンシャル・グループ
		{Market: "TYO", Code: "9984"},          // ソフトバンク
		{Market: "TYO", Code: "6178"},          // 日本郵政(株)
		{Market: "TYO", Code: "7182"},          // (株)ゆうちょ銀行
		{Market: "TYO", Code: "7267"},          // ホンダ
	}
	return portfolio
}
