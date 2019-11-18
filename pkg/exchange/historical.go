package exchange

import (
	"golang.org/x/text/currency"
)

// GetHistorical internal func to retrieve the foreign exchange reference rates for a specific date
func (c Client) GetHistorical(t string) (map[currency.Unit]float64, error) {
	return c.getHistorical(t, nil)
}

// GetHistoricalWithParams internal func to retrieve the foreign exchange reference rates for a specific date
func (c Client) GetHistoricalWithParams(p *Params, t string) (map[currency.Unit]float64, error) {
	return c.getHistorical(t, setParams(p))
}

// getHistorical internal func to retrieve the foreign exchange reference rates for a specific date
func (c Client) getHistorical(t string, p map[string]string) (map[currency.Unit]float64, error) {
	c.logger.Debug().Str("day", t).Msg("getting historical rate...")
	u := c.buildURL("/"+t, p)
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
