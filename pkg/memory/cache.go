package memory

import (
	"sync"

	"github.com/rs/zerolog"
)

// Cache is a key value in-memory store
type Cache struct {
	logger *zerolog.Logger

	storage map[string][]byte
	lock    sync.RWMutex
}

// New creates a new in-memory store
func New(logger *zerolog.Logger) *Cache {
	l := logger.With().Str("pkg", "memory-cache").Logger()
	return &Cache{
		logger:  &l,
		storage: make(map[string][]byte),
		lock:    sync.RWMutex{},
	}
}

// Save stores a key-value pair
func (c *Cache) Save(key string, value []byte) {
	c.logger.Debug().Str("key", key).Msg("saving key...")
	c.lock.Lock()
	defer c.lock.Unlock()
	c.storage[key] = value
	c.logger.Debug().Str("key", key).Msg("key saved")
}

// Load returns the value for a key if found
func (c *Cache) Load(key string) (bool, []byte) {
	c.logger.Debug().Str("key", key).Msg("loading key...")
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, found := c.storage[key]
	c.logger.Debug().Str("key", key).Bool("found", found).Msg("key loaded")
	return found, value
}
