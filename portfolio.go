package fugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os/user"

	"github.com/fatih/color"
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
		{Code: "7606"},  // ユナイテッドアローズ
	}
	return portfolio
}

func (portfolio *Portfolio) PrintStocks() {
	if len(portfolio.Stocks) > 0 {
		for _, stock := range portfolio.Stocks {
			fmt.Println("------------------------------------")
			headColor := color.New(color.FgGreen, color.Underline)
			priceColor := color.New(color.FgBlack, color.BgHiWhite, color.Bold)
			priceChangeColor := color.New(color.FgBlack, color.Bold)
			headColor.Printf(stock.Name+" (%s)\n", stock.Code)
			fmt.Println(stock.UpdatedAt)
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
				priceChangeColor.Print(" (")
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
