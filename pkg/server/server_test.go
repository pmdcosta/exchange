package server_test

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/pmdcosta/exchange/pkg/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	// create server
	l := zerolog.Nop()
	s := server.New(&l, "localhost:8080")

	// create handler
	var call bool
	r := func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			call = true
		})
	}

	// start and stop server
	go s.Start(r)
	defer s.Done()

	// execute request
	_, err := http.Get("http://localhost:8080/test")
	require.Nil(t, err)
	require.True(t, call)
}
