package exchange

import (
	"net/http"

	"golang.org/x/text/currency"
)

// Option is an optimal configuration that can be applied to a client
type Option func(c *Client)

// WithOptions creates a new client with optional parameters
func (c Client) WithOptions(options ...Option) Client {
	for _, opt := range options {
		opt(&c)
	}
	return c
}

// SetAddress sets the http address to connect to
func SetAddress(u string) Option {
	return func(c *Client) {
		c.address = u
	}
}

// SetCache sets the cache provider for the client
func SetCache(cache Cache) Option {
	return func(c *Client) {
		c.cache = cache
	}
}

// SetBaseCurrency currency to quote against
// Rates are quoted against the GBP by default
func SetBaseCurrency(currency currency.Unit) Option {
	return func(c *Client) {
		l := c.logger.With().Str("base-currency", currency.String()).Logger()
		c.logger = &l
		c.baseCurrency = currency
	}
}

// SetCurrencies sets specific exchange rates to be requested
func SetCurrencies(currencies ...currency.Unit) Option {
	return func(c *Client) {
		c.currencies = currencies
	}
}

// SetBackend overrides the http client for the requests
func SetBackend(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}
