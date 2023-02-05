package sliding

import (
	"math"
	"time"

	"github.com/onur1/ring"
)

// A Counter is a sliding window based counter and used for counting
// how frequently events are happening.
type Counter struct {
	count chan struct{}
	exit  chan struct{}
	peek  chan chan int
	dur   time.Duration
}

// NewCounter returns a new Counter with the given time interval.
func NewCounter(d time.Duration) *Counter {
	r, dur := newRing(d)

	c := &Counter{
		count: make(chan struct{}),
		exit:  make(chan struct{}, 1),
		peek:  make(chan chan int),
		dur:   dur,
	}

	framedur := time.Duration(dur.Milliseconds()/int64(r.Size()-1)) * time.Millisecond

	go c.loop(r, newClockTicker(framedur))

	return c
}

// Duration returns the calculated duration of time window.
func (c *Counter) Duration() time.Duration {
	return c.dur
}

// Stop stops a counter.
func (c *Counter) Stop() {
	close(c.exit)
}

// Peek returns the current count.
func (c *Counter) Peek() int {
	res := make(chan int, 1)
	c.peek <- res
	return <-res
}

// Inc increments a counter by 1.
func (c *Counter) Inc() {
	c.count <- struct{}{}
}

func (c *Counter) loop(r *ring.Ring[int], ticker ticker) {
	ptr, head := 0, 0
	len := r.Size() - 1
	C := ticker.getC()

LOOP:
	for {
		select {
		case <-c.exit:
			break LOOP
		case <-c.count:
			head = head + 1
			r.Put(ptr, head)
		case <-C:
			ptr = ptr + 1
			ptr = r.Put(ptr, head)
		case ch := <-c.peek:
			ch <- r.Get(ptr) - r.Get(ptr-len)
		}
	}

	ticker.stop()

	close(c.count)
	close(c.peek)
}

type ticker interface {
	getC() <-chan time.Time
	stop()
}

func newRing(d time.Duration) (*ring.Ring[int], time.Duration) {
	millis := int(d.Milliseconds())

	// ringsize is always a power of 2 and less than 64
	r := ring.NewRing[int](int(math.Ceil(math.Log(float64(millis))/math.Log(2))) + 1)

	len := r.Size() - 1

	// time window value needs to be rounded to the nearest millisecond
	// which is divisible by ringsize - 1
	diff := (millis % len)
	if diff <= len/2 {
		millis = millis - diff
	} else {
		millis = millis + (len - diff)
	}

	return r, time.Duration(millis) * time.Millisecond
}

type clockTicker struct {
	ticker *time.Ticker
}

func newClockTicker(d time.Duration) *clockTicker {
	ticker := time.NewTicker(d)
	return &clockTicker{
		ticker,
	}
}

func (t *clockTicker) getC() <-chan time.Time {
	return t.ticker.C
}

func (t *clockTicker) stop() {
	t.ticker.Stop()
	t.ticker = nil
}
