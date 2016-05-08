package fugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"
)

const fugorc = `/.fugorc`
const googleFinanceAPI = "http://www.google.com/finance/info?infotype=infoquoteall&q=%s"

type Portfolio struct {
	Stocks []Stock
}

// TODO: add portfolio from CLI
// func AddStock()  {
// }

// TODO: remove from CLI
// func RemoveStock()  {
// }

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

func (portfolio *Portfolio) fileName() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir + fugorc
}

func (portfolio *Portfolio) Update() *Portfolio {
	stockJSON := portfolio.fetch()
	newStocks := portfolio.parseToStock(stockJSON)

	var newPortfolio Portfolio
	codeStockMap := make(map[string]Stock)
	for _, newStock := range *newStocks {
		codeStockMap[newStock.Code] = newStock
	}

	for _, currentStock := range portfolio.Stocks {
		if newStock, found := codeStockMap[currentStock.Code]; found {
			newPortfolio.Stocks = append(newPortfolio.Stocks, newStock)
		} else {
			// make and return custom err
			fmt.Println("Not found in remote")
			fmt.Println("Misconfigured?")
		}
	}
	newPortfolio.saveToFile()
	return &newPortfolio
}

// Fetches stock price using Google Finance API
func (portfolio *Portfolio) fetch() []byte {
	res, err := http.Get(portfolio.buildFetchURL())
	if err != nil {
		fmt.Println("Failed to Fetch: " + err.Error())
	}
	dat, _ := ioutil.ReadAll(res.Body)
	return dat
}

func (portfolio *Portfolio) parseToStock(stockJson []byte) *[]Stock {
	// Parse JSON from remote
	stockJsonString := string(stockJson)[3:] // skip '//' chars
	s := bytes.NewReader([]byte(stockJsonString))
	var newStockData *[]Stock
	dec := json.NewDecoder(s)
	dec.Decode(&newStockData)
	return newStockData
}

func (portfolio *Portfolio) saveToFile() {
	dat, err := json.Marshal(portfolio)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(portfolio.fileName(), dat, 0644)
	return
}

func (portfolio *Portfolio) buildFetchURL() string {
	var stockCodes string
	for _, stock := range portfolio.Stocks {
		stockCodes += stock.Code + ","
	}
	return fmt.Sprintf(googleFinanceAPI, stockCodes)
}

// defaultPortfolio stock's are selected and ordered by market capitalization
func (portfolio *Portfolio) defaultPortfolio() *Portfolio {
	portfolio.Stocks = []Stock{
		{Code: "AAPL"},  // Apple Inc.
		{Code: "NI225"}, // 日経平均
		{Code: "7203"},  // トヨタ自動車(株)
		{Code: "9437"},  // (株)NTTドコモ
		{Code: "9432"},  // 日本電信電話(株)
		{Code: "2914"},  // JT
		{Code: "9433"},  // KDDI(株)
		{Code: "8306"},  // (株)三菱UFJフィナンシャル・グループ
		{Code: "9984"},  // ソフトバンク
		{Code: "6178"},  // 日本郵政(株)
		{Code: "7182"},  // (株)ゆうちょ銀行
		{Code: "7267"},  // ホンダ
	}
	return portfolio
}
