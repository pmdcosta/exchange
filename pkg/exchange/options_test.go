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

func TestClient_SetCache(t *testing.T) {
	t.Run("cache miss", testClient_SetCache_miss)
	t.Run("cache hit", testClient_SetCache_hit)

}

func testClient_SetCache_miss(t *testing.T) {
	now := time.Now()
	date := fmt.Sprintf("/%d-%d-%d", now.Year(), now.Month(), now.Day())

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, date, req.URL.String())
		_, _ = res.Write([]byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"EUR","date":"2019-11-15"}`))
	}))
	defer testServer.Close()

	// create client
	c := newClientWithCache(t, testServer.URL)
	defer c.Finish()

	// mock cache
	c.mockCache.EXPECT().Load(date).Times(1).Return(false, nil)
	c.mockCache.EXPECT().Save(date, []byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"EUR","date":"2019-11-15"}`)).Times(1).Return()

	// execute request
	rates, err := c.GetLatest()
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.CAD: 1.4608, currency.HKD: 8.6361}, rates)
}

func testClient_SetCache_hit(t *testing.T) {
	// create client
	c := newClientWithCache(t, "")
	defer c.Finish()

	now := time.Now()
	date := fmt.Sprintf("/%d-%d-%d", now.Year(), now.Month(), now.Day())

	// mock cache
	c.mockCache.EXPECT().Load(date).Times(1).Return(true, []byte(`{"rates":{"CAD":1.4608,"HKD":8.6361},"base":"EUR","date":"2019-11-15"}`))

	// execute request
	rates, err := c.GetLatest()
	require.Nil(t, err)
	require.Equal(t, map[currency.Unit]float64{currency.CAD: 1.4608, currency.HKD: 8.6361}, rates)
}
