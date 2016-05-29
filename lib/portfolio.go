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

type resourcer interface {
	GetStocks(codes interface{}) (*[]Stock, error)
}

// NewPortfolio is a Portfolio's Factory Method.
func NewPortfolio(path string) *Portfolio {
	return &Portfolio{path: path}
}

// GetPortfolio makes portfolio struct from fugorc.
func (pf *Portfolio) GetPortfolio() error {
	dat, err := ioutil.ReadFile(pf.path)
	if err != nil {
		return errors.New("portfolio file not found")
	}

	return json.Unmarshal(dat, pf)
}

// GetStocks is get from resourcer and returns stock pointer.
func GetStocks(res resourcer, stocks interface{}) (*[]Stock, error) {
	return res.GetStocks(stocks)
}

// Update updates portfolio by Stock Codes.
func (pf *Portfolio) Update(sts *[]Stock) error {
	codeStockMap := make(map[string]Stock)
	for _, s := range *sts {
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
func (pf *Portfolio) RemoveStock(code string) error {
	var otherStocks []Stock
	var removed bool

	for i := range pf.Stocks {
		if pf.Stocks[i].Code == code {
			removed = true
		} else {
			otherStocks = append(otherStocks, pf.Stocks[i])
		}
	}
	if !removed {
		return errors.New("stock not found in your portfolio")
	}
	pf.Stocks = otherStocks
	return pf.saveToFile()
}

// AddStock tries to add stocks to portfolio by the code like 'AAPL', '1234' etc.
func (pf *Portfolio) AddStock(stks *[]Stock) error {
	var err error
	if duplicated := pf.hasDuplicate(stks); duplicated {
		return errors.New("You have already had it in your portfolio")
	}
	pf.Stocks = append(pf.Stocks, *stks...)
	err = pf.saveToFile()
	return err
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
func (pf *Portfolio) hasDuplicate(stks *[]Stock) bool {
	pfMap := make(map[string]Stock)
	for _, s := range pf.Stocks {
		pfMap[s.Code] = s
	}

	for _, s := range *stks {
		if _, found := pfMap[s.Code]; found {
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
