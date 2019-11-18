package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/rs/zerolog"
	"golang.org/x/text/currency"
)

// Routes
const (
	RouteRate = "/v1/rates/{currency}"
)

// Handler is the http handler
type Handler struct {
	logger *zerolog.Logger

	exchange       Exchange
	recommendation Recommendation
}

// Recommendation whether this is a good time to exchange money or not
type Recommendation func(current float64, currency currency.Unit, data map[string]map[currency.Unit]float64) (bool, string)

//go:generate mockgen -destination ../../mocks/exchange_mock.go -package mocks github.com/pmdcosta/exchange/internal/handler Exchange
type Exchange interface {
	GetRangeWithParams(p *exchange.Params, start string, end string) (map[string]map[currency.Unit]float64, error)
}

// New creates a new http handler
func New(logger *zerolog.Logger, exchange Exchange, recommendation Recommendation) *Handler {
	l := logger.With().Str("pkg", "http-handler").Logger()
	return &Handler{
		logger:         &l,
		exchange:       exchange,
		recommendation: recommendation,
	}
}

// Routes creates the http routes
func (h *Handler) Routes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Get(RouteRate, h.GetRate)
	})
}

// writeJSON sends a JSON marshalled object to the response
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	if data == nil {
		data = StatusResponse{Status: http.StatusText(statusCode)}
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
	return nil
}

// writeError sends a JSON encoded error response
func writeError(w http.ResponseWriter, status int, err error) error {
	return writeJSON(w, status, ErrorResponse{Error: err.Error()})
}

// StatusResponse is a fallback response struct for successful requests
type StatusResponse struct {
	Status string `json:"status"`
}

// ErrorResponse is a fallback response struct for failed requests
type ErrorResponse struct {
	Error string `json:"error"`
}
