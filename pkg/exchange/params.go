package exchange

import (
	"strings"

	"golang.org/x/text/currency"
)

type Params struct {
	// currency to quote against
	BaseCurrency *currency.Unit
	// currencies to return
	Currencies []currency.Unit
}

// setParams encodes the Params to be added to the http request
func setParams(p *Params) map[string]string {
	// build url parameters
	params := make(map[string]string)
	// add base currency
	if p.BaseCurrency != nil {
		params["base"] = p.BaseCurrency.String()
	}
	// add currencies
	if p.Currencies != nil || len(p.Currencies) != 0 {
		var symbols = make([]string, 0, len(p.Currencies))
		for _, s := range p.Currencies {
			symbols = append(symbols, s.String())
		}
		params["symbols"] = strings.Join(symbols, ",")
	}
	return params
}
