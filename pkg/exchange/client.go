package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/text/currency"
)

const (
	// address of the exchange api
	address = "https://api.exchangeratesapi.io/"
	// latest endpoint path
	pathLatest = "/latest"
	// history endpoint path
	pathHistory = "/history"
)

// Client is the exchangeratesapi.io API client
type Client struct {
	logger *zerolog.Logger
	// internal http client
	client *http.Client

	// http address of the API
	address string
	// currency to quote against
	baseCurrency currency.Unit
	// currencies to return
	currencies []currency.Unit
	// cache stores the requests and responses
	cache Cache
}

// NewClient creates a new exchangeratesapi.io API client
func NewClient(logger *zerolog.Logger, options ...Option) *Client {
	l := logger.With().Str("pkg", "exchange-client").Logger()
	c := Client{
		logger:       &l,
		client:       &http.Client{Timeout: time.Second * 10},
		address:      address,
		baseCurrency: currency.GBP,
	}
	for _, opt := range options {
		opt(&c)
	}
	return &c
}

// Cache keeps the requests and responses from the API to avoid repeating requests
//go:generate mockgen -destination ../../mocks/cache_mock.go -package mocks github.com/pmdcosta/exchange/pkg/exchange Cache
type Cache interface {
	Save(request string, body []byte)
	Load(request string) (bool, []byte)
}

// buildURL builds the url for the http requests
func (c Client) buildURL(path string, params map[string]string) *url.URL {
	// build url
	u, _ := url.Parse(c.address)

	// build url parameters
	if params == nil {
		params = make(map[string]string)
	}
	// add base currency
	params["base"] = c.baseCurrency.String()
	// add currencies
	if c.currencies != nil || len(c.currencies) != 0 {
		var symbols = make([]string, 0, len(c.currencies))
		for _, s := range c.currencies {
			symbols = append(symbols, s.String())
		}
		params["symbols"] = strings.Join(symbols, ",")
	}
	// build url query
	q, _ := url.ParseQuery(u.RawQuery)
	for k, v := range params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	// set url path
	u.Path = path
	return u
}

// fetch executes an http request
func (c Client) fetch(u *url.URL, resp interface{}) error {
	c.logger.Info().Str("url", u.String()).Msg("executing http request...")

	// if available, retrieve from cache
	if c.cache != nil {
		if hit, body := c.cache.Load(u.RequestURI()); hit {
			if err := json.Unmarshal(body, &resp); err != nil {
				return err
			}
			return nil
		}
	}

	// execute http request
	r, err := c.client.Get(u.String())
	if err != nil {
		c.logger.Error().Str("url", u.String()).Err(err).Msg("failed to execute http request")
		return err
	}
	defer r.Body.Close()

	// check status code
	if r.StatusCode != http.StatusOK {
		c.logger.Error().Str("url", u.String()).Str("status", http.StatusText(r.StatusCode)).Msg("unexpected status code")
		return fmt.Errorf("unexpected status code: %s", http.StatusText(r.StatusCode))
	}

	// read body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.logger.Error().Str("url", u.String()).Err(err).Msg("failed read response body")
		return err
	}

	// handle response
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		c.logger.Error().Str("url", u.String()).Err(err).Msg("failed to unmarshal response")
		return err
	}

	// if available, save to cache
	if c.cache != nil {
		c.cache.Save(u.RequestURI(), bodyBytes)
	}

	return nil
}
