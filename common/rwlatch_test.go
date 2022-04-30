// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package common

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type counter struct {
	count int
	m     ReaderWriterLatch
}

func (c *counter) add(num int) {
	c.m.WLock()
	defer c.m.WUnlock()
	c.count += num
}

func (c *counter) read() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.count
}

func TestRWLatch(t *testing.T) {
	wg := sync.WaitGroup{}
	numCon := 100
	c := &counter{count: 0, m: NewRWLatch()}
	c.add(5)
	for i := 0; i < numCon; i++ {
		wg.Add(1)
		if i%2 == 0 {
			go func(cnt *counter) {
				defer wg.Done()
				cnt.read()
			}(c)
		} else {
			go func(cnt *counter) {
				defer wg.Done()
				cnt.add(1)
			}(c)
		}
	}
	wg.Wait()
	assert.Equal(t, 55, c.read())
}
