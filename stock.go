package fugo

import "time"

type Stock struct {
	Code         string    `json:"t"` // Depending on the market (f.g. code for TYO, ticker for NASDAQ...etc)
	Name         string    `json:"name"`
	Price        float64   `json:"l_fix,string"`
	ClosingPrice int       `json:"pcls_fix,string"`
	UpdatedAt    time.Time `json:"lt_dts,string"`
}
