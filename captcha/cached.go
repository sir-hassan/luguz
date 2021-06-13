package captcha

import (
	"fmt"
	"sync"
)

// Struct cachedGenerator builds a list of pre-rendered captcha images and servers
// them sequentially. At the same time it runs a go routine that builds new
// caches and replace the serving cache when the serving is depleted. During the
// very heavy load, the serving cache will deplete before the cache building go
// routine is finished building the next one. In this case, a random entry will
// be served from the serving cache until a new cache is ready.

type cachedGenerator struct {
	cache        []*Captcha
	currentIndex int
	lock         *sync.Mutex
	st           Generator
	newCaches    chan []*Captcha
}

// NewCachedGenerator func creates a new cachedGenerator instance.
func NewCachedGenerator(config GeneratorConfig, size int) Generator {
	c := &cachedGenerator{
		cache:     make([]*Captcha, size),
		st:        NewBasicGenerator(config),
		lock:      &sync.Mutex{},
		newCaches: make(chan []*Captcha),
	}

	go func() {
		for {
			newCache := make([]*Captcha, size)
			for i := range newCache {
				newCache[i], _ = c.st.Generate()
			}
			println("new cache generated")
			c.newCaches <- newCache
		}
	}()

	return c
}

// Generate func generates a new captcha.
func (c *cachedGenerator) Generate() (*Captcha, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// if its time to get new cache.
	if c.currentIndex > len(c.cache)-1 || c.cache[0] == nil {
		select {
		case newCache := <-c.newCaches:
			// release all the resources.
			if c.cache[0] != nil {
				for i := range c.cache {
					if err := c.st.Release(c.cache[i]); err != nil {
						return nil, fmt.Errorf("releaseing reasources error: %s", err)
					}
					c.cache[i] = nil
				}
			}
			c.cache = newCache
			c.currentIndex = 0
		default:
		}
	}

	if c.cache[0] == nil {
		// cache is not generated.
		return c.st.Generate()
	}
	if c.currentIndex <= len(c.cache)-1 {
		c.currentIndex++
		return c.cache[c.currentIndex-1], nil
	}
	return c.cache[randInMinMax(0, len(c.cache)-1)], nil

}

func (c *cachedGenerator) Release(captcha *Captcha) error {
	// we don't release captcha resource as this is a cached implementation.
	return nil
}

var _ Generator = &cachedGenerator{}
