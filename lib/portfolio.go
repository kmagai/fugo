package fugo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Fugorc is file name to specify fugo user setting file.
const Fugorc = "/.fugorc"

// Portfolio stores user configurable portfolio.
type Portfolio struct {
	Codes []string
	path  string
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

// RemoveStock tries to removes stock from portfolio by the code like 'AAPL', '1234' etc.
func (pf *Portfolio) RemoveStock(code string) error {
	var removed bool
	var codes []string

	for _, cd := range pf.Codes {
		if cd == code {
			removed = true
		} else {
			codes = append(codes, cd)
		}
	}
	if !removed {
		return errors.New("stock not found in your portfolio")
	}
	pf.Codes = codes

	return pf.saveToFile()
}

// AddStock tries to add a stock to portfolio by the code like 'AAPL', '1234' etc.
func (pf *Portfolio) AddStock(code string) error {
	var err error
	if duplicated := pf.hasDuplicate(code); duplicated {
		return errors.New("You have alread had it in your portfolio")
	}
	pf.Codes = append(pf.Codes, code)

	err = pf.saveToFile()
	return err
}

// AddStocks tries to add stocks to portfolio by the code like 'AAPL', '1234' etc.
func (pf *Portfolio) AddStocks(codes []string) error {
	var err error
	if duplicated := pf.hasDuplicate(codes); duplicated {
		return errors.New("You have already had it in your portfolio")
	}
	pf.Codes = append(pf.Codes, codes...)

	err = pf.saveToFile()
	return err
}

// SetDefaultPortfolio stock's are selected arbitrary.
func (pf *Portfolio) SetDefaultPortfolio() (*Portfolio, error) {
	defaultCodes := []string{"NI225", "7203", "9984", "6178", "AAPL"}
	pf.Codes = defaultCodes
	err := pf.saveToFile()
	return pf, err
}

// hasDuplicate return true if portfolio has any stock.
func (pf *Portfolio) hasDuplicate(codes interface{}) bool {
	codeMap := make(map[string]bool)
	for _, c := range pf.Codes {
		codeMap[c] = true
	}
	fmt.Println(codeMap)
	switch v := codes.(type) {
	case string:
		if found := codeMap[v]; found {
			return true
		}
	case []string:
		for _, s := range v {
			if found := codeMap[s]; found {
				return true
			}
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
	return ioutil.WriteFile(pf.path, dat, 0644)
}
