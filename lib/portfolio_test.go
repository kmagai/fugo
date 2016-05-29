package fugo

import (
	"log"
	"os/user"
	"reflect"
	"testing"
	"time"
)

var portfolio *Portfolio

var usr, _ = user.Current()
var path = usr.HomeDir + `/.tfugorc`

func setDefaultPortfolio() *Portfolio {
	portfolio := NewPortfolio(path)
	portfolio.Stocks = []Stock{
		Stock{Code: "TEST1", Name: "TEST1 Inc.", Price: 1000, Change: 500, ChangePercent: 100, UpdatedAt: time.Now()},
		Stock{Code: "TEST2", Name: "TEST2 Inc.", Price: 1500, Change: 1000, ChangePercent: 50, UpdatedAt: time.Now()},
		Stock{Code: "TEST3", Name: "TEST3 Inc.", Price: 5000, Change: 4000, ChangePercent: 25, UpdatedAt: time.Now()},
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

func TestUpdate(t *testing.T) {
	portfolio = setDefaultPortfolio()

	targetCode := portfolio.Stocks[1].Code
	newName := "UPDATED Inc."
	newPrice := 1000.0
	stock := &[]Stock{
		Stock{Code: targetCode, Name: newName, Price: newPrice, Change: 500, ChangePercent: 100, UpdatedAt: time.Now()},
	}
	beforeUpdateStockLen := len(portfolio.Stocks)
	err := portfolio.Update(stock)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	if beforeUpdateStockLen != len(portfolio.Stocks) {
		t.Errorf("got greater stock data than expected")
	}
	for _, stock := range portfolio.Stocks {
		if stock.Code == targetCode {
			if stock.Name != newName {
				t.Errorf("Name seems not to be updated")
			}
			if stock.Price != newPrice {
				t.Errorf("Name seems not to be updated")
			}
		}
	}
}

func TestRemoveStock(t *testing.T) {
	portfolio = setDefaultPortfolio()
	portfolio.GetPortfolio()
	defaultPortfolioStockLength := len(portfolio.Stocks)
	first := portfolio.Stocks[0]
	err := portfolio.RemoveStock(first.Code)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	if defaultPortfolioStockLength-1 != len(portfolio.Stocks) {
		t.Errorf("expected updated portfolio to have %d stocks but got %d", defaultPortfolioStockLength-1, len(portfolio.Stocks))
	}
}

func TestAddStock(t *testing.T) {
	portfolio = setDefaultPortfolio()
	portfolio.GetPortfolio()
	defaultPortfolioStockLength := len(portfolio.Stocks)

	stock := &[]Stock{
		Stock{Code: "ADD", Name: "ADD Inc.", Price: 1000, Change: 500, ChangePercent: 100, UpdatedAt: time.Now()},
	}

	err := portfolio.AddStock(stock)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	if defaultPortfolioStockLength+1 != len(portfolio.Stocks) {
		t.Errorf("expected updated portfolio to have %d stocks but got %d", defaultPortfolioStockLength+1, len(portfolio.Stocks))
	}
	setDefaultPortfolio()
}
