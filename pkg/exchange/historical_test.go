package exchange_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
)

func TestClient_GetHistorical(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/2019-11-14", req.URL.String())
		_, _ = res.Write([]byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"GBP","date":"2019-11-15"}`))
	}))
	defer testServer.Close()

	// create client
	c := newClient(t, testServer.URL)
	defer c.Finish()

	// execute request
	rates, err := c.GetHistorical("2019-11-14")
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.CAD: 1.4608, currency.HKD: 8.6361}, rates)
}
