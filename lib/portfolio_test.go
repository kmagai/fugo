package fugo

import (
	"log"
	"os/user"
	"reflect"
	"testing"
)

var portfolio *Portfolio

var usr, _ = user.Current()
var path = usr.HomeDir + `/.tfugorc`

func setDefaultPortfolio() *Portfolio {
	portfolio := NewPortfolio(path)
	portfolio.Codes = []string{
		"TEST1",
		"TEST2",
		"TEST3",
	}
	portfolio.saveToFile()
	return portfolio
}

func TestGetPortfolio(t *testing.T) {
	portfolio := NewPortfolio(path)
	portfolioSaved := setDefaultPortfolio()
	err := portfolio.GetPortfolio()
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	eq := reflect.DeepEqual(portfolioSaved, portfolio)

	if !eq {
		t.Errorf("Portfolio should properly be parsed")
	}
}

func TestRemoveStock(t *testing.T) {
	portfolio = setDefaultPortfolio()
	portfolio.GetPortfolio()
	codesNum := len(portfolio.Codes)
	first := portfolio.Codes[0]
	err := portfolio.RemoveStock(first)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	if codesNum-1 != len(portfolio.Codes) {
		t.Errorf("expected updated portfolio to have %d stocks but got %d", codesNum-1, len(portfolio.Codes))
	}
}

func TestAddStock(t *testing.T) {
	portfolio = setDefaultPortfolio()
	portfolio.GetPortfolio()
	codeNum := len(portfolio.Codes)

	code := "ADD"
	err := portfolio.AddStock(code)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	if codeNum+1 != len(portfolio.Codes) {
		t.Errorf("expected updated portfolio to have %d stocks but got %d", codeNum+1, len(portfolio.Codes))
	}
	setDefaultPortfolio()
}
