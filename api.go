package fugo

import "fmt"

const googleFinanceAPI = "http://www.google.com/finance/info?infotype=infoquoteall&q=%s"

// trimSlashes trims useless slashes in Google Finance API response
func trimSlashes(json []byte) []byte {
	return []byte(string(json)[3:])
}

// buildFetchURL builds Google Finance API url with the query specified
func buildFetchURL(query string) string {
	return fmt.Sprintf(googleFinanceAPI, query)
}

// buildQuery builds query specifically for Google Finance API with the stock
func buildQuery(stocks []Stock) string {
	var query string
	for _, stock := range stocks {
		query += stock.Code + ","
	}
	return query
}
