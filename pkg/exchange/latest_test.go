package exchange_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
)

func TestClient_GetLatest(t *testing.T) {
	now := time.Now()
	date := fmt.Sprintf("/%d-%d-%d", now.Year(), now.Month(), now.Day())
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, date, req.URL.String())
		_, _ = res.Write([]byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"GBP","date":"2019-11-15"}`))
	}))
	defer testServer.Close()

	// create client
	c := newClient(t, testServer.URL)
	defer c.Finish()

	// execute request
	rates, err := c.GetLatest()
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.CAD: 1.4608, currency.HKD: 8.6361}, rates)
}
