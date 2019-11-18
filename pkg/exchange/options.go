package exchange

// Option is an optimal configuration that can be applied to a client
type Option func(c *Client)

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
