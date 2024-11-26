package sliding

import (
	"sync"
	"testing"
	"time"

	"github.com/onur1/ring"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	c, r, ticker := newTestCounter(time.Millisecond * 80)

	dur := c.Duration().Milliseconds()
	framedur := ticker.dur

	assert.Equal(t, 77, int(dur))
	assert.Equal(t, 8, r.Size())
	assert.Equal(t, 11, int(framedur))

	go c.loop(r, ticker)

	var now int64

	assert.Equal(t, 0, c.Peek())

	c.Inc()
	c.Inc()

	assert.Equal(t, 2, c.Peek())

	for now < dur-(framedur*2) {
		now = ticker.advance()
		assert.Equal(t, 2, c.Peek())
	}

	assert.Equal(t, dur-(framedur*2), now)
	assert.Equal(t, 2, c.Peek())

	c.Inc()

	assert.Equal(t, 3, c.Peek())

	now = ticker.advance()

	assert.Equal(t, dur-framedur, now)
	assert.Equal(t, 3, c.Peek())

	now = ticker.advance()

	assert.Equal(t, 1, c.Peek())

	c.Inc()

	assert.Equal(t, 2, c.Peek())
	assert.Equal(t, dur, now)

	for now < dur+(dur-(framedur*3)) {
		now = ticker.advance()

		assert.Equal(t, 2, c.Peek())
	}

	assert.Equal(t, dur+(dur-(framedur*3)), now)

	now = ticker.advance()

	assert.Equal(t, dur+(dur-(framedur*2)), now)
	assert.Equal(t, 1, c.Peek())

	ticker.advance()

	assert.Equal(t, 1, c.Peek())

	now = ticker.advance()

	assert.Equal(t, dur*2, now)
	assert.Equal(t, 0, c.Peek())

	var wg sync.WaitGroup

	wg.Add(1)

	go func(t *fakeTicker) {
		defer wg.Done()
		<-t.stopch
	}(ticker)

	c.Stop()

	wg.Wait()
}

type fakeTicker struct {
	mu     sync.Mutex
	ch     chan time.Time
	now    int64
	dur    int64
	stopch chan struct{}
}

func newFakeTicker(d time.Duration) *fakeTicker {
	return &fakeTicker{
		ch:     make(chan time.Time),
		now:    0,
		dur:    d.Milliseconds(),
		stopch: make(chan struct{}),
	}
}

func (t *fakeTicker) advance() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now += t.dur
	t.ch <- time.Time{}

	return t.now
}

func (t *fakeTicker) getC() <-chan time.Time {
	return t.ch
}

func (t *fakeTicker) stop() {
	t.stopch <- struct{}{}
}

func newTestCounter(d time.Duration) (*Counter, *ring.Ring[int], *fakeTicker) {
	r, dur := newRing(d)

	c := &Counter{
		count: make(chan struct{}),
		exit:  make(chan struct{}, 1),
		peek:  make(chan chan int),
		dur:   dur,
	}

	framedur := time.Duration(dur.Milliseconds()/int64(r.Size()-1)) * time.Millisecond

	return c, r, newFakeTicker(framedur)
}
