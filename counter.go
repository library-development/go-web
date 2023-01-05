package web

import "sync"

type Counter struct {
	n    int
	lock sync.Mutex
}

func (c *Counter) Inc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.n++
}

func (c *Counter) Dec() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.n--
}

func (c *Counter) Value() int {
	return c.n
}
