package fugo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Fugorc is file name to specify fugo user setting file.
const Fugorc = "/.fugorc"

// Portfolio stores user configurable portfolio.
type Portfolio struct {
	Stocks []Stock
	path   string
}

type resource interface {
	GetStocks(codes interface{}) (*[]Stock, error)
}

// NewPortfolio is a Portfolio's Factory Method.
func NewPortfolio(path string) *Portfolio {
	return &Portfolio{path: path}
}

// GetPortfolio makes portfolio struct from fugorc.
func (pf *Portfolio) GetPortfolio() (*Portfolio, error) {
	dat, err := ioutil.ReadFile(pf.path)
	if err != nil {
		return pf, errors.New("portfolio file not found")
	}

	err = json.Unmarshal(dat, pf)
	return pf, err
}

// GetStocks is get from resource and returns stock pointer.
func GetStocks(res resource, stocks interface{}) (*[]Stock, error) {
	return res.GetStocks(stocks)
}

// Update updates portfolio by Stock Code.
func (pf *Portfolio) Update(updatedStock *[]Stock) error {
	codeStockMap := make(map[string]Stock)
	for _, s := range *updatedStock {
		codeStockMap[s.Code] = s
	}

	for i := range pf.Stocks {
		if stock, ok := codeStockMap[pf.Stocks[i].Code]; ok {
			pf.Stocks[i] = stock
		}
	}

	err := pf.saveToFile()
	return err
}

// RemoveStock tries to removes stock from portfolio by the code like 'AAPL', '1234' etc.
func (pf *Portfolio) RemoveStock(codeToRemove string) (removedStock *Stock, err error) {
	var otherStocks []Stock

	for i := range pf.Stocks {
		if pf.Stocks[i].Code == codeToRemove {
			removedStock = &pf.Stocks[i]
		} else {
			otherStocks = append(otherStocks, pf.Stocks[i])
		}
	}
	if removedStock == nil {
		return removedStock, errors.New("stock not found in your portfolio")
	}
	pf.Stocks = otherStocks
	err = pf.saveToFile()
	return removedStock, err
}

// AddStock tries to add stocks to portfolio by the code like 'AAPL', '1234' etc.
func (pf *Portfolio) AddStock(stocks *[]Stock) (*[]Stock, error) {
	var err error
	if duplicated := pf.hasDuplicate(stocks); duplicated {
		return nil, errors.New("You have already had it in your portfolio")
	}
	pf.Stocks = append(pf.Stocks, *stocks...)
	err = pf.saveToFile()
	return stocks, err
}

// SetDefaultPortfolio stock's are selected arbitrary.
func (pf *Portfolio) SetDefaultPortfolio() (*Portfolio, error) {
	pf.Stocks = []Stock{
		{Code: "NI225"}, // 日経平均
		{Code: "7203"},  // トヨタ自動車(株)
		{Code: "9984"},  // ソフトバンク
		{Code: "6178"},  // 日本郵政(株)
		{Code: "AAPL"},  // Apple Inc.
	}
	err := pf.saveToFile()
	return pf, err
}

// hasDuplicate return true if portfolio has any stock.
func (pf *Portfolio) hasDuplicate(stocks *[]Stock) bool {
	portfolioMap := make(map[string]Stock)
	for _, s := range pf.Stocks {
		portfolioMap[s.Code] = s
	}

	for _, s := range *stocks {
		if _, found := portfolioMap[s.Code]; found {
			return true
		}
	}
	return false
}

// saveToFile saves portfolio struct into fugorc.
func (pf *Portfolio) saveToFile() error {
	dat, err := json.Marshal(pf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(pf.path, dat, 0644)
	return err
}
