package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Entries  map[string]CacheEntry
	interval time.Duration
	mu       *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	newEntries := make(map[string]CacheEntry)
	return Cache{
		interval: interval,
		Entries:  newEntries,
		mu:       &sync.Mutex{},
	}
}

func (c *Cache) Add(key string, val []byte) {
	if _, exists := c.Entries[key]; !exists {
		newEntry := CacheEntry{
			createdAt: time.Now(),
			val:       val,
		}
		c.mu.Lock()
		c.Entries[key] = newEntry
		c.mu.Unlock()
	}
	c.reapLoop(key)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, exists := c.Entries[key]
	c.mu.Unlock()
	if !exists {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(key string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	expired := make(chan bool)
	go func() {
		time.Sleep(c.interval)
		expired <- true
	}()
	go func() {
		if <-expired {
			c.mu.Lock()
			delete(c.Entries, key)
			c.mu.Unlock()
		}
	}()

}
