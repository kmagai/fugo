package adapter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/kmagai/fugo"
)

const googleURL = "http://www.google.com/finance/info?infotype=infoquoteall&q=%s"

type GoogleFinanceAPI struct {
}

// NewGoogleAPI is a factory method for googleFinanceAPI
func NewGoogleAPI() GoogleFinanceAPI {
	return GoogleFinanceAPI{}
}

// GetStock gets stock struct from google API
func (g GoogleFinanceAPI) GetStock(stocks interface{}) ([]byte, error) {
	var query string
	switch s := stocks.(type) {
	case string:
		query = s
	case []fugo.Stock:
		query = buildQuery(s)
	}

	res, err := http.Get(buildFetchURL(query))
	if err != nil {
		return []byte{}, errors.New("failed to fetch")
	}
	defer res.Body.Close()

	stockJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, errors.New("couldn't properly read response. It could be a problem with a remote host")
	}

	return trimSlashes(stockJSON), nil
}

// trimSlashes trims useless slashes in Google Finance API response
func trimSlashes(json []byte) []byte {
	return []byte(string(json)[3:])
}

// buildFetchURL builds Google Finance API url with the query specified
func buildFetchURL(query string) string {
	return fmt.Sprintf(googleURL, query)
}

// buildQuery builds query specifically for Google Finance API with the stock
func buildQuery(stocks []fugo.Stock) string {
	var query string
	for _, stock := range stocks {
		query += stock.Code + ","
	}
	return query
}
