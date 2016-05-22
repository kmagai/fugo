package fugo

import (
	"log"
	"os"
	"os/user"
	"testing"
)

var portfolio *Portfolio

const testFugorc = `/.tfugorc`

func setDefaultPortfolio() *Portfolio {
	usr, _ := user.Current()
	portfolio = &Portfolio{}
	portfolio.Path = usr.HomeDir + testFugorc
	portfolio, _ = portfolio.SetDefaultPortfolio()
	defer os.Remove(portfolio.Path)
	return portfolio
}

func TestGetPortfolio(t *testing.T) {
	usr, err := user.Current()
	portfolio := &Portfolio{}
	portfolio.Path = usr.HomeDir + testFugorc
	defer os.Remove(portfolio.Path)
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
		t.Errorf("expected updated portfolio to have %i stocks but got %i", defaultPortfolioStockLength-1, len(portfolio.Stocks))
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
		t.Errorf("expected updated portfolio to have %i stocks but got %i", defaultPortfolioStockLength+1, len(portfolio.Stocks))
	}
}
