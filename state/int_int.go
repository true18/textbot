/*
 * Copyright (c) Zinnatullin Chingiz, 2022.
 */

package state

import "sync"

type intInt struct {
	mutex sync.Mutex
	m     map[int]int
}

func (c *intInt) Add(k, v int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.m[k] = v
}

func (c *intInt) Pop(k int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.m, k)
}

func (c *intInt) In(k int) (v int, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v, ok = c.m[k]

	return
}
