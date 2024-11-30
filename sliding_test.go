package sliding_test

import (
	"sync"
	"testing"
	"time"

	"github.com/tetsuo/sliding"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	c := sliding.NewCounter(time.Millisecond * 77)

	for range 100 {
		c.Increment()
	}

	assert.EqualValues(t, 100, c.Peek())

	time.Sleep(time.Millisecond * 100)
	assert.EqualValues(t, 0, c.Peek())

	c.Increment()
	assert.EqualValues(t, 1, c.Peek())

	time.Sleep(time.Millisecond * 40)

	for range 3 {
		c.Increment()
	}

	assert.EqualValues(t, 4, c.Peek())

	time.Sleep(time.Millisecond * 40)
	assert.EqualValues(t, 3, c.Peek())

	for range 5 {
		c.Increment()
	}

	time.Sleep(time.Millisecond * 40)
	assert.EqualValues(t, 5, c.Peek())

	time.Sleep(time.Millisecond * 40)
	assert.EqualValues(t, 0, c.Peek())
}

func TestCounterConcurrent(t *testing.T) {
	c := sliding.NewCounter(time.Second * 1)

	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100/4; i++ {
				c.Increment()
			}
		}()
	}

	wg.Wait()
	assert.EqualValues(t, 100, c.Peek())
}
