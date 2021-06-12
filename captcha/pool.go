package captcha

import "sync"

type cachedGenerator struct {
	cache        []*Captcha
	currentIndex int
	lock         *sync.Mutex
	st           *staticGenerator
	newCaches    chan []*Captcha
}

func NewCachedGenerator() *cachedGenerator {
	c := &cachedGenerator{
		cache:     make([]*Captcha, 10000),
		st:        NewStaticGenerator(),
		lock:      &sync.Mutex{},
		newCaches: make(chan []*Captcha),
	}

	go func() {
		for {
			newCache := make([]*Captcha, 10000)
			for i, _ := range newCache {
				newCache[i], _ = c.st.Generate()
			}
			println("new cache generated")
			c.newCaches <- newCache
		}
	}()

	return c
}

func (c *cachedGenerator) Generate() (*Captcha, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// if its time to get get.
	if c.currentIndex > len(c.cache)-1 || c.cache[0] == nil {
		select {
		case newCache := <-c.newCaches:
			// release all the resources.
			if c.cache[0] != nil {
				for i, _ := range c.cache {
					c.st.Release(c.cache[i])
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
	return c.cache[randBetween(0, len(c.cache)-1)], nil

}

func (c *cachedGenerator) Release(captcha *Captcha) error {
	return nil
}

var _ Generator = &cachedGenerator{}
