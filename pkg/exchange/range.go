package exchange

import (
	"fmt"
	"time"

	"golang.org/x/text/currency"
)

// GetRange retrieves the foreign exchange reference rates for a specific range of dates
func (c Client) GetRange(start time.Time, end time.Time) (map[time.Time]map[currency.Unit]float64, error) {
	s := fmt.Sprintf("%d-%d-%d", start.Year(), start.Month(), start.Day())
	e := fmt.Sprintf("%d-%d-%d", end.Year(), end.Month(), end.Day())
	c.logger.Debug().Str("start", s).Str("end", e).Msg("getting rates in range...")
	u := c.buildURL(pathHistory, map[string]string{"start_at": s, "end_at": e})
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
func (r *getRangeResponse) parse() map[time.Time]map[currency.Unit]float64 {
	layout := "2006-01-02"
	var rates = make(map[time.Time]map[currency.Unit]float64)
	for d, m := range r.Rates {
		// parse time
		t, err := time.Parse(layout, d)
		if err != nil {
			continue
		}
		tm := make(map[currency.Unit]float64)
		// parse values for time
		for s, v := range m {
			c, err := currency.ParseISO(s)
			if err != nil {
				continue
			}
			tm[c] = v
		}
		rates[t] = tm
	}
	return rates
}
