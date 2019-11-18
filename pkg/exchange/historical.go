package exchange

import (
	"fmt"
	"time"

	"golang.org/x/text/currency"
)

// GetHistorical retrieves the foreign exchange reference rates for a specific date
func (c Client) GetHistorical(t time.Time) (map[currency.Unit]float64, error) {
	u := c.buildURL(fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day()), nil)
	var response getHistoricalResponse
	if err := c.fetch(u, &response); err != nil {
		return nil, err
	}
	return response.parse(), nil
}

// getHistoricalResponse represents the http response from the GetHistorical http request
type getHistoricalResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// parse returns the mapping of currencies to values from the response of the API
func (r *getHistoricalResponse) parse() map[currency.Unit]float64 {
	var rates = make(map[currency.Unit]float64)
	for s, v := range r.Rates {
		c, err := currency.ParseISO(s)
		if err != nil {
			continue
		}
		rates[c] = v
	}
	return rates
}
