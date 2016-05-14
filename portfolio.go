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

type Portfolio struct {
	Stocks []Stock
}

func GetPortfolio() (*Portfolio, error) {
	portfolio := &Portfolio{}
	dat, err := ioutil.ReadFile(portfolio.fileName())
	if err != nil {
		portfolio := portfolio.defaultPortfolio()
		err = portfolio.saveToFile()
	} else {
		err = json.Unmarshal(dat, portfolio)
	}

	return portfolio, err
}

func (portfolio *Portfolio) Update() (*Portfolio, error) {
	res, err := http.Get(buildFetchURL(buildQuery(portfolio.Stocks)))
	if err != nil {
		fmt.Println("Failed to Fetch: " + err.Error())
		return portfolio, nil
	}

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't properly read response. It could be a problem with a remote host: " + err.Error())
		return portfolio, nil
	}

	newStocks, err := ParseToStocks(trimSlashes(stockJSON))
	if err != nil {
		fmt.Println("couldn't parse it!")
		return portfolio, err
	}

	var newPortfolio Portfolio
	codeStockMap := make(map[string]Stock)
	for _, s := range *newStocks {
		codeStockMap[s.Code] = s
	}

	for _, s := range portfolio.Stocks {
		if newStock, ok := codeStockMap[s.Code]; ok {
			newPortfolio.Stocks = append(newPortfolio.Stocks, newStock)
		} else {
			return portfolio, errors.New("Stock data not found in remote")
		}
	}
	err = newPortfolio.saveToFile()
	return &newPortfolio, nil
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
	err = newPortfolio.saveToFile()

	return removedStock, err
}

func (portfolio *Portfolio) AddStock(codeToAdd string) (*[]Stock, error) {
	var newPortfolio Portfolio
	var err error

	res, err := http.Get(buildFetchURL(codeToAdd))
	if err != nil {
		fmt.Println("Couldn't find any stock. You should check your code: " + err.Error())
		return nil, err
	}

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't properly read response. It could be a problem with a remote host: " + err.Error())
		return nil, err
	}
	fmt.Println(string(stockJSON))

	newStocks, err := ParseToStocks(trimSlashes(stockJSON))
	if err != nil {
		return nil, err
	}

	if duplicated := portfolio.hasDuplicate(newStocks); duplicated {
		fmt.Println("You have already had it in your portfolio")
		return nil, err
	}

	newPortfolio.Stocks = append(portfolio.Stocks, *newStocks...)
	err = newPortfolio.saveToFile()
	return newStocks, err
}

func (portfolio *Portfolio) hasDuplicate(stocks *[]Stock) bool {
	portfolioMap := make(map[string]Stock)
	for _, s := range portfolio.Stocks {
		portfolioMap[s.Code] = s
	}

	for _, s := range *stocks {
		if _, found := portfolioMap[s.Code]; found {
			return true
		}
	}
	return false
}

func (portfolio *Portfolio) saveToFile() (err error) {
	dat, err := json.Marshal(portfolio)
	if err != nil {
		return
	}
	ioutil.WriteFile(portfolio.fileName(), dat, 0644)
	return
}

func (portfolio *Portfolio) fileName() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir + fugorc
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
