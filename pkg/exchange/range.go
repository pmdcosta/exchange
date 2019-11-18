package exchange

import (
	"golang.org/x/text/currency"
)

// GetRange retrieves the foreign exchange reference rates for a specific range of dates
func (c Client) GetRange(start string, end string) (map[string]map[currency.Unit]float64, error) {
	return c.getRange(start, end, nil)
}

// GetRangeWithParams retrieves the foreign exchange reference rates for a specific range of dates
func (c Client) GetRangeWithParams(p *Params, start string, end string) (map[string]map[currency.Unit]float64, error) {
	return c.getRange(start, end, setParams(p))
}

// getRange internal func to retrieve the foreign exchange reference rates for a specific range of dates
func (c Client) getRange(s string, e string, p map[string]string) (map[string]map[currency.Unit]float64, error) {
	c.logger.Debug().Str("start", s).Str("end", e).Msg("getting rates in range...")
	if p == nil {
		p = make(map[string]string)
	}
	p["start_at"] = s
	p["end_at"] = e

	u := c.buildURL(pathHistory, p)
	var response getRangeResponse
	if err := c.fetch(u, &response); err != nil {
		return nil, err
	}
	return response.parse(), nil
}

// getRangeResponse represents the http response from the GetRange http request
type getRangeResponse struct {
	Rates map[string]map[string]float64 `json:"rates"`
}

// parse returns the mapping of currencies to values from the response of the API
func (r *getRangeResponse) parse() map[string]map[currency.Unit]float64 {
	var rates = make(map[string]map[currency.Unit]float64)
	for d, m := range r.Rates {
		tm := make(map[currency.Unit]float64)
		// parse values for time
		for s, v := range m {
			c, err := currency.ParseISO(s)
			if err != nil {
				continue
			}
			tm[c] = v
		}
		rates[d] = tm
	}
	return rates
}
