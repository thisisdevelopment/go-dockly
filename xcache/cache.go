package xcache

import (
	"runtime"
	"sync"
	"time"
)

const (
	// NoExpiration eg will never be deleted
	NoExpiration time.Duration = -1
	// DefaultExpiration as configured with New() eg 5 minutes
	DefaultExpiration time.Duration = 0
)

// Cache construct to control garbage collector
type Cache struct {
	*cache
}

type cache struct {
	expiration time.Duration
	items      map[string]Item
	mu         sync.RWMutex
	onEvicted  func(string, interface{})
	janitor    *janitor
}

// Item is a generic interface holding the cache object
type Item struct {
	Object     interface{}
	Expiration int64
}

// New returns a cache with a given default expiration and cleanup
// interval. If the expiration duration is less than one (or NoExpiration),
// the items in the cache never expire (by default)
func New(expiration, cleanupInterval time.Duration) *Cache {

	if expiration == 0 {
		expiration = -1
	}

	c := &cache{
		expiration: expiration,
		items:      make(map[string]Item),
	}
	// trick ensures that the janitor routine does not keep
	// C from being garbage collected.
	// On garbage collection, the finalizer stops the janitor routine,
	// and c will be collected.
	C := &Cache{c}
	if cleanupInterval > 0 {
		runJanitor(c, cleanupInterval)
		runtime.SetFinalizer(C, stopJanitor)
	}

	return C
}

// Set an item to the cache, replacing any existing item. If the duration is 0
// (DefaultExpiration), the cache's default expiration time is used. If it is -1
// (NoExpiration), the item never expires.
func (c *cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.expiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
	c.mu.Unlock()
}

// Get an item from the cache or nil, and a bool indicating
// whether the key was found.
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.items[k]
	if !exists {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

// Copies all unexpired items in the cache into a new map and returns it.
func (c *cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		// "Inlining" of Expired
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

// delete all expired items from the cache.
func (c *cache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.Delete(k)
		}
	}
	c.mu.Unlock()
}

// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
}

func (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

// Sets an (optional) function that is called with the key and value when an
// item is evicted from the cache. (Including when it is deleted manually, but
// not when it is overwritten.) Set to nil to disable.
func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}
