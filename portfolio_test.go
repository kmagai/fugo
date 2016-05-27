package fugo

import (
	"log"
	"os"
	"os/user"
	"testing"
)

var portfolio *Portfolio

var usr, _ = user.Current()
var path = usr.HomeDir + `/.tfugorc`

func setDefaultPortfolio() *Portfolio {
	portfolio := NewPortfolio(path)
	portfolio, _ = portfolio.SetDefaultPortfolio()
	defer os.Remove(path)
	return portfolio
}

func TestGetPortfolio(t *testing.T) {
	var err error
	portfolio := NewPortfolio(path)
	defer os.Remove(path)
	portfolio, _ = portfolio.SetDefaultPortfolio()
	portfolio, err = portfolio.GetPortfolio()
	if err != nil {
		t.Errorf("Failed to parse portfolio file to struct")
	}
}

func TestUpdate(t *testing.T) {
	portfolio = setDefaultPortfolio()
	portfolio, _ = portfolio.GetPortfolio()
	updatedPortfolio, err := portfolio.Update()
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	if len(updatedPortfolio.Stocks) > len(portfolio.Stocks) {
		t.Errorf("got greater stock data than expected")
	}
}

func TestRemoveStock(t *testing.T) {
	portfolio = setDefaultPortfolio()
	portfolio, _ = portfolio.GetPortfolio()
	defaultPortfolioStockLength := len(portfolio.Stocks)
	first := portfolio.Stocks[0]
	removedStock, err := portfolio.RemoveStock(first.Code)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	if removedStock == nil {
		t.Errorf("could not removed the stock")
	}
	if defaultPortfolioStockLength-1 != len(portfolio.Stocks) {
		t.Errorf("expected updated portfolio to have %d stocks but got %d", defaultPortfolioStockLength-1, len(portfolio.Stocks))
	}
}

func TestAddStock(t *testing.T) {
	portfolio = setDefaultPortfolio()
	portfolio, _ = portfolio.GetPortfolio()
	defaultPortfolioStockLength := len(portfolio.Stocks)
	addedStock, err := portfolio.AddStock("NFLX") // Netflix
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	if addedStock == nil {
		t.Errorf("could not add the stock")
	}

	if defaultPortfolioStockLength+1 != len(portfolio.Stocks) {
		t.Errorf("expected updated portfolio to have %d stocks but got %d", defaultPortfolioStockLength+1, len(portfolio.Stocks))
	}
}
