package exchange

import (
	"fmt"
	"time"

	"golang.org/x/text/currency"
)

// GetLatest retrieves the latest foreign exchange reference rates
func (c Client) GetLatest() (map[currency.Unit]float64, error) {
	return c.getHistorical(fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()), nil)
}

// GetLatestWithParams retrieves the latest foreign exchange reference rates
func (c Client) GetLatestWithParams(p *Params) (map[currency.Unit]float64, error) {
	return c.getHistorical(fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()), setParams(p))
}
