package memory

import "sync"

// Cache is a key value in-memory store
type Cache struct {
	storage map[string][]byte
	lock    sync.RWMutex
}

// New creates a new in-memory store
func New() *Cache {
	return &Cache{
		storage: make(map[string][]byte),
		lock:    sync.RWMutex{},
	}
}

// Save stores a key-value pair
func (c *Cache) Save(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.storage[key] = value
}

// Load returns the value for a key if found
func (c *Cache) Load(key string) (bool, []byte) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	found, value := c.storage[key]
	return value, found
}
