package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/pmdcosta/exchange/pkg/memory"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/currency"
)

func main() {
	// initialize clients
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	cache := memory.New(&l)
	client := exchange.NewClient(&l, exchange.SetCache(cache))

	// make request
	data, err := client.WithOptions(
		exchange.SetBaseCurrency(currency.EUR),
		exchange.SetCurrencies(currency.USD, currency.GBP),
	).GetLatest()
	if err != nil {
		l.Fatal().Err(err).Msg("failed to get latest rates")
	}
	fmt.Println(err, data)
}
