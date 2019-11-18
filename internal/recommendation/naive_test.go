package recommendation_test

import (
	"testing"

	"github.com/pmdcosta/exchange/internal/recommendation"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
)

var data = map[string]map[currency.Unit]float64{
	"1": {
		currency.EUR: 0.9,
	},
	"2": {
		currency.EUR: 0.8,
	},
	"3": {
		currency.EUR: 0.7,
	},
}

func TestNaive_False(t *testing.T) {
	rec, reason := recommendation.Naive(0.79, currency.EUR, data)
	require.False(t, rec)
	require.Equal(t, "the current exchange value is below the weekly average of 0.800000", reason)
}

func TestNaive_True(t *testing.T) {
	rec, reason := recommendation.Naive(0.81, currency.EUR, data)
	require.True(t, rec)
	require.Equal(t, "the current exchange value is above the weekly average of 0.800000", reason)
}
