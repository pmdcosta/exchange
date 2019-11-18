package exchange_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
)

func TestClient_GetRange(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/history?end_at=2019-11-04&start_at=2019-11-01&symbols=EUR%2CUSD", req.URL.String())
		_, _ = res.Write([]byte(`{"rates":{"2019-11-01":{"EUR":1.1626825412,"USD":1.2951120826},"2019-11-04":{"EUR":1.1578362356,"USD":1.2919136717}},"start_at":"2019-11-01","base":"GBP","end_at":"2019-11-04"}`))
	}))
	defer testServer.Close()

	// create client
	c := newClient(t, testServer.URL)
	defer c.Finish()

	// execute request
	rates, err := c.GetRangeWithParams(&exchange.Params{Currencies: []currency.Unit{currency.EUR, currency.USD}}, "2019-11-01", "2019-11-04")
	require.Nil(t, err)
	expected := map[string]map[currency.Unit]float64{
		"2019-11-01": {
			currency.EUR: 1.1626825412,
			currency.USD: 1.2951120826,
		},
		"2019-11-04": {
			currency.EUR: 1.1578362356,
			currency.USD: 1.2919136717,
		},
	}
	require.Equal(t, expected, rates)
}
