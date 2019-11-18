package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/namsral/flag"
	"github.com/pmdcosta/exchange/internal/handler"
	"github.com/pmdcosta/exchange/internal/recommendation"
	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/pmdcosta/exchange/pkg/memory"
	"github.com/pmdcosta/exchange/pkg/server"
	"github.com/rs/zerolog"
)

func main() {
	// parse flags
	var (
		debug     = flag.Bool("debug", false, "print debug level log messages")
		host      = flag.String("host", "0.0.0.0:8080", "http server host to bind to")
		recEngine = flag.String("recommendation", "naive", "recommendation engine to use")
	)
	flag.Parse()

	// initialize logger
	if !*debug {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	l := zerolog.New(os.Stdout).With().Logger()
	l.Info().Msg("starting up...")

	// set recommendation engine
	var rec handler.Recommendation
	switch *recEngine {
	case "naive":
		rec = recommendation.Naive
	}

	// bootstrap clients
	cache := memory.New(&l)
	exchangeClient := exchange.NewClient(&l, exchange.SetCache(cache))
	httpServer := server.New(&l, *host)
	httpHandler := handler.New(&l, exchangeClient, rec)
	go httpServer.Start(httpHandler.Routes)
	defer httpServer.Done()

	// wait for interrupt for shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	l.Info().Msg("shutting down...")
}
