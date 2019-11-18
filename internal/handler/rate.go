package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pmdcosta/exchange/pkg/exchange"
	"golang.org/x/text/currency"
)

// GetRate retrieves the latest exchange rate and evaluates whether this is a good time to exchange money
func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	s := chi.URLParam(r, "currency")

	// check requested currency
	c, err := currency.ParseISO(s)
	if err != nil {
		_ = writeError(w, http.StatusBadRequest, err)
		return
	}

	// get rates from the exchange api
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	last := now.AddDate(0, 0, -7)
	param := &exchange.Params{BaseCurrency: &c, Currencies: []currency.Unit{currency.EUR}}
	rates, err := h.exchange.GetRangeWithParams(param, getDate(last), getDate(now))
	if err != nil {
		_ = writeError(w, http.StatusInternalServerError, err)
		return
	}

	// get the latest value from the range
	latest := getLatestInRange(rates)
	if latest == nil {
		_ = writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to parse latest value"))
		return
	}

	//check whether this is a good time to exchange money
	good, reason := h.recommendation(latest[currency.EUR], currency.EUR, rates)

	// build response
	resp := GetRatesResponse{
		Rate: map[string]float64{
			c.String():            1,
			currency.EUR.String(): latest[currency.EUR],
		},
		Exchange: GetRatesRecommendationResponse{
			Good:   good,
			Reason: reason,
		},
	}
	_ = writeJSON(w, http.StatusOK, resp)
}

// GetRatesResponse is the http response for the GetRate request
type GetRatesResponse struct {
	Rate     map[string]float64             `json:"rate"`
	Exchange GetRatesRecommendationResponse `json:"exchange"`
}

// GetRatesRecommendationResponse is part of the GetRatesResponse
type GetRatesRecommendationResponse struct {
	Good   bool   `json:"recommended"`
	Reason string `json:"reason"`
}

// getLatestInRange retrieves the latest value from a range
// this is needed as there is not data on weekends
func getLatestInRange(r map[string]map[currency.Unit]float64) map[currency.Unit]float64 {
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	for i := 0; i < 4; i++ {
		date := now.AddDate(0, 0, -i)
		if values, found := r[getDate(date)]; found {
			return values
		}
	}
	return nil
}

func getDate(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
}
