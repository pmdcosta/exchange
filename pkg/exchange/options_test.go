package exchange_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
)

func TestClient_SetBase(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/latest?base=EUR", req.URL.String())
		_, _ = res.Write([]byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"EUR","date":"2019-11-15"}`))
	}))
	defer testServer.Close()

	// create client
	c := newClient(t, testServer.URL)
	defer c.Finish()

	// execute request
	rates, err := c.WithOptions(exchange.SetBaseCurrency(currency.EUR)).GetLatest()
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.CAD: 1.4608, currency.HKD: 8.6361}, rates)
}

func TestClient_SetCurrencies(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/latest?base=GBP&symbols=EUR%2CUSD", req.URL.String())
		_, _ = res.Write([]byte(`{"rates":{"EUR":1.4608,"USD":8.6361},"base":"GBP","date":"2019-11-15"}`))
	}))
	defer testServer.Close()

	// create client
	c := newClient(t, testServer.URL)
	defer c.Finish()

	// execute request
	rates, err := c.WithOptions(exchange.SetCurrencies(currency.EUR, currency.USD)).GetLatest()
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.EUR: 1.4608, currency.USD: 8.6361}, rates)
}

func TestClient_SetCache(t *testing.T) {
	t.Run("cache miss", testClient_SetCache_miss)
	t.Run("cache hit", testClient_SetCache_hit)

}

func testClient_SetCache_miss(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/latest?base=EUR", req.URL.String())
		_, _ = res.Write([]byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"EUR","date":"2019-11-15"}`))
	}))
	defer testServer.Close()

	// create client
	c := newClient(t, testServer.URL)
	defer c.Finish()

	// mock cache
	c.mockCache.EXPECT().Load("/latest?base=EUR").Times(1).Return(false, nil)
	c.mockCache.EXPECT().Save("/latest?base=EUR", []byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"EUR","date":"2019-11-15"}`)).Times(1).Return()

	// execute request
	rates, err := c.WithOptions(exchange.SetCache(c.mockCache), exchange.SetBaseCurrency(currency.EUR)).GetLatest()
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.CAD: 1.4608, currency.HKD: 8.6361}, rates)
}

func testClient_SetCache_hit(t *testing.T) {
	// create client
	c := newClient(t, "")
	defer c.Finish()

	// mock cache
	c.mockCache.EXPECT().Load("/latest?base=EUR").Times(1).Return(true, []byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"EUR","date":"2019-11-15"}`))

	// execute request
	rates, err := c.WithOptions(exchange.SetCache(c.mockCache), exchange.SetBaseCurrency(currency.EUR)).GetLatest()
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.CAD: 1.4608, currency.HKD: 8.6361}, rates)
}
