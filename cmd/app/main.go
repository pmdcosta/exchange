package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pmdcosta/exchange/internal/handler"
	"github.com/pmdcosta/exchange/internal/recommendation"
	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/pmdcosta/exchange/pkg/memory"
	"github.com/pmdcosta/exchange/pkg/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// initialize clients
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	cache := memory.New(&l)
	client := exchange.NewClient(&l, exchange.SetCache(cache))
	s := server.New(&l, "localhost:8080")
	h := handler.New(&l, client, recommendation.Naive)
	go s.Start(h.Routes)
	defer s.Done()
	<-sigs
}
