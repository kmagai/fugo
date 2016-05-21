package fugo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const Fugorc = `/.fugorc`

type Portfolio struct {
	Stocks []Stock
	path   string
}

// SetPortfolioFilePath returns portfolio setting file path
func (portfolio *Portfolio) SetPortfolioFilePath(dirname string, filename string) *Portfolio {
	portfolio.path = dirname + filename

	return portfolio
}

// GetPortfolio makes portfolio struct from fugorc
func (portfolio *Portfolio) GetPortfolio() (*Portfolio, error) {
	dat, err := ioutil.ReadFile(portfolio.path)
	if err != nil {
		return portfolio, errors.New("portfolio file not found")
	}

	err = json.Unmarshal(dat, portfolio)
	return portfolio, err
}

// Update updates portfolio with remote API
func (portfolio *Portfolio) Update() (*Portfolio, error) {
	newStocks, err := getRemoteStock(portfolio.Stocks)
	if err != nil {
		return portfolio, errors.New("Couldn't properly read the response. It could be a problem with a remote host")
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

// RemoveStock tries to removes stock from portfolio by the code like 'AAPL', '1234' etc
func (portfolio *Portfolio) RemoveStock(codeToRemove string) (*Stock, error) {
	var newPortfolio Portfolio
	var removedStock *Stock
	var err error

	for i, _ := range portfolio.Stocks {
		if portfolio.Stocks[i].Code == codeToRemove {
			removedStock = &portfolio.Stocks[i]
		} else {
			newPortfolio.Stocks = append(newPortfolio.Stocks, portfolio.Stocks[i])
		}
	}
	if removedStock == nil {
		return removedStock, errors.New("stock not found in your portfolio")
	}
	newPortfolio.path = portfolio.path
	err = newPortfolio.saveToFile()
	return removedStock, err
}

// AddStock tries to add stocks to portfolio by the code like 'AAPL', '1234' etc
func (portfolio *Portfolio) AddStock(codeToAdd string) (*[]Stock, error) {
	var newPortfolio Portfolio
	var err error

	newStocks, err := getRemoteStock(codeToAdd)
	if err != nil {
		return newStocks, errors.New("Couldn't properly read the response. It could be a problem with either the remote host or your typo")
	}

	if duplicated := portfolio.hasDuplicate(newStocks); duplicated {
		return nil, errors.New("You have already had it in your portfolio")
	}
	newPortfolio.path = portfolio.path
	newPortfolio.Stocks = append(portfolio.Stocks, *newStocks...)
	err = newPortfolio.saveToFile()
	return newStocks, err
}

// defaultPortfolio stock's are selected arbitrary
func (portfolio *Portfolio) SetDefaultPortfolio() (*Portfolio, error) {
	portfolio.Stocks = []Stock{
		{Code: "NI225"}, // 日経平均
		{Code: "7203"},  // トヨタ自動車(株)
		{Code: "9984"},  // ソフトバンク
		{Code: "6178"},  // 日本郵政(株)
		{Code: "AAPL"},  // Apple Inc.
	}
	err := portfolio.saveToFile()
	return portfolio, err
}

// getRemoteStock gets stock struct from remote
func getRemoteStock(stocks interface{}) (*[]Stock, error) {
	var query string
	switch s := stocks.(type) {
	case string:
		query = s
	case []Stock:
		query = buildQuery(s)
	}

	res, err := http.Get(buildFetchURL(query))
	if err != nil {
		return nil, errors.New("failed to fetch")
	}

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("couldn't properly read response. It could be a problem with a remote host")
	}

	return ParseToStocks(trimSlashes(stockJSON))
}

// hasDuplicate return true if portfolio has any stock
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

// saveToFile saves portfolio struct into fugorc
func (portfolio *Portfolio) saveToFile() error {
	dat, err := json.Marshal(portfolio)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(portfolio.path, dat, 0644)
	return err
}
