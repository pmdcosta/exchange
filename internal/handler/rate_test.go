package handler_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/currency"
)

func TestHandler_GetRate(t *testing.T) {
	f := func(current float64, currency currency.Unit, data map[string]map[currency.Unit]float64) (bool, string) {
		return true, "the current exchange value is pretty good"
	}
	h := New(t, f)
	defer h.Finish()

	// mock service
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	past := now.AddDate(0, 0, -7)
	param := &exchange.Params{BaseCurrency: &currency.USD, Currencies: []currency.Unit{currency.EUR}}
	values := map[string]map[currency.Unit]float64{
		getDate(now.AddDate(0, 0, -1)): {
			currency.EUR: 1.234,
		},
		getDate(now.AddDate(0, 0, -2)): {
			currency.EUR: 2.345,
		},
	}
	h.exchange.EXPECT().GetRangeWithParams(param, getDate(past), getDate(now)).Times(1).Return(values, nil)

	// build http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/rates/USD", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("currency", "USD")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// make http request
	h.Handler.GetRate(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"rate":{"EUR":1.234,"USD":1},"exchange":{"recommended":true,"reason":"the current exchange value is pretty good"}}`, string(body))
}

func getDate(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
}
