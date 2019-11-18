package exchange

import "golang.org/x/text/currency"

// GetLatest retrieves the latest foreign exchange reference rates
func (c Client) GetLatest() (map[currency.Unit]float64, error) {
	u := c.buildURL(pathLatest, nil)
	var response getLatestResponse
	if err := c.fetch(u, &response); err != nil {
		return nil, err
	}
	return response.parse(), nil
}

// getLatestResponse represents the http response from the GetLatest http request
type getLatestResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// parse returns the mapping of currencies to values from the response of the API
func (r *getLatestResponse) parse() map[currency.Unit]float64 {
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
