package fugo

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"
)

type Stock struct {
	Code          string    `json:"t"` // Depending on the market (f.g. integer code for TYO, ticker for NASDAQ...etc)
	Name          string    `json:"name"`
	Price         float64   `json:"l_fix,string"`
	Change        float64   `json:"c_fix,string"`
	ChangePercent float64   `json:"cp_fix,string"`
	UpdatedAt     time.Time `json:"lt_dts,string"`
}

// parse JSON from the API
func ParseToStocks(stockJson []byte) (*[]Stock, error) {
	s := bytes.NewReader(stockJson)
	var newStockData *[]Stock
	dec := json.NewDecoder(s)
	dec.Decode(&newStockData)
	if newStockData == nil {
		return nil, errors.New("failed to parse stock")
	} else {
		return newStockData, nil
	}
}
