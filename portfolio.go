package fugo

import (
	"encoding/json"
	"errors"
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
	res, err := http.Get(buildFetchURL(buildQuery(portfolio.Stocks)))
	if err != nil {
		fmt.Println("Failed to Fetch: " + err.Error())
	}

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't properly read response. It could be a problem with a remote host: " + err.Error())
	}

	stockJSON = []byte(string(stockJSON)[3:]) // trim '//'
	newStocks := ParseToStocks(stockJSON)

	var newPortfolio Portfolio
	codeStockMap := make(map[string]Stock)
	for _, newStock := range *newStocks {
		codeStockMap[newStock.Code] = newStock
	}

	for _, currentStock := range portfolio.Stocks {
		if newStock, ok := codeStockMap[currentStock.Code]; ok {
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

// func (portfolio *Portfolio) RemoveStock(codeToRemove string) *Portfolio, error {
func (portfolio *Portfolio) RemoveStock(codeToRemove string) (*Stock, error) {
	var newPortfolio Portfolio
	var removedStock *Stock
	var err error

	for _, stock := range portfolio.Stocks {
		if stock.Code == codeToRemove {
			removedStock = &stock
		} else {
			newPortfolio.Stocks = append(newPortfolio.Stocks, stock)
		}
	}
	if removedStock == nil {
		err = errors.New("stock not found in your portfolio")
		return removedStock, err
	}
	newPortfolio.saveToFile()

	return removedStock, err
}

func (portfolio *Portfolio) AddStock(codeToAdd string) (*[]Stock, error) {
	var newPortfolio Portfolio
	var addedStocks *[]Stock
	var err error

	res, err := http.Get(buildFetchURL(codeToAdd))
	if err != nil {
		fmt.Println("Couldn't find any stock. You should check your code: " + err.Error())
		return addedStocks, err
	}

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't properly read response. It could be a problem with a remote host: " + err.Error())
		return addedStocks, err
	}
	fmt.Println(string(stockJSON))

	stockJSON = []byte(string(stockJSON)[3:]) // trim '//'
	addedStocks = ParseToStocks(stockJSON)
	if addedStocks != nil {
		newPortfolio.Stocks = append(portfolio.Stocks, *addedStocks...)
		newPortfolio.saveToFile()
	}

	return addedStocks, err
}

func (portfolio *Portfolio) saveToFile() {
	dat, err := json.Marshal(portfolio)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(portfolio.fileName(), dat, 0644)
	return
}

func buildFetchURL(query string) string {
	return fmt.Sprintf(googleFinanceAPI, query)
}

func buildQuery(stocks []Stock) string {
	var query string
	for _, stock := range stocks {
		query += stock.Code + ","
	}
	return query
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

// TODO: 酷い
func (portfolio *Portfolio) PrintStocks() {
	if len(portfolio.Stocks) > 0 {
		for _, stock := range portfolio.Stocks {
			fmt.Println("------------------------------------")
			headColor := color.New(color.FgGreen, color.Underline)
			priceColor := color.New(color.FgBlack, color.BgHiWhite, color.Bold)
			priceChangeColor := color.New(color.FgBlack, color.Bold)
			headColor.Printf(stock.Name+" (%s)\n", stock.Code)
			const layout = "2006-01-02 15:04:05"
			fmt.Println(stock.UpdatedAt.Format(layout))
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
				priceChangeColor.Print("(")
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
