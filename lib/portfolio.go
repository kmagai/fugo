package fugo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Fugorc is file name to specify fugo user setting file
const Fugorc = `/.fugorc`

// Portfolio stores user configurable portfolio
type Portfolio struct {
	Stocks []Stock
	path   string
}

// Resource interface is for switching data resources
type Resource interface {
	GetStock(interface{}) ([]byte, error)
}

// NewPortfolio is a Portfolio's Factory Method
func NewPortfolio(path string) *Portfolio {
	return &Portfolio{path: path}
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

// GetStock is get from resource and returns stock pointer
func GetStock(resource Resource, stocks interface{}) (*[]Stock, error) {
	stockJSON, err := resource.GetStock(stocks)
	if err != nil {
		return nil, errors.New("Failed to fetch from remote")
	}
	stock, err := ParseToStocks(stockJSON)
	return stock, err
}

// Update updates portfolio by Stock Code
func (portfolio *Portfolio) Update(updatedStock *[]Stock) (*Portfolio, error) {
	codeStockMap := make(map[string]Stock)
	for _, s := range *updatedStock {
		codeStockMap[s.Code] = s
	}

	for i := range portfolio.Stocks {
		if stock, ok := codeStockMap[portfolio.Stocks[i].Code]; ok {
			portfolio.Stocks[i] = stock
		}
	}
	
	err := portfolio.saveToFile()
	return portfolio, err
}

// RemoveStock tries to removes stock from portfolio by the code like 'AAPL', '1234' etc
func (portfolio *Portfolio) RemoveStock(codeToRemove string) (*Stock, error) {
	var newStocks []Stock
	var removedStock *Stock
	var err error

	for i := range portfolio.Stocks {
		if portfolio.Stocks[i].Code == codeToRemove {
			removedStock = &portfolio.Stocks[i]
		} else {
			newStocks = append(newStocks, portfolio.Stocks[i])
		}
	}
	if removedStock == nil {
		return removedStock, errors.New("stock not found in your portfolio")
	}
	portfolio.Stocks = newStocks
	err = portfolio.saveToFile()
	return removedStock, err
}

// AddStock tries to add stocks to portfolio by the code like 'AAPL', '1234' etc
func (portfolio *Portfolio) AddStock(stocks *[]Stock) (*[]Stock, error) {
	var err error
	if duplicated := portfolio.hasDuplicate(stocks); duplicated {
		return nil, errors.New("You have already had it in your portfolio")
	}
	portfolio.Stocks = append(portfolio.Stocks, *stocks...)
	err = portfolio.saveToFile()
	return stocks, err
}

// SetDefaultPortfolio stock's are selected arbitrary
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
