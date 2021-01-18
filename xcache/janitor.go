package xcache

import "time"

// clean up loop
type janitor struct {
	interval time.Duration
	stop     chan bool
}

func (j *janitor) run(c *cache) {
	ticker := time.NewTicker(j.interval)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func runJanitor(c *cache, ci time.Duration) {
	c.janitor = &janitor{
		interval: ci,
		stop:     make(chan bool),
	}
	go c.janitor.run(c)
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}
