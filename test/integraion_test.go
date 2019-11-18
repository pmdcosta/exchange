//+build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Rates(t *testing.T) {
	current := getCurrentRate(t)
	average := getAverageRate(t)

	// get the exchange rate from the service
	r, err := http.Get("http://localhost:8080/v1/rates/USD")
	require.Nil(t, err)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	require.Nil(t, err)
	var resp getRatesResponse
	require.Nil(t, json.Unmarshal(bodyBytes, &resp))

	var reason string
	if current > average {
		reason = fmt.Sprintf("the current exchange value is above the weekly average of%9f", average)
	} else {
		reason = fmt.Sprintf("the current exchange value is below the weekly average of%9f", average)
	}
	expected := getRatesResponse{
		Rate: map[string]float64{
			"USD": 1,
			"EUR": current,
		},
		Exchange: &getRatesRecommendationResponse{
			Good:   current > average,
			Reason: reason,
		},
	}
	require.Equal(t, expected, resp)
}

type getRatesResponse struct {
	Rate     map[string]float64              `json:"rate"`
	Exchange *getRatesRecommendationResponse `json:"exchange,omitempty"`
}
type getRatesRecommendationResponse struct {
	Good   bool   `json:"recommended"`
	Reason string `json:"reason"`
}

// retrieves the current rate from exchangeratesapi.io
func getCurrentRate(t *testing.T) float64 {
	r, err := http.Get("https://api.exchangeratesapi.io/latest?base=USD&symbols=EUR")
	require.Nil(t, err)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	require.Nil(t, err)
	var resp getCurrentRateResponse
	require.Nil(t, json.Unmarshal(bodyBytes, &resp))
	return resp.Rates["EUR"]
}

type getCurrentRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// retrieves the current rate from exchangeratesapi.io
func getAverageRate(t *testing.T) float64 {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	past := today.AddDate(0, 0, -7)
	start := fmt.Sprintf("%d-%d-%d", past.Year(), past.Month(), past.Day())
	end := fmt.Sprintf("%d-%d-%d", today.Year(), today.Month(), today.Day())
	request := fmt.Sprintf("https://api.exchangeratesapi.io/history?start_at=%s&end_at=%s&base=USD&symbols=EUR", start, end)
	r, err := http.Get(request)
	require.Nil(t, err)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	require.Nil(t, err)
	var resp getRangeResponse
	require.Nil(t, json.Unmarshal(bodyBytes, &resp))

	var average float64
	for _, c := range resp.Rates {
		average += c["EUR"]
	}
	return average / float64(len(resp.Rates))
}

type getRangeResponse struct {
	Rates map[string]map[string]float64 `json:"rates"`
}
