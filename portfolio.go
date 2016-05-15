package fugo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os/user"
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
		err = errors.New("failed to fetch")
		return portfolio, nil
	}

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.New("couldn't properly read response. It could be a problem with a remote host")
		return portfolio, nil
	}

	newStocks, err := ParseToStocks(trimSlashes(stockJSON))
	if err != nil {
		err = errors.New("Couldn't properly read the response. It could be a problem with a remote host")
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
		err = errors.New("Couldn't find any stock. You should check your code")
		return nil, err
	}

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.New("Couldn't properly read response. It could be a problem with a remote host")
		return nil, err
	}

	newStocks, err := ParseToStocks(trimSlashes(stockJSON))
	if err != nil {
		err = errors.New("Couldn't properly read the response. It could be a problem with a remote host")
		return nil, err
	}

	if duplicated := portfolio.hasDuplicate(newStocks); duplicated {
		err = errors.New("You have already had it in your portfolio")
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
